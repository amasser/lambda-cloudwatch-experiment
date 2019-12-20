# Lambda to CloudWatch experiment

This experiment reproduces an issue of writing logs to AWS CloudWatch from AWS Lambda written in Go.

In small number of cases the last log written from the lambda will not be saved to CloudWatch.
In all observed occurrence of the issue the request which exhibits the faulty behavior is the last one written in CloudWatch stream. 

The simple [lambda](lambda.go) writes a log in JSON format which contains AWS request ID and message `begin`.
Then it sleeps for 10 to 600 miliseconds.
After which it again writes to log in JSON format which contains AWS request ID and message `end`.

Sometimes the second log (one with message `end`) will not be written to CloudWatch.

## Gathered data

We offer the following data as proof:

- 2019-12-09:
	- [Analysis records](findings/2019-12-09/analyze-logs-2019-12-09.txt) which lists number of request logs analyzed and list of the ones which have a missing pair.
	- Complete CloudWatch streams for bad occurrences:
	    - [Request c76d6233-7b64-489a-bc28-7f26a09446e6 stream c7116d03d8cd480e9139328dfe28a7ee](findings/2019-12-09/req_c76d6233-7b64-489a-bc28-7f26a09446e6_stream_c7116d03d8cd480e9139328dfe28a7ee.txt)
	    - [Request 83445a2c-f381-4f38-9091-c16cbf1eda33 stream a68cf9e178884125b5b38b5ee81b0f14](findings/2019-12-09/req_83445a2c-f381-4f38-9091-c16cbf1eda33_stream_a68cf9e178884125b5b38b5ee81b0f14.txt)
	    - [Request 728cae3b-6b12-41f9-a6eb-ed99b92d50b1 stream 01d9127a46934a5f806f4403fda64ffc](findings/2019-12-09/req_728cae3b-6b12-41f9-a6eb-ed99b92d50b1_stream_01d9127a46934a5f806f4403fda64ffc.txt)
	    - [Request e075102a-f773-484f-9ed3-f05026080b59 stream 9f192e37e3a14aed966ed4a5e1cc7288](findings/2019-12-09/req_e075102a-f773-484f-9ed3-f05026080b59_stream_9f192e37e3a14aed966ed4a5e1cc7288.txt)
	    - [Request 06058d99-3be3-4a1e-afda-e2c466050a75 stream 0385626bb44c4e3a8558a5ad5f4dff64](findings/2019-12-09/req_06058d99-3be3-4a1e-afda-e2c466050a75_stream_0385626bb44c4e3a8558a5ad5f4dff64.txt)
	- [Complete CloudWatch log group](findings/2019-12-09/complete-logs-2019-12-09.json)
- 2019-12-20:
	- [Analysis records](findings/2019-12-20/analyze_results.txt) which lists number of request logs analyzed and list of the ones which have a missing pair.
    - Complete CloudWatch streams for bad occurrences:
    	- [Request b2bc175f-71fd-4b2c-b473-c046ff8047aa steam 76d4930b335b4810b717fef20aab7de2](findings/2019-12-20/stream_76d4930b335b4810b717fef20aab7de2.txt)
    	- [Request 6648cefc-c975-46de-8215-57dae9bc6d46 steam d88f6e8b3dcf487a92ebdc3fa54c87f3](findings/2019-12-20/stream_d88f6e8b3dcf487a92ebdc3fa54c87f3.txt)
    - [Complete CloudWatch log group](findings/2019-12-20/complete-logs-2019-12-20.json)

## How to reproduce the experiment

Please install [Go](https://golang.org/dl/) and [Terraform](https://www.terraform.io/downloads.html) before proceeding.

Run `make install` to compile [lambda.go](lambda.go) binary and provision the infrastructure. 

Run `make run` to run experiment. The experiment will invoke the [lambda](lambda.go) with invocation type Event 600 times in 600 concurrent executions waiting 5 to 6 seconds between doing request (each of indviidual 600 concurrent executions).

Wait an ~hour for any delayed CloudWatch logs to be written (known problem).

Run `make analyze` to analyze logs. The output will contain number of request logs analyzed, how many are missing and AWS Request IDs of missing ones.

Use CloudWatch console or AWS CLI to analyze the issue.

For clean up run `make uninstall` to tear-down provision infrastructure, note this will also delete the CloudWatch group with the data.
