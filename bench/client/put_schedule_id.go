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

func (c *Client) PutScheduleId(ctx context.Context, schedule *model.Schedule, updates model.PutScheduleIdRequest) error {
	const endpoint = "PUT /schedule/:id"

	type requestBody struct {
		Attendees   []string          `json:"attendees,omitempty"`
		StartAt     int64             `json:"start_at,omitempty"`
		EndAt       int64             `json:"end_at,omitempty"`
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		MeetingRoom model.MeetingRoom `json:"meeting_room,omitempty"`
	}

	var attendees []string = nil
	if updates.Attendees != nil {
		attendees = updates.Attendees.IDList()
	}

	b, err := json.Marshal(requestBody{
		Attendees:   attendees,
		StartAt:     updates.StartAt,
		EndAt:       updates.EndAt,
		Title:       updates.Title,
		Description: updates.Description,
		MeetingRoom: updates.MeetingRoom,
	})
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	schedule.Lock()
	defer schedule.Unlock()

	req, err := c.ag.PUT(fmt.Sprintf("/schedule/%s", schedule.ID), bytes.NewBuffer(b))
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.ag.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", endpoint, err)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		schedule.Attendees = updates.Attendees
		schedule.StartAt = updates.StartAt
		schedule.EndAt = updates.EndAt
		schedule.Title = updates.Title
		schedule.Description = updates.Description
		schedule.MeetingRoom = updates.MeetingRoom
		return nil
	}

	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	switch res.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrBadRequest, rb)

	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrUnauthorized, rb)

	case http.StatusNotFound:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrNotFound, rb)

	case http.StatusConflict:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrConflict, rb)

	case http.StatusServiceUnavailable:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
