package nist

import "math"

// RandomExcursionsVariantTest implements the NIST Random Excursions Variant test.
// It returns the minimum p-value across the 18 states and whether it passes at Alpha.
func RandomExcursionsVariantTest(bitstream []byte) (float64, bool) {
	bits := expandBits(bitstream)
	n := len(bits)

	S := make([]int, n)
	S[0] = 2*int(bits[0]) - 1
	J := 0
	for i := 1; i < n; i++ {
		S[i] = S[i-1] + 2*int(bits[i]) - 1
		if S[i] == 0 {
			J++
		}
	}
	if S[n-1] != 0 {
		J++
	}

	constraint := int(math.Max(0.005*math.Sqrt(float64(n)), 500))
	if J < constraint {
		return 0, false
	}

	stateX := []int{-9, -8, -7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	minP := 1.0

	for _, x := range stateX {
		count := 0
		for i := 0; i < n; i++ {
			if S[i] == x {
				count++
			}
		}
		p := math.Erfc(math.Abs(float64(count)-float64(J)) / math.Sqrt(2*float64(J)*(4*math.Abs(float64(x))-2)))
		if p < minP {
			minP = p
		}
	}

	return minP, minP >= Alpha
}
