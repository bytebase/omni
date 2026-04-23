package catalog

import nodes "github.com/bytebase/omni/pg/ast"

// MaxHeapAttributeNumber is the maximum number of columns allowed in a table.
//
// pg: src/include/access/htup_details.h — MaxHeapAttributeNumber
const MaxHeapAttributeNumber = 1600

// TypeName identifies a type for resolution.
type TypeName struct {
	Schema  string // empty = use search path
	Name    string
	TypeMod int32 // -1 = no typmod
	IsArray bool
}

// ColumnDef describes a column in a CREATE TABLE statement.
type ColumnDef struct {
	Name             string
	Type             TypeName
	NotNull          bool
	Default          string          // opaque expression; empty = no default
	RawDefault       nodes.Node      // raw AST node for default expression (for analysis)
	IsSerial         byte            // 0=none, 2=smallserial, 4=serial, 8=bigserial
	Generated        byte            // 's' = stored generated, 0 = none
	GenerationExpr   string          // expression text for generated column
	RawGenExpr       nodes.Node      // raw AST node for generated expression (for analysis)
	Identity         byte            // 'a' = ALWAYS, 'd' = BY DEFAULT, 0 = none
	IdentityOptions  *nodes.List     // raw sequence options from GENERATED ... AS IDENTITY
	IdentitySequence *nodes.RangeVar // selected sequence name after transform/expansion
	IsLocal          bool            // true if defined locally (not only inherited)
	InhCount         int             // number of inheritance ancestors that define this column
	CollationName    string          // explicit COLLATE clause; empty = type default
	IsFromType       bool            // true if column came from OF TYPE (typed table)
}

// ConstraintType identifies the kind of table constraint.
type ConstraintType byte

const (
	ConstraintPK      ConstraintType = 'p'
	ConstraintUnique  ConstraintType = 'u'
	ConstraintFK      ConstraintType = 'f'
	ConstraintCheck   ConstraintType = 'c'
	ConstraintExclude ConstraintType = 'x'
	ConstraintTrigger ConstraintType = 't' // constraint trigger
)

// ConstraintDef describes a constraint in a CREATE TABLE statement.
type ConstraintDef struct {
	Name           string // empty = auto-generate
	Type           ConstraintType
	Columns        []string   // PK/UNIQUE/FK local columns
	IndexName      string     // PK/UNIQUE: user-specified USING INDEX name
	RefSchema      string     // FK only
	RefTable       string     // FK only
	RefColumns     []string   // FK only (empty = use PK of target)
	FKUpdAction    byte       // FK only: 'a', 'r', 'c', 'n', 'd'
	FKDelAction    byte       // FK only: 'a', 'r', 'c', 'n', 'd'
	FKMatchType    byte       // FK only: 's', 'f', 'p'
	Deferrable     bool       // DEFERRABLE
	Deferred       bool       // INITIALLY DEFERRED
	SkipValidation bool       // NOT VALID (CHECK/FK only)
	CheckExpr      string     // CHECK only
	RawCheckExpr   nodes.Node // raw AST node for CHECK expression (for analysis)
	ExclOps        []string   // EXCLUDE only: operator names per column
	AccessMethod   string     // EXCLUDE only: access method (default "gist")
}

// --- TRIGGER ---

// TriggerTiming indicates when a trigger fires.
type TriggerTiming byte

const (
	TriggerBefore    TriggerTiming = 'B'
	TriggerAfter     TriggerTiming = 'A'
	TriggerInsteadOf TriggerTiming = 'I'
)

// TriggerEvent is a bitmask of trigger events.
type TriggerEvent uint8

const (
	TriggerEventInsert TriggerEvent = 1 << iota
	TriggerEventUpdate
	TriggerEventDelete
	TriggerEventTruncate
)

// PARTITION_MAX_KEYS is the maximum number of partition key columns.
//
// pg: src/include/pg_config_manual.h — PARTITION_MAX_KEYS
const PARTITION_MAX_KEYS = 32

// INDEX_MAX_KEYS is the maximum number of columns in an index.
//
// pg: src/include/pg_config_manual.h — INDEX_MAX_KEYS
const INDEX_MAX_KEYS = 32

// Type kind constants (from pg_type.h TYPTYPE_*)
const (
	TYPTYPE_BASE       = 'b'
	TYPTYPE_COMPOSITE  = 'c'
	TYPTYPE_DOMAIN     = 'd'
	TYPTYPE_ENUM       = 'e'
	TYPTYPE_PSEUDO     = 'p'
	TYPTYPE_RANGE      = 'r'
	TYPTYPE_MULTIRANGE = 'm'
)

// Function kind constants (from pg_proc.h PROKIND_*)
const (
	PROKIND_FUNCTION  = 'f'
	PROKIND_PROCEDURE = 'p'
	PROKIND_AGGREGATE = 'a'
	PROKIND_WINDOW    = 'w'
)
