package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"golang.org/x/xerrors"
)

const (
	userIDKey = "user_id"
	jtiKey    = "jti"
)

var db *sqlx.DB
var orecocoReserveToken string
var tokenAuth *jwtauth.JWTAuth
var orecocoReserveURL string

var allowedFileTypeAndMime = map[string]string{
	"image/jpeg": "JPEG image data",
	"image/png":  "PNG image data",
	"image/gif":  "GIF image data",
	"image/bmp":  "PB bitmap",
}

type InitializeResponse struct {
	Language string `json:"language"`
}

type OrecocoReserveInitializeResponse struct {
	Token string `json:"token"`
}

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type ServerEnv struct {
	OrecocoReserveURL string
	Port              string
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
}

type SearchUserResponse struct {
	Users []*GetUserResponse `json:"users"`
}

type CreateScheduleRequest struct {
	Attendees   []string `json:"attendees"`
	StartAt     int64    `json:"start_at"`
	EndAt       int64    `json:"end_at"`
	Title       string   `json:"title"`
	MeetingRoom string   `json:"meeting_room"`
	Description string   `json:"description"`
}

type CreateScheduleResponse struct {
	ID string `json:"id"`
}

type UpdateScheduleRequest struct {
	Attendees   []string `json:"attendees"`
	StartAt     int64    `json:"start_at"`
	EndAt       int64    `json:"end_at"`
	Title       string   `json:"title"`
	MeetingRoom string   `json:"meeting_room"`
	Description string   `json:"description"`
}

type UpdateMeetingRoomRequest struct {
	ScheduleID    string `json:"schedule_id"`
	MeetingRoomID string `json:"meeting_room_id"`
	StartAt       int64  `json:"start_at"`
	EndAt         int64  `json:"end_at"`
}

type ReserveMeetingRoomRequest struct {
	ScheduleID    string `json:"schedule_id"`
	MeetingRoomID string `json:"meeting_room_id"`
	StartAt       int64  `json:"start_at"`
	EndAt         int64  `json:"end_at"`
}

type GetScheduleResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	StartAt     int64              `json:"start_at"`
	EndAt       int64              `json:"end_at"`
	Attendees   []*GetUserResponse `json:"attendees"`
	MeetingRoom string             `json:"meeting_room"`
}

type GetMeetingRoomReservationResponse struct {
	RoomID string `json:"room_id"`
}

type User struct {
	ID          string    `db:"id"`
	Email       string    `db:"email"`
	Name        string    `db:"name"`
	Password    string    `db:"password"`
	ImageBinary []byte    `db:"image_binary"`
	CreatedAt   time.Time `db:"created_at"`
}

type Schedule struct {
	ID               string    `db:"id"`
	Title            string    `db:"title"`
	Description      string    `db:"description"`
	ScheduleAttendee string    `db:"schedule_attendee"`
	StartAt          time.Time `db:"start_at"`
	EndAt            time.Time `db:"end_at"`
}

type SessionData struct{}

type UserResponses []*GetUserResponse

func (urs UserResponses) Len() int {
	return len(urs)
}

func (urs UserResponses) Less(i, j int) bool {
	return urs[i].Email < urs[j].Email
}

func (urs UserResponses) Swap(i, j int) {
	urs[i], urs[j] = urs[j], urs[i]
}

type GetScheduleResponses []*GetScheduleResponse

type GetCalendarResponse struct {
	Date      int64                `json:"date"`
	Schedules GetScheduleResponses `json:"schedules"`
}

type Session struct {
	sync.RWMutex
	m map[string]SessionData
}

var sess *Session = &Session{m: make(map[string]SessionData, 0)}

func (s *Session) SetNewSession(sessionKey string, sessionData SessionData) {
	s.Lock()
	defer s.Unlock()

	s.m[sessionKey] = sessionData
}

func (s *Session) ExistSession(sessionKey string) bool {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[sessionKey]; !ok {
		return false
	}
	return true
}

func (s *Session) Get(sessionKey string) (SessionData, error) {
	s.Lock()
	defer s.Unlock()

	if val, ok := s.m[sessionKey]; !ok {
		return SessionData{}, xerrors.New("Session key not found")
	} else {
		return val, nil
	}
}

func (s *Session) Delete(sessionKey string) {
	s.Lock()
	defer s.Unlock()

	delete(s.m, sessionKey)
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

func NewMySQLConnectionEnv() *MySQLConnectionEnv {
	return &MySQLConnectionEnv{
		Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		User:     getEnv("MYSQL_USER", "r-isucon"),
		DBName:   getEnv("MYSQL_DBNAME", "r-calendar"),
		Password: getEnv("MYSQL_PASS", "r-isucon"),
	}
}

func (mc *MySQLConnectionEnv) OpenDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=UTC", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
	return sqlx.Open("mysql", dsn)
}

func NewServerEnv() *ServerEnv {
	return &ServerEnv{
		OrecocoReserveURL: getEnv("ORECOCO_RESERVE_URL", "localhost:3003"),
		Port:              ":" + getEnv("CALENDAR_PORT", "3000"),
	}
}

func init() {
	time.Local = time.FixedZone("UTC", 0)

	tokenAuth = jwtauth.New("HS256", []byte("secret_key"), nil)
}

func main() {
	serverEnv := NewServerEnv()
	dbEnv := NewMySQLConnectionEnv()

	var err error
	db, err = dbEnv.OpenDB()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	authMiddlewares := []func(http.Handler) http.Handler{
		jwtauth.Verifier(tokenAuth),
		jwtauth.Authenticator,
		parseJWTToken,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/initialize", postInitialize(serverEnv.OrecocoReserveURL, dbEnv))

	r.Post("/login", postLogin)

	r.Group(func(r chi.Router) {
		r.Use(authMiddlewares...)

		r.Post("/logout", postLogout)
	})

	r.Route("/me", func(r chi.Router) {
		r.Use(authMiddlewares...)

		r.Get("/", getMe)
		r.Put("/", updateMe)
		r.Put("/icon", updateMyIcon)
	})

	r.Route("/user", func(r chi.Router) {

		r.Post("/", createUser)

		r.Group(func(r chi.Router) {
			r.Use(authMiddlewares...)

			r.Get("/", searchUser)
			r.Get("/{userID}", getUser)
			r.Get("/icon/{userID}", getUserIcon)
		})
	})

	r.Route("/schedule", func(r chi.Router) {

		r.Use(authMiddlewares...)

		r.Post("/", createNewSchedule)
		r.Get("/{scheduleID}", getSchedule)
		r.Put("/{scheduleID}", updateSchedule)
	})

	r.Route("/calendar", func(r chi.Router) {
		r.Use(authMiddlewares...)

		r.Get("/{userID}", getCalendar)
	})

	log.Printf("server start at : %v", time.Now())

	if err := http.ListenAndServe(serverEnv.Port, r); err != nil {
		log.Fatalf("failed to serve http server: %v", err)
	}
}

func generateSessionID() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var possibleSessionIDs []string

	for i := 0; i < 10000; i++ {
		b := make([]byte, 128)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}

		var code string
		for _, v := range b {
			code += string(letters[int(v)%len(letters)])
		}

		possibleSessionIDs = append(possibleSessionIDs, code)
	}

	return possibleSessionIDs[mathrand.Intn(len(possibleSessionIDs))], nil
}

func generateULID() string {
	id := ulid.MustNew(ulid.Now(), rand.Reader)
	return id.String()
}

func parseJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := jwtauth.TokenFromCookie(r)
		if tokenString == "" {
			http.Error(w, "ログインしてください", http.StatusUnauthorized)
			return
		}

		token, err := tokenAuth.Decode(tokenString)
		if err != nil {
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}

		jti, ok := token.Get("jti")
		if !ok || !sess.ExistSession(jti.(string)) {
			http.Error(w, "ログインしてください", http.StatusUnauthorized)
			return
		}

		c := context.WithValue(r.Context(), jtiKey, jti.(string))

		id, ok := token.Get(userIDKey)
		if !ok {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		c = context.WithValue(c, userIDKey, id)

		next.ServeHTTP(w, r.WithContext(c))
	})
}

func postInitialize(orecocoURL string, dbEnv *MySQLConnectionEnv) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDir := filepath.Join("..", "..", "sql")
		paths := []string{
			filepath.Join(sqlDir, "r-calendar-0_Schema.sql"),
			filepath.Join(sqlDir, "r-calendar-1_DummyUserData.sql"),
		}

		for _, p := range paths {
			sqlFile, _ := filepath.Abs(p)
			cmdStr := fmt.Sprintf(
				"mysql -h %v -u %v -p%v -P %v %v < %v",
				dbEnv.Host,
				dbEnv.User,
				dbEnv.Password,
				dbEnv.Port,
				dbEnv.DBName,
				sqlFile,
			)
			if err := exec.CommandContext(r.Context(), "bash", "-c", cmdStr).Run(); err != nil {
				log.Printf("Initialize script error: %v", err)
				http.Error(w, "initialize failed", http.StatusInternalServerError)
				return
			}
		}

		u, err := url.Parse(orecocoURL)
		if err != nil {
			log.Printf("url parse failed: %v", err)
			http.Error(w, "initialize failed", http.StatusInternalServerError)
			return
		}
		u.Path = path.Join(u.Path, "initialize")

		resp, err := http.Post(u.String(), "", nil)
		if err != nil {
			log.Printf("orecoco_reserve initialize failed: %v", err)
			http.Error(w, "initialize failed", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("orecoco_reserve initialize failed returned http status %v", resp.StatusCode)
			http.Error(w, "initialize failed", http.StatusInternalServerError)
			return
		}

		var orecocoToken OrecocoReserveInitializeResponse
		err = json.NewDecoder(resp.Body).Decode(&orecocoToken)
		if err != nil {
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}

		orecocoReserveURL = orecocoURL
		orecocoReserveToken = orecocoToken.Token

		sess = &Session{m: make(map[string]SessionData, 0)}

		w.WriteHeader(http.StatusOK)

		res := &InitializeResponse{
			Language: "go",
		}

		render.JSON(w, r, res)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var createUserReq CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserReq)
	if err != nil {
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	emailReg := regexp.MustCompile("^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$")
	if !emailReg.MatchString(createUserReq.Email) {
		http.Error(w, "メールアドレスの形式が不正です", http.StatusBadRequest)
		return
	}

	var u User
	err = db.Get(&u, "SELECT * FROM user WHERE email = ?", createUserReq.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("database error: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if u.Email != "" {
		http.Error(w, "指定されたメールアドレスはすでに利用されています", http.StatusBadRequest)
		return
	}

	passwordReg := regexp.MustCompile("^[a-zA-Z0-9.!@#$%^&-]{8,64}$")
	if !passwordReg.MatchString(createUserReq.Password) {
		http.Error(w, "パスワードは数字、アルファベット、記号(!@#$%^&-)から8~64文字以内で指定してください", http.StatusBadRequest)
		return
	}

	sum256 := sha256.Sum256([]byte(createUserReq.Password))
	sha256Hashed := hex.EncodeToString(sum256[:])

	uid := generateULID()

	_, err = db.Exec("INSERT INTO user(id, email, name, password) VALUES (?, ?, ?, ?)", uid, createUserReq.Email, createUserReq.Name, sha256Hashed)
	if err != nil {
		log.Printf("database error: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	newSessionID, err := generateSessionID()
	if err != nil {
		log.Printf("failed to generate new Session id: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": uid, "jti": newSessionID})
	if err != nil {
		log.Printf("failed to issue new jwt token: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	sess.SetNewSession(newSessionID, SessionData{})

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var user User
	err = db.Get(&user, "SELECT * FROM user WHERE email = ?", loginReq.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("database error: %v", err)
		}
		http.Error(w, "ユーザー名またはパスワードが不正です", http.StatusBadRequest)
		return
	}

	sum256 := sha256.Sum256([]byte(loginReq.Password))
	sha256Hashed := hex.EncodeToString(sum256[:])

	if sha256Hashed != user.Password {
		http.Error(w, "ユーザー名またはパスワードが不正です", http.StatusBadRequest)
		return
	}

	newSessionID, err := generateSessionID()
	if err != nil {
		log.Printf("failed to generate new Session id: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID, "jti": newSessionID})
	if err != nil {
		log.Printf("failed to issue new jwt token: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	sess.SetNewSession(newSessionID, SessionData{})

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
}

func postLogout(w http.ResponseWriter, r *http.Request) {
	sessionKey := r.Context().Value(jtiKey).(string)

	sess.Delete(sessionKey)

	w.WriteHeader(http.StatusCreated)
}

func getMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)

	var me User
	if err := db.Get(&me, "SELECT * FROM user WHERE id = ?", userID); err != nil {
		log.Printf("failed to get me, user id: %v", userID)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	userMe := &GetUserResponse{
		ID:    me.ID,
		Email: me.Email,
		Name:  me.Name,
	}

	if me.ImageBinary != nil {
		userMe.Icon = fmt.Sprintf("/icon/%s", userMe.ID)
	}

	render.JSON(w, r, userMe)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var user User
	err := db.Get(&user, "SELECT * FROM user WHERE id = ?", userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get user: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		} else {
			http.Error(w, "ユーザーが見つかりませんでした", http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	gotUser := GetUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	if user.ImageBinary != nil {
		gotUser.Icon = fmt.Sprintf("/icon/%s", user.ID)
	}

	render.JSON(w, r, gotUser)
}

func getUserIcon(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var user User
	err := db.Get(&user, "SELECT * FROM user WHERE id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "ユーザーが見つかりませんでした", http.StatusNotFound)
			return
		} else {
			log.Printf("failed to get user: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
	}

	if user.ImageBinary == nil {
		http.Error(w, "アイコンが登録されていません", http.StatusNotFound)
		return
	}

	tempFile, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		log.Printf("failed to create tmp file: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(user.ImageBinary)
	if err != nil {
		log.Printf("failed to write to file: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusBadRequest)
		return
	}

	if err = tempFile.Close(); err != nil {
		log.Printf("failed to save tempfile: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusBadRequest)
		return
	}

	mimeType, err := detectMimeType(tempFile.Name())
	if err != nil {
		log.Printf("failed to detect mime type: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Write(user.ImageBinary)
}

func updateMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)

	var user User
	err := db.Get(&user, "SELECT * FROM user WHERE id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "ユーザーが見つかりませんでした", http.StatusNotFound)
			return
		} else {
			log.Printf("failed to get user: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
	}

	var updateUserReq UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&updateUserReq)
	if err != nil {
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if updateUserReq.Name == "" {
		http.Error(w, "ユーザー名を設定してください", http.StatusBadRequest)
		return
	}

	if updateUserReq.Email == "" {
		http.Error(w, "メールアドレスを設定してください", http.StatusBadRequest)
		return
	}

	emailReg := regexp.MustCompile("^[a-zA-Z0-9_.+-]+@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$")
	if !emailReg.MatchString(updateUserReq.Email) {
		http.Error(w, "メールアドレスの形式が不正です", http.StatusBadRequest)
		return
	}

	passwordReg := regexp.MustCompile("^[a-zA-Z0-9.!@#$%^&-]{8,64}$")
	if !passwordReg.MatchString(updateUserReq.Password) {
		http.Error(w, "パスワードは数字、アルファベット、記号(!@#$%^&-)から8~64文字以内で指定してください", http.StatusBadRequest)
		return
	}

	sum256 := sha256.Sum256([]byte(updateUserReq.Password))
	sha256Hashed := hex.EncodeToString(sum256[:])

	_, err = db.Exec("UPDATE user SET name = ?, email = ?, password = ? WHERE id = ?", updateUserReq.Name, updateUserReq.Email, sha256Hashed, user.ID)
	if err != nil {
		log.Printf("failed to update user data: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateMyIcon(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)

	var user User
	err := db.Get(&user, "SELECT * FROM user WHERE id = ?", userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get user: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		} else {
			http.Error(w, "ユーザーが見つかりませんでした", http.StatusNotFound)
			return
		}
	}

	file, _, err := r.FormFile("icon")
	if err != nil {
		log.Printf("failed to get image from request: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		log.Printf("failed to create tmp file: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	img, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("failed to read file from request: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(img)
	if err != nil {
		log.Printf("failed to write to file: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusBadRequest)
		return
	}

	if err = tempFile.Close(); err != nil {
		log.Printf("failed to save tempfile: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusBadRequest)
		return
	}

	mimeType, err := detectMimeType(tempFile.Name())
	if err != nil {
		log.Printf("failed to detect mime type: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if mimeType == "" {
		http.Error(w, "アイコンに指定できるのはjpeg, png, gifまたはbmpの画像ファイルのみです", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE user SET image_binary = ? WHERE id = ?", img, userID)
	if err != nil {
		log.Printf("failed to update user icon: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func searchUser(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("query")

	if searchQuery == "" {
		http.Error(w, "検索条件を指定してください", http.StatusBadRequest)
		return
	}

	searchQuery = searchQuery + "%"

	var users []User
	// ユーザーが作成された順にソートして返す
	err := db.Select(&users, "SELECT * FROM user WHERE email LIKE ? OR name LIKE ? ORDER BY id", searchQuery, searchQuery)
	if err != nil {
		log.Printf("failed to get user: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "ユーザーが見つかりませんでした", http.StatusNoContent)
		return
	}

	var respUsers UserResponses
	for _, user := range users {
		respUser := GetUserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}
		if user.ImageBinary != nil {
			respUser.Icon = fmt.Sprintf("/icon/%s", user.ID)
		}
		respUsers = append(respUsers, &respUser)
	}

	w.WriteHeader(http.StatusOK)

	sort.Sort(respUsers)

	resp := &SearchUserResponse{Users: respUsers}

	render.JSON(w, r, resp)
}

func createNewSchedule(w http.ResponseWriter, r *http.Request) {
	var createReq CreateScheduleRequest
	err := json.NewDecoder(r.Body).Decode(&createReq)
	if err != nil {
		log.Printf("failed to decode create schedule request: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if createReq.Title == "" {
		http.Error(w, "タイトルを設定してください", http.StatusBadRequest)
		return
	}

	if createReq.StartAt == 0 || createReq.EndAt == 0 {
		http.Error(w, "時間の指定が不正です", http.StatusBadRequest)
		return
	}

	if createReq.StartAt >= createReq.EndAt {
		http.Error(w, "終了時間は開始時間よりも後に設定してください", http.StatusBadRequest)
		return
	}

	if len(createReq.Attendees) == 0 {
		http.Error(w, "参加者を指定してください", http.StatusBadRequest)
		return
	}

	attendees := map[string]struct{}{}
	for _, attendeeID := range createReq.Attendees {
		if _, ok := attendees[attendeeID]; ok {
			log.Printf("attendees duplicated")
			http.Error(w, "ユーザーIDに重複が存在しました", http.StatusBadRequest)
			return
		}
		var attendee User
		err := db.Get(&attendee, "SELECT * FROM user WHERE id = ?", attendeeID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "ユーザーを見つけることができませんでした", http.StatusBadRequest)
			} else {
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			}
			return
		}
		attendees[attendeeID] = struct{}{}
	}

	startAt := time.Unix(createReq.StartAt, 0)
	endAt := time.Unix(createReq.EndAt, 0)

	newScheduleID := generateULID()

	tx, err := db.BeginTx(r.Context(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Printf("failed to create new transaction: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	var scheduleAttendeesID string
	for _, attendeeID := range createReq.Attendees {
		scheduleAttendeesID += attendeeID + ","
	}
	scheduleAttendeesID = strings.TrimSuffix(scheduleAttendeesID, ",")

	_, err = tx.Exec("INSERT INTO schedule (id, title, description, schedule_attendee, start_at, end_at) VALUES (?, ?, ?, ?, ?, ?)", newScheduleID, createReq.Title, createReq.Description, scheduleAttendeesID, startAt, endAt)
	if err != nil {
		log.Printf("create new schedule error: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Printf("Rollback error: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if createReq.MeetingRoom != "" {
		createMeetingRoomRequest := ReserveMeetingRoomRequest{
			ScheduleID:    newScheduleID,
			MeetingRoomID: createReq.MeetingRoom,
			StartAt:       createReq.StartAt,
			EndAt:         createReq.EndAt,
		}
		result, err := callOrecocoReserve(http.MethodPost, "room", nil, createMeetingRoomRequest)
		if err != nil {
			log.Printf("failed to request orecoco-reserve: %v", err)
			if err := tx.Rollback(); err != nil {
				log.Printf("rollback error: %v", err)
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}

		if result.StatusCode != http.StatusCreated {
			if err := tx.Rollback(); err != nil {
				log.Printf("rollback error: %v", err)
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}

			http.Error(w, string(result.Body), result.StatusCode)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		if createReq.MeetingRoom != "" {
			deleteReservedRoom := &UpdateMeetingRoomRequest{
				ScheduleID:    newScheduleID,
				MeetingRoomID: "",
			}

			result, err := callOrecocoReserve(http.MethodPost, "room", nil, deleteReservedRoom)
			if err != nil {
				log.Printf("failed to delete meeting room: %v", err)
				log.Printf("result: %v", string(result.Body))
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}

			if result.StatusCode != http.StatusOK {
				log.Printf("failed to delete meeting room, result: %v", string(result.Body))
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}
		}
		log.Printf("new schedule commit error: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	createRes := CreateScheduleResponse{
		ID: newScheduleID,
	}

	render.JSON(w, r, createRes)
}

func getSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := chi.URLParam(r, "scheduleID")

	var schedule Schedule
	err := db.Get(&schedule, "SELECT * FROM schedule WHERE id = ?", scheduleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "スケジュールが見つかりませんでした", http.StatusNotFound)
		} else {
			log.Printf("failed to get schedule: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		}
		return
	}

	scheduleAttendeesID := strings.Split(schedule.ScheduleAttendee, ",")

	attendeeUsers := make(UserResponses, 0, len(scheduleAttendeesID))

	for _, attendeeID := range scheduleAttendeesID {
		var attendeeUser User
		err := db.Get(&attendeeUser, "SELECT * FROM user WHERE id = ?", attendeeID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("failed to get attendeeID id = %v, err = %v", attendeeID, err)
				http.Error(w, "参加者が見つかりませんでした", http.StatusNotFound)
			} else {
				log.Printf("failed to get attendeeID %v", err)
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			}
			return
		}
		respUser := GetUserResponse{
			ID:    attendeeUser.ID,
			Email: attendeeUser.Email,
			Name:  attendeeUser.Name,
		}
		if attendeeUser.ImageBinary != nil {
			respUser.Icon = fmt.Sprintf("/icon/%s", attendeeUser.ID)
		}
		attendeeUsers = append(attendeeUsers, &respUser)
	}

	sort.Sort(attendeeUsers)

	q := url.Values{}

	q.Add("scheduleId", scheduleID)

	result, err := callOrecocoReserve(http.MethodGet, "room", q, nil)
	if err != nil {
		log.Printf("failed to request orecoco-reserve: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	var meetingRoom string
	switch result.StatusCode {
	case http.StatusOK:
		var getMeetingRoomReservationResponse GetMeetingRoomReservationResponse
		err = json.Unmarshal(result.Body, &getMeetingRoomReservationResponse)
		if err != nil {
			log.Printf("failed to decode GetMeetingRoomReservationResponse: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		meetingRoom = getMeetingRoomReservationResponse.RoomID
	case http.StatusNotFound:
		// 予約なし
		meetingRoom = ""
	default:
		http.Error(w, string(result.Body), result.StatusCode)
		return
	}

	getScheduleRes := &GetScheduleResponse{
		ID:          schedule.ID,
		Title:       schedule.Title,
		Description: schedule.Description,
		StartAt:     schedule.StartAt.Unix(),
		EndAt:       schedule.EndAt.Unix(),
		Attendees:   attendeeUsers,
		MeetingRoom: meetingRoom,
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, getScheduleRes)
}

func updateSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := chi.URLParam(r, "scheduleID")

	var updateScheduleRequest UpdateScheduleRequest
	err := json.NewDecoder(r.Body).Decode(&updateScheduleRequest)
	if err != nil {
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if updateScheduleRequest.Title == "" {
		http.Error(w, "タイトルを設定してください", http.StatusBadRequest)
		return
	}

	if updateScheduleRequest.StartAt == 0 || updateScheduleRequest.EndAt == 0 {
		http.Error(w, "時間の指定が不正です", http.StatusBadRequest)
		return
	}

	if updateScheduleRequest.StartAt >= updateScheduleRequest.EndAt {
		http.Error(w, "終了時間は開始時間よりも後に設定してください", http.StatusBadRequest)
		return
	}

	if len(updateScheduleRequest.Attendees) == 0 {
		http.Error(w, "参加者を指定してください", http.StatusBadRequest)
		return
	}

	tx, err := db.BeginTx(r.Context(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Printf("failed to create new transaction: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	attendees := map[string]struct{}{}
	var attendeesID string
	for _, attendeeID := range updateScheduleRequest.Attendees {
		if _, ok := attendees[attendeeID]; ok {
			if err := tx.Rollback(); err != nil {
				log.Printf("rollback error: %v", err)
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}
			log.Printf("attendees duplicated")
			http.Error(w, "ユーザーIDに重複が存在しました", http.StatusBadRequest)
			return
		}
		_, err = tx.Exec("SELECT * FROM user WHERE id=?", attendeeID)
		if err != nil {
			log.Printf("failed to get user: %v", err)
			if err := tx.Rollback(); err != nil {
				log.Printf("rollback error: %v", err)
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "ユーザーが見つかりませんでした", http.StatusBadRequest)
				return
			} else {
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}
		}
		attendees[attendeeID] = struct{}{}
		attendeesID += attendeeID + ","
	}

	attendeesID = strings.TrimSuffix(attendeesID, ",")

	_, err = tx.Exec("UPDATE schedule SET title = ?, description = ?, schedule_attendee = ?, start_at = ?, end_at = ? WHERE id = ?", updateScheduleRequest.Title, updateScheduleRequest.Description, attendeesID, time.Unix(updateScheduleRequest.StartAt, 0), time.Unix(updateScheduleRequest.EndAt, 0), scheduleID)
	if err != nil {
		log.Printf("failed to update schedule: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	result, err := callOrecocoReserve(http.MethodPut, "room", nil, UpdateMeetingRoomRequest{
		ScheduleID:    scheduleID,
		MeetingRoomID: updateScheduleRequest.MeetingRoom,
		StartAt:       updateScheduleRequest.StartAt,
		EndAt:         updateScheduleRequest.EndAt,
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback error: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		log.Printf("failed to request orecoco-reserve: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if result.StatusCode != http.StatusOK {
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback error: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}

		http.Error(w, string(result.Body), result.StatusCode)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("schedule update commit error: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getCalendar(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	dateFromOriginStr := r.URL.Query().Get("date")

	var u User
	err := db.Get(&u, "SELECT * FROM user WHERE id = ?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "ユーザーが見つかりませんでした", http.StatusNotFound)
		} else {
			log.Printf("failed to get user: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		}
		return
	}

	dateFromOrigin, err := strconv.Atoi(dateFromOriginStr)
	if err != nil {
		log.Printf("invalid dateStr requested, requested: %v and error: %v", dateFromOriginStr, err)
		http.Error(w, "指定された日時は不正です", http.StatusBadRequest)
		return
	}

	var participationSchedules []Schedule
	err = db.Select(&participationSchedules, "SELECT * FROM schedule WHERE schedule_attendee LIKE ?", "%"+userID+"%")
	if err != nil {
		log.Printf("failed to get schedule: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	dateRaw := time.Unix(int64(dateFromOrigin*24*60*60), 0)
	startOfTheDay := time.Date(dateRaw.Year(), dateRaw.Month(), dateRaw.Day(), 0, 0, 0, 0, time.UTC)
	endOfTheDay := time.Date(dateRaw.Year(), dateRaw.Month(), dateRaw.Day()+1, 0, 0, 0, 0, time.UTC)

	var participationScheduleIDs []string
	for _, participationSchedule := range participationSchedules {
		if timeIn(participationSchedule.StartAt, participationSchedule.EndAt, startOfTheDay, endOfTheDay) {
			participationScheduleIDs = append(participationScheduleIDs, participationSchedule.ID)
		}
	}

	if len(participationScheduleIDs) == 0 {
		render.JSON(w, r, GetCalendarResponse{
			Date:      int64(dateFromOrigin),
			Schedules: make(GetScheduleResponses, 0),
		})
		return
	}

	var userSchedules []*Schedule
	query, args, err := sqlx.In("SELECT * FROM schedule WHERE id IN (?) ORDER BY start_at ASC, end_at DESC, id ASC", participationScheduleIDs)
	if err != nil {
		log.Printf("failed to create in query at get schedule err = %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	query = db.Rebind(query)
	rows, err := db.Query(query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("faield to get schedule err = %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var schedule Schedule
		err = rows.Scan(&schedule.ID, &schedule.Title, &schedule.Description, &schedule.ScheduleAttendee, &schedule.StartAt, &schedule.EndAt)
		if err != nil {
			log.Printf("failed to decode from rows err = %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		userSchedules = append(userSchedules, &schedule)
	}

	calendar := GetCalendarResponse{
		Date:      int64(dateFromOrigin),
		Schedules: make(GetScheduleResponses, 0, len(userSchedules)),
	}
	for _, userSchedule := range userSchedules {
		scheduleAttendees := strings.Split(userSchedule.ScheduleAttendee, ",")
		attendeeUsers := make(UserResponses, 0, len(scheduleAttendees))
		for _, attendee := range scheduleAttendees {
			var attendeeUser User
			err := db.Get(&attendeeUser, "SELECT * FROM user WHERE id = ?", attendee)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Printf("failed to get attendee id = %v, err = %v", attendee, err)
					http.Error(w, "参加者が見つかりませんでした", http.StatusNotFound)
				} else {
					log.Printf("failed to get attendee %v", err)
					http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				}
				return
			}
			respUser := GetUserResponse{
				ID:    attendeeUser.ID,
				Email: attendeeUser.Email,
				Name:  attendeeUser.Name,
			}
			if attendeeUser.ImageBinary != nil {
				respUser.Icon = fmt.Sprintf("/icon/%s", attendeeUser.ID)
			}
			attendeeUsers = append(attendeeUsers, &respUser)
		}

		sort.Sort(attendeeUsers)

		q := url.Values{}

		q.Add("scheduleId", userSchedule.ID)

		result, err := callOrecocoReserve(http.MethodGet, "room", q, nil)
		if err != nil {
			log.Printf("failed to request orecoco-reserve: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}

		var meetingRoom string
		switch result.StatusCode {
		case http.StatusOK:
			var getMeetingRoomReservationResponse GetMeetingRoomReservationResponse
			err = json.Unmarshal(result.Body, &getMeetingRoomReservationResponse)
			if err != nil {
				log.Printf("failed to decode GetMeetingRoomReservationResponse: %v", err)
				http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
				return
			}
			meetingRoom = getMeetingRoomReservationResponse.RoomID
		case http.StatusNotFound:
			// 予約なし
			meetingRoom = ""
		default:
			http.Error(w, string(result.Body), result.StatusCode)
			return
		}

		calendar.Schedules = append(calendar.Schedules, &GetScheduleResponse{
			ID:          userSchedule.ID,
			Title:       userSchedule.Title,
			Description: userSchedule.Description,
			StartAt:     userSchedule.StartAt.Unix(),
			EndAt:       userSchedule.EndAt.Unix(),
			Attendees:   attendeeUsers,
			MeetingRoom: meetingRoom,
		})
	}

	w.WriteHeader(http.StatusOK)

	render.JSON(w, r, calendar)
}

func timeIn(start, end, targetStartAt, targetEndAt time.Time) bool {
	if targetEndAt.Equal(start) || targetEndAt.Before(start) || targetStartAt.After(end) {
		return false
	}

	return true
}

func detectMimeType(fileName string) (string, error) {
	var mimeType string
	for mime, formatStr := range allowedFileTypeAndMime {
		reg := regexp.MustCompile(formatStr)
		out, err := exec.Command("file", fileName).Output()
		if err != nil {
			return "", err
		}

		if reg.MatchString(string(out)) {
			mimeType = mime
		}
	}

	return mimeType, nil
}

type orecocoReserveResult struct {
	StatusCode int
	Body       []byte
}

func callOrecocoReserve(httpMethod string, endPoint string, query url.Values, payload interface{}) (*orecocoReserveResult, error) {
	u, err := url.Parse(orecocoReserveURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, endPoint)

	var body io.Reader

	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequest(httpMethod, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	req.Header.Set("Authorization", "Bearer "+orecocoReserveToken)

	client := new(http.Client)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &orecocoReserveResult{
		StatusCode: res.StatusCode,
		Body:       responseBody,
	}, nil
}
