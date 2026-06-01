package rvar

import (
	"math"
	"testing"
)

func TestParseLiterals(t *testing.T) {
	cases := map[string]float64{
		"42":         42,
		"3.5":        3.5,
		"10s":        10,
		"5m":         300,
		"2h":         7200,
		"500ms":      0.5,
		"1d":         86400,
		"00:10":      10,
		"1:02:30":    3750,
		"1:00:00:00": 86400,
	}
	for expr, want := range cases {
		got, err := parseScalar(expr)
		if err != nil {
			t.Fatalf("parseScalar(%q) error: %v", expr, err)
		}
		if got != want {
			t.Errorf("parseScalar(%q) = %g, want %g", expr, got, want)
		}
	}
}

func TestConstantInt64(t *testing.T) {
	r := NewRvarr("10s")
	if got := r.GetInt64(); got != 10 {
		t.Errorf("GetInt64 = %d, want 10", got)
	}
}

func TestUnitsPerSecondScaling(t *testing.T) {
	SetUnitsPerSecond(100)
	defer SetUnitsPerSecond(1)

	r := NewRvarr("10s")
	if got := r.GetInt64(); got != 1000 {
		t.Errorf("GetInt64 with 100 units/s = %d, want 1000", got)
	}
}

func TestUniformInRange(t *testing.T) {
	SetGlobalSeed(1)
	r := NewRvarr("uniform(5s, 30s)")
	for i := 0; i < 1000; i++ {
		v := r.GetSeconds()
		if v < 5 || v >= 30 {
			t.Fatalf("uniform sample %g out of [5,30)", v)
		}
	}
}

func TestExponentialMean(t *testing.T) {
	SetGlobalSeed(2)
	r := NewRvarr("exp(10s)")
	var sum float64
	const n = 200000
	for i := 0; i < n; i++ {
		sum += r.GetSeconds()
	}
	mean := sum / n
	if math.Abs(mean-10) > 0.3 {
		t.Errorf("exponential mean = %g, want ~10", mean)
	}
}

func TestTriangularBounds(t *testing.T) {
	SetGlobalSeed(3)
	r := NewRvarr("tri(5s, 10s, 30s)")
	for i := 0; i < 1000; i++ {
		v := r.GetSeconds()
		if v < 5 || v > 30 {
			t.Fatalf("triangular sample %g out of [5,30]", v)
		}
	}
}

func TestNormalTruncated(t *testing.T) {
	SetGlobalSeed(4)
	r := NewRvarr("normal(60s, 20s, 50s, 70s)")
	for i := 0; i < 1000; i++ {
		v := r.GetSeconds()
		if v < 50 || v > 70 {
			t.Fatalf("truncated normal sample %g out of [50,70]", v)
		}
	}
}

func TestPoissonNonNegativeInt(t *testing.T) {
	SetGlobalSeed(5)
	r := NewRvarr("poisson(3)")
	for i := 0; i < 1000; i++ {
		v := r.GetSeconds()
		if v < 0 || v != math.Trunc(v) {
			t.Fatalf("poisson sample %g is not a non-negative integer", v)
		}
	}
}

func TestReproducibleStreams(t *testing.T) {
	SetGlobalSeed(99)
	a := NewRvarr("exp(10s)")
	first := []int64{a.GetInt64(), a.GetInt64(), a.GetInt64()}

	SetGlobalSeed(99)
	b := NewRvarr("exp(10s)")
	second := []int64{b.GetInt64(), b.GetInt64(), b.GetInt64()}

	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("streams diverged at %d: %d vs %d", i, first[i], second[i])
		}
	}
}

func TestIndependentPools(t *testing.T) {
	SetGlobalSeed(7)
	a := NewRvarr("uniform(0s, 100s)")
	b := NewRvarr("uniform(0s, 100s)")
	// Two distinct Rvars must draw from different streams.
	identical := true
	for i := 0; i < 50; i++ {
		if a.GetSeconds() != b.GetSeconds() {
			identical = false
			break
		}
	}
	if identical {
		t.Error("two Rvars produced identical streams; pools are not independent")
	}
}

func TestErrors(t *testing.T) {
	for _, expr := range []string{
		"uniform(30s, 5s)", // max < min
		"exp(0s)",          // mean must be > 0
		"uniform(1s)",      // wrong arity
		"bogus(1)",         // unknown distribution
		"not_a_time",       // unparseable literal
	} {
		if _, err := New(expr); err == nil {
			t.Errorf("New(%q) expected error, got nil", expr)
		}
	}
}
