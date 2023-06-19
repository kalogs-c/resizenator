.PHONY: run
run: ## Run the program
	go run main.go

.PHONY: run-race
run-race: ## Run the program with race detector
	go run -race main.go

.PHONY: test
test: ## Run the tests
	go test -v ./...

.PHONY: build
build: ## Build the program
	go build -o resizenator main.go
