package bench

import (
	"context"
	"errors"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/config"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/logger"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/report"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/snapshot"
	"github.com/recruit-tech/RISUCON2022Summer/bench/assets"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

func Setup(ctx context.Context, conf config.Config) (r *report.Report) {
	var language string
	defer func() {
		err := recover()
		if err != nil {
			switch err := err.(type) {
			case string:
				fails.Add(fails.Wrap(xerrors.New(err), fails.BenchmarkerErrorCode))

			case error:
				fails.Add(fails.Wrap(err, fails.BenchmarkerErrorCode))

			case interface{ String() string }:
				fails.Add(fails.Wrap(xerrors.New(err.String()), fails.BenchmarkerErrorCode))

			default:
				fails.Add(fails.Wrap(xerrors.Errorf("%#+v", err), fails.BenchmarkerErrorCode))
			}
		}

		fails.Done()
		defer fails.Reset()

		fr := fails.GetReport()
		if fails.IsFail() {
			r = &report.Report{
				Pass:       false,
				Reason:     "セットアップに失敗しました",
				FailReport: fr,
				Language:   language,
			}
			return
		}

		r = &report.Report{
			Pass:       true,
			Reason:     "セットアップに成功しました",
			FailReport: fr,
			Language:   language,
		}
	}()

	client.SetTargetUrl(conf.TargetUrl)
	logger.Init(conf.PublicWriter, conf.PrivateWriter, conf.Debug)

	eg := errgroup.Group{}
	eg.Go(func() error { return assets.Init() })
	eg.Go(func() error { return snapshot.LoadSnapshots(conf.SnapshotDir) })
	if err := eg.Wait(); err != nil {
		fails.Add(fails.Wrap(err, fails.CriticalErrorCode))
		return
	}

	c, err := client.New(ctx, client.InitializerType)
	if err != nil {
		fails.Add(fails.Wrap(err, fails.CriticalErrorCode))
		return
	}

	initRes, err := c.PostInitialize(ctx)
	if err != nil {
		fails.Add(fails.Wrap(err, fails.CriticalErrorCode))
		return
	}

	language = initRes.Language
	if language == "" {
		fails.Add(fails.Wrap(errors.New("languageが未設定です"), fails.CriticalErrorCode))
		return
	}

	return
}
