package model_test

import (
	"testing"

	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

var (
	schedule = model.Schedule{
		ID:          "1",
		Title:       "title1",
		Description: "description1",
		StartAt:     0,
		EndAt:       1,
	}

	getScheduleIdResponse1 = model.GetScheduleIdResponse{
		ID:          "1",
		Title:       "title1",
		Description: "description1",
		StartAt:     0,
		EndAt:       1,
	}
	getScheduleIdResponse2 = model.GetScheduleIdResponse{
		ID:          "2",
		Title:       "title2",
		Description: "description2",
		StartAt:     1,
		EndAt:       2,
	}
)

func TestSchedule_IsSame(t *testing.T) {
	testcases := []struct {
		Name string
		Arg  model.GetScheduleIdResponse
		Want bool
	}{
		{Name: "same schedule", Arg: getScheduleIdResponse1, Want: true},
		{Name: "different schedule", Arg: getScheduleIdResponse2, Want: false},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			got := schedule.IsSame(tc.Arg)
			if got != tc.Want {
				t.Errorf("want %v, got %v", tc.Want, got)
			}
		})
	}
}

func TestGetScheduleIdResponse_IsSame(t *testing.T) {
	testcases := []struct {
		Name string
		Arg  model.GetScheduleIdResponse
		Want bool
	}{
		{Name: "same schedule", Arg: getScheduleIdResponse1, Want: true},
		{Name: "different schedule", Arg: getScheduleIdResponse2, Want: false},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			got := getScheduleIdResponse1.IsSame(tc.Arg)
			if got != tc.Want {
				t.Errorf("want %v, got %v", tc.Want, got)
			}
		})
	}
}
