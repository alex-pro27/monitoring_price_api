package main

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/admin"
	"github.com/alex-pro27/monitoring_price_api/handlers/api"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	config.Load()
	logger.Init()

	databases.MigrateDefaultDB()

	router := mux.NewRouter()

	useSession := router.NewRoute().Subrouter()
	useDB := useSession.NewRoute().Subrouter()
	useTokenAuth := useDB.NewRoute().Subrouter()
	useBasicAuth := useDB.NewRoute().Subrouter()
	useSessionAuth := useDB.NewRoute().Subrouter()

	useSession.Use(middleware.SessionsMiddleware)
	useDB.Use(middleware.DefaultDBMiddleware)
	useTokenAuth.Use(middleware.TokenAuthMiddleware)
	useBasicAuth.Use(middleware.BasicAuthMiddleware)
	useSessionAuth.Use(middleware.SessionAuthMiddleware)

	useSessionAuth.HandleFunc("/api/check-auth", admin.CheckAuth).Methods("POST", "GET")

	useDB.HandleFunc("/api/admin/login", admin.Login).Methods("POST")
	useDB.HandleFunc("/api/admin/logout", admin.Logout).Methods("GET", "POST")

	useSessionAuth.HandleFunc("/api/admin/user/{id:[0-9]+}", admin.GetUser).Methods("GET")
	useSessionAuth.HandleFunc("/api/admin/user", admin.CreateUser).Methods("PUT")
	useSessionAuth.HandleFunc("/api/admin/user/{id:[0-9]+}", admin.UpdateUser).Methods("POST")
	useSessionAuth.HandleFunc("/api/admin/user/{id:[0-9]+}", admin.DeleteUser).Methods("DELETE")
	useSessionAuth.HandleFunc("/api/admin/users", admin.AllUsers).Methods("GET")

	useTokenAuth.HandleFunc("/api/monitoring-shops", api.GetMonitoringShops).Methods("GET")
	useTokenAuth.HandleFunc("/api/segments", api.GetSegments).Methods("GET")

	systemConf := config.Config.System
	fmt.Println("Server started", systemConf.Server)
	server := http.Server{
		Addr:    systemConf.Server,
		Handler: handlers.CORS()(router),
	}

	defer logger.Close()

	logger.HandleError(server.ListenAndServe())
}
