
GOOS := linux
GOARCH := arm64
CGO_ENABLED := 1

all:
.PHONY: build

vet:
	go vet ./...

lint: staticcheck
	staticcheck ./...

build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o dist/weddingphone main.go 

staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest 