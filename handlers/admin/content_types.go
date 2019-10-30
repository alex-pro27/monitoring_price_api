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
	Name         string                   `json:"name"`
	Label        string                   `json:"label"`
	Type         string                   `json:"type"`
	Disabled     bool                     `json:"disabled"`
	Required     bool                     `json:"required"`
	Value        interface{}              `json:"value"`
	ContentType  string                   `json:"content_type"`
	Options      []map[string]interface{} `json:"options"`
	Groups       map[uint]string          `json:"groups"`
	GroupBy      string                   `json:"group_by"`
	GroupByField string                   `json:"group_by_field"`
}

func addShortInfoToData(iobj reflect.Value, data *[]types.H) {
	el := iobj
	if iobj.Kind() == reflect.Ptr {
		el = iobj.Elem()
	}
	strMethod := iobj.MethodByName("String")
	if el.Kind() != reflect.Invalid {
		id := el.FieldByName("Model").FieldByName("ID").Interface()
		if id.(uint) > 0 {
			if strMethod.Kind() != reflect.Invalid {
				*data = append(*data, types.H{
					"value": id,
					"label": strMethod.Call(nil)[0].Interface().(string),
				})
			} else {
				*data = append(*data, types.H{
					"value": id,
					"label": "[Object]",
				})
			}
		}
	}
}

func addGroupShortInfoToData(obj reflect.Value, groupBy string, groups *map[uint]string, data *[]types.H) {
	el := obj
	if obj.Kind() == reflect.Ptr {
		el = obj.Elem()
	}

	addGroup := func(iobj, group reflect.Value) {
		var groupID uint
		var id interface{}
		var itemLabel string
		var groupName string
		groupID = group.FieldByName("Model").FieldByName("ID").Interface().(uint)
		strMethod := iobj.MethodByName("String")
		groupStrMethod := group.MethodByName("String")

		if groupStrMethod.Kind() != reflect.Invalid {
			groupName = groupStrMethod.Call(nil)[0].Interface().(string)
		} else {
			groupName = "[Object]"
		}
		if strMethod.Kind() != reflect.Invalid {
			itemLabel = strMethod.Call(nil)[0].Interface().(string)
		} else {
			itemLabel = "[Object]"
		}

		id = iobj.FieldByName("Model").FieldByName("ID").Interface()
		var item types.H

		item = types.H{
			"value":    id,
			"label":    itemLabel,
			"group_id": groupID,
		}
		if groups != nil {
			(*groups)[groupID] = groupName
		}

		*data = append(*data, item)
	}

	for i := 0; i < el.Len(); i++ {
		iobj := el.Index(i)
		if iobj.Kind() == reflect.Ptr {
			iobj = iobj.Elem()
		}
		if iobj.Kind() == reflect.Interface {
			iobj = reflect.ValueOf(iobj.Interface())
		}
		group := iobj.FieldByName(groupBy)
		if group.Kind() == reflect.Slice {
			for i := 0; i < group.Len(); i++ {
				addGroup(iobj, group.Index(i))
			}
		} else {
			addGroup(iobj, group)
		}
	}
}

func getShortInfo(db *gorm.DB, model interface{}, groupBy string, where ...interface{}) (groups map[uint]string, data []types.H) {
	groups = make(map[uint]string)
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		model = reflect.New(reflect.TypeOf(model)).Interface()
	}
	obj := reflect.ValueOf(model)
	var adminMetaMethod reflect.Value
	if obj.Elem().Kind() == reflect.Slice {
		itemModel := reflect.New(obj.Elem().Type().Elem())
		adminMetaMethod = itemModel.MethodByName("Admin")
	} else {
		adminMetaMethod = obj.MethodByName("Admin")
	}

	if adminMetaMethod.Kind() != reflect.Invalid {
		adminMeta := adminMetaMethod.Call(nil)[0].Interface().(types.AdminMeta)
		for _, preload := range adminMeta.Preload {
			db = db.Preload(preload)
		}
	}
	if groupBy != "" {
		db = db.Preload(groupBy)
	}
	db.Find(model, where...)

	if obj.Elem().Kind() == reflect.Slice {
		if groupBy != "" && obj.Elem().Len() > 0 && obj.Elem().Index(0).FieldByName(groupBy).Kind() != reflect.Invalid {
			addGroupShortInfoToData(obj, groupBy, &groups, &data)
		} else {
			for i := 0; i < obj.Elem().Len(); i++ {
				iobj := obj.Elem().Index(i)
				addShortInfoToData(iobj, &data)
			}
		}
	} else {
		addShortInfoToData(obj.Elem(), &data)
	}
	return groups, data
}

func getFieldsFromModel(db *gorm.DB, model interface{}, where ...interface{}) (interface{}, []Field) {
	var fields []Field
	model = reflect.New(reflect.TypeOf(model)).Interface()
	structFields := db.NewScope(model).GetStructFields()
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
		//for _, preload := range adminMeta.Preload {
		//	qs = db.Preload(preload)
		//}
		qs.Find(model, where...)
	}

	var ID uint = 0
	for i := 0; i < obj.NumField(); i++ {
		fieldObj := obj.Field(i)
		form := helpers.ParseTag(obj.Type().Field(i).Tag.Get("form"))
		disabled, _ := strconv.ParseBool(form["disabled"])
		required, _ := strconv.ParseBool(form["required"])
		empty, _ := strconv.ParseBool(form["empty"])

		var groupName string
		if form["group_by"] != "" && fieldObj.Kind() == reflect.Slice {
			elType := reflect.TypeOf(fieldObj.Interface()).Elem()
			groupField := reflect.New(elType).Elem().FieldByName(form["group_by"])
			groupFieldKind := groupField.Kind()
			if groupFieldKind == reflect.Struct {
				groupName = db.NewScope(groupField.Interface()).GetModelStruct().TableName(db)
			} else if groupFieldKind == reflect.Slice {
				elType := reflect.TypeOf(groupField.Interface()).Elem()
				groupField := reflect.New(elType).Elem()
				if groupField.Kind() == reflect.Struct {
					groupName = db.NewScope(groupField.Interface()).GetModelStruct().TableName(db)
				}
			}
		}

		field := Field{
			Name:         obj.Type().Field(i).Name,
			Label:        form["label"],
			Type:         form["type"],
			Required:     required,
			Disabled:     disabled,
			GroupBy:      groupName,
			GroupByField: form["group_by"],
		}

		if field.Name == "Model" {
			field.Name = "id"
			field.Type = "hidden"
			field.Disabled = true
			field.Required = true
			field.Value = fieldObj.FieldByName("ID").Interface()
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
			structField := structFields[i+3]
			defaultValue := structField.TagSettings["DEFAULT"]
			rel := structField.Relationship
			if !empty {
				value := obj.Field(i).Interface()
				field.Value = value
				if rel != nil {
					field.Type = rel.Kind
					field.ContentType = db.NewScope(value).GetModelStruct().TableName(db)
					if ID != 0 {
						switch rel.Kind {
						case "many_to_many":
							var idx []uint
							tableName := rel.JoinTableHandler.Table(db)
							foreignKey := rel.JoinTableHandler.DestinationForeignKeys()[0].DBName
							db.Table(tableName).Where(
								fmt.Sprintf("\"%s\" = ?", rel.ForeignDBNames[0]), ID,
							).Pluck(foreignKey, &idx)
							field.Groups, field.Value = getShortInfo(db, value, form["group_by"], "id IN (?)", idx)
							break
						case "has_many":
							field.Groups, field.Value = getShortInfo(
								db,
								value,
								form["group_by"],
								fmt.Sprintf("\"%s\" = ?", rel.ForeignDBNames[0]),
								ID,
							)
							break
						case "belongs_to":
							associationID := obj.FieldByName(rel.ForeignFieldNames[0]).Interface().(uint)
							if obj.Field(i).Kind() == reflect.Ptr {
								value = reflect.New(obj.Type().Field(i).Type.Elem()).Interface()
							}
							field.Groups, field.Value = getShortInfo(db, value, form["group_by"], "id = ?", associationID)
							if len(field.Value.([]types.H)) > 0 {
								field.Value = field.Value.([]types.H)[0]
							}
							break
						}
					} else {
						field.Value = nil
					}
				} else if form["choice"] != "" {
					methodChoice := obj.MethodByName(form["choice"])
					if methodChoice.Kind() != reflect.Invalid {
						field.Type = "choice"
						label := "No name"
						methodLabel := obj.MethodByName(fmt.Sprintf("Get%sName", obj.Type().Field(i).Name))
						if methodLabel.Kind() != reflect.Invalid {
							label = methodLabel.Call(nil)[0].Interface().(string)
						}
						field.Value = map[string]interface{}{
							"label": label,
							"value": value,
						}
						choices := methodChoice.Call(nil)[0]
						for _, k := range choices.MapKeys() {
							field.Options = append(field.Options, map[string]interface{}{
								"value": k.Interface(),
								"label": choices.MapIndex(k).Interface(),
							})
						}
					}
				} else {
					switch fieldObj.Interface().(type) {
					case int, uint, int32, int64, float32, float64:
						field.Type = "number"
						if ID == 0 {
							field.Value, _ = strconv.Atoi(defaultValue)
						}
						break
					case time.Time:
						if form["type"] == "date" {
							field.Type = "date"
						} else {
							field.Type = "datetime-local"
						}

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
						size, _ := strconv.Atoi(structField.TagSettings["SIZE"])
						t := structField.TagSettings["TYPE"]
						if form["type"] != "password" {
							if (size == 0 || size > 255) && strings.Index(t, "varchar") == -1 {
								field.Type = "text"
							} else {
								field.Type = "string"
							}
							if ID == 0 {
								field.Value = defaultValue
							}
						} else {
							field.Type = form["type"]
							field.Value = nil
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
						field.Type = "array"
						if field.Value != nil {
							list := make([]string, 0)
							for _, item := range field.Value.(pq.Int64Array) {
								list = append(list, strconv.FormatInt(item, 10))
							}
							field.Value = strings.Join(list, ",")
						}
						break
					}
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
		common.Forbidden(w, r)
	}
}

func CRUDContentType(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	contentTypeID, _ := strconv.Atoi(vars["content_type_id"])
	db := context.Get(r, "DB").(*gorm.DB)
	contentType := models.ContentType{}
	db.First(&contentType, contentTypeID)
	model := databases.FindModelByContentType(db, contentType.Table)

	if model == nil {
		common.ErrorResponse(w, r, "Model not found")
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
			common.ErrorResponse(w, r, err.Error())
			return
		}
	}

	if r.Method == "POST" || r.Method == "DELETE" {
		id, _ = strconv.Atoi(vars["id"])
		db.First(model, id)
		if obj.Elem().FieldByName("ID").Interface().(uint) == 0 {
			common.ErrorResponse(w, r, "Object not found")
			return
		}
	}
	if crud.Kind() == reflect.Invalid {
		if r.Method == "POST" || r.Method == "PUT" {
			if err = helpers.SetFieldsForModel(model, fields); err != nil {
				common.ErrorResponse(w, r, err.Error())
				return
			}
		}
		switch r.Method {
		case "PUT":
			if res := db.FirstOrCreate(model, model); res.Error != nil {
				common.ErrorResponse(w, r, "Ошибка добавления записи")
				return
			}
			helpers.SetManyToMany(db, model, fields)
			break
		case "POST":
			if res := db.Save(model); res.Error != nil {
				common.ErrorResponse(w, r, "Ошибка обновления записи")
				return
			}
			helpers.SetManyToMany(db, model, fields)
			break
		case "DELETE":
			if res := db.Delete(model, id); res.Error != nil {
				common.ErrorResponse(w, r, "Ошибка удаления записи")
				return
			}
			break
		}
	} else {
		manager := crud.Call([]reflect.Value{reflect.ValueOf(db)})[0].Interface().(types.CRUDManager)
		switch r.Method {
		case "PUT":
			if err := manager.Create(fields); err != nil {
				common.ErrorResponse(w, r, err.Error())
				return
			}
			break
		case "POST":
			if err := manager.Update(fields); err != nil {
				common.ErrorResponse(w, r, err.Error())
				return
			}
			break
		case "DELETE":
			if err := manager.Delete(); err != nil {
				common.ErrorResponse(w, r, err.Error())
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

func FilteredContentType(w http.ResponseWriter, r *http.Request) {
	contentTypeID, _ := strconv.Atoi(r.FormValue("content_type_id"))
	fieldName := r.FormValue("field_name")
	fieldValue := r.FormValue("value")
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
		field := obj.Elem().FieldByName(fieldName)
		fieldKind := field.Kind()
		modelSlice := reflect.New(reflect.SliceOf(obj.Elem().Type())).Interface()

		if fieldKind == reflect.Slice {
			elType := reflect.TypeOf(field.Interface()).Elem()
			if elType.Kind() == reflect.Struct {
				modelStructFields := db.NewScope(model).GetStructFields()
				fmt.Print(elType.Name())
				var groupStructField *gorm.StructField
				for _, f := range modelStructFields {
					if f.Relationship != nil && f.Relationship.Kind == "many_to_many" {
						joinTableHandler := f.Relationship.JoinTableHandler.(*gorm.JoinTableHandler)
						if joinTableHandler != nil && joinTableHandler.Destination.ModelType.Name() == elType.Name() {
							groupStructField = f
							break
						}
					}
				}
				if groupStructField != nil {
					tableName := groupStructField.Relationship.JoinTableHandler.Table(db)
					foreignDBName := groupStructField.Relationship.ForeignDBNames[0]
					assDBName := groupStructField.Relationship.AssociationForeignDBNames[0]
					groupID, _ := strconv.Atoi(fieldValue)
					db = db.Preload(r.FormValue("field_name"), "id = ?", groupID)
					db.Joins(
						fmt.Sprintf("INNER JOIN \"%s\" as t ON t.%s = id", tableName, foreignDBName),
					).Find(modelSlice, fmt.Sprintf("t.%s = ?", assDBName), groupID)
				}
			}
		} else if fieldKind == reflect.Struct {
			fieldName = helpers.ToSnakeCase(fieldName) + "_id"
			db = db.Preload(r.FormValue("field_name"))
			db.Find(modelSlice, fmt.Sprintf("\"%s\" = ?", fieldName), fieldValue)
		}

		data := make([]types.H, 0)
		nObj := reflect.ValueOf(modelSlice).Elem()
		addGroupShortInfoToData(nObj, r.FormValue("field_name"), nil, &data)
		common.JSONResponse(w, data)
	} else {
		common.Forbidden(w, r)
	}
}

func AllFieldsInModel(w http.ResponseWriter, r *http.Request) {
	contentTypeID, _ := strconv.Atoi(r.FormValue("content_type_id"))
	searchKeyWords := r.FormValue("keyword")
	groupBy := r.FormValue("group_by")
	orderBy := r.FormValue("order_by")
	pageSize, _ := strconv.Atoi(r.FormValue("page_size"))
	if pageSize == 0 {
		pageSize = 100
	} else if pageSize > 500 {
		pageSize = 500
	}
	//filters := r.FormValue("filters")
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
		qs := db
		if groupBy != "" {
			qs = qs.Preload(groupBy)
		}
		obj := reflect.New(reflect.TypeOf(model))
		adminMetaMethod := obj.MethodByName("Admin")
		var adminMeta types.AdminMeta
		if adminMetaMethod.Kind() != reflect.Invalid {
			adminMeta = adminMetaMethod.Call(nil)[0].Interface().(types.AdminMeta)
			keywords := strings.Split(searchKeyWords, " ")
			for _, fieldName := range adminMeta.SearchFields {
				isOr := len(keywords) < 2
				for _, keyword := range keywords {
					if keyword == "" {
						continue
					}
					fieldName = helpers.ToSnakeCase(fieldName)
					if !isOr {
						qs = qs.Where(fmt.Sprintf("%s ilike ?", fieldName), "%"+keyword+"%")
						isOr = true
					} else {
						qs = qs.Or(fmt.Sprintf("%s ilike ?", fieldName), "%"+keyword+"%")
					}
				}
			}

			if orderBy == "" {
				for _, fieldName := range adminMeta.OrderBy {
					fieldName = helpers.GetSortField(fieldName)
					qs = qs.Order(fieldName)
				}
			} else {
				for _, ord := range strings.Split(orderBy, ",") {
					for _, fieldInfo := range adminMeta.SortFields {
						if strings.Index(ord, helpers.ToSnakeCase(fieldInfo.Name)) > -1 {
							ord = helpers.GetSortField(ord)
							qs = qs.Order(ord)
						}
					}
				}
			}
			// TODO
			//for _, filter := range strings.Split(filters, ",") {
			//	for _, filterField := range adminMeta.FilterFields {
			//		if strings.Index(filter, helpers.ToSnakeCase(filterField.Name)) > -1 {
			//
			//		}
			//	}
			//}
		}
		modelSlice := reflect.New(reflect.SliceOf(obj.Elem().Type())).Interface()
		paginateData := common.Paginate(modelSlice, qs, page, 100, adminMeta.Preload, false)
		methodGetMeta := obj.MethodByName("Meta")
		meta := make(map[string]interface{})
		if methodGetMeta.Kind() != reflect.Invalid {
			modelMeta := methodGetMeta.Call(nil)[0].Interface().(types.ModelsMeta)
			meta["name"] = modelMeta.Name
			meta["plural"] = modelMeta.Plural
		} else {
			meta["name"] = obj.Type().Name()
			meta["plural"] = obj.Type().Name() + "s"
		}
		var result []types.H
		short, _ := strconv.ParseBool(r.FormValue("short"))
		isShort := short || !(len(adminMeta.SortFields) > 0 || len(adminMeta.ExtraFields) > 0)
		groups := make(map[uint]string)
		if groupBy != "" {
			addGroupShortInfoToData(reflect.ValueOf(paginateData.Result), groupBy, &groups, &result)
		} else {
			for _, item := range paginateData.Result {
				iobj := reflect.ValueOf(item)
				if isShort {
					addShortInfoToData(iobj, &result)
				} else {
					item := make(types.H)
					item["id"] = iobj.FieldByName("Model").FieldByName("ID").Interface()
					for _, fieldInfo := range adminMeta.SortFields {
						name := helpers.ToSnakeCase(fieldInfo.Name)
						val := iobj.FieldByName(fieldInfo.Name)
						if val.IsValid() {
							item[name] = iobj.FieldByName(fieldInfo.Name).Interface()
						} else {
							item[name] = nil
						}
					}
					for _, extraField := range adminMeta.ExtraFields {
						name, value := helpers.GetValue(iobj, extraField.Name)
						item[name] = value
					}
					result = append(result, item)
				}
			}
		}

		meta["short"] = isShort
		meta["toHTML"] = adminMeta.ShortToHtml
		meta["available_search"] = len(adminMeta.SearchFields) > 0
		data := types.H{
			"meta":     meta,
			"paginate": paginateData.Paginate,
			"result":   result,
			"groups":   groups,
		}

		if !isShort {
			var sortFields, extraFields []types.H
			for _, fieldInfo := range adminMeta.SortFields {
				structField, _ := obj.Elem().Type().FieldByName(fieldInfo.Name)
				form := helpers.ParseTag(structField.Tag.Get("form"))
				name := helpers.ToSnakeCase(fieldInfo.Name)
				label := name
				if fieldInfo.Label == "" && form["label"] != "" {
					label = form["label"]
				} else if fieldInfo.Label != "" {
					label = fieldInfo.Label
				}
				sortFields = append(sortFields, types.H{
					"name":   name,
					"label":  label,
					"toHTML": fieldInfo.ToHTML,
				})
			}
			for _, extraField := range adminMeta.ExtraFields {
				_fieldNames := strings.Split(extraField.Name, ".")
				_extraField := types.H{
					"label": extraField.Label,
				}
				if len(_fieldNames) > 1 {
					name := strings.Builder{}
					for i, fieldName := range _fieldNames {
						if i > 0 {
							name.Write([]byte("."))
						}
						name.Write([]byte(helpers.ToSnakeCase(fieldName)))
					}
					_extraField["name"] = name.String()
				} else {
					_extraField["name"] = helpers.ToSnakeCase(extraField.Name)
				}
				_extraField["toHTML"] = extraField.ToHTML
				extraFields = append(extraFields, _extraField)
			}
			data["sort_fields"] = sortFields
			data["extra_fields"] = extraFields
		}
		common.JSONResponse(w, data)
	} else {
		common.Forbidden(w, r)
	}
}
