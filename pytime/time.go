package pytime

import (
	"context"
	"os"
	"strconv"
	"time"
)

// Use the long enough past time as start time(one year ago), in case timex.Now() - lastTime equals 0
var startTime = time.Now().AddDate(-1, -1, -1)

// 计算相对时间差
func Duration2Now() time.Duration {
	return time.Since(startTime)
}

func Duration2Fixed(d time.Duration) time.Duration {
	return time.Since(startTime) - d
}

func GetRelativeCurrentTime() time.Time {
	return startTime.Add(Duration2Now())
}

func GetDuration(timestr string) time.Duration {
	dur, err := time.ParseDuration(timestr)
	if err != nil {
		panic(err)
	}

	return dur
}

// Duration like 1s, 500ms
type Duration time.Duration

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < time.Duration(d) {
			//比较d 和 ctimeout，当context的超时较小时取此值
			return Duration(ctimeout), c, func() {}
		}
	}
	//否则，基于d重新生成timeout的context返回
	ctx, cancel := context.WithTimeout(c, time.Duration(d))
	return d, ctx, cancel
}

/*
func main() {
	fmt.Println(int64(Duration2Now()),Duration2Now())	// 34128000000137802 9480h0m0.00013898s
}
*/

type TimeFormat string

// Format 格式化
func (ts TimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}

var (
	//1970-01-01 08:00:00 +0800 CST
	zeroTime = time.Unix(0, 0)

	TS TimeFormat = "2006-01-02 15:04:05"
)

func Duration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return duration
}

func Str2Duration(str string) (time.Duration, error) {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return time.Duration(0), err
	}
	return dur, nil
}

type TimeFormat string

// Format 格式化
func (ts TimeFormat) Format(t time.Time) string {
	return t.Format(string(ts))
}

// ParseInLocation parse time with location from env "TZ", if "TZ" hasn't been set then we use UTC by default.
func ParseInLocation(layout, value string) (time.Time, error) {
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(layout, value, loc)
}

// parseDate returns time.Date given string containing number of days since Jan 1, 1970,or zero value of time.Time if the string is empty.
func ParseDate(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, nil
	}

	days, err := strconv.Atoi(date)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(1970, 0, 0, 0, 0, 0, 0, time.UTC).AddDate(0, 0, days), nil
}

// IsZero reports whether t represents the zero time instant
func IsZero(t time.Time) bool {
	return t.IsZero() || zeroTime.Equal(t)
}
