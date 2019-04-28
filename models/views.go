package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type ViewType int

const (
	VIEW_TYPE_ALL ViewType = iota
	VIEW_TYPE_EDIT
	VIEW_TYPE_CUSTOM
)

type Views struct {
	gorm.Model
	Name          string   `gorm:"size:255"`
	Icon          string   `grom:"size:255"`
	ViewType      ViewType `gorm:"default:0" choice:"ViewTypeChoices"`
	ParentId      uint     `gorm:"default:null"`
	Parent        *Views
	RoutePath     string  `gorm:"size:255"`
	Children      []Views `gorm:"foreignkey:ParentId"`
	ContentTypeId uint
	ContentType   ContentType
}

func (Views) GetViewTypeChoices() map[ViewType]string {
	return map[ViewType]string{
		VIEW_TYPE_ALL:    "Для списка",
		VIEW_TYPE_EDIT:   "Для редактирования",
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

func (view Views) String() string {
	return view.Name
}
