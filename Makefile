.PHONY: build test test-pg test-mysql test-mysql-quick test-mysql-full test-mysql-containers test-mssql test-oracle proto proto-breaking clean

BUF := go run github.com/bufbuild/buf/cmd/buf@v1.72.0

build:
	go build ./...

test:
	go test ./...

test-pg:
	go test ./pg/...

test-mysql:
	go test ./mysql/...

test-mysql-quick:
	./scripts/test-mysql.sh quick

test-mysql-full:
	./scripts/test-mysql.sh full

test-mysql-containers:
	./scripts/test-mysql.sh container-shards

test-mssql:
	go test ./mssql/...

test-oracle:
	go test ./oracle/...

proto:
	cd proto && $(BUF) format -w && $(BUF) lint && $(BUF) generate

proto-breaking:
	cd proto && $(BUF) breaking --against 'https://github.com/bytebase/omni.git#branch=main,subdir=proto'

clean:
	go clean ./...
