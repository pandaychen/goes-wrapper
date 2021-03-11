package slidingwindow

// SWindow is a collection of SWinBucket

type SWindow struct {
	buckets    []*SWinBucket
	windowsize int
}

func NewSWindow(windowsize int) *SWindow {
	buckets := make([]*Bucket, size)
	for i := 0; i < windowsize; i++ {
		buckets[i] = new(SWinBucket)
	}

	return &SWindow{
		buckets:    buckets,
		windowsize: windowsize,
	}
}

func (w *SWindow) Add(index int, value float64) {
	//add value to index's bucket
	if w.buckets[index%w.windowsize] == nil {
		return
	}
	w.buckets[index%w.windowsize].add(value)
}

func (w *SWindow) ResetFixedBucket(index int) {
	if w.buckets[index%w.windowsize] == nil {
		return
	}
	w.buckets[index%w.windowsize].reset()
}

func (w *SWindow) Reset(index int) {
	for _, v := range buckets {
		if v == nil {
			continue
		}
		v.reset()
	}
}

// 统计滑动窗口的数据 start:起始位置 rcount:遍历长度 cb: 处理每个bucket的回调方法
func (w *SWindow) Reduce(start, rcount int, cb func(b *SWinBucket)) {
	for i := 0; i < rcount; i++ {
		cb(w.buckets[(start+i)%w.windowsize])
	}
}
