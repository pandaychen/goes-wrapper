package enums

import (
	"errors"
)

var (
	ERROR_GRPC_DIAL                         = errors.New("grpc dial error")
	ERROR_GRPC_NOT_SUPPORT_SERVICE_REGISTRY = errors.New("not support grpc service registry")

	ERROR_GRPC_NO_TRACING_DEFINE = errors.New("no tracing defined")

	ERROR_GRPC_ACCTOKEN_NULL = errors.New("grpc meta accesstoken null")
	ERROR_GRPC_ACCTOKEN_ILLEGAL = errors.New("grpc meta accesstoken invalid")

	//credential
	ERROR_CREDENTIAL_PARAM_ILLEGAL = errors.New("grpc meta token set invalid")
)
