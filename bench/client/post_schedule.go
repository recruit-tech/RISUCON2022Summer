package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

func (c *Client) PostSchedule(ctx context.Context, schedule *model.Schedule) (*model.PostScheduleResponse, error) {
	const endpoint = "POST /schedule"

	schedule.RLock()
	defer schedule.RUnlock()

	b, err := json.Marshal(model.PostScheduleRequest{
		Attendees:   schedule.Attendees.IDList(),
		StartAt:     schedule.StartAt,
		EndAt:       schedule.EndAt,
		Title:       schedule.Title,
		Description: schedule.Description,
		MeetingRoom: schedule.MeetingRoom,
	})
	if err != nil {
		return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	req, err := c.ag.POST("/schedule", bytes.NewBuffer(b))
	if err != nil {
		return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.ag.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", endpoint, err)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode == http.StatusCreated {
		var rb model.PostScheduleResponse
		err = json.NewDecoder(res.Body).Decode(&rb)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, fmt.Errorf("%s: %w", endpoint, ctxErr)
			}
			return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
		}

		return &rb, nil
	}

	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	switch res.StatusCode {
	case http.StatusBadRequest:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrBadRequest, rb)

	case http.StatusConflict:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrConflict, rb)

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrUnauthorized, rb)

	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return nil, fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
