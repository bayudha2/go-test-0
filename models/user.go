package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id" validate:"omitempty"`
	Username  string    `json:"username" validate:"required"`
	Fullname  string    `json:"fullname" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *User) CreateUser(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users(id, fullname, username, password, email, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", uuid.New().String(), p.Fullname, p.Username, p.Password, p.Email, time.Now())

	if err != nil {
		return err
	}
	return nil
}

func (p *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM users WHERE username=$1", p.Username).Scan(&p.ID, &p.Fullname, &p.Username, &p.Password, &p.Email, &p.CreatedAt)
}
