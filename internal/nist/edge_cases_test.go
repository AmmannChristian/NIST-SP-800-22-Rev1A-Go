package nist

import (
	"testing"
)

// TestAlpha tests the global Alpha constant
func TestConstants(t *testing.T) {
	if Alpha != 0.01 {
		t.Errorf("expected Alpha=0.01, got Alpha=%.2f", Alpha)
	}
}

// TestBitExtraction tests bit extraction from bytes
func TestBitPatterns(t *testing.T) {
	t.Run("all_patterns", func(t *testing.T) {
		// Test various bit patterns to increase coverage
		patterns := []byte{
			0x00, // All zeros
			0xFF, // All ones
			0xAA, // 10101010
			0x55, // 01010101
			0xF0, // 11110000
			0x0F, // 00001111
			0xCC, // 11001100
			0x33, // 00110011
		}

		for _, pattern := range patterns {
			data := make([]byte, 1000)
			for i := range data {
				data[i] = pattern
			}

			// Run a simple test to exercise the code paths
			if p, _ := FrequencyTest(data); p < 0 || p > 1 {
				t.Fatalf("frequency p-value out of range for pattern 0x%X: %.6f", pattern, p)
			}
			if p, _ := BlockFrequencyTest(data, 128); p < 0 || p > 1 {
				t.Fatalf("block frequency p-value out of range for pattern 0x%X: %.6f", pattern, p)
			}
			if p, _ := RunsTest(data); p < 0 || p > 1 {
				t.Fatalf("runs p-value out of range for pattern 0x%X: %.6f", pattern, p)
			}
		}
	})
}

// TestEdgeCases tests edge cases for various tests
func TestEdgeCases(t *testing.T) {
	t.Run("very_small_data", func(t *testing.T) {
		data := []byte{0xFF}
		if p, _ := FrequencyTest(data); p < 0 || p > 1 {
			t.Fatalf("frequency p-value out of range: %.6f", p)
		}
		if p, _ := RunsTest(data); p < 0 || p > 1 {
			t.Fatalf("runs p-value out of range: %.6f", p)
		}
		if p, _ := CumulativeSumsTest(data); p < 0 || p > 1 {
			t.Fatalf("cumulative sums p-value out of range: %.6f", p)
		}
	})

	t.Run("medium_data", func(t *testing.T) {
		data := make([]byte, 500)
		for i := range data {
			data[i] = byte(i % 256)
		}
		if p, _ := FrequencyTest(data); p < 0 || p > 1 {
			t.Fatalf("frequency p-value out of range: %.6f", p)
		}
		if p, _ := BlockFrequencyTest(data, 128); p < 0 || p > 1 {
			t.Fatalf("block frequency p-value out of range: %.6f", p)
		}
		if p, _ := RunsTest(data); p < 0 || p > 1 {
			t.Fatalf("runs p-value out of range: %.6f", p)
		}
		if p, _ := CumulativeSumsTest(data); p < 0 || p > 1 {
			t.Fatalf("cumulative sums p-value out of range: %.6f", p)
		}
		if p, _ := LongestRunOfOnesTest(data); p < 0 || p > 1 {
			t.Fatalf("longest run p-value out of range: %.6f", p)
		}
	})

	t.Run("large_data_patterns", func(t *testing.T) {
		data := make([]byte, 5000)
		for i := range data {
			// Mix of patterns
			if i%3 == 0 {
				data[i] = 0xAA
			} else if i%3 == 1 {
				data[i] = 0x55
			} else {
				data[i] = byte(i % 256)
			}
		}

		if p, _ := BinaryMatrixRankTest(data); p < 0 || p > 1 {
			t.Fatalf("binary matrix rank p-value out of range: %.6f", p)
		}
		if p, _ := DiscreteFourierTransformTest(data); p < 0 || p > 1 {
			t.Fatalf("DFT p-value out of range: %.6f", p)
		}
		if p, _ := NonOverlappingTemplateTest(data, 9); p < 0 || p > 1 {
			t.Fatalf("non-overlapping template p-value out of range: %.6f", p)
		}
		if p, _ := LinearComplexityTest(data, 500); p < 0 || p > 1 {
			t.Fatalf("linear complexity p-value out of range: %.6f", p)
		}
	})
}

// TestRunAllTestsVariousInputs tests RunAllTests with different inputs
func TestRunAllTestsVariousInputs(t *testing.T) {
	t.Run("minimum_size", func(t *testing.T) {
		// Exactly minimum size for Universal test
		data := make([]byte, 48480) // 387,840 bits
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
	})

	t.Run("alternating_pattern", func(t *testing.T) {
		data := make([]byte, 50000)
		for i := range data {
			if i%2 == 0 {
				data[i] = 0xAA
			} else {
				data[i] = 0x55
			}
		}

		results, _ := RunAllTests(data)
		if len(results) == 0 {
			t.Error("expected at least some results")
		}
	})

	t.Run("pseudo_random_pattern", func(t *testing.T) {
		data := make([]byte, 50000)
		for i := range data {
			data[i] = byte((i*i + i*7 + 13) % 256)
		}

		results, _ := RunAllTests(data)
		if len(results) == 0 {
			t.Error("expected at least some results")
		}
	})
}
