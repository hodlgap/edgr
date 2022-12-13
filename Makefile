
all: test vet lint

vet:
	go vet ./...

db:
	docker run --rm -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -p 5432:5432 -d postgres

.PHONY: format
format:
	@go install golang.org/x/tools/cmd/goimports@latest
	goimports -local "github.com/hodlgap/edgr" -w .
	gofmt -s -w .
	go mod tidy