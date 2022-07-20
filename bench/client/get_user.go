package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

func (c *Client) GetUser(ctx context.Context, query string) (*model.GetUserResponse, error) {
	const endpoint = "GET /user"

	req, err := c.ag.GET("/user")
	if err != nil {
		return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	q := url.Values{}
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")

	res, err := c.ag.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", endpoint, err)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		var rb model.GetUserResponse
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
	case http.StatusNoContent:
		return new(model.GetUserResponse), fmt.Errorf("%w: %q", ErrNoContent, rb)

	case http.StatusBadRequest:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrBadRequest, rb)

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrUnauthorized, rb)

	case http.StatusServiceUnavailable:
		return nil, fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return nil, fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
