package pyerrors

var (
	ErrWorkQueHandlerExists = _addWithMsg(-30001, "worker handler already registed")

	ErrWorkQueDriverExists = _addWithMsg(-30002, "driver already registed")

	ErrWorkQueBadDriver = _addWithMsg(-30003, "bad driver")

	ErrWorkQueBadDriverConfig = _addWithMsg(-30004, "bad driver config")
)
