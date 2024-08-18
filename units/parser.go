// Modified code from ParseDuration (time/format.go)

package units

import (
	"errors"
	"fmt"
)

var (
	errLeadingInt = errors.New("time: bad [0-9]*") // never printed
)

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func parseUnit[U any](s string, unitMap map[string]uint64, convert func(number int64) U) (U, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return convert(0), nil
	}
	if s == "" {
		return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return convert(0), fmt.Errorf("unit: missing unit from input '%s'", orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return convert(0), fmt.Errorf("unit: unknown unit '%s' in input '%s'", u, orig)
		}
		if v > 1<<63/unit {
			// overflow
			return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
			}
		}
		d += v
		if d > 1<<63 {
			return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
		}
	}
	if neg {
		return convert(int64(-d)), nil
	}
	if d > 1<<63-1 {
		return convert(0), fmt.Errorf("unit: invalid input string '%s'", orig)
	}
	return convert(int64(d)), nil
}

// leadingInt consumes the leading [0-9]* from s.
func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > 1<<63 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}
