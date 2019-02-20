package common

import (
	"encoding/json"
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

func Error404(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Page not found"))
	HandlerError(err)
}

func ErrorResponse(w http.ResponseWriter, message string) {
	JSONResponse(w, H{
		"error": true,
		"message": message,
	})
}
