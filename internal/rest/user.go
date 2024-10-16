package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/jwt"
	"github.com/goldenfealla/gear-manager/internal/middleware"
	"github.com/goldenfealla/gear-manager/internal/session"
	"github.com/goldenfealla/gear-manager/internal/validation"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	RegisterUser(f *domain.RegisterUserForm) (*domain.UserInfo, error)
	LoginUser(f *domain.LoginUserForm) (*domain.UserInfo, error)
}

type UserHandler struct {
	uc UserUsecase
	v  *validator.Validate
}

func NewUserHandler(e *echo.Echo, uc UserUsecase, v *validator.Validate) {
	handler := &UserHandler{
		uc,
		v,
	}

	group := e.Group("user")
	group.Use(middleware.IsAuth)
	group.GET("/refresh", handler.Refresh)
	group.GET("/check", handler.Check)
	group.GET("/logout", handler.Logout)
	group.GET("/test", handler.Test)
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
}

func (h *UserHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test user Ok")
}

func (h *UserHandler) Check(c echo.Context) error {
	_, ok := c.Get("user").(*domain.UserInfo)

	if !ok {
		return c.JSON(http.StatusUnauthorized, false)
	}

	return c.JSON(http.StatusOK, true)
}

func (h *UserHandler) Refresh(c echo.Context) error {
	user := c.Get("user").(*domain.UserInfo)

	refreshToken, err := jwt.GenerateRefreshToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	c.Response().Before(func() {
		session.DefaultSaveSession(c, &refreshToken)
	})

	accessToken, err := jwt.GenerateAccessToken(refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.UserCredential{
		Message: "Registered User",
		Token:   accessToken,
	})
}

func (h *UserHandler) Register(c echo.Context) error {
	var body domain.RegisterUserForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		ves := validation.GetValidationError(err.(validator.ValidationErrors))
		return c.JSON(http.StatusBadRequest, ves)
	}

	info, err := h.uc.RegisterUser(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	refreshToken, err := jwt.GenerateRefreshToken(info)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	c.Response().Before(func() {
		session.DefaultSaveSession(c, &refreshToken)
	})

	accessToken, err := jwt.GenerateAccessToken(refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.UserCredential{
		Message: "Registered User",
		Token:   accessToken,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var body domain.LoginUserForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		ves := validation.GetValidationError(err.(validator.ValidationErrors))
		return c.JSON(http.StatusBadRequest, ves)
	}

	user, err := h.uc.LoginUser(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	c.Response().Before(func() {
		session.DefaultSaveSession(c, &refreshToken)
	})

	accessToken, err := jwt.GenerateAccessToken(refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.UserCredential{
		Message: "Logged in User",
		Token:   accessToken,
	})
}

func (h *UserHandler) Logout(c echo.Context) error {
	err := session.DeleteSession(c)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
