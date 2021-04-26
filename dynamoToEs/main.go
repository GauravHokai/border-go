package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

type data struct {
	Month        string
	Cupcake      string
	Updated_time string
	Operation    string
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {

	for _, record := range e.Records {
		if record.EventName == "MODIFY" {
			r := record.Change.NewImage
			id := r["Month"].String()
			temp := data{
				Month:        r["Month"].String(),
				Cupcake:      r["Cupcake"].String(),
				Updated_time: r["Updated_time"].String(),
				Operation:    record.EventName,
			}

			body, err := json.Marshal(temp)

			if err != nil {
				fmt.Println(err)
				continue
			}

			url := "https://search-border-assign-es-lxz37crwep7yx3rpmnke3i4kea.us-east-1.es.amazonaws.com/datatable/doc/" + id

			client := &http.Client{}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(err)
				continue
			}
			req.SetBasicAuth("", "")
			req.Header.Add("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			bodyText, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(bodyText))
		}
	}
}
