package main

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers"
	"github.com/alex-pro27/monitoring_price_api/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func main()  {
	config.Load()
	databases.MigrateDefaultDB()
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/ping", handlers.Ping).Methods("GET")

	api.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")
	api.HandleFunc("/user", handlers.CreateUser).Methods("POST")
	api.HandleFunc("/user/{id}", handlers.UpdateUser).Methods("PUT")
	api.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")
	api.HandleFunc("/users", handlers.AllUsers).Methods("GET")

	router.Use(middlewares.DefaultDBMiddleware)
	systemConf := config.Config.System
	fmt.Println("Server started", systemConf.Server)
	server := http.Server{
		Addr: systemConf.Server,
		Handler: router,
	}

	common.HandlerError(server.ListenAndServe())
}