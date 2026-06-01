package rvar

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// lognormalGen is a log-normal distribution: exp(mu + sigma*N(0,1)). It is
// always non-negative, which suits service/processing times with a long tail.
type lognormalGen struct{ mu, sigma float64 }

func (l lognormalGen) Sample(rng *rand.Rand) float64 {
	return math.Exp(l.mu + l.sigma*rng.NormFloat64())
}

func init() {
	// lognormal(mu, sigma) — mu and sigma are the mean and std-dev of the
	// underlying normal (in log-seconds), NOT of the result.
	//   lognormal(2.3, 0.5)
	Register("lognormal", 2, 2, func(args []float64) (Generator, error) {
		mu, sigma := args[0], args[1]
		if sigma < 0 {
			return nil, fmt.Errorf("lognormal: sigma must be >= 0, got %g", sigma)
		}
		return lognormalGen{mu: mu, sigma: sigma}, nil
	})
}
