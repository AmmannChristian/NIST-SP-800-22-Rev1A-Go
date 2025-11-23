package nist

import (
	"testing"
)

func TestOverlappingTemplate(t *testing.T) {
	t.Run("valid_input", func(t *testing.T) {
		data := make([]byte, 100000)
		for i := range data {
			data[i] = 0xAA
		}
		p, _ := OverlappingTemplateTest(data, 9)
		if p == 0 {
			t.Errorf("expected non-zero p-value, got p=%.6f", p)
		}
	})
}
