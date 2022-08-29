package pytime

import "time"

var (
	//1970-01-01 08:00:00 +0800 CST
	zeroTime = time.Unix(0, 0)

	TS TimeFormat = "2006-01-02 15:04:05"
)
