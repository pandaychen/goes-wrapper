package pymath

import (
	"math/rand"
	"sync"
	"time"
)

type PyProbability struct {
	// 生产随机数
	Seed *rand.Rand
	// rand.New(...) returns a non thread safe object
	SeedLock sync.Mutex
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewPyProbability() *PyProbability {
	return &PyProbability{
		Seed: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GET random N floats
func (p *PyProbability) GetRandNFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + p.Seed.Float64()*(max-min)
	}
	return res
}

func (p *PyProbability) ProbablyTrue(proba float64) bool {
	//proba must (0,1)
	p.SeedLock.Lock()
	defer p.SeedLock.Unlock()
	ok := p.Seed.Float64() < proba

	return ok
}
