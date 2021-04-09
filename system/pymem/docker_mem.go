package pymem

import "github.com/shirou/gopsutil/mem"

type DockerMem struct{}

func NewDockerMem() *DockerMem {
	return &DockerMem{}
}

func (m *DockerMem) GetMemoryInfo() (Meminfo, error) {
	minfo, err := mem.VirtualMemory()
	if err != nil {
		return Meminfo{}, err
	}

	return Meminfo{
		Total:       minfo.Total,
		Available:   minfo.Available,
		UsedPercent: minfo.UsedPercent,
	}, nil
}
