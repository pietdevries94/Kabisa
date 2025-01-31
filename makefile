.PHONY: default
default: build

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY:	prepare-build
prepare-build: generate lint test

# build builds the application for the current platform in as bin/api
.PHONY: build
build: prepare-build
	CGO_ENABLED=0 go build -trimpath -o bin/api cmd/api/*.go

.PHONY: build-linux-amd64
build-linux-amd64: prepare-build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/api-linux-amd64 cmd/api/*.go

.PHONY: build-windows-amd64
build-windows-amd64: prepare-build
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/api-windows-amd64.exe cmd/api/*.go

.PHONY: build-darwin-amd64
build-darwin-amd64: prepare-build
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/api-darwin-amd64 cmd/api/*.go

.PHONY: build-darwin-arm64
build-darwin-arm64: prepare-build
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -o bin/api-darwin-arm64 cmd/api/*.go

.PHONY: build-all
build-all: build-linux-amd64 build-windows-amd64 build-darwin-amd64 build-darwin-arm64
