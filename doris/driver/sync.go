package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/doris/ast"
	"github.com/bytebase/omni/doris/parser"
	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// errUnknownDatabase is the client error number for a database the server does
// not have.
const errUnknownDatabase = 1049

const (
	baseTableType        = "BASE TABLE"
	viewTableType        = "VIEW"
	materializedViewType = "MATERIALIZED VIEW"
	autoIncrementSymbol  = "AUTO_INCREMENT"
)

// The COLUMN_KEY values information_schema.columns states for the key columns
// of a table, which are a leading prefix of its columns.
const (
	keyColumnAggregate = "AGG"
	keyColumnDuplicate = "DUP"
	keyColumnUnique    = "UNI"
)

// systemDatabaseExclusion returns the quoted, comma-joined system databases to
// exclude from instance discovery, matching the parser-side system-database
// classification. Doris has no `sys` metadatabase.
func systemDatabaseExclusion() string {
	return strings.Join([]string{"'__internal_schema'", "'_statistics_'", "'information_schema'", "'mysql'"}, ", ")
}

// SyncInstance reads the server version, lower_case_table_names, and
// non-system databases.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, _, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}

	var lowerCaseTableNames int32
	if text, err := d.getServerVariable(ctx, "lower_case_table_names"); err == nil {
		if v, err := strconv.ParseInt(text, 10, 32); err == nil {
			lowerCaseTableNames = int32(v)
		}
	}

	query := fmt.Sprintf(`
		SELECT
			SCHEMA_NAME,
			DEFAULT_CHARACTER_SET_NAME,
			DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE LOWER(SCHEMA_NAME) NOT IN (%s)
		ORDER BY SCHEMA_NAME`, systemDatabaseExclusion())
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var databases []*metadatapb.DatabaseSchemaMetadata
	for rows.Next() {
		database := &metadatapb.DatabaseSchemaMetadata{}
		if err := rows.Scan(
			&database.Name,
			&database.CharacterSet,
			&database.Collation,
		); err != nil {
			return nil, err
		}
		database.CharacterSet = ""
		database.Collation = ""
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &driver.Instance{
		Version:             version,
		Databases:           databases,
		LowerCaseTableNames: lowerCaseTableNames,
	}, nil
}

// isMaterializedView reports whether a table row is a materialized view, based on its
// table type and membership in the materialized-view set. Doris reports materialized
// views with TABLE_TYPE='BASE TABLE', and some setups with 'MATERIALIZED VIEW'; in
// either case, presence in materializedViewMap is the authoritative signal.
func isMaterializedView(tableType string, name string, materializedViewMap map[string]*metadatapb.MaterializedViewMetadata) bool {
	if _, ok := materializedViewMap[name]; !ok {
		return false
	}
	switch tableType {
	case baseTableType, viewTableType, materializedViewType:
		return true
	default:
		return false
	}
}

// getMaterializedViews returns the materialized views of the current database keyed by
// table name. The query is best-effort: a version that does not support it is logged
// and yields an empty set rather than failing the sync.
func (d *Driver) getMaterializedViews(ctx context.Context) (map[string]*metadatapb.MaterializedViewMetadata, error) {
	materializedViewMap := make(map[string]*metadatapb.MaterializedViewMetadata)

	// Doris exposes materialized views via the mv_infos() table-valued function,
	// whose RefreshInfo is the refresh clause SHOW CREATE MATERIALIZED VIEW
	// repeats. https://doris.apache.org/docs/sql-manual/sql-functions/table-valued-functions/mv_infos
	query := fmt.Sprintf(`SELECT Name, RefreshInfo, QuerySql FROM mv_infos("database"="%s") ORDER BY Name`, d.databaseName)

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		// The catalog may be unavailable on older engine versions; an empty set
		// leaves the sync without materialized views rather than failing it.
		return materializedViewMap, nil
	}
	defer rows.Close()

	for rows.Next() {
		materializedView := &metadatapb.MaterializedViewMetadata{}
		var refreshInfo sql.NullString
		if err := rows.Scan(&materializedView.Name, &refreshInfo, &materializedView.Definition); err != nil {
			return nil, err
		}
		materializedView.RefreshInfo = refreshInfo.String
		materializedViewMap[materializedView.Name] = materializedView
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return materializedViewMap, nil
}

// tableRow is one information_schema.TABLES row, held until the table loop
// because that loop asks the server more questions and the tape's pool holds
// one connection.
type tableRow struct {
	name, tableType, createOptions, comment string
	rowCount, dataSize, indexSize, dataFree int64
}

// tableClauses is what SHOW CREATE TABLE states about a table beyond the
// dictionaries: the distribution and properties of the table, and the column
// attributes information_schema loses (an aggregate function, a generated
// column's expression, an auto-increment start).
type tableClauses struct {
	distribution *metadatapb.TableDistributionMetadata
	properties   map[string]string
	columns      map[string]*columnClauses
	// indexProperties is keyed by index name.
	indexProperties map[string]map[string]string
}

type columnClauses struct {
	aggregateFunction string
	generation        *metadatapb.GenerationMetadata
	autoIncrement     *metadatapb.AutoIncrementMetadata
}

// SyncDatabase reads the full schema metadata of Config.Database, returning
// driver.ErrNotFound when that database does not exist on the server.
func (d *Driver) SyncDatabase(ctx context.Context) (meta *metadatapb.DatabaseSchemaMetadata, err error) {
	defer func() {
		if err != nil {
			err = wrapSyncDatabaseErr(err)
		}
	}()
	schemaMetadata := &metadatapb.SchemaMetadata{
		Name: "",
	}

	// Query column info. COLUMN_KEY names the key model on each key column;
	// the column dictionaries carry no other table-level fact.
	columnMap := make(map[string][]*metadatapb.ColumnMetadata)
	keyMap := make(map[string]*metadatapb.TableKeyMetadata)
	columnQuery := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			IFNULL(COLUMN_NAME, ''),
			ORDINAL_POSITION,
			IFNULL(COLUMN_KEY, ''),
			COLUMN_DEFAULT,
			IS_NULLABLE,
			COLUMN_TYPE,
			IFNULL(CHARACTER_SET_NAME, ''),
			IFNULL(COLLATION_NAME, ''),
			COLUMN_COMMENT,
			EXTRA
		FROM information_schema.columns
		WHERE TABLE_SCHEMA = '%s'
		ORDER BY TABLE_NAME, ORDINAL_POSITION`, d.databaseName)
	columnRows, err := d.db.QueryContext(ctx, columnQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}
	defer columnRows.Close()
	for columnRows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var tableName, columnKey, nullable, extra string
		var defaultStr sql.NullString
		if err := columnRows.Scan(
			&tableName,
			&column.Name,
			&column.Position,
			&columnKey,
			&defaultStr,
			&nullable,
			&column.Type,
			&column.CharacterSet,
			&column.Collation,
			&column.Comment,
			&extra,
		); err != nil {
			return nil, err
		}
		if model := keyModel(columnKey); model != metadatapb.TableKeyMetadata_MODEL_UNSPECIFIED {
			key := keyMap[tableName]
			if key == nil {
				key = &metadatapb.TableKeyMetadata{Model: model}
				keyMap[tableName] = key
			}
			if key.Model == model {
				key.Columns = append(key.Columns, column.Name)
			}
		}
		isNullBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		// Doris reports the empty value of an HLL or BITMAP column as a
		// one-byte NUL COLUMN_DEFAULT. It is the type's zero value, not a
		// default the definition states; a BITMAP column's BITMAP_EMPTY is
		// implied by its type.
		if defaultStr.Valid && isOpaqueType(column.Type) && defaultStr.String == "\x00" {
			defaultStr.Valid = false
		}
		if defaultStr.Valid {
			// Doris is MySQL-compatible, so use MySQL's default handling approach.
			unquotedDefault := unquoteMySQLString(defaultStr.String)
			switch {
			case isCurrentTimestampLike(unquotedDefault):
				column.Default = unquotedDefault
			case strings.Contains(extra, "DEFAULT_GENERATED"):
				unescapedDefault := unescapeExpressionDefault(unquotedDefault)
				column.Default = fmt.Sprintf("(%s)", unescapedDefault)
			default:
				// For non-generated and non CURRENT_XXX default value, preserve quotes for mysqldump compatibility
				column.Default = defaultStr.String
			}
		} else if strings.Contains(strings.ToUpper(extra), autoIncrementSymbol) {
			// Use the upper case to consistent with MySQL Dump.
			column.Default = autoIncrementSymbol
		} else if isNullBool {
			// This is NULL if the column has an explicit default of NULL,
			// or if the column definition includes no DEFAULT clause.
			// https://dev.mysql.com/doc/refman/8.0/en/information-schema-columns-table.html
			column.Default = "NULL"
		}
		column.Nullable = isNullBool

		columnMap[tableName] = append(columnMap[tableName], column)
	}
	if err := columnRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}
	_ = columnRows.Close()

	// Query view info.
	viewMap := make(map[string]*metadatapb.ViewMetadata)
	viewQuery := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			VIEW_DEFINITION
		FROM information_schema.VIEWS
		WHERE TABLE_SCHEMA = '%s'`, d.databaseName)
	viewRows, err := d.db.QueryContext(ctx, viewQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
	defer viewRows.Close()
	for viewRows.Next() {
		view := &metadatapb.ViewMetadata{}
		if err := viewRows.Scan(
			&view.Name,
			&view.Definition,
		); err != nil {
			return nil, err
		}
		view.Columns = columnMap[view.Name]
		viewMap[view.Name] = view
	}
	if err := viewRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
	_ = viewRows.Close()

	// Query materialized view info (engine-specific catalog; best-effort).
	materializedViewMap, err := d.getMaterializedViews(ctx)
	if err != nil {
		return nil, err
	}

	// Query table info.
	tableQuery := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			TABLE_TYPE,
			IFNULL(TABLE_ROWS, 0),
			IFNULL(DATA_LENGTH, 0),
			IFNULL(INDEX_LENGTH, 0),
			IFNULL(DATA_FREE, 0),
			IFNULL(CREATE_OPTIONS, ''),
			IFNULL(TABLE_COMMENT, '')
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = '%s'
		ORDER BY TABLE_NAME`, d.databaseName)
	tableRows, err := d.db.QueryContext(ctx, tableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}
	defer tableRows.Close()
	var tables []tableRow
	for tableRows.Next() {
		var t tableRow
		if err := tableRows.Scan(
			&t.name,
			&t.tableType,
			&t.rowCount,
			&t.dataSize,
			&t.indexSize,
			&t.dataFree,
			&t.createOptions,
			&t.comment,
		); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	if err := tableRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}
	_ = tableRows.Close()

	for _, table := range tables {
		// Doris reports materialized views as TABLE_TYPE='BASE TABLE';
		// route any row whose name is in the materialized-view set to MaterializedViews.
		if isMaterializedView(table.tableType, table.name, materializedViewMap) {
			materializedView := materializedViewMap[table.name]
			materializedView.Comment = table.comment
			schemaMetadata.MaterializedViews = append(schemaMetadata.MaterializedViews, materializedView)
			continue
		}
		switch table.tableType {
		case baseTableType:
			clauses, err := d.getTableClauses(ctx, table.name)
			if err != nil {
				return nil, err
			}
			columns := columnMap[table.name]
			for _, column := range columns {
				if c, ok := clauses.columns[column.Name]; ok {
					column.AggregateFunction = c.aggregateFunction
					column.Generation = c.generation
					column.AutoIncrement = c.autoIncrement
				}
			}
			indexes, err := d.getTableIndexes(ctx, table.name, clauses.indexProperties)
			if err != nil {
				return nil, err
			}
			schemaMetadata.Tables = append(schemaMetadata.Tables, &metadatapb.TableMetadata{
				Name:          table.name,
				Columns:       columns,
				Indexes:       indexes,
				RowCount:      table.rowCount,
				DataSize:      table.dataSize,
				IndexSize:     table.indexSize,
				DataFree:      table.dataFree,
				CreateOptions: table.createOptions,
				Comment:       table.comment,
				Key:           keyMap[table.name],
				Distribution:  clauses.distribution,
				Properties:    clauses.properties,
			})
		case viewTableType:
			if view, ok := viewMap[table.name]; ok {
				view.Comment = table.comment
				schemaMetadata.Views = append(schemaMetadata.Views, view)
			}
		default:
			// Includes 'MATERIALIZED VIEW' rows absent from the materialized-view set
			// and any other unrecognized type.
		}
	}

	procedures, err := d.getProcedures(ctx)
	if err != nil {
		return nil, err
	}
	schemaMetadata.Procedures = procedures

	databaseMetadata := &metadatapb.DatabaseSchemaMetadata{
		Name:    d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}
	// Query db info.
	databaseQuery := fmt.Sprintf(`
		SELECT
			DEFAULT_CHARACTER_SET_NAME,
			DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE SCHEMA_NAME = '%s'`, d.databaseName)
	if err := d.db.QueryRowContext(ctx, databaseQuery).Scan(
		&databaseMetadata.CharacterSet,
		&databaseMetadata.Collation,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: database %q not found", driver.ErrNotFound, d.databaseName)
		}
		return nil, fmt.Errorf("failed to execute query %q: %w", databaseQuery, err)
	}
	// "characterSet":"utf8\u0000", "collation":"utf8_general_ci\u0000".
	// ERROR: unsupported Unicode escape sequence (SQLSTATE 22P05).
	databaseMetadata.CharacterSet = ""
	databaseMetadata.Collation = ""

	return databaseMetadata, nil
}

// getTableClauses reads the clauses SHOW CREATE TABLE states for a table. The
// column dictionaries answer no aggregate function, generated expression or
// auto-increment start for Doris, and no distribution or PROPERTIES at all.
func (d *Driver) getTableClauses(ctx context.Context, tableName string) (*tableClauses, error) {
	query := fmt.Sprintf("SHOW CREATE TABLE %s.%s", quoteBacktick(d.databaseName), quoteBacktick(tableName))
	definition, err := d.showCreate(ctx, query, "Create Table")
	if err != nil {
		return nil, err
	}
	clauses := &tableClauses{columns: map[string]*columnClauses{}, indexProperties: map[string]map[string]string{}}
	// The definition is the server's own, so a shape omni does not parse yet
	// costs this table its clause facts rather than the whole sync, as
	// mysql/driver's partition count does. The golden pins the shapes the
	// corpus covers, so a parse regression still fails the test.
	stmt, err := parseCreateTable(definition)
	if err != nil {
		return clauses, nil
	}

	if dist := stmt.DistributedBy; dist != nil {
		clauses.distribution = &metadatapb.TableDistributionMetadata{
			Method:      distributionMethod(dist.Type),
			Columns:     dist.Columns,
			Buckets:     int32(dist.Buckets),
			BucketsAuto: dist.Auto,
		}
	}
	if len(stmt.Properties) > 0 {
		clauses.properties = make(map[string]string, len(stmt.Properties))
		for _, property := range stmt.Properties {
			clauses.properties[property.Key] = property.Value
		}
	}
	for _, index := range stmt.Indexes {
		if len(index.Properties) == 0 {
			continue
		}
		properties := make(map[string]string, len(index.Properties))
		for _, property := range index.Properties {
			properties[property.Key] = property.Value
		}
		clauses.indexProperties[index.Name] = properties
	}
	for _, column := range stmt.Columns {
		c := &columnClauses{aggregateFunction: column.AggType}
		if column.Generated != nil {
			loc := ast.NodeLoc(column.Generated)
			c.generation = &metadatapb.GenerationMetadata{Expression: definition[loc.Start:loc.End]}
		}
		if column.AutoInc {
			c.autoIncrement = &metadatapb.AutoIncrementMetadata{}
			if column.AutoIncStart != nil {
				c.autoIncrement.Start = *column.AutoIncStart
			}
		}
		if c.aggregateFunction != "" || c.generation != nil || c.autoIncrement != nil {
			clauses.columns[column.Name] = c
		}
	}
	return clauses, nil
}

// getTableIndexes reads the indexes SHOW INDEX states for a table. Doris's
// information_schema.statistics is empty, so SHOW INDEX is the only dictionary
// that lists them; the properties of an index come from the SHOW CREATE TABLE
// clause, which states them as a PROPERTIES list rather than as the single
// rendered string SHOW INDEX holds.
func (d *Driver) getTableIndexes(ctx context.Context, tableName string, properties map[string]map[string]string) ([]*metadatapb.IndexMetadata, error) {
	query := fmt.Sprintf("SHOW INDEX FROM %s.%s", quoteBacktick(d.databaseName), quoteBacktick(tableName))
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	position := make(map[string]int, len(columnNames))
	for i, name := range columnNames {
		position[strings.ToLower(name)] = i
	}
	for _, name := range []string{"key_name", "column_name", "index_type", "comment"} {
		if _, ok := position[name]; !ok {
			return nil, fmt.Errorf("the answer to %q holds no %q column", query, name)
		}
	}

	var indexes []*metadatapb.IndexMetadata
	byName := make(map[string]*metadatapb.IndexMetadata)
	for rows.Next() {
		values := make([]sql.NullString, len(columnNames))
		dest := make([]any, len(columnNames))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		name := values[position["key_name"]].String
		index := byName[name]
		if index == nil {
			index = &metadatapb.IndexMetadata{
				Name:       name,
				Type:       values[position["index_type"]].String,
				Comment:    values[position["comment"]].String,
				Properties: properties[name],
			}
			byName[name] = index
			indexes = append(indexes, index)
		}
		index.Expressions = append(index.Expressions, values[position["column_name"]].String)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return indexes, nil
}

// getProcedures reads the procedures of the current database.
// information_schema.routines lists them and holds no ROUTINE_DEFINITION, so
// each definition comes from SHOW CREATE PROCEDURE.
func (d *Driver) getProcedures(ctx context.Context) ([]*metadatapb.ProcedureMetadata, error) {
	query := fmt.Sprintf(`
		SELECT ROUTINE_NAME
		FROM information_schema.ROUTINES
		WHERE ROUTINE_SCHEMA = '%s' AND ROUTINE_TYPE = 'PROCEDURE'
		ORDER BY ROUTINE_NAME`, d.databaseName)
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	_ = rows.Close()

	var procedures []*metadatapb.ProcedureMetadata
	for _, name := range names {
		createQuery := fmt.Sprintf("SHOW CREATE PROCEDURE %s.%s", quoteBacktick(d.databaseName), quoteBacktick(name))
		definition, err := d.showCreate(ctx, createQuery, "Create Procedure")
		if err != nil {
			return nil, err
		}
		procedure := &metadatapb.ProcedureMetadata{Name: name, Definition: definition}
		// The definition is the server's own; a body omni does not parse loses
		// the comment the definition states, not the procedure.
		if stmt, err := parseCreateProcedure(definition); err == nil {
			procedure.Name = stmt.Name.String()
			procedure.Comment = stmt.Comment
		}
		procedures = append(procedures, procedure)
	}
	return procedures, nil
}

// showCreate returns the definition a SHOW CREATE statement states in its
// colName column, which is the one Doris fills with the CREATE statement.
func (d *Driver) showCreate(ctx context.Context, query, colName string) (string, error) {
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return "", err
	}
	definitionColumn := -1
	for i, name := range columnNames {
		if strings.EqualFold(name, colName) {
			definitionColumn = i
		}
	}
	if definitionColumn < 0 {
		return "", fmt.Errorf("the answer to %q holds no %s column", query, colName)
	}

	var definition sql.NullString
	var read bool
	for rows.Next() {
		dest := make([]any, len(columnNames))
		for i := range columnNames {
			if i == definitionColumn {
				dest[i] = &definition
			} else {
				dest[i] = new(sql.NullString)
			}
		}
		if err := rows.Scan(dest...); err != nil {
			return "", err
		}
		read = true
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if !read || !definition.Valid {
		return "", fmt.Errorf("the answer to %q states no %s", query, colName)
	}
	return definition.String, nil
}

func parseCreateTable(definition string) (*ast.CreateTableStmt, error) {
	file, err := parser.Parse(definition)
	if err != nil {
		return nil, err
	}
	if len(file.Stmts) != 1 {
		return nil, fmt.Errorf("holds %d statements, want 1", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		return nil, fmt.Errorf("holds a %T, want *ast.CreateTableStmt", file.Stmts[0])
	}
	return stmt, nil
}

func parseCreateProcedure(definition string) (*ast.CreateProcedureStmt, error) {
	file, err := parser.Parse(definition)
	if err != nil {
		return nil, err
	}
	if len(file.Stmts) != 1 {
		return nil, fmt.Errorf("holds %d statements, want 1", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateProcedureStmt)
	if !ok {
		return nil, fmt.Errorf("holds a %T, want *ast.CreateProcedureStmt", file.Stmts[0])
	}
	return stmt, nil
}

// keyModel maps a COLUMN_KEY value to the key model it names.
func keyModel(columnKey string) metadatapb.TableKeyMetadata_Model {
	switch columnKey {
	case keyColumnDuplicate:
		return metadatapb.TableKeyMetadata_MODEL_DUPLICATE
	case keyColumnUnique:
		return metadatapb.TableKeyMetadata_MODEL_UNIQUE
	case keyColumnAggregate:
		return metadatapb.TableKeyMetadata_MODEL_AGGREGATE
	default:
		return metadatapb.TableKeyMetadata_MODEL_UNSPECIFIED
	}
}

func distributionMethod(distributionType string) metadatapb.TableDistributionMetadata_Method {
	switch distributionType {
	case "HASH":
		return metadatapb.TableDistributionMetadata_METHOD_HASH
	case "RANDOM":
		return metadatapb.TableDistributionMetadata_METHOD_RANDOM
	default:
		return metadatapb.TableDistributionMetadata_METHOD_UNSPECIFIED
	}
}

// isOpaqueType reports whether a column's type holds a value Doris's dictionary
// can only render as its binary zero.
func isOpaqueType(columnType string) bool {
	switch strings.ToLower(columnType) {
	case "hll", "bitmap":
		return true
	default:
		return false
	}
}

// wrapSyncDatabaseErr reports a database the connection cannot reach as
// driver.ErrNotFound.
func wrapSyncDatabaseErr(err error) error {
	var mysqlErr *gomysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == errUnknownDatabase {
		return fmt.Errorf("%w: %w", driver.ErrNotFound, err)
	}
	return err
}

func unescapeExpressionDefault(s string) string {
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

func isCurrentTimestampLike(s string) bool {
	upper := strings.ToUpper(s)
	if strings.HasPrefix(upper, "CURRENT_TIMESTAMP") {
		return true
	}
	if strings.HasPrefix(upper, "CURRENT_DATE") {
		return true
	}
	return false
}

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

func stripSingleQuote(s string) string {
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	return s
}

func unquoteMySQLString(s string) string {
	if len(s) == 0 {
		return s
	}
	s = stripSingleQuote(s)
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '\\':
				result = append(result, '\\')
				i++
			case '\'':
				result = append(result, '\'')
				i++
			case '0':
				result = append(result, 0)
				i++
			case 'Z':
				result = append(result, 0x1A)
				i++
			case 'n':
				result = append(result, '\n')
				i++
			case 't':
				result = append(result, '\t')
				i++
			case 'r':
				result = append(result, '\r')
				i++
			case 'b':
				result = append(result, '\b')
				i++
			default:
				result = append(result, s[i])
			}
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}
