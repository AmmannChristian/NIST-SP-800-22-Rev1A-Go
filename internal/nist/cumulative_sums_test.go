package nist

import (
	"testing"
)

func TestCumulativeSums(t *testing.T) {
	t.Run("empty_input", func(t *testing.T) {
		p, pass := CumulativeSumsTest(nil)
		if p >= 0.01 || pass {
			t.Errorf("expected reject on empty input")
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 125)
		for i := range data {
			data[i] = 0xAA
		}
		p, _ := CumulativeSumsTest(data)
		if p == 0 {
			t.Errorf("expected non-zero p-value, got p=%.6f", p)
		}
	})
}
