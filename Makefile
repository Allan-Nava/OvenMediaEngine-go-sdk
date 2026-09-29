.PHONY: build test lint vuln integration

build:
	go build ./...

test:
	go vet ./...
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# read-only checks against a real OME: make integration OME_URL=http://localhost:8081 OME_ACCESS_TOKEN=admin:secret
integration:
	go test -count=1 -v ./test/
