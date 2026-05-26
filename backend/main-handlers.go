package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"metavida/backend/core"
	"metavida/backend/security"
)

var appHandlers = core.AppRouterType{}

func makeAppHandlers() *core.AppRouterType {
	if len(appHandlers) == 0 {
		for _, moduleHandlers := range appHandlersModules {
			for path, handler := range moduleHandlers {
				appHandlers[path] = handler
			}
		}
	}
	return &appHandlers
}

var apiNames = []string{"api", "go1", "go2", "go3", "go4", "go5"}

func LocalHandler(w http.ResponseWriter, request *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			core.Error(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error (Panic in LocalHandler): %v", r))
			fmt.Println(string(debug.Stack()))
		}
	}()

	const maxBodyBytes = int64(10 << 20)
	bodyReader := http.MaxBytesReader(w, request.Body, maxBodyBytes)
	bodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		core.Error(w, http.StatusRequestEntityTooLarge, "request body too large")
		return
	}
	body := string(bodyBytes)

	args := core.HandlerArgs{
		Body:           &body,
		Method:         strings.ToUpper(request.Method),
		Route:          request.URL.Path,
		ResponseWriter: w,
		Request:        request,
		StartTime:      time.Now().UnixMilli(),
		Query:          map[string]string{},
		Headers:        map[string]string{},
	}

	for key, values := range request.URL.Query() {
		args.Query[key] = strings.Join(values, ",")
	}
	for key, values := range request.Header {
		args.Headers[key] = strings.Join(values, ",")
	}

	mainHandler(&args)
}

func mainHandler(args *core.HandlerArgs) (response core.HandlerResponse) {
	defer func() {
		if r := recover(); r != nil {
			response = core.HandlerResponse{
				Error:      fmt.Sprintf("Internal Server Error (Panic): %v", r),
				StatusCode: http.StatusInternalServerError,
			}
			fmt.Println(string(debug.Stack()))
			prepareResponse(args, response)
		}
	}()

	route := strings.TrimPrefix(args.Route, "/")
	pathSegments := strings.Split(route, "/")
	if len(pathSegments) > 0 && contains(apiNames, pathSegments[0]) {
		route = strings.Join(pathSegments[1:], "/")
	}
	args.Route = route
	args.Authorization = headerValue(args.Headers, "Authorization", "authorization")

	funcPath := args.Method + "." + route
	isPublicPath := route == "health" || (len(route) > 2 && route[0:2] == "p-")
	if !isPublicPath {
		args.User = core.CheckUser(args, 0)
		if args.User.Error != "" {
			response = core.HandlerResponse{
				Error:      args.User.Error,
				StatusCode: http.StatusUnauthorized,
			}
			return prepareResponse(args, response)
		}
		if err := security.ValidateTokenUser(args.User); err != nil {
			response = core.HandlerResponse{
				Error:      err.Error(),
				StatusCode: http.StatusUnauthorized,
			}
			return prepareResponse(args, response)
		}
	} else {
		args.User = &core.UsuarioToken{}
	}

	handler, ok := appHandlers[funcPath]
	if !ok {
		response = core.HandlerResponse{
			Error:      "no handler for path: " + funcPath,
			StatusCode: http.StatusNotFound,
		}
		return prepareResponse(args, response)
	}

	response = handler(args)
	return prepareResponse(args, response)
}

func prepareResponse(args *core.HandlerArgs, response core.HandlerResponse) core.HandlerResponse {
	statusCode := response.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	if response.Error != "" {
		core.Error(args.ResponseWriter, statusCode, response.Error)
		return response
	}
	core.JSON(args.ResponseWriter, statusCode, response.Body)
	return response
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func headerValue(headers map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(headers[key]); value != "" {
			return value
		}
	}
	return ""
}

func healthHandler(_ *core.HandlerArgs) core.HandlerResponse {
	return core.HandlerResponse{
		Body: map[string]string{"status": "ok"},
	}
}
