package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/recruit-tech/RISUCON2022Summer/bench/assets"
	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
)

func (c *Client) PutMeIcon(ctx context.Context, me *model.User, icon assets.Icon) error {
	const endpoint = "PUT /me/icon"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("icon", "icon.jpg")
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	_, err = part.Write(icon)
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	contentType := writer.FormDataContentType()

	err = writer.Close()
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	me.Lock()
	defer me.Unlock()

	req, err := c.ag.PUT("/me/icon", &body)
	if err != nil {
		return fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	req.Header.Set("Content-Type", contentType)

	res, err := c.ag.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", endpoint, err)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	if res.StatusCode == http.StatusOK {
		me.SetHasIcon(true)
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

	case http.StatusServiceUnavailable:
		return fmt.Errorf("%s: %w with %q", endpoint, ErrServiceUnavailable, rb)

	default:
		return fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
	}
}
