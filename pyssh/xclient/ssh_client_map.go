package xclient

import "sync"

//
type GwClientsMap struct {
	sync.RWMutex
	GwClientMap map[string]*UserSSHClient
}

type UserSSHClient struct {
	ID   string // uniqid
	data map[*GwSSHClient]int64
}
