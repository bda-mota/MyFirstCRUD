package repository

import (
	"database/sql"
	"fmt"

	"github.com/bda-mota/MyFirstCRUD/myapp/config"
	"github.com/bda-mota/MyFirstCRUD/myapp/models"
)

// POST
func InsertProduct(p models.Product) (int64, error) {
	conn, err := config.OpenConn()
	if err != nil {
		return 0, fmt.Errorf("could not connect to the database: %v", err)
	}
	defer conn.Close()

	sql := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`

	var id int64
	err = conn.QueryRow(sql, p.Name, p.Price).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("could not insert product: %v", err)
	}

	return id, nil
}

// GET
func GetProductByID(id int64) (models.Product, error) {
	conn, err := config.OpenConn()
	if err != nil {
		return models.Product{}, fmt.Errorf("could not connect to the database: %v", err)
	}
	defer conn.Close()

	var getProduct models.Product
	row := `SELECT id, name, price FROM products WHERE id = $1`
	err = conn.QueryRow(row, id).Scan(&getProduct.ID, &getProduct.Name, &getProduct.Price)

	if err == sql.ErrNoRows {
		return models.Product{}, nil
	} else if err != nil {
		return models.Product{}, fmt.Errorf("could not retrieve product: %v", err)
	}

	return getProduct, nil
}

// DELETE
func DeleteProductByID(id int64) error {
	conn, err := config.OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	sql := `DELETE FROM products WHERE id = $1`
	res, err := conn.Exec(sql, id)
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
func UpdateProductByID(id int64, p models.Product) error {
	conn, err := config.OpenConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	sql := `UPDATE products SET name = $1, price = $2 WHERE id = $3`
	res, err := conn.Exec(sql, p.Name, p.Price, id)
	if err != nil {
		return fmt.Errorf("could not update product: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("no product found with ID %d", id)
	}
	return nil
}
