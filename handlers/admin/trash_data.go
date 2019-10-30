package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func RecoveryFromTrash(w http.ResponseWriter, r *http.Request) {
	contentTypeName := r.FormValue("content_type_name")
	db := context.Get(r, "DB").(*gorm.DB)
	id, _ := strconv.Atoi(r.FormValue("id"))
	contentType := new(models.ContentType)
	db.First(contentType, "\"table\" = ?", contentTypeName)
	model := databases.FindModelByContentType(db, contentType.Table)
	if model == nil {
		common.Forbidden(w, r)
		return
	}
	if !CheckPermission(w, r, models.ACCESS, model) {
		return
	}
	if err := db.Exec(fmt.Sprintf("UPDATE \"%s\" SET deleted_at = NULL WHERE id = ?", contentTypeName), id); err.Error != nil {
		common.ErrorResponse(w, r, "Не удалось восстановить объект")
	} else {
		common.JSONResponse(w, types.H{"success": true})
	}
}

func GetTrashData(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	results := make([]types.H, 0)
	groups := make(types.H)
	user := context.Get(r, "user").(*models.User)
	contentTypes := make([]models.ContentType, 0)
	if !user.IsSuperUser {
		db.Joins(
			"INNER JOIN views v ON v.content_type_id = content_types.id",
		).Joins(
			"INNER JOIN permissions p ON p.view_id = v.id",
		).Find(&contentTypes, "p.access = ?", models.ACCESS)
	} else {
		db.Find(&contentTypes)
	}
	if len(contentTypes) == 0 {
		common.Forbidden(w, r)
		return
	}
	for _, ct := range contentTypes {
		model := databases.FindModelByContentType(db, ct.Table)
		if model == nil {
			continue
		}
		obj := reflect.New(reflect.TypeOf(model))
		methodGetMeta := obj.MethodByName("Meta")
		if methodGetMeta.Kind() == reflect.Invalid {
			continue
		}
		modelMeta := methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta)
		modelSliceObj := reflect.New(reflect.SliceOf(obj.Elem().Type()))
		tableName := db.NewScope(model).GetModelStruct().TableName(db)
		db.Raw(fmt.Sprintf("SELECT * FROM \"%s\" WHERE deleted_at IS NOT NULL", tableName)).Scan(modelSliceObj.Interface())
		itemData := make([]types.H, 0)
		for i := 0; i < modelSliceObj.Elem().Len(); i++ {
			iobj := modelSliceObj.Elem().Index(i)
			strMethod := iobj.MethodByName("String")
			id := iobj.FieldByName("Model").FieldByName("ID").Interface()
			deletedAt := iobj.FieldByName("Model").FieldByName("DeletedAt").Interface().(*time.Time)
			if strMethod.Kind() != reflect.Invalid {
				itemData = append(itemData, types.H{
					"value":             id,
					"content_type_name": tableName,
					"deleted_at":        deletedAt,
					"label":             strMethod.Call(nil)[0].Interface().(string),
				})
			} else {
				itemData = append(itemData, types.H{
					"value":             id,
					"content_type_name": tableName,
					"deleted_at":        deletedAt,
					"label":             "[Object]",
				})
			}
		}
		if len(itemData) > 0 {
			groups[tableName] = modelMeta.Plural
			results = append(results, itemData...)
		}
	}
	_groups := types.H{}
	koazee.StreamOf(results).Sort(func(x, y types.H) int {
		res := -1
		if y["deleted_at"].(*time.Time).Unix() > x["deleted_at"].(*time.Time).Unix() {
			res = 1
		}
		return res
	}).ForEach(func(x types.H) {
		x["deleted_at"] = x["deleted_at"].(*time.Time).Format(helpers.ISO8601)
		if group, ok := _groups[x["content_type_name"].(string)]; ok {
			group.(types.H)["results"] = append(group.(types.H)["results"].([]types.H), x)
		} else {
			group := types.H{
				"title":             groups[x["content_type_name"].(string)],
				"content_type_name": x["content_type_name"],
				"results":           []types.H{x},
			}
			_groups[x["content_type_name"].(string)] = group
		}
	}).Do().Out()

	common.JSONResponse(w, _groups)
}
