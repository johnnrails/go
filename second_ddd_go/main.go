package main

import (
	"net/http"

	"github.com/johnnrails/ddd_go/second_ddd_go/handlers"
	"github.com/johnnrails/ddd_go/second_ddd_go/response"
	"github.com/julienschmidt/httprouter"
)

func Routes() *httprouter.Router {
	r := httprouter.New()

	r.GET("/", index)

	newsHandler := handlers.CreateNewsRoutesHandler()
	r.GET("/api/news", newsHandler.GetAll)
	r.GET("/api/news/:param", newsHandler.Get)
	r.POST("/api/news", newsHandler.Create)
	r.DELETE("/api/news/:id", newsHandler.Remove)
	r.PUT("/api/news/:id", newsHandler.Update)

	topicHandler := handlers.CreateTopicRoutesHandler()
	r.GET("/api/topic", topicHandler.GetAll)
	r.GET("/api/topic/:id", topicHandler.Get)
	r.POST("/api/topic", topicHandler.Create)
	r.DELETE("/api/topic/:id", topicHandler.Remove)
	r.PUT("/api/topic/:id", topicHandler.Update)

	return r
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response.JSON(w, http.StatusOK, "GO API")
}

func main() {
	http.ListenAndServe(":8080", Routes())
}
