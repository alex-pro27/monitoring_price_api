package types

import (
	"net/http"
)

type H map[string]interface{}

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

type CRUDManager interface {
	Create(fields H) error
	Update(fields H) error
	Delete() error
}

type ModelsMeta struct {
	Name   string
	Plural string
}

type ExtraField struct {
	Name  string
	Label string
}

type AdminMeta struct {
	ExcludeFields []string
	Fields        []string
	Preload       []string
	SortFields    []string
	OrderBy       []string
	SearchFields  []string
	ExtraFields   []ExtraField
}

type Model interface {
	Meta() ModelsMeta
	Admin() AdminMeta
	CRUD() *CRUDManager
	String() string
}
