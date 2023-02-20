package users

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/storage"
	"github.com/johnnrails/ddd_go/monitoring/opentelemetry/with-jaeger-2/trace"
	// oteltrace "go.opentelemetry.io/otel/trace"
)

type Controller struct {
	service Service
}

func NewController(storage storage.UserStorer) Controller {
	return Controller{
		service: Service{
			storage,
		},
	}
}

func (c Controller) Execute(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.NewSpan(r.Context(), "Controller.Create", nil)
	defer span.End()

	trace.AddSpanTags(span, map[string]string{
		"line_code":"31", 
		"app.tag_2":"val_2",
	})

	trace.AddSpanEvents(span, "event_test", map[string]string{
		"event_1":"val_1",
		"event_2":"val_2",
	})

	req := &Request{}
	if err := req.validate(ctx, r.Body); err != nil {
		span.RecordError(err)
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.service.Execute(ctx, req); err != nil {
		span.RecordError(err)
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postBody, _ := json.Marshal(map[string]string{
		"name":  req.Name,
 })
 responseBody := bytes.NewBuffer(postBody)
	w.Write(responseBody.Bytes())
	w.WriteHeader(http.StatusCreated)
}