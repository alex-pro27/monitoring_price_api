package handlers

import (
	"fmt"
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
	id, _ := strconv.Atoi(vars["id"])
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

func AllUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	user := context.Get(r, "user").(*models.User)
	fmt.Println(user.LastName, user.FirstName, user.Token.Key) // FIXME

	db := context.Get(r,"DB").(*gorm.DB)
	var users []models.User
	qs := db.Order("id")
	data := Paginate(&users, qs, page, 100)
	var idx []uint
	var results []H

	if len(data) == 0 {
		Error404(w)
		return
	}
	results = data["result"].([]H)

	for _, item := range results {
		idx = append(idx, item["id"].(uint))
	}
	var tokens []models.Token
	db.Where("user_id in (?)", idx).Find(&tokens)

	var usersRoles []struct{
		RoleID uint
		UserID uint
	}

	db.Table(
		"user_role",
	).Where(
		"user_id in (?)", idx,
	).Scan(&usersRoles)

	var rolesIDX []uint
	for _, item := range usersRoles{
		rolesIDX = append(rolesIDX, item.RoleID)
	}

	rolesIDX = Unique(rolesIDX)

	var roles []models.Role

	if len(rolesIDX) > 0 {
		db.Where("id in (?)", rolesIDX).Find(&roles)
	}

	for _, item := range results {
		inner:
		for i, token := range tokens {
			if item["id"] == token.UserID {
				tokens = append(tokens[:i], tokens[i+1:]...)
				item["token"] = token.Key
				break inner
			}
		}
		for _, role	:= range roles {
			for _, ur := range usersRoles {
				if role.ID == ur.RoleID && ur.UserID == item["id"] {
					item["roles"] = append(item["roles"].([]H), role.Serializer())
				}
			}
		}
	}

	JSONResponse(w, data)
}