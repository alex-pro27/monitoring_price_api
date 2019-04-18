package routes

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/admin"
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

	notAuthRoutes := []Route{
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

	usersRoutes := []Route{
		{
			Path:    "/users",
			Handler: admin.AllUsers,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
		{
			Path:    "/create-user",
			Handler: admin.CreateUser,
			Access:  models.WRITE,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/user/{id:[0-9]+}",
			Handler: admin.GetUser,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
		{
			Path:    "/user/{id:[0-9]+}",
			Handler: admin.UpdateUser,
			Access:  models.WRITE,
			Methods: []string{"POST"},
		},
		{
			Path:    "/user/{id:[0-9]+}",
			Handler: admin.DeleteUser,
			Access:  models.ACCESS,
			Methods: []string{"DELETE"},
		},
	}

	rolesRoutes := []Route{
		{
			Path:    "/roles",
			Handler: admin.AllRoles,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
		{
			Path:    "/create-role",
			Handler: admin.CreateRole,
			Access:  models.WRITE,
			Methods: []string{"PUT"},
		},
	}

	permissionsRoutes := []Route{
		{
			Path:    "/permissions",
			Handler: admin.GetPermissions,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
	}

	viewsRoutes := []Route{
		{
			Path:    "/views",
			Handler: admin.AllViews,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
		{
			Path:    "/create-view",
			Handler: admin.CreateView,
			Access:  models.WRITE,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/view/{id:[0-9]+}",
			Handler: admin.GetView,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
	}

	workGroupsRoutes := []Route{
		{
			Path:    "/work-groups",
			Handler: admin.AllWorkGroups,
			Access:  models.READ,
			Methods: []string{"GET"},
		},
	}

	noCheckPermissionsRoutes := []Route{
		{
			Path:    "/available-views",
			Handler: admin.GetAvailableViews,
			Methods: []string{"GET"},
		},
	}

	RegisterRoutes(router, notAuthRoutes, nil)
	RegisterRoutes(authRouter, usersRoutes, models.User{})
	RegisterRoutes(authRouter, rolesRoutes, models.Role{})
	RegisterRoutes(authRouter, permissionsRoutes, models.Permission{})
	RegisterRoutes(authRouter, viewsRoutes, models.Views{})
	RegisterRoutes(authRouter, workGroupsRoutes, models.WorkGroup{})
	RegisterRoutes(authRouter, noCheckPermissionsRoutes, nil)
}
