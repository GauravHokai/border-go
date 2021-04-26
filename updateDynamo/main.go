package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getRandomMonth() string {
	rand.Seed(time.Now().UnixNano())
	min := 2004
	max := 2021

	year := rand.Intn(max-min+1) + 2004
	month := rand.Intn(12) + 1

	str := fmt.Sprintf("%d-%02d", year, month)

	return str
}

func handler() {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})

	if err != nil {
		fmt.Println(err)
		//return err
	}

	svc := dynamodb.New(sess)

	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		tableName := "Data2"
		monthRand := getRandomMonth()
		cupcake := strconv.Itoa(rand.Intn(101))
		updateTime := strconv.FormatInt(time.Now().Unix(), 10)

		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":Updated_time": {
					S: aws.String(updateTime),
				},
				":Cupcake": {
					N: aws.String(cupcake),
				},
				":month_": {
					S: aws.String(monthRand),
				},
			},
			ExpressionAttributeNames: map[string]*string{
				"#m": aws.String("Month"),
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Month": {
					S: aws.String(monthRand),
				},
			},
			ReturnValues:        aws.String("UPDATED_NEW"),
			UpdateExpression:    aws.String("set Updated_time = :Updated_time, Cupcake = :Cupcake"),
			ConditionExpression: aws.String("#m = :month_"),
		}

		_, err := svc.UpdateItem(input)
		if err != nil {
			fmt.Println(err)
			//return err
		}

		fmt.Println(monthRand)
		fmt.Println(cupcake)
		fmt.Println(updateTime)

	}

}

func main() {
	lambda.Start(handler)
	//handler()
}
