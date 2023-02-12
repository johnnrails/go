package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func printDebugf(format string, args ...interface{}) {
	if env := os.Getenv("GO_SERVER_DEBUG"); len(env) != 0 {
		log.Printf("[DEBUG]"+format+"\n", args...)
	}
}

type ErrorResponse struct {
	Message string `json:"reason"`
	Error   error  `json:"-"`
}

func (e *ErrorResponse) String() string {
	return fmt.Sprintf("reason: %s, error: %s", e.Message, e.Error.Error())
}

func Respond(w http.ResponseWriter, code int, src interface{}) {
	var body []byte
	var err error

	switch s := src.(type) {
	// the type is []byte when the src comes as json
	case []byte:
		if !json.Valid(s) {
			Error(w, http.StatusInternalServerError, err, "invalid json")
			return
		}
		body = s
	case string:
		body = []byte(s)
	case *ErrorResponse, ErrorResponse:
		if body, err = json.Marshal(src); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"reason\":\"failed to parse json\"}"))
			return
		}
	default:
		if body, err = json.Marshal(src); err != nil {
			Error(w, http.StatusInternalServerError, err, "failed to parse json")
			return
		}
	}
	w.WriteHeader(code)
	w.Write(body)
}

func Error(w http.ResponseWriter, code int, err error, msg string) {
	e := &ErrorResponse{
		Message: msg,
		Error:   err,
	}
	printDebugf("%s", e.String())
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	Respond(w, code, e)
}

func JSON(w http.ResponseWriter, code int, src interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	Respond(w, code, src)
}
