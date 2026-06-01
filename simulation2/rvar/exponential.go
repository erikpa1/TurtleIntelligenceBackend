package rvar

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// expGen is an exponential distribution with the given mean, optionally
// truncated to [min, max] by rejection sampling.
type expGen struct {
	mean     float64
	min, max float64
	trunc    bool
}

func (e expGen) Sample(rng *rand.Rand) float64 {
	// rng.ExpFloat64 has rate 1 (mean 1); scaling by mean gives the target mean.
	for i := 0; i < 100; i++ {
		v := rng.ExpFloat64() * e.mean
		if !e.trunc || (v >= e.min && v <= e.max) {
			return v
		}
	}
	// Fall back to clamping if rejection failed to land in range.
	return math.Min(math.Max(rng.ExpFloat64()*e.mean, e.min), e.max)
}

func init() {
	// exp(mean)            — exponential with the given mean.
	// exp(mean, min)       — truncated below at min.
	// exp(mean, min, max)  — truncated to [min, max].
	//   exp(10s)   exp(10s, 2s)   exp(10s, 2s, 60s)
	ctor := func(args []float64) (Generator, error) {
		mean := args[0]
		if mean <= 0 {
			return nil, fmt.Errorf("exp: mean must be > 0, got %g", mean)
		}
		g := expGen{mean: mean, min: 0, max: math.Inf(1)}
		if len(args) >= 2 {
			g.trunc = true
			g.min = args[1]
		}
		if len(args) == 3 {
			g.max = args[2]
		}
		if g.max < g.min {
			return nil, fmt.Errorf("exp: max (%g) < min (%g)", g.max, g.min)
		}
		return g, nil
	}
	Register("exp", 1, 3, ctor)
	Register("exponential", 1, 3, ctor)
}
