package pyerrors

import (
	"strconv"

	"github.com/pkg/errors"
)

// New new a ecode.Codes by int value.
// NOTE: ecode must unique in global, the New will check repeat and then panic.
func New(e int) Code {
	if e <= 0 {
		panic("business ecode must greater than zero")
	}
	return _add(e)
}

// Int parse code int to error.
func Int(i int) Code {
	return Code(i)
}

// String parse code string to error.
func String(e string) Code {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		//注意：字符串错误统一返回ServerErr（服务器错误500）
		return ServerErr
	}
	return Code(i)
}

// Cause 方法：将error类型转换为项目Codes，生成错误必须调用errors包提供的方法生成！
func Cause(e error) Codes {
	if e == nil {
		return OK
	}

	//调用pkg/errors的Cause()方法
	//测试其是否可以转换为Codes类型
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return String(e.Error())
}

//判断错误码是否相等
// `Codes`与`Codes`判断使用：`Equal(ec1, ec2)`
// `Codes`与`error`类型判断使用：`EqualError(ec, err)`。先从error中剥离得到最原始的`Codes`再比较

// Equal equal a and b by code int.
func Equal(a, b Codes) bool {
	if a == nil {
		a = OK
	}
	if b == nil {
		b = OK
	}
	return a.Code() == b.Code()
}

// EqualError equal error
func EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}
