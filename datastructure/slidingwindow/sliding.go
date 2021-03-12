package slidingwindow

import (
	"sync"
	"time"
)

// must be a ring
type SlidingWindow struct {
	Lock              sync.RWMutex
	WindowSize        int
	Window            *SWindow
	TickInterval      time.Duration
	TickDurationTotal time.Duration

	CurIndex int
	LastTime time.Duration // start time of the last swinbucket
}

func NewSlidingWindow(size int, interval time.Duration) *SlidingWindow {
	if size < 1 {
		panic("WindowSize must greater than 0")
	}

	w := &SlidingWindow{
		WindowSize:   size,
		TickInterval: interval,
		LastTime:     time.Now(),
	}

	return w
}
