package middleware

import (
	"encoding/base64"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

/**
Basic Авторизация
*/
func BasicAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			authData := strings.Split(auth, " ")
			if len(authData) == 2 && authData[0] == "Basic" {
				data, err := base64.StdEncoding.DecodeString(authData[1])
				if err == nil {
					logpassw := strings.Split(string(data), ":")
					user := models.User{}
					db := context.Get(r, "DB").(*gorm.DB)
					user.GetByUserName(db, logpassw[0])
					if user.CheckPassword(logpassw[1]) {
						context.Set(r, "user", &user)
						h.ServeHTTP(w, r)
						return
					}
				}
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Not authorized, forbidden"))
		helpers.HandlerError(err)
	})
}

/**
Авторизация Токену
*/
func TokenAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			authData := strings.Split(auth, " ")
			if len(authData) == 2 && authData[0] == "Token" {
				user := models.User{}
				db := context.Get(r, "DB").(*gorm.DB)
				user.GetUserByToken(db, authData[1])
				if user.ID > 0 {
					context.Set(r, "user", &user)
					h.ServeHTTP(w, r)
					return
				}
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid token"))
		helpers.HandlerError(err)
	})
}

/**
Авторизация по сессии
*/
func SessionAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := context.Get(r, "sessions")
		if store != nil {
			store := store.(*sessions.CookieStore)
			session, _ := store.Get(r, "user")
			userID := session.Values["user_id"]
			if userID != nil {
				user := models.User{}
				db := context.Get(r, "DB").(*gorm.DB)
				user.GetById(db, userID.(uint))
				if user.ID != 0 {
					context.Set(r, "user", &user)
					h.ServeHTTP(w, r)
					return
				}
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Not authorized, forbidden"))
		helpers.HandlerError(err)
	})
}
