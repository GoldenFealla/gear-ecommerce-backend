package rest

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/middleware"
)

type OrderUsecase interface {
	GetCart(userID string) (*domain.FullOrder, error)
	AddGearToCart(userID string, gearID string) error
	SetGearQuantityCart(orderID string, gearID string, quantity int64) error
	RemoveGearFromCart(userID string, gearID string) error
	GetOrder(id string) (*domain.FullOrder, error)
	GetOrderList(userID string) ([]*domain.FullOrder, error)
}

type OrderHandler struct {
	ou OrderUsecase
	v  *validator.Validate
}

func NewOrderHandler(e *echo.Echo, ou OrderUsecase, v *validator.Validate) {
	handler := &OrderHandler{
		ou,
		v,
	}

	group := e.Group("order")
	group.Use(middleware.AuthenticatedWithConfig(&middleware.AuthenticatedConfig{
		Excludes: []string{},
	}))

	group.GET("/test", handler.Test)

	group.GET("", handler.GetOrder)
	group.GET("/list", handler.GetOrderList)
	group.GET("/cart", handler.GetCart)
	group.PUT("/add-to-cart", handler.AddGearToCart)
	group.PUT("/set-quantity", handler.SetGearQuantityCart)
	group.PUT("/remove-from-cart", handler.RemoveGearFromCart)
}

func (h *OrderHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test order Ok")
}

func (h *OrderHandler) GetCart(c echo.Context) error {
	user, ok := c.Get("user").(*domain.UserInfo)

	if user == nil || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.Response{
			Message: "You need to login",
			Data:    nil,
		})
	}

	cart, err := h.ou.GetCart(user.ID.String())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "OK",
		Data:    cart,
	})
}

func (h *OrderHandler) AddGearToCart(c echo.Context) error {
	user, ok := c.Get("user").(*domain.UserInfo)

	if user == nil || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.Response{
			Message: "You need to login",
			Data:    nil,
		})
	}

	if hasID := c.QueryParams().Has("gear_id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'gear_id' is required",
		})
	}

	gearID := c.QueryParam("gear_id")

	err := h.ou.AddGearToCart(user.ID.String(), gearID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "OK",
		Data:    nil,
	})
}

func (h *OrderHandler) SetGearQuantityCart(c echo.Context) error {
	user, ok := c.Get("user").(*domain.UserInfo)

	if user == nil || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.Response{
			Message: "You need to login",
			Data:    nil,
		})
	}

	if hasID := c.QueryParams().Has("gear_id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'gear_id' is required",
		})
	}

	gearID := c.QueryParam("gear_id")

	if hasID := c.QueryParams().Has("quantity"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'quantity' is required",
		})
	}

	quantity, err := strconv.ParseInt(c.QueryParam("quantity"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.ou.SetGearQuantityCart(user.ID.String(), gearID, quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "OK",
		Data:    nil,
	})
}

func (h *OrderHandler) RemoveGearFromCart(c echo.Context) error {
	user, ok := c.Get("user").(*domain.UserInfo)

	if user == nil || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.Response{
			Message: "You need to login",
			Data:    nil,
		})
	}

	if hasID := c.QueryParams().Has("gear_id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'gear_id' is required",
		})
	}

	gearID := c.QueryParam("gear_id")

	err := h.ou.RemoveGearFromCart(user.ID.String(), gearID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "OK",
		Data:    nil,
	})
}

func (h *OrderHandler) GetOrderList(c echo.Context) error {
	user, ok := c.Get("user").(*domain.UserInfo)

	if user == nil || !ok {
		return c.JSON(http.StatusUnauthorized, &domain.Response{
			Message: "You need to login",
			Data:    nil,
		})
	}

	result, err := h.ou.GetOrder(user.ID.String())

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

func (h *OrderHandler) GetOrder(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	result, err := h.ou.GetOrder(id)

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
