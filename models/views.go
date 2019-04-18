package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Views struct {
	gorm.Model
	Name          string `gorm:"size:255"`
	ParentID      uint
	Parent        *Views
	RoutePath     string  `gorm:"size:255"`
	Children      []Views `gorm:"foreignkey:ParentID"`
	ContentTypeID uint
	ContentType   ContentType
}

func (view Views) Serializer() types.H {
	var childrenIDX []uint
	for _, child := range view.Children {
		childrenIDX = append(childrenIDX, child.ID)
	}
	return types.H{
		"id":           view.ID,
		"name":         view.Name,
		"parent_id":    view.ParentID,
		"route_path":   view.RoutePath,
		"children_idx": childrenIDX,
		"content_type": view.ContentType,
	}
}
