package constant

import (
	"time"
)

const (
	// InitializeTimeout specifies a time limit for single request of `POST /initialize`
	InitializeTimeout time.Duration = 10 * time.Second
	// CompatibilityCheckRequestTimeout specifies a time limit for single request at compatibility check
	CompatibilityCheckRequestTimeout time.Duration = 2 * time.Second
	// CompatibilityCheckTimeout specifies a time limit for compatibility check
	CompatibilityCheckTimeout time.Duration = 5 * time.Second
	// LoadRequestTimeout specifies a time limit for single request at load
	LoadRequestTimeout time.Duration = 2 * time.Second
	// LoadTimeout specifies a time limit for load
	LoadTimeout time.Duration = 60 * time.Second

	// LoadErrorWaitMin specifies the minimum waiting duration when load worker return error
	LoadErrorWaitMin time.Duration = 1000 * time.Millisecond
	// LoadErrorWaitMax specifies the minimum maximum duration when load worker return error
	LoadErrorWaitMax time.Duration = 1500 * time.Millisecond

	// CreateScheduleWaitMin specifies the minimum waiting duration when `createSchedule`` load worker success
	CreateScheduleWaitMin time.Duration = 300 * time.Millisecond
	// CreateScheduleWaitMax specifies the minimum waiting duration when `createSchedule`` load worker success
	CreateScheduleWaitMax time.Duration = 500 * time.Millisecond

	// TypingIntervalMin specifies the minimum interval of typing
	TypingIntervalMin time.Duration = 5 * time.Millisecond
	// TypingIntervalMax specifies the maximum interval of typing
	TypingIntervalMax time.Duration = 15 * time.Millisecond

	// UpdatingIntervalMin specifies the minimum interval of updating user/schedule information
	UpdatingIntervalMin time.Duration = 20 * time.Millisecond
	// UpdatingIntervalMax specifies the maximum interval of updating user/schedule information
	UpdatingIntervalMaX time.Duration = 50 * time.Millisecond

	// LevelConcurrentPost specifies the lower limit of the level to create concurrent post users
	LevelConcurrentPost = 5
	// ConcurrentPostUserNum specifies count of the concurrent post users per single team
	ConcurrentPostUserNum = 2
	// RetryCountOnScheduleConflict specifies count of retry when creating / updating schedule is conflict
	RetryCountOnScheduleConflict = 3

	// SeeCalendarRate specifies the rate for visiting user's calendar page
	SeeCalendarRate = 0.06
	// GetUserIconRate specifies the rate for getting user's icon
	GetUserIconRate = 0.4
	// ScheduleSamplingRate specifies the rate for sampling schedule
	ScheduleSamplingRate = 0.3
)
