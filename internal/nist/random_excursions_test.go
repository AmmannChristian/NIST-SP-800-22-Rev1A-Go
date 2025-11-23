package nist

import (
	"testing"
)

func TestRandomExcursions(t *testing.T) {
	t.Run("insufficient_cycles", func(t *testing.T) {
		data := make([]byte, 100)
		for i := range data {
			data[i] = 0xAA
		}
		p, pass := RandomExcursionsTest(data)
		if p != 0 || pass {
			t.Fatalf("expected failure for insufficient cycles, got p=%.6f pass=%v", p, pass)
		}
	})

	t.Run("larger_input", func(t *testing.T) {
		// Alternating bits produce many zero crossings; should meet constraint
		data := make([]byte, 125) // 1000 bits
		for i := range data {
			data[i] = 0xAA
		}
		p, pass := RandomExcursionsTest(data)
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
		if pass {
			t.Fatalf("expecting periodic walk to fail uniformity, got p=%.6f", p)
		}
	})
}
