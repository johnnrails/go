package handlers

import (
	"net/http"
	"strconv"

	"github.com/johnnrails/ddd_go/microservices/products_api/helpers"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain/repositories"
	"github.com/johnnrails/ddd_go/second_ddd_go/response"
	"github.com/julienschmidt/httprouter"
)

type NewsRoutesHandler struct {
	repository repositories.NewsRepository
}

func CreateNewsRoutesHandler(repo repositories.NewsRepository) *NewsRoutesHandler {
	return &NewsRoutesHandler{repo}
}

func (handler *NewsRoutesHandler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slugOrID := ps.ByName("param")

	newsID, err := strconv.Atoi(slugOrID)
	// if error it means that param is a slug
	if err != nil {
		news, err2 := handler.repository.GetBySlug(slugOrID)
		if err2 != nil {
			response.Error(w, http.StatusNotFound, err2, err2.Error())
			return
		}
		response.JSON(w, http.StatusOK, news)
		return
	}

	news, err := handler.repository.GetByID(newsID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, news)
}

func (handler *NewsRoutesHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	status := queryValues.Get("status")

	if status == "draf" || status == "publish" {
		news, err := handler.repository.GetByStatus(status)
		if err != nil {
			response.Error(w, http.StatusNotFound, err, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, news)
		return
	}

	news, err := handler.repository.GetAll()
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, news)
	return
}

func (handler *NewsRoutesHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var n domain.News
	if err := helpers.FromJSON(&n, r.Body); err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
	}

	if err := handler.repository.Save(&n); err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, nil)
}

func (handler *NewsRoutesHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var n domain.News
	if err := helpers.FromJSON(&n, r.Body); err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	newsID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	if err = handler.repository.Update(newsID, n); err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, nil)
}

func (handler *NewsRoutesHandler) Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	newsID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	if err = handler.repository.Remove(newsID); err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, nil)
}
