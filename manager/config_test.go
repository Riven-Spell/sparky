package manager

import (
	"testing"
)

func TestNewDefaultConfig_Valid(t *testing.T) {
	cfg := NewDefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config validation failed: %v", err)
	}
}
