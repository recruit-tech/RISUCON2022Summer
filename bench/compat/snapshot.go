package compat

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"

	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech/RISUCON2022Summer/bench/snapshot"
)

func snapshotTest(ctx context.Context, team *model.Team) error {
	sss := snapshot.GetSnapshots()

	for _, ss := range sss {
		c, err := client.New(ctx, client.CompatibilityCheckerType)
		if err != nil {
			return err
		}

		if err := c.PostUser(ctx, &ss.User); err != nil {
			return err
		}

		me, err := c.GetMe(ctx)
		if err != nil {
			return err
		}
		snapshot.SetID(ss.User.ID, me.ID)
		ss.User.ID = me.ID
		if !ss.User.IsSame(*me) {
			return errors.New("GET /me: レスポンスがスナップショットにマッチしません")
		}

		if err := c.PostLogout(ctx); err != nil {
			return err
		}

		if err := team.Add(&ss.User); err != nil {
			return fails.Wrap(err, fails.BenchmarkerErrorCode)
		}
	}

	for _, ss := range sss {
		c, err := client.New(ctx, client.CompatibilityCheckerType)
		if err != nil {
			return err
		}

		if err := c.PostLogin(ctx, &ss.User); err != nil {
			return err
		}

		for _, schedule := range ss.Schedules {
			m := schedule.Attendees.Map()
			attendees := make(map[string]*model.User, len(m))
			for id, user := range m {
				user.ID = snapshot.GetID(id)
				attendees[user.ID] = user
			}
			schedule.Attendees = model.NewUserSet(attendees)
			r, err := c.PostSchedule(ctx, &schedule)
			if err != nil {
				return err
			}
			snapshot.SetID(schedule.ID, r.ID)
		}

		if err := c.PostLogout(ctx); err != nil {
			return err
		}
	}

	eg, childCtx := errgroup.WithContext(ctx)
	for _, ss := range sss {
		s := ss
		eg.Go(func() error {
			c, err := client.New(ctx, client.CompatibilityCheckerType)
			if err != nil {
				return err
			}

			if err := c.PostLogin(childCtx, &s.User); err != nil {
				return err
			}

			for _, v := range s.GetUserId {
				id := snapshot.GetID(v.UserId)
				r, err := c.GetUserId(childCtx, id)
				if err != nil {
					return err
				}
				if !convertUserResponseToUser(v.Response).IsSame(*r) {
					return errors.New("GET /user/:id: レスポンスがスナップショットにマッチしません")
				}
			}

			for _, v := range s.GetUser {
				r, err := c.GetUser(childCtx, v.Query)
				if err != nil {
					return err
				}
				if len(v.Response.Users) != len(r.Users) {
					return errors.New("GET /user: レスポンスがスナップショットにマッチしません")
				}
				for i := 0; i < len(v.Response.Users); i++ {
					if !convertUserResponseToUser(*v.Response.Users[i]).IsSame(*r.Users[i]) {
						return errors.New("GET /user: レスポンスがスナップショットにマッチしません")
					}
				}
			}

			for _, v := range s.GetScheduleId {
				v.Response.ID = snapshot.GetID(v.Response.ID)
				for i := 0; i < len(v.Response.Attendees); i++ {
					v.Response.Attendees[i].ID = snapshot.GetID(v.Response.Attendees[i].ID)
				}
				r, err := c.GetScheduleId(childCtx, v.Response.ID)
				if err != nil {
					return err
				}
				if !v.Response.IsSame(*r) {
					return errors.New("GET /schedule/:id: レスポンスがスナップショットにマッチしません")
				}
			}

			for _, v := range s.GetCalendarUserId {
				r, err := c.GetCalendarUserId(childCtx, snapshot.GetID(v.OwnerID), v.Date)
				if err != nil {
					return err
				}

				if v.Response.Date != (*r).Date {
					return errors.New("GET /calendar/:user_id: レスポンスがスナップショットにマッチしません")
				}

				if len(v.Response.Schedules) != len((*r).Schedules) {
					return errors.New("GET /calendar/:user_id: レスポンスがスナップショットにマッチしません")
				}

				for i := 0; i < len(v.Response.Schedules); i++ {
					v.Response.Schedules[i].ID = snapshot.GetID(v.Response.Schedules[i].ID)
					for j := 0; j < len(v.Response.Schedules[i].Attendees); j++ {
						v.Response.Schedules[i].Attendees[j].ID = snapshot.GetID(v.Response.Schedules[i].Attendees[j].ID)
					}
					if !v.Response.Schedules[i].IsSame((*r).Schedules[i]) {
						return errors.New("GET /calendar/:user_id: レスポンスがスナップショットにマッチしません")
					}
				}

			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func convertUserResponseToUser(u model.UserResponse) *model.User {
	user := &model.User{
		ID:    snapshot.GetID(u.ID),
		Email: u.Email,
		Name:  u.Name,
	}
	user.SetHasIcon(u.Icon != "")
	return user
}
