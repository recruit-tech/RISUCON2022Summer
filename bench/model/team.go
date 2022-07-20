package model

import (
	"bytes"
	"errors"
	"math/rand"
	"strconv"
	"sync"
)

type Team struct {
	Name    string
	members []*User
	mu      sync.RWMutex
}

func NewTeam(name string) *Team {
	return &Team{
		Name:    name,
		members: make([]*User, 0, 32),
		mu:      sync.RWMutex{},
	}
}

func (t *Team) String() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var b bytes.Buffer
	b.WriteString(strconv.Quote(t.Name))
	b.WriteString(" (")
	for i, member := range t.members {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Quote(member.Email))
	}
	b.WriteString(")")

	return b.String()
}

func (t *Team) Add(user *User) error {
	if user.ID == "" {
		return errors.New("User.ID is empty")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.members = append(t.members, user)
	return nil
}

func (t *Team) In(id string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, member := range t.members {
		if id == member.ID {
			return true
		}
	}

	return false
}

func (t *Team) Pick() *User {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.members[rand.Intn(len(t.members))]
}

func (t *Team) RandomUserSet() *UserSet {
	t.mu.RLock()
	defer t.mu.RUnlock()

	length := len(t.members)
	if length > 32 {
		length = 32
	}

	key := uint32(rand.Intn((1<<length)-1) + 1)
	m := make(map[string]*User)

	var u *User
	b := uint32(1)
	for i := 0; i < length; i++ {
		if b&key != 0 {
			u = t.members[i]
			m[u.ID] = u
		}

		b <<= 1
	}

	return &UserSet{
		m: m,
	}
}
