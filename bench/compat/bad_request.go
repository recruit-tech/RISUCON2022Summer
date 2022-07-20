package compat

import (
	"context"
	"errors"
	"time"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
)

func checkBadRequest(ctx context.Context, team *model.Team) error {
	c, err := client.New(ctx, client.CompatibilityCheckerType)
	if err != nil {
		return err
	}

	err = c.PostLogin(ctx, random.User())
	if !errors.Is(err, client.ErrBadRequest) {
		return errors.New("不正なログインに誤ったレスポンスが返されました")
	}

	err = c.PostUser(ctx, team.Pick())
	if !errors.Is(err, client.ErrBadRequest) {
		return errors.New("重複したユーザーの作成に誤ったレスポンスが返されました")
	}

	me := team.Pick()
	err = c.PostLogin(ctx, me)
	if err != nil {
		return err
	}

	badSchedule := random.Schedule(me, team)
	badSchedule.StartAt, badSchedule.EndAt = badSchedule.EndAt, badSchedule.StartAt
	_, err = c.PostSchedule(ctx, badSchedule)
	if !errors.Is(err, client.ErrBadRequest) {
		return errors.New("不適切なスケジュールの作成に誤ったレスポンスが返されました")
	}

	_, err = c.GetUser(ctx, "")
	if !errors.Is(err, client.ErrBadRequest) {
		return errors.New("空文字列でのユーザー検索に誤ったレスポンスが返されました")
	}

	schedule := random.Schedule(me, team)
	schedule.MeetingRoom = random.MeetingRoom()
	schedule.StartAt = time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	schedule.EndAt = time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC).Unix()
	_, err = c.PostSchedule(ctx, schedule)
	if err != nil {
		return err
	}
	_, err = c.PostSchedule(ctx, schedule)
	if !errors.Is(err, client.ErrConflict) {
		return errors.New("重複する会議室予約のスケジュール作成に誤ったレスポンスが返されました")
	}

	return nil
}
