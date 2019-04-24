package common

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"math"
	"reflect"
)

type PaginateInfo struct {
	CurrentPage int `json:"current_page"`
	Count       int `json:"count"`
	CountPage   int `json:"count_page"`
	Length      int `json:"length"`
}

type PaginateData struct {
	Paginate PaginateInfo  `json:"paginate"`
	Result   []interface{} `json:"result"`
}

func Paginate(
	model interface{},
	queryset *gorm.DB,
	page int,
	limit int,
	preloading []string,
	serialize bool,
) PaginateData {
	t := reflect.ValueOf(model)
	if t.Kind() == reflect.Ptr {
		obj := t.Elem()
		if obj.Kind() == reflect.Slice {
			if page == 0 {
				page = 1
			}
			count := 0
			queryset.Model(model).Count(&count)
			start := page*limit - limit
			for _, preload := range preloading {
				queryset = queryset.Preload(preload)
			}
			queryset.Offset(start).Limit(limit).Find(model)

			var result []interface{}

			for i := 0; i < obj.Len(); i++ {
				method := obj.Index(i).MethodByName("Serializer")
				if serialize && method.Kind() != reflect.Invalid {
					result = append(
						result,
						method.Call(nil)[0].Interface().(types.H),
					)
				} else {
					result = append(result, obj.Index(i).Interface())
				}
			}

			var length int
			if length = limit; limit != len(result) {
				length = len(result)
			}
			return PaginateData{
				Paginate: PaginateInfo{
					CurrentPage: page,
					Count:       count,
					CountPage:   int(math.Ceil(float64(count) / float64(limit))),
					Length:      length,
				},
				Result: result,
			}
		} else {
			panic("model not Slice")
		}
	} else {
		panic("model not Ptr")
	}
}
