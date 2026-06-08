# pg -- PostgreSQL Parser

Recursive descent parser for PostgreSQL 17, producing a full AST with position tracking.

## Public API

```go
import "github.com/bytebase/omni/pg"

stmts, err := pg.Parse(sql)
```

### `Parse(sql string) ([]Statement, error)`

Splits and parses a SQL string into individual statements.

### `Statement`

```go
type Statement struct {
    Text      string       // SQL text including trailing semicolon
    AST       ast.Node     // Inner statement node (e.g. *ast.SelectStmt)
    ByteStart int          // Inclusive start byte offset
    ByteEnd   int          // Exclusive end byte offset
    Start     Position     // Start line:column (1-based)
    End       Position     // End line:column (1-based)
}
```

## Packages

| Package | Description |
|---------|-------------|
| `pg/ast` | 210+ AST node types matching PostgreSQL internals |
| `pg/parser` | Recursive descent parser (~29,000 lines, 39 files) |
| `pg/catalog` | In-memory catalog simulation, DDL semantic analysis, type resolution |
| `pg/parsertest` | 746 test cases organized by SQL feature |
| `pg/pgregress` | PostgreSQL official regression test compatibility |
