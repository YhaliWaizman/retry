package command

import (
    "testing"
    "time"
    "github.com/yhaliwaizman/retry/internal/config"
    "github.com/yhaliwaizman/retry/internal/executor"
    "github.com/yhaliwaizman/retry/internal/logger"
)

func TestParseAttempts(t *testing.T) {
    // Create a runner instance for testing
    cfg := config.NewConfig()
    log := logger.NewLogger(false, false)
    exec := executor.NewExecutor()
    runner := NewRunner(cfg, log, exec)

    tests := []struct {
        name    string
        input   string
        want    int
        wantErr bool
    }{
        {"valid positive", "3", 3, false},
        {"valid large number", "100", 100, false},
        {"invalid zero", "0", 0, true},
        {"invalid negative", "-1", 0, true},
        {"invalid string", "abc", 0, true},
        {"invalid float", "3.5", 0, true},
        {"empty string", "", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := runner.parseAttempts(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("parseAttempts() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("parseAttempts() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestConfig(t *testing.T) {
    // Test config creation
    cfg := config.NewConfig()

    if cfg.Verbose {
        t.Error("Expected Verbose to be false by default")
    }
    if cfg.Quiet {
        t.Error("Expected Quiet to be false by default")
    }
    if cfg.Delay != time.Second {
        t.Errorf("Expected Delay to be 1s by default, got %v", cfg.Delay)
    }

    // Test setting config values
    cfg.Verbose = true
    cfg.Delay = 2 * time.Second
    cfg.CommandTimeout = 5 * time.Second

    if !cfg.Verbose {
        t.Error("Expected Verbose to be true after setting")
    }
    if cfg.Delay != 2*time.Second {
        t.Errorf("Expected Delay to be 2s, got %v", cfg.Delay)
    }
    if cfg.CommandTimeout != 5*time.Second {
        t.Errorf("Expected CommandTimeout to be 5s, got %v", cfg.CommandTimeout)
    }
}