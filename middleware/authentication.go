package middleware

import (
	"encoding/base64"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/utils"
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
					db := context.Get(r, "DB").(*gorm.DB)
					userManager := models.UserManager{db}
					user := userManager.GetByUserName(logpassw[0])
					if user.CheckPassword(logpassw[1]) {
						context.Set(r, "user", &user)
						h.ServeHTTP(w, r)
						return
					}
				}
			}
		}
		logger.Logger.Warning(
			fmt.Sprintf("Not authorized, forbidden: IP: %s, url: %s", utils.GetIPAddress(r), r.RequestURI),
		)
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Not authorized, forbidden"))
		logger.HandleError(err)
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
				db := context.Get(r, "DB").(*gorm.DB)
				userManager := models.UserManager{db}
				user := userManager.GetUserByToken(authData[1])
				if user.ID > 0 {
					context.Set(r, "user", &user)
					h.ServeHTTP(w, r)
					return
				}
			}
		}
		logger.Logger.Warning(
			fmt.Sprintf("Invalid token: IP: %s, url: %s", utils.GetIPAddress(r), r.RequestURI),
		)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid token"))
		logger.HandleError(err)
	})
}

/**
Авторизация по сессии
*/
func SessionAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store := context.Get(r, "sessions")
		if store != nil {
			store := store.(*sessions.FilesystemStore)
			session, _ := store.Get(r, "user")
			userID := session.Values["user_id"]
			if userID != nil {
				//session.Options.MaxAge = config.Config.Session.MaxAge
				db := context.Get(r, "DB").(*gorm.DB)
				userManager := models.UserManager{db}
				user := userManager.GetById(userID.(uint))
				if user.ID != 0 {
					context.Set(r, "user", &user)
					h.ServeHTTP(w, r)
					return
				}
			}
		}
		logger.Logger.Warning(
			fmt.Sprintf(
				"Session, not authorized, forbidden: IP: %s, url: %s", utils.GetIPAddress(r),
				r.RequestURI,
			),
		)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Not authorized, forbidden"))
		logger.HandleError(err)
	})
}
