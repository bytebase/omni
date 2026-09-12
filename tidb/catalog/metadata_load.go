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
// reference one that installs after it. A view may likewise name a view that
// installs after it.
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
		// A database already in the catalog keeps its own character set and
		// collation otherwise, and a table that overrides neither would resolve
		// its columns differently than the same snapshot loaded into an empty
		// catalog.
		stmt := metadataDatabaseStmt(meta)
		if c.GetDatabase(name) == nil {
			_ = c.DefineDatabase(stmt)
		} else if len(stmt.Options) > 0 {
			_ = c.alterDatabase(&nodes.AlterDatabaseStmt{Name: name, Options: stmt.Options})
		}
		c.SetCurrentDatabase(name)
	}
	fkChecks := c.ForeignKeyChecks()
	c.SetForeignKeyChecks(false)
	defer c.SetForeignKeyChecks(fkChecks)

	var views []*metadataObject
	for _, o := range collectMetadataObjects(meta) {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		err := c.installMetadataObject(o)
		if err == nil {
			if o.kind == metadataView {
				views = append(views, o)
			}
			continue
		}
		if standInErr := c.installStandIn(o); standInErr != nil {
			report.Missing[o.key()] = standInErr
			continue
		}
		report.Degraded[o.key()] = err
	}
	return report, c.reanalyzeViews(ctx, views)
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
		stmt, err := metadataViewStmt(o.view)
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
		return nil
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

func metadataDatabaseStmt(meta *metadata.DatabaseSchemaMetadata) *nodes.CreateDatabaseStmt {
	stmt := &nodes.CreateDatabaseStmt{Name: meta.GetName()}
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
		if len(idx.GetExpressions()) == 0 {
			continue
		}
		con := metadataIndexConstraint(idx)
		// Sync reports TIDB_PK_TYPE, CLUSTERED or NONCLUSTERED.
		if idx.GetPrimary() && t.GetPrimaryKeyType() != "" {
			clustered := strings.EqualFold(t.GetPrimaryKeyType(), "CLUSTERED")
			con.Clustered = &clustered
		}
		stmt.Constraints = append(stmt.Constraints, con)
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
	def.AutoIncrement = strings.EqualFold(col.GetDefault(), autoIncrementDefault)
	def.DefaultValue = columnDefault(col)
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

// columnDefault parses the default the snapshot records for a column, or
// returns nil where the catalog takes none: AUTO_INCREMENT is a flag rather
// than a default, a generated column has no default, and sync reports DEFAULT
// NULL even for the types that show no default.
func columnDefault(col *metadata.ColumnMetadata) nodes.ExprNode {
	value := col.GetDefault()
	switch {
	case value == "" || col.GetGeneration() != nil:
		return nil
	case strings.EqualFold(value, autoIncrementDefault):
		return nil
	case strings.EqualFold(value, "NULL") && !typeTakesDefault(col.GetType()):
		return nil
	}
	expr, err := parseExpr(value)
	if err != nil {
		return nil
	}
	return expr
}

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
// stays unwritten.
func metadataIndexConstraint(idx *metadata.IndexMetadata) *nodes.Constraint {
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
	default:
		c.Type, c.IndexType = nodes.ConstrIndex, indexType
	}
	c.Columns, c.IndexColumns = indexKeyParts(idx)
	// A primary key is always visible.
	if !idx.GetVisible() && !idx.GetPrimary() {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "INVISIBLE"})
	}
	if idx.GetComment() != "" {
		c.IndexOptions = append(c.IndexOptions, &nodes.IndexOption{Name: "COMMENT", Value: &nodes.StringLit{Value: idx.GetComment()}})
	}
	return c
}

// indexKeyParts splits an index's key parts: plain column names, or, where a
// part is an expression or has a length or descends, an IndexColumn for each.
func indexKeyParts(idx *metadata.IndexMetadata) ([]string, []*nodes.IndexColumn) {
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
	if plain {
		var columns []string
		for _, expr := range idx.GetExpressions() {
			columns = append(columns, unquoteIdent(expr))
		}
		return columns, nil
	}
	var keys []*nodes.IndexColumn
	for i, expr := range idx.GetExpressions() {
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
		keys = append(keys, key)
	}
	return nil, keys
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

// metadataViewStmt keeps the view's text as it is along with its parsed body.
// DefineView installs a view whose body does not analyze, so a view may name
// one that installs after it; reanalyzeViews retries it then. nameViewColumns
// names the columns once the body resolves.
func metadataViewStmt(v *metadata.ViewMetadata) (*nodes.CreateViewStmt, error) {
	sel, err := parseFirstStatement[*nodes.SelectStmt](v.GetDefinition())
	if err != nil {
		return nil, err
	}
	return &nodes.CreateViewStmt{OrReplace: true, Name: &nodes.TableRef{Name: v.GetName()}, Select: sel, SelectText: v.GetDefinition()}, nil
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
