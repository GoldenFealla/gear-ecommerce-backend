package domain

type ResponseSuccess struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Message string `json:"message"`
}
