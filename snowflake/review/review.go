// Package review implements Bytebase SQL Review V2 for Snowflake: the engine
// side of the contract in github.com/bytebase/omni/review.
//
// Review parses the SQL once, evaluates the rules that depend on the SQL
// alone once, and reports every finding against every target. The rules
// that read a target check the target's synced schema directly, without
// rebuilding it, and a finding equal across targets is reported once,
// listing them all. A syntax error ends the review: the result then holds
// that one finding and nothing else, because no other rule can be trusted
// on text the engine could not parse.
//
// A rule reports only what it can establish from the SQL and the synced
// schema, without simulating the change. Where they leave the answer open,
// the rule reports nothing.
//
// Snowflake has no local engine: no container image runs it, and CI cannot
// start one, so the tests are unit tests and the engine facts the rules rest
// on come from Snowflake's documentation:
//
//   - An unquoted identifier is stored and resolved in upper case, a quoted
//     one exactly as written, so names compare by their normalized form
//     (ast.Ident.Normalize). The comparison assumes the default
//     QUOTED_IDENTIFIERS_IGNORE_CASE = FALSE.
//     https://docs.snowflake.com/en/sql-reference/identifiers-syntax
//   - The = operator treats NULLs as unknown values: a comparison with NULL
//     is NULL, and a row it filters is not selected.
//     https://docs.snowflake.com/en/sql-reference/functions/equal_null
//   - Snowflake does not enforce a primary key on a standard table ("NOT
//     NULL and CHECK are the only types of constraints that are enforced"),
//     but Bytebase's rule still asks for one. A constraint created with the
//     non-default ENABLE or VALIDATE is not created at all, which the AST
//     does not record; such a key still counts as declared.
//     https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint
//   - CREATE TABLE ... AS SELECT infers only the column names and types from
//     its query; constraints are what its column list declares. CREATE TABLE
//     ... LIKE copies the source's constraints, CLONE its structure, and OR
//     REPLACE is the equivalent of DROP TABLE on the existing table. CREATE
//     OR ALTER TABLE adds, drops, or modifies constraints to match the
//     definition.
//     https://docs.snowflake.com/en/sql-reference/sql/create-table
//   - A temporary table exists only in the session that created it; a
//     transient table persists until explicitly dropped, like a permanent
//     table without Fail-safe.
//     https://docs.snowflake.com/en/user-guide/tables-temp-transient
//   - ALTER TABLE ... SWAP WITH "renames two tables in a single
//     transaction", and DROP { CONSTRAINT name | PRIMARY KEY | UNIQUE |
//     FOREIGN KEY } is the constraint drop.
//     https://docs.snowflake.com/en/sql-reference/sql/alter-table
//
// Bodies the parser keeps as text are not read: a function's or procedure's
// body, the string EXECUTE IMMEDIATE runs, and an alert's action the parser
// does not support. A Snowflake Scripting block, a task's body, and an
// alert's action the parser supports are parsed, and read.
package review

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/bytebase/omni/review"
)

// Rules lists the rules this package implements, in evaluation order.
// WalkThrough is not implemented.
var Rules = []review.Rule{
	review.Syntax,
	review.OnlineMigration,
	review.RequireIsNull,
	review.RequireWhere,
	review.DisallowDropObject,
	review.DisallowTruncate,
	review.DisallowRename,
	review.PriorBackup,
	review.DisallowDropConstraint,
	review.RequirePrimaryKey,
}

// Review evaluates opts.Rules over sql for every target. See the contract
// package for the meaning of each input and of the result. A panic in the
// parser or a rule is returned as an error: the review could not be
// evaluated, and the caller's goroutine survives.
func Review(ctx context.Context, sql string, opts review.Options, targets []review.Target) (result *review.Result, err error) {
	defer func() {
		if p := recover(); p != nil {
			result, err = nil, fmt.Errorf("snowflake/review: panic: %v", p)
		}
	}()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r := &reporter{targets: allTargets(len(targets))}
	ranges, stmts, syntax, err := parseReview(sql)
	if err != nil {
		return nil, err
	}
	result = &review.Result{Statements: ranges}
	if syntax != nil {
		r.add(*syntax)
		result.Findings = r.sorted()
		return result, nil
	}

	on := enabled(opts.Rules)
	if on[review.OnlineMigration] && opts.Change.OnlineMigration {
		checkOnlineMigration(r)
	}
	if on[review.PriorBackup] && opts.Change.PriorBackup {
		checkPriorBackup(r)
	}
	for i := range stmts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		s := &stmts[i]
		if on[review.RequireIsNull] {
			checkRequireIsNull(sql, s, r)
		}
		if on[review.RequireWhere] {
			checkRequireWhere(s, r)
		}
		if on[review.DisallowDropObject] {
			checkDisallowDropObject(s, r)
		}
		if on[review.DisallowTruncate] {
			checkDisallowTruncate(s, r)
		}
		if on[review.DisallowRename] {
			checkDisallowRename(s, r)
		}
	}
	var drops []droppedConstraint
	if on[review.DisallowDropConstraint] || on[review.RequirePrimaryKey] {
		drops = checkConstraints(stmts, on, r)
	}
	if len(drops) > 0 {
		findings, failures, err := reviewTargets(ctx, drops, on, targets)
		if err != nil {
			return nil, err
		}
		r.findings = append(r.findings, findings...)
		result.Failures = failures
	}
	result.Findings = r.sorted()
	return result, nil
}

// enabled resolves the requested rules to the set this package evaluates:
// nil means every rule in Rules, and a rule this package does not know is
// dropped. Syntax is always in the set.
func enabled(rules []review.Rule) map[review.Rule]bool {
	on := make(map[review.Rule]bool, len(Rules))
	if rules == nil {
		for _, rule := range Rules {
			on[rule] = true
		}
		return on
	}
	known := make(map[review.Rule]bool, len(Rules))
	for _, rule := range Rules {
		known[rule] = true
	}
	for _, rule := range rules {
		if known[rule] {
			on[rule] = true
		}
	}
	on[review.Syntax] = true
	return on
}

// allTargets is the target list of a finding that depends on the SQL
// alone: every index, ascending. It is never nil, so a review with no
// targets still reports an empty list rather than a missing one.
func allTargets(n int) []int {
	all := make([]int, n)
	for i := range all {
		all[i] = i
	}
	return all
}

// reporter collects findings. Rules that depend on the SQL alone report
// through report, which stamps every target on the finding.
type reporter struct {
	targets  []int
	findings []review.Finding
}

// report records a finding against every target. A statement of -1
// addresses the whole change.
func (r *reporter) report(rule review.Rule, statement int, rng review.Range, message string) {
	r.add(review.Finding{
		Rule:      rule,
		Statement: statement,
		Range:     rng,
		Message:   message,
	})
}

func (r *reporter) add(f review.Finding) {
	f.Targets = slices.Clone(r.targets)
	r.findings = append(r.findings, f)
}

// sorted returns the findings in the contract's order: by statement,
// range start, rule, and message. A whole-change finding (statement -1)
// sorts first.
func (r *reporter) sorted() []review.Finding {
	slices.SortStableFunc(r.findings, func(a, b review.Finding) int {
		return cmp.Or(
			cmp.Compare(a.Statement, b.Statement),
			cmp.Compare(a.Range.Start, b.Range.Start),
			cmp.Compare(a.Rule, b.Rule),
			cmp.Compare(a.Message, b.Message),
		)
	})
	return r.findings
}
