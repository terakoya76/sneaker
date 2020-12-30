package parser_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terakoya76/sneaker/parser"
)

func TestParseCrontab(t *testing.T) {
	crontab := `
# comment
* * * * * taskA

10 0 20 9 * taskB

*/5 */10 */3 */4 * taskC

7,27,47 23,0-7 10,20,30 2,3 * taskD
`

	cases := []struct {
		name     string
		str      string
		expected []*parser.Expression
		err      error
	}{
		{
			name: "crontab",
			str:  crontab,
			expected: []*parser.Expression{
				{
					Min:   "*",
					Hour:  "*",
					Day:   "*",
					Month: "*",
					Wday:  "*",
					Cmd:   "taskA",
				},

				{
					Min:   "10",
					Hour:  "0",
					Day:   "20",
					Month: "9",
					Wday:  "*",
					Cmd:   "taskB",
				},

				{
					Min:   "*/5",
					Hour:  "*/10",
					Day:   "*/3",
					Month: "*/4",
					Wday:  "*",
					Cmd:   "taskC",
				},

				{
					Min:   "7,27,47",
					Hour:  "23,0-7",
					Day:   "10,20,30",
					Month: "2,3",
					Wday:  "*",
					Cmd:   "taskD",
				},
			},
			err: nil,
		},
	}

	for _, c := range cases {
		actual := parser.ParseCrontab(c.str)
		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}

func TestEvaluate(t *testing.T) {
	cases := []struct {
		name     string
		exp      *parser.Expression
		expected parser.ExecutionSchedule
		err      error
	}{
		{
			name: "every *:* on */*",
			exp: &parser.Expression{
				Min:   "*",
				Hour:  "*",
				Day:   "*",
				Month: "*",
				Wday:  "*",
				Cmd:   "task",
			},
			expected: everyMinEveryHourEveryDayEveryMonthSched(1, 1, 1, 1),
			err:      nil,
		},

		{
			name: "every */4:*/5 on */3 */10",
			exp: &parser.Expression{
				Min:   "*/5",
				Hour:  "*/4",
				Day:   "*/10",
				Month: "*/3",
				Wday:  "*",
				Cmd:   "task",
			},
			expected: everyMinEveryHourEveryDayEveryMonthSched(5, 4, 10, 3),
			err:      nil,
		},

		{
			name: "each 7,8,9:10,11,12 on 4,5 20,21,22",
			exp: &parser.Expression{
				Min:   "10,11,12",
				Hour:  "7,8,9",
				Day:   "20,21,22",
				Month: "4,5",
				Wday:  "*",
				Cmd:   "task",
			},
			expected: specMinsSpecHoursSpecDaysSpecMonthesSched([]int{10, 11, 12}, []int{7, 8, 9}, []int{20, 21, 22}, []int{4, 5}),
			err:      nil,
		},

		{
			name: "each 7,8,9:10,11,12 on 4,5 20,21,22",
			exp: &parser.Expression{
				Min:   "10-12",
				Hour:  "7-9",
				Day:   "20-22",
				Month: "4-5",
				Wday:  "*",
				Cmd:   "task",
			},
			expected: specMinsSpecHoursSpecDaysSpecMonthesSched([]int{10, 11, 12}, []int{7, 8, 9}, []int{20, 21, 22}, []int{4, 5}),
			err:      nil,
		},

		{
			name: "every 07:05 on * 1",
			exp: &parser.Expression{
				Min:   "5",
				Hour:  "7",
				Day:   "1",
				Month: "10",
				Wday:  "*",
				Cmd:   "task",
			},
			expected: specMinSpecHourSpecDaySched(5, 7, 1, 10),
			err:      nil,
		},
	}

	for _, c := range cases {
		schedule := parser.InitSchedule()
		actual, err := c.exp.Evaluate(schedule)

		if !assert.Equal(t, c.err, err) {
			t.Errorf("err: %s is failed, expected: %s, actual: %s\n", c.name, c.err, err)
		}

		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}

func everyMinEveryHourEveryDayEveryMonthSched(min, hour, day, month int) parser.ExecutionSchedule {
	s := parser.InitSchedule()
	for l, mon := range s {
		if l%month == 0 {
			for k, d := range mon {
				if k%day == 0 {
					for j, h := range d {
						if j%hour == 0 {
							for i := range h {
								if i%min == 0 {
									s[l][k][j][i] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return s
}

func specMinsSpecHoursSpecDaysSpecMonthesSched(mins, hours, days, monthes []int) parser.ExecutionSchedule {
	s := parser.InitSchedule()

	for l, mon := range s {
		if contain(monthes, l) {
			for k, d := range mon {
				if contain(days, k) {
					for j, h := range d {
						if contain(hours, j) {
							for i := range h {
								if contain(mins, i) {
									s[l][k][j][i] = true
								}
							}
						}
					}
				}
			}
		}
	}

	return s
}

func specMinSpecHourSpecDaySched(min, hour, day, month int) parser.ExecutionSchedule {
	s := parser.InitSchedule()
	s[month][day][hour][min] = true

	return s
}

func contain(a []int, b int) bool {
	for _, item := range a {
		if b == item {
			return true
		}
	}

	return false
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
			name: "evaluate all",
			max:  parser.MaxMins,
			item: "*",
			expected: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
				10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
				20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
				30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
				40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
				50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
			},
			err: nil,
		},

		{
			name:     "evaluate num item",
			max:      parser.MaxMins,
			item:     "24",
			expected: []int{24},
			err:      nil,
		},

		{
			name:     "evaluate num item including value which greater than the max threshold",
			max:      parser.MaxHours,
			item:     "24",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate list item",
			max:      parser.MaxMins,
			item:     "7,24,47",
			expected: []int{7, 24, 47},
			err:      nil,
		},

		{
			name:     "evaluate list item including value which greater than the max threshold",
			max:      parser.MaxHours,
			item:     "7,24,47",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate range item",
			max:      parser.MaxMins,
			item:     "22-26",
			expected: []int{22, 23, 24, 25, 26},
			err:      nil,
		},

		{
			name:     "evaluate range item including value which is greater than the max threshold",
			max:      parser.MaxHours,
			item:     "22-26",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate step item",
			max:      parser.MaxMins,
			item:     "*/30",
			expected: []int{0, 30},
			err:      nil,
		},

		{
			name:     "evaluate step item including value which is greater than the max threshold",
			max:      parser.MaxHours,
			item:     "*/30",
			expected: []int{0},
			err:      nil,
		},

		{
			name: "evaluate combination includes all",
			max:  parser.MaxMins,
			item: "*,24",
			expected: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
				10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
				20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
				30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
				40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
				50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
			},
			err: nil,
		},

		{
			name:     "evaluate combination includes all and number which is greater than the max threshold",
			max:      parser.MaxHours,
			item:     "*,24",
			expected: []int{},
			err:      fmt.Errorf("%s", "given number is exceeded the max threshold"),
		},

		{
			name:     "evaluate combination of number, range, step",
			max:      parser.MaxMins,
			item:     "3,17-19,*/10",
			expected: []int{0, 3, 10, 17, 18, 19, 20, 30, 40, 50},
			err:      nil,
		},
	}

	for _, c := range cases {
		actual, err := parser.EvaluateItem(c.max, c.item)
		if !assert.Equal(t, c.err, err) {
			t.Errorf("err: %s is failed, expected: %s, actual: %s\n", c.name, c.err, err)
		}

		if !assert.Equal(t, c.expected, actual) {
			t.Errorf("case: %s is failed, expected: %+v, actual: %+v\n", c.name, c.expected, actual)
		}
	}
}
