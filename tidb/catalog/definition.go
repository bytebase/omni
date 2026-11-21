package catalog

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/tidb/ast"
)

const (
	disableUniqueAndForeignKeyCheckStmt = "SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;\nSET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;\n"
	restoreUniqueAndForeignKeyCheckStmt = "SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;\nSET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;\n"

	emptyCommentLine      = "--\n"
	tempViewHeader        = "-- Temporary view structure for "
	tableHeader           = "-- Table structure for "
	viewHeader            = "-- View structure for "
	sequenceHeader        = "-- Sequence structure for "
	functionHeader        = "-- Function structure for "
	procedureHeader       = "-- Procedure structure for "
	triggerHeader         = "-- Trigger structure for "
	eventHeader           = "-- Event structure for "
	setCharacterSetClient = "SET character_set_client = "
	setCharacterSetResult = "SET character_set_results = "
	setCollation          = "SET collation_connection = "
	setSQLMode            = "SET sql_mode = "
	setTimezone           = "SET time_zone = "
	delimiterDoubleSemi   = "DELIMITER ;;\n"
	delimiterSemi         = "DELIMITER ;\n"

	mysqlTypeBlob       = "blob"
	mysqlTypeTinyBob    = "tinyblob"
	mysqlTypeMediumBlob = "mediumblob"
	mysqlTypeLongBlob   = "longblob"
	mysqlTypeJSON       = "json"
	mysqlTypeGeometry   = "geometry"

	mysqlIndexFullText = "FULLTEXT"
	mysqlIndexSpatial  = "SPATIAL"

	mysqlNoAction = "NO ACTION"

	autoIncrementSymbol = "AUTO_INCREMENT"
	autoRandomSymbol    = "AUTO_RANDOM"
)

// DefinitionOptions selects the form DatabaseDefinition renders.
type DefinitionOptions struct {
	// PrintHeader wraps the output in statements that disable and restore unique
	// and foreign key checks.
	PrintHeader bool
}

// DatabaseDefinition renders a synced database snapshot as a mysqldump-style
// restore script. It renders only the first schema and returns "" for a snapshot
// with none.
func DatabaseDefinition(meta *metadata.DatabaseSchemaMetadata, opts DefinitionOptions) (string, error) {
	if len(meta.Schemas) == 0 {
		return "", nil
	}

	var buf strings.Builder

	if opts.PrintHeader {
		if _, err := buf.WriteString(disableUniqueAndForeignKeyCheckStmt); err != nil {
			return "", err
		}
	}

	schema := meta.Schemas[0]

	if err := writeDatabaseOptions(&buf, meta); err != nil {
		return "", err
	}

	// Construct temporal views.
	// Create a temporary view with the same name as the view and with columns of
	// the same name in order to satisfy views that depend on this view.
	// This temporary view will be removed when the actual view is created.
	// The properties of each column, are not preserved in this temporary
	// view. They are not necessary because other views only need to reference
	// the column name, thus we generate SELECT 1 AS colName1, 1 AS colName2.
	// This will not be necessary once we can determine dependencies
	// between views and can simply dump them in the appropriate order.
	// https://sourcegraph.com/github.com/mysql/mysql-server/-/blob/client/mysqldump.cc?L2781
	for _, view := range schema.Views {
		if len(view.Columns) == 0 {
			if err := writeInvalidTemporaryView(&buf, view); err != nil {
				return "", err
			}
			continue
		}
		if err := writeTemporaryView(&buf, view); err != nil {
			return "", err
		}
	}

	// Sequences come before tables: a default may read one.
	for _, sequence := range schema.Sequences {
		if err := writeSequence(&buf, sequence); err != nil {
			return "", err
		}
	}

	// Construct tables.
	for _, table := range schema.Tables {
		if err := writeTable(&buf, table, meta.Name); err != nil {
			return "", err
		}
	}

	// Construct views.
	for _, view := range schema.Views {
		if err := writeView(&buf, view); err != nil {
			return "", err
		}
	}

	// Construct functions.
	for _, function := range schema.Functions {
		if err := writeFunction(&buf, function); err != nil {
			return "", err
		}
	}

	// Construct procedures.
	for _, procedure := range schema.Procedures {
		if err := writeProcedure(&buf, procedure); err != nil {
			return "", err
		}
	}

	// Construct events.
	for _, event := range schema.Events {
		if err := writeEvent(&buf, event); err != nil {
			return "", err
		}
	}

	// Construct triggers.
	for _, table := range schema.Tables {
		for _, trigger := range table.Triggers {
			if err := writeTrigger(&buf, table.Name, trigger); err != nil {
				return "", err
			}
		}
	}

	if opts.PrintHeader {
		if _, err := buf.WriteString(restoreUniqueAndForeignKeyCheckStmt); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

// SequenceDefinition renders one sequence as its mysqldump-style section.
func SequenceDefinition(sequence *metadata.SequenceMetadata) (string, error) {
	var buf strings.Builder
	if err := writeSequence(&buf, sequence); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// writeSequence writes sequence as TiDB's own SHOW CREATE SEQUENCE does: every
// option, the engine, and the comment.
func writeSequence(out *strings.Builder, sequence *metadata.SequenceMetadata) error {
	for _, s := range []string{emptyCommentLine, sequenceHeader, quoteIdent(sequence.Name) + "\n", emptyCommentLine} {
		if _, err := io.WriteString(out, s); err != nil {
			return err
		}
	}
	return writeCreateSequence(out, sequence)
}

func writeCreateSequence(out *strings.Builder, sequence *metadata.SequenceMetadata) error {
	fmt.Fprintf(out, "CREATE SEQUENCE %s", quoteIdent(sequence.Name))
	for _, opt := range []struct{ name, value string }{
		{"start with", sequence.Start},
		{"minvalue", sequence.MinValue},
		{"maxvalue", sequence.MaxValue},
		{"increment by", sequence.Increment},
	} {
		if opt.value != "" {
			fmt.Fprintf(out, " %s %s", opt.name, opt.value)
		}
	}
	switch sequence.CacheSize {
	case "0":
		out.WriteString(" nocache")
	case "":
	default:
		fmt.Fprintf(out, " cache %s", sequence.CacheSize)
	}
	if sequence.Cycle {
		out.WriteString(" cycle")
	} else {
		out.WriteString(" nocycle")
	}
	out.WriteString(" ENGINE=InnoDB")
	if sequence.Comment != "" {
		fmt.Fprintf(out, " COMMENT='%s'", outputFormat(sequence.Comment))
	}
	_, err := out.WriteString(";\n\n")
	return err
}

// TableDefinition renders one table as its mysqldump-style section.
func TableDefinition(table *metadata.TableMetadata) (string, error) {
	var buf strings.Builder
	if err := writeTable(&buf, table, ""); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ViewDefinition renders one view as its mysqldump-style section.
func ViewDefinition(view *metadata.ViewMetadata) (string, error) {
	var buf strings.Builder
	if err := writeView(&buf, view); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FunctionDefinition renders one stored function as its mysqldump-style section.
func FunctionDefinition(function *metadata.FunctionMetadata) (string, error) {
	var buf strings.Builder
	if err := writeFunction(&buf, function); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ProcedureDefinition renders one stored procedure as its mysqldump-style section.
func ProcedureDefinition(procedure *metadata.ProcedureMetadata) (string, error) {
	var buf strings.Builder
	if err := writeProcedure(&buf, procedure); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeEvent(out io.Writer, event *metadata.EventMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, eventHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, event.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, sql mode and timezone when present.
	if err := writeSessionContextIfSet(out, event.CharacterSetClient, event.CollationConnection, event.SqlMode, event.TimeZone, true); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterDoubleSemi); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, event.Definition); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ";;\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterSemi); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n")
	return err
}

func writeTrigger(out io.Writer, tableName string, trigger *metadata.TriggerMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, triggerHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, quoteIdent(trigger.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, and sql mode when present.
	if err := writeSessionContextIfSet(out, trigger.CharacterSetClient, trigger.CollationConnection, trigger.SqlMode, "", false); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterDoubleSemi); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, "CREATE TRIGGER "+quoteIdent(trigger.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.Timing); err != nil {
		return err
	}
	if _, err := io.WriteString(out, " "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.Event); err != nil {
		return err
	}
	if _, err := io.WriteString(out, " ON "+quoteIdent(tableName)+" FOR EACH ROW\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.Body); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ";;\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterSemi); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n")
	return err
}

func writeProcedure(out io.Writer, procedure *metadata.ProcedureMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, procedureHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, procedure.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, and sql mode when present.
	if err := writeSessionContextIfSet(out, procedure.CharacterSetClient, procedure.CollationConnection, procedure.SqlMode, "", false); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterDoubleSemi); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, procedure.Definition); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ";;\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterSemi); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n")
	return err
}

func writeFunction(out io.Writer, function *metadata.FunctionMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, functionHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, function.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, and sql mode when present.
	if err := writeSessionContextIfSet(out, function.CharacterSetClient, function.CollationConnection, function.SqlMode, "", false); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterDoubleSemi); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, function.Definition); err != nil {
		return err
	}

	if _, err := io.WriteString(out, ";;\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, delimiterSemi); err != nil {
		return err
	}

	_, err := io.WriteString(out, "\n")
	return err
}

func writeView(out io.Writer, view *metadata.ViewMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, viewHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Drop temporary view.
	if _, err := io.WriteString(out, "DROP VIEW IF EXISTS `"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`;\n"); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, "CREATE"); err != nil {
		return err
	}
	if view.Algorithm != "" {
		if _, err := fmt.Fprintf(out, " ALGORITHM=%s", strings.ToUpper(view.Algorithm)); err != nil {
			return err
		}
	}
	if view.SqlSecurity != "" {
		if _, err := fmt.Fprintf(out, " SQL SECURITY %s", strings.ToUpper(view.SqlSecurity)); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, " VIEW `"+view.Name+"`"+viewColumnListClause(view)+" AS "+view.Definition); err != nil {
		return err
	}
	if view.CheckOption != "" && !strings.HasSuffix(strings.ToUpper(strings.TrimSpace(view.Definition)), "CHECK OPTION") {
		if _, err := fmt.Fprintf(out, " WITH %s CHECK OPTION", strings.ToUpper(view.CheckOption)); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

// viewColumnListClause returns the column list a CREATE VIEW must state for its
// body to yield the view's recorded columns, as " (`a`,`b`)", or "" when the
// body names them itself. MySQL and TiDB store a body under the select's own
// names and keep a declared list beside it (see nameViewColumns); MariaDB and
// OceanBase fold the list into the body's aliases when they create the view, so
// a snapshot they synced never needs one.
func viewColumnListClause(view *metadata.ViewMetadata) string {
	names := make([]string, 0, len(view.Columns))
	for _, column := range view.Columns {
		names = append(names, column.GetName())
	}
	if len(names) == 0 {
		return ""
	}
	if bodyNames, ok := viewBodyColumnNames(view.GetDefinition()); ok && slices.Equal(bodyNames, names) {
		return ""
	}
	quoted := make([]string, len(names))
	for i, name := range names {
		quoted[i] = quoteIdent(name)
	}
	return " (" + strings.Join(quoted, ",") + ")"
}

// viewBodyColumnNames lists the names a view's body gives its columns, in order:
// each select item's alias, else the column it selects. ok is false when the
// body does not say — it does not parse, an item is a star or an unaliased
// expression, or it is no select — and its caller must state the names instead.
func viewBodyColumnNames(body string) ([]string, bool) {
	// A stored body may end in WITH CHECK OPTION, so it parses only as a whole view.
	stmt, err := parseFirstStatement[*nodes.CreateViewStmt]("CREATE VIEW v AS " + body)
	if err != nil || stmt.Select == nil {
		return nil, false
	}
	// A set operation takes its columns from its left branch.
	sel := stmt.Select
	for sel.SetOp != 0 && sel.Left != nil {
		sel = sel.Left
	}
	names := make([]string, 0, len(sel.TargetList))
	for _, item := range sel.TargetList {
		// Only an aliased item is a ResTarget; a bare one is its expression.
		switch target := item.(type) {
		case *nodes.ResTarget:
			names = append(names, target.Name)
		case *nodes.ColumnRef:
			if target.Star {
				return nil, false
			}
			names = append(names, target.Column)
		default:
			return nil, false
		}
	}
	return names, true
}

func writeTable(out *strings.Builder, table *metadata.TableMetadata, dbName string) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tableHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, table.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Definition.
	createPrefix := "CREATE TABLE"
	if table.IsGlobalTemporary {
		createPrefix = "CREATE GLOBAL TEMPORARY TABLE"
	}
	if _, err := fmt.Fprintf(out, "%s `%s` (\n", createPrefix, table.Name); err != nil {
		return err
	}

	for i, column := range table.Columns {
		if i != 0 {
			if _, err := fmt.Fprint(out, ",\n"); err != nil {
				return err
			}
		}
		if err := printColumnClause(out, column, table); err != nil {
			return err
		}
	}

	if err := printPrimaryKeyClause(out, table); err != nil {
		return err
	}

	for _, index := range table.Indexes {
		if index.Primary {
			continue
		}
		if err := printIndexClause(out, index); err != nil {
			return err
		}
	}

	for _, fk := range table.ForeignKeys {
		if err := printForeignKeyClause(out, fk, dbName); err != nil {
			return err
		}
	}

	for _, check := range table.CheckConstraints {
		if err := printCheckClause(out, check); err != nil {
			return err
		}
	}

	if _, err := out.WriteString("\n)"); err != nil {
		return err
	}
	if table.Engine != "" {
		if _, err := fmt.Fprintf(out, " ENGINE=%s", table.Engine); err != nil {
			return err
		}
	}

	if table.Charset != "" {
		if _, err := fmt.Fprintf(out, " DEFAULT CHARSET=%s", table.Charset); err != nil {
			return err
		}
	}

	if table.Collation != "" {
		if _, err := fmt.Fprintf(out, " COLLATE=%s", table.Collation); err != nil {
			return err
		}
	}

	if table.AutoIdCache != 0 {
		if _, err := fmt.Fprintf(out, " /*T![auto_id_cache] AUTO_ID_CACHE=%d */", table.AutoIdCache); err != nil {
			return err
		}
	}

	if strings.HasPrefix(table.ShardingInfo, "SHARD_BITS=") {
		shardBits := strings.ReplaceAll(table.ShardingInfo, "SHARD_BITS", "SHARD_ROW_ID_BITS")
		if table.PreSplitRegions > 0 {
			if _, err := fmt.Fprintf(out, " /*T! %s PRE_SPLIT_REGIONS=%d */", shardBits, table.PreSplitRegions); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(out, " /*T! %s */", shardBits); err != nil {
				return err
			}
		}
	} else if table.PreSplitRegions > 0 {
		if _, err := fmt.Fprintf(out, " /*T! PRE_SPLIT_REGIONS=%d */", table.PreSplitRegions); err != nil {
			return err
		}
	}

	if comp := unquoteTableOption(createOption(table.CreateOptions, "compression")); comp != "" {
		if _, err := fmt.Fprintf(out, " COMPRESSION='%s'", escapeComment(comp)); err != nil {
			return err
		}
	}

	if table.Comment != "" {
		if _, err := fmt.Fprintf(out, " COMMENT='%s'", table.Comment); err != nil {
			return err
		}
	}

	if table.PlacementPolicy != "" {
		if _, err := fmt.Fprintf(out, " /*T![placement] PLACEMENT POLICY=%s */", quoteIdent(table.PlacementPolicy)); err != nil {
			return err
		}
	}

	if table.TtlColumn != "" {
		enable := "OFF"
		if table.TtlEnable {
			enable = "ON"
		}
		job := table.TtlJobInterval
		if job == "" {
			job = "1h"
		}
		if _, err := fmt.Fprintf(out, " /*T![ttl] TTL=%s + %s */ /*T![ttl] TTL_ENABLE='%s' */ /*T![ttl] TTL_JOB_INTERVAL='%s' */",
			quoteIdent(table.TtlColumn), ttlInterval(table.TtlInterval), enable, job); err != nil {
			return err
		}
	}

	if len(table.Partitions) > 0 {
		if err := printPartitionClause(out, table.Partitions); err != nil {
			return err
		}
	}

	if table.OnCommit != "" {
		if _, err := fmt.Fprintf(out, " ON COMMIT %s", table.OnCommit); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, ";\n\n"); err != nil {
		return err
	}

	for _, col := range table.Columns {
		if col.GetDefaultDropped() {
			if _, err := fmt.Fprintf(out, "ALTER TABLE `%s` ALTER COLUMN `%s` DROP DEFAULT;\n\n", table.Name, col.Name); err != nil {
				return err
			}
		}
	}

	// TiFlash replica precedes CACHE because a cached table rejects other ALTER TABLE statements (8242).
	if table.TiflashReplica > 0 || len(table.TiflashLocationLabels) > 0 {
		if _, err := fmt.Fprintf(out, "ALTER TABLE `%s` SET TIFLASH REPLICA %d", table.Name, table.TiflashReplica); err != nil {
			return err
		}
		if err := writeTiFlashLocationLabels(out, table.TiflashLocationLabels); err != nil {
			return err
		}
		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}
	}

	if isCachedCreateOptions(table.CreateOptions) {
		if _, err := fmt.Fprintf(out, "ALTER TABLE `%s` CACHE;\n\n", table.Name); err != nil {
			return err
		}
	}

	return nil
}

func writeDatabaseOptions(out *strings.Builder, meta *metadata.DatabaseSchemaMetadata) error {
	target := ""
	if meta.Name != "" {
		target = " " + quoteIdent(meta.Name)
	}
	if meta.PlacementPolicy != "" {
		if _, err := fmt.Fprintf(out, "ALTER DATABASE%s PLACEMENT POLICY=%s;\n\n", target, quoteIdent(meta.PlacementPolicy)); err != nil {
			return err
		}
	}
	if meta.TiflashReplica > 0 || len(meta.TiflashLocationLabels) > 0 {
		if _, err := fmt.Fprintf(out, "ALTER DATABASE%s SET TIFLASH REPLICA %d", target, meta.TiflashReplica); err != nil {
			return err
		}
		if err := writeTiFlashLocationLabels(out, meta.TiflashLocationLabels); err != nil {
			return err
		}
		if _, err := io.WriteString(out, ";\n\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeTiFlashLocationLabels(out io.Writer, locationLabels []string) error {
	if len(locationLabels) == 0 {
		return nil
	}
	labels := make([]string, len(locationLabels))
	for i, l := range locationLabels {
		labels[i] = "'" + outputFormat(l) + "'"
	}
	_, err := fmt.Fprintf(out, " LOCATION LABELS %s", strings.Join(labels, ", "))
	return err
}

func writeSessionContextIfSet(out io.Writer, charClient, collation, sqlMode, timeZone string, includeTimezone bool) error {
	for _, stmt := range []struct {
		prefix string
		value  string
		quote  bool
	}{
		{setCharacterSetClient, charClient, false},
		{setCharacterSetResult, charClient, false},
		{setCollation, collation, false},
		{setSQLMode, sqlMode, true},
		{setTimezone, timeZone, true},
	} {
		if stmt.prefix == setTimezone && !includeTimezone {
			continue
		}
		if stmt.value == "" {
			continue
		}
		if stmt.quote {
			if _, err := fmt.Fprintf(out, "%s'%s';\n", stmt.prefix, escapeComment(stmt.value)); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(out, "%s%s;\n", stmt.prefix, stmt.value); err != nil {
				return err
			}
		}
	}
	return nil
}

// printPartitionClause writes the PARTITION BY clause of a CREATE TABLE.
func printPartitionClause(buf *strings.Builder, partitions []*metadata.TablePartitionMetadata) error {
	if len(partitions) == 0 {
		return nil
	}
	vsc := getVersionSpecificComment(partitions)
	curComment := vsc
	if _, err := fmt.Fprintf(buf, "%s PARTITION BY ", curComment); err != nil {
		return err
	}
	switch partitions[0].Type {
	case metadata.TablePartitionMetadata_RANGE:
		if _, err := fmt.Fprintf(buf, "RANGE (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_RANGE_COLUMNS:
		if _, err := fmt.Fprintf(buf, "RANGE COLUMNS (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_LIST:
		if _, err := fmt.Fprintf(buf, "LIST (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_LIST_COLUMNS:
		if _, err := fmt.Fprintf(buf, "LIST COLUMNS (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_HASH:
		if _, err := fmt.Fprintf(buf, "HASH (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_KEY:
		if _, err := fmt.Fprintf(buf, "KEY%s (%s)", keyAlgorithmClause(partitions[0].Algorithm), partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_LINEAR_HASH:
		if _, err := fmt.Fprintf(buf, "LINEAR HASH (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_LINEAR_KEY:
		if _, err := fmt.Fprintf(buf, "LINEAR KEY%s (%s)", keyAlgorithmClause(partitions[0].Algorithm), partitions[0].Expression); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown partition type: %v", partitions[0].Type)
	}

	useDefault := int64(0)
	if partitions[0].UseDefault != "" {
		var err error
		useDefault, err = strconv.ParseInt(partitions[0].UseDefault, 10, 64)
		if err != nil {
			return err
		}
	}
	if useDefault != 0 {
		if _, err := fmt.Fprintf(buf, "\nPARTITIONS %d", useDefault); err != nil {
			return err
		}
	}
	if partitions[0].IntervalClause != "" {
		if _, err := fmt.Fprintf(buf, " %s", partitions[0].IntervalClause); err != nil {
			return err
		}
	}

	if len(partitions[0].Subpartitions) > 0 {
		if _, err := fmt.Fprint(buf, "\nSUBPARTITION BY "); err != nil {
			return err
		}
		switch partitions[0].Subpartitions[0].Type {
		case metadata.TablePartitionMetadata_HASH:
			if _, err := fmt.Fprintf(buf, "HASH (%s)", partitions[0].Subpartitions[0].Expression); err != nil {
				return err
			}
		case metadata.TablePartitionMetadata_LINEAR_HASH:
			if _, err := fmt.Fprintf(buf, "LINEAR HASH (%s)", partitions[0].Subpartitions[0].Expression); err != nil {
				return err
			}
		case metadata.TablePartitionMetadata_KEY:
			if _, err := fmt.Fprintf(buf, "KEY%s (%s)", keyAlgorithmClause(partitions[0].Subpartitions[0].Algorithm), partitions[0].Subpartitions[0].Expression); err != nil {
				return err
			}
		case metadata.TablePartitionMetadata_LINEAR_KEY:
			if _, err := fmt.Fprintf(buf, "LINEAR KEY%s (%s)", keyAlgorithmClause(partitions[0].Subpartitions[0].Algorithm), partitions[0].Subpartitions[0].Expression); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid subpartition type: %v", partitions[0].Subpartitions[0].Type)
		}
	}

	subUseDefault := 0
	if len(partitions[0].Subpartitions) > 0 && partitions[0].Subpartitions[0].UseDefault != "" {
		var err error
		subUseDefault, err = strconv.Atoi(partitions[0].Subpartitions[0].UseDefault)
		if err != nil {
			return err
		}
	}

	if subUseDefault != 0 {
		if _, err := fmt.Fprintf(buf, "\nSUBPARTITIONS %d", subUseDefault); err != nil {
			return err
		}
	}

	if useDefault == 0 && !(len(partitions) == 1 && partitions[0].Name == "" && partitions[0].IntervalClause != "") {
		if _, err := fmt.Fprint(buf, "\n("); err != nil {
			return err
		}
		preposition, err := getPrepositionByType(partitions[0].Type)
		if err != nil {
			return err
		}
		for i, partition := range partitions {
			if i != 0 {
				if _, err := fmt.Fprint(buf, ",\n "); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(buf, "PARTITION %s", partition.Name); err != nil {
				return err
			}
			if preposition != "" {
				if partition.Value != "MAXVALUE" {
					if _, err := fmt.Fprintf(buf, " VALUES %s (%s)", preposition, partition.Value); err != nil {
						return err
					}
				} else {
					if _, err := fmt.Fprintf(buf, " VALUES %s %s", preposition, partition.Value); err != nil {
						return err
					}
				}
			}

			if err := writePartitionOptions(buf, partition.Comment); err != nil {
				return err
			}

			if subUseDefault == 0 && len(partition.Subpartitions) > 0 {
				if _, err := fmt.Fprint(buf, "\n ("); err != nil {
					return err
				}
				for j, subPartition := range partition.Subpartitions {
					if _, err := fmt.Fprintf(buf, "SUBPARTITION %s", subPartition.Name); err != nil {
						return err
					}
					if err := writePartitionOptions(buf, subPartition.Comment); err != nil {
						return err
					}
					if j == len(partition.Subpartitions)-1 {
						if _, err := fmt.Fprint(buf, ")"); err != nil {
							return err
						}
					} else {
						if _, err := fmt.Fprint(buf, ",\n  "); err != nil {
							return err
						}
					}
				}
			}

			if i == len(partitions)-1 {
				if _, err := fmt.Fprint(buf, ")"); err != nil {
					return err
				}
			}
		}
	}

	if _, err := fmt.Fprint(buf, " */"); err != nil {
		return err
	}

	return nil
}

func getPrepositionByType(tp metadata.TablePartitionMetadata_Type) (string, error) {
	switch tp {
	case metadata.TablePartitionMetadata_RANGE:
		return "LESS THAN", nil
	case metadata.TablePartitionMetadata_RANGE_COLUMNS:
		return "LESS THAN", nil
	case metadata.TablePartitionMetadata_LIST:
		return "IN", nil
	case metadata.TablePartitionMetadata_LIST_COLUMNS:
		return "IN", nil
	case metadata.TablePartitionMetadata_HASH, metadata.TablePartitionMetadata_KEY, metadata.TablePartitionMetadata_LINEAR_HASH, metadata.TablePartitionMetadata_LINEAR_KEY:
		return "", nil
	default:
		return "", fmt.Errorf("unsupported partition type: %v", tp)
	}
}

func writePartitionOptions(buf io.StringWriter, comment string) error {
	/*
		int err = 0;
		err += add_space(fptr);
		if (p_elem->tablespace_name) {
			err += add_string(fptr, "TABLESPACE = ");
			err += add_ident_string(fptr, p_elem->tablespace_name);
			err += add_space(fptr);
		}
		if (p_elem->nodegroup_id != UNDEF_NODEGROUP)
			err += add_keyword_int(fptr, "NODEGROUP", (longlong)p_elem->nodegroup_id);
		if (p_elem->part_max_rows)
			err += add_keyword_int(fptr, "MAX_ROWS", (longlong)p_elem->part_max_rows);
		if (p_elem->part_min_rows)
			err += add_keyword_int(fptr, "MIN_ROWS", (longlong)p_elem->part_min_rows);
		if (!(current_thd->variables.sql_mode & MODE_NO_DIR_IN_CREATE)) {
			if (p_elem->data_file_name)
			err += add_keyword_path(fptr, "DATA DIRECTORY", p_elem->data_file_name);
			if (p_elem->index_file_name)
			err += add_keyword_path(fptr, "INDEX DIRECTORY", p_elem->index_file_name);
		}
		if (p_elem->part_comment)
			err += add_keyword_string(fptr, "COMMENT", true, p_elem->part_comment);
		return err + add_engine(fptr, p_elem->engine_type);
	*/
	// TODO(zp): Get all the partition options from the metadata is too complex, just write ENGINE=InnoDB for now.
	if _, err := buf.WriteString(" ENGINE=InnoDB"); err != nil {
		return err
	}
	if comment != "" {
		if _, err := buf.WriteString(" COMMENT '" + escapeComment(comment) + "'"); err != nil {
			return err
		}
	}

	return nil
}

func getVersionSpecificComment(partitions []*metadata.TablePartitionMetadata) string {
	if len(partitions) == 0 {
		return ""
	}
	partition := partitions[0]
	if partition.Type == metadata.TablePartitionMetadata_RANGE_COLUMNS || partition.Type == metadata.TablePartitionMetadata_LIST_COLUMNS {
		// MySQL introduce columns partitioning in 5.5+
		return "\n/*!50500"
	}
	return "\n/*!50100"
}

func keyAlgorithmClause(algorithm int32) string {
	if algorithm > 0 {
		return fmt.Sprintf(" ALGORITHM = %d", algorithm)
	}
	return ""
}

func printCheckClause(buf *strings.Builder, check *metadata.CheckConstraintMetadata) error {
	if _, err := fmt.Fprintf(buf, ",\n  CONSTRAINT `%s` CHECK %s", check.Name, check.Expression); err != nil {
		return err
	}
	if check.NotEnforced {
		if _, err := buf.WriteString(" NOT ENFORCED"); err != nil {
			return err
		}
	}
	return nil
}

func printForeignKeyClause(buf *strings.Builder, fk *metadata.ForeignKeyMetadata, dbName string) error {
	if _, err := fmt.Fprintf(buf, ",\n  CONSTRAINT `%s` FOREIGN KEY (", fk.Name); err != nil {
		return err
	}

	for i, column := range fk.Columns {
		if i != 0 {
			if _, err := fmt.Fprint(buf, ", "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(buf, "`%s`", column); err != nil {
			return err
		}
	}

	referencedTable := quoteIdent(fk.ReferencedTable)
	if fk.ReferencedSchema != "" && fk.ReferencedSchema != dbName {
		referencedTable = quoteIdent(fk.ReferencedSchema) + "." + referencedTable
	}
	if _, err := fmt.Fprintf(buf, ") REFERENCES %s (", referencedTable); err != nil {
		return err
	}

	for i, column := range fk.ReferencedColumns {
		if i != 0 {
			if _, err := fmt.Fprint(buf, ", "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(buf, "`%s`", column); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(buf, ")"); err != nil {
		return err
	}

	if fk.OnDelete != "" && !strings.EqualFold(fk.OnDelete, mysqlNoAction) {
		if _, err := fmt.Fprintf(buf, " ON DELETE %s", fk.OnDelete); err != nil {
			return err
		}
	}

	if fk.OnUpdate != "" && !strings.EqualFold(fk.OnUpdate, mysqlNoAction) {
		if _, err := fmt.Fprintf(buf, " ON UPDATE %s", fk.OnUpdate); err != nil {
			return err
		}
	}

	return nil
}

func printIndexClause(buf *strings.Builder, index *metadata.IndexMetadata) error {
	if index.Primary {
		return nil
	}

	if _, err := fmt.Fprint(buf, ",\n  "); err != nil {
		return err
	}

	switch {
	case index.Unique:
		if _, err := fmt.Fprintf(buf, "UNIQUE KEY `%s` (", index.Name); err != nil {
			return err
		}
	case index.IsVector || strings.EqualFold(index.Type, "HNSW"):
		if _, err := fmt.Fprintf(buf, "VECTOR INDEX `%s` (", index.Name); err != nil {
			return err
		}
	case strings.EqualFold(index.Type, mysqlIndexFullText):
		if _, err := fmt.Fprintf(buf, "FULLTEXT KEY `%s` (", index.Name); err != nil {
			return err
		}
	case strings.EqualFold(index.Type, mysqlIndexSpatial):
		if _, err := fmt.Fprintf(buf, "SPATIAL KEY `%s` (", index.Name); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprintf(buf, "KEY `%s` (", index.Name); err != nil {
			return err
		}
	}

	for i, expr := range index.Expressions {
		if i != 0 {
			if _, err := fmt.Fprint(buf, ", "); err != nil {
				return err
			}
		}
		keyLength := int64(-1)
		descending := false
		if len(index.KeyLength) > i {
			keyLength = index.KeyLength[i]
		}
		if len(index.Descending) > i {
			descending = index.Descending[i]
		}
		if err := printIndexKeyPart(buf, expr, keyLength, descending); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(buf, ")"); err != nil {
		return err
	}

	if index.KeyBlockSize > 0 {
		if _, err := fmt.Fprintf(buf, " KEY_BLOCK_SIZE=%d", index.KeyBlockSize); err != nil {
			return err
		}
	}

	if index.Comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT '%s'", escapeComment(index.Comment)); err != nil {
			return err
		}
	}

	if index.IsGlobal {
		if _, err := fmt.Fprint(buf, " /*T![global_index] GLOBAL */"); err != nil {
			return err
		}
	}

	return nil
}

func printIndexKeyPart(buf *strings.Builder, expr string, length int64, descending bool) error {
	if len(expr) == 0 {
		return errors.New("index key part expression is empty")
	}
	if expr[0] == '(' && expr[len(expr)-1] == ')' {
		if _, err := buf.WriteString(expr); err != nil {
			return err
		}
	} else {
		if _, err := buf.WriteString("`"); err != nil {
			return err
		}
		if _, err := buf.WriteString(expr); err != nil {
			return err
		}
		if _, err := buf.WriteString("`"); err != nil {
			return err
		}
	}
	if length > 0 {
		if _, err := buf.WriteString("("); err != nil {
			return err
		}
		if _, err := buf.WriteString(strconv.FormatInt(length, 10)); err != nil {
			return err
		}
		if _, err := buf.WriteString(")"); err != nil {
			return err
		}
	}
	if descending {
		if _, err := buf.WriteString(" DESC"); err != nil {
			return err
		}
	}
	return nil
}

func printPrimaryKeyClause(buf *strings.Builder, table *metadata.TableMetadata) error {
	for _, index := range table.GetIndexes() {
		if index.Primary {
			if _, err := fmt.Fprint(buf, ",\n  PRIMARY KEY ("); err != nil {
				return err
			}
			for i, column := range index.Expressions {
				if i != 0 {
					if _, err := fmt.Fprint(buf, ", "); err != nil {
						return err
					}
				}
				keyLength := int64(-1)
				descending := false
				if len(index.KeyLength) > i {
					keyLength = index.KeyLength[i]
				}
				if len(index.Descending) > i {
					descending = index.Descending[i]
				}
				if err := printIndexKeyPart(buf, column, keyLength, descending); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(buf, ")"); err != nil {
				return err
			}
			if index.KeyBlockSize > 0 {
				if _, err := fmt.Fprintf(buf, " KEY_BLOCK_SIZE=%d", index.KeyBlockSize); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintf(buf, " /*T![clustered_index] %s */", table.PrimaryKeyType); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func isAutoIncrement(column *metadata.ColumnMetadata) bool {
	return strings.EqualFold(column.GetDefault(), autoIncrementSymbol)
}

func isAutoRandom(column *metadata.ColumnMetadata) bool {
	return strings.HasPrefix(column.GetDefault(), autoRandomSymbol)
}

func printColumnClause(buf *strings.Builder, column *metadata.ColumnMetadata, table *metadata.TableMetadata) error {
	if _, err := fmt.Fprintf(buf, "  `%s` %s", column.Name, column.Type); err != nil {
		return err
	}

	if column.CharacterSet != "" && column.CharacterSet != table.Charset {
		if _, err := fmt.Fprintf(buf, " CHARACTER SET %s", column.CharacterSet); err != nil {
			return err
		}
	}

	if column.Collation != "" && column.Collation != table.Collation {
		if _, err := fmt.Fprintf(buf, " COLLATE %s", column.Collation); err != nil {
			return err
		}
	}

	if column.Generation != nil && column.Generation.Expression != "" {
		if _, err := fmt.Fprintf(buf, " GENERATED ALWAYS AS (%s) ", column.Generation.Expression); err != nil {
			return err
		}
		switch column.Generation.Type {
		case metadata.GenerationMetadata_TYPE_STORED:
			if _, err := fmt.Fprint(buf, "STORED"); err != nil {
				return err
			}
		case metadata.GenerationMetadata_TYPE_VIRTUAL:
			if _, err := fmt.Fprint(buf, "VIRTUAL"); err != nil {
				return err
			}
		default:
			// No special handling for other generation types
		}
	}

	if !column.Nullable {
		if _, err := fmt.Fprint(buf, " NOT NULL"); err != nil {
			return err
		}
	}

	switch {
	case isAutoIncrement(column):
		if _, err := fmt.Fprintf(buf, " %s", autoIncrementSymbol); err != nil {
			return err
		}
	case isAutoRandom(column):
		if _, err := fmt.Fprintf(buf, " /*T![auto_rand] %s */", column.GetDefault()); err != nil {
			return err
		}
	default:
		if err := printDefaultClause(buf, column); err != nil {
			return err
		}
	}

	if column.OnUpdate != "" {
		if _, err := fmt.Fprintf(buf, " ON UPDATE %s", column.OnUpdate); err != nil {
			return err
		}
	}
	if column.Comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT '%s'", column.Comment); err != nil {
			return err
		}
	}
	return nil
}

func printDefaultClause(buf *strings.Builder, column *metadata.ColumnMetadata) error {
	// Check if column has any default value
	hasDefault := column.Default != ""
	if !hasDefault {
		return nil
	}

	if column.Default == "NULL" {
		if !column.Nullable || !typeSupportsDefaultValue(column.Type) {
			// If the column is not nullable, then the default value should not be null.
			// For this case, we should not print the default clause.
			return nil
		}
		if column.Generation != nil && column.Generation.Expression != "" {
			return nil
		}
		if _, err := fmt.Fprint(buf, " DEFAULT NULL"); err != nil {
			return err
		}
		return nil
	}

	if column.Default != "" {
		if isAutoIncrement(column) {
			// If the default value is auto_increment, then we should not print the default clause.
			// We'll handle this in the following AUTO_INCREMENT clause.
			return nil
		}

		if isAutoRandom(column) {
			// If the default value is auto_random, then we should not print the default clause
			// We'll handle this in the following AUTO_RANDOM clause.
			return nil
		}

		if _, err := fmt.Fprintf(buf, " DEFAULT %s", column.Default); err != nil {
			return err
		}
	}

	return nil
}

func typeSupportsDefaultValue(tp string) bool {
	switch strings.ToLower(tp) {
	case mysqlTypeBlob, mysqlTypeTinyBob, mysqlTypeMediumBlob, mysqlTypeLongBlob, mysqlTypeJSON, mysqlTypeGeometry:
		return false
	default:
		return true
	}
}

func writeInvalidTemporaryView(out io.Writer, view *metadata.ViewMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tempViewHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "-- `"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	_, err := io.WriteString(out, "` references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them\n\n")
	return err
}

func writeTemporaryView(out io.Writer, view *metadata.ViewMetadata) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tempViewHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "`\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, "CREATE VIEW `"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "` AS SELECT\n  "); err != nil {
		return err
	}
	for i, column := range view.Columns {
		if i != 0 {
			if _, err := io.WriteString(out, ",\n  "); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, "1 AS `"); err != nil {
			return err
		}
		if _, err := io.WriteString(out, column.Name); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "`"); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}
