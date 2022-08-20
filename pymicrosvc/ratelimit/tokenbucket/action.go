package tokenbucket

import (
	"context"
	"time"

	"github.com/pandaychen/goes-wrapper/pymath"
	"github.com/pandaychen/goes-wrapper/pymicrosvc/common"
)

// 计算curTime到lastTime一共可以生产的token数（只允许1个goroutine运行，调用方需要加锁）
func (l *TokenbucketLimiter) howManyTokensElapsed(curNanoTime int64) int64 {
	diff := curNanoTime - l.lastTime
	if diff < time.Millisecond.Nanoseconds() {
		return 0
	}
	return diff / (l.windowSize * time.Millisecond.Nanoseconds()) * l.ratePerWindow //时间/单位时间 * 流速
}

//
func (l *TokenbucketLimiter) generateTokens(ctx context.Context, reqTokenNum, tmpCapacity int64) int64 {
	var (
		now, tokenNums, curTotal int64
	)
	now = time.Now().UnixNano()

	tokenNums = l.howManyTokensElapsed(now)
	//桶大小限制
	curTotal = pymath.Miner(tokenNums+l.curBucketNum /*不受桶大小限制的token数目*/, tmpCapacity /*1.桶size OR 2.token请求数*/)
	if curTotal >= reqTokenNum {
		//enough
		l.curBucketNum = curTotal - reqTokenNum
		l.lastTime = now
		return reqTokenNum
	} else {
		//not enough
		timeout := ctx.Value(common.TIMEOUT_KEY)
		if timeout == nil {
			//need blocking
			/*
				duration := l.wait(reqTokenNum - curTotal)
				l.curBucketNum = curTotal - reqTokenNum
				l.lastTime = now + duration
			*/
			l.wait(reqTokenNum - curTotal)
			return l.generateTokens(ctx, reqTokenNum, tmpCapacity)
		} else {
			//
			return -1
		}
	}

	return -1
}

// 计算需要sleep的时间
func (l *TokenbucketLimiter) wait(needTokenNum int64) time.Duration {
	var (
		waitWindowsNum int64
	)
	if needTokenNum <= 0 {
		return time.Duration(0)
	}

	waitWindowsNum = needTokenNum / l.ratePerWindow
	if waitWindowsNum <= 1 {
		waitWindowsNum = 1
	}

	//计算需要等待的时间
	needSleepDuration := time.Duration(waitWindowsNum * l.windowSize * time.Millisecond.Nanoseconds())

	time.Sleep(needSleepDuration)

	return needSleepDuration
}

// 申请token
func (l *TokenbucketLimiter) acquireTokens(ctx context.Context, reqTokenNum int64) int64 {
	var (
		maxToken, result int64
	)
	if l.capacity <= 0 {
		//不限速
		return reqTokenNum
	}

	if reqTokenNum < 1 {
		// no need
		return reqTokenNum
	}

	//reqTokenNum<容量 ==> 容量
	//reqTokenNum>容量 ==> reqTokenNum
	maxToken = pymath.Maxer(l.capacity, reqTokenNum)
	l.Lock()
	defer l.Unlock()
	result = l.generateTokens(ctx, reqTokenNum, maxToken)
	return result
}

func (l *TokenbucketLimiter) AcquireTokensWithBlock(ctx context.Context, token int64) int64 {
	return l.acquireTokens(ctx, token)
}

//TODO:
func (l *TokenbucketLimiter) AcquireTokensWithNonBlock(ctx context.Context, token int64) int64 {
	return -1
}
