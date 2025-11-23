package nist

import (
	"math"

	"gonum.org/v1/gonum/mathext"
)

// SerialTest implements the NIST Serial test with a fixed block length m.
// It returns the minimum p-value across the two computed statistics and whether it passes at Alpha.
func SerialTest(bitstream []byte, m int) (float64, bool) {
	bits := expandBits(bitstream)
	n := len(bits)
	if n == 0 || m < 2 {
		return 0, false
	}

	psim0 := psi2(bits, m)
	psim1 := psi2(bits, m-1)
	psim2 := psi2(bits, m-2)

	del1 := psim0 - psim1
	del2 := psim0 - 2.0*psim1 + psim2

	p1 := mathext.GammaIncRegComp(math.Pow(2, float64(m-1))/2.0, del1/2.0)
	p2 := mathext.GammaIncRegComp(math.Pow(2, float64(m-2))/2.0, del2/2.0)

	pValue := math.Min(p1, p2)

	return pValue, pValue >= Alpha
}
