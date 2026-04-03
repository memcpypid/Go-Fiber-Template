package handler

import (
	"go-fiber-template/internal/dto"
	"go-fiber-template/internal/service"
	"go-fiber-template/internal/utils"
	"go-fiber-template/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService service.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// GetUsers handles fetching all users with pagination
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	pagination := utils.GeneratePaginationFromRequest(c)

	users, total, err := h.userService.GetAll(c.UserContext(), pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, "Failed to fetch users", response.ErrorDetail{Field: "server", Message: err.Error()}))
	}

	return c.JSON(response.SuccessWithPagination("Users retrieved successfully", users, total, pagination.Limit, pagination.Page))
}
func (h *UserHandler) GetUserCount(c *fiber.Ctx) error {
	count, err := h.userService.GetUserCount(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, "Failed to fetch user count", response.ErrorDetail{Field: "server", Message: err.Error()}))
	}
	return c.JSON(response.Success("User count retrieved successfully", count))
}

// GetProfile retrieves the profile of the currently logged-in user
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	idStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, "Unauthorized"))
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, "Invalid ID format in token"))
	}

	user, err := h.userService.GetByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(fiber.StatusNotFound, "User not found"))
	}

	return c.JSON(response.Success("Profile retrieved successfully", user))
}

// UpdateProfile updates the profile of the currently logged-in user
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	idStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, "Unauthorized"))
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, "Invalid ID format in token"))
	}

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request payload", response.ErrorDetail{Field: "body", Message: err.Error()}))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err, nil))
	}

	// Map UpdateProfileRequest to UpdateUserRequest for service reuse
	updateReq := &dto.UpdateUserRequest{
		Name:     req.Name,
		Password: req.Password,
	}

	user, err := h.userService.Update(c.UserContext(), id, updateReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(response.Success("Profile updated successfully", user))
}

// UpdateUser handles user updates by ID (Admin only)
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request payload", response.ErrorDetail{Field: "body", Message: err.Error()}))
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ValidationError(err, nil))
	}

	user, err := h.userService.Update(c.UserContext(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(response.Success("User updated successfully", user))
}

// DeleteUser handles user deletion by ID (Admin only)
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid ID format"))
	}

	err = h.userService.Delete(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, "Failed to delete user", response.ErrorDetail{Field: "server", Message: err.Error()}))
	}

	return c.JSON(response.Success("User deleted successfully", nil))
}

// ActivateAccount activates a user account (Admin only)
func (h *UserHandler) ActivateAccount(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid ID format"))
	}

	if err := h.userService.ActivateAccount(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(response.Success("Account activated successfully", nil))
}

// DeactivateAccount deactivates a user account (Admin only)
func (h *UserHandler) DeactivateAccount(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid ID format"))
	}

	if err := h.userService.DeactivateAccount(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, err.Error()))
	}

	return c.JSON(response.Success("Account deactivated successfully", nil))
}
