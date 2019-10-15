package middleware

import (
	"encoding/base64"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

const (
	BASIC_AUTH   uint = 1
	TOKEN_AUTH   uint = 2
	SESSION_AUTH uint = 4
)

func basicAuth(r *http.Request) *models.User {
	authHeader := r.Header.Get("Authorization")
	authData := strings.Split(authHeader, " ")
	if len(authData) == 2 && authData[0] == "Basic" {
		data, err := base64.StdEncoding.DecodeString(authData[1])
		if err == nil {
			logpassw := strings.Split(string(data), ":")
			db := context.Get(r, "DB").(*gorm.DB)
			user := models.User{}
			user.Manager(db).GetByUserName(logpassw[0])
			if user.CheckPassword(logpassw[1]) {
				return &user
			}
		}
	}
	return nil
}

func tokenAuth(r *http.Request) *models.User {
	authHeader := r.Header.Get("Authorization")
	authData := strings.Split(authHeader, " ")
	if len(authData) == 2 && authData[0] == "Token" {
		db := context.Get(r, "DB").(*gorm.DB)
		user := models.User{}
		user.Manager(db).GetUserByToken(authData[1])
		if user.ID > 0 {
			return &user
		}
	}
	return nil
}

func sessionAuth(r *http.Request) *models.User {
	store := context.Get(r, "sessions")
	if store != nil {
		store := store.(*sessions.FilesystemStore)
		session, _ := store.Get(r, "user")
		userID := session.Values["user_id"]
		if userID != nil {
			//session.Options.MaxAge = config.Config.Session.MaxAge
			db := context.Get(r, "DB").(*gorm.DB)
			user := models.User{}
			user.Manager(db).GetById(userID.(uint))
			if user.ID != 0 {
				return &user
			}
		}
	}
	return nil
}

/**
Basic Авторизация
*/
func BasicAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := basicAuth(r)
		if user != nil {
			context.Set(r, "user", user)
			h.ServeHTTP(w, r)
			return
		}
		logger.Logger.Warningf(
			"Not authorized, forbidden: IP: %s, url: %s", utils.GetIPAddress(r), r.RequestURI,
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
		user := tokenAuth(r)
		if user != nil {
			context.Set(r, "user", user)
			h.ServeHTTP(w, r)
			return
		}
		logger.Logger.Warningf(
			"Invalid token: IP: %s, url: %s", utils.GetIPAddress(r), r.RequestURI,
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
		user := sessionAuth(r)
		if user != nil {
			context.Set(r, "user", user)
			h.ServeHTTP(w, r)
			return
		}
		logger.Logger.Warningf(
			"Session, not authorized, forbidden: IP: %s, url: %s", utils.GetIPAddress(r), r.RequestURI,
		)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Not authorized, forbidden"))
		logger.HandleError(err)
	})
}

func MixinAuthMiddle(authTypes uint) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allFlags := []uint{BASIC_AUTH, TOKEN_AUTH, SESSION_AUTH}
			flagsIndexes := helpers.GetFlags(authTypes, len(allFlags))
			var user *models.User
			isBasic := false
		CYCLE:
			for _, i := range flagsIndexes {
				switch allFlags[i] {
				case BASIC_AUTH:
					isBasic = true
					if user = basicAuth(r); user != nil {
						break CYCLE
					}
				case TOKEN_AUTH:
					if user = tokenAuth(r); user != nil {
						break CYCLE
					}
				case SESSION_AUTH:
					if user = sessionAuth(r); user != nil {
						break CYCLE
					}
				}
			}
			if user != nil {
				context.Set(r, "user", user)
				h.ServeHTTP(w, r)
				return
			}
			logger.Logger.Warningf(
				"Not authorized, forbidden: IP: %s, url: %s", utils.GetIPAddress(r), r.RequestURI,
			)
			if isBasic {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			}
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Not authorized, forbidden"))
			logger.HandleError(err)
		})
	}
}
