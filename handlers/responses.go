package handlers

import (
	"encoding/json"
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"log"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}) {
	body, err := json.MarshalIndent(data, "", "	")
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
	helpers.HandlerError(err)
}

func Forrbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write([]byte("Forbidden"))
	helpers.HandlerError(err)
}

func ErrorResponse(w http.ResponseWriter, message string) {
	JSONResponse(w, common.H{
		"error":   true,
		"message": message,
	})
}
