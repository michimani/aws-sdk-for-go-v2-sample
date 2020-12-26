package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SmapleItem struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	Message1  string `json:"message1"`
	Message2  string `json:"message2"`
	Message3  string `json:"message3"`
}

type SmapleItemMin struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	Message1  string `json:"message1"`
}

func main() {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	client := dynamodb.NewFromConfig(config)

	tableName := "sdk-for-go-sample-table"
	if err := createTable(client, tableName); err != nil {
		log.Fatal(err)
		return
	}

	for n := 1; n <= 5; n++ {
		item := SmapleItem{
			Name:      fmt.Sprintf("Sample Item %d", n),
			CreatedAt: time.Now().UnixNano(),
			Message1:  fmt.Sprintf("This is a sample message %d-1", n),
			Message2:  fmt.Sprintf("This is a sample message %d-2", n),
			Message3:  fmt.Sprintf("This is a sample message %d-3", n),
		}
		if err := putItem(client, tableName, item); err != nil {
			log.Fatal(err)
		}
	}

	itemsAllAttr, err := scan(client, tableName)
	if err != nil {
		log.Fatal(err)
		return
	}
	ja, err := json.Marshal(itemsAllAttr)
	if err != nil {
		log.Fatal(err)
		return
	}
	itemsAll := toJsonString(ja)

	itemsSomeAttr, err := scanWithSomeAttributes(client, tableName)
	if err != nil {
		log.Fatal(err)
		return
	}

	js, err := json.Marshal(itemsSomeAttr)
	if err != nil {
		log.Fatal(err)
		return
	}
	itemsSome := toJsonString(js)

	fmt.Printf("TableName: %s\n", tableName)
	fmt.Printf("Scan with all attributes:\n%s\n\n˜", itemsAll)
	fmt.Printf("Scan with some attributes:\n%s\n\n", itemsSome)
}

// createTable creates talbe if it does not exists.
func createTable(client *dynamodb.Client, name string) error {
	describeIn := dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}
	_, err := client.DescribeTable(context.TODO(), &describeIn)
	if err == nil {
		// table exists
		return nil
	}

	in := dynamodb.CreateTableInput{
		TableName: aws.String(name),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}

	if _, err := client.CreateTable(context.TODO(), &in); err != nil {
		return err
	}

	return nil
}

// putItem puts a item.
func putItem(client *dynamodb.Client, tableName string, item SmapleItem) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	in := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	_, perr := client.PutItem(context.TODO(), &in)
	if perr != nil {
		return err
	}

	fmt.Println(item.Name)

	return nil
}

// scan return list of SmapleItem with all attributes.
func scan(client *dynamodb.Client, tableName string) ([]SmapleItem, error) {
	in := dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	out, err := client.Scan(context.TODO(), &in)
	if err != nil {
		return nil, err
	}

	var sampleItems []SmapleItem
	aerr := attributevalue.UnmarshalListOfMaps(out.Items, &sampleItems)
	if aerr != nil {
		return nil, aerr
	}

	return sampleItems, nil
}

// scan return list of SmapleItem with some attributes.
func scanWithSomeAttributes(client *dynamodb.Client, tableName string) ([]SmapleItemMin, error) {
	proj := expression.NamesList(expression.Name("Name"), expression.Name("CreatedAt"), expression.Name("Message1"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		fmt.Println(err)
	}

	in := dynamodb.ScanInput{
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
		TableName:                aws.String(tableName),
	}

	out, err := client.Scan(context.TODO(), &in)
	if err != nil {
		return nil, err
	}

	var sampleItems []SmapleItemMin
	aerr := attributevalue.UnmarshalListOfMaps(out.Items, &sampleItems)
	if aerr != nil {
		return nil, aerr
	}

	return sampleItems, nil
}

func toJsonString(j []byte) string {
	var buf bytes.Buffer
	jerr := json.Indent(&buf, j, "", " ")

	if jerr != nil {
		fmt.Println(jerr.Error())
		return ""
	}

	return buf.String()
}
