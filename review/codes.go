package review

import (
	"fmt"
	"slices"
	"strings"
)

// FindingCode names the condition a finding reports, the same on every
// engine: a two-letter prefix for its rule and three digits, such as
// "WT101". A code is assigned by how a reader handles the condition, so
// one code may cover several object kinds or causes; the message names
// which. A published code is never renumbered or reused, and a code no
// longer reported stays in the table.
type FindingCode string

// CodeInfo describes a finding code for documentation and display.
type CodeInfo struct {
	Code FindingCode
	Rule Rule
	// Name is a readable kebab-case name, such as "object-not-found".
	Name string
	// Title is a short sentence-case English title.
	Title string
}

// The finding codes, by rule. WALK_THROUGH's hundreds digit is the
// condition: 1 the object doesn't exist, 2 the name is taken, 3 the object
// is of another kind, 4 another failure the engine establishes.
const (
	SyntaxError FindingCode = "SY101"

	ObjectNotFound         FindingCode = "WT101"
	NameTaken              FindingCode = "WT201"
	WrongObjectKind        FindingCode = "WT301"
	DuplicateName          FindingCode = "WT401"
	NoCurrentSchema        FindingCode = "WT402"
	NameTooLong            FindingCode = "WT403"
	ReservedName           FindingCode = "WT404"
	SchemaNotEmpty         FindingCode = "WT405"
	SystemColumn           FindingCode = "WT406"
	CrossDatabaseReference FindingCode = "WT407"
	UnsupportedForm        FindingCode = "WT408"
	NotFirstInBatch        FindingCode = "WT409"
	TriggerSchemaMismatch  FindingCode = "WT410"
	IndexBacksConstraint   FindingCode = "WT411"

	OnlineMigrationUnsupported FindingCode = "OM101"
	DirectiveInvalid           FindingCode = "OM102"
	UnsupportedChange          FindingCode = "OM103"
	GhostColumnRename          FindingCode = "OM104"
	ArchiveDatabaseMissing     FindingCode = "OM105"
	TableNotEligible           FindingCode = "OM106"

	PriorBackupUnsupported FindingCode = "PB101"
	ChangeTooLarge         FindingCode = "PB102"
	BackupLocationMissing  FindingCode = "PB103"
	NothingToBackUp        FindingCode = "PB104"
	StatementsShareBackup  FindingCode = "PB105"
	BackupReadsWrongTable  FindingCode = "PB106"
	TableNotBackupable     FindingCode = "PB107"
	StatementNotBackupable FindingCode = "PB108"
	BackupReadsNewObject   FindingCode = "PB109"

	NullComparison FindingCode = "NU101"

	MissingWhere FindingCode = "WH101"

	DropsData       FindingCode = "DO101"
	DropsDefinition FindingCode = "DO102"
	DropsAccess     FindingCode = "DO103"

	TruncatesData FindingCode = "TR101"

	DropsConstraint FindingCode = "DC101"

	RenamesObject FindingCode = "RN101"

	TableWithoutPrimaryKey FindingCode = "PK101"
	PrimaryKeyRemoved      FindingCode = "PK102"
)

var codeInfos = []CodeInfo{
	{SyntaxError, Syntax, "syntax-error", "Statement doesn't parse"},

	{ObjectNotFound, WalkThrough, "object-not-found", "Object doesn't exist"},
	{NameTaken, WalkThrough, "name-taken", "Name is already taken"},
	{WrongObjectKind, WalkThrough, "wrong-object-kind", "Object is of another kind"},
	{DuplicateName, WalkThrough, "duplicate-name", "Statement uses one name twice"},
	{NoCurrentSchema, WalkThrough, "no-current-schema", "No database or schema selected"},
	{NameTooLong, WalkThrough, "name-too-long", "Name is too long"},
	{ReservedName, WalkThrough, "reserved-name", "Name is reserved"},
	{SchemaNotEmpty, WalkThrough, "schema-not-empty", "Schema still holds objects"},
	{SystemColumn, WalkThrough, "system-column", "System column can't be changed"},
	{CrossDatabaseReference, WalkThrough, "cross-database-reference", "Name refers to another database"},
	{UnsupportedForm, WalkThrough, "unsupported-form", "Engine rejects this form of statement"},
	{NotFirstInBatch, WalkThrough, "not-first-in-batch", "Statement must start its batch"},
	{TriggerSchemaMismatch, WalkThrough, "trigger-schema-mismatch", "Trigger and table are in different schemas"},
	{IndexBacksConstraint, WalkThrough, "index-backs-constraint", "Index belongs to a constraint"},

	{OnlineMigrationUnsupported, OnlineMigration, "online-migration-unsupported", "Engine has no online migration"},
	{DirectiveInvalid, OnlineMigration, "directive-invalid", "gh-ost directive is invalid"},
	{UnsupportedChange, OnlineMigration, "unsupported-change", "Change doesn't fit online migration"},
	{GhostColumnRename, OnlineMigration, "ghost-column-rename", "gh-ost can't rename the column"},
	{ArchiveDatabaseMissing, OnlineMigration, "archive-database-missing", "Backup database doesn't exist"},
	{TableNotEligible, OnlineMigration, "table-not-eligible", "Table doesn't meet gh-ost's requirements"},

	{PriorBackupUnsupported, PriorBackup, "prior-backup-unsupported", "Engine has no prior backup"},
	{ChangeTooLarge, PriorBackup, "change-too-large", "Change exceeds the backup size limit"},
	{BackupLocationMissing, PriorBackup, "backup-location-missing", "Backup database or schema doesn't exist"},
	{NothingToBackUp, PriorBackup, "nothing-to-back-up", "Change has nothing to back up"},
	{StatementsShareBackup, PriorBackup, "statements-share-backup", "Statements can't share one backup"},
	{BackupReadsWrongTable, PriorBackup, "backup-reads-wrong-table", "Backup may read another table"},
	{TableNotBackupable, PriorBackup, "table-not-backupable", "Table can't be backed up"},
	{StatementNotBackupable, PriorBackup, "statement-not-backupable", "Statement form can't be backed up"},
	{BackupReadsNewObject, PriorBackup, "backup-reads-new-object", "Backup reads what the change creates"},

	{NullComparison, RequireIsNull, "null-comparison", "Comparison with NULL is never true"},

	{MissingWhere, RequireWhere, "missing-where", "UPDATE or DELETE has no WHERE clause"},

	{DropsData, DisallowDropObject, "drops-data", "Drops data"},
	{DropsDefinition, DisallowDropObject, "drops-definition", "Drops a definition"},
	{DropsAccess, DisallowDropObject, "drops-access", "Drops a user, role, or security object"},

	{TruncatesData, DisallowTruncate, "truncates-data", "Truncates data"},

	{DropsConstraint, DisallowDropConstraint, "drops-constraint", "Drops a constraint"},

	{RenamesObject, DisallowRename, "renames-object", "Renames a table, view, or column"},

	{TableWithoutPrimaryKey, RequirePrimaryKey, "table-without-primary-key", "Creates a table without a primary key"},
	{PrimaryKeyRemoved, RequirePrimaryKey, "primary-key-removed", "Removes a table's primary key"},
}

// Codes lists every finding code in code order.
func Codes() []CodeInfo {
	return slices.Clone(codeInfos)
}

// Info returns the code's description, and false for a code not in the
// table.
func (c FindingCode) Info() (CodeInfo, bool) {
	for _, info := range codeInfos {
		if info.Code == c {
			return info, true
		}
	}
	return CodeInfo{}, false
}

// Rule is the rule the code belongs to, "" for a code not in the table.
func (c FindingCode) Rule() Rule {
	info, _ := c.Info()
	return info.Rule
}

// DropCode is the DISALLOW_DROP_OBJECT code for dropping an object of the
// kind a message names it by, such as "table" or "materialized view": data
// is lost with a table, column, partition, sequence, schema, or database;
// access with a user, role, or security object; anything else is a
// definition the change could create again, including a table whose rows
// live outside it.
func DropCode(kind string) FindingCode {
	k := strings.ToLower(kind)
	switch {
	case dropsDefinition[k]:
		return DropsDefinition
	case dropsData[k] || strings.HasSuffix(k, " table") || strings.HasSuffix(k, "partition") || strings.HasSuffix(k, " database"):
		return DropsData
	case dropsAccess[k] || strings.HasSuffix(k, " policy") || strings.HasSuffix(k, " key") || strings.HasSuffix(k, " role") ||
		strings.HasSuffix(k, " credential") || strings.HasSuffix(k, " audit") || strings.HasSuffix(k, " audit specification") ||
		strings.HasSuffix(k, " profile"):
		return DropsAccess
	}
	return DropsDefinition
}

// dropsData are the kinds whose drop loses rows or the values that number
// them, beyond the "... table", "... partition", and "... database" kinds.
var dropsData = map[string]bool{
	"table": true, "column": true, "columns": true, "attribute": true, "partition": true, "subpartition": true,
	"schema": true, "database": true, "sequence": true, "system versioning": true,
	"objects owned by": true, "unused columns": true,
}

// dropsDefinition are the kinds a suffix rule would misfile: an external
// table's rows stay in the files it reads and a dynamic table's come back
// from its query, and a placement policy says where data lives, not who
// may read it.
var dropsDefinition = map[string]bool{"external table": true, "dynamic table": true, "placement policy": true}

// dropsAccess are the kinds whose drop takes away who may connect or what
// they may see, beyond the "... policy", "... key", "... role",
// "... credential", "... audit", and "... profile" kinds.
var dropsAccess = map[string]bool{
	"user": true, "role": true, "login": true, "profile": true, "certificate": true, "credential": true,
	"user mapping": true, "signature": true, "counter signature": true, "sensitivity classification": true,
	"cryptographic provider": true, "policy": true, "network rule": true, "share": true,
	"managed account": true, "security integration": true,
}

// CheckCodes returns an error unless f carries at least one code, every
// code is in the table and of f's rule, and none repeats.
func CheckCodes(f Finding) error {
	if len(f.Codes) == 0 {
		return fmt.Errorf("%s finding %q has no code", f.Rule, f.Message)
	}
	seen := make(map[FindingCode]bool, len(f.Codes))
	for _, c := range f.Codes {
		if rule := c.Rule(); rule != f.Rule {
			return fmt.Errorf("%s finding %q has code %s of rule %q", f.Rule, f.Message, c, rule)
		}
		if seen[c] {
			return fmt.Errorf("%s finding %q repeats code %s", f.Rule, f.Message, c)
		}
		seen[c] = true
	}
	return nil
}
