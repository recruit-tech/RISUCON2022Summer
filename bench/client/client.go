package client

import (
	"context"
	"fmt"

	"github.com/isucon/isucandar/agent"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
)

type Client struct {
	ag *agent.Agent
}

var targetUrl string

func SetTargetUrl(url string) {
	targetUrl = url
}

func New(ctx context.Context, ct ClientType) (*Client, error) {
	genOpts, found := clientOptsGeneratorMap[ct]
	if !found {
		err := fmt.Errorf("not supported ClientType: %s", ct)
		return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	ag, err := agent.NewAgent(genOpts()...)
	if err != nil {
		return nil, fails.Wrap(err, fails.BenchmarkerErrorCode)
	}

	return &Client{
		ag: ag,
	}, nil
}
