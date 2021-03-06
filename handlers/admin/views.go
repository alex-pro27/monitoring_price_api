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
	"sort"
	"strings"
)

type Permission struct {
	Name   string                  `json:"name"`
	Access models.PermissionAccess `json:"access"`
}

type View struct {
	ViewID          uint            `json:"view_id"`
	Path            string          `json:"path"`
	ContentTypeID   uint            `json:"content_type_id"`
	ContentTypeName string          `json:"content_type_name"`
	Name            string          `json:"name"`
	Plural          string          `json:"plural"`
	Icon            string          `json:"icon"`
	ParentID        uint            `json:"parent_id"`
	Children        []*View         `json:"children"`
	Permission      Permission      `json:"permission"`
	ViewType        models.ViewType `json:"view_type"`
	Menu            bool            `json:"menu"`
	Position        uint            `json:"-"`
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
		Name:   models.PermissionAccessChoices[models.ACCESS],
		Access: models.ACCESS,
	}
	if user.IsSuperUser {
		var views []models.Views
		db.Preload("ContentType").Preload("Children").Find(&views)
		var contentTypes []models.ContentType
		db.Find(&contentTypes)
		streamContentTypes := koazee.StreamOf(contentTypes)
		for _, item := range views {
			stream := koazee.StreamOf(item.Children)
			streamContentTypes = streamContentTypes.Filter(func(ct models.ContentType) bool {
				return ct.ID != item.ContentType.ID
			}).Do()
			view := &View{
				ViewID:          item.ID,
				ContentTypeID:   item.ContentType.ID,
				ContentTypeName: item.ContentType.Table,
				Name:            item.Name,
				Plural:          item.Name,
				Path:            item.RoutePath,
				Icon:            item.Icon,
				ParentID:        item.ParentId,
				Permission:      assesPermission,
				ViewType:        item.ViewType,
				Position:        item.PositionMenu,
			}
			for _, _item := range views {
				child := stream.Filter(func(v models.Views) bool { return v.ID == _item.ID }).Out().Val()
				if len(child.([]models.Views)) > 0 {
					_view := &View{
						ViewID:          _item.ID,
						ContentTypeID:   item.ContentType.ID,
						ContentTypeName: item.ContentType.Table,
						Name:            _item.Name,
						Plural:          _item.Name,
						Path:            _item.RoutePath,
						Icon:            _item.Icon,
						ParentID:        _item.ParentId,
						Permission:      assesPermission,
						ViewType:        _item.ViewType,
						Position:        _item.PositionMenu,
					}
					view.AddChild(_view)
				}
			}
			if view.ParentID == 0 {
				data = append(data, view)
			}
		}
		for _, contentType := range streamContentTypes.Out().Val().([]models.ContentType) {
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
	} else {
		var roles []models.Role
		db.Preload(
			"Permissions.View.ContentType",
		).Preload(
			"Permissions.View.Children",
		).Joins(
			"INNER JOIN users_roles ur ON ur.role_id = id",
		).Find(&roles, "ur.user_id = ? ", user.ID)

		for _, role := range roles {
			for _, permission := range role.Permissions {
				stream := koazee.StreamOf(permission.View.Children)
				view := &View{
					ViewID:          permission.ViewId,
					Path:            permission.View.RoutePath,
					ContentTypeID:   permission.View.ContentTypeId,
					ContentTypeName: permission.View.ContentType.Table,
					Name:            permission.View.Name,
					Icon:            permission.View.Icon,
					ParentID:        permission.View.ParentId,
					ViewType:        permission.View.ViewType,
					Menu:            permission.View.PositionMenu > 0,
					Position:        permission.View.PositionMenu,
					Permission: Permission{
						Name:   permission.GetPermissionName(),
						Access: permission.Access,
					},
				}
				for _, _permission := range role.Permissions {
					child := stream.Filter(func(v models.Views) bool { return v.ID == _permission.ViewId }).Out().Val()
					if len(child.([]models.Views)) > 0 {
						_view := &View{
							ViewID:          _permission.ViewId,
							Path:            _permission.View.RoutePath,
							ContentTypeID:   _permission.View.ContentTypeId,
							ContentTypeName: _permission.View.ContentType.Table,
							Name:            _permission.View.Name,
							Icon:            _permission.View.Icon,
							ParentID:        _permission.View.ParentId,
							ViewType:        _permission.View.ViewType,
							Menu:            _permission.View.PositionMenu > 0,
							Position:        _permission.View.PositionMenu,
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
		sort.Slice(data, func(i, j int) bool {
			return data[i].Position < data[j].Position
		})
		common.JSONResponse(w, data)
	} else {
		common.Forbidden(w, r)
	}
}
