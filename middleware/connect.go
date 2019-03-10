package middleware

import (
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/gorilla/context"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func DefaultDBMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.URL)
		db := databases.ConnectDefaultDB()
		context.Set(r, "DB", db)
		h.ServeHTTP(w, r)
		defer helpers.HandlerError(db.Close())
	})
}
