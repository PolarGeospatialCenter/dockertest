package vaulttest

import (
	"context"
	"fmt"
	"log"

	docker "docker.io/go-docker"
	dockertypes "docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/network"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/docker/go-connections/nat"
)

type TestContainer struct {
	Image       string
	Cmd         []string
	containerID string
}

func (t *TestContainer) Run(ctx context.Context) error {
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

	c, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, "")
	if err != nil {
		return err
	}

	t.containerID = c.ID
	err = cli.ContainerStart(ctx, c.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (t *TestContainer) GetPort(ctx context.Context, port string) (string, error) {
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

func (t *TestContainer) Stop(ctx context.Context) error {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	log.Printf("Killing container %s", t.containerID)
	err = cli.ContainerKill(ctx, t.containerID, "SIGINT")
	if err != nil {
		return fmt.Errorf("error killing container %s: %v", t.containerID, err)
	}

	log.Printf("Waiting for container %s", t.containerID)
	ok, errCh := cli.ContainerWait(ctx, t.containerID, container.WaitConditionNotRunning)
	select {
	case <-ok:
		log.Printf("Container stopped")
		break
	case e := <-errCh:
		return fmt.Errorf("error waiting for container to exit: %v", e)
	}

	log.Printf("Removing container %s", t.containerID)
	err = cli.ContainerRemove(ctx, t.containerID, dockertypes.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("error removing container %s: %v", t.containerID, err)
	}
	return nil
}

func startDynamoDB(ctx context.Context) (*dynamodb.DynamoDB, error) {
	c := &TestContainer{Image: "deangiberson/aws-dynamodb-local"}

	c.Run(ctx)
	defer c.Stop(ctx)
	port, err := c.GetPort(ctx, "8000/tcp")
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://localhost:%s", port)
	region := "us-east-2"
	db := dynamodb.New(session.New(&aws.Config{Endpoint: &endpoint, Region: &region}))
	return db, nil
}
