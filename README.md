# Omni

SQL toolchain for multiple database engines. Each engine provides a parser, AST, and additional components such as catalog simulation and semantic analysis -- all in pure Go with zero dependencies.

## Features

- **Zero dependencies** -- pure Go, no CGo, no generated code at runtime
- **Full AST** -- every parsed statement produces a complete abstract syntax tree
- **Position tracking** -- every AST node carries byte-offset location info
- **Beyond parsing** -- catalog simulation, DDL semantic analysis, and more per engine

## Engine Status

| Engine | Package | Batches | Status |
|--------|---------|---------|--------|
| PostgreSQL | `pg/` | 74/77 | Production-ready, 3 minor gaps |
| MySQL | `mysql/` | 97/97 | Complete |
| SQL Server | `mssql/` | 108/108 | Complete |
| Oracle | `oracle/` | 60/80 | In progress |

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

## Architecture

```
omni/
├── pg/                     PostgreSQL
│   ├── parse.go            Public API: Parse(sql) → []Statement
│   ├── ast/                210+ AST node types
│   ├── parser/             Recursive descent parser (~29K lines)
│   ├── catalog/            In-memory catalog simulation & DDL analysis
│   ├── parsertest/         746 test cases
│   └── pgregress/          PostgreSQL regression test compatibility
├── mysql/                  MySQL
│   ├── ast/                AST node types
│   ├── parser/             Recursive descent parser
│   └── parsertest/         Test cases
├── mssql/                  SQL Server (T-SQL)
│   ├── ast/                AST node types
│   ├── parser/             Recursive descent parser
│   └── parsertest/         Test cases
├── oracle/                 Oracle
│   ├── ast/                AST node types
│   └── parser/             Recursive descent parser
└── scripts/                Shared build & audit tooling
```

## Development

```bash
# Run all tests
make test

# Test a specific engine
make test-pg
make test-mysql
make test-mssql
make test-oracle

# Build everything
make build
```

## License

MIT -- see [LICENSE](LICENSE).
