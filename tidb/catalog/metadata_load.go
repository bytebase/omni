package catalog

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/tidb/ast"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// LoadMetadataReport lists the snapshot objects LoadMetadata could not install
// as the snapshot defines them. Keys are "table:<name>" and "view:<name>".
type LoadMetadataReport struct {
	// Degraded maps each object replaced by a stand-in to the error its
	// definition failed with.
	Degraded map[string]error
	// Missing maps each object the catalog has nothing for to the error that
	// left it out.
	Missing map[string]error
}

// LoadMetadata installs a Bytebase schema snapshot of a TiDB database into the
// database meta.Name, creating it if needed, and leaves that database
// selected. Foreign key checks are off during the load, so a table may
// reference one that installs after it.
//
// An object that fails to install does not stop the load: a table is replaced
// by a stand-in with the same column names, all TEXT, and a view by one that
// selects NULL under each of its column names. The error is non-nil only when
// ctx is done.
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
		if standInErr := c.installStandIn(o); standInErr != nil {
			report.Missing[o.key()] = standInErr
			continue
		}
		report.Degraded[o.key()] = err
	}
	return report, nil
}

// --- Objects ---

// metadataObjectKind lists the kinds in the order they install. Nothing else
// orders them: foreign key checks are off, and a view may name an object not
// yet installed.
type metadataObjectKind int

const (
	metadataTable metadataObjectKind = iota
	metadataView
)

// metadataObject is one install step. The metadata pointer matching kind is
// set.
type metadataObject struct {
	kind  metadataObjectKind
	name  string
	table *metadata.TableMetadata
	view  *metadata.ViewMetadata
}

// key names the object in the report.
func (o *metadataObject) key() string {
	if o.kind == metadataTable {
		return "table:" + o.name
	}
	return "view:" + o.name
}

// --- Collecting ---

// collectMetadataObjects lists the snapshot's objects in install order: by
// kind, then name.
func collectMetadataObjects(meta *metadata.DatabaseSchemaMetadata) []*metadataObject {
	var out []*metadataObject
	for _, s := range meta.GetSchemas() {
		for _, t := range s.GetTables() {
			if t.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataTable, name: t.GetName(), table: t})
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() != "" {
				out = append(out, &metadataObject{kind: metadataView, name: v.GetName(), view: v})
			}
		}
	}
	slices.SortFunc(out, func(a, b *metadataObject) int {
		return cmp.Or(cmp.Compare(a.kind, b.kind), cmp.Compare(a.name, b.name))
	})
	return out
}

// --- Installing ---

func (c *Catalog) installMetadataObject(o *metadataObject) error {
	if o.kind == metadataView {
		return c.DefineView(metadataViewStmt(o.view))
	}
	stmt, err := metadataTableStmt(o.table)
	if err != nil {
		return err
	}
	return c.DefineTable(stmt)
}

// --- Statements ---

// The snapshot carries types, defaults, and definitions as the text TiDB
// reports for them. These builders turn that text into the statement ASTs the
// Define* functions take.

// autoIncrementDefault is what sync stores as the default of an
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
// clause, which the BLOB types and JSON do not.
func typeTakesDefault(typ string) bool {
	head := strings.ToLower(strings.TrimSpace(typ))
	if i := strings.IndexAny(head, " ("); i >= 0 {
		head = head[:i]
	}
	switch head {
	case "blob", "tinyblob", "mediumblob", "longblob", "json":
		return false
	default:
		return true
	}
}

// metadataIndexConstraint builds an index from its key parts. The index type
// holds either the kind, FULLTEXT, or the access method. A key part that is an
// expression, or that has a length or descends, goes into IndexColumns; plain
// columns go into Columns.
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

// --- Stand-ins ---

// installStandIn installs o's stand-in.
func (c *Catalog) installStandIn(o *metadataObject) error {
	if o.kind == metadataView {
		return c.DefineView(standInViewStmt(o.view))
	}
	return c.DefineTable(standInTableStmt(o.table))
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
