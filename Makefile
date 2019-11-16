test: 
	go test -v -race ./...

update: 
	go get -u ./...

tidy: 
	go mod tidy
