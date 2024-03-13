package service

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type ServiceManager struct {
	Services map[string]*StyxService
	logger   flamingo.Logger
}

func (m *ServiceManager) Inject(logger flamingo.Logger) {
	m.logger = logger
}

func (m *ServiceManager) EmitEvent(eventName string, event interface{}) error {
	if len(m.Services) == 0 {
		return fmt.Errorf("no services defined")
	}
	for _, service := range m.Services {
		if !slices.Contains(service.Config.SubscribedEvents, eventName) {
			continue
		}
		m.logger.Info(fmt.Sprintf("Emitting event %s to %s", eventName, service.ServiceName))
		if err := service.EmitEvent(eventName, event); err != nil {
			return err
		}
	}
	return nil
}

type ConfigModule struct {
	injector *dingo.Injector
}

// DiscoverServicesByEnvironment evaluates the services from config and prepares service objects
func DiscoverServicesByEnvironment(config config.Map) []*StyxService {
	var services []*StyxService
	if serviceString, ok := config.Get("services"); serviceString != "" && ok {
		for _, serviceName := range strings.Split(serviceString.(string), ",") {
			s := &StyxService{ServiceName: serviceName}
			if strings.Contains(serviceName, ":") {
				// if service name contains port number, remove it from service name
				serviceSplit := strings.Split(serviceName, ":")
				s.ServiceName = serviceSplit[0]
				if portNumber, err := strconv.Atoi(serviceSplit[1]); err == nil {
					s.Port = portNumber
				}
			}
			s.IpAddress = s.ServiceName
			services = append(services, s)
		}
	}
	return services
}

func NewServiceManager(
	logger flamingo.Logger,
	cfg *struct {
		CompleteConfig config.Map `inject:"config:styx"`
	},
) *ServiceManager {
	m := &ServiceManager{}
	m.Inject(logger)
	m.Services = map[string]*StyxService{}

	prepareService := func(s *StyxService) {
		s.Inject(logger)
		if err := s.Init(); err != nil {
			logger.Error(err)
		}
		m.Services[s.ServiceName] = s
	}

	for _, s := range DiscoverServicesByDocker(logger) {
		prepareService(s)
	}
	for _, s := range DiscoverServicesByEnvironment(cfg.CompleteConfig) {
		prepareService(s)
	}
	return m
}

type Module struct{}

// Configure is the default Method a Module needs to implement
func (m *Module) Configure(injector *dingo.Injector) {
	// register our routes struct as a router Module - so that it is "known" to the router service
	web.BindRoutes(injector, new(routes))
	injector.Bind(new(ServiceManager)).ToProvider(NewServiceManager).In(dingo.Singleton)
}

// routes struct - our internal struct that gets the interface methods for router.Module
type routes struct {
	// controller - we will define routes that are handled by our HelloController - so we need this as a dependency
	controller     *Controller
	serviceManager *ServiceManager
}

// Inject dependencies - this is called by Dingo and gets an initializes instance of the HelloController passed automatically
func (r *routes) Inject(serviceController *Controller, serviceManager *ServiceManager) *routes {
	r.controller = serviceController
	r.serviceManager = serviceManager
	return r
}

// Routes method which defines all routes handlers in service
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/services", "services.index")
	registry.HandleGet("services.index", r.controller.Index)

	registry.MustRoute("/services/:service", "services.service")
	registry.HandleGet("services.service", r.controller.Service)

	registry.HandleData("serviceManager", func(ctx context.Context, req *web.Request, callParams web.RequestParams) interface{} {
		return r.serviceManager
	})
}
