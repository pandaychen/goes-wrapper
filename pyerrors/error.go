package pyerrors

import (
	"strconv"
	"sync/atomic"
)

var (
	_messages atomic.Value           // NOTE: stored map[int]string（全局错误码）
	_codes    = make(map[int]string) // register codes.
)

// Register register ecode message map.
func Register(cm map[int]string) {
	_messages.Store(cm)
}

// Codes ecode error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
}

// Code 是 Codes实例化类型，本项目的错误码就是int，因为Codes实现了Error()方法，所以可以直接当做error返回
// A Code is an int error code spec.
type Code int

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int {
	return int(e)
}

// Message return error message
func (e Code) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}

	//内置map中找不到就返回默认的Error()方法
	return e.Error()
}

// Details return details.
func (e Code) Details() []interface{} {
	return nil
}
