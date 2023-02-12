package handlers

import (
	"net/http"
	"strconv"

	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain/repositories"
	"github.com/johnnrails/ddd_go/second_ddd_go/helpers"
	"github.com/johnnrails/ddd_go/second_ddd_go/response"
	"github.com/julienschmidt/httprouter"
)

type TopicRoutesHandler struct {
	repository repositories.TopicRepository
}

func CreateTopicRoutesHandler(repo repositories.TopicRepository) *TopicRoutesHandler {
	return &TopicRoutesHandler{repo}
}

func (handler *TopicRoutesHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	topics, err := handler.repository.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, topics)
}

func (handler *TopicRoutesHandler) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	topicID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	topic, err := handler.repository.Get(topicID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, topic)
}

func (handler *TopicRoutesHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type payload struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	var p payload

	if err := helpers.FromJSON(&p, r.Body); err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	if err := handler.repository.Save(&domain.Topic{
		Name: p.Name,
		Slug: p.Slug,
	}); err != nil {
		response.Error(w, http.StatusConflict, err, err.Error())
	}
}

func (handler *TopicRoutesHandler) Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	topicID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
		return
	}

	// it's impossible to handle all the errors with this, buts works now.
	if err = handler.repository.Remove(topicID); err != nil {
		response.Error(w, http.StatusInternalServerError, err, err.Error())
	}

	response.JSON(w, http.StatusOK, nil)
}

func (handler *TopicRoutesHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t domain.Topic
	if err := helpers.FromJSON(&t, r.Body); err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
	}

	topicID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
	}

	t.ID = uint(topicID)

	if err = handler.repository.Update(&t); err != nil {
		response.Error(w, http.StatusInternalServerError, err, err.Error())
	}

	response.JSON(w, http.StatusOK, nil)
}
