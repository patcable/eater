package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Meals struct {
	Entries []Meal `json:"entries"`
	Style   template.CSS
}

type Meal struct {
	MealID           string  `json:"mealID"`
	EntryID          string  `json:"entryID"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	ImageURL         string  `json:"imageURL"`
	EatenAtLocalTime int64   `json:"eatenAtLocalTime"`
	FoodItemName     string  `json:"foodItemName"`
	FoodItemDetails  string  `json:"foodItemDetails"`
	ServingUnits     string  `json:"servingUnits"`
	ServingQuantity  float32 `json:"servingQuantity"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"

	uploadData, err := processMultipartUpload(request, "data")
	if err != nil {
		resp.StatusCode = 500
		resp.Body = err.Error()
		return resp, nil
	}

	var meals Meals
	json.Unmarshal(uploadData, &meals)
	css, _ := ioutil.ReadFile("style.css")
	meals.Style = template.CSS(css)
	templateFuncs := template.FuncMap{
		"convertTime": func(t int64) string {
			s := strconv.FormatInt(t, 10)
			layout := "20060102150405"
			tm, _ := time.Parse(layout, s)
			return tm.Format("Monday Jan 2 2006 @ 3:04pm")
		},
	}
	parsedTemplate, _ := template.New("foodlog.gohtml").Funcs(templateFuncs).ParseFiles("foodlog.gohtml")
	var tpl bytes.Buffer
	parsedTemplate.Execute(&tpl, meals)

	resp.Headers["Content-Type"] = "text/html"
	resp.Body = tpl.String()
	resp.StatusCode = 200
	return resp, nil
}

func processMultipartUpload(request events.APIGatewayProxyRequest, formField string) ([]byte, error) {
	r := http.Request{}
	r.Header = make(map[string][]string)
	for k, v := range request.Headers {
		if k == "content-type" || k == "Content-Type" {
			r.Header.Set(k, v)
		}
	}
	body, err := base64.StdEncoding.DecodeString(request.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("could not read request body: %s", err)
	}

	err = r.ParseMultipartForm(100000)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Multipart Form: %s", err)
	}

	file, _, err := r.FormFile(formField)
	if err != nil {
		return nil, fmt.Errorf("unable to pull file from field %s: %s", formField, err)
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to io.ReadAll: %s", err)
	}
	return buf, nil
}

func main() {
	lambda.Start(HandleRequest)
}
