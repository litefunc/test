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

// Item : the data structure we need to use
type Item struct {
	Name   string
	City   string
	Rating float64
}

func CreateItem(table string) {

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

	// Create item instance
	item := Item{
		Name:   "Louisa",
		City:   "Kaohsiung",
		Rating: 4.5,
	}

	// Marshall
	av, err := dynamodbattribute.MarshalMap(item)
	fmt.Println(av)
	if err != nil {
		fmt.Println("Got error while marshalling")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create PutItemInput instance
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	// Execute PutItem
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
