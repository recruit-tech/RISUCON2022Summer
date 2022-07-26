package load

import (
	"context"
	"errors"
	"math/rand"

	"github.com/recruit-tech/RISUCON2022Summer/bench/assets"
	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech/RISUCON2022Summer/bench/constant"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
)

func UpdateData(ctx context.Context, team *model.Team) error {
	user := team.Pick()
	c, err := client.New(ctx, client.LoaderType)
	if err != nil {
		return err
	}

	if err := c.PostLogin(ctx, user); err != nil {
		return err
	}

	switch rand.Intn(3) {
	case 0:
		err = c.PutMeIcon(ctx, user, assets.GetIcon())
		if err != nil {
			return err
		}
	case 1:
		err = updateMe(ctx, c, user)
		if err != nil {
			return err
		}
	case 2:
		schedule := model.ScheduleSample.Pick()
		if schedule == nil {
			break
		}

		var retryCount = 0
		for {
			err = updateSchedule(ctx, c, user, team, schedule)
			if err != nil {
				if errors.Is(err, client.ErrConflict) {
					if retryCount >= constant.RetryCountOnScheduleConflict {
						return nil
					}
					retryCount += 1
					continue
				}

				return err
			}
		}
	}

	return nil
}

func randomPutMeRequest(user *model.User) model.PutMeRequest {
	user.RLock()
	defer user.RUnlock()

	request := model.PutMeRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	randomUser := random.User()
	for i := 0; i < 3; i++ {
		switch rand.Intn(3) {
		case 0:
			request.Name = randomUser.Name
		case 1:
			request.Email = randomUser.Email
		case 2:
			request.Password = randomUser.Password
		}
	}
	return request
}

func updateMe(ctx context.Context, c *client.Client, me *model.User) error {
	updates := randomPutMeRequest(me)
	return c.PutMe(ctx, me, updates)
}

func randomPutScheduleIdRequest(schedule *model.Schedule, owner *model.User, team *model.Team) model.PutScheduleIdRequest {
	schedule.RLock()
	defer schedule.RUnlock()

	request := model.PutScheduleIdRequest{
		Attendees:   schedule.Attendees,
		StartAt:     schedule.StartAt,
		EndAt:       schedule.EndAt,
		Title:       schedule.Title,
		Description: schedule.Description,
		MeetingRoom: schedule.MeetingRoom,
	}

	randomSchedule := random.Schedule(owner, team)
	for i := 0; i < 3; i++ {
		switch rand.Intn(5) {
		case 0:
			request.Title = randomSchedule.Title
		case 1:
			request.Description = randomSchedule.Description
		case 2:
			request.Attendees = randomSchedule.Attendees
		case 3:
			request.StartAt = randomSchedule.StartAt
			request.EndAt = randomSchedule.EndAt
		case 4:
			request.MeetingRoom = randomSchedule.MeetingRoom
		}
	}

	return request
}

func updateSchedule(ctx context.Context, c *client.Client, owner *model.User, team *model.Team, schedule *model.Schedule) error {
	updates := randomPutScheduleIdRequest(schedule, owner, team)
	return c.PutScheduleId(ctx, schedule, updates)
}
