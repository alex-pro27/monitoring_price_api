package types

import (
	"net/http"
)

type H map[string]interface{}

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

type CRUDManager interface {
	Create(fields H) error
	Update(fields H) error
	Delete(fields H) error
}

type ModelsMeta struct {
	Name   string
	Plural string
}

type AdminMeta struct {
	ExcludeFields []string
	Fields        []string
}

type Model interface {
	Meta() ModelsMeta
	Admin() AdminMeta
	CRUD() *CRUDManager
	String() string
}
