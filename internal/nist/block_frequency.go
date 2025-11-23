package nist

import "gonum.org/v1/gonum/mathext"

// BlockFrequencyTest implements the NIST Block Frequency test.
// blockSize is the length of each block in bits (M in the NIST documentation).
// It returns the p-value and whether it passes at Alpha.
func BlockFrequencyTest(bitstream []byte, blockSize int) (float64, bool) {
	n := len(bitstream) * 8
	if blockSize <= 0 || n < blockSize {
		return 0, false
	}

	N := n / blockSize // number of complete blocks
	if N == 0 {
		return 0, false
	}

	var sum float64
	for block := 0; block < N; block++ {
		blockSum := 0
		offset := block * blockSize
		for j := 0; j < blockSize; j++ {
			bit := bitAt(bitstream, offset+j)
			blockSum += int(bit)
		}

		pi := float64(blockSum) / float64(blockSize)
		v := pi - 0.5
		sum += v * v
	}

	chiSquared := 4 * float64(blockSize) * sum
	pValue := mathext.GammaIncRegComp(float64(N)/2.0, chiSquared/2.0)

	return pValue, pValue >= Alpha
}
