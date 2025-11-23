package nist

import (
	"testing"
)

func TestBinaryMatrixRank(t *testing.T) {
	t.Run("insufficient_bits", func(t *testing.T) {
		data := make([]byte, 100)
		p, pass := BinaryMatrixRankTest(data)
		if p >= 0.01 || pass {
			t.Errorf("expected reject on insufficient bits")
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 1024)
		for i := range data {
			data[i] = byte(i % 256)
		}
		p, _ := BinaryMatrixRankTest(data)
		if p == 0 {
			t.Errorf("expected non-zero p-value, got p=%.6f", p)
		}
	})
}
