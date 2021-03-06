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
	authRouter.Use(middleware.AuthMiddleware(middleware.SESSION_AUTH | middleware.TOKEN_AUTH))

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
	}

	contentTypesRoutes := []Route{
		{
			Path:    "/check-auth",
			Handler: admin.CheckAuth,
			Methods: []string{"GET", "POST"},
		},
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
			Path:    "/content-types/filter",
			Handler: admin.FilteredContentType,
			Methods: []string{"GET"},
		},
		{
			Path:    "/content-type/{id:[0-9]+}",
			Handler: admin.GetContentTypeFields,
			Methods: []string{"GET"},
		},
		{
			Path:    "/content-type/{content_type_id:[0-9]+}/create",
			Handler: admin.CRUDContentType,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/content-type/{content_type_id:[0-9]+}/{id:[0-9]+}",
			Handler: admin.CRUDContentType,
			Methods: []string{"POST", "DELETE"},
		},
		{
			Path:    "/trash-data",
			Handler: admin.GetTrashData,
			Methods: []string{"GET"},
		},
		{
			Path:    "/recovery-from-trash",
			Handler: admin.RecoveryFromTrash,
			Methods: []string{"POST"},
		},
		{
			Path:    "/media/{name}",
			Handler: common.FileResponse,
			Methods: []string{"GET"},
		},
	}

	monitorinRoutes := []Route{
		{
			Path:    "/monitorings",
			Handler: admin.GetAllMonitoringList,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
		{
			Path:    "/update-monitorings",
			Handler: admin.UpdateMonitorings,
			Access:  models.WRITE,
			Methods: []string{"POST"},
		},
		{
			Path:    "/product-template-file",
			Handler: admin.GetTemplateBlank,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
	}

	completeWaresRoutes := []Route{
		{
			Path:    "/complete-wares",
			Handler: admin.GetCompletedWares,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
		{
			Path:    "/get-report",
			Handler: admin.GenerateReport,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
	}

	dashboardRoutes := []Route {
		{
			Path: 		"/online-users",
			Handler: 	admin.GetOnlineUsers,
			Methods:	[]string{"GET"},
		},
	}

	websockets := []Route {
		{
			Path:    "/ws",
			Handler: admin.HandleWebsocket,
			Methods: []string{"GET", "POST"},
		},
	}

	RegisterRoutes(router, noCheckAuthRoutes, nil)
	RegisterRoutes(router, websockets, nil)
	RegisterRoutes(authRouter, contentTypesRoutes, nil)
	RegisterRoutes(authRouter, monitorinRoutes, models.Monitoring{})
	RegisterRoutes(authRouter, completeWaresRoutes, models.CompletedWare{})
	RegisterRoutes(authRouter, dashboardRoutes, nil)
}
