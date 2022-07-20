package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
)

func (c *Client) PostLogout(ctx context.Context) error {
	const endpoint = "POST /logout"

	req, err := c.ag.POST("/logout", nil)
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

	if res.StatusCode == http.StatusCreated {
		return nil
	}

	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrUnauthorized, rb)

	case http.StatusServiceUnavailable:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
