package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          string `json:"id" validate:"omitempty"`
	UserId      string `json:"user_id" validate:"omitempty"`
	Description string `json:"description" validate:"required"`
	// Likes       uint      `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PayloadPosts struct {
	Data      []Post `json:"data"`
	TotalData int    `json:"total_data"`
}

func (p *Post) CreatePost(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO posts(id, user_id, description, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING id, user_id, description, created_at, updated_at`,
		uuid.New().String(), p.UserId, p.Description, time.Now(), time.Now()).Scan(&p.ID, &p.UserId, &p.Description, &p.CreatedAt, &p.UpdatedAt)

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	if err != nil {
		return err
	}
	return nil
}

func (p *Post) UpdatePost(db *sql.DB) error {
	err := db.QueryRow(`UPDATE posts SET description=$1, updated_at=$2
		WHERE id=$3 AND user_id=$4
		RETURNING id, user_id, description, created_at, updated_at`,
		p.Description, time.Now(), p.ID, p.UserId,
	).Scan(&p.ID, &p.UserId, &p.Description, &p.CreatedAt, &p.UpdatedAt)

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	if err != nil {
		return err
	}

	return nil
}

func (p *Post) DeletePost(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM posts WHERE id=$1 AND user_id=$2", p.ID, p.UserId)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Unauthorized request!")
	}
	return nil
}

func (p *Post) GetPost(db *sql.DB) error {
	err := db.QueryRow("SELECT * FROM posts WHERE id=$1", p.ID).Scan(&p.ID, &p.UserId, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetPosts(db *sql.DB, params Params) (PayloadPosts, error) {
	var result PayloadPosts
	query := fmt.Sprintf(`
		SELECT * From posts
		WHERE user_id = '%s'
		ORDER BY created_at %s LIMIT %v OFFSET ((%v - 1) * %v)
	`, p.UserId, params.Order, params.Limit, params.Page, params.Limit)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err.Error())
		return result, err
	}

	count := fmt.Sprintf(`
	SELECT COUNT(*) FROM posts WHERE user_id = '%s'
	`, p.UserId)
	if err := db.QueryRow(count).Scan(&result.TotalData); err != nil {
		return result, err
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserId, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return result, err
		}

		p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
		p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)
		posts = append(posts, p)
	}

	result.Data = posts
	return result, nil
}
