package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type ProcListener struct {
	Addr     string `json:"addr"`
	FD       int    `json:"fd"`
	Filename string `json:"filename"`
}

func NewProcListener(addr, filename string, fd int) *ProcListener {
	return &ProcListener{
		Addr:     addr,
		FD:       fd,
		Filename: filename,
	}
}

func (l *ProcListener) CreateTcpListener() (net.Listener, error) {
	ln, err := net.Listen("tcp", l.Addr)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

// Extract the encoded listener metadata from the environment
func (l *ProcListener) ImportEnvListener() (net.Listener, error) {
	listenerEnvStr := os.Getenv("LISTENER")
	if listenerEnvStr == "" {
		return nil, fmt.Errorf("unable to find LISTENER environment vars")
	}

	// Unmarshal the listener metadata.
	var (
		envListener ProcListener
		err         error
	)

	//create listener from env
	err = json.Unmarshal([]byte(listenerEnvStr), &envListener)
	if err != nil {
		return nil, err
	}
	if envListener.Addr /*from env*/ != l.Addr /*self or param*/ {
		return nil, fmt.Errorf("listener ip not equal env=%s,new=%s", envListener.Addr, l.Addr)
	}

	//更新成员
	l.Addr = envListener.Addr
	l.FD = envListener.FD
	l.Filename = envListener.Filename

	// The file has already been passed to this process, extract the file
	// descriptor and name from the metadata to rebuild/find the *os.File for
	// the listener
	listenerFile := os.NewFile(uintptr(l.FD), l.Filename)
	if listenerFile == nil {
		return nil, fmt.Errorf("unable to create listener file: %v", err)
	}
	defer listenerFile.Close()

	// Create a net.Listener from the *os.File
	ln, err := net.FileListener(listenerFile)
	if err != nil {
		fmt.Println(err, l.Filename, listenerFile)
		return nil, err
	}

	return ln, nil
}

func (l *ProcListener) RecreateListener(addr string) (net.Listener, error) {
	if addr != "" {
		l.Addr = addr
	}
	ln, err := l.ImportEnvListener()
	if err == nil {
		//retrieve from env succ
		return ln, nil
	}

	if addr != "" && l.Addr == "" {
		l.Addr = addr
	}
	//create a new listener
	ln, err = l.CreateTcpListener()
	if err != nil {
		return nil, err
	}

	return ln, nil
}
