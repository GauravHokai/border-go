package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type data struct {
	Month        string
	Cupcake      string
	Updated_time string
}

func main() {
	recordFile, err := os.Open("multiTimeline.csv")

	if err != nil {
		fmt.Println("Error: ", err)
	}

	reader := csv.NewReader(recordFile)

	records, _ := reader.ReadAll()

	var allData []data

	for _, element := range records {
		var temp data
		temp.Month = element[0]
		temp.Cupcake = element[1]
		temp.Updated_time = strconv.FormatInt(time.Now().Unix(), 10)

		allData = append(allData, temp)
	}

	//fmt.Println(allData)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(""),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})

	if err != nil {
		fmt.Println(err)
	}

	svc := dynamodb.New(sess)

	for _, item := range allData {
		//fmt.Println(item)
		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			fmt.Println("Marshalling Error: ", err)
		}

		//fmt.Println(av)

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Data2"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println("Got error calling PutItem: ", err)
		}

	}
}
