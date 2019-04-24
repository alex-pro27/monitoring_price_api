package helpers

import (
	"reflect"
	"regexp"
	"strings"
)

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

func SetFieldsOnModel(model interface{}, data map[string]interface{}) string {
	obj := reflect.ValueOf(model)
	var errs []error
	for key, value := range data {
		value := reflect.ValueOf(value)
		key := ToCamelCase(key)
		method := obj.MethodByName("Set" + key)
		if method.Kind() != reflect.Invalid {
			e := method.Call([]reflect.Value{value})[0].Interface()
			if e != nil {
				errs = append(errs, e.(error))
			}
		} else {
			if obj.Elem().FieldByName(key).Kind() != reflect.Invalid {
				obj.Elem().FieldByName(key).Set(value)
			}
		}

	}
	message := strings.Builder{}
	for _, e := range errs {
		if e != nil {
			message.Write([]byte(e.Error()))
			message.Write([]byte("\n"))
		}
	}
	return message.String()
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
