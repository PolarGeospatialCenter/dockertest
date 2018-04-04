package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestDynamoDB(t *testing.T) {
	ctx := context.Background()
	instance, err := Run(ctx)
	if err != nil {
		t.Fatalf("Unable to start dynamodb: %v", err)
	}
	defer instance.Stop(ctx)

	cli := dynamodb.New(session.New(instance.Config()))
	out, err := cli.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		t.Fatalf("Unable to list tables: %v", err)
	}
	if len(out.TableNames) > 0 {
		t.Errorf("Expecting 0 tables, got %d", len(out.TableNames))
	}
}
