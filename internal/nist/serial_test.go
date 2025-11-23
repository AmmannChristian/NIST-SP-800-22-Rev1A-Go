package nist

import (
	"testing"
)

func TestSerial(t *testing.T) {
	t.Run("m_too_small", func(t *testing.T) {
		data := make([]byte, 100)
		p, pass := SerialTest(data, 1)
		if p != 0 || pass {
			t.Errorf("expected reject when m < 2")
		}
	})

	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 100000)
		state := uint64(777)
		for i := range data {
			state = state*6364136223846793005 + 1442695040888963407
			data[i] = byte(state >> 56)
		}
		p, pass := SerialTest(data, 3)
		if p < 0 || p > 1 {
			t.Fatalf("p-value out of range: %.6f", p)
		}
		if pass != (p >= Alpha) {
			t.Fatalf("pass flag inconsistent with p-value %.6f", p)
		}
	})
}
