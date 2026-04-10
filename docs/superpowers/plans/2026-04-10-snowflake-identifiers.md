# Snowflake Identifiers + Qualified Names (T1.1) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the `Ident` value struct and `ObjectName` Node struct to `snowflake/ast`, plus `parseIdent` / `parseObjectName` parser helpers to `snowflake/parser`, so every Tier 1+ node has identifier parsing available.

**Architecture:** `Ident` is a value struct (not a Node) carrying `{Name, Quoted, Loc}` with Snowflake-specific normalization (unquoted → UPPERCASE, quoted → as-is). `ObjectName` is a Node with three named `Ident` fields (Database, Schema, Name) matching pg's `RangeVar` shape but with per-part quoted tracking. Parser helpers accept `tokIdent`, `tokQuotedIdent`, and non-reserved keyword tokens via the unexported `keywordReserved` map.

**Tech Stack:** Go 1.25, stdlib only (`strings` for normalization).

**Spec:** `docs/superpowers/specs/2026-04-10-snowflake-identifiers-design.md` (commit `e4f0f4f`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/identifiers` on branch `feat/snowflake/identifiers`

**Commit policy:** No commits during implementation. The user reviews the full diff at the end.

---

## File Structure

### Modified

| File | Changes | Approx delta |
|------|---------|-------------|
| `snowflake/ast/parsenodes.go` | Add `import "strings"`, append `Ident` struct + methods + `ObjectName` struct + methods + compile-time assertion | +120 |
| `snowflake/ast/nodetags.go` | Append `T_ObjectName` constant + update `String()` | +5 |
| `snowflake/ast/loc.go` | Append `*ObjectName` case to `NodeLoc` switch | +2 |
| `snowflake/ast/walk_generated.go` | Regenerated via `go generate` (no new cases — ObjectName has no Node children) | ~0 net |

### Created

| File | Purpose | Approx LOC |
|------|---------|-----------|
| `snowflake/ast/identifiers_test.go` | Ident + ObjectName unit tests (Normalize, String, IsEmpty, Parts, Matches, Tag) | 200 |
| `snowflake/parser/identifiers.go` | parseIdent, parseIdentStrict, parseObjectName, ParseObjectName helpers | 120 |
| `snowflake/parser/identifiers_test.go` | Parser helper tests (bare/quoted/keyword ident, 1/2/3-part names, errors) | 250 |

**Total: ~700 LOC** across 3 new + 4 modified files.

---

## Task 1: Add Ident struct + methods to parsenodes.go

**Files:**
- Modify: `snowflake/ast/parsenodes.go`

- [ ] **Step 1: Confirm worktree state**

Run: `pwd && git rev-parse --abbrev-ref HEAD`
Expected:
```
/Users/h3n4l/OpenSource/omni/.worktrees/identifiers
feat/snowflake/identifiers
```

- [ ] **Step 2: Add import and Ident to parsenodes.go**

The file currently has no imports. Use Edit to add the import and append the Ident type after the existing `File` block (after line 22 `var _ Node = (*File)(nil)`).

First, add the import. Replace the first line `package ast` with:

```go
package ast

import "strings"
```

Then append the Ident struct after `var _ Node = (*File)(nil)`:

```go

// ---------------------------------------------------------------------------
// Identifier types
// ---------------------------------------------------------------------------

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

- [ ] **Step 3: Verify build**

Run: `go build ./snowflake/ast/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/ast/...`
Expected: no output, exit 0.

---

## Task 2: Add ObjectName struct + nodetag + NodeLoc + regenerate walker

**Files:**
- Modify: `snowflake/ast/parsenodes.go`
- Modify: `snowflake/ast/nodetags.go`
- Modify: `snowflake/ast/loc.go`
- Regenerate: `snowflake/ast/walk_generated.go`

- [ ] **Step 1: Append ObjectName to parsenodes.go**

Use Edit to append the following after the `Ident.IsEmpty` function at the end of `snowflake/ast/parsenodes.go`:

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

- [ ] **Step 2: Add T_ObjectName to nodetags.go**

In `snowflake/ast/nodetags.go`, add the constant after `T_File`:

```go
	// T_ObjectName is the tag for *ObjectName, a qualified 1/2/3-part name.
	T_ObjectName
```

And add a case to the `String()` method:

```go
	case T_ObjectName:
		return "ObjectName"
```

- [ ] **Step 3: Add *ObjectName case to NodeLoc in loc.go**

In `snowflake/ast/loc.go`, add a case to the `NodeLoc` switch before the `default`:

```go
	case *ObjectName:
		return v.Loc
```

- [ ] **Step 4: Regenerate walk_generated.go**

Run: `go generate ./snowflake/ast/...`
Expected output:
```
Generated walk_generated.go: 1 cases, 1 child fields
```

**Important:** The output still says "1 cases, 1 child fields" — the same as before T1.1. This is CORRECT. The generator skips `ObjectName` because it has zero Node-typed child fields (all fields are `Ident` value types or `Loc`). The walker's `walkChildren` switch still has only the `*File` case. This is the expected behavior for leaf Nodes with no children.

- [ ] **Step 5: Verify build and generation round-trip**

Run: `go build ./snowflake/ast/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/ast/...`
Expected: no output, exit 0.

Verify the generated file is stable:

Run: `cp snowflake/ast/walk_generated.go /tmp/walk_gen_before.go && go generate ./snowflake/ast/... && diff /tmp/walk_gen_before.go snowflake/ast/walk_generated.go && echo "BYTE_IDENTICAL" && rm /tmp/walk_gen_before.go`
Expected: `BYTE_IDENTICAL` (no diff).

Run: `go test ./snowflake/ast/...`
Expected: all existing F1 tests still pass.

---

## Task 3: Write AST unit tests (identifiers_test.go)

**Files:**
- Create: `snowflake/ast/identifiers_test.go`

- [ ] **Step 1: Write the full test file**

Create `snowflake/ast/identifiers_test.go`:

```go
package ast

import "testing"

// ---------------------------------------------------------------------------
// Ident tests
// ---------------------------------------------------------------------------

func TestIdent_Normalize(t *testing.T) {
	cases := []struct {
		name   string
		ident  Ident
		want   string
	}{
		{"unquoted lowercase", Ident{Name: "foo"}, "FOO"},
		{"unquoted uppercase", Ident{Name: "FOO"}, "FOO"},
		{"unquoted mixed", Ident{Name: "FoO"}, "FOO"},
		{"quoted lowercase", Ident{Name: "foo", Quoted: true}, "foo"},
		{"quoted uppercase", Ident{Name: "FOO", Quoted: true}, "FOO"},
		{"quoted mixed", Ident{Name: "FoO", Quoted: true}, "FoO"},
		{"empty unquoted", Ident{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ident.Normalize(); got != c.want {
				t.Errorf("Normalize() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestIdent_String(t *testing.T) {
	cases := []struct {
		name  string
		ident Ident
		want  string
	}{
		{"unquoted", Ident{Name: "foo"}, "foo"},
		{"quoted simple", Ident{Name: "foo", Quoted: true}, `"foo"`},
		{"quoted with space", Ident{Name: "my table", Quoted: true}, `"my table"`},
		{"quoted with inner quote", Ident{Name: `a"b`, Quoted: true}, `"a""b"`},
		{"empty unquoted", Ident{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ident.String(); got != c.want {
				t.Errorf("String() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestIdent_IsEmpty(t *testing.T) {
	cases := []struct {
		name  string
		ident Ident
		want  bool
	}{
		{"zero value", Ident{}, true},
		{"has name", Ident{Name: "x"}, false},
		{"quoted empty name", Ident{Quoted: true}, false},
		{"has loc only", Ident{Loc: Loc{Start: 0, End: 1}}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ident.IsEmpty(); got != c.want {
				t.Errorf("IsEmpty() = %v, want %v", got, c.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ObjectName tests
// ---------------------------------------------------------------------------

func TestObjectName_Normalize(t *testing.T) {
	cases := []struct {
		name string
		obj  ObjectName
		want string
	}{
		{"1-part unquoted", ObjectName{Name: Ident{Name: "table"}}, "TABLE"},
		{"1-part quoted", ObjectName{Name: Ident{Name: "Table", Quoted: true}}, "Table"},
		{"2-part", ObjectName{
			Schema: Ident{Name: "schema"},
			Name:   Ident{Name: "table"},
		}, "SCHEMA.TABLE"},
		{"3-part", ObjectName{
			Database: Ident{Name: "db"},
			Schema:   Ident{Name: "schema"},
			Name:     Ident{Name: "table"},
		}, "DB.SCHEMA.TABLE"},
		{"3-part mixed quoting", ObjectName{
			Database: Ident{Name: "My DB", Quoted: true},
			Schema:   Ident{Name: "schema"},
			Name:     Ident{Name: "Quoted Table", Quoted: true},
		}, "My DB.SCHEMA.Quoted Table"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.obj.Normalize(); got != c.want {
				t.Errorf("Normalize() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestObjectName_String(t *testing.T) {
	cases := []struct {
		name string
		obj  ObjectName
		want string
	}{
		{"1-part", ObjectName{Name: Ident{Name: "table"}}, "table"},
		{"2-part", ObjectName{
			Schema: Ident{Name: "schema"},
			Name:   Ident{Name: "table"},
		}, "schema.table"},
		{"3-part quoted", ObjectName{
			Database: Ident{Name: "My DB", Quoted: true},
			Schema:   Ident{Name: "schema"},
			Name:     Ident{Name: "Quoted Table", Quoted: true},
		}, `"My DB".schema."Quoted Table"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.obj.String(); got != c.want {
				t.Errorf("String() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestObjectName_Parts(t *testing.T) {
	cases := []struct {
		name     string
		obj      ObjectName
		wantLen  int
	}{
		{"1-part", ObjectName{Name: Ident{Name: "t"}}, 1},
		{"2-part", ObjectName{Schema: Ident{Name: "s"}, Name: Ident{Name: "t"}}, 2},
		{"3-part", ObjectName{Database: Ident{Name: "d"}, Schema: Ident{Name: "s"}, Name: Ident{Name: "t"}}, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			parts := c.obj.Parts()
			if len(parts) != c.wantLen {
				t.Errorf("Parts() len = %d, want %d", len(parts), c.wantLen)
			}
		})
	}
}

func TestObjectName_Matches(t *testing.T) {
	mkObj := func(db, schema, name string, dbQ, schemaQ, nameQ bool) ObjectName {
		o := ObjectName{Name: Ident{Name: name, Quoted: nameQ}}
		if schema != "" {
			o.Schema = Ident{Name: schema, Quoted: schemaQ}
		}
		if db != "" {
			o.Database = Ident{Name: db, Quoted: dbQ}
		}
		return o
	}

	cases := []struct {
		name string
		a, b ObjectName
		want bool
	}{
		{"1-part matches same name (case-folded)",
			mkObj("", "", "foo", false, false, false),
			mkObj("db", "schema", "FOO", false, false, false),
			true},
		{"1-part does NOT match different name",
			mkObj("", "", "foo", false, false, false),
			mkObj("", "", "bar", false, false, false),
			false},
		{"2-part matches same schema.name",
			mkObj("", "schema", "table", false, false, false),
			mkObj("db", "SCHEMA", "TABLE", false, false, false),
			true},
		{"2-part does NOT match different schema",
			mkObj("", "schema1", "table", false, false, false),
			mkObj("", "schema2", "table", false, false, false),
			false},
		{"3-part exact match",
			mkObj("db", "schema", "table", false, false, false),
			mkObj("DB", "SCHEMA", "TABLE", false, false, false),
			true},
		{"3-part does NOT match different db",
			mkObj("db1", "schema", "table", false, false, false),
			mkObj("db2", "schema", "table", false, false, false),
			false},
		{"quoted vs unquoted are DIFFERENT",
			mkObj("", "", "foo", false, false, true),
			mkObj("", "", "foo", false, false, false),
			false},
		{"quoted vs unquoted same normalized value",
			mkObj("", "", "FOO", false, false, true),
			mkObj("", "", "foo", false, false, false),
			true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.a.Matches(c.b); got != c.want {
				t.Errorf("Matches() = %v, want %v\n  a=%+v\n  b=%+v", got, c.want, c.a, c.b)
			}
		})
	}
}

func TestObjectName_Tag(t *testing.T) {
	var n ObjectName
	if (&n).Tag() != T_ObjectName {
		t.Errorf("Tag() = %v, want T_ObjectName", (&n).Tag())
	}
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./snowflake/ast/...`
Expected: all existing F1 tests plus all new T1.1 AST tests pass.

- [ ] **Step 3: Verbose check**

Run: `go test -v ./snowflake/ast/... -run "TestIdent_|TestObjectName_" 2>&1 | grep -E "^(=== RUN|--- (PASS|FAIL))" | tail -40`
Expected: every subtest reported as PASS.

---

## Task 4: Write parser helpers (identifiers.go)

**Files:**
- Create: `snowflake/parser/identifiers.go`

- [ ] **Step 1: Write the full parser helpers file**

Create `snowflake/parser/identifiers.go`:

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
		// Non-reserved keyword used as an identifier.
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
//
//	table
//	schema.table
//	database.schema.table
//
// Starts by parsing one identifier, then greedily consumes up to two more
// dot-separated identifiers if present. Returns *ast.ObjectName with the
// correct parts populated and Loc spanning all parts.
func (p *Parser) parseObjectName() (*ast.ObjectName, error) {
	first, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	if p.cur.Type != '.' {
		// 1-part name.
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
		// 2-part name: schema.table.
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

	// 3-part name: database.schema.table.
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
//
//	ParseObjectName("my_table")
//	ParseObjectName("schema.table")
//	ParseObjectName(`"My DB".schema."Quoted Table"`)
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

- [ ] **Step 2: Verify build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go test ./snowflake/parser/...`
Expected: all existing F2/F3/F4 tests still pass.

---

## Task 5: Write parser helper tests (identifiers_test.go)

**Files:**
- Create: `snowflake/parser/identifiers_test.go`

- [ ] **Step 1: Write the full test file**

Create `snowflake/parser/identifiers_test.go`:

```go
package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseIdent is a helper that constructs a Parser from input and
// calls parseIdent. Since parseIdent is unexported, tests in the same
// package can call it directly through this helper.
func testParseIdent(input string) (ast.Ident, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseIdent()
}

// testParseIdentStrict wraps parseIdentStrict the same way.
func testParseIdentStrict(input string) (ast.Ident, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseIdentStrict()
}

// testParseObjectName wraps parseObjectName the same way.
func testParseObjectName(input string) (*ast.ObjectName, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseObjectName()
}

// ---------------------------------------------------------------------------
// parseIdent tests
// ---------------------------------------------------------------------------

func TestParseIdent_Bare(t *testing.T) {
	id, err := testParseIdent("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || id.Quoted {
		t.Errorf("got %+v, want Ident{foo, unquoted}", id)
	}
}

func TestParseIdent_Quoted(t *testing.T) {
	id, err := testParseIdent(`"foo"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || !id.Quoted {
		t.Errorf("got %+v, want Ident{foo, quoted}", id)
	}
}

func TestParseIdent_QuotedWithSpaces(t *testing.T) {
	id, err := testParseIdent(`"my table"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "my table" || !id.Quoted {
		t.Errorf("got %+v, want Ident{my table, quoted}", id)
	}
}

func TestParseIdent_QuotedEscape(t *testing.T) {
	// F2's lexer un-escapes "" → " in Token.Str.
	id, err := testParseIdent(`"a""b"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != `a"b` || !id.Quoted {
		t.Errorf("got %+v, want Ident{a\"b, quoted}", id)
	}
}

func TestParseIdent_NonReservedKeyword(t *testing.T) {
	// ALERT is a non-reserved keyword — should be accepted as identifier.
	id, err := testParseIdent("ALERT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "ALERT" || id.Quoted {
		t.Errorf("got %+v, want Ident{ALERT, unquoted}", id)
	}
}

func TestParseIdent_ReservedKeyword(t *testing.T) {
	// SELECT is a reserved keyword — should be rejected.
	_, err := testParseIdent("SELECT")
	if err == nil {
		t.Fatal("expected error for reserved keyword SELECT, got nil")
	}
	if pe, ok := err.(*ParseError); ok {
		if !strings.Contains(pe.Msg, "expected identifier") {
			t.Errorf("error Msg = %q, want 'expected identifier'", pe.Msg)
		}
	} else {
		t.Errorf("expected *ParseError, got %T", err)
	}
}

func TestParseIdent_EOF(t *testing.T) {
	_, err := testParseIdent("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestParseIdent_LocCorrect(t *testing.T) {
	// "   foo" — the identifier starts at byte 3.
	id, err := testParseIdent("   foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Loc.Start != 3 || id.Loc.End != 6 {
		t.Errorf("Loc = %+v, want {3, 6}", id.Loc)
	}
}

// ---------------------------------------------------------------------------
// parseIdentStrict tests
// ---------------------------------------------------------------------------

func TestParseIdentStrict_AcceptsBare(t *testing.T) {
	id, err := testParseIdentStrict("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || id.Quoted {
		t.Errorf("got %+v, want Ident{foo, unquoted}", id)
	}
}

func TestParseIdentStrict_AcceptsQuoted(t *testing.T) {
	id, err := testParseIdentStrict(`"foo"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id.Name != "foo" || !id.Quoted {
		t.Errorf("got %+v, want Ident{foo, quoted}", id)
	}
}

func TestParseIdentStrict_RejectsNonReservedKeyword(t *testing.T) {
	// parseIdentStrict rejects ALL keywords, even non-reserved ones.
	_, err := testParseIdentStrict("ALERT")
	if err == nil {
		t.Fatal("expected error for keyword ALERT in strict mode, got nil")
	}
}

// ---------------------------------------------------------------------------
// parseObjectName tests
// ---------------------------------------------------------------------------

func TestParseObjectName_OnePart(t *testing.T) {
	obj, err := testParseObjectName("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Name.Name != "foo" || !obj.Schema.IsEmpty() || !obj.Database.IsEmpty() {
		t.Errorf("got %+v, want 1-part name foo", obj)
	}
}

func TestParseObjectName_TwoParts(t *testing.T) {
	obj, err := testParseObjectName("schema.table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Schema.Name != "schema" || obj.Name.Name != "table" || !obj.Database.IsEmpty() {
		t.Errorf("got Schema=%q Name=%q DB=%q, want schema.table",
			obj.Schema.Name, obj.Name.Name, obj.Database.Name)
	}
	// Loc should span both parts.
	if obj.Loc.Start != 0 || obj.Loc.End != 12 {
		t.Errorf("Loc = %+v, want {0, 12}", obj.Loc)
	}
}

func TestParseObjectName_ThreeParts(t *testing.T) {
	obj, err := testParseObjectName("db.schema.table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Database.Name != "db" || obj.Schema.Name != "schema" || obj.Name.Name != "table" {
		t.Errorf("got DB=%q Schema=%q Name=%q, want db.schema.table",
			obj.Database.Name, obj.Schema.Name, obj.Name.Name)
	}
	if obj.Loc.Start != 0 || obj.Loc.End != 15 {
		t.Errorf("Loc = %+v, want {0, 15}", obj.Loc)
	}
}

func TestParseObjectName_Mixed(t *testing.T) {
	// "My DB".schema."Quoted Table"
	obj, err := testParseObjectName(`"My DB".schema."Quoted Table"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !obj.Database.Quoted || obj.Database.Name != "My DB" {
		t.Errorf("Database = %+v, want quoted 'My DB'", obj.Database)
	}
	if obj.Schema.Quoted || obj.Schema.Name != "schema" {
		t.Errorf("Schema = %+v, want unquoted 'schema'", obj.Schema)
	}
	if !obj.Name.Quoted || obj.Name.Name != "Quoted Table" {
		t.Errorf("Name = %+v, want quoted 'Quoted Table'", obj.Name)
	}
}

func TestParseObjectName_TrailingDot(t *testing.T) {
	// "foo." — error after the dot (no identifier follows).
	_, err := testParseObjectName("foo.")
	if err == nil {
		t.Fatal("expected error for trailing dot, got nil")
	}
}

func TestParseObjectName_LeadingDot(t *testing.T) {
	// ".foo" — the dot is not an identifier.
	_, err := testParseObjectName(".foo")
	if err == nil {
		t.Fatal("expected error for leading dot, got nil")
	}
}

func TestParseObjectName_DoubleDot(t *testing.T) {
	// "foo..bar" — error at the second dot.
	_, err := testParseObjectName("foo..bar")
	if err == nil {
		t.Fatal("expected error for double dot, got nil")
	}
}

func TestParseObjectName_Freestanding(t *testing.T) {
	// ParseObjectName is the public freestanding helper.
	obj, errs := ParseObjectName("db.schema.table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if obj.Database.Name != "db" || obj.Schema.Name != "schema" || obj.Name.Name != "table" {
		t.Errorf("ParseObjectName mismatch: %+v", obj)
	}
}

func TestParseObjectName_FreestandingTrailingTokens(t *testing.T) {
	// "foo bar" — parses "foo" then complains about trailing "bar".
	obj, errs := ParseObjectName("foo bar")
	if obj == nil {
		t.Fatal("expected non-nil ObjectName even with trailing tokens")
	}
	if obj.Name.Name != "foo" {
		t.Errorf("Name = %q, want foo", obj.Name.Name)
	}
	if len(errs) != 1 || !strings.Contains(errs[0].Msg, "unexpected token") {
		t.Errorf("expected 1 'unexpected token' error, got %+v", errs)
	}
}

func TestParseObjectName_FreestandingEmpty(t *testing.T) {
	_, errs := ParseObjectName("")
	if len(errs) == 0 {
		t.Error("expected error for empty input, got none")
	}
}

func TestParseObjectName_KeywordAsPart(t *testing.T) {
	// Non-reserved keyword ALERT should be usable as any part.
	obj, err := testParseObjectName("ALERT.ALERT.ALERT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Database.Name != "ALERT" || obj.Schema.Name != "ALERT" || obj.Name.Name != "ALERT" {
		t.Errorf("got %+v, want ALERT.ALERT.ALERT", obj)
	}
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./snowflake/...`
Expected: clean pass. All F1 + F2 + F3 + F4 + T1.1 tests green.

- [ ] **Step 3: Verbose check**

Run: `go test -v ./snowflake/parser/... -run "TestParseIdent|TestParseObjectName" 2>&1 | grep -E "^(=== RUN|--- (PASS|FAIL))" | tail -40`
Expected: every subtest reported as PASS.

---

## Task 6: Final acceptance criteria sweep

- [ ] **Step 1: Build**

Run: `go build ./snowflake/...`
Expected: no output, exit 0.

- [ ] **Step 2: Vet**

Run: `go vet ./snowflake/...`
Expected: no output, exit 0.

- [ ] **Step 3: Gofmt**

Run: `gofmt -l snowflake/`
Expected: no output.

If any file is listed, apply `gofmt -w snowflake/` and re-run.

- [ ] **Step 4: Test**

Run: `go test ./snowflake/...`
Expected:
```
ok  	github.com/bytebase/omni/snowflake/ast	(some duration)
ok  	github.com/bytebase/omni/snowflake/parser	(some duration)
```

Both packages pass.

- [ ] **Step 5: Walker generation round-trip**

Run: `cp snowflake/ast/walk_generated.go /tmp/walk_gen_check.go && go generate ./snowflake/ast/... && diff /tmp/walk_gen_check.go snowflake/ast/walk_generated.go && echo "BYTE_IDENTICAL" && rm /tmp/walk_gen_check.go`
Expected: `BYTE_IDENTICAL`.

- [ ] **Step 6: List all changed files**

Run: `git status snowflake/`
Expected:
```
modified:   snowflake/ast/loc.go
modified:   snowflake/ast/nodetags.go
modified:   snowflake/ast/parsenodes.go

Untracked files:
	snowflake/ast/identifiers_test.go
	snowflake/parser/identifiers.go
	snowflake/parser/identifiers_test.go
```

(`walk_generated.go` should NOT show as modified because the regen produced byte-identical output.)

- [ ] **Step 7: STOP and present for review**

Do NOT commit. Output a summary:

> T1.1 (snowflake identifiers + qualified names) implementation complete.
>
> All 6 tasks done. Acceptance criteria green:
> - `go build ./snowflake/...` ✓
> - `go vet ./snowflake/...` ✓
> - `gofmt -l snowflake/` clean ✓
> - `go test ./snowflake/...` ✓ (both ast and parser)
> - Walker generation round-trip ✓
>
> Files modified/created:
> - snowflake/ast/parsenodes.go (Ident + ObjectName types)
> - snowflake/ast/nodetags.go (T_ObjectName)
> - snowflake/ast/loc.go (*ObjectName case)
> - snowflake/ast/identifiers_test.go (NEW)
> - snowflake/parser/identifiers.go (NEW)
> - snowflake/parser/identifiers_test.go (NEW)

---

## Spec Coverage Checklist

| Spec section | Covered by |
|---|---|
| Ident struct (Name, Quoted, Loc) | Task 1 |
| Ident.Normalize / String / IsEmpty | Task 1 |
| ObjectName struct (Database, Schema, Name as Ident, Loc) | Task 2 |
| ObjectName.Tag / Normalize / String / Parts / Matches | Task 2 |
| T_ObjectName constant + String() | Task 2 |
| NodeLoc *ObjectName case | Task 2 |
| Walker regeneration | Task 2 |
| Ident unit tests (8 categories) | Task 3 |
| parseIdent / parseIdentStrict | Task 4 |
| parseObjectName | Task 4 |
| ParseObjectName freestanding | Task 4 |
| Parser helper tests (16 categories) | Task 5 |
| Acceptance criteria | Task 6 |
