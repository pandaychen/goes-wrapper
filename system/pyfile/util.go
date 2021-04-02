package pyfile

//quote from kratos

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// 获取文件所有内容
func ReadFileContent(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "os/stat: readfile(%s) error", path)
	}
	return strings.TrimSpace(string(contents)), nil
}

//解析文件内容到uint64类型
func ParseFileContent2Uint64(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, 10, 64)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, errors.Wrapf(err, "os/stat: parseUint(%s) error", s)
	}
	return v, nil
}

// ReadLines reads contents from a file and splits them by new lines.
// A convenience wrapper to ReadLinesOffsetN(filename, 0, -1).
func ReadAllLines(filename string) ([]string, error) {
	return ReadLinesWithOffsetN(filename, 0, -1)
}

// ReadLinesOffsetN reads contents from file and splits them by new line.
// The offset tells at which line number to start.
// The count determines the number of lines to read (starting from offset):
//   n >= 0: at most n lines
//   n < 0: whole file
func ReadLinesWithOffsetN(filename string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n+int(offset) || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < int(offset) {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}
