package postcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/bayudha2/go-test-0/config"
	"github.com/bayudha2/go-test-0/helper"
	"github.com/bayudha2/go-test-0/models"
	"github.com/bayudha2/go-test-0/utils"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var postInput models.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&postInput); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request payload")
		return
	}

	var userInfo config.JWTClaim
	utils.ParseToken(&userInfo, r)

	defer r.Body.Close()

	postInput.UserId = userInfo.Userid
	err := postInput.CreatePost(models.DB)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, postInput)
}
