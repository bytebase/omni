# Snowflake Identifiers + Qualified Names (T1.1) — Design

**DAG node:** T1.1 — identifiers + qualified names + normalization helpers
**Migration:** `docs/migration/snowflake/dag.md`
**Branch:** `feat/snowflake/identifiers`
**Status:** Approved, ready for plan + implementation.
**Depends on:** F4 (parser-entry, merged via PR #19).
**Unblocks:** T1.2 (data types), T1.3 (expressions), T1.4 (SELECT core), and transitively every Tier 1+ node that parses identifier-shaped tokens.

## Purpose

T1.1 adds the identifier and qualified-name building blocks that every subsequent Tier 1+ parser node needs. It ships:

1. **`Ident` struct** in `snowflake/ast` — a single identifier with Name, Quoted flag, and Loc. NOT a Node (no Tag method) — it's a value struct embedded in parent nodes.
2. **`ObjectName` struct** in `snowflake/ast` — a 1/2/3-part qualified name (database.schema.name) with three named Ident fields. IS a Node (`*ObjectName` satisfies Node via `Tag() NodeTag`).
3. **Normalization helpers** on Ident and ObjectName — `Normalize()`, `String()`, `IsEmpty()`, `Parts()`, `Matches()`.
4. **Parser helpers** in `snowflake/parser` — `parseIdent()`, `parseIdentStrict()`, `parseObjectName()`, `ParseObjectName()` (freestanding).

T1.1 does NOT replace any dispatch cases in F4's `parseStmt`. It adds helper functions that Tier 1.2/1.3/1.4/etc. will call when they parse concrete statements containing identifiers.

## Why Snowflake needs a dedicated Ident type (unlike pg/mysql)

pg and mysql don't have a per-identifier `Quoted` flag because their case-folding happens at lex time (pg lower-cases, mysql preserves). Snowflake's case-folding is different:

- **Unquoted identifiers**: source case is preserved through lexing and parsing; case-folding to UPPERCASE happens at **resolution time** (by the query planner).
- **Quoted identifiers**: case-sensitive, preserved as-is through resolution.

Bytebase's masking and query-span features need to know whether an identifier was quoted to correctly distinguish `foo` (resolves to `FOO`) from `"foo"` (resolves to `foo`). Storing just the string would lose this distinction. The `Quoted bool` flag preserves it at zero runtime cost.

## AST Type Definitions

### `Ident` (value struct, NOT a Node)

Added to `snowflake/ast/parsenodes.go`:

```go
// Ident represents a single identifier — a name used to reference a database
// object (table, column, schema, etc.).
//
// Name is the raw text from source: for quoted identifiers, the content
// between the double-quotes with "" un-escaped; for unquoted identifiers,
// the source bytes with case preserved.
//
// Quoted reports whether the source used "..." quoting. This matters because
// Snowflake case-folds unquoted identifiers to uppercase at resolution time,
// while quoted identifiers preserve case.
//
// Ident is a value struct, NOT a Node. It is embedded by value in parent
// nodes (e.g. ObjectName) and is not visited by the AST walker.
//
// The zero value (Name == "" && Quoted == false) represents an absent
// identifier — used by ObjectName for unused parts (e.g. a 1-part name
// has zero Database and Schema).
type Ident struct {
    Name   string
    Quoted bool
    Loc    Loc
}

// Normalize returns the canonical form of the identifier per Snowflake
// resolution rules:
//   - Quoted identifiers: returned as-is (case-sensitive)
//   - Unquoted identifiers: uppercased
func (i Ident) Normalize() string {
    if i.Quoted {
        return i.Name
    }
    return strings.ToUpper(i.Name)
}

// String returns the source form of the identifier, re-quoting if it was
// originally quoted. Inner " characters are escaped as "". Useful for
// deparse and error messages.
func (i Ident) String() string {
    if !i.Quoted {
        return i.Name
    }
    return `"` + strings.ReplaceAll(i.Name, `"`, `""`) + `"`
}

// IsEmpty reports whether the identifier is the zero value (absent).
// Used by ObjectName to check whether a part (Database, Schema) is present.
func (i Ident) IsEmpty() bool {
    return i.Name == "" && !i.Quoted
}
```

Note: `Ident` needs `import "strings"` in `parsenodes.go`. Since the current file doesn't have any imports, this will add an import block.

### `ObjectName` (Node)

Also in `snowflake/ast/parsenodes.go`:

```go
// ObjectName represents a qualified object name (1/2/3-part) like
// `table`, `schema.table`, or `database.schema.table`.
//
// For 1-part names, only Name is set (Database and Schema are zero Idents).
// For 2-part names, Schema and Name are set.
// For 3-part names, all three are set.
//
// ObjectName is a Node and can be used as a child in the AST tree. The
// walker visits *ObjectName but does NOT descend into the embedded Ident
// fields (they are value structs, not Nodes).
type ObjectName struct {
    Database Ident // may be zero (IsEmpty)
    Schema   Ident // may be zero (IsEmpty)
    Name     Ident // always present for a valid ObjectName
    Loc      Loc
}

// Tag implements Node.
func (n *ObjectName) Tag() NodeTag { return T_ObjectName }

// Compile-time assertion that *ObjectName satisfies Node.
var _ Node = (*ObjectName)(nil)

// Normalize returns the canonical dotted form with each non-empty part
// normalized per Snowflake resolution rules.
//
// Examples:
//   - 1-part: "TABLE"
//   - 2-part: "SCHEMA.TABLE"
//   - 3-part: "DB.SCHEMA.TABLE"
func (n ObjectName) Normalize() string {
    parts := n.Parts()
    normalized := make([]string, len(parts))
    for i, p := range parts {
        normalized[i] = p.Normalize()
    }
    return strings.Join(normalized, ".")
}

// String returns the source form with dots. Each part is re-quoted if it
// was originally quoted.
func (n ObjectName) String() string {
    parts := n.Parts()
    strs := make([]string, len(parts))
    for i, p := range parts {
        strs[i] = p.String()
    }
    return strings.Join(strs, ".")
}

// Parts returns the non-empty parts in order. Length is 1, 2, or 3.
func (n ObjectName) Parts() []Ident {
    switch {
    case !n.Database.IsEmpty():
        return []Ident{n.Database, n.Schema, n.Name}
    case !n.Schema.IsEmpty():
        return []Ident{n.Schema, n.Name}
    default:
        return []Ident{n.Name}
    }
}

// Matches reports whether this ObjectName suffix-matches other using
// normalized (case-folded) comparison. A 1-part name matches any other
// with the same normalized Name. A 2-part name matches any other with
// the same normalized Schema + Name. A 3-part name requires all three
// parts to match.
func (n ObjectName) Matches(other ObjectName) bool {
    if n.Name.Normalize() != other.Name.Normalize() {
        return false
    }
    if n.Schema.IsEmpty() {
        return true // 1-part match
    }
    if n.Schema.Normalize() != other.Schema.Normalize() {
        return false
    }
    if n.Database.IsEmpty() {
        return true // 2-part match
    }
    return n.Database.Normalize() == other.Database.Normalize()
}
```

### NodeTag additions

Add to `snowflake/ast/nodetags.go`:

```go
T_ObjectName
```

(After `T_File`. `T_Ident` is NOT added because Ident is not a Node.)

Also update `NodeTag.String()` to include the new tag.

### Walker regeneration

After adding `ObjectName` to `parsenodes.go`, run `go generate ./snowflake/ast/...` to regenerate `walk_generated.go`. The generator will add a case for `*ObjectName` with NO child walks (all fields are non-Node value types: `Ident`, `Loc`).

### `NodeLoc` update

Add a case for `*ObjectName` in `snowflake/ast/loc.go`'s `NodeLoc` switch:

```go
case *ObjectName:
    return v.Loc
```

## Parser Helpers

Added in a new file `snowflake/parser/identifiers.go`:

```go
package parser

import "github.com/bytebase/omni/snowflake/ast"

// parseIdent parses one identifier from the token stream. Accepts:
//   - tokIdent (bare identifier)
//   - tokQuotedIdent (double-quoted identifier)
//   - Any non-reserved keyword token (type >= 700 and not in keywordReserved)
//
// Returns a ParseError if the current token is none of the above.
func (p *Parser) parseIdent() (ast.Ident, error) {
    switch {
    case p.cur.Type == tokIdent:
        tok := p.advance()
        return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
    case p.cur.Type == tokQuotedIdent:
        tok := p.advance()
        return ast.Ident{Name: tok.Str, Quoted: true, Loc: tok.Loc}, nil
    case p.cur.Type >= 700 && !keywordReserved[p.cur.Type]:
        tok := p.advance()
        return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
    }
    return ast.Ident{}, &ParseError{
        Loc: p.cur.Loc,
        Msg: "expected identifier",
    }
}

// parseIdentStrict parses one identifier but ONLY accepts tokIdent or
// tokQuotedIdent — NOT non-reserved keywords. Used in contexts where
// keyword-as-identifier is not allowed (rare in Snowflake).
func (p *Parser) parseIdentStrict() (ast.Ident, error) {
    switch p.cur.Type {
    case tokIdent:
        tok := p.advance()
        return ast.Ident{Name: tok.Str, Quoted: false, Loc: tok.Loc}, nil
    case tokQuotedIdent:
        tok := p.advance()
        return ast.Ident{Name: tok.Str, Quoted: true, Loc: tok.Loc}, nil
    }
    return ast.Ident{}, &ParseError{
        Loc: p.cur.Loc,
        Msg: "expected identifier",
    }
}

// parseObjectName parses a dotted object name with 1, 2, or 3 parts:
//   table
//   schema.table
//   database.schema.table
//
// Starts by parsing one identifier, then greedily consumes up to two more
// .-separated identifiers if present. Returns *ast.ObjectName with the
// correct parts populated.
func (p *Parser) parseObjectName() (*ast.ObjectName, error) {
    first, err := p.parseIdent()
    if err != nil {
        return nil, err
    }

    if p.cur.Type != '.' {
        // 1-part name
        return &ast.ObjectName{
            Name: first,
            Loc:  first.Loc,
        }, nil
    }
    p.advance() // consume first dot

    second, err := p.parseIdent()
    if err != nil {
        return nil, err
    }

    if p.cur.Type != '.' {
        // 2-part name: schema.table
        return &ast.ObjectName{
            Schema: first,
            Name:   second,
            Loc:    ast.Loc{Start: first.Loc.Start, End: second.Loc.End},
        }, nil
    }
    p.advance() // consume second dot

    third, err := p.parseIdent()
    if err != nil {
        return nil, err
    }

    // 3-part name: database.schema.table
    return &ast.ObjectName{
        Database: first,
        Schema:   second,
        Name:     third,
        Loc:      ast.Loc{Start: first.Loc.Start, End: third.Loc.End},
    }, nil
}

// ParseObjectName parses an object name from a standalone string. Returns
// the ObjectName and any ParseErrors encountered. Useful for catalog
// lookups, tests, and callers that have a name string but not a token
// stream.
//
// Examples:
//   ParseObjectName("my_table")
//   ParseObjectName("schema.table")
//   ParseObjectName(`"My DB".schema."Quoted Table"`)
func ParseObjectName(input string) (*ast.ObjectName, []ParseError) {
    p := &Parser{
        lexer: NewLexer(input),
        input: input,
    }
    p.advance()

    obj, err := p.parseObjectName()
    if err != nil {
        if pe, ok := err.(*ParseError); ok {
            return nil, []ParseError{*pe}
        }
        return nil, []ParseError{{Msg: err.Error()}}
    }

    // Check for trailing tokens (should be EOF).
    if p.cur.Type != tokEOF {
        return obj, []ParseError{{
            Loc: p.cur.Loc,
            Msg: "unexpected token after object name",
        }}
    }

    return obj, nil
}
```

## File Layout

### Modified files

| File | Change |
|------|--------|
| `snowflake/ast/parsenodes.go` | Append `Ident` struct + methods, `ObjectName` struct + methods. Add `import "strings"` |
| `snowflake/ast/nodetags.go` | Append `T_ObjectName` constant + update `String()` |
| `snowflake/ast/loc.go` | Append `*ObjectName` case to `NodeLoc` switch |
| `snowflake/ast/walk_generated.go` | Regenerated via `go generate` (adds empty `*ObjectName` case) |

### Created files

| File | Purpose | Approx LOC |
|------|---------|-----------|
| `snowflake/parser/identifiers.go` | `parseIdent`, `parseIdentStrict`, `parseObjectName`, `ParseObjectName` | 120 |
| `snowflake/parser/identifiers_test.go` | Parser helper tests | 250 |
| `snowflake/ast/identifiers_test.go` | Ident/ObjectName unit tests (Normalize, String, Matches, etc.) | 200 |

**Estimated total: ~570 LOC** of new code across 3 new + 4 modified files.

## Testing

### `snowflake/ast/identifiers_test.go`

1. **TestIdent_Normalize** — `Ident{"foo", false, _}.Normalize() == "FOO"`; `Ident{"foo", true, _}.Normalize() == "foo"`
2. **TestIdent_String** — unquoted returns Name as-is; quoted re-wraps in `"..."` with `""` escape
3. **TestIdent_IsEmpty** — zero value is empty; non-empty Name is not empty
4. **TestObjectName_Normalize** — 1/2/3-part: `"TABLE"`, `"SCHEMA.TABLE"`, `"DB.SCHEMA.TABLE"`
5. **TestObjectName_String** — round-trip: `"My DB".schema.TABLE`
6. **TestObjectName_Parts** — correct count and order for 1/2/3-part names
7. **TestObjectName_Matches** — suffix match: 1-part matches 3-part with same name; 2-part matches only if schema matches; 3-part exact; quoted-vs-unquoted are DIFFERENT
8. **TestObjectName_Tag** — `(*ObjectName)(nil).Tag() == T_ObjectName` (compile-time + runtime assertion)

### `snowflake/parser/identifiers_test.go`

1. **TestParseIdent_Bare** — `foo` → `Ident{Name: "foo", Quoted: false}`
2. **TestParseIdent_Quoted** — `"foo"` → `Ident{Name: "foo", Quoted: true}`
3. **TestParseIdent_QuotedWithSpaces** — `"my table"` → `Ident{Name: "my table", Quoted: true}`
4. **TestParseIdent_QuotedEscape** — `"a""b"` → `Ident{Name: "a\"b", Quoted: true}` (F2 lexer already un-escapes)
5. **TestParseIdent_NonReservedKeyword** — `ALERT` → `Ident{Name: "ALERT", Quoted: false}`
6. **TestParseIdent_ReservedKeyword** — `SELECT` → ParseError "expected identifier"
7. **TestParseIdent_EOF** — empty → ParseError
8. **TestParseObjectName_OnePart** — `foo` → `ObjectName{Name: {"foo", false, _}}`
9. **TestParseObjectName_TwoParts** — `schema.table` → `ObjectName{Schema: _, Name: _}`
10. **TestParseObjectName_ThreeParts** — `db.schema.table` → all three set
11. **TestParseObjectName_Mixed** — `"My DB".schema."Quoted Table"` → per-part Quoted flags correct
12. **TestParseObjectName_TrailingDot** — `foo.` → error after the dot
13. **TestParseObjectName_LeadingDot** — `.foo` → error at the dot (not an identifier)
14. **TestParseObjectName_DoubleDot** — `foo..bar` → error at the second dot
15. **TestParseObjectName_FreestandingHelper** — `ParseObjectName("db.schema.table")` → correct ObjectName + no errors; `ParseObjectName("foo bar")` → ObjectName + trailing-token error
16. **TestParseIdentStrict_RejectsKeyword** — `ALERT` → error (even though ALERT is non-reserved)

Run scope: `go test ./snowflake/...`

## Out of Scope

| Feature | Where it lives |
|---|---|
| Column references (`table.column` in expressions) | T1.3 (expressions) |
| FROM clause alias handling (`FROM t AS alias`) | T1.4 (SELECT core) |
| Table function syntax (`FLATTEN(...)`) | T5.3 |
| Type names (`VARCHAR(100)`, `NUMBER(38,0)`) | T1.2 (data types) |
| Catalog/schema-resolution logic | `snowflake/catalog` (not yet in DAG) |
| Deparse (ObjectName → SQL string round-trip) | T3.2 |

## Acceptance Criteria

T1.1 is complete when:

1. `go build ./snowflake/...` succeeds.
2. `go vet ./snowflake/...` clean.
3. `gofmt -l snowflake/` clean.
4. `go test ./snowflake/...` passes — all F1/F2/F3/F4 tests still green + all T1.1 tests.
5. `go generate ./snowflake/ast/...` produces a byte-identical `walk_generated.go`.
6. `Ident.Normalize()` correctly uppercases unquoted identifiers and preserves quoted identifiers.
7. `ObjectName.Matches()` correctly handles suffix-match semantics with normalized comparison.
8. `parseIdent()` accepts `tokIdent`, `tokQuotedIdent`, and non-reserved keyword tokens — rejects reserved keywords.
9. `parseObjectName()` correctly parses 1/2/3-part dotted names with correct Loc spanning all parts.
10. `ParseObjectName(input)` freestanding helper works correctly for standalone string inputs.
11. After merge, `docs/migration/snowflake/dag.md` T1.1 status is flipped to `done`.

## Files Created / Modified

```
snowflake/ast/parsenodes.go          (MODIFIED: +Ident, +ObjectName, +methods, +import "strings")
snowflake/ast/nodetags.go            (MODIFIED: +T_ObjectName, +String case)
snowflake/ast/loc.go                 (MODIFIED: +*ObjectName case in NodeLoc)
snowflake/ast/walk_generated.go      (REGENERATED via go generate)
snowflake/ast/identifiers_test.go    (NEW: Ident/ObjectName unit tests)
snowflake/parser/identifiers.go      (NEW: parseIdent, parseIdentStrict, parseObjectName, ParseObjectName)
snowflake/parser/identifiers_test.go (NEW: parser helper tests)
docs/superpowers/specs/2026-04-10-snowflake-identifiers-design.md (this file)
```

Estimated total: ~570 LOC across 3 new + 4 modified files.
