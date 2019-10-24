package middleware

import (
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/gorilla/context"
	"net/http"
)

func DefaultDBMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := databases.ConnectDefaultDB()
		defer func() {
			logger.HandleError(db.Close())
		}()
		context.Set(r, "DB", db)
		h.ServeHTTP(w, r)
	})
}
