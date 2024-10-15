package models

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type RequestError struct {
	Message   string `json:"error"`
	ErrorCode int    `json:"errorCode"`
}
