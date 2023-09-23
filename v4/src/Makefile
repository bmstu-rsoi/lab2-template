.PHONY: test
generate:
	go generate -tags=testing ./...
	mockery

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o bin/apiserver ./cmd/apiserver/main.go

.PHONY: test
test:
	go test -tags=testing -v ./...
