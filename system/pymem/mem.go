package pymem

import (
	"github.com/shirou/gopsutil/mem"
)

type Meminfo struct {
	Total       uint64
	Available   uint64
	UsedPercent float64
}

type NormalMem struct{}

//must return a pointer
func NewNormalMem() *NormalMem {
	return &NormalMem{}
}

func (m *NormalMem) GetMemoryInfo() (Meminfo, error) {
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
