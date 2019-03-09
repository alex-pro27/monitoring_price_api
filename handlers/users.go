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
		FirstName: r.PostFormValue("first_name"),
		LastName:  r.PostFormValue("last_name"),
		UserName:  r.PostFormValue("username"),
		Password:  r.PostFormValue("password"),
		Email:     r.PostFormValue("email"),
		Phone:     r.PostFormValue("phone"),
	}
	db := context.Get(r, "DB").(*gorm.DB)
	err := user.Create(db)
	if err != nil {
		ErrorResponse(w, err.Error())
	} else {
		JSONResponse(w, user.Serializer())
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO
	JSONResponse(w, H{
		"message": "Update",
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO
	JSONResponse(w, H{
		"message": "Delete",
	})
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	regionID, _ := strconv.Atoi(r.FormValue("region"))
	workGroupsID, _ := strconv.Atoi(r.FormValue("work_groups"))

	user := context.Get(r, "user")

	if user != nil {
		user := user.(*models.User)
		fmt.Println(user.LastName, user.FirstName, user.Token.Key) // FIXME
	}

	db := context.Get(r, "DB").(*gorm.DB)
	qs := db
	var users []models.User

	if workGroupsID != 0 && regionID == 0 {
		qs = qs.Joins(
			"INNER JOIN user_workgroup uw ON users.id = uw.user_id",
		).Where(
			"uw.work_group_id = ?", workGroupsID,
		)
	}

	if regionID != 0 {
		qs = qs.Joins(
			"INNER JOIN user_workgroup uw ON users.id = uw.user_id",
		).Joins(
			"INNER JOIN workgroup_regions wr ON uw.work_group_id = wr.work_group_id",
		).Where(
			"wr.regions_id = ?", regionID,
		)
		if workGroupsID != 0 {
			qs = qs.Where("uw.work_group_id = ?", workGroupsID)
		}
	}
	qs = qs.Order("id")

	data := Paginate(&users, qs, page, 100, []string{
		"Token",
		"Roles",
		"WorkGroup",
		"WorkGroup.Regions",
	})

	//if len(data) == 0 {
	//	Error404(w)
	//	return
	//}

	JSONResponse(w, data)
}
