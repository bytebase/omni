package catalog

import (
	"sort"
	"strconv"
)

// Diff compares two catalogs, a desired target against the current schema,
// and reports the differences as a SchemaDiff; when the two are equal the diff
// is empty and no DDL is generated. Identity is the lower-cased (database,
// object name).
//
// What compares equal is what TiDB stores. The catalog already holds each
// object in that form, since Exec and LoadMetadata resolve types, defaults,
// charsets and keys the way the server does, and the SHOW CREATE renderers
// write the result as the server writes it (held to it by the container,
// simulate.tsv and LoadMetadata tests). So a column, a key, a foreign key, a
// partitioning, a table's options, a view and a sequence are compared by the
// text SHOW CREATE gives them. TiDB has no stored routines, triggers or
// events, so there is nothing to compare for them.

// DiffAction describes what happened to an object between two catalog states.
type DiffAction int

const (
	// DiffAdd means the object exists in `to` but not in `from`.
	DiffAdd DiffAction = iota + 1
	// DiffDrop means the object exists in `from` but not in `to`.
	DiffDrop
	// DiffModify means the object exists in both but has changed.
	DiffModify
)

// String renders a DiffAction for diagnostics.
func (a DiffAction) String() string {
	switch a {
	case DiffAdd:
		return "ADD"
	case DiffDrop:
		return "DROP"
	case DiffModify:
		return "MODIFY"
	default:
		return "UNKNOWN"
	}
}

// TableDiffEntry describes a table that was added, dropped, or modified, with
// the per-table sub-object diffs of a modified one. OptionsChanged reports a
// change to the table options SHOW CREATE TABLE ends the table with, TTL
// included.
type TableDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Table
	To       *Table

	Columns          []ColumnDiffEntry
	Indexes          []IndexDiffEntry
	ForeignKeys      []ForeignKeyDiffEntry
	Checks           []CheckDiffEntry
	PartitionChanged bool
	OptionsChanged   bool
}

// ColumnDiffEntry describes a column change within a table. Identity is the
// lower-cased column name.
type ColumnDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Column
	To     *Column
}

// IndexDiffEntry describes a key change within a table, the primary key
// (named PRIMARY) included.
type IndexDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Index
	To     *Index
}

// ForeignKeyDiffEntry describes a foreign-key change within a table.
type ForeignKeyDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// CheckDiffEntry describes a CHECK constraint change within a table.
type CheckDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// ViewDiffEntry describes a view change.
type ViewDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *View
	To       *View
}

// SequenceDiffEntry describes a sequence change.
type SequenceDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Sequence
	To       *Sequence
}

// SchemaDiff holds all differences between two catalog states, each list
// sorted by (database, name, action).
type SchemaDiff struct {
	Tables    []TableDiffEntry
	Views     []ViewDiffEntry
	Sequences []SequenceDiffEntry
}

// IsEmpty reports whether there are no differences.
func (d *SchemaDiff) IsEmpty() bool {
	return len(d.Tables) == 0 && len(d.Views) == 0 && len(d.Sequences) == 0
}

// Diff compares two catalog states (from = current, to = target) and returns
// all differences.
func Diff(from, to *Catalog) *SchemaDiff {
	return &SchemaDiff{
		Tables:    diffTables(from, to),
		Views:     diffViews(from, to),
		Sequences: diffSequences(from, to),
	}
}

// objectKey is the identity of an object: (database, name), lower-cased. A
// table's sub-objects leave the database empty.
type objectKey struct {
	db   string
	name string
}

// diffObjects compares two maps of one kind of object by identity: an object
// only in from is dropped, one only in to is added, and one in both is
// modified when changed says so. entry builds the result entry, given the
// object on each side (the zero value on the side it is missing from), and
// sortKey orders the result.
func diffObjects[T any, E any](from, to map[objectKey]T, changed func(a, b T) bool, entry func(DiffAction, objectKey, T, T) E, sortKey func(E) (string, string, DiffAction)) []E {
	var zero T
	var out []E
	for key, f := range from {
		if _, ok := to[key]; !ok {
			out = append(out, entry(DiffDrop, key, f, zero))
		}
	}
	for key, t := range to {
		f, ok := from[key]
		switch {
		case !ok:
			out = append(out, entry(DiffAdd, key, zero, t))
		case changed(f, t):
			out = append(out, entry(DiffModify, key, f, t))
		}
	}
	sort.Slice(out, func(i, j int) bool {
		di, ni, ai := sortKey(out[i])
		dj, nj, aj := sortKey(out[j])
		if di != dj {
			return di < dj
		}
		if ni != nj {
			return ni < nj
		}
		return ai < aj
	})
	return out
}

// diffNamed is diffObjects over a table's sub-objects, keyed by lower-cased
// name.
func diffNamed[T any, E any](from, to map[string]T, changed func(a, b T) bool, entry func(DiffAction, T, T) E, sortKey func(E) (string, DiffAction)) []E {
	f := make(map[objectKey]T, len(from))
	for k, v := range from {
		f[objectKey{name: k}] = v
	}
	t := make(map[objectKey]T, len(to))
	for k, v := range to {
		t[objectKey{name: k}] = v
	}
	return diffObjects(f, t, changed,
		func(a DiffAction, _ objectKey, x, y T) E { return entry(a, x, y) },
		func(e E) (string, string, DiffAction) {
			n, a := sortKey(e)
			return "", toLower(n), a
		})
}

// encodeKeyFields joins name/value pairs into one comparison key that no two
// different lists can share: each value is written with its length.
func encodeKeyFields(pairs ...string) string {
	var b []byte
	for i := 0; i+1 < len(pairs); i += 2 {
		b = append(b, pairs[i]...)
		b = append(b, ':')
		b = strconv.AppendInt(b, int64(len(pairs[i+1])), 10)
		b = append(b, ':')
		b = append(b, pairs[i+1]...)
		b = append(b, '|')
	}
	return string(b)
}
