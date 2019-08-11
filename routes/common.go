package routes

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/admin"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Path    string
	Methods []string
	Handler types.HTTPHandler
	Access  models.PermissionAccess
}

func RegisterRoutes(r *mux.Router, routes []Route, model interface{}) {
	for _, route := range routes {
		handler := route.Handler
		if model != nil {
			handler = CheckPermissionDecorator(route.Access, model)(handler)
		}
		r.HandleFunc(route.Path, handler).Methods(route.Methods...)
	}
}

func CheckPermissionDecorator(
	access models.PermissionAccess,
	model interface{},
) func(handler types.HTTPHandler) types.HTTPHandler {
	return func(handler types.HTTPHandler) types.HTTPHandler {
		return func(w http.ResponseWriter, r *http.Request) {
			if admin.CheckPermission(w, r, access, model) {
				handler(w, r)
			}
		}
	}
}
