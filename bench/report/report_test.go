package report_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
	"github.com/recruit-tech/RISUCON2022Summer/bench/report"
)

var buf bytes.Buffer

func init() {
	logger.Init(&buf, io.Discard, false)
	logger.Public.SetFlags(0)
}

func TestReport_PublicReport(t *testing.T) {
	r := report.Report{
		Pass:     true,
		Reason:   "OK",
		Level:    5,
		ScoreSum: 100,
		FailReport: fails.FailReport{
			Critical: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			},
			Application: []error{
				errors.New("error 4"),
				errors.New("error 5"),
				errors.New("error 6"),
			},
			Trivial: []error{
				errors.New("error 7"),
				errors.New("error 8"),
				errors.New("error 9"),
			},
		},
		Language: "Go",
	}

	_, got := r.PublicReport()
	want := report.PublicReport{
		Pass:  true,
		Score: 100,
		Messages: []string{
			"Load Level: 5",
			"Critical Error: error 1",
			"Critical Error: error 2",
			"Critical Error: error 3",
			"Application Error: error 4",
			"Application Error: error 5",
			"Application Error: error 6",
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func ExampleReport_Send() {
	r := report.Report{
		Pass:     true,
		Reason:   "OK",
		Level:    5,
		ScoreSum: 100,
		FailReport: fails.FailReport{
			Critical: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			},
			Application: []error{
				errors.New("error 4"),
				errors.New("error 5"),
				errors.New("error 6"),
			},
			Trivial: []error{
				errors.New("error 7"),
				errors.New("error 8"),
				errors.New("error 9"),
			},
		},
		Language: "Go",
	}

	r.Send()
	fmt.Println(buf.String())

	// Output:
	// {
	//   "pass": true,
	//   "score": 100,
	//   "messages": [
	//     "Load Level: 5",
	//     "Critical Error: error 1",
	//     "Critical Error: error 2",
	//     "Critical Error: error 3",
	//     "Application Error: error 4",
	//     "Application Error: error 5",
	//     "Application Error: error 6"
	//   ]
	// }
}
