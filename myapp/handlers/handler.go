package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bda-mota/MyFirstCRUD/myapp/models"
	"github.com/bda-mota/MyFirstCRUD/myapp/repository"
	"github.com/gorilla/mux"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}

	if newProduct.Name == "" {
		ResponseError(w, "Name is required", http.StatusBadRequest)
		return
	}
	if newProduct.Price <= 0 {
		ResponseError(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	id, _ := repository.InsertProduct(newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product created successfully",
		"id":      id,
	})
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	getProduct, err := repository.GetProductByID(convertedId)
	if err != nil {
		ResponseError(w, "could not retrieve product", http.StatusBadRequest)
		return
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

	err = repository.DeleteProductByID(convertedId)
	if err != nil {
		ResponseError(w, "product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// func UpdateProductByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	productID := vars["id"]

// 	convertedId, err := strconv.ParseInt(productID, 10, 64)
// 	if err != nil {
// 		ResponseError(w, "invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	var updateProduct models.Product
// 	json.NewDecoder(r.Body).Decode(&updateProduct)

// 	index := -1
// 	for i, value := range products {
// 		if value.ID == convertedId {
// 			index = i
// 			if updateProduct.Name != "" {
// 				products[i].Name = updateProduct.Name
// 			}
// 			if updateProduct.Price > 0 {
// 				products[i].Price = updateProduct.Price
// 			}
// 			break
// 		}
// 	}

// 	if index == -1 {
// 		ResponseError(w, "Product not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(products[index])
// }
