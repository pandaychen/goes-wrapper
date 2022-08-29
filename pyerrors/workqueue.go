package pyerrors

var (
	ErrWorkQueHandlerExists = _addWithMsg(-30001, "redislock: lock not held")
)
