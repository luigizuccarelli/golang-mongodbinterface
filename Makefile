.PHONY: all test build

all: clean build

build: 
	mkdir -p build
	go build -o build ./... 

test:
	go test -v -tags=test  -coverprofile=tests/results/cover.out ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/microservice 
	go clean ./...

image:
	podman build -t docker.io/lzuccarelli/golang-mongodbinterface:1.13.1 .

push:
	podman push docker.io/lzuccarelli/golang-mongodbinterface:1.13.1
