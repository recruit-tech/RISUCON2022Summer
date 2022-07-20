package random

import (
	"math/rand"
	"time"

	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

var (
	step        = int(1 * time.Hour / time.Second)
	minDuration = int(1 * time.Hour / time.Second)
	maxDuration = int(72 * time.Hour / time.Second)
	minDate     = int(time.Date(2022, 7, 3, 0, 0, 0, 0, time.UTC).Unix())
	maxDate     = int(time.Date(2022, 7, 10, 0, 0, 0, 0, time.UTC).Unix())

	specificMeetingRoomRate float32 = 0.6
)

func init() {
	if minDuration%step != 0 {
		panic("`minDuration` must be dividable by `step`")
	}
	if maxDuration%step != 0 {
		panic("`maxDuration` must be dividable by `step`")
	}
	if minDate%step != 0 {
		panic("`minDate` must be dividable by `step`")
	}
	if maxDate%step != 0 {
		panic("`maxDate` must be dividable by `step`")
	}
}

func Date() int64 {
	min := minDate / int(24*time.Hour/time.Second)
	max := maxDate / int(24*time.Hour/time.Second)
	return int64(rand.Intn(max-min) + min)
}

func duration() int64 {
	min := minDuration / step
	max := maxDuration / step
	return int64((min + rand.Intn(max-min)) * step)
}

func starttime() int64 {
	return int64(minDate + rand.Intn((maxDate-minDate-maxDuration)/step)*step)
}

func SetStep(s int) {
	if minDuration%s != 0 {
		panic("`minDuration` must be dividable by `step`")
	}
	if maxDuration%s != 0 {
		panic("`maxDuration` must be dividable by `step`")
	}
	if minDate%s != 0 {
		panic("`minDate` must be dividable by `step`")
	}
	if maxDate%s != 0 {
		panic("`maxDate` must be dividable by `step`")
	}
	step = s
}

func SetDateRange(min, max int) {
	if min%step != 0 {
		panic("`minDate` must be dividable by `step`")
	}
	if max%step != 0 {
		panic("`maxDate` must be dividable by `step`")
	}
	minDate, maxDate = min, max
}

func SetScheduleDuration(min, max int) {
	if min%step != 0 {
		panic("`minDuration` must be dividable by `step`")
	}
	if max%step != 0 {
		panic("`maxDuration` must be dividable by `step`")
	}
	minDuration, maxDuration = min, max
}

// Schedule generate random Schedule.
// ID field is empty string because ID is decided by the webapp server.
func Schedule(owner *model.User, team *model.Team) *model.Schedule {
	attendee := team.RandomUserSet()
	attendee.Add(owner)
	s := &model.Schedule{
		ID:          "",
		Title:       ID(),
		Description: ID(),
		Attendees:   attendee,
	}
	ChangeScheduleTimeRange(s)
	return s
}

func ChangeScheduleTimeRange(s *model.Schedule) {
	s.Lock()
	defer s.Unlock()

	s.StartAt = starttime()
	s.EndAt = s.StartAt + duration()
	s.MeetingRoom = ""
	if rand.Float32() < specificMeetingRoomRate {
		s.MeetingRoom = MeetingRoom()
		timeRange := TimeRange(s.MeetingRoom)
		s.StartAt = timeRange.StartAt
		s.EndAt = timeRange.EndAt
	}
}
