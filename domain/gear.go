package domain

import "github.com/google/uuid"

type GearType int

const (
	PC GearType = iota + 1
	Laptop
	MainBoard
	CPU
	GPU
	PSU
	RAM
	Fan
	Storage
	Monitor
)

type Gear struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Type     GearType  `json:"type" db:"type"`
	Price    float64   `json:"price" db:"price"`
	Discount float64   `json:"discount" db:"discount"`
	Quantity int64     `json:"quantity" db:"quantity"`
	ImageURL string    `json:"image_url" db:"image_url"`
}

type AddGearForm struct {
	Name        string   `json:"name,omitempty" validate:"required"`
	Type        GearType `json:"type,omitempty" validate:"required"`
	Price       float64  `json:"price,omitempty"`
	Discount    float64  `json:"discount,omitempty"`
	Quantity    int64    `json:"quantity,omitempty"`
	ImageBase64 *string  `json:"image_base64,omitempty"`
}

type UpdateGearForm struct {
	Name        *string   `json:"name,omitempty" db:"name" validate:"omitempty,gte=2,lte=64"`
	Type        *GearType `json:"type,omitempty" db:"type" validate:"omitempty,gte=2,lte=32"`
	Price       *float64  `json:"price,omitempty" db:"price" validate:"omitempty,gt=0"`
	Discount    *float64  `json:"discount,omitempty" db:"discount" validate:"omitempty,gt=0"`
	Quantity    *int64    `json:"quantity,omitempty" db:"quantity" validate:"omitempty,gt=0"`
	ImageBase64 string    `json:"image_base64,omitempty" db:"quantity" validate:"omitempty,gt=0"`
}
