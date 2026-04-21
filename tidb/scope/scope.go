// Package scope provides a shared, dependency-free scope data structure for
// tracking visible table references and resolving column names.  It is used by
// both the semantic analyzer (mysql/catalog) and the AST resolver/deparser
// (mysql/deparse).
package scope

import (
	"fmt"
	"strings"
)

// Column represents a column visible from a scope entry.
type Column struct {
	Name     string
	Position int // 1-based position in the table
}

// Table represents a table or virtual table in scope with its visible columns.
type Table struct {
	Name    string
	Columns []Column
}

// GetColumn returns a column by name (case-insensitive), or nil.
func (t *Table) GetColumn(name string) *Column {
	lower := strings.ToLower(name)
	for i := range t.Columns {
		if strings.ToLower(t.Columns[i].Name) == lower {
			return &t.Columns[i]
		}
	}
	return nil
}

// Entry is one named table reference visible in the current scope.
type Entry struct {
	Name  string // effective reference name (alias or table name), as registered
	Table *Table
}

// Scope tracks visible table references for column resolution.
type Scope struct {
	entries       []Entry
	byName        map[string]int  // lowered name -> index into entries
	coalescedCols map[string]bool // "tablename.colname" (lowered) -> hidden by USING/NATURAL
}

// New creates an empty scope.
func New() *Scope {
	return &Scope{
		byName:        make(map[string]int),
		coalescedCols: make(map[string]bool),
	}
}

// Add registers a table reference in the scope.
func (s *Scope) Add(name string, table *Table) {
	lower := strings.ToLower(name)
	s.byName[lower] = len(s.entries)
	s.entries = append(s.entries, Entry{Name: name, Table: table})
}

// ResolveColumn finds an unqualified column name across all scope entries.
// Returns (entry index, 1-based column position, error).
// Returns error if not found or ambiguous.
func (s *Scope) ResolveColumn(colName string) (int, int, error) {
	lower := strings.ToLower(colName)
	foundIdx := -1
	foundPos := 0
	for i, e := range s.entries {
		for j, c := range e.Table.Columns {
			if strings.ToLower(c.Name) == lower {
				if foundIdx >= 0 {
					return 0, 0, fmt.Errorf("column '%s' is ambiguous", colName)
				}
				foundIdx = i
				foundPos = j + 1 // 1-based
			}
		}
	}
	if foundIdx < 0 {
		return 0, 0, fmt.Errorf("unknown column '%s'", colName)
	}
	return foundIdx, foundPos, nil
}

// ResolveQualifiedColumn finds a column qualified by table name.
// Returns (entry index, 1-based column position, error).
func (s *Scope) ResolveQualifiedColumn(tableName, colName string) (int, int, error) {
	lowerTable := strings.ToLower(tableName)
	idx, ok := s.byName[lowerTable]
	if !ok {
		return 0, 0, fmt.Errorf("unknown table '%s'", tableName)
	}
	e := s.entries[idx]
	lowerCol := strings.ToLower(colName)
	for j, c := range e.Table.Columns {
		if strings.ToLower(c.Name) == lowerCol {
			return idx, j + 1, nil
		}
	}
	return 0, 0, fmt.Errorf("unknown column '%s.%s'", tableName, colName)
}

// GetTable returns the table for a named entry, or nil.
func (s *Scope) GetTable(name string) *Table {
	idx, ok := s.byName[strings.ToLower(name)]
	if !ok {
		return nil
	}
	return s.entries[idx].Table
}

// AllEntries returns all scope entries in registration order.
func (s *Scope) AllEntries() []Entry {
	return s.entries
}

// MarkCoalesced marks a column from a table as hidden by USING/NATURAL.
func (s *Scope) MarkCoalesced(tableName, colName string) {
	key := strings.ToLower(tableName) + "." + strings.ToLower(colName)
	s.coalescedCols[key] = true
}

// IsCoalesced returns true if the column is hidden by USING/NATURAL.
func (s *Scope) IsCoalesced(tableName, colName string) bool {
	key := strings.ToLower(tableName) + "." + strings.ToLower(colName)
	return s.coalescedCols[key]
}
