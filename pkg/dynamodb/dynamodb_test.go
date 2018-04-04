package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func TestDynamoDBWrite(t *testing.T) {
	ctx := context.Background()
	instance, err := Run(ctx)
	if err != nil {
		t.Fatalf("Unable to start dynamodb: %v", err)
	}
	defer instance.Stop(ctx)

	cli := dynamodb.New(session.New(instance.Config()))

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("last_update"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("last_update"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("test_table"),
	}

	_, err = cli.CreateTable(input)
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}

	item, _ := dynamodbattribute.MarshalMap(map[string]interface{}{"id": "test-entry", "last_update": 123})

	_, err = cli.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("test_table"),
	})
	if err != nil {
		t.Errorf("unable to put item: %v", err)
	}

}
