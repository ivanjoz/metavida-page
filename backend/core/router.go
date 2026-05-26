package core

import "net/http"

type HandlerFunc func(*HandlerArgs) HandlerResponse

type HandlerArgs struct {
	Body           *string
	Query          map[string]string
	Headers        map[string]string
	Route          string
	Method         string
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	StartTime      int64
	Authorization  string
	User           *UsuarioToken
}

type HandlerResponse struct {
	Body       any
	Error      string
	StatusCode int
}

type AppRouterType map[string]HandlerFunc
