.PHONY: build test test-verbose test-coverage clean install run

# Build the application
build:
	go build -o tada

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -cover

# Run tests with detailed coverage report
test-coverage-detail:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -f tada
	rm -f coverage.out coverage.html

# Install the application to GOPATH/bin
install:
	go install

# Run the application
run: build
	./tada

# Run the application with a custom todo file
run-dev: build
	./tada -f ./test-todo.txt
