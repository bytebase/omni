package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
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
	version, err := d.version(ctx)
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
	columnMap, err := getTableColumns(txn, schemaName, version, identityColumns)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns from database %q: %w", schemaName, err)
	}
	tableTriggerMap, viewTriggerMap, err := getTriggers(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get triggers from database %q: %w", schemaName, err)
	}
	indexMap, checkConstraintMap, foreignKeyMap, err := getIndexesAndConstraints(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices: %w", err)
	}
	tableMap, err := getTables(txn, schemaName, columnMap, indexMap, checkConstraintMap, foreignKeyMap, tableTriggerMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables from database %q: %w", schemaName, err)
	}
	viewMap, err := getViews(txn, schemaName, columnMap, viewTriggerMap)
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
	materializedViews, err := getMaterializedViews(txn, schemaName, indexMap, tableTriggerMap)
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

	// Get trigger comments
	triggerCommentMap, err := getTriggerComments(txn, schemaName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get trigger comments: %w", err)
	}

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
		// Add trigger comment if available
		triggerKey := tableKey{Schema: schemaName, Table: triggerName.String}
		if comment, ok := triggerCommentMap[triggerKey]; ok {
			trigger.Comment = comment
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
			if err := rows.Scan(&schemaName, &table.Name, &count); err != nil {
				return err
			}
			table.RowCount = int64(count.Float64)
			key := tableKey{Schema: schemaName, Table: table.Name}
			table.Columns = columnMap[key]
			table.Indexes = indexMap[key]
			table.CheckConstraints = checkConstraintMap[key]
			table.ForeignKeys = foreignKeyMap[key]
			table.Triggers = triggerMap[key]
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
		SELECT OWNER, TABLE_NAME, NUM_ROWS
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
		SELECT OWNER, TABLE_NAME, NUM_ROWS
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
			identity.seed, identity.increment = identityOptions(options.String)
		}
		identityColumns[columnKey{Schema: schemaName, Table: tableName, Column: columnName}] = identity
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return identityColumns, nil
}

// identityOptions reads START WITH and INCREMENT BY out of the option list Oracle
// writes for an identity, such as
// "START WITH: 5, INCREMENT BY: 3, MAX_VALUE: 999…, CYCLE_FLAG: N".
func identityOptions(options string) (seed, increment int64) {
	seed, increment = 1, 1
	for _, part := range strings.Split(options, ",") {
		name, value, ok := strings.Cut(part, ":")
		if !ok {
			continue
		}
		parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			continue
		}
		switch strings.TrimSpace(name) {
		case "START WITH":
			seed = parsed
		case "INCREMENT BY":
			increment = parsed
		}
	}
	return seed, increment
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

// getTableColumns gets the columns of a table.
func getTableColumns(txn *sql.Tx, schemaName string, version *oracleanalysis.Version, identityColumns map[columnKey]identityColumn) (map[tableKey][]*metadatapb.ColumnMetadata, error) {
	columnsMap := make(map[tableKey][]*metadatapb.ColumnMetadata)

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
			IDENTITY_COLUMN
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
			NULL
		FROM sys.all_tab_columns
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME, COLUMN_ID NULLS LAST, COLUMN_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get columns query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		column := &metadatapb.ColumnMetadata{}
		var schemaName, tableName, nullable string
		var defaultStr, collation, defaultOnNull, virtualColumn sql.NullString
		var dataLength, dataPrecision, dataScale, charLength sql.NullInt64
		var charUsed, identityColumnValue sql.NullString
		var columnID sql.NullInt32
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
		); err != nil {
			return nil, err
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
			return nil, err
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
		}
		if defaultOnNull.Valid && defaultOnNull.String == "YES" {
			column.DefaultOnNull = true
		}

		key := tableKey{Schema: schemaName, Table: tableName}
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
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return columnsMap, nil
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

func getConstraints(txn *sql.Tx, schemaName string) (
	map[tableKey][]*metadatapb.IndexMetadata,
	map[tableKey][]*metadatapb.CheckConstraintMetadata,
	map[tableKey][]*metadatapb.ForeignKeyMetadata,
	map[indexKey]bool,
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
			DELETE_RULE
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
	isConstraint := make(map[indexKey]bool)
	outerRColumnMap := make(map[constraintRef][]string)
	outerRTableMap := make(map[constraintRef]string)
	for constraintRows.Next() {
		var tableName, constraintName, constraintType string
		var searchCondition, rOwner, rConstraintName, deleteRule sql.NullString
		if err := constraintRows.Scan(&tableName, &constraintName, &constraintType, &searchCondition, &rOwner, &rConstraintName, &deleteRule); err != nil {
			return nil, nil, nil, nil, err
		}
		key := tableKey{Schema: schemaName, Table: tableName}
		constraintKey := constraintRef{Schema: schemaName, Constraint: constraintName}
		switch constraintType {
		case "P":
			index := &metadatapb.IndexMetadata{
				Name:         constraintName,
				Primary:      true,
				Unique:       true,
				IsConstraint: true,
			}
			if columns, ok := constraintColumnMap[constraintKey]; ok {
				index.Expressions = columns
			}
			indexMap[key] = append(indexMap[key], index)
			isConstraint[indexKey{Schema: schemaName, Table: tableName, Index: constraintName}] = true
		case "U":
			index := &metadatapb.IndexMetadata{
				Name:         constraintName,
				Unique:       true,
				IsConstraint: true,
			}
			if columns, ok := constraintColumnMap[constraintKey]; ok {
				index.Expressions = columns
			}
			indexMap[key] = append(indexMap[key], index)
			isConstraint[indexKey{Schema: schemaName, Table: tableName, Index: constraintName}] = true
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

// getIndexes gets all indices and constraints of a database.
func getIndexesAndConstraints(txn *sql.Tx, schemaName string) (map[tableKey][]*metadatapb.IndexMetadata, map[tableKey][]*metadatapb.CheckConstraintMetadata, map[tableKey][]*metadatapb.ForeignKeyMetadata, error) {
	indexMap, checkConstraintMap, foreignKeyMap, isConstraint, err := getConstraints(txn, schemaName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get constraints: %w", err)
	}

	// Get index comments
	indexCommentMap, err := getIndexComments(txn, schemaName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get index comments: %w", err)
	}

	queryColumn := fmt.Sprintf(`
		SELECT TABLE_OWNER, TABLE_NAME, INDEX_NAME, COLUMN_NAME, DESCEND
		FROM sys.all_ind_columns
		WHERE TABLE_OWNER = '%s'
		ORDER BY TABLE_NAME, INDEX_NAME, COLUMN_POSITION`, sqlLiteral(schemaName))

	slog.Debug("running get index column query")
	indexExpressionMap := make(map[indexKey][]string)
	indexColumnMap := make(map[indexKey][]string)
	descendingMap := make(map[indexKey][]bool)
	colRows, err := txn.Query(queryColumn)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryColumn, err)
	}
	defer func() { _ = colRows.Close() }()
	for colRows.Next() {
		var schemaName, tableName, indexName, columnName string
		var descend sql.NullString
		if err := colRows.Scan(&schemaName, &tableName, &indexName, &columnName, &descend); err != nil {
			return nil, nil, nil, err
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
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryColumn, err)
	}
	if err := colRows.Close(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	queryExpression := fmt.Sprintf(`
		SELECT TABLE_OWNER, TABLE_NAME, INDEX_NAME, COLUMN_EXPRESSION
		FROM sys.all_ind_expressions
		WHERE TABLE_OWNER = '%s'
		ORDER BY TABLE_NAME, INDEX_NAME, COLUMN_POSITION`, sqlLiteral(schemaName))

	slog.Debug("running get index expression query")
	expRows, err := txn.Query(queryExpression)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryExpression, err)
	}
	defer func() { _ = expRows.Close() }()
	for expRows.Next() {
		var schemaName, tableName, indexName, columnExpression string
		if err := expRows.Scan(&schemaName, &tableName, &indexName, &columnExpression); err != nil {
			return nil, nil, nil, err
		}
		key := indexKey{Schema: schemaName, Table: tableName, Index: indexName}
		indexExpressionMap[key] = append(indexExpressionMap[key], columnExpression)
	}
	if err := expRows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", queryExpression, err)
	}
	if err := expRows.Close(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, INDEX_NAME, UNIQUENESS, INDEX_TYPE, VISIBILITY
		FROM sys.all_indexes
		WHERE OWNER = '%s'
		ORDER BY TABLE_NAME, INDEX_NAME`, sqlLiteral(schemaName))

	slog.Debug("running get index query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		index := &metadatapb.IndexMetadata{}
		var schemaName, tableName, unique string
		var visibility sql.NullString
		// INDEX_TYPE is NORMAL, or FUNCTION-BASED NORMAL.
		if err := rows.Scan(&schemaName, &tableName, &index.Name, &unique, &index.Type, &visibility); err != nil {
			return nil, nil, nil, err
		}
		if isConstraint[indexKey{Schema: schemaName, Table: tableName, Index: index.Name}] {
			continue
		}

		// Skip system-generated indexes, such as those of LOBs, IOTs, and nested tables.
		if strings.HasPrefix(index.Name, "SYS_") {
			continue
		}

		index.Unique = unique == "UNIQUE"
		indexKey := indexKey{Schema: schemaName, Table: tableName, Index: index.Name}
		// Handle index expressions - preserve Oracle's classification
		columns := indexColumnMap[indexKey]
		expressions := indexExpressionMap[indexKey]

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
			if hasOnlyVirtualColumns {
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
		if desc, ok := descendingMap[indexKey]; ok && len(desc) > 0 {
			// Always set descending array if we have column information,
			// even if all columns are ascending (to maintain consistency)
			index.Descending = desc
		}
		index.Visible = true
		if visibility.Valid && visibility.String == "INVISIBLE" {
			index.Visible = false
		}

		// Add index comment if available
		if comment, ok := indexCommentMap[indexKey]; ok {
			index.Comment = comment
		}

		key := tableKey{Schema: schemaName, Table: tableName}
		indexMap[key] = append(indexMap[key], index)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return indexMap, checkConstraintMap, foreignKeyMap, nil
}

// getViews gets all views of a database.
func getViews(txn *sql.Tx, schemaName string, columnMap map[tableKey][]*metadatapb.ColumnMetadata, triggerMap map[tableKey][]*metadatapb.TriggerMetadata) (map[string][]*metadatapb.ViewMetadata, error) {
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
		view.Columns = columnMap[key]
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
func getMaterializedViews(txn *sql.Tx, schemaName string, indexMap map[tableKey][]*metadatapb.IndexMetadata, triggerMap map[tableKey][]*metadatapb.TriggerMetadata) ([]*metadatapb.MaterializedViewMetadata, error) {
	var materializedViews []*metadatapb.MaterializedViewMetadata

	// Get materialized view comments
	materializedViewCommentMap, err := getMaterializedViewComments(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get materialized view comments: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT OWNER, MVIEW_NAME, QUERY
		FROM sys.all_mviews
		WHERE OWNER = '%s'
		ORDER BY MVIEW_NAME
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
		if err := rows.Scan(&schemaName, &materializedView.Name, &materializedView.Definition); err != nil {
			return nil, err
		}
		materializedView.Definition = sanitizeOracleDefinition(
			schemaName,
			"MATERIALIZED VIEW",
			materializedView.Name,
			materializedView.Definition,
		)

		// Ensure the definition ends with a newline to match expected format
		if materializedView.Definition != "" && !strings.HasSuffix(materializedView.Definition, "\n") {
			materializedView.Definition += "\n"
		}

		// Add comment if available
		key := tableKey{Schema: schemaName, Table: materializedView.Name}
		if comment, ok := materializedViewCommentMap[key]; ok {
			materializedView.Comment = comment
		}
		// Indexes and triggers are on the container table, which has the view's name.
		materializedView.Indexes = indexMap[key]
		materializedView.Triggers = triggerMap[key]

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

// getSequences gets all sequences of a database.
func getSequences(txn *sql.Tx, schemaName string) ([]*metadatapb.SequenceMetadata, error) {
	var sequences []*metadatapb.SequenceMetadata

	// Get sequence comments
	sequenceCommentMap, err := getSequenceComments(txn, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to get sequence comments: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT SEQUENCE_NAME FROM ALL_SEQUENCES
		WHERE SEQUENCE_OWNER = '%s'
		ORDER BY SEQUENCE_NAME`, sqlLiteral(schemaName))

	slog.Debug("running get sequences query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		seq := &metadatapb.SequenceMetadata{}
		if err := rows.Scan(&seq.Name); err != nil {
			return nil, err
		}
		// Skip system-generated sequences (e.g., ISEQ$$_* for IDENTITY columns)
		if strings.HasPrefix(seq.Name, "ISEQ$$_") {
			continue
		}
		key := tableKey{Schema: schemaName, Table: seq.Name}
		if comment, ok := sequenceCommentMap[key]; ok {
			seq.Comment = comment
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

func getRoutines(txn *sql.Tx, schemaName string) ([]*metadatapb.FunctionMetadata, []*metadatapb.ProcedureMetadata, []*metadatapb.PackageMetadata, error) {
	var functions []*metadatapb.FunctionMetadata
	var procedures []*metadatapb.ProcedureMetadata
	var packages []*metadatapb.PackageMetadata

	// Get function comments
	functionCommentMap, err := getFunctionComments(txn, schemaName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get function comments: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			NAME,
			TYPE,
			TEXT
		FROM ALL_SOURCE
		WHERE
			TYPE IN ('FUNCTION', 'PROCEDURE', 'PACKAGE')
			AND
			OWNER = '%s'
		ORDER BY NAME, TYPE, LINE`, sqlLiteral(schemaName))

	slog.Debug("running get routines query")
	rows, err := txn.Query(query)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	defer func() { _ = rows.Close() }()

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
			// Oracle may include a trailing C-string null terminator via the go-ora driver.
			switch currentType {
			case "FUNCTION":
				function := &metadatapb.FunctionMetadata{
					Name:       currentName,
					Definition: sanitizeOracleDefinition(schemaName, "FUNCTION", currentName, strings.Join(defText, "")),
				}
				key := tableKey{Schema: schemaName, Table: currentName}
				if comment, ok := functionCommentMap[key]; ok {
					function.Comment = comment
				}
				functions = append(functions, function)
			case "PROCEDURE":
				procedures = append(procedures, &metadatapb.ProcedureMetadata{
					Name:       currentName,
					Definition: sanitizeOracleDefinition(schemaName, "PROCEDURE", currentName, strings.Join(defText, "")),
				})
			case "PACKAGE":
				packages = append(packages, &metadatapb.PackageMetadata{
					Name:       currentName,
					Definition: sanitizeOracleDefinition(schemaName, "PACKAGE", currentName, strings.Join(defText, "")),
				})
			default:
				// Ignore other types
			}
			currentName = name
			currentType = t
			defText = []string{def}
		}
	}
	// Oracle may include a trailing C-string null terminator via the go-ora driver.
	switch currentType {
	case "FUNCTION":
		function := &metadatapb.FunctionMetadata{
			Name:       currentName,
			Definition: sanitizeOracleDefinition(schemaName, "FUNCTION", currentName, strings.Join(defText, "")),
		}
		key := tableKey{Schema: schemaName, Table: currentName}
		if comment, ok := functionCommentMap[key]; ok {
			function.Comment = comment
		}
		functions = append(functions, function)
	case "PROCEDURE":
		procedures = append(procedures, &metadatapb.ProcedureMetadata{
			Name:       currentName,
			Definition: sanitizeOracleDefinition(schemaName, "PROCEDURE", currentName, strings.Join(defText, "")),
		})
	case "PACKAGE":
		packages = append(packages, &metadatapb.PackageMetadata{
			Name:       currentName,
			Definition: sanitizeOracleDefinition(schemaName, "PACKAGE", currentName, strings.Join(defText, "")),
		})
	default:
		// Ignore other types
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, nil, nil, err
	}

	return functions, procedures, packages, nil
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

// getViewComments gets comments for views from Oracle system views
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

// getSequenceComments gets comments for sequences from Oracle system views
func getSequenceComments(txn *sql.Tx, schemaName string) (map[tableKey]string, error) {
	sequenceCommentMap := make(map[tableKey]string)

	// Oracle doesn't have a direct sequence comments view, but we can try user_tab_comments
	// where TABLE_TYPE might be 'SEQUENCE' in some Oracle versions
	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE = 'SEQUENCE'
		ORDER BY OWNER, TABLE_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE = 'SEQUENCE'
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get sequence comments query")
	rows, err := txn.Query(query)
	if err != nil {
		// If the query fails (e.g., TABLE_TYPE = 'SEQUENCE' not supported), return empty map
		slog.Debug("sequence comments query failed, returning empty map", slog.String("error", err.Error()))
		return sequenceCommentMap, nil
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var schemaName, sequenceName, comment string
		if err := rows.Scan(&schemaName, &sequenceName, &comment); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: sequenceName}
		sequenceCommentMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return sequenceCommentMap, nil
}

// getFunctionComments gets comments for functions and procedures from Oracle system views
func getFunctionComments(txn *sql.Tx, schemaName string) (map[tableKey]string, error) {
	functionCommentMap := make(map[tableKey]string)

	// Oracle stores function/procedure comments in ALL_OBJECTS or USER_OBJECTS
	// However, Oracle doesn't have a standard system view for function/procedure comments
	// We can try to get them from ALL_TAB_COMMENTS where TABLE_TYPE might be 'FUNCTION' or 'PROCEDURE'
	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE IN ('FUNCTION', 'PROCEDURE')
		ORDER BY OWNER, TABLE_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE IN ('FUNCTION', 'PROCEDURE')
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get function comments query")
	rows, err := txn.Query(query)
	if err != nil {
		// If the query fails (e.g., TABLE_TYPE for functions not supported), return empty map
		slog.Debug("function comments query failed, returning empty map", slog.String("error", err.Error()))
		return functionCommentMap, nil
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var schemaName, functionName, comment string
		if err := rows.Scan(&schemaName, &functionName, &comment); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: functionName}
		functionCommentMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return functionCommentMap, nil
}

// getIndexComments gets comments for indexes from Oracle system views
func getIndexComments(_ *sql.Tx, schemaName string) (map[indexKey]string, error) {
	indexCommentMap := make(map[indexKey]string)

	// Oracle doesn't have a dedicated index comments view in all versions
	// We need to join ALL_INDEXES with ALL_TAB_COMMENTS to get index comments
	// However, Oracle stores index comments differently depending on the version
	// For now, return empty map since Oracle index comments are not consistently available
	slog.Debug("Oracle index comments not implemented - returning empty map", "schema", schemaName)
	return indexCommentMap, nil
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

// getTriggerComments gets comments for triggers from Oracle system views
func getTriggerComments(txn *sql.Tx, schemaName string) (map[tableKey]string, error) {
	triggerCommentMap := make(map[tableKey]string)

	// Oracle doesn't have a standard trigger comments view
	// Trigger comments might be stored in ALL_TAB_COMMENTS with TABLE_TYPE = 'TRIGGER'
	// or in some Oracle versions, they might not be available at all
	query := ""
	if schemaName == "" {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER NOT IN (%s) AND OWNER NOT LIKE 'APEX_%%' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE = 'TRIGGER'
		ORDER BY OWNER, TABLE_NAME`, systemSchema)
	} else {
		query = fmt.Sprintf(`
		SELECT OWNER, TABLE_NAME, COMMENTS
		FROM all_tab_comments
		WHERE OWNER = '%s' AND COMMENTS IS NOT NULL
		AND TABLE_TYPE = 'TRIGGER'
		ORDER BY TABLE_NAME`, sqlLiteral(schemaName))
	}

	slog.Debug("running get trigger comments query")
	rows, err := txn.Query(query)
	if err != nil {
		// If the query fails (e.g., TABLE_TYPE = 'TRIGGER' not supported), return empty map
		slog.Debug("trigger comments query failed, returning empty map", slog.String("error", err.Error()))
		return triggerCommentMap, nil
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var schemaName, triggerName, comment string
		if err := rows.Scan(&schemaName, &triggerName, &comment); err != nil {
			return nil, err
		}
		key := tableKey{Schema: schemaName, Table: triggerName}
		triggerCommentMap[key] = sanitizeUTF8String(comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", query, err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return triggerCommentMap, nil
}
