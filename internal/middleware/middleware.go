package middleware

import (
	"go-fiber-template/internal/config"
	"go.uber.org/zap"
)

// Middleware struct holds the dependencies for all middlewares.
// Methods attached to this struct act as standard Fiber handlers.
type Middleware struct {
	cfg    *config.Config
	logger *zap.Logger
}

// NewMiddleware constructor for DI
func NewMiddleware(cfg *config.Config, logger *zap.Logger) *Middleware {
	return &Middleware{
		cfg:    cfg,
		logger: logger,
	}
}
