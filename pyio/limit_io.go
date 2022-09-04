package pyio

import (
	"context"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/pandaychen/goes-wrapper/pymicrosvc/common"

	"github.com/pandaychen/goes-wrapper/pymicrosvc/ratelimit/tokenbucket"
)

// LimitIOReader read stream with RateLimiter.
type LimitIOReader struct {
	SrcReader io.Reader
	Limiter   common.RateLimiter
	Md5sum    hash.Hash //for MD5
}

// LimitIOWriter write stream with RateLimiter.
type LimitIOWriter struct {
	DestWriter io.Writer
	Limiter    common.RateLimiter
	Md5sum     hash.Hash //for MD5
}

// srcReader: reader
// rate: bytes/second
func NewLimitReaderWithLimiter(limiter common.RateLimiter, srcReader io.Reader, wantMd5 bool) *LimitIOReader {
	var (
		md5sum hash.Hash
	)
	if wantMd5 {
		md5sum = md5.New()
	}
	return &LimitIOReader{
		Limiter:   limiter,
		SrcReader: srcReader,
		Md5sum:    md5sum,
	}
}

func NewLimitWriterWithLimiter(limiter common.RateLimiter, destWriter io.Writer, wantMd5 bool) *LimitIOWriter {
	var (
		md5sum hash.Hash
	)
	if wantMd5 {
		md5sum = md5.New()
	}
	return &LimitIOWriter{
		Limiter:    limiter,
		DestWriter: destWriter,
		Md5sum:     md5sum,
	}
}

func newRateLimiterWithDefaultWindow(rate int) common.RateLimiter {
	return tokenbucket.NewRateLimiter(common.Trans2HumanRate(int64(rate)), int64(2))
}

func NewLimitReader(src io.Reader, rate int, calculateMd5 bool) *LimitIOReader {
	return NewLimitReaderWithLimiter(newRateLimiterWithDefaultWindow(rate), src, calculateMd5)
}

// 修改Reader实现，强制在读端限速
func (rl *LimitIOReader) Read(p []byte) (n int, err error) {
	n, e := rl.SrcReader.Read(p)
	if e != nil && e != io.EOF {
		return n, e
	}
	if n > 0 {
		if rl.Md5sum != nil {
			rl.Md5sum.Write(p[:n])
		}
		rl.Limiter.AcquireTokensWithBlock(context.TODO(), int64(n))
	}
	return n, e
}

// Write writes bytes from p.
func (wl *LimitIOWriter) Write(p []byte) (int, error) {
	if wl.Limiter == nil {
		return wl.DestWriter.Write(p)
	}
	n, err := wl.DestWriter.Write(p)
	if err != nil {
		//account error
		return n, err
	} else {
		if wl.Md5sum != nil {
			wl.Md5sum.Write(p[:n])
		}
		wl.Limiter.AcquireTokensWithBlock(context.TODO(), int64(n))
		return n, nil
	}
}

func copyFileWithLimit(sourceFile, destFile string) {
	file, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dstFile, err := os.Create("destFile")
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 2*1024)
	TotalLimit := tokenbucket.NewRateLimiter((1024 * 512), 2)
	//limitReader := NewLimitReaderWithLimiter(TotalLimit, file,false)
	limitReader := NewLimitReaderWithLimiter(TotalLimit, file, true)
	_, copyErr := io.CopyBuffer(dstFile, limitReader, buf)
	if copyErr != nil {
		fmt.Println(copyErr.Error())
		return
	} else {
		fmt.Printf("new md5=%x", limitReader.Md5sum.Sum(nil))
		return
	}
}

func copyFileWithWriterLimit(sourceFile, destFile string) {
	file, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dstFile, err := os.Create("destFile")
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 2*1024)
	TotalLimit := tokenbucket.NewRateLimiter((1024 * 512), 2)
	limitWriter := NewLimitWriterWithLimiter(TotalLimit, dstFile, true)
	_, copyErr := io.CopyBuffer(limitWriter, file, buf)
	if copyErr != nil {
		fmt.Println(copyErr.Error())
		return
	} else {
		fmt.Printf("new md5=%x", limitWriter.Md5sum.Sum(nil))
		return
	}
}

func main() {
	//copyFileWithLimit("./test_file", "test_file.bak")
	copyFileWithWriterLimit("./test_file", "test_file.bak")
}
