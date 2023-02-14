package productcontroller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bayudha2/go-test-0/helper"
	"github.com/bayudha2/go-test-0/helper/validation"
	"github.com/bayudha2/go-test-0/models"
	"github.com/gorilla/mux"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := models.Product{ID: id}
	if err := p.GetProduct(models.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			helper.RespondWithError(w, http.StatusNotFound, "Product not found")
		default:
			helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, p)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
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

	if by == "" {
		by = "name"
	}

	var params = models.Params{
		Page:   page,
		Limit:  limit,
		Order:  order,
		By:     by,
		Search: search,
	}

	products, err := models.GetProducts(models.DB, params)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if listErr, err := validation.Validate(&p); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	defer r.Body.Close()

	err := p.CreateProduct(models.DB)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, p)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if listErr, err := validation.Validate(&p); err != nil {
		helper.RespondWithMultiError(w, http.StatusBadRequest, listErr)
		return
	}

	defer r.Body.Close()
	p.ID = id

	if err := p.UpdateProduct(models.DB); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, p)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := models.Product{ID: id}
	if err := p.DeleteProduct(models.DB); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
