package common

/*
	客户端(no care struct)需要实现如下interface{}
	type PerRPCCredentials interface {
    	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
    	RequireTransportSecurity() bool
	}
在 gRPC 中默认定义了 PerRPCCredentials，是 gRPC 默认提供用于自定义认证的接口，作用是将所需的安全认证信息添加到每个 RPC 方法的上下文中。其包含 2 个方法：
GetRequestMetadata：获取当前请求认证所需的元数据（metadata）
RequireTransportSecurity：是否需要基于 TLS 认证进行安全传输
*/

import (
	"../enums"
	"golang.org/x/net/context"
	"fmt"
)

//this struct will be inject into metainfo
type RPCCredential struct {
	ServiceToken   string
	AppId          string
	CredentialType string
}

type WrapperPerRPCCredentials interface {
	GetRequestMetadata(context.Context, ...string) (map[string]string, error)
	RequireTransportSecurity() bool
	SetServiceToken(string)
	SetServiceTokenType(string)
}

//interface{} 和 struct{} 通信
func (c *RPCCredential) SetServiceToken(token string) {
	c.ServiceToken = token
}

func (c *RPCCredential) SetServiceTokenType(tokentype string) {
        c.CredentialType = tokentype
}

func (c *RPCCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if c.CredentialType == enums.GRPC_CREDENTIAL_HTTP_AUTH_TYPE {
		fmt.Println(c.ServiceToken)
		return map[string]string{
			"authorization": "basic " + c.ServiceToken,
		}, nil
	} else {
		return map[string]string{
			"AppId":  c.AppId,
			"Appkey": c.ServiceToken,
		}, nil
	}
}

func (c *RPCCredential) RequireTransportSecurity() bool {
	//return true	//if use tls
	return false
}
