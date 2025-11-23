package nist

import (
	"math"

	"gonum.org/v1/gonum/mathext"
)

// RandomExcursionsTest implements the NIST Random Excursions test.
// It returns the minimum p-value across the 8 states and whether it passes at Alpha.
func RandomExcursionsTest(bitstream []byte) (float64, bool) {
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

	cycle := make([]int, J+1)
	idx := 1
	for i := 0; i < n; i++ {
		if S[i] == 0 {
			cycle[idx] = i
			idx++
		}
	}
	cycle[J] = n - 1

	stateX := []int{-4, -3, -2, -1, 1, 2, 3, 4}
	pi := [][]float64{
		{0, 0, 0, 0, 0, 0},
		{0.5, 0.25, 0.125, 0.0625, 0.03125, 0.03125},
		{0.75, 0.0625, 0.046875, 0.03515625, 0.0263671875, 0.0791015625},
		{0.8333333333, 0.02777777778, 0.02314814815, 0.01929012346, 0.01607510288, 0.0803755143},
		{0.875, 0.015625, 0.013671875, 0.01196289063, 0.0104675293, 0.0732727051},
	}

	minP := 1.0
	cycleStart := 0
	cycleStop := cycle[1]

	nu := make([][]float64, 6)
	for i := range nu {
		nu[i] = make([]float64, 8)
	}

	for j := 1; j <= J; j++ {
		counter := make([]int, 8)
		for i := cycleStart; i <= cycleStop; i++ {
			val := S[i]
			if val >= 1 && val <= 4 {
				counter[val+3]++
			} else if val <= -1 && val >= -4 {
				idx := val + 4
				if idx >= 0 && idx < len(counter) {
					counter[idx]++
				}
			}
		}

		cycleStart = cycle[j] + 1
		if j < J {
			cycleStop = cycle[j+1]
		}

		for i := 0; i < 8; i++ {
			if counter[i] >= 0 && counter[i] <= 4 {
				nu[counter[i]][i]++
			} else if counter[i] >= 5 {
				nu[5][i]++
			}
		}
	}

	for i := 0; i < 8; i++ {
		x := stateX[i]
		idx := int(math.Abs(float64(x)))
		if idx >= len(pi) {
			continue
		}
		sum := 0.0
		for k := 0; k < 6; k++ {
			expected := float64(J) * pi[idx][k] // #nosec G602: idx guarded within pi bounds
			diff := nu[k][i] - expected
			sum += diff * diff / expected
		}
		p := mathext.GammaIncRegComp(2.5, sum/2.0)
		if p < minP {
			minP = p
		}
	}

	return minP, minP >= Alpha
}
