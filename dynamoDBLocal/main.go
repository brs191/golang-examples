//https://qiita.com/yumin/items/008588a4d6d828002ebc

package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tableName string
var sess *session.Session
var dbSvc *dynamodb.DynamoDB

type appContext struct {
	sess  *session.Session
	dbSvc *dynamodb.DynamoDB
}

// Item is a table item for Raja
type Item struct {
	Msisdn string
	Text   string
}

// User is list of msisdns
type User struct {
	ID string
}

func init() {
	log.SetPrefix("dynamoLocalxxx =>")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	// valid.SetFieldsRequiredByDefault(true)
}

func add(item Item) {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("MarshalMap err ", err)
		return
	}

	iteminput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dbSvc.PutItem(iteminput)
	if err != nil {
		log.Println("PutItemInput err ", err)
		return
	}
}

func get(msisdn string) {
	getresult, err := dbSvc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Msisdn": {
				S: aws.String(msisdn),
			},
		},
	})
	if err != nil {
		log.Println("GetItem err ", err)
	}

	if getresult.Item == nil {
		log.Println("Could not find " + msisdn)
		return
	}

	item := Item{}
	err = dynamodbattribute.UnmarshalMap(getresult.Item, &item)
	if err != nil {
		log.Println("unMarshalGet err ", err)
		return
	}

	log.Println("Found Item")
	log.Println("Year: ", item.Msisdn)
	log.Println("Year: ", item.Text)
}

func delete(msisdn string) {
	delInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Msisdn": {
				S: aws.String(msisdn),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.DeleteItem(delInput)
	if err != nil {
		log.Println("DeleteItem err ", err)
		return
	}

	log.Println("Delete ", msisdn+" from table "+tableName)
}

func main() {

	tableName = "raja-ex"

	creds := credentials.NewStaticCredentials("123", "123", "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:8000")}
	sess, err := session.NewSession(awsConfig)

	if err != nil {
		log.Println(err)
		return
	}
	dbSvc = dynamodb.New(sess)

	log.Println("dynamoDB Service create")
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Msisdn"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Msisdn"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err = dbSvc.CreateTable(input)
	if err != nil {
		log.Println("createtable err ", err)
		// return
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := dbSvc.ListTablesWithContext(timeoutCtx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		log.Println(*table)
	}

	item := Item{
		Msisdn: "123456789",
		Text:   "0123456789",
	}
	add(item)

	item = Item{
		Msisdn: "223456789",
		Text:   "0223456789",
	}
	add(item)

	item = Item{
		Msisdn: "323456789",
		Text:   "0323456789",
	}
	add(item)

	log.Println("added " + item.Msisdn + " in " + item.Text + " to table " + tableName)

	get("323456789")

	delete("323456789")
}
