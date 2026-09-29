package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/relentlessworks/phonekit/internal/api"
	"github.com/relentlessworks/phonekit/internal/config"
)

func main() {
	cfg := config.Load()

	// Generate a random secret if none provided
	if cfg.Secret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			log.Fatalf("failed to generate secret: %v", err)
		}
		cfg.Secret = hex.EncodeToString(b)
	}

	handler := api.NewHandler(cfg.Secret)

	mux := http.NewServeMux()
	handler.Register(mux)

	log.Printf("phonekit starting on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func init() {
	// Suppress unused import warning
	_ = fmt.Sprintf
}
