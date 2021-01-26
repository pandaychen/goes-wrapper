package pytime

import (
	"context"
	"time"
)

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
