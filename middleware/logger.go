package middleware

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/utils"
	"net/http"
)

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Info(fmt.Sprintf("IP:%s - %s: %s%s", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path))
		h.ServeHTTP(w, r)
	})
}
