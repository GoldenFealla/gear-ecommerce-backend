package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/validation"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	RegisterUser(f *domain.RegisterUserForm) error
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

	group.GET("/test", handler.Test)
	group.GET("/check", handler.Check)
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
	group.GET("/logout", handler.Logout)
}

func (h *UserHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test user Ok")
}

func (h *UserHandler) Check(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	auth, ok := sess.Values["authenticated"].(bool)

	if auth && ok {
		return c.JSON(http.StatusOK, &domain.ResponseSuccess{
			Message: "already logged in",
		})
	}

	return c.JSON(http.StatusUnauthorized, &domain.ResponseError{
		Message: "you haven't logged in",
	})
}

func (h *UserHandler) Register(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	auth, ok := sess.Values["authenticated"].(bool)

	if auth && ok {
		return c.JSON(http.StatusOK, &domain.ResponseSuccess{
			Message: "already logged in",
		})
	}

	var body domain.RegisterUserForm
	err = c.Bind(&body)

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

	err = h.uc.RegisterUser(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteDefaultMode,
	}

	sess.Values["authenticated"] = true
	sess.Values["username"] = body.Username

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, "Registered User")
}

func (h *UserHandler) Login(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	auth, ok := sess.Values["authenticated"].(bool)

	if auth && ok {
		return c.JSON(http.StatusOK, &domain.ResponseSuccess{
			Message: "already logged in",
		})
	}

	var body domain.LoginUserForm
	err = c.Bind(&body)

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

	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteDefaultMode,
	}

	sess.Values["authenticated"] = true
	sess.Values["username"] = user.Username

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	auth, ok := sess.Values["authenticated"].(bool)

	if !auth || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.ResponseError{
			Message: "you haven't logged in",
		})
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteDefaultMode,
	}

	delete(sess.Values, "authenticated")
	delete(sess.Values, "username")

	err = sess.Save(c.Request(), c.Response())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
