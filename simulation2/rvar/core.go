package rvar

import (
	"math"
	"math/rand/v2"
	"sync"
	"sync/atomic"
)

// ---------------------------------------------------------------------------
// Time-unit scaling
// ---------------------------------------------------------------------------

// unitsPerSecond is how many integer units make up one real second. It defaults
// to 1 (one second == one unit) but a simulation can rescale time resolution at
// startup, e.g. SetUnitsPerSecond(100) makes one second == 100 units.
//
// It is read on every GetInt64 call, so changing it affects all existing Rvars.
var unitsPerSecond atomic.Int64

func init() { unitsPerSecond.Store(1) }

// SetUnitsPerSecond sets the global time resolution. n must be >= 1; smaller
// values are clamped to 1. Call this once during simulation setup, before time
// values are consumed.
func SetUnitsPerSecond(n int64) {
	if n < 1 {
		n = 1
	}
	unitsPerSecond.Store(n)
}

// UnitsPerSecond returns the current time resolution (units per real second).
func UnitsPerSecond() int64 { return unitsPerSecond.Load() }

// ---------------------------------------------------------------------------
// Random seeding — each Rvar owns an independent stream ("its own pool")
// ---------------------------------------------------------------------------

const (
	defaultSeedHi uint64 = 0x2545F4914F6CDD1D
	defaultSeedLo uint64 = 0x9E3779B97F4A7C15
)

var (
	seedMu  sync.Mutex
	seedSrc = rand.New(rand.NewPCG(defaultSeedHi, defaultSeedLo))
)

// SetGlobalSeed reseeds the source from which every new Rvar draws its own
// stream seed. Call it once before constructing any Rvar to make a whole
// simulation reproducible; call it with different seeds for independent runs.
func SetGlobalSeed(seed uint64) {
	seedMu.Lock()
	defer seedMu.Unlock()
	seedSrc = rand.New(rand.NewPCG(seed, seed^defaultSeedLo))
}

// newRng hands out a fresh, independent random pool for one Rvar. Streams are
// decorrelated because their PCG seeds come from the shared source.
func newRng() *rand.Rand {
	seedMu.Lock()
	defer seedMu.Unlock()
	return rand.New(rand.NewPCG(seedSrc.Uint64(), seedSrc.Uint64()))
}

// ---------------------------------------------------------------------------
// Rvar
// ---------------------------------------------------------------------------

// Rvar is a random (or fixed) variable compiled from a string expression. It
// owns its own random number pool, so two Rvars built from the same expression
// produce independent streams. Rvar is NOT safe for concurrent use by multiple
// goroutines; give each goroutine its own Rvar.
type Rvar struct {
	expr string
	gen  Generator
	rng  *rand.Rand
}

// New compiles expr into an Rvar, returning an error for malformed expressions.
func New(expr string) (*Rvar, error) {
	gen, err := compile(expr)
	if err != nil {
		return nil, err
	}
	return &Rvar{expr: expr, gen: gen, rng: newRng()}, nil
}

// NewRvarr compiles expr and panics on any error. It is the convenient form for
// simulation setup, where a bad expression is a programming mistake that should
// fail fast:
//
//	self.delayTime = rvar.NewRvarr("exp(10s)")
func NewRvarr(expr string) *Rvar {
	r, err := New(expr)
	if err != nil {
		panic("rvar: " + err.Error())
	}
	return r
}

// GetSeconds returns the next sample in seconds, unscaled and unrounded.
func (r *Rvar) GetSeconds() float64 { return r.gen.Sample(r.rng) }

// GetFloat64 returns the next sample scaled to simulation units (still
// fractional). Use GetInt64 for whole simulation ticks.
func (r *Rvar) GetFloat64() float64 {
	return r.gen.Sample(r.rng) * float64(unitsPerSecond.Load())
}

// GetInt64 returns the next sample as a whole number of simulation units.
// Negative samples (possible for normal/lognormal tails) are clamped to 0,
// since time delays cannot be negative.
func (r *Rvar) GetInt64() int64 {
	v := math.Round(r.GetFloat64())
	if v < 0 {
		v = 0
	}
	return int64(v)
}

// Expr returns the original expression the Rvar was built from.
func (r *Rvar) Expr() string { return r.expr }
