package repository

import "github.com/bda-mota/MyFirstCRUD/myapp/models"

type MockManualProductRepository struct {
	InsertProductFunc     func(p models.Product) (int64, error)
	GetProductByIDFunc    func(id int64) (models.Product, error)
	DeleteProductByIDFunc func(id int64) error
	UpdateProductByIDFunc func(id int64, p models.Product) error
	GetAllProductsFunc    func() (sp []models.Product, err error)
}

func (m *MockManualProductRepository) InsertProduct(p models.Product) (int64, error) {
	return m.InsertProductFunc(p)
}

func (m *MockManualProductRepository) GetProductByID(id int64) (models.Product, error) {
	return m.GetProductByIDFunc(id)
}

func (m *MockManualProductRepository) DeleteProductByID(id int64) error {
	return m.DeleteProductByIDFunc(id)
}

func (m *MockManualProductRepository) UpdateProductByID(id int64, p models.Product) error {
	return m.UpdateProductByIDFunc(id, p)
}

func (m *MockManualProductRepository) GetAllProducts() (sp []models.Product, err error) {
	return m.GetAllProductsFunc()
}
