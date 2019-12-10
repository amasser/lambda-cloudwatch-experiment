package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

const logGroup = "/aws/lambda/lambda-cloudwatch-experiment-dev"

type logData struct {
	ReqID   string `json:"reqID"`
	Message string `json:"message"`
}

func main() {
	awsCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	cwsvc := cloudwatchlogs.New(awsCfg)

	reqInput := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(logGroup),
	}

	err = reqInput.Validate()
	if err != nil {
		log.Fatal(err)
	}

	events, err := fetchLogs(cwsvc, reqInput)
	if err != nil {
		log.Fatal(err)
	}

	requests := make(map[string]int)

logLoop:
	for _, reqEv := range events {
		if !strings.HasPrefix(*reqEv.Message, "{") {
			continue logLoop
		}

		var data logData
		err = json.Unmarshal([]byte(*reqEv.Message), &data)
		if err != nil {
			log.Fatal(err)
		}

		requests[data.ReqID]++
	}

	total := 0
	var missingRequests []string
	for reqID, cnt := range requests {
		total++
		if cnt < 2 {
			missingRequests = append(missingRequests, reqID)
		}
	}

	fmt.Printf("total %d, missing %d\n", total, len(missingRequests))
	for _, reqID := range missingRequests {
		fmt.Println(reqID)
	}
}

func fetchLogs(svc *cloudwatchlogs.Client, input *cloudwatchlogs.FilterLogEventsInput) ([]cloudwatchlogs.FilteredLogEvent, error) {
	logs := make([]cloudwatchlogs.FilteredLogEvent, 0)

	hasMore := true
	for hasMore {
		req := svc.FilterLogEventsRequest(input)

		res, err := req.Send(context.Background())
		if err != nil {
			return nil, err
		}

		logs = append(logs, res.Events...)

		if res.NextToken != nil {
			input.NextToken = res.NextToken
			continue
		}

		hasMore = false
	}

	sort.Slice(logs, func(i, j int) bool {
		return *logs[i].Timestamp < *logs[j].Timestamp
	})

	return logs, nil
}
