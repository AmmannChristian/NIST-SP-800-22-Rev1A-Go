package nist

import (
	"testing"
)

func TestLinearComplexity(t *testing.T) {
	t.Run("N_equals_zero", func(t *testing.T) {
		data := make([]byte, 10)
		p, pass := LinearComplexityTest(data, 1024)
		if p >= 0.01 || pass {
			t.Errorf("expected reject when N = 0")
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 10000)
		for i := range data {
			data[i] = byte(i % 256)
		}
		p, _ := LinearComplexityTest(data, 500)
		if p == 0 {
			t.Errorf("expected non-zero p-value, got p=%.6f", p)
		}
	})
}
