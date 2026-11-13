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

// markCoalesced marks a column from a table as coalesced (hidden during star expansion).
func (s *analyzerScope) markCoalesced(tableName, colName string) {
	s.base.MarkCoalesced(tableName, colName)
}

// isCoalesced returns true if the given table.column is coalesced away by USING/NATURAL.
func (s *analyzerScope) isCoalesced(tableName, colName string) bool {
	return s.base.IsCoalesced(tableName, colName)
}

// add registers a table reference in the scope.
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
	s.base.Add(name, &scope.Table{Name: name, Columns: scopeCols})
	s.rteMap = append(s.rteMap, rteIdx)
	s.cols = append(s.cols, visible)
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
