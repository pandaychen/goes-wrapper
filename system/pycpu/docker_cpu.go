package pycpu

//用于container or k8s pod的cpu信息采集

import (
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
)

/*
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpu]# pwd
/sys/fs/cgroup/cpu
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpu]# ls
cgroup.clone_children  cgroup.procs       cpu.cfs_quota_us          cpu.rt_period_us   cpu.shares  cpuacct.stat    cpuacct.usage         notify_on_release
cgroup.event_control   cpu.cfs_period_us  cpu.cfs_relax_thresh_sec  cpu.rt_runtime_us  cpu.stat    cpuacct.uptime  cpuacct.usage_percpu  tasks
*/

var (
	CGroupRootDir = "/sys/fs/cgroup"
)

type DockerCpuInfo struct {
	Interval time.Duration //采集周期
	Logger   *zap.Logger
	Pid      int //当前进程号
	IsPod    bool

	//sysinfo
	Frequency uint64 //主频
	Cores     uint64 //核心数
	Quota     float64
	CpuUsage  uint64 //CPU使用率

	//保存上一次的值，用于校正
	OldSysCpuUsage   uint64
	OldTotalCpuUsage uint64

	CgroupCpu *CgroupSystem
}

// 初始化container的CPU采集结构
func NewDockerCpuInfo(interval time.Duration, is_pod bool, logger *zap.Logger) (*DockerCpuInfo, error) {
	var (
		//quota    float64
		core_num int
		//frequency uint64
	)

	dc := &DockerCpuInfo{
		Interval: interval,
		Logger:   logger,
		Pid:      os.Getpid(),
		IsPod:    is_pod,
	}
	//dc.Pid = 30368
	dc_cgroup, err := NewCgroupSystem(dc.Pid, dc.IsPod)
	if err != nil {
		dc.Logger.Error("NewDockerCpuInfo-NewCgroupSystem-error", zap.Any("errmsg", err))
		return nil, err
	}
	dc.CgroupCpu = dc_cgroup

	//set core num
	core_num, err = cpu.Counts(true)
	if err != nil || core_num == 0 {
		var cpus []uint64
		cpus, err = dc_cgroup.GetCpuAcctUsagePerCPU()
		if err != nil {
			dc.Logger.Error("NewDockerCpuInfo-GetCpuAcctUsagePerCPU error", zap.Any("errmsg", err))
			return nil, err
		}
		core_num = len(cpus)
	}
	dc.Cores = uint64(core_num)

	//set quota num
	sets, err := dc_cgroup.GetCpuSetCPUs()
	if err != nil {
		dc.Logger.Error("NewDockerCpuInfo-GetCpuSetCPUs error", zap.Any("errmsg", err))
		return nil, err
	}
	dc.Quota = float64(len(sets))

	cgroup_quota, err := dc_cgroup.GetCpuCfsQuotaUs()
	if err == nil && cgroup_quota != -1 {
		// 容器限制的cpu值:= quota/peroid
		var period uint64
		period, err = dc_cgroup.GetCpuCfsPeriodUs()
		if err != nil {
			dc.Logger.Error("NewDockerCpuInfo-GetCpuCfsPeriodUs error", zap.Any("errmsg", err))
			return nil, err
		}

		//计算limiter配额，并和quota比较，取最小
		real_limiter := float64(cgroup_quota) / float64(period)
		if real_limiter < dc.Quota {
			dc.Quota = real_limiter
		}
	}

	dc.Frequency = GetCpuFreq()
	fmt.Println(dc.Frequency)
	dc.OldTotalCpuUsage, err = dc_cgroup.GetCpuAcctUsage()
	if err != nil {
		dc.Logger.Error("NewDockerCpuInfo-GetCpuAcctUsage error", zap.Any("errmsg", err))
		return nil, err
	}

	dc.OldSysCpuUsage, err = GetSystemCpuUsage()
	if err != nil {
		dc.Logger.Error("NewDockerCpuInfo-GetSystemCpuUsage error", zap.Any("errmsg", err))
		return nil, err
	}

	return dc, nil
}

//
func (c *DockerCpuInfo) GetCpuPercentage() (uint64, error) {
	var (
		cur_total  uint64
		cur_system uint64
		ret_usage  uint64
		err        error
	)
	cur_total, err = c.CgroupCpu.GetCpuAcctUsage()
	//fmt.Println(cur_total, c.OldTotalCpuUsage)
	if err != nil {
		c.Logger.Error("GetCpuPercentage-GetCpuAcctUsage error", zap.Any("errmsg", err))
		return 0, nil
	}
	cur_system, err = GetSystemCpuUsage()
	//fmt.Println(cur_system, c.OldSysCpuUsage)
	if err != nil {
		c.Logger.Error("GetCpuPercentage-GetSystemCpuUsage error", zap.Any("errmsg", err))
		return 0, nil
	}
	//校正
	if cur_system != c.OldSysCpuUsage {
		ret_usage = uint64(float64((cur_total-c.OldTotalCpuUsage)*c.Cores*1e3) / (float64(cur_system-c.OldSysCpuUsage) * c.Quota))
	}

	//update
	c.OldSysCpuUsage = cur_system
	c.OldTotalCpuUsage = cur_total
	return ret_usage, nil
}

/*
func main() {
	logger, _ := zap.NewProduction()
	c, _ := NewDockerCpuInfo(time.Second, false, logger)
	fmt.Println(c, c.CgroupCpu, c.Quota, c.Cores)
	fmt.Println(c.GetCpuPercentage())
}
*/
