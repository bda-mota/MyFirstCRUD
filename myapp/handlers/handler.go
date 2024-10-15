package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bda-mota/MyFirstCRUD/myapp/models"
	"github.com/gorilla/mux"
)

var products []models.Product
var id int64

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product

	json.NewDecoder(r.Body).Decode(&newProduct)

	if newProduct.Name == "" {
		ResponseError(w, "Name is required", http.StatusBadRequest)
		return
	}
	if newProduct.Value <= 0 {
		ResponseError(w, "Value must be greater than 0", http.StatusBadRequest)
		return
	}

	id++
	newProduct.ID = id
	products = append(products, newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	var getProduct models.Product
	for _, value := range products {
		if value.ID == convertedId {
			getProduct = value
			break
		}
	}

	if getProduct.ID == 0 {
		ResponseError(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getProduct)
}

func DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	index := -1
	for i, value := range products {
		if value.ID == convertedId {
			index = i
			break
		}
	}

	if index == -1 {
		ResponseError(w, "Product not found", http.StatusNotFound)
		return
	}

	products = append(products[:index], products[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var updateProduct models.Product
	json.NewDecoder(r.Body).Decode(&updateProduct)

	index := -1
	for i, value := range products {
		if value.ID == convertedId {
			index = i
			if updateProduct.Name != "" {
				products[i].Name = updateProduct.Name
			}
			if updateProduct.Value > 0 {
				products[i].Value = updateProduct.Value
			}
			break
		}
	}

	if index == -1 {
		ResponseError(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products[index])
}
