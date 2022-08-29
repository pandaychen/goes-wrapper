package pytime

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

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

// 时间戳转日期
func Parse2Timestr(p_timestamp interface{}) (string, error) {
	// golang 特定时间格式 2006-01-02 15:04:05
	// 按照美式时间格式非常容易记忆
	// 格式：月/日 小时：分钟：秒 年 外加时区
	// 01/02 03:04:05pm 06 -0700
	var (
		timestamp int64
	)

	switch t := p_timestamp.(type) {
	case int64:
		timestamp = t
	case int:
		timestamp = int64(t)
	case string:
		ti, err := strconv.Atoi(t)
		if err != nil {
			return "", err
		}
		timestamp = int64(ti)
	default:
		return "", errors.New("not support")
	}

	timeFormat := "2006-01-02 15:04:05"
	Loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}

	curt := time.Unix(timestamp, 0).In(Loc)
	//fmt.Println(curt.Unix())
	return curt.Format(timeFormat), nil
}

func GetCurrentNanoTimestampStr() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func GetCurrentMicroTimestampStr() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

//字符串转换为时间间隔
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
