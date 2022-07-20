package model_test

import (
	"testing"

	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

var (
	user = model.User{
		ID:       "1",
		Name:     "user1",
		Email:    "user1@example.com",
		Password: "password1",
	}

	userResponse1 = model.UserResponse{
		ID:    "1",
		Name:  "user1",
		Email: "user1@example.com",
		Icon:  "",
	}
	userResponse2 = model.UserResponse{
		ID:    "2",
		Name:  "user2",
		Email: "user2@example.com",
		Icon:  "",
	}
)

func TestUser_IsSame(t *testing.T) {
	testcases := []struct {
		Name string
		Arg  model.UserResponse
		Want bool
	}{
		{
			Name: "same user",
			Arg:  userResponse1,
			Want: true,
		},
		{
			Name: "different user",
			Arg:  userResponse2,
			Want: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			got := user.IsSame(tc.Arg)
			if got != tc.Want {
				t.Errorf("want %v, got %v", tc.Want, got)
			}
		})
	}
}
