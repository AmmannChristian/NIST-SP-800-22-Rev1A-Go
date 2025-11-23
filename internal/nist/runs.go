package nist

import (
	"math"
	"math/bits"
)

// RunsTest implements the NIST Runs test.
// It returns the p-value and whether it passes at Alpha.
func RunsTest(bitstream []byte) (float64, bool) {
	n := len(bitstream) * 8
	if n == 0 {
		return 0, false
	}

	ones := 0
	for _, b := range bitstream {
		ones += bits.OnesCount8(b)
	}

	pi := float64(ones) / float64(n)
	if math.Abs(pi-0.5) > 2.0/math.Sqrt(float64(n)) {
		// Precondition for the runs test is not met.
		return 0, false
	}

	runs := 1
	prev := bitAt(bitstream, 0)
	for i := 1; i < n; i++ {
		b := bitAt(bitstream, i)
		if b != prev {
			runs++
			prev = b
		}
	}

	erfcArg := math.Abs(float64(runs)-2.0*float64(n)*pi*(1-pi)) /
		(2.0 * math.Sqrt(2*float64(n)) * pi * (1 - pi))
	pValue := math.Erfc(erfcArg)

	return pValue, pValue >= Alpha
}
