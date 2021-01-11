package config

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadEnv loads environment variables into env
func LoadEnv(env interface{}) error {
	err := godotenv.Load()
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok || !errors.Is(perr.Unwrap(), os.ErrNotExist) {
			return err
		}
	}

	return envconfig.Process("", env)
}

// WithCancel replicates context.WithCancel but listens for Interrupt and SIGTERM signals
func WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		defer cancel()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
	}()

	return ctx, cancel
}
