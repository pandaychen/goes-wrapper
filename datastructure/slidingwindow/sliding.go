package slidingwindow

import (
	"sync"
	"time"

	"github.com/pandaychen/goes-wrapper/pytime"
)

const (
	DEFAULT_SLIDINGWINDOW_INTERVAL = 500 * time.Millisecond
)

// must be a ring
type SlidingWindow struct {
	WindowSize        int
	Window            *SWindow
	TickInterval      time.Duration //每个窗口（bucket）代表多长间隔
	TickDurationTotal time.Duration

	CurIndex    int
	LastTimeVal int64 // start time of the last swinbucket
	Lasttime    time.Duration
	Lock        sync.RWMutex
}

func NewSlidingWindow(size int, interval time.Duration) *SlidingWindow {
	if size < 1 {
		panic("WindowSize must greater than 0")
	}

	w := &SlidingWindow{
		WindowSize:   size,
		TickInterval: interval,
		Lasttime:     pytime.Duration2Now(), // 记录当前滑动窗口中最后一个位置的访问时间
	}

	//init window
	w.Window = NewSWindow(w.WindowSize)

	return w
}

//获取当前时间到上次slidingwindow被访问，这其中的差值跨越了几个window bucket
func (sw *SlidingWindow) get_sw_offset() int {
	span_bucket_size := int(pytime.Duration2Fixed(sw.Lasttime) / sw.TickInterval)
	if span_bucket_size >= 0 && span_bucket_size < sw.WindowSize {
		return span_bucket_size
	} else {
		// 超过了整个slidingwindow的最大长度
		return sw.WindowSize
	}
}

// 当有数据向滑动窗口写入时，更新本次写入的lasttime
func (sw *SlidingWindow) updata_sw_offset() error {
	span := sw.get_sw_offset()
	if span <= 0 {
		//不用更新
		return nil
	}

}

// 向slidingwindow的合适位置写入数据
func (sw *SlidingWindow) Add(val float64) {
	sw.Lock.Lock()
	defer sw.Lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}
