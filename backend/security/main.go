package security

import "metavida/backend/core"

var ModuleHandlers = core.AppRouterType{
	"POST.p-user-login": PostLogin,
	"POST.login":        PostLogin,
}
