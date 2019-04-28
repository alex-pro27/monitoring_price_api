package middleware

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/utils"
	"github.com/gorilla/context"
	"net/http"
)

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Logger.Errorf("IP:%s - %s: %s%s - %v", utils.GetIPAddress(r), r.Method, r.Host, r.URL.Path, rec)
				if config.Config.System.Debug {
					panic(rec)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write([]byte("500 Internal server error"))
					logger.HandleError(err)
				}
			}
		}()
		h.ServeHTTP(w, r)
		who := fmt.Sprintf("IP:%s", utils.GetIPAddress(r))
		if user := context.Get(r, "user"); user != nil {
			who = fmt.Sprintf("%s - %s", who, user.(*models.User).String())
		}
		logger.Logger.Infof("%s - %s: %s%s", who, r.Method, r.Host, r.URL.Path)
	})
}
