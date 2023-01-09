package authcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bayudha2/go-test-0/config"
	"github.com/bayudha2/go-test-0/helper"
	"github.com/bayudha2/go-test-0/helper/validation"
	"github.com/bayudha2/go-test-0/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type resp struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int    `json:"expires"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if listErr, err := validation.ValidateUser(&userInput); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	defer r.Body.Close()

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	if err := userInput.CreateUser(models.DB); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			helper.RespondWithError(w, http.StatusInternalServerError, "Username already used!")
			return
		} else {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.Login
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if listErr, err := validation.ValidateLogin(&userInput); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	defer r.Body.Close()

	var user models.User
	user.Username = userInput.Username

	if err := user.GetUser(models.DB); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			helper.RespondWithError(w, http.StatusUnauthorized, "Username or password is incorrect")
			return
		} else {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Your password is incorrect")
		return
	}

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken(user.Username, 15); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var refreshToken config.TokenPayload
	if err := refreshToken.CreateToken(user.Username, 30); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var payload models.AuthUserResponse
	payload.AccessToken = accessToken.Token
	payload.AccessTokenExpiresAt = accessToken.ExpTime
	payload.RefreshToken = refreshToken.Token
	payload.RefreshTokenExpiresAt = refreshToken.ExpTime
	payload.Username = user.Username

	var respPayload = resp{
		Username:     payload.Username,
		AccessToken:  payload.AccessToken,
		RefreshToken: payload.RefreshToken,
		Expires:      int(payload.AccessTokenExpiresAt.Unix()),
	}

	if err := payload.CreateSession(models.DB); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, respPayload)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mapToken); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	refreshToken := mapToken["refresh_token"]
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWT_KEY, nil
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
		refreshUsername, ok := claims["Username"].(string)

		if !ok {
			helper.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if claims["sub"] == "access_token" {
			helper.RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
			return
		}

		rsp := models.AuthUserResponse{Username: refreshUsername}
		if err := rsp.DeleteAuth(models.DB); err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		var accessToken config.TokenPayload
		if err := accessToken.CreateToken(refreshUsername, 15); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var refreshToken config.TokenPayload
		if err := refreshToken.CreateToken(refreshUsername, 30); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var payload models.AuthUserResponse
		payload.AccessToken = accessToken.Token
		payload.AccessTokenExpiresAt = accessToken.ExpTime
		payload.RefreshToken = refreshToken.Token
		payload.RefreshTokenExpiresAt = refreshToken.ExpTime
		payload.Username = refreshUsername

		var respPayload = resp{
			Username:     payload.Username,
			AccessToken:  payload.AccessToken,
			RefreshToken: payload.RefreshToken,
			Expires:      int(payload.AccessTokenExpiresAt.Unix()),
		}

		if err := payload.CreateSession(models.DB); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, respPayload)
	} else {
		helper.RespondWithError(w, http.StatusUnauthorized, "refresh expired")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		helper.RespondWithError(w, http.StatusUnauthorized, "Not Authorized!")
		return
	}

	getToken := strings.Split(bearer, " ")
	token, _ := jwt.Parse(getToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", "must bearer")
		}
		return config.JWT_KEY, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	username, _ := claims["Username"].(string)

	rsp := models.AuthUserResponse{Username: username}
	if err := rsp.DeleteAuth(models.DB); err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logout successfully"})
}
