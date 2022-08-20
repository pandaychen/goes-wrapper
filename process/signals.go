package process

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pandaychen/goes-wrapper/pynet"
)

func forkNewProcessWithListener(ln net.Listener, bind_addr string) (*os.Process, error) {
	//构建热重启需要的metadata，加入子进程的环境变量
	var (
		pl                    ProcListener
		listenerEnvStr        []byte
		err                   error
		lnFile                *os.File
		execDirPath, execName string
		newProcess            *os.Process
	)
	if lnFile, err = pynet.TransListener2File(ln); err != nil {
		return nil, err
	}
	defer lnFile.Close()

	pl.Addr = bind_addr
	pl.FD = 3
	pl.Filename = lnFile.Name()

	if listenerEnvStr, err = json.Marshal(pl); err != nil {
		return nil, err
	}

	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
		lnFile,
	}

	// Get current environment and add in the listener to it
	environment := append(os.Environ(), fmt.Sprintf("LISTENER=%s", string(listenerEnvStr)))

	// Get current process name and directory
	if execName, err = os.Executable(); err != nil {
		return nil, err
	}
	execDirPath = filepath.Dir(execName)

	// Spawn child process
	if newProcess, err = os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDirPath,
		Env:   environment,
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	}); err != nil {
		return nil, err
	}

	return newProcess, nil
}

func WaitSignals(addr string, ln net.Listener, serverSock interface{}) error {
	signalCh := make(chan os.Signal, 1024)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case s := <-signalCh:
			switch s {
			case syscall.SIGUSR2:
				/*
					tcpListener, ok := serverSock.(*net.TCPListener)
					if !ok {
					}
					tcpListener.SetDeadline(time.Now())
				*/
				//Spawn child process
				_, err := forkNewProcessWithListener(ln, addr)
				if err != nil {
					continue
				}

				switch serverSock.(type) {
				case *http.Server:
					serverSock.(*http.Server).Close()
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					// Return any errors during shutdown.
					return serverSock.(*http.Server).Shutdown(ctx)
				}
			case syscall.SIGINT, syscall.SIGQUIT:

			}
		}
	}
}
