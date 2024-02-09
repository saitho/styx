package module

import (
	"context"
	"fmt"
	"golang.org/x/exp/maps"
	"slices"

	"flamingo.me/flamingo/v3/framework/web"

	"saitho.me/styx-app/src/lib"
)

type (
	ModuleController struct {
		responder *web.Responder
		cfg       *lib.Config
	}
)

// Inject dependencies
func (controller *ModuleController) Inject(responder *web.Responder, config *lib.Config) *ModuleController {
	controller.responder = responder
	controller.cfg = config
	return controller
}

// Index is a controller action that renders Data
func (controller *ModuleController) Index(_ context.Context, r *web.Request) web.Result {
	return controller.responder.Render("modules", nil)
}

// Module is a controller action that renders Module details
func (controller *ModuleController) Module(_ context.Context, r *web.Request) web.Result {
	moduleName := r.Params["module"]
	if moduleName == "" {
		return controller.responder.Forbidden(fmt.Errorf("missing module name"))
	}
	if !slices.Contains(maps.Keys(controller.cfg.Modules), moduleName) {
		return controller.responder.Forbidden(fmt.Errorf("module \"%s\" not found", moduleName))
	}
	return controller.responder.Render("module_details", nil)
}
