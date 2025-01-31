package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/final_project1240930/api_gateway/internal/logs"
	services "gitlab.com/final_project1240930/api_gateway/internal/services/user_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) *userHandler {
	return &userHandler{userSrv: userSrv}
}

func createErrorResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func (h *userHandler) Register(c echo.Context) error {
	var req services.UserRequest

	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format during registration", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if req.Username == "" || req.Password == "" {
		logs.Error("Username and password are required during registration")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("username and password are required")))
	}

	resp, err := h.userSrv.Register(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to register user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) Login(c echo.Context) error {
	var req services.LoginRequest

	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format during login", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if req.Username == "" || req.Password == "" {
		logs.Error("Username and password are required during login")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("username and password are required")))
	}

	resp, err := h.userSrv.Login(c.Request().Context(), &req)
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			logs.Error("Invalid credentials provided", zap.String("username", req.Username))
			return c.JSON(http.StatusUnauthorized, createErrorResponse(errors.New("invalid credentials")))
		} else if status.Code(err) == codes.NotFound {
			logs.Error("User not found", zap.String("username", req.Username))
			return c.JSON(http.StatusNotFound, createErrorResponse(errors.New("user not found")))
		}
		logs.Error("Failed to login user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": resp.AccessToken,
	})
}

func (h *userHandler) GetAllUsers(c echo.Context) error {
	resp, err := h.userSrv.GetAllUsers(c.Request().Context())
	if err != nil {
		logs.Error("Failed to get all users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	req := services.UserIdentifier{Identifier: id}

	resp, err := h.userSrv.GetUser(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to get user", zap.String("userID", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateUserRole(c echo.Context) error {
	id := c.Param("id")
	var req services.UpdateRoleRequest

	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format during role update", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if req.Role == "" {
		logs.Error("Role is required during role update")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("role is required")))
	}

	req.Identifier = id
	resp, err := h.userSrv.UpdateUserRole(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update user role", zap.String("userID", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) UpdateUserPassword(c echo.Context) error {
	id := c.Param("id")
	var req services.UpdatePasswordRequest

	if err := c.Bind(&req); err != nil {
		logs.Error("Invalid request format during password update", zap.Error(err))
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("invalid request format")))
	}

	if len(req.Password) < 8 {
		logs.Error("Password must be at least 8 characters long")
		return c.JSON(http.StatusBadRequest, createErrorResponse(errors.New("password must be at least 8 characters long")))
	}

	req.Identifier = id
	resp, err := h.userSrv.UpdateUserPassword(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to update user password", zap.String("userID", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	req := services.UserIdentifier{Identifier: id}

	resp, err := h.userSrv.DeleteUser(c.Request().Context(), &req)
	if err != nil {
		logs.Error("Failed to delete user", zap.String("userID", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, createErrorResponse(err))
	}

	return c.JSON(http.StatusOK, resp)
}
