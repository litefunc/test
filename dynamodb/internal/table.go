package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateTable(tableName string) {

	os.Setenv("AWS_ACCESS_KEY_ID", "DUMMYIDEXAMPLE")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "DUMMYEXAMPLEKEY")

	// Create client
	cfg := aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}
	sess, err := session.NewSession(&cfg)
	if err != nil {
		log.Println(err)
		return
	}
	svc := dynamodb.New(sess)

	// Create CreateTableInput instance
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("City"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Rating"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("City"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Rating"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	// Execute CreateTable
	_, err = svc.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Created the table", tableName)
}
