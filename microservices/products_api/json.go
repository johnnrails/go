package main

import (
	"encoding/json"
	"io"
)

func ToJSON(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}
