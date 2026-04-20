// Package catalog — define.go
//
// AST-level Define API for catalog initialization.
//
// These entry points let callers install schema objects without going
// through the SQL lexer/parser, mirroring the shape of pg/catalog's
// Define{Relation,View,Enum,…} family. Each wrapper is a thin guard
// around the internal create* method reached today via Exec's
// processUtility dispatch; no additional validation or side effects
// are introduced beyond rejecting nil / empty-named stmts.
//
// # Design philosophy: loader, not validator
//
// This API is for cold-starting an in-memory catalog from structured
// metadata. It attempts a best-effort install of each object.
// Constraint checking (FK integrity, routine body validity, etc.) is
// explicitly NOT a responsibility of this API:
//
//   - FKs installed while ForeignKeyChecks()==false are stored as-is
//     and never revalidated.
//   - Views whose SELECT body fails AnalyzeSelectStmt are still
//     created, but with AnalyzedQuery == nil and a nil error.
//
// Callers who need validated state must either load with FK checks on
// (paying the topological-order cost) or run their own downstream
// checks.
//
// # Preconditions (per call)
//
//   - SetCurrentDatabase(name) MUST be called whenever the stmt does
//     not carry an explicit database qualifier. Missing currentDB
//     surfaces as *Error{Code: ErrNoDatabaseSelected}.
//   - SetForeignKeyChecks(false) is typically required while
//     bulk-loading schemas with forward FK references. Re-enabling it
//     after the load does not retroactively validate FKs already
//     installed.
//   - Topological ordering across kinds is the caller's responsibility:
//     DefineTrigger and DefineIndex on a not-yet-installed target
//     table return *Error{Code: ErrNoSuchTable}. DefineView tolerates
//     forward refs but yields AnalyzedQuery=nil.
//
// # Error contract
//
// Every Define* returns error, always of concrete type *Error when
// non-nil. Callers may inspect err.(*Error).Code (ErrDupTable,
// ErrNoDatabaseSelected, ErrWrongArguments, etc.) for idempotency and
// fallback decisions. On error, no catalog state is written; the call
// is atomic at the object level.
//
// # Concurrency
//
// *Catalog is NOT goroutine-safe. The underlying maps have no sync
// primitives. Callers MUST serialize Define* calls on a given Catalog.
//
// # Nil / empty-name guards
//
// Every Define* rejects nil stmt and empty required names with
// *Error{Code: ErrWrongArguments} rather than panicking. Per-kind
// required fields:
//
//	DefineDatabase:  stmt.Name
//	DefineTable:     stmt.Table.Name
//	DefineView:      stmt.Name.Name
//	DefineIndex:     stmt.Table.Name (and stmt.IndexName for the index itself)
//	DefineFunction,
//	DefineProcedure,
//	DefineRoutine:   stmt.Name.Name
//	DefineTrigger:   stmt.Name and stmt.Table.Name
//	DefineEvent:     stmt.Name
package catalog

import nodes "github.com/bytebase/omni/tidb/ast"

// DefineDatabase installs a database. stmt.Name must be non-empty.
func (c *Catalog) DefineDatabase(stmt *nodes.CreateDatabaseStmt) error {
	if stmt == nil || stmt.Name == "" {
		return errWrongArguments("DefineDatabase")
	}
	return c.createDatabase(stmt)
}

// DefineTable installs a table (including inline columns, indexes,
// foreign keys, CHECK constraints, and partitions). stmt.Table and
// stmt.Table.Name must be non-nil/non-empty.
//
// Foreign-key validity depends on the current foreign_key_checks
// session flag. See package doc.
func (c *Catalog) DefineTable(stmt *nodes.CreateTableStmt) error {
	if stmt == nil || stmt.Table == nil || stmt.Table.Name == "" {
		return errWrongArguments("DefineTable")
	}
	return c.createTable(stmt)
}

// DefineView installs a view. stmt.Name and stmt.Name.Name must be
// non-nil/non-empty.
//
// If AnalyzeSelectStmt on the view body fails (e.g. referenced table
// is not yet installed), DefineView still returns nil and the view is
// stored with AnalyzedQuery=nil. This is intentional loader behavior.
func (c *Catalog) DefineView(stmt *nodes.CreateViewStmt) error {
	if stmt == nil || stmt.Name == nil || stmt.Name.Name == "" {
		return errWrongArguments("DefineView")
	}
	return c.createView(stmt)
}

// DefineIndex installs an index on an existing table.
func (c *Catalog) DefineIndex(stmt *nodes.CreateIndexStmt) error {
	if stmt == nil || stmt.Table == nil || stmt.Table.Name == "" {
		return errWrongArguments("DefineIndex")
	}
	return c.createIndex(stmt)
}

// DefineFunction installs a stored function.
//
// Routing note: this function is a thin wrapper over createRoutine.
// The catalog routes the stmt to db.Functions or db.Procedures based
// on stmt.IsProcedure. Callers who set IsProcedure=true on a stmt
// passed to DefineFunction will land in db.Procedures — no kind-guard
// is applied. Use DefineFunction at call sites where intent is known
// to be a function, for readability; use DefineRoutine for generic
// paths.
func (c *Catalog) DefineFunction(stmt *nodes.CreateFunctionStmt) error {
	if stmt == nil || stmt.Name == nil || stmt.Name.Name == "" {
		return errWrongArguments("DefineFunction")
	}
	return c.createRoutine(stmt)
}

// DefineProcedure installs a stored procedure.
//
// See DefineFunction for the routing semantics — this wrapper is
// identical and exists purely for call-site clarity when the caller
// knows the intent is a procedure.
func (c *Catalog) DefineProcedure(stmt *nodes.CreateFunctionStmt) error {
	if stmt == nil || stmt.Name == nil || stmt.Name.Name == "" {
		return errWrongArguments("DefineProcedure")
	}
	return c.createRoutine(stmt)
}

// DefineRoutine installs a function or procedure, routed by
// stmt.IsProcedure. Use this when the caller does not statically know
// the kind (e.g. bulk loaders processing heterogeneous metadata).
func (c *Catalog) DefineRoutine(stmt *nodes.CreateFunctionStmt) error {
	if stmt == nil || stmt.Name == nil || stmt.Name.Name == "" {
		return errWrongArguments("DefineRoutine")
	}
	return c.createRoutine(stmt)
}

// DefineTrigger installs a trigger on an existing table. stmt.Name
// (trigger name) and stmt.Table.Name must be non-empty.
func (c *Catalog) DefineTrigger(stmt *nodes.CreateTriggerStmt) error {
	if stmt == nil || stmt.Name == "" || stmt.Table == nil || stmt.Table.Name == "" {
		return errWrongArguments("DefineTrigger")
	}
	return c.createTrigger(stmt)
}

// DefineEvent installs an event in the current or specified database.
func (c *Catalog) DefineEvent(stmt *nodes.CreateEventStmt) error {
	if stmt == nil || stmt.Name == "" {
		return errWrongArguments("DefineEvent")
	}
	return c.createEvent(stmt)
}
