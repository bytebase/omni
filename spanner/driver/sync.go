package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	metadatapb "github.com/bytebase/omni/metadata"

	"github.com/bytebase/omni/driver"
)

// tableKey is the map key for table metadata: a Spanner table's schema and name.
type tableKey struct {
	Schema string
	Table  string
}

// catalog is the sync query set of one dialect. GoogleSQL and PostgreSQL spell
// the system schemas' case and SPANNER_IS_MANAGED's type differently, and
// PostgreSQL states IS_UNIQUE as text, so the queries are written per dialect;
// every other query is the same text.
type catalog struct {
	tables      string
	columns     string
	indexes     string
	foreignKeys string
	views       string
}

// foreignKeyQuery is the same text for both dialects: TABLE_CONSTRAINTS,
// REFERENTIAL_CONSTRAINTS and KEY_COLUMN_USAGE spell their columns alike.
const foreignKeyQuery = `
    WITH t AS (
      SELECT DISTINCT
        tc.CONSTRAINT_SCHEMA,
        tc.CONSTRAINT_NAME,
        tc.TABLE_NAME,
        rc.UNIQUE_CONSTRAINT_SCHEMA,
        rc.UNIQUE_CONSTRAINT_NAME,
        ccu.TABLE_SCHEMA AS REFERENCED_TABLE_SCHEMA,
        ccu.TABLE_NAME AS REFERENCED_TABLE_NAME,
        rc.DELETE_RULE,
        rc.UPDATE_RULE,
        rc.MATCH_OPTION
      FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS tc
      JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS AS rc
      ON tc.CONSTRAINT_SCHEMA = rc.CONSTRAINT_SCHEMA AND tc.CONSTRAINT_NAME = rc.CONSTRAINT_NAME 
      JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE as ccu
      ON tc.CONSTRAINT_SCHEMA = ccu.CONSTRAINT_SCHEMA AND tc.CONSTRAINT_NAME = ccu.CONSTRAINT_NAME 
      WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY'
    )
    SELECT
      t.CONSTRAINT_SCHEMA,
      t.CONSTRAINT_NAME,
      t.TABLE_NAME,
      t.REFERENCED_TABLE_SCHEMA,
      t.REFERENCED_TABLE_NAME,
      ARRAY (
        SELECT
          COLUMN_NAME
        FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu
        WHERE t.CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA AND t.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
        ORDER BY kcu.ORDINAL_POSITION
      ) AS TABLE_COLUMNS,
      ARRAY (
        SELECT
          COLUMN_NAME
        FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu
        WHERE t.UNIQUE_CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA AND t.UNIQUE_CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
        ORDER BY kcu.ORDINAL_POSITION
      ) AS REFERENCED_TABLE_COLUMNS,
      t.DELETE_RULE,
      t.UPDATE_RULE,
      t.MATCH_OPTION
    FROM t
  `

// catalogFor returns the sync queries of d.
func catalogFor(d dialect) catalog {
	if d == postgresSQL {
		return catalog{
			tables: `
    SELECT
      TABLE_SCHEMA,
      TABLE_NAME
    FROM INFORMATION_SCHEMA.TABLES
    WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys') AND TABLE_TYPE = 'BASE TABLE'
    ORDER BY TABLE_SCHEMA, TABLE_NAME
  `,
			columns: `
    SELECT 
      TABLE_SCHEMA,
      TABLE_NAME,
      COLUMN_NAME,
      ORDINAL_POSITION,
      COLUMN_DEFAULT,
      IS_NULLABLE = 'YES',
      SPANNER_TYPE
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys')
    ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION
  `,
			indexes: `
    SELECT
      TABLE_SCHEMA,
      TABLE_NAME,
      INDEX_NAME,
      IS_UNIQUE = 'YES' AS IS_UNIQUE,
      INDEX_TYPE = 'PRIMARY_KEY' AS IS_PRIMARY,
      ARRAY (
        SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.INDEX_COLUMNS AS ic
        WHERE ic.TABLE_SCHEMA = i.TABLE_SCHEMA AND ic.TABLE_NAME = i.TABLE_NAME AND ic.INDEX_NAME = i.INDEX_NAME
        ORDER BY ic.ORDINAL_POSITION
      )
    FROM INFORMATION_SCHEMA.INDEXES as i
    WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys') AND SPANNER_IS_MANAGED = 'NO'
    ORDER BY TABLE_SCHEMA, TABLE_NAME, INDEX_NAME
  `,
			foreignKeys: foreignKeyQuery,
			views: `
  SELECT
    TABLE_SCHEMA,
    TABLE_NAME,
    VIEW_DEFINITION
  FROM INFORMATION_SCHEMA.VIEWS
  WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys')
  `,
		}
	}
	return catalog{
		tables: `
    SELECT
      TABLE_SCHEMA,
      TABLE_NAME
    FROM INFORMATION_SCHEMA.TABLES
    WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS') AND TABLE_TYPE = 'BASE TABLE'
    ORDER BY TABLE_SCHEMA, TABLE_NAME
  `,
		columns: `
    SELECT 
      TABLE_SCHEMA,
      TABLE_NAME,
      COLUMN_NAME,
      ORDINAL_POSITION,
      COLUMN_DEFAULT,
      IS_NULLABLE = 'YES',
      SPANNER_TYPE
    FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS')
    ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION
  `,
		indexes: `
    SELECT
      TABLE_SCHEMA,
      TABLE_NAME,
      INDEX_NAME,
      IS_UNIQUE,
      INDEX_TYPE = 'PRIMARY_KEY' AS IS_PRIMARY,
      ARRAY (
        SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.INDEX_COLUMNS AS ic
        WHERE ic.TABLE_SCHEMA = i.TABLE_SCHEMA AND ic.TABLE_NAME = i.TABLE_NAME AND ic.INDEX_NAME = i.INDEX_NAME
        ORDER BY ic.ORDINAL_POSITION
      )
    FROM INFORMATION_SCHEMA.INDEXES as i
    WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS') AND SPANNER_IS_MANAGED = FALSE
    ORDER BY TABLE_SCHEMA, TABLE_NAME, INDEX_NAME
  `,
		foreignKeys: foreignKeyQuery,
		views: `
  SELECT
    TABLE_SCHEMA,
    TABLE_NAME,
    VIEW_DEFINITION
  FROM INFORMATION_SCHEMA.VIEWS
  WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS')
  `,
	}
}

// SyncInstance reads the instance's databases with their names and nothing else:
// Spanner reports no version and no roles, and a database's schema is a
// SyncDatabase of its own. Every database is listed, whichever dialect it
// states: the dialect decides the SQL a caller may run, not whether the
// database is a database.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	databases, err := d.client.ListDatabases(ctx)
	if err != nil {
		return nil, err
	}
	metadata := make([]*metadatapb.DatabaseSchemaMetadata, 0, len(databases))
	for _, database := range databases {
		// database.Name is of the form
		// `projects/<project>/instances/<instance>/databases/<database>`, and a
		// sync reports the last part.
		databaseName, err := getDatabaseFromDSN(database.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get database name from %s: %w", database.Name, err)
		}
		metadata = append(metadata, &metadatapb.DatabaseSchemaMetadata{Name: databaseName})
	}
	return &driver.Instance{Databases: metadata}, nil
}

// SyncDatabase reads the schema of Config.Database: its tables with columns,
// indexes and foreign keys, and its views, from the INFORMATION_SCHEMA of one
// consistent snapshot. A database that is not there is driver.ErrNotFound.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	database, err := d.client.GetDatabase(ctx, d.databasePath(""))
	if err != nil {
		// A database the admin API does not have is driver.ErrNotFound, whatever
		// the client wrapped the server's answer in.
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("database %q not found: %w", d.databaseName, driver.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	tx, err := d.client.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	// A PostgreSQL database answers a catalog of its own spelling, and the sync
	// asks the one the server just named. The answer is also what a query splits
	// by afterwards, so the driver records it: Open may have found no database to
	// ask, and a sync is the later truth.
	parsed := dialectOf(database.Dialect)
	d.setDialect(parsed)
	queries := catalogFor(parsed)
	databaseMetadata := &metadatapb.DatabaseSchemaMetadata{
		Name: d.databaseName,
	}
	columnMap, err := getColumn(ctx, tx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns from database %q: %w", d.databaseName, err)
	}
	tableMap, err := getTable(ctx, tx, columnMap, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables from database %q: %w", d.databaseName, err)
	}
	viewMap, err := getView(ctx, tx, columnMap, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to get views from database %q: %w", d.databaseName, err)
	}

	// A schema is reported once, sorted by name, so two syncs of one database
	// write the same metadata.
	schemaNameSet := make(map[string]bool)
	for schemaName := range tableMap {
		schemaNameSet[schemaName] = true
	}
	for schemaName := range viewMap {
		schemaNameSet[schemaName] = true
	}
	schemaNames := make([]string, 0, len(schemaNameSet))
	for schemaName := range schemaNameSet {
		schemaNames = append(schemaNames, schemaName)
	}
	slices.Sort(schemaNames)
	for _, schemaName := range schemaNames {
		databaseMetadata.Schemas = append(databaseMetadata.Schemas, &metadatapb.SchemaMetadata{
			Name:   schemaName,
			Tables: tableMap[schemaName],
			Views:  viewMap[schemaName],
		})
	}

	return databaseMetadata, nil
}

func getTable(ctx context.Context, tx Reader, columnMap map[tableKey][]*metadatapb.ColumnMetadata, queries catalog) (map[string][]*metadatapb.TableMetadata, error) {
	indexMap, err := getIndex(ctx, tx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices: %w", err)
	}
	foreignKeyMap, err := getForeignKey(ctx, tx, queries)
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign keys: %w", err)
	}
	tableMap := make(map[string][]*metadatapb.TableMetadata)
	iter, err := tx.Query(ctx, queries.tables)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	for {
		cells, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		var table metadatapb.TableMetadata
		var schema string
		if err := scan(cells, &schema, &table.Name); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schema, Table: table.Name}
		table.Columns = columnMap[key]
		table.Indexes = indexMap[key]
		table.ForeignKeys = foreignKeyMap[key]

		tableMap[schema] = append(tableMap[schema], &table)
	}
	return tableMap, nil
}

func getColumn(ctx context.Context, tx Reader, queries catalog) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	columnsMap := make(map[tableKey][]*metadatapb.ColumnMetadata)
	iter, err := tx.Query(ctx, queries.columns)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	for {
		cells, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		column := &metadatapb.ColumnMetadata{}
		var (
			schemaName string
			tableName  string
			position   int64
			defaultStr nullString
		)
		if err := scan(cells, &schemaName, &tableName, &column.Name, &position, &defaultStr, &column.Nullable, &column.Type); err != nil {
			return nil, err
		}
		column.Position = int32(position)
		if defaultStr.Valid {
			column.Default = defaultStr.StringVal
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		columnsMap[key] = append(columnsMap[key], column)
	}
	return columnsMap, nil
}

func getIndex(ctx context.Context, tx Reader, queries catalog) (map[tableKey][]*metadatapb.IndexMetadata, error) {
	keyIndexMap := make(map[tableKey][]*metadatapb.IndexMetadata)

	iter, err := tx.Query(ctx, queries.indexes)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	for {
		cells, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		var idx metadatapb.IndexMetadata
		var schema, table string
		if err := scan(cells, &schema, &table, &idx.Name, &idx.Unique, &idx.Primary, &idx.Expressions); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schema, Table: table}
		keyIndexMap[key] = append(keyIndexMap[key], &idx)
	}

	return keyIndexMap, nil
}

func getForeignKey(ctx context.Context, tx Reader, queries catalog) (map[tableKey][]*metadatapb.ForeignKeyMetadata, error) {
	foreignKeyMap := make(map[tableKey][]*metadatapb.ForeignKeyMetadata)
	iter, err := tx.Query(ctx, queries.foreignKeys)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	for {
		cells, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		var fk metadatapb.ForeignKeyMetadata
		var schema, table string
		var columns, referencedColumns []string
		if err := scan(cells,
			&schema,
			&fk.Name,
			&table,
			&fk.ReferencedSchema,
			&fk.ReferencedTable,
			&columns,
			&referencedColumns,
			&fk.OnDelete,
			&fk.OnUpdate,
			&fk.MatchType,
		); err != nil {
			return nil, err
		}
		// The foreign key's own columns and the referenced ones are read but not
		// reported: the metadata holds the referenced table and the rules, not
		// the column pairs, which is what Bytebase's Spanner sync kept.
		_ = columns
		_ = referencedColumns
		key := tableKey{Schema: schema, Table: table}
		foreignKeyMap[key] = append(foreignKeyMap[key], &fk)
	}
	return foreignKeyMap, nil
}

func getView(ctx context.Context, tx Reader, columnMap map[tableKey][]*metadatapb.ColumnMetadata, queries catalog) (map[string][]*metadatapb.ViewMetadata, error) {
	viewMap := make(map[string][]*metadatapb.ViewMetadata)
	iter, err := tx.Query(ctx, queries.views)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	for {
		cells, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		var (
			schema     string
			name       string
			definition string
		)
		if err := scan(cells, &schema, &name, &definition); err != nil {
			return nil, err
		}

		key := tableKey{Schema: schema, Table: name}

		viewMap[schema] = append(viewMap[schema], &metadatapb.ViewMetadata{
			Name:       name,
			Definition: definition,
			Columns:    columnMap[key],
		})
	}
	return viewMap, nil
}
