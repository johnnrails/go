package json

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	type jsonResponse struct {
		Name string `json:"name"`
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := JSON(r.Context(), w, http.StatusOK, jsonResponse{"john"}); err != nil {
			t.Fatal(err)
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, req)
	header := w.Header()

	if header.Get("Content-Type") != "application/json" {
		t.Fatal("application/json not set in Content-Type")
	}
	cmp := bytes.Compare(w.Body.Bytes(), append([]byte(`{"name":"john"}`), 10))
	if cmp != 0 {
		t.Errorf("JSON Returned wrong body: %s | %d", w.Body.String(), cmp)
	}
}
