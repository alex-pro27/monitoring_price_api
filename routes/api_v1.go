package routes

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/api/v1"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/gorilla/mux"
)

func RegisterApiV1Routes(r *mux.Router) {
	router := r.NewRoute().PathPrefix("/api/v1").Subrouter()
	router.Use(middleware.DefaultDBMiddleware)
	routerTokenAuth := router.NewRoute().Subrouter()
	routerTokenAuth.Use(middleware.TokenAuthMiddleware)
	routerBasicAuth := router.NewRoute().Subrouter()
	routerBasicAuth.Use(middleware.BasicAuthMiddleware)
	notAuthRoutes := []Route{
		{
			Path:    "/user/{barcode}",
			Handler: v1.GetUser,
			Methods: []string{"GET"},
		},
		{
			Path:    "/check-pin",
			Handler: v1.CheckPin,
			Methods: []string{"POST"},
		},
	}
	apiRoutes := []Route{
		{
			Path:    "/rival/{region}/{shop}",
			Handler: v1.GetRivals,
			Methods: []string{"GET"},
		},
		{
			Path:    "/wares/{region}/{shop}",
			Handler: v1.GetWares,
			Methods: []string{"GET"},
		},
		{
			Path:    "/periods",
			Handler: v1.GetPeriods,
			Methods: []string{"GET"},
		},
		{
			Path:    "/unload-ware",
			Handler: v1.CompleteWare,
			Methods: []string{"POST"},
		},
		{
			Path:    "/unload-photo",
			Handler: v1.UploadPhoto,
			Methods: []string{"POST"},
		},
	}
	apiBasicAuth := []Route{
		{
			Path:    "/get-monitoring-data",
			Handler: v1.GetCompletedWares,
			Methods: []string{"GET"},
		},
		{
			Path:    "/media/{name}",
			Handler: common.FileResponse,
			Methods: []string{"GET"},
		},
		{
			Path:    "/monitoring-types",
			Handler: v1.GetPeriods,
			Methods: []string{"GET"},
		},
		{
			Path:    "/monitoring-shops",
			Handler: v1.GetMonitoringShops,
			Methods: []string{"GET"},
		},
		{
			Path:    "/regions",
			Handler: v1.GetRegions,
			Methods: []string{"GET"},
		},
	}
	RegisterRoutes(router, notAuthRoutes, nil)
	RegisterRoutes(routerTokenAuth, apiRoutes, nil)
	RegisterRoutes(routerBasicAuth, apiBasicAuth, nil)
}
