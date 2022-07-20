package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
)

func (c *Client) GetUserIconId(ctx context.Context, id string) error {
	const endpoint = "GET /user/icon/:id"

	req, err := c.ag.GET(fmt.Sprintf("/user/icon/%s", id))
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	res, err := c.ag.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", endpoint, err)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		return nil
	}

	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrNotFound, rb)

	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrUnauthorized, rb)

	case http.StatusServiceUnavailable:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
