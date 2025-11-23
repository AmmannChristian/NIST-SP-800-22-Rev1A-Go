package nist

import (
	"testing"
)

func TestFrequencyMonobit(t *testing.T) {
	t.Run("empty_input", func(t *testing.T) {
		p, pass := FrequencyTest(nil)
		if p != 0 || pass {
			t.Errorf("expected reject on empty input, got p=%.6f pass=%v", p, pass)
		}
	})

	t.Run("all_zeros", func(t *testing.T) {
		data := make([]byte, 125)
		p, pass := FrequencyTest(data)
		if pass || p >= Alpha {
			t.Fatalf("expected failure for biased stream, got p=%.6f pass=%v", p, pass)
		}
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
	})

	t.Run("all_ones", func(t *testing.T) {
		data := make([]byte, 125)
		for i := range data {
			data[i] = 0xFF
		}
		p, pass := FrequencyTest(data)
		if pass || p >= Alpha {
			t.Fatalf("expected failure for biased stream, got p=%.6f pass=%v", p, pass)
		}
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
	})

	t.Run("alternating_bits", func(t *testing.T) {
		data := make([]byte, 125)
		for i := range data {
			data[i] = 0xAA // 10101010
		}
		p, pass := FrequencyTest(data)
		if !pass || p < Alpha {
			t.Fatalf("expected alternating bits to pass, got p=%.6f pass=%v", p, pass)
		}
		if p <= 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
	})
}
