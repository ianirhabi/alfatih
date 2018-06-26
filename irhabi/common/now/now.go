// Copyright 2017 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package now

import (
	"errors"
	"regexp"
	"time"
)

var FirstDayMonday bool
var TimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006-1-2 15:4:5", "2006-1-2 15:4", "2006-1-2", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST"}

type Now struct {
	time.Time
}

func (now *Now) BeginningOfMinute() time.Time {
	return now.Truncate(time.Minute)
}

func (now *Now) BeginningOfHour() time.Time {
	return now.Truncate(time.Hour)
}

func (now *Now) BeginningOfDay() time.Time {
	d := time.Duration(-now.Hour()) * time.Hour
	return now.BeginningOfHour().Add(d)
}

func (now *Now) BeginningOfWeek() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if FirstDayMonday {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}

	d := time.Duration(-weekday) * 24 * time.Hour
	return t.Add(d)
}

func (now *Now) BeginningOfMonth() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.Day())+1) * 24 * time.Hour
	return t.Add(d)
}

func (now *Now) BeginningOfQuarter() time.Time {
	month := now.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

func (now *Now) BeginningOfYear() time.Time {
	t := now.BeginningOfDay()
	d := time.Duration(-int(t.YearDay())+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) EndOfMinute() time.Time {
	return now.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}

func (now *Now) EndOfHour() time.Time {
	return now.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

func (now *Now) EndOfDay() time.Time {
	return now.BeginningOfDay().Add(24*time.Hour - time.Nanosecond)
}

func (now *Now) EndOfWeek() time.Time {
	return now.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func (now *Now) EndOfMonth() time.Time {
	return now.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfQuarter() time.Time {
	return now.BeginningOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

func (now *Now) EndOfYear() time.Time {
	return now.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func (now *Now) Monday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+1) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func (now *Now) Sunday() time.Time {
	t := now.BeginningOfDay()
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	} else {
		d := time.Duration(7-weekday) * 24 * time.Hour
		return t.Truncate(time.Hour).Add(d)
	}
}

func (now *Now) EndOfSunday() time.Time {
	return now.Sunday().Add(24*time.Hour - time.Nanosecond)
}

func parseWithFormat(str string) (t time.Time, err error) {
	for _, format := range TimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}

func (now *Now) Parse(strs ...string) (t time.Time, err error) {
	var setCurrentTime bool
	var parseTime []int
	currentTime := []int{now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
	currentLocation := now.Location()

	for _, str := range strs {
		onlyTime := regexp.MustCompile(`^\s*\d+(:\d+)*\s*$`).MatchString(str) // match 15:04:05, 15

		t, err = parseWithFormat(str)
		location := t.Location()
		if location.String() == "UTC" {
			location = currentLocation
		}

		if err == nil {
			parseTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
			onlyTime = onlyTime && (parseTime[3] == 1) && (parseTime[4] == 1)

			for i, v := range parseTime {
				// Don't reset hour, minute, second if it is a time only string
				if onlyTime && i <= 2 {
					continue
				}

				// Fill up missed information with current time
				if v == 0 {
					if setCurrentTime {
						parseTime[i] = currentTime[i]
					}
				} else {
					setCurrentTime = true
				}

				// Default day and month is 1, fill up it if missing it
				if onlyTime {
					if i == 3 || i == 4 {
						parseTime[i] = currentTime[i]
						continue
					}
				}
			}
		}

		if len(parseTime) > 0 {
			t = time.Date(parseTime[5], time.Month(parseTime[4]), parseTime[3], parseTime[2], parseTime[1], parseTime[0], 0, location)
			currentTime = []int{t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
		}
	}
	return
}

func (now *Now) MustParse(strs ...string) (t time.Time) {
	t, err := now.Parse(strs...)
	if err != nil {
		panic(err)
	}
	return t
}

func (now *Now) Between(time1, time2 string) bool {
	restime := now.MustParse(time1)
	restime2 := now.MustParse(time2)
	return now.After(restime) && now.Before(restime2)
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

// MakeRange setup range date from string with format dd/mm/YYYY
func MakeRange(start, end string) (t *TimeRange, e error) {
	format := "02/01/2006"
	loc, _ := time.LoadLocation("Asia/Jakarta")

	t = new(TimeRange)

	if t.Start, e = time.ParseInLocation(format, start, loc); e == nil {
		if t.End, e = time.ParseInLocation(format, end, loc); e == nil {
			t.End = t.End.Add(24*time.Hour - 1*time.Second)
			if t.End.Before(t.Start) {
				e = errors.New("invalid date range, end date cannot be lower then start date")
			}
		}
	}

	return
}
