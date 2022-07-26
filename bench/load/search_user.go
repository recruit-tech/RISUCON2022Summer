package load

import (
	"context"
	"errors"
	"math/rand"

	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech/RISUCON2022Summer/bench/constant"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
	"golang.org/x/exp/utf8string"
	"golang.org/x/sync/errgroup"
)

func searchUser(ctx context.Context, team *model.Team) error {
	user := team.Pick()
	c, err := client.New(ctx, client.LoaderType)
	if err != nil {
		return err
	}

	if err := c.PostLogin(ctx, user); err != nil {
		return err
	}

	_, err = c.GetMe(ctx)
	if err != nil {
		return err
	}

	someone := team.Pick()
	someone.RLock()
	defer someone.RUnlock()

	var query string
	if rand.Float64() > 0.5 {
		query = randomSubstring(someone.Name)
	} else {
		query = randomSubstring(someone.Email)
	}

	searchResult, err := search(ctx, c, query)
	if err != nil {
		if errors.Is(err, client.ErrNoContent) {
			return errors.New("GET /user: 少なくとも1件は該当するはずのユーザー検索の結果が0件です")
		}
		return err
	}

	eg, childCtx := errgroup.WithContext(ctx)
	for _, user := range searchResult.Users {
		id := user.ID
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

	uid := searchResult.Users[rand.Intn(len(searchResult.Users))].ID
	date := random.Date()
	cal, err := c.GetCalendarUserId(ctx, uid, date)
	if err != nil {
		return err
	}

	eg, childCtx = errgroup.WithContext(ctx)
	for i := 0; i < 3; i++ {
		eg.Go(func() error {
			if len(cal.Schedules) == 0 {
				return nil
			}
			schedule := cal.Schedules[rand.Intn(len(cal.Schedules))]

			_, err := c.GetScheduleId(childCtx, schedule.ID)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if err := c.PostLogout(ctx); err != nil {
		return err
	}

	return nil
}

func search(ctx context.Context, c *client.Client, query string) (*model.GetUserResponse, error) {
	s := utf8string.NewString(query)
	l := s.RuneCount()
	for i := 1; i < l; i++ {
		err := func() error {
			childCtx, cancel := context.WithTimeout(ctx, random.Duration(constant.TypingIntervalMin, constant.TypingIntervalMax))
			defer cancel()
			_, err := c.GetUser(childCtx, s.Slice(0, i))
			if err := childCtx.Err(); err != nil {
				return nil
			}
			return err
		}()
		if err != nil && !errors.Is(err, context.Canceled) {
			return nil, err
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
	}

	users, err := c.GetUser(ctx, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func randomSubstring(str string) string {
	s := utf8string.NewString(str)
	l := rand.Intn(s.RuneCount()) + 1
	return s.Slice(0, l)
}
