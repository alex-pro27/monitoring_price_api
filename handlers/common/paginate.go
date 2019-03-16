package common

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"math"
	"reflect"
)

func Paginate(model interface{}, queryset *gorm.DB, page int, limit int, preloading []string) types.H {
	data := types.H{}
	t := reflect.ValueOf(model)
	if t.Kind() == reflect.Ptr {
		obj := t.Elem()
		if obj.Kind() == reflect.Slice {
			count := 0
			queryset.Model(model).Count(&count)
			start := page*limit - limit
			for _, preload := range preloading {
				queryset = queryset.Preload(preload)
			}
			queryset.Offset(start).Limit(limit).Find(model)

			var result []types.H

			for i := 0; i < obj.Len(); i++ {
				method := obj.Index(i).MethodByName("Serializer")
				if method.Kind() != reflect.Invalid {
					result = append(
						result,
						method.Call(nil)[0].Interface().(types.H),
					)
				}

			}
			var length int
			if length = limit; limit != len(result) {
				length = len(result)
			}
			if len(result) > 0 {
				data = types.H{
					"paginate": types.H{
						"current_page": page,
						"count":        count,
						"count_page":   math.Ceil(float64(count) / float64(limit)),
						"length":       length,
					},
					"result": result,
				}
			}
		}
	}

	return data
}
