# Omni

SQL toolchain for multiple database engines. Each engine provides a parser, AST, and additional components such as catalog simulation and semantic analysis -- all in pure Go with zero dependencies.

## Features

- **Zero dependencies** -- pure Go, no CGo, no generated code at runtime
- **Full AST** -- every parsed statement produces a complete abstract syntax tree
- **Position tracking** -- every AST node carries byte-offset location info
- **Beyond parsing** -- catalog simulation, DDL semantic analysis, and more per engine

## Engines

Every engine lives in its own top-level directory and follows the same shape: `ast/` holds the node types, `parser/` the hand-written recursive descent parser, and further packages add capabilities on top. The table lists the packages each engine ships today.

| Engine | Directory | Packages beyond `ast` + `parser` |
|--------|-----------|----------------------------------|
| PostgreSQL | `pg/` | `catalog`, `completion`, `plpgsql`, `pgregress`, `splittest` |
| Amazon Redshift | `redshift/` | `catalog`, `completion`, `analysis`, `plpgsql`, `compat`, `pgregress` |
| MySQL | `mysql/` | `catalog`, `completion`, `deparse`, `validate`, `scope`, `quality` |
| MariaDB | `mariadb/` | `catalog`, `completion`, `deparse`, `validate`, `scope`, `quality` |
| TiDB | `tidb/` | `catalog`, `completion`, `deparse`, `scope`, `quality` |
| SQL Server (T-SQL) | `mssql/` | `completion` |
| Oracle | `oracle/` | `quality` |
| Snowflake | `snowflake/` | `analysis`, `deparse`, `diagnostics`, `advisor` |
| Trino | `trino/` | `catalog`, `completion`, `analysis`, `deparse` |
| GoogleSQL (BigQuery, Spanner) | `googlesql/` | `analysis`, `diagnostics` |
| Apache Doris | `doris/` | `analysis` |
| StarRocks | `starrocks/` | `analysis` |
| PartiQL (DynamoDB) | `partiql/` | `catalog`, `completion`, `analysis` |
| Apache Cassandra (CQL) | `cassandra/` | -- |
| MongoDB (mongosh) | `mongo/` | `catalog`, `completion`, `analysis` |
| Azure Cosmos DB | `cosmosdb/` | `analysis` |
| Elasticsearch (Dev Console) | `elasticsearch/` | `analysis` |

Package roles:

- `catalog` -- in-memory schema simulation: apply DDL, diff two schemas, generate migrations
- `completion` -- parser-native SQL autocompletion
- `analysis` / `diagnostics` -- statement classification, query spans, changed-resource extraction, error diagnostics
- `deparse` -- AST back to SQL text
- `validate` / `scope` -- post-parse semantic checks and name resolution
- `quality` -- parser quality corpora and verifiers
- `compat` / `pgregress` / `splittest` -- compatibility and regression harnesses against the upstream engine

## Quick Start

```bash
go get github.com/bytebase/omni
```

### PostgreSQL

```go
package main

import (
    "fmt"
    "github.com/bytebase/omni/pg"
    "github.com/bytebase/omni/pg/ast"
)

func main() {
    stmts, err := pg.Parse("SELECT 1; CREATE TABLE t (id int);")
    if err != nil {
        panic(err)
    }
    for _, s := range stmts {
        fmt.Printf("%-20T  %s\n", s.AST, s.Text)
    }
    // *ast.SelectStmt       SELECT 1;
    // *ast.CreateStmt        CREATE TABLE t (id int);
}
```

Root-package entry points differ per engine:

| Engine | Root package exports |
|--------|----------------------|
| `pg`, `redshift`, `mssql`, `oracle`, `mongo`, `cassandra`, `cosmosdb` | `Parse` (plus `Split`, `CollectCompletion`, and engine-specific helpers where present) |
| `partiql` | `Split` only; parsing is in `partiql/parser` |
| `elasticsearch` | `ParseElasticsearchREST`, `SplitMultiSQL`, `ClassifyRequest`, `Diagnose`, `GetQuerySpan`, `GetStatementRanges` |
| all others | no root package; import `<engine>/parser` |

## Repository Layout

```
omni/
├── <engine>/               One directory per engine (see table above)
│   ├── ast/                AST node types
│   ├── parser/             Recursive descent parser + BNF reference material
│   ├── parsertest/         Parser test cases (where present)
│   └── ...                 catalog, completion, analysis, ... per engine
├── metadata/               Database schema snapshot types, generated from proto/
├── proto/                  Protobuf sources (`make proto` regenerates)
├── harness/                Standalone differential harnesses against real engines
│                           (separate Go modules / .NET project, not built by ./...)
├── scripts/                Build, test, and agent-pipeline tooling
│   └── prompts/            Prompt templates used by the pipeline scripts
└── docs/
    ├── engine-capability-guide.md   How an engine is built up, layer by layer
    ├── PARSER-DEFENSE-MATRIX.md     Defensive test coverage per engine
    ├── scenarios/<engine>/          Scenario checklists that drove each capability
    ├── plans/                       Dated implementation plans
    ├── specs/                       Dated design specs
    └── migration/<engine>/          Bytebase migration analyses
```

Some engines also keep working documents next to the code they describe (for example `pg/parser/PAREN_AUDIT.json` is read by a lint test, and `mysql/catalog/SCENARIOS-*.md` are the source for the catalog scenario tests).

## Development

```bash
# Every engine plus the conformance harness (several engines start database containers)
make test

# Spanner differential harness; needs a running emulator (SPANNER_EMULATOR_HOST).
# CI runs this too, so a green `make test` alone is not the full CI result.
make test-spanner

# One engine
make test-pg
make test-mysql
make test-redshift

# Build everything
make build
```

## License

MIT -- see [LICENSE](LICENSE).
