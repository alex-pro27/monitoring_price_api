package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

/**
Получить пользователя по ID
*/
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := context.Get(r, "DB").(*gorm.DB)
	userManager := models.UserManager{db}
	user := userManager.GetById(uint(id))
	if user.ID == 0 {
		common.ErrorResponse(w, "Пользователь не найден")
	} else {
		common.JSONResponse(w, user.Serializer())
	}
}

/**
Создание пользователя
*/
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	requiredData := map[string]interface{}{
		"UserName": r.PostFormValue("username"),
		"Email":    r.PostFormValue("email"),
		"Password": r.PostFormValue("password"),
	}

	extraData := map[string]interface{}{
		"FirstName": strings.Trim(r.PostFormValue("first_name"), ""),
		"LastName":  strings.Trim(r.PostFormValue("last_name"), ""),
		"Phone":     r.PostFormValue("phone"),
	}

	errs := helpers.SetFieldsOnModel(&user, requiredData, true)
	errs += helpers.SetFieldsOnModel(&user, extraData, false)

	if len(errs) > 0 {
		common.ErrorResponse(w, errs)
		return
	}
	db := context.Get(r, "DB").(*gorm.DB)
	userManager := models.UserManager{db}
	err := userManager.Create(&user)
	if err != nil {
		common.ErrorResponse(w, err.Error())
	} else {
		common.JSONResponse(w, user.Serializer())
	}
}

/**
Обновить информацию о пользователе по ID
*/
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := context.Get(r, "DB").(*gorm.DB)

	user := models.User{}
	db.First(&user, id)

	if user.ID == 0 {
		common.ErrorResponse(w, "Пользователь не найден")
		return
	}

	data := map[string]interface{}{
		"FirstName": strings.Trim(r.PostFormValue("first_name"), ""),
		"LastName":  strings.Trim(r.PostFormValue("last_name"), ""),
		"Email":     r.PostFormValue("email"),
		"Phone":     r.PostFormValue("phone"),
	}
	errs := helpers.SetFieldsOnModel(&user, data, false)

	if errs != "" {
		common.ErrorResponse(w, errs)
		return
	}
	_user := models.User{}
	db.First(&_user, "email = ?", user.Email)
	if _user.ID > 0 && _user.ID != user.ID {
		common.ErrorResponse(w, fmt.Sprintf("Email %s занят", user.Email))
		return
	}

	active, errParseBool := strconv.ParseBool(r.PostFormValue("active"))
	if errParseBool == nil {
		user.Active = active
	}
	tm := models.TokenManager{db}
	tm.NewToken(&user)
	db.Save(&user)
	common.JSONResponse(w, user.Serializer())
}

/**
Удаление юзера по ID
*/
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := context.Get(r, "DB").(*gorm.DB)
	user := models.User{}
	db.Delete(&user, id)
	common.JSONResponse(w, types.H{
		"error": false,
	})
}

/**
Список пользователей
*/
func AllUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	regionID, _ := strconv.Atoi(r.FormValue("region"))
	workGroupsID, _ := strconv.Atoi(r.FormValue("work_groups"))

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
	qs = qs.Order("last_name")

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
