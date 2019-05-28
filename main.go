package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	_year = iota
	_mon
	_day
	_hour
	_minute
	_second
)

type date struct {
	int
	flag int
	nxt  *[]date
}

func add_time(target *[]date, in time.Time, time_scale time.Duration) {

	add := func(target *[]date, input int, id int) *date {
		for _, v := range *target {
			if input == v.int {
				return &v
			}
		}
		n := date{int: input, nxt: &[]date{}, flag: id}
		*target = append(*target, n)
		return &n
	}

	year := add(target, in.Year(), _year)
	month := add(year.nxt, int(in.Month()), _mon)
	day := add(month.nxt, in.Day(), _day)
	hour := add(day.nxt, in.Hour(), _hour)
	minute := add(hour.nxt, in.Minute(), _minute)
	if time_scale == time.Second {
		add(minute.nxt, in.Second(), _second)
	}
}

func create_regex(dates *[]date, output []string) []string {

	timepend := func(input date) {

		pad := func(num int) string {
			if num < 10 {
				return fmt.Sprintf("0%d", num)
			}
			return fmt.Sprintf("%d", num)
		}

		output = append(output, fmt.Sprintf("%s", pad(input.int)))

	}

	skip_nested := func(input date) bool {
		switch input.flag {
			case _year:
				if len(*input.nxt) == 12 {
					return true
				}
			case _mon:
				switch input.int {
					case 2:
						if len(*input.nxt) >= 28 {
							return true
						}
					case 4:
						fallthrough
					case 6:
						fallthrough
					case 9:
						fallthrough
					case 11:
						if len(*input.nxt) == 30 {
							return true
						}
					default:
						if len(*input.nxt) == 31 {
							return true
						}
				}
				if len(*input.nxt) == 31 {
					return true
				}
			case _day:
				if len(*input.nxt) == 24 {
					return true
				}
			case _hour:
				fallthrough
			case _minute:
				if len(*input.nxt) == 60 {
					return true
				}
		}
		return false
	}

	date_len := len(*dates)

	for i, v := range *dates {
		switch i {
		case 0:
			switch v.flag {
			case _mon:
				fallthrough
			case _day:
				output = append(output, "-")
			case _hour:
				output = append(output, "(T| )")
			case _minute:
				fallthrough
			case _second:
				output = append(output, ":")
			}
			if date_len > 1 {
				output = append(output, "(")
			}
			timepend(v)
			if !skip_nested(v) {
				output = create_regex(v.nxt, output)
			}
			if i < date_len-1 {
				output = append(output, "|")
			}
		case date_len - 1:
			timepend(v)
			if !skip_nested(v) {
				output = create_regex(v.nxt, output)
			}
			if date_len > 1 {
				output = append(output, ")")
			}
		default:
			timepend(v)
			if i > 0 && i < date_len-1 {
				output = append(output, "|")
			}

		}
	}
	return output
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("syntax: %s <time duration>\n", os.Args[0])
		os.Exit(1)
	}

	time_sel, err := time.ParseDuration(os.Args[1])
	if err != nil {
		fmt.Printf("Err:%s\n", err.Error())
		os.Exit(0)
	}

	current := time.Now().UTC()
	past := current.Add(time_sel * -1)
	fmt.Printf("current time: %v\n", current)
	fmt.Printf("   past time: %v\n", past)

	var years []date
	var time_scale time.Duration

	if time_sel >= time.Minute {
		time_scale = time.Minute
	} else {
		time_scale = time.Second
	}

	for ; current.Sub(past) >= time_scale; past = past.Add(time_scale) {
		add_time(&years, past, time_scale)
	}

	add_time(&years, current, time_scale)

	fmt.Println(strings.Join(create_regex(&years, []string{}), ""))

}
