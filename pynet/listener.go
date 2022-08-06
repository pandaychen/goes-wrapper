package pynet

import (
	"fmt"
	"net"
	"os"
)

func TransListener2File(ln net.Listener) (*os.File, error) {
	switch t := ln.(type) {
	case *net.TCPListener:
		//TCP listener
		return t.File()
	case *net.UnixListener:
		//UDP listener
		return t.File()
	}
	return nil, fmt.Errorf("unsupported listener:%T", ln)
}
