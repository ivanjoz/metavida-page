package operations

import "metavida/backend/core"

var ModuleHandlers = core.AppRouterType{
	"POST.p-clients":       PostClients,
	"POST.clients":         PostClients,
	"POST.p-contact-email": PostContactEmail,
	"GET.clients":          GetClients,
}
