.PHONY: test
generate:
	go generate -tags=testing ./...
	mockery

.PHONY: build_apiserver
build_apiserver:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/apiserver ./cmd/apiserver/main.go

.PHONY: build_rating
build_rating:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/rating ./cmd/rating/main.go

.PHONY: build_library
build_library:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/library ./cmd/library/main.go

.PHONY: build_reservation
build_reservation:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/reservation ./cmd/reservation/main.go

.PHONY: test
test:
	go test -tags=testing -v ./...
