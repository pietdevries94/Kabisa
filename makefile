setenv = export
binname = api
ifeq ($(OS),Windows_NT)
	setenv = set
	binname = api.exe
endif

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
	${setenv} CGO_ENABLED=0
	go build -trimpath -o bin/${binname} cmd/api/main.go cmd/api/handlers.go cmd/api/types.go

.PHONY: build-linux-amd64
build-linux-amd64: prepare-build
	${setenv} CGO_ENABLED=0
	${setenv} GOOS=linux
	${setenv} GOARCH=amd64
	go build -trimpath -o bin/api-linux-amd64 cmd/api/main.go cmd/api/handlers.go cmd/api/types.go

.PHONY: build-windows-amd64
build-windows-amd64: prepare-build
	${setenv} CGO_ENABLED=0
	${setenv} GOOS=windows
	${setenv} GOARCH=amd64
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/api-windows-amd64.exe cmd/api/main.go cmd/api/handlers.go cmd/api/types.go

.PHONY: build-darwin-amd64
build-darwin-amd64: prepare-build
	${setenv} CGO_ENABLED=0
	${setenv} GOOS=darwin
	${setenv} GOARCH=amd64
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o bin/api-darwin-amd64 cmd/api/main.go cmd/api/handlers.go cmd/api/types.go

.PHONY: build-darwin-arm64
build-darwin-arm64: prepare-build
	${setenv} CGO_ENABLED=0
	${setenv} GOOS=darwin
	${setenv} GOARCH=arm64
	go build -trimpath -o bin/api-darwin-arm64 cmd/api/main.go cmd/api/handlers.go cmd/api/types.go

.PHONY: build-all
build-all: build-linux-amd64 build-windows-amd64 build-darwin-amd64 build-darwin-arm64
