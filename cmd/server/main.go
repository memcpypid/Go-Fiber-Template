package main

import (
	"fmt"
	"go-fiber-template/internal/config"
	"go-fiber-template/internal/delivery/http/handler"
	"go-fiber-template/internal/delivery/http/route"
	"go-fiber-template/internal/middleware"
	"go-fiber-template/internal/repository"
	"go-fiber-template/internal/service"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Configuration (Viper)
	v, err := config.NewViper(".")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize viper: %v", err))
	}

	cfg, err := config.NewConfig(v)
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// 2. Initialize Logger
	logger := config.NewLogger(cfg)
	defer logger.Sync()

	// 3. Initialize Database
	db := config.NewDatabase(cfg, logger)

	// 4. Initialize Validator
	validator := config.NewValidator()

	// 5. Initialize Repositories
	userRepo := repository.NewUserRepository(db, logger)
	tokenRepo := repository.NewRefreshTokenRepository(db, logger)

	// 6. Initialize Services
	userService := service.NewUserService(userRepo, logger)
	authService := service.NewAuthService(userRepo, tokenRepo, cfg, logger)

	// 7. Initialize Middlewares
	mid := middleware.NewMiddleware(cfg, logger)

	// 8. Initialize Handlers
	userHandler := handler.NewUserHandler(userService, validator)
	authHandler := handler.NewAuthHandler(authService, validator)

	// 9. Initialize Router & Get App
	r := route.NewRouter(mid, authHandler, userHandler)
	app := r.New()

	// 10. Start Server
	serverPort := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info("Starting Web Server", zap.String("port", serverPort))
	if err := app.Listen(serverPort); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
