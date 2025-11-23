package nist

import (
	"testing"
)

func TestBlockFrequency(t *testing.T) {
	t.Run("empty_input", func(t *testing.T) {
		p, pass := BlockFrequencyTest(nil, 128)
		if p >= 0.01 || pass {
			t.Errorf("expected reject on empty input")
		}
	})

	t.Run("insufficient_data", func(t *testing.T) {
		data := make([]byte, 10)
		p, pass := BlockFrequencyTest(data, 128)
		if p >= 0.01 || pass {
			t.Errorf("expected reject when data < block size")
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 256)
		for i := range data {
			data[i] = 0xAA
		}
		p, _ := BlockFrequencyTest(data, 128)
		if p == 0 {
			t.Errorf("expected non-zero p-value, got p=%.6f", p)
		}
	})
}
