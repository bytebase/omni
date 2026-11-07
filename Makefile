ENGINES := cassandra cosmosdb doris elasticsearch googlesql mariadb mongo mssql mysql oceanbase oracle partiql pg redshift snowflake starrocks tidb trino

.PHONY: build test test-spanner proto proto-breaking clean $(addprefix test-,$(ENGINES)) test-mysql-full test-mysql-containers

BUF := go run github.com/bufbuild/buf/cmd/buf@v1.72.0

build:
	go build ./...

# Full suite. Several engines start real database containers. googlesql/spanner-oracle
# has its own go.mod, so ./... does not reach it, and it needs a live emulator
# (SPANNER_EMULATOR_HOST): see test-spanner.
test:
	go test ./...

test-spanner:
	cd googlesql/spanner-oracle && go test ./...

# Per-engine targets: make test-pg, make test-mysql, ...
$(addprefix test-,$(ENGINES)): test-%:
	go test ./$*/...

test-mysql-full:
	./scripts/test-mysql.sh full

test-mysql-containers:
	./scripts/test-mysql.sh container-shards

proto:
	cd proto && $(BUF) format -w && $(BUF) lint && $(BUF) generate

proto-breaking:
	cd proto && $(BUF) breaking --against 'https://github.com/bytebase/omni.git#branch=main,subdir=proto'

clean:
	go clean ./...
