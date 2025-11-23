package nist

import (
	"math"

	"gonum.org/v1/gonum/mathext"
)

// ApproximateEntropyTest implements the NIST Approximate Entropy test.
// It returns the p-value and whether it passes at Alpha.
func ApproximateEntropyTest(bitstream []byte, m int) (float64, bool) {
	bits := expandBits(bitstream)
	n := len(bits)
	if n == 0 || m < 1 {
		return 0, false
	}

	var apEn [2]float64
	for r := 0; r < 2; r++ {
		blockSize := m + r
		numBlocks := n
		powLen := (1 << (blockSize + 1)) - 1
		P := make([]int, powLen)

		for i := 0; i < numBlocks; i++ {
			k := 1
			for j := 0; j < blockSize; j++ {
				k <<= 1
				if bits[(i+j)%n] == 1 {
					k++
				}
			}
			P[k-1]++
		}

		sum := 0.0
		index := (1 << blockSize) - 1
		for i := 0; i < (1 << blockSize); i++ {
			if P[index] > 0 {
				sum += float64(P[index]) * math.Log(float64(P[index])/float64(numBlocks))
			}
			index++
		}
		apEn[r] = sum / float64(numBlocks)
	}

	apen := apEn[0] - apEn[1]
	chiSquared := 2.0 * float64(n) * (math.Log(2) - apen)
	pValue := mathext.GammaIncRegComp(math.Pow(2, float64(m-1)), chiSquared/2.0)

	return pValue, pValue >= Alpha
}
