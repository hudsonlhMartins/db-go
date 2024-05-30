package main

import (
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	Name  string
	Price float64
	ID    string
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		Name:  name,
		Price: price,
		ID:    uuid.New().String(),
	}
}

func main() {
	db, err := sql.Open("mysql", "docker:docker@tcp(localhost:3306)/docker")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := selectAllProducts(db)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(res)

	if err != nil {
		panic(err)
	}

	println(string(json))
}

func insertProduct(db *sql.DB, product *Product) error {
	smt, err := db.Prepare("INSERT INTO products (id, name, price) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func selectproductByID(db *sql.DB, id string) (*Product, error) {
	smt, err := db.Prepare("SELECT id, name, price FROM products WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	var product Product
	err = smt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func selectAllProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func updateItemByID(db *sql.DB, product *Product) error {
	smt, err := db.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}
