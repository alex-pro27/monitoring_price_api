package handlers

import (
	. "github.com/alex-pro27/monitoring_price_api/common"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, H{
		"message": "PONG",
	})
}
