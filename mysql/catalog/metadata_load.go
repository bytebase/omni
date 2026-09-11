package catalog

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
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
}

// LoadMetadata installs a Bytebase schema snapshot of a MySQL database into
// the database meta.Name, creating it if needed, and leaves that database
// selected. Foreign key checks are off during the load, so a table may
// reference one that installs after it.
//
// An object that fails to install does not stop the load. A table is replaced
// by a stand-in with the same column names, all TEXT; a view by one that
// selects NULL under each of its column names; a trigger by one with an empty
// body; and an event by one with only its name. A function or procedure that
// fails is left out. The error is non-nil only when ctx is done.
func (c *Catalog) LoadMetadata(ctx context.Context, meta *metadata.DatabaseSchemaMetadata) (*LoadMetadataReport, error) {
	report := &LoadMetadataReport{
		Degraded: make(map[string]error),
		Missing:  make(map[string]error),
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	if name := meta.GetName(); name != "" {
		if c.GetDatabase(name) == nil {
			// A database without options always installs.
			_ = c.DefineDatabase(&nodes.CreateDatabaseStmt{Name: name})
		}
		c.SetCurrentDatabase(name)
	}
	fkChecks := c.ForeignKeyChecks()
	c.SetForeignKeyChecks(false)
	defer c.SetForeignKeyChecks(fkChecks)

	for _, o := range collectMetadataObjects(meta) {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		err := c.installMetadataObject(o)
		if err == nil {
			continue
		}
		switch standInErr := c.installStandIn(o); {
		case standInErr == nil:
			report.Degraded[o.key()] = err
		case errors.Is(standInErr, errNoStandIn):
			report.Missing[o.key()] = err
		default:
			report.Missing[o.key()] = standInErr
		}
	}
	return report, nil
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
	// parent is a trigger's table.
	parent string

	table     *metadata.TableMetadata
	view      *metadata.ViewMetadata
	function  *metadata.FunctionMetadata
	procedure *metadata.ProcedureMetadata
	trigger   *metadata.TriggerMetadata
	event     *metadata.EventMetadata
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
			out = append(out, &metadataObject{kind: metadataTable, name: t.GetName(), table: t})
			for _, trg := range t.GetTriggers() {
				if trg.GetName() != "" {
					out = append(out, &metadataObject{kind: metadataTrigger, name: trg.GetName(), parent: t.GetName(), trigger: trg})
				}
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataView, name: v.GetName(), view: v})
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
	case metadataTable:
		stmt, err := metadataTableStmt(o.table)
		if err != nil {
			return err
		}
		return c.DefineTable(stmt)
	case metadataView:
		return c.DefineView(metadataViewStmt(o.view))
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

// metadataTableStmt fails on a column whose type does not parse, so the table
// gets a stand-in rather than a partial definition. A default, ON UPDATE,
// generation, or CHECK expression that does not parse is dropped instead.
// Partitions are left out: they do not change what the columns resolve to.
func metadataTableStmt(t *metadata.TableMetadata) (*nodes.CreateTableStmt, error) {
	stmt := &nodes.CreateTableStmt{Table: &nodes.TableRef{Name: t.GetName()}}
	for _, col := range t.GetColumns() {
		if col.GetName() == "" {
			continue
		}
		def, err := metadataColumnDef(col)
		if err != nil {
			return nil, fmt.Errorf("column %q: %w", col.GetName(), err)
		}
		stmt.Columns = append(stmt.Columns, def)
	}
	for _, idx := range t.GetIndexes() {
		if len(idx.GetExpressions()) > 0 {
			stmt.Constraints = append(stmt.Constraints, metadataIndexConstraint(idx))
		}
	}
	for _, fk := range t.GetForeignKeys() {
		stmt.Constraints = append(stmt.Constraints, metadataForeignKey(fk))
	}
	for _, check := range t.GetCheckConstraints() {
		if expr, err := parseExpr(check.GetExpression()); err == nil {
			stmt.Constraints = append(stmt.Constraints, &nodes.Constraint{Type: nodes.ConstrCheck, Name: check.GetName(), Expr: expr})
		}
	}
	for _, opt := range []struct{ name, value string }{
		{"ENGINE", t.GetEngine()},
		{"CHARSET", t.GetCharset()},
		{"COLLATE", t.GetCollation()},
		{"COMMENT", t.GetComment()},
	} {
		if opt.value != "" {
			stmt.Options = append(stmt.Options, &nodes.TableOption{Name: opt.name, Value: opt.value})
		}
	}
	return stmt, nil
}

func metadataColumnDef(col *metadata.ColumnMetadata) (*nodes.ColumnDef, error) {
	typeName, err := parseTypeName(col.GetType())
	if err != nil {
		return nil, err
	}
	// The type text may carry its own character set and collation.
	if typeName.Charset == "" {
		typeName.Charset = col.GetCharacterSet()
	}
	if typeName.Collate == "" {
		typeName.Collate = col.GetCollation()
	}
	// SRID 0 is a valid SRID, so presence decides.
	if col.Srid != nil {
		typeName.SRID = int(col.GetSrid())
		typeName.HasSRID = true
	}
	def := &nodes.ColumnDef{Name: col.GetName(), TypeName: typeName, Comment: col.GetComment()}
	if !col.GetNullable() {
		def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrNotNull})
	}
	if col.GetIsInvisible() {
		def.Constraints = append(def.Constraints, &nodes.ColumnConstraint{Type: nodes.ColConstrInvisible})
	}
	autoIncrement := strings.EqualFold(col.GetDefault(), autoIncrementDefault)
	def.AutoIncrement = autoIncrement
	// Sync reports DEFAULT NULL for every nullable column, including types that
	// take no DEFAULT clause, which would fail the whole table.
	if !autoIncrement && col.GetDefault() != "" && col.GetGeneration() == nil &&
		(!strings.EqualFold(col.GetDefault(), "NULL") || typeTakesDefault(col.GetType())) {
		if expr, err := parseExpr(col.GetDefault()); err == nil {
			def.DefaultValue = expr
		}
	}
	if col.GetOnUpdate() != "" {
		if expr, err := parseExpr(col.GetOnUpdate()); err == nil {
			def.OnUpdate = expr
		}
	}
	if gen := col.GetGeneration(); gen.GetExpression() != "" {
		if expr, err := parseExpr(gen.GetExpression()); err == nil {
			def.Generated = &nodes.GeneratedColumn{Expr: expr, Stored: gen.GetType() == metadata.GenerationMetadata_TYPE_STORED}
		}
	}
	return def, nil
}

// typeTakesDefault reports whether a column type accepts a literal DEFAULT
// clause, which the BLOB types, JSON, and GEOMETRY do not.
func typeTakesDefault(typ string) bool {
	head := strings.ToLower(strings.TrimSpace(typ))
	if i := strings.IndexAny(head, " ("); i >= 0 {
		head = head[:i]
	}
	switch head {
	case "blob", "tinyblob", "mediumblob", "longblob", "json", "geometry":
		return false
	default:
		return true
	}
}

// metadataIndexConstraint builds an index from its key parts. The index type
// holds either the kind, FULLTEXT or SPATIAL, or the access method. A key part
// that is an expression, or that has a length or descends, goes into
// IndexColumns; plain columns go into Columns.
func metadataIndexConstraint(idx *metadata.IndexMetadata) *nodes.Constraint {
	indexType := strings.ToUpper(strings.TrimSpace(idx.GetType()))
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
	default:
		c.Type, c.IndexType = nodes.ConstrIndex, indexType
	}
	length := func(i int) int {
		if i < len(idx.GetKeyLength()) {
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
	for i, expr := range idx.GetExpressions() {
		if plain {
			c.Columns = append(c.Columns, unquoteIdent(expr))
			continue
		}
		key := &nodes.IndexColumn{Expr: &nodes.ColumnRef{Column: unquoteIdent(expr)}, Desc: descending(i)}
		if l := length(i); l > 0 {
			key.Length = l
		}
		if isExpressionKey(expr) {
			// On a parse failure the key stays a column, so the index still installs.
			if parsed, err := parseExpr(expr); err == nil {
				key.Expr = parsed
			}
		}
		c.IndexColumns = append(c.IndexColumns, key)
	}
	return c
}

// isExpressionKey reports whether a key part is a functional expression, as
// in "(lower(name))" or "lower(`name`)", rather than a column.
func isExpressionKey(s string) bool {
	return strings.Contains(s, "(")
}

func metadataForeignKey(fk *metadata.ForeignKeyMetadata) *nodes.Constraint {
	return &nodes.Constraint{
		Type:       nodes.ConstrForeignKey,
		Name:       fk.GetName(),
		Columns:    fk.GetColumns(),
		RefTable:   &nodes.TableRef{Schema: fk.GetReferencedSchema(), Name: fk.GetReferencedTable()},
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

// metadataViewStmt keeps the view's text as it is and its parsed body when it
// parses. DefineView installs a view whose body does not analyze, so a view
// may name one that installs after it.
func metadataViewStmt(v *metadata.ViewMetadata) *nodes.CreateViewStmt {
	stmt := &nodes.CreateViewStmt{OrReplace: true, Name: &nodes.TableRef{Name: v.GetName()}, SelectText: v.GetDefinition()}
	for _, col := range v.GetColumns() {
		if col.GetName() != "" {
			stmt.Columns = append(stmt.Columns, col.GetName())
		}
	}
	if sel, err := parseFirstStatement[*nodes.SelectStmt](v.GetDefinition()); err == nil {
		stmt.Select = sel
	}
	return stmt
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

func unquoteIdent(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return strings.ReplaceAll(s[1:len(s)-1], "``", "`")
	}
	return s
}
