package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/johnnrails/ddd_go/second_ddd_go/application"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain/repositories"
	"github.com/johnnrails/ddd_go/second_ddd_go/response"
	"github.com/julienschmidt/httprouter"
)

type TopicRoutesHandler struct {
	application *application.TopicApplication
}

func CreateTopicRoutesHandler(repo repositories.TopicRepository) *TopicRoutesHandler {
	a := application.CreateTopicApplication(repo)
	return &TopicRoutesHandler{
		application: a,
	}
}

func (handler *TopicRoutesHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	topics, err := handler.application.GetAll()
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
	topic, err := handler.application.Get(topicID)
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
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
		return
	}
	if err := handler.application.Add(p.Name, p.Slug); err != nil {
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
	if err = handler.application.Remove(topicID); err != nil {
		response.Error(w, http.StatusInternalServerError, err, err.Error())
	}
	response.JSON(w, http.StatusOK, nil)
}
func (handler *TopicRoutesHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
	}
	topicID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err, err.Error())
	}
	if err = handler.application.Update(t, topicID); err != nil {
		response.Error(w, http.StatusInternalServerError, err, err.Error())
	}
	response.JSON(w, http.StatusOK, nil)
}
