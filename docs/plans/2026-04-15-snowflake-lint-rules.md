# Plan: Snowflake Lint Rules (T2.7)

## Steps

1. Read advisor framework (done)
2. Read AST node types (done)
3. Read parser.IsReservedKeyword (done)
4. Write spec (done — specs/2026-04-15-snowflake-lint-rules-design.md)
5. Implement 14 rules in `snowflake/advisor/rules/`:
   - Step 5a: Column rules (column_no_null, column_require, column_maximum_varchar_length)
   - Step 5b: Table rules (table_require_pk, table_no_foreign_key)
   - Step 5c: Naming rules (naming_table, naming_table_no_keyword, naming_identifier_case, naming_identifier_no_keyword)
   - Step 5d: Query rules (select_no_select_all, where_require_select, where_require_update_delete)
   - Step 5e: Migration rules (migration_compatibility, table_drop_naming_convention)
6. Write tests for each rule (co-located *_test.go files)
7. Write integration test (rules_integration_test.go)
8. Run `go test ./snowflake/...`
9. gofmt -w
10. Commit in logical chunks
11. Push and open PR
