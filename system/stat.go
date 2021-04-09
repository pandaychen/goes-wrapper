package system

import (
	"errors"
	"fmt"
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
)

//create a memory wrapper
func NewPyMemorySystemMetrics(memtype string) (PyMemorySystemMetrics, error) {
	var memunit PyMemorySystemMetrics
	switch memtype {
	case "cpu":
		memunit = pycpu.NewNormalMem(mul, interval, logger)
	case "docker":
		return pycpu.NewDockerMem(mul, interval, logger)
	default:
		return nil, errors.New("not support type")
	}
}

//create a wrapper
func NewPyCpuSystemMetrics(cputype string, mul float64, interval time.Duration, logger *zap.Logger) (PyCpuSystemMetrics, error) {
	var cpunit PySystemMetrics
	switch cputype {
	case "cpu":
		cpunit = pycpu.NewCpuInfo(mul, interval, logger)
	case "docker":
		return pycpu.NewDockerCpuInfo(mul, interval, logger)
	default:
		return nil, errors.New("not support type")
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			//按照interval时间定时采集
			<-ticker.C
			usage, err := cpunit.GetCpuPercentage()
			if err != nil {
				logger.Error("NewPySystemMetrics-GetCpuPercentage", zap.Any("errmsg", err))
				continue
			}
			if usage == 0 {
				//logger.Error("NewPySystemMetrics-GetCpuPercentage usage invalid", zap.String("errmsg", "cpu usage zero"))
				continue
			}
			fmt.Println(cpunit.GetMemoryInfo())
			atomic.StoreUint64(&GlobalCpusage, usage)
		}
	}()

	return cpunit, nil
}

func GetCurrentCpuSystemMetrics() (uint64, uint64, error) {
	return atomic.LoadUint64(&GlobalCpusage), 0, nil
}

func main() {
	logger, _ := zap.NewProduction()
	NewPySystemMetrics("cpu", 10, DEFAULT_SYSTEM_COLLECT_INTERVAL, logger)
	for i := 0; i < 10; i++ {
		fmt.Println(GetCurrentCpuSystemMetrics())

		time.Sleep(1 * time.Second)
	}

	mem := NewPyMemorySystemMetrics()
	fmt.Println(mem.GetMemoryInfo())

	time.Sleep(100 * time.Second)
}
