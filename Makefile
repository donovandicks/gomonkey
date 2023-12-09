.PHONY: install-tools lint test test-verbose format run

help:
	@echo 'Targets:'
	@echo ' install-tools	Install the development dependencies'
	@echo ' init		Initialize a new go module'
	@echo ' lint		Lint the codebase with golangci-lint'
	@echo ' format		Format the codebase using gofumpt and golines'
	@echo ' test		Run all tests'
	@echo ' test-verbose	Run the tests with the verbose setting'
	@echo ' test-coverage	Run the tests with the coverage setting'

install-tools:
	# Install linting tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest

init: install-tools
	go mod init github.com/$(repo)

lint:
	golangci-lint run -c ./golangci.yml ./...

format:
	gofumpt -l -w -extra .
	golines . -w

test:
	go test -race -count=1 ./... -cover

test-fast:
	go test -count=1 ./...

test-verbose:
	go test -race -v -cover -count=1 ./... -cover

test-coverage:
	go test -race -v -cover -count=1 ./... -coverprofile=cover.out && go tool cover -html=cover.out
