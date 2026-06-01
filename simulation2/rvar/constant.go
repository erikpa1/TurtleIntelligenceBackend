package rvar

import "math/rand/v2"

// constGen always returns the same value. It backs both bare literals (handled
// by the parser) and the explicit const(x) form.
type constGen struct{ v float64 }

func (c constGen) Sample(*rand.Rand) float64 { return c.v }

// constant builds a fixed-value generator.
func constant(v float64) Generator { return constGen{v: v} }

func init() {
	// const(x) — a fixed number of seconds. Mostly useful for symmetry with
	// the random forms; "x" on its own works identically.
	Register("const", 1, 1, func(args []float64) (Generator, error) {
		return constant(args[0]), nil
	})
}
