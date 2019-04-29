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

type Field struct {
	Name        string      `json:"name"`
	Label       string      `json:"label"`
	Type        string      `json:"type"`
	Disabled    bool        `json:"disabled"`
	Required    bool        `json:"required"`
	Value       interface{} `json:"value"`
	ContentType string      `json:"content_type"`
}

func getShortInfo(db *gorm.DB, model interface{}, where ...interface{}) (data []types.H) {
	model = reflect.New(reflect.TypeOf(model)).Interface()
	db.Find(model, where...)
	obj := reflect.ValueOf(model).Elem()

	addToData := func(iobj reflect.Value) {
		el := iobj
		if iobj.Kind() == reflect.Ptr {
			el = iobj.Elem()
		}
		strMethod := iobj.MethodByName("String")
		if el.Kind() != reflect.Invalid {
			if strMethod.Kind() != reflect.Invalid {
				data = append(data, types.H{
					"value": el.FieldByName("Model").FieldByName("ID").Interface(),
					"label": strMethod.Call(nil)[0].Interface().(string),
				})
			} else {
				data = append(data, types.H{
					"value": el.FieldByName("Model").FieldByName("ID").Interface(),
					"label": "[Object]",
				})
			}
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

func getFieldsFromModel(db *gorm.DB, model interface{}, where ...interface{}) (interface{}, []Field) {
	var fields []Field
	model = reflect.New(reflect.TypeOf(model)).Interface()
	scope := db.NewScope(model).GetStructFields()
	obj := reflect.ValueOf(model).Elem()
	var adminMeta types.AdminMeta

	adminMetaMethod := obj.MethodByName("Admin")
	if adminMetaMethod.Kind() != reflect.Invalid {
		adminMeta = adminMetaMethod.Call(nil)[0].Interface().(types.AdminMeta)
	}

	_streamExclude := koazee.StreamOf(adminMeta.ExcludeFields)
	_streamFields := koazee.StreamOf(adminMeta.Fields)

	if len(where) > 0 {
		qs := db
		for _, preload := range adminMeta.Preload {
			qs = db.Preload(preload)
		}
		qs.Find(model, where...)
	}

	var ID uint = 0
	for i := 0; i < obj.NumField(); i++ {
		form := helpers.ParseTag(obj.Type().Field(i).Tag.Get("form"))

		disabled, _ := strconv.ParseBool(form["disabled"])
		required, _ := strconv.ParseBool(form["required"])

		field := Field{
			Name:     obj.Type().Field(i).Name,
			Label:    form["label"],
			Type:     form["type"],
			Required: required,
			Disabled: disabled,
		}

		if field.Name == "Model" {
			field.Name = "id"
			field.Type = "hidden"
			field.Disabled = true
			field.Required = true
			field.Value = obj.Field(i).FieldByName("ID").Interface()
			ID = field.Value.(uint)
		} else {
			index, _ := _streamExclude.IndexOf(field.Name)
			if index > -1 {
				continue
			}
			index, _ = _streamFields.IndexOf(field.Name)
			if len(adminMeta.Fields) > 0 && index == -1 {
				continue
			}
			dbInfo := scope[i+3]
			field.Value = obj.Field(i).Interface()
			defaultValue := dbInfo.TagSettings["DEFAULT"]
			rel := dbInfo.Relationship

			if rel != nil {
				if ID == 0 {
					continue
				}
				field.Type = rel.Kind
				value := obj.Field(i).Interface()
				field.ContentType = db.NewScope(value).GetModelStruct().TableName(db)
				switch rel.Kind {
				case "many_to_many":
					var idx []uint
					tableName := rel.JoinTableHandler.Table(db)
					foreignKey := rel.JoinTableHandler.DestinationForeignKeys()[0].DBName
					db.Table(tableName).Where(
						fmt.Sprintf("%s = ?", rel.ForeignDBNames[0]), ID,
					).Pluck(foreignKey, &idx)
					field.Value = getShortInfo(db, value, "id IN (?)", idx)
					break
				case "has_many":
					field.Value = getShortInfo(
						db, obj.Field(i).Interface(),
						fmt.Sprintf("%s = ?", rel.ForeignDBNames[0]),
						ID,
					)
					break
				case "belongs_to":
					associationID := obj.FieldByName(rel.ForeignFieldNames[0]).Interface().(uint)
					field.Value = getShortInfo(db, value, "id = ?", associationID)
					if len(field.Value.([]types.H)) > 0 {
						field.Value = field.Value.([]types.H)[0]
					}
					break
				}
			} else {
				switch obj.Field(i).Interface().(type) {
				case int, uint, int32, int64, float32, float64:
					field.Type = "number"
					if ID == 0 {
						field.Value, _ = strconv.Atoi(defaultValue)
					}
					break
				case time.Time:
					field.Type = "date"
					if ID == 0 {
						if strings.Index(defaultValue, "now") > -1 {
							field.Value = time.Now().Format(helpers.ISO8601)
						} else {
							field.Value = defaultValue
						}
					} else {
						field.Value = field.Value.(time.Time).Format(helpers.ISO8601)
					}
					break
				case string:
					size, _ := strconv.Atoi(dbInfo.TagSettings["SIZE"])
					t := dbInfo.TagSettings["TYPE"]
					if (size == 0 || size > 255) && strings.Index(t, "varchar") == -1 {
						field.Type = "text"
					} else {
						field.Type = "string"
					}
					if ID == 0 {
						field.Value = defaultValue
					}
					break
				case bool:
					if form["type"] == "switch" {
						field.Type = "switch"
					} else {
						field.Type = "checkbox"
					}
					if ID == 0 {
						field.Value, _ = strconv.ParseBool(defaultValue)
					}
					break
				case pq.Int64Array:
					field.Name = "array"
					break
				}
			}
		}

		field.Name = helpers.ToSnakeCase(field.Name)
		if field.Label == "" {
			field.Label = field.Name
		}

		fields = append(fields, field)
	}
	return model, fields
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
		model, fields := getFieldsFromModel(db, model, "id = ?", id)
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
				meta["title"] = "[Object]"
			}
		}
		data := types.H{
			"meta":   meta,
			"fields": fields,
		}
		common.JSONResponse(w, data)
	} else {
		common.Forbidden(w)
	}
}

func CRUDContentType(w http.ResponseWriter, r *http.Request) {
	var err error
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
			common.ErrorResponse(w, "Object not found")
			return
		}
	}
	if crud.Kind() == reflect.Invalid {
		if r.Method == "POST" || r.Method == "PUT" {
			if err = helpers.SetFieldsForModel(model, fields); err != nil {
				common.ErrorResponse(w, err.Error())
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
		obj := reflect.New(reflect.TypeOf(model))
		adminMetaMethod := obj.MethodByName("Admin")
		var adminMeta types.AdminMeta
		if adminMetaMethod.Kind() != reflect.Invalid {
			adminMeta = adminMetaMethod.Call(nil)[0].Interface().(types.AdminMeta)
		}
		modelSlice := reflect.New(reflect.SliceOf(obj.Elem().Type())).Interface()
		paginateData := common.Paginate(modelSlice, db, page, 100, adminMeta.Preload, false)

		methodGetMeta := obj.MethodByName("Meta")
		meta := make(map[string]string)
		if methodGetMeta.Kind() != reflect.Invalid {
			modelMeta := methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta)
			meta["name"] = modelMeta.Name
			meta["plural"] = modelMeta.Plural
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
					"value": iobj.FieldByName("Model").FieldByName("ID").Interface(),
					"label": strMethod.Call(nil)[0].Interface().(string),
				})
			} else {
				result = append(result, types.H{
					"value": iobj.FieldByName("Model").FieldByName("ID").Interface(),
					"label": "[Object]",
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
