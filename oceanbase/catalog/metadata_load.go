package catalog

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/oceanbase/ast"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// LoadMetadataReport lists the snapshot objects LoadMetadata could not install
// as the snapshot defines them. Keys are "table:<name>", "view:<name>",
// "function:<name>", "procedure:<name>", "trigger:<table>.<name>", and
// "event:<name>".
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

// LoadMetadata installs a Bytebase schema snapshot of a MySQL database into the
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
	for _, o := range collectMetadataObjects(meta) {
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
		for _, o := range installed {
			applyObjectSkipDump(db, o)
		}
	}
	return report, nil
}

func applyObjectSkipDump(db *Database, o *metadataObject) {
	if db == nil || o == nil {
		return
	}
	key := toLower(o.name)
	switch o.kind {
	case metadataTable:
		if tbl := db.Tables[key]; tbl != nil && o.table != nil {
			tbl.SkipDump = o.table.GetSkipDump()
		}
	case metadataView:
		if v := db.Views[key]; v != nil {
			if o.materialized && o.mview != nil {
				v.SkipDump = o.mview.GetSkipDump()
			} else if o.view != nil {
				v.SkipDump = o.view.GetSkipDump()
			}
		}
	case metadataFunction:
		if fn := db.Functions[key]; fn != nil && o.function != nil {
			fn.SkipDump = o.function.GetSkipDump()
		}
	case metadataProcedure:
		if p := db.Procedures[key]; p != nil && o.procedure != nil {
			p.SkipDump = o.procedure.GetSkipDump()
		}
	case metadataTrigger:
		if tr := db.Triggers[key]; tr != nil && o.trigger != nil {
			tr.SkipDump = o.trigger.GetSkipDump()
		}
	}
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
// explicit_defaults_for_timestamp=OFF rewrites a TIMESTAMP column. It returns a function restoring the caller's.
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

func (c *Catalog) installedTable(name string) *Table {
	db := c.GetDatabase(c.CurrentDatabase())
	if db == nil {
		return nil
	}
	return db.GetTable(name)
}

func applyDroppedColumnDefaults(tbl *Table, meta *metadata.TableMetadata) {
	if tbl == nil {
		return
	}
	dropped := make(map[string]bool)
	for _, col := range meta.GetColumns() {
		if col.GetDefaultDropped() {
			dropped[toLower(col.GetName())] = true
		}
	}
	if len(dropped) == 0 {
		return
	}
	for _, col := range tbl.Columns {
		if dropped[toLower(col.Name)] && col.Default == nil {
			col.DefaultDropped = true
		}
	}
}

// --- Objects ---

// metadataObjectKind lists the kinds in the order they install. Only a
// trigger needs another object first, its table: foreign key checks are off,
// a view may name an object not yet installed, and routine bodies are not
// analyzed.
type metadataObjectKind int

const (
	metadataTable metadataObjectKind = iota
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
	// parent is a trigger's table; order preserves the snapshot's action order.
	parent string
	order  int
	// snapshotDatabase is the database the snapshot's references qualify
	// with: a view's body, and a table's foreign keys.
	snapshotDatabase string

	// materialized marks a view the snapshot listed as a materialized view.
	materialized bool

	table     *metadata.TableMetadata
	view      *metadata.ViewMetadata
	mview     *metadata.MaterializedViewMetadata
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
func collectMetadataObjects(meta *metadata.DatabaseSchemaMetadata) []*metadataObject {
	var out []*metadataObject
	for _, s := range meta.GetSchemas() {
		for _, t := range s.GetTables() {
			if t.GetName() == "" {
				continue
			}
			out = append(out, &metadataObject{kind: metadataTable, name: t.GetName(), table: t, snapshotDatabase: meta.GetName()})
			for i, trg := range t.GetTriggers() {
				if trg.GetName() != "" {
					out = append(out, &metadataObject{kind: metadataTrigger, name: trg.GetName(), parent: t.GetName(), order: i, trigger: trg})
				}
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataView, name: v.GetName(), view: v, snapshotDatabase: meta.GetName()})
			}
		}
		// The catalog models no materialized view of its own, so a snapshot's materialized
		// view installs as the view its query is, marked materialized. Without it, a view
		// that reads the materialized view would name a relation the catalog does not hold;
		// the container table holds its rows, and the sync no longer lists that as a table.
		for _, v := range s.GetMaterializedViews() {
			if v.GetName() == "" || v.GetSkipDump() {
				continue
			}
			out = append(out, &metadataObject{
				kind:             metadataView,
				name:             v.GetName(),
				materialized:     true,
				view:             &metadata.ViewMetadata{Name: v.GetName(), Definition: v.GetDefinition()},
				mview:            v,
				snapshotDatabase: meta.GetName(),
			})
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
		if a.kind == metadataTrigger && b.kind == metadataTrigger {
			return cmp.Or(cmp.Compare(a.parent, b.parent), cmp.Compare(a.order, b.order), cmp.Compare(a.name, b.name))
		}
		return cmp.Or(cmp.Compare(a.kind, b.kind), cmp.Compare(a.name, b.name), cmp.Compare(a.parent, b.parent))
	})
	return out
}

// --- Installing ---

func (c *Catalog) installMetadataObject(o *metadataObject) error {
	switch o.kind {
	case metadataTable:
		stmt, omitted, err := metadataTableStmt(o.table, o.snapshotDatabase)
		if err != nil {
			return err
		}
		o.incomplete = omitted
		finishInstall := func() {
			applyDroppedColumnDefaults(c.installedTable(o.name), o.table)
			if db := c.GetDatabase(c.CurrentDatabase()); db != nil {
				key := toLower(o.table.GetName())
				db.MaterializedViewLogsKnown[key] = true
				if o.table.GetHasMaterializedViewLog() {
					db.MaterializedViewLogs[key] = true
				} else {
					delete(db.MaterializedViewLogs, key)
				}
			}
		}
		err = c.DefineTable(stmt)
		if err != nil && stmt.Partitions != nil {
			// The partitioning parsed but the catalog refuses it. The rest of the
			// table still installs.
			stmt.Partitions = nil
			if retry := c.DefineTable(stmt); retry == nil {
				if o.incomplete == nil {
					o.incomplete = fmt.Errorf("partitioning: %w", err)
				}
				finishInstall()
				return nil
			}
		}
		if err == nil {
			finishInstall()
		}
		return err
	case metadataView:
		stmt, err := metadataViewStmt(o.view, o.snapshotDatabase, c.CurrentDatabase())
		if err != nil {
			return err
		}
		if o.materialized {
			if err := c.DefineMaterializedView(stmt); err != nil {
				return err
			}
			if inst := c.installedView(o.name); inst != nil && o.mview != nil {
				inst.Comment = syncedComment(o.mview.GetComment())
				inst.Indexes = metadataMVIndexes(o.mview.GetIndexes(), &o.incomplete)
			}
		} else {
			if err := c.DefineView(stmt); err != nil {
				return err
			}
			if inst := c.installedView(o.name); inst != nil {
				if o.view.GetAlgorithm() != "" {
					inst.Algorithm = o.view.GetAlgorithm()
				}
				if o.view.GetSqlSecurity() != "" {
					inst.SqlSecurity = o.view.GetSqlSecurity()
				}
			}
		}
		nameViewColumns(c.installedView(o.name), metadataViewColumnNames(o.view))
		preserveMetadataViewColumns(c.installedView(o.name), o.view.GetColumns())
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

// The snapshot carries types, defaults, and definitions as the text MySQL
// reports for them. These builders turn that text into the statement ASTs the
// Define* functions take.

// autoIncrementDefault is what MySQL sync stores as the default of an
// AUTO_INCREMENT column.
const autoIncrementDefault = "AUTO_INCREMENT"

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
// left out instead, as is an index with a key expression that does not parse,
// which the catalog would reject as a column; omitted is the first such error.
func metadataTableStmt(t *metadata.TableMetadata, snapshotDatabase string) (stmt *nodes.CreateTableStmt, omitted error, err error) {
	omit := func(e error) {
		if omitted == nil {
			omitted = e
		}
	}
	stmt = &nodes.CreateTableStmt{
		Table:       &nodes.TableRef{Name: t.GetName()},
		ColumnGroup: t.GetColumnGroups(),
	}
	for _, col := range t.GetColumns() {
		if col.GetName() == "" {
			continue
		}
		def, colOmitted, err := metadataColumnDef(t, col)
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
		con, err := metadataIndexConstraint(idx, t.GetEngine())
		if err != nil {
			omit(fmt.Errorf("index %q: %w", idx.GetName(), err))
			continue
		}
		stmt.Constraints = append(stmt.Constraints, con)
	}
	for _, fk := range t.GetForeignKeys() {
		stmt.Constraints = append(stmt.Constraints, metadataForeignKey(fk, snapshotDatabase))
	}
	for _, check := range t.GetCheckConstraints() {
		expr, err := parseExpr(check.GetExpression())
		if err != nil {
			omit(fmt.Errorf("check %q: %w", check.GetName(), err))
			continue
		}
		stmt.Constraints = append(stmt.Constraints, &nodes.Constraint{
			Type:        nodes.ConstrCheck,
			Name:        check.GetName(),
			Expr:        expr,
			NotEnforced: check.GetNotEnforced(),
		})
	}
	if len(t.GetPartitions()) > 0 {
		partitions, err := metadataPartitionClause(t)
		if err != nil {
			omit(fmt.Errorf("partitioning: %w", err))
		}
		stmt.Partitions = partitions
	}
	opts := t.GetCreateOptions()
	for _, opt := range []struct{ name, value string }{
		{"ENGINE", t.GetEngine()},
		{"CHARSET", t.GetCharset()},
		{"COLLATE", t.GetCollation()},
		{"COMMENT", syncedComment(t.GetComment())},
		{"KEY_BLOCK_SIZE", createOption(opts, "key_block_size")},
		{"ROW_FORMAT", createOption(opts, "row_format")},
		{"COMPRESSION", unquoteTableOption(createOption(opts, "compression"))},
		{"ENCRYPTION", unquoteTableOption(createOption(opts, "encryption"))},
		{"STATS_PERSISTENT", createOption(opts, "stats_persistent")},
		{"STATS_AUTO_RECALC", createOption(opts, "stats_auto_recalc")},
		{"STATS_SAMPLE_PAGES", createOption(opts, "stats_sample_pages")},
		{"MIN_ROWS", createOption(opts, "min_rows")},
		{"MAX_ROWS", createOption(opts, "max_rows")},
		{"AVG_ROW_LENGTH", createOption(opts, "avg_row_length")},
		{"PACK_KEYS", createOption(opts, "pack_keys")},
		{"CHECKSUM", createOption(opts, "checksum")},
		{"DELAY_KEY_WRITE", createOption(opts, "delay_key_write")},
		{"TABLESPACE", unquoteTableOption(t.GetTablespace())},
	} {
		if opt.value != "" {
			stmt.Options = append(stmt.Options, &nodes.TableOption{Name: opt.name, Value: opt.value})
		}
	}
	return stmt, omitted, nil
}

// metadataPartitionClause parses the partitioning the snapshot records, as the
// definition renderer prints it, and attaches attributes that OceanBase DDL
// rejects inline (partition/subpartition comments) or forbids alongside explicit
// partition lists (SUBPARTITION TEMPLATE).
func metadataPartitionClause(t *metadata.TableMetadata) (*nodes.PartitionClause, error) {
	var buf strings.Builder
	if err := printPartitionClause(&buf, t.GetPartitions(), t.GetEngine(), false, false); err != nil {
		return nil, err
	}
	stmt, err := parseFirstStatement[*nodes.CreateTableStmt]("CREATE TABLE `probe` (`c` int)" + buf.String())
	if err != nil {
		return nil, err
	}
	if stmt.Partitions == nil {
		return nil, fmt.Errorf("no partitioning parsed from %q", buf.String())
	}
	for i, pd := range stmt.Partitions.Partitions {
		if i >= len(t.GetPartitions()) {
			break
		}
		partMeta := t.GetPartitions()[i]
		if partMeta.GetComment() != "" {
			pd.Options = append(pd.Options, &nodes.TableOption{Name: "COMMENT", Value: partMeta.GetComment()})
		}
		for j, spd := range pd.SubPartitions {
			if j >= len(partMeta.GetSubpartitions()) {
				break
			}
			if cmt := partMeta.GetSubpartitions()[j].GetComment(); cmt != "" {
				spd.Options = append(spd.Options, &nodes.TableOption{Name: "COMMENT", Value: cmt})
			}
		}
	}
	if len(t.GetPartitions()) > 0 && len(t.GetPartitions()[0].GetSubpartitionTemplate()) > 0 {
		tpl := t.GetPartitions()[0].GetSubpartitionTemplate()
		if len(stmt.Partitions.SubPartTemplate) == 0 {
			parsedTpl, err := metadataSubpartitionTemplate(t.GetPartitions()[0])
			if err != nil {
				return nil, err
			}
			stmt.Partitions.SubPartTemplate = parsedTpl
		}
		for i, spd := range stmt.Partitions.SubPartTemplate {
			if i < len(tpl) && tpl[i].GetComment() != "" {
				spd.Options = append(spd.Options, &nodes.TableOption{Name: "COMMENT", Value: tpl[i].GetComment()})
			}
		}
	}
	return stmt.Partitions, nil
}

func metadataSubpartitionTemplate(partMeta *metadata.TablePartitionMetadata) ([]*nodes.SubPartitionDef, error) {
	tpl := partMeta.GetSubpartitionTemplate()
	sub0 := tpl[0]
	if sub0.GetType() == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED && len(partMeta.GetSubpartitions()) > 0 {
		sub0 = partMeta.GetSubpartitions()[0]
	}
	var buf strings.Builder
	buf.WriteString("CREATE TABLE `probe` (`c` int) PARTITION BY HASH (`c`)")
	if err := writeSubpartitionBy(&buf, sub0); err != nil {
		return nil, err
	}
	subType := subPartitionTypeName(sub0.GetType())
	colCount := subPartitionColumnCount(sub0.GetExpression())
	buf.WriteString(" SUBPARTITION TEMPLATE (")
	for i, st := range tpl {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "SUBPARTITION %s", quoteIdent(st.GetName()))
		if vals := subpartitionValuesSQL(subType, colCount, st.GetValue()); vals != "" {
			buf.WriteString(vals)
		}
	}
	buf.WriteString(")")
	stmt, err := parseFirstStatement[*nodes.CreateTableStmt](buf.String())
	if err != nil {
		return nil, err
	}
	if stmt.Partitions == nil {
		return nil, fmt.Errorf("no subpartition template parsed from %q", buf.String())
	}
	for i, spd := range stmt.Partitions.SubPartTemplate {
		if i < len(tpl) && tpl[i].GetComment() != "" {
			spd.Options = append(spd.Options, &nodes.TableOption{Name: "COMMENT", Value: tpl[i].GetComment()})
		}
	}
	return stmt.Partitions.SubPartTemplate, nil
}

func metadataMVIndexes(idxs []*metadata.IndexMetadata, omitted *error) []*Index {
	if len(idxs) == 0 {
		return nil
	}
	out := make([]*Index, 0, len(idxs))
	for _, meta := range idxs {
		if len(meta.GetExpressions()) == 0 {
			continue
		}
		idxType := strings.ToUpper(strings.TrimSpace(meta.GetType()))
		if idxType == "BTREE" {
			idxType = ""
		}
		spatial := idxType == "SPATIAL"
		idx := &Index{
			Name:             meta.GetName(),
			Unique:           meta.GetUnique() || meta.GetPrimary(),
			Primary:          meta.GetPrimary(),
			Visible:          meta.GetVisible() || meta.GetPrimary(),
			Comment:          meta.GetComment(),
			KeyBlockSize:     int(meta.GetKeyBlockSize()),
			IndexType:        idxType,
			Fulltext:         idxType == "FULLTEXT",
			Spatial:          spatial,
			Vector:           meta.GetIsVector(),
			Scope:            meta.GetScope(),
			Storing:          slices.Clone(meta.GetStoringColumns()),
			VectorParams:     slices.Clone(meta.GetVectorParams()),
			Parser:           meta.GetParser(),
			ParserProperties: maps.Clone(meta.GetParserProperties()),
		}
		if idx.Scope == "" && meta.GetIsGlobal() {
			idx.Scope = "GLOBAL"
		}
		var parseErr error
		for i, expr := range meta.GetExpressions() {
			desc := i < len(meta.GetDescending()) && meta.GetDescending()[i]
			length := 0
			if !spatial && i < len(meta.GetKeyLength()) && meta.GetKeyLength()[i] > 0 {
				length = int(meta.GetKeyLength()[i])
			}
			col := &IndexColumn{
				Length:     length,
				Descending: desc,
			}
			if isExpressionKey(expr) {
				if _, err := parseExpr(unescapeFunctionalIndexQuotes(expr)); err != nil && parseErr == nil {
					parseErr = err
				}
				col.Expr = expr
			} else {
				col.Name = unquoteIdent(expr)
			}
			idx.Columns = append(idx.Columns, col)
		}
		if parseErr != nil {
			if *omitted == nil {
				*omitted = fmt.Errorf("index %q: %w", meta.GetName(), parseErr)
			}
			continue
		}
		out = append(out, idx)
	}
	return out
}

// createOption returns the value of one CREATE_OPTIONS entry; sync keeps them
// as space-separated name=value pairs.
func createOption(options, name string) string {
	for _, opt := range strings.Fields(options) {
		if k, v, ok := strings.Cut(opt, "="); ok && strings.EqualFold(k, name) {
			return v
		}
	}
	return ""
}

// metadataColumnDef fails on a type that does not parse. A default, ON UPDATE,
// or generation expression that does not parse is left out, and omitted is the
// first such error.
func metadataColumnDef(t *metadata.TableMetadata, col *metadata.ColumnMetadata) (def *nodes.ColumnDef, omitted error, err error) {
	typeName, err := parseTypeName(col.GetType())
	if err != nil {
		return nil, nil, err
	}
	// The type text may carry its own character set and collation.
	// When inherited flags are present, false marks an explicit column clause even
	// if it matches the table default; when absent (legacy live sync), matching the
	// table default stays inherited.
	charsetExplicit := col.CharsetInherited != nil && !*col.CharsetInherited ||
		col.CharsetInherited == nil && t != nil && !strings.EqualFold(col.GetCharacterSet(), t.GetCharset())
	collationExplicit := col.CollationInherited != nil && !*col.CollationInherited ||
		col.CollationInherited == nil && t != nil && !strings.EqualFold(col.GetCollation(), t.GetCollation())
	if typeName.Charset == "" && col.GetCharacterSet() != "" && charsetExplicit {
		typeName.Charset = col.GetCharacterSet()
	}
	if typeName.Collate == "" && col.GetCollation() != "" && collationExplicit {
		typeName.Collate = col.GetCollation()
	}
	// SRID 0 is a valid SRID, so presence decides.
	if col.Srid != nil {
		typeName.SRID = int(col.GetSrid())
		typeName.HasSRID = true
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
	def.DefaultValue, omitted = columnDefault(col)
	if col.GetOnUpdate() != "" {
		if expr, err := parseExpr(col.GetOnUpdate()); err == nil {
			def.OnUpdate = expr
		} else if omitted == nil {
			omitted = fmt.Errorf("ON UPDATE: %w", err)
		}
	}
	if gen := col.GetGeneration(); gen.GetExpression() != "" {
		if expr, err := parseExpr(gen.GetExpression()); err == nil {
			def.Generated = &nodes.GeneratedColumn{Expr: expr, Stored: gen.GetType() == metadata.GenerationMetadata_TYPE_STORED}
		} else if omitted == nil {
			omitted = fmt.Errorf("generation: %w", err)
		}
	}
	return def, omitted, nil
}

// columnDefault parses the default the snapshot records for a column, or
// returns nil where the catalog takes none: AUTO_INCREMENT is a flag rather
// than a default, a generated column has no default, and sync reports DEFAULT
// NULL even for the types that show no default.
func columnDefault(col *metadata.ColumnMetadata) (nodes.ExprNode, error) {
	value := col.GetDefault()
	switch {
	case value == "" || col.GetGeneration() != nil:
		return nil, nil
	case strings.EqualFold(value, autoIncrementDefault):
		return nil, nil
	case strings.EqualFold(value, "NULL") && !typeTakesDefault(col.GetType()):
		return nil, nil
	}
	expr, err := parseExpr(value)
	if err != nil {
		return nil, fmt.Errorf("default: %w", err)
	}
	return expr, nil
}

// typeTakesDefault reports whether a column type shows a DEFAULT NULL clause,
// which the BLOB, TEXT, JSON, and spatial types never do.
func typeTakesDefault(typ string) bool {
	head := strings.ToLower(strings.TrimSpace(typ))
	if i := strings.IndexAny(head, " ("); i >= 0 {
		head = head[:i]
	}
	return validateDefaultOfType("", head, &nodes.IntLit{}) == nil
}

// metadataIndexConstraint builds an index from its key parts. The index type
// holds either the kind, FULLTEXT or SPATIAL, or the access method. Sync reports
// the access method a key without USING gets too, so the engine's default one
// stays unwritten. The error is a key expression that did not parse; the index
// still comes back, with that key part as a column.
func metadataIndexConstraint(idx *metadata.IndexMetadata, engine string) (*nodes.Constraint, error) {
	indexType := strings.ToUpper(strings.TrimSpace(idx.GetType()))
	if indexType == defaultIndexType(engine) {
		indexType = ""
	}
	c := &nodes.Constraint{Name: idx.GetName()}
	// A key OceanBase named PRIMARY, after its first column, cannot be created by that
	// name (1280): unnamed, the catalog names it PRIMARY again (allocIndexName).
	if !idx.GetPrimary() && strings.EqualFold(c.Name, "PRIMARY") {
		c.Name = ""
	}
	switch {
	case idx.GetPrimary():
		c.Type, c.IndexType = nodes.ConstrPrimaryKey, indexType
	case idx.GetUnique():
		c.Type, c.IndexType = nodes.ConstrUnique, indexType
	case idx.GetIsVector():
		c.Type, c.IndexType = nodes.ConstrVector, indexType
	case indexType == "FULLTEXT":
		c.Type = nodes.ConstrFulltextIndex
	case indexType == "SPATIAL":
		c.Type = nodes.ConstrSpatialIndex
	default:
		c.Type, c.IndexType = nodes.ConstrIndex, indexType
	}
	var err error
	c.Columns, c.IndexColumns, err = indexKeyParts(idx, c.Type == nodes.ConstrSpatialIndex)
	// A primary key is always visible.
	if !idx.GetVisible() && !idx.GetPrimary() {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "INVISIBLE"})
	}
	if idx.GetComment() != "" {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "COMMENT", Value: &nodes.StringLit{Value: idx.GetComment()}})
	}
	if idx.GetKeyBlockSize() > 0 {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "KEY_BLOCK_SIZE", Value: &nodes.IntLit{Value: idx.GetKeyBlockSize()}})
	}
	if len(idx.GetStoringColumns()) > 0 {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "STORING", Columns: slices.Clone(idx.GetStoringColumns())})
	}
	if len(idx.GetVectorParams()) > 0 {
		props := make([]*nodes.ParserProperty, 0, len(idx.GetVectorParams()))
		for _, entry := range idx.GetVectorParams() {
			if k, v, ok := strings.Cut(entry, "="); ok {
				props = append(props, &nodes.ParserProperty{Name: k, Value: parsePropertyNode(v)})
			}
		}
		if len(props) > 0 {
			c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "WITH", Properties: props})
		}
	}
	if idx.GetParser() != "" {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "PARSER", Value: &nodes.StringLit{Value: idx.GetParser()}})
	}
	if len(idx.GetParserProperties()) > 0 {
		keys := make([]string, 0, len(idx.GetParserProperties()))
		for k := range idx.GetParserProperties() {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		props := make([]*nodes.ParserProperty, 0, len(keys))
		for _, k := range keys {
			props = append(props, &nodes.ParserProperty{Name: k, Value: parsePropertyNode(idx.GetParserProperties()[k])})
		}
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "PARSER_PROPERTIES", Properties: props})
	}
	scope := idx.GetScope()
	if scope == "" && idx.GetIsGlobal() {
		scope = "GLOBAL"
	}
	if scope != "" {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: strings.ToUpper(scope)})
	}
	return c, err
}

func parsePropertyNode(val string) nodes.Node {
	if n, err := strconv.ParseInt(val, 10, 64); err == nil {
		return &nodes.IntLit{Value: n}
	}
	return &nodes.StringLit{Value: val}
}

// indexKeyParts splits an index's key parts: plain column names, or, where a
// part is an expression or has a length or descends, an IndexColumn for each.
// A SPATIAL key takes no length, though sync reports one (32) for it.
func indexKeyParts(idx *metadata.IndexMetadata, spatial bool) ([]string, []*nodes.IndexColumn, error) {
	length := func(i int) int {
		if !spatial && i < len(idx.GetKeyLength()) {
			return int(idx.GetKeyLength()[i])
		}
		return 0
	}
	descending := func(i int) bool { return i < len(idx.GetDescending()) && idx.GetDescending()[i] }
	plain := true
	for i, expr := range idx.GetExpressions() {
		if isExpressionKey(expr) || length(i) > 0 || descending(i) {
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
		if isExpressionKey(expr) {
			// Sync reports the expression as information_schema stores it,
			// with its quotes escaped. Charset introducers stay: they can
			// change what the expression computes. On a parse failure the key
			// stays a column, and the caller leaves the index out.
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
// database, along with its parsed body. Sync qualifies every reference with
// that database, which the view resolves to without it.
// DefineView installs a view whose body does not analyze, so a view may name
// one that installs after it; reanalyzeViews retries it then. nameViewColumns
// names the columns once the body resolves.
func metadataViewStmt(v *metadata.ViewMetadata, database, target string) (*nodes.CreateViewStmt, error) {
	body := rebaseViewBodyDatabaseQualifier(v.GetDefinition(), database, target)
	sel, err := parseFirstStatement[*nodes.SelectStmt](body)
	if err != nil {
		return nil, err
	}
	return &nodes.CreateViewStmt{
		OrReplace:   true,
		Name:        &nodes.TableRef{Name: v.GetName()},
		Select:      sel,
		SelectText:  body,
		Algorithm:   v.GetAlgorithm(),
		SqlSecurity: v.GetSqlSecurity(),
		CheckOption: v.GetCheckOption(),
	}, nil
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

func preserveMetadataViewColumns(v *View, cols []*metadata.ColumnMetadata) {
	if v == nil || len(cols) == 0 {
		return
	}
	if len(v.ColumnMetadata) < len(cols) {
		grown := make([]ViewColumn, len(cols))
		copy(grown, v.ColumnMetadata)
		v.ColumnMetadata = grown
	}
	for i, col := range cols {
		vc := &v.ColumnMetadata[i]
		if vc.Name == "" {
			vc.Name = col.GetName()
		}
		if (v.AnalyzedQuery == nil || vc.Type == nil) && vc.ColumnType == "" {
			vc.ColumnType = col.GetType()
		}
		if v.AnalyzedQuery == nil {
			vc.Nullable = col.GetNullable()
		}
		if vc.Charset == "" {
			vc.Charset = col.GetCharacterSet()
		}
		if vc.Collation == "" {
			vc.Collation = col.GetCollation()
		}
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
		if o.materialized {
			return c.DefineMaterializedView(standInViewStmt(o.view))
		}
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
// MySQL requires a table or view to have a column.
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

// parseTypeName parses a column type as MySQL reports it, such as
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

// syncedComment is a table's or column's comment as sync records it: escaped
// by QUOTE(), which writes a backslash before a backslash and a quote, NUL as
// \0 and Ctrl-Z as \Z, without the quotes around it.
func syncedComment(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 == len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case '0':
			b.WriteByte(0)
		case 'Z':
			b.WriteByte(26)
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
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
	list, err := obparser.Parse(sql)
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
