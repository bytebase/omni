package driver

import (
	"context"
	"fmt"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

// tableKey identifies a table in one schema.
type tableKey struct {
	Schema string
	Table  string
}

// informationSchema is the name Trino gives every catalog's own metadata schema.
const informationSchema = "information_schema"

// SyncDatabase reads the full schema of Config.Database, as Trino reports it:
// every schema of the catalog, with its tables, views and their columns. The
// catalog's information_schema is the server's own metadata, not a schema of
// the catalog's data, so it is left out.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	schemaNames, err := d.getSchemaList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema list for catalog %s: %w", d.databaseName, err)
	}
	allTables, err := d.fetchAllTablesForCatalog(ctx, schemaNames)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tables for catalog %s: %w", d.databaseName, err)
	}
	allViews, err := d.fetchAllViewsForCatalog(ctx, schemaNames)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch views for catalog %s: %w", d.databaseName, err)
	}
	allColumns, err := d.fetchAllColumnsForCatalog(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns for catalog %s: %w", d.databaseName, err)
	}

	dbMetadata := &metadatapb.DatabaseSchemaMetadata{Name: d.databaseName}
	for _, schemaName := range schemaNames {
		tables := allTables[schemaName]
		for _, table := range tables {
			if columns, ok := allColumns[tableKey{Schema: schemaName, Table: table.Name}]; ok {
				table.Columns = columns
			}
		}
		views := allViews[schemaName]
		for _, view := range views {
			if columns, ok := allColumns[tableKey{Schema: schemaName, Table: view.Name}]; ok {
				view.Columns = columns
			}
		}
		dbMetadata.Schemas = append(dbMetadata.Schemas, &metadatapb.SchemaMetadata{
			Name:   schemaName,
			Tables: tables,
			Views:  views,
		})
	}
	return dbMetadata, nil
}

func (d *Driver) getSchemaList(ctx context.Context) ([]string, error) {
	schemas, err := d.queryStringValues(ctx,
		"SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = ? ORDER BY table_schem",
		d.databaseName)
	if err != nil {
		return nil, err
	}
	// The catalog's information_schema is the server's own metadata: JDBC lists
	// it among the catalog's schemas, its views as tables.
	kept := schemas[:0]
	for _, schema := range schemas {
		if schema != informationSchema {
			kept = append(kept, schema)
		}
	}
	return kept, nil
}

func (d *Driver) fetchAllTablesForCatalog(ctx context.Context, schemas []string) (map[string][]*metadatapb.TableMetadata, error) {
	rows, err := d.db.QueryContext(ctx,
		"SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_type = 'TABLE' ORDER BY table_schem, table_name",
		d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables for catalog: %w", err)
	}
	defer rows.Close()

	allTables := make(map[string][]*metadatapb.TableMetadata, len(schemas))
	for _, schema := range schemas {
		allTables[schema] = []*metadatapb.TableMetadata{}
	}
	for rows.Next() {
		var schemaName, tableName string
		if err := rows.Scan(&schemaName, &tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table row: %w", err)
		}
		allTables[schemaName] = append(allTables[schemaName], &metadatapb.TableMetadata{Name: tableName})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during table row iteration: %w", err)
	}
	return allTables, nil
}

// fetchAllViewsForCatalog reads every view's name and defining query from the
// catalog's information_schema.views. The catalog is quoted because it is an
// identifier in the query text, not a parameter.
func (d *Driver) fetchAllViewsForCatalog(ctx context.Context, schemas []string) (map[string][]*metadatapb.ViewMetadata, error) {
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf(
		"SELECT table_schema, table_name, view_definition FROM %s.information_schema.views ORDER BY table_schema, table_name",
		quoteIdentifier(d.databaseName)))
	if err != nil {
		return nil, fmt.Errorf("failed to query views for catalog: %w", err)
	}
	defer rows.Close()

	allViews := make(map[string][]*metadatapb.ViewMetadata, len(schemas))
	for _, schema := range schemas {
		allViews[schema] = []*metadatapb.ViewMetadata{}
	}
	for rows.Next() {
		var schemaName, viewName, definition string
		if err := rows.Scan(&schemaName, &viewName, &definition); err != nil {
			return nil, fmt.Errorf("failed to scan view row: %w", err)
		}
		allViews[schemaName] = append(allViews[schemaName], &metadatapb.ViewMetadata{
			Name:       viewName,
			Definition: definition,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during view row iteration: %w", err)
	}
	return allViews, nil
}

// quoteIdentifier writes name as a double-quoted Trino identifier.
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func (d *Driver) fetchAllColumnsForCatalog(ctx context.Context) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	rows, err := d.db.QueryContext(ctx,
		"SELECT table_schem, table_name, column_name, type_name, is_nullable FROM system.jdbc.columns WHERE table_cat = ? ORDER BY table_schem, table_name, ordinal_position",
		d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns for catalog: %w", err)
	}
	defer rows.Close()

	allColumns := make(map[tableKey][]*metadatapb.ColumnMetadata)
	for rows.Next() {
		var schemaName, tableName, columnName, dataType, isNullable string
		if err := rows.Scan(&schemaName, &tableName, &columnName, &dataType, &isNullable); err != nil {
			return nil, fmt.Errorf("failed to scan column row: %w", err)
		}
		nullable, err := convertYesNo(isNullable)
		if err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		allColumns[key] = append(allColumns[key], &metadatapb.ColumnMetadata{
			Name:     columnName,
			Type:     dataType,
			Nullable: nullable,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during column row iteration: %w", err)
	}
	return allColumns, nil
}

// convertYesNo converts the JDBC metadata's YES/NO to a bool.
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
