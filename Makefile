LDFLAGS = -ldflags "-s -X main.gitSHA=$(shell git rev-parse HEAD) -linkmode=external"
LINTERS= -E golint -E gosec -E interfacer -E unconvert -E dupl -E goconst -E gocyclo $\
-E gofmt -E goimports -E maligned -E depguard -E misspell -E lll -E unparam $\
-E nakedret -E prealloc -E scopelint -E gocritic -E gochecknoinits -E gochecknoglobals

.DEFAULT_GOAL := all

.PHONY: all
all: test build

.PHONY: build
build:
	go build -a $(LDFLAGS) cmd/idiogo/main.go

.PHONY: test
test:
	go test -count=1 -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

## Set up dev database with schema etc.
mssql-init:
	./scripts/configure_local_db.sh