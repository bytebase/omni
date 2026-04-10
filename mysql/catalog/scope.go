package catalog

import (
	"fmt"
	"strings"
)

// analyzerScope tracks the set of visible table/view references during
// semantic analysis, supporting column resolution by name.
type analyzerScope struct {
	entries       []scopeEntry
	byName        map[string]int    // lowered name -> index into entries
	coalescedCols map[string]bool   // "tablename.colname" (lowered) -> true; columns hidden by USING/NATURAL
}

// scopeEntry is one named table reference visible in the current scope.
type scopeEntry struct {
	name    string    // lowered effective reference name (alias or table name)
	rteIdx  int       // index into Query.RangeTable
	columns []*Column // columns available from this entry
}

func newScope() *analyzerScope {
	return &analyzerScope{
		byName:        make(map[string]int),
		coalescedCols: make(map[string]bool),
	}
}

// markCoalesced marks a column from a table as coalesced (hidden during star expansion).
func (s *analyzerScope) markCoalesced(tableName, colName string) {
	key := strings.ToLower(tableName) + "." + strings.ToLower(colName)
	s.coalescedCols[key] = true
}

// isCoalesced returns true if the given table.column is coalesced away by USING/NATURAL.
func (s *analyzerScope) isCoalesced(tableName, colName string) bool {
	key := strings.ToLower(tableName) + "." + strings.ToLower(colName)
	return s.coalescedCols[key]
}

// add registers a table reference in the scope.
func (s *analyzerScope) add(name string, rteIdx int, columns []*Column) {
	lower := strings.ToLower(name)
	s.byName[lower] = len(s.entries)
	s.entries = append(s.entries, scopeEntry{
		name:    lower,
		rteIdx:  rteIdx,
		columns: columns,
	})
}

// resolveColumn finds an unqualified column name across all scope entries.
// Returns the RTE index and 1-based attribute number.
// Error 1052 for ambiguous, 1054 for unknown.
func (s *analyzerScope) resolveColumn(colName string) (int, int, error) {
	lower := strings.ToLower(colName)
	var foundRTE, foundAtt int
	found := 0
	for _, e := range s.entries {
		for i, c := range e.columns {
			if strings.ToLower(c.Name) == lower {
				found++
				foundRTE = e.rteIdx
				foundAtt = i + 1 // 1-based
				if found > 1 {
					return 0, 0, &Error{
						Code:     1052,
						SQLState: "23000",
						Message:  fmt.Sprintf("Column '%s' in field list is ambiguous", colName),
					}
				}
			}
		}
	}
	if found == 0 {
		return 0, 0, errNoSuchColumn(colName, "field list")
	}
	return foundRTE, foundAtt, nil
}

// resolveQualifiedColumn finds a column qualified by table name or alias.
// Returns the RTE index and 1-based attribute number.
func (s *analyzerScope) resolveQualifiedColumn(tableName, colName string) (int, int, error) {
	lowerTable := strings.ToLower(tableName)
	idx, ok := s.byName[lowerTable]
	if !ok {
		return 0, 0, &Error{
			Code:     ErrUnknownTable,
			SQLState: sqlState(ErrUnknownTable),
			Message:  fmt.Sprintf("Unknown table '%s'", tableName),
		}
	}
	e := s.entries[idx]
	lowerCol := strings.ToLower(colName)
	for i, c := range e.columns {
		if strings.ToLower(c.Name) == lowerCol {
			return e.rteIdx, i + 1, nil
		}
	}
	return 0, 0, errNoSuchColumn(colName, "field list")
}

// getColumns returns the columns for a named table reference, or nil if not found.
func (s *analyzerScope) getColumns(tableName string) []*Column {
	idx, ok := s.byName[strings.ToLower(tableName)]
	if !ok {
		return nil
	}
	return s.entries[idx].columns
}

// allEntries returns all scope entries in registration order.
func (s *analyzerScope) allEntries() []scopeEntry {
	return s.entries
}
