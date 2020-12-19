package parser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/terakoya76/sneaker/enum"
)

const (
	CronExpressionPartsNum = 6

	MaxMonthes = 12
	MaxDays    = 31
	MaxHours   = 23
	MaxMins    = 59

	ExcludeMonth = 0
	ExcludeDay   = 0
)

func ParseCrontab(crontabs string) []*Expression {
	result := []*Expression{}
	lines := strings.Split(crontabs, "\n")
	for _, line := range lines {
		if filterNonExpression(line) {
			items := strings.Fields(line)
			exp := Expression{
				min:   items[0],
				hour:  items[1],
				day:   items[2],
				month: items[3],
				wday:  items[4],
				cmd:   strings.Join(items[5:], " "),
			}
			result = append(result, &exp)
		}
	}

	return result
}

func filterNonExpression(line string) bool {
	if strings.HasPrefix(line, "#") {
		return false
	}

	if len(strings.Fields(line)) < CronExpressionPartsNum {
		return false
	}

	return true
}

// NOTE:
// currently not support wday execution schedule
type ExecutionSchedule MonthlySchedule
type MonthlySchedule []DailySchedule
type DailySchedule []HourlySchedule
type HourlySchedule []MinutelySchedule
type MinutelySchedule []bool

func InitSchedule() ExecutionSchedule {
	monthes := make(ExecutionSchedule, MaxMonthes+1)

	for l := 0; l <= MaxMonthes; l++ {
		days := make(DailySchedule, MaxDays+1)

		for k := 0; k <= MaxDays; k++ {
			hours := make(HourlySchedule, MaxHours+1)

			for j := 0; j <= MaxHours; j++ {
				mins := make(MinutelySchedule, MaxMins+1)

				for i := 0; i <= MaxMins; i++ {
					mins[i] = false
				}

				hours[j] = mins
			}

			days[k] = hours
		}

		monthes[l] = days
	}

	return monthes
}

func (es ExecutionSchedule) String() string {
	mons := MonthlySchedule(es)
	return mons.String()
}

func (mons MonthlySchedule) String() string {
	var b strings.Builder

	for l, mon := range mons {
		if l == ExcludeMonth {
			continue
		}

		slice := strings.Split(mon.String(), "\n")
		slice = slice[0 : len(slice)-1]
		for _, str := range slice {
			fmt.Fprintf(&b, "%s ", enum.Month(l))
			fmt.Fprintln(&b, str)
		}
	}

	return b.String()
}

func (ds DailySchedule) String() string {
	var b strings.Builder

	for k, d := range ds {
		if k == ExcludeDay {
			continue
		}

		slice := strings.Split(d.String(), "\n")
		slice = slice[0 : len(slice)-1]
		for _, str := range slice {
			fmt.Fprintf(&b, "%02d, ", k)
			fmt.Fprintln(&b, str)
		}
	}

	return b.String()
}

func (hs HourlySchedule) String() string {
	var b strings.Builder

	for j, h := range hs {
		fmt.Fprintf(&b, "%02dH: ", j)
		for _, m := range h {
			if m {
				fmt.Fprint(&b, "■")
			} else {
				fmt.Fprint(&b, "□")
			}
		}
		fmt.Fprint(&b, "\n")
	}

	return b.String()
}

type Expression struct {
	min   string
	hour  string
	day   string
	month string
	wday  string
	cmd   string
}

func (e *Expression) Evaluate(schedule ExecutionSchedule) (ExecutionSchedule, error) {
	mons, err := EvaluateItem(MaxMonthes, e.month)
	if err != nil {
		return schedule, err
	}

	ds, err := EvaluateItem(MaxDays, e.day)
	if err != nil {
		return schedule, err
	}

	hs, err := EvaluateItem(MaxHours, e.hour)
	if err != nil {
		return schedule, err
	}

	ms, err := EvaluateItem(MaxMins, e.min)
	if err != nil {
		return schedule, err
	}

	for _, mon := range mons {
		for _, d := range ds {
			for _, h := range hs {
				for _, m := range ms {
					schedule[mon][d][h][m] = true
				}
			}
		}
	}

	return schedule, nil
}

func EvaluateItem(max int, item string) ([]int, error) {
	result := []int{}

	singlePartsNum := 1
	listItemPartsNum := 2

	all := false
	parts := strings.Split(item, ",")
	for _, part := range parts {
		arr := strings.Split(part, "/")
		switch len(arr) {
		case singlePartsNum:
			if arr[0] == "*" {
				all = true
			} else if strings.Contains(arr[0], "-") {
				nums, err := evaluteRange(max, arr[0])
				if err != nil {
					return []int{}, fmt.Errorf("%s", err.Error())
				}

				result = append(result, nums...)
			} else {
				num, err := evaluateNum(max, arr[0])
				if err != nil {
					return []int{}, fmt.Errorf("%s", err.Error())
				}

				result = append(result, num)
			}
		case listItemPartsNum:
			nums, err := evaluteStep(max, arr[0], arr[1])
			if err != nil {
				return []int{}, fmt.Errorf("%s", err.Error())
			}

			result = append(result, nums...)
		default:
			return []int{}, fmt.Errorf("%s", "Invalid cron expression")
		}
	}

	if all {
		return evaluteAll(max)
	}

	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result, nil
}

func evaluteAll(max int) ([]int, error) {
	all := make([]int, max+1)
	for i := 0; i <= max; i++ {
		all[i] = i
	}
	return all, nil
}

func evaluateNum(max int, item string) (int, error) {
	num, err := strconv.Atoi(item)
	if err != nil {
		return 0, err
	}

	if num <= max {
		return num, nil
	}

	return 0, fmt.Errorf("%s", "given number is exceeded the max threshold")
}

func evaluteRange(max int, item string) ([]int, error) {
	nums := strings.Split(item, "-")

	begin, err := strconv.Atoi(nums[0])
	if err != nil {
		return []int{}, err
	}

	end, err := strconv.Atoi(nums[1])
	if err != nil {
		return []int{}, err
	}

	rng := []int{}
	for i := begin; i <= end; i++ {
		if i >= max+1 {
			return []int{}, fmt.Errorf("%s", "given number is exceeded the max threshold")
		}
		rng = append(rng, i)
	}

	return rng, nil
}

func evaluteStep(max int, numerator, denominator string) ([]int, error) {
	den, err := strconv.Atoi(denominator)
	if err != nil {
		return []int{}, err
	}

	var rng []int
	if numerator == "*" {
		r, err := evaluteRange(max, fmt.Sprintf("0-%d", max))
		if err != nil {
			return []int{}, err
		}
		rng = r
	} else if strings.Contains(numerator, "-") {
		r, err := evaluteRange(max, numerator)
		if err != nil {
			return []int{}, err
		}
		rng = r
	} else {
		r, err := evaluteRange(max, fmt.Sprintf("%s-%d", numerator, max))
		if err != nil {
			return []int{}, err
		}
		rng = r
	}

	all := []int{}
	for i := rng[0]; i <= rng[len(rng)-1]; i += den {
		all = append(all, i)
	}

	return all, nil
}
