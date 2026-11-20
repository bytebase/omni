package catalog

import "sort"

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
	Action          DiffAction
	SchemaName      string
	Name            string
	From            *Relation
	To              *Relation
	Columns         []ColumnDiffEntry
	Constraints     []ConstraintDiffEntry
	Indexes         []IndexDiffEntry
	Triggers        []TriggerDiffEntry
	Policies        []PolicyDiffEntry
	Rules           []RuleDiffEntry
	RLSChanged      bool
	RLSEnabled      bool
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

// EventTriggerDiffEntry describes a database-level event trigger change.
type EventTriggerDiffEntry struct {
	Action DiffAction
	Name   string
	From   *EventTrigger
	To     *EventTrigger
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

// CompositeTypeDiffEntry describes a composite type change.
type CompositeTypeDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	From       *Relation // RelKind='c'
	To         *Relation // RelKind='c'
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

// RuleDiffEntry describes a rewrite rule change within a relation.
type RuleDiffEntry struct {
	Action DiffAction
	Name   string
	From   *Rule
	To     *Rule
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

// CollationDiffEntry describes a user collation change.
type CollationDiffEntry struct {
	Action     DiffAction
	SchemaName string
	Name       string
	From       *UserCollation
	To         *UserCollation
}

// AccessMethodDiffEntry describes a user-defined access method change.
type AccessMethodDiffEntry struct {
	Action DiffAction
	Name   string
	From   *AccessMethod
	To     *AccessMethod
}

// ---------------------------------------------------------------------------
// SchemaDiff — aggregate result of comparing two catalogs
// ---------------------------------------------------------------------------

// SchemaDiff holds all differences between two catalog states.
type SchemaDiff struct {
	Schemas        []SchemaDiffEntry
	Collations     []CollationDiffEntry
	AccessMethods  []AccessMethodDiffEntry
	Relations      []RelationDiffEntry
	Sequences      []SequenceDiffEntry
	Functions      []FunctionDiffEntry
	Enums          []EnumDiffEntry
	Domains        []DomainDiffEntry
	Ranges         []RangeDiffEntry
	CompositeTypes []CompositeTypeDiffEntry
	Extensions     []ExtensionDiffEntry
	EventTriggers  []EventTriggerDiffEntry
	Comments       []CommentDiffEntry
	Grants         []GrantDiffEntry
}

// IsEmpty returns true if there are no differences.
func (d *SchemaDiff) IsEmpty() bool {
	return len(d.Schemas) == 0 &&
		len(d.Collations) == 0 &&
		len(d.AccessMethods) == 0 &&
		len(d.Relations) == 0 &&
		len(d.Sequences) == 0 &&
		len(d.Functions) == 0 &&
		len(d.Enums) == 0 &&
		len(d.Domains) == 0 &&
		len(d.Ranges) == 0 &&
		len(d.CompositeTypes) == 0 &&
		len(d.Extensions) == 0 &&
		len(d.EventTriggers) == 0 &&
		len(d.Comments) == 0 &&
		len(d.Grants) == 0
}

// Diff compares two catalog states and returns all differences.
// The from catalog represents the old state and to represents the new state.
func Diff(from, to *Catalog) *SchemaDiff {
	relations := diffRelations(from, to)
	relations = append(relations, diffViews(from, to)...)
	return &SchemaDiff{
		Schemas:        diffSchemas(from, to),
		Collations:     diffCollations(from, to),
		AccessMethods:  diffAccessMethods(from, to),
		Relations:      relations,
		Sequences:      diffSequences(from, to),
		Functions:      diffFunctions(from, to),
		Enums:          diffEnums(from, to),
		Domains:        diffDomains(from, to),
		Ranges:         diffRanges(from, to),
		CompositeTypes: diffComposites(from, to),
		Extensions:     diffExtensions(from, to),
		EventTriggers:  diffEventTriggers(from, to),
		Comments:       diffComments(from, to),
		Grants:         diffGrants(from, to),
	}
}

func diffAccessMethods(from, to *Catalog) []AccessMethodDiffEntry {
	fromMap := map[string]*AccessMethod{}
	if from != nil {
		for _, am := range from.accessMethods {
			if am != nil && am.OID >= FirstNormalObjectId && !from.isExtensionMember(am.OID) {
				fromMap[am.Name] = am
			}
		}
	}
	toMap := map[string]*AccessMethod{}
	if to != nil {
		for _, am := range to.accessMethods {
			if am != nil && am.OID >= FirstNormalObjectId && !to.isExtensionMember(am.OID) {
				toMap[am.Name] = am
			}
		}
	}
	var out []AccessMethodDiffEntry
	for name, fam := range fromMap {
		if _, ok := toMap[name]; !ok {
			out = append(out, AccessMethodDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fam,
			})
		}
	}
	for name, tam := range toMap {
		fam, ok := fromMap[name]
		if !ok {
			out = append(out, AccessMethodDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     tam,
			})
			continue
		}
		if fam.Type != tam.Type || formatAccessMethodHandler(from, fam.Handler) != formatAccessMethodHandler(to, tam.Handler) {
			out = append(out, AccessMethodDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fam,
				To:     tam,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Action < out[j].Action
	})
	return out
}

func diffCollations(from, to *Catalog) []CollationDiffEntry {
	type key struct {
		schema string
		name   string
	}
	fromMap := map[key]*UserCollation{}
	for _, s := range from.schemas {
		if s.OID == PGCatalogNamespace || s.OID == PGToastNamespace {
			continue
		}
		for name, uc := range s.Collations {
			if uc != nil && !from.isExtensionMember(uc.OID) {
				fromMap[key{s.Name, name}] = uc
			}
		}
	}
	toMap := map[key]*UserCollation{}
	for _, s := range to.schemas {
		if s.OID == PGCatalogNamespace || s.OID == PGToastNamespace {
			continue
		}
		for name, uc := range s.Collations {
			if uc != nil && !to.isExtensionMember(uc.OID) {
				toMap[key{s.Name, name}] = uc
			}
		}
	}
	var out []CollationDiffEntry
	for k, fc := range fromMap {
		if _, ok := toMap[k]; !ok {
			out = append(out, CollationDiffEntry{
				Action:     DiffDrop,
				SchemaName: k.schema,
				Name:       k.name,
				From:       fc,
			})
		}
	}
	for k, tc := range toMap {
		fc, ok := fromMap[k]
		if !ok {
			out = append(out, CollationDiffEntry{
				Action:     DiffAdd,
				SchemaName: k.schema,
				Name:       k.name,
				To:         tc,
			})
			continue
		}
		if fc.DefClause != tc.DefClause || fc.OwnerName != tc.OwnerName {
			out = append(out, CollationDiffEntry{
				Action:     DiffModify,
				SchemaName: k.schema,
				Name:       k.name,
				From:       fc,
				To:         tc,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].SchemaName != out[j].SchemaName {
			return out[i].SchemaName < out[j].SchemaName
		}
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Action < out[j].Action
	})
	return out
}
