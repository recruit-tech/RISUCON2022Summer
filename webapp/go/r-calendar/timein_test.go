package main

import (
	"testing"
	"time"
)

func TestTimeIn(t *testing.T) {
	start := time.Date(2022, 1, 21, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC)

	testCases := map[string]struct {
		targetStartAt time.Time
		targetEndAt   time.Time

		expectVal bool
	}{
		"昨日の予定": {
			targetStartAt: time.Date(2022, 1, 20, 0, 0, 0, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 20, 0, 30, 0, 0, time.UTC),

			expectVal: false,
		},
		"ゆく年くる年": {
			targetStartAt: time.Date(2022, 1, 20, 23, 50, 0, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 21, 0, 10, 0, 0, time.UTC),

			expectVal: true,
		},
		"当日の予定": {
			targetStartAt: time.Date(2022, 1, 21, 0, 0, 0, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 21, 12, 44, 56, 0, time.UTC),

			expectVal: true,
		},
		"深夜作業": {
			targetStartAt: time.Date(2022, 1, 21, 23, 59, 59, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 22, 3, 30, 0, 0, time.UTC),

			expectVal: true,
		},
		"ぶっ続け": {
			targetStartAt: time.Date(2022, 1, 20, 12, 30, 0, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 22, 12, 30, 0, 0, time.UTC),

			expectVal: true,
		},
		"明日の予定": {
			targetStartAt: time.Date(2022, 1, 22, 12, 30, 0, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 22, 13, 30, 0, 0, time.UTC),

			expectVal: false,
		},
		"左端範囲外": {
			targetStartAt: time.Date(2022, 1, 20, 23, 59, 58, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 20, 23, 59, 59, 0, time.UTC),

			expectVal: false,
		},
		"左端範囲内": {
			targetStartAt: time.Date(2022, 1, 20, 23, 59, 59, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 21, 0, 0, 0, 0, time.UTC),

			expectVal: true,
		},
		"右端範囲内": {
			targetStartAt: time.Date(2022, 1, 21, 23, 59, 59, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC),

			expectVal: true,
		},
		"右端範囲外": {
			targetStartAt: time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC),
			targetEndAt:   time.Date(2022, 1, 22, 0, 0, 1, 0, time.UTC),

			expectVal: false,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			res := timeIn(tc.targetStartAt, tc.targetEndAt, start, end)
			if res != tc.expectVal {
				t.Errorf("unexpected output expect %v but got %v", tc.expectVal, res)
			}
		})
	}
}
