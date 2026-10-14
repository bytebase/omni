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

// targetGroup is the targets with the same synced schema, reviewed once.
type targetGroup struct {
	indices []int
	schema  *metadata.DatabaseSchemaMetadata
}

// groupTargets folds targets passing the same schema into one group each,
// in first-seen order. The one rule that reads a target reads nothing else
// of it.
func groupTargets(targets []review.Target) []targetGroup {
	var groups []targetGroup
	bySchema := make(map[*metadata.DatabaseSchemaMetadata]int)
	for i, t := range targets {
		g, ok := bySchema[t.Schema]
		if !ok {
			g = len(groups)
			bySchema[t.Schema] = g
			groups = append(groups, targetGroup{schema: t.Schema})
		}
		groups[g].indices = append(groups[g].indices, i)
	}
	return groups
}

// reviewTargets evaluates the DROP CONSTRAINT statements the SQL alone
// cannot classify against each distinct target once, concurrently, then
// merges the findings by (rule, statement, range, message) so equal
// findings from different targets become one finding listing them all. A
// group that could not be evaluated becomes a failure per target.
func reviewTargets(ctx context.Context, drops []droppedConstraint, on map[review.Rule]bool, targets []review.Target) ([]review.Finding, []review.TargetFailure, error) {
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
			findings, err := reviewGroup(drops, on, groups[i])
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

// reviewGroup evaluates one distinct target: each DROP CONSTRAINT of a name
// the change did not define, classified by the synced table's constraint of
// that name. A name the synced table does not hold is not judged. A panic
// is the group's failure, not the review's.
func reviewGroup(drops []droppedConstraint, on map[review.Rule]bool, g targetGroup) (findings []review.Finding, err error) {
	defer func() {
		if p := recover(); p != nil {
			findings, err = nil, fmt.Errorf("review panicked: %v", p)
		}
	}()
	r := &reporter{targets: g.indices}
	x := schemaIndex{schema: g.schema}
	for _, d := range drops {
		t := x.table(d.table)
		if t == nil {
			continue
		}
		if kind, ok := syncedConstraint(t, d.name.Normalize()); ok {
			reportDrop(d, kind, on, r)
		}
	}
	return r.findings, nil
}
