package domain

import "github.com/google/uuid"

type Gear struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Type     string    `json:"type" db:"type"`
	Price    float64   `json:"price" db:"price"`
	Discount float64   `json:"discount" db:"discount"`
	Quantity int64     `json:"quantity" db:"quantity"`
}

type AddGearForm struct {
	Name     string  `json:"name,omitempty" validate:"required"`
	Type     string  `json:"type,omitempty" validate:"required"`
	Price    float64 `json:"price,omitempty"`
	Discount float64 `json:"discount,omitempty"`
	Quantity int64   `json:"quantity,omitempty"`
}

type UpdateGearForm struct {
	Name     *string  `json:"name,omitempty"`
	Type     *string  `json:"type,omitempty"`
	Price    *float64 `json:"price,omitempty"`
	Discount *float64 `json:"discount,omitempty"`
	Quantity *int64   `json:"quantity,omitempty"`
}
