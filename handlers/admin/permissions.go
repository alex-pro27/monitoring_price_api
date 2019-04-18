package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
	"sort"
)

/**
Проверка на доступ к представлению
*/
func CheckPermission(
	w http.ResponseWriter,
	r *http.Request,
	accessCode models.PermissionAccess,
	model interface{},
) bool {
	user := context.Get(r, "user").(*models.User)
	access := false
	if user.IsSuperUser {
		access = true
	} else {
		db := context.Get(r, "DB").(*gorm.DB)
		tableName := db.NewScope(model).GetModelStruct().TableName(db)
		permission := models.Permission{}
		db.Joins(
			"INNER JOIN views v ON v.id = view_id",
		).Joins(
			"INNER JOIN content_types ct ON ct.id = v.content_type_id",
		).Find(&permission, "ct.table = ?", tableName)
		if permission.Access >= accessCode {
			access = true
		}
	}
	if !access {
		common.Forbidden(w)
	}
	return access
}

/**
Получить список разрешений
*/
func GetPermissions(w http.ResponseWriter, r *http.Request) {
	type Permission struct {
		Code models.PermissionAccess `json:"code"`
		Name string                  `json:"name"`
	}
	var permissions []Permission
	for code, name := range models.ChoiceAccess {
		permissions = append(permissions, Permission{code, name})
	}
	sort.Slice(permissions, func(i, j int) bool {
		return permissions[i].Code < permissions[j].Code
	})
	common.JSONResponse(w, permissions)
}
