package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
)

// TableRef is an engine-neutral schema-qualified table or object reference.
type TableRef struct {
	Schema string
	Name   string
}

// RenamePair captures the source and destination tables of a rename statement.
type RenamePair struct {
	Old TableRef
	New TableRef
}

// WalkThroughStatementKind classifies an AST statement for DDL walk-through
// precheck and cross-database rename target discovery.
type WalkThroughStatementKind int

const (
	// WalkThroughStatementOther is any statement not needing walk-through precheck.
	WalkThroughStatementOther WalkThroughStatementKind = iota
	// WalkThroughStatementCreateDatabase is a CREATE DATABASE / CREATE SCHEMA statement.
	WalkThroughStatementCreateDatabase
	// WalkThroughStatementAlterDatabase is an ALTER DATABASE / ALTER SCHEMA statement.
	WalkThroughStatementAlterDatabase
	// WalkThroughStatementDropDatabase is a DROP DATABASE / DROP SCHEMA statement.
	WalkThroughStatementDropDatabase
	// WalkThroughStatementCreateTable is a CREATE TABLE statement.
	WalkThroughStatementCreateTable
	// WalkThroughStatementDropTable is a DROP TABLE statement.
	WalkThroughStatementDropTable
	// WalkThroughStatementTableRefCheck is a DDL statement whose table or object
	// references must belong to the current database.
	WalkThroughStatementTableRefCheck
	// WalkThroughStatementRenameTable is a RENAME TABLE statement.
	WalkThroughStatementRenameTable
)

// String returns the canonical name of the walk-through statement kind.
func (k WalkThroughStatementKind) String() string {
	switch k {
	case WalkThroughStatementCreateDatabase:
		return "CREATE_DATABASE"
	case WalkThroughStatementAlterDatabase:
		return "ALTER_DATABASE"
	case WalkThroughStatementDropDatabase:
		return "DROP_DATABASE"
	case WalkThroughStatementCreateTable:
		return "CREATE_TABLE"
	case WalkThroughStatementDropTable:
		return "DROP_TABLE"
	case WalkThroughStatementTableRefCheck:
		return "TABLE_REF_CHECK"
	case WalkThroughStatementRenameTable:
		return "RENAME_TABLE"
	default:
		return "OTHER"
	}
}

// WalkThroughStatement holds the engine-neutral fields needed for walk-through
// precheck and cross-database rename target discovery. Source is optional
// caller-supplied SQL text for diagnostic formatting.
type WalkThroughStatement struct {
	Kind        WalkThroughStatementKind
	Source      string
	Database    string
	IfNotExists bool
	IsCTAS      bool
	TableRefs   []TableRef
	RenamePairs []RenamePair
}

// WalkThroughRequest is a script simulated against a database's synced schema.
type WalkThroughRequest struct {
	// Database is the reviewed database; its Metadata is the schema before the
	// script.
	Database *Database
	// Script is the whole script the statements were parsed from.
	Script string
	// SessionUser is the user the script runs as, or "" for the catalog's default.
	SessionUser string
	// SearchPath is the configured search path, or nil for the engine's default.
	SearchPath []string
}

// WalkThroughCode names the condition a walk-through finding reports.
type WalkThroughCode string

const (
	// WalkThroughDDLSimulationFailed is a failure the catalog reports with no
	// more specific code, or a failure to set the catalog up.
	WalkThroughDDLSimulationFailed WalkThroughCode = "ddl-simulation-failed"
	// WalkThroughNotCurrentDatabase is a reference to another database.
	WalkThroughNotCurrentDatabase WalkThroughCode = "not-current-database"
	// WalkThroughDatabaseIsDeleted is a drop of the reviewed database.
	WalkThroughDatabaseIsDeleted WalkThroughCode = "database-is-deleted"
	// WalkThroughCreateTableAs is a CREATE TABLE ... AS SELECT, whose columns
	// the catalog cannot know.
	WalkThroughCreateTableAs                        WalkThroughCode = "create-table-as"
	WalkThroughTableExists                          WalkThroughCode = "table-exists"
	WalkThroughTableNotExists                       WalkThroughCode = "table-not-exists"
	WalkThroughColumnExists                         WalkThroughCode = "column-exists"
	WalkThroughColumnNotExists                      WalkThroughCode = "column-not-exists"
	WalkThroughInvalidColumnDefault                 WalkThroughCode = "invalid-column-default"
	WalkThroughIndexExists                          WalkThroughCode = "index-exists"
	WalkThroughIndexNotExists                       WalkThroughCode = "index-not-exists"
	WalkThroughPrimaryKeyExists                     WalkThroughCode = "primary-key-exists"
	WalkThroughSpatialIndexKeyNullable              WalkThroughCode = "spatial-index-key-nullable"
	WalkThroughAutoIncrementExists                  WalkThroughCode = "auto-increment-exists"
	WalkThroughOnUpdateColumnNotDatetimeOrTimestamp WalkThroughCode = "on-update-column-not-datetime-or-timestamp"
	WalkThroughTableHasFK                           WalkThroughCode = "table-has-fk"
	WalkThroughRelationExists                       WalkThroughCode = "relation-exists"
	WalkThroughSchemaNotExists                      WalkThroughCode = "schema-not-exists"
	WalkThroughChangeColumnType                     WalkThroughCode = "change-column-type"
	WalkThroughColumnCannotNull                     WalkThroughCode = "column-cannot-null"
	WalkThroughSyntaxError                          WalkThroughCode = "syntax-error"
	WalkThroughNameIsKeywordIdentifier              WalkThroughCode = "name-is-keyword-identifier"
)

// WalkThroughFinding is the first problem a walk-through finds.
type WalkThroughFinding struct {
	Code WalkThroughCode
	// Error is true for an error and false for a warning.
	Error   bool
	Title   string
	Content string
	// Statement is the index of the statement the finding is at, or -1 when
	// Start places it.
	Statement int
	// Start is where the finding is when Statement is -1; line 0 means the
	// script as a whole.
	Start Position
}

// WalkThroughResult is a walk-through's outcome.
type WalkThroughResult struct {
	// Database is the catalog's schema after the script, or nil when a finding
	// stopped the simulation. An engine may also return nil when the script
	// leaves the schema unchanged.
	Database *Database
	Finding  *WalkThroughFinding
}

// WalkThrougher simulates a script's DDL against a database's schema.
type WalkThrougher interface {
	WalkThrough(ctx context.Context, request WalkThroughRequest, statements []ParsedStatement) WalkThroughResult
}

// MySQLWalkThrough is the walk-through MySQL, MariaDB, OceanBase and TiDB
// share. Each engine configures it with its catalog and dialect, and its
// WalkThrough method implements WalkThrougher.
type MySQLWalkThrough struct {
	// Classify returns a statement node's walk-through fields, or false for a
	// node of another engine.
	Classify func(node any) (WalkThroughStatement, bool)
	// IsCreateTableAs finds a CREATE TABLE ... AS among Split's segments by
	// its text, which covers the forms the parser builds no node for. A nil
	// IsCreateTableAs skips that check.
	IsCreateTableAs func(statement string) bool
	Split           func(script string) []Statement
	// ExecScript returns the script the catalog runs and, for each statement
	// of it, the index of the statement it came from. A nil ExecScript runs
	// the request's script, placing a failure by its line.
	ExecScript func(statements []ParsedStatement) (string, []int)
	NewCatalog func() MySQLCatalog
	// ErrorNumber returns a catalog error's MySQL error number and message.
	ErrorNumber func(err error) (number int, message string, ok bool)
	// TrackBatchTables lets CREATE TABLE IF NOT EXISTS ... AS skip the
	// CREATE TABLE AS warning for a table an earlier statement created.
	TrackBatchTables bool
	// TiDBErrors words catalog errors as TiDB does.
	TiDBErrors bool
	// KeepCaseSensitivity gives the result the reviewed database's object case
	// sensitivity; otherwise the result is case-sensitive.
	KeepCaseSensitivity bool
}

// MySQLCatalog is a MySQL-family catalog a walk-through runs on.
type MySQLCatalog interface {
	Load(ctx context.Context, database *metadata.DatabaseSchemaMetadata) error
	// Run executes a script, continuing past a statement that fails.
	Run(script string) ([]CatalogResult, error)
	DatabaseMetadata(database string) *metadata.DatabaseSchemaMetadata
}

// CatalogResult is the outcome of one statement a catalog ran.
type CatalogResult struct {
	// Index is the statement's index in the script.
	Index int
	// Line is the 1-based line the statement starts on.
	Line  int
	Error error
}

// orderedFinding is a finding with the statement index it sorts by, -1 when
// it has none.
type orderedFinding struct {
	order   int
	finding *WalkThroughFinding
}

// beforeOrSame reports whether f wins over g: the earlier statement, then the
// error over the warning, then f. A finding with no statement index sorts last.
func (f orderedFinding) beforeOrSame(g orderedFinding) bool {
	if f.finding == nil {
		return false
	}
	if g.finding == nil {
		return true
	}
	if f.order < 0 || g.order < 0 {
		return f.order >= 0 || g.order < 0
	}
	if f.order != g.order {
		return f.order < g.order
	}
	if f.finding.Error != g.finding.Error {
		return f.finding.Error
	}
	return true
}

// WalkThrough prechecks what the catalog would accept silently, then runs the
// script on a catalog loaded with the database's schema.
func (w MySQLWalkThrough) WalkThrough(ctx context.Context, request WalkThroughRequest, statements []ParsedStatement) WalkThroughResult {
	script, execStatements := request.Script, []int(nil)
	if w.ExecScript != nil {
		script, execStatements = w.ExecScript(statements)
	}
	if strings.TrimSpace(script) == "" {
		return WalkThroughResult{}
	}
	db := mysqlWalkThroughDatabase{request.Database}
	infos := make([]*WalkThroughStatement, len(statements))
	for i, statement := range statements {
		if info, ok := w.Classify(statement.Node); ok {
			info.Source = strings.TrimSpace(statement.Text)
			infos[i] = &info
		}
	}

	precheck := db.precheck(infos, w.TrackBatchTables)
	if w.IsCreateTableAs != nil {
		if raw := db.precheckRawCreateTableAs(script, infos, w.Split, w.IsCreateTableAs); raw.beforeOrSame(precheck) {
			precheck = raw
		}
	}

	c := w.NewCatalog()
	setup := "SET foreign_key_checks = 0;"
	for _, target := range db.renameTargetDatabases(infos) {
		setup += fmt.Sprintf("\nCREATE DATABASE IF NOT EXISTS `%s`;", strings.ReplaceAll(target, "`", "``"))
	}
	if _, err := c.Run(setup); err != nil {
		return walkThroughFailure("Failed to initialize catalog", err)
	}
	if err := c.Load(ctx, db.Metadata); err != nil {
		return walkThroughFailure("Failed to load schema", err)
	}

	results, err := c.Run(script)
	var failure orderedFinding
	if err != nil {
		failure = w.catalogFailure(db, -1, -1, Position{}, err)
	}
	for _, r := range results {
		if r.Error == nil {
			continue
		}
		statement, start := -1, Position{Line: r.Line}
		if execStatements != nil {
			start = Position{}
			if r.Index >= 0 && r.Index < len(execStatements) {
				statement = execStatements[r.Index]
			}
		}
		failure = w.catalogFailure(db, r.Index, statement, start, r.Error)
		break
	}
	if precheck.beforeOrSame(failure) {
		return WalkThroughResult{Finding: precheck.finding}
	}
	if failure.finding != nil {
		return WalkThroughResult{Finding: failure.finding}
	}
	return WalkThroughResult{Database: &Database{
		Name:                db.Name,
		Metadata:            c.DatabaseMetadata(db.Name),
		ObjectCaseSensitive: !w.KeepCaseSensitivity || db.ObjectCaseSensitive,
	}}
}

func walkThroughFailure(title string, err error) WalkThroughResult {
	return WalkThroughResult{Finding: &WalkThroughFinding{
		Code:      WalkThroughDDLSimulationFailed,
		Error:     true,
		Title:     title,
		Content:   err.Error(),
		Statement: -1,
	}}
}

func walkThroughWarning(code WalkThroughCode, content string) *WalkThroughFinding {
	return &WalkThroughFinding{Code: code, Title: content, Content: content, Statement: -1}
}

// mysqlWalkThroughDatabase is the reviewed database of a MySQL-family
// walk-through.
type mysqlWalkThroughDatabase struct{ *Database }

func (db mysqlWalkThroughDatabase) isCurrent(database string) bool {
	if !db.ObjectCaseSensitive {
		return strings.EqualFold(db.Name, database)
	}
	return db.Name == database
}

// inCurrent reports whether table names a table of the reviewed database.
func (db mysqlWalkThroughDatabase) inCurrent(table TableRef) bool {
	return table.Name != "" && (table.Schema == "" || db.isCurrent(table.Schema))
}

func (db mysqlWalkThroughDatabase) tableKey(name string) string {
	if !db.ObjectCaseSensitive {
		return strings.ToLower(name)
	}
	return name
}

func (db mysqlWalkThroughDatabase) tableExists(table TableRef) bool {
	return db.inCurrent(table) && db.Lookup(false).Schema("").Table(table.Name) != nil
}

// renameTargetDatabases lists, once each, the other databases that a rename of
// a table of the reviewed database moves it into.
func (db mysqlWalkThroughDatabase) renameTargetDatabases(infos []*WalkThroughStatement) []string {
	seen := make(map[string]bool)
	var targets []string
	for _, info := range infos {
		if info == nil {
			continue
		}
		for _, pair := range info.RenamePairs {
			if pair.Old.Schema != "" && !db.isCurrent(pair.Old.Schema) {
				continue
			}
			if pair.New.Schema == "" || db.isCurrent(pair.New.Schema) || seen[pair.New.Schema] {
				continue
			}
			seen[pair.New.Schema] = true
			targets = append(targets, pair.New.Schema)
		}
	}
	return targets
}

// precheck reports the first statement the catalog would accept silently but
// a review warns on: a reference to another database, a dropped database, or a
// CREATE TABLE ... AS SELECT.
func (db mysqlWalkThroughDatabase) precheck(infos []*WalkThroughStatement, trackBatchTables bool) orderedFinding {
	var created map[string]bool
	if trackBatchTables {
		created = make(map[string]bool)
	}
	for i, info := range infos {
		if info == nil {
			continue
		}
		if finding := db.precheckStatement(*info, created); finding != nil {
			finding.Statement = i
			return orderedFinding{order: i, finding: finding}
		}
		if created == nil {
			continue
		}
		switch info.Kind {
		case WalkThroughStatementCreateTable:
			if len(info.TableRefs) > 0 && db.inCurrent(info.TableRefs[0]) {
				created[db.tableKey(info.TableRefs[0].Name)] = true
			}
		case WalkThroughStatementDropTable:
			for _, t := range info.TableRefs {
				if db.inCurrent(t) {
					delete(created, db.tableKey(t.Name))
				}
			}
		default:
		}
	}
	return orderedFinding{}
}

func (db mysqlWalkThroughDatabase) precheckStatement(info WalkThroughStatement, created map[string]bool) *WalkThroughFinding {
	switch info.Kind {
	case WalkThroughStatementCreateDatabase, WalkThroughStatementAlterDatabase:
		if info.Database != "" && !db.isCurrent(info.Database) {
			return db.notCurrent(info.Database)
		}
	case WalkThroughStatementDropDatabase:
		if info.Database != "" && !db.isCurrent(info.Database) {
			return db.notCurrent(info.Database)
		}
		if info.Database != "" {
			return walkThroughWarning(WalkThroughDatabaseIsDeleted, fmt.Sprintf("Database `%s` is deleted", info.Database))
		}
	case WalkThroughStatementCreateTable:
		if len(info.TableRefs) > 0 {
			table := info.TableRefs[0]
			if table.Schema != "" && !db.isCurrent(table.Schema) {
				return db.notCurrent(table.Schema)
			}
			if info.IfNotExists && (db.tableExists(table) || (created != nil && db.inCurrent(table) && created[db.tableKey(table.Name)])) {
				return nil
			}
		}
		if info.IsCTAS {
			return walkThroughWarning(WalkThroughCreateTableAs, fmt.Sprintf("CREATE TABLE AS statement is used in \"%s\"", info.Source))
		}
	case WalkThroughStatementDropTable, WalkThroughStatementTableRefCheck, WalkThroughStatementRenameTable:
		for _, ref := range info.TableRefs {
			if ref.Schema != "" && !db.isCurrent(ref.Schema) {
				return db.notCurrent(ref.Schema)
			}
		}
	default:
	}
	return nil
}

func (db mysqlWalkThroughDatabase) notCurrent(database string) *WalkThroughFinding {
	return walkThroughWarning(WalkThroughNotCurrentDatabase, fmt.Sprintf("Database `%s` is not the current database `%s`", database, db.Name))
}

// precheckRawCreateTableAs finds a CREATE TABLE ... AS by its text, for a form
// the parser builds no node for. The n-th non-empty segment pairs with the
// n-th statement.
func (db mysqlWalkThroughDatabase) precheckRawCreateTableAs(script string, infos []*WalkThroughStatement, split func(string) []Statement, isCreateTableAs func(string) bool) orderedFinding {
	index := 0
	for _, segment := range split(script) {
		if segment.Empty {
			continue
		}
		current := index
		index++
		text := strings.TrimSpace(segment.Text)
		if !isCreateTableAs(text) {
			continue
		}
		if current < len(infos) {
			if info := infos[current]; info != nil && info.Kind == WalkThroughStatementCreateTable && info.IsCTAS && info.IfNotExists && len(info.TableRefs) > 0 && db.tableExists(info.TableRefs[0]) {
				continue
			}
		}
		finding := walkThroughWarning(WalkThroughCreateTableAs, fmt.Sprintf("CREATE TABLE AS statement is used in \"%s\"", text))
		finding.Start = PositionAt(script, segment.Start)
		return orderedFinding{order: current, finding: finding}
	}
	return orderedFinding{}
}

// MySQL-family error numbers the walk-through maps to a code.
const (
	errUnknownDatabase          = 1049
	errDupTable                 = 1050
	errUnknownTable             = 1051
	errNoSuchColumn             = 1054
	errDupColumn                = 1060
	errDupKeyName               = 1061
	errDupEntry                 = 1062
	errInvalidDefault           = 1067
	errMultiplePriKey           = 1068
	errKeyColumnDoesNotExist    = 1072
	errIncorrectTableDefinition = 1075
	errCantDropKey              = 1091
	errBlobCantHaveDefault      = 1101
	errNoSuchTable              = 1146
	errSpatialIndexKeyNullable  = 1252
	errInvalidOnUpdate          = 1294
	errFKMissingIndex           = 1822
	errFKNoRefTable             = 1824
	errDupIndex                 = 1831
	errFKCannotDropParent       = 3730
	errFKIncompatibleColumns    = 3780
)

// catalogFailure is the error finding for a statement the catalog refused,
// sorting by order.
func (w MySQLWalkThrough) catalogFailure(db mysqlWalkThroughDatabase, order, statement int, start Position, err error) orderedFinding {
	code, content := WalkThroughDDLSimulationFailed, err.Error()
	if number, message, ok := w.ErrorNumber(err); ok {
		code, content = db.mapErrorNumber(number, message, w.TiDBErrors)
	}
	if !w.TiDBErrors && strings.Contains(content, "BLOB, TEXT, GEOMETRY or JSON column ") {
		code = WalkThroughInvalidColumnDefault
		content = fmt.Sprintf("BLOB, TEXT, GEOMETRY or JSON column `%s` can't have a default value", quotedName(content))
	}
	return orderedFinding{order: order, finding: &WalkThroughFinding{
		Code:      code,
		Error:     true,
		Title:     content,
		Content:   content,
		Statement: statement,
		Start:     start,
	}}
}

func (db mysqlWalkThroughDatabase) mapErrorNumber(number int, message string, tidb bool) (WalkThroughCode, string) {
	name := quotedName(message)
	switch number {
	case errUnknownDatabase:
		return WalkThroughNotCurrentDatabase, fmt.Sprintf("Database `%s` is not the current database `%s`", name, db.Name)
	case errDupTable:
		return WalkThroughTableExists, fmt.Sprintf("Table `%s` already exists", unqualifiedName(name))
	case errUnknownTable, errNoSuchTable:
		return WalkThroughTableNotExists, fmt.Sprintf("Table `%s` does not exist", unqualifiedName(name))
	case errDupColumn:
		return WalkThroughColumnExists, fmt.Sprintf("Column `%s` already exists", unqualifiedName(name))
	case errNoSuchColumn:
		return WalkThroughColumnNotExists, fmt.Sprintf("Column `%s` does not exist", unqualifiedName(name))
	case errKeyColumnDoesNotExist:
		if tidb {
			return WalkThroughColumnNotExists, fmt.Sprintf("Column `%s` does not exist", strings.TrimPrefix(message, "column does not exist: "))
		}
	case errInvalidDefault:
		return WalkThroughInvalidColumnDefault, fmt.Sprintf("Invalid default value for column `%s`", unqualifiedName(name))
	case errBlobCantHaveDefault:
		if tidb {
			return WalkThroughInvalidColumnDefault, fmt.Sprintf("BLOB, TEXT or JSON column `%s` can't have a default value", unqualifiedName(name))
		}
	case errDupKeyName, errDupIndex, errDupEntry:
		return WalkThroughIndexExists, fmt.Sprintf("Index `%s` already exists", unqualifiedName(name))
	case errMultiplePriKey:
		if tidb {
			return WalkThroughPrimaryKeyExists, "Multiple primary key defined"
		}
		return WalkThroughPrimaryKeyExists, message
	case errSpatialIndexKeyNullable:
		if !tidb {
			return WalkThroughSpatialIndexKeyNullable, message
		}
	case errCantDropKey:
		if tidb && strings.HasPrefix(message, "Can't DROP") {
			return WalkThroughColumnNotExists, fmt.Sprintf("Column `%s` does not exist", tidbCantDropName(message))
		}
		if tidb {
			return WalkThroughIndexNotExists, fmt.Sprintf("Index `%s` does not exist", strings.TrimSuffix(strings.TrimPrefix(message, "index "), " doesn't exist"))
		}
		return WalkThroughIndexNotExists, fmt.Sprintf("Index `%s` does not exist", unqualifiedName(name))
	case errIncorrectTableDefinition:
		if tidb {
			return WalkThroughAutoIncrementExists, "There can be only one auto column and it must be defined as a key"
		}
		return WalkThroughAutoIncrementExists, message
	case errInvalidOnUpdate:
		if tidb {
			return WalkThroughOnUpdateColumnNotDatetimeOrTimestamp, fmt.Sprintf("Column `%s` use ON UPDATE but is not DATETIME or TIMESTAMP", unqualifiedName(name))
		}
		return WalkThroughOnUpdateColumnNotDatetimeOrTimestamp, message
	case errFKNoRefTable, errFKCannotDropParent, errFKMissingIndex, errFKIncompatibleColumns:
		if !tidb {
			return WalkThroughTableHasFK, message
		}
	default:
	}
	return WalkThroughDDLSimulationFailed, message
}

// quotedName returns the first '...' quoted text of a message, or the message.
func quotedName(message string) string {
	start := strings.Index(message, "'")
	if start < 0 {
		return message
	}
	end := strings.Index(message[start+1:], "'")
	if end < 0 {
		return message[start+1:]
	}
	return message[start+1 : start+1+end]
}

func unqualifiedName(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[i+1:]
	}
	return name
}

// tidbCantDropName extracts <col> from "Can't DROP '<col>'; check that ...".
func tidbCantDropName(message string) string {
	rest, ok := strings.CutPrefix(message, "Can't DROP '")
	if !ok {
		return message
	}
	if i := strings.IndexByte(rest, '\''); i >= 0 {
		return rest[:i]
	}
	return rest
}
