package file_utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func MD5(text string) []byte {
	h := md5.New()
	io.WriteString(h, text)
	return h.Sum(nil)
}

func MD5Str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
