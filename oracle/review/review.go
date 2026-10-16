// Package review implements Bytebase SQL Review V2 for Oracle: the engine
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
// the rule reports nothing. DDL inside PL/SQL runs only through EXECUTE
// IMMEDIATE, whose text no rule reads; the DML of a PL/SQL body, an
// anonymous block, or a trigger is read where a rule says so.
//
// Oracle folds an unquoted name to upper case and keeps a quoted one as
// written; the parser hands the rules names folded, and the rules compare
// them exactly, as Oracle does. A message writes a name as Oracle would
// need it: T for t, "Mixed" for "Mixed".
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
			result, err = nil, fmt.Errorf("oracle/review: panic: %v", p)
		}
	}()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r := &reporter{targets: allTargets(len(targets))}
	ranges, stmts, syntax := parseReview(sql)
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
	backup := on[review.PriorBackup] && opts.Change.PriorBackup
	if backup {
		checkPriorBackupSize(sql, opts.Change, r)
		checkPriorBackupStatements(stmts, r)
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
			checkDisallowDropObject(sql, s, r)
		}
		if on[review.DisallowTruncate] {
			checkDisallowTruncate(sql, s, r)
		}
		if on[review.DisallowRename] {
			checkDisallowRename(s, r)
		}
	}
	work := targetWork{
		sql:         sql,
		on:          on,
		backup:      backup,
		constraints: on[review.DisallowDropConstraint] || on[review.RequirePrimaryKey],
		renames:     on[review.DisallowRename] && hasRename(stmts),
	}
	if work.backup {
		work.made = madeByChange(stmts)
	}
	switch {
	case len(targets) == 0:
		// Without a target the rules that read one fall back to what the
		// change itself makes; target-dependent findings have no target to
		// report against.
		reviewSchemaRules(stmts, work, schemaIndex{}, r)
	case work.any():
		findings, failures, err := reviewTargets(ctx, stmts, work, targets)
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
