package routes

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/api/v2"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/gorilla/mux"
)

func RegisterApiV2Routes(r *mux.Router) {
	router := r.NewRoute().PathPrefix("/api/v2").Subrouter()
	router.Use(middleware.DefaultDBMiddleware)
	routerTokenAuth := router.NewRoute().Subrouter()
	routerTokenAuth.Use(middleware.TokenAuthMiddleware)
	apiRoutes := []Route{
		{
			Path:    "/monitoring-shops",
			Handler: v2.GetMonitoringShops,
			Methods: []string{"GET"},
		},
		{
			Path:    "/segments",
			Handler: v2.GetSegments,
			Methods: []string{"GET"},
		},
	}
	RegisterRoutes(routerTokenAuth, apiRoutes, nil)
}
