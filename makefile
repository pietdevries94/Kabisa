.PHONY: default
default: build

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

# build builds the application for the current platform
.PHONY: build
build: lint test
	CGO_ENABLED=0 go build -o bin/api cmd/api/*.go

.PHONY: build-linux-amd64
build-linux-amd64: lint test
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/api-linux-amd64 cmd/api/*.go

.PHONY: build-windows-amd64
build-windows-amd64: lint test
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/api-windows-amd64.exe cmd/api/*.go

.PHONY: build-darwin-amd64
build-darwin-amd64: lint test
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/api-darwin-amd64 cmd/api/*.go

.PHONY: build-darwin-arm64
build-darwin-arm64: lint test
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/api-darwin-arm64 cmd/api/*.go

.PHONY: build-all
build-all: build-linux-amd64 build-windows-amd64 build-darwin-amd64 build-darwin-arm64
