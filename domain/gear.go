package domain

import "github.com/google/uuid"

var GearTypeMap map[string]string = map[string]string{
	"all":       "ALL",
	"pc":        "PERSONAL_COMPUTER",
	"laptop":    "LAPTOP",
	"mainboard": "MAINBOARD",
	"cpu":       "CENTRAL_PROCESSING_UNIT",
	"gpu":       "GRAPHICS_PROCESSING_UNIT",
	"psu":       "POWER_SUPPLY_UNIT",
	"ram":       "RANDOM_ACCESS_MEMORY",
	"fan":       "FAN",
	"storage":   "STORAGE",
	"monitor":   "MONITOR",
}

type Gear struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Type     string    `json:"type" db:"type"`
	Brand    string    `json:"brand" db:"brand"`
	Variety  string    `json:"variety" db:"variety"`
	Price    float64   `json:"price" db:"price"`
	Discount float64   `json:"discount" db:"discount"`
	Quantity int64     `json:"quantity" db:"quantity"`
	ImageURL string    `json:"image_url" db:"image_url"`
}

type ListGearFilter struct {
	Page       *int64  `query:"page"`
	Limit      *int64  `query:"limit"`
	Category   *string `query:"category"`
	Brand      *string `query:"brand"`
	Variety    *string `query:"variety"`
	StartPrice *int64  `query:"start_price"`
	EndPrice   *int64  `query:"end_price"`
	Sort       *string `query:"sort"`
}

type AddGearForm struct {
	Name        string  `json:"name,omitempty"          conform:"trim" validate:"required"`
	Type        string  `json:"type,omitempty"          conform:"trim" validate:"required,is-gear"`
	Brand       string  `json:"brand"                   conform:"trim" validate:"required"`
	Variety     string  `json:"variety"                 conform:"trim" validate:"required"`
	Price       float64 `json:"price,omitempty"         conform:"trim" `
	Discount    float64 `json:"discount,omitempty"      conform:"trim" `
	Quantity    int64   `json:"quantity,omitempty"      conform:"trim" `
	ImageBase64 *string `json:"image_base64,omitempty"  conform:"trim" `
}

type UpdateGearForm struct {
	Name        *string  `json:"name,omitempty"         db:"name"       validate:"omitempty"`
	Type        *string  `json:"type,omitempty"         db:"type"       validate:"omitempty,is-gear"`
	Price       *float64 `json:"price,omitempty"        db:"price"      validate:"omitempty"`
	Discount    *float64 `json:"discount,omitempty"     db:"discount"   validate:"omitempty"`
	Quantity    *int64   `json:"quantity,omitempty"     db:"quantity"   validate:"omitempty"`
	ImageBase64 string   `json:"image_base64,omitempty" db:"quantity"   validate:"omitempty"`
}
