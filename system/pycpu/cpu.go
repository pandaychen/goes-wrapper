package pycpu

//非容器场景采集，由于psutils的cpu包是取interval进行采集，这里选择启动一个goroutine定时采集，外部程序只需要取当前的CPU数据即可

import (
	"errors"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
)

type CpuInfo struct {
	Interval  time.Duration //采集周期
	Cores     uint64        //核心数
	Cpu       int
	Frequency uint64 //主频
	Logger    *zap.Logger
}

func newCpuInfo(inter time.Duration, logger *zap.Logger) *CpuInfo {
	c := &CpuInfo{
		Interval: inter,
		Logger:   logger}
	return c
}

func (c *CpuInfo) GetCpuPercentage(multiple float64) (uint64, error) {
	var percents []float64
	var err error
	var usage uint64
	if multiple < 0 {
		return 0, errors.New("multiple illegal")
	}
	percents, err = cpu.Percent(c.Interval, false) //[6.0000000055879354]
	if err == nil {
		usage = uint64(percents[0] * multiple) // 扩大multiple倍
	}
	return usage, err
}

func (c *CpuInfo) UpdateCpuBasicInfo() {
	stats, err := cpu.Info()
	if err != nil {
		c.Logger.Error("UpdateCpuBasicInfo Info error", zap.Any("errmsg", err))
		return
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		c.Logger.Error("UpdateCpuBasicInfo Counts error", zap.Any("errmsg", err))
		return
	}
	c.Frequency = uint64(stats[0].Mhz)
	c.Cores = uint64(cores)

	return
}
