package rvar

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// normalGen is a Gaussian distribution, optionally truncated to [min, max].
type normalGen struct {
	mean, sd float64
	min, max float64
	trunc    bool
}

func (n normalGen) Sample(rng *rand.Rand) float64 {
	for i := 0; i < 100; i++ {
		v := rng.NormFloat64()*n.sd + n.mean
		if !n.trunc || (v >= n.min && v <= n.max) {
			return v
		}
	}
	return math.Min(math.Max(rng.NormFloat64()*n.sd+n.mean, n.min), n.max)
}

func init() {
	// normal(mean, sd)               — Gaussian.
	// normal(mean, sd, min, max)     — truncated to [min, max].
	//   normal(60s, 5s)   normal(60s, 5s, 30s, 90s)
	ctor := func(args []float64) (Generator, error) {
		mean, sd := args[0], args[1]
		if sd < 0 {
			return nil, fmt.Errorf("normal: sd must be >= 0, got %g", sd)
		}
		g := normalGen{mean: mean, sd: sd, min: math.Inf(-1), max: math.Inf(1)}
		if len(args) == 4 {
			g.trunc = true
			g.min, g.max = args[2], args[3]
			if g.max < g.min {
				return nil, fmt.Errorf("normal: max (%g) < min (%g)", g.max, g.min)
			}
		}
		return g, nil
	}
	Register("normal", 2, 4, ctor)
	Register("gauss", 2, 4, ctor)
}
