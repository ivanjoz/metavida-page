package main

import (
	"metavida/backend/core"
	"metavida/backend/operations"
	"metavida/backend/security"
)

var appHandlersModules = []core.AppRouterType{
	security.ModuleHandlers,
	operations.ModuleHandlers,
	systemHandlers,
}

var systemHandlers = core.AppRouterType{
	"GET.health": healthHandler,
}
