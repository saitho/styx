package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
	t "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	sdkClient "github.com/docker/docker/client"
)

type DockerClient struct {
}

func isRunningInDocker() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func FindLocalNetworks(dockerClient *sdkClient.Client, logger flamingo.Logger) ([]*t.NetworkResource, error) {
	var foundNetworks []*t.NetworkResource
	if !isRunningInDocker() {
		logger.Debug("Looking for local network but not running in Docker")
		return foundNetworks, nil
	}

	name, err := os.Hostname()
	if err != nil {
		logger.Error(err)
		return foundNetworks, err
	}
	addrs, err := net.LookupHost(name)
	if err != nil {
		logger.Error(err)
		return foundNetworks, err
	}
	localIp := addrs[0]

	networks, err := dockerClient.NetworkList(context.Background(), t.NetworkListOptions{})
	if err != nil {
		logger.Error(err)
		return foundNetworks, err
	}
	for _, network := range networks {
		// Network in NetworkList does not have containers attached to it, inspect again
		nw, _ := dockerClient.NetworkInspect(context.Background(), network.ID, t.NetworkInspectOptions{})
		for _, c := range nw.Containers {
			if strings.TrimSuffix(c.IPv4Address, "/24") == localIp {
				foundNetworks = append(foundNetworks, &nw)
			}
		}
	}
	return foundNetworks, nil
}

func ListContainers(logger flamingo.Logger) ([]t.Container, error) {
	dockerClient, err := sdkClient.NewClientWithOpts(
		sdkClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	filter := filters.NewArgs(filters.KeyValuePair{Key: "label", Value: "me.saitho.styx.service=1"})
	// if Styx runs on Docker as well, limit network to the same network as Styx runs on,
	// as it can not access containers outside its own network
	ownNetworks, _ := FindLocalNetworks(dockerClient, logger)
	if len(ownNetworks) > 0 {
		for _, n := range ownNetworks {
			filter.Add("network", n.Name)
		}
	} else {
		filter.Add("network", "bridge")
	}

	logger.Info(fmt.Sprintf("Connected to Docker v%s", dockerClient.ClientVersion()))
	return dockerClient.ContainerList(
		context.Background(),
		container.ListOptions{
			Filters: filter,
		})
}

// DiscoverServicesByDocker looks for running containers with styx label.
// The container name will be used as service name
func DiscoverServicesByDocker(logger flamingo.Logger) []*StyxService {
	var services []*StyxService
	containers, _ := ListContainers(logger)
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
		// Only use first name of container
		serviceName := strings.TrimPrefix(t2.Names[0], "/")
		logger.Info(fmt.Sprintf("Found tagged running container \"%s\" (id %s)", serviceName, t2.ID))
		for name, network := range t2.NetworkSettings.Networks {
			ipAddress := serviceName
			if name == "bridge" {
				ipAddress = network.IPAddress
			}
			s := &StyxService{
				ServiceName: serviceName,
				IpAddress:   ipAddress,
			}
			portFromLabel, ok := t2.Labels["me.saitho.styx.port"]
			if ok {
				if portNumber, err := strconv.Atoi(portFromLabel); err == nil {
					s.Port = portNumber
				}
			}
			services = append(services, s)
		}
	}
	return services
}
