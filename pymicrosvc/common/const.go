package common

import "time"

var (
	// grpc options
	GRPC_KeepAliveTime         = time.Duration(10) * time.Second
	GRPC_KeepAliveTimeout      = time.Duration(3) * time.Second
	GRPC_BackoffMaxDelay       = time.Duration(3) * time.Second
	GRPC_MaxSendMsgSize        = 1 << 24
	GRPC_MaxCallMsgSize        = 1 << 24
	GRPC_InitialWindowSize     = 1 << 30
	GRPC_InitialConnWindowSize = 1 << 30
)
