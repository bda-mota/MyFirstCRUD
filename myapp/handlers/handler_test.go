package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bda-mota/MyFirstCRUD/myapp/handlers"
	"github.com/bda-mota/MyFirstCRUD/myapp/models"
	"github.com/bda-mota/MyFirstCRUD/myapp/repository"
)

// ****** CREATE PRODUCT ******
func TestCreateProduct(t *testing.T) {
	//simulando a função InsertProduct
	mockRepo := &repository.MockManualProductRepository{
		InsertProductFunc: func(p models.Product) (int64, error) {
			//define o que a função irá retornar, simula o acesso ao banco de dados, nesse caso, sucesso
			return 1, nil
		},
	}

	//indicando para o handler quer ele irá chamar, nesse caso o handler irá chamar o mock ao invés da funCão real
	handler := handlers.ProductHandler{
		Repo: mockRepo,
	}

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
		t.Errorf("expected body %s, got %s", expectedResponse, rr.Body.String())
	}
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	mockRepo := &repository.MockManualProductRepository{
		InsertProductFunc: func(p models.Product) (int64, error) {
			return 1, nil
		},
	}

	handler := handlers.ProductHandler{
		Repo: mockRepo,
	}

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

	handler := handlers.ProductHandler{
		Repo: mockRepo,
	}

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

func TestCreateProduct_RepositoryError(t *testing.T) {
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

// ****** UPDATE PRODUCT *******
func TestGetProductByID(t *testing.T) {

}
