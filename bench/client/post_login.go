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

func (c *Client) PostLogin(ctx context.Context, user *model.User) error {
	const endpoint = "POST /login"

	user.RLock()
	defer user.RUnlock()

	b, err := json.Marshal(model.PostLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	req, err := c.ag.POST("/login", bytes.NewBuffer(b))
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
	case http.StatusBadRequest:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrBadRequest, rb)

	case http.StatusServiceUnavailable:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
