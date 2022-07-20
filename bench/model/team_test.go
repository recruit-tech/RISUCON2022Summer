package model_test

import (
	"fmt"
	"testing"

	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/random"
)

var (
	team = model.NewTeam("TEST")
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

func TestTeam_RandomUserSet(t *testing.T) {
	for i := 0; i < 1; i++ {
		us := team.RandomUserSet()
		ids := us.IDList()

		if len(ids) == 0 {
			t.Errorf("attendee of schedule is zero")
		}

		for _, id := range ids {
			if !team.In(id) {
				t.Errorf("user %q is not in team: %s", id, team.String())
			}
		}
	}
}
