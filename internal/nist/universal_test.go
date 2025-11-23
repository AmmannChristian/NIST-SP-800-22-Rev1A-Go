package nist

import (
	"testing"
)

func TestUniversalStatistical(t *testing.T) {
	t.Run("too_short", func(t *testing.T) {
		data := make([]byte, 1000)
		p, pass := UniversalStatisticalTest(data)
		if p >= 0.01 || pass {
			t.Errorf("expected reject on short input")
		}
	})

	t.Run("minimum_valid_input", func(t *testing.T) {
		data := make([]byte, 50000)
		for i := range data {
			data[i] = byte(i % 256)
		}
		p, _ := UniversalStatisticalTest(data)
		if p == 0 {
			t.Errorf("expected non-zero p-value, got p=%.6f", p)
		}
	})
}
