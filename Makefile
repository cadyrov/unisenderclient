lint:
	gofumpt -w ./ && gofmt -s -w ./ && golangci-lint run