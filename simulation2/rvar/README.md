# `rvar` — random variables for the simulation

`rvar` turns a small string expression into a **smart variable** that yields an
`int64` in the simulation's time domain. The expression can be a fixed time
literal or a random distribution. Every variable owns its own random number
pool, so two variables built from the same expression are statistically
independent.

```go
import "turtle/simulation2/rvar"

self.delayTime = rvar.NewRvarr("exp(10s)")   // compile once (panics on a bad expr)

wait := self.delayTime.GetInt64()            // draw a sample, in simulation units
```

## The value

| method          | returns  | meaning                                              |
| --------------- | -------- | ---------------------------------------------------- |
| `GetInt64()`    | `int64`  | next sample, rounded to whole units, clamped to ≥ 0  |
| `GetFloat64()`  | `float64`| next sample scaled to units, fractional              |
| `GetSeconds()`  | `float64`| next sample in raw seconds, unscaled                 |
| `Expr()`        | `string` | the original expression                              |

## Time units

By default **1 second == 1 unit**. The resolution is global and can be changed
once at startup:

```go
rvar.SetUnitsPerSecond(100)   // now 1 second == 100 units
rvar.NewRvarr("10s").GetInt64() // -> 1000
```

`GetSeconds()` is always in seconds regardless of the scale.

## Expressions

### Time literals

| form          | example      | seconds |
| ------------- | ------------ | ------- |
| plain number  | `"42"`       | 42      |
| seconds       | `"10s"`      | 10      |
| milliseconds  | `"500ms"`    | 0.5     |
| minutes       | `"5m"`       | 300     |
| hours         | `"2h"`       | 7200    |
| days          | `"1d"`       | 86400   |
| `mm:ss`       | `"00:10"`    | 10      |
| `hh:mm:ss`    | `"1:02:30"`  | 3750    |
| `d:hh:mm:ss`  | `"1:00:00:00"` | 86400 |

Colon strings read right-to-left as seconds, minutes, hours, days — matching
`tools.SecondsFromTimeString`. Arguments to distributions accept all of these
forms too, so `uniform(00:05, 00:30)` and `uniform(5s, 30s)` are equivalent.

### Distributions

| expression                       | description                                              |
| -------------------------------- | -------------------------------------------------------- |
| `const(x)`                       | fixed value (same as a bare literal `x`)                 |
| `uniform(min, max)`              | continuous uniform on `[min, max)`                       |
| `exp(mean)`                      | exponential with the given mean                          |
| `exp(mean, min)`                 | exponential, truncated below at `min`                    |
| `exp(mean, min, max)`            | exponential, truncated to `[min, max]`                   |
| `normal(mean, sd)`               | Gaussian (also `gauss(...)`)                             |
| `normal(mean, sd, min, max)`     | Gaussian truncated to `[min, max]`                       |
| `lognormal(mu, sigma)`           | log-normal; `mu`/`sigma` are of the underlying normal    |
| `tri(min, mode, max)`            | triangular (also `triangular(...)`)                      |
| `poisson(lambda)`                | Poisson **count** with mean `lambda` (a count, not a time) |

#### Examples

```go
rvar.NewRvarr("uniform(5s, 30s)")          // between 5 and 30 seconds
rvar.NewRvarr("exp(10s)")                  // mean 10s, occasional long tail
rvar.NewRvarr("exp(10s, 2s, 60s)")         // same, clamped to [2s, 60s]
rvar.NewRvarr("normal(60s, 5s)")           // ~60s, std-dev 5s
rvar.NewRvarr("normal(60s, 5s, 30s, 90s)") // truncated normal
rvar.NewRvarr("tri(5s, 10s, 30s)")         // low 5s, likely 10s, high 30s
rvar.NewRvarr("lognormal(2.3, 0.5)")       // long-tailed service time
rvar.NewRvarr("poisson(3)")                // 0,1,2,3,... arrivals per step
rvar.NewRvarr("00:10")                     // a fixed 10 seconds
```

## Reproducibility & random pools

Each `Rvar` is seeded from a shared source, giving every variable an
**independent stream**. To make a whole run reproducible, set the seed once
before constructing any variables:

```go
rvar.SetGlobalSeed(42)
a := rvar.NewRvarr("exp(10s)")   // same seed -> same sequence every run
```

Re-calling `SetGlobalSeed` with the same value and rebuilding the variables in
the same order reproduces the exact sequence. Different seeds give independent
runs. `Rvar` is **not** safe for concurrent use — give each goroutine its own.

## Error handling

`NewRvarr` panics on a malformed expression (fail-fast during setup). Use `New`
when you want to handle the error:

```go
v, err := rvar.New(userSuppliedExpr)
if err != nil {
    // report the bad expression
}
```

## Adding a distribution

The system is registration-based: drop a new file into the package and register
the distribution from its `init()`. No central file needs editing.

```go
// weibull.go
package rvar

import "math/rand/v2"

type weibullGen struct{ scale, shape float64 }

func (w weibullGen) Sample(rng *rand.Rand) float64 { /* ... */ }

func init() {
    // weibull(scale, shape)
    Register("weibull", 2, 2, func(args []float64) (Generator, error) {
        return weibullGen{scale: args[0], shape: args[1]}, nil
    })
}
```

`Register(name, minArgs, maxArgs, constructor)` — `maxArgs == -1` means
unbounded. Arguments arrive already parsed into `float64` seconds. The
`Generator` interface is a single method:

```go
type Generator interface {
    Sample(rng *rand.Rand) float64 // one observation, in seconds
}
```
