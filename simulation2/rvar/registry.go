package rvar

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

// Generator produces samples in the seconds domain. Implementations must be
// deterministic given the same rng stream and must not retain the rng.
type Generator interface {
	// Sample draws one observation, in seconds.
	Sample(rng *rand.Rand) float64
}

// GeneratorFunc adapts a plain function into a Generator.
type GeneratorFunc func(rng *rand.Rand) float64

func (f GeneratorFunc) Sample(rng *rand.Rand) float64 { return f(rng) }

// Constructor builds a Generator from the positional arguments parsed out of an
// expression. Each argument has already been converted to a float64 in the
// seconds domain (see parseScalar). Constructors should validate their input
// and return a descriptive error rather than panicking.
type Constructor func(args []float64) (Generator, error)

type registration struct {
	minArgs int
	maxArgs int // -1 means unbounded
	build   Constructor
}

// registry maps a function name (e.g. "exp") to its constructor. It is written
// only from init() functions, so no locking is required at runtime.
var registry = map[string]registration{}

// Register makes a distribution available under name. It accepts between
// minArgs and maxArgs positional arguments (maxArgs == -1 for unbounded).
// It panics on a duplicate name, which can only happen at init time.
func Register(name string, minArgs, maxArgs int, c Constructor) {
	if _, dup := registry[name]; dup {
		panic("rvar: distribution already registered: " + name)
	}
	registry[name] = registration{minArgs: minArgs, maxArgs: maxArgs, build: c}
}

// build looks up name, validates arity, and constructs the Generator.
func build(name string, args []float64) (Generator, error) {
	reg, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown distribution %q (known: %s)", name, knownNames())
	}
	if len(args) < reg.minArgs || (reg.maxArgs >= 0 && len(args) > reg.maxArgs) {
		return nil, fmt.Errorf("%s expects %s arguments, got %d", name, arityText(reg), len(args))
	}
	return reg.build(args)
}

func arityText(reg registration) string {
	switch {
	case reg.maxArgs < 0:
		return fmt.Sprintf("at least %d", reg.minArgs)
	case reg.minArgs == reg.maxArgs:
		return fmt.Sprintf("exactly %d", reg.minArgs)
	default:
		return fmt.Sprintf("%d-%d", reg.minArgs, reg.maxArgs)
	}
}

func knownNames() string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
