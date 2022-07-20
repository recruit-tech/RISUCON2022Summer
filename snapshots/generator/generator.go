package generator

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/exp/utf8string"

	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/snapshots/generator/model"
	"golang.org/x/sync/errgroup"
)

const (
	userNum         = 8
	schedulePerUser = 2

	getUserIdPerUser         = 1
	getSchedulePerUser       = 1
	getUserPerUser           = 2
	getCalendarUserIdPerUser = 2
)

func init() {
	if userNum < 2 {
		panic("`userNum` must be greater than or equal to 2")
	}

	if schedulePerUser < 2 {
		panic("`userNum` must be greater than or equal to 2")
	}
}

var (
	step        = int(30 * time.Minute / time.Second)
	minDuration = int(30 * time.Minute / time.Second)
	maxDuration = int(12 * time.Hour / time.Second)
	minDate     = int(time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC).Unix())
	maxDate     = int(time.Date(2022, 7, 3, 0, 0, 0, 0, time.UTC).Unix())
)

func randomSubstring(str string) string {
	s := utf8string.NewString(str)
	l := rand.Intn(s.RuneCount()) + 1
	return s.Slice(0, l)
}

func Run(targetUrl string, minify bool) {
	ctx := context.Background()

	client.SetTargetUrl(targetUrl)
	random.SetDateRange(minDate, maxDate)
	random.SetScheduleDuration(minDuration, maxDuration)
	random.SetStep(step)

	c, err := client.New(ctx, client.CompatibilityCheckerType)
	if err != nil {
		panic(err)
	}

	_, err = c.PostInitialize(ctx)
	if err != nil {
		panic(err)
	}

	var snapshots = make([]model.Snapshot, userNum)

	team := random.Team()
	for i := 0; i < userNum; i++ {
		u := random.User()

		c, err := client.New(ctx, client.CompatibilityCheckerType)
		if err != nil {
			panic(err)
		}

		if err := c.PostUser(ctx, u); err != nil {
			panic(err)
		}

		me, err := c.GetMe(ctx)
		if err != nil {
			panic(err)
		}
		u.ID = me.ID

		if err := c.PostLogout(ctx); err != nil {
			panic(err)
		}

		team.Add(u)
		snapshots[i].User = *u
		snapshots[i].GetMe = model.GetMe{Response: *me}
	}

	for i := 0; i < userNum-1; i++ {
		me := snapshots[i].User

		c, err := client.New(ctx, client.CompatibilityCheckerType)
		if err != nil {
			panic(err)
		}

		if err := c.PostLogin(ctx, &me); err != nil {
			panic(err)
		}

		for j := 0; j < schedulePerUser; j++ {
			schedule := random.Schedule(team)
			schedule.MeetingRoom = ""
			if i == 0 && j == 0 {
				schedule.MeetingRoom = random.MeetingRoom()
			}
			res, err := c.PostSchedule(ctx, schedule)
			if err != nil {
				panic(err)
			}
			schedule.ID = res.ID

			snapshots[i].Schedules = append(snapshots[i].Schedules, *schedule)
		}

		if err := c.PostLogout(ctx); err != nil {
			panic(err)
		}
	}

	err = func() error {
		me := snapshots[userNum-1].User

		c, err := client.New(ctx, client.CompatibilityCheckerType)
		if err != nil {
			return err
		}

		if err := c.PostLogin(ctx, &me); err != nil {
			return err
		}

		for j := 0; j < schedulePerUser; j++ {
			schedule := snapshots[0].Schedules[j]
			schedule.MeetingRoom = ""
			if j%2 == 0 {
				schedule.StartAt -= int64(step)
			} else {
				schedule.EndAt += int64(step)
			}

			res, err := c.PostSchedule(ctx, &schedule)
			if err != nil {
				return err
			}
			schedule.ID = res.ID

			snapshots[userNum-1].Schedules = append(snapshots[userNum-1].Schedules, schedule)
		}

		if err := c.PostLogout(ctx); err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		panic(err)
	}

	eg, childCtx := errgroup.WithContext(ctx)
	for i := 0; i < userNum; i++ {
		idx := i
		eg.Go(func() error {
			me := snapshots[idx].User

			c, err := client.New(childCtx, client.CompatibilityCheckerType)
			if err != nil {
				return err
			}

			if err := c.PostLogin(childCtx, &me); err != nil {
				return err
			}

			for j := 0; j < getUserIdPerUser; j++ {
				uid := team.Pick().ID
				res, err := c.GetUserId(childCtx, uid)
				if err != nil {
					return err
				}

				snapshots[idx].GetUserId = append(snapshots[idx].GetUserId, model.GetUserId{UserId: uid, Response: *res})
			}

			for j := 0; j < getUserPerUser; j++ {
				var query string
				if rand.Float64() > 0.5 {
					query = randomSubstring(team.Pick().Name)
				} else {
					query = randomSubstring(team.Pick().Email)
				}

				res, err := c.GetUser(childCtx, query)
				if err != nil {
					log.Println(query, res, err)
					return err
				}

				snapshots[idx].GetUser = append(snapshots[idx].GetUser, model.GetUser{Query: query, Response: *res})
			}

			for j := 0; j < getSchedulePerUser; j++ {
				sid := snapshots[rand.Intn(userNum)].Schedules[rand.Intn(schedulePerUser)].ID
				res, err := c.GetScheduleId(childCtx, sid)
				if err != nil {
					return err
				}

				snapshots[idx].GetScheduleId = append(snapshots[idx].GetScheduleId, model.GetScheduleId{ScheduleId: sid, Response: *res})
			}

			for j := 0; j < getCalendarUserIdPerUser; j++ {
				oid := snapshots[rand.Intn(userNum)].User.ID
				date := int64(rand.Intn(maxDate-minDate)+minDate) / int64(24*time.Hour/time.Second)
				res, err := c.GetCalendarUserId(childCtx, oid, date)
				if err != nil && !errors.Is(err, client.ErrNoContent) {
					log.Println(res, err)
					return err
				}

				snapshots[idx].GetCalendarUserId = append(snapshots[idx].GetCalendarUserId, model.GetCalendarUserId{OwnerID: oid, Date: date, Response: *res})
			}

			if err := c.PostLogout(childCtx); err != nil {
				return err
			}

			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	if !minify {
		encoder.SetIndent("", "\t")
	}

	if err := encoder.Encode(snapshots); err != nil {
		panic(err)
	}
}
