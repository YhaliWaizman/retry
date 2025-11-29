package config

import "time"

// Config holds all command-line flags
type Config struct {
    Verbose        bool
    Delay          time.Duration
    Quiet          bool
    Timeout        time.Duration
    CommandTimeout time.Duration
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
    return &Config{
        Delay: time.Second,
    }
}