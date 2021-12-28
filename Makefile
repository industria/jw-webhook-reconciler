.DEFAULT_GOAL := build

vet:
	go vet ./...
.PHONY:vet

build: vet
	go build -o reconsile
.PHONY:build

clean:
	go clean
.PHONY:clean
