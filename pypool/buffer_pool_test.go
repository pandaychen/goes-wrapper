package pypool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutBuffer(t *testing.T) {
	buf := make([]byte, 512)
	PutBuffer(buf)

	buf = make([]byte, 1024)
	PutBuffer(buf)

	buf = make([]byte, 2*1024)
	PutBuffer(buf)

	buf = make([]byte, 16*1024)
	PutBuffer(buf)
}

func TestGetBuffer(t *testing.T) {
	assert := assert.New(t)

	buf := GetBuffer(200)
	assert.Len(buf, 200)

	buf = GetBuffer(1025)
	assert.Len(buf, 1025)

	buf = GetBuffer(2 * 1024)
	assert.Len(buf, 2*1024)

	buf = GetBuffer(5 * 2000)
	assert.Len(buf, 5*2000)
}