package wheel_timer

import (
	"time"

	"github.com/pandaychen/goes-wrapper/pytime"
)

type TimingwheelConfig struct {
	SlotNum  int
	Interval time.Duration
}

func DefaultTimingwheelConfig() *TimingwheelConfig {
	return &TimingwheelConfig{
		SlotNum: 3600,
		//默认ticker周期为1s
		Interval: pytime.Duration("1s"),
	}
}
