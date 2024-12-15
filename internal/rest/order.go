package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/middleware"
)

type OrderUsecase interface {
	GetCart(ctx context.Context, userID string) (*domain.FullOrder, error)
	AddGearToCart(ctx context.Context, userID string, gearID string) error
	SetGearQuantityCart(ctx context.Context, orderID string, gearID string, quantity int64) error
	RemoveGearFromCart(ctx context.Context, userID string, gearID string) error
	PayCart(ctx context.Context, orderID string) error
	GetOrder(ctx context.Context, d string) (*domain.FullOrder, error)
	GetOrderList(ctx context.Context, orderID string, page int64, limit int64) ([]*domain.Order, error)
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
	group.PUT("/pay", handler.PayCart)
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

	ctx := c.Request().Context()
	cart, err := h.ou.GetCart(ctx, user.ID.String())

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

	ctx := c.Request().Context()
	err := h.ou.AddGearToCart(ctx, user.ID.String(), gearID)

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

	ctx := c.Request().Context()
	err = h.ou.SetGearQuantityCart(ctx, user.ID.String(), gearID, quantity)
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

	ctx := c.Request().Context()
	err := h.ou.RemoveGearFromCart(ctx, user.ID.String(), gearID)

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

func (h *OrderHandler) PayCart(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	orderID := c.QueryParam("id")

	ctx := c.Request().Context()
	err := h.ou.PayCart(ctx, orderID)
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

	pPage := c.QueryParams().Get("page")
	pLimit := c.QueryParams().Get("limit")

	var err error
	var page int64 = 1
	var limit int64 = 10

	if pPage != "" {
		page, err = strconv.ParseInt(pPage, 10, 64)

		if err != nil {
			return c.JSON(http.StatusBadRequest, &domain.Response{
				Message: err.Error(),
			})
		}
	}

	if pLimit != "" {
		limit, err = strconv.ParseInt(pLimit, 10, 64)

		if err != nil {
			return c.JSON(http.StatusBadRequest, &domain.Response{
				Message: err.Error(),
			})
		}
	}

	ctx := c.Request().Context()
	result, err := h.ou.GetOrderList(ctx, user.ID.String(), page, limit)

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

	ctx := c.Request().Context()
	result, err := h.ou.GetOrder(ctx, id)

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
