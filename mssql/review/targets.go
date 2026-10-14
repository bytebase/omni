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
// this change: the prior backup a change asks for, the DROP CONSTRAINT and
// sp_rename statements the SQL alone cannot classify. The other rules
// depend on the SQL alone and run once for every target.
type targetWork struct {
	sql     string
	on      map[review.Rule]bool
	backup  bool
	made    changeMade
	drops   []droppedConstraint
	renames []renamedObject
}

func (w targetWork) any() bool {
	return w.backup || len(w.drops) > 0 || len(w.renames) > 0
}

// targetGroup is the targets with the same inputs, reviewed once.
type targetGroup struct {
	indices []int
	target  review.Target
}

// groupTargets folds targets passing the same schema and backup database
// presence into one group each, in first-seen order. The rules read
// nothing else of a target.
func groupTargets(targets []review.Target) []targetGroup {
	type key struct {
		schema *metadata.DatabaseSchemaMetadata
		backup bool
	}
	var groups []targetGroup
	byKey := make(map[key]int)
	for i, t := range targets {
		k := key{t.Schema, t.BackupDatabaseExists}
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
func reviewTargets(ctx context.Context, stmts []statement, work targetWork, targets []review.Target) ([]review.Finding, []review.TargetFailure, error) {
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
			findings, err := reviewGroup(stmts, work, groups[i])
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
func reviewGroup(stmts []statement, work targetWork, g targetGroup) (findings []review.Finding, err error) {
	defer func() {
		if p := recover(); p != nil {
			findings, err = nil, fmt.Errorf("review panicked: %v", p)
		}
	}()
	r := &reporter{targets: g.indices}
	x := schemaIndex{schema: g.target.Schema}
	if work.backup {
		checkPriorBackupGroups(work.sql, stmts, g.target.Schema.GetName(), r)
		checkPriorBackupTarget(stmts, g.target.BackupDatabaseExists, x, work.made, r)
	}
	if x.schema != nil {
		for _, d := range work.drops {
			// A USE in a block of the statement itself counts too.
			if usesDatabase(stmts, d.position+1) {
				continue
			}
			t := x.table(d.alter.Name)
			if t == nil {
				continue
			}
			kind, ok := constraintKind(t, d.name)
			if !ok {
				continue
			}
			reportDrop(stmts, d.position, &stmts[d.position], d.loc, d.alter, kind, d.name, d.nested, work.on, r)
		}
		for _, rn := range work.renames {
			if usesDatabase(stmts, rn.position+1) {
				continue
			}
			if kind := x.relationKind(rn.name); kind != "" {
				r.report(review.DisallowRename, rn.statement, rangeOf(rn.loc), "renames "+kind+" "+object(rn.name)+" to "+ident(rn.newName))
			}
		}
	}
	return r.findings, nil
}
