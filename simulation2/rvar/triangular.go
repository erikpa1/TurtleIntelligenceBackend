package rvar

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// triangularGen draws from a triangular distribution over [min, max] with the
// given mode (peak). Handy when only a low/likely/high estimate is known.
type triangularGen struct {
	min, mode, max float64
	fc             float64 // (mode-min)/(max-min), the split point of the CDF
}

func (t triangularGen) Sample(rng *rand.Rand) float64 {
	u := rng.Float64()
	if u < t.fc {
		return t.min + math.Sqrt(u*(t.max-t.min)*(t.mode-t.min))
	}
	return t.max - math.Sqrt((1-u)*(t.max-t.min)*(t.max-t.mode))
}

func init() {
	// tri(min, mode, max) — triangular distribution.
	//   tri(5s, 10s, 30s)
	ctor := func(args []float64) (Generator, error) {
		min, mode, max := args[0], args[1], args[2]
		if !(min <= mode && mode <= max) {
			return nil, fmt.Errorf("tri: require min <= mode <= max, got %g, %g, %g", min, mode, max)
		}
		if max == min {
			return constant(min), nil
		}
		return triangularGen{min: min, mode: mode, max: max, fc: (mode - min) / (max - min)}, nil
	}
	Register("tri", 3, 3, ctor)
	Register("triangular", 3, 3, ctor)
}
