package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/johnnrails/ddd_go/second_ddd_go/application"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/infra/persistence"
	"github.com/johnnrails/ddd_go/second_ddd_go/response"
	"github.com/julienschmidt/httprouter"
)

type NewsRoutesHandler struct {
	application *application.NewsApplication
}

func CreateNewsRoutesHandler() *NewsRoutesHandler {
	nr, _ := persistence.CreateNewsRepository()
	na := application.CreateNewsApplication(nr)
	return &NewsRoutesHandler{
		application: na,
	}
}

func (handler *NewsRoutesHandler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slugOrID := ps.ByName("param")

	newsID, err := strconv.Atoi(slugOrID)
	// if error it means that param is a slug
	if err != nil {
		news, err2 := handler.application.GetBySlug(slugOrID)
		if err2 != nil {
			response.Error(w, http.StatusNotFound, err2, err2.Error())
			return
		}
		response.JSON(w, http.StatusOK, news)
		return
	}

	news, err := handler.application.GetByID(newsID)
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
		news, err := handler.application.GetByStatus(status)
		if err != nil {
			response.Error(w, http.StatusNotFound, err, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, news)
		return
	}

	limit, err1 := strconv.Atoi(queryValues.Get("limit"))
	page, err2 := strconv.Atoi(queryValues.Get("page"))

	if err1 != nil || err2 != nil {
		response.Error(w, http.StatusBadRequest, errors.New("could not parse limit or page query values"), "could not parse limit or page query values")
	}

	if limit != 0 && page != 0 {
		news, err := handler.application.GetAll(limit, page)
		if err != nil {
			response.Error(w, http.StatusNotFound, err, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, news)
		return
	}
}

func (handler *NewsRoutesHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var n domain.News
	if err := decoder.Decode(&n); err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
	}
	err := handler.application.Add(n)
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, nil)
}

func (handler *NewsRoutesHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var n domain.News
	err := decoder.Decode(&n)
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	newsID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	err = handler.application.Update(newsID, n)
	if err != nil {
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
	err = handler.application.Remove(newsID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, nil)
}
