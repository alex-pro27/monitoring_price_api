package admin

import (
	"encoding/json"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

/**
Получить все роли
*/
func AllRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.Role
	db := context.Get(r, "DB").(*gorm.DB)
	db.Find(&roles)
	var data []types.H
	for _, item := range roles {
		data = append(data, item.Serializer())
	}
	common.JSONResponse(w, roles)
}

/**
Получить роль по id
*/
func GetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := context.Get(r, "DB").(*gorm.DB)
	role := models.Role{}
	db.First(&role, "id = ?", id)
	if role.ID > 0 {
		common.JSONResponse(w, role.Serializer())
	} else {
		common.Error404(w)
	}
}

/**
Создание роли
@method PUT
@param {
	"name": string
	"view_id": int
	"permissions_views": []{
		"view_id": int
		"permission": int
	}
*/
func CreateRole(w http.ResponseWriter, r *http.Request) {
	type PermissionView struct {
		ViewID     int                     `json:"view_id"`
		Permission models.PermissionAccess `json:"permission"`
	}

	var permissionsViews []PermissionView

	err := json.Unmarshal([]byte(r.PostFormValue("permissions_views")), &permissionsViews)

	if err != nil {
		common.ErrorResponse(w, "ошибка парсинга данных")
		return
	}

	db := context.Get(r, "DB").(*gorm.DB)

	var permissions []models.Permission

	for _, pw := range permissionsViews {
		permission := models.Permission{
			Access: pw.Permission,
			ViewID: uint(pw.ViewID),
		}
		permissions = append(permissions, permission)
	}

	role := models.Role{
		Name:        r.PostFormValue("name"),
		Permissions: permissions,
	}
	db.FirstOrCreate(&role)
	db.NewRecord(role)
	common.JSONResponse(w, role.Serializer())
}
