package compat

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/recruit-tech/RISUCON2022Summer/bench/constant"
	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/model"
	"github.com/recruit-tech/RISUCON2022Summer/bench/report"
)

func Check(ctx context.Context) (r *report.Report) {
	ctx, cancel := context.WithTimeout(ctx, constant.CompatibilityCheckTimeout)
	defer cancel()

	fails.SetCriticalHook(cancel)
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
				Reason:     "互換性チェックに失敗しました",
				FailReport: fr,
			}
			return
		}

		r = &report.Report{
			Pass:       true,
			Reason:     "互換性チェックに成功しました",
			FailReport: fr,
		}
	}()

	t := model.NewTeam("CCG: Compatibility Checking Group")
	err := snapshotTest(ctx, t)
	if err != nil {
		if err := ctx.Err(); err != nil {
			fails.Add(fails.Wrap(errors.New("互換性チェックがタイムアウトしました"), fails.CriticalErrorCode))
			return
		}
		fails.Add(fails.Wrap(err, fails.CriticalErrorCode))
		return
	}

	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error { return checkSuccess(childCtx, t) })
	eg.Go(func() error { return checkBadRequest(childCtx, t) })
	eg.Go(func() error { return checkUnauthorized(childCtx, t) })
	eg.Go(func() error { return checkTimeIn(childCtx, t) })

	if err := eg.Wait(); err != nil {
		if err := ctx.Err(); err != nil {
			fails.Add(fails.Wrap(errors.New("互換性チェックがタイムアウトしました"), fails.CriticalErrorCode))
			return
		}
		fails.Add(fails.Wrap(err, fails.CriticalErrorCode))
		return
	}

	return
}
