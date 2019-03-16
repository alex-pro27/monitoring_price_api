package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
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
	userManager := models.UserManager{db}
	user := userManager.GetById(uint(id))
	if user.ID == 0 {
		common.ErrorResponse(w, "пользователь не найден")
	} else {
		common.JSONResponse(w, user.Serializer())
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	userManager := models.UserManager{db}
	user, err := userManager.Create(
		&models.User{
			FirstName: r.PostFormValue("first_name"),
			LastName:  r.PostFormValue("last_name"),
			UserName:  r.PostFormValue("username"),
			Password:  r.PostFormValue("password"),
			Email:     r.PostFormValue("email"),
			Phone:     r.PostFormValue("phone"),
		},
	)
	if err != nil {
		common.ErrorResponse(w, err.Error())
	} else {
		common.JSONResponse(w, user.Serializer())
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO
	common.JSONResponse(w, types.H{
		"message": "Update",
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO
	common.JSONResponse(w, types.H{
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
			"INNER JOIN users_work_groups uw ON users.id = uw.user_id",
		).Where(
			"uw.work_group_id = ?", workGroupsID,
		)
	}

	if regionID != 0 {
		qs = qs.Joins(
			"INNER JOIN users_work_groups uw ON users.id = uw.user_id",
		).Joins(
			"INNER JOIN work_groups_regions wr ON uw.work_group_id = wr.work_group_id",
		).Where(
			"wr.regions_id = ?", regionID,
		)
		if workGroupsID != 0 {
			qs = qs.Where("uw.work_group_id = ?", workGroupsID)
		}
	}
	qs = qs.Order("id")

	data := common.Paginate(&users, qs, page, 100, []string{
		"Token",
		"Roles",
		"WorkGroup",
		"WorkGroup.Regions",
	})

	if len(data) == 0 {
		common.Error404(w)
		return
	}

	common.JSONResponse(w, data)
}
