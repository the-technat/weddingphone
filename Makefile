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
	go run main.go

build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o dist/weddingphone main.go 
	echo "Binary available at ${PWD}/dist/weddingphone"

staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest 
