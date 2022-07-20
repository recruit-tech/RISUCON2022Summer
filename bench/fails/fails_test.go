package fails_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
)

func init() {
	logger.Init(io.Discard, io.Discard, true)
}

type reportCount struct {
	c, a, t int
}

func newReportCount(fr fails.FailReport) reportCount {
	return reportCount{
		c: len(fr.Critical),
		a: len(fr.Application),
		t: len(fr.Trivial),
	}
}

func (rc reportCount) String() string {
	return fmt.Sprintf("(c: %d, a: %d, t: %d)", rc.c, rc.a, rc.t)
}

type timeoutError struct{}
type temporaryError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return false }

func (temporaryError) Error() string   { return "temporary" }
func (temporaryError) Timeout() bool   { return false }
func (temporaryError) Temporary() bool { return true }

var (
	errCritical          = fails.Wrap(errors.New("critical error"), fails.CriticalErrorCode)
	errApplication       = fails.Wrap(errors.New("application error"), fails.ApplicationErrorCode)
	errBenchmarker       = fails.Wrap(errors.New("benchmarker error"), fails.BenchmarkerErrorCode)
	errNestedBenchmarker = fails.Wrap(errBenchmarker, fails.ApplicationErrorCode)
	errCanceled          = context.Canceled
	errTimeout           = new(timeoutError)
	errTemporary         = new(temporaryError)
	errUntagged          = errors.New("untagged error")
)

func TestFails_FailReport(t *testing.T) {
	testcases := []struct {
		Name                string
		Errs                []error
		Want                reportCount
		IsAllBenchmarkerErr bool
	}{
		{
			Name:                "critical error",
			Errs:                []error{errCritical},
			Want:                reportCount{c: 1, a: 0, t: 0},
			IsAllBenchmarkerErr: false,
		},
		{
			Name:                "application error",
			Errs:                []error{errApplication},
			Want:                reportCount{c: 0, a: 1, t: 0},
			IsAllBenchmarkerErr: false,
		},
		{
			Name:                "benchmarker error",
			Errs:                []error{errBenchmarker},
			Want:                reportCount{c: 1, a: 0, t: 0},
			IsAllBenchmarkerErr: true,
		},
		{
			Name:                "nested benchmarker error",
			Errs:                []error{errNestedBenchmarker},
			Want:                reportCount{c: 1, a: 0, t: 0},
			IsAllBenchmarkerErr: true,
		},
		{
			Name:                "canceled error",
			Errs:                []error{errCanceled},
			Want:                reportCount{c: 0, a: 0, t: 1},
			IsAllBenchmarkerErr: false,
		},
		{
			Name:                "timeout error",
			Errs:                []error{errTimeout},
			Want:                reportCount{c: 0, a: 0, t: 1},
			IsAllBenchmarkerErr: false,
		},
		{
			Name:                "temporary error",
			Errs:                []error{errTemporary},
			Want:                reportCount{c: 0, a: 0, t: 1},
			IsAllBenchmarkerErr: false,
		},
		{
			Name:                "untagged error",
			Errs:                []error{errUntagged},
			Want:                reportCount{c: 0, a: 0, t: 1},
			IsAllBenchmarkerErr: false,
		},
		{
			Name: "multiple error",
			Errs: []error{
				errCritical, errCritical,
				errApplication, errApplication, errApplication,
				errTemporary, errTimeout, errCanceled},
			Want:                reportCount{c: 2, a: 3, t: 3},
			IsAllBenchmarkerErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			defer fails.Reset()
			for _, err := range tc.Errs {
				fails.Add(err)
			}
			fails.Done()
			fr := fails.GetReport()
			got := newReportCount(fr)

			if got.c != tc.Want.c || got.a != tc.Want.a || got.t != tc.Want.t {
				t.Errorf("want %s, got %s", tc.Want, got)
			}

			if tc.IsAllBenchmarkerErr {
				for _, err := range fr.Critical {
					if !errors.Is(err, fails.ErrBenchmarker) {
						t.Errorf("non ErrBenchmarker is in critical: %v", err)
					}
				}
				for _, err := range fr.Application {
					if !errors.Is(err, fails.ErrBenchmarker) {
						t.Errorf("non ErrBenchmarker is in application: %v", err)
					}
				}
				for _, err := range fr.Trivial {
					if !errors.Is(err, fails.ErrBenchmarker) {
						t.Errorf("non ErrBenchmarker is in trivial: %v", err)
					}
				}
			}
		})
	}
}
