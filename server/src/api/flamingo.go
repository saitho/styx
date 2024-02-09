package api

import (
	"context"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"

	"saitho.me/styx-app/src/lib"
)

type Module struct{}

// Configure is the default Method a Module needs to implement
func (m *Module) Configure(injector *dingo.Injector) {
	// register our routes struct as a router Module - so that it is "known" to the router module
	web.BindRoutes(injector, new(routes))
}

// routes struct - our internal struct that gets the interface methods for router.Module
type routes struct {
	// controller - we will defined routes that are handled by our HelloController - so we need this as a dependency
	controller *ApiController
	config     *lib.Config
}

// Inject dependencies - this is called by Dingo and gets an initializes instance of the HelloController passed automatically
func (r *routes) Inject(controller *ApiController, config *lib.Config) *routes {
	r.controller = controller
	r.config = config
	return r
}

// Routes method which defines all routes handlers in module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/api", "api.index")
	registry.HandleGet("api.index", r.controller.Index)

	registry.HandleData("config", func(ctx context.Context, req *web.Request, callParams web.RequestParams) interface{} {
		return r.config
	})
}
