package xclient

import "sync"

//
type SSHClientsManager struct {
	sync.RWMutex
	gwClientMap map[string]*UserSSHClient

	//lockless
	addChan    chan *SShClient2Add
	reqChan    chan string // reqId
	resultChan chan *SSHClient
	searchChan chan string // prefix

}

func NewSSHClientsManager() *SSHClientsManager {
	m := SSHClientsManager{
		storeChan:  make(chan *storeClient),
		reqChan:    make(chan string),
		resultChan: make(chan *SSHClient),
		searchChan: make(chan string),
	}
	go m.run()
	return &m
}

func (m *SSHClientsManager) GetClientsMap() map[string]*UserSSHClient {
	return m.gwClientMap
}
