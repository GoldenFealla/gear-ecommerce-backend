package domain

import (
	"github.com/google/uuid"
)

type OrderStatus string

const (
	CART       OrderStatus = "CART"
	PAID       OrderStatus = "PAID"
	DELIVERING OrderStatus = "DELIVERING"
	DONE       OrderStatus = "DONE"
)

type Order struct {
	ID     uuid.UUID   `json:"id" db:"id"`
	Status OrderStatus `json:"status" db:"status"`
	UserID uuid.UUID   `json:"user_id" db:"user_id"`
	Total  int64       `json:"total" db:"total"`
}

type OrderGear struct {
	Gear     *Gear `json:"gear"`
	Quantity int64 `json:"quantity"`
}

type FullOrder struct {
	Order     *Order       `json:"order"`
	OrderGear []*OrderGear `json:"order_gear"`
}

type AddOrderForm struct {
	Name       string      `json:"name"`
	Status     OrderStatus `json:"status"`
	GearIDList []string    `json:"gear_id_list"`
}
