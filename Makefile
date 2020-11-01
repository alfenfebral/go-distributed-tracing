GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

mock: 
	mockery --dir repository --all --output mocks/repository
	mockery --dir services --all --output mocks/services
init-realize:
	realize start --path="." --run main.go
run:
	realize start
test:
	go test ./...
build:
	go build -o go-clean-architecture main.go
.PHONY: test/cover
test/cover:
	mkdir -p coverage
	$(GOTEST) -v -coverprofile=coverage/coverage.out ./...
	$(GOCOVER) -func=coverage/coverage.out
	$(GOCOVER) -html=coverage/coverage.out -o coverage/coverage.html