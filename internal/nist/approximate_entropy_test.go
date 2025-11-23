package nist

import (
	"testing"
)

func TestApproximateEntropy(t *testing.T) {
	t.Run("m_too_small", func(t *testing.T) {
		data := make([]byte, 100)
		p, pass := ApproximateEntropyTest(data, 0)
		if p != 0 || pass {
			t.Errorf("expected reject when m < 1")
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		// Deterministic pseudo-random stream for stability
		data := make([]byte, 10000)
		state := uint64(42)
		for i := range data {
			state = state*6364136223846793005 + 1442695040888963407
			data[i] = byte(state >> 56)
		}
		p, pass := ApproximateEntropyTest(data, 5)
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
		if !pass {
			t.Fatalf("expected structured but varied input to pass, got p=%.6f", p)
		}
	})
}
