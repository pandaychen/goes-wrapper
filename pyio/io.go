package pyio

import (
	"io"
	"sync"

	"github.com/pandaychen/goes-wrapper/pypool"
)

//io.ReadWriteCloser: like net.Conn：https://pkg.go.dev/net#Conn

var (
	_DEFAULT_BUFFER_SIZE = 32 * 1024 //32K
)

// link two io.ReadWriteCloser with sync.Pool
func BridgeIoConns(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) (inCount int64, outCount int64) {
	var (
		wg sync.WaitGroup
	)

	pipeFunc := func(toStream io.ReadWriteCloser, fromStream io.ReadWriteCloser, count *int64) {
		defer func() {
			toStream.Close()
			fromStream.Close()
			wg.Done()
		}()

		buf := pypool.GetBuffer(_DEFAULT_BUFFER_SIZE)
		defer pypool.PutBuffer(buf)

		//https://pkg.go.dev/io#CopyBuffer
		*count, _ = io.CopyBuffer(toStream, fromStream, buf)
	}

	wg.Add(2)
	go pipeFunc(c1, c2, &inCount)
	go pipeFunc(c2, c1, &outCount)
	wg.Wait()
	return
}
