package pyerrors

import "github.com/go-redis/redis/v8"

var (
	// ErrInvalidParams  is returned when parameters is invalid.
	ErrRedisInvalidParams = _addWithMsg(-10002, "invalid params")

	// ErrNotObtained is returned when a Lock cannot be obtained.
	ErrRedisNotObtained = _addWithMsg(-10003, "redislock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive Lock.
	ErrRedisLockNotHeld = _addWithMsg(-10004, "redislock: lock not held")
)

const (
	//Nil reply returned by Redis when key does not exist.
	RedisNil = redis.Nil
)
