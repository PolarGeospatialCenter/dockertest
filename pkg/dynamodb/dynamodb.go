package dynamodb

import (
	"context"
	"fmt"

	"github.com/PolarGeospatialCenter/dockertest/pkg/docker"
	"github.com/aws/aws-sdk-go/aws"
)

type Instance struct {
	config *aws.Config
	*docker.Container
}

func Run(ctx context.Context) (*Instance, error) {
	i := &Instance{
		Container: &docker.Container{Image: "docker.io/deangiberson/aws-dynamodb-local"},
	}

	err := i.Container.Run(ctx)

	port, err := i.GetPort(ctx, "8000/tcp")
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://localhost:%s", port)
	region := "us-east-2"
	i.config = &aws.Config{Endpoint: &endpoint, Region: &region}
	return i, nil
}

func (i *Instance) Config() *aws.Config {
	return i.config
}
