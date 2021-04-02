package pycpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pandaychen/goes-wrapper/system/pyfile"
	"github.com/pkg/errors"
)

const nanoSecondsPerSecond = 1e9

// ErrNoCFSLimit is no quota limit
var ErrNoCFSLimit = errors.Errorf("no quota limit")

// 暂时设定getClockTicks = 100
var clockTicksPerSecond = uint64(100)

//get cpu MHz         : 2399.988
func GetCpuFreq() uint64 {
	lines, err := pyfile.ReadAllLines("/proc/cpuinfo")
	if err != nil {
		return 0
	}
	for _, line := range lines {
		array := strings.Split(line, ":")
		if len(array) < 2 {
			continue
		}
		key := strings.TrimSpace(array[0])
		value := strings.TrimSpace(array[1])
		if key == "cpu MHz" {
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				return uint64(t * 1000.0 * 1000.0)
			}
		}
	}
	return 0
}

// GetSystemCpuUsage returns the host system's cpu usage in
// nanoseconds. An error is returned if the format of the underlying
// file does not match.
//
// Uses /proc/stat defined by POSIX. Looks for the cpu
// statistics line and then sums up the first seven fields
// provided. See man 5 proc for details on specific field
// information.

/*
[root@bf-serv-cli-test-56d6495d95-jfbm2 /]# cat /proc/stat
cpu  3467887829 7402144 670669262 65943628489 1051360 14218 253095025 0 0 0
cpu0 247812177 438929 51372570 4073750664 780116 1458 39979362 0 0 0
cpu1 235169218 434668 45947692 4094986373 10697 1639 26968916 0 0 0
cpu2 225781683 439578 43285824 4107724817 9592 1485 20862824 0 0 0
cpu3 220676697 513232 42289613 4113805579 9143 1655 18261785 0 0 0
cpu4 218076124 499334 41674599 4117533972 8491 1644 16610866 0 0 0
cpu5 215340998 487522 41179452 4120552851 8926 1738 15768888 0 0 0
cpu6 213499268 478080 40969862 4122913060 8723 1690 15081638 0 0 0
cpu7 210787169 471204 40527378 4124205430 150382 2903 15216436 0 0 0
cpu8 211982698 464440 40708345 4131111431 8645 1 10894860 0 0 0
cpu9 211332535 464665 40571501 4131870642 8526 0 10808146 0 0 0
cpu10 211359772 462076 40641212 4131862495 8403 0 10665528 0 0 0
cpu11 210828977 456705 40264459 4132867468 8112 0 10532617 0 0 0
cpu12 209268162 450942 40216814 4134700931 7863 0 10376140 0 0 0
cpu13 209149626 448272 40424134 4134634762 7823 1 10453089 0 0 0
cpu14 209033168 447126 40380819 4134762085 7891 0 10326323 0 0 0
cpu15 207789552 445363 40214984 4136345926 8021 0 10287602 0 0 0
intr 77186006610 24 114 0 0 229 0 0 0 0 0 0 0 0 0 43133294 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 69858620 0 34419764 0 329450373 161493 3882562 191001 3891312254 178012 4147213687 170135 3985072830 166723 581183279 163935 4031634070 163040 2484392448 161672 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 1475123891404
btime 1573183603
processes 3236443852
procs_running 2
procs_blocked 0
*/
func GetSystemCpuUsage() (uint64, error) {
	var (
		line  string
		f     *os.File
		usage uint64 = 0
		err   error
	)
	if f, err = os.Open("/proc/stat"); err != nil {
		return 0, err
	}
	bufReader := bufio.NewReaderSize(nil, 128)
	defer func() {
		bufReader.Reset(nil)
		f.Close()
	}()
	bufReader.Reset(f)
	for err == nil {
		if line, err = bufReader.ReadString('\n'); err != nil {
			err = errors.WithStack(err)
			return 0, err
		}
		array := strings.Fields(line)
		switch array[0] {
		case "cpu": //只统计cpu那行的数据
			if len(array) < 8 {
				err = errors.WithStack(fmt.Errorf("bad format of cpu stats"))
				return 0, err
			}
			var totalClockTicks uint64
			for _, i := range array[1:8] {
				var v uint64
				if v, err = strconv.ParseUint(i, 10, 64); err != nil {
					err = errors.WithStack(fmt.Errorf("error parsing cpu stats"))
					return 0, err
				}
				totalClockTicks += v
			}
			usage = (totalClockTicks * nanoSecondsPerSecond) / clockTicksPerSecond
			return usage, nil
		}
	}
	err = errors.Errorf("bad stats format")
	return 0, err
}
