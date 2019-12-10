package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context) error {
	var reqID string
	if lctx, ok := lambdacontext.FromContext(ctx); ok {
		reqID = lctx.AwsRequestID
	}

	_, err := fmt.Printf("%s\n", mustJsonEncode(logData{ReqID: reqID, Message: "begin"}))
	if err != nil {
		panic(err)
	}

	min, max := 10, 600
	randomValue := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(randomValue) * time.Millisecond)

	_, err = fmt.Printf("%s\n", mustJsonEncode(logData{ReqID: reqID, Message: "end"}))
	if err != nil {
		panic(err)
	}

	return nil
}

type logData struct {
	ReqID   string `json:"reqID"`
	Message string `json:"message"`
}

func mustJsonEncode(data logData) []byte {
	v, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return v
}
