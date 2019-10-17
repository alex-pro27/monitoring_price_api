package main

import (
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/middleware"
	"github.com/alex-pro27/monitoring_price_api/routes"
	"github.com/gorilla/mux"
	"github.com/xlab/closer"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
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
	logger.Logger.Infof("Server started: %s", config.Config.System.Server)
	logger.HandleError(Server.ListenAndServe())
	closer.Close()
}

func CloseServer() {
	logger.HandleError(Server.Close())
	memProfile, _ := os.Create(config.Config.System.MemProfiler)
	logger.HandleError(pprof.WriteHeapProfile(memProfile))
	defer logger.HandleError(memProfile.Close())
}

func main() {
	Init()
	runtime.GC()
	closer.Bind(CloseServer)
	go StartServer()
	closer.Hold()
}
