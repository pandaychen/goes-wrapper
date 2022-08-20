package tokenbucket

//令牌桶实现

import (
	"sync"
	"time"

	"github.com/pandaychen/goes-wrapper/pymicrosvc/common"
)

var (
	METRIC_SECOND_PER_MILLISECOND int64 = 1000 //1s=1000ms
)

type TokenbucketLimiter struct {
	sync.Mutex

	capacity      int64 //令牌桶容量
	curBucketNum  int64 //令牌桶当前
	rate          int64 //每秒（second）产生多少token数目
	ratePerWindow int64
	windowSize    int64 //每隔多少ms放一次token

	//window * (rate/1000) 就是每经过window的时间，产生的token总数
	lastTime int64
}

func NewRateLimiter(rate int64, window_size int64) *common.RateLimiter {
	l := TokenbucketLimiter{}

	l.lastTime = time.Now().UnixNano()

	l.capacity = rate
	l.curBucketNum = 0
	l.rate = rate
	l.setWinSpanSize(window_size)
	//计算流速
	l.calculationRate()
	return &l
}

//设置限流窗口大小
func (l *TokenbucketLimiter) setWinSpanSize(window_size int64) {
	if window_size <= 1 {
		l.windowSize = window_size
		return
	}
	if window_size >= 1000 {
		l.windowSize = METRIC_SECOND_PER_MILLISECOND
		return
	}

	l.windowSize = window_size
	return
}

func (l *TokenbucketLimiter) calculationRate() {
	var (
		ratePerWindow int64
	)
	if l.rate == 0 {
		//不限制
		return
	}
	if ratePerWindow = int64(l.rate) / 1000 /*转为ms单位*/ * int64(l.windowSize); ratePerWindow > 0 {
		l.ratePerWindow = ratePerWindow
	} else {
		//ratePerWindow == 0
		l.ratePerWindow = 1
		l.setWinSpanSize(int64(l.ratePerWindow * 1000 / l.rate))
	}
	return
}
