package model

type PostInitializeResponse struct {
	Language string `json:"language"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Icon  string `json:"icon"`
}

type GetMeResponse = UserResponse

type getUserResponseItem = UserResponse

type GetUserResponse struct {
	Users []*getUserResponseItem `json:"users"`
}

type GetUserIdResponse = UserResponse

type PostScheduleResponse struct {
	ID string `json:"id"`
}

type GetScheduleIdResponse struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	StartAt     int64          `json:"start_at"`
	EndAt       int64          `json:"end_at"`
	Attendees   []UserResponse `json:"attendees"`
	MeetingRoom MeetingRoom    `json:"meeting_room"`
}

func (s *GetScheduleIdResponse) IsSame(sr GetScheduleIdResponse) bool {
	return s.ID == sr.ID &&
		s.Title == sr.Title &&
		s.Description == sr.Description &&
		s.StartAt == sr.StartAt &&
		s.EndAt == sr.EndAt &&
		s.MeetingRoom == sr.MeetingRoom
}

type GetCalendarUserIdResponseItem = GetScheduleIdResponse

type GetCalendarUserIdResponse struct {
	Date      int64                           `json:"date"`
	Schedules []GetCalendarUserIdResponseItem `json:"schedules"`
}
