package catalog

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mariadbparser "github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
)

// LoadMetadataReport lists the snapshot objects LoadMetadata could not install
// as the snapshot defines them. Keys are "table:<name>", "view:<name>",
// "function:<name>", "procedure:<name>", "trigger:<table>.<name>",
// "event:<name>", and "sequence:<name>".
type LoadMetadataReport struct {
	// Degraded maps each object replaced by a stand-in to the error its
	// definition failed with.
	Degraded map[string]error
	// Missing maps each object the catalog has nothing for to the error that
	// left it out.
	Missing map[string]error
	// Incomplete maps each table installed without part of its definition to
	// the error of the first part left out: a default, ON UPDATE, generation,
	// or CHECK expression, an index, or the partitioning.
	Incomplete map[string]error
}

// LoadMetadata installs a Bytebase schema snapshot of a MariaDB database into the
// database meta.Name, creating it if needed with the snapshot's character set
// and collation, and leaves that database selected.
//
// The settings that would rewrite what it installs are off for the length of
// the load, and the caller's are restored afterward: a table may reference one
// that installs after it, and each table installs as the snapshot records it. A
// view may likewise name a view that installs after it, since views whose
// bodies did not analyze are reinstalled once everything is in. Each routine,
// trigger, and event the load installs takes the session context the snapshot
// records for it, as ApplySessionContext stamps it.
//
// An object that fails to install does not stop the load. A table is replaced
// by a stand-in with the same column names, all TEXT; a view by one that
// selects NULL under each of its column names; a trigger by one with an empty
// body; and an event by one with only its name. A function or procedure that
// fails is left out. A table that installs without part of its definition is
// reported Incomplete. The error is non-nil only when ctx is done.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadata.DatabaseSchemaMetadata) (*LoadMetadataReport, error) {
	return c.LoadMetadataInto(ctx, meta, meta.GetName())
}

// LoadMetadataInto is LoadMetadata into the database named database rather
// than meta.Name, so that two snapshots of differently named databases can be
// loaded under one name and diffed. View bodies drop the qualifier of
// meta.Name, which sync puts on every reference.
func (c *Catalog) LoadMetadataInto(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, database string) (*LoadMetadataReport, error) {
	report := &LoadMetadataReport{
		Degraded:   make(map[string]error),
		Missing:    make(map[string]error),
		Incomplete: make(map[string]error),
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	c.selectMetadataDatabase(meta, database)
	restore := c.loadSessionSettings()
	defer restore()

	var views, installed []*metadataObject
	for _, o := range collectMetadataObjects(meta, c.CurrentDatabase()) {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		err := c.installMetadataObject(o)
		if err == nil {
			if o.incomplete != nil {
				report.Incomplete[o.key()] = o.incomplete
			}
			installed = append(installed, o)
			if o.kind == metadataView {
				views = append(views, o)
			}
			continue
		}
		switch standInErr := c.installStandIn(o); {
		case standInErr == nil:
			installed = append(installed, o)
			report.Degraded[o.key()] = err
		case errors.Is(standInErr, errNoStandIn):
			report.Missing[o.key()] = err
		default:
			report.Missing[o.key()] = standInErr
		}
	}
	if err := c.reanalyzeViews(ctx, views); err != nil {
		return report, err
	}
	if db := c.GetDatabase(c.CurrentDatabase()); db != nil {
		applyDatabaseSessionContext(db, sessionContexts(installed))
	}
	return report, nil
}

// selectMetadataDatabase creates the database named name where the catalog
// lacks it, gives it the snapshot's character set and collation, and selects
// it.
func (c *Catalog) selectMetadataDatabase(meta *metadata.DatabaseSchemaMetadata, name string) {
	if name == "" {
		return
	}
	stmt := metadataDatabaseStmt(meta, name)
	if c.GetDatabase(name) == nil {
		_ = c.DefineDatabase(stmt)
	} else if len(stmt.Options) > 0 {
		_ = c.alterDatabase(&nodes.AlterDatabaseStmt{Name: name, Options: stmt.Options})
	}
	c.SetCurrentDatabase(name)
}

// loadSessionSettings turns off the session settings that rewrite what a table
// installs as: foreign key checks reject a forward reference, and
// explicit_defaults_for_timestamp=OFF rewrites a TIMESTAMP column. It returns a
// function restoring the caller's.
func (c *Catalog) loadSessionSettings() func() {
	fkChecks := c.ForeignKeyChecks()
	explicitDefaults := c.session.ExplicitDefaultsForTimestamp
	c.SetForeignKeyChecks(false)
	c.session.ExplicitDefaultsForTimestamp = true
	return func() {
		c.SetForeignKeyChecks(fkChecks)
		c.session.ExplicitDefaultsForTimestamp = explicitDefaults
	}
}

// sessionContexts collects the creation context the snapshot records for each
// routine, trigger, and event in objects. An object with none recorded is left
// out, so it keeps a bare recreate.
func sessionContexts(objects []*metadataObject) SessionContextMap {
	m := SessionContextMap{
		Functions:  make(map[string]SessionContext),
		Procedures: make(map[string]SessionContext),
		Triggers:   make(map[string]SessionContext),
		Events:     make(map[string]SessionContext),
	}
	add := func(into map[string]SessionContext, name string, ctx SessionContext) {
		if ctx != (SessionContext{}) {
			into[toLower(name)] = ctx
		}
	}
	for _, o := range objects {
		switch o.kind {
		case metadataFunction:
			add(m.Functions, o.name, SessionContext{SQLMode: o.function.GetSqlMode(), CharacterSetClient: o.function.GetCharacterSetClient(), CollationConnection: o.function.GetCollationConnection()})
		case metadataProcedure:
			add(m.Procedures, o.name, SessionContext{SQLMode: o.procedure.GetSqlMode(), CharacterSetClient: o.procedure.GetCharacterSetClient(), CollationConnection: o.procedure.GetCollationConnection()})
		case metadataTrigger:
			add(m.Triggers, o.name, SessionContext{SQLMode: o.trigger.GetSqlMode(), CharacterSetClient: o.trigger.GetCharacterSetClient(), CollationConnection: o.trigger.GetCollationConnection()})
		case metadataEvent:
			add(m.Events, o.name, SessionContext{SQLMode: o.event.GetSqlMode(), CharacterSetClient: o.event.GetCharacterSetClient(), CollationConnection: o.event.GetCollationConnection(), TimeZone: o.event.GetTimeZone()})
		}
	}
	return m
}

// reanalyzeViews reinstalls each view whose body did not analyze while a pass
// analyzes more of them: a body analyzes only once the views it names are
// installed, and those may install after it.
func (c *Catalog) reanalyzeViews(ctx context.Context, views []*metadataObject) error {
	for progress := true; progress; {
		progress = false
		for _, o := range views {
			if c.viewAnalyzed(o.name) {
				continue
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			if c.installMetadataObject(o) == nil && c.viewAnalyzed(o.name) {
				progress = true
			}
		}
	}
	return nil
}

func (c *Catalog) viewAnalyzed(name string) bool {
	v := c.installedView(name)
	return v != nil && v.AnalyzedQuery != nil
}

// installedView returns the view of that name in the database being loaded.
func (c *Catalog) installedView(name string) *View {
	db := c.GetDatabase(c.CurrentDatabase())
	if db == nil {
		return nil
	}
	return db.Views[toLower(name)]
}

// --- Objects ---

// metadataObjectKind lists the kinds in the order they install. Only a
// trigger needs another object first, its table: foreign key checks are off, a
// sequence a default reads is in place before its table,
// a view may name an object not yet installed, and routine bodies are not
// analyzed.
type metadataObjectKind int

const (
	metadataSequence metadataObjectKind = iota
	metadataTable
	metadataView
	metadataFunction
	metadataProcedure
	metadataTrigger
	metadataEvent
)

// metadataObject is one install step. The metadata pointer matching kind is
// set.
type metadataObject struct {
	kind metadataObjectKind
	name string
	// parent is a trigger's table.
	parent string
	// snapshotDatabase is the database the snapshot's references qualify
	// with: a view's body, and a table's foreign keys.
	snapshotDatabase string
	// database is the database the object installs into.
	database string

	sequence  *metadata.SequenceMetadata
	table     *metadata.TableMetadata
	view      *metadata.ViewMetadata
	function  *metadata.FunctionMetadata
	procedure *metadata.ProcedureMetadata
	trigger   *metadata.TriggerMetadata
	event     *metadata.EventMetadata

	// incomplete is the first part of a table's definition its install left out.
	incomplete error
}

// key names the object in the report.
func (o *metadataObject) key() string {
	switch o.kind {
	case metadataSequence:
		return "sequence:" + o.name
	case metadataTable:
		return "table:" + o.name
	case metadataView:
		return "view:" + o.name
	case metadataFunction:
		return "function:" + o.name
	case metadataProcedure:
		return "procedure:" + o.name
	case metadataTrigger:
		return "trigger:" + o.parent + "." + o.name
	default:
		return "event:" + o.name
	}
}

// --- Collecting ---

// collectMetadataObjects lists the snapshot's objects in install order: by
// kind, then name.
func collectMetadataObjects(meta *metadata.DatabaseSchemaMetadata, database string) []*metadataObject {
	var out []*metadataObject
	for _, s := range meta.GetSchemas() {
		for _, seq := range s.GetSequences() {
			if seq.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataSequence, name: seq.GetName(), sequence: seq})
			}
		}
		for _, t := range s.GetTables() {
			if t.GetName() == "" {
				continue
			}
			out = append(out, &metadataObject{kind: metadataTable, name: t.GetName(), table: t, snapshotDatabase: meta.GetName(), database: database})
			for _, trg := range t.GetTriggers() {
				if trg.GetName() != "" {
					out = append(out, &metadataObject{kind: metadataTrigger, name: trg.GetName(), parent: t.GetName(), trigger: trg})
				}
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataView, name: v.GetName(), view: v, snapshotDatabase: meta.GetName(), database: database})
			}
		}
		for _, fn := range s.GetFunctions() {
			if fn.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataFunction, name: fn.GetName(), function: fn})
			}
		}
		for _, p := range s.GetProcedures() {
			if p.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataProcedure, name: p.GetName(), procedure: p})
			}
		}
		for _, e := range s.GetEvents() {
			if e.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataEvent, name: e.GetName(), event: e})
			}
		}
	}
	slices.SortFunc(out, func(a, b *metadataObject) int {
		return cmp.Or(cmp.Compare(a.kind, b.kind), cmp.Compare(a.name, b.name), cmp.Compare(a.parent, b.parent))
	})
	return out
}

// --- Installing ---

func (c *Catalog) installMetadataObject(o *metadataObject) error {
	switch o.kind {
	case metadataSequence:
		stmt, err := metadataSequenceStmt(o.sequence)
		if err != nil {
			return err
		}
		return c.createSequence(stmt)
	case metadataTable:
		stmt, omitted, err := metadataTableStmt(o.table, o.snapshotDatabase, o.database)
		if err != nil {
			return err
		}
		o.incomplete = omitted
		return c.DefineTable(stmt)
	case metadataView:
		stmt, err := metadataViewStmt(o.view, o.snapshotDatabase, o.database)
		if err != nil {
			return err
		}
		if err := c.DefineView(stmt); err != nil {
			return err
		}
		nameViewColumns(c.installedView(o.name), metadataViewColumnNames(o.view))
		return nil
	case metadataFunction:
		stmt, err := parseRoutine(o.function.GetDefinition(), false)
		if err != nil {
			return err
		}
		return c.DefineFunction(stmt)
	case metadataProcedure:
		stmt, err := parseRoutine(o.procedure.GetDefinition(), true)
		if err != nil {
			return err
		}
		return c.DefineProcedure(stmt)
	case metadataTrigger:
		if o.trigger.GetTiming() == "" || o.trigger.GetEvent() == "" {
			return errors.New("trigger has no timing or event")
		}
		return c.DefineTrigger(metadataTriggerStmt(o.trigger, o.parent))
	default:
		stmt, err := parseFirstStatement[*nodes.CreateEventStmt](o.event.GetDefinition())
		if err != nil {
			return err
		}
		return c.DefineEvent(stmt)
	}
}

// --- Statements ---

// The snapshot carries types, defaults, and definitions as the text MariaDB
// reports for them. These builders turn that text into the statement ASTs the
// Define* functions take.

// autoIncrementDefault is what sync stores as the default of an
// AUTO_INCREMENT column.
const autoIncrementDefault = "AUTO_INCREMENT"

// sequenceLiteral is a sequence option value as the parser builds it, its
// digits kept whole where they pass int64, as BIGINT UNSIGNED limits do.
func sequenceLiteral(text string) (nodes.ExprNode, error) {
	text = strings.TrimSpace(text)
	digits, negative := strings.CutPrefix(text, "-")
	if digits == "" || strings.Trim(digits, "0123456789") != "" {
		return nil, fmt.Errorf("sequence value %q is no integer", text)
	}
	var lit nodes.ExprNode = &nodes.FloatLit{Value: digits}
	if n, err := strconv.ParseInt(digits, 10, 64); err == nil {
		lit = &nodes.IntLit{Value: n}
	}
	if negative {
		return &nodes.UnaryExpr{Op: nodes.UnaryMinus, Operand: lit}, nil
	}
	return lit, nil
}

// metadataSequenceStmt builds a sequence from the attributes sync reads off
// its own row.
func metadataSequenceStmt(seq *metadata.SequenceMetadata) (*nodes.CreateSequenceStmt, error) {
	stmt := &nodes.CreateSequenceStmt{Name: &nodes.TableRef{Name: seq.GetName()}}
	if dt := strings.TrimSpace(seq.GetDataType()); dt != "" && !strings.EqualFold(dt, "bigint") {
		name, unsigned := strings.CutSuffix(strings.ToLower(dt), " unsigned")
		stmt.DataType = &nodes.DataType{Name: name, Unsigned: unsigned}
	}
	for _, v := range []struct {
		text string
		into *nodes.ExprNode
	}{
		{seq.GetStart(), &stmt.Start},
		{seq.GetIncrement(), &stmt.Increment},
		{seq.GetMinValue(), &stmt.MinValue},
		{seq.GetMaxValue(), &stmt.MaxValue},
		{seq.GetCacheSize(), &stmt.Cache},
	} {
		if v.text == "" {
			continue
		}
		expr, err := sequenceLiteral(v.text)
		if err != nil {
			return nil, err
		}
		*v.into = expr
	}
	cycle := seq.GetCycle()
	stmt.Cycle = &cycle
	if seq.GetComment() != "" {
		stmt.Options = append(stmt.Options, &nodes.TableOption{Name: "COMMENT", Value: syncedComment(seq.GetComment())})
	}
	return stmt, nil
}

func metadataDatabaseStmt(meta *metadata.DatabaseSchemaMetadata, name string) *nodes.CreateDatabaseStmt {
	stmt := &nodes.CreateDatabaseStmt{Name: name}
	for _, opt := range []*nodes.DatabaseOption{
		{Name: "CHARACTER SET", Value: meta.GetCharacterSet()},
		{Name: "COLLATE", Value: meta.GetCollation()},
	} {
		if opt.Value != "" {
			stmt.Options = append(stmt.Options, opt)
		}
	}
	return stmt
}

// metadataTableStmt fails on a column whose type does not parse, so the table
// gets a stand-in rather than a partial definition. A default, ON UPDATE,
// generation, or CHECK expression, or a partitioning that does not parse is
// left out instead, and an index key expression stays a column so the index
// still installs; omitted is the first such error. References to the
// snapshot's database, in foreign keys and in the sequence a default reads,
// move to database.
func metadataTableStmt(t *metadata.TableMetadata, snapshotDatabase, database string) (stmt *nodes.CreateTableStmt, omitted error, err error) {
	omit := func(e error) {
		if omitted == nil {
			omitted = e
		}
	}
	stmt = &nodes.CreateTableStmt{Table: &nodes.TableRef{Name: t.GetName()}}
	columnChecks, tableChecks := splitColumnChecks(t)
	for _, col := range t.GetColumns() {
		if col.GetName() == "" {
			continue
		}
		def, colOmitted, err := metadataColumnDef(col, columnChecks[strings.ToLower(col.GetName())], snapshotDatabase, database)
		if err != nil {
			return nil, nil, fmt.Errorf("column %q: %w", col.GetName(), err)
		}
		if colOmitted != nil {
			omit(fmt.Errorf("column %q: %w", col.GetName(), colOmitted))
		}
		stmt.Columns = append(stmt.Columns, def)
	}
	for _, idx := range t.GetIndexes() {
		if len(idx.GetExpressions()) == 0 {
			continue
		}
		con, err := metadataIndexConstraint(idx, t)
		if err != nil {
			omit(fmt.Errorf("index %q: %w", idx.GetName(), err))
		}
		stmt.Constraints = append(stmt.Constraints, con)
	}
	for _, fk := range t.GetForeignKeys() {
		stmt.Constraints = append(stmt.Constraints, metadataForeignKey(fk, snapshotDatabase))
	}
	for _, check := range tableChecks {
		expr, err := parseExpr(check.GetExpression())
		if err != nil {
			omit(fmt.Errorf("check %q: %w", check.GetName(), err))
			continue
		}
		stmt.Constraints = append(stmt.Constraints, &nodes.Constraint{Type: nodes.ConstrCheck, Name: check.GetName(), Expr: expr})
	}
	if len(t.GetPartitions()) > 0 {
		partitions, err := metadataPartitionClause(t)
		if err != nil {
			omit(fmt.Errorf("partitioning: %w", err))
		}
		stmt.Partitions = partitions
	}
	charset := t.GetCharset()
	if charset == "" && t.GetCollation() != "" {
		// Sync leaves the character set out on 11.x, where it cannot join a
		// UCA 14.0.0 collation to its character set.
		charset = charsetForCollation(t.GetCollation())
	}
	for _, opt := range []struct{ name, value string }{
		{"ENGINE", t.GetEngine()},
		{"CHARSET", charset},
		{"COLLATE", t.GetCollation()},
		{"COMMENT", syncedComment(t.GetComment())},
	} {
		if opt.value != "" {
			stmt.Options = append(stmt.Options, &nodes.TableOption{Name: opt.name, Value: opt.value})
		}
	}
	createOptions, err := metadataCreateOptions(t.GetCreateOptions())
	if err != nil {
		omit(fmt.Errorf("create options: %w", err))
	}
	stmt.Options = append(stmt.Options, createOptions...)
	// Sync records a system-versioning period's columns as ROW START and ROW
	// END generations and not the period itself, which they declare.
	for _, def := range stmt.Columns {
		switch {
		case def.Generated == nil:
		case def.Generated.RowBound == nodes.RowBoundStart:
			stmt.PeriodStartCol = def.Name
		case def.Generated.RowBound == nodes.RowBoundEnd:
			stmt.PeriodEndCol = def.Name
		}
	}
	return stmt, omitted, nil
}

// syncedComment is a table's, column's or sequence's comment as sync records
// it: escaped by QUOTE(), without the quotes around it.
func syncedComment(comment string) string {
	if raw, ok := unquoteSQLString("'" + comment + "'"); ok {
		return raw
	}
	return comment
}

// metadataCreateOptions parses information_schema.TABLES.CREATE_OPTIONS, which
// sync keeps with WITH SYSTEM VERSIONING appended for a system-versioned
// table: "stats_persistent=1 row_format=DYNAMIC `PAGE_COMPRESSED`='1'". The
// word "partitioned" marks a partitioned table, whose partitioning the
// snapshot holds on its own.
func metadataCreateOptions(options string) ([]*nodes.TableOption, error) {
	var kept []string
	for _, word := range strings.Fields(options) {
		if !strings.EqualFold(word, "partitioned") {
			kept = append(kept, word)
		}
	}
	if len(kept) == 0 {
		return nil, nil
	}
	text := strings.Join(kept, " ")
	stmt, err := parseFirstStatement[*nodes.CreateTableStmt]("CREATE TABLE `probe` (`c` int) " + text)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", text, err)
	}
	return stmt.Options, nil
}

// splitColumnChecks tells a table's column checks from its table checks.
// MariaDB names a check declared on a column after its column, and the
// json_valid check of a JSON column too, while a table check named after a
// column is rare: sync records both kinds alike, so a check that shares its
// name with a column is taken as that column's.
func splitColumnChecks(t *metadata.TableMetadata) (column map[string]*metadata.CheckConstraintMetadata, table []*metadata.CheckConstraintMetadata) {
	names := make(map[string]bool)
	for _, col := range t.GetColumns() {
		names[strings.ToLower(col.GetName())] = true
	}
	column = make(map[string]*metadata.CheckConstraintMetadata)
	for _, check := range t.GetCheckConstraints() {
		key := strings.ToLower(check.GetName())
		if names[key] && column[key] == nil {
			column[key] = check
			continue
		}
		table = append(table, check)
	}
	return column, table
}

// metadataPartitionClause parses the partitioning the snapshot records, as the
// definition renderer prints it.
func metadataPartitionClause(t *metadata.TableMetadata) (*nodes.PartitionClause, error) {
	var buf strings.Builder
	if err := printPartitionClause(&buf, t.GetPartitions(), t.GetEngine()); err != nil {
		return nil, err
	}
	stmt, err := parseFirstStatement[*nodes.CreateTableStmt]("CREATE TABLE `probe` (`c` int)" + buf.String())
	if err != nil {
		return nil, err
	}
	if stmt.Partitions == nil {
		return nil, fmt.Errorf("no partitioning parsed from %q", buf.String())
	}
	return stmt.Partitions, nil
}

// metadataColumnDef fails on a type that does not parse. A default, ON UPDATE,
// generation, or check expression that does not parse is left out, and
// omitted is the first such error. check is the column's own check, if any.
func metadataColumnDef(col *metadata.ColumnMetadata, check *metadata.CheckConstraintMetadata, snapshotDatabase, database string) (def *nodes.ColumnDef, omitted error, err error) {
	typeName, err := parseTypeName(col.GetType())
	if err != nil {
		return nil, nil, err
	}
	// The type text may carry its own character set and collation.
	if typeName.Charset == "" {
		typeName.Charset = col.GetCharacterSet()
	}
	if typeName.Collate == "" {
		typeName.Collate = col.GetCollation()
	}
	if isJSONColumn(col, check) {
		// MariaDB keeps JSON as this longtext and its json_valid check, and
		// declaring it JSON again installs both as the server holds them.
		typeName = &nodes.DataType{Name: "JSON"}
		check = nil
	}
	def = &nodes.ColumnDef{Name: col.GetName(), TypeName: typeName, Comment: syncedComment(col.GetComment())}
	switch {
	case !col.GetNullable():
		def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrNotNull})
	case strings.EqualFold(typeName.Name, "timestamp"):
		// A nullable TIMESTAMP says NULL, as SHOW CREATE TABLE prints it, so
		// explicit_defaults_for_timestamp=OFF cannot read it as NOT NULL.
		def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrNull})
	}
	if col.GetIsInvisible() {
		def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrInvisible})
	}
	def.AutoIncrement = strings.EqualFold(col.GetDefault(), autoIncrementDefault)
	def.DefaultValue, omitted = columnDefault(col, snapshotDatabase, database)
	if col.GetOnUpdate() != "" {
		if expr, err := parseExpr(col.GetOnUpdate()); err == nil {
			def.OnUpdate = expr
		} else if omitted == nil {
			omitted = fmt.Errorf("ON UPDATE: %w", err)
		}
	}
	if gen := col.GetGeneration(); gen.GetExpression() != "" {
		switch bound := strings.ToUpper(strings.TrimSpace(gen.GetExpression())); bound {
		case "ROW START", "ROW END":
			// Sync reports a system-versioning period column as a stored
			// generated column whose expression is its bound.
			def.Generated = &nodes.GeneratedColumn{RowBound: nodes.RowBoundStart}
			if bound == "ROW END" {
				def.Generated.RowBound = nodes.RowBoundEnd
			}
		default:
			if expr, err := parseExpr(gen.GetExpression()); err == nil {
				def.Generated = &nodes.GeneratedColumn{Expr: expr, Stored: gen.GetType() == metadata.GenerationMetadata_TYPE_STORED}
			} else if omitted == nil {
				omitted = fmt.Errorf("generation: %w", err)
			}
		}
	}
	if check != nil {
		if expr, err := parseExpr(check.GetExpression()); err == nil {
			def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrCheck, Expr: expr})
		} else if omitted == nil {
			omitted = fmt.Errorf("check: %w", err)
		}
	}
	return def, omitted, nil
}

// isJSONColumn reports whether a column is the longtext a JSON column becomes:
// utf8mb4_bin, with the json_valid check named after it.
func isJSONColumn(col *metadata.ColumnMetadata, check *metadata.CheckConstraintMetadata) bool {
	if check == nil || !strings.EqualFold(strings.TrimSpace(col.GetType()), "longtext") || !strings.EqualFold(col.GetCollation(), "utf8mb4_bin") {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(check.GetExpression()), "json_valid("+quoteIdent(col.GetName())+")")
}

// columnDefault parses the default the snapshot records for a column, or
// returns nil where the catalog takes none: AUTO_INCREMENT is a flag rather
// than a default, and a generated column has no default. MariaDB shows DEFAULT
// NULL on every nullable column without a default, BLOB, TEXT, JSON and
// spatial ones included, which is what sync reports for them.
//
// Sync records MariaDB's information_schema.COLUMNS.COLUMN_DEFAULT, which is
// the default as SHOW CREATE TABLE writes it ('x', 1.50, NULL, curdate(),
// nextval(`db`.`s`)), wrapped in quotes by QUOTE(); only a CURRENT_TIMESTAMP
// form, and the NULL of a column with no default, stay bare. A quoted default
// is therefore read as the expression it quotes, and as the string itself
// where that does not parse.
func columnDefault(col *metadata.ColumnMetadata, snapshotDatabase, database string) (nodes.ExprNode, error) {
	value := col.GetDefault()
	switch {
	case value == "" || col.GetGeneration() != nil:
		return nil, nil
	case strings.EqualFold(value, autoIncrementDefault):
		return nil, nil
	}
	if inner, expr, ok := syncedDefault(value); ok {
		if rebased := requalifySequences(inner, snapshotDatabase, database); rebased != inner {
			if e, err := parseExpr(rebased); err == nil {
				expr = e
			}
		}
		return expressionDefault(expr), nil
	}
	expr, err := parseExpr(value)
	if err != nil {
		return nil, fmt.Errorf("default: %w", err)
	}
	return expressionDefault(expr), nil
}

// expressionDefault parenthesizes a default that is not a literal, as DEFAULT
// (expr) is written, so the column takes it as an expression: parseExpr drops
// the parentheses COLUMN_DEFAULT gives `a` + 1, and gives none to `a`.
func expressionDefault(expr nodes.ExprNode) nodes.ExprNode {
	switch e := expr.(type) {
	case *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.BoolLit, *nodes.NullLit,
		*nodes.HexLit, *nodes.BitLit, *nodes.TemporalLit, *nodes.ParenExpr:
		return expr
	case *nodes.UnaryExpr:
		switch e.Operand.(type) {
		case *nodes.IntLit, *nodes.FloatLit:
			return expr
		}
	}
	return &nodes.ParenExpr{Expr: expr}
}

// syncedDefault reads a default sync recorded as QUOTE(COLUMN_DEFAULT): the
// default's own text and its expression. ok is false for a default in any
// other form, such as the quoted string a hand-written snapshot gives a string
// default ('hello'). COLUMN_DEFAULT backquotes every column it names, so text
// that parses only by reading a bare word as a column is such a string.
func syncedDefault(value string) (string, nodes.ExprNode, bool) {
	inner, ok := unquoteSQLString(value)
	if !ok {
		return "", nil, false
	}
	expr, err := parseExpr(inner)
	if err != nil {
		return "", nil, false
	}
	bare := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if ref, ok := n.(*nodes.ColumnRef); ok && !strings.Contains(inner, quoteIdent(ref.Column)) {
			bare = true
		}
		return !bare
	})
	if bare {
		return "", nil, false
	}
	return inner, expr, true
}

// unquoteSQLString undoes QUOTE(): s between single quotes, with \\, \',
// NUL (\0) and Ctrl-Z (\Z) escaped. ok is false for anything else.
func unquoteSQLString(s string) (string, bool) {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return "", false
	}
	var b strings.Builder
	body := s[1 : len(s)-1]
	for i := 0; i < len(body); i++ {
		c := body[i]
		switch {
		case c == '\\' && i+1 < len(body):
			i++
			switch body[i] {
			case '0':
				b.WriteByte(0)
			case 'Z':
				b.WriteByte(26)
			default:
				b.WriteByte(body[i])
			}
		case c == '\'':
			if i+1 < len(body) && body[i+1] == '\'' {
				i++
				b.WriteByte('\'')
				continue
			}
			return "", false
		default:
			b.WriteByte(c)
		}
	}
	return b.String(), true
}

// requalifySequences rewrites the database the sequence calls of a default
// name, nextval(`from`.`s`) and the like, to the database to, or leaves it out
// where to is empty. MariaDB stores those calls qualified. Only a call's
// sequence changes: a string that spells the qualifier keeps it.
func requalifySequences(expr, from, to string) string {
	if from == "" || from == to || !strings.Contains(expr, quoteIdent(from)+".") {
		return expr
	}
	const prefix = "SELECT ("
	list, err := mariadbparser.Parse(prefix + expr + ")")
	if err != nil || len(list.Items) != 1 {
		return expr
	}
	qualifier := quoteIdent(from) + "."
	var starts []int
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		call, ok := n.(*nodes.FuncCallExpr)
		if !ok || call.Schema != "" || len(call.Args) == 0 {
			return true
		}
		switch strings.ToUpper(call.Name) {
		case "NEXTVAL", "LASTVAL", "SETVAL":
		default:
			return true
		}
		if ref, ok := call.Args[0].(*nodes.ColumnRef); ok && ref.Table == from {
			if start := ref.Loc.Start - len(prefix); start >= 0 && strings.HasPrefix(expr[start:], qualifier) {
				starts = append(starts, start)
			}
		}
		return true
	})
	replacement := ""
	if to != "" {
		replacement = quoteIdent(to) + "."
	}
	var out strings.Builder
	prev := 0
	for _, start := range starts {
		out.WriteString(expr[prev:start])
		out.WriteString(replacement)
		prev = start + len(qualifier)
	}
	out.WriteString(expr[prev:])
	return out.String()
}

// metadataIndexConstraint builds an index from its key parts. The index type
// holds either the kind, FULLTEXT or SPATIAL, or the access method. Sync reports
// the access method a key without USING gets too, so the engine's default one
// stays unwritten. The error is a key expression that did not parse; the index
// still comes back, with that key part as a column.
func metadataIndexConstraint(idx *metadata.IndexMetadata, t *metadata.TableMetadata) (*nodes.Constraint, error) {
	engine := t.GetEngine()
	indexType := strings.ToUpper(strings.TrimSpace(idx.GetType()))
	if indexType == defaultIndexType(engine) {
		indexType = ""
	}
	c := &nodes.Constraint{Name: idx.GetName()}
	switch {
	case idx.GetPrimary():
		c.Type, c.IndexType = nodes.ConstrPrimaryKey, indexType
	case idx.GetUnique():
		c.Type, c.IndexType = nodes.ConstrUnique, indexType
	case indexType == "FULLTEXT":
		c.Type = nodes.ConstrFulltextIndex
	case indexType == "SPATIAL":
		c.Type = nodes.ConstrSpatialIndex
	case indexType == "VECTOR":
		c.Type = nodes.ConstrVector
	default:
		c.Type, c.IndexType = nodes.ConstrIndex, indexType
	}
	var err error
	c.Columns, c.IndexColumns, err = indexKeyParts(idx, t.GetColumns(), c.Type == nodes.ConstrSpatialIndex)
	// MariaDB's invisible index is an IGNORED one; a primary key cannot be.
	if !idx.GetVisible() && !idx.GetPrimary() {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "IGNORED"})
	}
	if idx.GetComment() != "" {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "COMMENT", Value: &nodes.StringLit{Value: idx.GetComment()}})
	}
	return c, err
}

// indexKeyParts splits an index's key parts: plain column names, or, where a
// part is an expression or has a length or descends, an IndexColumn for each.
// A SPATIAL key takes no length, though sync reports one (32) for it.
func indexKeyParts(idx *metadata.IndexMetadata, tableColumns []*metadata.ColumnMetadata, spatial bool) ([]string, []*nodes.IndexColumn, error) {
	columns := make(map[string]bool)
	for _, col := range tableColumns {
		columns[toLower(col.GetName())] = true
	}
	// A column's name may hold parentheses, as `price(usd)` does.
	isExpression := func(expr string) bool {
		return isExpressionKey(expr) && !columns[toLower(unquoteIdent(expr))]
	}
	length := func(i int) int {
		if !spatial && i < len(idx.GetKeyLength()) {
			return int(idx.GetKeyLength()[i])
		}
		return 0
	}
	descending := func(i int) bool { return i < len(idx.GetDescending()) && idx.GetDescending()[i] }
	plain := true
	for i, expr := range idx.GetExpressions() {
		if isExpression(expr) || length(i) > 0 || descending(i) {
			plain = false
		}
	}
	if plain {
		var columns []string
		for _, expr := range idx.GetExpressions() {
			columns = append(columns, unquoteIdent(expr))
		}
		return columns, nil, nil
	}
	var keys []*nodes.IndexColumn
	var firstErr error
	for i, expr := range idx.GetExpressions() {
		key := &nodes.IndexColumn{Expr: &nodes.ColumnRef{Column: unquoteIdent(expr)}, Desc: descending(i)}
		if l := length(i); l > 0 {
			key.Length = l
		}
		if isExpression(expr) {
			// Sync reports the expression as information_schema stores it,
			// with its quotes escaped. Charset introducers stay: they can
			// change what the expression computes. On a parse failure the key
			// stays a column, so the index still installs.
			if parsed, err := parseExpr(unescapeFunctionalIndexQuotes(expr)); err == nil {
				key.Expr = parsed
			} else if firstErr == nil {
				firstErr = err
			}
		}
		keys = append(keys, key)
	}
	return nil, keys, firstErr
}

// defaultIndexType is the access method a key without USING gets.
func defaultIndexType(engine string) string {
	if strings.EqualFold(engine, "MEMORY") || strings.EqualFold(engine, "HEAP") {
		return "HASH"
	}
	return "BTREE"
}

// isExpressionKey reports whether a key part is a functional expression, as
// in "(lower(name))" or "lower(`name`)", rather than a column.
func isExpressionKey(s string) bool {
	return strings.Contains(s, "(")
}

// metadataForeignKey leaves a reference to the snapshot's own database
// unqualified, so it follows the snapshot into the database it installs into.
func metadataForeignKey(fk *metadata.ForeignKeyMetadata, snapshotDatabase string) *nodes.Constraint {
	refSchema := fk.GetReferencedSchema()
	if refSchema == snapshotDatabase {
		refSchema = ""
	}
	return &nodes.Constraint{
		Type:       nodes.ConstrForeignKey,
		Name:       fk.GetName(),
		Columns:    fk.GetColumns(),
		RefTable:   &nodes.TableRef{Schema: refSchema, Name: fk.GetReferencedTable()},
		RefColumns: fk.GetReferencedColumns(),
		OnUpdate:   referenceAction(fk.GetOnUpdate()),
		OnDelete:   referenceAction(fk.GetOnDelete()),
		Match:      strings.ToUpper(fk.GetMatchType()),
	}
}

func referenceAction(action string) nodes.ReferenceAction {
	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "RESTRICT":
		return nodes.RefActRestrict
	case "CASCADE":
		return nodes.RefActCascade
	case "SET NULL":
		return nodes.RefActSetNull
	case "SET DEFAULT":
		return nodes.RefActSetDefault
	case "NO ACTION":
		return nodes.RefActNoAction
	default:
		return nodes.RefActNone
	}
}

// metadataViewStmt keeps the view's text, less the qualifier of the snapshot's
// database on its tables and columns and with database's on its function calls,
// along with its parsed body. Sync qualifies every reference with that
// database, which the view resolves to without it.
// DefineView installs a view whose body does not analyze, so a view may name
// one that installs after it; reanalyzeViews retries it then. nameViewColumns
// names the columns once the body resolves.
func metadataViewStmt(v *metadata.ViewMetadata, snapshotDatabase, database string) (*nodes.CreateViewStmt, error) {
	body := stripBodyQualifier(v.GetDefinition(), snapshotDatabase, database)
	body, checkOption := cutViewCheckOption(body)
	sel, err := parseFirstStatement[*nodes.SelectStmt](body)
	if err != nil {
		return nil, err
	}
	return &nodes.CreateViewStmt{OrReplace: true, Name: &nodes.TableRef{Name: v.GetName()}, Select: sel, SelectText: body, CheckOption: checkOption}, nil
}

// cutViewCheckOption takes the WITH CHECK OPTION MariaDB ends a view's
// information_schema.VIEWS.VIEW_DEFINITION with off the body: "select ... WITH
// CASCADED CHECK OPTION". Measured on mariadb:11.8.8.
func cutViewCheckOption(body string) (string, string) {
	for _, option := range []string{"CASCADED", "LOCAL"} {
		if rest, ok := strings.CutSuffix(body, " WITH "+option+" CHECK OPTION"); ok {
			return rest, option
		}
	}
	return body, ""
}

// metadataViewColumnNames lists the column names the snapshot records for a
// view.
func metadataViewColumnNames(v *metadata.ViewMetadata) []string {
	var names []string
	for _, col := range v.GetColumns() {
		if col.GetName() != "" {
			names = append(names, col.GetName())
		}
	}
	return names
}

// nameViewColumns names a view's columns from the snapshot: the ones its
// resolved body left unnamed, and, where the body named as many by other names,
// the snapshot's as a column list. A stored definition drops the column list of
// CREATE VIEW v (a, b), and only a view created with one has names its body
// does not give.
func nameViewColumns(v *View, names []string) {
	if v == nil || len(names) == 0 || len(v.Columns) > len(names) {
		return
	}
	// Only a body that named every column of its own can disagree with the
	// snapshot over a name. One that left a column unnamed did not resolve: a
	// star over a missing relation expands to fewer columns, or to a single
	// unnamed one, and names them all from the snapshot instead.
	bodyNamedAll := len(v.Columns) == len(names) && !slices.Contains(v.Columns, "")
	switch {
	case !bodyNamedAll:
		v.Columns = names
	case !slices.EqualFunc(v.Columns, names, strings.EqualFold):
		v.Columns, v.ExplicitColumns = names, true
	default:
		return
	}
	nameViewColumnMetadata(v, names)
}

// nameViewColumnMetadata carries the view's column names into the metadata
// DefineView inferred alongside them. A body with a column for each of the
// snapshot's kept them in that order, so only the names change. One with fewer
// expanded a star the catalog could not resolve, which shifts every column
// after it: what the body inferred belongs to no column the snapshot names,
// and the names stand on their own.
func nameViewColumnMetadata(v *View, names []string) {
	if len(v.ColumnMetadata) != len(names) {
		v.ColumnMetadata = make([]ViewColumn, len(names))
	}
	for i, name := range names {
		v.ColumnMetadata[i].Name = name
	}
}

// metadataTriggerStmt builds a trigger from its parts. Sync stores the body
// without the CREATE TRIGGER header, and a body that does not parse is kept as
// text only.
func metadataTriggerStmt(tm *metadata.TriggerMetadata, table string) *nodes.CreateTriggerStmt {
	stmt := &nodes.CreateTriggerStmt{
		Name:     tm.GetName(),
		Timing:   tm.GetTiming(),
		Event:    tm.GetEvent(),
		Table:    &nodes.TableRef{Name: table},
		BodyText: tm.GetBody(),
	}
	if strings.TrimSpace(tm.GetBody()) != "" {
		probe := fmt.Sprintf("CREATE TRIGGER `probe` %s %s ON %s FOR EACH ROW %s", tm.GetTiming(), tm.GetEvent(), quoteIdent(table), tm.GetBody())
		if parsed, err := parseFirstStatement[*nodes.CreateTriggerStmt](probe); err == nil {
			stmt.Body = parsed.Body
		}
	}
	return stmt
}

// --- Stand-ins ---

var errNoStandIn = errors.New("no stand-in for this kind of object")

// installStandIn installs o's stand-in, or returns errNoStandIn for a kind
// without one.
func (c *Catalog) installStandIn(o *metadataObject) error {
	switch o.kind {
	case metadataTable:
		return c.DefineTable(standInTableStmt(o.table))
	case metadataView:
		return c.DefineView(standInViewStmt(o.view))
	case metadataTrigger:
		timing, event := o.trigger.GetTiming(), o.trigger.GetEvent()
		if timing == "" {
			timing = "BEFORE"
		}
		if event == "" {
			event = "INSERT"
		}
		return c.DefineTrigger(&nodes.CreateTriggerStmt{
			Name:     o.name,
			Timing:   timing,
			Event:    event,
			Table:    &nodes.TableRef{Name: o.parent},
			BodyText: "BEGIN END",
		})
	case metadataEvent:
		return c.DefineEvent(&nodes.CreateEventStmt{Name: o.name})
	default:
		return errNoStandIn
	}
}

// standInPlaceholder names the one column of a stand-in with no column names;
// MariaDB requires a table or view to have a column.
const standInPlaceholder = "__stand_in"

func standInTableStmt(t *metadata.TableMetadata) *nodes.CreateTableStmt {
	stmt := &nodes.CreateTableStmt{Table: &nodes.TableRef{Name: t.GetName()}}
	var names []string
	for _, col := range t.GetColumns() {
		names = appendUniqueName(names, col.GetName())
	}
	if len(names) == 0 {
		names = []string{standInPlaceholder}
	}
	for _, name := range names {
		stmt.Columns = append(stmt.Columns, &nodes.ColumnDef{Name: name, TypeName: &nodes.DataType{Name: "TEXT"}})
	}
	return stmt
}

// standInViewStmt selects NULL under the view's own column names, else under
// the names of the columns it reads.
func standInViewStmt(v *metadata.ViewMetadata) *nodes.CreateViewStmt {
	var names []string
	for _, col := range v.GetColumns() {
		names = appendUniqueName(names, col.GetName())
	}
	if len(names) == 0 {
		for _, dep := range v.GetDependencyColumns() {
			names = appendUniqueName(names, dep.GetColumn())
		}
	}
	if len(names) == 0 {
		names = []string{standInPlaceholder}
	}
	targets := make([]nodes.ExprNode, 0, len(names))
	for _, name := range names {
		targets = append(targets, &nodes.ResTarget{Name: name, Val: &nodes.NullLit{}})
	}
	return &nodes.CreateViewStmt{
		OrReplace: true,
		Name:      &nodes.TableRef{Name: v.GetName()},
		Select:    &nodes.SelectStmt{TargetList: targets},
	}
}

func appendUniqueName(names []string, name string) []string {
	if name == "" || slices.Contains(names, name) {
		return names
	}
	return append(names, name)
}

// --- Parsing snapshot text ---

// parseTypeName parses a column type as MariaDB reports it, such as
// "varchar(255)", "int unsigned", or "enum('a','b')".
func parseTypeName(typ string) (*nodes.DataType, error) {
	if strings.TrimSpace(typ) == "" {
		return nil, errors.New("empty type")
	}
	stmt, err := parseFirstStatement[*nodes.CreateTableStmt]("CREATE TABLE `probe` (`c` " + typ + ")")
	if err != nil {
		return nil, fmt.Errorf("type %q: %w", typ, err)
	}
	if len(stmt.Columns) == 0 || stmt.Columns[0].TypeName == nil {
		return nil, fmt.Errorf("type %q: no type parsed", typ)
	}
	return stmt.Columns[0].TypeName, nil
}

// parseExpr parses a default, ON UPDATE, generation, or CHECK expression.
func parseExpr(expr string) (nodes.ExprNode, error) {
	if strings.TrimSpace(expr) == "" {
		return nil, errors.New("empty expression")
	}
	sel, err := parseFirstStatement[*nodes.SelectStmt]("SELECT (" + expr + ") AS `probe`")
	if err != nil {
		return nil, fmt.Errorf("expression %q: %w", expr, err)
	}
	if len(sel.TargetList) == 0 {
		return nil, fmt.Errorf("expression %q: no target parsed", expr)
	}
	target, ok := sel.TargetList[0].(*nodes.ResTarget)
	if !ok {
		return nil, fmt.Errorf("expression %q: parsed as %T", expr, sel.TargetList[0])
	}
	// The probe's parentheses are not part of the expression.
	if p, ok := target.Val.(*nodes.ParenExpr); ok {
		return p.Expr, nil
	}
	return target.Val, nil
}

// parseRoutine parses a function's or procedure's SHOW CREATE text.
func parseRoutine(definition string, procedure bool) (*nodes.CreateFunctionStmt, error) {
	stmt, err := parseFirstStatement[*nodes.CreateFunctionStmt](definition)
	if err != nil {
		return nil, err
	}
	stmt.IsProcedure = procedure
	return stmt, nil
}

// parseFirstStatement parses sql and returns its first statement of type T.
func parseFirstStatement[T nodes.Node](sql string) (T, error) {
	var zero T
	if strings.TrimSpace(sql) == "" {
		return zero, errors.New("empty definition")
	}
	list, err := mariadbparser.Parse(sql)
	if err != nil {
		return zero, err
	}
	if list != nil {
		for _, item := range list.Items {
			if stmt, ok := item.(T); ok {
				return stmt, nil
			}
		}
	}
	return zero, fmt.Errorf("no %T statement", zero)
}

func unquoteIdent(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return strings.ReplaceAll(s[1:len(s)-1], "``", "`")
	}
	return s
}
