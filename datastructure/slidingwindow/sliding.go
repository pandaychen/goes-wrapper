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

	CurIndex       int   //
	LastTimeVal    int64 // start time of the last swinbucket
	Lasttime       time.Duration
	Lock           sync.RWMutex
	ReduceCallback func(b *SWinBucket) float64
}

func NewSlidingWindow(size int, interval time.Duration, cb func(b *SWinBucket) float64) *SlidingWindow {
	if size < 1 {
		panic("WindowSize must greater than 0")
	}

	w := &SlidingWindow{
		CurIndex:     0,
		WindowSize:   size,
		TickInterval: interval,
		Lasttime:     pytime.Duration2Now(), // 记录当前滑动窗口中最后一个位置的访问时间
	}

	//init window
	w.Window = NewSWindow(w.WindowSize)
	w.ReduceCallback = cb
	return w
}

//获取当前时间到上次slidingwindow被访问，这其中的差值跨越了几个window bucket
func (sw *SlidingWindow) get_sw_offset() int {
	span_bucket_size := int(pytime.Duration2Fixed(sw.Lasttime) / sw.TickInterval)
	if span_bucket_size >= 0 && span_bucket_size < sw.WindowSize {
		return span_bucket_size
	} else {
		// 超过了整个slidingwindow的最大长度，包含了span_bucket_size<0的情况
		return sw.WindowSize
	}
}

// 当有数据向滑动窗口写入时，更新本次写入的lasttime
func (sw *SlidingWindow) updata_sw_offset() error {
	//取当前时间离最后（滑动窗口）一次更新之间经过了多少时间跨度
	span := sw.get_sw_offset()
	if span <= 0 {
		//不用更新
		return nil
	}

	var curindex int = sw.CurIndex

	var i int
	//无论如何，lasttime到lasttine+span这个区间的节点只存在两种情况：1.旧数据 2.空节点
	for i = 0; i < span; i++ {
		// 依次重置过期的bucket
		sw.Window.ResetFixedBucket((curindex + i + 1) % sw.WindowSize)
	}
	//更新本地要写入的index及lasttime，注意后者要按照interval对齐
	sw.CurIndex = (curindex + span) % sw.WindowSize

	curRtime := pytime.Duration2Now()

	//调整curtime和sw.Lasttime的差值能够是sw.TickInterval的倍数关系
	sw.Lasttime = curRtime - (curRtime-sw.Lasttime)%sw.TickInterval
	return nil
}

// 向slidingwindow的合适位置写入数据
func (sw *SlidingWindow) Add(val float64) {
	sw.Lock.Lock()
	defer sw.Lock.Unlock()

	//更新游标
	sw.updata_sw_offset()

	//想curIndex位置写入数据（多个？累计的情况？）
	sw.Window.Add(sw.CurIndex, val)
}

// 遍历指定范围内滑动窗口，得到累计的值
func (sw *SlidingWindow) Reduce() float64 {
	sw.Lock.RLock()
	defer sw.Lock.RUnlock()

	//srange为需要range的节点个数
	var srange int

	offset := sw.get_sw_offset()

	if offset == 0 {
		//落在当前的视角窗口也计入统计范围
		srange = sw.WindowSize
	} else {
		srange = sw.WindowSize - offset //这个offset的部分是不需要的（过期或者为空）
	}

	if srange > 0 {
		//2种位置（前或者后），计算出首地址（即当前时间戳对应位置的地址，从这里开始遍历，这里加1是当前指向的节点以及被reset了，跳过当前这个节点，从下个节点开始）
		startindex := (sw.CurIndex + offset + 1) % sw.WindowSize
		//使用global在回调方法中获取结构
		res := sw.Window.Reduce(startindex, srange, sw.ReduceCallback)
		return res
	}
	return -1
}
