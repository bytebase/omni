# redshift -- Amazon Redshift Parser

Recursive descent parser for Amazon Redshift SQL, producing a full AST with position tracking. The engine is a fork of `pg/` that tracks Redshift's PostgreSQL 8.0 heritage plus its own extensions (COPY/UNLOAD, DISTKEY/SORTKEY, external schemas, and so on).

## Public API

```go
import "github.com/bytebase/omni/redshift"

stmts, err := redshift.Parse(sql)
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

### Bytebase runtime surfaces

The root package also exposes the helpers Bytebase wires into its product features:

| Function | Purpose |
|----------|---------|
| `Split` | Statement splitting without building an AST |
| `StatementRanges` | Byte ranges of each statement |
| `Diagnose` | Syntax diagnostics for the SQL editor |
| `ValidateSQLForEditor` | Read-only check for the SQL editor |
| `ClassifyStatement`, `GetStatementTypes` | Statement classification and report |
| `ExtractChangedResources` | Tables touched by DDL/DML |
| `CollectCompletion` | Parser-native autocompletion |

## Packages

| Package | Description |
|---------|-------------|
| `redshift/ast` | AST node types, forked from `pg/ast` |
| `redshift/parser` | Recursive descent parser |
| `redshift/catalog` | In-memory catalog simulation, DDL semantic analysis |
| `redshift/completion` | Parser-native C3-style SQL completion |
| `redshift/analysis` | Statement classification, query span, changed resources |
| `redshift/plpgsql` | PL/pgSQL parser for function bodies |
| `redshift/compat` | Compatibility report against the AWS command manifest |
| `redshift/parsertest` | Parser test cases organized by SQL feature |
| `redshift/pgregress` | PostgreSQL regression test compatibility |

Parser maintenance conventions are shared with `pg/parser`; see `redshift/parser/CLAUDE.md`.
