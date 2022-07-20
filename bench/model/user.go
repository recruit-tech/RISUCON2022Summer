package model

import (
	"encoding/json"
	"sync"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	hasIcon bool
	mu      sync.RWMutex
}

func (u *User) Lock()                   { u.mu.Lock() }
func (u *User) Unlock()                 { u.mu.Unlock() }
func (u *User) RLock()                  { u.mu.RLock() }
func (u *User) RUnlock()                { u.mu.RUnlock() }
func (u *User) SetHasIcon(hasIcon bool) { u.hasIcon = hasIcon }

func (u *User) GetName() string {
	u.RLock()
	defer u.RUnlock()

	return u.Name
}

func (u *User) GetEmail() string {
	u.RLock()
	defer u.RUnlock()

	return u.Email
}

func (u *User) IsSame(ur UserResponse) bool {
	u.RLock()
	defer u.RUnlock()

	return u.ID == ur.ID &&
		u.Name == ur.Name &&
		u.Email == ur.Email &&
		u.hasIcon == (ur.Icon != "")
}

type UserSet struct {
	m map[string]*User
}

func NewUserSet(m map[string]*User) *UserSet {
	return &UserSet{m: m}
}

func (us UserSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(us.m)
}

func (us *UserSet) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &us.m)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserSet) Add(u *User) {
	if u != nil {
		us.m[u.ID] = u
	}
}

func (us UserSet) Map() map[string]*User {
	return us.m
}

func (us UserSet) IDList() []string {
	l := make([]string, 0, len(us.m))
	for id := range us.m {
		l = append(l, id)
	}
	return l
}
