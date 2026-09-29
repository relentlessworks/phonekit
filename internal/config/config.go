package config

import (
	"flag"
	"os"
)

// Config holds all configuration for the phonekit service.
type Config struct {
	Addr   string
	Secret string
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Addr:   ":8080",
		Secret: "",
	}
}

// Load parses config from defaults < env < flags.
func Load() *Config {
	c := Default()

	// Environment variables
	if v := os.Getenv("PHONEKIT_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("PHONEKIT_SECRET"); v != "" {
		c.Secret = v
	}

	// Flags (override env)
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.StringVar(&c.Secret, "secret", c.Secret, "auth token secret (auto-generated if empty)")
	flag.Parse()

	return c
}
