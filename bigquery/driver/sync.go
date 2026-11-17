package driver

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// SyncDatabase reads the schema of the dataset the driver is open on: every
// table's columns from the dataset's INFORMATION_SCHEMA, and each table's comment,
// row count and data size from its metadata.
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
	columns, err := d.columns(ctx)
	if err != nil {
		return nil, notFoundError(err)
	}
	tables, err := d.tables(ctx, columns)
	if err != nil {
		return nil, notFoundError(err)
	}
	return &metadatapb.DatabaseSchemaMetadata{
		Name:    d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{{Tables: tables}},
	}, nil
}

// columnRow is one row of INFORMATION_SCHEMA.COLUMNS as the client decodes it.
type columnRow struct {
	TableName       string              `bigquery:"table_name"`
	ColumnName      string              `bigquery:"column_name"`
	OrdinalPosition int32               `bigquery:"ordinal_position"`
	IsNullable      string              `bigquery:"is_nullable"`
	DataType        string              `bigquery:"data_type"`
	CollationName   bigquery.NullString `bigquery:"collation_name"`
	ColumnDefault   bigquery.NullString `bigquery:"column_default"`
}

// columns reads every column of every table in the dataset, grouped by table
// name and ordered by table and position.
func (d *Driver) columns(ctx context.Context) (map[string][]*metadatapb.ColumnMetadata, error) {
	q := d.query(fmt.Sprintf(`
		SELECT
			table_name,
			column_name,
			ordinal_position,
			is_nullable,
			data_type,
			collation_name,
			column_default
		FROM %s.INFORMATION_SCHEMA.COLUMNS ORDER BY table_name, ordinal_position;`, d.databaseName))
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}
	columns := make(map[string][]*metadatapb.ColumnMetadata)
	for {
		var row columnRow
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		column, err := columnMetadata(row)
		if err != nil {
			return nil, err
		}
		columns[row.TableName] = append(columns[row.TableName], column)
	}
	return columns, nil
}

// columnMetadata maps one INFORMATION_SCHEMA row onto a synced column.
//
// A column BigQuery states no default for reports the text "NULL", which is what
// Bytebase's driver reported and what a caller reads as no default.
func columnMetadata(row columnRow) (*metadatapb.ColumnMetadata, error) {
	nullable, err := convertYesNo(row.IsNullable)
	if err != nil {
		return nil, err
	}
	column := &metadatapb.ColumnMetadata{
		Name:     row.ColumnName,
		Position: row.OrdinalPosition,
		Nullable: nullable,
		Type:     row.DataType,
		Default:  "NULL",
	}
	if row.CollationName.Valid {
		column.Collation = row.CollationName.String()
	}
	if row.ColumnDefault.Valid {
		column.Default = row.ColumnDefault.String()
	}
	return column, nil
}

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

// tables reads every table of the dataset with its columns, comment, row count
// and data size, ordered by name.
func (d *Driver) tables(ctx context.Context, columns map[string][]*metadatapb.ColumnMetadata) ([]*metadatapb.TableMetadata, error) {
	var tables []*metadatapb.TableMetadata
	it := d.client.Dataset(d.databaseName).Tables(ctx)
	for {
		table, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		metadata, err := table.Metadata(ctx)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &metadatapb.TableMetadata{
			Name:     table.TableID,
			Columns:  columns[table.TableID],
			RowCount: int64(metadata.NumRows),
			Comment:  metadata.Description,
			DataSize: metadata.NumBytes,
		})
	}
	sort.Slice(tables, func(i, j int) bool { return tables[i].Name < tables[j].Name })
	return tables, nil
}
