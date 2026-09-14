ENGINES := cassandra cosmosdb doris elasticsearch googlesql mariadb mongo mssql mysql oracle partiql pg redshift snowflake starrocks tidb trino

.PHONY: build test test-spanner proto proto-breaking clean $(addprefix test-,$(ENGINES)) test-mysql-full test-mysql-containers

BUF := go run github.com/bufbuild/buf/cmd/buf@v1.72.0

build:
	go build ./...

# Full suite. Several engines start real database containers. harness/ holds
# nested Go modules that ./... does not reach; conformance runs here, the
# spanner harness needs a live emulator (SPANNER_EMULATOR_HOST), see test-spanner.
test:
	go test ./...
	cd harness/conformance && go test ./...

test-spanner:
	cd harness/googlesql-spanner && go test ./...

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
