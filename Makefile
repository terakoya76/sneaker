lint:
	@if [ -z `which golangci-lint 2> /dev/null` ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.52.2; \
	fi
	@gofmt -s -w .
	@golangci-lint run --timeout 3m

test: lint
	go test -race -v ./...
