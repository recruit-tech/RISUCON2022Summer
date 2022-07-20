package worker

import (
	"context"
)

func Process(ctx context.Context, f func(context.Context)) {
	for {
		c := make(chan struct{})

		go func() {
			f(ctx)
			c <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			return

		case <-c:
		}
	}
}
