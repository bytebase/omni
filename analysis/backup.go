package analysis

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytebase/omni/metadata"
)

// BackupRequest is a script whose DML statements are backed up before they
// run.
type BackupRequest struct {
	// DefaultDatabase is the database the script runs in.
	DefaultDatabase string
	// SourceDatabase is the database the backed-up tables are read from.
	SourceDatabase string
	// BackupDatabase holds the backup tables: a database, or a schema for an
	// engine that backs up into a schema.
	BackupDatabase string
	// TablePrefix starts each backup table's name.
	TablePrefix   string
	CaseSensitive bool
	// Version is the server version as the engine's own type, or nil.
	Version any
	Catalog Catalog
	// ReferencingDatabases lists the databases that may hold a foreign key
	// into the named one, or is nil to read every database.
	ReferencingDatabases func(database string) ([]string, error)
}

// Backup is a statement that saves the rows one or more statements change.
type Backup struct {
	Statement    string
	SourceSchema string
	SourceTable  string
	BackupTable  string
	// Start and End are where the backed-up statements lie in the script, in
	// the convention the engine's RestoreSQL reads back, and zero when the
	// engine cannot place them.
	Start Position
	End   Position
}

// BackupGenerator returns the backups a script's DML statements need.
type BackupGenerator interface {
	Backups(ctx context.Context, request BackupRequest, script string) ([]Backup, error)
}

// RestoreRequest is one backup to restore.
type RestoreRequest struct {
	// Script is the script the backup was taken for.
	Script         string
	SourceDatabase string
	SourceSchema   string
	SourceTable    string
	BackupDatabase string
	BackupSchema   string
	BackupTable    string
	// Start and End are the backup's positions, as Backups returned them.
	Start         Position
	End           Position
	CaseSensitive bool
	Catalog       Catalog
}

// RestoreGenerator returns the SQL that restores a backup's rows.
type RestoreGenerator interface {
	RestoreSQL(ctx context.Context, request RestoreRequest) (string, error)
}

// RestoreTable is the synced table a restore writes rows back into.
type RestoreTable struct {
	Schema     string
	Table      *metadata.TableMetadata
	PrimaryKey *metadata.IndexMetadata
}

// LookUpRestoreTable finds the restore's source table in its synced database.
// An engine with schemas passes the schema an unqualified table is in; an
// engine without passes "" and reads the unnamed schema. notSynced, when not
// nil, makes the error for an unsynced database. Index names follow
// detailCaseSensitive.
func LookUpRestoreTable(ctx context.Context, request RestoreRequest, defaultSchema string, detailCaseSensitive bool, notSynced func(database string) error) (RestoreTable, error) {
	if request.Catalog.GetDatabase == nil {
		return RestoreTable{}, errors.New("the catalog cannot look databases up")
	}
	database, err := request.Catalog.GetDatabase(ctx, request.SourceDatabase)
	if err != nil {
		return RestoreTable{}, fmt.Errorf("failed to get database metadata for %s: %w", request.SourceDatabase, err)
	}
	if database == nil {
		if notSynced != nil {
			return RestoreTable{}, notSynced(request.SourceDatabase)
		}
		return RestoreTable{}, fmt.Errorf("database metadata for database %q not found (database not synced)", request.SourceDatabase)
	}

	schemaName, scope, missingPrefix, missingSuffix := "", request.SourceDatabase, "no ", ""
	if defaultSchema != "" {
		schemaName = request.SourceSchema
		if schemaName == "" {
			schemaName = defaultSchema
		}
		scope, missingPrefix, missingSuffix = schemaName, "", " not found"
	}
	schema := database.Lookup(detailCaseSensitive).Schema(schemaName)
	if schema == nil {
		return RestoreTable{}, fmt.Errorf("%sschema metadata%s for %s", missingPrefix, missingSuffix, scope)
	}
	table := schema.Table(request.SourceTable)
	if table == nil {
		return RestoreTable{}, fmt.Errorf("%stable metadata%s for %s.%s", missingPrefix, missingSuffix, scope, request.SourceTable)
	}
	return RestoreTable{Schema: schemaName, Table: table, PrimaryKey: primaryKey(table, detailCaseSensitive)}, nil
}

// primaryKey returns the table's primary key. Index names are unique under
// detailCaseSensitive, a later index replacing an earlier one of the same name,
// as the caller's model has it.
func primaryKey(table *metadata.TableMetadata, detailCaseSensitive bool) *metadata.IndexMetadata {
	indexes := map[string]*metadata.IndexMetadata{}
	for _, index := range table.GetIndexes() {
		indexes[fold(index.GetName(), detailCaseSensitive)] = index
	}
	for _, index := range table.GetIndexes() {
		if index.GetPrimary() && indexes[fold(index.GetName(), detailCaseSensitive)] == index {
			return index
		}
	}
	return nil
}

// InRange reports whether [start, end] lies within [from, to], comparing
// lines, then columns.
func InRange(start, end, from, to Position) bool {
	return !before(start, from) && !before(to, end)
}

func before(a, b Position) bool {
	return a.Line < b.Line || a.Line == b.Line && a.Column < b.Column
}

// StatementNotFoundError reports that no statement lies in the backup's range.
func StatementNotFoundError(request RestoreRequest) error {
	return fmt.Errorf("could not find statement at position (line %d:%d - %d:%d)",
		request.Start.Line, request.Start.Column, request.End.Line, request.End.Column)
}

// TruncateComment cuts sql to limit characters for a restore's comment,
// marking a cut with "...".
func TruncateComment(sql string, limit int) string {
	characters := 0
	for i := range sql {
		if characters >= limit {
			return sql[:i] + "..."
		}
		characters++
	}
	return sql
}
