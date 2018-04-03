package vault

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PolarGeospatialCenter/dockertest/pkg/docker"
	vault "github.com/hashicorp/vault/api"
)

const vaultTestRootToken = "701432d1-00e7-7c94-10c4-8450ab3c4b31"

type Instance struct {
	config *vault.Config
	*docker.Container
}

func Run(ctx context.Context) (*Instance, error) {
	instance := &Instance{
		Container: &docker.Container{
			Image: "docker.io/library/vault",
			Cmd:   []string{"vault", "server", "-dev", "-dev-root-token-id", vaultTestRootToken, "-dev-listen-address", "0.0.0.0:8200"},
		},
	}

	err := instance.Container.Run(ctx)
	if err != nil {
		return nil, err
	}

	port, err := instance.GetPort(ctx, "8200/tcp")
	if err != nil {
		return nil, err
	}

	instance.config = vault.DefaultConfig()
	instance.config.Address = fmt.Sprintf("http://127.0.0.1:%s", port)

	timeout := time.After(10 * time.Second)
	checkInterval := time.Tick(50 * time.Millisecond)
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("vault failed to start after 10 seconds")
		case <-checkInterval:
			if instance.running() {
				return instance, nil
			}
		}
	}
}

func (i *Instance) running() bool {
	c := http.Client{}
	resp, err := c.Get(fmt.Sprintf("%s/v1/sys/seal-status", i.Config().Address))
	return err == nil && resp.StatusCode == 200
}

func (i *Instance) Config() *vault.Config {
	return i.config
}

func (i *Instance) RootToken() string {
	return vaultTestRootToken
}
