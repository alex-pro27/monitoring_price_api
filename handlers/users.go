package handlers

import (
	. "github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)


func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ErrorResponse(w, "id not a number")
		return
	}
	db := context.Get(r, "DB").(*gorm.DB)
	user := models.User{}
	user.GetById(db, id)
	if user.ID == 0 {
		ErrorResponse(w, "пользователь не найден")
	} else {
		JSONResponse(w, user.Serializer())
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		FirstName: 	r.FormValue("first_name"),
		LastName:  	r.FormValue("last_name"),
		UserName:  	r.FormValue("username"),
		Password:  	r.FormValue("password"),
		Email:		r.FormValue("email"),
		Phone:		r.FormValue("phone"),
	}
	db := context.Get(r,"DB").(*gorm.DB)
	err := user.Create(db)
	if err != nil {
		ErrorResponse(w, err.Error())
	} else {
		JSONResponse(w, user.Serializer())
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w,  H{
		"message": "Update",
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w,  H{
		"message": "Delete",
	})
}

func AllUsers(w http.ResponseWriter, r *http.Request)  {
	JSONResponse(w,  H{
		"message": "AllUsers",
	})
}