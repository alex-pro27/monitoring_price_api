package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const ISO8601 = "2006-01-02T15:04:05"

func ToCamelCase(str string) string {
	link := regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}

func GetSortField(fieldName string) string {
	if match, _ := regexp.MatchString("^-", fieldName); match {
		fieldName = ToSnakeCase(fieldName[1:]) + " desc"
	} else {
		fieldName = ToSnakeCase(fieldName)
	}
	return fieldName
}

func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func IsZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func GetValue(value reflect.Value, fieldName string) (string, interface{}) {
	_fieldNames := strings.Split(fieldName, ".")
	name := strings.Builder{}
	var _getValue func(string, reflect.Value) (string, reflect.Value)
	_getValue = func(name string, obj reflect.Value) (string, reflect.Value) {
		if obj.IsValid() {
			method := obj.MethodByName(name)
			if method.Kind() != reflect.Invalid {
				obj = method.Call(nil)[0]
			} else {
				if obj.Kind() == reflect.Ptr {
					obj = obj.Elem()
				}
				obj = obj.FieldByName(name)
			}
		}
		return ToSnakeCase(name), obj
	}
	for i, _fieldName := range _fieldNames {
		_name, _value := _getValue(_fieldName, value)
		value = _value
		if i > 0 {
			name.Write([]byte("."))
		}
		name.Write([]byte(_name))
	}
	var val interface{}
	if value.IsValid() {
		val = value.Interface()
	}
	return name.String(), val
}

func SetManyToMany(db *gorm.DB, model interface{}, data map[string]interface{}) {
	obj := reflect.ValueOf(model)
	for key, value := range data {
		sl := reflect.ValueOf(value)
		name := ToCamelCase(key)
		if sl.IsValid() && sl.Kind() != reflect.Slice {
			continue
		}
		field := obj.Elem().FieldByName(name)
		if field.Kind() != reflect.Slice {
			continue
		}
		relatedModels := reflect.New(field.Type()).Interface()
		db.Find(relatedModels, "id IN (?)", value)
		db.Model(model).Association(name).Replace(relatedModels)
		field.Set(reflect.ValueOf(relatedModels).Elem())
	}
}

func SetFieldsForModel(model interface{}, data map[string]interface{}) error {
	obj := reflect.ValueOf(model)
	errs := make(map[string]string)
	for key, value := range data {
		value := reflect.ValueOf(value)
		name := ToCamelCase(key)
		kind := obj.Elem().FieldByName(name).Kind()

		if kind != reflect.Invalid {
			if value.Kind() != obj.Elem().FieldByName(name).Kind() && kind != reflect.Slice {
				var _value interface{}
				var strValue string
				if value.IsValid() {
					strValue = fmt.Sprintf("%v", value.Interface())
				}
				switch obj.Elem().FieldByName(name).Interface().(type) {
				case uint:
					_value, _ = strconv.Atoi(strValue)
					_value = uint(_value.(int))
					break
				case int:
					_value, _ = strconv.Atoi(strValue)
					break
				case int32:
					_value, _ = strconv.ParseInt(strValue, 10, 32)
					break
				case int64:
					_value, _ = strconv.ParseInt(strValue, 10, 64)
					break
				case float32:
					_value, _ = strconv.ParseFloat(strValue, 32)
					break
				case float64:
					_value, _ = strconv.ParseFloat(strValue, 64)
					break
				case time.Time:
					_value, _ = time.Parse(ISO8601, strValue)
					break
				default:
					continue
				}
				value = reflect.ValueOf(_value)
			}
			method := obj.MethodByName("Set" + name)
			if method.Kind() != reflect.Invalid {
				e := method.Call([]reflect.Value{value})[0].Interface()
				if e != nil {
					errs[key] = e.(error).Error()
				}
			} else {
				if kind != reflect.Slice {
					obj.Elem().FieldByName(name).Set(value)
				}
			}
		}
	}

	for i := 0; i < obj.Elem().NumField(); i++ {
		fieldType := obj.Elem().Type().Field(i)
		field := obj.Elem().Field(i)
		key := ToSnakeCase(fieldType.Name)
		if IsZero(field) {
			form := ParseTag(fieldType.Tag.Get("form"))
			required, _ := strconv.ParseBool(form["required"])
			if required {
				errs[key] = "обязательное поле"
			}
		}
	}

	var err error
	if len(errs) > 0 {
		messageByte, _ := json.Marshal(errs)
		message := string(messageByte)
		err = errors.New(message)
	}
	return err
}

func ParseTag(tag string) map[string]string {
	data := make(map[string]string)
	options := strings.Split(tag, ";")
	for _, item := range options {
		opt := strings.Split(item, ":")
		if len(opt) > 1 {
			data[opt[0]] = opt[1]
		} else if len(opt) == 1 {
			data[opt[0]] = "true"
		}
	}
	return data
}
