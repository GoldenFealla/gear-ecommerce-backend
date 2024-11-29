package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
	"github.com/labstack/echo/v4"
)

type GearUsecase interface {
	GetGearBrandList(category string) ([]string, error)
	GetGearList(filter domain.ListGearFilter) ([]*domain.Gear, error)
	GetGearByID(id string) (*domain.Gear, error)
	AddGear(g *domain.AddGearForm) error
	UpdateGear(id string, g *domain.UpdateGearForm) error
	DeleteGear(id string) error
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
	group.GET("/list-brand", handler.GetGearBrandList)
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

	result, err := h.uc.GetGearBrandList(category)

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

	result, err := h.uc.GetGearList(filter)

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

	result, err := h.uc.GetGearByID(id)

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

	err = h.v.Struct(body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	err = h.uc.AddGear(&body)

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

	err = h.uc.UpdateGear(id, &body)

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

	err := h.uc.DeleteGear(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.Response{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.Response{
		Message: "Successfully delete document",
	})
}
