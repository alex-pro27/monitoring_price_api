package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Авторизация по логину и паролю
*/
func Login(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	db := context.Get(r, "DB").(*gorm.DB)
	userManager := models.UserManager{db}
	user := userManager.GetByUserName(username)
	if user.ID != 0 && user.CheckPassword(password) {
		context.Set(r, "user", &user)
		if common.Login(w, r) == nil {
			common.JSONResponse(w, user.Serializer())
			return
		}
	}
	common.Forbidden(w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	common.Logout(w, r)
	common.JSONResponse(w, types.H{
		"error": false,
	})

}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	common.JSONResponse(w, types.H{
		"error": false,
	})
}
