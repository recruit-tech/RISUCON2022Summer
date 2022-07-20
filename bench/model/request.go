package model

type PostLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PutMeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostScheduleRequest struct {
	Attendees   []string    `json:"attendees"`
	StartAt     int64       `json:"start_at"`
	EndAt       int64       `json:"end_at"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	MeetingRoom MeetingRoom `json:"meeting_room"`
}

type PutScheduleIdRequest struct {
	Attendees   *UserSet    `json:"attendees"`
	StartAt     int64       `json:"start_at"`
	EndAt       int64       `json:"end_at"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	MeetingRoom MeetingRoom `json:"meeting_room"`
}
