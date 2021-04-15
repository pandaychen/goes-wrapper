package pymath

//获取最大公约数

func GetGCD(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
