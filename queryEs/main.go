package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type data struct {
	Month   string
	Cupcake string
}

type bigdata2 struct {
	Datarows [][]interface{} `json:"datarows"`
	//Datarows []data
}

// type responseData struct {
// 	Schema   []string
// 	total    int
// 	datarows [][]data
// 	size     int
// 	status   int
// }

// type Response struct {
// 	Body string `json: body`
// }

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

	fmt.Println(string(temp))

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

	//fmt.Println(string(bodyText))

	temp1 := &bigdata2{}

	x := json.Unmarshal(bodyText, temp1)

	if x != nil {
		fmt.Println(err)
	}

	//fmt.Println(temp1.Datarows)

	var finalData []data
	for _, dd := range temp1.Datarows {
		result := data{dd[0].(string), dd[1].(string)}

		finalData = append(finalData, result)
	}

	ress, err := json.Marshal(finalData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(ress))

	m := make(map[string]string)

	m["content-type"] = "application/json"
	m["Access-Control-Allow-Origin"] = "*"

	return events.APIGatewayProxyResponse{Body: string(ress), Headers: m, StatusCode: 200}, nil

}
