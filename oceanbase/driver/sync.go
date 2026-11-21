package driver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	gomysql "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	oceanbaseast "github.com/bytebase/omni/oceanbase/ast"
	oceanbaseparser "github.com/bytebase/omni/oceanbase/parser"
)

const (
	baseTableType     = "BASE TABLE"
	viewTableType     = "VIEW"
	externalTableType = "EXTERNAL TABLE"
	autoIncrement     = "AUTO_INCREMENT"
	virtualGenToken   = "VIRTUAL GENERATED"
	storedGenToken    = "STORED GENERATED"
	erBadDBError      = 1049
)

var (
	systemDatabases = []string{
		"LBACSYS",
		"ORAAUDITOR",
		"SYS",
		"__public",
		"information_schema",
		"mysql",
		"oceanbase",
		"performance_schema",
		"sys",
	}

	systemDatabaseClause = func() string {
		names := slices.Clone(systemDatabases)
		slices.Sort(names)
		quoted := make([]string, len(names))
		for i, k := range names {
			quoted[i] = fmt.Sprintf("'%s'", k)
		}
		return strings.Join(quoted, ", ")
	}()

	viewDefMatcher = regexp.MustCompile("CREATE ALGORITHM=(UNDEFINED|MERGE|TEMPTABLE) DEFINER=`([^`]+)`@`([^`]+)` SQL SECURITY (DEFINER|INVOKER) VIEW `([^`]+)`(?:\\.`([^`]+)`)?( \\(`(?:[^`]|``)+`(?:, ?`(?:[^`]|``)+`)*\\))? AS (?P<def>.+)")

	columnTypeCanonicalSynonyms = map[string]string{
		"tinyint(4)":     "tinyint",
		"smallint(6)":    "smallint",
		"mediumint(9)":   "mediumint",
		"int(11)":        "int",
		"bigint(20)":     "bigint",
		"char(1)":        "char",
		"integer":        "int",
		"boolean":        "tinyint(1)",
		"decimal(10, 0)": "decimal",
		"dec":            "decimal",
	}

	temporalKeywordDefaultRe = regexp.MustCompile(`^(?i)CURRENT_(?:TIMESTAMP|DATE|TIME)(?:\(\d+\))?$`)
	temporalTypeFspRe        = regexp.MustCompile(`(?i)^(?:datetime|timestamp|time)\((\d+)\)$`)
	onUpdateCurrentTSRe      = regexp.MustCompile(`(?i)on update CURRENT_TIMESTAMP(?:\((\d+)\))?`)
)

func wrapSyncDatabaseErr(err error) error {
	var mysqlErr *gomysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == erBadDBError {
		return fmt.Errorf("%w: %w", driver.ErrNotFound, err)
	}
	return err
}

type semverTriple struct {
	major, minor, patch int
}

func parseSemver(v string) (semverTriple, error) {
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return semverTriple{}, fmt.Errorf("invalid semantic version %q", v)
	}
	maj, err := strconv.Atoi(parts[0])
	if err != nil {
		return semverTriple{}, err
	}
	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return semverTriple{}, err
	}
	pat, err := strconv.Atoi(parts[2])
	if err != nil {
		return semverTriple{}, err
	}
	return semverTriple{maj, min, pat}, nil
}

func (s semverTriple) atLeast(maj, min, pat int) bool {
	if s.major != maj {
		return s.major > maj
	}
	if s.minor != min {
		return s.minor > min
	}
	return s.patch >= pat
}

func isStockMySQL(versionSuffix string) bool {
	return !strings.Contains(versionSuffix, "MariaDB") &&
		!strings.Contains(versionSuffix, "OceanBase") &&
		!strings.Contains(versionSuffix, "TiDB")
}

// SyncInstance reads the server version, lower_case_table_names, user accounts
// with their global grants, and non-system databases.
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

	roles := d.getInstanceRoles(ctx)

	query := `
		SELECT
			SCHEMA_NAME,
			DEFAULT_CHARACTER_SET_NAME,
			DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE LOWER(SCHEMA_NAME) NOT IN (` + systemDatabaseClause + `)
		ORDER BY SCHEMA_NAME`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var databases []*metadatapb.DatabaseSchemaMetadata
	for rows.Next() {
		database := &metadatapb.DatabaseSchemaMetadata{}
		if err := rows.Scan(&database.Name, &database.CharacterSet, &database.Collation); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &driver.Instance{
		Version:             version,
		Databases:           databases,
		Roles:               roles,
		LowerCaseTableNames: lowerCaseTableNames,
	}, nil
}

// SyncDatabase reads the full schema metadata of Config.Database, returning
// driver.ErrNotFound when that database does not exist on the server.
func (d *Driver) SyncDatabase(ctx context.Context) (meta *metadatapb.DatabaseSchemaMetadata, err error) {
	defer func() {
		if err != nil {
			err = wrapSyncDatabaseErr(err)
		}
	}()
	schemaMetadata := &metadatapb.SchemaMetadata{Name: ""}

	version, rest, err := d.getVersion(ctx)
	if err != nil {
		if d.engineVersion == "" {
			return nil, err
		}
		version, rest, err = parseVersion(d.engineVersion)
		if err != nil {
			return nil, err
		}
	}

	sv, err := parseSemver(version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OceanBase version %s to semantic version: %w", version, err)
	}
	atLeast800 := sv.atLeast(8, 0, 0)
	atLeast803 := sv.atLeast(8, 0, 3)
	atLeast8013 := sv.atLeast(8, 0, 13)
	atLeast8016 := sv.atLeast(8, 0, 16)
	atLeast570 := sv.atLeast(5, 7, 0)
	isMariaDB := strings.Contains(rest, "MariaDB")
	isOceanBase := strings.Contains(rest, "OceanBase")
	stockMySQL := isStockMySQL(rest)

	binaryDefaultFmt := binaryDefaultVerbatim
	if stockMySQL {
		if atLeast800 {
			binaryDefaultFmt = binaryDefaultHexNotation
		} else {
			binaryDefaultFmt = binaryDefaultRawBytes
		}
	}

	indexMap := make(map[string]map[string]*metadatapb.IndexMetadata)
	indexQuery := `
		SELECT
			TABLE_NAME,
			INDEX_NAME,
			COLUMN_NAME,
			COLLATION,
			IFNULL(SUB_PART, -1),
			'',
			SEQ_IN_INDEX,
			INDEX_TYPE,
			NON_UNIQUE,
			1,
			INDEX_COMMENT
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`
	// OceanBase reports a 5.7 version but has EXPRESSION and IS_VISIBLE, which
	// a functional key part needs: without EXPRESSION the sync would name the
	// hidden column holding it.
	if (atLeast8013 && !isMariaDB) || isOceanBase {
		indexQuery = `
			SELECT
				TABLE_NAME,
				INDEX_NAME,
				COLUMN_NAME,
				COLLATION,
				IFNULL(SUB_PART, -1),
				EXPRESSION,
				SEQ_IN_INDEX,
				INDEX_TYPE,
				NON_UNIQUE,
				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
				INDEX_COMMENT
			FROM information_schema.STATISTICS
			WHERE TABLE_SCHEMA = ?
			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`
	}
	indexRows, err := d.db.QueryContext(ctx, indexQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", indexQuery, err)
	}
	defer indexRows.Close()
	for indexRows.Next() {
		var tableName, indexName, indexType, comment, expression string
		var columnName, expressionName, collation sql.NullString
		var position int
		var subPart int64
		var nonUnique, visible bool
		if err := indexRows.Scan(
			&tableName,
			&indexName,
			&columnName,
			&collation,
			&subPart,
			&expressionName,
			&position,
			&indexType,
			&nonUnique,
			&visible,
			&comment,
		); err != nil {
			return nil, err
		}
		// A functional key part has its expression in EXPRESSION; oceanbase-ce
		// 4.3.5.6 also names the hidden column that holds it in COLUMN_NAME
		// (SYS_NC17$), which a restore cannot index.
		if expressionName.Valid && expressionName.String != "" {
			expression = fmt.Sprintf("(%s)", expressionName.String)
		} else if columnName.Valid {
			expression = columnName.String
		}
		desc := collation.Valid && collation.String == "D"
		if _, ok := indexMap[tableName]; !ok {
			indexMap[tableName] = make(map[string]*metadatapb.IndexMetadata)
		}
		if _, ok := indexMap[tableName][indexName]; !ok {
			indexMap[tableName][indexName] = &metadatapb.IndexMetadata{
				Name:    indexName,
				Type:    indexType,
				Unique:  !nonUnique,
				Primary: indexName == "PRIMARY",
				Visible: visible,
				Comment: comment,
			}
		}
		idx := indexMap[tableName][indexName]
		idx.Expressions = append(idx.Expressions, expression)
		idx.KeyLength = append(idx.KeyLength, subPart)
		idx.Descending = append(idx.Descending, desc)
	}
	if err := indexRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", indexQuery, err)
	}

	columnMap := make(map[string][]*metadatapb.ColumnMetadata)
	sridSelect := "NULL"
	if atLeast803 && stockMySQL {
		sridSelect = "SRS_ID"
	}
	genSelect := "convert(GENERATION_EXPRESSION using BINARY)"
	if !atLeast570 {
		genSelect = "NULL"
	}
	columnQuery := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			IFNULL(COLUMN_NAME, ''),
			ORDINAL_POSITION,
			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
			IS_NULLABLE,
			COLUMN_TYPE,
			IFNULL(CHARACTER_SET_NAME, ''),
			IFNULL(COLLATION_NAME, ''),
			QUOTE(COLUMN_COMMENT),
			%s,
			EXTRA,
			%s
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME, ORDINAL_POSITION`, genSelect, sridSelect)
	columnRows, err := d.db.QueryContext(ctx, columnQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}
	defer columnRows.Close()
	for columnRows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var tableName, nullable, extra, tp string
		var defaultStr sql.NullString
		var generationExpr []byte
		var srid sql.NullInt64
		if err := columnRows.Scan(
			&tableName,
			&column.Name,
			&column.Position,
			&defaultStr,
			&nullable,
			&tp,
			&column.CharacterSet,
			&column.Collation,
			&column.Comment,
			&generationExpr,
			&extra,
			&srid,
		); err != nil {
			return nil, err
		}
		column.Comment = stripSingleQuote(column.Comment)
		nullableBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		column.Type = canonicalColumnTypeSynonym(tp)
		column.Nullable = nullableBool
		setColumnMetadataDefault(column, defaultStr, nullableBool, extra, binaryDefaultFmt)
		if srid.Valid {
			v := uint32(srid.Int64)
			column.Srid = &v
		}
		if stockMySQL && isInvisibleColumnExtra(extra) {
			column.IsInvisible = true
		}
		columnMap[tableName] = append(columnMap[tableName], column)

		text := string(generationExpr)
		if containsInvisibleChars(generationExpr) {
			if iso, err := utf8ToISO88591(text); err == nil {
				text = iso
			}
		}
		text = strings.ReplaceAll(text, `\'`, `'`)
		if genType, ok := generatedColumnType(extra); ok && len(generationExpr) != 0 {
			column.Generation = &metadatapb.GenerationMetadata{
				Type:       genType,
				Expression: text,
			}
			// OceanBase fills COLUMN_DEFAULT with the generation expression too. The
			// generation clause is the only statement of it, and a column carries one or
			// the other, never both.
			column.Default = ""
		}
	}
	if err := columnRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}

	checkMap := make(map[string][]*metadatapb.CheckConstraintMetadata)
	if atLeast8016 || isOceanBase {
		checkQuery := getCheckConstraintQuery(isMariaDB)
		checkRows, err := d.db.QueryContext(ctx, checkQuery, d.databaseName)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", checkQuery, err)
		}
		defer checkRows.Close()
		for checkRows.Next() {
			check := &metadatapb.CheckConstraintMetadata{}
			var tableName, enforced string
			if err := checkRows.Scan(&tableName, &check.Name, &check.Expression, &enforced); err != nil {
				return nil, err
			}
			check.Expression = unescapeCheckClause(check.Expression)
			check.NotEnforced = strings.EqualFold(enforced, "NO")
			checkMap[tableName] = append(checkMap[tableName], check)
		}
		if err := checkRows.Err(); err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", checkQuery, err)
		}
	}

	viewMap := make(map[string]*metadatapb.ViewMetadata)
	var viewNames []string
	const viewQuery = `
		SELECT
			TABLE_NAME,
			VIEW_DEFINITION,
			IFNULL(CHECK_OPTION, '')
		FROM information_schema.VIEWS
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME`
	viewRows, err := d.db.QueryContext(ctx, viewQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
	defer viewRows.Close()
	for viewRows.Next() {
		view := &metadatapb.ViewMetadata{}
		var checkOption string
		if err := viewRows.Scan(&view.Name, &view.Definition, &checkOption); err != nil {
			return nil, err
		}
		if checkOption != "" && !strings.EqualFold(checkOption, "NONE") {
			view.CheckOption = strings.ToUpper(checkOption)
		}
		view.Definition = stripViewCheckOptionSuffix(view.Definition)
		view.Columns = columnMap[view.Name]
		viewMap[view.Name] = view
		viewNames = append(viewNames, view.Name)
	}
	if err := viewRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
	_ = viewRows.Close()

	for _, viewName := range viewNames {
		def, err := d.reconcileViewDefinition(ctx, d.databaseName, viewName)
		if err != nil {
			return nil, err
		}
		if def != "" {
			viewMap[viewName].Definition = def
		}
	}

	triggerMap, err := d.getTriggerList(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}

	events, err := d.getEventList(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	schemaMetadata.Events = events

	foreignKeysMap, err := d.getForeignKeyList(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}

	functions, procedures, err := d.syncRoutines(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	schemaMetadata.Functions = functions
	schemaMetadata.Procedures = procedures

	const tableQuery = `
		SELECT
			TABLES.TABLE_NAME,
			TABLES.TABLE_TYPE,
			IFNULL(TABLES.ENGINE, ''),
			IFNULL(TABLES.TABLE_COLLATION, ''),
			IFNULL(TABLES.TABLE_ROWS, 0),
			IFNULL(TABLES.DATA_LENGTH, 0),
			IFNULL(TABLES.INDEX_LENGTH, 0),
			IFNULL(TABLES.DATA_FREE, 0),
			IFNULL(TABLES.CREATE_OPTIONS, ''),
			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
			IFNULL(CCSA.CHARACTER_SET_NAME, '')
		FROM information_schema.TABLES TABLES
		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME`
	tableRows, err := d.db.QueryContext(ctx, tableQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}
	defer tableRows.Close()
	for tableRows.Next() {
		var tableName, tableType, engine, collation, createOptions, comment, charset string
		var rowCount, dataSize, indexSize, dataFree int64
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
			&charset,
		); err != nil {
			return nil, err
		}
		comment = stripSingleQuote(comment)

		switch tableType {
		case baseTableType:
			tableMetadata := &metadatapb.TableMetadata{
				Name:             tableName,
				Columns:          columnMap[tableName],
				ForeignKeys:      foreignKeysMap[tableName],
				Engine:           engine,
				Collation:        collation,
				RowCount:         rowCount,
				DataSize:         dataSize,
				IndexSize:        indexSize,
				DataFree:         dataFree,
				CreateOptions:    createOptions,
				Comment:          comment,
				CheckConstraints: checkMap[tableName],
				Charset:          charset,
				Triggers:         triggerMap[tableName],
			}
			if indexes, ok := indexMap[tableName]; ok {
				indexNames := make([]string, 0, len(indexes))
				for indexName := range indexes {
					indexNames = append(indexNames, indexName)
				}
				slices.Sort(indexNames)
				for _, indexName := range indexNames {
					tableMetadata.Indexes = append(tableMetadata.Indexes, indexes[indexName])
				}
			}
			schemaMetadata.Tables = append(schemaMetadata.Tables, tableMetadata)
		case viewTableType:
			if view, ok := viewMap[tableName]; ok {
				schemaMetadata.Views = append(schemaMetadata.Views, view)
			}
		case externalTableType:
			schemaMetadata.ExternalTables = append(schemaMetadata.ExternalTables, &metadatapb.ExternalTableMetadata{
				Name:    tableName,
				Columns: columnMap[tableName],
			})
		}
	}
	if err := tableRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}
	_ = tableRows.Close()

	// An external table's location, pattern and format are in its own DDL,
	// which information_schema does not carry.
	for _, externalTable := range schemaMetadata.ExternalTables {
		location, pattern, autoRefresh, format, err := d.getExternalTableFacts(ctx, d.databaseName, externalTable.Name)
		if err != nil {
			return nil, err
		}
		externalTable.Location = location
		externalTable.Pattern = pattern
		externalTable.AutoRefresh = autoRefresh
		externalTable.Format = format
	}

	materializedViews, containerTables, err := d.listMaterializedViews(ctx, d.databaseName, schemaMetadata.Tables)
	if err != nil {
		return nil, err
	}
	schemaMetadata.MaterializedViews = materializedViews
	// A materialized view's container table holds its rows; it is not a table of the
	// database, and the view's own CREATE makes a container again.
	schemaMetadata.Tables = withoutTables(schemaMetadata.Tables, containerTables)

	hasMVLog, mvLogTables, err := d.listMaterializedViewLogs(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	if hasMVLog != nil {
		// Each materialized view log creates an mlog$_<table> backing table that is
		// recreated automatically by CREATE MATERIALIZED VIEW LOG ON <table>.
		schemaMetadata.Tables = withoutTables(schemaMetadata.Tables, mvLogTables)
		for _, table := range schemaMetadata.Tables {
			table.HasMaterializedViewLog = proto.Bool(hasMVLog[table.Name])
		}
	}
	if err := d.reconcileTablesFromShowCreate(ctx, d.databaseName, schemaMetadata.Tables); err != nil {
		return nil, err
	}

	databaseMetadata := &metadatapb.DatabaseSchemaMetadata{
		Name:    d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}
	const databaseQuery = `
		SELECT
			DEFAULT_CHARACTER_SET_NAME,
			DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE SCHEMA_NAME = ?`
	if err := d.db.QueryRowContext(ctx, databaseQuery, d.databaseName).Scan(
		&databaseMetadata.CharacterSet,
		&databaseMetadata.Collation,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("database %q: %w", d.databaseName, driver.ErrNotFound)
		}
		return nil, err
	}
	return databaseMetadata, nil
}

// listMaterializedViews reads materialized views from oceanbase.DBA_MVIEWS,
// copying container table comments and indexes and leaving the list empty when
// the account lacks SELECT on the oceanbase catalog. The second result names the
// container table of each view read: it is the view's storage, and the server
// lists it as a table of the database too.
func (d *Driver) listMaterializedViews(ctx context.Context, databaseName string, tables []*metadatapb.TableMetadata) ([]*metadatapb.MaterializedViewMetadata, []string, error) {
	const query = `
		SELECT MVIEW_NAME, CONTAINER_NAME, QUERY
		FROM oceanbase.DBA_MVIEWS
		WHERE OWNER = ?
		ORDER BY MVIEW_NAME`
	rows, err := d.db.QueryContext(ctx, query, databaseName)
	if err != nil {
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		return nil, nil, nil
	}
	defer rows.Close()

	containers := make(map[string]*metadatapb.TableMetadata, len(tables))
	for _, table := range tables {
		containers[table.Name] = table
	}
	var views []*metadatapb.MaterializedViewMetadata
	var containerNames []string
	for rows.Next() {
		var name, containerName, definition string
		if err := rows.Scan(&name, &containerName, &definition); err != nil {
			return nil, nil, err
		}
		view := &metadatapb.MaterializedViewMetadata{
			Name:       name,
			Definition: definition,
		}
		containerNames = append(containerNames, containerName)
		if container, ok := containers[containerName]; ok {
			view.Comment = container.Comment
			for _, index := range container.Indexes {
				view.Indexes = append(view.Indexes, proto.Clone(index).(*metadatapb.IndexMetadata))
			}
		}
		views = append(views, view)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return views, containerNames, nil
}

// listMaterializedViewLogs reads materialized view logs from oceanbase.DBA_MVIEW_LOGS,
// returning nil maps when the account lacks SELECT on the oceanbase catalog so
// callers leave HasMaterializedViewLog unset (unknown).
func (d *Driver) listMaterializedViewLogs(ctx context.Context, databaseName string) (map[string]bool, []string, error) {
	const query = `
		SELECT MASTER, LOG_TABLE
		FROM oceanbase.DBA_MVIEW_LOGS
		WHERE LOG_OWNER = ?
		ORDER BY MASTER`
	rows, err := d.db.QueryContext(ctx, query, databaseName)
	if err != nil {
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		if isCatalogPrivilegeOrMissingErr(err) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	hasLog := make(map[string]bool)
	var logTables []string
	for rows.Next() {
		var master, logTable string
		if err := rows.Scan(&master, &logTable); err != nil {
			return nil, nil, err
		}
		hasLog[master] = true
		if logTable != "" {
			logTables = append(logTables, logTable)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return hasLog, logTables, nil
}

func isCatalogPrivilegeOrMissingErr(err error) bool {
	var mysqlErr *gomysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}
	switch mysqlErr.Number {
	case 1044, // ER_DBACCESS_DENIED_ERROR
		1045, // ER_ACCESS_DENIED_ERROR
		1049, // ER_BAD_DB_ERROR
		1109, // ER_UNKNOWN_TABLE
		1142, // ER_TABLEACCESS_DENIED_ERROR
		1143, // ER_COLUMNACCESS_DENIED_ERROR
		1146: // ER_NO_SUCH_TABLE
		return true
	default:
		return false
	}
}

// withoutTables drops the named tables, keeping the order of the rest.
func withoutTables(tables []*metadatapb.TableMetadata, names []string) []*metadatapb.TableMetadata {
	if len(names) == 0 {
		return tables
	}
	drop := make(map[string]bool, len(names))
	for _, name := range names {
		drop[name] = true
	}
	kept := make([]*metadatapb.TableMetadata, 0, len(tables))
	for _, table := range tables {
		if !drop[table.Name] {
			kept = append(kept, table)
		}
	}
	return kept
}

// eventRow is the part of information_schema.EVENTS that states an event.
type eventRow struct {
	name, eventType, executeAt, intervalValue, intervalField string
	starts, ends, status, onCompletion, comment              string
	timeZone, sqlMode, definition                            string
}

// getEventList reads events from information_schema.EVENTS, which states the schedule,
// whether the event is enabled, and the body. OceanBase rejects SHOW CREATE EVENT, so the
// statement the dump writes is built here from those columns.
func (d *Driver) getEventList(ctx context.Context, databaseName string) ([]*metadatapb.EventMetadata, error) {
	// A caller may open its connection with parseTime, which returns a temporal column as
	// a time.Time; the cast keeps the scan on a string.
	const query = `
		SELECT
			EVENT_NAME,
			EVENT_TYPE,
			IFNULL(CAST(EXECUTE_AT AS CHAR), ''),
			IFNULL(INTERVAL_VALUE, ''),
			IFNULL(INTERVAL_FIELD, ''),
			IFNULL(CAST(STARTS AS CHAR), ''),
			IFNULL(CAST(ENDS AS CHAR), ''),
			IFNULL(STATUS, ''),
			IFNULL(ON_COMPLETION, ''),
			IFNULL(EVENT_COMMENT, ''),
			IFNULL(TIME_ZONE, ''),
			IFNULL(SQL_MODE, ''),
			IFNULL(EVENT_DEFINITION, '')
		FROM information_schema.EVENTS
		WHERE EVENT_SCHEMA = ?
		ORDER BY EVENT_NAME`
	rows, err := d.db.QueryContext(ctx, query, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var events []*metadatapb.EventMetadata
	// A database created again under a dropped one's name lists that database's
	// event until OceanBase forgets it, and information_schema cannot tell the two
	// apart. The event's name is the same, so the sync keeps the first.
	seen := make(map[string]bool)
	for rows.Next() {
		r := &eventRow{}
		if err := rows.Scan(
			&r.name,
			&r.eventType,
			&r.executeAt,
			&r.intervalValue,
			&r.intervalField,
			&r.starts,
			&r.ends,
			&r.status,
			&r.onCompletion,
			&r.comment,
			&r.timeZone,
			&r.sqlMode,
			&r.definition,
		); err != nil {
			return nil, err
		}
		if seen[r.name] {
			continue
		}
		seen[r.name] = true
		events = append(events, &metadatapb.EventMetadata{
			Name:       r.name,
			Definition: createEventStatement(r),
			TimeZone:   r.timeZone,
			SqlMode:    r.sqlMode,
			Comment:    r.comment,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return events, nil
}

// createEventStatement builds the CREATE EVENT statement one information_schema.EVENTS row
// states, which is what the dump writes back. It answers "" for a row without a body: an
// event the account may not read has nothing faithful to write.
func createEventStatement(r *eventRow) string {
	body := strings.TrimSuffix(strings.TrimSpace(r.definition), ";")
	if body == "" {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE EVENT %s ON SCHEDULE ", quoteBacktick(r.name))
	if strings.EqualFold(r.eventType, "ONE TIME") {
		fmt.Fprintf(&b, "AT '%s'", r.executeAt)
	} else {
		fmt.Fprintf(&b, "EVERY %s %s", r.intervalValue, r.intervalField)
		if r.starts != "" {
			fmt.Fprintf(&b, " STARTS '%s'", r.starts)
		}
		if r.ends != "" {
			fmt.Fprintf(&b, " ENDS '%s'", r.ends)
		}
	}
	if r.onCompletion != "" {
		fmt.Fprintf(&b, " ON COMPLETION %s", r.onCompletion)
	}
	if strings.EqualFold(r.status, "DISABLED") {
		b.WriteString(" DISABLE")
	}
	if r.comment != "" {
		fmt.Fprintf(&b, " COMMENT '%s'", escapeEventComment(r.comment, r.sqlMode))
	}
	fmt.Fprintf(&b, " DO %s", body)
	return b.String()
}

// escapeEventComment writes the comment information_schema.EVENTS states as a string literal
// the statement's own sql_mode reads back: a backslash escapes the byte after it unless the
// mode says otherwise, and a quote is always doubled.
func escapeEventComment(comment, sqlMode string) string {
	if !strings.Contains(strings.ToUpper(sqlMode), "NO_BACKSLASH_ESCAPES") {
		comment = strings.ReplaceAll(comment, `\`, `\\`)
	}
	return strings.ReplaceAll(comment, "'", "''")
}

func (d *Driver) getTriggerList(ctx context.Context, databaseName string) (map[string][]*metadatapb.TriggerMetadata, error) {
	const triggersQuery = `
	SELECT
		TRIGGER_NAME,
		EVENT_OBJECT_TABLE,
		EVENT_MANIPULATION,
		ACTION_TIMING,
		ACTION_STATEMENT,
		SQL_MODE,
		CHARACTER_SET_CLIENT,
		COLLATION_CONNECTION
	FROM INFORMATION_SCHEMA.TRIGGERS
	WHERE EVENT_OBJECT_SCHEMA = ?
	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
	`
	rows, err := d.db.QueryContext(ctx, triggersQuery, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", triggersQuery, err)
	}
	defer rows.Close()

	triggerMap := make(map[string][]*metadatapb.TriggerMetadata)
	for rows.Next() {
		var name, table, event, timing, statement, sqlMode, charsetClient, collationConnection string
		if err := rows.Scan(&name, &table, &event, &timing, &statement, &sqlMode, &charsetClient, &collationConnection); err != nil {
			return nil, err
		}
		triggerMap[table] = append(triggerMap[table], &metadatapb.TriggerMetadata{
			Name:                name,
			Event:               event,
			Timing:              timing,
			Body:                triggerBody(statement),
			SqlMode:             sqlMode,
			CharacterSetClient:  charsetClient,
			CollationConnection: collationConnection,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", triggersQuery, err)
	}
	return triggerMap, nil
}

// triggerBody returns the body of the CREATE TRIGGER statement OceanBase answers in
// ACTION_STATEMENT. MySQL answers the body alone; OceanBase repeats the statement's own
// header as `TRIGGER <name> <timing> <event> ON <table> FOR EACH ROW`, which the DDL writer
// would state after its own.
func triggerBody(statement string) string {
	const head = "TRIGGER "
	if len(statement) < len(head) || !strings.EqualFold(statement[:len(head)], head) {
		return statement
	}
	// The header ends at the first FOR EACH ROW outside a quoted name: no trigger nests in
	// another, so a body holds none.
	header := []string{"FOR", "EACH", "ROW"}
	next := 0
	for i := len(head); i < len(statement); {
		switch c := statement[i]; {
		case c == '\'' || c == '"' || c == '`':
			i = skipQuoted(statement, i)
		case isNameByte(c):
			start := i
			for i < len(statement) && isNameByte(statement[i]) {
				i++
			}
			word := statement[start:i]
			switch {
			case next < len(header) && strings.EqualFold(word, header[next]):
				next++
				if next == len(header) {
					return strings.TrimSpace(statement[i:])
				}
			case strings.EqualFold(word, header[0]):
				next = 1
			default:
				next = 0
			}
		default:
			i++
		}
	}
	return statement
}

// isNameByte reports whether c can appear in an unquoted keyword or identifier.
func isNameByte(c byte) bool {
	return c == '_' || c == '$' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}

// skipQuoted returns the offset just past the quoted run that starts at i, whose opening
// quote is statement[i]. A doubled quote escapes itself, and a backslash escapes the byte
// after it outside a backtick-quoted name.
func skipQuoted(statement string, i int) int {
	quote := statement[i]
	for j := i + 1; j < len(statement); j++ {
		switch statement[j] {
		case '\\':
			if quote != '`' {
				j++
			}
		case quote:
			if j+1 < len(statement) && statement[j+1] == quote {
				j++
				continue
			}
			return j + 1
		}
	}
	return len(statement)
}

func (d *Driver) syncRoutines(ctx context.Context, databaseName string) ([]*metadatapb.FunctionMetadata, []*metadatapb.ProcedureMetadata, error) {
	const routinesQuery = `
		SELECT
			ROUTINE_NAME,
			ROUTINE_TYPE,
			SQL_MODE,
			CHARACTER_SET_CLIENT,
			COLLATION_CONNECTION,
			DATABASE_COLLATION,
			IFNULL(ROUTINE_COMMENT, '')
		FROM
			INFORMATION_SCHEMA.ROUTINES
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
	`
	rows, err := d.db.QueryContext(ctx, routinesQuery, databaseName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", routinesQuery, err)
	}
	defer rows.Close()

	type routineRow struct {
		name, routineType, comment                             string
		sqlMode, charsetClient, collationConn, databaseCollate string
	}
	var rawRoutines []routineRow
	for rows.Next() {
		var name, routineType, comment string
		var sqlMode, charsetClient, collationConn, databaseCollate sql.NullString
		if err := rows.Scan(&name, &routineType, &sqlMode, &charsetClient, &collationConn, &databaseCollate, &comment); err != nil {
			return nil, nil, err
		}
		rawRoutines = append(rawRoutines, routineRow{
			name:            name,
			routineType:     routineType,
			comment:         comment,
			sqlMode:         sqlMode.String,
			charsetClient:   charsetClient.String,
			collationConn:   collationConn.String,
			databaseCollate: databaseCollate.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", routinesQuery, err)
	}
	_ = rows.Close()

	var functions []*metadatapb.FunctionMetadata
	var procedures []*metadatapb.ProcedureMetadata
	for _, r := range rawRoutines {
		if strings.EqualFold(r.routineType, "PROCEDURE") {
			def, err := d.getCreateProcedureStmt(ctx, databaseName, r.name)
			if err != nil {
				return nil, nil, err
			}
			procedures = append(procedures, &metadatapb.ProcedureMetadata{
				Name:                r.name,
				Definition:          def,
				SqlMode:             r.sqlMode,
				CharacterSetClient:  r.charsetClient,
				CollationConnection: r.collationConn,
				DatabaseCollation:   r.databaseCollate,
				Comment:             r.comment,
			})
		} else {
			def, err := d.getCreateFunctionStmt(ctx, databaseName, r.name)
			if err != nil {
				return nil, nil, err
			}
			functions = append(functions, &metadatapb.FunctionMetadata{
				Name:                r.name,
				Definition:          def,
				SqlMode:             r.sqlMode,
				CharacterSetClient:  r.charsetClient,
				CollationConnection: r.collationConn,
				DatabaseCollation:   r.databaseCollate,
				Comment:             r.comment,
			})
		}
	}
	return functions, procedures, nil
}

func (d *Driver) getCreateFunctionStmt(ctx context.Context, databaseName, functionName string) (string, error) {
	f, err := d.showCreateColumn(ctx, fmt.Sprintf("SHOW CREATE FUNCTION %s.%s", quoteBacktick(databaseName), quoteBacktick(functionName)), "Create Function")
	if err != nil || f == "" {
		return "", err
	}
	if idx := strings.Index(f, " FUNCTION "); idx >= 0 {
		f = "CREATE" + f[idx:]
	}
	return stripReturnsCharset(f), nil
}

func (d *Driver) getCreateProcedureStmt(ctx context.Context, databaseName, procName string) (string, error) {
	p, err := d.showCreateColumn(ctx, fmt.Sprintf("SHOW CREATE PROCEDURE %s.%s", quoteBacktick(databaseName), quoteBacktick(procName)), "Create Procedure")
	if err != nil || p == "" {
		return "", err
	}
	if idx := strings.Index(p, " PROCEDURE "); idx >= 0 {
		p = "CREATE" + p[idx:]
	}
	return p, nil
}

func (d *Driver) showCreateColumn(ctx context.Context, query, colName string) (string, error) {
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	defIdx := -1
	for i, c := range columns {
		if strings.EqualFold(c, colName) {
			defIdx = i
			break
		}
	}
	if defIdx == -1 {
		return "", fmt.Errorf("failed to find column %s", colName)
	}

	var val sql.NullString
	for rows.Next() {
		dests := make([]any, len(columns))
		for i := range columns {
			if i == defIdx {
				dests[i] = &val
			} else {
				dests[i] = new(sql.NullString)
			}
		}
		if err := rows.Scan(dests...); err != nil {
			return "", err
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if val.Valid {
		return val.String, nil
	}
	return "", nil
}

func (d *Driver) reconcileViewDefinition(ctx context.Context, databaseName, viewName string) (string, error) {
	query := fmt.Sprintf("SHOW CREATE VIEW %s.%s", quoteBacktick(databaseName), quoteBacktick(viewName))
	var viewCol, createStmt, csClient, collConn sql.NullString
	if err := d.db.QueryRowContext(ctx, query).Scan(&viewCol, &createStmt, &csClient, &collConn); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to scan row for query: %s: %w", query, err)
	}
	def, err := getViewDefFromCreateView(createStmt.String)
	if err != nil {
		return "", nil
	}
	return def, nil
}

func getViewDefFromCreateView(createView string) (string, error) {
	matches := viewDefMatcher.FindStringSubmatch(createView)
	if len(matches) == 0 {
		return "", fmt.Errorf("failed to match view definition, %s", createView)
	}
	for i, name := range viewDefMatcher.SubexpNames() {
		if name == "def" && i < len(matches) {
			return matches[i], nil
		}
	}
	return "", fmt.Errorf("failed to match view definition, %s", createView)
}

func (d *Driver) getForeignKeyList(ctx context.Context, databaseName string) (map[string][]*metadatapb.ForeignKeyMetadata, error) {
	const fkQuery = `
		SELECT
			TABLE_NAME,
			CONSTRAINT_NAME,
			REFERENCED_TABLE_NAME,
			DELETE_RULE,
			UPDATE_RULE,
			MATCH_OPTION
		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
		WHERE CONSTRAINT_SCHEMA = ?;
	`
	const kcuQuery = `
		SELECT
			TABLE_SCHEMA,
			TABLE_NAME,
			CONSTRAINT_NAME,
			COLUMN_NAME,
			REFERENCED_TABLE_SCHEMA,
			REFERENCED_COLUMN_NAME
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
	`

	type fkKey struct{ table, name string }
	fkRows, err := d.db.QueryContext(ctx, fkQuery, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", fkQuery, err)
	}
	defer fkRows.Close()
	fkMap := make(map[fkKey]*metadatapb.ForeignKeyMetadata)
	for fkRows.Next() {
		var tableName string
		var fk metadatapb.ForeignKeyMetadata
		if err := fkRows.Scan(&tableName, &fk.Name, &fk.ReferencedTable, &fk.OnDelete, &fk.OnUpdate, &fk.MatchType); err != nil {
			return nil, err
		}
		fkMap[fkKey{tableName, fk.Name}] = &fk
	}
	if err := fkRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", fkQuery, err)
	}

	kcuRows, err := d.db.QueryContext(ctx, kcuQuery, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", kcuQuery, err)
	}
	defer kcuRows.Close()
	for kcuRows.Next() {
		var tableSchema, tableName, fkName, column, referencedSchema, referencedColumn string
		if err := kcuRows.Scan(&tableSchema, &tableName, &fkName, &column, &referencedSchema, &referencedColumn); err != nil {
			return nil, err
		}
		if fk, ok := fkMap[fkKey{tableName, fkName}]; ok {
			if referencedSchema != tableSchema {
				fk.ReferencedSchema = referencedSchema
			}
			fk.Columns = append(fk.Columns, column)
			fk.ReferencedColumns = append(fk.ReferencedColumns, referencedColumn)
		}
	}
	if err := kcuRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", kcuQuery, err)
	}

	result := make(map[string][]*metadatapb.ForeignKeyMetadata)
	for key, fk := range fkMap {
		result[key.table] = append(result[key.table], fk)
	}
	for tbl := range result {
		slices.SortFunc(result[tbl], func(a, b *metadatapb.ForeignKeyMetadata) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	return result, nil
}

func canonicalColumnTypeSynonym(columnType string) string {
	if canonical, ok := columnTypeCanonicalSynonyms[strings.ToLower(columnType)]; ok {
		return canonical
	}
	return columnType
}

func unescapeCheckClause(clause string) string {
	return strings.ReplaceAll(clause, `\'`, `'`)
}

func stripViewCheckOptionSuffix(def string) string {
	trimmed := strings.TrimSpace(def)
	upper := strings.ToUpper(trimmed)
	for _, suffix := range []string{
		" WITH CASCADED CHECK OPTION",
		" WITH LOCAL CHECK OPTION",
		" WITH CHECK OPTION",
	} {
		if strings.HasSuffix(upper, suffix) {
			return strings.TrimSpace(trimmed[:len(trimmed)-len(suffix)])
		}
	}
	return def
}

func getCheckConstraintQuery(isMariaDB bool) string {
	if isMariaDB {
		return `
			SELECT
				TABLE_NAME,
				CONSTRAINT_NAME,
				CHECK_CLAUSE,
				'YES'
			FROM information_schema.CHECK_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA = ?
			ORDER BY TABLE_NAME, CONSTRAINT_NAME
		`
	}
	return `
			SELECT
				tc.TABLE_NAME,
				cc.CONSTRAINT_NAME,
				cc.CHECK_CLAUSE,
				tc.ENFORCED
			FROM information_schema.CHECK_CONSTRAINTS cc
				JOIN information_schema.TABLE_CONSTRAINTS tc
					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
		`
}

func stripReturnsCharset(def string) string {
	returnsIdx := strings.Index(def, " RETURNS ")
	if returnsIdx == -1 {
		return def
	}
	lineEnd := strings.Index(def[returnsIdx:], "\n")
	if lineEnd == -1 {
		return def
	}
	charsetIdx := strings.Index(def[returnsIdx:returnsIdx+lineEnd], " CHARSET ")
	if charsetIdx == -1 {
		return def
	}
	return def[:returnsIdx+charsetIdx] + def[returnsIdx+lineEnd:]
}

type binaryDefaultFormat int

const (
	binaryDefaultVerbatim binaryDefaultFormat = iota
	binaryDefaultRawBytes
	binaryDefaultHexNotation
)

// generatedColumnType returns the kind of generated column EXTRA states, if any.
func generatedColumnType(extra string) (metadatapb.GenerationMetadata_Type, bool) {
	upper := strings.ToUpper(extra)
	switch {
	case strings.Contains(upper, virtualGenToken):
		return metadatapb.GenerationMetadata_TYPE_VIRTUAL, true
	case strings.Contains(upper, storedGenToken):
		return metadatapb.GenerationMetadata_TYPE_STORED, true
	}
	return metadatapb.GenerationMetadata_TYPE_UNSPECIFIED, false
}

// temporalDefaultWithPrecision gives a CURRENT_TIMESTAMP or CURRENT_TIME default the
// fractional-second precision its column type declares. OceanBase answers COLUMN_DEFAULT
// without it, and the server refuses a bare CURRENT_TIMESTAMP on a DATETIME(6) column.
func temporalDefaultWithPrecision(value, columnType string) string {
	if strings.HasSuffix(value, ")") || strings.EqualFold(value, "CURRENT_DATE") {
		return value
	}
	fsp := temporalTypeFspRe.FindStringSubmatch(columnType)
	if fsp == nil {
		return value
	}
	return value + "(" + fsp[1] + ")"
}

func setColumnMetadataDefault(column *metadatapb.ColumnMetadata, defaultStr sql.NullString, nullableBool bool, extra string, binaryFormat binaryDefaultFormat) {
	if defaultStr.Valid {
		unquotedDefault := unquoteMySQLString(defaultStr.String)
		isBinaryDefault := binaryFormat != binaryDefaultVerbatim && isBinaryFamilyColumnType(column.Type)
		switch {
		case isTemporalColumnType(column.Type) && temporalKeywordDefaultRe.MatchString(unquotedDefault):
			column.Default = temporalDefaultWithPrecision(unquotedDefault, column.Type)
		case strings.Contains(extra, "DEFAULT_GENERATED"):
			column.Default = fmt.Sprintf("(%s)", unescapeExpressionDefault(unquotedDefault))
		case isBinaryDefault:
			if canonical, ok := canonicalBinaryDefault(unquotedDefault, column.Type, binaryFormat); ok {
				column.Default = canonical
			} else {
				column.Default = defaultStr.String
			}
		default:
			column.Default = defaultStr.String
		}
	} else if strings.Contains(strings.ToUpper(extra), autoIncrement) {
		column.Default = autoIncrement
	} else if nullableBool {
		column.Default = "NULL"
	}

	if match := onUpdateCurrentTSRe.FindStringSubmatch(extra); len(match) > 0 {
		if match[1] != "" {
			column.OnUpdate = fmt.Sprintf("CURRENT_TIMESTAMP(%s)", match[1])
		} else {
			column.OnUpdate = "CURRENT_TIMESTAMP"
		}
	}
}

func isBinaryFamilyColumnType(columnType string) bool {
	lower := strings.ToLower(strings.TrimSpace(columnType))
	switch {
	case lower == "binary" || strings.HasPrefix(lower, "binary("):
		return true
	case lower == "varbinary" || strings.HasPrefix(lower, "varbinary("):
		return true
	case lower == "tinyblob" || lower == "mediumblob" || lower == "longblob":
		return true
	case lower == "blob" || strings.HasPrefix(lower, "blob("):
		return true
	default:
		return false
	}
}

func isFixedBinaryColumnType(columnType string) bool {
	lower := strings.ToLower(strings.TrimSpace(columnType))
	return lower == "binary" || strings.HasPrefix(lower, "binary(")
}

func canonicalBinaryDefault(unquoted, columnType string, format binaryDefaultFormat) (string, bool) {
	var value []byte
	switch {
	case format == binaryDefaultHexNotation && unquoted == "":
		value = nil
	case format == binaryDefaultHexNotation:
		if len(unquoted) < 2 || unquoted[0] != '0' || (unquoted[1] != 'x' && unquoted[1] != 'X') {
			return "", false
		}
		decoded, err := hex.DecodeString(unquoted[2:])
		if err != nil {
			return "", false
		}
		value = decoded
	default:
		value = []byte(unquoted)
	}
	if isFixedBinaryColumnType(columnType) {
		value = bytes.TrimRight(value, "\x00")
	}
	if len(value) == 0 {
		return "''", true
	}
	if isCleanDefaultText(value) {
		escaped := strings.ReplaceAll(string(value), `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `'`, `\'`)
		return "'" + escaped + "'", true
	}
	return "0x" + hex.EncodeToString(value), true
}

func isCleanDefaultText(value []byte) bool {
	if !utf8.Valid(value) {
		return false
	}
	for _, b := range value {
		if b < 0x20 || b == 0x7F {
			return false
		}
	}
	return true
}

func isInvisibleColumnExtra(extra string) bool {
	for _, tok := range strings.Fields(extra) {
		if strings.EqualFold(tok, "INVISIBLE") {
			return true
		}
	}
	return false
}

func unescapeExpressionDefault(s string) string {
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

func isTemporalColumnType(columnType string) bool {
	lower := strings.ToLower(strings.TrimSpace(columnType))
	for _, prefix := range []string{"timestamp", "datetime", "date", "time"} {
		if lower == prefix || strings.HasPrefix(lower, prefix+"(") {
			return true
		}
	}
	return false
}

func containsInvisibleChars(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			return true
		}
		if !unicode.IsPrint(r) {
			return true
		}
		data = data[size:]
	}
	return false
}

func utf8ToISO88591(utf8Str string) (string, error) {
	out := make([]byte, 0, len(utf8Str))
	for i := 0; i < len(utf8Str); {
		r, size := utf8.DecodeRuneInString(utf8Str[i:])
		if r == utf8.RuneError && size == 1 {
			return "", errors.New("invalid UTF-8 sequence")
		}
		if r > 0xFF {
			return "", errors.New("rune out of ISO-8859-1 range")
		}
		out = append(out, byte(r))
		i += size
	}
	return string(out), nil
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

// getExternalTableFacts reads the location and FORMAT clause SHOW CREATE TABLE
// states for an external table.
func (d *Driver) getExternalTableFacts(ctx context.Context, database, table string) (location, pattern, autoRefresh, format string, err error) {
	var tableName, ddl string
	query := "SHOW CREATE TABLE " + quoteBacktick(database) + "." + quoteBacktick(table)
	if err := d.db.QueryRowContext(ctx, query).Scan(&tableName, &ddl); err != nil {
		return "", "", "", "", fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	location, pattern, autoRefresh, format = externalTableFacts(ddl)
	return location, pattern, autoRefresh, format, nil
}

// externalTableFacts reads LOCATION and FORMAT from the CREATE EXTERNAL TABLE
// statement SHOW CREATE TABLE states. Both are top-level clauses, after the
// column list, so the scan skips quoted text and parenthesized expressions.
func externalTableFacts(ddl string) (location, pattern, autoRefresh, format string) {
	depth := 0
	for i := 0; i < len(ddl); {
		switch ddl[i] {
		case '\'', '"', '`':
			i = skipQuotedText(ddl, i, ddl[i])
			continue
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
		if depth == 0 {
			if keywordAt(ddl, i, "LOCATION") {
				j := i + len("LOCATION")
				j = skipSpace(ddl, j)
				if j < len(ddl) && ddl[j] == '=' {
					j = skipSpace(ddl, j+1)
				}
				if j < len(ddl) && ddl[j] == '\'' {
					value, next := readQuotedText(ddl, j, '\'')
					location = value
					i = next
					continue
				}
			}
			if keywordAt(ddl, i, "PATTERN") {
				j := i + len("PATTERN")
				j = skipSpace(ddl, j)
				if j < len(ddl) && ddl[j] == '=' {
					j = skipSpace(ddl, j+1)
				}
				if j < len(ddl) && ddl[j] == '\'' {
					value, next := readQuotedText(ddl, j, '\'')
					pattern = value
					i = next
					continue
				}
			}
			if keywordAt(ddl, i, "AUTO_REFRESH") {
				j := i + len("AUTO_REFRESH")
				j = skipSpace(ddl, j)
				if j < len(ddl) && ddl[j] == '=' {
					j = skipSpace(ddl, j+1)
				}
				start := j
				for j < len(ddl) && ddl[j] != '\n' && ddl[j] != '\r' {
					j++
				}
				autoRefresh = strings.TrimSpace(ddl[start:j])
				i = j
				continue
			}
			if keywordAt(ddl, i, "FORMAT") {
				j := i + len("FORMAT")
				j = skipSpace(ddl, j)
				if j < len(ddl) && ddl[j] == '=' {
					j = skipSpace(ddl, j+1)
				}
				if j < len(ddl) && ddl[j] == '(' {
					value, next, ok := readParenthesized(ddl, j)
					if ok {
						format = normalizeExternalFormat(value)
						i = next
						continue
					}
				}
			}
		}
		i++
	}
	return location, pattern, autoRefresh, format
}

func keywordAt(text string, i int, keyword string) bool {
	if !strings.HasPrefix(text[i:], keyword) {
		return false
	}
	if i > 0 && isWordByte(text[i-1]) {
		return false
	}
	end := i + len(keyword)
	return end >= len(text) || !isWordByte(text[end])
}

func isWordByte(c byte) bool {
	return c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func skipSpace(text string, i int) int {
	for i < len(text) && (text[i] == ' ' || text[i] == '\t' || text[i] == '\n' || text[i] == '\r') {
		i++
	}
	return i
}

// skipQuotedText returns the index after the quoted text that starts at i.
func skipQuotedText(text string, i int, quote byte) int {
	_, next := readQuotedText(text, i, quote)
	return next
}

// readQuotedText reads text[i] as its opening quote and returns the unescaped
// content and the index after the closing quote. A doubled quote is literal, as
// is a MySQL backslash escape in a string.
func readQuotedText(text string, i int, quote byte) (string, int) {
	var b strings.Builder
	for j := i + 1; j < len(text); j++ {
		if text[j] == '\\' && quote != '`' && j+1 < len(text) {
			b.WriteString(unescapeMySQLByte(text[j+1]))
			j++
			continue
		}
		if text[j] != quote {
			b.WriteByte(text[j])
			continue
		}
		if j+1 < len(text) && text[j+1] == quote {
			b.WriteByte(quote)
			j++
			continue
		}
		return b.String(), j + 1
	}
	return b.String(), len(text)
}

// unescapeMySQLByte is the value MySQL writes for a backslash escape.
func unescapeMySQLByte(c byte) string {
	switch c {
	case '0':
		return "\x00"
	case 'b':
		return "\b"
	case 'n':
		return "\n"
	case 'r':
		return "\r"
	case 't':
		return "\t"
	case 'Z':
		return "\x1a"
	default:
		return string(c)
	}
}

// readParenthesized reads the parenthesized text that starts at i and returns
// its content without the parentheses.
func readParenthesized(text string, i int) (string, int, bool) {
	depth := 0
	for j := i; j < len(text); j++ {
		switch text[j] {
		case '\'', '"', '`':
			j = skipQuotedText(text, j, text[j]) - 1
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return text[i+1 : j], j + 1, true
			}
		}
	}
	return "", len(text), false
}

// normalizeExternalFormat collapses the whitespace SHOW CREATE TABLE writes
// between FORMAT's entries, keeping the spacing inside quoted values.
func normalizeExternalFormat(format string) string {
	var b strings.Builder
	inQuote := false
	pendingSpace := false
	for i := 0; i < len(format); i++ {
		c := format[i]
		if inQuote && c == '\\' && i+1 < len(format) {
			if pendingSpace {
				if b.Len() > 0 {
					b.WriteByte(' ')
				}
				pendingSpace = false
			}
			b.WriteByte(c)
			b.WriteByte(format[i+1])
			i++
			continue
		}
		if c == '\'' {
			if inQuote && i+1 < len(format) && format[i+1] == '\'' {
				if pendingSpace {
					if b.Len() > 0 {
						b.WriteByte(' ')
					}
					pendingSpace = false
				}
				b.WriteByte(c)
				b.WriteByte(format[i+1])
				i++
				continue
			}
			inQuote = !inQuote
		}
		if !inQuote && (c == ' ' || c == '\t' || c == '\n' || c == '\r') {
			pendingSpace = true
			continue
		}
		if pendingSpace {
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			pendingSpace = false
		}
		b.WriteByte(c)
	}
	return strings.TrimSpace(b.String())
}

// reconcileTablesFromShowCreate reads what only a table's SHOW CREATE TABLE
// states on oceanbase-ce 4.3.5.6: its WITH COLUMN GROUP list, which
// information_schema does not carry anywhere, and its vector indexes.
// information_schema.STATISTICS lists a vector index over the hidden columns
// that hold it (__vid_<n>, __type_<n>), as a BTREE, beside the hidden tables
// that support it (rowkey_vid_table, vid_rowkey_table, <index>_index_id_table,
// <index>_index_snapshot_data_table), none of which a restore can create. So an
// index SHOW CREATE TABLE does not name is dropped, and a VECTOR KEY takes its
// column and WITH (...) parameters from there.
func (d *Driver) reconcileTablesFromShowCreate(ctx context.Context, databaseName string, tables []*metadatapb.TableMetadata) error {
	for _, t := range tables {
		query := "SHOW CREATE TABLE " + quoteBacktick(databaseName) + "." + quoteBacktick(t.Name)
		var name, createTable string
		if err := d.db.QueryRowContext(ctx, query).Scan(&name, &createTable); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
		list, err := oceanbaseparser.Parse(createTable)
		if err != nil || list == nil || len(list.Items) != 1 {
			continue
		}
		stmt, ok := list.Items[0].(*oceanbaseast.CreateTableStmt)
		if !ok {
			continue
		}
		t.ColumnGroups = strings.Join(strings.Fields(stmt.ColumnGroup), " ")
		reconcileIndexesFromShowCreate(t, stmt.Constraints)
	}
	return nil
}

func reconcileIndexesFromShowCreate(t *metadatapb.TableMetadata, constraints []*oceanbaseast.Constraint) {
	named := make(map[string]*oceanbaseast.Constraint, len(constraints))
	for _, con := range constraints {
		if con == nil || con.Type == oceanbaseast.ConstrCheck || con.Type == oceanbaseast.ConstrForeignKey {
			continue
		}
		key := strings.ToLower(con.Name)
		if con.Type == oceanbaseast.ConstrPrimaryKey {
			key = "primary"
		}
		if key != "" {
			named[key] = con
		}
	}
	kept := t.Indexes[:0]
	for _, idx := range t.Indexes {
		con, ok := named[strings.ToLower(idx.Name)]
		if !ok {
			continue
		}
		// STATISTICS lists a FULLTEXT or VECTOR key over the hidden columns
		// that hold it (__word_segment_<n>, __doc_id_<n>; __vid_<n>,
		// __type_<n>); its own columns are in SHOW CREATE TABLE.
		if con.Type == oceanbaseast.ConstrFulltextIndex || con.Type == oceanbaseast.ConstrVector {
			idx.Expressions = idx.Expressions[:0]
			idx.KeyLength = idx.KeyLength[:0]
			idx.Descending = idx.Descending[:0]
			for _, col := range con.IndexColumns {
				if ref, ok := col.Expr.(*oceanbaseast.ColumnRef); ok {
					idx.Expressions = append(idx.Expressions, ref.Column)
					idx.KeyLength = append(idx.KeyLength, -1)
					idx.Descending = append(idx.Descending, false)
				}
			}
		}
		if con.Type == oceanbaseast.ConstrFulltextIndex {
			idx.Parser, idx.ParserProperties = fulltextParser(con.IndexOptions)
		}
		// A GLOBAL index spans the partitions; a restore that drops the word
		// makes it local, and a unique key without the partition key then
		// fails with 1503. oceanbase-ce 4.3.5.6 prints LOCAL for every other
		// index, written or not, so only GLOBAL is kept, as the catalog
		// exports an index written without a scope.
		for _, opt := range con.IndexOptions {
			if opt != nil && strings.EqualFold(opt.Name, "GLOBAL") {
				idx.Scope = "GLOBAL"
				idx.IsGlobal = true
			}
		}
		if con.Type == oceanbaseast.ConstrVector {
			idx.IsVector = true
			idx.Type = ""
			for _, opt := range con.IndexOptions {
				if opt != nil && strings.EqualFold(opt.Name, "WITH") {
					idx.VectorParams = vectorParams(opt)
				}
			}
		}
		kept = append(kept, idx)
	}
	t.Indexes = kept
}

// vectorParams is a VECTOR KEY's WITH (...) list as NAME=VALUE entries, in the
// order SHOW CREATE TABLE prints them with every default filled in:
// DISTANCE=L2, TYPE=HNSW, LIB=VSAG, M=16, EF_CONSTRUCTION=200, EF_SEARCH=64.
func vectorParams(opt *oceanbaseast.IndexOption) []string {
	params := make([]string, 0, len(opt.Properties))
	for _, prop := range opt.Properties {
		if prop == nil {
			continue
		}
		value := ""
		switch v := prop.Value.(type) {
		case *oceanbaseast.IntLit:
			value = strconv.FormatInt(v.Value, 10)
		case *oceanbaseast.StringLit:
			value = strings.ToUpper(v.Value)
		}
		params = append(params, strings.ToUpper(prop.Name)+"="+value)
	}
	return params
}

// fulltextParser is a FULLTEXT KEY's WITH PARSER and PARSER_PROPERTIES as SHOW
// CREATE TABLE prints them, with every default filled in. oceanbase-ce 4.3.5.6
// prints an index written with neither as WITH PARSER space
// PARSER_PROPERTIES=(min_token_size=3,max_token_size=84), which the catalog
// exports as no parser at all, so that default is left empty.
func fulltextParser(opts []*oceanbaseast.IndexOption) (string, map[string]string) {
	parser := ""
	var props map[string]string
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		switch strings.ToUpper(opt.Name) {
		case "PARSER":
			if v, ok := opt.Value.(*oceanbaseast.StringLit); ok {
				parser = strings.ToLower(v.Value)
			}
		case "PARSER_PROPERTIES":
			props = make(map[string]string, len(opt.Properties))
			for _, prop := range opt.Properties {
				if prop == nil {
					continue
				}
				key := strings.ToLower(prop.Name)
				switch v := prop.Value.(type) {
				case *oceanbaseast.IntLit:
					props[key] = strconv.FormatInt(v.Value, 10)
				case *oceanbaseast.StringLit:
					props[key] = v.Value
					if key == "ik_mode" {
						props[key] = strings.ToLower(v.Value)
					}
				}
			}
		}
	}
	if parser == "space" && len(props) == 2 && props["min_token_size"] == "3" && props["max_token_size"] == "84" {
		return "", nil
	}
	return parser, props
}
