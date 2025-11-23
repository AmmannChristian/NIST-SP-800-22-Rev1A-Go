package nist

import (
	"testing"
)

func TestNonOverlappingTemplate(t *testing.T) {
	t.Run("wrong_template_size", func(t *testing.T) {
		data := make([]byte, 1000)
		p, pass := NonOverlappingTemplateTest(data, 5)
		if p != 0 || pass {
			t.Fatalf("expected reject on wrong template size, got p=%.6f pass=%v", p, pass)
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 10000)
		state := uint64(99)
		for i := range data {
			state = state*6364136223846793005 + 1442695040888963407
			data[i] = byte(state >> 56)
		}
		p, pass := NonOverlappingTemplateTest(data, 9)
		if p < 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
		if pass != (p >= Alpha) {
			t.Fatalf("pass flag inconsistent with p-value %.6f", p)
		}
	})
}
