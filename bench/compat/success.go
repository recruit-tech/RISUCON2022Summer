package compat

import (
	"context"
	"errors"
	"time"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/bench/assets"
)

func checkSuccess(ctx context.Context, team *model.Team) error {
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

	someone := team.Pick()
	resSomeone, err := c.GetUserId(ctx, someone.ID)
	if err != nil {
		return err
	}

	randomUser := random.User()
	c.PutMe(ctx, user, model.PutMeRequest{
		Name:     randomUser.Name,
		Email:    randomUser.Email,
		Password: randomUser.Password,
	})

	me, err = c.GetMe(ctx)
	if err != nil {
		return err
	}

	if !user.IsSame(*me) {
		return errors.New("GET /me: 意図しないユーザーを取得しました")
	}

	err = c.GetUserIconId(ctx, me.ID)
	if err == nil {
		return errors.New("GET /user/icon/:id: 作成直後のユーザーにアイコンが設定されています")
	}
	if !errors.Is(err, client.ErrNotFound) {
		return err
	}

	err = c.PutMeIcon(ctx, user, assets.GetIcon())
	if err != nil {
		return err
	}

	err = c.GetUserIconId(ctx, me.ID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return errors.New("GET /user/icon/:id: アイコンをセットしたはずのユーザーのアイコンが見つかりません")
		}
		return err
	}

	err = c.PostLogout(ctx)
	if err != nil {
		return err
	}

	err = c.PostLogin(ctx, user)
	if err != nil {
		return err
	}

	me, err = c.GetMe(ctx)
	if err != nil {
		return err
	}

	if !user.IsSame(*me) {
		return errors.New("GET /me: 意図しないユーザーを取得しました")
	}

	if !someone.IsSame(*resSomeone) {
		return errors.New("GET /user/:id: 意図しないユーザーの情報を取得しました")
	}

	schedule := random.Schedule(user, team)
	ps, err := c.PostSchedule(ctx, schedule)
	if err != nil {
		return nil
	}

	if ps.ID == "" {
		return errors.New("POST /schedule: IDが空文字列です")
	}
	schedule.ID = ps.ID

	s, err := c.GetScheduleId(ctx, schedule.ID)
	if err != nil {
		return err
	}

	if !schedule.IsSame(*s) {
		return errors.New("GET /schedule/:id: 意図しないスケジュールの情報を取得しました")
	}

	randomSchedule := random.Schedule(user, team)
	randomSchedule.MeetingRoom = random.MeetingRoom()
	putScheduleReq := model.PutScheduleIdRequest{
		Attendees:   randomSchedule.Attendees,
		StartAt:     randomSchedule.StartAt,
		EndAt:       randomSchedule.EndAt,
		Title:       randomSchedule.Title,
		Description: randomSchedule.Description,
		MeetingRoom: randomSchedule.MeetingRoom,
	}
	err = c.PutScheduleId(ctx, schedule, putScheduleReq)
	if err != nil {
		return err
	}

	s, err = c.GetScheduleId(ctx, schedule.ID)
	if err != nil {
		return err
	}

	if !schedule.IsSame(*s) {
		return errors.New("GET /schedule/:id: 意図しないスケジュールの情報を取得しました")
	}

	err = c.PutScheduleId(ctx, schedule, putScheduleReq)
	if err != nil {
		return errors.New("PUT /schedule/:id: スケジュールの更新に失敗しました")
	}

	searchResult, err := c.GetUser(ctx, someone.Name)
	if errors.Is(err, client.ErrNoContent) {
		return errors.New("GET /user: 少なくとも1件は該当するはずのユーザー検索の結果が0件です")
	}
	if err != nil {
		return err
	}

	include := false
	for _, u := range searchResult.Users {
		if someone.IsSame(*u) {
			include = true
			break
		}
	}
	if !include {
		return errors.New("GET /user: 意図したユーザーが含まれていません")
	}

	date := schedule.StartAt / (int64(24 * time.Hour / time.Second))
	calendar, err := c.GetCalendarUserId(ctx, me.ID, date)
	if err != nil {
		return err
	}

	include = false
	if calendar.Date != date {
		return errors.New("GET /calendar/:user_id: 意図しないカレンダーの情報を取得しました")
	}
	for _, s := range calendar.Schedules {
		if schedule.IsSame(s) {
			include = true
			break
		}
	}
	if !include {
		return errors.New("GET /calendar/:user_id: 意図したスケジュールが含まれていません")
	}

	err = c.PostLogout(ctx)
	if err != nil {
		return err
	}

	return nil
}
