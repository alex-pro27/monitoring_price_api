package main

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/alex-pro27/monitoring_price_api/routes"
	"github.com/gorilla/mux"
	"net/http"
)

var Server *http.Server

func Init() {
	config.Load()
	logger.Init()
	databases.MigrateDefaultDB()
	router := mux.NewRouter()
	router.Use(middleware.LoggerMiddleware)
	routes.RegisterAdminRoutes(router)
	routes.RegisterApiV1Routes(router)
	routes.RegisterApiV2Routes(router)
	router.NotFoundHandler = http.HandlerFunc(common.Error404)
	router.MethodNotAllowedHandler = http.HandlerFunc(common.Error405)
	Server = &http.Server{
		Addr:    config.Config.System.Server,
		Handler: router,
	}
}

func StartServer() {
	fmt.Println("Server started", config.Config.System.Server)
	defer logger.Close()
	logger.HandleError(Server.ListenAndServe())
}

func main() {
	Init()
	StartServer()
}
