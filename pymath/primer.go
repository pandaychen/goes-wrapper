package pymath

//封装的素数算法

type PyPrime struct {
	Number uint64
}

func NewPyPrime(num uint64) *PyPrime {
	return &PyPrime{
		Number: num,
	}
}

func (p *PyPrime) NormalCheck() bool {
	if p.Number < 2 {
		return false
	}
	// sqrt(p.Number)
	for i := 2; uint64(i*i) <= p.Number; i++ {
		if p.Number%uint64(i) == 0 {
			return false
		}
	}
	return true
}

// TODO: millerranbin素性检测
//https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test
func (p *PyPrime) MillerRabinCheck() bool {
	return true
}
