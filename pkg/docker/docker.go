package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Container struct {
	Image       string
	Cmd         []string
	containerID string
}

func (t *Container) Run(ctx context.Context) error {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	containerConfig := &container.Config{
		Image: t.Image,
		Cmd:   t.Cmd,
	}

	hostConfig := &container.HostConfig{
		PublishAllPorts: true,
	}

	networkConfig := &network.NetworkingConfig{}

	status, err := cli.ImagePull(ctx, containerConfig.Image, dockertypes.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("unable to pull image: %v", err)
	}

	_, err = ioutil.ReadAll(status)
	if err != nil {
		return fmt.Errorf("unable to read status from image pull: %v", err)
	}
	status.Close()

	c, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, "")
	if err != nil {
		return fmt.Errorf("unable to create container: %v", err)
	}

	t.containerID = c.ID
	err = cli.ContainerStart(ctx, c.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (t *Container) GetPort(ctx context.Context, port string) (string, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return "", err
	}

	data, err := cli.ContainerInspect(ctx, t.containerID)
	if err != nil {
		return "", err
	}

	ports := data.NetworkSettings.Ports
	hostPort := ports[nat.Port(port)][0].HostPort

	return hostPort, nil
}

func (t *Container) Stop(ctx context.Context) error {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	err = cli.ContainerKill(ctx, t.containerID, "SIGKILL")
	if err != nil {
		return fmt.Errorf("error killing container %s: %v", t.containerID, err)
	}

	toCtx, _ := context.WithTimeout(ctx, 1*time.Second)
	_, err = cli.ContainerWait(toCtx, t.containerID)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, t.containerID, dockertypes.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("error removing container %s: %v", t.containerID, err)
	}
	return nil
}
