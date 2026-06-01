package rvar

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// poissonGen returns Poisson-distributed counts with mean lambda, using Knuth's
// algorithm. Unlike the other generators its output is a whole number, so it is
// typically used for "how many" rather than "how long".
type poissonGen struct {
	lambda float64
	limit  float64 // precomputed e^-lambda
}

func (p poissonGen) Sample(rng *rand.Rand) float64 {
	k := 0
	prod := 1.0
	for {
		prod *= rng.Float64()
		if prod <= p.limit {
			return float64(k)
		}
		k++
	}
}

func init() {
	// poisson(lambda) — number of events with mean lambda (a count, not a time).
	//   poisson(3)
	Register("poisson", 1, 1, func(args []float64) (Generator, error) {
		lambda := args[0]
		if lambda <= 0 {
			return nil, fmt.Errorf("poisson: lambda must be > 0, got %g", lambda)
		}
		return poissonGen{lambda: lambda, limit: math.Exp(-lambda)}, nil
	})
}
