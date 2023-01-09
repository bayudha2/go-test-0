package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        string    `json:"id" validate:"omitempty"`
	Name      string    `json:"name" validate:"required,min=10"`
	Price     float64   `json:"price" validate:"required,number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductParams struct {
	Page   string
	Limit  string
	Order  string
	By     string
	Search string
}

type payload struct {
	Data      []Product `json:"data"`
	TotalData int       `json:"total_data"`
}

var TotalData int

func (p *Product) GetProduct(db *sql.DB) error {
	err := db.QueryRow("SELECT name, price, created_at, updated_at FROM products WHERE id=$1", p.ID).Scan(&p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	return nil

}

func (p *Product) UpdateProduct(db *sql.DB) error {
	err := db.QueryRow("UPDATE products SET name=$1, price=$2, updated_at=$3 WHERE id=$4 RETURNING id, name, price, created_at, updated_at", p.Name, p.Price, time.Now(), p.ID).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	return nil
}

func (p *Product) DeleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *Product) CreateProduct(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO products(id, name, price, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id, name, price, created_at, updated_at", uuid.New().String(), p.Name, p.Price, time.Now(), time.Now()).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	return nil
}

func GetProducts(db *sql.DB, param ProductParams) (payload, error) {
	var respPayload payload
	query := fmt.Sprintf(`
	SELECT id, name, price, created_at, updated_at
		FROM products
		WHERE name ILIKE '%%%s%%'
		ORDER BY %s %s LIMIT %v OFFSET ((%v - 1) * %v)
	`, param.Search, param.By, param.Order, param.Limit, param.Page, param.Limit)
	rows, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		return respPayload, err
	}

	count := fmt.Sprintf(`
		SELECT COUNT(*) FROM products WHERE name ILIKE '%%%s%%'
	`, param.Search)
	if err := db.QueryRow(count).Scan(&respPayload.TotalData); err != nil {
		return respPayload, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return respPayload, err
		}
		p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
		p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)
		products = append(products, p)
	}

	respPayload.Data = products
	return respPayload, nil
}
