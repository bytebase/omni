package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"

	"cloud.google.com/go/bigquery"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// SyncDatabase reads the schema of the dataset the driver is open on: every
// table's and view's columns, and each view's definition, from the dataset's
// INFORMATION_SCHEMA, and each table's comment, row count and data size from its
// metadata.
//
// The dataset is one unnamed schema in the synced metadata, which is what
// Bytebase's driver reported: a BigQuery dataset has no schema layer below it.
// Tables are reported in name order rather than in the order the service lists
// them, because the service's order is not part of an answer and a caller
// compares two syncs.
//
// A dataset that is not there is driver.ErrNotFound, which is how a caller tells
// a database it may drop from a request that failed.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	if d.databaseName == "" {
		return nil, fmt.Errorf("the driver is open on no dataset: %w", driver.ErrNotFound)
	}
	// The dataset is asked for before its columns, because the service answers a
	// read of a missing dataset's INFORMATION_SCHEMA with accessDenied: it states a
	// table's absence only to a caller that may read it. The dataset's own
	// metadata does say not found, so a caller that may list the project learns
	// which of the two happened, and one that may not gets accessDenied for both.
	if _, err := d.client.DatasetMetadata(ctx, d.databaseName); err != nil {
		return nil, notFoundError(err)
	}
	columns, err := d.columns(ctx)
	if err != nil {
		return nil, notFoundError(err)
	}
	views, err := d.viewDefinitions(ctx)
	if err != nil {
		return nil, notFoundError(err)
	}
	schema, err := d.relations(ctx, columns, views)
	if err != nil {
		return nil, notFoundError(err)
	}
	return &metadatapb.DatabaseSchemaMetadata{
		Name:    d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{schema},
	}, nil
}

// columnRow is one row of INFORMATION_SCHEMA.COLUMNS as the client decodes it.
//
// The view states no position for the pseudo-columns of an ingestion-time
// partitioned table, so the position is nullable: the client refuses to decode
// a NULL into a plain integer, and one such table would fail the whole sync.
type columnRow struct {
	TableName       string              `bigquery:"table_name"`
	ColumnName      string              `bigquery:"column_name"`
	OrdinalPosition bigquery.NullInt64  `bigquery:"ordinal_position"`
	IsNullable      string              `bigquery:"is_nullable"`
	DataType        string              `bigquery:"data_type"`
	CollationName   bigquery.NullString `bigquery:"collation_name"`
	ColumnDefault   bigquery.NullString `bigquery:"column_default"`
}

// columns reads every column of every table in the dataset, grouped by table
// name and ordered by table and position.
func (d *Driver) columns(ctx context.Context) (map[string][]*metadatapb.ColumnMetadata, error) {
	// is_hidden = 'NO' drops pseudo-columns such as _PARTITIONTIME /
	// _PARTITIONDATE at the service; those rows also state a NULL
	// ordinal_position, which the loop below skips as a second guard.
	statement := fmt.Sprintf(`
		SELECT
			table_name,
			column_name,
			ordinal_position,
			is_nullable,
			data_type,
			collation_name,
			column_default
		FROM %s.INFORMATION_SCHEMA.COLUMNS
		WHERE is_hidden = 'NO'
		ORDER BY table_name, ordinal_position;`, quoteIdentifier(d.databaseName))
	it, err := d.client.Read(ctx, d.databaseName, statement)
	if err != nil {
		return nil, err
	}
	columns := make(map[string][]*metadatapb.ColumnMetadata)
	for {
		values, err := it.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		row, err := decodeColumnRow(values)
		if err != nil {
			return nil, err
		}
		if !row.OrdinalPosition.Valid {
			// A pseudo-column of an ingestion-time partitioned table
			// (_PARTITIONTIME, _PARTITIONDATE) is not one the table declares.
			continue
		}
		column, err := columnMetadata(row)
		if err != nil {
			return nil, err
		}
		columns[row.TableName] = append(columns[row.TableName], column)
	}
	return columns, nil
}

// decodeColumnRow reads the seven columns selected from INFORMATION_SCHEMA.COLUMNS
// into a columnRow.
func decodeColumnRow(values []bigquery.Value) (columnRow, error) {
	if len(values) < 7 {
		return columnRow{}, fmt.Errorf("INFORMATION_SCHEMA.COLUMNS row holds %d values, want 7", len(values))
	}
	tableName, ok := values[0].(string)
	if !ok {
		return columnRow{}, fmt.Errorf("unexpected table_name value %T", values[0])
	}
	columnName, ok := values[1].(string)
	if !ok {
		return columnRow{}, fmt.Errorf("unexpected column_name value %T", values[1])
	}
	var ordinalPosition bigquery.NullInt64
	if values[2] != nil {
		pos, ok := values[2].(int64)
		if !ok {
			return columnRow{}, fmt.Errorf("unexpected ordinal_position value %T", values[2])
		}
		ordinalPosition = bigquery.NullInt64{Int64: pos, Valid: true}
	}
	isNullable, ok := values[3].(string)
	if !ok {
		return columnRow{}, fmt.Errorf("unexpected is_nullable value %T", values[3])
	}
	dataType, ok := values[4].(string)
	if !ok {
		return columnRow{}, fmt.Errorf("unexpected data_type value %T", values[4])
	}
	var collationName bigquery.NullString
	if values[5] != nil {
		s, ok := values[5].(string)
		if !ok {
			return columnRow{}, fmt.Errorf("unexpected collation_name value %T", values[5])
		}
		collationName = bigquery.NullString{StringVal: s, Valid: true}
	}
	var columnDefault bigquery.NullString
	if values[6] != nil {
		s, ok := values[6].(string)
		if !ok {
			return columnRow{}, fmt.Errorf("unexpected column_default value %T", values[6])
		}
		columnDefault = bigquery.NullString{StringVal: s, Valid: true}
	}
	return columnRow{
		TableName:       tableName,
		ColumnName:      columnName,
		OrdinalPosition: ordinalPosition,
		IsNullable:      isNullable,
		DataType:        dataType,
		CollationName:   collationName,
		ColumnDefault:   columnDefault,
	}, nil
}

// columnMetadata maps one INFORMATION_SCHEMA row onto a synced column.
//
// A column BigQuery states no default or no collation for states the text "NULL"
// in that view, and the client decodes that text as a value rather than an absence:
// it is what a caller reads as none, and what Bytebase's driver reported.
func columnMetadata(row columnRow) (*metadatapb.ColumnMetadata, error) {
	nullable, err := convertYesNo(row.IsNullable)
	if err != nil {
		return nil, err
	}
	column := &metadatapb.ColumnMetadata{
		Name:     row.ColumnName,
		Position: int32(row.OrdinalPosition.Int64),
		Nullable: nullable,
		Type:     row.DataType,
		Default:  noValueText,
	}
	if row.CollationName.Valid && row.CollationName.String() != noValueText {
		column.Collation = row.CollationName.String()
	}
	if row.ColumnDefault.Valid {
		column.Default = row.ColumnDefault.String()
	}
	return column, nil
}

// noValueText is the text BigQuery's INFORMATION_SCHEMA states for a value a column
// does not have.
const noValueText = "NULL"

// convertYesNo converts an INFORMATION_SCHEMA YES or NO to a bool.
func convertYesNo(value string) (bool, error) {
	switch value {
	case "YES":
		return true, nil
	case "NO":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognized is_nullable value %q", value)
	}
}

// viewDefinitions reads the definition of every logical view in the dataset,
// keyed by view name. A materialized view is not listed in
// INFORMATION_SCHEMA.VIEWS, so it stays a table: it stores rows, and a query
// reads them rather than its definition.
func (d *Driver) viewDefinitions(ctx context.Context) (map[string]string, error) {
	statement := fmt.Sprintf(`
		SELECT
			table_name,
			view_definition
		FROM %s.INFORMATION_SCHEMA.VIEWS ORDER BY table_name;`, quoteIdentifier(d.databaseName))
	it, err := d.client.Read(ctx, d.databaseName, statement)
	if err != nil {
		return nil, err
	}
	views := make(map[string]string)
	for {
		values, err := it.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(values) < 2 {
			return nil, fmt.Errorf("INFORMATION_SCHEMA.VIEWS row holds %d values, want 2", len(values))
		}
		name, ok := values[0].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected table_name value %T", values[0])
		}
		definition, ok := values[1].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected view_definition value %T", values[1])
		}
		views[name] = definition
	}
	return views, nil
}

// relations reads every table and view of the dataset, ordered by name: a table
// with its columns, comment, row count and data size, and a view with its
// columns, comment and definition. Table IDs are sorted before reading each
// one's metadata so metadata requests happen in deterministic name order.
func (d *Driver) relations(ctx context.Context, columns map[string][]*metadatapb.ColumnMetadata, views map[string]string) (*metadatapb.SchemaMetadata, error) {
	tableIDs, err := d.client.ListTables(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	sort.Strings(tableIDs)
	schema := &metadatapb.SchemaMetadata{}
	for _, tableID := range tableIDs {
		metadata, err := d.client.TableMetadata(ctx, d.databaseName, tableID)
		if err != nil {
			return nil, err
		}
		if definition, ok := views[tableID]; ok {
			schema.Views = append(schema.Views, &metadatapb.ViewMetadata{
				Name:       tableID,
				Definition: definition,
				Columns:    columns[tableID],
				Comment:    metadata.Description,
			})
			continue
		}
		schema.Tables = append(schema.Tables, &metadatapb.TableMetadata{
			Name:     tableID,
			Columns:  columns[tableID],
			RowCount: int64(metadata.NumRows),
			Comment:  metadata.Description,
			DataSize: metadata.NumBytes,
		})
	}
	return schema, nil
}

// quoteIdentifier quotes name as a BigQuery path expression, which is what a
// dataset a query names is: a dataset may be called a reserved keyword, and the
// path it opens then has to be quoted.
func quoteIdentifier(name string) string {
	return "`" + name + "`"
}
