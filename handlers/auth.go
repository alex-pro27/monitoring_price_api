package handlers

import (
	"errors"
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Авторизация по логину и паролю
*/
func Auth(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	db := context.Get(r, "DB").(*gorm.DB)
	user := models.User{}
	user.GetByUserName(db, username)
	if user.ID != 0 && user.CheckPassword(password) {
		context.Set(r, "user", &user)
		if Login(w, r) == nil {
			JSONResponse(w, user.Serializer())
			return
		}
	}
	Forrbidden(w)
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, common.H{
		"error": false,
	})
}

/**
Сохраняет в сессию Пользователя
*/
func Login(w http.ResponseWriter, r *http.Request) error {
	user := context.Get(r, "user")
	if user != nil {
		user := user.(*models.User)
		store := context.Get(r, "sessions")
		if store != nil {
			store := store.(*sessions.CookieStore)
			session, _ := store.Get(r, "user")
			session.Values["user_id"] = user.ID
			helpers.HandlerError(session.Save(r, w))
			return nil
		}
	}
	return errors.New("not auth")
}

/**
Удаление пользователя из сесси
*/
func Logout(w http.ResponseWriter, r *http.Request) {
	store := context.Get(r, "sessions").(*sessions.CookieStore)
	session, _ := store.Get(r, "user")
	session.Values["user_id"] = 0
	session.Options.MaxAge = -1
	helpers.HandlerError(session.Save(r, w))
	JSONResponse(w, common.H{
		"error": false,
	})
}
