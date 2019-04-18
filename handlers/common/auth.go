package common

import (
	"errors"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"net/http"
)

/**
Сохраняет в сессию Пользователя
*/
func Login(w http.ResponseWriter, r *http.Request) error {
	user := context.Get(r, "user")
	if user != nil {
		user := user.(*models.User)
		store := context.Get(r, "sessions")
		if store != nil {
			store := store.(*sessions.FilesystemStore)
			session, _ := store.Get(r, "user")
			session.Values["user_id"] = user.ID
			logger.HandleError(session.Save(r, w))
			return nil
		}
	}
	return errors.New("not auth")
}

/**
Удаление пользователя из сессии
*/
func Logout(w http.ResponseWriter, r *http.Request) {
	store := context.Get(r, "sessions").(*sessions.FilesystemStore)
	session, _ := store.Get(r, "user")
	session.Values["user_id"] = 0
	session.Options.MaxAge = -1
	logger.HandleError(session.Save(r, w))
}
