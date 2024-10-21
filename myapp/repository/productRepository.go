package repository

import (
	"database/sql"
	"fmt"

	"github.com/bda-mota/MyFirstCRUD/myapp/models"
)

type ProductRepository interface {
	InsertProduct(p models.Product) (int64, error)
	GetProductByID(id int64) (models.Product, error)
	DeleteProductByID(id int64) error
	UpdateProductByID(id int64, p models.Product) error
	GetAllProducts() (sp []models.Product, err error)
}
type PostgresProductRepository struct {
	DB *sql.DB
}

// POST
func (r *PostgresProductRepository) InsertProduct(p models.Product) (int64, error) {
	sql := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`

	var id int64
	if err := r.DB.QueryRow(sql, p.Name, p.Price).Scan(&id); err != nil {
		return 0, fmt.Errorf("could not insert product: %v", err)
	}

	return id, nil
}

// GET
func (r *PostgresProductRepository) GetProductByID(id int64) (models.Product, error) {
	var getProduct models.Product
	row := `SELECT id, name, price FROM products WHERE id = $1`
	err := r.DB.QueryRow(row, id).Scan(&getProduct.ID, &getProduct.Name, &getProduct.Price)

	if err == sql.ErrNoRows {
		return models.Product{}, nil
	} else if err != nil {
		return models.Product{}, fmt.Errorf("could not retrieve product: %v", err)
	}

	return getProduct, nil
}

// DELETE
func (r *PostgresProductRepository) DeleteProductByID(id int64) error {
	sql := `DELETE FROM products WHERE id = $1`
	res, err := r.DB.Exec(sql, id)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rowws affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product found with ID %d", id)
	}
	return nil
}

// PUT
func (r *PostgresProductRepository) UpdateProductByID(id int64, p models.Product) error {
	sql := `UPDATE products SET name = $1, price = $2 WHERE id = $3`
	res, err := r.DB.Exec(sql, p.Name, p.Price, id)
	if err != nil {
		return fmt.Errorf("could not update product: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// GET ALL
func (r *PostgresProductRepository) GetAllProducts() (sp []models.Product, err error) {
	rows, err := r.DB.Query(`SELECT * FROM products`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product

		if err = rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			continue
		}

		sp = append(sp, p)
	}

	if len(sp) == 0 {
		return nil, fmt.Errorf("no products found")
	}
	return sp, nil
}
