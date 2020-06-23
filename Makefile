default: build

GIT_SHA:= $(shell git rev-parse --short HEAD || echo "GitNotFound")
BINARY:= "xdscli"

GO_LDFLAGS:= "-X main._gitSHA=$(GIT_SHA)"
build:
	go build -ldflags $(GO_LDFLAGS) -o $(BINARY)

lint:
	@go fmt $(shell go list ./...)
