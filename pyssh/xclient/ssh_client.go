package xclient

//基于multiplexing的SSH客户端封装

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type GwSSHClient struct {
	sync.Mutex
	*ssh.Client
	Config *SSHClientOptions

	//代理连接
	ProxyClient *SSHClient

	//复用session存储（注意：必须采用指针存储）
	traceSessionMap map[*ssh.Session]time.Time

	//复用个数
	refCount int32

	logger *zap.Logger
}

// 初始化网关用的SSH客户端
func NewGwSSHClient(opts ...SSHClientOption) (*GwSSHClient, error) {
	cfg := &SSHClientOptions{
		Host: "127.0.0.1",
		Port: "22", //default
	}
	for _, setter := range opts {
		setter(cfg)
	}
	return NewSSHClientWithCfg(cfg)
}

func getAvailableProxyClient(cfgs ...SSHClientOptions) (*GwSSHClient, error) {
	for i := range cfgs {
		if proxyClient, err := NewSSHClientWithCfg(&cfgs[i]); err == nil {
			return proxyClient, nil
		}
	}
	return nil, errors.New("bad proxy")
}

//初始化网关用的SSH客户端
func NewSSHClientWithCfg(cfg *SSHClientOptions) (*GwSSHClient, error) {
	sshCfg := ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            cfg.AuthMethods(),
		Timeout:         time.Duration(cfg.Timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config:          InitSSHConfig(),
	}
	destAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	if len(cfg.proxySSHClientOptions) > 0 {
		//TODO
	}

	gosshClient, err := ssh.Dial("tcp", destAddr, &sshCfg)
	if err != nil {
		//Log
		return nil, err
	}

	//初始化SSH客户端配置
	return &GwSSHClient{
		Client:          gosshClient,
		Config:          cfg,
		traceSessionMap: make(map[*ssh.Session]time.Time)}, nil
}

func (s *GwSSHClient) String() string {
	return fmt.Sprintf("%s@%s:%s", s.Config.Username, s.Config.Host, s.Config.Port)
}

// 获取连接引用计数
func (s *GwSSHClient) GetReferenceCount() int32 {
	return atomic.LoadInt32(&s.refCount)
}

func (s *GwSSHClient) Close() error {
	if s.ProxyClient != nil {
		s.ProxyClient.Close()
	}
	return s.Client.Close()
}

// 获取复用连接
func (s *GwSSHClient) AcquireSession() (*ssh.Session, error) {
	//fork session
	atomic.AddInt32(&s.refCount, 1)
	sshSession, err := s.Client.NewSession()
	if err != nil {
		//log
		atomic.AddInt32(&s.refCount, -1)
		return nil, err
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.traceSessionMap[sshSession] = time.Now()
	return sshSession, nil
}

func (s *GwSSHClient) ReleaseSession(sshSession *ssh.Session) {
	atomic.AddInt32(&s.refCount, -1)
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.traceSessionMap, sshSession)
}
