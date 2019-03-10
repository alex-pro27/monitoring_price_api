package main

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	config.Load()
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

	if config.Config.System.Debug {
		router.HandleFunc("/ping", handlers.Ping).Methods("GET")
	}

	useSessionAuth.HandleFunc("/api/check-auth/", handlers.CheckAuth).Methods("POST")

	useDB.HandleFunc("/api/admin/login", handlers.Auth).Methods("POST")
	useDB.HandleFunc("/api/admin/logout", handlers.Logout).Methods("GET")

	useSessionAuth.HandleFunc("/api/admin/user/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	useSessionAuth.HandleFunc("/api/admin/user", handlers.CreateUser).Methods("PUT")
	useSessionAuth.HandleFunc("/api/admin/user/{id:[0-9]+}", handlers.UpdateUser).Methods("POST")
	useSessionAuth.HandleFunc("/api/admin/user/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")
	useSessionAuth.HandleFunc("/api/admin/users", handlers.AllUsers).Methods("GET")

	useTokenAuth.HandleFunc("/api/monitoring-shops", handlers.GetMonitoringShops).Methods("GET")
	useTokenAuth.HandleFunc("/api/segments", handlers.GetSegments).Methods("GET")

	systemConf := config.Config.System
	fmt.Println("Server started", systemConf.Server)
	server := http.Server{
		Addr:    systemConf.Server,
		Handler: muxHandlers.CORS()(router),
	}
	helpers.HandlerError(server.ListenAndServe())
}
