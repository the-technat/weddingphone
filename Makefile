GOOS := linux
GOARCH := arm64
CGO_ENABLED := 1
PWD := $(shell pwd)

all: build

vet:
	go vet ./...

lint: staticcheck
	staticcheck ./...

run: 
	mkdir -p ${PWD}/dist/recordings
	SAVE_PATH=${PWD}/dist/recordings go run cmd/main.go

build: go 
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o dist/weddingphone cmd/main.go 
	echo "Binary available at ${PWD}/dist/weddingphone"

staticcheck: go
	go install honnef.co/go/tools/cmd/staticcheck@latest 

go:
	which go 