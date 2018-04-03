package vaulttest

import (
	"context"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
)

const vaultTestRootToken = "701432d1-00e7-7c94-10c4-8450ab3c4b31"

// Run starts a vault instance and returns a *vault.Config and the root token
func Run(ctx context.Context) (*vault.Config, string) {
	c := &TestContainer{
		Image: "vault",
		Cmd:   []string{"vault", "server", "-dev", "-dev-root-token-id", vaultTestRootToken, "-dev-listen-address", "0.0.0.0:8200"},
	}

	c.Run(ctx)
	defer c.Stop(ctx)
	port, err := c.GetPort(ctx, "8200/tcp")
	if err != nil {
		return nil, ""
	}

	config := vault.DefaultConfig()
	config.Address = fmt.Sprintf("http://127.0.0.1:%s", port)
	time.Sleep(500 * time.Millisecond)
	return config, vaultTestRootToken
}
