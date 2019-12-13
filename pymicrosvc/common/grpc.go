package common

type GrpcClientConfig struct {
	Scheme           string //resovler-scheme
	RootName         string //服务名1
	SvcName          string //服务名2
	MetaToken        string //客户端认证
	RegistyType      string //服务发现方法	consul/etcdv2/etcdv3
	RegistryAddrs    string //服务注册地址
	ResovleScheme    string //reslover-名字
	LoadBalancerName string //负载均衡方法
	TlsCommonName    string //
	TlsCertPath      string
	KeepAliveTtlOn   bool //keep-live on
}

type GrpcServerConfig struct {
	RootName     string //服务名1
	SvcName      string //服务名
	BindAddr     string //服务地址
	Port         int
	NodeId       string //服务节点
	RegistryAddr string //服务注册地址
}
