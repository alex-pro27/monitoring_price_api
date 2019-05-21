package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type ViewType int

const (
	VIEW_TYPE_CUSTOM ViewType = iota
	VIEW_TYPE_ALL
	VIEW_TYPE_EDIT
)

type Views struct {
	gorm.Model
	Name          string      `gorm:"size:255" form:"label:Название;required"`
	PositionMenu  uint        `gorm:"default:0" form:"label:Позиция в меню"`
	Icon          string      `grom:"size:255" form:"label:Иконка"`
	ViewType      ViewType    `gorm:"default:0" form:"label:Тип;choice:GetViewTypeChoices;required"`
	ParentId      uint        `gorm:"default:null"`
	Parent        *Views      `form:"label: Родитель;"`
	RoutePath     string      `gorm:"size:255" form:"label:URL path;required"`
	Children      []Views     `gorm:"foreignkey:ParentId" form:"label: Дети;"`
	ContentTypeId uint        `gorm:"default:null"`
	ContentType   ContentType `form:"label:Таблица;"`
}

func (Views) GetViewTypeChoices() map[ViewType]string {
	return map[ViewType]string{
		VIEW_TYPE_ALL:    "Список",
		VIEW_TYPE_EDIT:   "Редактор",
		VIEW_TYPE_CUSTOM: "Другое",
	}
}

func (view Views) GetViewTypeName() string {
	return view.GetViewTypeChoices()[view.ViewType]
}

func (view Views) Serializer() types.H {
	var childrenIDX []uint
	for _, child := range view.Children {
		childrenIDX = append(childrenIDX, child.ID)
	}
	return types.H{
		"id":           view.ID,
		"name":         view.Name,
		"icon":         view.Icon,
		"parent_id":    view.ParentId,
		"route_path":   view.RoutePath,
		"children_idx": childrenIDX,
		"content_type": view.ContentType,
	}
}

func (Views) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Представление",
		Plural: "Представления",
	}
}

func (Views) Admin() types.AdminMeta {
	return types.AdminMeta{
		ExcludeFields: []string{"ContentTypeId", "ParentId"},
		OrderBy:       []string{"PositionMenu"},
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
			{Name: "PositionMenu"},
		},
	}
}

func (view *Views) CRUD(db *gorm.DB) types.CRUDManager {
	return &ViewManager{db, view}
}

func (view Views) String() string {
	return view.Name
}
