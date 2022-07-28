package common

import "golang.org/x/net/context"

type RateLimiter interface {
	//阻塞获取 tokens
	AcquireTokensWithBlock(ctx context.Context, token int64) int64

	//阻塞获取 tokens，最大超时时间为ctx定义的时间
	AcquireTokensWithNonBlock(ctx context.Context, token int64) int64
}
