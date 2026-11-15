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

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	mysqlast "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
	gomysql "github.com/go-sql-driver/mysql"
)

const erBadDBError = 1049

func wrapSyncDatabaseErr(err error) error {
	var mysqlErr *gomysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == erBadDBError {
		return fmt.Errorf("%w: %w", driver.ErrNotFound, err)
	}
	return err
}

const (
	baseTableType   = "BASE TABLE"
	viewTableType   = "VIEW"
	autoIncrement   = "AUTO_INCREMENT"
	virtualGenToken = "VIRTUAL GENERATED"
	storedGenToken  = "STORED GENERATED"
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
	currentTSPrecision       = regexp.MustCompile(`CURRENT_TIMESTAMP\((\d+)\)`)
)

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
		return nil, fmt.Errorf("failed to parse MySQL version %s to semantic version: %w", version, err)
	}
	atLeast800 := sv.atLeast(8, 0, 0)
	atLeast803 := sv.atLeast(8, 0, 3)
	atLeast8013 := sv.atLeast(8, 0, 13)
	atLeast8016 := sv.atLeast(8, 0, 16)
	atLeast570 := sv.atLeast(5, 7, 0)
	isMariaDB := strings.Contains(rest, "MariaDB")
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
	if atLeast8013 && !isMariaDB {
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
		if columnName.Valid {
			expression = columnName.String
		} else if expressionName.Valid {
			expression = fmt.Sprintf("(%s)", expressionName.String)
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
		upperExtra := strings.ToUpper(extra)
		if extra != "" && strings.Contains(upperExtra, virtualGenToken) && len(generationExpr) != 0 {
			column.Generation = &metadatapb.GenerationMetadata{
				Type:       metadatapb.GenerationMetadata_TYPE_VIRTUAL,
				Expression: text,
			}
		} else if extra != "" && strings.Contains(upperExtra, storedGenToken) && len(generationExpr) != 0 {
			column.Generation = &metadatapb.GenerationMetadata{
				Type:       metadatapb.GenerationMetadata_TYPE_STORED,
				Expression: text,
			}
		}
	}
	if err := columnRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}

	checkMap := make(map[string][]*metadatapb.CheckConstraintMetadata)
	if atLeast8016 {
		checkQuery := getCheckConstraintQuery(isMariaDB)
		checkRows, err := d.db.QueryContext(ctx, checkQuery, d.databaseName)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", checkQuery, err)
		}
		defer checkRows.Close()
		for checkRows.Next() {
			check := &metadatapb.CheckConstraintMetadata{}
			var tableName string
			if err := checkRows.Scan(&tableName, &check.Name, &check.Expression); err != nil {
				return nil, err
			}
			check.Expression = unescapeCheckClause(check.Expression)
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
			VIEW_DEFINITION
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
		if err := viewRows.Scan(&view.Name, &view.Definition); err != nil {
			return nil, err
		}
		view.Columns = columnMap[view.Name]
		viewMap[view.Name] = view
		viewNames = append(viewNames, view.Name)
	}
	if err := viewRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}
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

	partitionTables, err := d.listPartitionTables(ctx, d.databaseName)
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
				Partitions:       partitionTables[tableName],
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
		}
	}
	if err := tableRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
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

func (d *Driver) getEventList(ctx context.Context, databaseName string) ([]*metadatapb.EventMetadata, error) {
	const listEventsQuery = `
	SELECT
		EVENT_NAME,
		TIME_ZONE,
		SQL_MODE,
		CHARACTER_SET_CLIENT,
		COLLATION_CONNECTION
	FROM INFORMATION_SCHEMA.EVENTS
	WHERE EVENT_SCHEMA = ?
	ORDER BY EVENT_NAME ASC;
	`
	rows, err := d.db.QueryContext(ctx, listEventsQuery, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", listEventsQuery, err)
	}
	defer rows.Close()

	type eventRow struct {
		name, timeZone, sqlMode, charsetClient, collationConnection string
	}
	var rawEvents []eventRow
	for rows.Next() {
		var r eventRow
		if err := rows.Scan(&r.name, &r.timeZone, &r.sqlMode, &r.charsetClient, &r.collationConnection); err != nil {
			return nil, err
		}
		rawEvents = append(rawEvents, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", listEventsQuery, err)
	}
	_ = rows.Close()

	var events []*metadatapb.EventMetadata
	for _, r := range rawEvents {
		def, err := d.showCreateColumn(ctx, fmt.Sprintf("SHOW CREATE EVENT %s.%s", quoteBacktick(databaseName), quoteBacktick(r.name)), "Create Event")
		if err != nil {
			return nil, err
		}
		events = append(events, &metadatapb.EventMetadata{
			Name:                r.name,
			TimeZone:            r.timeZone,
			Definition:          def,
			SqlMode:             r.sqlMode,
			CharacterSetClient:  r.charsetClient,
			CollationConnection: r.collationConnection,
		})
	}
	return events, nil
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
			Body:                statement,
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
	// Close routineRows before running SHOW CREATE on a 1-connection pool (tape).
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

func (d *Driver) listPartitionTables(ctx context.Context, databaseName string) (map[string][]*metadatapb.TablePartitionMetadata, error) {
	const query = `
		SELECT
			TABLE_NAME,
			PARTITION_NAME,
			SUBPARTITION_NAME,
			PARTITION_METHOD,
			SUBPARTITION_METHOD,
			PARTITION_EXPRESSION,
			SUBPARTITION_EXPRESSION,
			PARTITION_DESCRIPTION
		FROM INFORMATION_SCHEMA.PARTITIONS
		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
	`
	rows, err := d.db.QueryContext(ctx, query, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	type partitionKey struct {
		tableName     string
		partitionName string
	}
	partitionMap := make(map[partitionKey]int)
	result := make(map[string][]*metadatapb.TablePartitionMetadata)

	for rows.Next() {
		var tableName, partitionName, partitionMethod string
		var subpartitionName, subpartitionMethod, subpartitionExpression, partitionExpression, partitionDescription sql.NullString
		if err := rows.Scan(
			&tableName,
			&partitionName,
			&subpartitionName,
			&partitionMethod,
			&subpartitionMethod,
			&partitionExpression,
			&subpartitionExpression,
			&partitionDescription,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		pk := partitionKey{tableName: tableName, partitionName: partitionName}
		if _, ok := partitionMap[pk]; !ok {
			tp := convertPartitionType(partitionMethod)
			if tp == metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
				continue
			}
			partition := &metadatapb.TablePartitionMetadata{
				Name:          partitionName,
				Type:          tp,
				Expression:    partitionExpression.String,
				Value:         partitionDescription.String,
				Subpartitions: []*metadatapb.TablePartitionMetadata{},
			}
			partitionMap[pk] = len(result[tableName])
			result[tableName] = append(result[tableName], partition)
		}
		if subpartitionName.Valid {
			tp := convertPartitionType(subpartitionMethod.String)
			if tp == metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
				continue
			}
			sub := &metadatapb.TablePartitionMetadata{
				Name:          subpartitionName.String,
				Type:          tp,
				Expression:    subpartitionExpression.String,
				Subpartitions: []*metadatapb.TablePartitionMetadata{},
			}
			if idx, ok := partitionMap[pk]; ok {
				result[tableName][idx].Subpartitions = append(result[tableName][idx].Subpartitions, sub)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}
	_ = rows.Close()

	tableNames := make([]string, 0, len(result))
	for tbl, parts := range result {
		if len(parts) > 0 {
			tableNames = append(tableNames, tbl)
		}
	}
	slices.Sort(tableNames)

	for _, tbl := range tableNames {
		partitions := result[tbl]
		showQuery := fmt.Sprintf("SHOW CREATE TABLE %s.%s", quoteBacktick(databaseName), quoteBacktick(tbl))
		showRows, err := d.db.QueryContext(ctx, showQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", showQuery, err)
		}
		for showRows.Next() {
			var name, createTable string
			if err := showRows.Scan(&name, &createTable); err != nil {
				_ = showRows.Close()
				return nil, fmt.Errorf("failed to scan row for %q: %w", showQuery, err)
			}
			partitionNum, subPartitionNum := parsePartitionCounts(createTable)
			for _, p := range partitions {
				if partitionNum != 0 {
					p.UseDefault = strconv.Itoa(partitionNum)
				}
				for _, sp := range p.Subpartitions {
					if subPartitionNum != 0 {
						sp.UseDefault = strconv.Itoa(subPartitionNum)
					}
				}
			}
		}
		if err := showRows.Err(); err != nil {
			_ = showRows.Close()
			return nil, fmt.Errorf("failed to iterate %q: %w", showQuery, err)
		}
		_ = showRows.Close()
	}
	return result, nil
}

// parsePartitionCounts extracts explicit PARTITIONS N and SUBPARTITIONS N
// counts from a SHOW CREATE TABLE definition the server already accepted; a
// definition the parser cannot read leaves both counts unset.
func parsePartitionCounts(createTable string) (partitions, subpartitions int) {
	list, err := mysqlparser.Parse(createTable)
	if err != nil || list == nil || len(list.Items) != 1 {
		return 0, 0
	}
	stmt, ok := list.Items[0].(*mysqlast.CreateTableStmt)
	if !ok || stmt.Partitions == nil {
		return 0, 0
	}
	return stmt.Partitions.NumParts, stmt.Partitions.NumSubParts
}

func convertPartitionType(tp string) metadatapb.TablePartitionMetadata_Type {
	switch strings.ToUpper(tp) {
	case "RANGE":
		return metadatapb.TablePartitionMetadata_RANGE
	case "RANGE COLUMNS":
		return metadatapb.TablePartitionMetadata_RANGE_COLUMNS
	case "LIST":
		return metadatapb.TablePartitionMetadata_LIST
	case "LIST COLUMNS":
		return metadatapb.TablePartitionMetadata_LIST_COLUMNS
	case "HASH":
		return metadatapb.TablePartitionMetadata_HASH
	case "KEY":
		return metadatapb.TablePartitionMetadata_KEY
	case "LINEAR HASH":
		return metadatapb.TablePartitionMetadata_LINEAR_HASH
	case "LINEAR KEY":
		return metadatapb.TablePartitionMetadata_LINEAR_KEY
	default:
		return metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED
	}
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

func getCheckConstraintQuery(isMariaDB bool) string {
	if isMariaDB {
		return `
			SELECT
				TABLE_NAME,
				CONSTRAINT_NAME,
				CHECK_CLAUSE
			FROM information_schema.CHECK_CONSTRAINTS
			WHERE CONSTRAINT_SCHEMA = ?
			ORDER BY TABLE_NAME, CONSTRAINT_NAME
		`
	}
	return `
			SELECT
				tc.TABLE_NAME,
				cc.CONSTRAINT_NAME,
				cc.CHECK_CLAUSE
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

func setColumnMetadataDefault(column *metadatapb.ColumnMetadata, defaultStr sql.NullString, nullableBool bool, extra string, binaryFormat binaryDefaultFormat) {
	if defaultStr.Valid {
		unquotedDefault := unquoteMySQLString(defaultStr.String)
		isBinaryDefault := binaryFormat != binaryDefaultVerbatim && isBinaryFamilyColumnType(column.Type)
		switch {
		case isTemporalColumnType(column.Type) && temporalKeywordDefaultRe.MatchString(unquotedDefault):
			column.Default = unquotedDefault
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

	if strings.Contains(extra, "on update CURRENT_TIMESTAMP") {
		if match := currentTSPrecision.FindStringSubmatch(extra); len(match) > 1 {
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
