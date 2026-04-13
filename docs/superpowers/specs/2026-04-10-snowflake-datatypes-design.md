# Snowflake Data Types (T1.2) — Design

**DAG node:** T1.2 — data types
**Migration:** `docs/migration/snowflake/dag.md`
**Branch:** `feat/snowflake/datatypes`
**Status:** Approved, ready for plan + implementation.
**Depends on:** T1.1 (identifiers, merged via PR #22).
**Unblocks:** T1.3 (expressions — needs type parsing for CAST/TRY_CAST, :: operator).

## Purpose

T1.2 adds data type parsing — the building blocks for CREATE TABLE column definitions, CAST/TRY_CAST expressions, function signatures, and type annotations throughout Snowflake SQL. It ships:

1. **`TypeName` Node struct** in `snowflake/ast` — a single generic struct with a `TypeKind` enum for fast dispatch, source text `Name` for round-tripping, numeric `Params`, recursive `ElementType`, and VECTOR-specific `VectorDim`.
2. **`TypeKind` enum** with 22 values covering every Snowflake data type category.
3. **`parseDataType()` parser helper** in `snowflake/parser` that Tier 1.3+ nodes call when they encounter type positions.
4. **F2 keyword additions** — 3 new `kw*` constants and 6 keyword map entries for timestamp variant keywords missing from the extraction.

## AST Types

### `TypeKind` enum

Added to `snowflake/ast/parsenodes.go`:

```go
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
    TypeVarchar                // VARCHAR, CHAR VARYING, NCHAR VARYING, NVARCHAR, NVARCHAR2, STRING, TEXT — may have (length)
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
func (k TypeKind) String() string
```

### `TypeName` struct

```go
// TypeName represents a Snowflake data type as it appears in SQL source.
//
// Examples:
//   INT                      → Kind=TypeInt, Name="INT", Params=nil
//   NUMBER(38, 0)            → Kind=TypeNumber, Name="NUMBER", Params=[38, 0]
//   VARCHAR(100)             → Kind=TypeVarchar, Name="VARCHAR", Params=[100]
//   TIMESTAMP_LTZ(9)         → Kind=TypeTimestampLTZ, Name="TIMESTAMP_LTZ", Params=[9]
//   DOUBLE PRECISION         → Kind=TypeFloat, Name="DOUBLE PRECISION", Params=nil
//   ARRAY(VARCHAR)           → Kind=TypeArray, Name="ARRAY", ElementType=&TypeName{Kind:TypeVarchar}
//   VECTOR(INT, 256)         → Kind=TypeVector, Name="VECTOR", ElementType=&TypeName{Kind:TypeInt}, VectorDim=256
//
// TypeName is a Node. The walker descends into ElementType when non-nil.
type TypeName struct {
    Kind        TypeKind   // classified type category
    Name        string     // source text of the type name for round-tripping
    Params      []int      // numeric type parameters; nil if absent. [precision] or [precision, scale].
    ElementType *TypeName  // element type for ARRAY and VECTOR; nil otherwise
    VectorDim   int        // dimension for VECTOR(type, dim); -1 if not VECTOR
    Loc         Loc
}

// Tag implements Node.
func (n *TypeName) Tag() NodeTag { return T_TypeName }

// Compile-time assertion.
var _ Node = (*TypeName)(nil)
```

### NodeTag and walker

- `T_TypeName` added to `nodetags.go`
- `*TypeName` case added to `NodeLoc` in `loc.go`
- Walker regeneration: `walk_generated.go` gains a new case for `*TypeName` that walks `ElementType` when non-nil (since `ElementType *TypeName` is a pointer-to-known-struct)

## F2 Keyword Additions

The legacy grammar defines timestamp variants via non-standard ANTLR rules (`'TIMESTAMP' '_'? 'LTZ'`) that F2's grep extraction missed. T1.2 adds them:

### New constants in `tokens.go`

```go
kwTIMESTAMP_LTZ
kwTIMESTAMP_NTZ
kwTIMESTAMP_TZ
```

### New entries in `keywords.go` keywordMap

```go
"timestamp_ltz":  kwTIMESTAMP_LTZ,
"timestampltz":   kwTIMESTAMP_LTZ,   // Snowflake accepts both forms
"timestamp_ntz":  kwTIMESTAMP_NTZ,
"timestampntz":   kwTIMESTAMP_NTZ,
"timestamp_tz":   kwTIMESTAMP_TZ,
"timestamptz":    kwTIMESTAMP_TZ,
```

F2's lexer reads `TIMESTAMP_LTZ` as one identifier (since `_` is an ident-cont character), looks it up in keywordMap, and now finds a match. Both the with-underscore (`TIMESTAMP_LTZ`) and without-underscore (`TIMESTAMPLTZ`) forms are supported — two map entries per variant.

## Parser Helpers

New file `snowflake/parser/datatypes.go`:

```go
// parseDataType parses a Snowflake data type. Called from type positions
// like CAST(x AS type), CREATE TABLE column definitions, function
// signatures, etc.
//
// Handles all forms from the legacy SnowflakeParser.g4 data_type rule:
//   - Simple types: INT, BOOLEAN, DATE, VARIANT, GEOGRAPHY, GEOMETRY
//   - Parametric: NUMBER(38,0), VARCHAR(100), TIME(9)
//   - Multi-word (fused via parser lookahead): DOUBLE PRECISION, CHAR VARYING, NCHAR VARYING
//   - Timestamp variants: TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, TIMESTAMP_TZ
//   - Compound: ARRAY, ARRAY(VARCHAR), VECTOR(INT, 256)
func (p *Parser) parseDataType() (*ast.TypeName, error)
```

### Multi-word type handling

F2 does NOT fuse multi-word types (`DOUBLE PRECISION` comes as `kwDOUBLE` + `tokIdent("PRECISION")`). The parser handles this via lookahead:

```go
case kwDOUBLE:
    p.advance()
    if p.cur.Type == tokIdent && strings.ToUpper(p.cur.Str) == "PRECISION" {
        endTok := p.advance()
        return &ast.TypeName{Kind: ast.TypeFloat, Name: "DOUBLE PRECISION",
            Loc: ast.Loc{Start: tok.Loc.Start, End: endTok.Loc.End}}, nil
    }
    return &ast.TypeName{Kind: ast.TypeFloat, Name: "DOUBLE", Loc: tok.Loc}, nil
```

Same pattern for `kwCHAR` + `"VARYING"` → CHAR VARYING (TypeVarchar) and `kwNCHAR` + `"VARYING"` → NCHAR VARYING (TypeVarchar).

### Optional type parameters

```go
// parseOptionalTypeParams parses an optional parenthesized parameter list
// after a type keyword: (n) or (n, m). Returns nil if no opening paren.
func (p *Parser) parseOptionalTypeParams() ([]int, error)
```

### ARRAY element type

The legacy grammar uses `ARRAY '(' data_type ')'` (parentheses, NOT angle brackets). The parser:

```go
case kwARRAY:
    p.advance()
    if p.cur.Type == '(' {
        p.advance() // consume (
        elem, err := p.parseDataType()
        if err != nil { return nil, err }
        if _, err := p.expect(')'); err != nil { return nil, err }
        return &ast.TypeName{Kind: ast.TypeArray, Name: "ARRAY",
            ElementType: elem, VectorDim: -1, ...}, nil
    }
    return &ast.TypeName{Kind: ast.TypeArray, Name: "ARRAY", VectorDim: -1, ...}, nil
```

### VECTOR element type + dimension

Per the grammar's `vector_element_type` rule, the element type is restricted to INT/INTEGER/FLOAT/FLOAT4/FLOAT8:

```go
case kwVECTOR:
    p.advance()
    p.expect('(')
    elemType, err := p.parseVectorElementType() // INT|INTEGER|FLOAT|FLOAT4|FLOAT8
    p.expect(',')
    dim := p.parseIntLiteral()
    p.expect(')')
    return &ast.TypeName{Kind: ast.TypeVector, Name: "VECTOR",
        ElementType: elemType, VectorDim: dim, ...}, nil
```

### Freestanding helper

```go
// ParseDataType parses a data type from a standalone string. Useful for
// tests and catalog integration.
func ParseDataType(input string) (*ast.TypeName, []ParseError)
```

## Dispatch Table

The full keyword → TypeKind mapping:

| Keywords | TypeKind | Params | Notes |
|---|---|---|---|
| INT, INTEGER, SMALLINT, TINYINT, BYTEINT, BIGINT | TypeInt | none | |
| NUMBER, NUMERIC, DECIMAL | TypeNumber | (p) or (p, s) | |
| FLOAT, FLOAT4, FLOAT8, REAL | TypeFloat | none | |
| DOUBLE | TypeFloat | none | Lookahead for PRECISION |
| BOOLEAN | TypeBoolean | none | |
| DATE | TypeDate | none | |
| DATETIME | TypeDateTime | (p) | |
| TIME | TypeTime | (p) | |
| TIMESTAMP | TypeTimestamp | (p) | |
| TIMESTAMP_LTZ | TypeTimestampLTZ | (p) | NEW F2 keyword |
| TIMESTAMP_NTZ | TypeTimestampNTZ | (p) | NEW F2 keyword |
| TIMESTAMP_TZ | TypeTimestampTZ | (p) | NEW F2 keyword |
| CHAR, NCHAR, CHARACTER | TypeChar | (n) | Lookahead for VARYING → TypeVarchar |
| VARCHAR, NVARCHAR, NVARCHAR2, STRING, TEXT | TypeVarchar | (n) | |
| BINARY | TypeBinary | (n) | |
| VARBINARY | TypeVarbinary | (n) | |
| VARIANT | TypeVariant | none | |
| OBJECT | TypeObject | none | |
| ARRAY | TypeArray | none | Optional (element_type) |
| GEOGRAPHY | TypeGeography | none | |
| GEOMETRY | TypeGeometry | none | |
| VECTOR | TypeVector | (elem, dim) | ElementType + VectorDim |

## File Layout

### Modified

| File | Changes |
|------|---------|
| `snowflake/ast/parsenodes.go` | Append TypeKind enum + TypeName struct + methods |
| `snowflake/ast/nodetags.go` | Append T_TypeName + String() case |
| `snowflake/ast/loc.go` | Append *TypeName case to NodeLoc |
| `snowflake/ast/walk_generated.go` | Regenerated — gains *TypeName case with ElementType child walk |
| `snowflake/parser/tokens.go` | Append kwTIMESTAMP_LTZ/NTZ/TZ constants |
| `snowflake/parser/keywords.go` | Append 6 keywordMap entries for timestamp variants |

### Created

| File | Purpose | Approx LOC |
|------|---------|-----------|
| `snowflake/parser/datatypes.go` | parseDataType, parseOptionalTypeParams, parseVectorElementType, ParseDataType | 200 |
| `snowflake/parser/datatypes_test.go` | Table-driven tests for every type form | 300 |
| `snowflake/ast/datatypes_test.go` | TypeKind.String, TypeName methods | 50 |

**Estimated total: ~550 LOC** across 3 new + 6 modified files.

## Testing

### `snowflake/parser/datatypes_test.go`

Table-driven using a `testParseDataType(input) (*ast.TypeName, error)` helper:

1. **Simple types**: INT, INTEGER, SMALLINT, TINYINT, BYTEINT, BIGINT, BOOLEAN, DATE, VARIANT, OBJECT, GEOGRAPHY, GEOMETRY — each returns the correct Kind + Name, nil Params
2. **Number with params**: NUMBER(38,0), NUMERIC(10), DECIMAL(18,2) — correct Params slices
3. **Float aliases**: FLOAT, FLOAT4, FLOAT8, DOUBLE, REAL — all Kind=TypeFloat
4. **DOUBLE PRECISION**: two tokens fused → Kind=TypeFloat, Name="DOUBLE PRECISION"
5. **CHAR VARYING / NCHAR VARYING**: fused via lookahead → Kind=TypeVarchar
6. **String types**: VARCHAR, VARCHAR(100), NVARCHAR, NVARCHAR2, STRING, TEXT, CHARACTER, CHAR(10)
7. **Binary**: BINARY, BINARY(16), VARBINARY, VARBINARY(256)
8. **Time with precision**: TIME(9), TIME (no precision)
9. **Timestamp variants**: TIMESTAMP, TIMESTAMP(9), TIMESTAMP_LTZ, TIMESTAMP_LTZ(9), TIMESTAMP_NTZ, TIMESTAMP_NTZ(6), TIMESTAMP_TZ, TIMESTAMP_TZ(3), TIMESTAMPLTZ (no underscore form)
10. **DATETIME**: DATETIME, DATETIME(3)
11. **ARRAY untyped**: ARRAY → Kind=TypeArray, ElementType=nil
12. **ARRAY typed**: ARRAY(VARCHAR) → ElementType=TypeName{TypeVarchar}, ARRAY(NUMBER(10,2)) → nested
13. **VECTOR**: VECTOR(INT, 256), VECTOR(FLOAT, 768)
14. **Error cases**: missing closing paren, VECTOR missing dimension, unknown type keyword, empty input
15. **Freestanding**: ParseDataType("VARCHAR(100)") → correct result + no errors
16. **Loc spanning**: DOUBLE PRECISION Loc spans both tokens; ARRAY(VARCHAR) Loc spans from ARRAY to closing paren

### `snowflake/ast/datatypes_test.go`

- TypeKind.String() returns correct names for all 22 values + Unknown
- TypeName.Tag() == T_TypeName

Run scope: `go test ./snowflake/...`

## Out of Scope

| Feature | Where |
|---|---|
| Column definitions (DEFAULT, NOT NULL, constraints) | T2.2 (CREATE TABLE) |
| CAST / TRY_CAST expressions | T1.3 (expressions) |
| `::` cast operator | T1.3 |
| Type validation / compatibility checking | semantic package |
| Type resolution (finding catalog entries for types) | catalog package |

## Acceptance Criteria

1. `go build ./snowflake/...` succeeds.
2. `go vet ./snowflake/...` clean.
3. `gofmt -l snowflake/` clean.
4. `go test ./snowflake/...` passes — all prior tests still green + all T1.2 tests.
5. `go generate ./snowflake/ast/...` produces correct `walk_generated.go` with a new `*TypeName` case that walks `ElementType`.
6. Every data type form from the legacy `data_type` rule parses correctly.
7. Timestamp variant keywords (`TIMESTAMP_LTZ`, `TIMESTAMPLTZ`, etc.) are recognized by F2's lexer.
8. Multi-word types (`DOUBLE PRECISION`, `CHAR VARYING`, `NCHAR VARYING`) are fused correctly by the parser.
9. After merge, `docs/migration/snowflake/dag.md` T1.2 status is flipped to `done`.
