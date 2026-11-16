package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"

	"github.com/bytebase/omni/driver"
)

// SyncInstance reads the server version and its databases. Hive has no role a
// sync reads, so Instance.Roles stays empty.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	var instanceMetadata driver.Instance

	// version.
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}

	// databases.
	databaseNames, err := d.getDatabaseNames(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range databaseNames {
		instanceMetadata.Databases = append(instanceMetadata.Databases, &metadatapb.DatabaseSchemaMetadata{
			Name: v,
		})
	}

	instanceMetadata.Version = version
	return &instanceMetadata, nil
}

// SyncDatabase reads the full schema of d.databaseName, or of the server's
// default database when none was named.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	dbName := d.databaseName
	if dbName == "" {
		dbName = "default"
	}

	schemaMetadata, err := d.getDatabaseInfoByName(ctx, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to get info from database %s: %w", dbName, err)
	}
	return &metadatapb.DatabaseSchemaMetadata{
		Name:    dbName,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}, nil
}

func (d *Driver) getVersion(ctx context.Context) (string, error) {
	result, err := d.queryStatement(ctx, "SELECT VERSION()", 0)
	if err != nil {
		return "", fmt.Errorf("failed to get version from instance: %w", err)
	}

	if len(result.Rows) == 0 || len(result.Rows[0]) == 0 {
		return "", errors.New("invalid version result")
	}
	// rawVersion has the format of "1.2.3 commitID".
	rawVersion := result.Rows[0][0].String
	tokens := strings.Split(rawVersion, " ")
	if len(tokens) == 0 {
		return "", fmt.Errorf("invalid version %q", rawVersion)
	}
	return tokens[0], nil
}

func (d *Driver) getDatabaseNames(ctx context.Context) ([]string, error) {
	var databaseNames []string
	result, err := d.queryStatement(ctx, "SHOW DATABASES", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get database names from instance: %w", err)
	}
	for _, row := range result.Rows {
		if len(row) == 0 {
			return nil, errors.New("row values have zero length")
		}
		databaseNames = append(databaseNames, row[0].String)
	}
	return databaseNames, nil
}

func (d *Driver) listTablesNames(ctx context.Context, databaseName string) ([]string, error) {
	result, err := d.queryStatement(ctx, fmt.Sprintf("SHOW TABLES FROM %s", quoteBacktick(databaseName)), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables from database %s: %w", databaseName, err)
	}

	var tableNames []string
	for _, row := range result.Rows {
		if len(row) == 0 {
			return nil, errors.New("row values have zero length")
		}
		tableNames = append(tableNames, row[0].String)
	}
	return tableNames, nil
}

// getTables fetches table info and returns structed table data.
func (d *Driver) getTables(ctx context.Context, databaseName string) (
	[]*metadatapb.TableMetadata,
	[]*metadatapb.ExternalTableMetadata,
	[]*metadatapb.ViewMetadata,
	[]*metadatapb.MaterializedViewMetadata,
	error,
) {
	var (
		tableMetadatas    []*metadatapb.TableMetadata
		extTableMetadatas []*metadatapb.ExternalTableMetadata
		viewMetadatas     []*metadatapb.ViewMetadata
		mtViewMetadatas   []*metadatapb.MaterializedViewMetadata
	)

	// list tables' names.
	tableNames, err := d.listTablesNames(ctx, databaseName)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to list tables: %w", err)
	}

	// iterations in tables of certain database.
	for _, tableName := range tableNames {
		// filter out index table names.
		if strings.HasSuffix(tableName, "__") {
			continue
		}

		tableInfo, err := d.getTableInfo(ctx, tableName, databaseName)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to describe table %s's type: %w", tableName, err)
		}

		// different processing way according to the type of the table.
		switch tableInfo.tableType {
		case "MATERIALIZED_VIEW":
			mtViewMetadatas = append(mtViewMetadatas, &metadatapb.MaterializedViewMetadata{
				Name:       tableName,
				Definition: tableInfo.viewDef,
				Comment:    tableInfo.comment,
			})
		case "VIRTUAL_VIEW":
			viewMetadatas = append(viewMetadatas, &metadatapb.ViewMetadata{
				Name:       tableName,
				Definition: tableInfo.viewDef,
				Comment:    tableInfo.comment,
			})
		case "EXTERNAL_TABLE":
			extTableMetadatas = append(extTableMetadatas, &metadatapb.ExternalTableMetadata{
				Name:    tableName,
				Columns: tableInfo.colMetadatas,
			})
		case "MANAGED_TABLE":
			// A table without partitions answers SHOW PARTITIONS with an error,
			// which is why a table carries no partitions rather than no rows.
			partitions := d.getPartitions(ctx, databaseName, tableName)
			tableMetadatas = append(tableMetadatas, &metadatapb.TableMetadata{
				Engine:     "HDFS",
				Comment:    tableInfo.comment,
				Columns:    tableInfo.colMetadatas,
				DataSize:   int64(tableInfo.totalSize),
				RowCount:   int64(tableInfo.numRows),
				Partitions: partitions,
				Name:       tableName,
			})
		default:
			// ignore other types of tables.
		}
	}

	return tableMetadatas, extTableMetadatas, viewMetadatas, mtViewMetadatas, nil
}

// getPartitions reads a managed table's partitions. Hive answers a table that
// is not partitioned with an error, which is the same as having none.
func (d *Driver) getPartitions(ctx context.Context, databaseName, tableName string) []*metadatapb.TablePartitionMetadata {
	result, err := d.queryStatement(ctx, fmt.Sprintf("SHOW PARTITIONS `%s`.`%s`", databaseName, tableName), 0)
	if err != nil {
		slog.Debug("failed to get partitions", "error", err)
		return nil
	}
	partitions := make([]*metadatapb.TablePartitionMetadata, 0, len(result.Rows))
	for _, row := range result.Rows {
		if len(row) == 0 {
			continue
		}
		partitions = append(partitions, &metadatapb.TablePartitionMetadata{
			Name: row[0].String,
		})
	}
	return partitions
}

// This function gets certain database info by name.
func (d *Driver) getDatabaseInfoByName(ctx context.Context, databaseName string) (*metadatapb.SchemaMetadata, error) {
	// fetch table metadata.
	tableMetadata, extTabMetadata, viewMetadata, mtViewMetadata, err := d.getTables(ctx, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table metadata from database %s: %w", databaseName, err)
	}

	return &metadatapb.SchemaMetadata{
		Name:              databaseName,
		Tables:            tableMetadata,
		ExternalTables:    extTabMetadata,
		Views:             viewMetadata,
		MaterializedViews: mtViewMetadata,
	}, nil
}

// TableInfo is what one DESCRIBE FORMATTED answer says about a table.
type TableInfo struct {
	tableType    string
	colMetadatas []*metadatapb.ColumnMetadata
	numRows      int
	viewDef      string
	totalSize    int
	comment      string
}

const (
	columnSection = iota
	tableInfoSection
	storageSection
	viewSection
)

func (d *Driver) getTableInfo(ctx context.Context, tableName string, databaseName string) (
	*TableInfo,
	error,
) {
	tableInfo := &TableInfo{}

	query := fmt.Sprintf("DESCRIBE FORMATTED `%s`.`%s`", databaseName, tableName)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to describe table %s: %w", tableName, err)
	}
	defer rows.Close()

	section := columnSection
	for rows.Next() {
		// DESCRIBE FORMATTED answers NULL for the cells a section's headings and
		// blank separator rows do not fill, so each cell is read as a value that
		// may be absent.
		var (
			colName, dataType, comment sql.NullString
		)
		if err := rows.Scan(&colName, &dataType, &comment); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		colNameValue, dataTypeValue, commentValue := colName.String, dataType.String, comment.String

		// The first rows are column metadata, followed by "# Detailed Table Information" and "# Storage Information"
		switch {
		case strings.HasPrefix(colNameValue, "# col_name"):
			continue
		case strings.HasPrefix(colNameValue, "# Detailed Table Information"):
			section = tableInfoSection
		case strings.HasPrefix(colNameValue, "# Storage Information"):
			section = storageSection
		case strings.HasPrefix(colNameValue, "# View Information"):
			section = viewSection
		default:
			// No action needed for other cases
		}
		switch section {
		case columnSection:
			if colNameValue != "" && dataTypeValue != "" {
				// Column metadata.
				position := len(tableInfo.colMetadatas) + 1
				tableInfo.colMetadatas = append(tableInfo.colMetadatas, &metadatapb.ColumnMetadata{
					Name:     colNameValue,
					Type:     dataTypeValue,
					Comment:  commentValue,
					Position: int32(position),
				})
			}
		case tableInfoSection:
			switch {
			case trimField(colNameValue) == "Table Type:":
				tableInfo.tableType = trimField(dataTypeValue)
			case trimField(dataTypeValue) == "numRows":
				n, err := strconv.Atoi(trimField(commentValue))
				if err != nil {
					return nil, fmt.Errorf("failed to parse row count: %w", err)
				}
				tableInfo.numRows = n
			case trimField(dataTypeValue) == "totalSize":
				n, err := strconv.Atoi(trimField(commentValue))
				if err != nil {
					return nil, fmt.Errorf("failed to parse row count: %w", err)
				}
				tableInfo.totalSize = n
			case trimField(dataTypeValue) == "comment":
				tableInfo.comment = trimField(commentValue)
			default:
				// No action needed for other table info fields
			}
		case viewSection:
			switch {
			case strings.HasPrefix(colNameValue, "View Original Text:"):
				tableInfo.viewDef = trimField(dataTypeValue)
			case strings.HasPrefix(colNameValue, "Original Query:"):
				tableInfo.viewDef = trimField(dataTypeValue)
			default:
				// No action needed for other view fields
			}
		default:
			// No action needed for other sections
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tableInfo, nil
}

func trimField(field string) string {
	return strings.Trim(field, " ")
}
