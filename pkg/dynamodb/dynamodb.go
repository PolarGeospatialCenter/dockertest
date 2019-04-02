package dynamodb

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PolarGeospatialCenter/dockertest/pkg/docker"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Instance struct {
	config *aws.Config
	*docker.Container
}

func Run(ctx context.Context) (*Instance, error) {
	i := &Instance{
		Container: &docker.Container{Image: "docker.io/amazon/dynamodb-local"},
	}

	err := i.Container.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start dynamodb: %v", err)
	}

	port, err := i.GetPort(ctx, "8000/tcp")
	if err != nil {
		return nil, err
	}

	i.config = &aws.Config{}
	i.config.WithEndpoint(fmt.Sprintf("http://localhost:%s", port))
	i.config.WithRegion("us-east-2")
	i.config.WithCredentials(credentials.NewStaticCredentials("fake_id", "bad_secret", "bad_token"))

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
	endpoint := i.Config().Endpoint
	resp, err := c.Get(*endpoint)
	return err == nil && resp.StatusCode == http.StatusBadRequest
}

func (i *Instance) Config() *aws.Config {
	return i.config
}
