package admin

import (
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/wesovilabs/koazee"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func getShortInfo(db *gorm.DB, model interface{}, where ...interface{}) (data []types.H) {
	model = reflect.New(reflect.TypeOf(model)).Interface()
	db.Find(model, where...)
	obj := reflect.ValueOf(model).Elem()

	addToData := func(iobj reflect.Value) {
		strMethod := iobj.MethodByName("String")
		if strMethod.Kind() != reflect.Invalid {
			data = append(data, types.H{
				"id":    iobj.FieldByName("Model").FieldByName("ID").Interface(),
				"title": strMethod.Call(nil)[0].Interface().(string),
			})
		} else {
			data = append(data, types.H{
				"id":    iobj.FieldByName("Model").FieldByName("ID").Interface(),
				"title": "<Object>",
			})
		}
	}

	if obj.Kind() == reflect.Slice {
		for i := 0; i < obj.Len(); i++ {
			iobj := obj.Index(i)
			addToData(iobj)
		}
	} else {
		addToData(obj)
	}
	return data
}

func modelByFields(db *gorm.DB, model interface{}, where ...interface{}) (interface{}, [][]types.H) {
	var data [][]types.H
	model = reflect.New(reflect.TypeOf(model)).Interface()
	scope := db.NewScope(model).GetStructFields()
	obj := reflect.ValueOf(model).Elem()
	var (
		excludeFields []string
		fields        []string
	)

	adminMetaMethod := obj.MethodByName("Admin")
	if adminMetaMethod.Kind() != reflect.Invalid {
		adminMeta := adminMetaMethod.Call(nil)[0].Interface().(types.AdminMeta)
		excludeFields = adminMeta.ExcludeFields
		fields = adminMeta.Fields
	}

	_streamExclude := koazee.StreamOf(excludeFields)
	_streamFields := koazee.StreamOf(fields)
	getFields := func(obj reflect.Value) (data []types.H) {
		var ID uint = 0

		for i := 0; i < obj.NumField(); i++ {

			fieldName := obj.Type().Field(i).Name

			var (
				value       interface{}
				contentType string
			)

			form := helpers.ParseTag(obj.Type().Field(i).Tag.Get("form"))
			disabled, _ := strconv.ParseBool(form["disabled"])
			required, _ := strconv.ParseBool(form["required"])
			fieldType := form["type"]
			label := form["label"]

			if fieldName == "Model" {
				fieldName = "id"
				fieldType = "hidden"
				disabled = true
				required = true
				value = obj.Field(i).FieldByName("ID").Interface()
				ID = value.(uint)
			} else {

				index, _ := _streamExclude.IndexOf(fieldName)

				if index > -1 {
					continue
				}
				index, _ = _streamFields.IndexOf(fieldName)
				if len(fields) > 0 && index == -1 {
					continue
				}

				dbInfo := scope[i+3]
				value = obj.Field(i).Interface()
				defaultValue := dbInfo.TagSettings["DEFAULT"]

				switch obj.Field(i).Interface().(type) {
				case int, uint, int32, int64, float32, float64:
					fieldType = "number"
					if ID == 0 {
						value, _ = strconv.Atoi(defaultValue)
					}
					break
				case time.Time:
					fieldType = "date"
					if ID == 0 {
						if strings.Index(defaultValue, "now") > -1 {
							value = time.Now().Format(time.RFC3339Nano)
						} else {
							value = defaultValue
						}
					}
					break
				case string:
					size, _ := strconv.Atoi(dbInfo.TagSettings["SIZE"])
					t := dbInfo.TagSettings["TYPE"]
					if (size == 0 || size > 255) && strings.Index(t, "varchar") == -1 {
						fieldType = "text"
					} else {
						fieldType = "string"
					}
					if ID == 0 {
						value = defaultValue
					}
					break
				case bool:
					if form["type"] == "switch" {
						fieldType = "switch"
					} else {
						fieldType = "checkbox"
					}
					if ID == 0 {
						value, _ = strconv.ParseBool(defaultValue)
					}
					break
				case pq.Int64Array:
					fieldName = "array"
					break
				default:
					if reflect.Slice == obj.Field(i).Kind() {
						if ID > 0 {
							fieldType = "array_rel"
							contentType = dbInfo.DBName
							rel := dbInfo.Relationship
							if rel.Kind == "many_to_many" {
								var idx []uint
								tableName := rel.JoinTableHandler.Table(db)
								foreignKey := rel.JoinTableHandler.DestinationForeignKeys()[0].DBName
								db.Table(tableName).Where(
									fmt.Sprintf("%s = ?", rel.ForeignDBNames[0]), ID,
								).Pluck(foreignKey, &idx)
								value = getShortInfo(db, obj.Field(i).Interface(), "id IN (?)", idx)
							} else {
								continue
							}
						}
					} else if reflect.Struct == obj.Field(i).Kind() {
						fieldType = "rel"
						contentType = dbInfo.DBName
						if ID > 0 {
							value = getShortInfo(db, obj.Field(i).Interface(), "id = ?", ID)
							if len(value.([]types.H)) > 0 {
								value = value.([]types.H)[0]
							}
						} else {
							value = nil
						}
					}
				}
			}

			name := helpers.ToSnakeCase(fieldName)
			if label == "" {
				label = name
			}

			field := types.H{
				"name":     name,
				"label":    label,
				"type":     fieldType,
				"disabled": disabled,
				"required": required,
				"value":    value,
			}
			if fieldType == "rel" || fieldType == "array_rel" {
				field["content_type"] = contentType
			}
			data = append(data, field)
		}
		return data
	}
	if len(where) > 0 {
		db.Find(model, where...)
	}
	if obj.Kind() == reflect.Slice {
		for i := 0; i < obj.Len(); i++ {
			data = append(data, getFields(obj.Index(i)))
		}
	} else {
		data = append(data, getFields(obj))
	}
	return model, data
}

func GetContentTypeFields(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	contentTypeID, _ := strconv.Atoi(r.FormValue("content_type_id"))
	db := context.Get(r, "DB").(*gorm.DB)
	contentType := models.ContentType{}
	db.First(&contentType, contentTypeID)
	model := databases.FindModelByContentType(db, contentType.Table)

	if !CheckPermission(w, r, models.READ, model) {
		return
	}

	if model != nil {
		model, fields := modelByFields(db, model, "id = ?", id)
		if len(fields) > 0 {
			obj := reflect.ValueOf(model)
			methodGetMeta := obj.MethodByName("Meta")
			meta := make(map[string]string)
			if methodGetMeta.Kind() != reflect.Invalid {
				meta["name"] = methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta).Name
				meta["plural"] = methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta).Plural
			} else {
				meta["name"] = obj.Type().Name()
				meta["plural"] = obj.Type().Name() + "s"
			}
			strMethod := obj.MethodByName("String")
			ID := obj.Elem().FieldByName("Model").FieldByName("ID").Interface().(uint)
			if ID > 0 {
				if strMethod.Kind() != reflect.Invalid {
					meta["title"] = strMethod.Call(nil)[0].Interface().(string)
				} else {
					meta["title"] = "<Object>"
				}
			}
			data := types.H{
				"meta":   meta,
				"fields": fields[0],
			}
			common.JSONResponse(w, data)
			return
		}
	}
	common.Forbidden(w)
}

func CRUDContentType(w http.ResponseWriter, r *http.Request) {
	contentTypeID, _ := strconv.Atoi(r.PostFormValue("content_type_id"))
	db := context.Get(r, "DB").(*gorm.DB)
	contentType := models.ContentType{}
	db.First(&contentType, contentTypeID)
	model := databases.FindModelByContentType(db, contentType.Table)

	if model == nil {
		common.ErrorResponse(w, "Model not found")
		return
	}

	access := models.ACCESS
	if r.Method == "PUT" || r.Method == "POST" {
		access = models.WRITE
	}

	if !CheckPermission(w, r, access, model) {
		return
	}

	model = reflect.New(reflect.TypeOf(model)).Interface()
	obj := reflect.ValueOf(model)
	crud := obj.MethodByName("CRUD")
	var (
		fields types.H
		id     int
	)

	if r.Method == "PUT" || r.Method == "POST" {
		fields = make(types.H)
		if err := json.Unmarshal([]byte(r.PostFormValue("fields")), &fields); err != nil {
			common.ErrorResponse(w, err.Error())
			return
		}
	}

	if r.Method == "POST" || r.Method == "DELETE" {
		vars := mux.Vars(r)
		id, _ = strconv.Atoi(vars["id"])
		db.First(model, id)
		if obj.Elem().FieldByName("ID").Interface().(uint) == 0 {
			common.ErrorResponse(w, "")
			return
		}
	}

	if crud.Kind() == reflect.Invalid {
		if r.Method == "POST" || r.Method == "PUT" {
			errs := helpers.SetFieldsOnModel(model, fields)
			if errs != "" {
				common.ErrorResponse(w, errs)
				return
			}
		}
		switch r.Method {
		case "PUT":
			db.FirstOrCreate(model, model)
			break
		case "POST":
			db.Save(model)
			break
		case "DELETE":
			db.Delete(model, id)
			break
		}
	} else {
		manager := crud.Call([]reflect.Value{reflect.ValueOf(db)})[0].Interface().(types.CRUDManager)
		switch r.Method {
		case "PUT":
			if err := manager.Create(fields); err != nil {
				common.ErrorResponse(w, err.Error())
				return
			}
			break
		case "POST":
			if err := manager.Update(fields); err != nil {
				common.ErrorResponse(w, err.Error())
				return
			}
			break
		case "DELETE":
			if err := manager.Delete(fields); err != nil {
				common.ErrorResponse(w, err.Error())
				return
			}
			break
		}
	}
	if r.Method == "POST" || r.Method == "PUT" {
		methodSerialize := obj.MethodByName("Serializer")
		data := model
		if methodSerialize.Kind() != reflect.Invalid {
			data = methodSerialize.Call(nil)[0].Interface().(types.H)
		}
		common.JSONResponse(w, data)
	} else {
		common.JSONResponse(w, types.H{"error": false})
	}
}

func AllFieldsInModel(w http.ResponseWriter, r *http.Request) {
	contentTypeID, _ := strconv.Atoi(r.FormValue("content_type_id"))
	page, _ := strconv.Atoi(r.FormValue("page"))
	db := context.Get(r, "DB").(*gorm.DB)
	contentType := models.ContentType{}
	if contentTypeID == 0 {
		name := r.FormValue("content_type_name")
		db.First(&contentType, "\"table\" = ?", name)
	} else {
		db.First(&contentType, contentTypeID)
	}

	model := databases.FindModelByContentType(db, contentType.Table)

	if model != nil {
		if !CheckPermission(w, r, models.READ, model) {
			return
		}
		modelSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(model))).Interface()
		paginateData := common.Paginate(modelSlice, db, page, 100, []string{}, false)
		obj := reflect.ValueOf(model)
		methodGetMeta := obj.MethodByName("Meta")
		meta := make(map[string]string)
		if methodGetMeta.Kind() != reflect.Invalid {
			meta["name"] = methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta).Name
			meta["plural"] = methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta).Plural
		} else {
			meta["name"] = obj.Type().Name()
			meta["plural"] = obj.Type().Name() + "s"
		}
		var result []types.H
		for _, item := range paginateData.Result {
			iobj := reflect.ValueOf(item)
			strMethod := iobj.MethodByName("String")
			if strMethod.Kind() != reflect.Invalid {
				result = append(result, types.H{
					"id":    iobj.FieldByName("Model").FieldByName("ID").Interface(),
					"title": strMethod.Call(nil)[0].Interface().(string),
				})
			} else {
				result = append(result, types.H{
					"id":    iobj.FieldByName("Model").FieldByName("ID").Interface(),
					"title": "<Object>",
				})
			}
		}
		data := types.H{
			"meta":     meta,
			"paginate": paginateData.Paginate,
			"result":   result,
		}
		common.JSONResponse(w, data)
	} else {
		common.Forbidden(w)
	}
}
