package consul

import (
	"context"
	"testing"
	"time"

	consul "github.com/hashicorp/consul/api"
)

func TestConsul(t *testing.T) {
	ctx := context.Background()
	instance, err := Run(ctx)
	if err != nil {
		t.Fatalf("Unable to start consul: %s", err)
	}
	defer instance.Stop(ctx)

	client, err := consul.NewClient(instance.Config())
	if err != nil {
		t.Fatalf("Unable to create consul client: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	_, err = client.Status().Leader()
	if err != nil {
		t.Fatalf("Unable to get leader from consul: %s", err)
	}
}
