package system

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/pandaychen/goes-wrapper/system/pycpu"
	"github.com/pandaychen/goes-wrapper/system/pymem"
	"go.uber.org/zap"
)

const (
	DEFAULT_SYSTEM_COLLECT_INTERVAL time.Duration = 600 * time.Millisecond
)

//for cpu and docker-cpu
type PyCpuSystemMetrics interface {
	GetCpuPercentage() (uint64, error)
	GetCpuBasicInfo() pycpu.CpuBasic
}

//for mem and docker-mem
type PyMemorySystemMetrics interface {
	GetMemoryInfo() (pymem.Meminfo, error)
}

var (
	GlobalCpusage            uint64
	GlobalPyCpuSystemMetrics PyCpuSystemMetrics
	GlobalMeminfo            pymem.Meminfo
)

//create a memory wrapper
func NewPyMemorySystemMetrics(memtype string) (PyMemorySystemMetrics, error) {
	var (
		memunit PyMemorySystemMetrics
		err     error
	)
	switch memtype {
	case "cpu":
		memunit = pymem.NewNormalMem()
	case "docker":
		memunit = pymem.NewDockerMem()
	default:
		return nil, errors.New("not support type")
	}

	//for memory
	go func() {
		ticker := time.NewTicker(2 * DEFAULT_SYSTEM_COLLECT_INTERVAL)
		defer ticker.Stop()
		for {
			<-ticker.C
			GlobalMeminfo, err = memunit.GetMemoryInfo()
			if err != nil {
				logger.Error("NewPyMemorySystemMetrics-GetMemoryInfo", zap.Any("errmsg", err))
				continue
			}
		}
	}()

	return memunit, nil
}

//create a wrapper
func NewPyCpuSystemMetrics(cputype string, mul float64, interval time.Duration, logger *zap.Logger) (PyCpuSystemMetrics, error) {
	var cpunit PyCpuSystemMetrics
	switch cputype {
	case "cpu":
		cpunit = pycpu.NewCpuInfo(mul, interval, logger)
	case "docker":
		return pycpu.NewDockerCpuInfo(mul, interval, logger)
	default:
		return nil, errors.New("not support type")
	}

	//for cpu
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			//按照interval时间定时采集
			<-ticker.C
			usage, err := cpunit.GetCpuPercentage()
			if err != nil {
				logger.Error("NewPyCpuSystemMetrics-GetCpuPercentage", zap.Any("errmsg", err))
				continue
			}
			if usage == 0 {
				continue
			}
			atomic.StoreUint64(&GlobalCpusage, usage)
		}
	}()

	return cpunit, nil
}

func GetCurrentCpuSystemMetrics() (uint64, uint64, error) {
	return atomic.LoadUint64(&GlobalCpusage), 0, nil
}

func GetCurrentMemorySystemMetrics() *pymem.Meminfo {
	return &GlobalMeminfo
}

/*
func main() {
	logger, _ := zap.NewProduction()
	NewPyCpuSystemMetrics("cpu", 10, DEFAULT_SYSTEM_COLLECT_INTERVAL, logger)
	for i := 0; i < 10; i++ {
		fmt.Println(GetCurrentCpuSystemMetrics())

		time.Sleep(1 * time.Second)
	}

	mem, _ := NewPyMemorySystemMetrics("cpu")
	fmt.Println(mem.GetMemoryInfo())

	time.Sleep(100 * time.Second)
}
*/
