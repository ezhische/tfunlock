package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	client *dynamodb.Client
}

func NewClient(staticAccessKey, staticSecretKey, region, endpoint string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(staticAccessKey, staticSecretKey, "")))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)

	}

	// Создание клиента DynamoDB c CustomEndpoint
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &Client{client: client}, nil
}

func (c *Client) FullScanTable(tableName string) (map[string]string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Выполнение запроса на получение всех данных
	result, err := c.client.Scan(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}
	out := make(map[string]string, len(result.Items))
	// Выводим все найденные элементы
	for _, item := range result.Items {
		Digest := item["Digest"].(*types.AttributeValueMemberS)
		LockID := item["LockID"].(*types.AttributeValueMemberS)
		out[LockID.Value] = Digest.Value
	}
	return out, nil // Возвращаем все найденные элементы
}

func (c *Client) UpdateDigest(tableName, lockID, newDigest string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"LockID": &types.AttributeValueMemberS{Value: lockID}, // Используем LockID как ключ
		},
		UpdateExpression: aws.String("SET #D = :d"),
		ExpressionAttributeNames: map[string]string{
			"#D": "Digest", // Обновляемое поле
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":d": &types.AttributeValueMemberS{Value: newDigest}, // Новое значение Digest
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}

	// Выполняем запрос на обновление
	_, err := c.client.UpdateItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to update Digest: %w", err)
	}
	return nil
}

func (c *Client) DeleteItemByLockID(tableName, lockID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"LockID": &types.AttributeValueMemberS{Value: lockID}, // Ключевое поле LockID
		},
	}

	// Выполняем запрос на удаление
	_, err := c.client.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete item with LockID %s: %w", lockID, err)
	}
	return nil
}

func (c *Client) GetItemByLockID(tableName, LockID string) (map[string]string, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"LockID": &types.AttributeValueMemberS{Value: LockID}, // Ключевое поле LockID
		},
	}
	result, err := c.client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item %s from %s table: %w", LockID, tableName, err)
	}
	out := make(map[string]string, len(result.Item))
	if result.Item == nil {
		return out, fmt.Errorf("item with LockID %s not found", LockID)
	}
	Digest := result.Item["Digest"].(*types.AttributeValueMemberS)
	id := result.Item["LockID"].(*types.AttributeValueMemberS)
	out[id.Value] = Digest.Value
	return out, nil
}
