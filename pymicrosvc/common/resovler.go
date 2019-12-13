package common

type ResolverOptions struct {
	Endpoints   string //注册地址
	RootName    string
	ServiceName string
	//EtcdPrefixKey = /RootName/ServiceName
	Scheme   string
	UserName string //`yaml:"user_name"`
	Pass     string //`yaml:"pass"`
}
