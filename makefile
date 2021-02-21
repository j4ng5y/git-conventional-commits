all: test build

.PHONY: test
test:
	@go test --coverprofile=coverage.out ./...

.PHONY: build
build:
	@go build -a -o bin/git-cc cmd/cc/git-cc.go

.PHONY: coverage
coverage:
	@go tool cover -func=coverage.out

.PHONY: install
install: bin/git-cc
	@cp bin/git-cc /usr/local/bin/git-cc