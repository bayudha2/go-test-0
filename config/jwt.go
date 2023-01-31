package config

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bayudha2/go-test-0/helper"
	"github.com/golang-jwt/jwt/v4"
)

var JWT_KEY = []byte("asdads901023910239asdd")

type JWTClaim struct {
	Username string
	Userid   string
	jwt.RegisteredClaims
}

type TokenPayload struct {
	Token   string
	ExpTime time.Time
}

func (p *TokenPayload) CreateToken(userid string, username string, duration int) error {
	var subject string
	if duration == 30 {
		subject = "refresh_token"
	} else {
		subject = "access_token"
	}

	expTime := time.Now().Add(time.Duration(duration) * time.Minute)
	claims := JWTClaim{
		Username: username,
		Userid:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
			Subject:   subject,
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgo.SignedString(JWT_KEY)

	p.Token = token
	p.ExpTime = expTime

	return err
}

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			helper.RespondWithError(w, http.StatusUnauthorized, "Not Authorized!")
			return
		}

		getToken := strings.Split(bearer, " ")

		token, err := jwt.Parse(getToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", "must bearer")
			}
			return JWT_KEY, nil
		})

		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			if claims["sub"] == "refresh_token" {
				fmt.Println(claims)
				helper.RespondWithError(w, http.StatusUnauthorized, "Not Authorized!")
				return
			}

			next.ServeHTTP(w, r)
		} else {
			helper.RespondWithError(w, http.StatusUnauthorized, "Token expired!!!")
			return
		}
	})
}
