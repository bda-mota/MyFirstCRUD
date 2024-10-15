package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bda-mota/MyFirstCRUD/myapp/models"
)

func ResponseError(w http.ResponseWriter, message string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)

	responseError := models.RequestError{
		Message:   message,
		ErrorCode: errorCode,
	}

	json.NewEncoder(w).Encode(responseError)
}
