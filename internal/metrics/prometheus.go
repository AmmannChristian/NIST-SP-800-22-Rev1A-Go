package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TestsTotal counts the total number of individual tests run
	TestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nist_tests_total",
			Help: "Total number of NIST statistical tests run",
		},
		[]string{"test", "status"},
	)

	// TestDuration tracks the duration of individual tests
	TestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "nist_test_duration_seconds",
			Help:    "Duration of individual NIST tests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"test"},
	)

	// OverallDuration tracks the duration of the entire test suite
	OverallDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "nist_overall_duration_seconds",
			Help:    "Duration of the entire NIST test suite in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// LastOverallPassRate stores the last overall pass rate
	LastOverallPassRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "nist_last_overall_pass_rate",
			Help: "Last overall pass rate of NIST tests (0.0-1.0)",
		},
	)

	// PValue stores the p-value for each test
	PValue = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nist_p_value",
			Help: "P-value of individual NIST tests",
		},
		[]string{"test"},
	)

	// RequestsTotal counts total gRPC requests
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nist_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)
)

// RecordTestDuration records the duration of a test
func RecordTestDuration(testName string, durationSeconds float64) {
	TestDuration.WithLabelValues(testName).Observe(durationSeconds)
}

// IncrementTestsTotal increments the total tests counter
func IncrementTestsTotal(testName, status string) {
	TestsTotal.WithLabelValues(testName, status).Inc()
}

// RecordPValue records the p-value of a test
func RecordPValue(testName string, pValue float64) {
	PValue.WithLabelValues(testName).Set(pValue)
}

// IncrementRequestsTotal increments the total requests counter
func IncrementRequestsTotal(method, status string) {
	RequestsTotal.WithLabelValues(method, status).Inc()
}
