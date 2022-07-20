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

func (c *Client) PutMe(ctx context.Context, me *model.User, updates model.PutMeRequest) error {
	const endpoint = "PUT /me"

	b, err := json.Marshal(updates)
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	me.Lock()
	defer me.Unlock()

	req, err := c.ag.PUT("/me", bytes.NewBuffer(b))
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
		me.Name = updates.Name
		me.Email = updates.Email
		me.Password = updates.Password
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

	case http.StatusServiceUnavailable:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
