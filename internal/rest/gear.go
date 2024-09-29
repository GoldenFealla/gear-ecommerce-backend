package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
	"github.com/labstack/echo/v4"
)

type GearUsecase interface {
	GetGearList() ([]*domain.Gear, error)
	GetGearByID(id string) (*domain.Gear, error)
	AddGear(g *domain.AddGearForm) error
	UpdateGear(g *domain.UpdateGearForm) error
	DeleteGear(id string) error
}

type GearHandler struct {
	uc GearUsecase
}

func NewGearHandler(e *echo.Echo, uc GearUsecase) {
	handler := &GearHandler{
		uc,
	}

	group := e.Group("gear")

	group.GET("/test", handler.Test)
	group.GET("/", handler.GetGearByID)
	group.GET("/list", handler.GetGearList)
	group.POST("/create", handler.AddGear)
	group.PUT("/update", handler.UpdateGear)
	group.DELETE("/delete", handler.DeleteGear)
}

func (h *GearHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test gear Ok")
}

func (h *GearHandler) GetGearList(c echo.Context) error {
	result, err := h.uc.GetGearList()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *GearHandler) GetGearByID(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.ResponseError{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	result, err := h.uc.GetGearByID(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *GearHandler) AddGear(c echo.Context) error {
	var body domain.AddGearForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, body)
}

func (h *GearHandler) UpdateGear(c echo.Context) error {
	var body domain.UpdateGearForm
	err := c.Bind(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, body)
}

func (h *GearHandler) DeleteGear(c echo.Context) error {
	if hasID := c.QueryParams().Has("id"); !hasID {
		return c.JSON(http.StatusBadRequest, &domain.ResponseError{
			Message: "query param 'id' is required",
		})
	}

	id := c.QueryParams().Get("id")

	err := h.uc.DeleteGear(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &domain.ResponseError{
		Message: "Successfully delete document",
	})
}
