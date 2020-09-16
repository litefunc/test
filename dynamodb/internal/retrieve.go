package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Retrieve(tableName string) {

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

	// Item feature
	city := "Kaohsiung"
	rating := "4.5"

	// Create GetItemInput instance
	feature := &dynamodb.GetItemInput{

		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{

			"City": {

				S: aws.String(city),
			},
			"Rating": {

				N: aws.String(rating),
			},
		},
	}

	// Get result
	result, err := svc.GetItem(feature)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Hold data
	holder := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &holder)
	if err != nil {

		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	if holder.Name == "" {
		fmt.Println("Could not find data")
		return
	}

	// Display data
	fmt.Println("Found item:", holder)
}
