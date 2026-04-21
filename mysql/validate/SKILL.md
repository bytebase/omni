# MySQL Validate

Static semantic analysis of MySQL ASTs produced by `mysql/parser`.

## Scope

Parser is grammar-only; validator covers what MySQL's `sp_head::parse` /
`sp_pcontext` layer checks after parse but before runtime. Diagnostics
are structured data (`[]Diagnostic`) — no panics, no parser-coupled
errors.

## Usage

```go
list, err := parser.Parse(sql)
if err != nil {
    return err // grammar error
}
diags := validate.Validate(list, validate.Options{})
```

## Diagnostic codes (stable API)

| Code | Meaning |
|---|---|
| `duplicate_variable` | Two `DECLARE` for the same name in one scope |
| `duplicate_condition` | Two `DECLARE CONDITION` for the same name |
| `duplicate_cursor` | Two `DECLARE CURSOR` for the same name |
| `duplicate_label` | Same begin-label in a sibling/nested scope (barrier: HANDLER body) |
| `duplicate_handler_condition` | Same condition listed twice in one HANDLER FOR list |
| `undeclared_condition` | `HANDLER FOR <name>` without matching `DECLARE CONDITION` |
| `undeclared_label` | `LEAVE <name>` with no enclosing label |
| `undeclared_loop_label` | `ITERATE <name>` with no enclosing loop label |
| `undeclared_cursor` | `OPEN`/`FETCH`/`CLOSE <name>` without matching `DECLARE CURSOR` |
| `return_outside_function` | `RETURN` outside a function body |
| `function_missing_return` | Function body with no `RETURN` in any path |
| `undeclared_variable` | Bare-name SET target that isn't a declared local nor a known sysvar |

## Adding a rule

1. Add or extend a walker case in `validate.go`.
2. Use `v.emit("<code>", "<message>", pos)`.
3. Add a test — direct-AST in `validate_test.go` or SQL-fed in `semantic_sql_test.go`.
4. Document the code in the table above.

## System-variable table

`sysvars.go` embeds the MySQL 8.0 known system variable set. Refresh
procedure lives in `sysvars_source.md`.
