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
	Name   string
	Label  string
	ToHTML string // date, datetime, image
}

type AdminMeta struct {
	ShortToHtml   string // date, datetime, image
	ExcludeFields []string
	Fields        []string
	Preload       []string
	SortFields    []AdminMetaField
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
