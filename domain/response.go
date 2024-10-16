package domain

type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
