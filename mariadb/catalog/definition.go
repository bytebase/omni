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

	nodes "github.com/bytebase/omni/mariadb/ast"
	mariadbparser "github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
)

const (
	// explicit_defaults_for_timestamp is ON while the script runs, so a TIMESTAMP
	// column takes only the default it is written with: OFF, which 10.7 has by
	// default, gives the first NOT NULL one DEFAULT current_timestamp() ON
	// UPDATE current_timestamp().
	disableUniqueAndForeignKeyCheckStmt = "SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;\nSET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;\nSET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1;\n"
	restoreUniqueAndForeignKeyCheckStmt = "SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP;\nSET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;\nSET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;\n"

	emptyCommentLine      = "--\n"
	tempViewHeader        = "-- Temporary view structure for "
	tableHeader           = "-- Table structure for "
	sequenceHeader        = "-- Sequence structure for "
	viewHeader            = "-- View structure for "
	functionHeader        = "-- Function structure for "
	procedureHeader       = "-- Procedure structure for "
	triggerHeader         = "-- Trigger structure for "
	eventHeader           = "-- Event structure for "
	setCharacterSetClient = "SET character_set_client = "
	setCharacterSetResult = "SET character_set_results = "
	setCollation          = "SET collation_connection = "
	setSQLMode            = "SET sql_mode = "
	setTimezone           = "SET time_zone = "
	// A function is created under its own session settings, which must not
	// carry over to the views written after it: a view's string literals, for
	// one, read differently under NO_BACKSLASH_ESCAPES.
	saveSessionStmt     = "SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode;\n"
	restoreSessionStmt  = "SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode;\n"
	delimiterDoubleSemi = "DELIMITER ;;\n"
	delimiterSemi       = "DELIMITER ;\n"

	mysqlIndexFullText = "FULLTEXT"
	mysqlIndexSpatial  = "SPATIAL"

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
// table comments are written between quotes as they are: sync stores them
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

	// Construct functions before views: CREATE VIEW checks that the functions it calls
	// exist, while a function body is checked only when it runs. Each function sets
	// the session it was created under, and the views after it take the script's.
	if len(schema.Functions) > 0 {
		if _, err := buf.WriteString(saveSessionStmt); err != nil {
			return "", err
		}
	}
	for _, function := range schema.Functions {
		if err := writeFunction(&buf, function); err != nil {
			return "", err
		}
	}
	if len(schema.Functions) > 0 {
		if _, err := buf.WriteString(restoreSessionStmt); err != nil {
			return "", err
		}
	}

	// Construct views.
	for _, view := range schema.Views {
		if err := writeView(&buf, view); err != nil {
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

func writeSequence(out *strings.Builder, sequence *metadata.SequenceMetadata) error {
	for _, s := range []string{emptyCommentLine, sequenceHeader, quoteIdent(sequence.Name) + "\n", emptyCommentLine} {
		if _, err := io.WriteString(out, s); err != nil {
			return err
		}
	}
	return writeCreateSequence(out, sequence)
}

// writeCreateSequence writes a sequence as SHOW CREATE SEQUENCE does, from the
// attributes sync reads off its row. AS appears only for a type other than
// bigint, which MariaDB accepts from 11.5 on.
func writeCreateSequence(out *strings.Builder, sequence *metadata.SequenceMetadata) error {
	fmt.Fprintf(out, "CREATE SEQUENCE %s", quoteIdent(sequence.Name))
	if dt := strings.TrimSpace(sequence.DataType); dt != "" && !strings.EqualFold(dt, "bigint") {
		fmt.Fprintf(out, " as %s", strings.ToLower(dt))
	}
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
	case "":
	case "0":
		out.WriteString(" nocache")
	default:
		fmt.Fprintf(out, " cache %s", sequence.CacheSize)
	}
	if sequence.Cycle {
		out.WriteString(" cycle")
	} else {
		out.WriteString(" nocycle")
	}
	if sequence.Comment != "" {
		fmt.Fprintf(out, " COMMENT='%s'", sequence.Comment)
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
	// An event synced without its CREATE EVENT, as on OceanBase, has nothing faithful to write.
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

	// Set charset, collation, sql mode and timezone when recorded on the snapshot.
	if err := writeAdditionalEventsIfSet(out, event.CharacterSetClient, event.CharacterSetClient, event.CollationConnection, event.SqlMode); err != nil {
		return err
	}
	if event.TimeZone != "" {
		if _, err := io.WriteString(out, setTimezone); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "'"); err != nil {
			return err
		}
		if _, err := io.WriteString(out, event.TimeZone); err != nil {
			return err
		}
		if _, err := io.WriteString(out, "';\n"); err != nil {
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

	// Set charset, collation, and sql mode when recorded on the snapshot.
	if err := writeAdditionalEventsIfSet(out, trigger.CharacterSetClient, trigger.CharacterSetClient, trigger.CollationConnection, trigger.SqlMode); err != nil {
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
	if _, err := io.WriteString(out, "CREATE OR REPLACE VIEW "+quoteIdent(view.Name)+viewColumnListClause(view)+" AS "); err != nil {
		return err
	}
	if _, err := io.WriteString(out, view.Definition); err != nil {
		return err
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

	columnChecks, tableChecks := splitColumnChecks(table)
	for i, column := range table.Columns {
		if i != 0 {
			if _, err := fmt.Fprint(buf, ",\n"); err != nil {
				return err
			}
		}
		if err := printColumnClause(buf, column, table, columnChecks[strings.ToLower(column.Name)], dbName); err != nil {
			return err
		}
	}
	writeSystemTimePeriod(buf, table)

	if err := printPrimaryKeyClause(buf, table); err != nil {
		return err
	}

	secondaryIndexes := make([]*metadata.IndexMetadata, 0, len(table.Indexes))
	for _, index := range table.Indexes {
		if !index.Primary {
			secondaryIndexes = append(secondaryIndexes, index)
		}
	}
	foreignKeys := slices.Clone(table.ForeignKeys)
	checks := slices.Clone(tableChecks)
	if sorted {
		slices.SortFunc(secondaryIndexes, func(a, b *metadata.IndexMetadata) int { return cmp.Compare(a.Name, b.Name) })
		slices.SortFunc(foreignKeys, func(a, b *metadata.ForeignKeyMetadata) int { return cmp.Compare(a.Name, b.Name) })
		slices.SortFunc(checks, func(a, b *metadata.CheckConstraintMetadata) int { return cmp.Compare(a.Name, b.Name) })
	}
	for _, index := range secondaryIndexes {
		if err := printIndexClause(buf, index, table); err != nil {
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
	writeSystemVersioningOption(buf, table)

	if len(table.Partitions) > 0 {
		if err := printPartitionClause(buf, table.Partitions, table.Engine); err != nil {
			return err
		}
	}

	_, err := io.WriteString(buf, ";\n\n")
	return err
}

func printPartitionClause(buf *strings.Builder, partitions []*metadata.TablePartitionMetadata, engine string) error {
	if len(partitions) == 0 {
		return nil
	}
	// MariaDB writes the clause plainly, as SHOW CREATE TABLE does: the parser
	// skips a MySQL /*!50100 */ comment, so a table loaded from one would lose
	// its partitioning.
	if _, err := fmt.Fprint(buf, "\nPARTITION BY "); err != nil {
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
		if _, err := fmt.Fprintf(buf, "KEY (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_LINEAR_HASH:
		if _, err := fmt.Fprintf(buf, "LINEAR HASH (%s)", partitions[0].Expression); err != nil {
			return err
		}
	case metadata.TablePartitionMetadata_LINEAR_KEY:
		if _, err := fmt.Fprintf(buf, "LINEAR KEY (%s)", partitions[0].Expression); err != nil {
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
			if _, err := fmt.Fprintf(buf, "KEY (%s)", partitions[0].Subpartitions[0].Expression); err != nil {
				return err
			}
		case metadata.TablePartitionMetadata_LINEAR_KEY:
			if _, err := fmt.Fprintf(buf, "LINEAR KEY (%s)", partitions[0].Subpartitions[0].Expression); err != nil {
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

	if useDefault == 0 {
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
			switch {
			case preposition == "IN" && partition.Value == "DEFAULT":
				// The catch-all partition of a LIST partitioning.
				if _, err := fmt.Fprint(buf, " DEFAULT"); err != nil {
					return err
				}
			case preposition != "":
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

			if subUseDefault == 0 && len(partition.Subpartitions) > 0 {
				if _, err := fmt.Fprint(buf, "\n ("); err != nil {
					return err
				}
				for j, subPartition := range partition.Subpartitions {
					if _, err := fmt.Fprintf(buf, "SUBPARTITION %s", quoteIdent(subPartition.Name)); err != nil {
						return err
					}
					if err := writePartitionOptions(buf, engine, subPartition.Comment); err != nil {
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
				if err := writePartitionOptions(buf, engine, partition.Comment); err != nil {
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

func writePartitionOptions(buf io.StringWriter, engine string, comment string) error {
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
	// TODO(zp): Get all the partition options from the metadata is too complex, just write
	// the engine for now. It must be the table's own engine: the server rejects a
	// partitioned table whose partitions name a different one with 1497.
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

// printCheckClause writes a table check. MariaDB's information_schema keeps
// the expression without the parentheses CHECK needs, and drops redundant ones
// when it stores a check, so they are always written.
func printCheckClause(buf *strings.Builder, check *metadata.CheckConstraintMetadata) error {
	_, err := fmt.Fprintf(buf, ",\n  CONSTRAINT %s CHECK (%s)", quoteIdent(check.Name), check.Expression)
	return err
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

	// RESTRICT is what an unwritten action is; NO ACTION is stored apart.
	if fk.OnDelete != "" && !strings.EqualFold(fk.OnDelete, "RESTRICT") {
		if _, err := fmt.Fprintf(buf, " ON DELETE %s", fk.OnDelete); err != nil {
			return err
		}
	}

	if fk.OnUpdate != "" && !strings.EqualFold(fk.OnUpdate, "RESTRICT") {
		if _, err := fmt.Fprintf(buf, " ON UPDATE %s", fk.OnUpdate); err != nil {
			return err
		}
	}

	return nil
}

func printIndexClause(buf *strings.Builder, index *metadata.IndexMetadata, table *metadata.TableMetadata) error {
	engine := table.GetEngine()
	if index.Primary {
		return nil
	}

	if _, err := fmt.Fprint(buf, ",\n  "); err != nil {
		return err
	}

	isFullText := strings.EqualFold(index.Type, mysqlIndexFullText)
	isSpatial := strings.EqualFold(index.Type, mysqlIndexSpatial)
	if index.Unique {
		if _, err := fmt.Fprint(buf, "UNIQUE "); err != nil {
			return err
		}
	} else if isFullText {
		if _, err := fmt.Fprint(buf, "FULLTEXT "); err != nil {
			return err
		}
	} else if isSpatial {
		if _, err := fmt.Fprint(buf, "SPATIAL "); err != nil {
			return err
		}
	} else if strings.EqualFold(index.Type, "VECTOR") {
		if _, err := fmt.Fprint(buf, "VECTOR "); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(buf, "KEY %s (", quoteIdent(index.Name)); err != nil {
		return err
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
		// SPATIAL and FULLTEXT index parts do not take a key-part prefix length. The sync
		// records SUB_PART for spatial parts (e.g. 32 for a POINT) even though SHOW CREATE
		// renders none; emitting `(32)` would diverge from the user form and DROP+ADD the
		// index on every no-op. Suppress the prefix for those index types.
		if err := printIndexKeyPart(buf, table, expr, keyLength, descending, isSpatial || isFullText); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(buf, ")", indexUsingClause(index.Type, engine)); err != nil {
		return err
	}

	if index.Comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT '%s'", escapeComment(index.Comment)); err != nil {
			return err
		}
	}

	// MariaDB's invisible index is an IGNORED one.
	if !index.Visible {
		if _, err := fmt.Fprint(buf, " IGNORED"); err != nil {
			return err
		}
	}

	return nil
}

func printIndexKeyPart(buf *strings.Builder, table *metadata.TableMetadata, expr string, length int64, descending bool, suppressPrefix bool) error {
	if len(expr) == 0 {
		return errors.New("index key part expression is empty")
	}
	// A column's name may itself be parenthesized, as `(n)` is.
	if expr[0] == '(' && expr[len(expr)-1] == ')' && !hasColumn(table, expr) {
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

func printPrimaryKeyClause(buf *strings.Builder, table *metadata.TableMetadata) error {
	indexes, engine := table.GetIndexes(), table.GetEngine()
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
				if err := printIndexKeyPart(buf, table, column, keyLength, descending, false); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(buf, ")", indexUsingClause(index.Type, engine)); err != nil {
				return err
			}
			if index.Comment != "" {
				if _, err := fmt.Fprintf(buf, " COMMENT '%s'", escapeComment(index.Comment)); err != nil {
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

// printColumnClause writes a column in the order MariaDB's SHOW CREATE TABLE
// does. check is the column's own check, which MariaDB writes on the column.
// dbName is the database being dumped, empty when unknown.
func printColumnClause(buf *strings.Builder, column *metadata.ColumnMetadata, table *metadata.TableMetadata, check *metadata.CheckConstraintMetadata, dbName string) error {
	if _, err := fmt.Fprintf(buf, "  %s %s", quoteIdent(column.Name), normalizeColumnType(column.Type)); err != nil {
		return err
	}

	tableCharset := table.Charset
	if tableCharset == "" {
		tableCharset = charsetForCollation(table.Collation)
	}
	if column.CharacterSet != "" && column.CharacterSet != tableCharset {
		if _, err := fmt.Fprintf(buf, " CHARACTER SET %s", column.CharacterSet); err != nil {
			return err
		}
	}

	if column.Collation != "" && column.Collation != table.Collation {
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
	} else if strings.HasPrefix(strings.ToLower(strings.TrimSpace(column.Type)), "timestamp") {
		// A nullable TIMESTAMP says NULL, as SHOW CREATE TABLE writes it, so a
		// session with explicit_defaults_for_timestamp off does not make it
		// NOT NULL.
		if _, err := fmt.Fprint(buf, " NULL"); err != nil {
			return err
		}
	}

	if err := printDefaultClause(buf, column, dbName); err != nil {
		return err
	}

	if column.OnUpdate != "" {
		if _, err := fmt.Fprintf(buf, " ON UPDATE %s", column.OnUpdate); err != nil {
			return err
		}
	}

	if isAutoIncrement(column) {
		if _, err := fmt.Fprintf(buf, " %s", autoIncrementSymbol); err != nil {
			return err
		}
	}

	if column.IsInvisible {
		if _, err := fmt.Fprint(buf, " INVISIBLE"); err != nil {
			return err
		}
	}

	if column.Comment != "" {
		if _, err := fmt.Fprintf(buf, " COMMENT '%s'", column.Comment); err != nil {
			return err
		}
	}

	if check != nil {
		if _, err := fmt.Fprintf(buf, " CHECK (%s)", check.Expression); err != nil {
			return err
		}
	}
	return nil
}

// normalizeColumnType spells MySQL 8.0's geomcollection, which a snapshot taken
// from MySQL may carry, as the geometrycollection MariaDB accepts.
func normalizeColumnType(columnType string) string {
	if strings.EqualFold(columnType, "geomcollection") {
		return "geometrycollection"
	}
	return columnType
}

func printDefaultClause(buf *strings.Builder, column *metadata.ColumnMetadata, dbName string) error {
	if column.Default == "" || isAutoIncrement(column) {
		return nil
	}
	if column.Generation != nil && column.Generation.Expression != "" {
		return nil
	}
	// MariaDB qualifies the sequence a default reads with its database; the
	// dump leaves that out so the default follows it into the database it
	// restores to.
	text := requalifySequences(renderColumnDefault(column), dbName, "")
	if strings.EqualFold(text, "NULL") && !column.Nullable {
		return nil
	}
	_, err := fmt.Fprintf(buf, " DEFAULT %s", text)
	return err
}

// renderColumnDefault returns the DEFAULT clause's text for a column. Sync
// records information_schema's COLUMN_DEFAULT, which is the default as SHOW
// CREATE TABLE writes it, inside the quotes QUOTE() adds (see columnDefault),
// so the quoted text is written as it is. MariaDB writes an expression other
// than a function call in parentheses, which the DEFAULT grammar needs.
func renderColumnDefault(column *metadata.ColumnMetadata) string {
	inner, expr, ok := syncedDefault(column.Default)
	if !ok {
		return column.Default
	}
	switch e := expr.(type) {
	case *nodes.StringLit, *nodes.IntLit, *nodes.FloatLit, *nodes.NullLit, *nodes.BoolLit,
		*nodes.HexLit, *nodes.BitLit, *nodes.TemporalLit, *nodes.FuncCallExpr:
		return inner
	case *nodes.UnaryExpr:
		switch e.Operand.(type) {
		case *nodes.IntLit, *nodes.FloatLit:
			return inner
		}
	}
	return "(" + inner + ")"
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
	// An empty sql_mode is a mode too: a synced object carries its character
	// set with it, and restores under exactly the mode it was created with.
	if sqlMode != "" || characterSetClient != "" {
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

	sequences := make([]*metadata.SequenceMetadata, 0, len(schema.Sequences))
	for _, sequence := range schema.Sequences {
		if !sequence.SkipDump {
			sequences = append(sequences, sequence)
		}
	}
	slices.SortFunc(sequences, func(a, b *metadata.SequenceMetadata) int { return cmp.Compare(a.Name, b.Name) })
	for _, sequence := range sequences {
		if err := writeCreateSequence(&buf, sequence); err != nil {
			return "", err
		}
	}

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
// `partitioned` token removed. information_schema.TABLES.CREATE_OPTIONS folds the
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
// table and column references of a view body.
//
// information_schema.VIEWS qualifies the references of a view body with their
// database, while users write them unqualified. The omni SDL diff folds only the
// qualifier of the database it loads into, so a body still naming the synced database
// would differ from the user's forever. The body is parsed so that only a qualifier in the database
// position goes: a table alias or table spelled like the database is a qualifier of
// its columns and stays. A body that does not parse is returned unchanged.
func stripViewBodyDatabaseQualifier(body, dbName string) string {
	return stripBodyQualifier(body, dbName, dbName)
}

// stripBodyQualifier is stripViewBodyDatabaseQualifier, with the function calls
// dbName qualifies qualified with functionDB instead. A call keeps a qualifier:
// MariaDB stores a call as it was written, and an unqualified one may name a
// builtin. Measured on mariadb:11.8.8: `db`.`upper`(1) stays the stored function
// and upper(1) is stored as ucase(1).
func stripBodyQualifier(body, dbName, functionDB string) string {
	if dbName == "" || body == "" {
		return body
	}
	qualifier := quoteIdent(dbName) + "."
	if !strings.Contains(body, qualifier) {
		return body
	}
	// A stored body may end in WITH CHECK OPTION, so it parses only as a whole view.
	const prefix = "CREATE VIEW v AS "
	list, err := mariadbparser.Parse(prefix + body)
	if err != nil || len(list.Items) != 1 {
		return body
	}
	replacements := make(map[int]string)
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		var schema, replacement string
		var loc nodes.Loc
		switch n := n.(type) {
		case *nodes.TableRef:
			schema, loc = n.Schema, n.Loc
		case *nodes.ColumnRef:
			schema, loc = n.Schema, n.Loc
		case *nodes.FuncCallExpr:
			if functionDB == "" || functionDB == dbName {
				return true
			}
			schema, loc, replacement = n.Schema, n.Loc, quoteIdent(functionDB)+"."
		default:
			return true
		}
		start := loc.Start - len(prefix)
		if schema == dbName && start >= 0 && strings.HasPrefix(body[start:], qualifier) {
			replacements[start] = replacement
		}
		return true
	})

	var out strings.Builder
	out.Grow(len(body))
	prev := 0
	for _, start := range slices.Sorted(maps.Keys(replacements)) {
		out.WriteString(body[prev:start])
		out.WriteString(replacements[start])
		prev = start + len(qualifier)
	}
	out.WriteString(body[prev:])
	return out.String()
}

// writeRoutineSDL emits a stored function or procedure. The synced Definition is
// already a complete CREATE FUNCTION/PROCEDURE statement (sync rewrites SHOW CREATE
// to drop the DEFINER and the RETURNS charset attribute), so the
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
// Sync stores only the action statement (TriggerMetadata.Body) plus the timing,
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

func hasColumn(table *metadata.TableMetadata, name string) bool {
	for _, col := range table.GetColumns() {
		if strings.EqualFold(col.GetName(), name) {
			return true
		}
	}
	return false
}
