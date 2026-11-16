package driver

import (
	"context"
	"errors"
	"slices"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"

	"github.com/bytebase/omni/driver"
)

// SyncInstance syncs the instance: DynamoDB's pseudo database, named
// {account_id}-{region} because DynamoDB has no databases of its own and
// Bytebase names one per account and region. DynamoDB reports no version and no
// roles.
func (d *Driver) SyncInstance(context.Context) (*driver.Instance, error) {
	if d.accountID == "" {
		return nil, errors.New("account id is empty in the driver config")
	}
	if d.region == "" {
		return nil, errors.New("region is empty in the AWS config")
	}
	return &driver.Instance{
		Databases: []*metadatapb.DatabaseSchemaMetadata{
			{
				Name: formatDatabaseName(d.accountID, d.region),
			},
		},
	}, nil
}

// SyncDatabase syncs the database schema: one schema holding one table per
// table of the region, with the table's key schema as its columns and indexes.
// DynamoDB has no databases, so Config.Database names nothing the sync reads
// and a sync never reports driver.ErrNotFound.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	schemaMetadata := &metadatapb.SchemaMetadata{}

	tableNames, err := d.listAllTables(ctx)
	if err != nil {
		return nil, err
	}

	// DynamoDB does not support batch describe table, so the driver describes
	// one table at a time. Those round trips are why a region with many tables
	// syncs slowly.
	for _, tableName := range tableNames {
		tableMetadata, err := d.syncTable(ctx, tableName)
		if err != nil {
			return nil, err
		}
		if tableMetadata != nil {
			schemaMetadata.Tables = append(schemaMetadata.Tables, tableMetadata)
		}
	}

	return &metadatapb.DatabaseSchemaMetadata{
		Name: d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{
			schemaMetadata,
		},
	}, nil
}

func (d *Driver) syncTable(ctx context.Context, tableName string) (*metadatapb.TableMetadata, error) {
	tableMetadata := &metadatapb.TableMetadata{
		Name: tableName,
	}
	table, err := d.client.DescribeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}
	// A table that is not fully available yet, or no longer is, is left out
	// rather than synced half-made.
	if table.Status != tableStatusActive {
		return nil, nil
	}
	if len(table.KeySchema) > 0 {
		var hashKeyAttributes, rangeKeyAttributes []string
		for _, keySchema := range table.KeySchema {
			switch keySchema.KeyType {
			case keyTypeHash:
				hashKeyAttributes = append(hashKeyAttributes, keySchema.AttributeName)
			case keyTypeRange:
				rangeKeyAttributes = append(rangeKeyAttributes, keySchema.AttributeName)
			default:
				// Unknown key type, ignore
			}
		}
		tableMetadata.Indexes = append(tableMetadata.Indexes, &metadatapb.IndexMetadata{
			Name:        strings.Join(hashKeyAttributes, ","),
			Expressions: append([]string{}, hashKeyAttributes...),
			Type:        keyTypeHash,
		})
		tableMetadata.Indexes = append(tableMetadata.Indexes, &metadatapb.IndexMetadata{
			Name:        strings.Join(rangeKeyAttributes, ","),
			Expressions: append([]string{}, rangeKeyAttributes...),
			Type:        keyTypeRange,
		})
		columnsMap := make(map[string]bool)
		for _, key := range hashKeyAttributes {
			columnsMap[key] = true
		}
		for _, key := range rangeKeyAttributes {
			columnsMap[key] = true
		}
		sortedColumns := make([]string, 0, len(columnsMap))
		for key := range columnsMap {
			sortedColumns = append(sortedColumns, key)
		}
		slices.Sort(sortedColumns)
		tableMetadata.Columns = make([]*metadatapb.ColumnMetadata, 0, len(sortedColumns))
		for _, key := range sortedColumns {
			tableMetadata.Columns = append(tableMetadata.Columns, &metadatapb.ColumnMetadata{
				Name: key,
			})
		}
	}
	tableMetadata.Name = table.Name
	// DynamoDB reports the item count as an estimate it refreshes periodically,
	// so a table that was just written may still report 0.
	tableMetadata.RowCount = table.ItemCount

	return tableMetadata, nil
}

func (d *Driver) listAllTables(ctx context.Context) ([]string, error) {
	var result []string
	var exclusiveStartTableName string
	for {
		tableNames, lastEvaluatedTableName, err := d.client.ListTables(ctx, exclusiveStartTableName, 100)
		if err != nil {
			return nil, err
		}
		result = append(result, tableNames...)
		if lastEvaluatedTableName == "" {
			break
		}
		exclusiveStartTableName = lastEvaluatedTableName
	}
	return result, nil
}

func formatDatabaseName(accountID string, region string) string {
	return accountID + "-" + region
}
