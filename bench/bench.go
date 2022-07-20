package bench

import (
	"context"

	"github.com/recruit-tech/RISUCON2022Summer/bench/compat"
	"github.com/recruit-tech/RISUCON2022Summer/bench/config"
	"github.com/recruit-tech/RISUCON2022Summer/bench/load"
	"github.com/recruit-tech/RISUCON2022Summer/bench/report"
	"github.com/recruit-tech/RISUCON2022Summer/bench/score"
)

func Run(conf config.Config) bool {
	var (
		r    *report.Report
		lang string
	)
	defer func() {
		if r == nil {
			panic("report is not found")
		}

		r.Language = lang
		r.Level = score.Level()
		r.ScoreSum = score.Sum()
		if r.ScoreSum < 0 || !r.Pass {
			r.ScoreSum = 0
			r.Pass = false
		}
		r.Webhook = conf.Webhook

		err := r.Send()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	r = Setup(ctx, conf)
	if r == nil || !r.Pass {
		return false
	}
	lang = r.Language

	r = compat.Check(ctx)
	if r == nil || !r.Pass {
		return false
	}

	r = load.Run(ctx)
	return r == nil || r.Pass
}
