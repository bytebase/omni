package review

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
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
			// A canceled review starts no more targets.
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()
			if ctx.Err() != nil {
				return
			}
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
			// A constraint the change renamed to the name is not the synced
			// one of that name.
			if renamedTo(stmts, d.position, d.loc.Start, d.name) {
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
			// An object the change dropped before may have been replaced
			// by one of another kind.
			if droppedBefore(stmts, rn.position, rn.loc.Start, rn.name) {
				continue
			}
			if kind := x.relationKind(rn.name); kind != "" {
				r.report(review.DisallowRename, rn.statement, rangeOf(rn.loc), "renames "+kind+" "+object(rn.name)+" to "+ident(rn.newName))
			}
		}
	}
	return r.findings, nil
}

// droppedBefore reports whether the change, before offset at in statement
// pos, drops an object the reference may mean, or renames it away, so
// that the synced schema no longer tells what the name holds.
func droppedBefore(stmts []statement, pos, at int, ref *ast.TableRef) bool {
	dropped := false
	for i := 0; i <= pos && !dropped; i++ {
		ast.Inspect(stmts[i].node, func(n ast.Node) bool {
			if i == pos && ast.NodeLoc(n).End > at {
				// Not before the rename: only what it holds may be.
				return !dropped
			}
			switch v := n.(type) {
			case *ast.DropStmt:
				if v.Names != nil {
					for _, item := range v.Names.Items {
						if t, ok := item.(*ast.TableRef); ok && sameObject(t, ref) {
							dropped = true
						}
					}
				}
			case *ast.RenameStmt:
				dropped = dropped || v.ColumnName == "" && sameObject(v.Name, ref)
			case *ast.ExecStmt:
				if name, _, kind, ok := spRename(v); ok && (kind == "" || strings.EqualFold(kind, "OBJECT")) {
					if t := parseObjectName(name); t != nil && sameObject(t, ref) {
						dropped = true
					}
				}
			}
			return !dropped
		})
	}
	return dropped
}

// renamedTo reports whether the change, before offset at in statement pos,
// gives something the name with sp_rename.
func renamedTo(stmts []statement, pos, at int, name string) bool {
	renamed := false
	for i := 0; i <= pos && !renamed; i++ {
		ast.Inspect(stmts[i].node, func(n ast.Node) bool {
			if i == pos && ast.NodeLoc(n).End > at {
				return !renamed
			}
			if v, ok := n.(*ast.ExecStmt); ok {
				if _, newName, _, ok := spRename(v); ok && sameName(newName, name) {
					renamed = true
				}
			}
			return !renamed
		})
	}
	return renamed
}
