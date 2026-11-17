package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// systemSchemas are the schemas every database holds and no user owns.
var systemSchemas = map[string]bool{
	"information_schema": true,
}

// tableKey names a table or a view of one schema, which a sync reads columns
// for before it reads the object itself.
type tableKey struct {
	schema string
	table  string
}

// SyncInstance reads the account's version, databases and users.
//
// The user list is SNOWFLAKE.ACCOUNT_USAGE.USERS with GRANTS_TO_USERS beside it,
// so the login needs a role that can read ACCOUNT_USAGE — Bytebase's instance
// sync read the same two views and failed without one — and the view lags the
// account by up to an hour, so a user created just now may not be listed.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}
	databases, err := d.getDatabaseNames(ctx)
	if err != nil {
		return nil, err
	}

	metadata := make([]*metadatapb.DatabaseSchemaMetadata, 0, len(databases))
	for _, database := range databases {
		metadata = append(metadata, &metadatapb.DatabaseSchemaMetadata{Name: database})
	}
	return &driver.Instance{
		Version:   version,
		Databases: metadata,
		Roles:     roles,
	}, nil
}

// SyncDatabase reads the schema of the database the driver is open on: its
// schemas and the tables, views, streams and tasks in them.
//
// The database is looked up in the account's database list first, so a database
// that is not there is driver.ErrNotFound. A client that validates the
// configured database when it connects refuses a database that does not exist
// before this runs, so that check answers for a database dropped or unshared
// after the driver connected.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	databases, err := d.getDatabaseNames(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(databases, d.databaseName) {
		return nil, fmt.Errorf("database %q: %w", d.databaseName, driver.ErrNotFound)
	}

	schemaList, err := d.getSchemaList(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	tableMap, viewMap, err := d.getTableSchema(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	streamMap, err := d.getStreamSchema(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	taskMap, err := d.getTaskSchema(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}

	databaseMetadata := &metadatapb.DatabaseSchemaMetadata{Name: d.databaseName}
	for _, schemaName := range schemaList {
		databaseMetadata.Schemas = append(databaseMetadata.Schemas, &metadatapb.SchemaMetadata{
			Name:    schemaName,
			Tables:  tableMap[schemaName],
			Views:   viewMap[schemaName],
			Streams: streamMap[schemaName],
			Tasks:   taskMap[schemaName],
		})
	}
	return databaseMetadata, nil
}

// getSchemaList reads the database's schemas, without the system ones.
func (d *Driver) getSchemaList(ctx context.Context, database string) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT
			SCHEMA_NAME
		FROM %s.SCHEMATA
		WHERE %s
		ORDER BY SCHEMA_NAME`, informationSchema(database), excludedSchemaPredicate("SCHEMA_NAME"))

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		result = append(result, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return result, nil
}

// excludedSchemaPredicate writes the WHERE clause that leaves the system
// schemas out of a read of column columnName.
func excludedSchemaPredicate(columnName string) string {
	excluded := make([]string, 0, len(systemSchemas))
	for name := range systemSchemas {
		excluded = append(excluded, sqlLiteral(name))
	}
	slices.Sort(excluded)
	return fmt.Sprintf("LOWER(%s) NOT IN (%s)", columnName, strings.Join(excluded, ", "))
}

// getTableSchema reads a database's tables with their columns, and its views
// with the columns they read.
func (d *Driver) getTableSchema(ctx context.Context, database string) (map[string][]*metadatapb.TableMetadata, map[string][]*metadatapb.ViewMetadata, error) {
	tableMap, viewMap := make(map[string][]*metadatapb.TableMetadata), make(map[string][]*metadatapb.ViewMetadata)

	columnMap, err := d.getColumnMap(ctx, database)
	if err != nil {
		return nil, nil, err
	}

	tableQuery := fmt.Sprintf(`
		SELECT
			TABLE_SCHEMA,
			TABLE_NAME,
			ROW_COUNT,
			BYTES,
			IFNULL(COMMENT, '')
		FROM %s.TABLES
		WHERE TABLE_TYPE = 'BASE TABLE' AND %s
		ORDER BY TABLE_SCHEMA, TABLE_NAME`, informationSchema(database), excludedSchemaPredicate("TABLE_SCHEMA"))
	tableRows, err := d.db.QueryContext(ctx, tableQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}
	defer tableRows.Close()
	for tableRows.Next() {
		var schemaName string
		// A table the account has never scanned states no row count and no size.
		var rowCount, dataSize sql.NullInt64
		table := &metadatapb.TableMetadata{}
		if err := tableRows.Scan(
			&schemaName,
			&table.Name,
			&rowCount,
			&dataSize,
			&table.Comment,
		); err != nil {
			return nil, nil, err
		}
		table.RowCount = rowCount.Int64
		table.DataSize = dataSize.Int64
		table.Columns = columnMap[tableKey{schema: schemaName, table: table.Name}]
		tableMap[schemaName] = append(tableMap[schemaName], table)
	}
	if err := tableRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}

	viewQuery := fmt.Sprintf(`
		SELECT
			TABLE_SCHEMA,
			TABLE_NAME,
			IFNULL(VIEW_DEFINITION, ''),
			IFNULL(COMMENT, '')
		FROM %s.VIEWS
		WHERE %s
		ORDER BY TABLE_SCHEMA, TABLE_NAME`, informationSchema(database), excludedSchemaPredicate("TABLE_SCHEMA"))
	viewRows, err := d.db.QueryContext(ctx, viewQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
	defer viewRows.Close()
	for viewRows.Next() {
		view := &metadatapb.ViewMetadata{}
		var schemaName string
		if err := viewRows.Scan(
			&schemaName,
			&view.Name,
			&view.Definition,
			&view.Comment,
		); err != nil {
			return nil, nil, err
		}
		columns := columnMap[tableKey{schema: schemaName, table: view.Name}]
		view.Columns = columns
		for _, column := range columns {
			// INFORMATION_SCHEMA.COLUMNS states a view's columns without the
			// database they came from, and one view column may read another
			// database's, so a dependency column names the view itself.
			view.DependencyColumns = append(view.DependencyColumns, &metadatapb.DependencyColumn{
				Schema: schemaName,
				Table:  view.Name,
				Column: column.Name,
			})
		}
		viewMap[schemaName] = append(viewMap[schemaName], view)
	}
	if err := viewRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
	return tableMap, viewMap, nil
}

// getColumnMap reads every column of the database, by table.
func (d *Driver) getColumnMap(ctx context.Context, database string) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	columnQuery := fmt.Sprintf(`
		SELECT
			TABLE_SCHEMA,
			TABLE_NAME,
			IFNULL(COLUMN_NAME, ''),
			ORDINAL_POSITION,
			COLUMN_DEFAULT,
			IS_NULLABLE,
			DATA_TYPE,
			IFNULL(CHARACTER_SET_NAME, ''),
			IFNULL(COLLATION_NAME, ''),
			IFNULL(COMMENT, '')
		FROM %s.COLUMNS
		WHERE %s
		ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION`, informationSchema(database), excludedSchemaPredicate("TABLE_SCHEMA"))
	columnRows, err := d.db.QueryContext(ctx, columnQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}
	defer columnRows.Close()

	columnMap := make(map[tableKey][]*metadatapb.ColumnMetadata)
	for columnRows.Next() {
		var schemaName, tableName, nullable string
		var defaultStr sql.NullString
		column := &metadatapb.ColumnMetadata{}
		if err := columnRows.Scan(
			&schemaName,
			&tableName,
			&column.Name,
			&column.Position,
			&defaultStr,
			&nullable,
			&column.Type,
			&column.CharacterSet,
			&column.Collation,
			&column.Comment,
		); err != nil {
			return nil, err
		}
		column.Default = defaultStr.String
		isNullBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		column.Nullable = isNullBool

		key := tableKey{schema: schemaName, table: tableName}
		columnMap[key] = append(columnMap[key], column)
	}
	if err := columnRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}
	return columnMap, nil
}

// getStreamSchema reads a database's streams, with the definition the server
// states for each.
func (d *Driver) getStreamSchema(ctx context.Context, database string) (map[string][]*metadatapb.StreamMetadata, error) {
	streamMap := make(map[string][]*metadatapb.StreamMetadata)

	streamQuery := fmt.Sprintf(`SHOW STREAMS IN DATABASE %s`, quoteIdent(database))
	streamRows, err := d.db.QueryContext(ctx, streamQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", streamQuery, err)
	}
	defer streamRows.Close()

	columns, err := streamRows.Columns()
	if err != nil {
		return nil, fmt.Errorf("cannot get streams from %q query: %w", streamQuery, err)
	}
	// https://docs.snowflake.com/en/sql-reference/sql/show-streams#output
	streamNameIndex := columnIndex(columns, "name")
	schemaNameIndex := columnIndex(columns, "schema_name")
	ownerIndex := columnIndex(columns, "owner")
	commentIndex := columnIndex(columns, "comment")
	tableNameIndex := columnIndex(columns, "table_name")
	typeIndex := columnIndex(columns, "type")
	staleIndex := columnIndex(columns, "stale")
	modeIndex := columnIndex(columns, "mode")
	for streamRows.Next() {
		// Every column is read as nullable: a stream states no comment, and a
		// server may leave a column this does not know empty.
		var streamName, schemaName, owner, comment, tableName, streamType, mode sql.NullString
		var stale sql.NullBool
		var unused any
		cols := make([]any, len(columns))
		for i := range cols {
			cols[i] = &unused
		}
		setColumn(cols, streamNameIndex, &streamName)
		setColumn(cols, schemaNameIndex, &schemaName)
		setColumn(cols, ownerIndex, &owner)
		setColumn(cols, commentIndex, &comment)
		setColumn(cols, tableNameIndex, &tableName)
		setColumn(cols, typeIndex, &streamType)
		setColumn(cols, staleIndex, &stale)
		setColumn(cols, modeIndex, &mode)
		if err := streamRows.Scan(cols...); err != nil {
			return nil, err
		}
		streamMap[schemaName.String] = append(streamMap[schemaName.String], &metadatapb.StreamMetadata{
			Name:      streamName.String,
			TableName: tableName.String,
			Owner:     owner.String,
			Comment:   comment.String,
			Type:      streamTypeOf(streamType.String),
			Stale:     stale.Bool,
			Mode:      streamModeOf(mode.String),
		})
	}
	if err := streamRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", streamQuery, err)
	}

	for schemaName, streams := range streamMap {
		for _, stream := range streams {
			definition, err := d.getObjectDefinition(ctx, "STREAM", database, schemaName, stream.Name)
			if err != nil {
				return nil, err
			}
			stream.Definition = definition
		}
		slices.SortFunc(streams, func(a, b *metadatapb.StreamMetadata) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	return streamMap, nil
}

// setColumn points the scan at index at dest, and leaves it unread when the
// server's answer states no such column: a value read off the wrong column is
// worse than a field left empty.
func setColumn(cols []any, index int, dest any) {
	if index >= 0 && index < len(cols) {
		cols[index] = dest
	}
}

// streamTypeOf maps what SHOW STREAMS states into the metadata's type.
func streamTypeOf(streamType string) metadatapb.StreamMetadata_Type {
	if streamType == "DELTA" {
		return metadatapb.StreamMetadata_TYPE_DELTA
	}
	return metadatapb.StreamMetadata_TYPE_UNSPECIFIED
}

// streamModeOf maps what SHOW STREAMS states into the metadata's mode.
func streamModeOf(mode string) metadatapb.StreamMetadata_Mode {
	switch mode {
	case "DEFAULT":
		return metadatapb.StreamMetadata_MODE_DEFAULT
	case "APPEND_ONLY":
		return metadatapb.StreamMetadata_MODE_APPEND_ONLY
	case "INSERT_ONLY":
		return metadatapb.StreamMetadata_MODE_INSERT_ONLY
	default:
		return metadatapb.StreamMetadata_MODE_UNSPECIFIED
	}
}

// getTaskSchema reads a database's tasks, with the definition the server states
// for each.
func (d *Driver) getTaskSchema(ctx context.Context, database string) (map[string][]*metadatapb.TaskMetadata, error) {
	taskMap := make(map[string][]*metadatapb.TaskMetadata)

	taskQuery := fmt.Sprintf(`SHOW TASKS IN DATABASE %s`, quoteIdent(database))
	taskRows, err := d.db.QueryContext(ctx, taskQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", taskQuery, err)
	}
	defer taskRows.Close()

	columns, err := taskRows.Columns()
	if err != nil {
		return nil, fmt.Errorf("cannot get tasks from %q query: %w", taskQuery, err)
	}
	// https://docs.snowflake.com/en/sql-reference/sql/show-tasks#output
	taskNameIndex := columnIndex(columns, "name")
	taskIDIndex := columnIndex(columns, "id")
	schemaNameIndex := columnIndex(columns, "schema_name")
	ownerIndex := columnIndex(columns, "owner")
	commentIndex := columnIndex(columns, "comment")
	warehouseIndex := columnIndex(columns, "warehouse")
	scheduleIndex := columnIndex(columns, "schedule")
	predecessorsIndex := columnIndex(columns, "predecessors")
	stateIndex := columnIndex(columns, "state")
	conditionIndex := columnIndex(columns, "condition")
	for taskRows.Next() {
		// A task states no warehouse, schedule or condition when it uses the
		// session's, is manually triggered or runs unconditionally, so every
		// column is read as nullable.
		var taskName, taskID, schemaName, owner, comment, state sql.NullString
		var warehouse, schedule, condition sql.NullString
		var predecessors ArrayString
		var unused any
		cols := make([]any, len(columns))
		for i := range cols {
			cols[i] = &unused
		}
		setColumn(cols, taskNameIndex, &taskName)
		setColumn(cols, taskIDIndex, &taskID)
		setColumn(cols, schemaNameIndex, &schemaName)
		setColumn(cols, ownerIndex, &owner)
		setColumn(cols, commentIndex, &comment)
		setColumn(cols, warehouseIndex, &warehouse)
		setColumn(cols, scheduleIndex, &schedule)
		setColumn(cols, predecessorsIndex, &predecessors)
		setColumn(cols, stateIndex, &state)
		setColumn(cols, conditionIndex, &condition)
		if err := taskRows.Scan(cols...); err != nil {
			return nil, err
		}
		taskMap[schemaName.String] = append(taskMap[schemaName.String], &metadatapb.TaskMetadata{
			Name:         taskName.String,
			Id:           taskID.String,
			Owner:        owner.String,
			Comment:      comment.String,
			Warehouse:    warehouse.String,
			Schedule:     schedule.String,
			Predecessors: predecessors,
			State:        taskStateOf(state.String),
			Condition:    condition.String,
		})
	}
	if err := taskRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", taskQuery, err)
	}

	for schemaName, tasks := range taskMap {
		for _, task := range tasks {
			definition, err := d.getObjectDefinition(ctx, "TASK", database, schemaName, task.Name)
			if err != nil {
				return nil, err
			}
			task.Definition = definition
		}
		slices.SortFunc(tasks, func(a, b *metadatapb.TaskMetadata) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	return taskMap, nil
}

// taskStateOf maps what SHOW TASKS states into the metadata's state.
func taskStateOf(state string) metadatapb.TaskMetadata_State {
	switch state {
	case "started":
		return metadatapb.TaskMetadata_STATE_STARTED
	case "suspended":
		return metadatapb.TaskMetadata_STATE_SUSPENDED
	default:
		return metadatapb.TaskMetadata_STATE_UNSPECIFIED
	}
}

// getObjectDefinition reads the DDL the server states for one object.
func (d *Driver) getObjectDefinition(ctx context.Context, objectType, database, schema, name string) (string, error) {
	qualified := fmt.Sprintf("%s.%s.%s", quoteIdent(database), quoteIdent(schema), quoteIdent(name))
	query := fmt.Sprintf("SELECT GET_DDL(%s, %s, TRUE)", sqlLiteral(objectType), sqlLiteral(qualified))
	var definition string
	if err := d.db.QueryRowContext(ctx, query).Scan(&definition); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return definition, nil
}

// ArrayString scans the JSON array of strings a SHOW TASKS answer writes a
// task's predecessors as.
type ArrayString []string

// Scan implements the sql.Scanner interface.
func (a *ArrayString) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		// A task that is no other task's predecessor states an empty array, and
		// one the server left NULL has none.
		*a = nil
		return nil
	case string:
		return json.Unmarshal([]byte(v), a)
	case []byte:
		return json.Unmarshal(v, a)
	default:
		return errors.New("invalid type")
	}
}

// convertYesNo reads the YES/NO answer INFORMATION_SCHEMA.COLUMNS gives for a
// column's nullability.
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

// informationSchema is the database's INFORMATION_SCHEMA as a path: the
// database as one quoted identifier and the schema's own name, which is what a
// query names. Quoting it here keeps the caller's name one identifier and does
// not quote it twice.
func informationSchema(database string) string {
	return quoteIdent(database) + ".INFORMATION_SCHEMA"
}
