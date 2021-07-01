all: build

build:
	go build -o bin/rmqhc cmd/rmqhc/main.go
