package catalog

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

// ---------------------------------------------------------------------------
// Per-object diff entry types
// ---------------------------------------------------------------------------

// SchemaDiffEntry describes a schema that was added, dropped, or modified.
type SchemaDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Schema
	To     *Schema
}

// RelationDiffEntry describes a relation (table/view/matview) change.
type RelationDiffEntry struct {
	Action      DiffAction
	SchemaName  string
	Name        string
	From        *Relation
	To          *Relation
	Columns     []ColumnDiffEntry
	Constraints []ConstraintDiffEntry
	Indexes     []IndexDiffEntry
	Triggers    []TriggerDiffEntry
	Policies    []PolicyDiffEntry
	RLSChanged  bool
	RLSEnabled  bool
	ForceRLSEnabled bool
}

// ColumnDiffEntry describes a column change within a relation.
type ColumnDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Column
	To     *Column
}

// ConstraintDiffEntry describes a constraint change within a relation.
type ConstraintDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Constraint
	To     *Constraint
}

// IndexDiffEntry describes an index change within a relation.
type IndexDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Index
	To     *Index
}

// SequenceDiffEntry describes a standalone sequence change.
type SequenceDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	From       *Sequence
	To         *Sequence
}

// FunctionDiffEntry describes a function or procedure change.
type FunctionDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	Identity   string
	From       *UserProc
	To         *UserProc
}

// TriggerDiffEntry describes a trigger change within a relation.
type TriggerDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Trigger
	To     *Trigger
}

// EnumDiffEntry describes an enum type change.
type EnumDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	FromValues []string
	ToValues   []string
}

// DomainDiffEntry describes a domain type change.
type DomainDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	From       *DomainType
	To         *DomainType
}

// RangeDiffEntry describes a range type change.
type RangeDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	From       *RangeType
	To         *RangeType
}

// ExtensionDiffEntry describes an extension change.
type ExtensionDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Extension
	To     *Extension
}

// PolicyDiffEntry describes a policy change within a relation.
type PolicyDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Policy
	To     *Policy
}

// CommentDiffEntry describes a comment change.
type CommentDiffEntry struct {
	Action         DiffAction
	ObjType        byte
	ObjDescription string
	SubID          int16
	From           string
	To             string
}

// GrantDiffEntry describes a grant change.
type GrantDiffEntry struct {
	Action DiffAction
	From   Grant
	To     Grant
}

// ---------------------------------------------------------------------------
// SchemaDiff — aggregate result of comparing two catalogs
// ---------------------------------------------------------------------------

// SchemaDiff holds all differences between two catalog states.
type SchemaDiff struct {
	Schemas    []SchemaDiffEntry
	Relations  []RelationDiffEntry
	Sequences  []SequenceDiffEntry
	Functions  []FunctionDiffEntry
	Enums      []EnumDiffEntry
	Domains    []DomainDiffEntry
	Ranges     []RangeDiffEntry
	Extensions []ExtensionDiffEntry
	Comments   []CommentDiffEntry
	Grants     []GrantDiffEntry
}

// IsEmpty returns true if there are no differences.
func (d *SchemaDiff) IsEmpty() bool {
	return len(d.Schemas) == 0 &&
		len(d.Relations) == 0 &&
		len(d.Sequences) == 0 &&
		len(d.Functions) == 0 &&
		len(d.Enums) == 0 &&
		len(d.Domains) == 0 &&
		len(d.Ranges) == 0 &&
		len(d.Extensions) == 0 &&
		len(d.Comments) == 0 &&
		len(d.Grants) == 0
}

// Diff compares two catalog states and returns all differences.
// The from catalog represents the old state and to represents the new state.
func Diff(from, to *Catalog) *SchemaDiff {
	return &SchemaDiff{
		Schemas:   diffSchemas(from, to),
		Relations: diffRelations(from, to),
	}
}
