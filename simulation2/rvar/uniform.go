package rvar

import (
	"fmt"
	"math/rand/v2"
)

// uniformGen draws uniformly from the closed-open interval [min, max).
type uniformGen struct{ min, span float64 }

func (u uniformGen) Sample(rng *rand.Rand) float64 {
	return u.min + rng.Float64()*u.span
}

func init() {
	// uniform(min, max) — continuous uniform between two times.
	//   uniform(5s, 30s)
	Register("uniform", 2, 2, func(args []float64) (Generator, error) {
		min, max := args[0], args[1]
		if max < min {
			return nil, fmt.Errorf("uniform: max (%g) < min (%g)", max, min)
		}
		return uniformGen{min: min, span: max - min}, nil
	})
}
