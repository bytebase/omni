# Snowflake Data Types (T1.2) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `TypeName` Node struct with `TypeKind` enum (22 values) to `snowflake/ast`, plus `parseDataType()` parser helper to `snowflake/parser`, covering every Snowflake data type form from the legacy grammar's `data_type` rule.

**Architecture:** Single `TypeName` struct with `Kind TypeKind` enum for fast dispatch, `Name string` for round-tripping, `Params []int` for numeric parameters, recursive `ElementType *TypeName` for ARRAY/VECTOR, and `VectorDim int` for VECTOR dimensions. Multi-word types (DOUBLE PRECISION, CHAR VARYING, NCHAR VARYING) fused via parser lookahead. Timestamp variants (TIMESTAMP_LTZ/NTZ/TZ) added to F2's keyword map as scope creep.

**Tech Stack:** Go 1.25, stdlib only (`strings` for case-folding in multi-word lookahead).

**Spec:** `docs/superpowers/specs/2026-04-10-snowflake-datatypes-design.md` (commit `594c578`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/datatypes` on branch `feat/snowflake/datatypes`

**Commit policy:** No commits during implementation. User reviews the full diff at the end.

---

## File Structure

### Modified

| File | Changes |
|------|---------|
| `snowflake/parser/tokens.go` | Append 3 keyword constants (`kwTIMESTAMP_LTZ/NTZ/TZ`) after `kwTIMESTAMP` block |
| `snowflake/parser/keywords.go` | Append 6 keywordMap entries for timestamp variants |
| `snowflake/ast/parsenodes.go` | Append TypeKind enum (22 values) + TypeName struct + methods |
| `snowflake/ast/nodetags.go` | Append T_TypeName + String() case |
| `snowflake/ast/loc.go` | Append *TypeName case to NodeLoc |
| `snowflake/ast/walk_generated.go` | Regenerated — gains *TypeName case with ElementType child walk |

### Created

| File | Purpose | Approx LOC |
|------|---------|-----------|
| `snowflake/parser/datatypes.go` | parseDataType, parseOptionalTypeParams, parseVectorElementType, ParseDataType | 200 |
| `snowflake/parser/datatypes_test.go` | 16 test categories | 300 |
| `snowflake/ast/datatypes_test.go` | TypeKind.String + TypeName.Tag tests | 50 |

Total: ~550 LOC across 3 new + 6 modified files.

---

## Task 1: F2 enhancement — add timestamp variant keywords

**Files:**
- Modify: `snowflake/parser/tokens.go`
- Modify: `snowflake/parser/keywords.go`

- [ ] **Step 1: Confirm worktree state**

Run: `pwd && git rev-parse --abbrev-ref HEAD`
Expected:
```
/Users/h3n4l/OpenSource/omni/.worktrees/datatypes
feat/snowflake/datatypes
```

- [ ] **Step 2: Add 3 keyword constants to tokens.go**

In `snowflake/parser/tokens.go`, find the line `kwTIMESTAMP` (around line 837) and add the three new constants immediately after it:

```go
	kwTIMESTAMP_LTZ
	kwTIMESTAMP_NTZ
	kwTIMESTAMP_TZ
```

These must appear AFTER `kwTIMESTAMP` in the iota block. The exact position among the alphabetically-sorted constants matters for the iota values to be unique. Insert them right after `kwTIMESTAMP` (before `kwTIMESTAMP_ABORT_ON_ERROR` or whatever follows).

- [ ] **Step 3: Add 6 keyword map entries to keywords.go**

In `snowflake/parser/keywords.go`, find the `"timestamp"` entry in the `keywordMap` and add the 6 new entries after it, maintaining alphabetical order:

```go
	"timestamp_ltz":                                kwTIMESTAMP_LTZ,
	"timestamp_ntz":                                kwTIMESTAMP_NTZ,
	"timestamp_tz":                                 kwTIMESTAMP_TZ,
	"timestampltz":                                 kwTIMESTAMP_LTZ,
	"timestampntz":                                 kwTIMESTAMP_NTZ,
	"timestamptz":                                  kwTIMESTAMP_TZ,
```

The no-underscore aliases (`timestampltz`, etc.) support Snowflake's legacy form where the underscore is optional.

- [ ] **Step 4: Verify build and existing tests**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go test ./snowflake/parser/...`
Expected: all existing F2/F3/F4/T1.1 tests still pass.

---

## Task 2: Add TypeKind enum + TypeName struct to AST

**Files:**
- Modify: `snowflake/ast/parsenodes.go`
- Modify: `snowflake/ast/nodetags.go`
- Modify: `snowflake/ast/loc.go`
- Regenerate: `snowflake/ast/walk_generated.go`

- [ ] **Step 1: Append TypeKind enum and TypeName struct to parsenodes.go**

Use Edit to append the following after the `ObjectName.Matches` method at the end of `snowflake/ast/parsenodes.go`:

```go

// ---------------------------------------------------------------------------
// Data type types
// ---------------------------------------------------------------------------

// TypeKind classifies Snowflake data types into categories for fast
// switch dispatch by downstream consumers. The Name field on TypeName
// carries the exact source text for round-tripping; Kind carries the
// category for semantic analysis.
type TypeKind int

const (
	TypeUnknown      TypeKind = iota
	TypeInt                    // INT, INTEGER, SMALLINT, TINYINT, BYTEINT, BIGINT
	TypeNumber                 // NUMBER, NUMERIC, DECIMAL — may have (precision, scale)
	TypeFloat                  // FLOAT, FLOAT4, FLOAT8, DOUBLE, DOUBLE PRECISION, REAL
	TypeBoolean                // BOOLEAN
	TypeDate                   // DATE
	TypeDateTime               // DATETIME — may have (precision)
	TypeTime                   // TIME — may have (precision)
	TypeTimestamp               // TIMESTAMP — may have (precision)
	TypeTimestampLTZ            // TIMESTAMP_LTZ — may have (precision)
	TypeTimestampNTZ            // TIMESTAMP_NTZ — may have (precision)
	TypeTimestampTZ             // TIMESTAMP_TZ — may have (precision)
	TypeChar                   // CHAR, NCHAR, CHARACTER — may have (length)
	TypeVarchar                // VARCHAR, CHAR VARYING, NCHAR VARYING, NVARCHAR, NVARCHAR2, STRING, TEXT
	TypeBinary                 // BINARY — may have (length)
	TypeVarbinary              // VARBINARY — may have (length)
	TypeVariant                // VARIANT
	TypeObject                 // OBJECT
	TypeArray                  // ARRAY — may have ElementType
	TypeGeography              // GEOGRAPHY
	TypeGeometry               // GEOMETRY
	TypeVector                 // VECTOR — has ElementType + VectorDim
)

// String returns the human-readable name of the TypeKind.
func (k TypeKind) String() string {
	switch k {
	case TypeUnknown:
		return "Unknown"
	case TypeInt:
		return "Int"
	case TypeNumber:
		return "Number"
	case TypeFloat:
		return "Float"
	case TypeBoolean:
		return "Boolean"
	case TypeDate:
		return "Date"
	case TypeDateTime:
		return "DateTime"
	case TypeTime:
		return "Time"
	case TypeTimestamp:
		return "Timestamp"
	case TypeTimestampLTZ:
		return "TimestampLTZ"
	case TypeTimestampNTZ:
		return "TimestampNTZ"
	case TypeTimestampTZ:
		return "TimestampTZ"
	case TypeChar:
		return "Char"
	case TypeVarchar:
		return "Varchar"
	case TypeBinary:
		return "Binary"
	case TypeVarbinary:
		return "Varbinary"
	case TypeVariant:
		return "Variant"
	case TypeObject:
		return "Object"
	case TypeArray:
		return "Array"
	case TypeGeography:
		return "Geography"
	case TypeGeometry:
		return "Geometry"
	case TypeVector:
		return "Vector"
	default:
		return "Unknown"
	}
}

// TypeName represents a Snowflake data type as it appears in SQL source.
//
// Examples:
//
//	INT                  → Kind=TypeInt, Name="INT", Params=nil
//	NUMBER(38, 0)        → Kind=TypeNumber, Name="NUMBER", Params=[38, 0]
//	VARCHAR(100)         → Kind=TypeVarchar, Name="VARCHAR", Params=[100]
//	TIMESTAMP_LTZ(9)     → Kind=TypeTimestampLTZ, Name="TIMESTAMP_LTZ", Params=[9]
//	DOUBLE PRECISION     → Kind=TypeFloat, Name="DOUBLE PRECISION", Params=nil
//	ARRAY(VARCHAR)       → Kind=TypeArray, Name="ARRAY", ElementType=&TypeName{...}
//	VECTOR(INT, 256)     → Kind=TypeVector, Name="VECTOR", ElementType=&TypeName{...}, VectorDim=256
//
// TypeName is a Node. The walker descends into ElementType when non-nil.
type TypeName struct {
	Kind        TypeKind  // classified type category
	Name        string    // source text of the type name for round-tripping
	Params      []int     // numeric type parameters; nil if absent
	ElementType *TypeName // element type for ARRAY and VECTOR; nil otherwise
	VectorDim   int       // dimension for VECTOR(type, dim); -1 if not VECTOR
	Loc         Loc
}

// Tag implements Node.
func (n *TypeName) Tag() NodeTag { return T_TypeName }

// Compile-time assertion that *TypeName satisfies Node.
var _ Node = (*TypeName)(nil)
```

- [ ] **Step 2: Add T_TypeName to nodetags.go**

In `snowflake/ast/nodetags.go`, add the constant after `T_ObjectName`:

```go
	// T_TypeName is the tag for *TypeName, a data type reference.
	T_TypeName
```

And add a case to the `String()` method:

```go
	case T_TypeName:
		return "TypeName"
```

- [ ] **Step 3: Add *TypeName case to NodeLoc in loc.go**

In `snowflake/ast/loc.go`, add a case to the `NodeLoc` switch before the `default`:

```go
	case *TypeName:
		return v.Loc
```

- [ ] **Step 4: Regenerate walk_generated.go**

Run: `go generate ./snowflake/ast/...`
Expected output:
```
Generated walk_generated.go: 2 cases, 2 child fields
```

**Important:** Unlike T1.1 (which produced "1 cases, 1 child fields"), T1.2 adds a REAL case. The `TypeName` struct has an `ElementType *TypeName` field — a pointer-to-known-struct that the walker generator recognizes as a child node. The generated case should look like:

```go
case *TypeName:
    if n.ElementType != nil {
        Walk(v, n.ElementType)
    }
```

Verify by reading the generated file:

Run: `cat snowflake/ast/walk_generated.go`
Expected: two cases — `*File` (walkNodes on Stmts) and `*TypeName` (Walk on ElementType).

- [ ] **Step 5: Verify build and tests**

Run: `go build ./snowflake/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/...`
Expected: no output, exit 0.

Run: `go test ./snowflake/ast/...`
Expected: all existing F1 + T1.1 AST tests pass.

---

## Task 3: Write AST unit tests for TypeKind and TypeName

**Files:**
- Create: `snowflake/ast/datatypes_test.go`

- [ ] **Step 1: Write the test file**

Create `snowflake/ast/datatypes_test.go`:

```go
package ast

import "testing"

func TestTypeKind_String(t *testing.T) {
	cases := []struct {
		kind TypeKind
		want string
	}{
		{TypeUnknown, "Unknown"},
		{TypeInt, "Int"},
		{TypeNumber, "Number"},
		{TypeFloat, "Float"},
		{TypeBoolean, "Boolean"},
		{TypeDate, "Date"},
		{TypeDateTime, "DateTime"},
		{TypeTime, "Time"},
		{TypeTimestamp, "Timestamp"},
		{TypeTimestampLTZ, "TimestampLTZ"},
		{TypeTimestampNTZ, "TimestampNTZ"},
		{TypeTimestampTZ, "TimestampTZ"},
		{TypeChar, "Char"},
		{TypeVarchar, "Varchar"},
		{TypeBinary, "Binary"},
		{TypeVarbinary, "Varbinary"},
		{TypeVariant, "Variant"},
		{TypeObject, "Object"},
		{TypeArray, "Array"},
		{TypeGeography, "Geography"},
		{TypeGeometry, "Geometry"},
		{TypeVector, "Vector"},
		{TypeKind(999), "Unknown"},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			if got := c.kind.String(); got != c.want {
				t.Errorf("TypeKind(%d).String() = %q, want %q", c.kind, got, c.want)
			}
		})
	}
}

func TestTypeName_Tag(t *testing.T) {
	var n TypeName
	if (&n).Tag() != T_TypeName {
		t.Errorf("Tag() = %v, want T_TypeName", (&n).Tag())
	}
}

func TestTypeName_WalkerVisitsElementType(t *testing.T) {
	// Verify the walker descends into ElementType.
	inner := &TypeName{Kind: TypeVarchar, Name: "VARCHAR"}
	outer := &TypeName{Kind: TypeArray, Name: "ARRAY", ElementType: inner, VectorDim: -1}

	var visited []NodeTag
	Inspect(outer, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	// Should visit outer (TypeName) then inner (TypeName).
	if len(visited) != 2 {
		t.Fatalf("visited %d nodes, want 2: %+v", len(visited), visited)
	}
	if visited[0] != T_TypeName || visited[1] != T_TypeName {
		t.Errorf("visited = %+v, want [T_TypeName, T_TypeName]", visited)
	}
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./snowflake/ast/...`
Expected: all tests pass.

---

## Task 4: Write parseDataType parser helper

**Files:**
- Create: `snowflake/parser/datatypes.go`

- [ ] **Step 1: Write the full parser helper file**

Create `snowflake/parser/datatypes.go`:

```go
package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// parseDataType parses a Snowflake data type. Called from type positions
// like CAST(x AS type), CREATE TABLE column definitions, function
// signatures, etc.
//
// Handles all forms from the legacy SnowflakeParser.g4 data_type rule.
func (p *Parser) parseDataType() (*ast.TypeName, error) {
	tok := p.cur

	switch tok.Type {
	// Integer types — no parameters.
	case kwINT, kwINTEGER, kwSMALLINT, kwTINYINT, kwBYTEINT, kwBIGINT:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeInt, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// Numeric types — optional (precision [, scale]).
	case kwNUMBER, kwNUMERIC, kwDECIMAL:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeNumber, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Float types — no parameters. DOUBLE has PRECISION lookahead.
	case kwFLOAT, kwFLOAT4, kwFLOAT8, kwREAL:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeFloat, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwDOUBLE:
		p.advance()
		// Lookahead for PRECISION (comes as tokIdent since it's not a keyword).
		if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "PRECISION" {
			precTok := p.advance()
			return &ast.TypeName{
				Kind:      ast.TypeFloat,
				Name:      "DOUBLE PRECISION",
				VectorDim: -1,
				Loc:       ast.Loc{Start: tok.Loc.Start, End: precTok.Loc.End},
			}, nil
		}
		return &ast.TypeName{Kind: ast.TypeFloat, Name: "DOUBLE", VectorDim: -1, Loc: tok.Loc}, nil

	// Boolean.
	case kwBOOLEAN:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeBoolean, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// Date.
	case kwDATE:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeDate, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// DateTime — optional (precision).
	case kwDATETIME:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeDateTime, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Time — optional (precision).
	case kwTIME:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeTime, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Timestamp variants — optional (precision).
	case kwTIMESTAMP:
		return p.parseTimestampType(ast.TypeTimestamp, tok)
	case kwTIMESTAMP_LTZ:
		return p.parseTimestampType(ast.TypeTimestampLTZ, tok)
	case kwTIMESTAMP_NTZ:
		return p.parseTimestampType(ast.TypeTimestampNTZ, tok)
	case kwTIMESTAMP_TZ:
		return p.parseTimestampType(ast.TypeTimestampTZ, tok)

	// Char types — optional (length). Lookahead for VARYING → TypeVarchar.
	case kwCHAR, kwNCHAR:
		p.advance()
		// Lookahead for VARYING (comes as tokIdent).
		if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "VARYING" {
			varyTok := p.advance()
			name := tok.Str + " VARYING"
			params, endLoc, err := p.parseOptionalTypeParams()
			if err != nil {
				return nil, err
			}
			loc := ast.Loc{Start: tok.Loc.Start, End: varyTok.Loc.End}
			if endLoc.End > loc.End {
				loc.End = endLoc.End
			}
			return &ast.TypeName{Kind: ast.TypeVarchar, Name: name, Params: params, VectorDim: -1, Loc: loc}, nil
		}
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeChar, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	case kwCHARACTER:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeChar, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Varchar-family types — optional (length).
	case kwVARCHAR, kwNVARCHAR, kwNVARCHAR2, kwSTRING, kwTEXT:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeVarchar, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Binary types — optional (length).
	case kwBINARY:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeBinary, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	case kwVARBINARY:
		p.advance()
		params, endLoc, err := p.parseOptionalTypeParams()
		if err != nil {
			return nil, err
		}
		loc := tok.Loc
		if endLoc.End > loc.End {
			loc.End = endLoc.End
		}
		return &ast.TypeName{Kind: ast.TypeVarbinary, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil

	// Semi-structured types.
	case kwVARIANT:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeVariant, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwOBJECT:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeObject, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwARRAY:
		p.advance()
		if p.cur.Type == '(' {
			openTok := p.advance() // consume (
			_ = openTok
			elem, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			closeTok, err := p.expect(')')
			if err != nil {
				return nil, err
			}
			return &ast.TypeName{
				Kind:        ast.TypeArray,
				Name:        "ARRAY",
				ElementType: elem,
				VectorDim:   -1,
				Loc:         ast.Loc{Start: tok.Loc.Start, End: closeTok.Loc.End},
			}, nil
		}
		return &ast.TypeName{Kind: ast.TypeArray, Name: "ARRAY", VectorDim: -1, Loc: tok.Loc}, nil

	// Geospatial types.
	case kwGEOGRAPHY:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeGeography, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	case kwGEOMETRY:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeGeometry, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil

	// VECTOR(element_type, dimensions).
	case kwVECTOR:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		elemType, err := p.parseVectorElementType()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		dimTok := p.cur
		if dimTok.Type != tokInt {
			return nil, &ParseError{Loc: dimTok.Loc, Msg: "expected integer dimension for VECTOR"}
		}
		p.advance()
		closeTok, err := p.expect(')')
		if err != nil {
			return nil, err
		}
		return &ast.TypeName{
			Kind:        ast.TypeVector,
			Name:        "VECTOR",
			ElementType: elemType,
			VectorDim:   int(dimTok.Ival),
			Loc:         ast.Loc{Start: tok.Loc.Start, End: closeTok.Loc.End},
		}, nil
	}

	return nil, &ParseError{Loc: tok.Loc, Msg: "expected data type"}
}

// parseTimestampType handles TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ,
// TIMESTAMP_TZ — all with optional (precision).
func (p *Parser) parseTimestampType(kind ast.TypeKind, tok Token) (*ast.TypeName, error) {
	p.advance()
	params, endLoc, err := p.parseOptionalTypeParams()
	if err != nil {
		return nil, err
	}
	loc := tok.Loc
	if endLoc.End > loc.End {
		loc.End = endLoc.End
	}
	return &ast.TypeName{Kind: kind, Name: tok.Str, Params: params, VectorDim: -1, Loc: loc}, nil
}

// parseOptionalTypeParams parses an optional parenthesized parameter list
// after a type keyword: (n) or (n, m). Returns nil if no opening paren.
// The returned ast.Loc is the span of the parenthesized expression (or
// NoLoc if no params).
func (p *Parser) parseOptionalTypeParams() ([]int, ast.Loc, error) {
	if p.cur.Type != '(' {
		return nil, ast.NoLoc(), nil
	}
	openTok := p.advance() // consume (

	if p.cur.Type != tokInt {
		return nil, ast.NoLoc(), &ParseError{Loc: p.cur.Loc, Msg: "expected integer type parameter"}
	}
	first := int(p.cur.Ival)
	p.advance()

	params := []int{first}

	if p.cur.Type == ',' {
		p.advance() // consume ,
		if p.cur.Type != tokInt {
			return nil, ast.NoLoc(), &ParseError{Loc: p.cur.Loc, Msg: "expected integer type parameter"}
		}
		second := int(p.cur.Ival)
		p.advance()
		params = append(params, second)
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	_ = openTok
	return params, ast.Loc{Start: openTok.Loc.Start, End: closeTok.Loc.End}, nil
}

// parseVectorElementType parses the element type for VECTOR, which is
// restricted to INT | INTEGER | FLOAT | FLOAT4 | FLOAT8 per the grammar.
func (p *Parser) parseVectorElementType() (*ast.TypeName, error) {
	tok := p.cur
	switch tok.Type {
	case kwINT, kwINTEGER:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeInt, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil
	case kwFLOAT, kwFLOAT4, kwFLOAT8:
		p.advance()
		return &ast.TypeName{Kind: ast.TypeFloat, Name: tok.Str, VectorDim: -1, Loc: tok.Loc}, nil
	}
	return nil, &ParseError{Loc: tok.Loc, Msg: "expected VECTOR element type (INT, INTEGER, FLOAT, FLOAT4, or FLOAT8)"}
}

// ParseDataType parses a data type from a standalone string. Useful for
// tests and catalog integration.
func ParseDataType(input string) (*ast.TypeName, []ParseError) {
	p := &Parser{
		lexer: NewLexer(input),
		input: input,
	}
	p.advance()

	dt, err := p.parseDataType()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}

	if p.cur.Type != tokEOF {
		return dt, []ParseError{{
			Loc: p.cur.Loc,
			Msg: "unexpected token after data type",
		}}
	}

	return dt, nil
}

```

- [ ] **Step 2: Verify build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 5: Write parser data type tests

**Files:**
- Create: `snowflake/parser/datatypes_test.go`

- [ ] **Step 1: Write the full test file**

Create `snowflake/parser/datatypes_test.go`:

```go
package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseDataType constructs a Parser from input and calls parseDataType.
func testParseDataType(input string) (*ast.TypeName, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseDataType()
}

func TestParseDataType_SimpleInt(t *testing.T) {
	for _, kw := range []string{"INT", "INTEGER", "SMALLINT", "TINYINT", "BYTEINT", "BIGINT"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeInt {
				t.Errorf("Kind = %v, want TypeInt", dt.Kind)
			}
			if strings.ToUpper(dt.Name) != kw {
				t.Errorf("Name = %q, want %q", dt.Name, kw)
			}
			if dt.Params != nil {
				t.Errorf("Params = %v, want nil", dt.Params)
			}
		})
	}
}

func TestParseDataType_NumberWithParams(t *testing.T) {
	cases := []struct {
		input      string
		wantName   string
		wantParams []int
	}{
		{"NUMBER(38, 0)", "NUMBER", []int{38, 0}},
		{"NUMERIC(10)", "NUMERIC", []int{10}},
		{"DECIMAL(18, 2)", "DECIMAL", []int{18, 2}},
		{"NUMBER", "NUMBER", nil},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			dt, err := testParseDataType(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeNumber {
				t.Errorf("Kind = %v, want TypeNumber", dt.Kind)
			}
			if len(dt.Params) != len(c.wantParams) {
				t.Fatalf("Params = %v, want %v", dt.Params, c.wantParams)
			}
			for i, p := range dt.Params {
				if p != c.wantParams[i] {
					t.Errorf("Params[%d] = %d, want %d", i, p, c.wantParams[i])
				}
			}
		})
	}
}

func TestParseDataType_FloatAliases(t *testing.T) {
	for _, kw := range []string{"FLOAT", "FLOAT4", "FLOAT8", "REAL"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeFloat {
				t.Errorf("Kind = %v, want TypeFloat", dt.Kind)
			}
		})
	}
}

func TestParseDataType_DoublePrecision(t *testing.T) {
	dt, err := testParseDataType("DOUBLE PRECISION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeFloat {
		t.Errorf("Kind = %v, want TypeFloat", dt.Kind)
	}
	if dt.Name != "DOUBLE PRECISION" {
		t.Errorf("Name = %q, want 'DOUBLE PRECISION'", dt.Name)
	}
}

func TestParseDataType_DoubleAlone(t *testing.T) {
	dt, err := testParseDataType("DOUBLE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeFloat || dt.Name != "DOUBLE" {
		t.Errorf("got Kind=%v Name=%q, want TypeFloat DOUBLE", dt.Kind, dt.Name)
	}
}

func TestParseDataType_CharVarying(t *testing.T) {
	for _, input := range []string{"CHAR VARYING", "NCHAR VARYING"} {
		t.Run(input, func(t *testing.T) {
			dt, err := testParseDataType(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeVarchar {
				t.Errorf("Kind = %v, want TypeVarchar", dt.Kind)
			}
			if dt.Name != input {
				t.Errorf("Name = %q, want %q", dt.Name, input)
			}
		})
	}
}

func TestParseDataType_CharVaryingWithLength(t *testing.T) {
	dt, err := testParseDataType("CHAR VARYING(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeVarchar || dt.Name != "CHAR VARYING" {
		t.Errorf("got Kind=%v Name=%q", dt.Kind, dt.Name)
	}
	if len(dt.Params) != 1 || dt.Params[0] != 100 {
		t.Errorf("Params = %v, want [100]", dt.Params)
	}
}

func TestParseDataType_StringTypes(t *testing.T) {
	for _, kw := range []string{"VARCHAR", "NVARCHAR", "NVARCHAR2", "STRING", "TEXT"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeVarchar {
				t.Errorf("Kind = %v, want TypeVarchar", dt.Kind)
			}
		})
	}
}

func TestParseDataType_VarcharWithLength(t *testing.T) {
	dt, err := testParseDataType("VARCHAR(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeVarchar {
		t.Errorf("Kind = %v, want TypeVarchar", dt.Kind)
	}
	if len(dt.Params) != 1 || dt.Params[0] != 100 {
		t.Errorf("Params = %v, want [100]", dt.Params)
	}
}

func TestParseDataType_CharTypes(t *testing.T) {
	for _, kw := range []string{"CHAR", "NCHAR", "CHARACTER"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeChar {
				t.Errorf("Kind = %v, want TypeChar", dt.Kind)
			}
		})
	}
}

func TestParseDataType_CharWithLength(t *testing.T) {
	dt, err := testParseDataType("CHAR(10)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeChar || len(dt.Params) != 1 || dt.Params[0] != 10 {
		t.Errorf("got Kind=%v Params=%v, want TypeChar [10]", dt.Kind, dt.Params)
	}
}

func TestParseDataType_Binary(t *testing.T) {
	for _, tc := range []struct {
		input string
		kind  ast.TypeKind
	}{
		{"BINARY", ast.TypeBinary},
		{"BINARY(16)", ast.TypeBinary},
		{"VARBINARY", ast.TypeVarbinary},
		{"VARBINARY(256)", ast.TypeVarbinary},
	} {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseDataType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != tc.kind {
				t.Errorf("Kind = %v, want %v", dt.Kind, tc.kind)
			}
		})
	}
}

func TestParseDataType_TimestampVariants(t *testing.T) {
	cases := []struct {
		input    string
		wantKind ast.TypeKind
	}{
		{"TIMESTAMP", ast.TypeTimestamp},
		{"TIMESTAMP(9)", ast.TypeTimestamp},
		{"TIMESTAMP_LTZ", ast.TypeTimestampLTZ},
		{"TIMESTAMP_LTZ(9)", ast.TypeTimestampLTZ},
		{"TIMESTAMPLTZ", ast.TypeTimestampLTZ},
		{"TIMESTAMP_NTZ", ast.TypeTimestampNTZ},
		{"TIMESTAMP_NTZ(6)", ast.TypeTimestampNTZ},
		{"TIMESTAMPNTZ", ast.TypeTimestampNTZ},
		{"TIMESTAMP_TZ", ast.TypeTimestampTZ},
		{"TIMESTAMP_TZ(3)", ast.TypeTimestampTZ},
		{"TIMESTAMPTZ", ast.TypeTimestampTZ},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			dt, err := testParseDataType(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != c.wantKind {
				t.Errorf("Kind = %v, want %v", dt.Kind, c.wantKind)
			}
		})
	}
}

func TestParseDataType_TimeWithPrecision(t *testing.T) {
	dt, err := testParseDataType("TIME(9)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeTime || len(dt.Params) != 1 || dt.Params[0] != 9 {
		t.Errorf("got Kind=%v Params=%v, want TypeTime [9]", dt.Kind, dt.Params)
	}
}

func TestParseDataType_DateTime(t *testing.T) {
	dt, err := testParseDataType("DATETIME")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeDateTime {
		t.Errorf("Kind = %v, want TypeDateTime", dt.Kind)
	}
}

func TestParseDataType_Boolean(t *testing.T) {
	dt, err := testParseDataType("BOOLEAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeBoolean {
		t.Errorf("Kind = %v, want TypeBoolean", dt.Kind)
	}
}

func TestParseDataType_SemiStructured(t *testing.T) {
	for _, tc := range []struct {
		input string
		kind  ast.TypeKind
	}{
		{"VARIANT", ast.TypeVariant},
		{"OBJECT", ast.TypeObject},
		{"GEOGRAPHY", ast.TypeGeography},
		{"GEOMETRY", ast.TypeGeometry},
	} {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseDataType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != tc.kind {
				t.Errorf("Kind = %v, want %v", dt.Kind, tc.kind)
			}
		})
	}
}

func TestParseDataType_ArrayUntyped(t *testing.T) {
	dt, err := testParseDataType("ARRAY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeArray {
		t.Errorf("Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType != nil {
		t.Errorf("ElementType = %+v, want nil", dt.ElementType)
	}
}

func TestParseDataType_ArrayTyped(t *testing.T) {
	dt, err := testParseDataType("ARRAY(VARCHAR)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeArray {
		t.Errorf("Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType == nil {
		t.Fatal("ElementType is nil, want *TypeName")
	}
	if dt.ElementType.Kind != ast.TypeVarchar {
		t.Errorf("ElementType.Kind = %v, want TypeVarchar", dt.ElementType.Kind)
	}
}

func TestParseDataType_ArrayNested(t *testing.T) {
	dt, err := testParseDataType("ARRAY(NUMBER(10, 2))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeArray || dt.ElementType == nil {
		t.Fatalf("got Kind=%v ElementType=%v", dt.Kind, dt.ElementType)
	}
	if dt.ElementType.Kind != ast.TypeNumber {
		t.Errorf("ElementType.Kind = %v, want TypeNumber", dt.ElementType.Kind)
	}
	if len(dt.ElementType.Params) != 2 || dt.ElementType.Params[0] != 10 || dt.ElementType.Params[1] != 2 {
		t.Errorf("ElementType.Params = %v, want [10, 2]", dt.ElementType.Params)
	}
}

func TestParseDataType_Vector(t *testing.T) {
	cases := []struct {
		input    string
		elemKind ast.TypeKind
		dim      int
	}{
		{"VECTOR(INT, 256)", ast.TypeInt, 256},
		{"VECTOR(FLOAT, 768)", ast.TypeFloat, 768},
		{"VECTOR(FLOAT4, 128)", ast.TypeFloat, 128},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			dt, err := testParseDataType(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeVector {
				t.Errorf("Kind = %v, want TypeVector", dt.Kind)
			}
			if dt.ElementType == nil || dt.ElementType.Kind != c.elemKind {
				t.Errorf("ElementType = %+v, want Kind=%v", dt.ElementType, c.elemKind)
			}
			if dt.VectorDim != c.dim {
				t.Errorf("VectorDim = %d, want %d", dt.VectorDim, c.dim)
			}
		})
	}
}

func TestParseDataType_Errors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"unknown keyword", "FOOBAR"},
		{"missing close paren", "NUMBER(38"},
		{"vector missing dim", "VECTOR(INT)"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := testParseDataType(c.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestParseDataType_Freestanding(t *testing.T) {
	dt, errs := ParseDataType("VARCHAR(100)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if dt.Kind != ast.TypeVarchar || len(dt.Params) != 1 || dt.Params[0] != 100 {
		t.Errorf("got Kind=%v Params=%v", dt.Kind, dt.Params)
	}
}

func TestParseDataType_FreestandingTrailing(t *testing.T) {
	dt, errs := ParseDataType("INT foo")
	if dt == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if len(errs) == 0 || !strings.Contains(errs[0].Msg, "unexpected token") {
		t.Errorf("expected 'unexpected token' error, got %+v", errs)
	}
}

func TestParseDataType_LocSpanning(t *testing.T) {
	// "DOUBLE PRECISION" — Loc should span both words.
	dt, err := testParseDataType("DOUBLE PRECISION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "DOUBLE" = 6 chars at 0..5, " " at 6, "PRECISION" at 7..15.
	// Loc should be {0, 16}.
	if dt.Loc.Start != 0 || dt.Loc.End != 16 {
		t.Errorf("Loc = %+v, want {0, 16}", dt.Loc)
	}
}

func TestParseDataType_LocWithParams(t *testing.T) {
	// "VARCHAR(100)" — Loc should span from V to closing ).
	dt, err := testParseDataType("VARCHAR(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Loc.Start != 0 || dt.Loc.End != 12 {
		t.Errorf("Loc = %+v, want {0, 12}", dt.Loc)
	}
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./snowflake/...`
Expected: all tests pass (F1 + T1.1 + F2/F3/F4 + T1.2).

- [ ] **Step 3: Verbose check**

Run: `go test -v ./snowflake/parser/... -run TestParseDataType 2>&1 | grep -E "^(=== RUN|--- (PASS|FAIL))" | tail -60`
Expected: every subtest reported PASS.

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
Expected: no output. If any file listed, `gofmt -w snowflake/` and re-check.

- [ ] **Step 4: Test**

Run: `go test ./snowflake/...`
Expected: both ast and parser pass.

- [ ] **Step 5: Walker generation round-trip**

Run: `cp snowflake/ast/walk_generated.go /tmp/wg_check.go && go generate ./snowflake/ast/... && diff /tmp/wg_check.go snowflake/ast/walk_generated.go && echo "BYTE_IDENTICAL" && rm /tmp/wg_check.go`
Expected: `Generated walk_generated.go: 2 cases, 2 child fields` and `BYTE_IDENTICAL`.

- [ ] **Step 6: List files**

Run: `git status snowflake/`
Expected:
```
modified:   snowflake/ast/loc.go
modified:   snowflake/ast/nodetags.go
modified:   snowflake/ast/parsenodes.go
modified:   snowflake/ast/walk_generated.go
modified:   snowflake/parser/keywords.go
modified:   snowflake/parser/tokens.go

Untracked:
	snowflake/ast/datatypes_test.go
	snowflake/parser/datatypes.go
	snowflake/parser/datatypes_test.go
```

Note `walk_generated.go` IS modified this time (unlike T1.1) because it gained the `*TypeName` case.

- [ ] **Step 7: STOP and present for review**

Do NOT commit.

---

## Spec Coverage Checklist

| Spec section | Covered by |
|---|---|
| TypeKind enum (22 values) | Task 2 |
| TypeName struct (Kind, Name, Params, ElementType, VectorDim, Loc) | Task 2 |
| T_TypeName constant + String() | Task 2 |
| NodeLoc *TypeName case | Task 2 |
| Walker with ElementType child | Task 2 (regen) |
| F2 timestamp variant keywords (3 constants + 6 map entries) | Task 1 |
| parseDataType dispatch (full keyword table) | Task 4 |
| Multi-word types (DOUBLE PRECISION, CHAR/NCHAR VARYING) | Task 4 |
| Optional type params (precision, scale, length) | Task 4 |
| ARRAY with optional element type | Task 4 |
| VECTOR with element type + dimension | Task 4 |
| ParseDataType freestanding helper | Task 4 |
| TypeKind.String tests | Task 3 |
| Walker traversal of ElementType test | Task 3 |
| Parser tests (16 categories) | Task 5 |
| Acceptance criteria | Task 6 |
