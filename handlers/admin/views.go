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
Получить все вьюхи доступные юзеру
*/
func GetAvailableViews(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)

	user := context.Get(r, "user").(*models.User)
	type Permission struct {
		Name   string                  `json:"name"`
		Access models.PermissionAccess `json:"access"`
	}

	type View struct {
		ID         uint       `json:"id"`
		Name       string     `json:"name"`
		ParentID   uint       `json:"parent_id"`
		Children   []uint     `json:"children"`
		Permission Permission `json:"permission"`
	}

	data := make(map[string]View)

	if user.IsSuperUser {
		var views []models.Views
		db.Find(&views)

		for _, view := range views {
			var childrenIDX []uint
			for _, child := range view.Children {
				childrenIDX = append(childrenIDX, child.ID)
			}
			data[view.RoutePath] = View{
				ID:       view.ID,
				Name:     view.Name,
				ParentID: view.ParentID,
				Children: childrenIDX,
				Permission: Permission{
					Name:   models.ChoiceAccess[models.ACCESS],
					Access: models.ACCESS,
				},
			}
		}
	} else {
		var roles []models.Role
		db.Preload(
			"Permissions.View.Children",
		).Joins(
			"INNER JOIN users_roles ur ON ur.role_id = id",
		).Find(&roles, "ur.user_id = ?", user.ID)

		for _, role := range roles {
			for _, permission := range role.Permissions {
				var childrenIDX []uint
				for _, child := range permission.View.Children {
					childrenIDX = append(childrenIDX, child.ID)
				}
				if permission.Access > data[permission.View.RoutePath].Permission.Access {
					data[permission.View.RoutePath] = View{
						ID:       permission.ViewID,
						Name:     permission.View.Name,
						ParentID: permission.View.ParentID,
						Children: childrenIDX,
						Permission: Permission{
							Name:   permission.GetPermissionName(),
							Access: permission.Access,
						},
					}
				}
			}
		}
	}
	if user.IsSuperUser || len(data) > 0 {
		common.JSONResponse(w, data)
	} else {
		common.Forbidden(w)
	}
}

func GetView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	db := context.Get(r, "DB").(*gorm.DB)
	view := models.Views{}
	db.Preloads("ContentType").First(&view, id)
	common.JSONResponse(w, view.Serializer())
}

/**
Получить все представления в адмике
*/
func AllViews(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	var views []models.Views
	db.Find(&views)
	var data []types.H
	for _, item := range views {
		data = append(data, item.Serializer())
	}
	common.JSONResponse(w, data)
}

/**
Записать информацию о доступном представлении в админке
@method PUT
@param {
	"name": string
	"route_path": string
	"parent_id": int
	"children_idx": []int
}
*/
func CreateView(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*models.User)
	if user.IsSuperUser {
		db := context.Get(r, "DB").(*gorm.DB)
		view := models.Views{
			Name:      r.PostFormValue("name"),
			RoutePath: r.PostFormValue("route_path"),
		}
		parentID, _ := strconv.Atoi(r.FormValue("prent_id"))
		if parentID > 0 {
			view.ParentID = uint(parentID)
		}
		var childrenIDX []int
		err := json.Unmarshal([]byte(r.PostFormValue("children_idx")), &childrenIDX)

		db.Create(&view)
		db.NewRecord(view)

		if err != nil {
			var children []models.Views
			db.Model(&children).Where("id IN (?)", childrenIDX).Update("parent_id", view.ID)
			view.Children = children
		}

		common.JSONResponse(w, view.Serializer())
	} else {
		common.Forbidden(w)
	}
}
