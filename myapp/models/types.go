package models

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type RequestError struct {
	Message   string `json:"error"`
	ErrorCode int    `json:"errorCode"`
}
