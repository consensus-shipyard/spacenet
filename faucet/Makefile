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

.PHONY: cert
cert:
	mkcert -cert-file "./_cert/cert.pem" -key-file "./_cert/key.pem" faucet.com 127.0.0.1 localhost

.PHONY: demo
demo:
	rm -rf ./_db_data
	go run ./cmd/main.go --web-allowed-origins "https://localhost" --tls-cert-file "./_cert/cert.pem" --tls-key-file "./_cert/key.pem"