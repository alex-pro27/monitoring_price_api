package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
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

func SetFieldsForModel(model interface{}, data map[string]interface{}) error {
	obj := reflect.ValueOf(model)
	errs := make(map[string]string)
	for key, value := range data {
		value := reflect.ValueOf(value)
		name := ToCamelCase(key)
		kind := obj.Elem().FieldByName(name).Kind()
		if kind != reflect.Invalid && kind != reflect.Slice {
			if value.Kind() != obj.Elem().FieldByName(name).Kind() {
				var _value interface{}
				strValue := fmt.Sprintf("%v", value.Interface())
				switch obj.Elem().FieldByName(name).Interface().(type) {
				case uint:
					_value, _ = strconv.Atoi(strValue)
					_value = uint(_value.(int))
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
				obj.Elem().FieldByName(name).Set(value)
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
