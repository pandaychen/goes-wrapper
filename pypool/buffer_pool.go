package pypool

//分配[]byte的多级字节池

import "sync"

var (
	bufPool32k     sync.Pool
	bufPool16k     sync.Pool
	bufPool8k      sync.Pool
	bufPool2k      sync.Pool
	bufPool1k      sync.Pool
	bufPoolDefault sync.Pool
)

func GetBuffer(size int) []byte {
	var (
		buf interface{}
	)
	if size >= 32*1024 {
		buf = bufPool32k.Get()
	} else if size >= 16*1024 {
		buf = bufPool16k.Get()
	} else if size >= 8*1024 {
		buf = bufPool8k.Get()
	} else if size >= 2*1024 {
		buf = bufPool2k.Get()
	} else if size >= 1*1024 {
		buf = bufPool1k.Get()
	} else {
		buf = bufPoolDefault.Get()
	}
	if buf == nil {
		return make([]byte, size)
	}
	bufferer := buf.([]byte)
	if cap(bufferer) < size {
		//WARING：从sync.Pool中取出的不一定能满足要求
		return make([]byte, size)
	}
	return bufferer[:size]
}

func PutBuffer(buffer []byte) {
	var (
		size int = cap(buffer)
	)

	//Need to reset buffer?

	if size >= 32*1024 {
		bufPool32k.Put(buffer)
	} else if size >= 16*1024 {
		bufPool16k.Put(buffer)
	} else if size >= 8*1024 {
		bufPool8k.Put(buffer)
	} else if size >= 2*1024 {
		bufPool2k.Put(buffer)
	} else if size >= 1*1024 {
		bufPool1k.Put(buffer)
	} else {
		bufPoolDefault.Put(buffer)
	}
}
