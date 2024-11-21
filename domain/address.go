package domain

import "github.com/google/uuid"

type Address struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`

	Address string `json:"address"`
	Country string `json:"country"`
}

type AddAddressForm struct {
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required,gt=0"`

	Address string `json:"address" db:"address" validate:"required,gt=0"`
	Country string `json:"country" db:"country" validate:"required,gt=0"`
}

type UpdateAddressForm struct {
	Address *string `json:"address,omitempty" db:"address" validate:"required,gt=0"`
	Country *string `json:"country,omitempty" db:"country" validate:"omitempty,gt=0"`
}
