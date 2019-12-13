package common

import (
	"time"
)

/*
	注册到KV系统中的VALUE
	1. ETCD -- value
	2. consul -- tags
*/
type RegistyNodeValue struct {
	Addr     string //IP:port
	Weight   int
	Metadata map[string]string
}

/*
	注册到KV系统中的VALUE
	1. ETCD -- 同类型的前缀(以ServiceName区分) /RootName/ServiceName/uniqid
	2. consul -- 相同的ServiceName,但要求NodeUniqID必须唯一(长度限制)
*/
type RegistyOptions struct {
	Endpoints   string //注册地址
	RootName    string //
	ServiceName string
	NodeUniqID  string
	NodeValue   RegistyNodeValue
	Interval    time.Duration
	Address     string
	Port        int
	UserName    string //`yaml:"user_name"`
	Pass        string //`yaml:"pass"`
}
