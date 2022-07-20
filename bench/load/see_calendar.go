package load

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/constant"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"golang.org/x/sync/errgroup"
)

func seeCalendar(ctx context.Context, team *model.Team) error {
	user := team.Pick()
	c, err := client.New(ctx, client.LoaderType)
	if err != nil {
		return err
	}

	if err := c.PostLogin(ctx, user); err != nil {
		return err
	}

	schedule := random.Schedule(user, team)
	date := schedule.StartAt / int64(24*time.Hour/time.Second)
	for userID := range schedule.Attendees.Map() {
		if rand.Float32() < constant.SeeCalendarRate {
			_, err := c.GetUserId(ctx, userID)
			if err != nil {
				return err
			}

			eg, childCtx := errgroup.WithContext(ctx)

			fetch := func(date int64) error {
				_, err := c.GetCalendarUserId(childCtx, userID, date)
				if err != nil {
					return err
				}

				return nil
			}

			eg.Go(func() error { return fetch(date - 1) })
			eg.Go(func() error { return fetch(date) })
			eg.Go(func() error { return fetch(date + 1) })

			if err := eg.Wait(); err != nil {
				return err
			}
		}
	}

	var retryCount = 0
	for {
		res, err := c.PostSchedule(ctx, schedule)
		if err != nil {
			if errors.Is(err, client.ErrConflict) {
				if retryCount >= constant.RetryCountOnScheduleConflict {
					return nil
				}
				retryCount += 1
				random.ChangeScheduleTimeRange(schedule)
				continue
			}

			return err
		}

		schedule.ID = res.ID
		break
	}

	got, err := c.GetScheduleId(ctx, schedule.ID)
	if err != nil {
		return err
	}
	if !schedule.IsSame(*got) {
		return errors.New("GET /schedule/:id: 意図しないスケジュールの情報を取得しました")
	}

	if rand.Float32() < constant.ScheduleSamplingRate {
		model.ScheduleSample.Add(schedule)
	}

	eg, childCtx := errgroup.WithContext(ctx)
	for _, userID := range schedule.Attendees.IDList() {
		id := userID
		if rand.Float32() < constant.GetUserIconRate {
			eg.Go(func() error {
				err := c.GetUserIconId(childCtx, id)
				if err != nil && !errors.Is(err, client.ErrNotFound) {
					return err
				}

				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if err := c.PostLogout(ctx); err != nil {
		return err
	}

	return nil
}
