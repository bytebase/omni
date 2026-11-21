package catalog

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	mysqlast "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
	nodes "github.com/bytebase/omni/tidb/ast"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// LoadMetadataReport lists the snapshot objects LoadMetadata could not install
// as the snapshot defines them. Keys are "table:<name>", "view:<name>",
// "sequence:<name>", "function:<name>", "procedure:<name>",
// "trigger:<table>.<name>", and "event:<name>".
type LoadMetadataReport struct {
	// Degraded maps each object replaced by a stand-in to the error its
	// definition failed with.
	Degraded map[string]error
	// Missing maps each object the catalog has nothing for to the error that
	// left it out.
	Missing map[string]error
	// Incomplete maps each table installed without part of its definition to
	// the error of the first part left out: a default, AUTO_RANDOM, ON UPDATE,
	// generation, or CHECK expression, an index key expression, or the
	// partitioning.
	Incomplete map[string]error
}

// LoadMetadata installs a Bytebase schema snapshot of a TiDB database into the
// database meta.Name, creating it if needed, and leaves that database
// selected. Foreign key checks are off during the load, so a table may
// reference one that installs after it. A view may likewise name a view that
// installs after it.
//
// An object that fails to install does not stop the load: a table is replaced
// by a stand-in with the same column names, all TEXT, and a view by one that
// selects NULL under each of its column names. A table that installs without
// part of its definition is reported Incomplete. The error is non-nil only when
// ctx is done.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadata.DatabaseSchemaMetadata) (*LoadMetadataReport, error) {
	return c.LoadMetadataInto(ctx, meta, meta.GetName())
}

// LoadMetadataInto is LoadMetadata into the database named database rather
// than meta.Name, so that two snapshots of differently named databases can be
// loaded under one name and diffed. View bodies and foreign keys drop the
// qualifier of meta.Name, which TiDB puts on every table a view reads.
func (c *Catalog) LoadMetadataInto(ctx context.Context, meta *metadata.DatabaseSchemaMetadata, database string) (*LoadMetadataReport, error) {
	report := &LoadMetadataReport{
		Degraded:   make(map[string]error),
		Missing:    make(map[string]error),
		Incomplete: make(map[string]error),
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	if name := database; name != "" {
		// A database already in the catalog keeps its own character set and
		// collation otherwise, and a table that overrides neither would resolve
		// its columns differently than the same snapshot loaded into an empty
		// catalog.
		stmt := metadataDatabaseStmt(meta, name)
		if c.GetDatabase(name) == nil {
			_ = c.DefineDatabase(stmt)
		} else if len(stmt.Options) > 0 {
			_ = c.alterDatabase(&nodes.AlterDatabaseStmt{Name: name, Options: stmt.Options})
		}
		if db := c.GetDatabase(name); db != nil {
			db.PlacementPolicy = resolvePolicyRef(meta.GetPlacementPolicy())
			db.TiFlashReplica = int(meta.GetTiflashReplica())
			db.TiFlashLocationLabels = slices.Clone(meta.GetTiflashLocationLabels())
		}
		c.SetCurrentDatabase(name)
	}
	fkChecks := c.ForeignKeyChecks()
	c.SetForeignKeyChecks(false)
	defer c.SetForeignKeyChecks(fkChecks)
	// A snapshot's CHECK constraints exist on the server it came from, so they
	// are kept whatever tidb_enable_check_constraint says for this session.
	checks := c.CheckConstraints()
	c.SetCheckConstraints(true)
	defer c.SetCheckConstraints(checks)

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
		if o.kind == metadataSequence {
			// A sequence has no stand-in.
			report.Missing[o.key()] = err
			continue
		}
		if standInErr := c.installStandIn(o); standInErr != nil {
			if errors.Is(standInErr, errNoStandIn) {
				report.Missing[o.key()] = err
			} else {
				report.Missing[o.key()] = standInErr
			}
			continue
		}
		installed = append(installed, o)
		report.Degraded[o.key()] = err
	}
	if err := c.reanalyzeViews(ctx, views); err != nil {
		return report, err
	}
	if db := c.GetDatabase(c.CurrentDatabase()); db != nil {
		for _, o := range installed {
			applyObjectSessionContext(db, o)
			applyObjectSkipDump(db, o)
		}
	}
	return report, nil
}

func applyObjectSessionContext(db *Database, o *metadataObject) {
	if db == nil || o == nil {
		return
	}
	key := toLower(o.name)
	switch o.kind {
	case metadataFunction:
		if fn := db.Functions[key]; fn != nil && o.function != nil {
			sqlMode, cs, col := o.function.GetSqlMode(), o.function.GetCharacterSetClient(), o.function.GetCollationConnection()
			if sqlMode != "" || cs != "" || col != "" {
				fn.SQLMode = sqlMode
				fn.CharacterSetClient = cs
				fn.CollationConnection = col
				fn.HasSessionContext = true
			}
		}
	case metadataProcedure:
		if p := db.Procedures[key]; p != nil && o.procedure != nil {
			sqlMode, cs, col := o.procedure.GetSqlMode(), o.procedure.GetCharacterSetClient(), o.procedure.GetCollationConnection()
			if sqlMode != "" || cs != "" || col != "" {
				p.SQLMode = sqlMode
				p.CharacterSetClient = cs
				p.CollationConnection = col
				p.HasSessionContext = true
			}
		}
	case metadataTrigger:
		if tr := db.Triggers[key]; tr != nil && o.trigger != nil {
			sqlMode, cs, col := o.trigger.GetSqlMode(), o.trigger.GetCharacterSetClient(), o.trigger.GetCollationConnection()
			if sqlMode != "" || cs != "" || col != "" {
				tr.SQLMode = sqlMode
				tr.CharacterSetClient = cs
				tr.CollationConnection = col
				tr.HasSessionContext = true
			}
		}
	case metadataEvent:
		if ev := db.Events[key]; ev != nil && o.event != nil {
			sqlMode, cs, col, tz := o.event.GetSqlMode(), o.event.GetCharacterSetClient(), o.event.GetCollationConnection(), o.event.GetTimeZone()
			if sqlMode != "" || cs != "" || col != "" || tz != "" {
				ev.SQLMode = sqlMode
				ev.CharacterSetClient = cs
				ev.CollationConnection = col
				ev.TimeZone = tz
				ev.HasSessionContext = true
			}
		}
	}
}

func applyObjectSkipDump(db *Database, o *metadataObject) {
	if db == nil || o == nil {
		return
	}
	key := toLower(o.name)
	switch o.kind {
	case metadataSequence:
		if s := db.Sequences[key]; s != nil && o.sequence != nil {
			s.SkipDump = o.sequence.GetSkipDump()
		}
	case metadataTable:
		if tbl := db.Tables[key]; tbl != nil && o.table != nil {
			tbl.SkipDump = o.table.GetSkipDump()
		}
	case metadataView:
		if v := db.Views[key]; v != nil && o.view != nil {
			v.SkipDump = o.view.GetSkipDump()
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

// installedTable returns the table LoadMetadata just installed, or nil.
func (c *Catalog) installedTable(name string) *Table {
	db := c.GetDatabase(c.CurrentDatabase())
	if db == nil {
		return nil
	}
	return db.Tables[toLower(name)]
}

// isCachedCreateOptions reports whether information_schema.TABLES'
// CREATE_OPTIONS marks the table cached, which TiDB writes as `cached=on`
// (read on v8.5.5).
func isCachedCreateOptions(options string) bool {
	for _, opt := range strings.Split(options, " ") {
		if strings.EqualFold(strings.TrimSpace(opt), "cached=on") {
			return true
		}
	}
	return false
}

// --- Objects ---

// metadataObjectKind lists the kinds in the order they install. Nothing else
// orders them: foreign key checks are off, and a view may name an object not
// yet installed.
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
	kind   metadataObjectKind
	name   string
	parent string
	order  int
	// snapshotDatabase is the database the snapshot's references qualify
	// with: a view's body, and a table's foreign keys.
	snapshotDatabase string

	table     *metadata.TableMetadata
	view      *metadata.ViewMetadata
	sequence  *metadata.SequenceMetadata
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
		for _, q := range s.GetSequences() {
			if q.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataSequence, name: q.GetName(), sequence: q})
			}
		}
		for _, fn := range s.GetFunctions() {
			if fn.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataFunction, name: fn.GetName(), function: fn})
			}
		}
		for _, proc := range s.GetProcedures() {
			if proc.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataProcedure, name: proc.GetName(), procedure: proc})
			}
		}
		for _, ev := range s.GetEvents() {
			if ev.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataEvent, name: ev.GetName(), event: ev})
			}
		}
	}
	slices.SortFunc(out, func(a, b *metadataObject) int {
		if a.kind == metadataTrigger && b.kind == metadataTrigger {
			return cmp.Or(cmp.Compare(a.parent, b.parent), cmp.Compare(a.order, b.order), cmp.Compare(a.name, b.name))
		}
		return cmp.Or(cmp.Compare(a.kind, b.kind), cmp.Compare(a.name, b.name))
	})
	return out
}

// --- Installing ---

func (c *Catalog) installMetadataObject(o *metadataObject) error {
	switch o.kind {
	case metadataSequence:
		return c.installMetadataSequence(o.sequence)
	case metadataView:
		stmt, err := metadataViewStmt(o.view, o.snapshotDatabase, c.CurrentDatabase())
		if err != nil {
			return err
		}
		if err := c.DefineView(stmt); err != nil {
			return err
		}
		// A * or an unaliased expression leaves the body naming fewer columns
		// than the snapshot records, so the snapshot names the rest. They are
		// not a column list the view was created with, so they stay inferred.
		nameViewColumns(c.installedView(o.name), metadataViewColumnNames(o.view))
		preserveMetadataViewColumns(c.installedView(o.name), o.view.GetColumns())
		return nil
	case metadataFunction:
		stmt, err := parseTiDBRoutine(o.function.GetDefinition(), false)
		if err != nil {
			return err
		}
		return c.DefineFunction(stmt)
	case metadataProcedure:
		stmt, err := parseTiDBRoutine(o.procedure.GetDefinition(), true)
		if err != nil {
			return err
		}
		return c.DefineProcedure(stmt)
	case metadataTrigger:
		if o.trigger.GetTiming() == "" || o.trigger.GetEvent() == "" {
			return errors.New("trigger has no timing or event")
		}
		return c.DefineTrigger(&nodes.CreateTriggerStmt{
			Name:     o.trigger.GetName(),
			Timing:   o.trigger.GetTiming(),
			Event:    o.trigger.GetEvent(),
			Table:    &nodes.TableRef{Name: o.parent},
			BodyText: o.trigger.GetBody(),
		})
	case metadataEvent:
		stmt, err := parseTiDBEvent(o.event.GetDefinition())
		if err != nil {
			return err
		}
		return c.DefineEvent(stmt)
	}
	stmt, omitted, err := metadataTableStmt(o.table, o.snapshotDatabase)
	if err != nil {
		return err
	}
	o.incomplete = omitted
	err = c.DefineTable(stmt)
	if err == nil {
		applyTableMetadataExtras(c.installedTable(o.name), o.table)
	}
	if err != nil && stmt.Partitions != nil {
		// The partitioning parsed but the catalog refuses it, as it does a RANGE
		// whose MAXVALUE partition sync lists before another. The rest of the
		// table still installs.
		stmt.Partitions = nil
		if retry := c.DefineTable(stmt); retry == nil {
			if o.incomplete == nil {
				o.incomplete = fmt.Errorf("partitioning: %w", err)
			}
			applyTableMetadataExtras(c.installedTable(o.name), o.table)
			return nil
		}
	}
	return err
}

func applyTableMetadataExtras(tbl *Table, meta *metadata.TableMetadata) {
	if tbl == nil {
		return
	}
	// CACHE is not a CREATE TABLE option the grammar takes, so the synced
	// table's `cached=on` in CREATE_OPTIONS is read here rather than
	// parsed: SHOW CREATE TABLE writes the attribute back, and a plan that
	// follows it has to NOCACHE the table before any other DDL (8242).
	if isCachedCreateOptions(meta.GetCreateOptions()) {
		tbl.Cached = true
	}
	if meta.GetPlacementPolicy() != "" {
		tbl.PlacementPolicy = resolvePolicyRef(meta.GetPlacementPolicy())
	}
	if meta.GetTiflashReplica() > 0 || len(meta.GetTiflashLocationLabels()) > 0 {
		tbl.TiFlashReplica = int(meta.GetTiflashReplica())
		tbl.TiFlashLocationLabels = slices.Clone(meta.GetTiflashLocationLabels())
	}
	applyDroppedColumnDefaults(tbl, meta)
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

// installMetadataSequence installs a synced sequence as information_schema
// reports it. The values are the server's own, so they are not validated
// again; a value that is not an integer leaves the sequence out.
func (c *Catalog) installMetadataSequence(q *metadata.SequenceMetadata) error {
	db, err := c.resolveDatabase("")
	if err != nil {
		return err
	}
	key := toLower(q.GetName())
	if db.nameTaken(key) {
		return errDupTableQualified(db.Name, q.GetName())
	}
	seq := &Sequence{Name: q.GetName(), Database: db, Cycle: q.GetCycle(), Comment: q.GetComment()}
	for _, f := range []struct {
		name  string
		value string
		into  *int64
	}{
		{"start", q.GetStart(), &seq.Start},
		{"minimum value", q.GetMinValue(), &seq.MinValue},
		{"maximum value", q.GetMaxValue(), &seq.MaxValue},
		{"increment", q.GetIncrement(), &seq.Increment},
		{"cache size", q.GetCacheSize(), &seq.Cache},
	} {
		n, err := strconv.ParseInt(strings.TrimSpace(f.value), 10, 64)
		if err != nil {
			return fmt.Errorf("sequence %s: %s %q is not an integer", q.GetName(), f.name, f.value)
		}
		*f.into = n
	}
	db.Sequences[key] = seq
	return nil
}

// --- Statements ---

// The snapshot carries types, defaults, and definitions as the text TiDB
// reports for them. These builders turn that text into the statement ASTs the
// Define* functions take.

// autoIncrementDefault is what sync stores as the default of an
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
// gets a stand-in rather than a partial definition. A default, AUTO_RANDOM, ON
// UPDATE, generation, or CHECK expression, or a partitioning that does not
// parse is left out instead, as is an index with a key expression that does not
// parse, since TiDB has no column to key in its place; omitted is the first such
// error.
func metadataTableStmt(t *metadata.TableMetadata, snapshotDatabase string) (stmt *nodes.CreateTableStmt, omitted error, err error) {
	omit := func(e error) {
		if omitted == nil {
			omitted = e
		}
	}
	stmt = &nodes.CreateTableStmt{
		Table:    &nodes.TableRef{Name: t.GetName()},
		Global:   t.GetIsGlobalTemporary(),
		OnCommit: t.GetOnCommit(),
	}
	for _, col := range t.GetColumns() {
		if col.GetName() == "" {
			continue
		}
		def, colOmitted, err := metadataColumnDef(col)
		if err != nil {
			return nil, nil, fmt.Errorf("column %q: %w", col.GetName(), err)
		}
		if colOmitted != nil {
			omit(fmt.Errorf("column %q: %w", col.GetName(), colOmitted))
		}
		stmt.Columns = append(stmt.Columns, def)
	}
	tableCols := tableColumnSet(t.GetColumns())
	for _, idx := range t.GetIndexes() {
		if len(idx.GetExpressions()) == 0 {
			continue
		}
		con, err := metadataIndexConstraint(idx, tableCols)
		if err != nil {
			omit(fmt.Errorf("index %q: %w", idx.GetName(), err))
			continue
		}
		// Sync reports TIDB_PK_TYPE, CLUSTERED or NONCLUSTERED.
		if idx.GetPrimary() && t.GetPrimaryKeyType() != "" {
			clustered := strings.EqualFold(t.GetPrimaryKeyType(), "CLUSTERED")
			con.Clustered = &clustered
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
	for _, opt := range []struct{ name, value string }{
		{"ENGINE", t.GetEngine()},
		{"CHARSET", t.GetCharset()},
		{"COLLATE", t.GetCollation()},
		{"COMMENT", syncedComment(t.GetComment())},
		{"SHARD_ROW_ID_BITS", shardRowIDBits(t.GetShardingInfo())},
		{"COMPRESSION", unquoteTableOption(createOption(t.GetCreateOptions(), "compression"))},
	} {
		if opt.value != "" {
			stmt.Options = append(stmt.Options, &nodes.TableOption{Name: opt.name, Value: opt.value})
		}
	}
	if t.GetAutoIdCache() != 0 {
		stmt.Options = append(stmt.Options, &nodes.TableOption{Name: "AUTO_ID_CACHE", Value: strconv.FormatInt(t.GetAutoIdCache(), 10)})
	}
	if t.GetPreSplitRegions() != 0 {
		stmt.Options = append(stmt.Options, &nodes.TableOption{Name: "PRE_SPLIT_REGIONS", Value: strconv.FormatInt(int64(t.GetPreSplitRegions()), 10)})
	}
	if t.GetTtlColumn() != "" && t.GetTtlInterval() != "" {
		ttlEnable := "OFF"
		if t.GetTtlEnable() {
			ttlEnable = "ON"
		}
		stmt.Options = append(stmt.Options,
			&nodes.TableOption{Name: "TTL", Value: fmt.Sprintf("%s + %s", quoteIdent(t.GetTtlColumn()), t.GetTtlInterval())},
			&nodes.TableOption{Name: "TTL_ENABLE", Value: ttlEnable},
		)
		if t.GetTtlJobInterval() != "" {
			stmt.Options = append(stmt.Options, &nodes.TableOption{Name: "TTL_JOB_INTERVAL", Value: t.GetTtlJobInterval()})
		}
	}
	return stmt, omitted, nil
}

func createOption(options, name string) string {
	for _, opt := range splitCreateOptions(options) {
		if k, v, ok := strings.Cut(opt, "="); ok && strings.EqualFold(k, name) {
			return v
		}
	}
	return ""
}

func splitCreateOptions(options string) []string {
	var out []string
	start := -1
	var quote byte
	for i := 0; i < len(options); i++ {
		ch := options[i]
		if quote != 0 {
			if ch == '\\' && quote != '`' && i+1 < len(options) {
				i++
				continue
			}
			if ch == quote {
				if i+1 < len(options) && options[i+1] == quote {
					i++
					continue
				}
				quote = 0
			}
			continue
		}
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			if start >= 0 {
				out = append(out, options[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
		if ch == '\'' || ch == '"' || ch == '`' {
			quote = ch
		}
	}
	if start >= 0 {
		out = append(out, options[start:])
	}
	return out
}

func unquoteTableOption(s string) string {
	if len(s) < 2 {
		return s
	}
	q := s[0]
	if (q != '\'' && q != '"' && q != '`') || s[len(s)-1] != q {
		return s
	}
	inner := s[1 : len(s)-1]
	var b strings.Builder
	b.Grow(len(inner))
	for i := 0; i < len(inner); i++ {
		ch := inner[i]
		if ch == q && i+1 < len(inner) && inner[i+1] == q {
			b.WriteByte(q)
			i++
			continue
		}
		if ch == '\\' && q != '`' && i+1 < len(inner) {
			i++
			switch inner[i] {
			case '0':
				b.WriteByte(0)
			case 'Z':
				b.WriteByte(26)
			default:
				b.WriteByte(inner[i])
			}
			continue
		}
		b.WriteByte(ch)
	}
	return b.String()
}

// shardRowIDBits returns n from the TIDB_ROW_ID_SHARDING_INFO sync records for
// a table created with SHARD_ROW_ID_BITS=n, "SHARD_BITS=n", and "" for any
// other: NOT_SHARDED, NOT_SHARDED(PK_IS_HANDLE), or PK_AUTO_RANDOM_BITS=n,
// which the AUTO_RANDOM column carries.
func shardRowIDBits(info string) string {
	bits, ok := strings.CutPrefix(info, "SHARD_BITS=")
	if !ok || bits == "0" {
		return ""
	}
	return bits
}

// metadataPartitionClause parses the partitioning the snapshot records, as the
// definition renderer prints it.
func metadataPartitionClause(t *metadata.TableMetadata) (*nodes.PartitionClause, error) {
	var buf strings.Builder
	if err := printPartitionClause(&buf, t.GetPartitions()); err != nil {
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

// metadataColumnDef fails on a type that does not parse. A default,
// AUTO_RANDOM, ON UPDATE, or generation expression that does not parse is left
// out, and omitted is the first such error.
func metadataColumnDef(col *metadata.ColumnMetadata) (def *nodes.ColumnDef, omitted error, err error) {
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
	def = &nodes.ColumnDef{Name: col.GetName(), TypeName: typeName, Comment: syncedComment(col.GetComment())}
	if !col.GetNullable() {
		def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrNotNull})
	}
	// IsInvisible is ignored. TiDB v8.5.5 has no column visibility: it answers
	// 1064 for INVISIBLE after a column type, even inside /*!80023 ... */, and
	// information_schema.COLUMNS.EXTRA never says INVISIBLE, so Bytebase's TiDB
	// sync never sets the flag. Keeping it would make SHOW CREATE TABLE write
	// text TiDB refuses.
	def.AutoIncrement = strings.EqualFold(col.GetDefault(), autoIncrementDefault)
	if isAutoRandomDefault(col.GetDefault()) {
		if omitted = setAutoRandom(def, col.GetDefault()); omitted != nil {
			omitted = fmt.Errorf("AUTO_RANDOM: %w", omitted)
		}
	} else {
		def.DefaultValue, omitted = columnDefault(col)
	}
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

// isAutoRandomDefault reports whether a default is the AUTO_RANDOM attribute
// sync stores there: AUTO_RANDOM(shard_bits) or AUTO_RANDOM(shard_bits,
// range_bits), from the table's PK_AUTO_RANDOM_BITS and RANGE BITS.
func isAutoRandomDefault(value string) bool {
	return len(value) >= len(autoRandomSymbol) && strings.EqualFold(value[:len(autoRandomSymbol)], autoRandomSymbol)
}

// setAutoRandom gives def the AUTO_RANDOM attribute value spells.
func setAutoRandom(def *nodes.ColumnDef, value string) error {
	stmt, err := parseFirstStatement[*nodes.CreateTableStmt]("CREATE TABLE `probe` (`c` bigint " + value + ")")
	if err != nil {
		return err
	}
	if len(stmt.Columns) != 1 || !stmt.Columns[0].AutoRandom {
		return fmt.Errorf("no AUTO_RANDOM parsed from %q", value)
	}
	probe := stmt.Columns[0]
	def.AutoRandom, def.AutoRandomShardBits, def.AutoRandomRangeBits = true, probe.AutoRandomShardBits, probe.AutoRandomRangeBits
	return nil
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
	// information_schema reports a BIT default as its literal, b'101', which
	// sync quotes like any other text.
	if lit, ok := expr.(*nodes.StringLit); ok && strings.HasPrefix(strings.ToLower(strings.TrimSpace(col.GetType())), "bit") && bitLiteralText.MatchString(lit.Value) {
		return parseExpr(lit.Value)
	}
	return expr, nil
}

var bitLiteralText = regexp.MustCompile(`^[bB]'[01]*'$`)

// typeTakesDefault reports whether a column type shows a DEFAULT NULL clause,
// which the BLOB, TEXT, and JSON types never do.
func typeTakesDefault(typ string) bool {
	head := strings.ToLower(strings.TrimSpace(typ))
	if i := strings.IndexAny(head, " ("); i >= 0 {
		head = head[:i]
	}
	switch head {
	case "blob", "tinyblob", "mediumblob", "longblob", "text", "tinytext", "mediumtext", "longtext", "json":
		return false
	default:
		return true
	}
}

// metadataIndexConstraint builds an index from its key parts. The index type
// holds either the kind, FULLTEXT, or the access method. Sync reports the
// access method a key without USING gets too, and TiDB's is BTREE, so that one
// stays unwritten. It fails on a key expression that does not parse.
func metadataIndexConstraint(idx *metadata.IndexMetadata, tableCols map[string]bool) (*nodes.Constraint, error) {
	indexType := strings.ToUpper(strings.TrimSpace(idx.GetType()))
	if indexType == "BTREE" {
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
	case idx.GetIsVector() || indexType == "HNSW":
		// A vector index syncs as an index of type HNSW. Left an ordinary index,
		// SHOW CREATE TABLE would write KEY for it and a change could not tell it
		// from one.
		if indexType == "" {
			indexType = "HNSW"
		}
		c.Type, c.IndexType = nodes.ConstrVector, indexType
	default:
		c.Type, c.IndexType = nodes.ConstrIndex, indexType
	}
	var err error
	if c.Columns, c.IndexColumns, err = indexKeyParts(idx, tableCols); err != nil {
		return nil, err
	}
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
	if idx.GetIsGlobal() {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "GLOBAL"})
	}
	return c, nil
}

// indexKeyParts splits an index's key parts: plain column names, or, where a
// part is an expression or has a length or descends, an IndexColumn for each.
// It fails on a key expression that does not parse.
func indexKeyParts(idx *metadata.IndexMetadata, tableCols map[string]bool) ([]string, []*nodes.IndexColumn, error) {
	length := func(i int) int {
		if i < len(idx.GetKeyLength()) {
			return int(idx.GetKeyLength()[i])
		}
		return 0
	}
	descending := func(i int) bool { return i < len(idx.GetDescending()) && idx.GetDescending()[i] }
	plain := true
	for i, expr := range idx.GetExpressions() {
		if isExpressionKey(expr, tableCols) || length(i) > 0 || descending(i) {
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
	for i, expr := range idx.GetExpressions() {
		key := &nodes.IndexColumn{Expr: &nodes.ColumnRef{Column: unquoteIdent(expr)}, Desc: descending(i)}
		if l := length(i); l > 0 {
			key.Length = l
		}
		if isExpressionKey(expr, tableCols) {
			// Unlike MySQL's, TiDB's EXPRESSION is SQL as written: quotes are
			// doubled, not backslash-escaped, so it parses as it is.
			parsed, err := parseExpr(expr)
			if err != nil {
				return nil, nil, err
			}
			key.Expr, key.Functional = parsed, true
		}
		keys = append(keys, key)
	}
	return nil, keys, nil
}

func tableColumnSet(cols []*metadata.ColumnMetadata) map[string]bool {
	set := make(map[string]bool, len(cols))
	for _, col := range cols {
		if name := col.GetName(); name != "" {
			set[strings.ToLower(name)] = true
		}
	}
	return set
}

// isExpressionKey reports whether a key part is a functional expression, as
// in "(lower(name))" or "lower(`name`)", rather than an existing column.
func isExpressionKey(s string, tableCols map[string]bool) bool {
	if tableCols[strings.ToLower(unquoteIdent(s))] {
		return false
	}
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

// referenceAction maps a synced rule to the action written. information_schema
// reports NO ACTION for a foreign key without the clause too, which is the far
// more common one, so NO ACTION stays unwritten, as the definition renderer
// leaves it; TiDB's SHOW CREATE TABLE writes it only where it was written.
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
	default:
		return nodes.RefActNone
	}
}

// metadataViewStmt keeps the view's text, less the qualifier of the snapshot's
// database, along with its parsed body. TiDB's VIEW_DEFINITION qualifies every
// table with that database, which the view resolves to without it.
// DefineView installs a view whose body does not analyze, so a view may name
// one that installs after it; reanalyzeViews retries it then. nameViewColumns
// names the columns once the body resolves.
func metadataViewStmt(v *metadata.ViewMetadata, database, target string) (*nodes.CreateViewStmt, error) {
	body := stripViewBodyDatabaseQualifier(v.GetDefinition(), database, target)
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

// stripViewBodyDatabaseQualifier removes the qualifier of database dbName from
// a view body. The body is parsed so that only a qualifier in the database
// position goes: a table alias or table spelled like the database qualifies
// its columns and stays. A body that does not parse is returned unchanged.
// A table whose name a CTE, another table's or a derived table's alias or another database's table
// also has keeps a qualifier, target's, since without one a reference to it
// would resolve to the CTE or the alias, or be ambiguous.
func stripViewBodyDatabaseQualifier(body, dbName, target string) string {
	if dbName == "" || body == "" {
		return body
	}
	qualifier := quoteIdent(dbName) + "."
	if !strings.Contains(body, qualifier) {
		return body
	}
	// A stored body may end in WITH CHECK OPTION, so it parses only as a whole view.
	const prefix = "CREATE VIEW v AS "
	list, err := tidbparser.Parse(prefix + body)
	if err != nil || list == nil || len(list.Items) != 1 {
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
		case *nodes.TableRef:
			schema, table, loc, tableRef = n.Schema, n.Name, n.Loc, true
			if n.Alias != "" && !ctes[strings.ToLower(n.Name)] {
				// Aliased, the table is known by its alias, which nothing else
				// can shadow here.
				table = ""
			}
		case *nodes.ColumnRef:
			schema, table, loc = n.Schema, n.Table, n.Loc
		case *nodes.FuncCallExpr:
			schema, loc = n.Schema, n.Loc
		case *nodes.DefaultExpr:
			// The node starts at DEFAULT; its qualifier is the first inside it.
			schema, table, loc = n.Schema, n.Table, n.Loc
			if begin, end := loc.Start-len(prefix), loc.End-len(prefix); schema == dbName && begin >= 0 && end <= len(body) && begin < end {
				if i := strings.Index(body[begin:end], qualifier); i >= 0 {
					loc.Start += i
				}
			}
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
// resolved body left unnamed, and, where the body names as many by other names,
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
	}
	if len(v.ColumnMetadata) != len(v.Columns) {
		v.ColumnMetadata = make([]ViewColumn, len(v.Columns))
	}
	for i, name := range v.Columns {
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
		if v.AnalyzedQuery == nil && vc.ColumnType == "" {
			vc.ColumnType = col.GetType()
		}
		if vc.Charset == "" {
			vc.Charset = col.GetCharacterSet()
		}
		if vc.Collation == "" {
			vc.Collation = col.GetCollation()
		}
	}
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
// a table or view needs a column.
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

// parseTypeName parses a column type as TiDB reports it, such as
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

// parseFirstStatement parses sql and returns its first statement of type T.
func parseFirstStatement[T nodes.Node](sql string) (T, error) {
	var zero T
	if strings.TrimSpace(sql) == "" {
		return zero, errors.New("empty definition")
	}
	list, err := tidbparser.Parse(sql)
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

func parseMySQLFirst[T mysqlast.Node](sql string) (T, error) {
	var zero T
	if strings.TrimSpace(sql) == "" {
		return zero, errors.New("empty definition")
	}
	list, err := mysqlparser.Parse(sql)
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

func formatTiDBDefiner(definer string) string {
	if definer == "" || strings.HasPrefix(definer, "`") || strings.HasPrefix(definer, "'") || strings.EqualFold(definer, "CURRENT_USER") {
		return definer
	}
	if user, host, ok := strings.Cut(definer, "@"); ok {
		return fmt.Sprintf("`%s`@`%s`", strings.ReplaceAll(user, "`", "``"), strings.ReplaceAll(host, "`", "``"))
	}
	return definer
}

func toTiDBDataType(dt *mysqlast.DataType) *nodes.DataType {
	if dt == nil {
		return nil
	}
	return &nodes.DataType{
		Name:       dt.Name,
		Length:     dt.Length,
		Scale:      dt.Scale,
		HasScale:   dt.Scale > 0,
		Unsigned:   dt.Unsigned,
		Zerofill:   dt.Zerofill,
		Charset:    dt.Charset,
		Collate:    dt.Collate,
		Binary:     dt.Binary,
		EnumValues: append([]string(nil), dt.EnumValues...),
		ArrayDim:   dt.ArrayDim,
		SRID:       dt.SRID,
	}
}

func parseTiDBRoutine(definition string, isProcedure bool) (*nodes.CreateFunctionStmt, error) {
	src, err := parseMySQLFirst[*mysqlast.CreateFunctionStmt](definition)
	if err != nil {
		return nil, err
	}
	var name *nodes.TableRef
	if src.Name != nil {
		name = &nodes.TableRef{Name: src.Name.Name}
	}
	params := make([]*nodes.FuncParam, 0, len(src.Params))
	for _, p := range src.Params {
		params = append(params, &nodes.FuncParam{
			Direction: p.Direction,
			Name:      p.Name,
			TypeName:  toTiDBDataType(p.TypeName),
		})
	}
	chars := make([]*nodes.RoutineCharacteristic, 0, len(src.Characteristics))
	for _, ch := range src.Characteristics {
		chars = append(chars, &nodes.RoutineCharacteristic{Name: ch.Name, Value: ch.Value})
	}
	return &nodes.CreateFunctionStmt{
		IsProcedure:     isProcedure,
		Definer:         formatTiDBDefiner(src.Definer),
		Name:            name,
		Params:          params,
		Returns:         toTiDBDataType(src.Returns),
		BodyText:        src.BodyText,
		Characteristics: chars,
	}, nil
}

func parseTiDBEvent(definition string) (*nodes.CreateEventStmt, error) {
	src, err := parseMySQLFirst[*mysqlast.CreateEventStmt](definition)
	if err != nil {
		return nil, err
	}
	var schedule *nodes.EventSchedule
	if src.Schedule != nil {
		schedule = &nodes.EventSchedule{RawText: src.Schedule.RawText}
	}
	return &nodes.CreateEventStmt{
		Definer:      formatTiDBDefiner(src.Definer),
		Name:         src.Name,
		Schedule:     schedule,
		OnCompletion: src.OnCompletion,
		Enable:       src.Enable,
		Comment:      src.Comment,
		BodyText:     src.BodyText,
	}, nil
}
