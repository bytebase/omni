package catalog

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/oceanbase/scope"
)

// analyzerScope wraps scope.Scope and adds analyzer-specific features:
// parent chain for correlated subqueries, and rteIdx mapping.
type analyzerScope struct {
	base   *scope.Scope
	rteMap []int       // parallel to base entries: entry index -> RTE index in Query.RangeTable
	cols   [][]*Column // parallel to base entries: entry index -> catalog Column pointers
	parent *analyzerScope

	// firstByName maps each name an entry was registered under to the first
	// such entry, for entries [0, indexed); see entryIndex.
	firstByName map[string]int
	indexed     int

	// quals records each relation in scope with its catalog/schema qualifiers,
	// preserving duplicates and remote relations whose columns are unknown.
	quals      []rteQualifier
	hasUnknown bool

	// fullJoinExprs maps `<entry>.<column>` to the expression a FULL join's
	// merged key takes: COALESCE(left, right), nested as the joins nest.
	fullJoinExprs map[string]AnalyzedExpr

	// markedJoins maps a join to the lowered USING column names every base
	// table under it is already marked coalesced for (markCoalescedInNode).
	markedJoins map[*JoinExprNodeQ]map[string]bool

	// aliases maps the lowered names of the select list to their entries, once
	// GROUP BY, HAVING and ORDER BY are being read: those clauses may name a
	// select item by its alias, as WHERE may not. A name a table in scope has
	// keeps meaning that column (selectAlias).
	aliases map[string]*TargetEntryQ

	// ctes is the CTEs this query sees, its own and its enclosing queries'.
	ctes map[string]*CommonTableExprQ

	// namesUncertain is set when a derived table or CTE in scope has a column the
	// analyzer cannot name as OceanBase does (selectDerivesNames). A name that then
	// fails to resolve is not known to be missing, and its error is no refusal.
	namesUncertain bool
}

// uncertain reports whether this scope or one around it has uncertain names.
func (s *analyzerScope) uncertain() bool {
	for ; s != nil; s = s.parent {
		if s.namesUncertain {
			return true
		}
	}
	return false
}

// unsure turns err into a plain error when names in scope are uncertain, so it
// reads as an analysis that could not finish rather than OceanBase's refusal.
func (s *analyzerScope) unsure(err error) error {
	if err != nil && s.uncertain() {
		return fmt.Errorf("%v (a derived column's name is not known)", err)
	}
	return err
}

// selectAliases maps the names of the select list's items to them, the first of
// each name.
func selectAliases(tl []*TargetEntryQ) map[string]*TargetEntryQ {
	m := make(map[string]*TargetEntryQ, len(tl))
	for _, te := range tl {
		if te.ResJunk || te.ResName == "" {
			continue
		}
		if _, ok := m[strings.ToLower(te.ResName)]; !ok {
			m[strings.ToLower(te.ResName)] = te
		}
	}
	return m
}

// selectAlias returns the select item an unqualified name that no table in scope
// has reads, when the clause being read may name one: this query's select list,
// then those of the queries around it, as OceanBase's resolver searches the outer
// select list for a name a subquery of GROUP BY, HAVING or ORDER BY leaves open
// (ob_select_resolver.cpp). A scope holds its aliases only from GROUP BY on, so a
// subquery in WHERE, read before that, still sees none.
func (s *analyzerScope) selectAlias(name string) (*TargetEntryQ, int) {
	lower := strings.ToLower(name)
	for levels := 0; s != nil; s, levels = s.parent, levels+1 {
		if te, ok := s.aliases[lower]; ok {
			return te, levels
		}
	}
	return nil, 0
}

func newScope() *analyzerScope {
	return &analyzerScope{
		base: scope.New(),
	}
}

// newScopeWithParent creates a new scope with a parent scope for correlated subquery resolution.
func newScopeWithParent(parent *analyzerScope) *analyzerScope {
	return &analyzerScope{
		base:   scope.New(),
		parent: parent,
	}
}

// addFullJoinExpr records the expression a FULL join's merged column takes in
// the star expansion of a scope entry.
func (s *analyzerScope) addFullJoinExpr(entryName, column string, expr AnalyzedExpr) {
	if s.fullJoinExprs == nil {
		s.fullJoinExprs = map[string]AnalyzedExpr{}
	}
	s.fullJoinExprs[strings.ToLower(entryName)+"."+strings.ToLower(column)] = expr
}

// fullJoinExpr returns the expression a FULL join's merged column takes in the
// star expansion of that entry and column, if the scope holds one.
func (s *analyzerScope) fullJoinExpr(entryName, column string) (AnalyzedExpr, bool) {
	expr, ok := s.fullJoinExprs[strings.ToLower(entryName)+"."+strings.ToLower(column)]
	return expr, ok
}

// markCoalesced marks a column from a table as coalesced (hidden during star expansion).
func (s *analyzerScope) markCoalesced(tableName, colName string) {
	s.base.MarkCoalesced(tableName, colName)
}

// isCoalesced returns true if the given table.column is coalesced away by USING/NATURAL.
func (s *analyzerScope) isCoalesced(tableName, colName string) bool {
	return s.base.IsCoalesced(tableName, colName)
}

type rteQualifier struct {
	catalog  string
	schema   string
	name     string // effective reference name (alias if present, else table name)
	hasAlias bool
	unknown  bool
	rteIdx   int
	entryIdx int // index in cols when !unknown, -1 when unknown
}

// addRelation registers a base, view or remote relation in the scope, keeping
// its catalog and schema qualifiers so qualified references distinguish
// `c1.d.t` from `c2.d.t` and never overwrite a local table of the same name.
func (s *analyzerScope) addRelation(rte *RangeTableEntryQ, rteIdx int, columns []*Column, unknown bool) {
	if unknown {
		s.hasUnknown = true
		s.quals = append(s.quals, rteQualifier{
			catalog:  rte.Catalog,
			schema:   rte.DBName,
			name:     rte.ERef,
			hasAlias: rte.Alias != "",
			unknown:  true,
			rteIdx:   rteIdx,
			entryIdx: -1,
		})
		return
	}
	s.add(rte.ERef, rteIdx, columns)
	q := &s.quals[len(s.quals)-1]
	q.catalog = rte.Catalog
	q.schema = rte.DBName
	q.hasAlias = rte.Alias != ""
}

// soleUnknownColumnRTE returns the RTE index of a relation whose columns are
// unknown when it is the only scope entry: an unqualified column then belongs to
// it, and any other entry would make the name ambiguous.
func (s *analyzerScope) soleUnknownColumnRTE() (int, bool) {
	if len(s.quals) == 1 && s.quals[0].unknown {
		return s.quals[0].rteIdx, true
	}
	return 0, false
}

// soleUnknownColumnInChain returns the RTE index of the only relation whose
// columns are unknown, searching this scope and then the queries around it, and
// how many scopes out it sits. `SELECT (SELECT a) FROM cat.db.t` reads the
// remote relation's a through the correlation, so an unqualified name no scope
// resolves belongs to it.
func (s *analyzerScope) soleUnknownColumnInChain() (int, int, bool) {
	for cur, levels := s, 0; cur != nil; cur, levels = cur.parent, levels+1 {
		if idx, ok := cur.soleUnknownColumnRTE(); ok {
			return idx, levels, true
		}
	}
	return 0, 0, false
}

func (s *analyzerScope) hasUnknownInChain() bool {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.hasUnknown {
			return true
		}
	}
	return false
}

// resolveQualifiedRefFull resolves a table-qualified column reference
// (`t.c`, `d.t.c` or `cat.d.t.c`), honoring catalog and schema qualifiers and
// rejecting ambiguous matches across local and remote relations.
func (s *analyzerScope) resolveQualifiedRefFull(catalog, schema, tableName, colName string) (*VarExprQ, error) {
	if catalog == "" && schema == "" && !s.hasUnknownInChain() {
		rteIdx, attNum, levelsUp, err := s.resolveQualifiedColumnFull(tableName, colName)
		if err != nil {
			return nil, err
		}
		return &VarExprQ{RangeIdx: rteIdx, AttNum: attNum, LevelsUp: levelsUp}, nil
	}

	for cur, levelsUp := s, 0; cur != nil; cur, levelsUp = cur.parent, levelsUp+1 {
		var matches []rteQualifier
		for _, q := range cur.quals {
			if !strings.EqualFold(q.name, tableName) {
				continue
			}
			if catalog != "" {
				if q.hasAlias || !strings.EqualFold(q.catalog, catalog) || !strings.EqualFold(q.schema, schema) {
					continue
				}
			} else if schema != "" {
				if q.hasAlias || !strings.EqualFold(q.schema, schema) {
					continue
				}
			}
			matches = append(matches, q)
		}
		if len(matches) > 1 {
			return nil, s.unsure(&Error{
				Code:     1052,
				SQLState: "23000",
				Message:  fmt.Sprintf("Column '%s' in field list is ambiguous", colName),
			})
		}
		if len(matches) == 1 {
			m := matches[0]
			if m.unknown {
				return &VarExprQ{RangeIdx: m.rteIdx, Name: colName, LevelsUp: levelsUp}, nil
			}
			lowerCol := strings.ToLower(colName)
			for j, c := range cur.cols[m.entryIdx] {
				if strings.ToLower(c.Name) == lowerCol {
					return &VarExprQ{RangeIdx: m.rteIdx, AttNum: j + 1, LevelsUp: levelsUp}, nil
				}
			}
			return nil, s.unsure(errNoSuchColumn(colName, "field list"))
		}
	}

	return nil, s.unsure(&Error{
		Code:     ErrUnknownTable,
		SQLState: sqlState(ErrUnknownTable),
		Message:  fmt.Sprintf("Unknown table '%s'", tableName),
	})
}

func (s *analyzerScope) add(name string, rteIdx int, columns []*Column) {
	visible := make([]*Column, 0, len(columns))
	for _, c := range columns {
		if c.Hidden == ColumnHiddenNone {
			visible = append(visible, c)
		}
	}
	scopeCols := make([]scope.Column, len(visible))
	for i, c := range visible {
		scopeCols[i] = scope.Column{Name: c.Name, Position: i + 1}
	}
	entryIdx := len(s.cols)
	s.base.Add(name, &scope.Table{Name: name, Columns: scopeCols})
	s.rteMap = append(s.rteMap, rteIdx)
	s.cols = append(s.cols, visible)
	s.quals = append(s.quals, rteQualifier{
		name:     name,
		rteIdx:   rteIdx,
		entryIdx: entryIdx,
	})
}

// resolveColumn finds an unqualified column name across all scope entries.
// Returns the RTE index and 1-based attribute number.
// Error 1052 for ambiguous, 1054 for unknown.
func (s *analyzerScope) resolveColumn(colName string) (int, int, error) {
	entryIdx, pos, err := s.base.ResolveColumn(colName)
	if err != nil {
		if strings.Contains(err.Error(), "ambiguous") {
			return 0, 0, &Error{
				Code:     1052,
				SQLState: "23000",
				Message:  fmt.Sprintf("Column '%s' in field list is ambiguous", colName),
			}
		}
		return 0, 0, errNoSuchColumn(colName, "field list")
	}
	return s.rteMap[entryIdx], pos, nil
}

// resolveQualifiedColumn finds a column qualified by table name or alias.
// Returns the RTE index and 1-based attribute number.
func (s *analyzerScope) resolveQualifiedColumn(tableName, colName string) (int, int, error) {
	entryIdx, pos, err := s.base.ResolveQualifiedColumn(tableName, colName)
	if err != nil {
		if strings.Contains(err.Error(), "unknown table") {
			return 0, 0, &Error{
				Code:     ErrUnknownTable,
				SQLState: sqlState(ErrUnknownTable),
				Message:  fmt.Sprintf("Unknown table '%s'", tableName),
			}
		}
		return 0, 0, errNoSuchColumn(colName, "field list")
	}
	return s.rteMap[entryIdx], pos, nil
}

// getColumns returns the columns for a named table reference, or nil if not found.
func (s *analyzerScope) getColumns(tableName string) []*Column {
	entries := s.base.AllEntries()
	lower := strings.ToLower(tableName)
	for i, e := range entries {
		if strings.ToLower(e.Name) == lower {
			return s.cols[i]
		}
	}
	return nil
}

// resolveColumnFull resolves an unqualified column, trying parent scopes
// if not found locally. Returns (rteIdx, attNum, levelsUp, error).
func (s *analyzerScope) resolveColumnFull(colName string) (int, int, int, error) {
	rteIdx, attNum, err := s.resolveColumn(colName)
	if err == nil {
		return rteIdx, attNum, 0, nil
	}
	if s.parent != nil {
		rteIdx, attNum, parentLevels, parentErr := s.parent.resolveColumnFull(colName)
		if parentErr == nil {
			return rteIdx, attNum, parentLevels + 1, nil
		}
	}
	return 0, 0, 0, s.unsure(err)
}

// resolveQualifiedColumnFull resolves a qualified column, trying parent scopes
// if not found locally. Returns (rteIdx, attNum, levelsUp, error).
func (s *analyzerScope) resolveQualifiedColumnFull(tableName, colName string) (int, int, int, error) {
	rteIdx, attNum, err := s.resolveQualifiedColumn(tableName, colName)
	if err == nil {
		return rteIdx, attNum, 0, nil
	}
	if s.parent != nil {
		rteIdx, attNum, parentLevels, parentErr := s.parent.resolveQualifiedColumnFull(tableName, colName)
		if parentErr == nil {
			return rteIdx, attNum, parentLevels + 1, nil
		}
	}
	return 0, 0, 0, s.unsure(err)
}

// allEntries returns all scope entries in registration order.
// This returns a slice of the internal scopeEntry type for backward compatibility
// with the analyzer code that needs rteIdx and catalog column pointers.
func (s *analyzerScope) allEntries() []scopeEntry {
	entries := s.base.AllEntries()
	result := make([]scopeEntry, len(entries))
	for i, e := range entries {
		result[i] = scopeEntry{
			name:    e.Name,
			rteIdx:  s.rteMap[i],
			columns: s.cols[i],
		}
	}
	return result
}

// entryIndex returns the index of the first entry registered under exactly
// name, or -1. `SELECT t1.*, t2.*, ...` looks up one entry per item, and
// scanning (or copying, as allEntries does) every entry for each made a
// select list over a long join chain quadratic: 4.6 GB allocated at 10,000.
func (s *analyzerScope) entryIndex(name string) int {
	entries := s.base.AllEntries()
	if s.firstByName == nil {
		s.firstByName = make(map[string]int, len(entries))
	}
	for ; s.indexed < len(entries); s.indexed++ {
		if _, ok := s.firstByName[entries[s.indexed].Name]; !ok {
			s.firstByName[entries[s.indexed].Name] = s.indexed
		}
	}
	if i, ok := s.firstByName[name]; ok {
		return i
	}
	return -1
}

// entry returns entry i in the form allEntries does.
func (s *analyzerScope) entry(i int) scopeEntry {
	return scopeEntry{
		name:    s.base.AllEntries()[i].Name,
		rteIdx:  s.rteMap[i],
		columns: s.cols[i],
	}
}

// scopeEntry is one named table reference visible in the current scope.
// Kept for backward compatibility with analyzer code that accesses rteIdx and columns.
type scopeEntry struct {
	name    string    // effective reference name (alias or table name)
	rteIdx  int       // index into Query.RangeTable
	columns []*Column // columns available from this entry
}
