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

type AdminMetaField struct {
	Name  string
	Label string
	Type  string
}

type AdminMeta struct {
	ExcludeFields []string
	Fields        []string
	Preload       []string
	SortFields    []string
	OrderBy       []string
	SearchFields  []string
	ExtraFields   []AdminMetaField
	FilterFields  []AdminMetaField
}

type Model interface {
	Meta() ModelsMeta
	Admin() AdminMeta
	CRUD() *CRUDManager
	String() string
}
