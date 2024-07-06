.PHONY: start format run

# Format all Go files
format:
	find . -name '*.go' -exec gofmt -w {} +

# Run the Go application
run:
	go run cmd/app/main.go

# Combine both commands
start: format run
