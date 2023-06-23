package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/dummynotes/notes/internal/models"
)

type DynamoDBClient struct {
	db *dynamodb.Client
}

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
		TableName: aws.String("Notes"),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("PutItem: %v\n", err)
	}

	return item, nil
}

//	func (c *DynamoDBClient) List() (interface{}, error) {
//		// Implementation of listing items in DynamoDB
//
// c.JSON(http.StatusCreated, gin.H{"notes": []interface{}{note}})
// }
func (c *DynamoDBClient) Get(id string) (interface{}, error) {
	note := models.Note{}

	selectedKeys := map[string]string{
		"NoteId": id,
	}

	key, _ := attributevalue.MarshalMap(selectedKeys)

	data, err := c.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Notes"),
		Key:       key,
	})

	if err != nil {
		return note, fmt.Errorf("GetItem: %v\n", err)
	}

	if data.Item == nil {
		return note, fmt.Errorf("GetItem: Data not found.\n")
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
		"NoteId": id,
	}

	key, _ := attributevalue.MarshalMap(selectedKeys)

	_, err := c.db.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Notes"),
		Key:       key,
	})

	if err != nil {
		return false, fmt.Errorf("DeleteItem: %v\n", err)
	}

	return true, nil
}
