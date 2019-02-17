package middlewares

import (
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/gorilla/context"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

func DefaultDBMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.URL)
		db := databases.ConnectDefaultDB()
		//defer common.HandlerError(db.Close())
		context.Set(r, "DB", db)
		h.ServeHTTP(w, r)
	})
}
