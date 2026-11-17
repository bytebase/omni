package driver

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// tableKey names a table or a view inside one schema, which is how a sync
// gathers the columns, indexes and foreign keys read by separate queries.
type tableKey struct {
	Schema string
	Table  string
}

// SyncInstance reads the server version, its databases other than the ones the
// cloud provider and PostgreSQL keep for themselves, and its roles.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}
	databases, err := d.getDatabases(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %w", err)
	}

	var userDatabases []*metadatapb.DatabaseSchemaMetadata
	for _, database := range databases {
		if excludedDatabaseList[database.Name] {
			continue
		}
		userDatabases = append(userDatabases, database)
	}
	return &driver.Instance{
		Version:   version,
		Databases: userDatabases,
		Roles:     roles,
	}, nil
}

// SyncDatabase reads the schema of the database the driver is connected to: its
// schemas, and inside each the tables with their columns, indexes and foreign
// keys, and the views with their columns. A datashare database is read through
// the SVV_ALL_* views, because it publishes no pg_catalog entries.
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

	txn, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = txn.Rollback() }()

	schemaList, err := d.getSchemas(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas from database %q: %w", d.databaseName, err)
	}
	var tableMap map[string][]*metadatapb.TableMetadata
	var viewMap map[string][]*metadatapb.ViewMetadata
	if d.datashare {
		tableMap, err = d.getDatashareTables(txn)
		if err != nil {
			return nil, fmt.Errorf("failed to get tables from datashare database %q: %w", d.databaseName, err)
		}
	} else {
		columnMap, err := getTableColumns(txn)
		if err != nil {
			return nil, fmt.Errorf("failed to get table columns: %w", err)
		}
		tableMap, err = getTables(txn, columnMap)
		if err != nil {
			return nil, fmt.Errorf("failed to get tables from database %q: %w", d.databaseName, err)
		}
		viewMap, err = getViews(txn, columnMap)
		if err != nil {
			return nil, fmt.Errorf("failed to get views from database %q: %w", d.databaseName, err)
		}
	}
	if err := txn.Commit(); err != nil {
		return nil, err
	}
	for _, schemaName := range schemaList {
		databaseMetadata.Schemas = append(databaseMetadata.Schemas, &metadatapb.SchemaMetadata{
			Name:   schemaName,
			Tables: tableMap[schemaName],
			Views:  viewMap[schemaName],
		})
	}
	return databaseMetadata, nil
}

// getSchemas reads the database's schema names in name order.
func (d *Driver) getSchemas(txn *sql.Tx) ([]string, error) {
	const query = `
		SELECT
			schema_name
		FROM
			SVV_ALL_SCHEMAS
		WHERE
			database_name = $1
		ORDER BY
			schema_name;
	`
	rows, err := txn.Query(query, d.databaseName)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return result, nil
}

// getTables reads the database's tables with their row estimate and comment,
// and gathers the columns, indexes and foreign keys read by other queries.
func getTables(txn *sql.Tx, columnMap map[tableKey][]*metadatapb.ColumnMetadata) (map[string][]*metadatapb.TableMetadata, error) {
	indexMap, err := getIndexes(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes: %w", err)
	}
	foreignKeysMap, err := getForeignKeys(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign keys: %w", err)
	}

	tableMap := make(map[string][]*metadatapb.TableMetadata)
	const query = `
	SELECT
		ptbl.schemaname,
		ptbl.tablename,
		0, -- data size
		0, -- index size
		GREATEST(pc.reltuples::bigint, 0::bigint) AS estimate,
		obj_description(pc.oid) AS comment
	FROM pg_catalog.pg_tables AS ptbl
	JOIN pg_namespace AS pns ON pns.nspname = ptbl.schemaname
	LEFT JOIN pg_class AS pc ON pc.relname = ptbl.tablename AND pns.oid = pc.relnamespace
	WHERE ptbl.schemaname NOT IN ('pg_catalog', 'information_schema');`
	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		table := &metadatapb.TableMetadata{}
		var schemaName string
		var comment sql.NullString
		if err := rows.Scan(&schemaName, &table.Name, &table.DataSize, &table.IndexSize, &table.RowCount, &comment); err != nil {
			return nil, err
		}
		if comment.Valid {
			table.Comment = comment.String
		}
		key := tableKey{Schema: schemaName, Table: table.Name}
		table.Columns = columnMap[key]
		table.Indexes = indexMap[key]
		table.ForeignKeys = foreignKeysMap[key]
		tableMap[schemaName] = append(tableMap[schemaName], table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tableMap, nil
}

// getTableColumns reads every column of every table and view, in schema, table
// and position order.
func getTableColumns(txn *sql.Tx) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	columnsMap := make(map[tableKey][]*metadatapb.ColumnMetadata)

	const query = `
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
		pg_catalog.col_description(pc.oid, cols.ordinal_position::int) as column_comment
	FROM 
		INFORMATION_SCHEMA.COLUMNS AS cols
		JOIN pg_namespace AS pns ON pns.nspname = cols.table_schema
		LEFT JOIN pg_class AS pc ON pc.relname = cols.table_name AND pns.oid = pc.relnamespace
	WHERE cols.table_schema NOT IN ('pg_catalog', 'information_schema')
	ORDER BY cols.table_schema, cols.table_name, cols.ordinal_position;`
	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var schemaName, tableName, nullable string
		var defaultStr, collation, udtSchema, udtName, comment sql.NullString
		var characterMaxLength, numericPrecision, numericScale sql.NullInt64
		if err := rows.Scan(&schemaName, &tableName, &column.Name, &column.Type, &characterMaxLength, &numericPrecision, &numericScale, &column.Position, &defaultStr, &nullable, &collation, &udtSchema, &udtName, &comment); err != nil {
			return nil, err
		}
		if defaultStr.Valid {
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
		default:
			column.Type = columnType(column.Type, characterMaxLength, numericPrecision, numericScale)
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

// getDatashareTables reads the tables inside a datashare database.
func (d *Driver) getDatashareTables(txn *sql.Tx) (map[string][]*metadatapb.TableMetadata, error) {
	columnMap, err := d.getDatashareTableColumns(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}

	tableMap := make(map[string][]*metadatapb.TableMetadata)
	const query = `
	SELECT
		schema_name,
		table_name,
		table_type
	FROM SVV_ALL_TABLES
	WHERE database_name = $1
	ORDER BY schema_name, table_name;`
	rows, err := txn.Query(query, d.databaseName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		table := &metadatapb.TableMetadata{}
		var schemaName, tableType string
		if err := rows.Scan(&schemaName, &table.Name, &tableType); err != nil {
			return nil, err
		}
		if strings.EqualFold(tableType, "VIEW") {
			// A shared view is listed here too, and it is not a table: the sync
			// has no view definition for a datashare, so VERIFY.md says what a
			// recording has to settle before one is reported.
			continue
		}
		key := tableKey{Schema: schemaName, Table: table.Name}
		table.Columns = columnMap[key]
		tableMap[schemaName] = append(tableMap[schemaName], table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tableMap, nil
}

// getDatashareTableColumns reads the columns of the tables inside a datashare
// database.
func (d *Driver) getDatashareTableColumns(txn *sql.Tx) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	columnsMap := make(map[tableKey][]*metadatapb.ColumnMetadata)

	const query = `
	SELECT
		schema_name,
		table_name,
		column_name,
		data_type,
		ordinal_position,
		column_default,
		is_nullable,
		character_maximum_length,
		numeric_precision,
		numeric_scale
	FROM SVV_ALL_COLUMNS
	WHERE database_name = $1
	ORDER BY schema_name, table_name, ordinal_position;`
	rows, err := txn.Query(query, d.databaseName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var schemaName, tableName, nullable string
		var defaultStr sql.NullString
		var varcharMaxLength, numericPrecision, numericScale sql.NullInt32
		if err := rows.Scan(&schemaName, &tableName, &column.Name, &column.Type, &column.Position, &defaultStr, &nullable, &varcharMaxLength, &numericPrecision, &numericScale); err != nil {
			return nil, err
		}
		if defaultStr.Valid {
			column.Default = defaultStr.String
		}
		isNullBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		column.Nullable = isNullBool
		if column.Type == "character varying" && varcharMaxLength.Valid {
			column.Type = fmt.Sprintf("varchar(%d)", varcharMaxLength.Int32)
		}
		if column.Type == "numeric" && numericPrecision.Valid && numericScale.Valid {
			column.Type = fmt.Sprintf("numeric(%d,%d)", numericPrecision.Int32, numericScale.Int32)
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		columnsMap[key] = append(columnsMap[key], column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return columnsMap, nil
}

// getViews reads the database's views with their definitions and comments.
func getViews(txn *sql.Tx, columnMap map[tableKey][]*metadatapb.ColumnMetadata) (map[string][]*metadatapb.ViewMetadata, error) {
	viewMap := make(map[string][]*metadatapb.ViewMetadata)

	const query = `
SELECT
	pgv.schemaname,
	pgv.viewname,
	obj_description(pc.oid) AS comment
FROM pg_catalog.pg_views AS pgv
JOIN pg_namespace AS pns ON pns.nspname = pgv.schemaname
JOIN pg_class AS pc ON pc.relname = pgv.viewname AND pns.oid = pc.relnamespace
WHERE pgv.schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY pgv.schemaname, pgv.viewname;`
	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	type viewRow struct {
		schemaName string
		view       *metadatapb.ViewMetadata
	}
	var views []viewRow
	for rows.Next() {
		view := &metadatapb.ViewMetadata{}
		var schemaName string
		var comment sql.NullString
		if err := rows.Scan(&schemaName, &view.Name, &comment); err != nil {
			rows.Close()
			return nil, err
		}
		if comment.Valid {
			view.Comment = comment.String
		}
		views = append(views, viewRow{schemaName: schemaName, view: view})
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}

	// PG_VIEWS states a definition in a varchar(256), so a long view's text is
	// cut short there; SHOW VIEW states it whole. The views are read in the
	// order the catalog listed them, so a recording holds one order.
	for _, row := range views {
		definition, err := showView(txn, row.schemaName, row.view.Name)
		if err != nil {
			return nil, err
		}
		row.view.Definition = definition
		key := tableKey{Schema: row.schemaName, Table: row.view.Name}
		row.view.Columns = columnMap[key]
		viewMap[row.schemaName] = append(viewMap[row.schemaName], row.view)
	}
	return viewMap, nil
}

// showView reads a view's whole definition, which is what a synced view needs to
// be recreated. PG_VIEWS truncates it, and a view the account cannot read
// states none: the sync refuses rather than record a view it cannot recreate.
func showView(txn *sql.Tx, schemaName, viewName string) (string, error) {
	rows, err := txn.Query(fmt.Sprintf("SHOW VIEW %s.%s;", quoteIdent(schemaName), quoteIdent(viewName)))
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		return "", fmt.Errorf("schema %q view %q has empty definition; please check whether proper privileges have been granted to Bytebase", schemaName, viewName)
	}
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	values := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	if err := rows.Scan(dest...); err != nil {
		return "", err
	}
	for _, value := range values {
		switch text := value.(type) {
		case string:
			if text != "" {
				return text, nil
			}
		case []byte:
			if len(text) > 0 {
				return string(text), nil
			}
		}
	}
	return "", fmt.Errorf("schema %q view %q has empty definition; please check whether proper privileges have been granted to Bytebase", schemaName, viewName)
}

// getIndexes reads the database's indexes with their key expressions. Redshift
// stores a key column's expression in the index definition, and
// generate_series stands in for PostgreSQL's generate_subscripts to read one
// column's definition at a time.
func getIndexes(txn *sql.Tx) (map[tableKey][]*metadatapb.IndexMetadata, error) {
	indexMap := make(map[tableKey][]*metadatapb.IndexMetadata)

	const query = `
	WITH index_positions AS (
		SELECT generate_series(0, 31) as pos
	)
	SELECT 
		pgidx.schemaname,
		pgidx.tablename, 
		pgidx.indexname,
		pgidx.indexdef,
		ix.indisunique,
		ix.indisprimary,
		p.pos + 1 as position,
		pg_get_indexdef(pc.oid, p.pos + 1, true) as column_expression,
		obj_description(pc.oid) AS comment
	FROM
		pg_indexes AS pgidx 
		JOIN pg_namespace AS pns ON pns.nspname = pgidx.schemaname
		JOIN pg_class AS pc ON pc.relname = pgidx.indexname AND pns.oid = pc.relnamespace
		JOIN pg_index AS ix ON ix.indexrelid = pc.oid
		CROSS JOIN index_positions p
	WHERE 
		pgidx.schemaname NOT IN ('pg_catalog', 'information_schema')
		AND ix.indkey[p.pos] > 0  -- Only process positions that have actual columns
	ORDER BY
		pgidx.schemaname, pgidx.tablename, pgidx.indexname, p.pos;`

	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type indexKey struct {
		schema string
		table  string
		index  string
	}
	indexData := make(map[indexKey]*metadatapb.IndexMetadata)
	indexColumns := make(map[indexKey][]string)

	for rows.Next() {
		var schemaName, tableName, indexName, indexDef string
		var isUnique, isPrimary bool
		var position int
		var columnExpression, comment sql.NullString
		if err := rows.Scan(&schemaName, &tableName, &indexName, &indexDef,
			&isUnique, &isPrimary, &position, &columnExpression, &comment); err != nil {
			return nil, err
		}

		key := indexKey{schema: schemaName, table: tableName, index: indexName}
		// One row per key position: the index itself is read from the first.
		if _, exists := indexData[key]; !exists {
			indexData[key] = &metadatapb.IndexMetadata{
				Name:       indexName,
				Definition: indexDef,
				Type:       getIndexMethodType(indexDef),
				Unique:     isUnique,
				Primary:    isPrimary,
			}
			if comment.Valid {
				indexData[key].Comment = comment.String
			}
		}
		if columnExpression.Valid {
			indexColumns[key] = append(indexColumns[key], columnExpression.String)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for key, index := range indexData {
		if columns, ok := indexColumns[key]; ok {
			index.Expressions = columns
		}
		tableKey := tableKey{Schema: key.schema, Table: key.table}
		indexMap[tableKey] = append(indexMap[tableKey], index)
	}
	return indexMap, nil
}

// getIndexMethodType reads the access method an index definition names.
func getIndexMethodType(stmt string) string {
	matches := indexMethodRegexp.FindStringSubmatch(stmt)
	if len(matches) > 1 {
		return matches[1]
	}
	// Redshift's default.
	return "btree"
}

var indexMethodRegexp = regexp.MustCompile(`USING (\w+)`)

// getForeignKeys reads the database's foreign keys, keyed by the table that
// holds each.
func getForeignKeys(txn *sql.Tx) (map[tableKey][]*metadatapb.ForeignKeyMetadata, error) {
	const query = `
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
		JOIN pg_namespace ref_nsp ON ref_nsp.oid = ref_cls.relnamespace
	WHERE
		n.nspname NOT IN('pg_catalog', 'information_schema')
		AND c.contype = 'f'
	ORDER BY fk_schema, fk_table, fk_name;
	`
	foreignKeysMap := make(map[tableKey][]*metadatapb.ForeignKeyMetadata)
	rows, err := txn.Query(query)
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

		var err error
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

// convertForeignKeyMatchType spells out the match option pg_constraint writes
// as one character.
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

// convertForeignKeyActionCode spells out the delete and update action
// pg_constraint writes as one character.
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

// getForeignKeyColumnsAndReferencedColumns reads the two column lists out of a
// constraint definition such as
// "FOREIGN KEY (a, b) REFERENCES s.t(x, y)". The lists are found by their
// parentheses rather than by a pattern, because a quoted identifier may hold a
// parenthesis or a comma.
func getForeignKeyColumnsAndReferencedColumns(definition string) ([]string, []string, error) {
	columnList, rest, err := parenthesizedList(definition, "FOREIGN KEY")
	if err != nil {
		return nil, nil, err
	}
	const references = "REFERENCES"
	rest = strings.TrimSpace(rest)
	if len(rest) < len(references) || !strings.EqualFold(rest[:len(references)], references) {
		return nil, nil, fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	rest = rest[len(references):]
	open := indexOutsideQuotes(rest, '(')
	if open < 0 {
		return nil, nil, fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	referencedList, _, err := parentheses(rest[open:])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid foreign key definition %q: %w", definition, err)
	}

	columns, err := getColumnList(columnList)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid foreign key definition %q: %w", definition, err)
	}
	referencedColumns, err := getColumnList(referencedList)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid foreign key definition %q: %w", definition, err)
	}
	return columns, referencedColumns, nil
}

// parenthesizedList reads the parenthesis group a definition opens with after
// prefix.
func parenthesizedList(definition, prefix string) (string, string, error) {
	trimmed := strings.TrimSpace(definition)
	if len(trimmed) < len(prefix) || !strings.EqualFold(trimmed[:len(prefix)], prefix) {
		return "", "", fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	inner, rest, err := parentheses(trimmed[len(prefix):])
	if err != nil {
		return "", "", fmt.Errorf("invalid foreign key definition: %q", definition)
	}
	return inner, rest, nil
}

// parentheses reads a balanced parenthesis group and returns what it holds and
// what follows it, skipping the bytes a quoted identifier holds.
func parentheses(s string) (string, string, error) {
	s = strings.TrimSpace(s)
	if s == "" || s[0] != '(' {
		return "", "", fmt.Errorf("no parenthesized list: %q", s)
	}
	depth := 0
	quoted := false
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '"':
			if quoted && i+1 < len(s) && s[i+1] == '"' {
				i++
				continue
			}
			quoted = !quoted
		case quoted:
		case c == '(':
			depth++
		case c == ')':
			depth--
			if depth == 0 {
				return s[1:i], s[i+1:], nil
			}
		}
	}
	return "", "", fmt.Errorf("unbalanced parenthesis: %q", s)
}

// indexOutsideQuotes returns the offset of target in s, skipping the bytes a
// quoted identifier holds, or -1 when it is not there.
func indexOutsideQuotes(s string, target byte) int {
	quoted := false
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '"':
			if quoted && i+1 < len(s) && s[i+1] == '"' {
				i++
				continue
			}
			quoted = !quoted
		case quoted:
		case c == target:
			return i
		}
	}
	return -1
}

// getColumnList splits a comma-separated column list, where a name may be
// quoted and a quoted name may hold a comma or a doubled quote.
func getColumnList(definition string) ([]string, error) {
	var result []string
	var name strings.Builder
	quoted := false
	for i := 0; i < len(definition); i++ {
		switch c := definition[i]; {
		case c == '"':
			if quoted && i+1 < len(definition) && definition[i+1] == '"' {
				name.WriteByte('"')
				i++
				continue
			}
			quoted = !quoted
		case c == ',' && !quoted:
			result = append(result, name.String())
			name.Reset()
		case !quoted && isSpaceByte(c):
			// The spaces between names are separators; a space inside a quoted
			// name is part of it.
		default:
			name.WriteByte(c)
		}
	}
	if quoted {
		return nil, fmt.Errorf("invalid column list definition: %q", definition)
	}
	result = append(result, name.String())
	return result, nil
}

// isSpaceByte reports whether c separates one name from another.
func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// convertYesNo reads the is_nullable answers Redshift's INFORMATION_SCHEMA and
// SVV_ALL_COLUMNS give.
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

// columnType spells the type INFORMATION_SCHEMA.COLUMNS reports, whose data_type
// names the base type and whose other columns carry its modifiers: a character
// or bit type's length, and a numeric type's precision and scale.
func columnType(dataType string, characterMaxLength, numericPrecision, numericScale sql.NullInt64) string {
	switch dataType {
	case "character", "character varying", "bit", "bit varying":
		// For character varying(n), the character maximum length is n. For
		// character without a length, character_maximum_length is null and no
		// length is appended.
		if characterMaxLength.Valid {
			return fmt.Sprintf("%s(%d)", dataType, characterMaxLength.Int64)
		}
	case "numeric":
		// numeric and decimal have no modifier in data_type, so the precision
		// and scale are what tells numeric(12,2) apart.
		if numericPrecision.Valid && numericScale.Valid {
			return fmt.Sprintf("%s(%d,%d)", dataType, numericPrecision.Int64, numericScale.Int64)
		}
	}
	return dataType
}
