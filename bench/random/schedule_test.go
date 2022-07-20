package random_test

import (
	"fmt"
	"testing"

	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
)

var (
	team = random.Team()
)

func init() {
	for i := 0; i < 32; i++ {
		u := random.User()
		u.ID = fmt.Sprintf("user%d", i)
		err := team.Add(u)
		if err != nil {
			panic(err)
		}
	}
}

func TestSchedule(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := random.Schedule(team.Pick(), team)

		if s.StartAt >= s.EndAt {
			t.Errorf("start of schedule (%d) is before the schedule end (%d)", s.StartAt, s.EndAt)
		}

		if len(s.Attendees.IDList()) == 0 {
			t.Errorf("attendee of schedule is zero")
		}
	}
}
