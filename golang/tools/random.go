package tools

import "math/rand"

func RandomFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
