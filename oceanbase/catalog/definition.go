package catalog

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/oceanbase/ast"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

const (
	disableUniqueAndForeignKeyCheckStmt = "SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;\nSET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;\n"
	restoreUniqueAndForeignKeyCheckStmt = "SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;\nSET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;\n"

	emptyCommentLine       = "--\n"
	tempViewHeader         = "-- Temporary view structure for "
	tableHeader            = "-- Table structure for "
	viewHeader             = "-- View structure for "
	materializedViewHeader = "-- Materialized view structure for "
	functionHeader         = "-- Function structure for "
	procedureHeader        = "-- Procedure structure for "
	triggerHeader          = "-- Trigger structure for "
	eventHeader            = "-- Event structure for "
	setCharacterSetClient  = "SET character_set_client = "
	setCharacterSetResult  = "SET character_set_results = "
	setCollation           = "SET collation_connection = "
	setSQLMode             = "SET sql_mode = "
	setTimezone            = "SET time_zone = "
	delimiterDoubleSemi    = "DELIMITER ;;\n"
	delimiterSemi          = "DELIMITER ;\n"

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
)

// DefinitionOptions selects the form DatabaseDefinition renders.
type DefinitionOptions struct {
	// PrintHeader wraps the mysqldump-style output in statements that disable and
	// restore unique and foreign key checks. SDL output ignores it.
	PrintHeader bool
	// SDLFormat renders the declarative form LoadSDL reads instead of the
	// mysqldump-style restore script.
	SDLFormat bool
}

// DatabaseDefinition renders a synced database snapshot as DDL: a mysqldump-style
// restore script, or with opts.SDLFormat the canonical SDL that LoadSDL reads back.
// It renders only the first schema and returns "" for a snapshot with none. Column and
// table comments are written between quotes as they are: the OceanBase sync stores them
// escaped, as QUOTE() returns them without the outer quotes.
func DatabaseDefinition(meta *metadata.DatabaseSchemaMetadata, opts DefinitionOptions) (string, error) {
	if len(meta.Schemas) == 0 {
		return "", nil
	}

	if opts.SDLFormat {
		return getSDLFormat(meta)
	}

	var buf strings.Builder

	if opts.PrintHeader {
		if _, err := buf.WriteString(disableUniqueAndForeignKeyCheckStmt); err != nil {
			return "", err
		}
	}

	schema := meta.Schemas[0]

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

	// Construct tables.
	for _, table := range schema.Tables {
		if err := writeTable(&buf, table, meta.Name); err != nil {
			return "", err
		}
	}

	// Construct external tables after the tables they may read.
	for _, externalTable := range schema.ExternalTables {
		if err := writeExternalTable(&buf, externalTable); err != nil {
			return "", err
		}
	}

	// Construct functions before views: CREATE VIEW checks that the functions it calls
	// exist, while a function body is checked only when it runs.
	for _, function := range schema.Functions {
		if err := writeFunction(&buf, function); err != nil {
			return "", err
		}
	}

	// Construct views and materialized views, each after the ones it reads.
	if err := writeViewsInDependencyOrder(&buf, meta.Name, schema); err != nil {
		return "", err
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

// writeMaterializedView writes a materialized view and the indexes its container table
// carries. OceanBase answers no SHOW CREATE MATERIALIZED VIEW, so the synced query and the
// container table's indexes are the only statements of the view; a refresh policy has no
// field in the metadata, and a plain CREATE gives the server's own default.
func writeMaterializedView(out *strings.Builder, view *metadata.MaterializedViewMetadata) error {
	if view.SkipDump || strings.TrimSpace(view.Definition) == "" {
		return nil
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, materializedViewHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, quoteIdent(view.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// The comment arrives QUOTE()-escaped, so a backslash in it is an escape character: say
	// so, because a stored routine's section before this one may have left a mode where it is
	// not.
	if _, err := io.WriteString(out, "SET sql_mode = DEFAULT;\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "CREATE MATERIALIZED VIEW %s", quoteIdent(view.Name)); err != nil {
		return err
	}
	for _, index := range view.Indexes {
		if !index.Primary && !strings.EqualFold(index.Name, "PRIMARY") {
			continue
		}
		if len(index.Expressions) == 0 {
			break
		}
		cols := make([]string, len(index.Expressions))
		for i, col := range index.Expressions {
			cols[i] = quoteIdent(col)
		}
		if _, err := fmt.Fprintf(out, " (PRIMARY KEY (%s)%s", strings.Join(cols, ", "), indexUsingClause(index.Type, "")); err != nil {
			return err
		}
		if index.Comment != "" {
			if _, err := fmt.Fprintf(out, " COMMENT '%s'", escapeComment(index.Comment)); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, ")"); err != nil {
			return err
		}
		break
	}
	if view.Comment != "" {
		if _, err := fmt.Fprintf(out, " COMMENT '%s'", view.Comment); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, " AS %s;\n", strings.TrimSpace(view.Definition)); err != nil {
		return err
	}

	for _, index := range view.Indexes {
		if err := writeIndex(out, view.Name, index); err != nil {
			return err
		}
	}

	_, err := io.WriteString(out, "\n")
	return err
}

// writeIndex writes one standalone CREATE INDEX statement, which a materialized view needs
// for an index its container table carries: a view has no CREATE TABLE to hold it.
func writeIndex(out *strings.Builder, tableName string, index *metadata.IndexMetadata) error {
	if index.Primary {
		return nil
	}
	if _, err := fmt.Fprintf(out, "\nCREATE %sINDEX %s ON %s (", indexQualifier(index), quoteIdent(index.Name), quoteIdent(tableName)); err != nil {
		return err
	}
	if err := printIndexKeyParts(out, index, ""); err != nil {
		return err
	}
	_, err := io.WriteString(out, ";\n")
	return err
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
	// An event synced without its CREATE EVENT has nothing faithful to write.
	if event.Definition == "" {
		return nil
	}
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, eventHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, quoteIdent(event.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, sql mode and timezone. OceanBase states no character set or
	// collation for an event, so those lines are written only when the metadata carries
	// them; an empty sql_mode is a value the server states and is written as one.
	if err := writeAdditionalEventsIfSet(out, event.CharacterSetClient, event.CharacterSetClient, event.CollationConnection, ""); err != nil {
		return err
	}
	if _, err := io.WriteString(out, setSQLMode+"'"+event.SqlMode+"';\n"); err != nil {
		return err
	}
	// OceanBase refuses the literal SYSTEM as a session zone, and SYSTEM is the zone the
	// session starts in: give it back what an earlier event's explicit zone may have changed.
	switch {
	case strings.EqualFold(event.TimeZone, "SYSTEM"):
		if _, err := io.WriteString(out, setTimezone+"DEFAULT;\n"); err != nil {
			return err
		}
	case event.TimeZone != "":
		if _, err := io.WriteString(out, setTimezone+"'"+event.TimeZone+"';\n"); err != nil {
			return err
		}
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

	// Set charset, collation, and sql mode.
	if _, err := io.WriteString(out, setCharacterSetClient); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.CharacterSetClient); err != nil {
		return err
	}
	if _, err := io.WriteString(out, ";\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, setCharacterSetResult); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.CharacterSetClient); err != nil {
		return err
	}
	if _, err := io.WriteString(out, ";\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, setCollation); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.CollationConnection); err != nil {
		return err
	}
	if _, err := io.WriteString(out, ";\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, setSQLMode); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "'"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, trigger.SqlMode); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "';\n"); err != nil {
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
	if _, err := io.WriteString(out, quoteIdent(procedure.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, and sql mode.
	if err := writeAdditionalEventsIfSet(out, procedure.CharacterSetClient, procedure.CharacterSetClient, procedure.CollationConnection, procedure.SqlMode); err != nil {
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
	if _, err := io.WriteString(out, quoteIdent(function.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Set charset, collation, and sql mode.
	if err := writeAdditionalEventsIfSet(out, function.CharacterSetClient, function.CharacterSetClient, function.CollationConnection, function.SqlMode); err != nil {
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
	if _, err := io.WriteString(out, quoteIdent(view.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, "CREATE OR REPLACE"); err != nil {
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
	if _, err := io.WriteString(out, " VIEW "+quoteIdent(view.Name)+viewColumnListClause(view)+" AS "+view.Definition); err != nil {
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

// printExternalColumnClause writes one external table column as SHOW CREATE
// TABLE states it: the type, the file column it reads through a generation
// expression, nullability and comment. The engine stores a generated external
// column's DEFAULT and character set too, but its own DDL states neither.
func printExternalColumnClause(buf *strings.Builder, column *metadata.ColumnMetadata) error {
	if _, err := fmt.Fprintf(buf, "  %s %s", quoteIdent(column.Name), normalizeColumnType(column.Type)); err != nil {
		return err
	}
	if column.Generation != nil && column.Generation.Expression != "" {
		if _, err := fmt.Fprintf(buf, " GENERATED ALWAYS AS (%s)", column.Generation.Expression); err != nil {
			return err
		}
	}
	if !column.Nullable {
		if _, err := fmt.Fprint(buf, " NOT NULL"); err != nil {
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

// quoteExternalTableString escapes a value for an OceanBase string literal,
// which honors backslash escapes as well as doubled quotes.
func quoteExternalTableString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, "'", "''")
}

// writeExternalTable writes an external table: its columns and the LOCATION,
// PATTERN and FORMAT clauses SHOW CREATE TABLE states for it. A table the sync
// read without any still renders, so the columns are not lost.
func writeExternalTable(out *strings.Builder, table *metadata.ExternalTableMetadata) error {
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tableHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, quoteIdent(table.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "CREATE EXTERNAL TABLE %s (\n", quoteIdent(table.Name)); err != nil {
		return err
	}
	for i, column := range table.Columns {
		if i != 0 {
			if _, err := io.WriteString(out, ",\n"); err != nil {
				return err
			}
		}
		if err := printExternalColumnClause(out, column); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, "\n)"); err != nil {
		return err
	}
	if table.Location != "" {
		if _, err := fmt.Fprintf(out, "\nLOCATION = '%s'", quoteExternalTableString(table.Location)); err != nil {
			return err
		}
	}
	if table.Pattern != "" {
		if _, err := fmt.Fprintf(out, "\nPATTERN = '%s'", quoteExternalTableString(table.Pattern)); err != nil {
			return err
		}
	}
	if table.AutoRefresh != "" {
		if _, err := fmt.Fprintf(out, "\nAUTO_REFRESH = %s", table.AutoRefresh); err != nil {
			return err
		}
	}
	if table.Format != "" {
		if _, err := fmt.Fprintf(out, "\nFORMAT (%s)", table.Format); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

func writeTable(out *strings.Builder, table *metadata.TableMetadata, dbName string) error {
	// Header.
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	if _, err := io.WriteString(out, tableHeader); err != nil {
		return err
	}
	if _, err := io.WriteString(out, quoteIdent(table.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	return writeCreateTable(out, table, dbName, false)
}

// writeCreateTable writes the CREATE TABLE statement both output forms share. sorted
// orders secondary indexes, foreign keys, and checks by name, as the SDL form does;
// otherwise they keep their snapshot order. dbName is the database being dumped, empty
// when unknown.
func writeCreateTable(buf *strings.Builder, table *metadata.TableMetadata, dbName string, sorted bool) error {
	if _, err := fmt.Fprintf(buf, "CREATE TABLE %s (\n", quoteIdent(table.Name)); err != nil {
		return err
	}

	for i, column := range table.Columns {
		if i != 0 {
			if _, err := fmt.Fprint(buf, ",\n"); err != nil {
				return err
			}
		}
		if err := printColumnClause(buf, column, table); err != nil {
			return err
		}
	}
	writeSystemTimePeriod(buf, table)

	if err := printPrimaryKeyClause(buf, table.Indexes, table.Engine); err != nil {
		return err
	}

	secondaryIndexes := make([]*metadata.IndexMetadata, 0, len(table.Indexes))
	for _, index := range table.Indexes {
		if !index.Primary {
			secondaryIndexes = append(secondaryIndexes, index)
		}
	}
	foreignKeys := slices.Clone(table.ForeignKeys)
	checks := slices.Clone(table.CheckConstraints)
	if sorted {
		slices.SortFunc(secondaryIndexes, func(a, b *metadata.IndexMetadata) int { return cmp.Compare(a.Name, b.Name) })
		slices.SortFunc(foreignKeys, func(a, b *metadata.ForeignKeyMetadata) int { return cmp.Compare(a.Name, b.Name) })
		slices.SortFunc(checks, func(a, b *metadata.CheckConstraintMetadata) int { return cmp.Compare(a.Name, b.Name) })
	}
	for _, index := range secondaryIndexes {
		if err := printIndexClause(buf, index, table.Engine); err != nil {
			return err
		}
	}
	for _, fk := range foreignKeys {
		if err := printForeignKeyClause(buf, fk, dbName); err != nil {
			return err
		}
	}
	for _, check := range checks {
		if err := printCheckClause(buf, check); err != nil {
			return err
		}
	}

	if _, err := buf.WriteString("\n)"); err != nil {
		return err
	}
	if table.Engine != "" {
		if _, err := fmt.Fprintf(buf, " ENGINE=%s", table.Engine); err != nil {
			return err
		}
	}
	if opts := filterCreateOptions(table.CreateOptions); opts != "" {
		if _, err := fmt.Fprintf(buf, " %s", opts); err != nil {
			return err
		}
	}
	if table.Charset != "" {
		if _, err := fmt.Fprintf(buf, " DEFAULT CHARSET=%s", table.Charset); err != nil {
			return err
		}
	}
	if table.Collation != "" {
		if _, err := fmt.Fprintf(buf, " COLLATE=%s", table.Collation); err != nil {
			return err
		}
	}
	if table.Comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT='%s'", table.Comment); err != nil {
			return err
		}
	}
	if ts := unquoteTableOption(table.Tablespace); ts != "" {
		if _, err := fmt.Fprintf(buf, " /*!50100 TABLESPACE %s */", quoteIdent(ts)); err != nil {
			return err
		}
	}
	writeSystemVersioningOption(buf, table)

	// The partition clause is a /*!50100 PARTITION BY … */ executable comment, which
	// omni's loader splices back as live SQL.
	if len(table.Partitions) > 0 {
		if err := printPartitionClause(buf, table.Partitions, table.Engine, true, true); err != nil {
			return err
		}
	}
	if table.ColumnGroups != "" {
		if _, err := fmt.Fprintf(buf, " WITH COLUMN GROUP(%s)", table.ColumnGroups); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(buf, ";\n\n"); err != nil {
		return err
	}
	for _, col := range table.Columns {
		if col.GetDefaultDropped() {
			if _, err := fmt.Fprintf(buf, "ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT;\n\n", quoteIdent(table.Name), quoteIdent(col.Name)); err != nil {
				return err
			}
		}
	}
	if table.HasMaterializedViewLog {
		if _, err := fmt.Fprintf(buf, "CREATE MATERIALIZED VIEW LOG ON %s;\n\n", quoteIdent(table.Name)); err != nil {
			return err
		}
	}
	return nil
}

func printPartitionClause(buf *strings.Builder, partitions []*metadata.TablePartitionMetadata, engine string, includeComments, includeTemplate bool) error {
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

	var sub0 *metadata.TablePartitionMetadata
	if len(partitions[0].Subpartitions) > 0 {
		sub0 = partitions[0].Subpartitions[0]
	} else if len(partitions[0].SubpartitionTemplate) > 0 {
		sub0 = partitions[0].SubpartitionTemplate[0]
	}
	if sub0 != nil {
		if err := writeSubpartitionBy(buf, sub0); err != nil {
			return err
		}
	}

	hasTemplate := includeTemplate && len(partitions[0].SubpartitionTemplate) > 0
	if hasTemplate {
		tpl := partitions[0].SubpartitionTemplate
		tplHead := tpl[0]
		if tplHead.Type == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED && sub0 != nil {
			tplHead = sub0
		}
		if err := writeSubpartitionTemplate(buf, tpl, tplHead, engine, includeComments); err != nil {
			return err
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

	if partitions[0].SplitSize != "" {
		if _, err := fmt.Fprintf(buf, " SIZE (%s)", quoteSplitSize(partitions[0].SplitSize)); err != nil {
			return err
		}
	}
	if partitions[0].AutoSize != "" {
		if _, err := fmt.Fprintf(buf, " PARTITION SIZE %s PARTITIONS AUTO", partitions[0].AutoSize); err != nil {
			return err
		}
	}

	if useDefault == 0 && partitions[0].AutoSize == "" && (len(partitions) > 1 || partitions[0].Name != "") {
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
			if _, err := fmt.Fprintf(buf, "PARTITION %s", quoteIdent(partition.Name)); err != nil {
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

			partComment := ""
			if includeComments {
				partComment = partition.Comment
			}
			if subUseDefault == 0 && len(partition.Subpartitions) > 0 && !hasTemplate {
				if err := writePartitionOptions(buf, engine, partComment); err != nil {
					return err
				}
				if _, err := fmt.Fprint(buf, "\n ("); err != nil {
					return err
				}
				for j, subPartition := range partition.Subpartitions {
					if _, err := fmt.Fprintf(buf, "SUBPARTITION %s", quoteIdent(subPartition.Name)); err != nil {
						return err
					}
					if subVals := subpartitionValuesSQL(subPartitionTypeName(subPartition.Type), subPartitionColumnCount(subPartition.Expression), subPartition.Value); subVals != "" {
						if _, err := buf.WriteString(subVals); err != nil {
							return err
						}
					}
					subComment := ""
					if includeComments {
						subComment = subPartition.Comment
					}
					if err := writePartitionOptions(buf, engine, subComment); err != nil {
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
			} else {
				if err := writePartitionOptions(buf, engine, partComment); err != nil {
					return err
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

func writePartitionOptions(buf io.StringWriter, engine, comment string) error {
	if engine == "" {
		engine = "InnoDB"
	}
	if _, err := buf.WriteString(" ENGINE=" + engine); err != nil {
		return err
	}
	if comment != "" {
		if _, err := buf.WriteString(" COMMENT = '" + escapeComment(comment) + "'"); err != nil {
			return err
		}
	}

	return nil
}

func writeSubpartitionTemplate(buf *strings.Builder, tpl []*metadata.TablePartitionMetadata, head *metadata.TablePartitionMetadata, engine string, includeComments bool) error {
	subType := subPartitionTypeName(head.Type)
	colCount := subPartitionColumnCount(head.Expression)
	if _, err := fmt.Fprint(buf, "\nSUBPARTITION TEMPLATE ("); err != nil {
		return err
	}
	for i, st := range tpl {
		if i > 0 {
			if _, err := fmt.Fprint(buf, ",\n "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(buf, "SUBPARTITION %s", quoteIdent(st.Name)); err != nil {
			return err
		}
		if subVals := subpartitionValuesSQL(subType, colCount, st.Value); subVals != "" {
			if _, err := buf.WriteString(subVals); err != nil {
				return err
			}
		}
		cmt := ""
		if includeComments {
			cmt = st.Comment
		}
		if err := writePartitionOptions(buf, engine, cmt); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(buf, ")")
	return err
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

func writeSubpartitionBy(buf *strings.Builder, sub0 *metadata.TablePartitionMetadata) error {
	if _, err := fmt.Fprint(buf, "\nSUBPARTITION BY "); err != nil {
		return err
	}
	switch sub0.Type {
	case metadata.TablePartitionMetadata_RANGE:
		_, err := fmt.Fprintf(buf, "RANGE (%s)", sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_RANGE_COLUMNS:
		_, err := fmt.Fprintf(buf, "RANGE COLUMNS (%s)", sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_LIST:
		_, err := fmt.Fprintf(buf, "LIST (%s)", sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_LIST_COLUMNS:
		_, err := fmt.Fprintf(buf, "LIST COLUMNS (%s)", sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_HASH:
		_, err := fmt.Fprintf(buf, "HASH (%s)", sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_LINEAR_HASH:
		_, err := fmt.Fprintf(buf, "LINEAR HASH (%s)", sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_KEY:
		_, err := fmt.Fprintf(buf, "KEY%s (%s)", keyAlgorithmClause(sub0.Algorithm), sub0.Expression)
		return err
	case metadata.TablePartitionMetadata_LINEAR_KEY:
		_, err := fmt.Fprintf(buf, "LINEAR KEY%s (%s)", keyAlgorithmClause(sub0.Algorithm), sub0.Expression)
		return err
	default:
		return fmt.Errorf("invalid subpartition type: %v", sub0.Type)
	}
}

func subPartitionTypeName(tp metadata.TablePartitionMetadata_Type) string {
	switch tp {
	case metadata.TablePartitionMetadata_RANGE:
		return "RANGE"
	case metadata.TablePartitionMetadata_RANGE_COLUMNS:
		return "RANGE COLUMNS"
	case metadata.TablePartitionMetadata_LIST:
		return "LIST"
	case metadata.TablePartitionMetadata_LIST_COLUMNS:
		return "LIST COLUMNS"
	default:
		return ""
	}
}

func subPartitionColumnCount(expr string) int {
	if strings.TrimSpace(expr) == "" {
		return 1
	}
	return strings.Count(expr, ",") + 1
}

func printCheckClause(buf *strings.Builder, check *metadata.CheckConstraintMetadata) error {
	if _, err := fmt.Fprintf(buf, ",\n  CONSTRAINT %s CHECK %s", quoteIdent(check.Name), check.Expression); err != nil {
		return err
	}
	if check.NotEnforced {
		if _, err := buf.WriteString(" NOT ENFORCED"); err != nil {
			return err
		}
	}
	return nil
}

// printForeignKeyClause qualifies the referenced table with its database when that
// is not dbName, the database being dumped.
func printForeignKeyClause(buf *strings.Builder, fk *metadata.ForeignKeyMetadata, dbName string) error {
	if _, err := fmt.Fprintf(buf, ",\n  CONSTRAINT %s FOREIGN KEY (", quoteIdent(fk.Name)); err != nil {
		return err
	}

	for i, column := range fk.Columns {
		if i != 0 {
			if _, err := fmt.Fprint(buf, ", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(quoteIdent(column)); err != nil {
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
		if _, err := buf.WriteString(quoteIdent(column)); err != nil {
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

func printIndexClause(buf *strings.Builder, index *metadata.IndexMetadata, engine string) error {
	if index.Primary {
		return nil
	}

	if _, err := fmt.Fprint(buf, ",\n  "); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(buf, "%sKEY %s (", indexQualifier(index), quoteIdent(index.Name)); err != nil {
		return err
	}
	return printIndexKeyParts(buf, index, engine)
}

// indexQualifier is the keyword that opens an index, in a CREATE TABLE clause and in a
// standalone CREATE INDEX alike.
func indexQualifier(index *metadata.IndexMetadata) string {
	switch {
	case index.Unique:
		return "UNIQUE "
	case strings.EqualFold(index.Type, mysqlIndexFullText):
		return "FULLTEXT "
	case strings.EqualFold(index.Type, mysqlIndexSpatial):
		return "SPATIAL "
	case index.IsVector:
		return "VECTOR "
	}
	return ""
}

// printIndexKeyParts writes the key parts, type and options of an index, from its name and
// opening parenthesis to the end of its clause.
func printIndexKeyParts(buf *strings.Builder, index *metadata.IndexMetadata, engine string) error {
	isFullText := strings.EqualFold(index.Type, mysqlIndexFullText)
	isSpatial := strings.EqualFold(index.Type, mysqlIndexSpatial)

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
		// SPATIAL and FULLTEXT index parts do not take a key-part prefix length. The sync
		// records SUB_PART for spatial parts (e.g. 32 for a POINT) even though SHOW CREATE
		// renders none; emitting `(32)` would diverge from the user form and DROP+ADD the
		// index on every no-op. Suppress the prefix for those index types.
		if err := printIndexKeyPart(buf, expr, keyLength, descending, isSpatial || isFullText); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(buf, ")"); err != nil {
		return err
	}

	if isFullText {
		if index.Parser != "" {
			if _, err := fmt.Fprintf(buf, " WITH PARSER %s", index.Parser); err != nil {
				return err
			}
		}
		if len(index.ParserProperties) > 0 {
			parts := make([]string, 0, len(index.ParserProperties))
			for _, key := range parserPropertyOrder {
				v, ok := index.ParserProperties[key]
				if !ok {
					continue
				}
				if stringParserProperties[key] {
					v = quoteStringLiteral(v)
				}
				parts = append(parts, key+"="+v)
			}
			if _, err := fmt.Fprintf(buf, " PARSER_PROPERTIES=(%s)", strings.Join(parts, ",")); err != nil {
				return err
			}
		}
	}

	if len(index.StoringColumns) > 0 {
		names := make([]string, 0, len(index.StoringColumns))
		for _, col := range index.StoringColumns {
			names = append(names, quoteIdent(col))
		}
		if _, err := fmt.Fprintf(buf, " STORING (%s)", strings.Join(names, ", ")); err != nil {
			return err
		}
	}

	if index.IsVector && len(index.VectorParams) > 0 {
		if _, err := fmt.Fprintf(buf, " WITH (%s)", strings.Join(index.VectorParams, ", ")); err != nil {
			return err
		}
	}

	if using := indexUsingClause(index.Type, engine); using != "" {
		if _, err := buf.WriteString(using); err != nil {
			return err
		}
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

	// INVISIBLE secondary index (MySQL 8.0). The version-gated executable comment is parsed
	// back by omni's loader, so the dumped `from` carries invisibility and the diff against
	// an INVISIBLE user `to` is a no-op instead of a DROP+ADD.
	if !index.Visible {
		if _, err := fmt.Fprint(buf, " /*!80000 INVISIBLE */"); err != nil {
			return err
		}
	}

	scope := index.Scope
	if scope == "" && index.IsGlobal {
		scope = "GLOBAL"
	}
	if scope != "" {
		if _, err := fmt.Fprintf(buf, " %s", strings.ToUpper(scope)); err != nil {
			return err
		}
	}

	return nil
}

func printIndexKeyPart(buf *strings.Builder, expr string, length int64, descending bool, suppressPrefix bool) error {
	if len(expr) == 0 {
		return errors.New("index key part expression is empty")
	}
	if expr[0] == '(' && expr[len(expr)-1] == ')' {
		if _, err := buf.WriteString(normalizeFunctionalIndexExpr(expr)); err != nil {
			return err
		}
	} else {
		if _, err := buf.WriteString(quoteIdent(expr)); err != nil {
			return err
		}
	}
	if length > 0 && !suppressPrefix {
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

// normalizeFunctionalIndexExpr rewrites a synced functional (expression) index key part into
// the form the omni SDL loader can parse. The information_schema.STATISTICS.EXPRESSION text
// MySQL records for an `((expr))` key part:
//   - backslash-escapes the single quotes around string literals (e.g. json_extract(`tags`,
//     _utf8mb4\'$.ids\')), which the omni lexer rejects (`syntax error at or near "\"`); and
//   - prefixes string literals with a charset introducer (`_utf8mb4'…'`, `_latin1'…'`, …)
//     reflecting the connection charset, which the omni expression parser does not accept in
//     this position (`unexpected token`).
//
// Both make LoadSDL fail outright on the dumped source even though omni accepts the canonical
// user spelling (e.g. `((cast(json_extract(`tags`,'$.ids') as unsigned array)))`). This
// unescapes the quotes and strips the introducer so the emitted expression is omni-parseable;
// omni's functional-index normalizer then canonicalizes the dumped and user forms identically,
// yielding an empty no-op. The transform is string-literal aware: the introducer is removed
// only when it immediately precedes a string literal, and literal contents are left untouched.
func normalizeFunctionalIndexExpr(expr string) string {
	return stripIndexExprCharsetIntroducers(unescapeFunctionalIndexQuotes(expr))
}

// unescapeFunctionalIndexQuotes undoes the backslash escaping that
// information_schema.STATISTICS.EXPRESSION applies to single quotes (\' -> ') and backslashes
// (\\ -> \). MySQL stores the expression text escaped; the omni lexer wants the raw SQL.
func unescapeFunctionalIndexQuotes(expr string) string {
	if !strings.Contains(expr, `\`) {
		return expr
	}
	var b strings.Builder
	b.Grow(len(expr))
	for i := 0; i < len(expr); i++ {
		if expr[i] == '\\' && i+1 < len(expr) && (expr[i+1] == '\'' || expr[i+1] == '\\') {
			b.WriteByte(expr[i+1])
			i++
			continue
		}
		b.WriteByte(expr[i])
	}
	return b.String()
}

// stripIndexExprCharsetIntroducers removes a MySQL charset introducer token (an underscore-prefixed
// identifier such as `_utf8mb4` or `_latin1`) when it immediately precedes a single-quoted
// string literal, e.g. `_utf8mb4'$.ids'` -> `'$.ids'`. The scan tracks whether it is inside a
// '…'/"…" string literal (honoring backslash escapes and the doubled-quote escape) so that an
// underscore sequence occurring inside literal text is never mistaken for an introducer and
// literal contents are preserved byte for byte. Backtick-quoted identifiers are passed through
// untouched. Only an introducer attached to a string literal is stripped; underscores that are
// part of an ordinary identifier (`my_col`, a leading-underscore name not followed by a quote)
// are left intact.
func stripIndexExprCharsetIntroducers(expr string) string {
	if !strings.Contains(expr, "_") {
		return expr
	}
	var b strings.Builder
	b.Grow(len(expr))
	var quote byte // open delimiter of the literal currently being copied, or 0 when outside one
	for i := 0; i < len(expr); {
		c := expr[i]
		if quote != 0 {
			b.WriteByte(c)
			switch {
			case c == '\\' && i+1 < len(expr):
				b.WriteByte(expr[i+1])
				i += 2
				continue
			case c == quote:
				if i+1 < len(expr) && expr[i+1] == quote {
					b.WriteByte(expr[i+1])
					i += 2
					continue
				}
				quote = 0
			default:
			}
			i++
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
			b.WriteByte(c)
			i++
		case '`':
			// Copy a backtick-quoted identifier verbatim so an underscore inside it is never
			// treated as an introducer.
			b.WriteByte(c)
			i++
			for i < len(expr) {
				b.WriteByte(expr[i])
				if expr[i] == '`' {
					i++
					break
				}
				i++
			}
		case '_':
			if n := charsetIntroducerLen(expr[i:]); n > 0 {
				// Drop the `_<charset>` introducer; the following string literal is copied on the
				// next iterations.
				i += n
				continue
			}
			b.WriteByte(c)
			i++
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String()
}

// charsetIntroducerLen returns the byte length of a leading charset-introducer token in s
// when s begins with `_<charset>` immediately followed by a single-quoted string literal,
// or 0 otherwise. The charset name accepts the ASCII letters/digits MySQL charset names use
// (e.g. utf8mb4, latin1, big5); the terminating quote is required so a bare leading-underscore
// identifier is not stripped.
func charsetIntroducerLen(s string) int {
	if len(s) < 2 || s[0] != '_' {
		return 0
	}
	i := 1
	for i < len(s) && isCharsetNameByte(s[i]) {
		i++
	}
	if i == 1 || i >= len(s) || s[i] != '\'' {
		return 0
	}
	return i
}

func isCharsetNameByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

func printPrimaryKeyClause(buf *strings.Builder, indexes []*metadata.IndexMetadata, engine string) error {
	for _, index := range indexes {
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
				// Primary keys are never SPATIAL/FULLTEXT, so prefix lengths are valid here.
				if err := printIndexKeyPart(buf, column, keyLength, descending, false); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(buf, ")", indexUsingClause(index.Type, engine)); err != nil {
				return err
			}
			if index.KeyBlockSize > 0 {
				if _, err := fmt.Fprintf(buf, " KEY_BLOCK_SIZE=%d", index.KeyBlockSize); err != nil {
					return err
				}
			}
			return nil
		}
	}

	return nil
}

// indexUsingClause returns the USING clause for an index whose BTREE or HASH access
// method is not the one engine gives a key written without USING.
func indexUsingClause(indexType, engine string) string {
	indexType = strings.ToUpper(indexType)
	if (indexType != "BTREE" && indexType != "HASH") || indexType == defaultIndexType(engine) {
		return ""
	}
	return " USING " + indexType
}

func isAutoIncrement(column *metadata.ColumnMetadata) bool {
	return strings.EqualFold(column.GetDefault(), autoIncrementSymbol)
}

func printColumnClause(buf *strings.Builder, column *metadata.ColumnMetadata, table *metadata.TableMetadata) error {
	if _, err := fmt.Fprintf(buf, "  %s %s", quoteIdent(column.Name), normalizeColumnType(column.Type)); err != nil {
		return err
	}

	if column.CharacterSet != "" && (column.CharsetInherited != nil && !*column.CharsetInherited || column.CharacterSet != table.Charset) {
		if _, err := fmt.Fprintf(buf, " CHARACTER SET %s", column.CharacterSet); err != nil {
			return err
		}
	}

	if column.Collation != "" && (column.CollationInherited != nil && !*column.CollationInherited || column.Collation != table.Collation) {
		if _, err := fmt.Fprintf(buf, " COLLATE %s", column.Collation); err != nil {
			return err
		}
	}

	if isSystemTimeColumn(column) {
		writeSystemTimeColumnAttribute(buf, column)
	} else if column.Generation != nil && column.Generation.Expression != "" {
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
			// Unknown generation type
		}
	}

	if !column.Nullable {
		if _, err := fmt.Fprint(buf, " NOT NULL"); err != nil {
			return err
		}
	}

	// SRID for spatial columns (MySQL 8.0). Emitted after NOT NULL and before DEFAULT to
	// match the SHOW CREATE / omni-canonical position; the version-gated executable comment
	// is spliced back as live SQL by omni's loader so the dumped `from` carries the SRID and
	// the diff against a user `to` with the same SRID is a no-op.
	writeColumnSRIDAttribute(buf, column)

	if err := printDefaultClause(buf, column); err != nil {
		return err
	}

	// Handle auto_increment.
	if isAutoIncrement(column) {
		if _, err := fmt.Fprintf(buf, " %s", autoIncrementSymbol); err != nil {
			return err
		}
	}

	if column.OnUpdate != "" {
		if _, err := fmt.Fprintf(buf, " ON UPDATE %s", column.OnUpdate); err != nil {
			return err
		}
	}

	// INVISIBLE column (MySQL 8.0.23+). Emitted before COMMENT to match the SHOW CREATE /
	// omni-canonical order (verified against 8.0.32 for both regular and generated columns);
	// like SRID the version-gated comment is parsed back by omni's loader so the dumped `from`
	// carries invisibility and the diff against an INVISIBLE user `to` is a no-op.
	writeColumnInvisibleAttribute(buf, column)

	if column.Comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT '%s'", column.Comment); err != nil {
			return err
		}
	}
	return nil
}

// writeColumnSRIDAttribute emits the ` /*!80003 SRID n */` column attribute when the
// column carries an explicit SRID — including the valid SRID 0, which presence (not a
// zero sentinel) distinguishes from "no SRID".
func writeColumnSRIDAttribute(buf *strings.Builder, column *metadata.ColumnMetadata) {
	if column.Srid != nil {
		_, _ = fmt.Fprintf(buf, " /*!80003 SRID %d */", column.GetSrid())
	}
}

// writeColumnInvisibleAttribute emits the ` /*!80023 INVISIBLE */` column attribute for
// invisible columns (MySQL 8.0.23+). MySQL's canonical SHOW CREATE order places INVISIBLE
// before COMMENT for both regular and generated columns (verified against 8.0.32), so callers
// emit it before COMMENT.
func writeColumnInvisibleAttribute(buf *strings.Builder, column *metadata.ColumnMetadata) {
	if column.IsInvisible {
		_, _ = buf.WriteString(" /*!80023 INVISIBLE */")
	}
}

// normalizeColumnType rewrites engine-specific column-type spellings the omni SDL parser
// does not accept into the canonical spelling it does. MySQL 8.0's SHOW CREATE /
// information_schema render GEOMETRYCOLLECTION as the synonym `geomcollection`, which the
// omni loader rejects with "expected data type"; the canonical `geometrycollection`
// (what 5.7 already dumps) round-trips. The two spellings denote the same type, so this is
// a pure rename with no semantic change.
func normalizeColumnType(columnType string) string {
	if strings.EqualFold(columnType, "geomcollection") {
		return "geometrycollection"
	}
	return columnType
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
		if _, err := fmt.Fprintf(buf, " DEFAULT %s", renderColumnDefault(column)); err != nil {
			return err
		}
	}

	return nil
}

// renderColumnDefault returns the DEFAULT expression text to emit for a column.
//
// For a BIT column, the synced COLUMN_DEFAULT is the bit-literal text (e.g. b'0'),
// which the sync stores QUOTE()-escaped as a string literal in ColumnMetadata.Default
// (the bit literal wrapped in single quotes with the inner quotes backslash-escaped).
// Emitting that verbatim produces a quoted string default that the omni loader takes as a
// string — never matching the user's BIT literal b'0', so the no-op phantom-re-emits
// MODIFY ... DEFAULT b'0' forever. Real mysqldump renders the bit literal unquoted
// (DEFAULT b'0'); match it by recovering the inner literal and emitting it without the
// surrounding quotes when the column is BIT and the recovered text is a bit/hex literal
// (b'…'/0x…/x'…'). Any other default (including a genuine string) is emitted verbatim.
func renderColumnDefault(column *metadata.ColumnMetadata) string {
	if isBitColumnType(column.Type) {
		if lit, ok := bitLiteralFromQuotedDefault(column.Default); ok {
			return lit
		}
	}
	return column.Default
}

// isBitColumnType reports whether a column type spelling is BIT (with or without a
// display width, e.g. "bit" or "bit(8)").
func isBitColumnType(columnType string) bool {
	lower := strings.ToLower(strings.TrimSpace(columnType))
	return lower == "bit" || strings.HasPrefix(lower, "bit(")
}

// bitLiteralFromQuotedDefault recovers a bit/hex literal from a stored column default and
// reports whether the default denotes one. The default may already be a bare literal
// (b'0', 0x0A, x'1f') or the QUOTE()-escaped string form the MySQL sync stores for static
// defaults (the literal wrapped in single quotes with inner quotes backslash-escaped). It
// strips at most one layer of surrounding single quotes and the QUOTE() backslash escaping,
// then recognizes the MySQL bit/hex literal spellings
// (b'…' / B'…' bit, 0x… hex, x'…' / X'…' hex). On success it returns the canonical literal
// text to emit unquoted; otherwise it returns ("", false) and the caller emits the default
// verbatim.
func bitLiteralFromQuotedDefault(def string) (string, bool) {
	if isBitOrHexLiteral(def) {
		return def, true
	}
	if len(def) >= 2 && def[0] == '\'' && def[len(def)-1] == '\'' {
		inner := unescapeQuotedDefault(def[1 : len(def)-1])
		if isBitOrHexLiteral(inner) {
			return inner, true
		}
	}
	return "", false
}

// isBitOrHexLiteral reports whether s is a MySQL bit literal (b'…' / B'…') or hex literal
// (0x… / 0X… / x'…' / X'…') with only valid digits between the delimiters.
func isBitOrHexLiteral(s string) bool {
	switch {
	case len(s) >= 3 && (s[0] == 'b' || s[0] == 'B') && s[1] == '\'' && s[len(s)-1] == '\'':
		return allDigitsIn(s[2:len(s)-1], "01")
	case len(s) >= 3 && (s[0] == 'x' || s[0] == 'X') && s[1] == '\'' && s[len(s)-1] == '\'':
		return allDigitsIn(s[2:len(s)-1], "0123456789abcdefABCDEF")
	case len(s) >= 3 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X'):
		return allDigitsIn(s[2:], "0123456789abcdefABCDEF")
	default:
		return false
	}
}

func allDigitsIn(s, allowed string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !strings.ContainsRune(allowed, rune(s[i])) {
			return false
		}
	}
	return true
}

// unescapeQuotedDefault undoes the QUOTE() backslash escaping (\' -> ', \\ -> \) the MySQL
// sync applies to static defaults. Only these two sequences appear inside a bit/hex literal,
// so the full QUOTE() escape set is intentionally not reproduced here.
func unescapeQuotedDefault(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && (s[i+1] == '\'' || s[i+1] == '\\') {
			b.WriteByte(s[i+1])
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
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
	if _, err := io.WriteString(out, quoteIdent(view.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}
	_, err := io.WriteString(out, "-- "+quoteIdent(view.Name)+" references invalid table(s) or column(s) or function(s) or definer/invoker of view lack rights to use them\n\n")
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
	if _, err := io.WriteString(out, quoteIdent(view.Name)+"\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, emptyCommentLine); err != nil {
		return err
	}

	// Definition.
	if _, err := io.WriteString(out, "CREATE VIEW "+quoteIdent(view.Name)+" AS SELECT\n  "); err != nil {
		return err
	}
	for i, column := range view.Columns {
		if i != 0 {
			if _, err := io.WriteString(out, ",\n  "); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, "1 AS "+quoteIdent(column.Name)); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, ";\n\n")
	return err
}

func writeAdditionalEventsIfSet(out io.Writer, characterSetClient, characterSetResult, collationConnection, sqlMode string) error {
	events := []struct {
		condition bool
		prefix    string
		value     string
	}{
		{characterSetClient != "", setCharacterSetClient, characterSetClient},
		{characterSetResult != "", setCharacterSetResult, characterSetResult},
		{collationConnection != "", setCollation, collationConnection},
	}

	for _, event := range events {
		if event.condition {
			if _, err := io.WriteString(out, event.prefix); err != nil {
				return err
			}
			if _, err := io.WriteString(out, event.value); err != nil {
				return err
			}
			if _, err := io.WriteString(out, ";\n"); err != nil {
				return err
			}
		}
	}
	if sqlMode != "" {
		if _, err := io.WriteString(out, setSQLMode); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "'"); err != nil {
			return err
		}
		if _, err := io.WriteString(out, sqlMode); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "';\n"); err != nil {
			return err
		}
	}

	return nil
}

// getSDLFormat emits a canonical, deterministic SDL dump of the database metadata
// such that the omni catalog LoadSDL round-trips it. Unlike the mysqldump-style
// output DatabaseDefinition produces by default, the SDL dump:
//   - contains only the constructive CREATE statements LoadSDL accepts (no SET
//     headers, no DELIMITER blocks, no temporary-view scaffolding, no version-gated
//     /*!50100 PARTITION */ executable comments);
//   - emits objects and intra-table clauses in a stable order (tables and views
//     sorted by name; secondary indexes, foreign keys, and check constraints sorted
//     by name) so the same metadata always produces byte-identical SDL;
//   - strips the volatile AUTO_INCREMENT=N table option (a live counter, never
//     schema);
//   - emits stored routines, triggers, and events as their canonical CREATE forms
//     (DEFINER omitted, no DELIMITER blocks — omni's LoadSDL splitter is BEGIN…END
//     aware and the routine/trigger/event differs ignore DEFINER).
//
// Emission order is dependency-safe so every CREATE resolves against objects loaded
// before it: tables → functions → procedures → views → triggers → events. Functions
// precede views/triggers because a view or trigger body may call a stored function;
// triggers and events follow their tables because the omni loader requires a
// trigger's owning table (and any table an event/trigger body touches) to exist.
//
// Type, charset/collation, and default canonicalization are intentionally NOT done
// here: the omni Diff path resolves both the dumped source and the user target
// through the same Normalizer (CanonicalColumn), so faithful emission is sufficient
// for the no-op idempotence property — the canonical comparison happens in omni.
func getSDLFormat(meta *metadata.DatabaseSchemaMetadata) (string, error) {
	var buf strings.Builder

	schema := meta.Schemas[0]

	tables := make([]*metadata.TableMetadata, 0, len(schema.Tables))
	for _, table := range schema.Tables {
		if table.SkipDump {
			continue
		}
		tables = append(tables, table)
	}
	slices.SortFunc(tables, func(a, b *metadata.TableMetadata) int { return cmp.Compare(a.Name, b.Name) })

	for _, table := range tables {
		if err := writeTableSDL(&buf, meta.Name, table); err != nil {
			return "", err
		}
	}

	functions := make([]*metadata.FunctionMetadata, 0, len(schema.Functions))
	for _, function := range schema.Functions {
		if function.SkipDump {
			continue
		}
		functions = append(functions, function)
	}
	slices.SortFunc(functions, func(a, b *metadata.FunctionMetadata) int { return cmp.Compare(a.Name, b.Name) })
	for _, function := range functions {
		if err := writeRoutineSDL(&buf, function.Definition); err != nil {
			return "", err
		}
	}

	procedures := make([]*metadata.ProcedureMetadata, 0, len(schema.Procedures))
	for _, procedure := range schema.Procedures {
		if procedure.SkipDump {
			continue
		}
		procedures = append(procedures, procedure)
	}
	slices.SortFunc(procedures, func(a, b *metadata.ProcedureMetadata) int { return cmp.Compare(a.Name, b.Name) })
	for _, procedure := range procedures {
		if err := writeRoutineSDL(&buf, procedure.Definition); err != nil {
			return "", err
		}
	}

	views := make([]*metadata.ViewMetadata, 0, len(schema.Views))
	for _, view := range schema.Views {
		if view.SkipDump {
			continue
		}
		views = append(views, view)
	}
	slices.SortFunc(views, func(a, b *metadata.ViewMetadata) int { return cmp.Compare(a.Name, b.Name) })

	for _, view := range views {
		if err := writeViewSDL(&buf, meta.Name, view); err != nil {
			return "", err
		}
	}

	materializedViews := make([]*metadata.MaterializedViewMetadata, 0, len(schema.MaterializedViews))
	for _, view := range schema.MaterializedViews {
		if view.SkipDump {
			continue
		}
		materializedViews = append(materializedViews, view)
	}
	slices.SortFunc(materializedViews, func(a, b *metadata.MaterializedViewMetadata) int { return cmp.Compare(a.Name, b.Name) })
	for _, view := range materializedViews {
		if err := writeMaterializedViewSDL(&buf, meta.Name, view); err != nil {
			return "", err
		}
	}

	// Triggers hang off tables (TableMetadata.Triggers); emit them after every table
	// in deterministic (table, trigger) order.
	for _, table := range tables {
		triggers := make([]*metadata.TriggerMetadata, 0, len(table.Triggers))
		for _, trigger := range table.Triggers {
			if trigger.SkipDump {
				continue
			}
			triggers = append(triggers, trigger)
		}
		slices.SortFunc(triggers, func(a, b *metadata.TriggerMetadata) int { return cmp.Compare(a.Name, b.Name) })
		for _, trigger := range triggers {
			if err := writeTriggerSDL(&buf, table.Name, trigger); err != nil {
				return "", err
			}
		}
	}

	events := make([]*metadata.EventMetadata, 0, len(schema.Events))
	events = append(events, schema.Events...)
	slices.SortFunc(events, func(a, b *metadata.EventMetadata) int { return cmp.Compare(a.Name, b.Name) })
	for _, event := range events {
		if err := writeEventSDL(&buf, event); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

// writeTableSDL emits a single CREATE TABLE statement in canonical SDL form, with
// secondary indexes, foreign keys, and checks sorted by name.
func writeTableSDL(buf *strings.Builder, dbName string, table *metadata.TableMetadata) error {
	return writeCreateTable(buf, table, dbName, true)
}

// filterCreateOptions returns the table CREATE_OPTIONS string with the synthetic
// `partitioned` token removed. MySQL's information_schema.TABLES.CREATE_OPTIONS folds the
// fact that a table is partitioned into a literal `partitioned` token alongside the real
// declared options (e.g. `partitioned`, or `row_format=DYNAMIC partitioned`). Partitioning
// is emitted separately via the PARTITION BY clause from TableMetadata.Partitions, so the
// token must not be rendered as a table option (it is not valid table-option DDL). The
// remaining tokens (row_format=…, key_block_size=…, …) are already valid `name=value`
// table-option syntax that the omni loader parses case-insensitively, so they are returned
// verbatim and joined by single spaces.
func filterCreateOptions(createOptions string) string {
	if createOptions == "" {
		return ""
	}
	kept := make([]string, 0, 4)
	// The system-versioning marker is written separately by writeSystemVersioningOption.
	for _, tok := range strings.Fields(strings.Replace(createOptions, systemVersioningOption, "", 1)) {
		if strings.EqualFold(tok, "partitioned") {
			continue
		}
		kept = append(kept, tok)
	}
	return strings.Join(kept, " ")
}

// writeViewsInDependencyOrder writes the regular views and the materialized views, each
// after the ones it reads. A view may read a materialized view and a materialized view may
// read a view, and a materialized view must not be created while the temporary placeholder
// stands in for a view it reads: the view it would materialize has the placeholder's column
// types, `1` for a `varchar` column among them.
func writeViewsInDependencyOrder(buf *strings.Builder, dbName string, schema *metadata.SchemaMetadata) error {
	// Definitions are identified by their exact name: a tenant with
	// lower_case_table_names=0 holds views that differ only in case, and both are written.
	type viewDefinition struct {
		name  string
		refs  []relationRef
		write func(*strings.Builder) error
	}
	ownDB := toLower(dbName)
	definitions := make([]*viewDefinition, 0, len(schema.Views)+len(schema.MaterializedViews))
	for _, view := range schema.Views {
		definitions = append(definitions, &viewDefinition{
			name: view.Name,
			refs: viewReferences(view.Definition, ownDB),
			write: func(out *strings.Builder) error {
				return writeView(out, view)
			},
		})
	}
	for _, view := range schema.MaterializedViews {
		if view.SkipDump {
			continue
		}
		definitions = append(definitions, &viewDefinition{
			name: view.Name,
			refs: viewReferences(view.Definition, ownDB),
			write: func(out *strings.Builder) error {
				return writeMaterializedView(out, view)
			},
		})
	}
	slices.SortFunc(definitions, func(a, b *viewDefinition) int { return cmp.Compare(a.name, b.name) })
	// An exact name wins, so two views that differ only in case keep their own edges. A
	// reference whose case differs from the name the definition was synced under falls back
	// to the folded name, and only where that names one view: choosing between two
	// case-distinct views could order one against what it reads.
	byName := make(map[string][]*viewDefinition, len(definitions))
	byFoldedName := make(map[string][]*viewDefinition, len(definitions))
	for _, definition := range definitions {
		byName[definition.name] = append(byName[definition.name], definition)
		folded := toLower(definition.name)
		byFoldedName[folded] = append(byFoldedName[folded], definition)
	}

	written := make(map[string]bool, len(definitions))
	var write func(view *viewDefinition) error
	write = func(view *viewDefinition) error {
		if written[view.name] {
			return nil
		}
		// Marked before the walk so a cycle ends it; OceanBase refuses one, and the mark
		// only keeps a malformed snapshot from recursing.
		written[view.name] = true
		var dependencies []*viewDefinition
		seen := make(map[*viewDefinition]bool, len(view.refs))
		for _, ref := range view.refs {
			if ref.db != ownDB {
				continue
			}
			candidates := byName[ref.name]
			if len(candidates) == 0 {
				if folded := byFoldedName[toLower(ref.name)]; len(folded) == 1 {
					candidates = folded
				}
			}
			for _, dependency := range candidates {
				if dependency != view && !seen[dependency] {
					seen[dependency] = true
					dependencies = append(dependencies, dependency)
				}
			}
		}
		// In name order, so a definition reading several others emits them the same way
		// every time: a map's iteration order is not stable between runs.
		slices.SortFunc(dependencies, func(a, b *viewDefinition) int { return cmp.Compare(a.name, b.name) })
		for _, dependency := range dependencies {
			if err := write(dependency); err != nil {
				return err
			}
		}
		return view.write(buf)
	}
	for _, view := range definitions {
		if err := write(view); err != nil {
			return err
		}
	}
	return nil
}

// relationRef is a relation a definition reads, with the case the definition spells.
type relationRef struct {
	db   string
	name string
}

// viewReferences returns the relations a definition reads in ownDB, spelled as the
// definition spells them. referencedRelations answers the same question for a diff, keyed by
// folded name, which a dump cannot use: a tenant with lower_case_table_names=0 holds views
// that differ only in case, and folding them together would order one against another.
func viewReferences(definition, ownDB string) []relationRef {
	// A stored body may end in WITH CHECK OPTION, so it parses only as a whole view. A body
	// the parser cannot read yields no edge and the object is written in name order: the
	// definition came from the server, which accepted it, so the dump still writes it.
	stmt, err := parseFirstStatement[*nodes.CreateViewStmt]("CREATE VIEW `v` AS " + definition)
	if err != nil || stmt.Select == nil {
		return nil
	}
	var refs []relationRef
	for _, ref := range selectTableRefs(stmt.Select) {
		if ref == nil || ref.Name == "" {
			continue
		}
		db := ref.Schema
		if db == "" {
			db = ownDB
		}
		refs = append(refs, relationRef{db: toLower(db), name: ref.Name})
	}
	return refs
}

// writeMaterializedViewSDL emits a materialized view in the one form the loader holds: a
// view. The catalog models no materialized view, so a `CREATE MATERIALIZED VIEW` here would
// not load and a dependent view would read a relation SDL never wrote; the container's
// indexes and the refresh policy have no place in this form either.
func writeMaterializedViewSDL(buf *strings.Builder, dbName string, view *metadata.MaterializedViewMetadata) error {
	body := formatViewBodySDL(stripViewBodyDatabaseQualifier(view.Definition, dbName))
	if body == "" {
		return nil
	}
	_, err := fmt.Fprintf(buf, "CREATE OR REPLACE VIEW %s AS\n%s;\n\n", quoteIdent(view.Name), body)
	return err
}

// writeViewSDL emits a single CREATE OR REPLACE VIEW statement. LoadSDL accepts
// CREATE OR REPLACE VIEW and is order-tolerant, so no temporary-view scaffolding is
// needed. dbName is the database being dumped; its own-database qualifier is stripped
// from the body so the emitted form matches the user's unqualified spelling (see
// stripViewBodyDatabaseQualifier). The stored one-line body is then
// pretty-printed by formatViewBodySDL — a whitespace-only, deterministic rewrite, so
// the omni no-op invariant and the dump-cycle byte stability both hold.
func writeViewSDL(buf *strings.Builder, dbName string, view *metadata.ViewMetadata) error {
	body := formatViewBodySDL(stripViewBodyDatabaseQualifier(view.Definition, dbName))
	if body == "" {
		// Defensive: an empty definition keeps the historical inline framing.
		if _, err := fmt.Fprintf(buf, "CREATE OR REPLACE VIEW %s AS ;\n\n", quoteIdent(view.Name)); err != nil {
			return err
		}
		return nil
	}
	if _, err := fmt.Fprintf(buf, "CREATE OR REPLACE VIEW %s%s AS\n%s;\n\n", quoteIdent(view.Name), viewColumnListClause(view), body); err != nil {
		return err
	}
	return nil
}

// stripViewBodyDatabaseQualifier removes the dumped database's own qualifier from the
// table, column, and function references of a view body.
//
// information_schema.VIEWS qualifies the references of a view body with their
// database, while users write them unqualified. The omni SDL diff folds only the
// qualifier of the database it loads into, so a body still naming the synced database
// would differ from the user's forever. The body is parsed so that only a qualifier in the database
// position goes: a table alias or table spelled like the database is a qualifier of
// its columns and stays. A body that does not parse is returned unchanged. A table
// known by a name another relation of the body (a CTE, an alias, a derived table,
// another database's table) also has keeps its qualifier, since without one a
// reference to it would resolve to that relation, or be ambiguous.
func stripViewBodyDatabaseQualifier(body, dbName string) string {
	return rebaseViewBodyDatabaseQualifier(body, dbName, dbName)
}

// rebaseViewBodyDatabaseQualifier is stripViewBodyDatabaseQualifier for a body
// installed into database target: a table that keeps a qualifier is qualified
// with target instead.
func rebaseViewBodyDatabaseQualifier(body, dbName, target string) string {
	if dbName == "" || body == "" {
		return body
	}
	qualifier := quoteIdent(dbName) + "."
	if !strings.Contains(body, qualifier) {
		return body
	}
	// A stored body may end in WITH CHECK OPTION, so it parses only as a whole view.
	const prefix = "CREATE VIEW v AS "
	list, err := obparser.Parse(prefix + body)
	if err != nil || len(list.Items) != 1 {
		return body
	}
	// Each relation the body names is known by its alias, or else by its name. A
	// reference through dbName stays qualified when another relation, or a CTE, is
	// known by the same name: without the qualifier it would resolve to that one,
	// or be ambiguous.
	exposed := map[string]int{}
	ctes := map[string]bool{}
	expose := func(name string) {
		if name != "" {
			exposed[strings.ToLower(name)]++
		}
	}
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		switch n := n.(type) {
		case *nodes.CommonTableExpr:
			ctes[strings.ToLower(n.Name)] = true
			expose(n.Name)
		case *nodes.TableRef:
			switch {
			case n.Loc.Start < len(prefix):
				// The view's own name, in the prefix.
			case n.Alias != "":
				expose(n.Alias)
			default:
				expose(n.Name)
			}
		case *nodes.SubqueryExpr:
			expose(n.Alias)
		case *nodes.JsonTableExpr:
			expose(n.Alias)
		case *nodes.ExternalTableRef:
			expose(n.Alias)
		case *nodes.RbIterateExpr:
			expose(n.Alias)
		case *nodes.UnnestExpr:
			expose(n.Alias)
		case *nodes.TableCollectionExpr:
			expose(n.Alias)
		}
		return true
	})
	// keepQualifier reports whether a reference through dbName to table, by its
	// name in a FROM clause (tableRef) or as a column's qualifier, needs one.
	keepQualifier := func(table string, tableRef bool) bool {
		table = strings.ToLower(table)
		if tableRef && ctes[table] {
			return true
		}
		return exposed[table] > 1
	}
	replacement := make(map[int]string)
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		var schema, table string
		var loc nodes.Loc
		tableRef := false
		switch n := n.(type) {
		case *nodes.RemoteSequenceExpr:
			return false
		case *nodes.TableRef:
			if n.Catalog != "" || n.Dblink != "" {
				return true
			}
			schema, table, loc, tableRef = n.Schema, n.Name, n.Loc, true
			if n.Alias != "" && !ctes[strings.ToLower(n.Name)] {
				// Aliased, the table is known by its alias, which nothing else
				// can shadow here.
				table = ""
			}
		case *nodes.ColumnRef:
			if n.Catalog != "" {
				return true
			}
			schema, table, loc = n.Schema, n.Table, n.Loc
		case *nodes.FuncCallExpr:
			schema, loc = n.Schema, n.Loc
		default:
			return true
		}
		start := loc.Start - len(prefix)
		if schema == dbName && start >= 0 && strings.HasPrefix(body[start:], qualifier) {
			replacement[start] = ""
			if table != "" && keepQualifier(table, tableRef) {
				replacement[start] = quoteIdent(target) + "."
			}
		}
		return true
	})
	starts := slices.Sorted(maps.Keys(replacement))

	var out strings.Builder
	out.Grow(len(body))
	prev := 0
	for _, start := range starts {
		out.WriteString(body[prev:start])
		out.WriteString(replacement[start])
		prev = start + len(qualifier)
	}
	out.WriteString(body[prev:])
	return out.String()
}

// writeRoutineSDL emits a stored function or procedure. The synced Definition is
// already a complete, canonical CREATE FUNCTION/PROCEDURE statement (the MySQL sync
// rewrites SHOW CREATE to drop the DEFINER and the RETURNS charset attribute), so the
// only framing the SDL form adds is a terminating ";". The body is emitted verbatim:
// omni's loader splitter is BEGIN…END aware, so internal ";" inside a compound body
// do not need DELIMITER wrapping, and omni's routine differ compares the body byte for
// byte after trimming — reformatting would phantom-diff. A leading DEFINER, if one
// somehow survives, is stripped to match omni's DEFINER-agnostic routine identity.
func writeRoutineSDL(buf *strings.Builder, definition string) error {
	def := strings.TrimSpace(stripLeadingDefiner(definition))
	if def == "" {
		return nil
	}
	if _, err := fmt.Fprintf(buf, "%s;\n\n", def); err != nil {
		return err
	}
	return nil
}

// writeTriggerSDL constructs a canonical CREATE TRIGGER from the trigger's metadata.
// MySQL stores only the action statement (TriggerMetadata.Body) plus the timing,
// event, and owning table, so the CREATE wrapper is rebuilt here exactly as omni's
// loader expects: `CREATE TRIGGER <name> <timing> <event> ON <table> FOR EACH ROW
// <body>`. DEFINER is omitted (omni's trigger differ ignores it) and the body is
// emitted verbatim (compared byte for byte after trimming).
func writeTriggerSDL(buf *strings.Builder, tableName string, trigger *metadata.TriggerMetadata) error {
	if _, err := fmt.Fprintf(buf, "CREATE TRIGGER %s %s %s ON %s FOR EACH ROW\n%s;\n\n",
		quoteIdent(trigger.Name), trigger.Timing, trigger.Event, quoteIdent(tableName), strings.TrimSpace(trigger.Body)); err != nil {
		return err
	}
	return nil
}

// writeEventSDL emits a scheduled event. The synced Definition is the full
// SHOW CREATE EVENT text, which (unlike functions/procedures) still carries the
// DEFINER clause, so it is stripped here to match omni's DEFINER-agnostic event
// identity. The schedule's auto-injected STARTS '<create-time>' is left intact: omni's
// event differ normalizes STARTS out of the canonical schedule key, so it does not
// perturb the no-op.
func writeEventSDL(buf *strings.Builder, event *metadata.EventMetadata) error {
	def := strings.TrimSpace(stripLeadingDefiner(event.Definition))
	if def == "" {
		return nil
	}
	if _, err := fmt.Fprintf(buf, "%s;\n\n", def); err != nil {
		return err
	}
	return nil
}

// stripLeadingDefiner removes a leading `CREATE DEFINER=<account> ` clause from a
// CREATE statement, collapsing it back to a bare `CREATE <OBJECT> …`. The account is
// parsed properly — `user`@`host` with each part backtick, single- or double-quoted
// (quoted parts may contain spaces, '@', and doubled/escaped quotes), an unquoted
// user@host, or CURRENT_USER[()] — so a legal quoted account with a space
// (DEFINER=`my user`@`%`) does not get cut at its first space. Statements without a
// DEFINER, or with an account this scan cannot parse, are returned unchanged (omni
// tolerates a DEFINER and never diffs on it, so keeping one is safe while corrupting
// the statement is not). Stripping keeps the dump canonical and matches the
// routine/trigger emission, which carry no DEFINER at all.
func stripLeadingDefiner(definition string) string {
	const createKW = "CREATE "
	trimmed := strings.TrimLeft(definition, " \t\r\n")
	if !strings.HasPrefix(strings.ToUpper(trimmed), createKW) {
		return definition
	}
	rest := strings.TrimLeft(trimmed[len(createKW):], " \t")
	const definerKW = "DEFINER"
	if !strings.HasPrefix(strings.ToUpper(rest), definerKW) {
		return definition
	}
	pos := len(definerKW)
	pos = skipSpacesAndTabs(rest, pos)
	if pos >= len(rest) || rest[pos] != '=' {
		return definition
	}
	pos = skipSpacesAndTabs(rest, pos+1)
	end, ok := scanDefinerAccount(rest, pos)
	if !ok {
		return definition
	}
	// The clause must be followed by whitespace and the object keyword.
	afterAccount := skipSpacesAndTabs(rest, end)
	if afterAccount == end || afterAccount >= len(rest) {
		return definition
	}
	return createKW + rest[afterAccount:]
}

// skipSpacesAndTabs returns the first offset >= pos in s that is not a space or tab.
func skipSpacesAndTabs(s string, pos int) int {
	for pos < len(s) && (s[pos] == ' ' || s[pos] == '\t') {
		pos++
	}
	return pos
}

// scanDefinerAccount scans a MySQL account name starting at pos: CURRENT_USER[()] or
// <part>[@<part>] where each part is backtick/single/double quoted (doubled-delimiter
// and, for '/'"', backslash escapes) or an unquoted run. Returns the offset just past
// the account and whether the scan succeeded.
func scanDefinerAccount(s string, pos int) (int, bool) {
	const currentUser = "CURRENT_USER"
	if len(s)-pos >= len(currentUser) && strings.EqualFold(s[pos:pos+len(currentUser)], currentUser) {
		end := pos + len(currentUser)
		if strings.HasPrefix(s[end:], "()") {
			end += 2
		}
		return end, true
	}
	end, ok := scanAccountPart(s, pos)
	if !ok {
		return 0, false
	}
	if end < len(s) && s[end] == '@' {
		hostEnd, ok := scanAccountPart(s, end+1)
		if !ok {
			return 0, false
		}
		end = hostEnd
	}
	return end, true
}

// scanAccountPart scans one user or host part of an account name starting at pos.
func scanAccountPart(s string, pos int) (int, bool) {
	if pos >= len(s) {
		return 0, false
	}
	switch q := s[pos]; q {
	case '`', '\'', '"':
		for i := pos + 1; i < len(s); {
			switch {
			case q != '`' && s[i] == '\\' && i+1 < len(s):
				// Backslash escapes apply inside '…' and "…" (not backticks).
				i += 2
			case s[i] == q:
				if i+1 < len(s) && s[i+1] == q {
					// Doubled delimiter stays inside the part.
					i += 2
					continue
				}
				return i + 1, true
			default:
				i++
			}
		}
		return 0, false // Unterminated quote.
	default:
		// Unquoted part: a run of bytes up to the '@' separator or whitespace.
		i := pos
		for i < len(s) && s[i] != '@' && s[i] != ' ' && s[i] != '\t' && s[i] != '\r' && s[i] != '\n' {
			i++
		}
		if i == pos {
			return 0, false
		}
		return i, true
	}
}
