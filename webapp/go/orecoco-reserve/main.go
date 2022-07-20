package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	scheduleIDKey = "scheduleId"
)

var db *sqlx.DB
var token string

type InitializeResponse struct {
	Token string `json:"token"`
}

type ServerEnv struct {
	Port string
}

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type ReserveMeetingRoomRequest struct {
	ScheduleID    string `json:"schedule_id"`
	MeetingRoomID string `json:"meeting_room_id"`
	StartAt       int64  `json:"start_at"`
	EndAt         int64  `json:"end_at"`
}

type GetMeetingRoomReservationResponse struct {
	RoomID string `json:"room_id"`
}

type UpdateMeetingRoomRequest struct {
	ScheduleID    string `json:"schedule_id"`
	MeetingRoomID string `json:"meeting_room_id"`
	StartAt       int64  `json:"start_at"`
	EndAt         int64  `json:"end_at"`
}

type MeetingRoom struct {
	ID      string    `db:"id"`
	RoomID  string    `db:"room_id"`
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
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
		Port:     getEnv("MYSQL_PORT", "3308"),
		User:     getEnv("MYSQL_USER", "r-isucon"),
		DBName:   getEnv("MYSQL_DBNAME", "orecoco-reserve"),
		Password: getEnv("MYSQL_PASS", "r-isucon"),
	}
}

func (mc *MySQLConnectionEnv) OpenDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=UTC", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
	return sqlx.Open("mysql", dsn)
}

func NewServerEnv() *ServerEnv {
	return &ServerEnv{
		Port: ":" + getEnv("ORECOCO_PORT", "3003"),
	}
}

func generateOrecocoToken() string {
	out, err := exec.Command("bash", "-c", "cat /dev/urandom | base64 | fold -w 32 | head -n 1").Output()
	if err != nil {
		log.Fatalf("orecoco-reserve token generate failed: %v", err)
	}
	return string(out)
}

func main() {
	serverEnv := NewServerEnv()
	dbEnv := NewMySQLConnectionEnv()

	var err error
	db, err = dbEnv.OpenDB()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/initialize", postInitialize(dbEnv))

	r.Route("/room", func(r chi.Router) {
		r.Use(parseOrecocoToken)

		r.Post("/", reserveMeetingRoom)
		r.Get("/", getMeetingRoomReservation)
		r.Put("/", updateMeetingRoom)
	})

	log.Printf("server start at: %v", time.Now())

	if err := http.ListenAndServe(serverEnv.Port, r); err != nil {
		log.Fatalf("failed to server http server: %v", err)
	}
}

func parseOrecocoToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		splitedToken := strings.Split(tokenHeader, "Bearer ")
		if len(splitedToken) != 2 {
			log.Printf("invalid token request")
			http.Error(w, "不正なアクセスです", http.StatusUnauthorized)
			return
		}

		if token != splitedToken[1] {
			log.Printf("invalid token")
			http.Error(w, "不正なトークンです", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func postInitialize(dbEnv *MySQLConnectionEnv) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDir := filepath.Join("..", "..", "sql")
		paths := []string{
			filepath.Join(sqlDir, "orecoco-reserve_0_Schema.sql"),
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
				http.Error(w, "initizlize failed", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)

		token = strings.TrimRight(generateOrecocoToken(), "\n")

		res := &InitializeResponse{
			Token: token,
		}

		render.JSON(w, r, res)
	}
}

func reserveMeetingRoom(w http.ResponseWriter, r *http.Request) {
	var reserveMeetingRoomReq ReserveMeetingRoomRequest
	err := json.NewDecoder(r.Body).Decode(&reserveMeetingRoomReq)
	if err != nil {
		log.Printf("failed to decode ReserveMeetingRoomRequest: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	expectRoomExist, err := meetingRoomExist(reserveMeetingRoomReq.MeetingRoomID)
	if err != nil {
		log.Printf("failed to find expected meeting room: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	} else if !expectRoomExist {
		http.Error(w, "会議室が見つかりません", http.StatusBadRequest)
		return
	}

	startTime := time.Unix(reserveMeetingRoomReq.StartAt, 0)
	endTime := time.Unix(reserveMeetingRoomReq.EndAt, 0)

	tx, err := db.BeginTxx(r.Context(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Printf("failed to create new transaction: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	var conflictSchedules []MeetingRoom
	err = tx.Select(&conflictSchedules, "SELECT * FROM meeting_room WHERE room_id = ? AND NOT (start_at >= ? OR end_at <= ?) FOR UPDATE", reserveMeetingRoomReq.MeetingRoomID, endTime, startTime)
	if err == nil && len(conflictSchedules) != 0 {
		if err := tx.Rollback(); err != nil {
			log.Printf("Rollback error: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		http.Error(w, "すでに予定が入っています", http.StatusConflict)
		return
	} else if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("Rollback error: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		log.Printf("failed to get conflicted schedules: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO meeting_room (id, room_id, start_at, end_at) VALUES (?, ?, ?, ?)", reserveMeetingRoomReq.ScheduleID, reserveMeetingRoomReq.MeetingRoomID, startTime, endTime)
	if err != nil {
		log.Printf("create new meeting room error: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Printf("Rollback error: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("meeting room reserve commit error: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	return
}

func getMeetingRoomReservation(w http.ResponseWriter, r *http.Request) {
	scheduleID := r.URL.Query().Get(scheduleIDKey)
	if scheduleID == "" {
		http.Error(w, "スケジュールを指定してください", http.StatusBadRequest)
		return
	}

	var meetingRoom MeetingRoom
	err := db.Get(&meetingRoom, "SELECT * FROM meeting_room WHERE id = ?", scheduleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "指定したスケジュールには部屋が予約されていません", http.StatusNotFound)
			return
		} else {
			log.Printf("failed to get meeting room reservation: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	getMeetingRoomReservationRes := GetMeetingRoomReservationResponse{
		RoomID: meetingRoom.RoomID,
	}

	render.JSON(w, r, getMeetingRoomReservationRes)
}

func updateMeetingRoom(w http.ResponseWriter, r *http.Request) {
	var updateMeetingRoomRequest UpdateMeetingRoomRequest
	err := json.NewDecoder(r.Body).Decode(&updateMeetingRoomRequest)
	if err != nil {
		log.Printf("failed to decode UpdateMeetingRoomRequest: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	if updateMeetingRoomRequest.MeetingRoomID == "" {
		_, err := db.Exec("DELETE FROM meeting_room WHERE id = ?", updateMeetingRoomRequest.ScheduleID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to delete meeting_room: %v", err)
			http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	expectRoomExist, err := meetingRoomExist(updateMeetingRoomRequest.MeetingRoomID)
	if err != nil {
		log.Printf("failed to find expected meeting room: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	} else if !expectRoomExist {
		http.Error(w, "会議室が見つかりません", http.StatusBadRequest)
		return
	}

	tx, err := db.BeginTxx(r.Context(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		log.Printf("failed to create new transaction: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	startTime := time.Unix(updateMeetingRoomRequest.StartAt, 0)
	endTime := time.Unix(updateMeetingRoomRequest.EndAt, 0)

	var conflictSchedules []MeetingRoom
	err = tx.Select(&conflictSchedules, "SELECT * FROM meeting_room WHERE room_id = ? AND NOT (start_at >= ? OR end_at <= ?) AND id != ? FOR UPDATE", updateMeetingRoomRequest.MeetingRoomID, endTime, startTime, updateMeetingRoomRequest.ScheduleID)
	if err == nil && len(conflictSchedules) != 0 {
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback error: %v", err)
		}
		http.Error(w, "すでに予定が入っています", http.StatusConflict)
		return
	} else if err != nil {
		log.Printf("failed to get conflicted schedules: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback error: %v", err)
		}
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM meeting_room WHERE id = ?", updateMeetingRoomRequest.ScheduleID)
	if err != nil {
		log.Printf("failed to delete meeting room: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback error: %v", err)
		}
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	startTime = time.Unix(updateMeetingRoomRequest.StartAt, 0)
	endTime = time.Unix(updateMeetingRoomRequest.EndAt, 0)

	_, err = tx.Exec("INSERT INTO meeting_room (id, room_id, start_at, end_at) VALUES (?, ?, ?, ?)", updateMeetingRoomRequest.ScheduleID, updateMeetingRoomRequest.MeetingRoomID, startTime, endTime)
	if err != nil {
		log.Printf("failed to insert meeting room: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Printf("rollback error: %v", err)
		}
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("meeting room update commit error: %v", err)
		http.Error(w, "サーバー側のエラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func meetingRoomExist(name string) (bool, error) {
	f, err := os.Open("./meeting_room.txt")
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var exist bool
	for scanner.Scan() {
		roomName := strings.TrimRight(scanner.Text(), "\n")
		if roomName == name {
			exist = true
		}
	}

	if err = scanner.Err(); err != nil {
		return false, err
	}

	return exist, nil
}
