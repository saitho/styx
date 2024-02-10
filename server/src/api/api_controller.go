package api

import (
	"context"
	"saitho.me/styx-app/src/service"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	ApiController struct {
		responder      *web.Responder
		serviceManager *service.ServiceManager
	}

	apiViewData struct {
		Services map[string]*service.StyxService
	}
)

// Inject dependencies
func (controller *ApiController) Inject(responder *web.Responder, serviceManager *service.ServiceManager) *ApiController {
	controller.responder = responder
	controller.serviceManager = serviceManager
	return controller
}

// Index is a controller action that renders Data
func (controller *ApiController) Index(_ context.Context, r *web.Request) web.Result {
	return controller.responder.Data(apiViewData{
		Services: controller.serviceManager.Services,
	})
}
