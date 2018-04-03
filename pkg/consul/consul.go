package consul

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	timeout := time.After(10 * time.Second)
	checkInterval := time.Tick(50 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("consul failed to start after 10 seconds")
		case <-checkInterval:
			if i.running() {
				return i, nil
			}
		}
	}
}

func (i *Instance) running() bool {
	c := http.Client{}
	resp, err := c.Get(fmt.Sprintf("http://%s", i.Config().Address))
	return err == nil && resp.StatusCode == 200
}

func (i *Instance) Config() *consul.Config {
	return i.config
}
