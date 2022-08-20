package common

import "golang.org/x/net/context"

const (
	TIMEOUT_KEY = "timeout"
)

type RateLimiter interface {
	//阻塞获取 tokens
	AcquireTokensWithBlock(ctx context.Context, token int64) int64

	//阻塞获取 tokens，最大超时时间为ctx定义的时间
	AcquireTokensWithNonBlock(ctx context.Context, token int64) int64
}

// TransRate trans the rate to multiples of 1000,the production of rate should be division by 1000.
func Trans2HumanRate(rate int64) int64 {
	if rate <= 0 {
		//default rate
		rate = 10 * 1024 * 1024
	}
	rate = (rate/1000 + 1) * 1000
	return int64(rate)
}
