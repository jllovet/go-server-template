build:
	@go build -o bin/server cmd/main.go

run:
	@./bin/server

test:
	@go test ./...