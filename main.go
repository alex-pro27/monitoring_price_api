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
	"os"
	"os/signal"
	"runtime"
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
	runtime.GC()
	defer func() {
		fmt.Println("Server closed")
		//memProfile, _ := os.Create(config.Config.System.MemProfiler)
		//logger.HandleError(pprof.WriteHeapProfile(memProfile))
		//logger.HandleError(memProfile.Close())
		if databases.DB != nil {
			logger.HandleError(databases.DB.Close())
		}
	}()
	logger.Logger.Infof("Server started: %s", config.Config.System.Server)
	logger.HandleError(Server.ListenAndServe())
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = Server.Close()
	}()
	Init()
	StartServer()
}
