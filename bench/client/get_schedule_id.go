package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

func (c *Client) GetScheduleId(ctx context.Context, id string) (*model.GetScheduleIdResponse, error) {
	const endpoint = "GET /schedule/:id"

	req, err := c.ag.GET(fmt.Sprintf("/schedule/%s", id))
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

	if res.StatusCode == http.StatusOK {
		var rb model.GetScheduleIdResponse
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
	case http.StatusNotFound:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrNotFound, rb)

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrUnauthorized, rb)

	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return nil, fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
