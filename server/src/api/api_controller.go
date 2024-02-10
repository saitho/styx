package api

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"

	"saitho.me/styx-app/src/lib"
)

type (
	ApiController struct {
		responder *web.Responder
		cfg       *lib.Config
	}

	apiViewData struct {
		Services map[string]*lib.StyxService
	}
)

// Inject dependencies
func (controller *ApiController) Inject(responder *web.Responder, config *lib.Config) *ApiController {
	controller.responder = responder
	controller.cfg = config
	return controller
}

// Index is a controller action that renders Data
func (controller *ApiController) Index(_ context.Context, r *web.Request) web.Result {
	return controller.responder.Data(apiViewData{
		Services: controller.cfg.Services,
	})
}
