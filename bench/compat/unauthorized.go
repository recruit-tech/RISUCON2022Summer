package compat

import (
	"context"
	"errors"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/bench/assets"
)

func checkUnauthorized(ctx context.Context, team *model.Team) error {
	c, err := client.New(ctx, client.CompatibilityCheckerType)
	if err != nil {
		return err
	}

	_, err = c.GetMe(ctx)
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("GET /me: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	err = c.PutMe(ctx, team.Pick(), model.PutMeRequest{})
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("PUT /me: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	err = c.PutMeIcon(ctx, team.Pick(), assets.GetIcon())
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("PUT /me/icon: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	err = c.PostLogout(ctx)
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("POST /logout: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	_, err = c.GetUserId(ctx, "ID")
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("GET /user/:id: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	_, err = c.GetUser(ctx, "QUERY")
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("GET /user: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	_, err = c.PostSchedule(ctx, random.Schedule(team.Pick(), team))
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("POST /schedule: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	_, err = c.GetScheduleId(ctx, "ID")
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("GET /schedule/:id: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	err = c.PutScheduleId(ctx, random.Schedule(team.Pick(), team), model.PutScheduleIdRequest{})
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("PUT /schedule/:id: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	_, err = c.GetCalendarUserId(ctx, "USER_ID", 0)
	if !errors.Is(err, client.ErrUnauthorized) {
		return errors.New("GET /calendar/:user_id: 認証していないリクエストに401以外のstatus codeが返されています")
	}

	return nil
}
