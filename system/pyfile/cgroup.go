package pyfile

import (
	"strconv"
	"strings"
)

// ParseUintList parses and validates the specified string as the value
// found in some cgroup file (e.g. cpuset.cpus, cpuset.mems), which could be
// one of the formats below. Note that duplicates are actually allowed in the
// input string. It returns a map[int]bool with available elements from val
// set to true.
// Supported formats:
// 7
// 1-6
// 0,3-4,7,8-10
// 0-0,0,1-7
// 03,1-3 <- this is gonna get parsed as [1,2,3]
// 3,2,1
// 0-2,3,1

// 解析group中cpu的组织格式
func ParseCgroupCpuToUintList(val string) (map[int]bool, error) {
	if val == "" {
		return map[int]bool{}, nil
	}

	availableInts := make(map[int]bool)
	split := strings.Split(val, ",")
	errInvalidFormat := errors.Errorf("os/stat: invalid format: %s", val)
	for _, r := range split {
		if !strings.Contains(r, "-") {
			v, err := strconv.Atoi(r)
			if err != nil {
				return nil, errInvalidFormat
			}
			availableInts[v] = true
		} else {
			split := strings.SplitN(r, "-", 2)
			min, err := strconv.Atoi(split[0])
			if err != nil {
				return nil, errInvalidFormat
			}
			max, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, errInvalidFormat
			}
			if max < min {
				return nil, errInvalidFormat
			}
			for i := min; i <= max; i++ {
				availableInts[i] = true
			}
		}
	}
	return availableInts, nil
}
