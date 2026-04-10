package catalog

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/mysql/scope"
)

// analyzerScope wraps scope.Scope and adds analyzer-specific features:
// parent chain for correlated subqueries, and rteIdx mapping.
type analyzerScope struct {
	base   *scope.Scope
	rteMap []int       // parallel to base entries: entry index -> RTE index in Query.RangeTable
	cols   [][]*Column // parallel to base entries: entry index -> catalog Column pointers
	parent *analyzerScope
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
	scopeCols := make([]scope.Column, len(columns))
	for i, c := range columns {
		scopeCols[i] = scope.Column{Name: c.Name, Position: i + 1}
	}
	s.base.Add(name, &scope.Table{Name: name, Columns: scopeCols})
	s.rteMap = append(s.rteMap, rteIdx)
	s.cols = append(s.cols, columns)
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
	return 0, 0, 0, err
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
	return 0, 0, 0, err
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

// scopeEntry is one named table reference visible in the current scope.
// Kept for backward compatibility with analyzer code that accesses rteIdx and columns.
type scopeEntry struct {
	name    string    // effective reference name (alias or table name)
	rteIdx  int       // index into Query.RangeTable
	columns []*Column // columns available from this entry
}
