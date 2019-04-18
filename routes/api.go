package routes

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/api"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/gorilla/mux"
)

func RegisterApiRoutes(r *mux.Router) {
	router := r.NewRoute().PathPrefix("/api").Subrouter()
	router.Use(middleware.DefaultDBMiddleware)
	routerTokenAuth := router.NewRoute().Subrouter()
	routerTokenAuth.Use(middleware.TokenAuthMiddleware)
	apiRoutes := []Route{
		{
			Path:    "/monitoring-shops",
			Handler: api.GetMonitoringShops,
			Methods: []string{"GET"},
		},
		{
			Path:    "/segments",
			Handler: api.GetSegments,
			Methods: []string{"GET"},
		},
	}
	RegisterRoutes(routerTokenAuth, apiRoutes, nil)
}
