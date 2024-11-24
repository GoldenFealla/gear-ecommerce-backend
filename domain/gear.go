package domain

import "github.com/google/uuid"

var GearTypeMap map[string]string = map[string]string{
	"PC":        "PERSONAL_COMPUTER",
	"Laptop":    "LAPTOP",
	"MainBoard": "MAINBOARD",
	"CPU":       "CENTRAL_PROCESSING_UNIT",
	"GPU":       "GRAPHICS_PROCESSING_UNIT",
	"PSU":       "POWER_SUPPLY_UNIT",
	"RAM":       "RANDOM_ACCESS_MEMORY",
	"Fan":       "FAN",
	"Storage":   "STORAGE",
	"Monitor":   "MONITOR",
}

type Gear struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Type     string    `json:"type" db:"type"`
	Price    float64   `json:"price" db:"price"`
	Discount float64   `json:"discount" db:"discount"`
	Quantity int64     `json:"quantity" db:"quantity"`
	ImageURL string    `json:"image_url" db:"image_url"`
}

type AddGearForm struct {
	Name        string  `json:"name,omitempty" validate:"required"`
	Type        string  `json:"type,omitempty" validate:"required,is-gear"`
	Price       float64 `json:"price,omitempty"`
	Discount    float64 `json:"discount,omitempty"`
	Quantity    int64   `json:"quantity,omitempty"`
	ImageBase64 *string `json:"image_base64,omitempty"`
}

type UpdateGearForm struct {
	Name        *string  `json:"name,omitempty" db:"name" validate:"omitempty"`
	Type        *string  `json:"type,omitempty" db:"type" validate:"omitempty,is-gear"`
	Price       *float64 `json:"price,omitempty" db:"price" validate:"omitempty"`
	Discount    *float64 `json:"discount,omitempty" db:"discount" validate:"omitempty"`
	Quantity    *int64   `json:"quantity,omitempty" db:"quantity" validate:"omitempty"`
	ImageBase64 string   `json:"image_base64,omitempty" db:"quantity" validate:"omitempty"`
}
