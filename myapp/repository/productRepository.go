package repository

import (
	"fmt"

	"github.com/bda-mota/MyFirstCRUD/myapp/config"
	"github.com/bda-mota/MyFirstCRUD/myapp/models"
)

func InsertProduct(p models.Product) (int64, error) {
	conn, err := config.OpenConn()
	if err != nil {
		return 0, fmt.Errorf("could not connect to the database: %v", err)
	}
	defer conn.Close()

	sql := `INSERT INTO products (name, value) VALUES ($1, $2) RETURNING id`

	var id int64
	err = conn.QueryRow(sql, p.Name, p.Value).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("could not insert product: %v", err)
	}

	return id, nil
}
