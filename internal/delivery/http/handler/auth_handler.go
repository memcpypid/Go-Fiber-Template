package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-fiber-template/internal/dto"
	"go-fiber-template/internal/service"
	"go-fiber-template/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService service.AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request payload", response.ErrorDetail{Field: "body", Message: err.Error()}))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err, nil))
	}

	resp, err := h.authService.Login(c.UserContext(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, err.Error()))
	}

	return c.JSON(response.Success("Login successful", resp))
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request payload", response.ErrorDetail{Field: "body", Message: err.Error()}))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err, nil))
	}

	resp, err := h.authService.Register(c.UserContext(), &req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(response.Error(fiber.StatusConflict, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success("Registration successful", resp))
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request payload", response.ErrorDetail{Field: "body", Message: err.Error()}))
	}

	resp, err := h.authService.RefreshToken(c.UserContext(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, err.Error()))
	}

	return c.JSON(response.Success("Token refreshed", resp))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request payload", response.ErrorDetail{Field: "body", Message: err.Error()}))
	}

	err := h.authService.Logout(c.UserContext(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(response.Success("Logout successful", nil))
}
