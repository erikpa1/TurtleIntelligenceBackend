package rvar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// callRe matches a distribution call like "exp(10s, 30s)". Group 1 is the
// function name, group 2 is the raw argument list (possibly empty).
var callRe = regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\((.*)\)\s*$`)

// compile turns an expression into a Generator. It first tries to read the
// expression as a function call; otherwise it falls back to a fixed literal
// (a plain number or a time string such as "00:10" / "10s").
func compile(expr string) (Generator, error) {
	if m := callRe.FindStringSubmatch(expr); m != nil {
		name := m[1]
		args, err := parseArgs(m[2])
		if err != nil {
			return nil, fmt.Errorf("in %q: %w", expr, err)
		}
		return build(name, args)
	}

	secs, err := parseScalar(expr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %q: %w", expr, err)
	}
	return constant(secs), nil
}

// parseArgs splits a comma-separated argument list and converts each to seconds.
// An empty (or whitespace-only) list yields no arguments.
func parseArgs(raw string) ([]float64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	args := make([]float64, len(parts))
	for i, p := range parts {
		v, err := parseScalar(p)
		if err != nil {
			return nil, fmt.Errorf("argument %d: %w", i+1, err)
		}
		args[i] = v
	}
	return args, nil
}

// suffixRe matches a number with an optional time-unit suffix, e.g. "10s",
// "2.5h", "500ms", or just "42".
var suffixRe = regexp.MustCompile(`^([0-9]*\.?[0-9]+)\s*([a-zA-Z]*)$`)

// parseScalar converts a single token into seconds. It understands:
//
//   - plain numbers:        "42", "3.5"   -> taken as seconds
//   - suffix units:         "10s", "5m", "2h", "500ms", "1d"
//   - colon time strings:   "00:10" (mm:ss), "1:02:30" (hh:mm:ss), "d:h:m:s"
//
// For colon strings the rightmost component is always seconds, matching the
// rest of the codebase (tools.SecondsFromTimeString).
func parseScalar(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty value")
	}

	if strings.Contains(s, ":") {
		return parseColon(s)
	}

	m := suffixRe.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("not a number or duration: %q", s)
	}
	n, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("bad number %q: %w", m[1], err)
	}

	switch strings.ToLower(m[2]) {
	case "", "s", "sec", "secs":
		return n, nil
	case "ms":
		return n / 1000, nil
	case "m", "min", "mins":
		return n * 60, nil
	case "h", "hr", "hrs":
		return n * 3600, nil
	case "d", "day", "days":
		return n * 86400, nil
	default:
		return 0, fmt.Errorf("unknown time unit %q", m[2])
	}
}

// parseColon converts a colon-separated time string to seconds, reading
// components right-to-left as seconds, minutes, hours, days.
func parseColon(s string) (float64, error) {
	parts := strings.Split(s, ":")
	mults := []float64{1, 60, 3600, 86400}
	if len(parts) > len(mults) {
		return 0, fmt.Errorf("too many time components in %q", s)
	}
	var total float64
	for i := 0; i < len(parts); i++ {
		// parts[len-1] is seconds (mult 1), parts[len-2] minutes, etc.
		comp := strings.TrimSpace(parts[len(parts)-1-i])
		if comp == "" {
			comp = "0"
		}
		n, err := strconv.ParseFloat(comp, 64)
		if err != nil {
			return 0, fmt.Errorf("bad time component %q: %w", comp, err)
		}
		total += n * mults[i]
	}
	return total, nil
}
