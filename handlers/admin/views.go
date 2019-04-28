package admin

import (
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"net/http"
	"reflect"
	"strings"
)

type Permission struct {
	Name   string                  `json:"name"`
	Access models.PermissionAccess `json:"access"`
}

type View struct {
	ViewID        uint       `json:"view_id"`
	Path          string     `json:"path"`
	ContentTypeID uint       `json:"content_type_id"`
	Name          string     `json:"name"`
	Plural        string     `json:"plural"`
	Icon          string     `json:"icon"`
	ParentID      uint       `json:"parent_id"`
	Children      []*View    `json:"children"`
	Permission    Permission `json:"permission"`
}

func (view *View) AddChild(child *View) {
	view.Children = append(view.Children, child)
}

/**
Получить все вьюхи доступные юзеру
*/
func GetAvailableViews(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)

	user := context.Get(r, "user").(*models.User)

	var data []*View
	assesPermission := Permission{
		Name:   models.Permission{}.GetChoiceAccess()[models.ACCESS],
		Access: models.ACCESS,
	}
	if user.IsSuperUser {
		var views []models.Views
		db.Preload("ContentType").Preload("Children").Find(&views)
		var contentTypes []models.ContentType
		db.Find(&contentTypes)

		for _, contentType := range contentTypes {
			model := databases.FindModelByContentType(db, contentType.Table)
			model = reflect.New(reflect.TypeOf(model)).Interface()
			obj := reflect.ValueOf(model)
			methodGetMeta := obj.MethodByName("Meta")
			var name, plural string
			if methodGetMeta.Kind() != reflect.Invalid {
				modelMeta := methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta)
				name = modelMeta.Name
				plural = modelMeta.Plural
			} else {
				name = obj.Elem().Type().Name()
				plural = name + "s"
			}
			path := strings.Replace(contentType.Table, "_", "-", -1)
			view := &View{
				ContentTypeID: contentType.ID,
				Name:          name,
				Plural:        plural,
				Path:          "/" + path,
				Permission:    assesPermission,
			}
			data = append(data, view)
		}

		for _, item := range views {
			stream := koazee.StreamOf(item.Children)
			view := &View{
				ViewID:        item.ID,
				ContentTypeID: item.ContentType.ID,
				Name:          item.Name,
				Plural:        item.Name,
				Path:          item.RoutePath,
				Icon:          item.Icon,
				ParentID:      item.ParentId,
				Permission:    assesPermission,
			}
			for _, _item := range views {
				child := stream.Filter(func(v models.Views) bool { return v.ID == _item.ID }).Out().Val()
				if len(child.([]models.Views)) > 0 {
					_view := &View{
						ViewID:        _item.ID,
						ContentTypeID: item.ContentType.ID,
						Name:          _item.Name,
						Plural:        _item.Name,
						Path:          _item.RoutePath,
						Icon:          _item.Icon,
						ParentID:      _item.ParentId,
						Permission:    assesPermission,
					}
					view.AddChild(_view)
				}
			}
			if view.ParentID == 0 {
				data = append(data, view)
			}
		}
	} else {
		var roles []models.Role
		db.Preload(
			"Permissions.View",
		).Preload(
			"Permissions.View.Children",
		).Joins(
			"INNER JOIN users_roles ur ON ur.role_id = id",
		).Find(&roles, "ur.user_id = ? ", user.ID)

		for _, role := range roles {
			for _, permission := range role.Permissions {
				stream := koazee.StreamOf(permission.View.Children)
				view := &View{
					ViewID:        permission.ViewId,
					Path:          permission.View.RoutePath,
					ContentTypeID: permission.View.ContentTypeId,
					Name:          permission.View.Name,
					ParentID:      permission.View.ParentId,
					Permission: Permission{
						Name:   permission.GetPermissionName(),
						Access: permission.Access,
					},
				}
				for _, _permission := range role.Permissions {
					child := stream.Filter(func(v models.Views) bool { return v.ID == _permission.ViewId }).Out().Val()
					if len(child.([]models.Views)) > 0 {
						_view := &View{
							ViewID:        _permission.ViewId,
							Path:          _permission.View.RoutePath,
							ContentTypeID: _permission.View.ContentTypeId,
							Name:          _permission.View.Name,
							ParentID:      _permission.View.ParentId,
							Permission: Permission{
								Name:   _permission.GetPermissionName(),
								Access: _permission.Access,
							},
						}
						view.AddChild(_view)
					}
				}
				if view.ParentID == 0 {
					data = append(data, view)
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
