package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gomysql "github.com/go-sql-driver/mysql"

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

// systemDatabaseExclusion returns the quoted, comma-joined system databases to
// exclude from instance discovery, matching the parser-side system-database
// classification. Doris has no `sys` metadatabase.
func systemDatabaseExclusion() string {
	return strings.Join([]string{"'__internal_schema'", "'_statistics_'", "'information_schema'", "'mysql'"}, ", ")
}

// SyncInstance reads the server version, lower_case_table_names, the
// non-system databases, and the accounts and roles with their privileges.
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

	roles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}

	return &driver.Instance{
		Version:             version,
		Databases:           databases,
		Roles:               roles,
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

	// Doris exposes materialized views via the mv_infos() table-valued function.
	// https://doris.apache.org/docs/sql-manual/sql-functions/table-valued-functions/mv_infos
	query := fmt.Sprintf(`SELECT Name, QuerySql FROM mv_infos("database"="%s")`, d.databaseName)

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		// The catalog may be unavailable on older engine versions; an empty set
		// leaves the sync without materialized views rather than failing it.
		return materializedViewMap, nil
	}
	defer rows.Close()

	for rows.Next() {
		materializedView := &metadatapb.MaterializedViewMetadata{}
		if err := rows.Scan(&materializedView.Name, &materializedView.Definition); err != nil {
			return nil, err
		}
		materializedViewMap[materializedView.Name] = materializedView
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return materializedViewMap, nil
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

	// There is not yet a way to list indexes from information_schema.
	// Query column info.
	columnMap := make(map[string][]*metadatapb.ColumnMetadata)
	columnQuery := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			IFNULL(COLUMN_NAME, ''),
			ORDINAL_POSITION,
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
		var tableName, nullable, extra string
		var defaultStr sql.NullString
		if err := columnRows.Scan(
			&tableName,
			&column.Name,
			&column.Position,
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
		isNullBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
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
			IFNULL(ENGINE, ''),
			IFNULL(TABLE_COLLATION, ''),
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
	for tableRows.Next() {
		var tableName, tableType, engine, collation, createOptions, comment string
		var rowCount, dataSize, indexSize, dataFree int64
		// Workaround TiDB bug https://github.com/pingcap/tidb/issues/27970
		var tableCollation sql.NullString
		if err := tableRows.Scan(
			&tableName,
			&tableType,
			&engine,
			&collation,
			&rowCount,
			&dataSize,
			&indexSize,
			&dataFree,
			&createOptions,
			&comment,
		); err != nil {
			return nil, err
		}

		// Doris reports materialized views as TABLE_TYPE='BASE TABLE';
		// route any row whose name is in the materialized-view set to MaterializedViews.
		if isMaterializedView(tableType, tableName, materializedViewMap) {
			materializedView := materializedViewMap[tableName]
			materializedView.Comment = comment
			schemaMetadata.MaterializedViews = append(schemaMetadata.MaterializedViews, materializedView)
			continue
		}
		switch tableType {
		case baseTableType:
			tableMetadata := &metadatapb.TableMetadata{
				Name:          tableName,
				Columns:       columnMap[tableName],
				Engine:        engine,
				Collation:     collation,
				RowCount:      rowCount,
				DataSize:      dataSize,
				IndexSize:     indexSize,
				DataFree:      dataFree,
				CreateOptions: createOptions,
				Comment:       comment,
			}
			if tableCollation.Valid {
				tableMetadata.Collation = tableCollation.String
			}
			schemaMetadata.Tables = append(schemaMetadata.Tables, tableMetadata)
		case viewTableType:
			if view, ok := viewMap[tableName]; ok {
				view.Comment = comment
				schemaMetadata.Views = append(schemaMetadata.Views, view)
			}
		default:
			// Includes 'MATERIALIZED VIEW' rows absent from the materialized-view set
			// and any other unrecognized type.
		}
	}
	if err := tableRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}

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
