package review

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"sync"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

// targetWork is what the rules that read a target have to look at in
// this change: the prior backup and online migration a change asks for,
// and the DROP INDEX and DROP COLUMN statements the SQL alone cannot
// classify. The other rules depend on the SQL alone and run once for
// every target.
type targetWork struct {
	backup  bool
	made    changeMade
	online  bool
	drops   []droppedIndex
	columns columnDrops
}

func (w targetWork) any() bool {
	return w.backup || w.online || len(w.drops) > 0 || w.columns.any()
}

// schemaTarget is what the rules read of a target.
type schemaTarget struct {
	schema              *metadata.DatabaseSchemaMetadata
	backupExists        bool
	lowerCaseTableNames int
}

// targetGroup is the targets with the same inputs, reviewed once.
type targetGroup struct {
	indices []int
	target  review.Target
}

// groupTargets folds targets passing the same schema, backup database
// presence, and lower_case_table_names into one group each, in first-seen
// order. The rules read nothing else of a target. Equal schemas passed as
// different values are reviewed apart: comparing them would cost more
// than reviewing them.
func groupTargets(targets []review.Target) []targetGroup {
	var groups []targetGroup
	byKey := make(map[schemaTarget]int)
	for i, t := range targets {
		k := schemaTarget{t.Schema, t.BackupDatabaseExists, t.LowerCaseTableNames}
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
func reviewTargets(ctx context.Context, sql string, stmts []statement, work targetWork, targets []review.Target) ([]review.Finding, []review.TargetFailure, error) {
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
			findings, err := reviewGroup(sql, stmts, work, groups[i])
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
	// A group may report a finding twice; each target is listed once.
	for i := range findings {
		slices.Sort(findings[i].Targets)
		findings[i].Targets = slices.Compact(findings[i].Targets)
	}
	slices.SortFunc(failures, func(a, b review.TargetFailure) int { return a.Target - b.Target })
	return findings, failures, nil
}

// reviewGroup evaluates one distinct target. A panic is the group's
// failure, not the review's.
func reviewGroup(sql string, stmts []statement, work targetWork, g targetGroup) (findings []review.Finding, err error) {
	defer func() {
		if p := recover(); p != nil {
			findings, err = nil, fmt.Errorf("review panicked: %v", p)
		}
	}()
	r := &reporter{targets: g.indices}
	t := g.target
	x := newSchemaIndex(schemaTarget{t.Schema, t.BackupDatabaseExists, t.LowerCaseTableNames})
	if work.online {
		checkOnlineMigrationTarget(stmts, t, x, r)
	}
	if work.backup {
		checkPriorBackupMixed(stmts, t.Schema.GetName(), r)
		checkPriorBackupTarget(sql, stmts, t, x, work.made, r)
	}
	for _, d := range work.drops {
		checkDroppedIndex(stmts, d, x, r)
	}
	checkColumnDrops(work.columns, x, r)
	return r.findings, nil
}
