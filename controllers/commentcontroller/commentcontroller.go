package commentcontroller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bayudha2/go-test-0/config"
	"github.com/bayudha2/go-test-0/helper"
	"github.com/bayudha2/go-test-0/helper/validation"
	"github.com/bayudha2/go-test-0/models"
	"github.com/bayudha2/go-test-0/utils"
	"github.com/gorilla/mux"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var commentInput models.Comment
	if r.Body == nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&commentInput); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if listErr, err := validation.Validate(&commentInput); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	var userInfo config.JWTClaim
	utils.ParseToken(&userInfo, r)

	defer r.Body.Close()

	commentInput.UserId = userInfo.Userid
	err := commentInput.CreateComment(models.DB)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, commentInput)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Id")
		return
	}

	var commentInput models.Comment
	if r.Body == nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&commentInput); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if listErr, err := validation.Validate(&commentInput); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	var userInfo config.JWTClaim
	utils.ParseToken(&userInfo, r)

	defer r.Body.Close()

	commentInput.UserId = userInfo.Userid
	commentInput.ID = id
	err := commentInput.UpdateComment(models.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helper.RespondWithError(w, http.StatusNotFound, "Not Found!")
		default:
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, commentInput)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Id")
		return
	}

	var comment models.Comment
	var userInfo config.JWTClaim
	utils.ParseToken(&userInfo, r)

	comment.ID = id
	comment.UserId = userInfo.Userid

	err := comment.DeleteComment(models.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helper.RespondWithError(w, http.StatusNotFound, err.Error())
		default:
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

func GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&comment); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if comment.PostId == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Post id required!")
		return
	}

	defer r.Body.Close()

	comments, err := comment.GetAllCommentByPost(models.DB)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, comments)
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request id")
		return
	}

	var comment models.Comment
	comment.ID = id

	err := comment.GetComment(models.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helper.RespondWithError(w, http.StatusNotFound, "Comment not found!")
		default:
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, comment)
}
