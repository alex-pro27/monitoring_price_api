package routes

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/admin"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(r *mux.Router) {
	prefix := "/api/admin"
	router := r.NewRoute().PathPrefix(prefix).Subrouter()
	router.Use(
		middleware.SessionsMiddleware,
		middleware.DefaultDBMiddleware,
	)
	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(middleware.SessionAuthMiddleware)

	noCheckAuthRoutes := []Route{
		{
			Path:    "/login",
			Handler: admin.Login,
			Methods: []string{"POST"},
		},
		{
			Path:    "/logout",
			Handler: admin.Logout,
			Methods: []string{"GET", "POST"},
		},
		{
			Path:    "/check-auth",
			Handler: admin.CheckAuth,
			Methods: []string{"GET", "POST"},
		},
	}

	contentTypesRoutes := []Route{
		{
			Path:    "/available-views",
			Handler: admin.GetAvailableViews,
			Methods: []string{"GET"},
		},
		{
			Path:    "/content-types",
			Handler: admin.AllFieldsInModel,
			Methods: []string{"GET"},
		},
		{
			Path:    "/content-type/{id:[0-9]+}",
			Handler: admin.GetContentTypeFields,
			Methods: []string{"GET"},
		},
		{
			Path:    "/content-type/create",
			Handler: admin.CRUDContentType,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/content-type/{id:[0-9]+}",
			Handler: admin.CRUDContentType,
			Methods: []string{"POST", "DELETE"},
		},
		{
			Path:    "/media/{name}",
			Handler: common.FileResponse,
			Methods: []string{"GET"},
		},
	}

	usersRoutes := []Route{
		{
			Path:    "/users",
			Handler: admin.AllUsers,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
	}

	RegisterRoutes(router, noCheckAuthRoutes, nil)
	RegisterRoutes(authRouter, usersRoutes, models.User{})
	RegisterRoutes(authRouter, contentTypesRoutes, nil)
}
