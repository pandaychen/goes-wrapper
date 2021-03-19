package pytime

import (
	"context"
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
