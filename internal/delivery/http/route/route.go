package route

import (
	"go-fiber-template/internal/delivery/http/handler"
	"go-fiber-template/internal/middleware"
	"go-fiber-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	mw          *middleware.Middleware
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
}

func NewRouter(
	mw *middleware.Middleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *Router {
	return &Router{
		mw:          mw,
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

func (r *Router) New() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:           "Go Fiber Clean Architecture",
		EnablePrintRoutes: true,
	})

	// Global middleware
	if r.mw != nil {
		app.Use(r.mw.CORSMiddleware())
		app.Use(r.mw.LoggingMiddleware())
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(response.Success("server is healthy", nil))
	})

	api := app.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	{
		auth.Post("/register", r.authHandler.Register)
		auth.Post("/login", r.authHandler.Login)
		auth.Post("/refresh", r.authHandler.RefreshToken)
		auth.Post("/logout", r.authHandler.Logout)
	}

	// Protected Routes
	protected := api.Group("/")
	protected.Use(r.mw.AuthMiddleware())
	{
		protected.Get("/users/me", r.userHandler.GetProfile)
		protected.Put("/users/me", r.userHandler.UpdateProfile)

		// Admin Only Routes
		admin := protected.Group("/users")
		admin.Use(r.mw.RoleMiddleware("admin"))
		{
			admin.Get("", r.userHandler.GetUsers)
			admin.Get("/count", r.userHandler.GetUserCount)
			admin.Put("/:id", r.userHandler.UpdateUser)
			admin.Delete("/:id", r.userHandler.DeleteUser)
			admin.Patch("/:id/activate", r.userHandler.ActivateAccount)
			admin.Patch("/:id/deactivate", r.userHandler.DeactivateAccount)
		}
	}

	return app
}
