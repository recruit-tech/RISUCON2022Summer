package load

import (
	"context"
	"net"
	"time"

	"github.com/recruit-tech/RISUCON2022Summer/bench/client"
	"github.com/recruit-tech/RISUCON2022Summer/bench/constant"
	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
	"github.com/recruit-tech/RISUCON2022Summer/bench/report"
	"github.com/recruit-tech/RISUCON2022Summer/bench/score"
	"github.com/recruit-tech/RISUCON2022Summer/bench/worker"
	"golang.org/x/xerrors"
)

func Run(ctx context.Context) (r *report.Report) {
	ctx, cancel := context.WithTimeout(ctx, constant.LoadTimeout)
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
				Reason:     "負荷走行に失敗しました",
				FailReport: fr,
			}
			return
		}

		r = &report.Report{
			Pass:       true,
			Reason:     "負荷走行に成功しました",
			FailReport: fr,
		}
	}()

	RunNewLoader(ctx, 1)

	func() {
		for {
			select {
			case level := <-score.LevelUp():
				RunNewLoader(ctx, level)
				logger.Private.Println("負荷レベルが上昇しました")
				continue

			case <-ctx.Done():
				return
			}
		}
	}()

	return
}

func RunNewLoader(ctx context.Context, level int64) {
	team, err := com.establishTeam(ctx, level >= constant.LevelConcurrentPost)
	if err != nil {
		if err := ctx.Err(); err != nil {
			return
		}
		var nerr net.Error
		if xerrors.As(err, &nerr) && (nerr.Timeout()) {
			fails.Add(nerr)
			return
		}
		if xerrors.Is(err, client.ErrServiceUnavailable) {
			fails.Add(fails.Wrap(err, fails.TrivialErrorCode))
			return
		}
		fails.Add(fails.Wrap(err, fails.CriticalErrorCode))
		return
	}

	var (
		seeCalendarFunc = func(ctx context.Context) {
			if err := seeCalendar(ctx, team); err != nil {
				if err := ctx.Err(); err != nil {
					return
				}
				var nerr net.Error
				if xerrors.As(err, &nerr) && (nerr.Timeout()) {
					fails.Add(nerr)
					return
				}
				if xerrors.Is(err, client.ErrServiceUnavailable) {
					fails.Add(fails.Wrap(err, fails.TrivialErrorCode))
					return
				}
				fails.Add(fails.Wrap(err, fails.ApplicationErrorCode))
				wait(ctx, random.Duration(constant.LoadErrorWaitMin, constant.LoadErrorWaitMax))
				return
			}

			score.Increment()
		}
		searchUserFunc = func(ctx context.Context) {
			if err := searchUser(ctx, team); err != nil {
				if err := ctx.Err(); err != nil {
					return
				}
				var nerr net.Error
				if xerrors.As(err, &nerr) && (nerr.Timeout()) {
					fails.Add(nerr)
					return
				}
				if xerrors.Is(err, client.ErrServiceUnavailable) {
					fails.Add(fails.Wrap(err, fails.TrivialErrorCode))
					return
				}
				fails.Add(fails.Wrap(err, fails.ApplicationErrorCode))
				wait(ctx, random.Duration(constant.LoadErrorWaitMin, constant.LoadErrorWaitMax))
				return
			}

			score.Increment()
		}
		updateDataFunc = func(ctx context.Context) {
			if err := UpdateData(ctx, team); err != nil {
				if err := ctx.Err(); err != nil {
					return
				}
				var nerr net.Error
				if xerrors.As(err, &nerr) && (nerr.Timeout()) {
					fails.Add(nerr)
					return
				}
				if xerrors.Is(err, client.ErrServiceUnavailable) {
					fails.Add(fails.Wrap(err, fails.TrivialErrorCode))
					return
				}
				fails.Add(fails.Wrap(err, fails.ApplicationErrorCode))
				return
			}
		}
		createScheduleFunc = func(ctx context.Context) {
			if err := createSchedule(ctx, team); err != nil {
				if err := ctx.Err(); err != nil {
					return
				}
				var nerr net.Error
				if xerrors.As(err, &nerr) && (nerr.Timeout()) {
					fails.Add(nerr)
					return
				}
				if xerrors.Is(err, client.ErrServiceUnavailable) {
					fails.Add(fails.Wrap(err, fails.TrivialErrorCode))
					return
				}
				fails.Add(fails.Wrap(err, fails.ApplicationErrorCode))
				wait(ctx, random.Duration(constant.LoadErrorWaitMin, constant.LoadErrorWaitMax))
				return
			}

			wait(ctx, random.Duration(constant.CreateScheduleWaitMin, constant.CreateScheduleWaitMax))
		}
	)

	go worker.Process(ctx, seeCalendarFunc)
	go worker.Process(ctx, searchUserFunc)
	go worker.Process(ctx, updateDataFunc)
	go worker.Process(ctx, createScheduleFunc)
}

func wait(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	select {
	case <-t.C:
		return

	case <-ctx.Done():
		t.Stop()
		return
	}
}
