package service

import (
	"context"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"

	"saitho.me/styx-app/src/lib"
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
	controller *Controller
	config     *lib.Config
}

// Inject dependencies - this is called by Dingo and gets an initializes instance of the HelloController passed automatically
func (r *routes) Inject(serviceController *Controller, config *lib.Config) *routes {
	r.controller = serviceController
	r.config = config
	return r
}

// Routes method which defines all routes handlers in service
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/services", "services.index")
	registry.HandleGet("services.index", r.controller.Index)

	registry.MustRoute("/services/:service", "services.service")
	registry.HandleGet("services.service", r.controller.Service)

	registry.HandleData("config", func(ctx context.Context, req *web.Request, callParams web.RequestParams) interface{} {
		return r.config
	})
}
