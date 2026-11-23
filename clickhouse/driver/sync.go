package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"

	chparser "github.com/bytebase/omni/clickhouse/parser"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

var (
	systemDatabases = map[string]bool{
		"system":             true,
		"information_schema": true,
		"INFORMATION_SCHEMA": true,
	}
	// unknownDatabase is ClickHouse's UNKNOWN_DATABASE.
	unknownDatabase int32 = 81
	// systemDatabaseClause is built in sorted order so a recording's query
	// text does not move with Go's map iteration.
	systemDatabaseClause = func() string {
		names := make([]string, 0, len(systemDatabases))
		for name := range systemDatabases {
			names = append(names, name)
		}
		sort.Strings(names)
		quoted := make([]string, 0, len(names))
		for _, name := range names {
			quoted = append(quoted, fmt.Sprintf("'%s'", name))
		}
		return strings.Join(quoted, ", ")
	}()
)

// SyncInstance syncs the instance: the server version, its user databases
// without their tables, and its users and roles, each with the access
// statements the server states for it (getInstanceRoles).
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		if d.engineVersion == "" {
			return nil, err
		}
		version = d.engineVersion
	}

	instanceRoles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}

	var databases []*metadatapb.DatabaseSchemaMetadata
	where := fmt.Sprintf("schema_name NOT IN (%s)", systemDatabaseClause)
	query := `
		SELECT
			schema_name
		FROM information_schema.schemata
		WHERE ` + where + `
		ORDER BY schema_name`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	for rows.Next() {
		database := &metadatapb.DatabaseSchemaMetadata{}
		if err := rows.Scan(&database.Name); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &driver.Instance{
		Version:   version,
		Databases: databases,
		Roles:     instanceRoles,
	}, nil
}

// SyncDatabase syncs the database's schema: one schema holding its tables and
// views, each with its columns, data-skipping indexes and sorting keys. A
// database that is not there is driver.ErrNotFound.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	metadata, err := d.syncDatabase(ctx)
	if err != nil {
		return nil, wrapDatabaseNotFound(err)
	}
	return metadata, nil
}

// wrapDatabaseNotFound reports a database ClickHouse says it does not have as
// driver.ErrNotFound: a sync and a dump name the database they read, so its
// absence is the object the call asked for and not a failure of the read. The
// client answers UNKNOWN_DATABASE for the database its connection opens on.
func wrapDatabaseNotFound(err error) error {
	var engineError *clickhouse.Exception
	if errors.As(err, &engineError) && engineError.Code == unknownDatabase {
		return fmt.Errorf("%w: %s", driver.ErrNotFound, engineError.Message)
	}
	return err
}

func (d *Driver) syncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	databaseName, err := d.syncDatabaseName(ctx)
	if err != nil {
		return nil, err
	}
	schemaMetadata := &metadatapb.SchemaMetadata{
		Name: "",
	}

	// tableName -> columns
	columnMap := make(map[string][]*metadatapb.ColumnMetadata)
	columnQuery := `
		SELECT
			table_name,
			column_name,
			ordinal_position,
			column_default,
			is_nullable,
			column_type,
			ifNull(character_set_name, ''),
			ifNull(collation_name, ''),
			column_comment
		FROM information_schema.columns
		WHERE table_schema = $1
		ORDER BY table_name, ordinal_position`
	columnRows, err := d.db.QueryContext(ctx, columnQuery, databaseName)
	if err != nil {
		return nil, formatErrorWithQuery(err, columnQuery)
	}
	defer columnRows.Close()
	for columnRows.Next() {
		column := &metadatapb.ColumnMetadata{}
		// https://clickhouse.com/docs/en/operations/system-tables/information_schema#columns
		// defaultValueExpression is an expression for the default value, or an
		// empty string if it is not defined.
		var tableName, nullable, defaultValueExpression string
		if err := columnRows.Scan(
			&tableName,
			&column.Name,
			&column.Position,
			&defaultValueExpression,
			&nullable,
			&column.Type,
			&column.CharacterSet,
			&column.Collation,
			&column.Comment,
		); err != nil {
			return nil, err
		}
		nullableBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		column.Nullable = nullableBool
		if defaultValueExpression == "" {
			column.Default = "NULL"
		} else {
			column.Default = defaultValueExpression
		}
		columnMap[tableName] = append(columnMap[tableName], column)
	}
	if err := columnRows.Err(); err != nil {
		return nil, formatErrorWithQuery(err, columnQuery)
	}
	if err := columnRows.Close(); err != nil {
		return nil, err
	}

	// system.tables is used rather than information_schema.tables because it
	// carries the engine. metadata_modification_time is not selected: Bytebase
	// read it into a field the metadata never held, and its value would make a
	// recording differ from the next.
	tableQuery := `
		SELECT
			name,
			engine,
			ifNull(total_rows, 0),
			ifNull(total_bytes, 0),
			create_table_query,
			comment,
			sorting_key,
			primary_key
		FROM system.tables
		WHERE database = $1
		ORDER BY name`
	type tableRow struct {
		name, engine, definition, comment, sortingKey, primaryKey string
		rowCount, totalBytes                                      int64
	}
	tableRows, err := d.db.QueryContext(ctx, tableQuery, databaseName)
	if err != nil {
		return nil, formatErrorWithQuery(err, tableQuery)
	}
	var tableList []tableRow
	for tableRows.Next() {
		var row tableRow
		if err := tableRows.Scan(
			&row.name,
			&row.engine,
			&row.rowCount,
			&row.totalBytes,
			&row.definition,
			&row.comment,
			&row.sortingKey,
			&row.primaryKey,
		); err != nil {
			_ = tableRows.Close()
			return nil, err
		}
		tableList = append(tableList, row)
	}
	if err := tableRows.Err(); err != nil {
		_ = tableRows.Close()
		return nil, formatErrorWithQuery(err, tableQuery)
	}
	if err := tableRows.Close(); err != nil {
		return nil, err
	}
	for _, row := range tableList {
		// For a view, the engine is "View".
		if row.engine == "View" {
			schemaMetadata.Views = append(schemaMetadata.Views, &metadatapb.ViewMetadata{
				Name:       row.name,
				Columns:    columnMap[row.name],
				Definition: row.definition,
				Comment:    row.comment,
			})
			continue
		}
		table := &metadatapb.TableMetadata{
			Name:     row.name,
			Columns:  columnMap[row.name],
			Engine:   row.engine,
			RowCount: row.rowCount,
			DataSize: row.totalBytes,
			Comment:  row.comment,
		}
		indexes, err := d.getDataSkippingIndices(ctx, databaseName, row.name)
		if err != nil {
			return nil, err
		}
		table.Indexes = indexes
		if row.primaryKey != "" {
			primaryKeys := splitKeyExpressions(row.primaryKey)
			// ClickHouse keeps primary keys in system.tables rather than in an
			// index, so they are recorded as one to fit the metadata design.
			table.Indexes = append(table.Indexes, &metadatapb.IndexMetadata{
				Primary:     true,
				Expressions: primaryKeys,
			})
		}
		if row.sortingKey != "" {
			table.SortingKeys = splitKeyExpressions(row.sortingKey)
		}
		schemaMetadata.Tables = append(schemaMetadata.Tables, table)
	}

	return &metadatapb.DatabaseSchemaMetadata{
		Name:    databaseName,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}, nil
}

// splitKeyExpressions splits ClickHouse's comma-separated sorting key, primary
// key and index expressions at top-level commas only: a function or tuple keeps
// its own commas.
func splitKeyExpressions(s string) []string {
	l := chparser.NewLexer(s)
	var expressions []string
	start, depth := 0, 0
	for {
		tok := l.NextToken()
		if tok.Type == chparser.TokEOF {
			break
		}
		switch tok.Type {
		case chparser.TokLParen, chparser.TokLBracket, chparser.TokLBrace:
			depth++
		case chparser.TokRParen, chparser.TokRBracket, chparser.TokRBrace:
			if depth > 0 {
				depth--
			}
		case chparser.TokComma:
			if depth == 0 {
				if expr := strings.TrimSpace(s[start:tok.Pos]); expr != "" {
					expressions = append(expressions, expr)
				}
				start = tok.End
			}
		}
	}
	if expr := strings.TrimSpace(s[start:]); expr != "" {
		expressions = append(expressions, expr)
	}
	return expressions
}

// syncDatabaseName is Config.Database, or the connection's own database when
// the caller left it empty, which is the engine's default.
func (d *Driver) syncDatabaseName(ctx context.Context) (string, error) {
	if d.databaseName != "" {
		return d.databaseName, nil
	}
	const query = "SELECT currentDatabase()"
	var name string
	if err := d.db.QueryRowContext(ctx, query).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("query %q returned empty row", query)
		}
		return "", formatErrorWithQuery(err, query)
	}
	return name, nil
}

func (d *Driver) getDataSkippingIndices(ctx context.Context, database, table string) ([]*metadatapb.IndexMetadata, error) {
	query := `
		SELECT
			name,
			type,
			expr,
			granularity
		FROM system.data_skipping_indices
		WHERE database = $1 AND table = $2
		ORDER BY name`
	rows, err := d.db.QueryContext(ctx, query, database, table)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()

	var indices []*metadatapb.IndexMetadata
	for rows.Next() {
		var expr string
		index := &metadatapb.IndexMetadata{}
		if err := rows.Scan(
			&index.Name,
			&index.Type,
			&expr,
			&index.Granularity,
		); err != nil {
			return nil, err
		}
		index.Expressions = splitKeyExpressions(expr)
		indices = append(indices, index)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return indices, nil
}

func (d *Driver) getVersion(ctx context.Context) (string, error) {
	const query = "SELECT VERSION()"
	var version string
	if err := d.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("query %q returned empty row", query)
		}
		return "", formatErrorWithQuery(err, query)
	}
	return version, nil
}

// convertYesNo converts ClickHouse's 0/1 and YES/NO spellings to a bool.
func convertYesNo(s string) (bool, error) {
	switch s {
	case "YES", "Y", "1":
		return true, nil
	case "NO", "N", "0":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognized isNullable type %q", s)
	}
}

// formatErrorWithQuery is Bytebase's db/util.FormatErrorWithQuery.
func formatErrorWithQuery(err error, query string) error {
	return fmt.Errorf("failed to execute query %q: %w", query, err)
}
