package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	omnidriver "github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	oracleanalysis "github.com/bytebase/omni/oracle/analysis"
)

// Keys that name one object of a schema.
type tableKey struct {
	Schema string
	Table  string
}

type columnKey struct {
	Schema string
	Table  string
	Column string
}

type constraintRef struct {
	Schema     string
	Constraint string
}

type indexKey struct {
	Schema string
	Table  string
	Index  string
}

// parseVersion reads the release out of the version string, such as "23.26.3.0.0".
func parseVersion(banner string) (*oracleanalysis.Version, error) {
	re := regexp.MustCompile(`(\d+)\.(\d+)`)
	match := re.FindStringSubmatch(banner)
	if len(match) < 3 {
		return nil, fmt.Errorf("failed to parse version from banner: %s", banner)
	}
	first, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse first version from banner %s: %w", banner, err)
	}
	second, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse second version from banner %s: %w", banner, err)
	}
	return &oracleanalysis.Version{First: first, Second: second}, nil
}

// convertYesNo reads the YES/NO spelling Oracle reports for a nullable column.
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

// sanitizeUTF8Message replaces invalid UTF-8 byte sequences in every populated
// string field of m, recursing into nested messages, repeated fields, and maps
// (both keys and values). The driver can hand back raw non-UTF-8 bytes for
// dictionary strings on a non-UTF-8 character set, and the metadata type refuses
// to marshal those.
func sanitizeUTF8Message(m proto.Message) {
	if m == nil {
		return
	}
	sanitizeUTF8Value(m.ProtoReflect())
}

func sanitizeUTF8Value(m protoreflect.Message) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsMap():
			sanitizeUTF8Map(fd, v.Map())
		case fd.IsList():
			list := v.List()
			switch fd.Kind() {
			case protoreflect.StringKind:
				for i := 0; i < list.Len(); i++ {
					s := list.Get(i).String()
					if !utf8.ValidString(s) {
						list.Set(i, protoreflect.ValueOfString(sanitizeUTF8String(s)))
					}
				}
			case protoreflect.MessageKind, protoreflect.GroupKind:
				for i := 0; i < list.Len(); i++ {
					sanitizeUTF8Value(list.Get(i).Message())
				}
			default:
			}
		case fd.Kind() == protoreflect.StringKind:
			if s := v.String(); !utf8.ValidString(s) {
				m.Set(fd, protoreflect.ValueOfString(sanitizeUTF8String(s)))
			}
		case fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind:
			sanitizeUTF8Value(v.Message())
		default:
		}
		return true
	})
}

func sanitizeUTF8Map(fd protoreflect.FieldDescriptor, mp protoreflect.Map) {
	valueKind := fd.MapValue().Kind()
	keyIsString := fd.MapKey().Kind() == protoreflect.StringKind

	type rekey struct {
		oldKey protoreflect.MapKey
		newKey protoreflect.MapKey
		value  protoreflect.Value
	}
	var rekeys []rekey
	mp.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		switch valueKind {
		case protoreflect.StringKind:
			if s := v.String(); !utf8.ValidString(s) {
				mp.Set(k, protoreflect.ValueOfString(sanitizeUTF8String(s)))
			}
		case protoreflect.MessageKind, protoreflect.GroupKind:
			sanitizeUTF8Value(v.Message())
		default:
		}
		if keyIsString {
			if s := k.String(); !utf8.ValidString(s) {
				rekeys = append(rekeys, rekey{
					oldKey: k,
					newKey: protoreflect.ValueOfString(sanitizeUTF8String(s)).MapKey(),
					value:  mp.Get(k),
				})
			}
		}
		return true
	})
	for _, r := range rekeys {
		mp.Clear(r.oldKey)
		mp.Set(r.newKey, r.value)
	}
}

// sanitizeUTF8String drops bytes that are not valid UTF-8, which the metadata
// type refuses.
func sanitizeUTF8String(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	var b strings.Builder
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			s = s[1:]
			continue
		}
		b.WriteString(s[:size])
		s = s[size:]
	}
	return b.String()
}

const systemSchema = "'ANONYMOUS','APPQOSSYS','AUDSYS','CTXSYS','DBSFWUSER','DBSNMP','DGPDB_INT','DIP','DVF','DVSYS','GGSYS','GSMADMIN_INTERNAL','GSMCATUSER','GSMROOTUSER','GSMUSER','LBACSYS','MDDATA','MDSYS','OPS$ORACLE','ORACLE_OCM','OUTLN','REMOTE_SCHEDULER_AGENT','SYS','SYS$UMF','SYSBACKUP','SYSDG','SYSKM','SYSRAC','SYSTEM','XDB','XS$NULL','XS$$NULL','FLOWS_FILES','HR','MDSYS','EXFSYS','MGMT_VIEW','OLAPSYS','ORDDATA','ORDPLUGINS','ORDSYS','OWBSYS','OWBSYS_AUDIT','SCOTT','SI_INFORMTN_SCHEMA','SPATIAL_CSW_ADMIN_USR','SPATIAL_WFS_ADMIN_USR','SYSMAN','WMSYS','OJVMSYS'"

var (
	semVersionRegex       = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)
	canonicalVersionRegex = regexp.MustCompile(`[0-9][0-9][a-z]`)
)

// SyncInstance syncs the instance.
func (d *Driver) SyncInstance(ctx context.Context) (*omnidriver.Instance, error) {
	var fullVersion string
	queryVersion := "SELECT BANNER FROM v$version WHERE banner LIKE 'Oracle%'"
	if err := d.db.QueryRowContext(ctx, queryVersion).Scan(&fullVersion); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryVersion, err)
	}
	tokens := strings.Fields(fullVersion)
	var version, canonicalVersion string
	for _, token := range tokens {
		if semVersionRegex.MatchString(token) {
			version = token
			continue
		}
		if canonicalVersionRegex.MatchString(token) {
			canonicalVersion = token
			continue
		}
	}
	if canonicalVersion != "" {
		version = fmt.Sprintf("%s (%s)", version, canonicalVersion)
	}

	txn, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = txn.Rollback() }()

	schemas, err := getSchemas(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas from database %q: %w", d.databaseName, err)
	}
	var databases []*metadatapb.DatabaseSchemaMetadata
	for _, schema := range schemas {
		databases = append(databases, &metadatapb.DatabaseSchemaMetadata{
			// The driver can hand back raw non-UTF-8 bytes for dictionary
			// strings on non-UTF-8 charsets (e.g. ZHS16GBK values ending in a
			// truncated multi-byte character); proto string fields reject them
			// at marshal time, so sanitize at the source.
			Name:        sanitizeUTF8String(schema),
			ServiceName: "",
		})
	}

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	return &omnidriver.Instance{
		Version:   version,
		Databases: databases,
	}, nil
}

// schema returns the schema the driver reads — Config.Database, or the account's
// own when the caller named none — resolving the latter once, since a driver is
// shared and its sessions start in that schema.
func (d *Driver) schema(ctx context.Context) (string, error) {
	d.schemaMu.Lock()
	defer d.schemaMu.Unlock()
	if d.schemaName != "" {
		return d.schemaName, nil
	}
	name := d.databaseName
	if name == "" {
		// A failed lookup is not remembered: the caller may reach the server on
		// the next call.
		if err := d.db.QueryRowContext(ctx, "SELECT SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA') FROM DUAL").Scan(&name); err != nil {
			return "", fmt.Errorf("failed to read the current schema: %w", err)
		}
	}
	d.schemaName = name
	return name, nil
}

// SyncDatabase reads the full schema metadata of Config.Database, or of the
// account's own schema when Config names none.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	schemaName, err := d.schema(ctx)
	if err != nil {
		return nil, err
	}
	// The version comes from the server on a connection of its own: the
	// dictionary reads below hold the transaction's.
	version, err := d.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	txn, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = txn.Rollback() }()

	identityColumns, err := getIdentityColumns(txn, schemaName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity columns from database %q: %w", schemaName, err)
	}
	columnMap, internalPos, err := getTableColumns(txn, schemaName, version, identityColumns)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns from database %q: %w", schemaName, err)
	}
	tableTriggerMap, viewTriggerMap, err := getTriggers(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get triggers from database %q: %w", schemaName, err)
	}
	indexMap, otherIndexes, checkConstraintMap, foreignKeyMap, err := getIndexesAndConstraints(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices: %w", err)
	}
	partitionMap, err := getTablePartitions(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table partitions from database %q: %w", schemaName, err)
	}
	tableMap, err := getTables(txn, schemaName, columnMap, indexMap, checkConstraintMap, foreignKeyMap, tableTriggerMap, partitionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables from database %q: %w", schemaName, err)
	}
	viewMap, err := getViews(txn, schemaName, columnMap, internalPos, viewTriggerMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get views from database %q: %w", schemaName, err)
	}
	sequences, err := getSequences(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequences from database %q: %w", schemaName, err)
	}
	dbLinks, err := getDBLinks(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to get db links from database %q: %w", schemaName, err)
	}
	functions, procedures, packages, err := getRoutines(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get routines from database %q: %w", schemaName, err)
	}
	materializedViews, err := getMaterializedViews(txn, schemaName, columnMap, indexMap, tableTriggerMap, partitionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get materialized views from database %q: %w", schemaName, err)
	}
	synonyms, err := getSynonyms(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get synonyms from database %q: %w", schemaName, err)
	}
	publicSynonyms, err := getSynonyms(txn, publicSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to get public synonyms: %w", err)
	}

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	databaseMetadata := &metadatapb.DatabaseSchemaMetadata{
		Name:            schemaName,
		ServiceName:     d.serviceName,
		LinkedDatabases: dbLinks,
		PublicSynonyms:  publicSynonyms,
	}
	databaseMetadata.Schemas = append(databaseMetadata.Schemas, &metadatapb.SchemaMetadata{
		Name:              "",
		Tables:            tableMap[schemaName],
		Views:             viewMap[schemaName],
		MaterializedViews: materializedViews,
		Sequences:         sequences,
		Functions:         functions,
		Procedures:        procedures,
		Packages:          packages,
		Synonyms:          synonyms,
		OtherIndexes:      otherIndexes,
	})
	// Deep-sanitize every string field (names, definitions, defaults, …), not
	// just comments: the go-ora driver returns wholly-unconverted raw bytes
	// for any GBK/SJIS/Big5 value that ends in a dangling multi-byte lead
	// byte, and one invalid string anywhere fails proto marshaling for the
	// entire database metadata (BYT-9916).
	sanitizeUTF8Message(databaseMetadata)
	return databaseMetadata, nil
}

// publicSchema owns the public synonyms. No user can take its name (ORA-01920).
const publicSchema = "PUBLIC"

func getSynonyms(txn *sql.Tx, owner string) ([]*metadatapb.SynonymMetadata, error) {
	query := `
		SELECT SYNONYM_NAME, TABLE_OWNER, TABLE_NAME, DB_LINK
		FROM ALL_SYNONYMS
		WHERE OWNER = :1
		ORDER BY SYNONYM_NAME`
	rows, err := txn.Query(query, owner)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	var synonyms []*metadatapb.SynonymMetadata
	for rows.Next() {
		var name string
		var targetSchema, targetName, targetDBLink sql.NullString
		if err := rows.Scan(&name, &targetSchema, &targetName, &targetDBLink); err != nil {
			return nil, err
		}
		synonyms = append(synonyms, &metadatapb.SynonymMetadata{
			Name:         name,
			TargetSchema: targetSchema.String,
			TargetName:   targetName.String,
			TargetDbLink: targetDBLink.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return synonyms, nil
}

func getDBLinks(txn *sql.Tx) ([]*metadatapb.LinkedDatabaseMetadata, error) {
	query := `
	SELECT DB_LINK, HOST, USERNAME
	FROM all_db_links
	ORDER BY DB_LINK`
	slog.Debug("running get db link query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	var result []*metadatapb.LinkedDatabaseMetadata
	for rows.Next() {
		dbLink := &metadatapb.LinkedDatabaseMetadata{}
		var name, host, username sql.NullString
		if err := rows.Scan(&name, &host, &username); err != nil {
			return nil, err
		}
		if !name.Valid {
			continue
		}
		dbLink.Name = name.String
		if host.Valid {
			dbLink.Host = host.String
		}
		if username.Valid {
			dbLink.Username = username.String
		}
		result = append(result, dbLink)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return result, nil
}

func getSchemas(txn *sql.Tx) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT username FROM all_users
		WHERE username NOT IN (%s) AND username NOT LIKE 'APEX_%%' ORDER BY username`,
		systemSchema)
	slog.Debug("running get schemas query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	var result []string
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		result = append(result, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return result, nil
}

func getTriggers(txn *sql.Tx, schemaName string) (map[tableKey][]*metadatapb.TriggerMetadata, map[tableKey][]*metadatapb.TriggerMetadata, error) {
	tableTriggerMap := make(map[tableKey][]*metadatapb.TriggerMetadata)
	viewTriggerMap := make(map[tableKey][]*metadatapb.TriggerMetadata)

	// A trigger the schema owns on another schema's table is left out: keyed by
	// TABLE_NAME alone, it would land on this schema's table of the same name.
	query := fmt.Sprintf(`
		SELECT TRIGGER_NAME, TABLE_NAME, DESCRIPTION, TRIGGER_BODY, BASE_OBJECT_TYPE, TRIGGER_TYPE, TRIGGERING_EVENT
		FROM sys.all_triggers
		WHERE OWNER = '%s' AND TABLE_OWNER = '%s' AND TABLE_NAME is NOT NULL AND (BASE_OBJECT_TYPE = 'TABLE' OR BASE_OBJECT_TYPE = 'VIEW')
		ORDER BY TRIGGER_NAME`, sqlLiteral(schemaName), sqlLiteral(schemaName))

	slog.Debug("running get triggers query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var triggerName, tableName, description, triggerBody, baseObjectType, triggerType, triggeringEvent sql.NullString
		if err := rows.Scan(&triggerName, &tableName, &description, &triggerBody, &baseObjectType, &triggerType, &triggeringEvent); err != nil {
			return nil, nil, err
		}
		if !triggerName.Valid || !tableName.Valid || !description.Valid || !triggerBody.Valid || !baseObjectType.Valid {
			continue
		}
		key := tableKey{Schema: schemaName, Table: tableName.String}
		trigger := &metadatapb.TriggerMetadata{
			Name: triggerName.String,
			Body: sanitizeOracleMetadataString(
				schemaName,
				"TRIGGER",
				triggerName.String,
				"body",
				constructTriggerBody(description.String, triggerBody.String),
			),
		}
		if triggerType.Valid {
			trigger.Timing = triggerType.String
		}
		if triggeringEvent.Valid {
			trigger.Event = triggeringEvent.String
		}

		if baseObjectType.String == "TABLE" {
			tableTriggerMap[key] = append(tableTriggerMap[key], trigger)
		} else {
			viewTriggerMap[key] = append(viewTriggerMap[key], trigger)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}

	return tableTriggerMap, viewTriggerMap, nil
}

func constructTriggerBody(description, triggerBody string) string {
	var buf strings.Builder
	_, _ = buf.WriteString("CREATE OR REPLACE TRIGGER ")
	_, _ = buf.WriteString(description)
	_, _ = buf.WriteString(triggerBody)
	return buf.String()
}

// getTables gets all tables of a database.
func getTables(
	txn *sql.Tx,
	schemaName string,
	columnMap map[tableKey][]*metadatapb.ColumnMetadata,
	indexMap map[tableKey][]*metadatapb.IndexMetadata,
	checkConstraintMap map[tableKey][]*metadatapb.CheckConstraintMetadata,
	foreignKeyMap map[tableKey][]*metadatapb.ForeignKeyMetadata,
	triggerMap map[tableKey][]*metadatapb.TriggerMetadata,
	partitionMap map[tableKey][]*metadatapb.TablePartitionMetadata,
) (map[string][]*metadatapb.TableMetadata, error) {
	columnCommentMap, err := getTableColumnComments(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table column comments: %w", err)
	}
	for key, columns := range columnMap {
		for _, column := range columns {
			comment, ok := columnCommentMap[columnKey{Schema: key.Schema, Table: key.Table, Column: column.Name}]
			if ok {
				column.Comment = comment
			}
		}
	}
	tableCommentMap, err := getTableComments(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table comments: %w", err)
	}
	tableMap := make(map[string][]*metadatapb.TableMetadata)
	appendTables := func(query string, skipDump bool) error {
		rows, err := txn.Query(query)
		if err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			table := &metadatapb.TableMetadata{SkipDump: skipDump}
			var schemaName string
			// https://github.com/rana/ora/issues/57#issuecomment-179909837
			// NUMBER in Oracle can hold 38 decimal digits, so int64 is not enough with its 19 decimal digits.
			// float64 is a little bit better - not precise enough, but won't overflow.
			var count sql.NullFloat64
			var tablespace sql.NullString
			if err := rows.Scan(&schemaName, &table.Name, &count, &tablespace); err != nil {
				return err
			}
			table.RowCount = int64(count.Float64)
			// Oracle gives no tablespace to a partitioned, index-organized,
			// or external table, so one says the table is a heap table of
			// no partitions, which WALK_THROUGH reads (TestSyncTablespace).
			table.Tablespace = tablespace.String
			key := tableKey{Schema: schemaName, Table: table.Name}
			table.Columns = columnMap[key]
			table.Indexes = indexMap[key]
			table.CheckConstraints = checkConstraintMap[key]
			table.ForeignKeys = foreignKeyMap[key]
			table.Triggers = triggerMap[key]
			table.Partitions = partitionMap[key]
			if comment, ok := tableCommentMap[tableKey{Schema: schemaName, Table: table.Name}]; ok {
				table.Comment = comment
			}

			tableMap[schemaName] = append(tableMap[schemaName], table)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
		return rows.Close()
	}

	query := fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, NUM_ROWS, TABLESPACE_NAME
		FROM all_tables
		WHERE OWNER = '%s'
		AND TABLE_NAME NOT IN (
			SELECT MVIEW_NAME FROM all_mviews WHERE OWNER = '%s'
		)
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName), sqlLiteral(schemaName))
	slog.Debug("running get tables query")
	if err := appendTables(query, false); err != nil {
		return nil, err
	}

	// ALL_TABLES omits object and XMLType tables (CREATE TABLE ... OF), which
	// DROP TABLE drops all the same, so they are listed too, with SkipDump. The
	// ALL_MVIEWS subquery drops an object materialized view's container and keeps
	// the plan on the owner index: without it Oracle full-scans OBJ$, which holds
	// every schema's objects.
	objectTableQuery := fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, NUM_ROWS, TABLESPACE_NAME
		FROM all_object_tables
		WHERE OWNER = '%s'
		AND TABLE_NAME NOT IN (
			SELECT MVIEW_NAME FROM all_mviews WHERE OWNER = '%s'
		)
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName), sqlLiteral(schemaName))
	slog.Debug("running get object tables query")
	if err := appendTables(objectTableQuery, true); err != nil {
		return nil, err
	}

	return tableMap, nil
}

// getTableComments reads the comments of the schema's tables. It is one of the
// four comment lookups the sync makes: Oracle's COMMENT statement names a table
// (ALL_TAB_COMMENTS, the view of a view's comment too), a materialized view
// (ALL_MVIEW_COMMENTS), and a column (ALL_COL_COMMENTS). It names no index,
// sequence, routine, or trigger, and COMMENT ON any of those fails with
// ORA-32594. TestSyncCommentableObjects pins the refusals.
func getTableComments(txn *sql.Tx, schemaName string) (map[tableKey]string, error) {
	tableCommentMap := make(map[tableKey]string)

	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		ORDER BY OWNER, TABLE_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName))
	}
	slog.Debug("running get table comments query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var schemaName, tableName, comment string
		if err := rows.Scan(&schemaName, &tableName, &comment); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		tableCommentMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return tableCommentMap, nil
}

func getTableColumnComments(txn *sql.Tx, schemaName string) (map[columnKey]string, error) {
	columnCommentsMap := make(map[columnKey]string)

	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COLUMN_NAME, COMMENTS
		FROM all_col_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		ORDER BY OWNER, TABLE_NAME, COLUMN_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COLUMN_NAME, COMMENTS
		FROM all_col_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		ORDER BY TABLE_NAME, COLUMN_NAME`, sqlLiteral(schemaName))
	}
	slog.Debug("running get table column comments query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var schemaName, tableName, columnName, comment string
		if err := rows.Scan(&schemaName, &tableName, &columnName, &comment); err != nil {
			return nil, err
		}
		key := columnKey{Schema: schemaName, Table: tableName, Column: columnName}
		columnCommentsMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return columnCommentsMap, nil
}

// identityColumn is what ALL_TAB_IDENTITY_COLS states about an identity column
// beyond its default.
type identityColumn struct {
	generation metadatapb.ColumnMetadata_IdentityGeneration
	seed       int64
	increment  int64
	options    []string
}

// getIdentityColumns reads the identity columns of a schema. Identity arrived in
// Oracle 12.1, and the generation and seed live in a dictionary view of their
// own, so this is a query of its own rather than another column of the one above.
func getIdentityColumns(txn *sql.Tx, schemaName string, version *oracleanalysis.Version) (map[columnKey]identityColumn, error) {
	identityColumns := make(map[columnKey]identityColumn)
	if version.First < 12 {
		return identityColumns, nil
	}
	query := fmt.Sprintf(`
		SELECT TABLE_NAME, COLUMN_NAME, GENERATION_TYPE, IDENTITY_OPTIONS
		FROM all_tab_identity_cols
		WHERE OWNER = '%s'`, sqlLiteral(schemaName))
	rows, err := txn.Query(query)
	if err != nil {
		// The view needs a privilege the account may not hold; the columns then
		// keep the default Oracle wrote for them, as before.
		slog.Warn("failed to list Oracle identity columns", slog.String("schema", schemaName), slog.String("error", err.Error()))
		return identityColumns, nil
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var tableName, columnName, generationType string
		var options sql.NullString
		if err := rows.Scan(&tableName, &columnName, &generationType, &options); err != nil {
			return nil, err
		}
		identity := identityColumn{
			generation: metadatapb.ColumnMetadata_IDENTITY_GENERATION_UNSPECIFIED,
			seed:       1,
			increment:  1,
		}
		switch generationType {
		case "ALWAYS":
			identity.generation = metadatapb.ColumnMetadata_ALWAYS
		case "BY DEFAULT":
			identity.generation = metadatapb.ColumnMetadata_BY_DEFAULT
		}
		if options.Valid {
			identity.seed, identity.increment, identity.options = identityOptions(options.String)
		}
		identityColumns[columnKey{Schema: schemaName, Table: tableName, Column: columnName}] = identity
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return identityColumns, nil
}

const (
	oracleDefaultSeqMaxAsc  = "9999999999999999999999999999"
	oracleDefaultSeqMinAsc  = "1"
	oracleDefaultSeqMaxDesc = "-1"
	oracleDefaultSeqMinDesc = "-999999999999999999999999999"
	oracleDefaultSeqCache   = "20"
)

// identityOptions reads START WITH, INCREMENT BY, and any non-default sequence
// options out of the option list Oracle writes for an identity, such as
// "START WITH: 5, INCREMENT BY: 3, MAX_VALUE: 999…, CYCLE_FLAG: N".
func identityOptions(options string) (seed, increment int64, seqOptions []string) {
	seed, increment = 1, 1
	var maxVal, minVal, cacheSize string
	var cycle, order, keep, scale, extend, session, shard bool
	for _, part := range strings.Split(options, ",") {
		name, value, ok := strings.Cut(part, ":")
		if !ok {
			continue
		}
		key := strings.TrimSpace(name)
		val := strings.TrimSpace(value)
		switch key {
		case "START WITH":
			if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
				seed = parsed
			}
		case "INCREMENT BY":
			if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
				increment = parsed
			}
		case "MAX_VALUE":
			maxVal = val
		case "MIN_VALUE":
			minVal = val
		case "CACHE_SIZE":
			cacheSize = val
		case "CYCLE_FLAG":
			cycle = val == "Y"
		case "ORDER_FLAG":
			order = val == "Y"
		case "KEEP_VALUE":
			keep = val == "Y"
		case "SCALE_FLAG":
			scale = val == "Y"
		case "EXTEND_FLAG":
			extend = val == "Y"
		case "SESSION_FLAG":
			session = val == "Y"
		case "SHARDED_FLAG":
			shard = val == "Y"
		}
	}
	defaultMax := oracleDefaultSeqMaxAsc
	defaultMin := oracleDefaultSeqMinAsc
	if increment < 0 {
		defaultMax = oracleDefaultSeqMaxDesc
		defaultMin = oracleDefaultSeqMinDesc
	}
	if maxVal != "" && maxVal != defaultMax {
		seqOptions = append(seqOptions, "MAXVALUE="+maxVal)
	}
	if minVal != "" && minVal != defaultMin {
		seqOptions = append(seqOptions, "MINVALUE="+minVal)
	}
	if cacheSize != "" && cacheSize != oracleDefaultSeqCache {
		seqOptions = append(seqOptions, "CACHE="+cacheSize)
	}
	if cycle {
		seqOptions = append(seqOptions, "CYCLE=true")
	}
	if order {
		seqOptions = append(seqOptions, "ORDER=true")
	}
	if keep {
		seqOptions = append(seqOptions, "KEEP=true")
	}
	if scale {
		seqOptions = append(seqOptions, "SCALE=true")
		if extend {
			seqOptions = append(seqOptions, "EXTEND=true")
		}
	}
	if session {
		seqOptions = append(seqOptions, "SESSION=true")
	}
	if shard {
		seqOptions = append(seqOptions, "SHARD=true")
	}
	return seed, increment, seqOptions
}

// isGeneratedNotNull reports whether a check constraint's condition is the
// NOT NULL Oracle records for a column, which it names itself.
func isGeneratedNotNull(searchCondition string) bool {
	return strings.HasSuffix(strings.TrimSpace(searchCondition), "IS NOT NULL") && !strings.Contains(searchCondition, " AND ")
}

// sqlLiteral writes s as the body of an SQL string literal, doubling its
// apostrophes: an Oracle schema is a user, and a quoted user name may hold one.
func sqlLiteral(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// getTableColumns gets the columns of a table or view. It also returns each
// column's INTERNAL_COLUMN_ID so views can keep their projection order even
// when an invisible column appears before a visible one.
func getTableColumns(txn *sql.Tx, schemaName string, version *oracleanalysis.Version, identityColumns map[columnKey]identityColumn) (map[tableKey][]*metadatapb.ColumnMetadata, map[columnKey]int32, error) {
	columnsMap := make(map[tableKey][]*metadatapb.ColumnMetadata)
	internalPos := make(map[columnKey]int32)

	// Invisible columns have no COLUMN_ID. They sort after the visible columns
	// and take the positions that follow. ALL_TAB_COLUMNS is ALL_TAB_COLS with
	// USER_GENERATED = 'YES', which leaves out hidden system columns; the newer
	// query reads ALL_TAB_COLS under that filter to get VIRTUAL_COLUMN too.
	query := ""
	var virtualColumns map[columnKey]bool
	// https://docs.oracle.com/en/database/oracle/oracle-database/12.2/refrn/ALL_TAB_COLS.html#GUID-85036F42-140A-406B-BE11-0AC49A00DBA3
	equalOrHigherThan12c2release := version.First > 12 || (version.First == 12 && version.Second >= 2)
	if equalOrHigherThan12c2release {
		query = fmt.Sprintf(`
		SELECT
			OWNER,
			TABLE_NAME,
			COLUMN_NAME,
			DATA_TYPE,
			DATA_LENGTH,
			DATA_PRECISION,
			DATA_SCALE,
			CHAR_USED,
			CHAR_LENGTH,
			COLUMN_ID,
			DATA_DEFAULT,
			NULLABLE,
			COLLATION,
			DEFAULT_ON_NULL,
			VIRTUAL_COLUMN,
			IDENTITY_COLUMN,
			INTERNAL_COLUMN_ID
		FROM sys.all_tab_cols
		WHERE OWNER = '%s' AND USER_GENERATED = 'YES'
		ORDER BY TABLE_NAME, COLUMN_ID NULLS LAST, COLUMN_NAME`, sqlLiteral(schemaName))
	} else {
		// Failing to tell virtual columns apart only keeps their expressions
		// in Default, as before, so it must not fail the sync.
		var err error
		virtualColumns, err = getVirtualColumns(txn, schemaName)
		if err != nil {
			slog.Warn("failed to list Oracle virtual columns", slog.String("schema", schemaName), slog.String("error", err.Error()))
		}
		query = fmt.Sprintf(`
		SELECT
			OWNER,
			TABLE_NAME,
			COLUMN_NAME,
			DATA_TYPE,
			DATA_LENGTH,
			DATA_PRECISION,
			DATA_SCALE,
			CHAR_USED,
			CHAR_LENGTH,
			COLUMN_ID,
			DATA_DEFAULT,
			NULLABLE,
			NULL,
			NULL,
			NULL,
			NULL,
			NULL
		FROM sys.all_tab_columns
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME, COLUMN_ID NULLS LAST, COLUMN_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get columns query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var schemaName, tableName, nullable string
		var defaultStr, collation, defaultOnNull, virtualColumn sql.NullString
		var dataLength, dataPrecision, dataScale, charLength sql.NullInt64
		var charUsed, identityColumnValue sql.NullString
		var columnID, internalColumnID sql.NullInt32
		if err := rows.Scan(
			&schemaName,
			&tableName,
			&column.Name,
			&column.Type,
			&dataLength,
			&dataPrecision,
			&dataScale,
			&charUsed,
			&charLength,
			&columnID,
			&defaultStr,
			&nullable,
			&collation,
			&defaultOnNull,
			&virtualColumn,
			&identityColumnValue,
			&internalColumnID,
		); err != nil {
			return nil, nil, err
		}
		column.Type = getTypeString(column.Type, dataLength, dataPrecision, dataScale, charLength, charUsed)
		if defaultStr.Valid {
			// Clean up the default expression to remove current schema name and fix NEXTVAL syntax
			cleanedDefault := cleanDefaultExpression(defaultStr.String, schemaName)
			// DATA_DEFAULT holds a virtual column's expression.
			if virtualColumn.String == "YES" || virtualColumns[columnKey{Schema: schemaName, Table: tableName, Column: column.Name}] {
				column.Generation = &metadatapb.GenerationMetadata{
					Type:       metadatapb.GenerationMetadata_TYPE_VIRTUAL,
					Expression: cleanedDefault,
				}
			} else {
				column.Default = cleanedDefault
			}
		}
		isNullBool, err := convertYesNo(nullable)
		if err != nil {
			return nil, nil, err
		}
		column.Nullable = isNullBool
		if collation.Valid {
			column.Collation = collation.String
		}
		if identityColumnValue.String == "YES" {
			identity := identityColumns[columnKey{Schema: schemaName, Table: tableName, Column: column.Name}]
			column.IsIdentity = true
			column.IdentityGeneration = identity.generation
			column.IdentitySeed = identity.seed
			column.IdentityIncrement = identity.increment
			column.Options = identity.options
		}
		if defaultOnNull.Valid && defaultOnNull.String == "YES" {
			column.DefaultOnNull = true
		}

		key := tableKey{Schema: schemaName, Table: tableName}
		if internalColumnID.Valid {
			internalPos[columnKey{Schema: schemaName, Table: tableName, Column: column.Name}] = internalColumnID.Int32
		}
		if columnID.Valid {
			column.Position = columnID.Int32
		} else {
			column.IsInvisible = true
			column.Position = 1
			if columns := columnsMap[key]; len(columns) > 0 {
				column.Position = columns[len(columns)-1].Position + 1
			}
		}
		columnsMap[key] = append(columnsMap[key], column)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return columnsMap, internalPos, nil
}

// getVirtualColumns lists the virtual columns for versions whose column query
// reads ALL_TAB_COLUMNS, which has no VIRTUAL_COLUMN.
func getVirtualColumns(txn *sql.Tx, schemaName string) (map[columnKey]bool, error) {
	query := fmt.Sprintf(`
		SELECT TABLE_NAME, COLUMN_NAME
		FROM sys.all_tab_cols
		WHERE OWNER = '%s' AND VIRTUAL_COLUMN = 'YES'`, sqlLiteral(schemaName))
	slog.Debug("running get virtual columns query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	virtualColumns := make(map[columnKey]bool)
	for rows.Next() {
		var tableName, columnName string
		if err := rows.Scan(&tableName, &columnName); err != nil {
			return nil, err
		}
		virtualColumns[columnKey{Schema: schemaName, Table: tableName, Column: columnName}] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return virtualColumns, nil
}

func getTypeString(dataType string, dataLength, dataPrecision, dataScale, charLength sql.NullInt64, charUsed sql.NullString) string {
	switch dataType {
	case "VARCHAR2", "CHAR":
		// A CHAR declaration writes its length in characters or in bytes, and
		// the dictionary states which: reading DATA_LENGTH as bytes would turn
		// VARCHAR2(50 CHAR) into VARCHAR2(200 BYTE) on a 4-byte character set.
		if charUsed.String == "C" && charLength.Valid {
			return fmt.Sprintf("%s(%d CHAR)", dataType, charLength.Int64)
		}
		return fmt.Sprintf("%s(%d BYTE)", dataType, dataLength.Int64)
	case "NVARCHAR2", "NCHAR":
		// The national character types declare their length in characters.
		if charLength.Valid {
			return fmt.Sprintf("%s(%d)", dataType, charLength.Int64)
		}
		return fmt.Sprintf("%s(%d)", dataType, dataLength.Int64)
	case "RAW", "UROWID":
		return fmt.Sprintf("%s(%d)", dataType, dataLength.Int64)
	case "NUMBER":
		// A wildcard precision, NUMBER(*, s), keeps its scale and has none.
		// A scale may be negative, which rounds to the left of the point.
		switch {
		case !dataPrecision.Valid:
			if dataScale.Valid {
				return fmt.Sprintf("%s(*,%d)", dataType, dataScale.Int64)
			}
		case dataPrecision.Int64 > 0:
			if dataScale.Valid && dataScale.Int64 != 0 {
				return fmt.Sprintf("%s(%d,%d)", dataType, dataPrecision.Int64, dataScale.Int64)
			}
			return fmt.Sprintf("%s(%d)", dataType, dataPrecision.Int64)
		default:
			// do nothing
		}
	case "FLOAT":
		switch {
		case !dataPrecision.Valid || dataPrecision.Int64 == 0:
		// do nothing
		case dataPrecision.Valid && dataPrecision.Int64 > 0:
			return fmt.Sprintf("%s(%d)", dataType, dataPrecision.Int64)
		default:
			// do nothing
		}
	default:
		// return dataType as-is for unhandled types
	}
	return dataType
}

func getOuterSchemaRColumns(txn *sql.Tx, outerRTableMap map[constraintRef]string, outerRColumnMap map[constraintRef][]string, schemaName, constraintName string) (string, []string, error) {
	queryColumns := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			CONSTRAINT_NAME,
			COLUMN_NAME
		FROM
			SYS.ALL_CONS_COLUMNS
		WHERE
			OWNER = '%s'
		ORDER BY TABLE_NAME, CONSTRAINT_NAME, POSITION`, sqlLiteral(schemaName))

	slog.Debug("running get outer schema reference columns query")
	rows, err := txn.Query(queryColumns)
	if err != nil {
		return "", nil, fmt.Errorf("failed to execute query %q: %w", queryColumns, err)
	}
	defer func() { _ = rows.Close() }()

	// The query reads every constraint of the owner, and a constraint two
	// foreign keys name is read once per call, so each constraint's columns start
	// again in this call and then take its rows in position order.
	cleared := make(map[constraintRef]bool, 8)
	for rows.Next() {
		var tableName, constraintName string
		var columnName sql.NullString
		if err := rows.Scan(&tableName, &constraintName, &columnName); err != nil {
			return "", nil, err
		}
		if !columnName.Valid {
			continue
		}
		key := constraintRef{Schema: schemaName, Constraint: constraintName}
		outerRTableMap[key] = tableName
		if !cleared[key] {
			cleared[key] = true
			outerRColumnMap[key] = nil
		}
		outerRColumnMap[key] = append(outerRColumnMap[key], columnName.String)
	}
	if err := rows.Err(); err != nil {
		return "", nil, fmt.Errorf("failed to execute query %q: %w", queryColumns, err)
	}

	constraintKey := constraintRef{Schema: schemaName, Constraint: constraintName}
	return outerRTableMap[constraintKey], outerRColumnMap[constraintKey], nil
}

// quoteOracleColumns writes column names as the quoted identifiers an index's
// expressions hold, which is how the index path writes them too.
func quoteOracleColumns(columns []string) []string {
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = fmt.Sprintf(`"%s"`, column)
	}
	return quoted
}

// constraintIndex is what a schema's index enforces: the constraints that draw
// on it. Oracle keeps a constraint and its index as two objects with two names —
// a constraint may reuse an index the schema already had, and the index's own
// name is what the dictionary, and so the DDL, calls it — and one index may
// enforce several constraints, each on keys of its own.
type constraintIndex struct {
	// Constraints are the constraints the index enforces, the primary key's
	// first, each with the keys it declares.
	Constraints []*metadatapb.IndexConstraintMetadata
	// Primary is whether one of them is the primary key.
	Primary bool
}

func getConstraints(txn *sql.Tx, schemaName string) (
	map[tableKey][]*metadatapb.IndexMetadata,
	map[tableKey][]*metadatapb.CheckConstraintMetadata,
	map[tableKey][]*metadatapb.ForeignKeyMetadata,
	map[indexKey]constraintIndex,
	error,
) {
	queryConstraintColumns := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			CONSTRAINT_NAME,
			COLUMN_NAME
		FROM SYS.ALL_CONS_COLUMNS
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME, CONSTRAINT_NAME, POSITION`, sqlLiteral(schemaName))

	slog.Debug("running get constraint columns query")
	constraintColumnRows, err := txn.Query(queryConstraintColumns)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryConstraintColumns, err)
	}
	defer func() { _ = constraintColumnRows.Close() }()
	constraintColumnMap := make(map[constraintRef][]string)
	constraintTableMap := make(map[constraintRef]string)
	for constraintColumnRows.Next() {
		var tableName, constraintName string
		var columnName sql.NullString
		if err := constraintColumnRows.Scan(&tableName, &constraintName, &columnName); err != nil {
			return nil, nil, nil, nil, err
		}
		key := constraintRef{Schema: schemaName, Constraint: constraintName}
		if columnName.Valid {
			constraintColumnMap[key] = append(constraintColumnMap[key], columnName.String)
		}
		constraintTableMap[key] = tableName
	}
	if err := constraintColumnRows.Err(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryConstraintColumns, err)
	}

	queryConstraints := fmt.Sprintf(`
		SELECT
			TABLE_NAME,
			CONSTRAINT_NAME,
			CONSTRAINT_TYPE,
			SEARCH_CONDITION,
			R_OWNER,
			R_CONSTRAINT_NAME,
			DELETE_RULE,
			INDEX_OWNER,
			INDEX_NAME
		FROM SYS.ALL_CONSTRAINTS
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME, CONSTRAINT_NAME`, sqlLiteral(schemaName))

	slog.Debug("running get constraints query")
	constraintRows, err := txn.Query(queryConstraints)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryConstraints, err)
	}
	defer func() { _ = constraintRows.Close() }()
	indexMap := make(map[tableKey][]*metadatapb.IndexMetadata)
	checkConstraintMap := make(map[tableKey][]*metadatapb.CheckConstraintMetadata)
	foreignKeyMap := make(map[tableKey][]*metadatapb.ForeignKeyMetadata)
	isConstraint := make(map[indexKey]constraintIndex)
	outerRColumnMap := make(map[constraintRef][]string)
	outerRTableMap := make(map[constraintRef]string)
	for constraintRows.Next() {
		var tableName, constraintName, constraintType string
		var searchCondition, rOwner, rConstraintName, deleteRule, indexOwner, indexName sql.NullString
		if err := constraintRows.Scan(&tableName, &constraintName, &constraintType, &searchCondition, &rOwner, &rConstraintName, &deleteRule, &indexOwner, &indexName); err != nil {
			return nil, nil, nil, nil, err
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		constraintKey := constraintRef{Schema: schemaName, Constraint: constraintName}
		switch constraintType {
		case "P", "U":
			// An enabled primary key or unique constraint is enforced by an index
			// the schema owns, which may be one the schema already had under its
			// own name: ALL_CONSTRAINTS.INDEX_NAME says which, and the index is
			// read as an index below, with the constraint laid over it. A
			// constraint with no index of the schema's is carried as the
			// constraint it is: a disabled one has none at all (INDEX_NAME is
			// null), and one whose index another schema owns has no index in the
			// schema being synced.
			constraint := &metadatapb.IndexConstraintMetadata{
				Name:    constraintName,
				Columns: quoteOracleColumns(constraintColumnMap[constraintKey]),
				Primary: constraintType == "P",
			}
			if indexName.Valid && indexOwner.Valid && indexOwner.String == schemaName {
				index := indexKey{Schema: indexOwner.String, Table: tableName, Index: indexName.String}
				enforced := isConstraint[index]
				if constraint.Primary {
					enforced.Constraints = append([]*metadatapb.IndexConstraintMetadata{constraint}, enforced.Constraints...)
					enforced.Primary = true
				} else {
					enforced.Constraints = append(enforced.Constraints, constraint)
				}
				isConstraint[index] = enforced
				break
			}
			metadata := &metadatapb.IndexMetadata{
				Name:         constraintName,
				Primary:      constraint.Primary,
				Unique:       true,
				IsConstraint: true,
				Constraints:  []*metadatapb.IndexConstraintMetadata{constraint},
			}
			metadata.Expressions = constraint.Columns
			indexMap[key] = append(indexMap[key], metadata)
		case "C":
			// Oracle records NOT NULL as a check constraint it named SYS_C<n>
			// itself. A user's unnamed CHECK gets a SYS_C name as well, so the
			// condition says which one this is.
			if strings.HasPrefix(constraintName, "SYS_C") && isGeneratedNotNull(searchCondition.String) {
				continue
			}
			constraint := &metadatapb.CheckConstraintMetadata{
				Name: constraintName,
			}
			if searchCondition.Valid {
				constraint.Expression = searchCondition.String
			}
			checkConstraintMap[key] = append(checkConstraintMap[key], constraint)
		case "R":
			if rOwner.Valid && rConstraintName.Valid {
				foreignKey := &metadatapb.ForeignKeyMetadata{
					Name:    constraintName,
					Columns: constraintColumnMap[constraintKey],
				}
				if deleteRule.Valid && deleteRule.String != "NO ACTION" {
					foreignKey.OnDelete = deleteRule.String
				}
				if rOwner.String == schemaName {
					// Same-schema reference - don't set ReferencedSchema for portability
					rConstraintKey := constraintRef{Schema: rOwner.String, Constraint: rConstraintName.String}
					foreignKey.ReferencedTable = constraintTableMap[rConstraintKey]
					foreignKey.ReferencedColumns = constraintColumnMap[rConstraintKey]
				} else {
					// Cross-schema reference - set ReferencedSchema
					foreignKey.ReferencedSchema = rOwner.String
					foreignKey.ReferencedTable, foreignKey.ReferencedColumns, err = getOuterSchemaRColumns(txn, outerRTableMap, outerRColumnMap, rOwner.String, rConstraintName.String)
					if err != nil {
						return nil, nil, nil, nil, fmt.Errorf("failed to get outer schema reference columns: %w", err)
					}
				}
				foreignKeyMap[key] = append(foreignKeyMap[key], foreignKey)
			}
		default:
			// Ignore other constraint types
		}
	}
	if err := constraintRows.Err(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryConstraints, err)
	}

	return indexMap, checkConstraintMap, foreignKeyMap, isConstraint, nil
}

// getIndexesAndConstraints reads the schema's indexes and their constraints:
// the indexes of its tables by table, and the indexes it owns on another
// schema's table or on a cluster, which are its too and no table of it holds
// (TestSyncOtherIndexes). An index carries no comment: Oracle names no index
// in its COMMENT statement, and TestSyncCommentableObjects pins the refusal.
func getIndexesAndConstraints(txn *sql.Tx, schemaName string) (map[tableKey][]*metadatapb.IndexMetadata, []*metadatapb.SchemaIndexMetadata, map[tableKey][]*metadatapb.CheckConstraintMetadata, map[tableKey][]*metadatapb.ForeignKeyMetadata, error) {
	indexMap, checkConstraintMap, foreignKeyMap, constraintIndexMap, err := getConstraints(txn, schemaName)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get constraints: %w", err)
	}

	// The columns are read by the index's owner, which is the schema, so an
	// index it owns on another schema's table has its columns, and one another
	// schema owns on its table, which is not the schema's, has none.
	queryColumn := fmt.Sprintf(`
		SELECT INDEX_OWNER, TABLE_NAME, INDEX_NAME, COLUMN_NAME, DESCEND
		FROM sys.all_ind_columns
		WHERE INDEX_OWNER = '%s'
		ORDER BY TABLE_NAME, INDEX_NAME, COLUMN_POSITION`, sqlLiteral(schemaName))

	slog.Debug("running get index column query")
	indexExpressionMap := make(map[indexKey][]string)
	indexColumnMap := make(map[indexKey][]string)
	descendingMap := make(map[indexKey][]bool)
	colRows, err := txn.Query(queryColumn)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryColumn, err)
	}
	defer func() { _ = colRows.Close() }()
	for colRows.Next() {
		var schemaName, tableName, indexName, columnName string
		var descend sql.NullString
		if err := colRows.Scan(&schemaName, &tableName, &indexName, &columnName, &descend); err != nil {
			return nil, nil, nil, nil, err
		}
		key := indexKey{Schema: schemaName, Table: tableName, Index: indexName}

		// Clean up column name but preserve the original if it's already clean
		cleanColumnName := columnName
		if strings.HasPrefix(columnName, "\"") && strings.HasSuffix(columnName, "\"") {
			cleanColumnName = columnName[1 : len(columnName)-1]
		}

		indexColumnMap[key] = append(indexColumnMap[key], cleanColumnName)
		descendingMap[key] = append(descendingMap[key], descend.String == "DESC")
	}
	if err := colRows.Err(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryColumn, err)
	}
	if err := colRows.Close(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	queryExpression := fmt.Sprintf(`
		SELECT INDEX_OWNER, TABLE_NAME, INDEX_NAME, COLUMN_EXPRESSION
		FROM sys.all_ind_expressions
		WHERE INDEX_OWNER = '%s'
		ORDER BY TABLE_NAME, INDEX_NAME, COLUMN_POSITION`, sqlLiteral(schemaName))

	slog.Debug("running get index expression query")
	expRows, err := txn.Query(queryExpression)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryExpression, err)
	}
	defer func() { _ = expRows.Close() }()
	for expRows.Next() {
		var schemaName, tableName, indexName, columnExpression string
		if err := expRows.Scan(&schemaName, &tableName, &indexName, &columnExpression); err != nil {
			return nil, nil, nil, nil, err
		}
		key := indexKey{Schema: schemaName, Table: tableName, Index: indexName}
		indexExpressionMap[key] = append(indexExpressionMap[key], columnExpression)
	}
	if err := expRows.Err(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryExpression, err)
	}
	if err := expRows.Close(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT OWNER, TABLE_OWNER, TABLE_NAME, TABLE_TYPE, INDEX_NAME, UNIQUENESS, INDEX_TYPE, VISIBILITY
		FROM sys.all_indexes
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME, INDEX_NAME`, sqlLiteral(schemaName))

	slog.Debug("running get index query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	var otherIndexes []*metadatapb.SchemaIndexMetadata
	for rows.Next() {
		index := &metadatapb.IndexMetadata{}
		var schemaName, tableOwner, tableName, tableType, unique string
		var visibility sql.NullString
		// INDEX_TYPE is NORMAL, or FUNCTION-BASED NORMAL.
		if err := rows.Scan(&schemaName, &tableOwner, &tableName, &tableType, &index.Name, &unique, &index.Type, &visibility); err != nil {
			return nil, nil, nil, nil, err
		}
		key := indexKey{Schema: schemaName, Table: tableName, Index: index.Name}
		constraint, isConstraintIndex := constraintIndexMap[key]

		// Skip system-generated indexes, such as those of LOBs, IOTs, and nested
		// tables. A constraint's index is the schema's own whatever it is called:
		// Oracle names the index of a constraint the user did not name SYS_C<n>.
		if !isConstraintIndex && strings.HasPrefix(index.Name, "SYS_") {
			continue
		}

		// The index's own row describes it, a constraint's or not: a constraint
		// that reuses an index the schema already had makes the two names, and
		// the two uniquenesses, differ.
		index.Unique = unique == "UNIQUE"
		if isConstraintIndex {
			index.Primary = constraint.Primary
			index.IsConstraint = true
			index.Constraints = constraint.Constraints
		}
		columns := indexColumnMap[key]
		expressions := indexExpressionMap[key]

		// For function-based indexes (including those with DESC columns)
		if len(expressions) > 0 {
			// Combine regular columns with expressions
			// Oracle reports ASC columns in ALL_IND_COLUMNS and DESC columns as expressions
			var combinedExpressions []string
			exprIdx := 0

			for i := 0; i < len(columns); i++ {
				if strings.HasPrefix(columns[i], "SYS_NC") {
					// This is a virtual column for a DESC column, use the expression
					if exprIdx < len(expressions) {
						combinedExpressions = append(combinedExpressions, expressions[exprIdx])
						exprIdx++
					}
				} else {
					// Regular column
					combinedExpressions = append(combinedExpressions, fmt.Sprintf(`"%s"`, columns[i]))
				}
			}

			index.Expressions = combinedExpressions
		} else {
			// Simple column index without expressions
			// Skip indexes that reference only system-generated virtual columns
			hasOnlyVirtualColumns := true
			for _, col := range columns {
				if !strings.HasPrefix(col, "SYS_NC") {
					hasOnlyVirtualColumns = false
					break
				}
			}
			if hasOnlyVirtualColumns && !isConstraintIndex {
				continue
			}

			// For column-based indexes, quote the column names to match expected format
			quotedColumns := make([]string, len(columns))
			for i, col := range columns {
				quotedColumns[i] = fmt.Sprintf(`"%s"`, col)
			}
			index.Expressions = quotedColumns
		}
		// Note: Keep index.Type as reported by Oracle (don't modify it)

		// Set descending flags for all indexes (both function-based and normal)
		if desc, ok := descendingMap[key]; ok && len(desc) > 0 {
			// Always set descending array if we have column information,
			// even if all columns are ascending (to maintain consistency)
			index.Descending = desc
		}
		index.Visible = true
		if visibility.Valid && visibility.String == "INVISIBLE" {
			index.Visible = false
		}

		// An index on another schema's table, or on a cluster (TABLE_TYPE
		// CLUSTER), is the schema's and on none of its tables.
		if tableOwner != schemaName || tableType != "TABLE" {
			otherIndexes = append(otherIndexes, &metadatapb.SchemaIndexMetadata{Schema: tableOwner, Table: tableName, Index: index})
			continue
		}
		table := tableKey{Schema: schemaName, Table: tableName}
		indexMap[table] = append(indexMap[table], index)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	// ALL_INDEXES lists no cluster's index, even to its owner, where
	// ALL_OBJECTS and ALL_IND_COLUMNS do (TestSyncOtherIndexes), so the
	// schema's indexes ALL_INDEXES leaves out are read from those.
	queryHidden := fmt.Sprintf(`
		SELECT o.OBJECT_NAME, c.TABLE_OWNER, c.TABLE_NAME, k.CLUSTER_NAME
		FROM sys.all_objects o
		LEFT JOIN (
			SELECT DISTINCT INDEX_NAME, TABLE_OWNER, TABLE_NAME
			FROM sys.all_ind_columns
			WHERE INDEX_OWNER = '%[1]s'
		) c ON c.INDEX_NAME = o.OBJECT_NAME
		LEFT JOIN sys.all_clusters k ON k.OWNER = c.TABLE_OWNER AND k.CLUSTER_NAME = c.TABLE_NAME
		WHERE o.OWNER = '%[1]s' AND o.OBJECT_TYPE = 'INDEX'
		AND o.OBJECT_NAME NOT IN (SELECT INDEX_NAME FROM sys.all_indexes WHERE OWNER = '%[1]s')
		ORDER BY o.OBJECT_NAME`, sqlLiteral(schemaName))
	slog.Debug("running get hidden index query")
	hiddenRows, err := txn.Query(queryHidden)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryHidden, err)
	}
	defer func() { _ = hiddenRows.Close() }()
	for hiddenRows.Next() {
		var name string
		var tableOwner, tableName, cluster sql.NullString
		if err := hiddenRows.Scan(&name, &tableOwner, &tableName, &cluster); err != nil {
			return nil, nil, nil, nil, err
		}
		if strings.HasPrefix(name, "SYS_") {
			continue
		}
		index := &metadatapb.IndexMetadata{Name: name, Visible: true}
		if cluster.Valid {
			index.Type = "CLUSTER"
		}
		for _, col := range indexColumnMap[indexKey{Schema: schemaName, Table: tableName.String, Index: name}] {
			index.Expressions = append(index.Expressions, fmt.Sprintf(`"%s"`, col))
		}
		otherIndexes = append(otherIndexes, &metadatapb.SchemaIndexMetadata{Schema: tableOwner.String, Table: tableName.String, Index: index})
	}
	if err := hiddenRows.Err(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryHidden, err)
	}
	if err := hiddenRows.Close(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}
	slices.SortFunc(otherIndexes, func(a, b *metadatapb.SchemaIndexMetadata) int {
		return strings.Compare(a.Index.Name, b.Index.Name)
	})

	return indexMap, otherIndexes, checkConstraintMap, foreignKeyMap, nil
}

// getViews gets all views of a database.
func getViews(txn *sql.Tx, schemaName string, columnMap map[tableKey][]*metadatapb.ColumnMetadata, internalPos map[columnKey]int32, triggerMap map[tableKey][]*metadatapb.TriggerMetadata) (map[string][]*metadatapb.ViewMetadata, error) {
	viewMap := make(map[string][]*metadatapb.ViewMetadata)

	// Get view comments
	viewCommentMap, err := getViewComments(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get view comments: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT OWNER, VIEW_NAME, TEXT
		FROM sys.all_views
		WHERE OWNER = '%s'
		ORDER BY view_name
	`, sqlLiteral(schemaName))

	slog.Debug("running get view query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		view := &metadatapb.ViewMetadata{}
		var schemaName string
		if err := rows.Scan(&schemaName, &view.Name, &view.Definition); err != nil {
			return nil, err
		}
		view.Definition = sanitizeOracleDefinition(schemaName, "VIEW", view.Name, view.Definition)
		key := tableKey{Schema: schemaName, Table: view.Name}
		view.Columns = orderViewColumns(schemaName, view.Name, columnMap[key], internalPos)
		view.Triggers = triggerMap[key]
		if comment, ok := viewCommentMap[key]; ok {
			view.Comment = comment
		}

		viewMap[schemaName] = append(viewMap[schemaName], view)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	// Fetch all view dependencies in a single query to avoid N+1 problem
	viewDependenciesMap, err := getAllViewDependencies(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get view dependencies: %w", err)
	}

	// Populate view dependencies from the map
	for schemaName, list := range viewMap {
		for _, view := range list {
			key := tableKey{Schema: schemaName, Table: view.Name}
			if dependencies, ok := viewDependenciesMap[key]; ok {
				view.DependencyColumns = dependencies
			}
		}
	}

	return viewMap, nil
}

// getAllViewDependencies gets all view dependencies for a schema in a single query.
func getAllViewDependencies(txn *sql.Tx, schemaName string) (map[tableKey][]*metadatapb.DependencyColumn, error) {
	dependenciesMap := make(map[tableKey][]*metadatapb.DependencyColumn)

	query := fmt.Sprintf(`
		SELECT 
			NAME as view_name,
			REFERENCED_OWNER as source_schema,
			REFERENCED_NAME as source_table,
			'*' as column_name
		FROM ALL_DEPENDENCIES
		WHERE OWNER = '%s'
			AND TYPE = 'VIEW'
			AND REFERENCED_TYPE IN ('TABLE', 'VIEW')
		ORDER BY NAME, REFERENCED_OWNER, REFERENCED_NAME
	`, sqlLiteral(schemaName))

	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var viewName string
		dependencyColumn := &metadatapb.DependencyColumn{}
		if err := rows.Scan(&viewName, &dependencyColumn.Schema, &dependencyColumn.Table, &dependencyColumn.Column); err != nil {
			return nil, err
		}

		key := tableKey{Schema: schemaName, Table: viewName}
		dependenciesMap[key] = append(dependenciesMap[key], dependencyColumn)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dependenciesMap, nil
}

// getAllMaterializedViewDependencies gets all materialized view dependencies for a schema in a single query.
func getAllMaterializedViewDependencies(txn *sql.Tx, schemaName string) (map[tableKey][]*metadatapb.DependencyColumn, error) {
	dependenciesMap := make(map[tableKey][]*metadatapb.DependencyColumn)

	query := fmt.Sprintf(`
		SELECT 
			NAME as mview_name,
			REFERENCED_OWNER as source_schema,
			REFERENCED_NAME as source_table,
			'*' as column_name
		FROM ALL_DEPENDENCIES
		WHERE OWNER = '%s'
			AND TYPE = 'MATERIALIZED VIEW'
			AND REFERENCED_TYPE IN ('TABLE', 'VIEW')
		ORDER BY NAME, REFERENCED_OWNER, REFERENCED_NAME
	`, sqlLiteral(schemaName))

	rows, err := txn.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var mviewName string
		dependencyColumn := &metadatapb.DependencyColumn{}
		if err := rows.Scan(&mviewName, &dependencyColumn.Schema, &dependencyColumn.Table, &dependencyColumn.Column); err != nil {
			return nil, err
		}

		key := tableKey{Schema: schemaName, Table: mviewName}
		dependenciesMap[key] = append(dependenciesMap[key], dependencyColumn)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dependenciesMap, nil
}

// getMaterializedViews gets all materialized views of a database.
func getMaterializedViews(txn *sql.Tx, schemaName string, columnMap map[tableKey][]*metadatapb.ColumnMetadata, indexMap map[tableKey][]*metadatapb.IndexMetadata, triggerMap map[tableKey][]*metadatapb.TriggerMetadata, partitionMap map[tableKey][]*metadatapb.TablePartitionMetadata) ([]*metadatapb.MaterializedViewMetadata, error) {
	var materializedViews []*metadatapb.MaterializedViewMetadata

	// Get materialized view comments
	materializedViewCommentMap, err := getMaterializedViewComments(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get materialized view comments: %w", err)
	}

	// The container table has the view's name; ALL_TABLES gives its
	// tablespace as it gives a table's.
	query := fmt.Sprintf(`
		SELECT m.OWNER, m.MVIEW_NAME, m.QUERY, m.BUILD_MODE, m.REFRESH_METHOD, m.REFRESH_MODE, t.TABLESPACE_NAME
		FROM sys.all_mviews m
		LEFT JOIN sys.all_tables t ON t.OWNER = m.OWNER AND t.TABLE_NAME = m.MVIEW_NAME
		WHERE m.OWNER = '%s'
		ORDER BY m.MVIEW_NAME
	`, sqlLiteral(schemaName))

	slog.Debug("running get materialized views query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		materializedView := &metadatapb.MaterializedViewMetadata{}
		var schemaName string
		var buildMode, refreshMethod, refreshMode, tablespace sql.NullString
		if err := rows.Scan(&schemaName, &materializedView.Name, &materializedView.Definition, &buildMode, &refreshMethod, &refreshMode, &tablespace); err != nil {
			return nil, err
		}
		// As a table's: Oracle gives no tablespace to a partitioned or
		// index-organized container (TestSyncMaterializedViewTablespace).
		materializedView.Tablespace = tablespace.String
		materializedView.Definition = sanitizeOracleDefinition(
			schemaName,
			"MATERIALIZED VIEW",
			materializedView.Name,
			materializedView.Definition,
		)
		if prefix := materializedViewRefreshPrefix(buildMode.String, refreshMethod.String, refreshMode.String); prefix != "" && materializedView.Definition != "" {
			materializedView.Definition = prefix + " AS " + strings.TrimSpace(materializedView.Definition)
		}

		// Ensure the definition ends with a newline to match expected format
		if materializedView.Definition != "" && !strings.HasSuffix(materializedView.Definition, "\n") {
			materializedView.Definition += "\n"
		}

		// Add comment if available
		key := tableKey{Schema: schemaName, Table: materializedView.Name}
		if comment, ok := materializedViewCommentMap[key]; ok {
			materializedView.Comment = comment
		}
		// Columns, indexes, triggers and partitions are the container
		// table's, which has the view's name.
		materializedView.Columns = columnMap[key]
		materializedView.Indexes = indexMap[key]
		materializedView.Triggers = triggerMap[key]
		materializedView.Partitions = partitionMap[key]

		materializedViews = append(materializedViews, materializedView)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	// Fetch all materialized view dependencies in a single query to avoid N+1 problem
	materializedViewDependenciesMap, err := getAllMaterializedViewDependencies(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get materialized view dependencies: %w", err)
	}

	// Populate dependencies from the map
	for _, materializedView := range materializedViews {
		key := tableKey{Schema: schemaName, Table: materializedView.Name}
		if dependencies, ok := materializedViewDependenciesMap[key]; ok {
			materializedView.DependencyColumns = dependencies
		}
	}

	return materializedViews, nil
}

// sequenceFacts are the facts ALL_SEQUENCES states about a sequence, in the
// order getSequences reads them. number marks a NUMBER column, which is read as
// text: MAX_VALUE reaches 28 digits, which no Go integer holds.
var sequenceFacts = []struct {
	column string
	number bool
}{
	{"SEQUENCE_NAME", false},
	{"MIN_VALUE", true},
	{"MAX_VALUE", true},
	{"INCREMENT_BY", true},
	{"CYCLE_FLAG", false},
	{"ORDER_FLAG", false},
	{"CACHE_SIZE", true},
	{"LAST_NUMBER", true},
	{"KEEP_VALUE", false},
	{"SCALE_FLAG", false},
	{"EXTEND_FLAG", false},
	{"SESSION_FLAG", false},
	{"SHARDED_FLAG", false},
}

// getSequences gets all sequences of a database. Every fact ALL_SEQUENCES
// states is carried, except the declared START WITH: Oracle does not retain it,
// ALL_SEQUENCES has no column for it, and DBMS_METADATA reports LAST_NUMBER in
// its place, so start stays empty and last_value carries LAST_NUMBER.
// TestSyncSequenceValues pins that, on a sequence a NEXTVAL has advanced.
func getSequences(txn *sql.Tx, schemaName string) ([]*metadatapb.SequenceMetadata, error) {
	var sequences []*metadatapb.SequenceMetadata

	columns, err := sequenceColumns(txn)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT %s
		FROM ALL_SEQUENCES
		WHERE SEQUENCE_OWNER = '%s'
		ORDER BY SEQUENCE_NAME`, sequenceProjection(columns), sqlLiteral(schemaName))

	slog.Debug("running get sequences query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		seq := &metadatapb.SequenceMetadata{}
		var cycle, order, keep, scale, extend, session, sharded sql.NullString
		// The scan follows sequenceFacts.
		if err := rows.Scan(
			&seq.Name,
			&seq.MinValue,
			&seq.MaxValue,
			&seq.Increment,
			&cycle,
			&order,
			&seq.CacheSize,
			&seq.LastValue,
			&keep,
			&scale,
			&extend,
			&session,
			&sharded,
		); err != nil {
			return nil, err
		}
		// Skip system-generated sequences (e.g., ISEQ$$_* for IDENTITY columns)
		if strings.HasPrefix(seq.Name, "ISEQ$$_") {
			continue
		}
		for _, flag := range []struct {
			column string
			value  sql.NullString
			target *bool
		}{
			{"CYCLE_FLAG", cycle, &seq.Cycle},
			{"ORDER_FLAG", order, &seq.Order},
			{"KEEP_VALUE", keep, &seq.Keep},
			{"SCALE_FLAG", scale, &seq.Scale},
			{"EXTEND_FLAG", extend, &seq.Extend},
			{"SESSION_FLAG", session, &seq.Session},
			{"SHARDED_FLAG", sharded, &seq.Sharded},
		} {
			// A release with no column for the flag says nothing about it.
			if !flag.value.Valid {
				continue
			}
			on, err := convertYesNo(flag.value.String)
			if err != nil {
				return nil, fmt.Errorf("sequence %s: %s: %w", seq.Name, flag.column, err)
			}
			*flag.target = on
		}
		sequences = append(sequences, seq)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return sequences, nil
}

// sequenceColumns returns the columns ALL_SEQUENCES has. The view gained them
// release by release — SCALE_FLAG, EXTEND_FLAG and SHARDED_FLAG with sharded
// sequences, SESSION_FLAG with session sequences, KEEP_VALUE with KEEP — and a
// release without one answers ORA-00904 for it, so the facts this server can
// state are the view's own column list rather than a release this code guesses
// at.
func sequenceColumns(txn *sql.Tx) (map[string]bool, error) {
	query := `
		SELECT COLUMN_NAME
		FROM ALL_TAB_COLS
		WHERE OWNER = 'SYS' AND TABLE_NAME = 'ALL_SEQUENCES'`

	slog.Debug("running get sequence columns query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	columns := make(map[string]bool)
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, err
		}
		columns[column] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	return columns, nil
}

// sequenceProjection returns the select list that reads every fact of a
// sequence out of a view with columns: a fact the view has no column for is
// read as NULL, so the query holds for every release.
func sequenceProjection(columns map[string]bool) string {
	selected := make([]string, 0, len(sequenceFacts))
	for _, fact := range sequenceFacts {
		switch {
		case !columns[fact.column]:
			selected = append(selected, "NULL")
		case fact.number:
			selected = append(selected, "TO_CHAR("+fact.column+")")
		default:
			selected = append(selected, fact.column)
		}
	}
	return strings.Join(selected, ", ")
}

func getRoutines(txn *sql.Tx, schemaName string) ([]*metadatapb.FunctionMetadata, []*metadatapb.ProcedureMetadata, []*metadatapb.PackageMetadata, error) {
	var functions []*metadatapb.FunctionMetadata
	var procedures []*metadatapb.ProcedureMetadata
	var packages []*metadatapb.PackageMetadata
	packageIndex := make(map[string]int)

	query := fmt.Sprintf(`
		SELECT
			NAME,
			TYPE,
			TEXT
		FROM ALL_SOURCE
		WHERE
			TYPE IN ('FUNCTION', 'PROCEDURE', 'PACKAGE', 'PACKAGE BODY')
			AND
			OWNER = '%s'
		ORDER BY NAME, TYPE, LINE`, sqlLiteral(schemaName))

	slog.Debug("running get routines query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	flush := func(currentName, currentType string, defText []string) {
		if currentType == "" {
			return
		}
		text := sanitizeOracleDefinition(schemaName, currentType, currentName, strings.Join(defText, ""))
		switch currentType {
		case "FUNCTION":
			functions = append(functions, &metadatapb.FunctionMetadata{
				Name:       currentName,
				Definition: text,
			})
		case "PROCEDURE":
			procedures = append(procedures, &metadatapb.ProcedureMetadata{
				Name:       currentName,
				Definition: text,
			})
		case "PACKAGE":
			packageIndex[currentName] = len(packages)
			packages = append(packages, &metadatapb.PackageMetadata{
				Name:       currentName,
				Definition: text,
			})
		case "PACKAGE BODY":
			if idx, ok := packageIndex[currentName]; ok {
				spec := ensurePLSQLUnitCreateAndSemicolon(packages[idx].Definition)
				body := ensurePLSQLUnitCreateAndSemicolon(text)
				packages[idx].Definition = spec + "\n\n" + body
			} else {
				packageIndex[currentName] = len(packages)
				packages = append(packages, &metadatapb.PackageMetadata{
					Name:       currentName,
					Definition: text,
				})
			}
		}
	}

	var currentName, currentType string
	var defText []string
	for rows.Next() {
		var name, t, def string
		if err := rows.Scan(&name, &t, &def); err != nil {
			return nil, nil, nil, err
		}
		if name == currentName && t == currentType {
			defText = append(defText, def)
		} else {
			flush(currentName, currentType, defText)
			currentName = name
			currentType = t
			defText = []string{def}
		}
	}
	flush(currentName, currentType, defText)
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, nil, nil, err
	}

	return functions, procedures, packages, nil
}

func ensurePLSQLUnitCreateAndSemicolon(unit string) string {
	trimmed := strings.TrimSpace(unit)
	if trimmed == "" {
		return ""
	}
	if !strings.HasPrefix(strings.ToUpper(trimmed), "CREATE") {
		trimmed = "CREATE OR REPLACE " + trimmed
	}
	if !strings.HasSuffix(trimmed, ";") {
		trimmed += ";"
	}
	return trimmed
}

func materializedViewRefreshPrefix(buildMode, refreshMethod, refreshMode string) string {
	build := strings.ToUpper(strings.TrimSpace(buildMode))
	method := strings.ToUpper(strings.TrimSpace(refreshMethod))
	mode := strings.ToUpper(strings.TrimSpace(refreshMode))

	// Materialized view logs (CREATE MATERIALIZED VIEW LOG) are not synced or
	// rendered yet, so emitting REFRESH FAST would fail when restoring into an
	// empty schema (ORA-23413). Fall back to Oracle's default FORCE method.
	if method == "FAST" {
		method = "FORCE"
	}

	// Omit the prefix when the effective clauses match Oracle's default
	// CREATE MATERIALIZED VIEW AS ... settings (BUILD IMMEDIATE REFRESH FORCE
	// ON DEMAND). BUILD_MODE='PREBUILT' is intentionally not emitted as
	// ON PREBUILT TABLE because getTables excludes tables listed in ALL_MVIEWS,
	// so a dump restored into an empty schema does not recreate the prebuilt
	// backing table first.
	if build != "DEFERRED" && (method == "" || method == "FORCE") && (mode == "" || mode == "DEMAND") {
		return ""
	}
	var parts []string
	if build == "DEFERRED" {
		parts = append(parts, "BUILD DEFERRED")
	}
	if method == "NEVER" {
		parts = append(parts, "NEVER REFRESH")
	} else if method != "" || mode != "" {
		refresh := "REFRESH"
		if method != "" {
			refresh += " " + method
		}
		if mode != "" && mode != "NEVER" {
			refresh += " ON " + mode
		}
		parts = append(parts, refresh)
	}
	return strings.Join(parts, " ")
}

func orderViewColumns(schemaName, viewName string, columns []*metadatapb.ColumnMetadata, internalPos map[columnKey]int32) []*metadatapb.ColumnMetadata {
	if len(columns) <= 1 || len(internalPos) == 0 {
		return columns
	}
	ordered := make([]*metadatapb.ColumnMetadata, len(columns))
	copy(ordered, columns)
	slices.SortStableFunc(ordered, func(a, b *metadatapb.ColumnMetadata) int {
		posA := internalPos[columnKey{Schema: schemaName, Table: viewName, Column: a.GetName()}]
		posB := internalPos[columnKey{Schema: schemaName, Table: viewName, Column: b.GetName()}]
		if posA != posB {
			if posA < posB {
				return -1
			}
			return 1
		}
		return int(a.GetPosition() - b.GetPosition())
	})
	for i, col := range ordered {
		col.Position = int32(i + 1)
	}
	return ordered
}

func oraclePartitionType(name string) metadatapb.TablePartitionMetadata_Type {
	switch strings.ToUpper(strings.TrimSpace(name)) {
	case "RANGE":
		return metadatapb.TablePartitionMetadata_RANGE
	case "LIST":
		return metadatapb.TablePartitionMetadata_LIST
	case "HASH":
		return metadatapb.TablePartitionMetadata_HASH
	default:
		return metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED
	}
}

func formatOraclePartitionKeyColumn(name string) string {
	if isSimpleOracleIdentifier(name) {
		return name
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func isSimpleOracleIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		isLetter := 'A' <= c && c <= 'Z'
		isDigit := '0' <= c && c <= '9'
		if i == 0 && !isLetter {
			return false
		}
		if !isLetter && !isDigit && c != '_' && c != '$' && c != '#' {
			return false
		}
	}
	return true
}

func getTablePartitions(txn *sql.Tx, schemaName string) (map[tableKey][]*metadatapb.TablePartitionMetadata, error) {
	result := make(map[tableKey][]*metadatapb.TablePartitionMetadata)

	type partTableInfo struct {
		partType    metadatapb.TablePartitionMetadata_Type
		subpartType metadatapb.TablePartitionMetadata_Type
		interval    string
	}
	partTables := make(map[string]partTableInfo)
	hasSubpartitions := false

	queryTables := fmt.Sprintf(`
		SELECT TABLE_NAME, PARTITIONING_TYPE, SUBPARTITIONING_TYPE, INTERVAL
		FROM all_part_tables
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName))
	slog.Debug("running get partitioned tables query")
	rows, err := txn.Query(queryTables)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryTables, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var tableName, partType, subpartType string
		var interval sql.NullString
		if err := rows.Scan(&tableName, &partType, &subpartType, &interval); err != nil {
			return nil, err
		}
		pt := oraclePartitionType(partType)
		if pt == metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
			continue
		}
		spt := oraclePartitionType(subpartType)
		if spt != metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
			hasSubpartitions = true
		}
		partTables[tableName] = partTableInfo{
			partType:    pt,
			subpartType: spt,
			interval:    strings.TrimSpace(interval.String),
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryTables, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}
	if len(partTables) == 0 {
		return result, nil
	}

	partKeys := make(map[string][]string)
	queryKeys := fmt.Sprintf(`
		SELECT NAME, COLUMN_NAME
		FROM all_part_key_columns
		WHERE OWNER = '%s' AND OBJECT_TYPE = 'TABLE'
		ORDER BY NAME, COLUMN_POSITION`, sqlLiteral(schemaName))
	keyRows, err := txn.Query(queryKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryKeys, err)
	}
	defer func() { _ = keyRows.Close() }()
	for keyRows.Next() {
		var tableName, colName string
		if err := keyRows.Scan(&tableName, &colName); err != nil {
			return nil, err
		}
		partKeys[tableName] = append(partKeys[tableName], formatOraclePartitionKeyColumn(colName))
	}
	if err := keyRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryKeys, err)
	}
	if err := keyRows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	type partRef struct {
		table     string
		partition string
	}
	subpartMap := make(map[partRef][]*metadatapb.TablePartitionMetadata)
	if hasSubpartitions {
		subpartKeys := make(map[string][]string)
		querySubKeys := fmt.Sprintf(`
			SELECT NAME, COLUMN_NAME
			FROM all_subpart_key_columns
			WHERE OWNER = '%s' AND OBJECT_TYPE = 'TABLE'
			ORDER BY NAME, COLUMN_POSITION`, sqlLiteral(schemaName))
		subKeyRows, err := txn.Query(querySubKeys)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", querySubKeys, err)
		}
		defer func() { _ = subKeyRows.Close() }()
		for subKeyRows.Next() {
			var tableName, colName string
			if err := subKeyRows.Scan(&tableName, &colName); err != nil {
				return nil, err
			}
			subpartKeys[tableName] = append(subpartKeys[tableName], formatOraclePartitionKeyColumn(colName))
		}
		if err := subKeyRows.Err(); err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", querySubKeys, err)
		}
		if err := subKeyRows.Close(); err != nil {
			return nil, fmt.Errorf("failed to close rows: %w", err)
		}

		querySubparts := fmt.Sprintf(`
			SELECT TABLE_NAME, PARTITION_NAME, SUBPARTITION_NAME, HIGH_VALUE
			FROM all_tab_subpartitions
			WHERE TABLE_OWNER = '%s'
			ORDER BY TABLE_NAME, PARTITION_NAME, SUBPARTITION_POSITION`, sqlLiteral(schemaName))
		subRows, err := txn.Query(querySubparts)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", querySubparts, err)
		}
		defer func() { _ = subRows.Close() }()
		for subRows.Next() {
			var tableName, partName, subpartName string
			var highValue sql.NullString
			if err := subRows.Scan(&tableName, &partName, &subpartName, &highValue); err != nil {
				return nil, err
			}
			info, ok := partTables[tableName]
			if !ok || info.subpartType == metadatapb.TablePartitionMetadata_TYPE_UNSPECIFIED {
				continue
			}
			ref := partRef{table: tableName, partition: partName}
			subpartMap[ref] = append(subpartMap[ref], &metadatapb.TablePartitionMetadata{
				Name:       subpartName,
				Type:       info.subpartType,
				Expression: strings.Join(subpartKeys[tableName], ", "),
				Value:      strings.TrimSpace(highValue.String),
			})
		}
		if err := subRows.Err(); err != nil {
			return nil, fmt.Errorf("failed to execute query %q: %w", querySubparts, err)
		}
		if err := subRows.Close(); err != nil {
			return nil, fmt.Errorf("failed to close rows: %w", err)
		}
	}

	queryParts := fmt.Sprintf(`
		SELECT TABLE_NAME, PARTITION_NAME, HIGH_VALUE, INTERVAL
		FROM all_tab_partitions
		WHERE TABLE_OWNER = '%s'
		ORDER BY TABLE_NAME, PARTITION_POSITION`, sqlLiteral(schemaName))
	partRows, err := txn.Query(queryParts)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryParts, err)
	}
	defer func() { _ = partRows.Close() }()
	for partRows.Next() {
		var tableName, partName string
		var highValue, intervalFlag sql.NullString
		if err := partRows.Scan(&tableName, &partName, &highValue, &intervalFlag); err != nil {
			return nil, err
		}
		info, ok := partTables[tableName]
		if !ok {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(intervalFlag.String), "YES") {
			continue
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		result[key] = append(result[key], &metadatapb.TablePartitionMetadata{
			Name:           partName,
			Type:           info.partType,
			Expression:     strings.Join(partKeys[tableName], ", "),
			Value:          strings.TrimSpace(highValue.String),
			IntervalClause: info.interval,
			Subpartitions:  subpartMap[partRef{table: tableName, partition: partName}],
		})
	}
	if err := partRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", queryParts, err)
	}
	if err := partRows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}
	return result, nil
}

func sanitizeOracleDefinition(schemaName, objectType, objectName, definition string) string {
	return sanitizeOracleMetadataString(schemaName, objectType, objectName, "definition", definition)
}

func sanitizeOracleMetadataString(schemaName, objectType, objectName, field, value string) string {
	// Oracle may include a trailing C-string null terminator via the go-ora driver.
	value = strings.TrimRight(value, "\x00")
	if utf8.ValidString(value) {
		return value
	}
	slog.Warn("sanitized invalid UTF-8 in Oracle metadata",
		slog.String("schema", schemaName),
		slog.String("object_type", objectType),
		slog.String("object_name", objectName),
		slog.String("field", field),
	)

	var sanitized strings.Builder
	sanitized.Grow(len(value))
	for i := 0; i < len(value); {
		r, width := utf8.DecodeRuneInString(value[i:])
		if r == utf8.RuneError && width == 1 {
			_, _ = fmt.Fprintf(&sanitized, "\\x%02x", value[i])
			i++
			continue
		}
		_, _ = sanitized.WriteString(value[i : i+width])
		i += width
	}
	return sanitized.String()
}

// cleanDefaultExpression cleans up Oracle default expressions to make them more portable
// and fix common syntax issues like schema-qualified sequence NEXTVAL calls
func cleanDefaultExpression(defaultExpr, currentSchema string) string {
	if defaultExpr == "" {
		return defaultExpr
	}

	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(defaultExpr)

	// A string literal is a value, not an expression: its text may hold the
	// qualification a sequence default carries, which it keeps.
	if strings.HasPrefix(cleaned, "'") {
		return cleaned
	}

	// Fix schema-qualified sequence NEXTVAL calls
	// Pattern: "SCHEMA"."SEQUENCE_NAME"."NEXTVAL" -> "SEQUENCE_NAME".NEXTVAL
	// Keep sequence name quoted for compatibility
	if strings.Contains(cleaned, ".\"NEXTVAL\"") {
		// Use regex to match quoted schema and sequence names followed by quoted NEXTVAL
		re := regexp.MustCompile(`"` + regexp.QuoteMeta(currentSchema) + `"\."([^"]+)"\."NEXTVAL"`)
		cleaned = re.ReplaceAllString(cleaned, `"$1".NEXTVAL`)
	}

	// Also handle unquoted NEXTVAL (less common)
	if strings.Contains(cleaned, ".NEXTVAL") {
		re := regexp.MustCompile(`"` + regexp.QuoteMeta(currentSchema) + `"\."([^"]+)"\.NEXTVAL`)
		cleaned = re.ReplaceAllString(cleaned, `"$1".NEXTVAL`)
	}

	return cleaned
}

// getViewComments gets comments for views: COMMENT ON TABLE is the form that
// comments a view, and it reports here with TABLE_TYPE = 'VIEW'.
func getViewComments(txn *sql.Tx, schemaName string) (map[tableKey]string, error) {
	viewCommentMap := make(map[tableKey]string)

	// Oracle stores view comments in ALL_TAB_COMMENTS where TABLE_TYPE = 'VIEW'
	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE = 'VIEW'
		ORDER BY OWNER, TABLE_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE = 'VIEW'
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get view comments query")
	rows, err := txn.Query(query)
	if err != nil {
		// If the query fails (e.g., TABLE_TYPE = 'VIEW' not supported), return empty map
		slog.Debug("view comments query failed, returning empty map", slog.String("error", err.Error()))
		return viewCommentMap, nil
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var schemaName, viewName, comment string
		if err := rows.Scan(&schemaName, &viewName, &comment); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: viewName}
		viewCommentMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return viewCommentMap, nil
}

// getMaterializedViewComments gets comments for materialized views from Oracle system views
func getMaterializedViewComments(txn *sql.Tx, schemaName string) (map[tableKey]string, error) {
	materializedViewCommentMap := make(map[tableKey]string)

	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, MVIEW_NAME, COMMENTS
		FROM all_mview_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		ORDER BY OWNER, MVIEW_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, MVIEW_NAME, COMMENTS
		FROM all_mview_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		ORDER BY MVIEW_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get materialized view comments query")
	rows, err := txn.Query(query)
	if err != nil {
		// If the view doesn't exist, return empty map
		slog.Debug("materialized view comments query failed, returning empty map", slog.String("error", err.Error()))
		return materializedViewCommentMap, nil
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var schemaName, mvName, comment string
		if err := rows.Scan(&schemaName, &mvName, &comment); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: mvName}
		materializedViewCommentMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return materializedViewCommentMap, nil
}
