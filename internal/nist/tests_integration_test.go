package nist

import (
	"testing"
)

// TestRunAllTests tests the RunAllTests function
func TestRunAllTests(t *testing.T) {
	t.Run("insufficient bits", func(t *testing.T) {
		data := make([]byte, 100) // Too small
		_, err := RunAllTests(data)
		if err == nil {
			t.Error("expected error for insufficient bits")
		}
	})

	t.Run("valid input returns 15 results", func(t *testing.T) {
		// Minimum 387,840 bits = 48,480 bytes
		data := make([]byte, 50000)
		for i := range data {
			data[i] = byte(i % 256)
		}

		results, err := RunAllTests(data)
		if err != nil {
			t.Fatalf("RunAllTests failed: %v", err)
		}

		if len(results) != 15 {
			t.Errorf("expected 15 results, got %d", len(results))
		}

		// Check all results have names
		for _, r := range results {
			if r.Name == "" {
				t.Errorf("result missing name: %+v", r)
			}
		}
	})

	t.Run("all test names are unique", func(t *testing.T) {
		data := make([]byte, 50000)
		for i := range data {
			data[i] = byte(i % 256)
		}

		results, err := RunAllTests(data)
		if err != nil {
			t.Fatalf("RunAllTests failed: %v", err)
		}

		names := make(map[string]bool)
		for _, r := range results {
			if names[r.Name] {
				t.Errorf("duplicate test name: %s", r.Name)
			}
			names[r.Name] = true
		}
	})
}
