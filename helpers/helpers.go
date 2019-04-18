package helpers

import (
	"reflect"
	"regexp"
	"strings"
)

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

func SetFieldsOnModel(model interface{}, data map[string]interface{}, required bool) string {
	obj := reflect.ValueOf(model)
	var errs []error
	for key, value := range data {
		value := reflect.ValueOf(value)
		if required || !IsZero(value) {
			method := obj.MethodByName("Set" + key)
			if method.Kind() != reflect.Invalid {
				e := method.Call([]reflect.Value{value})[0].Interface()
				if e != nil {
					errs = append(errs, e.(error))
				}
			} else {
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
