package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        string    `json:"id" validate:"omitempty"`
	PostId    string    `json:"post_id" validate:"required"`
	UserId    string    `json:"user_id" validate:"omitempty"`
	Content   string    `json:"content" validate:"required"`
	CommentId *string   `json:"comment_id" validate:"omitempty"`
	HasChild  bool      `json:"has_child"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PayloadComments struct {
	Data      []Comment `json:"data"`
	TotalData int       `json:"total_data"`
}

func (p *Comment) GetComment(db *sql.DB) error {
	err := db.QueryRow("SELECT * FROM comments WHERE id=$1", p.ID).Scan(&p.ID, &p.PostId, &p.UserId, &p.Content, &p.CommentId, &p.CreatedAt, &p.UpdatedAt)

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	if err != nil {
		return err
	}
	return nil
}

func (p *Comment) GetAllCommentByPost(db *sql.DB) (PayloadComments, error) {
	var result PayloadComments
	query := fmt.Sprintf(`
	SELECT * FROM comments
	WHERE post_id = '%s'
	`, p.PostId)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err.Error())
		return result, err
	}

	count := fmt.Sprintf(`
	SELECT COUNT(*) FROM comments WHERE post_id = '%s'`,
		p.PostId)

	if err := db.QueryRow(count).Scan(&result.TotalData); err != nil {
		return result, err
	}

	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostId, &c.UserId, &c.Content, &c.CommentId, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return result, err
		}

		var child string
		err := db.QueryRow("SELECT id FROM comments WHERE parent_id = $1", c.ID).Scan(&child)

		c.HasChild = true
		if err != nil && err == sql.ErrNoRows {
			c.HasChild = false
		}

		c.CreatedAt = c.CreatedAt.UTC().Add(time.Hour * 7)
		c.UpdatedAt = c.UpdatedAt.UTC().Add(time.Hour * 7)

		if c.CommentId != nil {
			continue
		}

		comments = append(comments, c)
	}

	result.Data = comments
	return result, nil
}

func (p *Comment) CreateComment(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO comments(id, post_id, user_id, content, parent_id, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, post_id, user_id, content, parent_id, created_at, updated_at`,
		uuid.New().String(), p.PostId, p.UserId, p.Content, p.CommentId, time.Now(), time.Now(),
	).Scan(&p.ID, &p.PostId, &p.UserId, &p.Content, &p.CommentId, &p.CreatedAt, &p.UpdatedAt)

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	if err != nil {
		return err
	}

	return nil
}

func (p *Comment) UpdateComment(db *sql.DB) error {
	err := db.QueryRow(`UPDATE comments SET content=$1, updated_at=$2
		WHERE id=$3 AND user_id=$4
		RETURNING id, post_id, user_id, content, parent_id, created_at, updated_at`,
		p.Content, time.Now(), p.ID, p.UserId,
	).Scan(&p.ID, &p.PostId, &p.UserId, &p.Content, &p.CommentId, &p.CreatedAt, &p.UpdatedAt)

	p.CreatedAt = p.CreatedAt.UTC().Add(time.Hour * 7)
	p.UpdatedAt = p.UpdatedAt.UTC().Add(time.Hour * 7)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (p *Comment) DeleteComment(db *sql.DB) error {
	res, err := db.Exec("DELETE FROM comments WHERE id=$1 AND user_id=$2", p.ID, p.UserId)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("Not found !")
	}

	return nil
}
