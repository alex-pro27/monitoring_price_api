package common

import "net/http"

type HttpHandler func(http.ResponseWriter, *http.Request)

type H map[string]interface{}

type CallBack func(args... interface{})

type Model interface {
	Serializer() H
}