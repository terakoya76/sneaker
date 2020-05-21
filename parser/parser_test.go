package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	crontab string = `
# comment
5 0 * * * taskA

10 0 * * * taskB

30 18 * * * taskC

7,27,47 23,0-7 * * * taskD
`
)

func TestParseCrontab(t *testing.T) {
	cases := []struct {
		name     string
		str      string
		expected []*Expression
		err      error
	}{
		{
			name: "crontab",
			str:  crontab,
			expected: []*Expression{
				{
					min:   "5",
					hour:  "0",
					day:   "*",
					month: "*",
					wday:  "*",
					cmd:   "taskA",
				},
				{
					min:   "10",
					hour:  "0",
					day:   "*",
					month: "*",
					wday:  "*",
					cmd:   "taskB",
				},
				{
					min:   "30",
					hour:  "18",
					day:   "*",
					month: "*",
					wday:  "*",
					cmd:   "taskC",
				},
				{
					min:   "7,27,47",
					hour:  "23,0-7",
					day:   "*",
					month: "*",
					wday:  "*",
					cmd:   "taskD",
				},
			},
			err: nil,
		},
	}

	for _, c := range cases {
		actual := ParseCrontab(c.str)
		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}

func TestEvaluate(t *testing.T) {
	s1 := InitSchedule()
	for _, h := range s1 {
		h[5] = true
	}

	s2 := InitSchedule()
	for _, h := range s2 {
		for i := range h {
			if i%5 == 0 {
				h[i] = true
			}
		}
	}

	s3 := InitSchedule()
	for i, h := range s3 {
		if i == 7 {
			h[5] = true
		}
	}

	s4 := InitSchedule()
	for i, h := range s4 {
		if i == 7 {
			for j := range h {
				if j%5 == 0 {
					h[j] = true
				}
			}
		}
	}

	cases := []struct {
		name     string
		exp      *Expression
		expected ExecutionSchedule
		err      error
	}{
		{
			name: "every xx:05",
			exp: &Expression{
				min:   "5",
				hour:  "*",
				day:   "*",
				month: "*",
				wday:  "*",
				cmd:   "task",
			},
			expected: s1,
			err:      nil,
		},

		{
			name: "every 5min",
			exp: &Expression{
				min:   "*/5",
				hour:  "*",
				day:   "*",
				month: "*",
				wday:  "*",
				cmd:   "task",
			},
			expected: s2,
			err:      nil,
		},

		{
			name: "every 07:05",
			exp: &Expression{
				min:   "5",
				hour:  "7",
				day:   "*",
				month: "*",
				wday:  "*",
				cmd:   "task",
			},
			expected: s3,
			err:      nil,
		},

		{
			name: "every 5min on 07:xx",
			exp: &Expression{
				min:   "*/5",
				hour:  "7",
				day:   "*",
				month: "*",
				wday:  "*",
				cmd:   "task",
			},
			expected: s4,
			err:      nil,
		},
	}

	for _, c := range cases {
		schedule := InitSchedule()
		actual, err := c.exp.Evaluate(schedule)
		if !assert.Equal(t, c.err, err) {
			t.Errorf("err: %s is failed, expected: %s, actual: %s\n", c.name, c.err, err)
		}

		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}

		t.Log(actual)
	}
}

func TestEvaluateItem(t *testing.T) {
	cases := []struct {
		name     string
		max      int
		item     string
		expected []int
		err      error
	}{
		{
			name:     "evaluate all",
			max:      60,
			item:     "*",
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
			err:      nil,
		},

		{
			name:     "evaluate num item",
			max:      60,
			item:     "24",
			expected: []int{24},
			err:      nil,
		},

		{
			name:     "evaluate num item including value which greater than the max threshold",
			max:      24,
			item:     "24",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate list item",
			max:      60,
			item:     "7,24,47",
			expected: []int{7, 24, 47},
			err:      nil,
		},

		{
			name:     "evaluate list item including value which greater than the max threshold",
			max:      24,
			item:     "7,24,47",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate range item",
			max:      60,
			item:     "22-26",
			expected: []int{22, 23, 24, 25, 26},
			err:      nil,
		},

		{
			name:     "evaluate range item including value which is greater than the max threshold",
			max:      24,
			item:     "22-26",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate step item",
			max:      60,
			item:     "*/30",
			expected: []int{0, 30},
			err:      nil,
		},

		{
			name:     "evaluate step item including value which is greater than the max threshold",
			max:      24,
			item:     "*/30",
			expected: []int{0},
			err:      nil,
		},

		{
			name:     "evaluate combination includes all",
			max:      60,
			item:     "*,24",
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59},
			err:      nil,
		},

		{
			name:     "evaluate combination includes all and number which is greater than the max threshold",
			max:      24,
			item:     "*,24",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate combination of number, range, step",
			max:      60,
			item:     "3,17-19,*/10",
			expected: []int{0, 3, 10, 17, 18, 19, 20, 30, 40, 50},
			err:      nil,
		},
	}

	for _, c := range cases {
		actual, err := EvaluateItem(c.max, c.item)
		if !assert.Equal(t, c.err, err) {
			t.Errorf("err: %s is failed, expected: %s, actual: %s\n", c.name, c.err, err)
		}

		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}
