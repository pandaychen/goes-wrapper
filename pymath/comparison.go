package pymath

func Maxer(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Miner(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
