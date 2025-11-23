package nist

import (
	"testing"
)

func TestDiscreteFourierTransform(t *testing.T) {
	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 1000)
		for i := range data {
			data[i] = 0xCC
		}
		p, pass := DiscreteFourierTransformTest(data)
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
		if pass {
			t.Fatalf("expected periodic pattern to fail DFT test, got p=%.6f", p)
		}
	})
}
