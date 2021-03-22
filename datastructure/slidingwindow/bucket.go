package slidingwindow

//sliding window basic unit

type SWinBucket struct {
	Sum   float64 //可扩展，存储多维度数据
	Count int64
}

func (b *SWinBucket) add(v float64) {
	b.Sum += v
	b.Count++ // 只是累加窗口中的值（rollingcounter）
}

func (b *SWinBucket) reset() {
	b.Sum = 0
	b.Count = 0
}
