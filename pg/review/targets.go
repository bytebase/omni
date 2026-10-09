package review

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"sync"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// targetWork reports which of the rules that read a target have
// something to look at in this change: the prior backup check when the
// change asks for a backup and has an UPDATE or DELETE to back up, and
// the scan when a statement can give it a finding. The others depend on
// the SQL alone and run once for every target.
func targetWork(stmts []statement, on map[review.Rule]bool, change review.Change) (backup, scans bool) {
	if on[review.PriorBackup] && change.PriorBackup {
		for i := range stmts {
			switch stmts[i].node.(type) {
			case *ast.UpdateStmt, *ast.DeleteStmt:
				backup = true
			}
		}
	}
	return backup, needsScan(stmts, on)
}

// needsScan reports whether a statement of the change can give the scan
// something to report: a CREATE TABLE, alone or in CREATE SCHEMA, for
// RequirePrimaryKey, or an ALTER TABLE dropping a constraint or a column
// for either rule.
func needsScan(stmts []statement, on map[review.Rule]bool) bool {
	if !on[review.DisallowDropConstraint] && !on[review.RequirePrimaryKey] {
		return false
	}
	for i := range stmts {
		switch v := stmts[i].node.(type) {
		case *ast.CreateStmt:
			if on[review.RequirePrimaryKey] {
				return true
			}
		case *ast.CreateSchemaStmt:
			if on[review.RequirePrimaryKey] && v.SchemaElts != nil && slices.ContainsFunc(v.SchemaElts.Items, func(n ast.Node) bool {
				_, ok := n.(*ast.CreateStmt)
				return ok
			}) {
				return true
			}
		case *ast.AlterTableStmt:
			if v.Cmds == nil {
				continue
			}
			for _, item := range v.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch ast.AlterTableType(cmd.Subtype) {
				case ast.AT_DropConstraint:
					return true
				case ast.AT_DropColumn:
					if on[review.RequirePrimaryKey] {
						return true
					}
				}
			}
		}
	}
	return false
}

// targetGroup is the targets with the same inputs, reviewed once.
type targetGroup struct {
	indices []int
	target  review.Target
}

// groupTargets folds targets passing the same schema and session user into
// one group each, in first-seen order. The rules read nothing else of a
// target. Equal schemas passed as different values are reviewed apart:
// comparing them would cost more than reviewing them.
func groupTargets(targets []review.Target) []targetGroup {
	type key struct {
		schema *metadata.DatabaseSchemaMetadata
		user   string
	}
	var groups []targetGroup
	byKey := make(map[key]int)
	for i, t := range targets {
		k := key{t.Schema, t.SessionUser}
		g, ok := byKey[k]
		if !ok {
			g = len(groups)
			byKey[k] = g
			groups = append(groups, targetGroup{target: t})
		}
		groups[g].indices = append(groups[g].indices, i)
	}
	return groups
}

// reviewTargets evaluates the rules that read a target: each distinct
// target once, concurrently, then the findings merged by (rule,
// statement, range, message) so equal findings from different targets
// become one finding listing them all. A group that could not be
// evaluated becomes a failure per target.
func reviewTargets(ctx context.Context, stmts []statement, on map[review.Rule]bool, backup, scans bool, targets []review.Target) ([]review.Finding, []review.TargetFailure, error) {
	groups := groupTargets(targets)
	type outcome struct {
		findings []review.Finding
		err      error
	}
	outcomes := make([]outcome, len(groups))
	sem := make(chan struct{}, max(1, min(len(groups), runtime.GOMAXPROCS(0))))
	var wg sync.WaitGroup
	for i := range groups {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			findings, err := reviewGroup(ctx, stmts, on, backup, scans, groups[i])
			outcomes[i] = outcome{findings: findings, err: err}
		}()
	}
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	type key struct {
		rule      review.Rule
		statement int
		rng       review.Range
		message   string
	}
	merged := make(map[key]int)
	var findings []review.Finding
	var failures []review.TargetFailure
	for i, o := range outcomes {
		if o.err != nil {
			for _, t := range groups[i].indices {
				failures = append(failures, review.TargetFailure{Target: t, Err: o.err})
			}
			continue
		}
		for _, f := range o.findings {
			k := key{f.Rule, f.Statement, f.Range, f.Message}
			if j, ok := merged[k]; ok {
				findings[j].Targets = append(findings[j].Targets, f.Targets...)
				continue
			}
			merged[k] = len(findings)
			findings = append(findings, f)
		}
	}
	for i := range findings {
		slices.Sort(findings[i].Targets)
	}
	slices.SortFunc(failures, func(a, b review.TargetFailure) int { return a.Target - b.Target })
	return findings, failures, nil
}

// reviewGroup evaluates one distinct target. A panic is the group's
// failure, not the review's.
func reviewGroup(ctx context.Context, stmts []statement, on map[review.Rule]bool, backup, scans bool, g targetGroup) (findings []review.Finding, err error) {
	defer func() {
		if p := recover(); p != nil {
			findings, err = nil, fmt.Errorf("review panicked: %v", p)
		}
	}()
	r := &reporter{targets: g.indices}
	if backup {
		checkPriorBackupTarget(stmts, g.target.Schema, r)
	}
	if scans {
		if err := scanTarget(ctx, stmts, on, g.target, r); err != nil {
			return nil, err
		}
	}
	return r.findings, nil
}
