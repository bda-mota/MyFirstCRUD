package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bda-mota/MyFirstCRUD/myapp/models"
	"github.com/bda-mota/MyFirstCRUD/myapp/repository"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	Repo repository.ProductRepository
}

// POST
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product

	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
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

	id, err := h.Repo.InsertProduct(newProduct)
	if err != nil {
		ResponseError(w, "Could not insert the product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product created successfully",
		"id":      id,
	})
}

// GET
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	getProduct, err := h.Repo.GetProductByID(convertedId)
	if err != nil {
		ResponseError(w, "could not retrieve product", http.StatusInternalServerError)
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

// DELETE
func (h *ProductHandler) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteProductByID(convertedId)
	if err != nil {
		if err == sql.ErrNoRows {
			ResponseError(w, "Product not found", http.StatusNotFound)
		} else {
			ResponseError(w, "could not delete product", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}

// PUT
func (h *ProductHandler) UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	convertedId, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		ResponseError(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var updateProduct models.Product
	if err = json.NewDecoder(r.Body).Decode(&updateProduct); err != nil {
		ResponseError(w, "Could not update product", http.StatusInternalServerError)
		return
	}

	err = h.Repo.UpdateProductByID(convertedId, updateProduct)
	if err != nil {
		if err.Error() == "product not found" {
			ResponseError(w, "Product not found", http.StatusNotFound)
		} else {
			ResponseError(w, "could not update product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
}

// GET ALL
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	list, err := h.Repo.GetAllProducts()
	if err != nil {
		ResponseError(w, "products not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}
