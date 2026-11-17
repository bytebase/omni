package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

const (
	baseTableType          = "BASE TABLE"
	viewTableType          = "VIEW"
	autoIncrement          = "AUTO_INCREMENT"
	autoRandSymbol         = "AUTO_RANDOM"
	pkAutoRandomBitsSymbol = "PK_AUTO_RANDOM_BITS"
	virtualGenToken        = "VIRTUAL GENERATED"
	storedGenToken         = "STORED GENERATED"
	erBadDBError           = 1049
)

var (
	systemDatabases = []string{
		"information_schema",
		"metrics_schema",
		"mysql",
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

	pkAutoRandomBitsRegex    = regexp.MustCompile(`PK_AUTO_RANDOM_BITS=(\d+)`)
	rangeBitsRegex           = regexp.MustCompile(`RANGE BITS=(\d+)`)
	tidbSemverRegex          = regexp.MustCompile(`v\d+\.\d+\.\d+`)
	mysqlCompatVersionRegex  = regexp.MustCompile(`^\d+\.\d+\.\d+`)
	temporalKeywordDefaultRe = regexp.MustCompile(`^(?i)CURRENT_(?:TIMESTAMP|DATE|TIME)(?:\(\d+\))?$`)
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

// SyncInstance reads the verbatim TiDB version string, lower_case_table_names,
// user accounts with their global grants, and non-system databases.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
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

	version, err := d.getVersion(ctx)
	if err != nil {
		if d.engineVersion == "" {
			return nil, err
		}
		version = d.engineVersion
	}

	indexMap := make(map[string]map[string]*metadatapb.IndexMetadata)
	const indexQuery = `
		SELECT
			TABLE_NAME,
			KEY_NAME,
			COLUMN_NAME,
			SEQ_IN_INDEX,
			NON_UNIQUE,
			INDEX_COMMENT,
			SUB_PART,
			EXPRESSION
		FROM information_schema.TIDB_INDEXES
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME, KEY_NAME, SEQ_IN_INDEX`
	indexRows, err := d.db.QueryContext(ctx, indexQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", indexQuery, err)
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var tableName, keyName, comment string
		var columnName, expressionName sql.NullString
		var position int
		var nonUnique bool
		var subPart sql.NullInt64
		if err := indexRows.Scan(
			&tableName,
			&keyName,
			&columnName,
			&position,
			&nonUnique,
			&comment,
			&subPart,
			&expressionName,
		); err != nil {
			return nil, err
		}

		var expression string
		if expressionName.Valid {
			expression = fmt.Sprintf("(%s)", expressionName.String)
		} else if columnName.Valid {
			expression = columnName.String
		}

		if _, ok := indexMap[tableName]; !ok {
			indexMap[tableName] = make(map[string]*metadatapb.IndexMetadata)
		}
		if _, ok := indexMap[tableName][keyName]; !ok {
			indexMap[tableName][keyName] = &metadatapb.IndexMetadata{
				Name:    keyName,
				Type:    "BTREE",
				Unique:  !nonUnique,
				Primary: keyName == "PRIMARY",
				Visible: true,
				Comment: comment,
			}
		}
		idx := indexMap[tableName][keyName]
		idx.Expressions = append(idx.Expressions, expression)
		if subPart.Valid {
			idx.KeyLength = append(idx.KeyLength, subPart.Int64)
		} else {
			idx.KeyLength = append(idx.KeyLength, -1)
		}
	}
	if err := indexRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", indexQuery, err)
	}

	columnMap := make(map[string][]*metadatapb.ColumnMetadata)
	const columnQuery = `
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
			GENERATION_EXPRESSION,
			EXTRA
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME, ORDINAL_POSITION`
	columnRows, err := d.db.QueryContext(ctx, columnQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}
	defer columnRows.Close()

	for columnRows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var tableName, nullable, extra string
		var defaultStr, generationExpr sql.NullString
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
			&generationExpr,
			&extra,
		); err != nil {
			return nil, err
		}
		column.Comment = stripSingleQuote(column.Comment)
		nullableBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, err
		}
		column.Nullable = nullableBool
		setColumnMetadataDefault(column, defaultStr, nullableBool, extra)

		if generationExpr.Valid && generationExpr.String != "" {
			upperExtra := strings.ToUpper(extra)
			if strings.Contains(upperExtra, virtualGenToken) {
				column.Generation = &metadatapb.GenerationMetadata{
					Type:       metadatapb.GenerationMetadata_TYPE_VIRTUAL,
					Expression: generationExpr.String,
				}
			} else if strings.Contains(upperExtra, storedGenToken) {
				column.Generation = &metadatapb.GenerationMetadata{
					Type:       metadatapb.GenerationMetadata_TYPE_STORED,
					Expression: generationExpr.String,
				}
			}
		}
		columnMap[tableName] = append(columnMap[tableName], column)
	}
	if err := columnRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", columnQuery, err)
	}

	viewMap := make(map[string]*metadatapb.ViewMetadata)
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
	}
	if err := viewRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", viewQuery, err)
	}

	foreignKeysMap, err := d.getForeignKeyList(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}

	var checkMap map[string][]*metadatapb.CheckConstraintMetadata
	if supportsCheckConstraintTable(version) {
		checkMap, err = d.getCheckConstraintList(ctx, d.databaseName)
		if err != nil {
			return nil, err
		}
	}

	sequences, err := d.getSequenceList(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}

	partitionTables, err := d.listPartitionTables(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}

	const tableQuery = `
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
			QUOTE(IFNULL(TABLE_COMMENT, '')),
			IFNULL(TIDB_ROW_ID_SHARDING_INFO, ''),
			IFNULL(TIDB_PK_TYPE, '')
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ?
		ORDER BY TABLE_NAME`
	tableRows, err := d.db.QueryContext(ctx, tableQuery, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", tableQuery, err)
	}
	defer tableRows.Close()

	for tableRows.Next() {
		var tableName, tableType, engine, collation, createOptions, comment, shardingInfo, pkType string
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
			&shardingInfo,
			&pkType,
		); err != nil {
			return nil, err
		}
		comment = stripSingleQuote(comment)

		switch tableType {
		case baseTableType:
			columns := columnMap[tableName]
			if strings.Contains(shardingInfo, pkAutoRandomBitsSymbol) {
				autoRandText := autoRandSymbol
				if randomBitsMatch := pkAutoRandomBitsRegex.FindStringSubmatch(shardingInfo); len(randomBitsMatch) > 1 {
					if rangeBitsMatch := rangeBitsRegex.FindStringSubmatch(shardingInfo); len(rangeBitsMatch) > 1 {
						autoRandText += fmt.Sprintf("(%s, %s)", randomBitsMatch[1], rangeBitsMatch[1])
					} else {
						autoRandText += fmt.Sprintf("(%s)", randomBitsMatch[1])
					}
				}
				if indexes, ok := indexMap[tableName]; ok {
					for _, index := range indexes {
						if index.Primary && len(index.Expressions) > 0 {
							colName := index.Expressions[0]
							for i, col := range columns {
								if col.Name == colName {
									columns[i].Default = autoRandText
									break
								}
							}
						}
					}
				}
			}

			tableMetadata := &metadatapb.TableMetadata{
				Name:             tableName,
				Columns:          columns,
				ForeignKeys:      foreignKeysMap[tableName],
				Engine:           engine,
				Collation:        collation,
				RowCount:         rowCount,
				DataSize:         dataSize,
				IndexSize:        indexSize,
				DataFree:         dataFree,
				CreateOptions:    createOptions,
				Comment:          comment,
				Charset:          convertCollationToCharset(collation),
				Partitions:       partitionTables[tableName],
				CheckConstraints: checkMap[tableName],
				ShardingInfo:     shardingInfo,
				PrimaryKeyType:   pkType,
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
	schemaMetadata.Sequences = sequences

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

func convertCollationToCharset(collation string) string {
	tokens := strings.Split(collation, "_")
	return tokens[0]
}

func isTimeConstant(s string) bool {
	if s == "0000-00-00 00:00:00" {
		return true
	}
	_, err := time.Parse(time.DateTime, s)
	return err == nil
}

func setColumnMetadataDefault(column *metadatapb.ColumnMetadata, defaultStr sql.NullString, nullableBool bool, extra string) {
	if defaultStr.Valid {
		unquotedDefault := unquoteMySQLString(defaultStr.String)
		switch {
		case isTemporalColumnType(column.Type) && temporalKeywordDefaultRe.MatchString(unquotedDefault):
			column.Default = unquotedDefault
		case strings.Contains(extra, "DEFAULT_GENERATED"):
			if isTimeConstant(defaultStr.String) {
				column.Default = defaultStr.String
			} else {
				column.Default = fmt.Sprintf("(%s)", unescapeExpressionDefault(unquotedDefault))
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
			PARTITION_DESCRIPTION,
			PARTITION_COMMENT
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
		var subpartitionName, subpartitionMethod, subpartitionExpression, partitionExpression, partitionDescription, partitionComment sql.NullString
		if err := rows.Scan(
			&tableName,
			&partitionName,
			&subpartitionName,
			&partitionMethod,
			&subpartitionMethod,
			&partitionExpression,
			&subpartitionExpression,
			&partitionDescription,
			&partitionComment,
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
				Comment:       partitionComment.String,
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
				Comment:       partitionComment.String,
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
	return result, nil
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
			kcu.TABLE_SCHEMA,
			fks.TABLE_NAME,
			fks.CONSTRAINT_NAME,
			kcu.COLUMN_NAME,
			kcu.REFERENCED_TABLE_SCHEMA,
			fks.REFERENCED_TABLE_NAME,
			kcu.REFERENCED_COLUMN_NAME,
			fks.DELETE_RULE,
			fks.UPDATE_RULE,
			fks.MATCH_OPTION
		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
				AND fks.TABLE_NAME = kcu.TABLE_NAME
				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION;
	`
	fkRows, err := d.db.QueryContext(ctx, fkQuery, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", fkQuery, err)
	}
	defer fkRows.Close()

	foreignKeysMap := make(map[string][]*metadatapb.ForeignKeyMetadata)
	var buildingFk *metadatapb.ForeignKeyMetadata
	var buildingTable string
	for fkRows.Next() {
		var tableSchema, tableName, referencedSchema, column, referencedColumn string
		var fk metadatapb.ForeignKeyMetadata
		if err := fkRows.Scan(
			&tableSchema,
			&tableName,
			&fk.Name,
			&column,
			&referencedSchema,
			&fk.ReferencedTable,
			&referencedColumn,
			&fk.OnDelete,
			&fk.OnUpdate,
			&fk.MatchType,
		); err != nil {
			return nil, err
		}
		if referencedSchema != tableSchema {
			fk.ReferencedSchema = referencedSchema
		}
		fk.Columns = append(fk.Columns, column)
		fk.ReferencedColumns = append(fk.ReferencedColumns, referencedColumn)
		if buildingFk == nil {
			buildingTable = tableName
			buildingFk = &fk
		} else if tableName == buildingTable && buildingFk.Name == fk.Name {
			buildingFk.Columns = append(buildingFk.Columns, fk.Columns[0])
			buildingFk.ReferencedColumns = append(buildingFk.ReferencedColumns, fk.ReferencedColumns[0])
		} else {
			foreignKeysMap[buildingTable] = append(foreignKeysMap[buildingTable], buildingFk)
			buildingTable = tableName
			buildingFk = &fk
		}
	}
	if err := fkRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", fkQuery, err)
	}
	if buildingFk != nil {
		foreignKeysMap[buildingTable] = append(foreignKeysMap[buildingTable], buildingFk)
	}
	return foreignKeysMap, nil
}

func (d *Driver) getCheckConstraintList(ctx context.Context, databaseName string) (map[string][]*metadatapb.CheckConstraintMetadata, error) {
	const checkQuery = `
		SELECT
			tc.TABLE_NAME,
			cc.CONSTRAINT_NAME,
			cc.CHECK_CLAUSE
		FROM information_schema.CHECK_CONSTRAINTS cc
			JOIN information_schema.TABLE_CONSTRAINTS tc
				ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
				AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
		WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ? AND cc.CONSTRAINT_SCHEMA = ?
		ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
	`
	checkRows, err := d.db.QueryContext(ctx, checkQuery, databaseName, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", checkQuery, err)
	}
	defer checkRows.Close()

	checkMap := make(map[string][]*metadatapb.CheckConstraintMetadata)
	for checkRows.Next() {
		check := &metadatapb.CheckConstraintMetadata{}
		var tableName string
		if err := checkRows.Scan(&tableName, &check.Name, &check.Expression); err != nil {
			return nil, err
		}
		checkMap[tableName] = append(checkMap[tableName], check)
	}
	if err := checkRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", checkQuery, err)
	}
	return checkMap, nil
}

func supportsCheckConstraintTable(rawVersion string) bool {
	if loc := tidbSemverRegex.FindStringIndex(rawVersion); loc != nil {
		if sv, err := parseSemver(rawVersion[loc[0]+1 : loc[1]]); err == nil {
			return sv.atLeast(7, 4, 0)
		}
	}
	loc := mysqlCompatVersionRegex.FindStringIndex(rawVersion)
	if loc == nil {
		return false
	}
	sv, err := parseSemver(rawVersion[loc[0]:loc[1]])
	if err != nil {
		return false
	}
	return sv.atLeast(8, 0, 0)
}

func (d *Driver) getSequenceList(ctx context.Context, databaseName string) ([]*metadatapb.SequenceMetadata, error) {
	const query = `
		SELECT
			SEQUENCE_NAME,
			START,
			MIN_VALUE,
			MAX_VALUE,
			INCREMENT,
			CYCLE,
			CACHE_VALUE,
			IFNULL(COMMENT, '')
		FROM information_schema.SEQUENCES
		WHERE SEQUENCE_SCHEMA = ?
		ORDER BY SEQUENCE_NAME`
	rows, err := d.db.QueryContext(ctx, query, databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer rows.Close()

	var sequences []*metadatapb.SequenceMetadata
	for rows.Next() {
		seq := &metadatapb.SequenceMetadata{DataType: "BIGINT"}
		var cycleOption int64
		if err := rows.Scan(
			&seq.Name,
			&seq.Start,
			&seq.MinValue,
			&seq.MaxValue,
			&seq.Increment,
			&cycleOption,
			&seq.CacheSize,
			&seq.Comment,
		); err != nil {
			return nil, err
		}
		seq.Cycle = cycleOption != 0
		sequences = append(sequences, seq)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return sequences, nil
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

func unescapeExpressionDefault(s string) string {
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
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
