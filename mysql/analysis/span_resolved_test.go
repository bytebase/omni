package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/enginetest/span"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
)

// This file reads the query MySQL reports it resolved: the text EXPLAIN
// returns (the "query" of EXPLAIN FORMAT=JSON), after MySQL merged views and
// derived tables into the query blocks that read them. MySQL prints it in
// one form: every column reference qualified as `db`.`table`.`column`, or
// `alias`.`column` for a table read through an alias and for a derived
// table or CTE it materialized;
// every table as `db`.`table` with its alias; every select item followed by
// AS `name`. This reader knows that form only.

type rtok struct {
	kind byte // '`' identifier, '\'' string, 'w' word or number, 'c' comment, else punctuation
	text string
}

func lexResolved(s string) ([]rtok, error) {
	var out []rtok
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == ' ' || c == '\n' || c == '\t' || c == '\r':
			i++
		case c == '`':
			var b strings.Builder
			j := i + 1
			for ; j < len(s); j++ {
				if s[j] == '`' {
					if j+1 < len(s) && s[j+1] == '`' {
						b.WriteByte('`')
						j++
						continue
					}
					break
				}
				b.WriteByte(s[j])
			}
			if j >= len(s) {
				return nil, fmt.Errorf("an unclosed identifier at %d", i)
			}
			out = append(out, rtok{'`', b.String()})
			i = j + 1
		case c == '\'' || c == '"':
			j := i + 1
			for ; j < len(s); j++ {
				if s[j] == '\\' {
					j++
					continue
				}
				if s[j] == c {
					if j+1 < len(s) && s[j+1] == c {
						j++
						continue
					}
					break
				}
			}
			if j >= len(s) {
				return nil, fmt.Errorf("an unclosed string at %d", i)
			}
			out = append(out, rtok{'\'', s[i : j+1]})
			i = j + 1
		case strings.HasPrefix(s[i:], "/*"):
			j := strings.Index(s[i:], "*/")
			if j < 0 {
				return nil, fmt.Errorf("an unclosed comment at %d", i)
			}
			out = append(out, rtok{'c', s[i : i+j+2]})
			i += j + 2
		case isWordByte(c):
			j := i
			for j < len(s) && (isWordByte(s[j]) || (s[j] == '.' && j+1 < len(s) && s[j+1] >= '0' && s[j+1] <= '9' && s[i] >= '0' && s[i] <= '9')) {
				j++
			}
			out = append(out, rtok{'w', strings.ToLower(s[i:j])})
			i = j
		default:
			out = append(out, rtok{c, string(c)})
			i++
		}
	}
	return out, nil
}

func isWordByte(c byte) bool {
	return c == '_' || c == '$' || c == '@' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}

func (t rtok) is(word string) bool { return t.kind == 'w' && t.text == word }

// A resolvedQuery is a query expression: its branches, each a query block,
// and the CTEs its WITH clause defines.
type resolvedQuery struct {
	ctes     map[string]*derivedRel
	branches []*resolvedBlock
	parent   *resolvedScope
}

// A resolvedBlock is one query block.
type resolvedBlock struct {
	items  [][]rtok // each select item, without its AS name
	names  []string // each select item's AS name
	tables map[string]string
	// tables maps each alias in FROM to db.table; derived maps each derived
	// table's alias to it.
	derived map[string]*derivedRel
}

// A derivedRel is a derived table or a CTE.
type derivedRel struct {
	columns []string // its column list, when it names one
	body    []rtok
	query   *resolvedQuery
}

// resolvedScope is what names in a query block resolve against: the block,
// then the WITH clauses around it.
type resolvedScope struct {
	block  *resolvedBlock
	query  *resolvedQuery
	parent *resolvedScope
}

// split returns toks cut at each depth-0 token cut accepts.
func split(toks []rtok, cut func(rtok) bool) [][]rtok {
	var out [][]rtok
	depth, start := 0, 0
	for i, t := range toks {
		switch t.kind {
		case '(':
			depth++
		case ')':
			depth--
		}
		if depth == 0 && cut(t) {
			out = append(out, toks[start:i])
			start = i + 1
		}
	}
	return append(out, toks[start:])
}

// closing returns the index of the parenthesis that closes toks[open].
func closing(toks []rtok, open int) int {
	depth := 0
	for i := open; i < len(toks); i++ {
		switch toks[i].kind {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func stripComments(toks []rtok) []rtok {
	var out []rtok
	for _, t := range toks {
		if t.kind != 'c' {
			out = append(out, t)
		}
	}
	return out
}

// startsQuery reports whether toks, the inside of parentheses, is a query.
func startsQuery(toks []rtok) bool {
	for len(toks) > 0 && toks[0].kind == '(' {
		toks = toks[1:]
	}
	return len(toks) > 0 && (toks[0].is("select") || toks[0].is("with") || toks[0].is("values") || toks[0].is("table"))
}

func parseResolved(toks []rtok) (*resolvedQuery, error) {
	toks = stripComments(toks)
	q := &resolvedQuery{ctes: map[string]*derivedRel{}}
	if len(toks) > 0 && toks[0].is("with") {
		i := 1
		if i < len(toks) && toks[i].is("recursive") {
			i++
		}
		for {
			if i >= len(toks) || toks[i].kind != '`' {
				return nil, fmt.Errorf("a CTE without a name")
			}
			d := &derivedRel{}
			name := toks[i].text
			i++
			if i < len(toks) && toks[i].kind == '(' {
				end := closing(toks, i)
				for _, c := range toks[i+1 : end] {
					if c.kind == '`' {
						d.columns = append(d.columns, c.text)
					}
				}
				i = end + 1
			}
			if i+1 >= len(toks) || !toks[i].is("as") || toks[i+1].kind != '(' {
				return nil, fmt.Errorf("CTE %s without AS (", name)
			}
			end := closing(toks, i+1)
			d.body = toks[i+2 : end]
			q.ctes[name] = d
			i = end + 1
			if i < len(toks) && toks[i].kind == ',' {
				i++
				continue
			}
			break
		}
		toks = toks[i:]
	}
	for _, branch := range split(toks, func(t rtok) bool { return t.is("union") || t.is("intersect") || t.is("except") }) {
		for len(branch) > 0 && (branch[0].is("all") || branch[0].is("distinct")) {
			branch = branch[1:]
		}
		if len(branch) > 0 && branch[0].kind == '(' && closing(branch, 0) == len(branch)-1 {
			inner, err := parseResolved(branch[1 : len(branch)-1])
			if err != nil {
				return nil, err
			}
			for name, d := range inner.ctes {
				q.ctes[name] = d
			}
			q.branches = append(q.branches, inner.branches...)
			continue
		}
		b, err := parseBlock(branch)
		if err != nil {
			return nil, err
		}
		q.branches = append(q.branches, b)
	}
	return q, nil
}

var clauseWords = map[string]bool{"where": true, "group": true, "having": true, "window": true, "order": true, "limit": true, "into": true, "for": true, "qualify": true}

func parseBlock(toks []rtok) (*resolvedBlock, error) {
	b := &resolvedBlock{tables: map[string]string{}, derived: map[string]*derivedRel{}}
	if len(toks) > 2 && toks[0].is("values") && toks[1].is("row") && toks[2].kind == '(' {
		// VALUES names its columns column_0, column_1, ...; the first row
		// gives their sources, and the others add theirs.
		for _, row := range split(toks[1:], func(t rtok) bool { return t.kind == ',' }) {
			if len(row) < 3 || !row[0].is("row") || row[1].kind != '(' || closing(row, 1) != len(row)-1 {
				return nil, fmt.Errorf("a VALUES row that is not ROW(...)")
			}
			for i, item := range split(row[2:len(row)-1], func(t rtok) bool { return t.kind == ',' }) {
				if i < len(b.items) {
					b.items[i] = append(append(append([]rtok{}, b.items[i]...), rtok{'+', "+"}), item...)
					continue
				}
				b.items = append(b.items, item)
				b.names = append(b.names, fmt.Sprintf("column_%d", i))
			}
		}
		return b, nil
	}
	if len(toks) == 0 || !toks[0].is("select") {
		return nil, fmt.Errorf("a query block that is not a SELECT")
	}
	i := 1
	for i < len(toks) && toks[i].kind == 'w' && (toks[i].text == "distinct" || toks[i].text == "all" || toks[i].text == "distinctrow" || toks[i].text == "straight_join" || toks[i].text == "high_priority" || strings.HasPrefix(toks[i].text, "sql_")) {
		i++
	}
	// The select list runs to the depth-0 FROM or the first clause.
	end := len(toks)
	depth := 0
	for j := i; j < len(toks); j++ {
		switch toks[j].kind {
		case '(':
			depth++
		case ')':
			depth--
		}
		if depth == 0 && (toks[j].is("from") || clauseWords[toks[j].text] && toks[j].kind == 'w') {
			end = j
			break
		}
	}
	for _, item := range split(toks[i:end], func(t rtok) bool { return t.kind == ',' }) {
		name := ""
		if n := len(item); n >= 2 && item[n-2].is("as") && item[n-1].kind == '`' {
			name, item = item[n-1].text, item[:n-2]
		}
		b.items = append(b.items, item)
		b.names = append(b.names, name)
	}
	if end < len(toks) && toks[end].is("from") {
		from := toks[end+1:]
		depth := 0
		for j, t := range from {
			switch t.kind {
			case '(':
				depth++
			case ')':
				depth--
			}
			if depth == 0 && t.kind == 'w' && clauseWords[t.text] {
				from = from[:j]
				break
			}
		}
		if err := b.readFrom(from); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// readFrom records the tables and derived tables a FROM clause names, at any
// depth of its joins.
func (b *resolvedBlock) readFrom(toks []rtok) error {
	for i := 0; i < len(toks); i++ {
		t := toks[i]
		switch {
		case t.kind == '(' && startsQuery(toks[i+1:]):
			end := closing(toks, i)
			d := &derivedRel{body: toks[i+1 : end]}
			j := end + 1
			if j < len(toks) && toks[j].is("as") {
				j++
			}
			if j >= len(toks) || toks[j].kind != '`' {
				return fmt.Errorf("a derived table without an alias")
			}
			name := toks[j].text
			j++
			if j < len(toks) && toks[j].kind == '(' {
				cend := closing(toks, j)
				for _, c := range toks[j+1 : cend] {
					if c.kind == '`' {
						d.columns = append(d.columns, c.text)
					}
				}
				j = cend + 1
			}
			b.derived[name] = d
			i = j - 1
		case t.kind == '(':
			// A join's parentheses or an ON condition: read on inside them.
		case t.is("on") || t.is("using"):
			// Skip the condition: it names columns, not tables.
			if i+1 < len(toks) && toks[i+1].kind == '(' {
				i = closing(toks, i+1)
			}
		case t.kind == '`' && i+2 < len(toks) && toks[i+1].kind == '.' && toks[i+2].kind == '`' && (i+3 >= len(toks) || toks[i+3].kind != '.'):
			db, table := t.text, toks[i+2].text
			alias := table
			j := i + 3
			if j < len(toks) && toks[j].kind == '`' {
				alias = toks[j].text
				j++
			}
			b.tables[alias] = db + "." + table
			i = j - 1
		}
	}
	return nil
}

// resolvedReader reads the lineage of a resolved query, against the schema a
// fixture stored.
type resolvedReader struct {
	views map[string]testcontainer.SpanRelation // db.view
	// stack holds the derived tables, CTEs and views being read, so a
	// recursive CTE's reference to itself adds nothing.
	stack map[any]bool
}

func newResolvedReader(fx testcontainer.SpanFixture) *resolvedReader {
	r := &resolvedReader{views: map[string]testcontainer.SpanRelation{}, stack: map[any]bool{}}
	for _, d := range fx.Databases {
		for _, v := range d.Views {
			r.views[d.Name+"."+v.Name] = v
		}
	}
	return r
}

// lineage is a set of source columns and whether MySQL leaves more unstated.
type lineage struct {
	sources  map[string]bool
	unstated bool
}

func (l *lineage) add(o lineage) {
	for s := range o.sources {
		l.sources[s] = true
	}
	l.unstated = l.unstated || o.unstated
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

func unknown() lineage { return lineage{sources: map[string]bool{}, unstated: true} }

// columns returns the lineage of each output column of q.
func (r *resolvedReader) columns(q *resolvedQuery, parent *resolvedScope) ([]lineage, error) {
	var out []lineage
	for k, b := range q.branches {
		scope := &resolvedScope{block: b, query: q, parent: parent}
		if k > 0 && len(b.items) != len(out) {
			return nil, fmt.Errorf("a branch of %d columns after one of %d", len(b.items), len(out))
		}
		for i, item := range b.items {
			l := r.item(item, scope)
			if k == 0 {
				out = append(out, l)
			} else {
				out[i].add(l)
			}
		}
	}
	return out, nil
}

// item is the lineage of one select item. A lone column reference states its
// column exactly, and an item that names no column states it has none. For
// any other expression MySQL states the columns it names and no more: the
// columns of a subquery inside it are not read here.
func (r *resolvedReader) item(toks []rtok, scope *resolvedScope) lineage {
	if l, ok := r.ref(toks, scope); ok {
		return l
	}
	out := lineage{sources: map[string]bool{}}
	named := false
	for i := 0; i < len(toks); i++ {
		t := toks[i]
		if t.kind == '(' && startsQuery(toks[i+1:]) {
			out.unstated = true
			i = closing(toks, i)
			continue
		}
		if t.kind != '`' || (i > 0 && toks[i-1].kind == '.') {
			continue
		}
		n := 1
		for i+2*n < len(toks) && toks[i+2*n-1].kind == '.' && toks[i+2*n].kind == '`' {
			n++
		}
		if n < 2 {
			// A bare identifier: a window's name, never a column in this form.
			continue
		}
		l, _ := r.ref(toks[i:i+2*n-1], scope)
		out.add(l)
		named = true
		i += 2*n - 2
	}
	if named {
		out.unstated = true
	}
	return out
}

// ref reads toks as one qualified column reference; ok is false when toks is
// anything else.
func (r *resolvedReader) ref(toks []rtok, scope *resolvedScope) (lineage, bool) {
	var parts []string
	for i, t := range toks {
		if i%2 == 1 {
			if t.kind != '.' {
				return lineage{}, false
			}
			continue
		}
		if t.kind != '`' {
			return lineage{}, false
		}
		parts = append(parts, t.text)
	}
	if len(toks)%2 == 0 || len(parts) < 2 || len(parts) > 3 {
		return lineage{}, false
	}
	column := parts[len(parts)-1]
	alias := parts[len(parts)-2]
	for s := scope; s != nil; s = s.parent {
		if s.block != nil {
			// A table MySQL reads by its own name is `db`.`table`.`column`;
			// one it reads through an alias is `alias`.`column`.
			if table, ok := s.block.tables[alias]; ok && (len(parts) == 2 || strings.HasPrefix(table, parts[0]+".")) {
				return r.tableColumn(table, column), true
			}
			if d, ok := s.block.derived[alias]; ok && len(parts) == 2 {
				return r.derivedColumn(d, column, s), true
			}
		}
		if d, ok := s.query.ctes[alias]; ok && len(parts) == 2 {
			return r.derivedColumn(d, column, s), true
		}
	}
	return unknown(), true
}

// tableColumn is a column of a table, or of a view MySQL did not merge,
// which reads through the view's definition.
func (r *resolvedReader) tableColumn(table, column string) lineage {
	v, ok := r.views[table]
	if !ok {
		return lineage{sources: map[string]bool{table + "." + column: true}}
	}
	if r.stack[table] {
		return unknown()
	}
	r.stack[table] = true
	defer delete(r.stack, table)
	toks, err := lexResolved(v.Definition)
	if err != nil {
		return unknown()
	}
	q, err := parseResolved(toks)
	if err != nil {
		return unknown()
	}
	cols, err := r.columns(q, nil)
	if err != nil {
		return unknown()
	}
	for i, c := range v.Columns {
		if c == column && i < len(cols) {
			return cols[i]
		}
	}
	return unknown()
}

// derivedColumn is a column of a derived table or CTE MySQL materialized.
func (r *resolvedReader) derivedColumn(d *derivedRel, column string, scope *resolvedScope) lineage {
	if r.stack[d] {
		return lineage{sources: map[string]bool{}}
	}
	r.stack[d] = true
	defer delete(r.stack, d)
	if d.query == nil {
		q, err := parseResolved(d.body)
		if err != nil {
			return unknown()
		}
		d.query = q
	}
	cols, err := r.columns(d.query, scope)
	if err != nil {
		return unknown()
	}
	names := d.columns
	if names == nil && len(d.query.branches) > 0 {
		names = d.query.branches[0].names
	}
	for i, c := range names {
		if c == column && i < len(cols) {
			return cols[i]
		}
	}
	return unknown()
}

// resolvedLineage is each output column's lineage, as MySQL's resolved query
// states it.
func resolvedLineage(fx testcontainer.SpanFixture, text string) ([][]string, error) {
	toks, err := lexResolved(text)
	if err != nil {
		return nil, err
	}
	q, err := parseResolved(toks)
	if err != nil {
		return nil, err
	}
	cols, err := newResolvedReader(fx).columns(q, nil)
	if err != nil {
		return nil, err
	}
	out := make([][]string, len(cols))
	for i, c := range cols {
		out[i] = c.list()
	}
	return out, nil
}
