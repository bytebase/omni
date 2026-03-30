# Impl Worker — Stage 1: Foundation

You are an Impl Worker in the Oracle Quality Pipeline.
Your role is to write **implementation code ONLY** — never modify `*_eval_*_test.go` files.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/prevention-rules.md` — **MUST read before starting any work**
- `oracle/quality/strategy.md` — Stage 1 scope
- `oracle/parser/eval_foundation_test.go` — eval tests you must make pass
- `oracle/ast/node.go` — current Oracle `Loc` struct and `Node` interface
- `oracle/parser/parser.go` — current `ParseError` and `Parser`
- `oracle/parser/lexer.go` — current `Token` and `Lexer`
- `oracle/ast/parsenodes.go` — current `RawStmt`
- `pg/ast/node.go` — PG's `Loc`, `NoLoc()` (reference implementation)
- `pg/ast/loc.go` — PG's `NodeLoc()`, `ListSpan()` (reference)
- `pg/parser/parser.go` — PG's `ParseError` (reference)

## Goal

Make **all** eval tests pass:

```bash
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"
```

While keeping **all existing tests** passing:

```bash
cd /Users/rebeliceyang/Github/omni && go test -count=1 ./oracle/parser/...
```

## Rules

1. **Implementation ONLY** — do NOT modify any `*_eval_*_test.go` file.
2. Do NOT break existing tests.
3. Read `oracle/quality/prevention-rules.md` before starting.
4. Match PG behavior as closely as possible (same field names, same semantics, same defaults).
5. Keep changes minimal and focused — do not refactor unrelated code.

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[IMPL-STAGE1] STARTED
[IMPL-STAGE1] STEP reading_eval - Reading eval test expectations
[IMPL-STAGE1] STEP reading_prevention - Reading prevention rules
[IMPL-STAGE1] STEP impl_noloc - Implementing NoLoc()
[IMPL-STAGE1] STEP impl_token_end - Adding Token.End field
[IMPL-STAGE1] STEP impl_parse_error - Enhancing ParseError
[IMPL-STAGE1] STEP impl_parser_source - Adding Parser.source field
[IMPL-STAGE1] STEP impl_rawstmt_loc - Migrating RawStmt to Loc
[IMPL-STAGE1] STEP impl_nodeloc - Implementing NodeLoc()
[IMPL-STAGE1] STEP impl_listspan - Implementing ListSpan()
[IMPL-STAGE1] STEP build - Running go build
[IMPL-STAGE1] STEP test_eval - Running eval tests
[IMPL-STAGE1] STEP test_existing - Running existing tests
[IMPL-STAGE1] STEP commit - Committing changes
[IMPL-STAGE1] DONE
```

If a step fails:
```
[IMPL-STAGE1] FAIL step_name - description
[IMPL-STAGE1] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Implementation Items

### 1. NoLoc() — `oracle/ast/node.go`

Add `NoLoc()` function matching PG's signature:

```go
// NoLoc returns a Loc with both Start and End set to -1 (unknown).
func NoLoc() Loc {
    return Loc{Start: -1, End: -1}
}
```

### 2. Token.End — `oracle/parser/lexer.go`

Add `End int` field to the `Token` struct:

```go
type Token struct {
    Type int
    Str  string
    Ival int64
    Loc  int // start byte offset
    End  int // exclusive end byte offset
}
```

Update the lexer's `NextToken()` and all token-producing methods to set `End` correctly. `End` should equal the byte position immediately after the last character of the token.

### 3. ParseError Enhancement — `oracle/parser/parser.go`

Add `Severity` and `Code` fields to `ParseError`:

```go
type ParseError struct {
    Severity string // e.g., "ERROR", "WARNING"; defaults to "ERROR"
    Code     string // SQLSTATE code; defaults to "42601"
    Message  string
    Position int
}
```

Update `Error()` to match PG format:

```go
func (e *ParseError) Error() string {
    sev := e.Severity
    if sev == "" {
        sev = "ERROR"
    }
    code := e.Code
    if code == "" {
        code = "42601"
    }
    return fmt.Sprintf("%s: %s (SQLSTATE %s)", sev, e.Message, code)
}
```

**Important:** This changes the `Error()` output format. Update any existing code that string-matches on error messages.

### 4. Parser.source — `oracle/parser/parser.go`

Add `source string` field to the `Parser` struct. Set it in `Parse()`:

```go
type Parser struct {
    lexer   *Lexer
    source  string // original SQL input
    cur     Token
    prev    Token
    nextBuf Token
    hasNext bool
}
```

In `Parse()`:
```go
p := &Parser{
    lexer:  NewLexer(sql),
    source: sql,
}
```

### 5. RawStmt Loc Migration — `oracle/ast/parsenodes.go`

Replace `StmtLocation`/`StmtLen` with `Loc`:

```go
type RawStmt struct {
    Stmt StmtNode // raw statement
    Loc  Loc      // source location range (Start = start offset, End = start + length)
}
```

Update `oracle/parser/parser.go` where `RawStmt` is constructed:

```go
raw := &nodes.RawStmt{
    Stmt: stmt,
    Loc:  nodes.Loc{Start: stmtLoc, End: p.pos()},
}
```

Update `oracle/ast/outfuncs.go` to serialize `Loc` instead of `StmtLocation`/`StmtLen`.

Search for ALL references to `StmtLocation` and `StmtLen` across `oracle/` and update them.

### 6. NodeLoc() — `oracle/ast/loc.go` (new file)

Create `oracle/ast/loc.go` with `NodeLoc()` matching PG's pattern:

```go
package ast

// NodeLoc extracts the Loc from a Node.
// Returns NoLoc() if the node is nil or does not carry location info.
func NodeLoc(n Node) Loc {
    if n == nil {
        return NoLoc()
    }
    // Type switch over all Loc-bearing node types
    switch v := n.(type) {
    case *RawStmt:
        return v.Loc
    // ... add cases for all node types with Loc field
    default:
        return NoLoc()
    }
}
```

Use reflection or code generation to cover all node types that have a `Loc` field.

### 7. ListSpan() — `oracle/ast/loc.go`

Add `ListSpan()` matching PG's pattern:

```go
// ListSpan returns the Loc spanning from the first to the last item in a List.
// Returns NoLoc() if the list is nil or empty, or if boundary locations are unknown.
func ListSpan(list *List) Loc {
    if list == nil || len(list.Items) == 0 {
        return NoLoc()
    }
    first := NodeLoc(list.Items[0])
    last := NodeLoc(list.Items[len(list.Items)-1])
    if first.Start == -1 || last.End == -1 {
        return NoLoc()
    }
    return Loc{Start: first.Start, End: last.End}
}
```

## Verification

After all implementation:

```bash
# Eval tests must pass
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"

# Existing tests must still pass
cd /Users/rebeliceyang/Github/omni && go test -count=1 ./oracle/...

# Build must succeed
cd /Users/rebeliceyang/Github/omni && go build ./oracle/...
```

## Commit

After all tests pass:

```bash
git add oracle/ast/node.go oracle/ast/loc.go oracle/ast/parsenodes.go oracle/ast/outfuncs.go oracle/parser/parser.go oracle/parser/lexer.go
git commit -m "feat(oracle): implement stage 1 foundation infrastructure

Add NoLoc(), Token.End, ParseError Severity/Code fields, Parser.source,
RawStmt Loc migration, NodeLoc(), and ListSpan() to match PG patterns."
```

## Important Notes

- The `RawStmt` migration is the riskiest change — it touches serialization (`outfuncs.go`) and every place that reads `StmtLocation`/`StmtLen`. Search thoroughly.
- The `ParseError.Error()` format change may break existing test assertions that match on error strings. Find and update those.
- `Token.End` requires changes throughout the lexer — every code path that produces a token must set `End`.
