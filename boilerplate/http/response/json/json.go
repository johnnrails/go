package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	httperrors "github.com/johnnrails/ddd_go/boilerplate/http/errors"
	"github.com/johnnrails/ddd_go/boilerplate/http/response"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := hf(w, r); err != nil {
		code := http.StatusInternalServerError
		err = JSONError(r.Context(), w, code, err)
		if err != nil {
			panic(err)
		}
	}
}

func JSON(ctx context.Context, w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if payload == nil {
		_, err := w.Write([]byte("{}"))
		return err
	}

	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	encoder.SetIndent("", "")

	if err := encoder.Encode(payload); err != nil {
		return err
	}

	response.Flush(w)
	return nil
}

func JSONError(ctx context.Context, w http.ResponseWriter, code int, err error) error {
	if err == nil {
		return errors.New("JSONError called with err nil")
	}

	httperr := httperrors.NewHttpError(ctx, code, err)

	if err := JSON(ctx, w, httperr.Code, httperr.Message); err != nil {
		return err
	}

	return nil
}

func NotFound() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		httpError := &httperrors.HttpError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Route %s %s", r.URL.Path, http.StatusText(http.StatusNotFound)),
		}
		return JSON(r.Context(), w, httpError.Code, httpError)
	}
	return HandlerFunc(fn)
}

func NowAllowed() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		httpError := &httperrors.HttpError{
			Code:    http.StatusMethodNotAllowed,
			Message: fmt.Sprintf("Route %s %s", r.URL.Path, http.StatusText(http.StatusNotFound)),
		}
		return JSON(r.Context(), w, httpError.Code, httpError)
	}
	return HandlerFunc(fn)
}
