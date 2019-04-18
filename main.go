package main

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/routes"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	config.Load()
	logger.Init()
	databases.MigrateDefaultDB()
	systemConf := &config.Config.System
	//CORSMiddleware := handlers.CORS()
	router := mux.NewRouter()
	routes.RegisterAdminRoutes(router)
	routes.RegisterApiRoutes(router)
	fmt.Println("Server started", systemConf.Server)
	defer logger.Close()
	logger.HandleError(http.ListenAndServe(systemConf.Server, router))
}
