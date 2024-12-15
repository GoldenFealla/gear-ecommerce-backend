package rest

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
	"github.com/labstack/echo/v4"
	"github.com/leebenson/conform"
)

type GearUsecase interface {
	GetGearVarietyList(ctx context.Context, category string) ([]string, error)
	GetGearBrandList(ctx context.Context, category string) ([]string, error)
	GetGearListCount(ctx context.Context, filter domain.ListGearFilter) (int64, error)
	GetGearList(ctx context.Context, filter domain.ListGearFilter) ([]*domain.Gear, error)
	GetGearByID(ctx context.Context, id string) (*domain.Gear, error)
	AddGear(ctx context.Context, g *domain.AddGearForm) error
	UpdateGear(ctx context.Context, id string, g *domain.UpdateGearForm) error
	DeleteGear(ctx context.Context, id string) error
}

type GearHandler struct {
	uc GearUsecase
	v  *validator.Validate
}

func NewGearHandler(e *echo.Echo, uc GearUsecase, v *validator.Validate) {
	handler := &GearHandler{
		uc,
		v,
	}

	group := e.Group("gear")

	group.GET("/test", handler.Test)
	group.GET("/", handler.GetGearByID)
	group.GET("/list-count", handler.GetGearListCount)
	group.GET("/list-brand", handler.GetGearBrandList)
	group.GET("/list-variety", handler.GetGearVarietyList)
	group.GET("/list", handler.GetGearList)
	group.POST("/create", handler.AddGear)
	group.PUT("/update", handler.UpdateGear)
	group.DELETE("/delete", handler.DeleteGear)
}

func (h *GearHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test gear Ok")
}

func (h *GearHandler) GetGearBrandList(c echo.Context) error {
	if hasCategory := c.QueryParams().Has("category"); !hasCategory {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'category' is required",
		})
	}

	category := c.QueryParams().Get("category")

	ctx := c.Request().Context()
	result, err := h.uc.GetGearBrandList(ctx, category)

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

func (h *GearHandler) GetGearVarietyList(c echo.Context) error {
	if hasCategory := c.QueryParams().Has("category"); !hasCategory {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'category' is required",
		})
	}

	category := c.QueryParams().Get("category")

	ctx := c.Request().Context()
	result, err := h.uc.GetGearVarietyList(ctx, category)

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

func (h *GearHandler) GetGearListCount(c echo.Context) error {
	if hasCategory := c.QueryParams().Has("category"); !hasCategory {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'category' is required",
		})
	}

	filter := domain.ListGearFilter{}

	err := c.Bind(&filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()
	result, err := h.uc.GetGearListCount(ctx, filter)

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

func (h *GearHandler) GetGearList(c echo.Context) error {
	if hasCategory := c.QueryParams().Has("category"); !hasCategory {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'category' is required",
		})
	}

	defaultPage := int64(1)
	defaultLimit := int64(10)

	filter := domain.ListGearFilter{
		Page:  &defaultPage,
		Limit: &defaultLimit,
	}

	err := c.Bind(&filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()
	result, err := h.uc.GetGearList(ctx, filter)

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

func (h *GearHandler) GetGearByID(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	ctx := c.Request().Context()
	result, err := h.uc.GetGearByID(ctx, id)

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

func (h *GearHandler) AddGear(c echo.Context) error {
	var body domain.AddGearForm
	err := c.Bind(&body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	err = conform.Strings(&body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.v.Struct(body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()
	err = h.uc.AddGear(ctx, &body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Created Gear",
	})
}

func (h *GearHandler) UpdateGear(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	var body domain.UpdateGearForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	ctx := c.Request().Context()
	err = h.uc.UpdateGear(ctx, id, &body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &domain.Response{
		Message: "Updated gear",
	})
}

func (h *GearHandler) DeleteGear(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.Response{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	ctx := c.Request().Context()
	err := h.uc.DeleteGear(ctx, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "Successfully delete document",
	})
}
