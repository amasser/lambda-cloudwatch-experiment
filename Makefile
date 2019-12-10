.PHONY: run analyze install uninstall

run:
	go run run-test/run-test.go 600 600 5 6

analyze:
	go run analyze-logs/analyze-logs.go

install:
	GOOS=linux GOARCH=amd64 go build -o lambda
	zip lambda.zip lambda
	cd infrastructure && \
	terraform init -upgrade && \
	terraform apply

uninstall:
	cd infrastructure && \
	terraform destroy
