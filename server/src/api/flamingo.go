package api

import (
	"context"
	"saitho.me/styx-app/src/service"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type Module struct{}

// Configure is the default Method a Module needs to implement
func (m *Module) Configure(injector *dingo.Injector) {
	// register our routes struct as a router Module - so that it is "known" to the router service
	web.BindRoutes(injector, new(routes))
}

// routes struct - our internal struct that gets the interface methods for router.Module
type routes struct {
	// controller - we will defined routes that are handled by our HelloController - so we need this as a dependency
	controller     *ApiController
	serviceManager *service.ServiceManager
}

// Inject dependencies - this is called by Dingo and gets an initializes instance of the HelloController passed automatically
func (r *routes) Inject(controller *ApiController, serviceManager *service.ServiceManager) *routes {
	r.controller = controller
	r.serviceManager = serviceManager
	return r
}

// Routes method which defines all routes handlers in service
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/api", "api.index")
	registry.HandleGet("api.index", r.controller.Index)

	registry.HandleData("serviceManager", func(ctx context.Context, req *web.Request, callParams web.RequestParams) interface{} {
		return r.serviceManager
	})
}
