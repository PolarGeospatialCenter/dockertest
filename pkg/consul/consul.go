package consul

import (
	"context"
	"fmt"

	"github.com/PolarGeospatialCenter/dockertest/pkg/docker"
	consul "github.com/hashicorp/consul/api"
)

type Instance struct {
	config *consul.Config
	*docker.Container
}

func Run(ctx context.Context) (*Instance, error) {
	i := &Instance{
		Container: &docker.Container{
			Image: "docker.io/library/consul",
			Cmd:   []string{"consul", "agent", "-dev", "-client", "0.0.0.0"},
		},
	}

	err := i.Container.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start consul: %v", err)
	}

	port, err := i.GetPort(ctx, "8500/tcp")
	if err != nil {
		return nil, err
	}

	i.config = consul.DefaultConfig()
	i.config.Address = fmt.Sprintf("localhost:%s", port)
	return i, nil
}

func (i *Instance) Config() *consul.Config {
	return i.config
}
