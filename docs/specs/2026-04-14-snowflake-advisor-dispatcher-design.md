# Snowflake Advisor Dispatcher — Design Spec (T2.6)

Date: 2026-04-14
Author: Claude (agent)
Status: Accepted

---

## Context

T2.7 will implement 14 lint rules for Snowflake SQL. Before those rules can land, we need
the framework that runs rules against a parsed AST and collects findings. This is T2.6.

The omni repo has no existing advisor framework for any engine — this is the first one. We
design it to be reusable but keep the scope narrow (YAGNI): just enough that T2.7 can add
14 rules by implementing one interface each.

---

## Goals

1. A `Rule` interface that rules implement (ID, Severity, Check).
2. A per-instance `Advisor` that holds a rule set and runs them over an AST.
3. A `Context` struct that rules can inspect (source text, config knobs — minimal for now).
4. A `Finding` value type (rule ID, severity, location, message).
5. A `Severity` type with ERROR / WARNING / INFO constants matching bytebase conventions.
6. A stub example rule (`NoSelectStar`) to prove the wiring works.
7. Tests: registration, invocation, finding collection, location accuracy.

---

## Non-Goals

- No global registry (per-instance is more testable; global can be added later).
- No rule config / payload beyond what Rule.Check receives in Context.
- No database schema / catalog in Context (T2.7 rules don't need it yet).
- The 14 T2.7 rules themselves.

---

## Package layout

```
snowflake/advisor/
  advisor.go         — Advisor, Context, Finding, Severity
  rule.go            — Rule interface
  example_rule.go    — NoSelectStar (trivial example rule)
  advisor_test.go    — framework tests + example rule tests
```

---

## Key types

### Severity

```go
type Severity int

const (
    SeverityInfo    Severity = iota
    SeverityWarning
    SeverityError
)
```

Matches bytebase's INFO / WARNING / ERROR ordering (lower = less severe).

### Finding

```go
type Finding struct {
    RuleID   string
    Severity Severity
    Loc      ast.Loc
    Message  string
}
```

`Loc` carries the byte range from the AST node that triggered the finding. Callers (e.g.
bytebase) can map byte offsets to line/column using a line table.

### Rule interface

```go
type Rule interface {
    // ID returns the canonical rule identifier, e.g. "snowflake.select.no-select-star".
    ID() string

    // Severity is the default severity for findings emitted by this rule.
    // Individual findings may carry the same or a different severity if the
    // rule implementation wishes to vary it.
    Severity() Severity

    // Check is called for each AST node during a depth-first walk.
    // It must be non-blocking and side-effect-free outside the returned slice.
    Check(ctx *Context, node ast.Node) []*Finding
}
```

Why a method-per-node approach was rejected: the AST has ~30 node types; storing one
callback per node type in the registry would add complexity with no gain for 14 rules.

### Context

```go
type Context struct {
    // SQL is the original source text. Rules may use it for substring extraction.
    SQL string
}
```

Deliberately minimal. If T2.7 needs catalog data, Context will grow then.

### Advisor

```go
type Advisor struct {
    rules []Rule
}

func New(rules ...Rule) *Advisor

func (a *Advisor) Check(ctx *Context, root ast.Node) []*Finding
```

`Check` performs a depth-first `ast.Walk` using itself as the `ast.Visitor`. On each
`Visit(node)` call it fans out to every registered rule and collects their findings.

The `Advisor` implements `ast.Visitor`:

```go
func (a *Advisor) Visit(node ast.Node) ast.Visitor {
    if node == nil {
        return nil  // post-order signal; nothing to do
    }
    for _, r := range a.rules {
        a.findings = append(a.findings, r.Check(a.ctx, node)...)
    }
    return a  // continue descent
}
```

Note: `findings` is stored on a per-run scratchpad, not on the Advisor struct itself, so
the Advisor is safe to reuse across calls.

---

## Example rule: NoSelectStar

Rule ID: `"snowflake.select.no-select-star"`
Severity: WARNING
Logic: emit a finding when `*SelectTarget` has `Star == true` and `Qualifier == nil`.

This exercises the full pipeline without depending on any real bytebase rule catalog.

---

## Testing strategy

1. Unit test: construct an Advisor with a counting Rule, walk a trivial AST, assert invocation count equals node count.
2. Integration test: parse a real SELECT * statement, run NoSelectStar rule, assert one WARNING finding with a correct Loc.
3. Location test: Loc.Start / Loc.End map to the expected source bytes.
4. Severity filtering test: construct a multi-rule Advisor, verify findings carry the rule's severity.
5. Empty input test: no stmts → no findings.

---

## Alternatives considered

**Functional rule (CheckFunc instead of interface)**
Simpler for trivial rules, but forces callers to store ID/severity separately.
Rejected in favour of the interface which bundles all metadata.

**Global registry (init()-based)**
Mirrors bytebase's advisor.Register pattern. Not needed for omni's current design where the
caller explicitly wires rules. Can be layered on top later.

**Node-type callbacks (per-type OnEnter/OnExit)**
The bytebase legacy pattern (GenericChecker.EnterEveryRule). Useful when rules need both
enter and exit events. None of the T2.7 rules need exit events; a single Check() is cleaner.
