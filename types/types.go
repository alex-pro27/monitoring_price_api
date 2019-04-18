package types

import (
	"net/http"
)

type H map[string]interface{}

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

type ModelsMeta struct {
	Name string
}
