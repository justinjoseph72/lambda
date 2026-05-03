package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	util "first-excercise/internal/util"
)

var client *dynamodb.Client

type orderDbItem struct {
	OrderId string  `dynamodbav:"orderId"`
	Amount  float64 `dynamodbav:"amount"`
	Item    string  `dynamodbav:"item"`
}

func Init(cfg aws.Config) {
	client = dynamodb.NewFromConfig(cfg)
	log.Println("DynamoDB client initialized")
}

func PutOrder(ctx context.Context, order *util.Order) error {
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return fmt.Errorf("DYNAMODB_TABLE_NAME environment variable is not set")
	}

	orderDbItem := orderDbItem{
		OrderId: order.OrderId,
		Amount:  order.Amount,
		Item:    order.Item,
	}

	item, err := attributevalue.MarshalMap(orderDbItem)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %v", err)
	}

	log.Printf("Putting item %v\n", item)

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	log.Printf("Order %s saved to DynamoDB", order.OrderId)
	return nil
}

func ScanOrders(ctx context.Context) ([]util.Order, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE_NAME environment variable is not set")
	}

	output, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan DynamoDB table: %v", err)
	}

	var orderDbItems []orderDbItem
	err = attributevalue.UnmarshalListOfMaps(output.Items, &orderDbItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB items: %v", err)
	}

	var orders []util.Order
	for _, item := range orderDbItems {
		orders = append(orders, util.Order{
			OrderId: item.OrderId,
			Amount:  item.Amount,
			Item:    item.Item,
		})
	}

	log.Printf("Scanned %d orders from DynamoDB", len(orders))
	return orders, nil
}

func GetOrder(ctx context.Context, orderId string) (*util.Order, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE_NAME environment variable is not set")
	}

	output, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"orderId": &types.AttributeValueMemberS{Value: orderId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if output.Item == nil {
		return nil, fmt.Errorf("order with ID %s not found", orderId)
	}

	var orderDbItem orderDbItem
	err = attributevalue.UnmarshalMap(output.Item, &orderDbItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB item: %v", err)
	}

	order := &util.Order{
		OrderId: orderDbItem.OrderId,
		Amount:  orderDbItem.Amount,
		Item:    orderDbItem.Item,
	}

	log.Printf("Retrieved order %s from DynamoDB", orderId)
	return order, nil
}
