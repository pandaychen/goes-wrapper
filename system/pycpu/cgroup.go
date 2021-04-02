package pycpu

//解析某个docker中运行进程的cgroup-system的数据

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pandaychen/goes-wrapper/system/pyfile"
)

var (
//CGroupRootDir string = "/sys/fs/cgroup"
)

/*
cgroup 子系统
[root@bf-serv-cli-test-56d6495d95-jfbm2 /]# ll /sys/fs/cgroup/
total 0
drwxr-xr-x 2 root root  0 Mar 29 07:59 blkio
lrwxrwxrwx 1 root root 11 Mar 29 07:59 cpu -> cpu,cpuacct
drwxr-xr-x 2 root root  0 Mar 29 07:59 cpu,cpuacct
lrwxrwxrwx 1 root root 11 Mar 29 07:59 cpuacct -> cpu,cpuacct
drwxr-xr-x 2 root root  0 Mar 29 07:59 cpuset
drwxr-xr-x 2 root root  0 Mar 29 07:59 devices
drwxr-xr-x 2 root root  0 Mar 29 07:59 freezer
drwxr-xr-x 2 root root  0 Mar 29 07:59 hugetlb
drwxr-xr-x 2 root root  0 Mar 29 07:59 memory
drwxr-xr-x 2 root root  0 Mar 29 07:59 net_cls
drwxr-xr-x 2 root root  0 Nov  8  2019 oom
drwxr-xr-x 2 root root  0 Mar 29 07:59 perf_event
drwxr-xr-x 2 root root  0 Mar 29 07:59 pids
drwxr-xr-x 2 root root  0 Mar 29 07:59 systemd
*/

// 存储当前进程的cgroup基础信息
type CgroupSystem struct {
	Name  string //docker name
	Pid   int    //当前的进程号
	IsPod bool   //是否为pod
	//CgroupSet map[string]interface{}
	CgroupSet map[string]string
}

//get current process's cgroup info,eg:
/*
	in a tke pod:
[root@bf-serv-cli-test-56d6495d95-jfbm2 cgroup]# cat /proc/6/cgroup
12:oom:/
11:net_cls:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
10:pids:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
9:hugetlb:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
8:cpuset:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
7:blkio:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
6:cpuacct,cpu:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
5:memory:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
4:freezer:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
3:devices:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
2:perf_event:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
1:name=systemd:/kubepods/burstable/podb60126d8-9064-11eb-836e-7a2eb63c648a/bde707e9e3043b4492f3a7d9be42b70d85c583c4792a3c3968194369163874e7
----------------------
	in a container:
[root@VM_0_7_centos system]# cat /proc/30368/cgroup
11:hugetlb:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
10:cpuset:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
9:freezer:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
8:devices:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
7:pids:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
6:cpuacct,cpu:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
5:memory:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
4:net_prio,net_cls:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
3:blkio:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
2:perf_event:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
1:name=systemd:/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
----------------------
	normal:
[root@VM_0_7_centos system]# cat /proc/4429/cgroup
11:hugetlb:/
10:cpuset:/
9:freezer:/
8:devices:/user.slice
7:pids:/user.slice
6:cpuacct,cpu:/user.slice
5:memory:/user.slice
4:net_prio,net_cls:/
3:blkio:/user.slice
2:perf_event:/
*/
func NewCgroupSystem(cur_pid int, ispod bool) (*CgroupSystem, error) {
	cg_fpath := fmt.Sprintf("/proc/%d/cgroup", cur_pid)

	//create a cgroup
	cur_cgroup := &CgroupSystem{
		Name:      cg_fpath,
		Pid:       cur_pid,
		IsPod:     ispod,
		CgroupSet: make(map[string]string),
	}

	fp, err := os.Open(cg_fpath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf := bufio.NewReader(fp)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		array := strings.Split(line, ":")
		if len(array) != 3 {
			return nil, fmt.Errorf("invalid cgroup file format %s", line)
		}
		dirpath := array[2]
		if dirpath != "/" && !ispod {
			//must in a docker,like /sys/fs/cgroup/cpu/system.slice/docker-1e939f40fb6a516434beab6ab3d2c91eb7f3100937a8b25d7576d0da8e05df2f.scope
			cur_cgroup.CgroupSet[array[1]] = path.Join(CGroupRootDir, array[1])
			if strings.Contains(array[1], ",") {
				for _, k := range strings.Split(array[1], ",") {
					cur_cgroup.CgroupSet[k] = path.Join(CGroupRootDir, k)
				}
			}
		} else {
			//like /sys/fs/cgroup/hugetlb/
			cur_cgroup.CgroupSet[array[1]] = path.Join(CGroupRootDir, array[1], array[2])
			if strings.Contains(array[1], ",") {
				for _, k := range strings.Split(array[1], ",") {
					cur_cgroup.CgroupSet[k] = path.Join(CGroupRootDir, k, array[2])
				}
			}
		}
	}
	return cur_cgroup, nil
}

// 获取CPU配额信息cpu.cfs_quota_us(可能为负数-1)
/*
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpu]# cat cpu.cfs_quota_us
200000
*/
func (c *CgroupSystem) GetCpuCfsQuotaUs() (int64, error) {
	if _, exists := c.CgroupSet["cpu"]; !exists {
		return -1, errors.New("param error")
	}
	data, err := pyfile.ReadFileContent(path.Join(c.CgroupSet["cpu"], "cpu.cfs_quota_us"))
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(data, 10, 64)
}

/* cpu.cfs_period_us
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpu]# cat cpu.cfs_period_us
100000
*/
func (c *CgroupSystem) GetCpuCfsPeriodUs() (uint64, error) {
	if _, exists := c.CgroupSet["cpu"]; !exists {
		return 0, errors.New("param error")
	}
	data, err := pyfile.ReadFileContent(path.Join(c.CgroupSet["cpu"], "cpu.cfs_period_us"))
	if err != nil {
		return 0, err
	}

	return pyfile.ParseFileContent2Uint64(data)
}

/*cpuacct.usage
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpuacct]# cat cpuacct.usage
71723686392
*/
func (c *CgroupSystem) GetCpuAcctUsage() (uint64, error) {
	if _, exists := c.CgroupSet["cpuacct"]; !exists {
		return 0, errors.New("param error")
	}
	data, err := pyfile.ReadFileContent(path.Join(c.CgroupSet["cpuacct"], "cpuacct.usage"))
	if err != nil {
		return 0, err
	}
	return pyfile.ParseFileContent2Uint64(data)
}

/*cpuacct.usage_percpu：每个cpu使用的时间
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpuacct]# cat cpuacct.usage_percpu
3856785637 9388740792 3249991024 9274368598 3909220940 4252675051 4999399867 3482721584 3265393334 2792999228 3690055360 3654456859 3307438423 4643994810 4450415726 3517266852
*/
func (c *CgroupSystem) GetCpuAcctUsagePerCPU() ([]uint64, error) {
	if _, exists := c.CgroupSet["cpuacct"]; !exists {
		return nil, errors.New("param error")
	}
	data, err := pyfile.ReadFileContent(path.Join(c.CgroupSet["cpuacct"], "cpuacct.usage_percpu"))
	if err != nil {
		return nil, err
	}
	var array []uint64
	for _, v := range strings.Fields(string(data)) {
		var u uint64
		if u, err = pyfile.ParseFileContent2Uint64(v); err != nil {
			return nil, err
		}
		if u != 0 {
			array = append(array, u)
		}
	}
	return array, nil
}

/*cpuset.cpus
[root@bf-serv-cli-test-56d6495d95-jfbm2 cpuset]# cat cpuset.cpus
0-15
*/
func (c *CgroupSystem) GetCpuSetCPUs() ([]uint64, error) {
	if _, exists := c.CgroupSet["cpuset"]; !exists {
		return nil, errors.New("param error")
	}
	data, err := pyfile.ReadFileContent(path.Join(c.CgroupSet["cpuset"], "cpuset.cpus"))
	if err != nil {
		return nil, err
	}
	cpus, err := pyfile.ParseCgroupCpuToUintList(data)
	if err != nil {
		return nil, err
	}
	var sets []uint64
	for k := range cpus {
		sets = append(sets, uint64(k))
	}
	return sets, nil
}
