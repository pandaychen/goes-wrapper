package common

import (
	opentracing "github.com/opentracing/opentracing-go"
	"time"
)

type GrpcClientConfig struct {
	Scheme string //resovler-scheme
	RootName         string //服务名1 -- 对应registry:RegistyNodeKey
	SvcName          string //服务名2 -- 对应registry:ServiceName
	MetaToken        string //客户端认证
	RegistyType      string //服务发现方法	consul/etcdv2/etcdv3
	RegistryAddrs    string //服务注册地址
	ResovleScheme    string //reslover-名字
	LoadBalancerName string //负载均衡方法
	TlsCommonName    string //
	TlsCertPath      string
	KeepAliveTtlOn   bool //keep-live on

	//OPENTracing配置
	Tracer     opentracing.Tracer //服务tracer
	TracerAddr string             //服务tracer-地址
	TracerType string             //tracer类型(zipkin|jaeger)

	//METAINFO认证
	AuthType          string //auth-basic(http标准认证)|APPID+APPSECKEY
	AccessToken       string //服务认证票据|AppSeckey
	AppID             string
	AuthRpcCredential WrapperPerRPCCredentials //notice this is a interface{}

	//客户端重试选项
	RetryOn			bool			//重试开关
	RetryInterval	time.Duration	//重试间隔
	RetryMax		int				//重试次数
	RetryTimeout	time.Duration	//重试超时时间

}

type GrpcServerConfig struct {
	RootName      string //服务名1
	SvcName       string //服务名
	BindAddr      string //服务地址
	Port          int
	NodeId        string //服务节点
	RegistryAddrs string //服务注册地址

	//OPENTracing配置
	Tracer     opentracing.Tracer //服务tracer
	TracerAddr string             //服务tracer-地址
	TracerType string             //tracer类型(zipkin|jaeger)

	//METAINFO认证
	AuthType    string //auth-basic(http标准认证)|APPID+APPSECKEY
	AccessToken string //服务认证票据
	AppID       string
}
