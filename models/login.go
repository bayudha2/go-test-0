package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthUserResponse struct {
	Username              string    `json:"username"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (p *AuthUserResponse) CreateSession(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO sessions(id, username, refresh_token, expires_at, created_at) VALUES($1, $2, $3, $4, $5) RETURNING username, refresh_token, expires_at, created_at", uuid.New().String(), p.Username, p.RefreshToken, p.RefreshTokenExpiresAt, time.Now())
	return err
}

func (p *AuthUserResponse) DeleteAuth(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM sessions WHERE username=$1", p.Username)
	return err
}
