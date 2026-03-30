# Eval Worker — Stage 1: Foundation

You are an Eval Worker in the Oracle Quality Pipeline.
Your role is to write **tests ONLY** — never implementation code.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/strategy.md` — Stage 1 scope (10 infrastructure items)
- `oracle/ast/node.go` — current Oracle `Loc` struct and `Node` interface
- `oracle/parser/parser.go` — current `ParseError` struct and `Parser` struct
- `oracle/parser/lexer.go` — current `Token` struct
- `oracle/ast/parsenodes.go` — current `RawStmt` struct
- `pg/ast/node.go` — PG's `Loc`, `NoLoc()` (reference for expected behavior)
- `pg/ast/loc.go` — PG's `NodeLoc()`, `ListSpan()` (reference)
- `pg/parser/parser.go` — PG's `ParseError` struct (reference)

## Output Files

- **Test file:** `oracle/parser/eval_foundation_test.go`
- **Coverage report:** `oracle/quality/coverage/stage1-foundation.json`

## Rules

1. **Tests ONLY** — do NOT modify any non-test `.go` file.
2. Every test function MUST be named `TestEvalStage1_*`.
3. Tests should fail clearly with descriptive messages (not just compile errors).
4. Use `reflect` where possible to check field existence so tests compile even when fields are missing.
5. Reference PG behavior to determine expected behavior for each item.
6. Do NOT import packages that do not exist yet; use reflection to probe for functions/fields.

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[EVAL-STAGE1] STARTED
[EVAL-STAGE1] STEP reading_refs - Reading PG reference files
[EVAL-STAGE1] STEP writing_tests - Writing eval tests (10 items)
[EVAL-STAGE1] STEP build_check - Running go build on test file
[EVAL-STAGE1] STEP coverage_report - Generating stage1-foundation.json
[EVAL-STAGE1] DONE
```

If a step fails:
```
[EVAL-STAGE1] FAIL step_name - description
[EVAL-STAGE1] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Enumerable Test Items (10 total)

Each item below becomes one `TestEvalStage1_*` function.

### Item 1: NoLoc()

**What:** `oracle/ast` package exports a `NoLoc()` function returning `Loc{Start: -1, End: -1}`.
**PG reference:** `pg/ast/node.go` — `func NoLoc() Loc { return Loc{Start: -1, End: -1} }`
**Test strategy:** Call `ast.NoLoc()` and assert `Start == -1` and `End == -1`.
**Function name:** `TestEvalStage1_NoLoc`

### Item 2: Loc Sentinel Consistency

**What:** The zero-value `Loc{}` must NOT be used as "unknown". Unknown positions use `-1`.
**PG reference:** PG Loc uses `-1` for unknown in both fields.
**Test strategy:** Verify `Loc{}.Start == 0` and `Loc{}.End == 0` (zero-value is valid position 0, not unknown). Verify `NoLoc().Start == -1` and `NoLoc().End == -1`.
**Function name:** `TestEvalStage1_LocSentinel`

### Item 3: Token.End

**What:** The `Token` struct has an `End` field (int) representing the exclusive end byte offset.
**Test strategy:** Use `reflect` to check that `Token` struct has field `End` of type `int`. Then lex a simple token and verify `End > Loc` (start).
**Function name:** `TestEvalStage1_TokenEnd`

### Item 4: ParseError.Severity

**What:** `ParseError` has a `Severity` field (string).
**PG reference:** `pg/parser/parser.go` — `Severity string` field.
**Test strategy:** Use `reflect` to check field existence. Verify default behavior: when Severity is empty, `Error()` should still produce output containing "ERROR".
**Function name:** `TestEvalStage1_ParseErrorSeverity`

### Item 5: ParseError.Code

**What:** `ParseError` has a `Code` field (string) for SQLSTATE codes.
**PG reference:** `pg/parser/parser.go` — `Code string` field, defaults to `"42601"`.
**Test strategy:** Use `reflect` to check field existence. Verify that `Error()` output contains the code (default `"42601"` for syntax errors).
**Function name:** `TestEvalStage1_ParseErrorCode`

### Item 6: ParseError.Error() Format

**What:** `ParseError.Error()` returns `"SEVERITY: message (SQLSTATE code)"` format.
**PG reference:** PG returns `fmt.Sprintf("%s: %s (SQLSTATE %s)", sev, e.Message, code)`.
**Test strategy:** Create a `ParseError` with known fields, call `Error()`, verify format matches `"ERROR: msg (SQLSTATE 42601)"`.
**Function name:** `TestEvalStage1_ParseErrorFormat`

### Item 7: Parser.source

**What:** The `Parser` struct stores the original SQL input in a `source` field (string).
**Test strategy:** Use `reflect` to check that `Parser` struct has field `source` of type `string`. Parse a SQL string, then use `reflect` on the parser (if accessible) or verify indirectly that the parser can reference source text.
**Function name:** `TestEvalStage1_ParserSource`

### Item 8: RawStmt Loc

**What:** `RawStmt` uses a `Loc` field (type `ast.Loc`) instead of separate `StmtLocation`/`StmtLen` fields.
**PG reference:** PG `RawStmt` has `Loc Loc` field.
**Test strategy:** Use `reflect` to check that `RawStmt` has field `Loc` of type `Loc`. Also check that `StmtLocation` and `StmtLen` fields are absent (migration complete).
**Function name:** `TestEvalStage1_RawStmtLoc`

### Item 9: NodeLoc()

**What:** `oracle/ast` exports a `NodeLoc(Node) Loc` function that extracts location from any node.
**PG reference:** `pg/ast/loc.go` — returns `NoLoc()` for nil, dispatches via type switch.
**Test strategy:** Call `ast.NodeLoc(nil)` and verify it returns `NoLoc()`. Call with a `RawStmt` and verify it returns the `Loc` field.
**Function name:** `TestEvalStage1_NodeLoc`

### Item 10: ListSpan()

**What:** `oracle/ast` exports a `ListSpan(*List) Loc` function returning the span from first to last item.
**PG reference:** `pg/ast/loc.go` — returns `NoLoc()` for nil/empty, otherwise `Loc{first.Start, last.End}`.
**Test strategy:** Call `ast.ListSpan(nil)` → `NoLoc()`. Call with empty list → `NoLoc()`. Call with populated list → verify span covers first to last element.
**Function name:** `TestEvalStage1_ListSpan`

## Coverage Report Format

After writing tests, generate `oracle/quality/coverage/stage1-foundation.json`:

```json
{
  "stage": "1-foundation",
  "total_items": 10,
  "tested_items": 10,
  "items": [
    {"id": "noloc", "test": "TestEvalStage1_NoLoc", "status": "written"},
    {"id": "loc_sentinel", "test": "TestEvalStage1_LocSentinel", "status": "written"},
    {"id": "token_end", "test": "TestEvalStage1_TokenEnd", "status": "written"},
    {"id": "parse_error_severity", "test": "TestEvalStage1_ParseErrorSeverity", "status": "written"},
    {"id": "parse_error_code", "test": "TestEvalStage1_ParseErrorCode", "status": "written"},
    {"id": "parse_error_format", "test": "TestEvalStage1_ParseErrorFormat", "status": "written"},
    {"id": "parser_source", "test": "TestEvalStage1_ParserSource", "status": "written"},
    {"id": "rawstmt_loc", "test": "TestEvalStage1_RawStmtLoc", "status": "written"},
    {"id": "nodeloc", "test": "TestEvalStage1_NodeLoc", "status": "written"},
    {"id": "listspan", "test": "TestEvalStage1_ListSpan", "status": "written"}
  ]
}
```

The `status` field transitions: `"written"` → `"passing"` (once impl worker makes it pass) → `"verified"` (once insight worker reviews).

## Verification

After writing the test file:

```bash
# Must compile (tests may fail, but must compile)
cd /Users/rebeliceyang/Github/omni && go build ./oracle/parser/

# Run eval tests to see current state (failures expected before impl)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"
```

## Important Notes

- Tests MUST compile even if the features are not yet implemented. Use `reflect` to probe for fields/functions that may not exist.
- For functions that may not exist (like `NoLoc`, `NodeLoc`, `ListSpan`), use a build-tag guarded approach or conditional compilation. If that is not possible, write the test assuming the function exists — the Impl Worker will add it.
- Prefer clear failure messages: `t.Fatalf("RawStmt missing Loc field; has StmtLocation/StmtLen instead — migration needed")`.
- Never modify existing test files (`compare_test.go`, etc.).
