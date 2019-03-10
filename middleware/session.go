package middleware

import (
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"net/http"
)

func SessionsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var store = sessions.NewCookieStore([]byte(config.Config.Session.Key))
		store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   config.Config.Session.MaxAge,
			HttpOnly: true,
		}
		context.Set(r, "sessions", store)
		h.ServeHTTP(w, r)
	})
}
