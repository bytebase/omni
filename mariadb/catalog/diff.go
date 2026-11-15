package catalog

// The SDL diff compares two catalogs, a desired target against the synced
// current schema, and reports the differences as a SchemaDiff. When the two are
// equal the diff is empty (IsEmpty) and no DDL is generated. Identity is the
// lower-cased (database, object name); the Normalizer (normalize.go) decides what
// compares equal. Primary and unique keys are compared as indexes, so
// diffConstraints reports nothing.

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

// ---------------------------------------------------------------------------
// Per-object diff entry types
// ---------------------------------------------------------------------------

// TableDiffEntry describes a table that was added, dropped, or modified, with
// the per-table sub-object diffs. Constraints is always empty: primary and unique
// keys are compared as indexes.
type TableDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Table
	To       *Table

	Columns []ColumnDiffEntry

	Indexes          []IndexDiffEntry
	Constraints      []ConstraintDiffEntry
	ForeignKeys      []ForeignKeyDiffEntry
	Checks           []CheckDiffEntry
	PartitionChanged bool
}

// ColumnDiffEntry describes a column change within a table. Identity is the
// lower-cased column name.
type ColumnDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Column
	To     *Column
}

// IndexDiffEntry describes an index change within a table.
type IndexDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Index
	To     *Index
}

// ConstraintDiffEntry is a constraint change within a table. It is never
// populated; see TableDiffEntry.
type ConstraintDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// ForeignKeyDiffEntry describes a foreign-key change within a table.
type ForeignKeyDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// CheckDiffEntry describes a table-level CHECK change within a table.
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

// RoutineDiffEntry describes a stored function or procedure change.
type RoutineDiffEntry struct {
	Action      DiffAction
	Database    string
	Name        string
	IsProcedure bool
	From        *Routine
	To          *Routine
}

// TriggerDiffEntry describes a trigger change.
type TriggerDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Trigger
	To       *Trigger
}

// EventDiffEntry describes a scheduled-event change.
type EventDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Event
	To       *Event
}

// ---------------------------------------------------------------------------
// SchemaDiff — aggregate result of comparing two catalogs
// ---------------------------------------------------------------------------

// SchemaDiff holds all differences between two catalog states.
type SchemaDiff struct {
	Tables []TableDiffEntry

	Views      []ViewDiffEntry
	Functions  []RoutineDiffEntry
	Procedures []RoutineDiffEntry
	Triggers   []TriggerDiffEntry
	Events     []EventDiffEntry
	Sequences  []SequenceDiffEntry
}

// IsEmpty reports whether there are no differences. An empty diff means the
// release emits no DDL.
func (d *SchemaDiff) IsEmpty() bool {
	return len(d.Tables) == 0 &&
		len(d.Views) == 0 &&
		len(d.Functions) == 0 &&
		len(d.Procedures) == 0 &&
		len(d.Triggers) == 0 &&
		len(d.Events) == 0 &&
		len(d.Sequences) == 0
}

// ---------------------------------------------------------------------------
// Top-level dispatcher
// ---------------------------------------------------------------------------

// Diff compares two catalog states (from = old/current, to = new/target) and
// returns all differences. The version and explicit_defaults_for_timestamp come
// from the `to` catalog's session.
func Diff(from, to *Catalog) *SchemaDiff {
	return DiffWithNormalizer(from, to, defaultNormalizer(to))
}

// DiffWithNormalizer compares two catalogs using an explicit Normalizer. A nil
// normalizer falls back to the default for `to`.
func DiffWithNormalizer(from, to *Catalog, n *Normalizer) *SchemaDiff {
	if n == nil {
		n = defaultNormalizer(to)
	}
	return &SchemaDiff{
		Tables: diffTables(from, to, n),
		// Database-level objects.
		Views:      diffViews(from, to, n),
		Functions:  diffFunctions(from, to, n),
		Procedures: diffProcedures(from, to, n),
		Triggers:   diffTriggers(from, to, n),
		Events:     diffEvents(from, to, n),
		Sequences:  diffSequences(from, to),
	}
}

// defaultNormalizer builds the Normalizer for Diff from the `to` catalog's
// session version and explicit_defaults_for_timestamp. A nil `to` gets the
// MariaDB118 defaults.
func defaultNormalizer(to *Catalog) *Normalizer {
	if to == nil {
		return NormalizerFor(MariaDB118)
	}
	n := NormalizerFor(to.session.Version)
	n.ExplicitDefaultsForTimestamp = to.session.ExplicitDefaultsForTimestamp
	return n
}
