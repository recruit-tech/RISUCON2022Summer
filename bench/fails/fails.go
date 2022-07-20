package fails

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/isucon/isucandar/failure"
	"github.com/pkg/errors"
	"golang.org/x/xerrors"

	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
)

const (
	CriticalErrorCode    failure.StringCode = "Critical Error"
	ApplicationErrorCode failure.StringCode = "Application Error"
	BenchmarkerErrorCode failure.StringCode = "Benchmarker Error"
	TrivialErrorCode     failure.StringCode = "Trivial Error"
)

var (
	ErrBenchmarker = xerrors.New("ベンチマーカーの不具合です。運営に連絡してください。")
)

func Wrap(err error, code failure.StringCode) error {
	return failure.NewError(code, errors.WithStack(err))
}

type fails struct {
	errors       *failure.Errors
	appErrCount  uint32
	criticalHook func()
	status       atomic.Value
}

var (
	noop          = func() {}
	statusNotFail = struct{ s string }{"FAIL"}
	statusFail    = struct{ s string }{"NOT FAIL"}

	f = &fails{
		errors:       failure.NewErrors(context.TODO()),
		appErrCount:  0,
		criticalHook: noop,
	}
)

func init() {
	f.status.Store(statusNotFail)
}

func Add(err error) {
	f.errors.Add(err)
	if failure.IsCode(err, CriticalErrorCode) || failure.IsCode(err, BenchmarkerErrorCode) {
		f.status.Store(statusFail)
		f.criticalHook()
	} else if failure.IsCode(err, ApplicationErrorCode) {
		cnt := atomic.AddUint32(&f.appErrCount, 1)
		if cnt >= 10 {
			f.status.Store(statusFail)
			f.criticalHook()
		}
	}
}

func SetCriticalHook(hook func()) {
	f.criticalHook = hook
}

func IsFail() bool {
	return f.status.Load() == statusFail
}

func Done() {
	f.errors.Done()
}

func Reset() {
	f.errors = failure.NewErrors(context.TODO())
	atomic.StoreUint32(&f.appErrCount, 0)
	f.criticalHook = noop
	f.status.Store(statusNotFail)
}

type FailReport struct {
	Critical    []error
	Application []error
	Trivial     []error
}

func GetReport() FailReport {
	critical, application, trivial := make([]error, 0), make([]error, 0), make([]error, 0)

	for _, err := range f.errors.All() {
		if failure.IsCode(err, BenchmarkerErrorCode) {
			id := random.ID()
			logger.Private.Printf("%s: %+v\n", id, err)
			critical = append(critical, fmt.Errorf("%w (%s)", ErrBenchmarker, id))
			continue
		}

		switch failure.StringCode(failure.GetErrorCode(err)) {
		case CriticalErrorCode, BenchmarkerErrorCode:
			critical = append(critical, xerrors.Unwrap(err))

		case ApplicationErrorCode:
			application = append(application, xerrors.Unwrap(err))

		case TrivialErrorCode, failure.CanceledErrorCode, failure.TimeoutErrorCode, failure.TemporaryErrorCode, failure.UnknownErrorCode:
			trivial = append(trivial, xerrors.Unwrap(err))

		default:
			trivial = append(trivial, xerrors.Unwrap(err))
		}
	}

	return FailReport{
		Critical:    critical,
		Application: application,
		Trivial:     trivial,
	}
}
