package compat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
	"golang.org/x/sync/errgroup"
)

func checkTimeIn(ctx context.Context, team *model.Team) error {
	c, err := client.New(ctx, client.CompatibilityCheckerType)
	if err != nil {
		return err
	}

	user := random.User()
	err = c.PostUser(ctx, user)
	if err != nil {
		return err
	}

	me, err := c.GetMe(ctx)
	if err != nil {
		return err
	}
	user.ID = me.ID

	if !user.IsSame(*me) {
		return errors.New("GET /me: 意図しないユーザーを取得しました")
	}

	date := random.Date()
	year, month, day := time.Unix(date*24*60*60, 0).Date()

	shouldNotIncludeSchedules := map[string]*model.Schedule{
		"昨日の予定": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day-1, 0, 0, 0, 0, time.UTC),
			time.Date(year, month, day-1, 0, 30, 0, 0, time.UTC),
		),
		"明日の予定": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day+1, 12, 30, 0, 0, time.UTC),
			time.Date(year, month, day+1, 13, 30, 0, 0, time.UTC),
		),
		"左端範囲外": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day-1, 23, 59, 58, 0, time.UTC),
			time.Date(year, month, day-1, 23, 59, 59, 0, time.UTC),
		),
		"右端範囲外": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC),
			time.Date(year, month, day+1, 0, 0, 1, 0, time.UTC),
		),
	}

	shouldIncludeSchedules := map[string]*model.Schedule{
		"ゆく年くる年": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day-1, 23, 50, 0, 0, time.UTC),
			time.Date(year, month, day, 0, 10, 0, 0, time.UTC),
		),
		"当日の予定": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
			time.Date(year, month, day, 12, 44, 56, 0, time.UTC),
		),
		"深夜作業": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day, 23, 59, 59, 0, time.UTC),
			time.Date(year, month, day+1, 3, 30, 0, 0, time.UTC),
		),
		"ぶっ続け": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day-1, 12, 30, 0, 0, time.UTC),
			time.Date(year, month, day+1, 12, 30, 0, 0, time.UTC),
		),
		"左端範囲内": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day-1, 23, 59, 59, 0, time.UTC),
			time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		),
		"右端範囲内": generateRandomScheduleWithSpecificTime(user, team,
			time.Date(year, month, day, 23, 59, 59, 0, time.UTC),
			time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC),
		),
	}

	eg, childCtx := errgroup.WithContext(ctx)
	for name, schedule := range shouldNotIncludeSchedules {
		n := name
		s := schedule
		s.Title = n
		eg.Go(func() error {
			res, err := c.PostSchedule(childCtx, s)
			if err != nil {
				return err
			}
			s.ID = res.ID
			return nil
		})
	}

	for name, schedule := range shouldIncludeSchedules {
		n := name
		s := schedule
		s.Title = n
		eg.Go(func() error {
			res, err := c.PostSchedule(childCtx, s)
			if err != nil {
				return err
			}
			s.ID = res.ID
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	calendar, err := c.GetCalendarUserId(ctx, user.ID, date)
	if err != nil {
		return err
	}

	includedScheduleIDSet := make(map[string]struct{}, len(calendar.Schedules))
	for _, schedule := range calendar.Schedules {
		includedScheduleIDSet[schedule.ID] = struct{}{}
	}

	for _, schedule := range shouldNotIncludeSchedules {
		if _, include := includedScheduleIDSet[schedule.ID]; include {
			logger.Private.Println(schedule.Title)
			return fmt.Errorf("GET /calendar/:user_id: 意図しないスケジュールが含まれています (date: %d, schedule_id: %s)", date, schedule.ID)
		}
	}

	for _, schedule := range shouldIncludeSchedules {
		if _, include := includedScheduleIDSet[schedule.ID]; !include {
			logger.Private.Println(schedule.Title)
			return fmt.Errorf("GET /calendar/:user_id: 意図したスケジュールが含まれていません (date: %d, schedule_id: %s)", date, schedule.ID)
		}
	}

	err = c.PostLogout(ctx)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomScheduleWithSpecificTime(owner *model.User, team *model.Team, startAt, endAt time.Time) *model.Schedule {
	s := random.Schedule(owner, team)
	s.MeetingRoom = ""
	s.StartAt = int64(startAt.Unix())
	s.EndAt = int64(endAt.Unix())
	return s
}
