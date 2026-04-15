# Snowflake Deparse (T3.2) — Implementation Plan

**Date:** 2026-04-15
**Node:** T3.2 deparse-core
**Status:** COMPLETE

## Steps

1. [x] Read AST package (parsenodes.go) to enumerate all node types and fields
2. [x] Study reference implementations (pg/deparse/ for patterns)
3. [x] Read parser tests to understand accepted SQL syntax
4. [x] Write spec: docs/superpowers/specs/2026-04-15-snowflake-deparse-design.md
5. [x] Write plan: this file
6. [x] Implement snowflake/deparse/ package (writer, expr, select, DML, DDL)
7. [x] Write round-trip test suite (130+ test cases)
8. [x] Fix bugs discovered by tests:
   - StarExpr handling in SelectTarget
   - Inline FK constraint syntax (FOREIGN KEY REFERENCES)
   - INSERT FIRST ELSE branch detection
9. [x] Run go test ./snowflake/... — all pass
10. [x] gofmt -w snowflake/deparse/
11. [x] Commit in logical chunks
12. [x] Push and open PR

## Implementation notes

- ~1400 LOC in deparse package (excluding tests)
- ~200 LOC test file with 130+ round-trip test cases
- All 6 snowflake/* packages pass with no regressions
