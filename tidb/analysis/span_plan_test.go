package analysis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
)

// This file reads what TiDB states about a query's lineage and tables.
//
// TiDB's plan (EXPLAIN FORMAT = 'brief') names every column it reads as
// db.table.column, the table's own name and never an alias or a view, and
// every column it computes as Column#N, where an operator defines it with
// expr->Column#N. The top projection of a plan lists the result set's
// columns in order. A plan without one (a plain read, a join, an aggregate
// or a union on top) lists no output columns.
//
// A column definition's org_table and org_name name the column a plain
// column reference reads, through derived tables and CTEs. Through a view
// TiDB keeps the table but takes the column's name from the view's select
// list and the database from the view (BuildDataSourceFromView in
// pkg/planner/core/logical_plan_builder.go), and through a CTE it takes the
// session's database, so they are read only where they name a column of a
// table the fixture has, and the definition's table is no view.

var (
	planColumnRef = regexp.MustCompile(`(?:^|[^\w$.#])([A-Za-z_$][\w$]*\.[\w$]+\.[\w$]+)`)
	planComputed  = regexp.MustCompile(`Column#\d+`)
	planDefine    = regexp.MustCompile(`^(.*)->(Column#\d+)(.*)$`)
)

// planNode is one operator of a plan.
type planNode struct {
	name   string // the operator, without its tree, role and number
	task   string
	info   string
	indent int
}

// passThrough are the operators whose output columns are their child's, in
// order, so the top projection may lie below them.
var passThrough = map[string]bool{
	"TableReader": true, "IndexReader": true, "IndexLookUp": true,
	"Sort": true, "TopN": true, "Limit": true, "Selection": true,
}

func readPlan(rows []string) ([]planNode, error) {
	var out []planNode
	for _, row := range rows {
		f := strings.Split(row, "\t")
		if len(f) != 4 {
			return nil, fmt.Errorf("a plan row %q", row)
		}
		id := f[0]
		indent := 0
		for _, r := range id {
			if strings.ContainsRune("└├│─ ", r) {
				indent++
				continue
			}
			break
		}
		name := strings.TrimLeft(id, "└├│─ ")
		if i := strings.IndexByte(name, '('); i >= 0 {
			name = name[:i]
		}
		if i := strings.LastIndexByte(name, '_'); i > 0 && strings.Trim(name[i+1:], "0123456789") == "" {
			name = name[:i]
		}
		out = append(out, planNode{name: name, task: f[1], info: f[3], indent: indent})
	}
	return out, nil
}

// splitTop cuts s at each comma outside parentheses and brackets.
func splitTop(s string) []string {
	var out []string
	depth, start := 0, 0
	for i, c := range s {
		switch c {
		case '(', '[':
			depth++
		case ')', ']':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	return append(out, strings.TrimSpace(s[start:]))
}

// planReader resolves a plan's computed columns to the columns they read.
type planReader struct {
	defs map[string]string // Column#N to the expression that defines it
}

func newPlanReader(plan []planNode) *planReader {
	r := &planReader{defs: map[string]string{}}
	for _, n := range plan {
		for _, item := range splitTop(n.info) {
			if m := planDefine.FindStringSubmatch(item); m != nil {
				expr := m[1] + m[3]
				if i := strings.IndexByte(expr, ':'); i >= 0 && !strings.ContainsAny(expr[:i], "( ") {
					expr = expr[i+1:] // funcs:, group by:
				}
				r.defs[m[2]] = expr
			}
		}
	}
	return r
}

// exact returns the one column expr is, through computed columns that are
// each one column, if it is one.
func (r *planReader) exact(expr string, seen map[string]bool) (string, bool) {
	expr = strings.TrimSpace(expr)
	if planComputed.FindString(expr) == expr {
		def, ok := r.defs[expr]
		if !ok || seen[expr] {
			return "", false
		}
		seen[expr] = true
		return r.exact(def, seen)
	}
	if m := planColumnRef.FindStringSubmatch(expr); m != nil && m[1] == expr {
		return expr, true
	}
	return "", false
}

// sources adds the columns expr reads to l, and marks it unstated where a
// computed column has no definition in the plan.
func (r *planReader) sources(expr string, l *lineage, seen map[string]bool) {
	for _, m := range planColumnRef.FindAllStringSubmatch(expr, -1) {
		l.sources[m[1]] = true
	}
	for _, c := range planComputed.FindAllString(expr, -1) {
		if seen[c] {
			continue
		}
		seen[c] = true
		def, ok := r.defs[c]
		if !ok {
			l.unstated = true
			continue
		}
		r.sources(def, l, seen)
	}
}

// topProjection returns the expressions of the projection that makes a
// plan's output, if it has one.
func topProjection(plan []planNode) ([]string, bool) {
	for i := 0; i < len(plan); {
		n := plan[i]
		if n.name == "Projection" {
			var items []string
			for _, item := range splitTop(n.info) {
				if m := planDefine.FindStringSubmatch(item); m != nil {
					item = m[2]
				}
				items = append(items, item)
			}
			return items, true
		}
		if !passThrough[n.name] || i+1 >= len(plan) || plan[i+1].indent <= n.indent {
			return nil, false
		}
		i++
	}
	return nil, false
}

// lineage is a set of source columns and whether TiDB leaves more unstated.
type lineage struct {
	sources  map[string]bool
	unstated bool
}

func (l lineage) list() []string {
	var out []string
	for s := range l.sources {
		out = append(out, s)
	}
	if l.unstated {
		out = append(out, span.Unstated)
	}
	return out
}

// fixtureNames maps each of a fixture's databases, tables and views, and
// their columns, by lower-case name to the name it stored.
type fixtureNames struct {
	tables  map[string]string // db.table, lower case, to its stored name
	views   map[string]string
	columns map[string]string // db.table.column, lower case
}

func newFixtureNames(fx testcontainer.SpanFixture) fixtureNames {
	n := fixtureNames{tables: map[string]string{}, views: map[string]string{}, columns: map[string]string{}}
	for _, d := range fx.Databases {
		for _, t := range d.Tables {
			key := strings.ToLower(d.Name + "." + t.Name)
			n.tables[key] = d.Name + "." + t.Name
			for _, c := range t.Columns {
				n.columns[key+"."+strings.ToLower(c)] = d.Name + "." + t.Name + "." + c
			}
		}
		for _, v := range d.Views {
			n.views[strings.ToLower(d.Name+"."+v.Name)] = d.Name + "." + v.Name
		}
	}
	return n
}

// stored spells a plan's db.table.column as the fixture stored it.
func (n fixtureNames) stored(col string) string {
	if s, ok := n.columns[strings.ToLower(col)]; ok {
		return s
	}
	return col
}

// planLineage is each output column's lineage as TiDB states it.
func planLineage(fx testcontainer.SpanFixture, a testcontainer.SpanAnswer) ([][]string, error) {
	plan, err := readPlan(a.Plan)
	if err != nil {
		return nil, err
	}
	names := newFixtureNames(fx)
	r := newPlanReader(plan)
	items, projected := topProjection(plan)
	if projected && len(items) != len(a.Columns) {
		return nil, fmt.Errorf("the top projection has %d columns and the result set %d", len(items), len(a.Columns))
	}
	out := make([][]string, len(a.Columns))
	for i, c := range a.Columns {
		l := lineage{sources: map[string]bool{}}
		switch {
		case projected:
			if col, ok := r.exact(items[i], map[string]bool{}); ok {
				l.sources[col] = true
				break
			}
			// An expression states the columns it names, not which of them
			// decide its value. Nor does a constant state that it reads
			// nothing: TiDB evaluates an uncorrelated subquery while it
			// plans, and the plan holds its value.
			r.sources(items[i], &l, map[string]bool{})
			l.unstated = true
		case c.OrgTable != "" && names.views[strings.ToLower(c.Schema+"."+c.Table)] == "":
			col := c.Schema + "." + c.OrgTable + "." + c.OrgName
			if _, ok := names.columns[strings.ToLower(col)]; ok {
				l.sources[col] = true
			} else {
				l.unstated = true
			}
		default:
			l.unstated = true
		}
		var list []string
		for _, s := range l.list() {
			if s != span.Unstated {
				s = names.stored(s)
			}
			list = append(list, s)
		}
		out[i] = list
	}
	return out, nil
}

// summaryTables reads the statement summary's TABLE_NAMES: the tables and
// views the statement names, as the fixture stored them. A view stands for
// tables TiDB does not name, which are Unstated.
func summaryTables(fx testcontainer.SpanFixture, list string) ([]string, error) {
	names := newFixtureNames(fx)
	var out []string
	for _, t := range strings.Split(list, ",") {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		key := strings.ToLower(t)
		if v, ok := names.views[key]; ok {
			out = append(out, v, span.Unstated)
			continue
		}
		s, ok := names.tables[key]
		if !ok {
			return nil, fmt.Errorf("the statement summary names %s, which the fixture does not have", t)
		}
		out = append(out, s)
	}
	return out, nil
}
