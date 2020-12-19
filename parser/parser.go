package parser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	CronExpressionPartsNum = 6
	MaxHour                = 23
	MaxMin                 = 59
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
// currently only support daily execution schedule
type ExecutionSchedule [][]bool

func InitSchedule() ExecutionSchedule {
	hours := make([][]bool, MaxHour+1)
	for i := 0; i <= MaxHour; i++ {
		mins := make([]bool, MaxMin+1)
		for j := 0; j <= MaxMin; j++ {
			mins[j] = false
		}
		hours[i] = mins
	}

	return hours
}

func (es ExecutionSchedule) String() string {
	var b strings.Builder

	for i, h := range es {
		fmt.Fprintf(&b, "%02d: ", i)
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
	hs, err := EvaluateItem(MaxHour, e.hour)
	if err != nil {
		return schedule, err
	}

	ms, err := EvaluateItem(MaxMin, e.min)
	if err != nil {
		return schedule, err
	}

	for _, h := range hs {
		for _, m := range ms {
			schedule[h][m] = true
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
