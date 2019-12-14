package enums



const (
	SERVICE_REGISTRY_ETCDV3 = "etcdv3"
	SERVICE_REGISTRY_CONSUL = "consul"
	SERVICE_REGISTRY_ETCDV2 = "etcdv2"

	TRACING_NAME_JAEGER = "jaeger"
	TRACING_NAME_ZIPKIN = "zipkin"
	
	//etcd-kv-item
	GRPC_SERVICE_NAME_ROOT	 = "root"
	GRPC_SERVICE_NAME_TEST = "testsvc"

	REGISTRY_ETCDV3_DEFAULT_ADDRS = "http://127.0.0.1:2379"
	REGISTRY_CONSUL_DEFAULT_ADDRS = "http://127.0.0.1:8500"
	
	TRACING_JAEGER_DEFAULT_ADDR =  "127.0.0.1:6831"

	//credential type
	GRPC_CREDENTIAL_HTTP_AUTH_TYPE = "basic"
	GRPC_CREDENTIAL_APP_KEY_TYPE	  =  "appkey"

	GRPC_CREDENTIAL_DEFAULT_TOKEN = "abcedefgh"
)