package main

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	config.Load()
	databases.MigrateDefaultDB()

	router := mux.NewRouter()

	useDB := router.NewRoute().Subrouter()
	useTokenAuth := useDB.NewRoute().Subrouter()
	useBasicAuth := useDB.NewRoute().Subrouter()

	useDB.Use(middleware.DefaultDBMiddleware)
	useTokenAuth.Use(middleware.TokenAuthMiddleware)
	useBasicAuth.Use(middleware.BasicAuthMiddleware)

	if config.Config.System.Debug {
		router.HandleFunc("/ping", handlers.Ping).Methods("GET")
	}

	useTokenAuth.HandleFunc("/api/user/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	useDB.HandleFunc("/api/user", handlers.CreateUser).Methods("PUT")
	useTokenAuth.HandleFunc("/api/user/{id:[0-9]+}", handlers.UpdateUser).Methods("POST")
	useTokenAuth.HandleFunc("/api/user/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")
	useDB.HandleFunc("/api/users", handlers.AllUsers).Methods("GET")

	useTokenAuth.HandleFunc("/api/monitoring-shops", handlers.GetMonitoringShops).Methods("GET")
	useTokenAuth.HandleFunc("/api/segments", handlers.GetSegments).Methods("GET")

	systemConf := config.Config.System
	fmt.Println("Server started", systemConf.Server)
	server := http.Server{
		Addr:    systemConf.Server,
		Handler: muxHandlers.CORS()(router),
	}
	common.HandlerError(server.ListenAndServe())
}
