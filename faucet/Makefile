all: tidy format lint test

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v -shuffle=on -count=1 -race -timeout 20m ./...

.PHONY: format
format:
	gofmt -w -s .
	goimports -w -local "github.com/filecoin-project/" .

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: vulncheck
vulncheck:
	govulncheck -v ./...
