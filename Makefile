.PHONY: build test test-pg test-mysql test-mysql-quick test-mysql-full test-mysql-containers test-mssql test-oracle clean

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

clean:
	go clean ./...
