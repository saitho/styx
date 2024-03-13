package service

import (
	"context"
	"fmt"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
	t "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	sdkClient "github.com/docker/docker/client"
)

type DockerClient struct {
	api sdkClient.CommonAPIClient
}

func (client *DockerClient) ListContainers(logger flamingo.Logger) ([]t.Container, error) {
	dockerClient, err := sdkClient.NewClientWithOpts(
		sdkClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(fmt.Sprintf("Connected to Docker v%s", dockerClient.ClientVersion()))
	return dockerClient.ContainerList(
		context.Background(),
		container.ListOptions{
			Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: "me.saitho.styx.service=1"}),
		})
}

// DiscoverServicesByDocker looks for running containers with styx label.
// The container name will be used as service name
func DiscoverServicesByDocker(logger flamingo.Logger) []StyxService {
	var services []StyxService
	containers, _ := new(DockerClient).ListContainers(logger)
	if len(containers) == 0 {
		logger.Debugf("No tagged containers found")
		return services
	}
	for _, t2 := range containers {
		if t2.State != "running" {
			// For now only support services running at startup of this application
			logger.Debugf("Skip container with image \"%s\" (id %s) as it is not running", t2.Image, t2.ID)
			continue
		}
		logger.Info(fmt.Sprintf("Found tagged running container with image \"%s\" (id %s)", t2.Image, t2.ID))
		// Only use first name of container
		// todo: make port configurable via label
		for name, network := range t2.NetworkSettings.Networks {
			// for now only support bridge network which is Docker default
			if name != "bridge" {
				continue
			}
			services = append(services, StyxService{
				ServiceName: strings.TrimPrefix(t2.Names[0], "/"),
				IpAddress:   network.IPAddress,
			})
		}
	}
	return services
}
