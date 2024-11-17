package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/middleware"
)

type AddressUsecase interface {
	GetAddressList(userID string) ([]*domain.Address, error)
	GetAddressByID(id string) (*domain.Address, error)
	AddAddress(userID string, g *domain.AddAddressForm) error
	UpdateAddress(id string, g *domain.UpdateAddressForm) error
	DeleteAddress(id string) error
}

type AddressHandler struct {
	ac AddressUsecase
	v  *validator.Validate
}

func NewAddressHandler(e *echo.Echo, ac AddressUsecase, v *validator.Validate) {
	handler := &AddressHandler{
		ac,
		v,
	}

	group := e.Group("address")
	group.Use(middleware.AuthenticatedWithConfig(&middleware.AuthenticatedConfig{
		Excludes: []string{},
	}))

	group.GET("/test", handler.Test)

	group.GET("", handler.GetAddress)
	group.GET("/list", handler.GetListAddress)
	group.POST("/add", handler.AddAddress)
	group.PUT("/update", handler.UpdateAddress)
	group.PUT("/delete", handler.DeleteAddress)
}

func (h *AddressHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test address Ok")
}

func (h *AddressHandler) GetAddress(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	result, err := h.ac.GetAddressByID(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "OK",
		Data:    result,
	})
}

func (h *AddressHandler) GetListAddress(c echo.Context) error {
	if hasID := c.QueryParams().Has("user_id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'user_id' is required",
		})
	}

	userID := c.QueryParams().Get("user_id")

	result, err := h.ac.GetAddressList(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "OK",
		Data:    result,
	})
}

func (h *AddressHandler) AddAddress(c echo.Context) error {
	u, ok := c.Get("user").(*domain.UserInfo)

	if u == nil || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.Response{
			Message: "You need to login to add",
			Data:    nil,
		})
	}

	var body domain.AddAddressForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.ac.AddAddress(u.ID.String(), &body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Created address",
	})
}

func (h *AddressHandler) UpdateAddress(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	var body domain.UpdateAddressForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.ac.UpdateAddress(id, &body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Updated address",
	})
}

func (h *AddressHandler) DeleteAddress(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	err := h.ac.DeleteAddress(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "Successfully delete address",
	})
}