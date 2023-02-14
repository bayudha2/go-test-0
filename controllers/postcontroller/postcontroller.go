package postcontroller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bayudha2/go-test-0/config"
	"github.com/bayudha2/go-test-0/helper"
	"github.com/bayudha2/go-test-0/helper/validation"
	"github.com/bayudha2/go-test-0/models"
	"github.com/bayudha2/go-test-0/utils"
	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var postInput models.Post
	if r.Body == nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&postInput); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if listErr, err := validation.Validate(&postInput); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
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

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if r.Body == nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&post); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if listErr, err := validation.Validate(&post); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	var useInfo config.JWTClaim
	utils.ParseToken(&useInfo, r)

	defer r.Body.Close()
	post.UserId = useInfo.Userid
	post.ID = id

	err := post.UpdatePost(models.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized request!")
		default:
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	var userInfo config.JWTClaim
	utils.ParseToken(&userInfo, r)

	post.UserId = userInfo.Userid
	post.ID = id

	err := post.DeletePost(models.DB)
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized request") {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		} else {
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	order := r.URL.Query().Get("order")
	by := r.URL.Query().Get("by")
	search := r.URL.Query().Get("search")

	if page == "" {
		page = strconv.Itoa(1)
	}

	if limit == "" {
		limit = strconv.Itoa(10)
	}

	if order == "" {
		order = "asc"
	}

	var params = models.Params{
		Page:   page,
		Limit:  limit,
		Order:  order,
		By:     by,
		Search: search,
	}

	var post models.Post
	var userInfo config.JWTClaim

	utils.ParseToken(&userInfo, r)
	post.UserId = userInfo.Userid

	posts, err := post.GetPosts(models.DB, params)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	post.ID = id

	if err := post.GetPost(models.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		default:
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, post)
}
