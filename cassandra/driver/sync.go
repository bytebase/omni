package driver

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// versionStatement is the query that answers a cluster's release version, used
// by Ping and by SyncInstance.
const versionStatement = "SELECT release_version FROM system.local"

// getVersion returns the release version the cluster reports.
func (d *Driver) getVersion(ctx context.Context) (string, error) {
	c, err := d.conn()
	if err != nil {
		return "", err
	}
	columns, rows, _, err := c.Query(ctx, versionStatement, 0, nil)
	if err != nil {
		return "", fmt.Errorf("failed to query: %w", err)
	}
	if len(columns) == 0 || len(rows) == 0 || len(rows[0].Values) == 0 {
		return "", fmt.Errorf("failed to query: %s answered no version", versionStatement)
	}
	version, ok := rows[0].Values[0].(*string)
	if !ok || version == nil {
		return "", fmt.Errorf("failed to query: %s answered a version that is not a string", versionStatement)
	}
	return *version, nil
}

// https://docs.datastax.com/en/cql/hcd/reference/system-virtual-tables/system-keyspace-tables.html
func isSystemDatabase(database string) bool {
	switch database {
	case
		"HiveMetaStore",
		"dse_analytics",
		"dse_insights",
		"dse_insights_local",
		"dse_leases",
		"dse_perf",
		"dse_security",
		"dse_system",
		"dse_system_local",
		"dsefs",
		"solr_admin",
		"system_traces",
		"system_auth",
		"system_distributed",
		"system_schema",
		"system":
		return true
	default:
		return false
	}
}

// SyncInstance reads the server version and its user keyspaces, without their
// tables. Cassandra has no roles this driver reads, so Instance.Roles is empty.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	databases, err := d.getDatabases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}

	var filteredDatabases []*metadatapb.DatabaseSchemaMetadata
	for _, database := range databases {
		if isSystemDatabase(database.Name) {
			continue
		}
		filteredDatabases = append(filteredDatabases, database)
	}

	return &driver.Instance{
		Version:   version,
		Databases: filteredDatabases,
	}, nil
}

type primaryKey struct {
	name            string
	clusteringOrder string
	kind            string
	position        int
}

type columnOrderInfo struct {
	column   *metadatapb.ColumnMetadata
	kind     string // "partition_key", "clustering", or "regular"
	position int    // Position within partition/clustering keys
}

// SyncDatabase reads the full schema of Config.Database: its tables with their
// columns in the order Cassandra reads them, the primary key and any secondary
// indexes. A keyspace the server does not have answers an empty schema, as
// Bytebase's did: Cassandra's catalog queries report no rows rather than an
// error, and this driver does not ask whether the keyspace exists first.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	c, err := d.conn()
	if err != nil {
		return nil, err
	}
	schemaMetadata := &metadatapb.SchemaMetadata{
		Name: "",
	}

	tablePKMap := map[string][]primaryKey{}
	columnOrderMap := map[string][]columnOrderInfo{}
	_, columnRows, _, err := c.Query(ctx, fmt.Sprintf(`
		SELECT
			table_name,
			column_name,
			kind,
			position,
			clustering_order,
			type
		FROM system_schema.columns
		WHERE keyspace_name = %s
	`, sqlLiteral(d.config.Database)), 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to scan column: %w", err)
	}
	for _, row := range columnRows {
		if len(row.Values) != 6 {
			return nil, fmt.Errorf("failed to scan column: want 6 values, got %d", len(row.Values))
		}
		tableName := stringValue(row.Values[0])
		columnName := stringValue(row.Values[1])
		kind := stringValue(row.Values[2])
		position := intValue(row.Values[3])
		clusteringOrder := stringValue(row.Values[4])
		columnType := stringValue(row.Values[5])
		columnOrderMap[tableName] = append(columnOrderMap[tableName], columnOrderInfo{
			column: &metadatapb.ColumnMetadata{
				Name:     columnName,
				Type:     columnType,
				Nullable: true,
			},
			kind:     kind,
			position: position,
		})
		if kind != "regular" {
			tablePKMap[tableName] = append(tablePKMap[tableName], primaryKey{
				name:            columnName,
				kind:            kind,
				clusteringOrder: clusteringOrder,
				position:        position,
			})
		}
	}

	// Sort columns to match Cassandra's SELECT * order:
	// 1. Partition keys (by position)
	// 2. Clustering keys (by position)
	// 3. Regular columns (alphabetically)
	// This ensures masking is applied to the correct columns.
	for tableName, columns := range columnOrderMap {
		slices.SortFunc(columns, func(a, b columnOrderInfo) int {
			// Primary key columns come before regular columns
			aPrimary := (a.kind == "partition_key" || a.kind == "clustering")
			bPrimary := (b.kind == "partition_key" || b.kind == "clustering")

			if aPrimary && !bPrimary {
				return -1
			}
			if !aPrimary && bPrimary {
				return 1
			}

			// Within primary key columns, sort by type then position
			if aPrimary && bPrimary {
				// Partition keys before clustering keys
				if a.kind == "partition_key" && b.kind == "clustering" {
					return -1
				}
				if a.kind == "clustering" && b.kind == "partition_key" {
					return 1
				}
				// Within same type, sort by position
				if a.position < b.position {
					return -1
				}
				if a.position > b.position {
					return 1
				}
			}

			// Regular columns are sorted alphabetically
			if !aPrimary && !bPrimary {
				return strings.Compare(a.column.Name, b.column.Name)
			}

			return 0
		})
		columnOrderMap[tableName] = columns
	}

	indexMap := map[string][]*metadatapb.IndexMetadata{}
	_, indexRows, _, err := c.Query(ctx, fmt.Sprintf(`
		SELECT
			table_name,
			index_name,
			kind,
			toJson(options)
		FROM system_schema.indexes
		WHERE keyspace_name = %s
		ORDER BY table_name, index_name
	`, sqlLiteral(d.config.Database)), 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to scan index: %w", err)
	}
	for _, row := range indexRows {
		if len(row.Values) != 4 {
			return nil, fmt.Errorf("failed to scan index: want 4 values, got %d", len(row.Values))
		}
		tableName := stringValue(row.Values[0])
		indexName := stringValue(row.Values[1])
		kind := stringValue(row.Values[2])
		options := stringValue(row.Values[3])
		indexMap[tableName] = append(indexMap[tableName], &metadatapb.IndexMetadata{
			Name:        indexName,
			Type:        kind,
			Expressions: []string{options},
			Definition:  options,
		})
	}

	_, tableRows, _, err := c.Query(ctx, fmt.Sprintf(`
		SELECT
			table_name,
			comment
		FROM system_schema.tables
		WHERE keyspace_name = %s
		ORDER BY table_name
	`, sqlLiteral(d.config.Database)), 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}
	for _, row := range tableRows {
		if len(row.Values) != 2 {
			return nil, fmt.Errorf("failed to scan table: want 2 values, got %d", len(row.Values))
		}
		tableName := stringValue(row.Values[0])
		comment := stringValue(row.Values[1])

		pks := tablePKMap[tableName]
		slices.SortFunc(pks, func(a, b primaryKey) int {
			if a.kind == "partition_key" && b.kind == "clustering" {
				return -1
			}
			if a.kind == "clustering" && b.kind == "partition_key" {
				return 1
			}
			if a.position < b.position {
				return -1
			}
			return 1
		})
		pk := &metadatapb.IndexMetadata{
			Name:        "PRIMARY KEY",
			Expressions: getPKExpressions(pks),
			Descending:  getPKDescending(pks),
			Primary:     true,
			Definition:  getPKDefinition(pks),
		}

		// Extract the sorted column metadata
		var sortedColumns []*metadatapb.ColumnMetadata
		for _, colInfo := range columnOrderMap[tableName] {
			sortedColumns = append(sortedColumns, colInfo.column)
		}

		table := &metadatapb.TableMetadata{
			Name:    tableName,
			Comment: comment,
			Columns: sortedColumns,
		}
		table.Indexes = append(table.Indexes, pk)
		table.Indexes = append(table.Indexes, indexMap[tableName]...)

		schemaMetadata.Tables = append(schemaMetadata.Tables, table)
	}

	return &metadatapb.DatabaseSchemaMetadata{
		Name:    d.config.Database,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}, nil
}

// getDatabases reads the cluster's keyspace names, system ones included;
// SyncInstance filters them.
func (d *Driver) getDatabases(ctx context.Context) ([]*metadatapb.DatabaseSchemaMetadata, error) {
	c, err := d.conn()
	if err != nil {
		return nil, err
	}
	columns, rows, _, err := c.Query(ctx, "SELECT keyspace_name FROM system_schema.keyspaces", 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}
	if len(columns) != 1 {
		return nil, fmt.Errorf("failed to scan: want 1 column, got %d", len(columns))
	}
	databases := make([]*metadatapb.DatabaseSchemaMetadata, 0, len(rows))
	for _, row := range rows {
		databases = append(databases, &metadatapb.DatabaseSchemaMetadata{Name: stringValue(row.Values[0])})
	}
	return databases, nil
}

// sqlLiteral writes a value as a CQL string literal, so the sync queries name
// the keyspace in their own text rather than in a bound parameter.
func sqlLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

// stringValue reads a text cell the way the sync queries hand it back.
func stringValue(v any) string {
	if s, ok := v.(*string); ok && s != nil {
		return *s
	}
	return fmt.Sprintf("%v", v)
}

// intValue reads an integer cell the way the sync queries hand it back.
func intValue(v any) int {
	switch v := v.(type) {
	case *int:
		if v != nil {
			return *v
		}
	case *int32:
		if v != nil {
			return int(*v)
		}
	case *int64:
		if v != nil {
			return int(*v)
		}
	}
	return 0
}

func getPKExpressions(pks []primaryKey) []string {
	var partition, clustering []string
	for _, k := range pks {
		switch k.kind {
		case "partition_key":
			partition = append(partition, k.name)
		case "clustering":
			clustering = append(clustering, k.name)
		default:
			// Ignore other key types
		}
	}

	res := []string{strings.Join(partition, ",")}
	if len(partition) > 1 {
		res[0] = "(" + res[0] + ")"
	}
	res = append(res, clustering...)
	return res
}

func getPKDescending(pks []primaryKey) []bool {
	res := []bool{false}
	for _, pk := range pks {
		if pk.kind == "partition_key" {
			continue
		}
		res = append(res, pk.clusteringOrder == "desc")
	}
	return res
}

func getPKDefinition(pks []primaryKey) string {
	var partition, clustering []string
	for _, k := range pks {
		switch k.kind {
		case "partition_key":
			partition = append(partition, k.name)
		case "clustering":
			clustering = append(clustering, k.name)
		default:
			// Ignore other key types
		}
	}
	if len(partition) == 1 {
		return "(" + strings.Join(append(partition, clustering...), ",") + ")"
	}
	return "(" +
		strings.Join(append([]string{"(" + strings.Join(partition, ",") + ")"}, clustering...), ",") +
		")"
}
