package catalog

// MySQL declarative-schema (SDL) diff engine — tables + columns.
//
// This is the MySQL analog of the PostgreSQL reference at pg/catalog/diff.go. It
// compares two in-memory catalogs (a desired TARGET against the synced CURRENT
// schema read back via SHOW CREATE TABLE) and reports the structural differences as
// a SchemaDiff. generate-core turns a SchemaDiff into DDL; the defining correctness
// property is idempotence: when target == current, Diff is empty (IsEmpty), so the
// release emits no DDL.
//
// Shape note (vs PG): PostgreSQL is OID/namespace-keyed and threads relOID/SchemaName
// through every differ. MySQL has no OIDs and no schema/namespace layer — "schema" IS
// the database — so identity here is (database, table, column) by lower-cased name.
// The PG *shape* transfers (a thin top-level Diff dispatcher calling one diffX per
// object kind, packing the results into SchemaDiff), but the plumbing is re-expressed
// in database.table terms.
//
// Every canonicalization-sensitive comparison is routed through normalize.go (the
// single owner of MySQL's stored-form canonical key). Column equality is decided by
// Normalizer.CanonicalColumn; table-level engine/charset/collation comparisons go
// through CanonicalEngine / foldCharset / effectiveTableCollation. The table-level
// AUTO_INCREMENT=N counter is never compared (IgnoreTableAutoIncrement is constant true);
// ROW_FORMAT is compared only when normalize.go's IgnoreRowFormat says it is significant.
// diff-core never re-implements canonicalization.
//
// Scope of THIS node (omni:diff-core): the dispatcher, the SchemaDiff aggregate, the
// table differ (diff_table.go), and the column differ (diff_column.go). The breadth
// object kinds (indexes, foreign keys, constraints, checks, views, triggers, routines,
// events, partitions) have their extension points wired here — empty diff sections and
// per-kind seams — but are populated by later nodes, mirroring PG's diff_<obj>.go layout.

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

// TableDiffEntry describes a table that was added, dropped, or modified. It is the
// MySQL analog of PG's RelationDiffEntry: it carries the per-table sub-object diffs.
// Columns is populated by this node (diff-core). Indexes / Constraints / ForeignKeys
// / Checks / PartitionChanged are extension points for the breadth nodes — they are
// always present (empty) so generate-core and the breadth differs can append into a
// stable shape without a struct change.
type TableDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Table
	To       *Table

	// Populated by diff-core.
	Columns []ColumnDiffEntry

	// Extension points for breadth nodes (left empty by diff-core). A modified table
	// is reported whenever ANY of these is non-empty, so wiring them here means a
	// breadth node only has to fill its own slice — compareTable already folds them
	// into the change decision via tableSubdiffsChanged.
	Indexes          []IndexDiffEntry
	Constraints      []ConstraintDiffEntry
	ForeignKeys      []ForeignKeyDiffEntry
	Checks           []CheckDiffEntry
	PartitionChanged bool
}

// ColumnDiffEntry describes a column change within a table. Identity is the column
// name (lower-cased); equality is decided by the normalize-core CanonicalColumn key.
type ColumnDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Column
	To     *Column
}

// IndexDiffEntry describes a secondary-index change within a table.
// Extension point: populated by the index breadth node, not diff-core.
type IndexDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Index
	To     *Index
}

// ConstraintDiffEntry describes a PRIMARY KEY / UNIQUE constraint change.
// Extension point: populated by a breadth node, not diff-core.
type ConstraintDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// ForeignKeyDiffEntry describes a foreign-key change within a table.
// Extension point: populated by the FK breadth node, not diff-core.
type ForeignKeyDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// CheckDiffEntry describes a CHECK-constraint change within a table (8.0 only — on
// 5.7 CHECK is unrepresentable; see normalize.go CheckSupported).
// Extension point: populated by a breadth node, not diff-core.
type CheckDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// ViewDiffEntry describes a view change.
// Extension point: populated by the view breadth node, not diff-core.
type ViewDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *View
	To       *View
}

// RoutineDiffEntry describes a stored function or procedure change.
// Extension point: populated by the routine breadth node, not diff-core.
type RoutineDiffEntry struct {
	Action      DiffAction
	Database    string
	Name        string
	IsProcedure bool
	From        *Routine
	To          *Routine
}

// TriggerDiffEntry describes a trigger change.
// Extension point: populated by the trigger breadth node, not diff-core.
type TriggerDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Trigger
	To       *Trigger
}

// EventDiffEntry describes a scheduled-event change.
// Extension point: populated by the event breadth node, not diff-core.
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

// SchemaDiff holds all differences between two catalog states. It is intentionally
// extensible: diff-core populates Tables (and the per-table Columns inside each
// entry); every other slice is an extension point a breadth node fills later. The
// dispatcher (Diff) calls one diffX per object kind, so a breadth node adds its kind
// by implementing diff_<obj>.go and appending here — no other node changes.
type SchemaDiff struct {
	// Populated by diff-core.
	Tables []TableDiffEntry

	// Extension points for breadth nodes (database-level objects).
	Views      []ViewDiffEntry
	Functions  []RoutineDiffEntry
	Procedures []RoutineDiffEntry
	Triggers   []TriggerDiffEntry
	Events     []EventDiffEntry
}

// IsEmpty reports whether there are no differences. This is the idempotence gate read
// by bytebase's DiffSDLMigration: an empty diff means the release emits no DDL. A
// non-empty diff on a no-op release is a normalization bug, not a tolerable nit.
func (d *SchemaDiff) IsEmpty() bool {
	return len(d.Tables) == 0 &&
		len(d.Views) == 0 &&
		len(d.Functions) == 0 &&
		len(d.Procedures) == 0 &&
		len(d.Triggers) == 0 &&
		len(d.Events) == 0
}

// ---------------------------------------------------------------------------
// Top-level dispatcher
// ---------------------------------------------------------------------------

// Diff compares two catalog states (from = old/current, to = new/target) and returns
// all differences, using the modern MySQL 8.0 canonical form by default. The
// explicit_defaults_for_timestamp value is taken from the `to` catalog's session when
// it carries one (the synced schema records the source server's effective setting),
// so bare-TIMESTAMP nullability canonicalizes correctly without a version override.
//
// Diff is the contract entry point (pg/catalog/diff.go:210). Callers that know the
// target server version — bytebase's mysqlDiffSDLMigration and the oracle tests, where
// 5.7-vs-8.0 stored forms genuinely diverge — should use DiffWithNormalizer to select
// the version whose stored form the canonicalizer must reproduce.
func Diff(from, to *Catalog) *SchemaDiff {
	return DiffWithNormalizer(from, to, defaultNormalizer(to))
}

// DiffWithNormalizer compares two catalogs using an explicit Normalizer, which fixes
// the target MySQL version (and explicit_defaults_for_timestamp) whose stored form the
// canonical comparison must reproduce. This is the version-aware workhorse; Diff
// delegates to it with a default normalizer. A nil normalizer falls back to the
// default for `to`.
//
// The dispatcher mirrors PG: one diffX per object kind, packed into SchemaDiff. Only
// the table differ (which descends into columns) is implemented in diff-core; the
// breadth slices are left nil for later nodes.
func DiffWithNormalizer(from, to *Catalog, n *Normalizer) *SchemaDiff {
	if n == nil {
		n = defaultNormalizer(to)
	}
	return &SchemaDiff{
		Tables: diffTables(from, to, n),
		// Breadth extension points — populated by later nodes:
		//   Views:      diffViews(from, to, n),
		//   Functions:  diffFunctions(from, to, n),
		//   Procedures: diffProcedures(from, to, n),
		//   Triggers:   diffTriggers(from, to, n),
		//   Events:     diffEvents(from, to, n),
	}
}

// defaultNormalizer builds the Normalizer used by the parameterless Diff: MySQL 8.0
// stored form, seeded from the `to` catalog's captured explicit_defaults_for_timestamp
// when available. Using `to` (the target/desired schema, which on the release path is
// re-loaded from the source's synced metadata) keeps the EDFT value aligned with the
// server whose stored form we must match.
func defaultNormalizer(to *Catalog) *Normalizer {
	n := NormalizerFor(MySQL80)
	if to != nil {
		n.ExplicitDefaultsForTimestamp = to.session.ExplicitDefaultsForTimestamp
	}
	return n
}
