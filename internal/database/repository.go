package database

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dummynotes/notes/internal/models"
)

type DynamoDBClient struct {
	db *dynamodb.Client
}

const tablename = "Notes"

func GetDynamodbClient(region string) *dynamodb.Client {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(awsConfig, func(opt *dynamodb.Options) {
		opt.Region = region
	})
}

func (c *DynamoDBClient) Configure(db *dynamodb.Client) {
	c.db = db
}

func (c *DynamoDBClient) Create(entity interface{}) (interface{}, error) {
	av, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Record, %w", err)
	}

	item, err := c.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tablename),
		Item:      av,
	})
	if err != nil {
		fmt.Errorf("entity: ", entity)
		return nil, fmt.Errorf("PutItem: %v\n", err)
	}

	return item, nil
}

func (c *DynamoDBClient) List(userId string) (interface{}, error) {
	items, err := c.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tablename),
		IndexName:              aws.String("NoteIdUserIdIndex"),
		KeyConditionExpression: aws.String("UserId = :gsi1pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1pk": &types.AttributeValueMemberS{Value: userId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ListItems: %v\n", err)
	}

	return items, nil
}

func (c *DynamoDBClient) Get(id string) (interface{}, error) {
	note := models.Note{}

	selectedKeys := map[string]string{
		"NoteId": id,
	}

	key, _ := attributevalue.MarshalMap(selectedKeys)

	data, err := c.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key:       key,
	})

	if err != nil {
		return note, fmt.Errorf("GetItem: %v\n", err)
	}

	if data.Item == nil {
		return nil, nil
	}

	err = attributevalue.UnmarshalMap(data.Item, &note)
	if err != nil {
		return note, fmt.Errorf("UnmarshalMap: %v\n", err)
	}

	return note, nil
}

//	func (c *DynamoDBClient) Update(id interface{}, entity interface{}) (bool, error) {
//		// Implementation of updating an item in DynamoDB
//	}
func (c *DynamoDBClient) Delete(id string) (bool, error) {
	selectedKeys := map[string]string{
		"NoteID": id,
	}

	key, _ := attributevalue.MarshalMap(selectedKeys)

	_, err := c.db.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tablename),
		Key:       key,
	})

	if err != nil {
		return false, fmt.Errorf("DeleteItem: %v\n", err)
	}

	return true, nil
}
