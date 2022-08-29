package pyerrors

import "fmt"

// 全局错误码，可以被外部引用！NotModified/TemporaryRedirect...
var (
	OK = _add(0) // 正确

	NotModified        = _add(-304) // 木有改动
	TemporaryRedirect  = _add(-307) // 302跳转
	RequestErr         = _add(-400) // 请求错误
	Unauthorized       = _add(-401) // 未认证
	AccessDenied       = _add(-403) // 访问权限不足
	NothingFound       = _add(-404) // 啥都木有
	MethodNotAllowed   = _add(-405) // 不支持该方法
	Conflict           = _add(-409) // 冲突
	Canceled           = _add(-498) // 客户端取消请求
	ServerErr          = _add(-500) // 服务器错误
	ServiceUnavailable = _add(-503) // 过载保护,服务暂不可用（客户端熔断返回此错误）
	Deadline           = _add(-504) // 服务调用超时
	LimitExceed        = _add(-509) // 超出限制             （服务端限流返回此错误）

)

//业务全局错误码
var (
	ErrRedisSample = _addWithMsg(-10001, "go-redis test error msg")
)

// 注册全局错误码
func _add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = ""
	return Int(e)
}

func _addWithMsg(e int, msg string) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = msg
	return Int(e)
}
