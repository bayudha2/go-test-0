package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          string    `json:"id" validate:"omitempty"`
	UserId      string    `json:"user_id" validate:"omitempty"`
	Description string    `json:"description" validate:"required"`
	Likes       uint      `json:"likes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Post) CreatePost(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO posts(id, user_id, description, likes, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6) 
		RETURNING id, user_id, description, likes, created_at, updated_at`,
		uuid.New().String(), p.UserId, p.Description, 0, time.Now(), time.Now()).Scan(&p.ID, &p.UserId, &p.Description, &p.Likes, &p.CreatedAt, &p.UpdatedAt)

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	if err != nil {
		return err
	}
	return nil
}
