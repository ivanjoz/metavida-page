package operations

import "metavida/backend/core"

var ModuleHandlers = core.AppRouterType{
	"POST.p-clients": PostClients,
	"POST.clients":   PostClients,
	"GET.clients":    GetClients,
}
