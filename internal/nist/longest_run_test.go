package nist

import (
	"testing"
)

func TestLongestRunOfOnes(t *testing.T) {
	t.Run("too_short", func(t *testing.T) {
		data := make([]byte, 10)
		p, pass := LongestRunOfOnesTest(data)
		if p != 0 || pass {
			t.Errorf("expected reject on short input, got p=%.6f pass=%v", p, pass)
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 1000)
		for i := range data {
			data[i] = 0xF0
		}
		p, pass := LongestRunOfOnesTest(data)
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
		if pass {
			t.Fatalf("expected biased input to fail, got p=%.6f", p)
		}
	})
}
