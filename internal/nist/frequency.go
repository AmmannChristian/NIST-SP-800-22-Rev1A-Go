package nist

import (
	"math"
	"math/bits"
)

// FrequencyTest implements the NIST Monobit (Frequency) test.
// It returns the p-value and whether it passes at Alpha.
func FrequencyTest(bitstream []byte) (float64, bool) {
	n := len(bitstream) * 8
	if n == 0 {
		return 0, false
	}

	ones := 0
	for _, b := range bitstream {
		ones += bits.OnesCount8(b)
	}

	sum := float64(2*ones - n)
	sObs := math.Abs(sum) / math.Sqrt(float64(n))
	pValue := math.Erfc(sObs / math.Sqrt2)

	return pValue, pValue >= Alpha
}
