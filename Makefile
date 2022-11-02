GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

mock: 
	mockery --dir repository --all --output mocks/repository
	mockery --dir services --all --output mocks/services
run:
	air
test:
	go test ./...
mock-test:
	make mock
	make test
build:
	go build -o go-clean-architecture cmds/app/main.go
.PHONY: test/cover
test/cover:
	mkdir -p coverage
	$(GOTEST) -v -coverprofile=coverage/coverage.out ./...
	$(GOCOVER) -func=coverage/coverage.out
	$(GOCOVER) -html=coverage/coverage.out -o coverage/coverage.html