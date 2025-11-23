package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsRegistration(t *testing.T) {
	// Ensure the collectors can be used without panic and are registered.
	if _, err := TestsTotal.GetMetricWithLabelValues("frequency", "pass"); err != nil {
		t.Fatalf("TestsTotal missing labels: %v", err)
	}
	if _, err := TestDuration.GetMetricWithLabelValues("frequency"); err != nil {
		t.Fatalf("TestDuration missing labels: %v", err)
	}
	if _, err := PValue.GetMetricWithLabelValues("frequency"); err != nil {
		t.Fatalf("PValue missing labels: %v", err)
	}
	if _, err := RequestsTotal.GetMetricWithLabelValues("RunTests", "success"); err != nil {
		t.Fatalf("RequestsTotal missing labels: %v", err)
	}

	// Gather to assert metrics exist.
	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}
	required := map[string]bool{
		"nist_tests_total":              false,
		"nist_test_duration_seconds":    false,
		"nist_overall_duration_seconds": false,
		"nist_last_overall_pass_rate":   false,
		"nist_p_value":                  false,
		"nist_requests_total":           false,
	}
	for _, mf := range mfs {
		if _, ok := required[mf.GetName()]; ok {
			required[mf.GetName()] = true
		}
	}
	for name, seen := range required {
		if !seen {
			t.Fatalf("expected metric %s to be registered", name)
		}
	}
}

func TestMetricWrappers(t *testing.T) {
	// Test wrapper functions to ensure they don't panic and record something
	RecordTestDuration("test_test", 1.0)
	IncrementTestsTotal("test_test", "pass")
	RecordPValue("test_test", 0.5)
	IncrementRequestsTotal("TestRPC", "ok")

	// We can't easily check the exact values without more complex setup,
	// but running them ensures coverage and no panics.
}
