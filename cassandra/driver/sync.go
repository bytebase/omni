package driver

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

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
	kind     string // "partition_key", "clustering", "static", or "regular"
	position int    // Position within partition/clustering keys
}

// SyncDatabase reads the full schema of Config.Database: its tables with their
// columns in the order Cassandra reads them, the primary key and any secondary
// indexes. A keyspace the server no longer has is driver.ErrNotFound: the
// driver asks the catalog whether the keyspace exists before reading it, so a
// keyspace dropped after the driver opened is reported as gone rather than as
// an empty schema. A driver opened on a keyspace the server never had fails
// earlier, when gocql builds the session.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	c, err := d.conn()
	if err != nil {
		return nil, err
	}
	exists, err := d.keyspaceExists(ctx, c)
	if err != nil {
		if isMissingKeyspace(err) {
			return nil, fmt.Errorf("%w: keyspace %s does not exist", driver.ErrNotFound, d.config.Database)
		}
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("%w: keyspace %s does not exist", driver.ErrNotFound, d.config.Database)
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
		// A static column is not a key column: only the partition and
		// clustering keys take part in the primary key's expressions and
		// descending flags.
		if kind == "partition_key" || kind == "clustering" {
			tablePKMap[tableName] = append(tablePKMap[tableName], primaryKey{
				name:            columnName,
				kind:            kind,
				clusteringOrder: clusteringOrder,
				position:        position,
			})
		}
	}

	// Sort columns to match Cassandra's SELECT * order: the partition key
	// columns by position, the clustering columns by position, then the static
	// columns and last the regular columns, each by name. A masking policy is
	// applied by position, so this order is what keeps it on its column.
	kindRank := func(kind string) int {
		switch kind {
		case "partition_key":
			return 0
		case "clustering":
			return 1
		case "static":
			return 2
		}
		return 3
	}
	for _, columns := range columnOrderMap {
		slices.SortFunc(columns, func(a, b columnOrderInfo) int {
			if ra, rb := kindRank(a.kind), kindRank(b.kind); ra != rb {
				return ra - rb
			}
			if a.kind == "partition_key" || a.kind == "clustering" {
				return a.position - b.position
			}
			return strings.Compare(a.column.Name, b.column.Name)
		})
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

	// A materialized view is listed in system_schema.views, never in
	// system_schema.tables, and its columns are in system_schema.columns under
	// its own name. It can only select columns of its base table, under their
	// own names, so each of its columns depends on the base-table column of
	// the same name: what a query of the view reads.
	_, viewRows, _, err := c.Query(ctx, fmt.Sprintf(`
		SELECT
			view_name,
			base_table_name
		FROM system_schema.views
		WHERE keyspace_name = %s
		ORDER BY view_name
	`, sqlLiteral(d.config.Database)), 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to scan materialized view: %w", err)
	}
	for _, row := range viewRows {
		if len(row.Values) != 2 {
			return nil, fmt.Errorf("failed to scan materialized view: want 2 values, got %d", len(row.Values))
		}
		viewName := stringValue(row.Values[0])
		baseTable := stringValue(row.Values[1])
		view := &metadatapb.MaterializedViewMetadata{Name: viewName}
		for _, colInfo := range columnOrderMap[viewName] {
			view.Columns = append(view.Columns, colInfo.column)
			view.DependencyColumns = append(view.DependencyColumns, &metadatapb.DependencyColumn{
				Table:  baseTable,
				Column: colInfo.column.Name,
			})
		}
		schemaMetadata.MaterializedViews = append(schemaMetadata.MaterializedViews, view)
	}

	return &metadatapb.DatabaseSchemaMetadata{
		Name:    d.config.Database,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}, nil
}

// isMissingKeyspace reports whether the server answered that the session's
// keyspace does not exist. A session outlives its keyspace: gocql refuses a
// query with either no keyspaces or a message naming the one that is gone.
func isMissingKeyspace(err error) bool {
	var requestErr gocql.RequestError
	if !errors.As(err, &requestErr) {
		return false
	}
	message := requestErr.Message()
	return strings.Contains(message, "does not exist") || strings.Contains(message, "Keyspace")
}

// keyspaceExists reports whether the server has the driver's keyspace, so a
// dropped keyspace is reported as gone rather than as an empty schema.
func (d *Driver) keyspaceExists(ctx context.Context, c client) (bool, error) {
	_, rows, _, err := c.Query(ctx, fmt.Sprintf(
		"SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name = %s", sqlLiteral(d.config.Database)), 0, nil)
	if err != nil {
		return false, fmt.Errorf("failed to query keyspace: %w", err)
	}
	return len(rows) > 0, nil
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
