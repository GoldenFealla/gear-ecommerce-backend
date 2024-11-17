package domain

import "github.com/google/uuid"

type Address struct {
	ID      uuid.UUID `json:"id"`
	Street  string    `json:"street"`
	Region  string    `json:"region"` // Can be State, Provine or District
	City    string    `json:"city"`
	Postal  string    `json:"postal"`
	Country string    `json:"country"`
	UserID  uuid.UUID `json:"user_id"`
}

type AddAddressForm struct {
	Street  string    `json:"street" db:"street" validate:"required,gt=0"`
	Region  string    `json:"region" db:"region" validate:"required,gt=0"`
	City    string    `json:"city" db:"city" validate:"required,gt=0"`
	Postal  string    `json:"postal" db:"postal" validate:"required,gt=0"`
	Country string    `json:"country" db:"country" validate:"required,gt=0"`
	UserID  uuid.UUID `json:"user_id" db:"user_id" validate:"required,gt=0"`
}

type UpdateAddressForm struct {
	Street  *string `json:"street,omitempty" db:"street" validate:"omitempty,gt=0"`
	Region  *string `json:"region,omitempty" db:"region" validate:"omitempty,gt=0"`
	City    *string `json:"city,omitempty" db:"city" validate:"omitempty,gt=0"`
	Postal  *string `json:"postal,omitempty" db:"postal" validate:"omitempty,gt=0"`
	Country *string `json:"country,omitempty" db:"country" validate:"omitempty,gt=0"`
}
