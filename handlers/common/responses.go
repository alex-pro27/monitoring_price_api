package common

import (
	"encoding/json"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/types"
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

func Error404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Page not found"))
	logger.HandleError(err)
}

func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write([]byte("Forbidden"))
	logger.HandleError(err)
}

func ErrorResponse(w http.ResponseWriter, message string) {
	JSONResponse(w, types.H{
		"error":   true,
		"message": message,
	})
}
