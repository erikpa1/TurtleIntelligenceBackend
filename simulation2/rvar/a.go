// Package rvar implements "random variables" for the event-driven simulation.
//
// An Rvar is a smart value built from a string expression. The expression can
// be a fixed time literal ("00:10", "10s", "2h", "42") or a random distribution
// written as a function call ("exp(10s)", "uniform(5s, 30s)", "normal(60s, 5s)").
//
// Every Rvar evaluates to an int64 in the simulation time domain. By default one
// real second equals one integer unit, but this can be rescaled globally (e.g.
// 1 second == 100 units) with SetUnitsPerSecond — see core.go.
//
// Typical use:
//
//	self.delayTime = rvar.NewRvarr("exp(10s)")
//	...
//	wait := self.delayTime.GetInt64()
//
// Architecture
//
//   - core.go        Rvar type, constructors, GetInt64, time-unit scaling, seeding.
//   - registry.go    Generator interface + name->constructor registry (Register).
//   - parse.go       expression parser and time-literal/duration parsing.
//   - constant.go    fixed values and time literals.
//   - uniform.go     uniform(min, max)
//   - exponential.go exp(mean[, min[, max]])
//   - normal.go      normal(mean, sd[, min, max])
//   - lognormal.go   lognormal(mu, sigma)
//   - triangular.go  tri(min, mode, max)
//   - poisson.go     poisson(lambda)
//
// Each distribution lives in its own file and registers itself in an init()
// function, so adding a new distribution is a single self-contained file — no
// central switch to edit. See README.md for the full reference and examples.
package rvar
