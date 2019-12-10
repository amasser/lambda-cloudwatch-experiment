package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

const functionName = "lambda-cloudwatch-experiment-dev"

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) < 5 {
		log.Fatal("no args provided: total, concurrent, minSleep, maxSleep")
	}

	total := mustAtoi(os.Args[1])
	concurrent := mustAtoi(os.Args[2])
	minSleep := mustAtoi(os.Args[3])
	maxSleep := mustAtoi(os.Args[4])

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	svc := lambda.New(cfg)

	ch := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			for {
				_, ok := <-ch
				if !ok {
					return
				}
				doReq(svc)

				randomValue := rand.Intn(maxSleep-minSleep+1) + minSleep
				time.Sleep(time.Duration(randomValue) * time.Second)
			}
		}()
	}

	for i := 0; i < total; i++ {
		ch <- struct{}{}
	}
	close(ch)

	wg.Wait()
}

func doReq(svc *lambda.Client) {
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: lambda.InvocationTypeEvent,
	}
	req := svc.InvokeRequest(input)

	_, err := req.Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
