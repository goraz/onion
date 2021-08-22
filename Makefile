export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
test: 
	go test -v -race ./...

update: 
	go get -u -t ./...

tidy: 
	go mod tidy

bin/golangci-lint:
	mkdir -p $(ROOT)/bin
	echo "*"> $(ROOT)/bin/.gitignore
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(ROOT)/bin v1.27.0

lint: bin/golangci-lint
	bin/golangci-lint run

