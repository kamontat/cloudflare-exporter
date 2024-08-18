package units

import (
	"strconv"
)

func format(in int64, buf *[32]byte, unitMap map[string]uint64, order []string) int {
	w := 0

	u := uint64(in)
	neg := in < 0
	if neg {
		u = -u
		buf[w] = '-'
		w++
	}

	for _, key := range order {
		unitS := unitMap[key]

		if u > unitS {
			// Add number
			w += fmtInt(buf[w:], u/unitS)
			// Add unit
			w += fmtUnit(buf[w:], key)
			u %= unitS
		}

		if u <= 0 {
			break
		}
	}

	if u > 0 && false {
		// FIXME: update frac for unit contains fraction number
		w += fmtFrac(buf[w:], u, 9)
	}

	return w
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros. It omits the decimal
// point too when the fraction is 0. It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) int {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) (w int) {
	if v == 0 {
		buf[w] = '0'
		w++
		return
	}

	vs := strconv.FormatUint(v, 10)
	newBuf := []byte(vs)
	copy(buf, newBuf)
	w += len(vs)

	return
}

func fmtUnit(buf []byte, unit string) (w int) {
	for _, char := range unit {
		buf[w] = byte(char)
		w++
	}
	return
}
