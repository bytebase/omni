package review

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"slices"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/review"
)

// targetGroup is the targets whose inputs are equal, reviewed once.
type targetGroup struct {
	indices []int
	target  review.Target
}

// groupTargets folds targets with equal inputs into one group each, in
// first-seen order. PostgreSQL reads Schema and SessionUser; the other
// fields do not change its findings.
func groupTargets(targets []review.Target) ([]targetGroup, error) {
	var groups []targetGroup
	byKey := make(map[string]int)
	for i, t := range targets {
		key, err := targetKey(t)
		if err != nil {
			return nil, fmt.Errorf("target %d: %w", i, err)
		}
		g, ok := byKey[key]
		if !ok {
			g = len(groups)
			byKey[key] = g
			groups = append(groups, targetGroup{target: t})
		}
		groups[g].indices = append(groups[g].indices, i)
	}
	return groups, nil
}

func targetKey(t review.Target) (string, error) {
	schema := "nil"
	if t.Schema != nil {
		b, err := proto.MarshalOptions{Deterministic: true}.Marshal(t.Schema)
		if err != nil {
			return "", err
		}
		sum := sha256.Sum256(b)
		schema = hex.EncodeToString(sum[:])
	}
	return schema + "\x00" + t.SessionUser, nil
}

// reviewTargets evaluates the rules that read a target: each distinct
// target once, concurrently, then the findings merged by (rule,
// statement, range, message) so equal findings from different targets
// become one finding listing them all. A group that could not be
// evaluated becomes a failure per target.
func reviewTargets(ctx context.Context, stmts []statement, on map[review.Rule]bool, targets []review.Target) ([]review.Finding, []review.TargetFailure, error) {
	groups, err := groupTargets(targets)
	if err != nil {
		return nil, nil, err
	}
	type outcome struct {
		findings []review.Finding
		err      error
	}
	outcomes := make([]outcome, len(groups))
	workers := min(len(groups), runtime.GOMAXPROCS(0))
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	for i := range groups {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			findings, err := reviewGroup(ctx, stmts, on, groups[i])
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

// reviewGroup evaluates one distinct target. A panic in the catalog is
// the group's failure, not the review's.
func reviewGroup(ctx context.Context, stmts []statement, on map[review.Rule]bool, g targetGroup) (findings []review.Finding, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("review panicked: %v", p)
		}
	}()
	s, err := newSession(ctx, g.target)
	if err != nil {
		return nil, err
	}
	r := &reporter{targets: g.indices}
	if err := s.walk(ctx, stmts, on, r); err != nil {
		return nil, err
	}
	return r.findings, nil
}
