package service

import (
	"context"
	"fmt"
	"golang.org/x/exp/maps"
	"slices"

	"flamingo.me/flamingo/v3/framework/web"

	"saitho.me/styx-app/src/lib"
)

type (
	Controller struct {
		responder *web.Responder
		cfg       *lib.Config
	}
)

// Inject dependencies
func (controller *Controller) Inject(responder *web.Responder, config *lib.Config) *Controller {
	controller.responder = responder
	controller.cfg = config
	return controller
}

// Index is a controller action that renders Data
func (controller *Controller) Index(_ context.Context, r *web.Request) web.Result {
	return controller.responder.Render("services", nil)
}

// Service is a controller action that renders Service details
func (controller *Controller) Service(_ context.Context, r *web.Request) web.Result {
	serviceName := r.Params["service"]
	if serviceName == "" {
		return controller.responder.Forbidden(fmt.Errorf("missing service name"))
	}
	if !slices.Contains(maps.Keys(controller.cfg.Services), serviceName) {
		return controller.responder.Forbidden(fmt.Errorf("service \"%s\" not found", serviceName))
	}
	return controller.responder.Render("service_details", nil)
}
