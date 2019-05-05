package common

import (
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/alex-pro27/monitoring_price_api/utils"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}) {
	var (
		body []byte
		err  error
	)
	if config.Config.System.Debug {
		body, err = json.MarshalIndent(data, "", "	")
	} else {
		body, err = json.Marshal(data)
	}

	if err != nil {
		log.Printf("Failed to encode a JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v", err)
		return
	}
}

func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Page not found"))
	logger.Logger.Warningf("404 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	logger.HandleError(err)
}

func Error405(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte("Method not allowed"))
	logger.Logger.Warningf("405 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	logger.HandleError(err)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	logger.Logger.Warningf("403 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	_, err := w.Write([]byte("Forbidden"))
	logger.HandleError(err)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	logger.Logger.Warningf("401 - IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path)
	if message == "" {
		message = "Unauthorized"
	}
	_, err := w.Write([]byte(message))
	logger.HandleError(err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.Logger.Warningf("IP:%s - %s: %s%s - %s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path, message)
	w.WriteHeader(http.StatusOK)
	JSONResponse(w, types.H{
		"error":   true,
		"message": message,
	})
}
