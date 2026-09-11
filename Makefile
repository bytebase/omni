ENGINES := cassandra cosmosdb doris elasticsearch googlesql mariadb mongo mssql mysql oracle partiql pg redshift snowflake starrocks tidb trino

.PHONY: build test test-short proto clean $(addprefix test-,$(ENGINES)) test-mysql-quick test-mysql-full test-mysql-containers

BUF := go run github.com/bufbuild/buf/cmd/buf@v1.72.0

build:
	go build ./...

# Full suite. Several engines start real database containers unless -short is set.
test:
	go test ./...

# What CI runs on every PR: hermetic, no containers.
test-short:
	go test -short ./...

# Per-engine targets: make test-pg, make test-mysql, ...
$(addprefix test-,$(ENGINES)): test-%:
	go test ./$*/...

test-mysql-quick:
	./scripts/test-mysql.sh quick

test-mysql-full:
	./scripts/test-mysql.sh full

test-mysql-containers:
	./scripts/test-mysql.sh container-shards

proto:
	cd proto && $(BUF) format -w && $(BUF) lint && $(BUF) generate

clean:
	go clean ./...
