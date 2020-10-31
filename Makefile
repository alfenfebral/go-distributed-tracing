mock: 
	mockery -dir "repository" -all -output mocks \
	mockery -dir "services" -all -output mocks
init-realize:
	realize start --path="." --run main.go
run:
	realize start
test:
	go test ./...
build:
	go build -o go-clean-architecture main.go