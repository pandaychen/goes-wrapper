package xzap

//日志的context定义

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
)
