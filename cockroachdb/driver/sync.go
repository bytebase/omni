package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
	metadatapb "github.com/bytebase/omni/metadata"

	crrawanalysis "github.com/bytebase/omni/cockroachdb/analysis"
	"github.com/bytebase/omni/driver"
)

type tableKey struct {
	Schema string
	Table  string
}

// SyncInstance reads the server version, its databases other than the
// system ones, and its roles.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}

	instanceRoles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}

	databases, err := d.getDatabases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}

	var filteredDatabases []*metadatapb.DatabaseSchemaMetadata
	for _, database := range databases {
		if crrawanalysis.IsSystemDatabase(database.Name) {
			continue
		}
		filteredDatabases = append(filteredDatabases, database)
	}

	return &driver.Instance{
		Version:   version,
		Databases: filteredDatabases,
		Roles:     instanceRoles,
	}, nil
}

// SyncDatabase reads the schema of the database the driver is connected to.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	databases, err := d.getDatabases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}

	var databaseMetadata *metadatapb.DatabaseSchemaMetadata
	for _, database := range databases {
		if database.Name == d.databaseName {
			databaseMetadata = database
			break
		}
	}
	if databaseMetadata == nil {
		return nil, fmt.Errorf("database %q: %w", d.databaseName, driver.ErrNotFound)
	}
	if err := d.db.QueryRowContext(ctx, "SHOW search_path").Scan(&databaseMetadata.SearchPath); err != nil {
		return nil, fmt.Errorf("failed to get search path for database %q: %w", d.databaseName, err)
	}
	isAtLeastPG10 := isAtLeastPG10(d.engineVersion)

	txn, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	schemas, err := getSchemas(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas from database %q: %w", d.databaseName, err)
	}
	columnMap, err := getTableColumns(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns from database %q: %w", d.databaseName, err)
	}
	indexMap, err := getIndexes(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes from database %q: %w", d.databaseName, err)
	}
	checksMap, err := getChecks(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get check constraints from database %q: %w", d.databaseName, err)
	}
	tableMap, externalTableMap, err := getTables(txn, isAtLeastPG10, columnMap, indexMap, checksMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables from database %q: %w", d.databaseName, err)
	}
	var tablePartitionMap map[tableKey][]*metadatapb.TablePartitionMetadata
	if isAtLeastPG10 {
		tablePartitionMap, err = getTablePartitions(txn)
		if err != nil {
			return nil, fmt.Errorf("failed to get table partitions from database %q: %w", d.databaseName, err)
		}
	}
	viewMap, err := getViews(txn, columnMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get views from database %q: %w", d.databaseName, err)
	}
	materializedViewMap, err := getMaterializedViews(txn, indexMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get materialized views from database %q: %w", d.databaseName, err)
	}
	functionMap, err := getFunctions(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get functions from database %q: %w", d.databaseName, err)
	}
	sequenceMap, err := getSequences(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequences from database %q: %w", d.databaseName, err)
	}
	enumMap, err := getEnumTypes(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get enum types from database %q: %w", d.databaseName, err)
	}
	compositeMap, err := getCompositeTypes(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get composite types from database %q: %w", d.databaseName, err)
	}

	extensions, err := getExtensions(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get extensions from database %q: %w", d.databaseName, err)
	}

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	for _, schemaName := range schemas {
		tables := tableMap[schemaName]
		for _, table := range tables {
			if isAtLeastPG10 {
				table.Partitions = warpTablePartitions(tablePartitionMap, schemaName, table.Name)
			}
		}
		databaseMetadata.Schemas = append(databaseMetadata.Schemas, &metadatapb.SchemaMetadata{
			Name:              schemaName,
			Tables:            tables,
			ExternalTables:    externalTableMap[schemaName],
			Views:             viewMap[schemaName],
			Functions:         functionMap[schemaName],
			Sequences:         sequenceMap[schemaName],
			MaterializedViews: materializedViewMap[schemaName],
			EnumTypes:         enumMap[schemaName],
			CompositeTypes:    compositeMap[schemaName],
		})
	}
	databaseMetadata.Extensions = extensions

	return databaseMetadata, err
}

func warpTablePartitions(m map[tableKey][]*metadatapb.TablePartitionMetadata, schemaName, tableName string) []*metadatapb.TablePartitionMetadata {
	key := tableKey{Schema: schemaName, Table: tableName}
	if partitions, exists := m[key]; exists {
		defer delete(m, key)
		for _, partition := range partitions {
			partition.Subpartitions = warpTablePartitions(m, schemaName, partition.Name)
		}
		return partitions
	}
	return []*metadatapb.TablePartitionMetadata{}
}

var listForeignKeyQuery = `
SELECT
	n.nspname AS fk_schema,
	fk_cls.relname AS fk_table,
	conname AS fk_name,
	ref_nsp.nspname AS fk_ref_schema,
	ref_cls.relname AS fk_ref_table,
	confdeltype AS delete_option,
	confupdtype AS update_option,
	confmatchtype AS match_option,
	pg_get_constraintdef(c.oid) AS fk_def
FROM
	pg_constraint c
	JOIN pg_namespace n ON n.oid = c.connamespace
	JOIN pg_class fk_cls ON fk_cls.oid = c.conrelid
	JOIN pg_namespace fk_nsp ON fk_nsp.oid = fk_cls.relnamespace
	JOIN pg_class ref_cls ON ref_cls.oid = c.confrelid
	JOIN pg_namespace ref_nsp ON ref_nsp.oid = ref_cls.relnamespace` + fmt.Sprintf(`
WHERE
	n.nspname NOT IN(%s)
	AND c.contype = 'f'
ORDER BY fk_schema, fk_table, fk_name;`, crrawanalysis.SystemSchemaWhereClause)

func getForeignKeys(txn *sql.Tx) (map[tableKey][]*metadatapb.ForeignKeyMetadata, error) {
	foreignKeysMap := make(map[tableKey][]*metadatapb.ForeignKeyMetadata)
	rows, err := txn.Query(listForeignKeyQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fkMetadata metadatapb.ForeignKeyMetadata
		var fkSchema, fkTable, fkDefinition string
		if err := rows.Scan(
			&fkSchema,
			&fkTable,
			&fkMetadata.Name,
			&fkMetadata.ReferencedSchema,
			&fkMetadata.ReferencedTable,
			&fkMetadata.OnDelete,
			&fkMetadata.OnUpdate,
			&fkMetadata.MatchType,
			&fkDefinition,
		); err != nil {
			return nil, err
		}

		fkMetadata.OnDelete = convertForeignKeyActionCode(fkMetadata.OnDelete)
		fkMetadata.OnUpdate = convertForeignKeyActionCode(fkMetadata.OnUpdate)
		fkMetadata.MatchType = convertForeignKeyMatchType(fkMetadata.MatchType)

		if fkMetadata.Columns, fkMetadata.ReferencedColumns, err = getForeignKeyColumnsAndReferencedColumns(fkDefinition); err != nil {
			return nil, err
		}
		key := tableKey{Schema: fkSchema, Table: fkTable}
		foreignKeysMap[key] = append(foreignKeysMap[key], &fkMetadata)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return foreignKeysMap, nil
}

func convertForeignKeyMatchType(in string) string {
	switch in {
	case "f":
		return "FULL"
	case "p":
		return "PARTIAL"
	case "s":
		return "SIMPLE"
	default:
		return in
	}
}

func convertForeignKeyActionCode(in string) string {
	switch in {
	case "a":
		return "NO ACTION"
	case "r":
		return "RESTRICT"
	case "c":
		return "CASCADE"
	case "n":
		return "SET NULL"
	case "d":
		return "SET DEFAULT"
	default:
		return in
	}
}

func getForeignKeyColumnsAndReferencedColumns(definition string) ([]string, []string, error) {
	keyIndex := keywordOutsideQuotes(definition, "FOREIGN KEY", 0)
	if keyIndex < 0 {
		return nil, nil, fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	columns, next, ok := parenthesized(definition, keyIndex+len("FOREIGN KEY"))
	if !ok {
		return nil, nil, fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	referencesIndex := keywordOutsideQuotes(definition, "REFERENCES", next)
	if referencesIndex < 0 {
		return nil, nil, fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	referencedColumns, _, ok := parenthesized(definition, referencesIndex+len("REFERENCES"))
	if !ok {
		return nil, nil, fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	columnList, err := getColumnList(columns)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid foreign key definition %q: %w", definition, err)
	}
	referencedColumnList, err := getColumnList(referencedColumns)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid foreign key definition %q: %w", definition, err)
	}
	return columnList, referencedColumnList, nil
}

// keywordOutsideQuotes finds keyword in s at or after start, skipping quoted
// identifiers so a name such as "references" is not mistaken for the keyword.
func keywordOutsideQuotes(s, keyword string, start int) int {
	for i := start; i+len(keyword) <= len(s); i++ {
		if s[i] == '"' {
			i++
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						i++
					} else {
						break
					}
				}
				i++
			}
			continue
		}
		if strings.EqualFold(s[i:i+len(keyword)], keyword) {
			return i
		}
	}
	return -1
}

// parenthesized returns the text inside the first parenthesis at or after
// start, honoring quoted identifiers so a parenthesis inside a name does not
// close the list, and the offset just past the closing parenthesis.
func parenthesized(s string, start int) (string, int, bool) {
	i := start
	for i < len(s) && s[i] != '(' {
		if s[i] == '"' {
			i++
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						i += 2
						continue
					}
					break
				}
				i++
			}
		}
		i++
	}
	if i >= len(s) {
		return "", start, false
	}
	begin := i + 1
	depth := 1
	i = begin
	for i < len(s) {
		switch s[i] {
		case '"':
			i++
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						i += 2
						continue
					}
					break
				}
				i++
			}
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return s[begin:i], i + 1, true
			}
		}
		i++
	}
	return "", start, false
}

func getColumnList(definition string) ([]string, error) {
	var result []string
	for i := 0; i < len(definition); {
		for i < len(definition) && (definition[i] == ',' || definition[i] == ' ' || definition[i] == '\t' || definition[i] == '\n' || definition[i] == '\r') {
			i++
		}
		if i >= len(definition) {
			break
		}
		if definition[i] == '"' {
			i++
			var name strings.Builder
			for i < len(definition) {
				if definition[i] != '"' {
					name.WriteByte(definition[i])
					i++
					continue
				}
				if i+1 < len(definition) && definition[i+1] == '"' {
					name.WriteByte('"')
					i += 2
					continue
				}
				i++
				break
			}
			result = append(result, name.String())
			continue
		}
		start := i
		for i < len(definition) && definition[i] != ',' {
			i++
		}
		result = append(result, strings.TrimSpace(definition[start:i]))
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("invalid column list definition: %q", definition)
	}
	return result, nil
}

var listSchemaQuery = fmt.Sprintf(`
SELECT nspname
FROM pg_catalog.pg_namespace
WHERE nspname NOT IN (%s)
ORDER BY nspname;
`, crrawanalysis.SystemSchemaWhereClause)

func getSchemas(txn *sql.Tx) ([]string, error) {
	rows, err := txn.Query(listSchemaQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemaNames []string
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		if crrawanalysis.IsSystemSchema(schemaName) {
			continue
		}
		schemaNames = append(schemaNames, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schemaNames, nil
}

func getListForeignTableQuery() string {
	return `SELECT
		foreign_table.foreign_table_schema,
		foreign_table.foreign_table_name,
		foreign_table.foreign_server_catalog,
		foreign_table.foreign_server_name
	FROM information_schema.foreign_tables AS foreign_table
	ORDER BY foreign_table.foreign_table_schema, foreign_table.foreign_table_name;`
}

func getListTableQuery(isAtLeastPG10 bool) string {
	relisPartition := ""
	if isAtLeastPG10 {
		relisPartition = " AND (pc.relispartition IS NULL OR pc.relispartition IS FALSE)"
	}
	// CockroachDB does not support pg_table_size and pg_indexes_size now.
	return `
	SELECT tbl.schemaname, tbl.tablename,
		GREATEST(pc.reltuples::bigint, 0::BIGINT) AS estimate,
		obj_description(pc.oid, 'pg_class') AS comment
	FROM pg_catalog.pg_tables tbl
	LEFT JOIN pg_namespace n ON n.nspname = tbl.schemaname
	LEFT JOIN pg_class pc ON pc.relname = tbl.tablename AND pc.relnamespace = n.oid AND pc.relkind IN ('r', 'p')` + fmt.Sprintf(`
	WHERE tbl.schemaname NOT IN (%s)%s
	ORDER BY tbl.schemaname, tbl.tablename;`, crrawanalysis.SystemSchemaWhereClause, relisPartition)
}

// getTables gets all tables of a database.
func getTables(txn *sql.Tx, isAtLeastPG10 bool, columnMap map[tableKey][]*metadatapb.ColumnMetadata, indexMap map[tableKey][]*metadatapb.IndexMetadata, checksMap map[tableKey][]*metadatapb.CheckConstraintMetadata) (map[string][]*metadatapb.TableMetadata, map[string][]*metadatapb.ExternalTableMetadata, error) {
	foreignKeysMap, err := getForeignKeys(txn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get foreign keys: %w", err)
	}
	foreignTablesMap, err := getForeignTables(txn, columnMap)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get foreign tables: %w", err)
	}

	tableMap := make(map[string][]*metadatapb.TableMetadata)
	query := getListTableQuery(isAtLeastPG10)
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		table := &metadatapb.TableMetadata{}
		var schemaName string
		var comment sql.NullString
		if err := rows.Scan(&schemaName, &table.Name, &table.RowCount, &comment); err != nil {
			return nil, nil, err
		}
		if comment.Valid {
			table.Comment = comment.String
		}
		key := tableKey{Schema: schemaName, Table: table.Name}
		table.Columns = columnMap[key]
		table.Indexes = indexMap[key]
		table.ForeignKeys = foreignKeysMap[key]
		table.CheckConstraints = checksMap[key]

		tableMap[schemaName] = append(tableMap[schemaName], table)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return tableMap, foreignTablesMap, nil
}

func getForeignTables(txn *sql.Tx, columnMap map[tableKey][]*metadatapb.ColumnMetadata) (map[string][]*metadatapb.ExternalTableMetadata, error) {
	query := getListForeignTableQuery()
	rows, err := txn.Query(query)
	if err != nil {
		// Experimental feature, log error and return.
		slog.Error("failed to query foreign table", "error", err)
		return nil, nil
	}
	defer rows.Close()

	foreignTablesMap := make(map[string][]*metadatapb.ExternalTableMetadata)

	for rows.Next() {
		var schemaName, tableName, foreignServerCatalog, foreignServerName string
		if err := rows.Scan(&schemaName, &tableName, &foreignServerCatalog, &foreignServerName); err != nil {
			slog.Error("failed to scan foreign table", "error", err)
			return nil, nil
		}
		externalTable := &metadatapb.ExternalTableMetadata{
			Name:                 tableName,
			ExternalServerName:   foreignServerName,
			ExternalDatabaseName: foreignServerCatalog,
		}
		key := tableKey{Schema: schemaName, Table: externalTable.Name}
		externalTable.Columns = columnMap[key]

		foreignTablesMap[schemaName] = append(foreignTablesMap[schemaName], externalTable)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to scan foreign table", "error", err)
		return nil, nil
	}

	return foreignTablesMap, nil
}

var listTablePartitionQuery = `
SELECT
	n.nspname AS schema_name,
	c.relname AS table_name,
	i2.nspname AS inh_schema_name,
	i2.relname AS inh_table_name,
	i2.partstrat AS partition_type,
	pg_get_expr(c.relpartbound, c.oid) AS rel_part_bound
FROM
	pg_catalog.pg_class c
	LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
    LEFT JOIN (
		pg_inherits i 
		INNER JOIN pg_class c2 ON i.inhparent = c2.oid 
		LEFT JOIN pg_namespace n2 ON n2.oid = c2.relnamespace
		LEFT JOIN pg_partitioned_table p ON p.partrelid = c2.oid
	) i2 ON i2.inhrelid = c.oid 
WHERE
	((c.relkind = 'r'::"char") OR (c.relkind = 'f'::"char") OR (c.relkind = 'p'::"char"))
	AND c.relispartition IS TRUE ` + fmt.Sprintf(`
	AND n.nspname NOT IN (%s)
ORDER BY c.oid;`, crrawanalysis.SystemSchemaWhereClause)

func getTablePartitions(txn *sql.Tx) (map[tableKey][]*metadatapb.TablePartitionMetadata, error) {
	result := make(map[tableKey][]*metadatapb.TablePartitionMetadata)
	rows, err := txn.Query(listTablePartitionQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schemaName, tableName, inhSchemaName, inhTableName, partitionType, relPartBound string
		if err := rows.Scan(&schemaName, &tableName, &inhSchemaName, &inhTableName, &partitionType, &relPartBound); err != nil {
			return nil, err
		}
		key := tableKey{Schema: inhSchemaName, Table: inhTableName}
		metadata := &metadatapb.TablePartitionMetadata{
			Name:       tableName,
			Expression: relPartBound,
		}
		switch strings.ToLower(partitionType) {
		case "l":
			metadata.Type = metadatapb.TablePartitionMetadata_LIST
		case "r":
			metadata.Type = metadatapb.TablePartitionMetadata_RANGE
		case "h":
			metadata.Type = metadatapb.TablePartitionMetadata_HASH
		default:
			return nil, fmt.Errorf("invalid partition type %q", partitionType)
		}
		result[key] = append(result[key], metadata)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

var listColumnQuery = `
SELECT
	cols.table_schema,
	cols.table_name,
	cols.column_name,
	cols.data_type,
	cols.character_maximum_length,
	cols.numeric_precision,
	cols.numeric_scale,
	cols.ordinal_position,
	cols.column_default,
	cols.is_nullable,
	cols.collation_name,
	cols.udt_schema,
	cols.udt_name,
	pg_catalog.col_description(c.oid, cols.ordinal_position::int) as column_comment,
	cols.is_hidden,
	cols.generation_expression,
	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
FROM INFORMATION_SCHEMA.COLUMNS AS cols
LEFT JOIN pg_namespace n ON n.nspname = cols.table_schema
LEFT JOIN pg_class c ON c.relname = cols.table_name AND c.relnamespace = n.oid AND c.relkind IN ('r', 'v', 'm', 'p', 'f')` + fmt.Sprintf(`
WHERE cols.table_schema NOT IN (%s)
ORDER BY cols.table_schema, cols.table_name, cols.ordinal_position;`, crrawanalysis.SystemSchemaWhereClause)

// getTableColumns gets the columns of a table.
func getTableColumns(txn *sql.Tx) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	columnsMap := make(map[tableKey][]*metadatapb.ColumnMetadata)
	rows, err := txn.Query(listColumnQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var schemaName, tableName, nullable string
		var characterMaxLength, defaultStr, collation, udtSchema, udtName, comment sql.NullString
		var numericPrecision, numericScale sql.NullInt64
		var hidden string
		var generation, generated sql.NullString
		if err := rows.Scan(&schemaName, &tableName, &column.Name, &column.Type, &characterMaxLength, &numericPrecision, &numericScale, &column.Position, &defaultStr, &nullable, &collation, &udtSchema, &udtName, &comment, &hidden, &generation, &generated); err != nil {
			return nil, err
		}
		// A computed column states its expression apart from its default;
		// pg_attribute.attgenerated tells a stored one ('s') from a virtual
		// one ('v').
		if generation.String != "" {
			column.Generation = &metadatapb.GenerationMetadata{Type: metadatapb.GenerationMetadata_TYPE_STORED, Expression: generation.String}
			if generated.String == "v" {
				column.Generation.Type = metadatapb.GenerationMetadata_TYPE_VIRTUAL
			}
		}
		// A hidden column, such as the rowid CockroachDB adds to a table with
		// no primary key, is one * leaves out.
		column.IsInvisible = hidden == "YES"
		if defaultStr.Valid {
			// Store in Default field (migration from DefaultExpression to Default)
			column.Default = defaultStr.String
		}
		isNullBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		column.Nullable = isNullBool
		switch column.Type {
		case "USER-DEFINED":
			column.Type = fmt.Sprintf("%s.%s", udtSchema.String, udtName.String)
		case "ARRAY":
			column.Type = udtName.String
		case "character", "character varying", "bit", "bit varying":
			if characterMaxLength.Valid {
				// For character varying(n), the character maximum length is n.
				// For character without length specifier, key character_maximum_length is null,
				// we don't need to append the length.
				// https://www.postgresql.org/docs/current/infoschema-columns.html.
				column.Type = fmt.Sprintf("%s(%s)", column.Type, characterMaxLength.String)
			}
		case "numeric", "decimal":
			if numericPrecision.Valid {
				if numericScale.Valid && numericScale.Int64 != 0 {
					column.Type = fmt.Sprintf("%s(%d,%d)", column.Type, numericPrecision.Int64, numericScale.Int64)
				} else {
					column.Type = fmt.Sprintf("%s(%d)", column.Type, numericPrecision.Int64)
				}
			}
		default:
			// Keep the type as is
		}
		column.Collation = collation.String
		column.Comment = comment.String

		key := tableKey{Schema: schemaName, Table: tableName}
		columnsMap[key] = append(columnsMap[key], column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columnsMap, nil
}

var listMaterializedViewQuery = `
SELECT m.schemaname, m.matviewname, m.definition, obj_description(c.oid, 'pg_class')
FROM pg_catalog.pg_matviews m
LEFT JOIN pg_namespace n ON n.nspname = m.schemaname
LEFT JOIN pg_class c ON c.relname = m.matviewname AND c.relnamespace = n.oid AND c.relkind = 'm'` + fmt.Sprintf(`
WHERE m.schemaname NOT IN (%s)
ORDER BY m.schemaname, m.matviewname;`, crrawanalysis.SystemSchemaWhereClause)

func getMaterializedViews(txn *sql.Tx, indexMap map[tableKey][]*metadatapb.IndexMetadata) (map[string][]*metadatapb.MaterializedViewMetadata, error) {
	matviewMap := make(map[string][]*metadatapb.MaterializedViewMetadata)

	rows, err := txn.Query(listMaterializedViewQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		matview := &metadatapb.MaterializedViewMetadata{}
		var schemaName string
		var def, comment sql.NullString
		if err := rows.Scan(&schemaName, &matview.Name, &def, &comment); err != nil {
			return nil, err
		}
		// Return error on NULL view definition.
		if !def.Valid {
			return nil, fmt.Errorf("schema %q materialized view %q has empty definition; please check whether proper privileges have been granted to Bytebase", schemaName, matview.Name)
		}
		matview.Definition = def.String
		if comment.Valid {
			matview.Comment = comment.String
		}

		matview.Indexes = indexMap[tableKey{Schema: schemaName, Table: matview.Name}]
		matviewMap[schemaName] = append(matviewMap[schemaName], matview)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for schemaName, list := range matviewMap {
		for _, matview := range list {
			dependencies, err := getViewDependencies(txn, schemaName, matview.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get materialized view %q dependencies: %w", matview.Name, err)
			}
			matview.DependencyColumns = dependencies
		}
	}

	return matviewMap, nil
}

var listViewQuery = `
SELECT v.schemaname, v.viewname, v.definition, obj_description(c.oid, 'pg_class')
FROM pg_catalog.pg_views v
LEFT JOIN pg_namespace n ON n.nspname = v.schemaname
LEFT JOIN pg_class c ON c.relname = v.viewname AND c.relnamespace = n.oid AND c.relkind = 'v'` + fmt.Sprintf(`
WHERE v.schemaname NOT IN (%s)
ORDER BY v.schemaname, v.viewname;`, crrawanalysis.SystemSchemaWhereClause)

// getViews gets all views of a database.
func getViews(txn *sql.Tx, columnMap map[tableKey][]*metadatapb.ColumnMetadata) (map[string][]*metadatapb.ViewMetadata, error) {
	viewMap := make(map[string][]*metadatapb.ViewMetadata)

	rows, err := txn.Query(listViewQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		view := &metadatapb.ViewMetadata{}
		var schemaName string
		var def, comment sql.NullString
		if err := rows.Scan(&schemaName, &view.Name, &def, &comment); err != nil {
			return nil, err
		}
		// Return error on NULL view definition.
		// https://github.com/bytebase/bytebase/issues/343
		if !def.Valid {
			return nil, fmt.Errorf("schema %q view %q has empty definition; please check whether proper privileges have been granted to Bytebase", schemaName, view.Name)
		}
		view.Definition = def.String
		if comment.Valid {
			view.Comment = comment.String
		}

		key := tableKey{Schema: schemaName, Table: view.Name}
		view.Columns = columnMap[key]

		viewMap[schemaName] = append(viewMap[schemaName], view)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for schemaName, list := range viewMap {
		for _, view := range list {
			dependencies, err := getViewDependencies(txn, schemaName, view.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get view %q dependencies: %w", view.Name, err)
			}
			view.DependencyColumns = dependencies
		}
	}

	return viewMap, nil
}

// getViewDependencies gets the dependencies of a view.
func getViewDependencies(txn *sql.Tx, schemaName, viewName string) ([]*metadatapb.DependencyColumn, error) {
	var result []*metadatapb.DependencyColumn

	query := `
		SELECT source_ns.nspname as source_schema,
	  		source_table.relname as source_table,
	  		pg_attribute.attname as column_name
	  	FROM pg_depend 
	  		JOIN pg_rewrite ON pg_depend.objid = pg_rewrite.oid 
	  		JOIN pg_class as dependency_view ON pg_rewrite.ev_class = dependency_view.oid 
	  		JOIN pg_class as source_table ON pg_depend.refobjid = source_table.oid 
	  		JOIN pg_attribute ON pg_depend.refobjid = pg_attribute.attrelid 
	  		    AND pg_depend.refobjsubid = pg_attribute.attnum 
	  		JOIN pg_namespace dependency_ns ON dependency_ns.oid = dependency_view.relnamespace
	  		JOIN pg_namespace source_ns ON source_ns.oid = source_table.relnamespace
	  	WHERE 
	  		dependency_ns.nspname = $1
	  		AND dependency_view.relname = $2
	  		AND pg_attribute.attnum > 0
	  		-- Only consider SELECT rules (view definitions), not INSERT/UPDATE/DELETE rules
	  		AND pg_rewrite.ev_type = '1'
	  	ORDER BY 1,2,3;
	`

	rows, err := txn.Query(query, schemaName, viewName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		dependencyColumn := &metadatapb.DependencyColumn{}
		if err := rows.Scan(&dependencyColumn.Schema, &dependencyColumn.Table, &dependencyColumn.Column); err != nil {
			return nil, err
		}
		result = append(result, dependencyColumn)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// getExtensions gets all extensions of a database.
func getExtensions(txn *sql.Tx) ([]*metadatapb.ExtensionMetadata, error) {
	var extensions []*metadatapb.ExtensionMetadata

	query := `
		SELECT e.extname, e.extversion, n.nspname, c.description
		FROM pg_catalog.pg_extension e
		LEFT JOIN pg_catalog.pg_namespace n ON n.oid = e.extnamespace
		LEFT JOIN pg_catalog.pg_description c ON c.objoid = e.oid AND c.classoid = 'pg_catalog.pg_extension'::pg_catalog.regclass
		WHERE n.nspname != 'pg_catalog'
		ORDER BY e.extname;`
	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		e := &metadatapb.ExtensionMetadata{}
		var description sql.NullString
		if err := rows.Scan(&e.Name, &e.Version, &e.Schema, &description); err != nil {
			return nil, err
		}
		if description.Valid {
			e.Description = description.String
		}
		extensions = append(extensions, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return extensions, nil
}

var listCheckConstraintQuery = `
SELECT nsp.nspname, rel.relname, con.conname, pg_get_constraintdef(con.oid), obj_description(con.oid, 'pg_constraint')
FROM pg_catalog.pg_constraint con
JOIN pg_catalog.pg_class rel ON rel.oid = con.conrelid
JOIN pg_catalog.pg_namespace nsp ON nsp.oid = con.connamespace` + fmt.Sprintf(`
WHERE con.contype = 'c'
	AND nsp.nspname NOT IN (%s)
ORDER BY nsp.nspname, rel.relname, con.conname;`, crrawanalysis.SystemSchemaWhereClause)

func getChecks(txn *sql.Tx) (map[tableKey][]*metadatapb.CheckConstraintMetadata, error) {
	checksMap := make(map[tableKey][]*metadatapb.CheckConstraintMetadata)
	rows, err := txn.Query(listCheckConstraintQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		check := &metadatapb.CheckConstraintMetadata{}
		var schemaName, tableName, definition string
		var comment sql.NullString
		if err := rows.Scan(&schemaName, &tableName, &check.Name, &definition, &comment); err != nil {
			return nil, err
		}
		check.Expression = strings.TrimPrefix(definition, "CHECK ")
		if comment.Valid {
			check.Comment = comment.String
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		checksMap[key] = append(checksMap[key], check)
	}
	return checksMap, rows.Err()
}

var listEnumTypeQuery = `
SELECT pn.nspname, pt.typname, pe.enumlabel
FROM pg_enum pe
JOIN pg_type pt ON pe.enumtypid = pt.oid
JOIN pg_namespace pn ON pt.typnamespace = pn.oid` + fmt.Sprintf(`
WHERE pn.nspname NOT IN (%s)
ORDER BY pn.nspname, pt.typname, pe.enumsortorder;`, crrawanalysis.SystemSchemaWhereClause)

func getEnumTypes(txn *sql.Tx) (map[string][]*metadatapb.EnumTypeMetadata, error) {
	enumMap := make(map[string][]*metadatapb.EnumTypeMetadata)
	rows, err := txn.Query(listEnumTypeQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var currentSchema, currentName string
	var values []string
	flush := func() {
		if currentSchema == "" {
			return
		}
		enumMap[currentSchema] = append(enumMap[currentSchema], &metadatapb.EnumTypeMetadata{Name: currentName, Values: values})
	}
	for rows.Next() {
		var schemaName, typeName, value string
		if err := rows.Scan(&schemaName, &typeName, &value); err != nil {
			return nil, err
		}
		if currentSchema != schemaName || currentName != typeName {
			flush()
			currentSchema, currentName, values = schemaName, typeName, nil
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	flush()
	return enumMap, nil
}

var listCompositeTypeQuery = `
SELECT pn.nspname, pt.typname, a.attname, pg_catalog.format_type(a.atttypid, a.atttypmod)
FROM pg_type pt
JOIN pg_class pc ON pc.oid = pt.typrelid AND pc.relkind = 'c'
JOIN pg_namespace pn ON pn.oid = pt.typnamespace
JOIN pg_attribute a ON a.attrelid = pc.oid AND a.attnum > 0 AND NOT a.attisdropped` + fmt.Sprintf(`
WHERE pn.nspname NOT IN (%s)
ORDER BY pn.nspname, pt.typname, a.attnum;`, crrawanalysis.SystemSchemaWhereClause)

func getCompositeTypes(txn *sql.Tx) (map[string][]*metadatapb.CompositeTypeMetadata, error) {
	compositeMap := make(map[string][]*metadatapb.CompositeTypeMetadata)
	rows, err := txn.Query(listCompositeTypeQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var currentSchema, currentName string
	var current *metadatapb.CompositeTypeMetadata
	for rows.Next() {
		var schemaName, typeName, attributeName, attributeType string
		if err := rows.Scan(&schemaName, &typeName, &attributeName, &attributeType); err != nil {
			return nil, err
		}
		if current == nil || currentSchema != schemaName || currentName != typeName {
			current = &metadatapb.CompositeTypeMetadata{Name: typeName}
			compositeMap[schemaName] = append(compositeMap[schemaName], current)
			currentSchema, currentName = schemaName, typeName
		}
		current.Attributes = append(current.Attributes, &metadatapb.CompositeTypeAttribute{Name: attributeName, Type: attributeType})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return compositeMap, nil
}

// getSequences gets all sequences of a database.
func getSequences(txn *sql.Tx) (map[string][]*metadatapb.SequenceMetadata, error) {
	query := `
		SELECT sequence_schema, sequence_name, data_type, start_value, minimum_value, maximum_value, increment, cycle_option
		FROM information_schema.sequences
		ORDER BY sequence_schema, sequence_name;`
	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sequenceMap := make(map[string][]*metadatapb.SequenceMetadata)
	for rows.Next() {
		sequence := &metadatapb.SequenceMetadata{}
		var schemaName, cycleOption string
		if err := rows.Scan(&schemaName, &sequence.Name, &sequence.DataType, &sequence.Start, &sequence.MinValue, &sequence.MaxValue, &sequence.Increment, &cycleOption); err != nil {
			return nil, err
		}
		sequence.Cycle = strings.EqualFold(cycleOption, "YES")
		sequenceMap[schemaName] = append(sequenceMap[schemaName], sequence)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sequenceMap, nil
}

var listIndexQuery = `
SELECT idx.schemaname, idx.tablename, idx.indexname, idx.indexdef, (SELECT constraint_type
	FROM information_schema.table_constraints
	WHERE constraint_schema = idx.schemaname
	AND constraint_name = idx.indexname
	AND table_schema = idx.schemaname
	AND table_name = idx.tablename
	AND constraint_type = 'PRIMARY KEY') AS constraint_type,
	obj_description(c.oid, 'pg_class') AS comment
FROM pg_indexes AS idx
LEFT JOIN pg_namespace n ON n.nspname = idx.schemaname
LEFT JOIN pg_class c ON c.relname = idx.indexname AND c.relnamespace = n.oid AND c.relkind = 'i'` + fmt.Sprintf(`
WHERE idx.schemaname NOT IN (%s)
ORDER BY idx.schemaname, idx.tablename, idx.indexname;`, crrawanalysis.SystemSchemaWhereClause)

// getIndexes gets all indices of a database.
func getIndexes(txn *sql.Tx) (map[tableKey][]*metadatapb.IndexMetadata, error) {
	indexMap := make(map[tableKey][]*metadatapb.IndexMetadata)

	rows, err := txn.Query(listIndexQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		index := &metadatapb.IndexMetadata{}
		var schemaName, tableName, statement string
		var constraintType sql.NullString
		var comment sql.NullString
		if err := rows.Scan(&schemaName, &tableName, &index.Name, &statement, &constraintType, &comment); err != nil {
			return nil, err
		}

		indexDef, err := crrawanalysis.ParseIndexDefinition(statement)
		if err != nil {
			return nil, err
		}
		index.Expressions = indexDef.Expressions
		index.Descending = indexDef.Descending
		index.Definition = statement
		index.Unique = indexDef.Unique

		// FIXME(zp): Get expression.
		index.Type = getIndexMethodType(statement)
		if constraintType.Valid {
			switch constraintType.String {
			case "PRIMARY KEY":
				index.Primary = true
				index.Unique = true
			case "UNIQUE":
				index.Unique = true
			default:
				// Other constraint types
			}
		}
		if comment.Valid {
			index.Comment = comment.String
		}

		key := tableKey{Schema: schemaName, Table: tableName}
		indexMap[key] = append(indexMap[key], index)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return indexMap, nil
}

func getIndexMethodType(stmt string) string {
	re := regexp.MustCompile(`USING (\w+) `)
	matches := re.FindStringSubmatch(stmt)
	if len(matches) == 0 {
		return ""
	}
	return matches[1]
}

// listFunctionQuery reads each routine's CREATE statement. CockroachDB's
// pg_get_functiondef returns only the body of a procedure, so a procedure is
// read from crdb_internal.create_procedure_statements, whose procedure_id is
// its pg_proc oid less 100000.
var listFunctionQuery = `
select n.nspname as function_schema,
	p.proname as function_name,
	p.proname || '(' || pg_get_function_identity_arguments(p.oid) || ')' as function_signature,
	case when l.lanname = 'internal' then p.prosrc
			when p.prokind = 'p' then (
				select s.create_statement from crdb_internal.create_procedure_statements as s
				where s.procedure_id + 100000 = p.oid::INT8)
			else pg_get_functiondef(p.oid)
			end as definition
from pg_proc p
left join pg_namespace n on p.pronamespace = n.oid
left join pg_language l on p.prolang = l.oid
left join pg_type t on t.oid = p.prorettype ` + fmt.Sprintf(`
where n.nspname not in (%s)
order by function_schema, function_name;`, crrawanalysis.SystemSchemaWhereClause)

// getFunctions gets all functions of a database.
func getFunctions(txn *sql.Tx) (map[string][]*metadatapb.FunctionMetadata, error) {
	functionMap := make(map[string][]*metadatapb.FunctionMetadata)

	rows, err := txn.Query(listFunctionQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		function := &metadatapb.FunctionMetadata{}
		var schemaName string
		if err := rows.Scan(&schemaName, &function.Name, &function.Signature, &function.Definition); err != nil {
			return nil, err
		}
		functionMap[schemaName] = append(functionMap[schemaName], function)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return functionMap, nil
}

func isAtLeastPG10(version string) bool {
	v, err := semver.ParseTolerant(version)
	if err != nil {
		slog.Error("invalid postgres version", slog.String("version", version))
		// Assume the version is at least 10.0 for any error.
		return true
	}
	return v.Major >= 10
}
