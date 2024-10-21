package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bda-mota/MyFirstCRUD/myapp/handlers"
	"github.com/bda-mota/MyFirstCRUD/myapp/models"
	"github.com/bda-mota/MyFirstCRUD/myapp/repository"
	"github.com/gorilla/mux"
)

// ****** CREATE ******
func TestCreateProduct_Success(t *testing.T) {
	//simulando a função InsertProduct
	mockRepo := &repository.MockManualProductRepository{
		InsertProductFunc: func(p models.Product) (int64, error) {
			//define o que a função irá retornar, simula o acesso ao banco de dados, nesse caso, sucesso
			return 1, nil
		},
	}

	//indicando para o handler quer ele irá chamar, nesse caso o handler irá chamar o mock ao invés da funCão real
	handler := handlers.ProductHandler{Repo: mockRepo}

	//criando o produto da requisição
	product := models.Product{Name: "Test Product", Price: 10}
	//mudando ele para formato json, para ser o r.body
	body, _ := json.Marshal(product)
	//cria uma nova requisição do tipo POST para o endpoint /products
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	//para capturar a resposta da requisição, falso responseWriter
	rr := httptest.NewRecorder()

	handler.CreateProduct(rr, req)

	//verifica status da resposta, e se for diferente retorna um erro
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	//verificando o corpo da resposta
	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"id":1,"message":"Product created successfully"}`
	if actualResponse != expectedResponse {
		t.Errorf("expected body %s, got %s", expectedResponse, actualResponse)
	}
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		InsertProductFunc: func(p models.Product) (int64, error) {
			return 1, nil
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	product := models.Product{Name: "Test Product", Price: 0}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CreateProduct(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"Price must be greater than 0","errorCode":400}`
	if actualResponse != expectedResponse {
		t.Errorf("expected body %s, got %s", expectedResponse, rr.Body.String())
	}
}

func TestCreateProduct_InvalidName(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		InsertProductFunc: func(p models.Product) (int64, error) {
			return 1, nil
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	product := models.Product{Name: "", Price: 10}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CreateProduct(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"Name is required","errorCode":400}`
	if actualResponse != expectedResponse {
		t.Errorf("expected body %s, got %s", expectedResponse, rr.Body.String())
	}
}

func TestCreateProduct_ServerError(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		InsertProductFunc: func(p models.Product) (int64, error) {
			return 0, fmt.Errorf("database error")
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	product := models.Product{Name: "Test Product", Price: 10.0}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CreateProduct(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, status)
	}
	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"Could not insert the product","errorCode":500}`
	if actualResponse != expectedResponse {
		t.Errorf("expected body %s, got %s", expectedResponse, actualResponse)
	}
}

// ****** GET *******
func TestGetProductByID_Success(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		GetProductByIDFunc: func(id int64) (models.Product, error) {
			return models.Product{ID: 1, Name: "Test Product", Price: 10.0}, nil
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	req, _ := http.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.GetProductByID).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	expectedProduct := models.Product{ID: 1, Name: "Test Product", Price: 10.0}
	var actualProduct models.Product
	if err := json.NewDecoder(rr.Body).Decode(&actualProduct); err != nil {
		t.Fatalf("failed to decode response %v", err)
	}

	if actualProduct != expectedProduct {
		t.Errorf("expected product %v, got %v", expectedProduct, actualProduct)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		GetProductByIDFunc: func(id int64) (models.Product, error) {
			return models.Product{}, nil
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	req, _ := http.NewRequest("GET", "/products/0", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.GetProductByID).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"Product not found","errorCode":404}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

func TestGetProductByID_ServerError(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		GetProductByIDFunc: func(id int64) (models.Product, error) {
			return models.Product{}, fmt.Errorf("could not retrieve product")
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	req, _ := http.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.GetProductByID).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"could not retrieve product","errorCode":500}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

// ****** DELETE *******
func TestDeleteProductByID_Success(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		DeleteProductByIDFunc: func(id int64) error { return nil },
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.DeleteProductByID).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"message":"Product deleted successfully"}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

func TestDeleteProductByID_NotFound(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		DeleteProductByIDFunc: func(id int64) error {
			return sql.ErrNoRows
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.DeleteProductByID).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"Product not found","errorCode":404}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

func TestDeleteProductByID_ServerError(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		DeleteProductByIDFunc: func(id int64) error {
			return fmt.Errorf("could not delete product")
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.DeleteProductByID).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"could not delete product","errorCode":500}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

// ****** UPDATE *******
func TestUpdateProductByID_Success(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		UpdateProductByIDFunc: func(id int64, p models.Product) error {
			return nil
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	product := models.Product{Name: "New name", Price: 10}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.UpdateProductByID).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"message":"Product updated successfully"}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

func TestUpdateProductByID_NotFound(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		UpdateProductByIDFunc: func(id int64, p models.Product) error {
			return fmt.Errorf("product not found")
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	product := models.Product{Name: "non-existent name", Price: 10}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.UpdateProductByID).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"Product not found","errorCode":404}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

func TestUpdateProductByID_ServerError(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		UpdateProductByIDFunc: func(id int64, p models.Product) error {
			return fmt.Errorf("internal server error")
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}

	product := models.Product{Name: "Test name", Price: 10}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.UpdateProductByID).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	actualResponse := strings.TrimSpace(rr.Body.String())
	expectedResponse := `{"error":"could not update product","errorCode":500}`
	if actualResponse != expectedResponse {
		t.Errorf("expected product %v, got %v", expectedResponse, actualResponse)
	}
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		GetAllProductsFunc: func() (sp []models.Product, err error) {
			return []models.Product{
				{ID: 1, Name: "Product 1", Price: 10.0},
				{ID: 2, Name: "Product 2", Price: 15.0},
			}, nil
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}
	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products", handler.GetAllProducts).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	expectedResponse := `[{"id":1,"name":"Product 1","price":10},{"id":2,"name":"Product 2","price":15}]`
	actualResponse := strings.TrimSpace(rr.Body.String())
	if actualResponse != expectedResponse {
		t.Errorf("expected body %s, got %s", expectedResponse, actualResponse)
	}
}

func TestGetAll_NotFound(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		GetAllProductsFunc: func() (sp []models.Product, err error) {
			return []models.Product{}, fmt.Errorf("no products found")
		},
	}

	handler := handlers.ProductHandler{Repo: mockRepo}
	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/products", handler.GetAllProducts).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, status)
	}

	expectedResponse := `{"error":"no products found","errorCode":404}`
	actualResponse := strings.TrimSpace(rr.Body.String())
	if actualResponse != expectedResponse {
		t.Errorf("expected body %s, got %s", expectedResponse, actualResponse)
	}
}
