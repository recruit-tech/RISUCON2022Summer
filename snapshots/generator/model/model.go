package model

import "github.com/recruit-tech/RISUCON2022Summer/bench/model"

type GetMe struct {
	Response model.GetMeResponse `json:"response"`
}

type GetUserId struct {
	UserId   string                  `json:"user_id"`
	Response model.GetUserIdResponse `json:"response"`
}

type GetUser struct {
	Query    string                `json:"query"`
	Response model.GetUserResponse `json:"response"`
}

type GetScheduleId struct {
	ScheduleId string                      `json:"schedule_id"`
	Response   model.GetScheduleIdResponse `json:"response"`
}

type GetCalendarUserId struct {
	OwnerID  string                          `json:"owner_id"`
	Date     int64                           `json:"date"`
	Response model.GetCalendarUserIdResponse `json:"response"`
}

type Snapshot struct {
	User      model.User       `json:"user"`
	Schedules []model.Schedule `json:"schedules"`

	GetMe             GetMe               `json:"get_me"`
	GetUserId         []GetUserId         `json:"get_user_id"`
	GetUser           []GetUser           `json:"get_user"`
	GetScheduleId     []GetScheduleId     `json:"get_schedule"`
	GetCalendarUserId []GetCalendarUserId `json:"get_calendar_user_id"`
}
