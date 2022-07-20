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

func (c *Client) PostInitialize(ctx context.Context) (*model.PostInitializeResponse, error) {
	const endpoint = "POST /initialize"

	req, err := c.ag.POST("/initialize", nil)
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
		var rb model.PostInitializeResponse
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

	return nil, fmt.Errorf("%w: %q", newUnexpectedStatusCodeError(endpoint, res.StatusCode), rb)
}
