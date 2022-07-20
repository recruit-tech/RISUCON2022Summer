package model

import (
	"math/rand"
	"sync"
)

type MeetingRoom = string

const (
	AlphaMeetingRoom        MeetingRoom = "ALPHA"
	BravoMeetingRoom        MeetingRoom = "BRAVO"
	CharlieMeetingRoom      MeetingRoom = "CHARLIE"
	DeltaMeetingRoom        MeetingRoom = "DELTA"
	EchoMeetingRoom         MeetingRoom = "ECHO"
	FoxtrotMeetingRoom      MeetingRoom = "FOXTROT"
	GolfMeetingRoom         MeetingRoom = "GOLF"
	HotelMeetingRoom        MeetingRoom = "HOTEL"
	IndiaMeetingRoom        MeetingRoom = "INDIA"
	JulietMeetingRoom       MeetingRoom = "JULIET"
	KiloMeetingRoom         MeetingRoom = "KILO"
	LimaMeetingRoom         MeetingRoom = "LIMA"
	MikenovemberMeetingRoom MeetingRoom = "MIKENOVEMBER"
	NovemberMeetingRoom     MeetingRoom = "NOVEMBER"
	OscarMeetingRoom        MeetingRoom = "OSCAR"
	PapaMeetingRoom         MeetingRoom = "PAPA"
	QuebecMeetingRoom       MeetingRoom = "QUEBEC"
	RomeoMeetingRoom        MeetingRoom = "ROMEO"
	SierraMeetingRoom       MeetingRoom = "SIERRA"
	TangoMeetingRoom        MeetingRoom = "TANGO"
	UniformMeetingRoom      MeetingRoom = "UNIFORM"
	VictorMeetingRoom       MeetingRoom = "VICTOR"
	WhiskeyMeetingRoom      MeetingRoom = "WHISKEY"
	XrayMeetingRoom         MeetingRoom = "XRAY"
	YankeeMeetingRoom       MeetingRoom = "YANKEE"
	ZuluMeetingRoom         MeetingRoom = "ZULU"
)

type Schedule struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Attendees   *UserSet    `json:"attendees"`
	StartAt     int64       `json:"start_at"`
	EndAt       int64       `json:"end_at"`
	MeetingRoom MeetingRoom `json:"meeting_room"`

	mu sync.RWMutex
}

func (s *Schedule) Lock()    { s.mu.Lock() }
func (s *Schedule) Unlock()  { s.mu.Unlock() }
func (s *Schedule) RLock()   { s.mu.RLock() }
func (s *Schedule) RUnlock() { s.mu.RUnlock() }

func (s *Schedule) IsSame(sr GetScheduleIdResponse) bool {
	s.RLock()
	defer s.RUnlock()

	return s.ID == sr.ID &&
		s.Title == sr.Title &&
		s.Description == sr.Description &&
		s.StartAt == sr.StartAt &&
		s.EndAt == sr.EndAt &&
		s.MeetingRoom == sr.MeetingRoom
}

type scheduleStore struct {
	mu    sync.RWMutex
	items []*Schedule
}

var ScheduleSample = &scheduleStore{
	items: nil,
}

func (ss *scheduleStore) Add(s *Schedule) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.items = append(ss.items, s)
}

func (ss *scheduleStore) Pick() *Schedule {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if ss.items == nil {
		return nil
	}

	return ss.items[rand.Intn(len(ss.items))]
}
