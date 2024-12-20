package rest

import (
	"context"
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
	RegisterUser(ctx context.Context, f *domain.RegisterUserForm) (*domain.UserInfo, error)
	LoginUser(ctx context.Context, f *domain.LoginUserForm) (*domain.UserInfo, error)
	UpdateUser(ctx context.Context, id string, f *domain.UpdateUserForm) (*domain.UserInfo, error)
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
	group.Use(middleware.AuthenticatedWithConfig(&middleware.AuthenticatedConfig{
		Excludes: []string{
			"/user/test",
			"/user/login",
			"/user/register",
		},
	}))

	group.GET("/test", handler.Test)

	group.GET("/check", handler.Check)
	group.GET("/refresh", handler.Refresh)

	group.POST("/register", handler.Register)
	group.PUT("/update", handler.Update)

	group.POST("/login", handler.Login)
	group.GET("/logout", handler.Logout)
}

func (h *UserHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test user Ok")
}

func (h *UserHandler) Check(c echo.Context) error {
	u, ok := c.Get("user").(*domain.UserInfo)

	if u != nil && ok {
		return c.JSON(http.StatusOK, &domain.Response{
			Message: "OK",
			Data:    u,
		})
	}

	return c.JSON(http.StatusUnauthorized, &domain.Response{
		Message: "Not log in",
		Data:    nil,
	})
}

func (h *UserHandler) Refresh(c echo.Context) error {
	user := c.Get("user").(*domain.UserInfo)

	refreshToken, err := jwt.GenerateRefreshToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	c.Response().Before(func() {
		session.DefaultSaveSession(c, &refreshToken)
	})

	accessToken, err := jwt.GenerateAccessToken(refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Registered User",
		Data: &domain.UserCredential{
			Token:    accessToken,
			UserInfo: user,
		},
	})
}

func (h *UserHandler) Register(c echo.Context) error {
	var body domain.RegisterUserForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		ves := validation.GetValidationError(err.(validator.ValidationErrors))
		return c.JSON(http.StatusBadRequest, ves)
	}

	ctx := c.Request().Context()
	info, err := h.uc.RegisterUser(ctx, &body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	refreshToken, err := jwt.GenerateRefreshToken(info)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	c.Response().Before(func() {
		session.DefaultSaveSession(c, &refreshToken)
	})

	accessToken, err := jwt.GenerateAccessToken(refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Registered User",
		Data: &domain.UserCredential{
			Token:    accessToken,
			UserInfo: info,
		},
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	u, ok := c.Get("user").(*domain.UserInfo)

	if u != nil && ok {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "You already logged in",
		})
	}

	var body domain.LoginUserForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		ves := validation.GetValidationError(err.(validator.ValidationErrors))
		return c.JSON(http.StatusBadRequest, ves)
	}

	ctx := c.Request().Context()
	user, err := h.uc.LoginUser(ctx, &body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	c.Response().Before(func() {
		session.DefaultSaveSession(c, &refreshToken)
	})

	accessToken, err := jwt.GenerateAccessToken(refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Logged in User",
		Data: &domain.UserCredential{
			Token:    accessToken,
			UserInfo: user,
		},
	})
}

func (h *UserHandler) Update(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	var body domain.UpdateUserForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		ves := validation.GetValidationError(err.(validator.ValidationErrors))
		return c.JSON(http.StatusBadRequest, ves)
	}

	ctx := c.Request().Context()
	info, err := h.uc.UpdateUser(ctx, id, &body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "Registered User",
		Data:    info,
	})
}

func (h *UserHandler) Logout(c echo.Context) error {
	err := session.DeleteSession(c)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.NoContent(http.StatusOK)
}
