# Engine test database matrix

The required baseline for live tests is the engine's own database or its official local emulator. Redshift has the temporary PostgreSQL exception recorded below; OceanBase retains its existing mixed baseline while its migration is paused. Container startup, connection setup, image selection, and cleanup belong in that engine's `internal/testcontainer` package. Tests call the engine-owned helper. Each test package shares a server and isolates fixtures with databases or schemas, except when a test requires a different server configuration.

The matrix covers all 18 engine directories. It records the selected images and the current coverage gaps. The gaps remain open until the native oracle suites pass.

| Engine | Selected database image | Coverage and current gaps |
| --- | --- | --- |
| `cassandra` | `cassandra:5.0.9` | CQL acceptance and rejection through the server's `cqlsh`. |
| `cosmosdb` | `mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:vnext-EN20261008` | Official NoSQL emulator, with SELECT acceptance and rejection probes. Emulator limitations apply. |
| `doris` | `dyrnq/doris:4.0.8` | Native parser oracle, including the repaired 4.0 string-function, query-tail, and load/job grammar. |
| `elasticsearch` | `docker.elastic.co/elasticsearch/elasticsearch:9.5.6` | Parsed REST requests execute against Elasticsearch, including indexing and searches. |
| `googlesql` | `gcr.io/cloud-spanner-emulator/emulator:1.5.54` | Official Spanner emulator. Live coverage targets the GoogleSQL features Spanner implements. |
| `mariadb` | `mariadb:11.8.8` | Native parser, catalog, review, and routine-validation oracles pass, including temporal defaults, column metadata, constraints, index types, and stored DDL. Configured external full-text parser plugins also have metadata unit fixtures. |
| `mongo` | `mongo:9.0.4` | Shell command acceptance and rejection through the server image's `mongosh`. |
| `mssql` | `mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04` | Native parser and review oracles. ScriptDOM supplies an additional grammar reference. |
| `mysql` | `mysql:8.4.13` | Sole live MySQL version. Legacy normalization models use unit fixtures. Upgrade drift is checked against the 8.4 server. |
| `oceanbase` | Native OceanBase `4.3.5.6-106000012026040916` plus inherited `mysql:5.7` / `mysql:8.0` oracles | Baseline unchanged. Native-only migration and compatibility repairs are paused and preserved for a separate follow-up. |
| `oracle` | `gvenzl/oracle-free:23.26.3-slim-faststart` | Native parser and review oracles. Reference tests support explicit Oracle DSNs. |
| `partiql` | `amazon/dynamodb-local:3.3.1` | Official DynamoDB Local tests its supported PartiQL subset, including SELECT, INSERT, UPDATE, DELETE, and invalid statements. |
| `pg` | `postgres:18.6-alpine` | Sole live PostgreSQL version for parser, catalog, analysis, review, and splitter oracles. |
| `redshift` | `postgres:18.6-alpine` | Explicit temporary PostgreSQL baseline, matching `pg`. PostgreSQL 18's not-null catalog rows are excluded from constraint comparisons; column nullability is compared separately, preserving Redshift's production model. Real Redshift checks use `REDSHIFT_COMPAT_DSN` through the compatibility report command. |
| `snowflake` | Pending | A native test account or suitable licensed emulator is required for live SQL coverage. |
| `starrocks` | `starrocks/allin1-ubuntu:4.1.7` | Native parser acceptance and rejection oracle. |
| `tidb` | `pingcap/tidb:v8.5.5` | Full parser, catalog, completion, deparser, and review suites pass. Native connections verify `tidb_version()`. Constraint scenarios explicitly enable `tidb_enable_check_constraint` on both sides; baseline probes also cover its default OFF behavior. Expression-index probes use the default function restrictions; metadata unit fixtures additionally cover an explicitly configured `allow-expression-index`. |
| `trino` | `trinodb/trino:482` | Native parser, deparser, analysis, and completion oracles. |

## Helper layout

Every engine with a local database or emulator uses the same filenames:

```text
<engine>/internal/testcontainer/
  container.go
  container_test.go  # when the helper has tests
```

`container.go` owns the image pin, startup, connections, database isolation, protocol helpers, and cleanup. Helper tests belong in `container_test.go`. The nested Spanner harness follows the same layout. Keep implementation inside its engine; do not add empty files for engines without helper tests. Snowflake has no container helper while its native backend is unavailable. Import the engine-owned package without an alias and call it as `testcontainer.*`.

## Image policy

Pin a complete release version or a dated vendor build. Use an immutable digest when the vendor provides only a moving tag. MySQL uses only the selected 8.4 release, and PostgreSQL and the temporary Redshift baseline use the same selected PostgreSQL 18 patch.

An image upgrade includes running its engine's oracle suite, investigating acceptance and metadata drift, and fixing the parser, catalog, or test fixtures according to the new engine's behavior. Keep this matrix and the image constants in the same change.

TiDB's test servers must be TiDB. Its review image override rejects other image repositories, and parser, catalog, and review connections verify the TiDB-specific version function.

Doris reports `5.7.99` through `SELECT VERSION()` for MySQL protocol compatibility. Read `SELECT @@version_comment` for the Doris build. Select the Doris image by its Doris release version. [Doris version function](https://doris.apache.org/docs/4.x/sql-manual/sql-functions/scalar-functions/system-functions/version-function/)

## Emulator scope

DynamoDB Local provides a reference for the DynamoDB PartiQL subset. The general PartiQL parser also supports constructs outside that subset. GoogleSQL's Spanner emulator and Cosmos DB's emulator likewise cover their respective service implementations. [DynamoDB PartiQL reference](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ql-reference.html), [Cosmos DB emulator](https://learn.microsoft.com/en-us/azure/cosmos-db/emulator-linux)

LocalStack's Redshift implementation uses PostgreSQL. Real Redshift is required to establish Redshift-specific SQL fidelity. Snowflake's LocalStack emulator requires a license, and Snowflake's official Snowpark local framework has limited SQL support. [LocalStack Redshift implementation](https://blog.localstack.cloud/localstack-for-aws-release-2026-08-0/), [LocalStack Snowflake licensing](https://docs.localstack.cloud/snowflake/getting-started/auth-token/), [Snowpark local testing](https://docs.snowflake.com/en/developer-guide/snowpark/python/testing-locally)

## Additional harnesses

`harness/mssql-scriptdom` uses `mcr.microsoft.com/dotnet/sdk:8.0.425` to run Microsoft's ScriptDOM parser. This supplements the native SQL Server oracle.

`harness/googlesql-spanner` uses the same Spanner emulator release as `googlesql`. `harness/conformance/start_tidb.sh` uses the selected TiDB release. These are nested Go modules and require their own test commands.
