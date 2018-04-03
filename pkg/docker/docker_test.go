package docker

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	docker "github.com/docker/docker/client"
)

func TestContainer(t *testing.T) {
	c := &Container{
		Image: "docker.io/crccheck/hello-world",
	}

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		t.Errorf("unable to create docker client: %v", err)
	}

	err = c.Run(ctx)
	if err != nil {
		t.Fatalf("unable to start docker container: %v", err)
	}

	_, err = cli.ContainerInspect(ctx, c.containerID)
	if err != nil {
		t.Errorf("unable to inspect container: %v", err)
	}

	hostPort, err := c.GetPort(ctx, "8000/tcp")
	if err != nil {
		t.Errorf("Unable to get host port for container: %v", err)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%s", hostPort))
	if err != nil {
		t.Errorf("Unable to get hello world: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Got %d, expecting %d", resp.StatusCode, http.StatusOK)
	}

	err = c.Stop(ctx)
	if err != nil {
		t.Fatalf("unable to stop docker container: %v", err)
	}

	_, err = cli.ContainerInspect(ctx, c.containerID)
	if err == nil {
		t.Errorf("expected container inspect to fail")
	}

}
