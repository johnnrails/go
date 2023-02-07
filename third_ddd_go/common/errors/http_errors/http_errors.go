package http_errors

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/johnnrails/ddd_go/third_ddd_go/common/errors"
)

type ErrResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}

func logError() {
	// do something to log this error later.
}

func basicResponseWithErr(w http.ResponseWriter, r *http.Request, slug string, status int) {
	resp := ErrResponse{slug, status}
	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

func InternalServerError(slug string, w http.ResponseWriter, r *http.Request) {
	basicResponseWithErr(w, r, slug, http.StatusInternalServerError)
}

func Unauthorized(slug string, w http.ResponseWriter, r *http.Request) {
	basicResponseWithErr(w, r, slug, http.StatusUnauthorized)
}

func BadRequest(slug string, w http.ResponseWriter, r *http.Request) {
	basicResponseWithErr(w, r, slug, http.StatusBadRequest)
}

func RespondWithSlugError(w http.ResponseWriter, r *http.Request, err error) {
	slugErr, ok := err.(errors.SlugError)
	if !ok {
		InternalServerError("unable-to-create-slug-error", w, r)
	}

	switch slugErr.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorized("unauthorized", w, r)
	case errors.ErrorTypeIncorrectInput:
		BadRequest("incorrect-input", w, r)
	default:
		InternalServerError("internal-server-error", w, r)
	}
}
