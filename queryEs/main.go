package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type data struct {
	Month   string
	Cupcake int
}

type insideSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type responseData struct {
	Schema   []insideSchema `json:"schema"`
	Total    int            `json:"total"`
	Datarows [][]string     `json:"datarows"`
	Size     int            `json:"size"`
	Status   int            `json:"status"`
}

func main() {
	lambda.Start(handler)
	//handler()
}
func handler() (events.APIGatewayProxyResponse, error) {
	q := map[string]string{
		"query": "Select Month, Cupcake from bor;",
	}
	url := "https://search-border-assign-es-lxz37crwep7yx3rpmnke3i4kea.us-east-1.es.amazonaws.com/_opendistro/_sql"
	fmt.Println(q)

	temp, err := json.Marshal(q)

	//fmt.Println("temp: ", string(temp))

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(temp))
	if err != nil {
		fmt.Println(err)
		//continue
	}
	req.SetBasicAuth("", "")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("bodyText: ", string(bodyText))

	var respo responseData

	json.Unmarshal(bodyText, &respo)

	//fmt.Println(respo.Datarows)

	var finalData []data

	for _, val := range respo.Datarows {
		t, _ := strconv.Atoi(val[1])
		result := data{
			Month:   val[0],
			Cupcake: t,
		}

		finalData = append(finalData, result)
	}

	//fmt.Println(finalData)

	finalResponse, err := json.Marshal(finalData)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(finalResponse))
	m := make(map[string]string)

	m["content-type"] = "application/json"
	m["Access-Control-Allow-Origin"] = "*"

	return events.APIGatewayProxyResponse{Body: string(finalResponse), Headers: m, StatusCode: 200}, nil

}
