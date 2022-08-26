package xclient

//基于唯一ID的ssh客户端管理封装

import "time"

type UserSSHClient struct {
	UniqId       string                 // uniqid
	sshClientMap map[*GwSSHClient]int64 //value：session计数

	username string
}

type SShClient2Add struct {
	UniqId string
	*GwSSHClient
}

func NewUserSSHClient(uniq_id string) *UserSSHClient {
	if uniq_id == "" {
		return nil
	}

	return &UserSSHClient{
		UniqId:       uniq_id,
		sshClientMap: make(map[*GwSSHClient]int64),
	}
}

//加入ssh-client
func (c *UserSSHClient) AddSshClient(client *GwSSHClient) {
	c.sshClientMap[client] = time.Now().UnixNano()
}

func (c *UserSSHClient) Count() int {
	return len(c.sshClientMap)
}

func (c *UserSSHClient) RemoveSshClient(client *GwSSHClient) {
	if _, exists := c.sshClientMap[client]; exists {
		delete(c.sshClientMap, client)
		if client != nil {
			client.Close()
		}
	}
}

// 获取可用的client
func (c *UserSSHClient) GetSshClient() (*GwSSHClient, error) {
	var (
		client   *GwSSHClient
		refCount int32
	)

	// 取引用最少的 SSHClient
	for clientItem := range c.sshClientMap {
		if refCount <= clientItem.GetReferenceCount() {
			refCount = clientItem.GetReferenceCount()
			client = clientItem
		}
	}
	return client, nil
}

// 回收引用数为0的clients
func (c *UserSSHClient) Recycle() {
	//
	needRemovedClients := make([]*GwSSHClient, 0, c.Count())
	for cli := range c.sshClientMap {
		if cli.RefCount() <= 0 {
			//回收
			needRemovedClients = append(needRemovedClients, cli)
			if cli != nil {
				cli.Close()
			}
		}
	}
	if len(needRemovedClients) > 0 {
		for i := range needRemovedClients {
			delete(c.sshClientMap, needRemovedClients[i])
		}
	}
}
