package utils

import (
	"net/http"
	"strings"

	"github.com/bayudha2/go-test-0/config"
	"github.com/golang-jwt/jwt"
)

func ParseToken(jwtclaim *config.JWTClaim, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	getToken := strings.Split(bearer, " ")

	token, _ := jwt.Parse(getToken[1], func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	jwtclaim.Username = claims["Username"].(string)
	jwtclaim.Userid = claims["Userid"].(string)
}
