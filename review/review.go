// Package review is the contract between Bytebase SQL Review V2 and the
// engine review packages (mysql/review, pg/review, ...). It holds the
// types that cross that boundary and nothing else: no engine imports, no
// parser, no AST.
//
// The unit of a review is one SQL text against every target it will run
// on. An engine package exposes
//
//	var Rules = []review.Rule{...}
//	func Review(ctx context.Context, sql string, opts review.Options, targets []review.Target) (*review.Result, error)
//
// and the caller makes one call per (change, engine). Splitting, parsing,
// evaluating the rules that depend on the SQL alone once, deduplicating
// targets with equal inputs, concurrency, and merging findings across
// targets are all internal to the engine package. The caller supplies
// inputs and reads findings.
//
// The cut sits between facts and judgment. The caller resolves what needs
// a connection or its own configuration (the synced schema, whether the
// change asks for online migration or prior backup, whether the backup
// database exists) and passes plain values; the engine package derives
// every finding from the SQL plus those values. A syntax error is a
// finding, never an error. An error means the call could not be
// evaluated; a per-target problem is a Result.Failures entry.
package review

import (
	"context"

	"github.com/bytebase/omni/metadata"
)

// Rule names a review rule. The value equals the bytebase.v1.ReviewRuleType
// enum name so the caller maps by name, which tolerates the two sides
// being on different versions: a rule an engine does not know is ignored,
// a rule the caller does not know is dropped.
type Rule string

const (
	// Syntax: the statements do not parse for the engine.
	Syntax Rule = "SYNTAX"
	// WalkThrough: applying the statements to the pre-change schema fails.
	WalkThrough Rule = "WALK_THROUGH"
	// OnlineMigration: the change requests online migration but is not eligible.
	OnlineMigration Rule = "ONLINE_MIGRATION"
	// PriorBackup: the change enables prior backup but the backup cannot be taken.
	PriorBackup Rule = "PRIOR_BACKUP"
	// RequireIsNull: = NULL or <> NULL in a predicate, which is always false.
	RequireIsNull Rule = "REQUIRE_IS_NULL"
	// RequireWhere: UPDATE or DELETE without WHERE.
	RequireWhere Rule = "REQUIRE_WHERE"
	// DisallowDropObject: DROP TABLE, COLUMN, SCHEMA, or DATABASE.
	DisallowDropObject Rule = "DISALLOW_DROP_OBJECT"
	// DisallowTruncate: TRUNCATE.
	DisallowTruncate Rule = "DISALLOW_TRUNCATE"
	// DisallowDropConstraint: dropping a PRIMARY KEY, FOREIGN KEY, UNIQUE, or CHECK constraint.
	DisallowDropConstraint Rule = "DISALLOW_DROP_CONSTRAINT"
	// DisallowRename: renaming a table or column.
	DisallowRename Rule = "DISALLOW_RENAME"
	// RequirePrimaryKey: the change creates a table without a primary key,
	// or drops one without adding it back.
	RequirePrimaryKey Rule = "REQUIRE_PRIMARY_KEY"
)

// ReviewFunc is the signature of every engine package's Review.
type ReviewFunc func(ctx context.Context, sql string, opts Options, targets []Target) (*Result, error)

// Options applies to the whole change: the same for every target.
type Options struct {
	// Rules to evaluate. nil means every rule the engine implements. A
	// rule the engine does not implement is ignored, not an error. Syntax
	// is always evaluated whether or not it is listed: parsing comes
	// before every other rule, and a statement that does not parse is
	// reported rather than silently skipped.
	Rules []Rule
	// Change is what the change asks for.
	Change Change
}

// Change is the change's requested behaviors, read from the spec and the
// sheet.
type Change struct {
	// OnlineMigration is true when the sheet carries a gh-ost directive.
	OnlineMigration bool
	// PriorBackup is true when the spec enables prior backup.
	PriorBackup bool
	// MaxBackupSize is the largest SQL text prior backup can handle, in
	// bytes; 0 means no limit.
	MaxBackupSize int
}

// Target is one database's input. Findings refer to targets by their
// index in the slice passed to Review; the value carries no name of its
// own beyond Schema.Name. Two targets with equal values produce equal
// findings, which is what lets an engine review each distinct value once.
type Target struct {
	// Schema is the synced pre-change schema. nil is allowed: WalkThrough
	// is then not evaluated, and the rules that consult the schema fall
	// back to what the change itself creates. The current database is
	// Schema.Name and the PostgreSQL search path is Schema.SearchPath.
	Schema *metadata.DatabaseSchemaMetadata
	// BackupDatabaseExists reports whether the engine's backup database
	// (bbdataarchive) exists on the target's instance. PriorBackup stores
	// its copies there and OnlineMigration stages its gh-ost tables there,
	// so the one fact serves both rules. Engines whose backup location is
	// a schema inside the database (PostgreSQL) read it from Schema and
	// ignore this field.
	BackupDatabaseExists bool
	// SessionUser is the role the change runs as; PostgreSQL ownership
	// checks in WalkThrough use it. Empty disables those checks: the
	// engine has no connection to ask, so the caller resolves the role or
	// opts out.
	SessionUser string
	// LowerCaseTableNames is the MySQL family's lower_case_table_names
	// server variable (0: names are case-sensitive; 1 or 2: not), an
	// instance property the schema does not carry. Other engines ignore it.
	LowerCaseTableNames int
}

// Result is one review of one SQL text against its targets.
type Result struct {
	// Statements lists every non-empty statement's byte range in the
	// input, in order, whether or not it parsed. Findings refer to them by
	// index.
	Statements []Range
	// Findings, already merged across targets, sorted by (Statement,
	// Range.Start, Rule, Message).
	Findings []Finding
	// Failures lists the targets that could not be evaluated, with the
	// reason. A target listed here has no target-dependent findings; the
	// findings that depend on the SQL alone still apply to it.
	Failures []TargetFailure
}

// Range is a half-open byte range [Start, End) in the input SQL.
type Range struct {
	Start int
	End   int
}

// Finding is one merged result.
type Finding struct {
	Rule Rule
	// Statement is the index into Result.Statements, or -1 when the
	// finding addresses the whole change rather than one statement.
	Statement int
	// Range is the anchor in absolute byte offsets of the input. A zero
	// Range with Statement >= 0 anchors the whole statement.
	Range Range
	// Message is one line naming the object as the SQL wrote it. It never
	// names the target, so equal findings from different targets merge.
	Message string
	// Targets are the indices of the targets this finding applies to,
	// ascending. A finding that depends on the SQL alone lists every
	// target.
	Targets []int
}

// TargetFailure is one target that could not be evaluated.
type TargetFailure struct {
	Target int
	Err    error
}
