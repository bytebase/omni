# Snowflake Advisor Dispatcher — Implementation Plan (T2.6)

Date: 2026-04-14
Status: In progress

---

## Steps

### Step 1 — Core types (advisor.go + rule.go)

Files: `snowflake/advisor/advisor.go`, `snowflake/advisor/rule.go`

- Define `Severity` (int, ERROR/WARNING/INFO constants + String())
- Define `Finding` struct (RuleID, Severity, Loc, Message)
- Define `Context` struct (SQL string)
- Define `Rule` interface (ID, Severity, Check)
- Define `Advisor` struct with per-run state
- Implement `New(rules ...Rule) *Advisor`
- Implement `Advisor.Check(ctx *Context, root ast.Node) []*Finding`
- Implement `Advisor` as `ast.Visitor` (Visit method)

### Step 2 — Example rule (example_rule.go)

File: `snowflake/advisor/example_rule.go`

- Implement `NoSelectStarRule` (Rule interface)
- ID: `"snowflake.select.no-select-star"`
- Severity: WARNING
- Check: match `*ast.SelectTarget` with Star==true, Qualifier==nil

### Step 3 — Tests (advisor_test.go)

File: `snowflake/advisor/advisor_test.go`

- Test: empty rules, empty AST → no findings
- Test: counting rule is called once per node
- Test: NoSelectStar fires on `SELECT *`
- Test: NoSelectStar does NOT fire on `SELECT a, b`
- Test: NoSelectStar fires on `SELECT t.*` (qualified star — rule should NOT flag)
- Test: location accuracy — Loc.Start/End bracket the * token
- Test: severity on finding matches rule's Severity()

### Step 4 — Run tests + fmt

```
cd /Users/h3n4l/OpenSource/omni/.worktrees/snowflake-advisor
go test ./snowflake/advisor/... ./snowflake/... -count=1
gofmt -w snowflake/
```

### Step 5 — Commit + push + PR

Commit message: `feat(snowflake): T2.6 advisor dispatcher framework`

---

## Dependency tree

- advisor.go depends on snowflake/ast (Loc, Node, Walk)
- rule.go depends on advisor.go (Context, Finding, Severity)
- example_rule.go depends on snowflake/ast (SelectTarget, StarExpr) and advisor.go
- advisor_test.go depends on all of the above + snowflake/parser (Parse)
