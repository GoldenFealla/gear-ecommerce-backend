package domain

type Gear struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
	Quantity int64   `json:"quantity"`
}

type AddGearForm struct {
	Name     string  `json:"name,omitempty" validate:"required"`
	Type     string  `json:"type,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Discount float64 `json:"discount,omitempty"`
	Quantity int64   `json:"quantity,omitempty"`
}

type UpdateGearForm struct {
	Name     string  `json:"name,omitempty"`
	Type     string  `json:"type,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Discount float64 `json:"discount,omitempty"`
	Quantity int64   `json:"quantity,omitempty"`
}
