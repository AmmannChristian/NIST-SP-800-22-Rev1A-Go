package config

import (
	"testing"
)

func TestLoadWithEnvOverrides(t *testing.T) {
	t.Setenv("GRPC_PORT", "5000")
	t.Setenv("METRICS_PORT", "6000")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.GRPCPort != 5000 || cfg.MetricsPort != 6000 {
		t.Fatalf("unexpected ports: %+v", cfg)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("unexpected log level: %s", cfg.LogLevel)
	}
}

func TestValidateFailures(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{"bad grpc port", Config{GRPCPort: 0, MetricsPort: 9000, LogLevel: "info"}},
		{"bad metrics port", Config{GRPCPort: 9000, MetricsPort: 70000, LogLevel: "info"}},
		{"bad log level", Config{GRPCPort: 9000, MetricsPort: 9001, LogLevel: "verbose"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); err == nil {
				t.Fatalf("expected error for %s", tt.name)
			}
		})
	}

	// getEnvInt falls back on parse error
	t.Setenv("SOME_INT", "notanint")
	if v := getEnvInt("SOME_INT", 42); v != 42 {
		t.Fatalf("expected default on parse error, got %d", v)
	}
}
