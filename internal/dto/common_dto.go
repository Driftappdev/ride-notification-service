package dto

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Result[T any] struct {
	Data T `json:"data"`
}
