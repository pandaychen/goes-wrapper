package pymath

import "math/rand"

func PowerOfTwoChoices(rander *rand.Rand, length int) (int, int) {
	var (
		a int
		b int
	)
	a = rander.Intn(length)
	b = rander.Intn(length)

	if a == b {
		b = (b + 1) % length
	}

	return a, b
}
