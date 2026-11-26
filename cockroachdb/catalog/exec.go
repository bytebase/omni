package catalog

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/cockroachdb/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

// Catalog is one database's schema, changed by DDL as CockroachDB changes it.
// It starts with the schema public, as a new database does.
//
// Exec takes what CockroachDB takes. It refuses a statement with
// CockroachDB's SQLSTATE where the object a statement names is missing, taken
// or still needed by another, the cases TestSimulation pins against the
// server; a form the catalog does not model is refused with an Error that
// carries no code.
type Catalog struct {
	l *sdlLoader
	// pending is the table ALTER TABLE is changing, whose columns a foreign
	// key it adds names.
	pending *metadatapb.TableMetadata
}

// New returns the catalog of an empty database.
func New() *Catalog {
	l := newSDLLoader()
	l.schema("public")
	return &Catalog{l: l}
}

// Metadata is what the catalog holds, as a sync of the database would state
// it.
func (c *Catalog) Metadata() *metadatapb.DatabaseSchemaMetadata { return c.l.meta }

// ExecResult is one statement Exec ran: its text, the 1-based line it starts
// on, and the error CockroachDB would refuse it with, or nil.
type ExecResult struct {
	SQL   string
	Line  int
	Error error
}

// Error is a statement CockroachDB refuses. Code is the SQLSTATE CockroachDB
// refuses it with, or empty for a form the catalog does not model.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }

// SQLSTATEs CockroachDB refuses DDL with.
const (
	codeDuplicateColumn   = "42701"
	codeDuplicateObject   = "42710"
	codeDuplicateRelation = "42P07"
	codeDuplicateSchema   = "42P06"
	codeUndefinedColumn   = "42703"
	codeUndefinedObject   = "42704"
	codeUndefinedTable    = "42P01"
	codeUndefinedSchema   = "3F000"
	codeWrongObjectType   = "42809"
	codeDependentObjects  = "2BP01"
	codeInvalidTableDef   = "42P16"
	codeFeatureNotSupport = "0A000"
	codeInvalidParameter  = "22023"
	// CockroachDB refuses a foreign key no unique key of the referenced
	// table covers with 23503, where PostgreSQL uses 42830.
	codeForeignKeyViolation = "23503"
)

// Exec runs a script on the catalog, one statement at a time, and stops at
// the first statement CockroachDB would refuse, which changes nothing. It
// returns a syntax error, and changes nothing, when the script does not parse.
func (c *Catalog) Exec(sql string) ([]ExecResult, error) {
	type parsed struct {
		text string
		line int
		ast  tree.Statement
	}
	var stmts []parsed
	for _, seg := range analysis.SplitSegments(sql) {
		if seg.Empty() {
			continue
		}
		list, err := parser.Parse(seg.Text)
		if err != nil {
			return nil, err
		}
		line := strings.Count(sql[:seg.ByteStart], "\n") + 1
		line += strings.Count(seg.Text[:len(seg.Text)-len(strings.TrimLeft(seg.Text, " \t\r\n"))], "\n")
		for _, s := range list {
			stmts = append(stmts, parsed{text: s.SQL, line: line, ast: s.AST})
		}
	}
	var results []ExecResult
	for _, s := range stmts {
		saved := proto.Clone(c.l.meta).(*metadatapb.DatabaseSchemaMetadata)
		err := c.exec(s.ast)
		if err == nil {
			err = c.l.resolveKeyless()
		}
		results = append(results, ExecResult{SQL: s.text, Line: s.line, Error: err})
		if err != nil {
			c.restore(saved)
			break
		}
	}
	return results, nil
}

// restore puts back what the catalog held before a refused statement.
func (c *Catalog) restore(meta *metadatapb.DatabaseSchemaMetadata) {
	c.l.meta = meta
	c.l.keyless = nil
	c.l.schemas = make(map[string]*metadatapb.SchemaMetadata)
	for _, s := range meta.Schemas {
		c.l.schemas[s.Name] = s
	}
}

func ptr[T any](v T) *T { return &v }

// refuse is an Error with CockroachDB's SQLSTATE.
func refuse(code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

func (c *Catalog) exec(stmt tree.Statement) error {
	switch n := stmt.(type) {
	case *tree.CreateSchema:
		return c.createSchema(n)
	case *tree.AlterSchema:
		return c.alterSchema(n)
	case *tree.DropSchema:
		return c.dropSchema(n)
	case *tree.SetVar:
		if strings.EqualFold(n.Name, "search_path") {
			// The catalog resolves every name in public.
			return &Error{Message: "SET search_path is not supported"}
		}
	case *tree.CommentOnSchema:
		if c.l.schemas[string(n.Name.SchemaName)] == nil {
			return refuse(codeUndefinedSchema, "unknown schema %q", n.Name.SchemaName)
		}
	case *tree.CreateTable:
		return c.createTable(n)
	case *tree.CreateView:
		return c.createView(n)
	case *tree.CreateSequence:
		return c.createSequence(n)
	case *tree.AlterSequence:
		return c.alterSequence(n)
	case *tree.CreateIndex:
		return c.createIndex(n)
	case *tree.CreateType:
		return c.createType(n)
	case *tree.AlterType:
		return c.alterType(n)
	case *tree.DropType:
		return c.dropType(n)
	case *tree.RenameTable:
		return c.renameRelation(n)
	case *tree.AlterTableSetSchema:
		return c.setRelationSchema(n)
	case *tree.DropTable:
		return c.dropRelations(n.Names, n.IfExists, n.DropBehavior, relTable, false)
	case *tree.DropView:
		kind := relView
		if n.IsMaterialized {
			kind = relMaterializedView
		}
		return c.dropRelations(n.Names, n.IfExists, n.DropBehavior, kind, false)
	case *tree.DropSequence:
		return c.dropRelations(n.Names, n.IfExists, n.DropBehavior, relSequence, false)
	case *tree.Truncate:
		for i := range n.Tables {
			table, err := c.table(&n.Tables[i])
			if err != nil {
				return err
			}
			if fks := c.referencing(schemaOf(&n.Tables[i]), table.Name, ""); len(fks) > 0 && n.DropBehavior != tree.DropCascade {
				return refuse("XXUUU", "%q is referenced by foreign key %q", table.Name, fks[0].Name)
			}
		}
		return nil
	case *tree.AlterTable:
		return c.alterTable(n)
	case *tree.RenameIndex:
		return c.renameIndex(n)
	case *tree.AlterIndexVisible:
		return c.alterIndexVisible(n)
	case *tree.DropIndex:
		return c.dropIndex(n)
	case *tree.CommentOnTable:
		if _, err := c.table(ptr(n.Table.ToTableName())); err != nil {
			return err
		}
	case *tree.CommentOnColumn:
		table, err := c.table(ptr(n.TableName.ToTableName()))
		if err != nil {
			return err
		}
		if findColumn(table, string(n.ColumnName)) == nil {
			return refuse(codeUndefinedColumn, "column %q does not exist", n.ColumnName)
		}
	case *tree.CommentOnIndex:
		if _, _, err := c.index(&n.Index, false); err != nil {
			return err
		}
	case *tree.CommentOnConstraint:
		table, err := c.table(ptr(n.Table.ToTableName()))
		if err != nil {
			return err
		}
		if !constraintNamed(table, string(n.Constraint)) {
			return refuse(codeUndefinedObject, "constraint %q of relation %q does not exist", n.Constraint, table.Name)
		}
	case *tree.CommentOnType:
		if _, _, _, err := c.userType(n.Name); err != nil {
			return err
		}
	}
	return c.l.load(stmt)
}

// The kinds of relation a schema holds, which share its namespace.
type relKind int

const (
	relNone relKind = iota
	relTable
	relView
	relMaterializedView
	relSequence
)

func (k relKind) String() string {
	switch k {
	case relTable:
		return "table"
	case relView:
		return "view"
	case relMaterializedView:
		return "materialized view"
	case relSequence:
		return "sequence"
	}
	return "relation"
}

// relation reports the kind of the relation a schema holds by name.
func relation(s *metadatapb.SchemaMetadata, name string) relKind {
	switch {
	case s == nil:
		return relNone
	case schemaTable(s, name) != nil:
		return relTable
	case findView(s, name) != nil:
		return relView
	case findMaterializedView(s, name) != nil:
		return relMaterializedView
	case findSequence(s, name) != nil:
		return relSequence
	}
	return relNone
}

// typeTaken reports whether a schema holds a type by name; a relation is
// one, since each has a row type.
func typeTaken(s *metadatapb.SchemaMetadata, name string) bool {
	return slices.ContainsFunc(s.EnumTypes, func(e *metadatapb.EnumTypeMetadata) bool { return e.Name == name }) ||
		slices.ContainsFunc(s.CompositeTypes, func(t *metadatapb.CompositeTypeMetadata) bool { return t.Name == name })
}

// schemaFor returns the schema a name is in, refused when there is none.
func (c *Catalog) schemaFor(schema string) (*metadatapb.SchemaMetadata, error) {
	if schema == "" {
		schema = "public"
	}
	s := c.l.schemas[schema]
	if s == nil {
		return nil, refuse(codeUndefinedSchema, "schema %q does not exist", schema)
	}
	return s, nil
}

func objectSchema(name *tree.UnresolvedObjectName) string {
	if s := name.Schema(); s != "" {
		return s
	}
	return "public"
}

// table finds a table, refused when there is none.
func (c *Catalog) table(tn *tree.TableName) (*metadatapb.TableMetadata, error) {
	s := c.l.schemas[schemaOf(tn)]
	if table := schemaTable(s, string(tn.ObjectName)); s != nil && table != nil {
		return table, nil
	}
	if k := relation(s, string(tn.ObjectName)); k != relNone {
		return nil, refuse(codeWrongObjectType, "%q is not a table", tn.ObjectName)
	}
	return nil, refuse(codeUndefinedTable, "relation %q does not exist", tn.String())
}

func (c *Catalog) createSchema(n *tree.CreateSchema) error {
	name := string(n.Schema.SchemaName)
	if name == "" {
		name = n.AuthRole.Name
	}
	if c.l.schemas[name] != nil {
		if n.IfNotExists {
			return nil
		}
		return refuse(codeDuplicateSchema, "schema %q already exists", name)
	}
	return c.l.load(n)
}

func (c *Catalog) alterSchema(n *tree.AlterSchema) error {
	s, err := c.schemaFor(string(n.Schema.SchemaName))
	if err != nil {
		return err
	}
	rename, ok := n.Cmd.(*tree.AlterSchemaRename)
	if !ok {
		return c.l.load(n)
	}
	to := string(rename.NewName)
	if c.l.schemas[to] != nil {
		return refuse(codeDuplicateSchema, "schema %q already exists", to)
	}
	from := s.Name
	delete(c.l.schemas, from)
	s.Name = to
	c.l.schemas[to] = s
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, fk := range table.ForeignKeys {
			if fk.ReferencedSchema == from {
				fk.ReferencedSchema = to
			}
		}
	})
	// CockroachDB holds a column's type by ID, so the column follows it.
	for _, e := range s.EnumTypes {
		c.retype(from, e.Name, to, e.Name)
	}
	for _, t := range s.CompositeTypes {
		c.retype(from, t.Name, to, t.Name)
	}
	return nil
}

func (c *Catalog) dropSchema(n *tree.DropSchema) error {
	for _, prefix := range n.Names {
		name := string(prefix.SchemaName)
		s := c.l.schemas[name]
		if s == nil {
			if n.IfExists {
				continue
			}
			return refuse(codeUndefinedSchema, "unknown schema %q", name)
		}
		if name == "public" {
			return refuse(codeUndefinedSchema, "cannot drop schema %q", name)
		}
		empty := len(s.Tables)+len(s.Views)+len(s.MaterializedViews)+len(s.Sequences)+len(s.EnumTypes)+len(s.CompositeTypes)+len(s.Functions) == 0
		if !empty && n.DropBehavior != tree.DropCascade {
			return refuse(codeDependentObjects, "schema %q is not empty and CASCADE was not specified", name)
		}
		var names tree.TableNames
		for _, table := range s.Tables {
			names = append(names, tree.MakeTableNameWithSchema("", tree.Name(name), tree.Name(table.Name)))
		}
		if err := c.dropRelations(names, true, tree.DropCascade, relTable, true); err != nil {
			return err
		}
		delete(c.l.schemas, name)
		c.l.meta.Schemas = slices.DeleteFunc(c.l.meta.Schemas, func(x *metadatapb.SchemaMetadata) bool { return x == s })
	}
	return nil
}

// checkType refuses a type reference to a user-defined type the catalog
// does not hold.
func (c *Catalog) checkType(ref tree.ResolvableTypeReference) error {
	if _, ok := tree.GetStaticallyKnownType(ref); ok {
		return nil
	}
	switch r := ref.(type) {
	case *tree.ArrayTypeReference:
		return c.checkType(r.ElementType)
	case *tree.UnresolvedObjectName:
		_, _, _, err := c.userType(r)
		return err
	}
	return nil
}

// userType finds an enum or composite type.
func (c *Catalog) userType(name *tree.UnresolvedObjectName) (*metadatapb.SchemaMetadata, *metadatapb.EnumTypeMetadata, *metadatapb.CompositeTypeMetadata, error) {
	s := c.l.schemas[objectSchema(name)]
	if s != nil {
		for _, e := range s.EnumTypes {
			if e.Name == name.Object() {
				return s, e, nil, nil
			}
		}
		for _, t := range s.CompositeTypes {
			if t.Name == name.Object() {
				return s, nil, t, nil
			}
		}
	}
	return nil, nil, nil, refuse(codeUndefinedObject, "type %q does not exist", name.String())
}

func (c *Catalog) createTable(n *tree.CreateTable) error {
	s, err := c.schemaFor(schemaOf(&n.Table))
	if err != nil {
		return err
	}
	if n.Persistence.IsTemporary() {
		return refuse("XCEXF", "temporary tables are only supported experimentally")
	}
	name := string(n.Table.ObjectName)
	if relation(s, name) != relNone || typeTaken(s, name) {
		if n.IfNotExists {
			return nil
		}
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	columns := map[string]bool{}
	for _, def := range n.Defs {
		if d, ok := def.(*tree.ColumnTableDef); ok {
			if columns[string(d.Name)] {
				return refuse(codeDuplicateColumn, "duplicate column name: %q", d.Name)
			}
			columns[string(d.Name)] = true
			if err := c.checkType(d.Type); err != nil {
				return err
			}
			if err := c.checkSequences(d.DefaultExpr.Expr); err != nil {
				return err
			}
		}
	}
	// The new table's columns, for the foreign keys it declares.
	pending := &metadatapb.TableMetadata{Name: name}
	for _, def := range n.Defs {
		if d, ok := def.(*tree.ColumnTableDef); ok {
			typ, _ := sdlType(d.Type)
			pending.Columns = append(pending.Columns, &metadatapb.ColumnMetadata{Name: string(d.Name), Type: typ})
		}
	}
	c.pending = pending
	defer func() { c.pending = nil }()
	for _, def := range n.Defs {
		if err := c.checkTableDef(name, columns, def); err != nil {
			return err
		}
	}
	return c.l.load(n)
}

// checkTableDef refuses a column or constraint definition that names a
// column, table or key the catalog does not hold.
func (c *Catalog) checkTableDef(table string, columns map[string]bool, def tree.TableDef) error {
	switch d := def.(type) {
	case *tree.ColumnTableDef:
		for _, check := range d.CheckExprs {
			if err := checkColumns(columns, check.Expr); err != nil {
				return err
			}
		}
		if d.References.Table != nil {
			fk := &tree.ForeignKeyConstraintTableDef{Table: *d.References.Table, FromCols: tree.NameList{d.Name}}
			if d.References.Col != "" {
				fk.ToCols = tree.NameList{d.References.Col}
			}
			return c.checkForeignKey(table, columns, fk)
		}
	case *tree.CheckConstraintTableDef:
		return checkColumns(columns, d.Expr)
	case *tree.UniqueConstraintTableDef:
		return checkIndexColumns(columns, d.Columns)
	case *tree.IndexTableDef:
		return checkIndexColumns(columns, d.Columns)
	case *tree.ForeignKeyConstraintTableDef:
		return c.checkForeignKey(table, columns, d)
	}
	return nil
}

// checkSequences refuses a default that calls nextval on a sequence the
// catalog does not hold.
func (c *Catalog) checkSequences(expr tree.Expr) error {
	if expr == nil {
		return nil
	}
	for _, name := range usedSequences(tree.AsString(expr)) {
		schema, sequence, _ := strings.Cut(name, ".")
		if relation(c.l.schemas[schema], sequence) == relNone {
			return refuse(codeUndefinedTable, "relation %q does not exist", sequence)
		}
	}
	return nil
}

func checkIndexColumns(columns map[string]bool, elems tree.IndexElemList) error {
	for _, elem := range elems {
		if elem.Expr == nil && !columns[string(elem.Column)] {
			return refuse(codeUndefinedColumn, "column %q does not exist", elem.Column)
		}
		if elem.Expr != nil {
			if err := checkColumns(columns, elem.Expr); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkColumns refuses an expression that reads a column the table does not
// have.
func checkColumns(columns map[string]bool, expr tree.Expr) error {
	var missing string
	tree.WalkExprConst(&columnCollector{visit: func(column string) {
		if missing == "" && !columns[column] {
			missing = column
		}
	}}, expr)
	if missing != "" {
		return refuse(codeUndefinedColumn, "column %q does not exist", missing)
	}
	return nil
}

// checkForeignKey refuses a foreign key to a table or columns the catalog
// does not hold, or to columns no unique key of the referenced table covers.
func (c *Catalog) checkForeignKey(table string, columns map[string]bool, d *tree.ForeignKeyConstraintTableDef) error {
	for _, col := range d.FromCols {
		if !columns[string(col)] {
			return refuse(codeUndefinedColumn, "column %q does not exist", col)
		}
	}
	if schemaOf(&d.Table) == "public" && string(d.Table.ObjectName) == table && c.l.schemas["public"] != nil && schemaTable(c.l.schemas["public"], table) == nil {
		// A table's foreign key to itself, in its CREATE TABLE.
		return nil
	}
	referenced, err := c.table(&d.Table)
	if err != nil {
		return err
	}
	refColumns := map[string]bool{}
	for _, column := range referenced.Columns {
		refColumns[column.Name] = true
	}
	var keys [][]string
	for _, index := range referenced.Indexes {
		if index.Unique && !strings.Contains(index.Definition, " WHERE ") {
			keys = append(keys, index.Expressions)
		}
	}
	to := make([]string, len(d.ToCols))
	for i, col := range d.ToCols {
		if !refColumns[string(col)] {
			return refuse(codeUndefinedColumn, "column %q does not exist", col)
		}
		to[i] = string(col)
	}
	if len(to) == 0 {
		for _, index := range referenced.Indexes {
			if index.Primary {
				to = index.Expressions
			}
		}
		if len(to) == 1 && isImplicitRowIDName(referenced, to[0]) {
			return refuse(codeForeignKeyViolation, "there is no unique constraint matching given keys for referenced table %s", referenced.Name)
		}
	}
	if from := c.columnFamilies(table, d.FromCols); from != nil {
		for i, col := range to {
			ref := findColumn(referenced, col)
			if i < len(from) && from[i] != "" && ref != nil && family(ref.Type) != "" && from[i] != family(ref.Type) {
				return refuse("42804", "type of %q does not match foreign key %q.%q", d.FromCols[i], referenced.Name, col)
			}
		}
	}
	for _, key := range keys {
		if sameSet(key, to) {
			return nil
		}
	}
	return refuse(codeForeignKeyViolation, "there is no unique constraint matching given keys for referenced table %s", referenced.Name)
}

// columnFamilies returns the type families of a foreign key's columns, from
// the table CREATE TABLE or ALTER TABLE is declaring it on.
func (c *Catalog) columnFamilies(table string, cols tree.NameList) []string {
	t := c.pending
	if t == nil || t.Name != table {
		return nil
	}
	out := make([]string, len(cols))
	for i, col := range cols {
		if column := findColumn(t, string(col)); column != nil {
			out[i] = family(column.Type)
		}
	}
	return out
}

// family is the type family of a built-in type as a column states it.
func family(typ string) string {
	ref, err := parser.GetTypeFromValidSQLSyntax(typ)
	if err != nil {
		return ""
	}
	t, ok := tree.GetStaticallyKnownType(ref)
	if !ok {
		return ""
	}
	return t.Family().String()
}

func isImplicitRowIDName(table *metadatapb.TableMetadata, name string) bool {
	column := findColumn(table, name)
	return column != nil && isImplicitRowID(column)
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		if !slices.Contains(b, x) {
			return false
		}
	}
	return true
}

func (c *Catalog) createView(n *tree.CreateView) error {
	s, err := c.schemaFor(schemaOf(&n.Name))
	if err != nil {
		return err
	}
	name := string(n.Name.ObjectName)
	kind := relation(s, name)
	switch {
	case kind == relNone:
	case n.IfNotExists:
		return nil
	case n.Replace && kind == relView && !n.Materialized:
		old := findView(s, name)
		names := outputNames(n.AsSource)
		if len(n.ColumnNames) > 0 {
			names = nil
			for _, col := range n.ColumnNames {
				names = append(names, string(col))
			}
		}
		if err := keepsColumns(old.Columns, names); err != nil {
			return err
		}
		if names == nil {
			// A star: the columns are known once the loader expands it,
			// and Exec puts the view back if they do not keep the old ones.
			oldColumns := old.Columns
			if err := c.l.load(n); err != nil {
				return err
			}
			var now []string
			for _, column := range findView(s, name).Columns {
				now = append(now, column.Name)
			}
			if now == nil {
				return &Error{Message: fmt.Sprintf("the columns of view %q are not known", name)}
			}
			return keepsColumns(oldColumns, now)
		}
	case n.Replace:
		return refuse(codeWrongObjectType, "%q is not a view", name)
	default:
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	var missing error
	storedRelations(n.AsSource, func(tn *tree.TableName) {
		if missing != nil {
			return
		}
		if relation(c.l.schemas[schemaOf(tn)], string(tn.ObjectName)) == relNone {
			missing = refuse(codeUndefinedTable, "relation %q does not exist", tn.String())
		}
	})
	if missing != nil {
		return missing
	}
	if err := c.checkViewColumns(n.AsSource); err != nil {
		return err
	}
	return c.l.load(n)
}

// checkViewColumns refuses a query on one table that names a column the
// table does not have. A query of another shape is taken as written.
func (c *Catalog) checkViewColumns(query *tree.Select) error {
	sel, ok := query.Select.(*tree.SelectClause)
	if !ok || query.With != nil || len(sel.From.Tables) != 1 {
		return nil
	}
	aliased, ok := sel.From.Tables[0].(*tree.AliasedTableExpr)
	if !ok || len(aliased.As.Cols) > 0 {
		return nil
	}
	tn, ok := aliased.Expr.(*tree.TableName)
	if !ok {
		return nil
	}
	table := schemaTable(c.l.schemas[schemaOf(tn)], string(tn.ObjectName))
	if table == nil {
		return nil
	}
	columns := map[string]bool{}
	for _, column := range table.Columns {
		columns[column.Name] = true
	}
	for _, expr := range sel.Exprs {
		if err := checkColumns(columns, expr.Expr); err != nil {
			return err
		}
	}
	return nil
}

// keepsColumns refuses a replacement view whose columns drop or rename the
// old view's, as CockroachDB does. names is nil when they are not known.
func keepsColumns(old []*metadatapb.ColumnMetadata, names []string) error {
	if names == nil {
		return nil
	}
	if len(names) < len(old) {
		return refuse(codeInvalidTableDef, "cannot drop columns from view")
	}
	for i, column := range old {
		if names[i] != column.Name {
			return refuse(codeInvalidTableDef, "cannot change name of view column %q to %q", column.Name, names[i])
		}
	}
	return nil
}

func (c *Catalog) createSequence(n *tree.CreateSequence) error {
	s, err := c.schemaFor(schemaOf(&n.Name))
	if err != nil {
		return err
	}
	name := string(n.Name.ObjectName)
	if relation(s, name) != relNone {
		if n.IfNotExists {
			return nil
		}
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	if err := c.checkSequenceOptions(n.Options); err != nil {
		return err
	}
	if err := c.l.load(n); err != nil {
		return err
	}
	return checkBounds(findSequence(s, name))
}

// checkBounds refuses a sequence whose start is outside its bounds.
func checkBounds(sequence *metadatapb.SequenceMetadata) error {
	minimum, err1 := strconv.ParseInt(sequence.MinValue, 10, 64)
	maximum, err2 := strconv.ParseInt(sequence.MaxValue, 10, 64)
	start, err3 := strconv.ParseInt(sequence.Start, 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil
	}
	switch {
	case minimum > maximum:
		return refuse(codeInvalidParameter, "MINVALUE (%d) must be less than MAXVALUE (%d)", minimum, maximum)
	case start > maximum:
		return refuse(codeInvalidParameter, "START value (%d) cannot be greater than MAXVALUE (%d)", start, maximum)
	case start < minimum:
		return refuse(codeInvalidParameter, "START value (%d) cannot be less than MINVALUE (%d)", start, minimum)
	}
	return nil
}

func (c *Catalog) checkSequenceOptions(options tree.SequenceOptions) error {
	for _, option := range options {
		if option.Name == tree.SeqOptCycle {
			return refuse(codeFeatureNotSupport, "unimplemented: CYCLE option is not supported")
		}
		if option.Name != tree.SeqOptOwnedBy || option.ColumnItemVal == nil {
			continue
		}
		table, err := c.table(ptr(option.ColumnItemVal.TableName.ToTableName()))
		if err != nil {
			return err
		}
		if findColumn(table, string(option.ColumnItemVal.ColumnName)) == nil {
			return refuse(codeUndefinedColumn, "column %q does not exist", option.ColumnItemVal.ColumnName)
		}
	}
	return nil
}

func (c *Catalog) alterSequence(n *tree.AlterSequence) error {
	s := c.l.schemas[objectSchema(n.Name)]
	if relation(s, n.Name.Object()) != relSequence {
		if n.IfExists {
			return nil
		}
		return refuse(codeUndefinedTable, "relation %q does not exist", n.Name.String())
	}
	if err := c.checkSequenceOptions(n.Options); err != nil {
		return err
	}
	// RESTART moves the sequence's value, which the catalog does not hold.
	n.Options = slices.DeleteFunc(n.Options, func(o tree.SequenceOption) bool { return o.Name == tree.SeqOptRestart })
	if err := c.l.load(n); err != nil {
		return err
	}
	return checkBounds(findSequence(s, n.Name.Object()))
}

func (c *Catalog) createIndex(n *tree.CreateIndex) error {
	s := c.l.schemas[schemaOf(&n.Table)]
	name := string(n.Table.ObjectName)
	var columns []*metadatapb.ColumnMetadata
	var indexes []*metadatapb.IndexMetadata
	switch relation(s, name) {
	case relTable:
		table := schemaTable(s, name)
		columns, indexes = table.Columns, table.Indexes
	case relMaterializedView:
		view := findMaterializedView(s, name)
		columns, indexes = view.Columns, view.Indexes
	case relNone:
		return refuse(codeUndefinedTable, "relation %q does not exist", n.Table.String())
	default:
		return refuse(codeWrongObjectType, "%q is not a table or materialized view", name)
	}
	if n.Name != "" && slices.ContainsFunc(indexes, func(i *metadatapb.IndexMetadata) bool { return i.Name == string(n.Name) }) {
		if n.IfNotExists {
			return nil
		}
		return refuse(codeDuplicateRelation, "index with name %q already exists", n.Name)
	}
	if n.Name != "" && c.indexNamed(s, string(n.Name)) != nil {
		if n.IfNotExists {
			return nil
		}
		return refuse(codeDuplicateRelation, "relation %q already exists", n.Name)
	}
	names := map[string]bool{}
	for _, column := range columns {
		names[column.Name] = true
	}
	if err := checkIndexColumns(names, n.Columns); err != nil {
		return err
	}
	for _, col := range n.Storing {
		if !names[string(col)] {
			return refuse(codeUndefinedColumn, "column %q does not exist", col)
		}
	}
	if n.Predicate != nil {
		if err := checkColumns(names, n.Predicate); err != nil {
			return err
		}
	}
	return c.l.load(n)
}

// indexNamed finds an index of any relation of a schema.
func (c *Catalog) indexNamed(s *metadatapb.SchemaMetadata, name string) *metadatapb.IndexMetadata {
	for _, table := range s.Tables {
		for _, index := range table.Indexes {
			if index.Name == name {
				return index
			}
		}
	}
	for _, view := range s.MaterializedViews {
		for _, index := range view.Indexes {
			if index.Name == name {
				return index
			}
		}
	}
	return nil
}

func (c *Catalog) createType(n *tree.CreateType) error {
	s, err := c.schemaFor(n.TypeName.Schema())
	if err != nil {
		return err
	}
	name := n.TypeName.Object()
	if typeTaken(s, name) || relation(s, name) != relNone {
		if n.IfNotExists {
			return nil
		}
		if relation(s, name) != relNone {
			return refuse(codeDuplicateRelation, "relation %q already exists", name)
		}
		return refuse(codeDuplicateObject, "type %q already exists", name)
	}
	seen := map[tree.EnumValue]bool{}
	for _, label := range n.EnumLabels {
		if seen[label] {
			return refuse("42P17", "enum definition contains duplicate value %q", label)
		}
		seen[label] = true
	}
	for _, elem := range n.CompositeTypeList {
		if err := c.checkType(elem.Type); err != nil {
			return err
		}
		if _, ok := tree.GetStaticallyKnownType(elem.Type); !ok {
			return refuse(codeFeatureNotSupport, "composite types that reference user-defined types not yet supported")
		}
	}
	return c.l.load(n)
}

func (c *Catalog) alterType(n *tree.AlterType) error {
	s, enum, composite, err := c.userType(n.Type)
	if err != nil {
		return err
	}
	name := n.Type.Object()
	switch cmd := n.Cmd.(type) {
	case *tree.AlterTypeAddValue:
		if enum == nil {
			return refuse(codeWrongObjectType, "%q is not an enum", name)
		}
		if slices.Contains(enum.Values, string(cmd.NewVal)) {
			if cmd.IfNotExists {
				return nil
			}
			return refuse(codeDuplicateObject, "enum value %q already exists", cmd.NewVal)
		}
		at := len(enum.Values)
		if cmd.Placement != nil {
			at = slices.Index(enum.Values, string(cmd.Placement.ExistingVal))
			if at < 0 {
				return refuse(codeInvalidParameter, "%q is not an existing enum value", cmd.Placement.ExistingVal)
			}
			if !cmd.Placement.Before {
				at++
			}
		}
		enum.Values = slices.Insert(enum.Values, at, string(cmd.NewVal))
	case *tree.AlterTypeRenameValue:
		if enum == nil {
			return refuse(codeWrongObjectType, "%q is not an enum", name)
		}
		at := slices.Index(enum.Values, string(cmd.OldVal))
		if at < 0 {
			return refuse(codeInvalidParameter, "%q is not an existing enum value", cmd.OldVal)
		}
		if slices.Contains(enum.Values, string(cmd.NewVal)) {
			return refuse(codeDuplicateObject, "enum value %q already exists", cmd.NewVal)
		}
		enum.Values[at] = string(cmd.NewVal)
	case *tree.AlterTypeDropValue:
		if enum == nil {
			return refuse(codeWrongObjectType, "%q is not an enum", name)
		}
		at := slices.Index(enum.Values, string(cmd.Val))
		if at < 0 {
			return refuse(codeUndefinedObject, "enum value %q does not exist", cmd.Val)
		}
		if user := c.labelUser(typeName(s.Name, name), string(cmd.Val)); user != "" {
			return refuse(codeDependentObjects, "could not remove enum value %q as it is being used in a default expression of %q", cmd.Val, user)
		}
		enum.Values = slices.Delete(enum.Values, at, at+1)
	case *tree.AlterTypeRename:
		to := string(cmd.NewName)
		if typeTaken(s, to) || relation(s, to) != relNone {
			return refuse(codeDuplicateObject, "type %q already exists", to)
		}
		c.retype(s.Name, name, s.Name, to)
		if enum != nil {
			enum.Name = to
		} else {
			composite.Name = to
		}
	case *tree.AlterTypeSetSchema:
		target, err := c.schemaFor(string(cmd.Schema))
		if err != nil {
			return err
		}
		if typeTaken(target, name) || relation(target, name) != relNone {
			return refuse(codeDuplicateObject, "type %q already exists", name)
		}
		c.retype(s.Name, name, target.Name, name)
		if enum != nil {
			s.EnumTypes = slices.DeleteFunc(s.EnumTypes, func(e *metadatapb.EnumTypeMetadata) bool { return e == enum })
			target.EnumTypes = append(target.EnumTypes, enum)
		} else {
			s.CompositeTypes = slices.DeleteFunc(s.CompositeTypes, func(t *metadatapb.CompositeTypeMetadata) bool { return t == composite })
			target.CompositeTypes = append(target.CompositeTypes, composite)
		}
	default:
		return &Error{Message: fmt.Sprintf("ALTER TYPE %s is not supported", tree.AsString(n.Cmd))}
	}
	return nil
}

// typeName is a user-defined type as a column states it.
func typeName(schema, name string) string {
	return tree.NameString(schema) + "." + tree.NameString(name)
}

// retype renames a user-defined type in the columns that hold it.
func (c *Catalog) retype(fromSchema, from, toSchema, to string) {
	old, now := typeName(fromSchema, from), typeName(toSchema, to)
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, column := range table.Columns {
			switch column.Type {
			case old:
				column.Type = now
			case old + "[]":
				column.Type = now + "[]"
			}
		}
	})
}

// labelUser returns a table a column default of which holds an enum label.
func (c *Catalog) labelUser(typ, label string) string {
	user := ""
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, column := range table.Columns {
			if user == "" && column.Type == typ && column.Default != "" {
				expr, err := parser.ParseExpr(column.Default)
				if err == nil && readsLabel(expr, label) {
					user = table.Name
				}
			}
		}
	})
	return user
}

func readsLabel(expr tree.Expr, label string) bool {
	found := false
	_, _ = tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		if str, ok := e.(*tree.StrVal); ok && str.RawString() == label {
			found = true
		}
		return true, e, nil
	})
	return found
}

// typeUsers returns the tables whose columns hold a type.
func (c *Catalog) typeUsers(schema, name string) []string {
	typ := typeName(schema, name)
	var users []string
	c.eachTable(func(s *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, column := range table.Columns {
			if column.Type == typ || column.Type == typ+"[]" {
				users = append(users, table.Name)
				return
			}
		}
	})
	for _, s := range c.l.meta.Schemas {
		for _, composite := range s.CompositeTypes {
			for _, attr := range composite.Attributes {
				if attr.Type == typ {
					users = append(users, composite.Name)
				}
			}
		}
	}
	return users
}

func (c *Catalog) dropType(n *tree.DropType) error {
	type target struct {
		s    *metadatapb.SchemaMetadata
		name string
	}
	var targets []target
	for _, name := range n.Names {
		s, _, _, err := c.userType(name)
		if err != nil {
			if relation(c.l.schemas[objectSchema(name)], name.Object()) != relNone {
				// A relation's row type is one its relation needs.
				return refuse(codeDependentObjects, "cannot drop type %q because relation %q requires it", name.Object(), name.Object())
			}
			if n.IfExists {
				continue
			}
			return err
		}
		targets = append(targets, target{s, name.Object()})
	}
	for _, t := range targets {
		if users := c.typeUsers(t.s.Name, t.name); len(users) > 0 {
			return refuse(codeDependentObjects, "cannot drop type %q because other objects (%s) still depend on it", t.name, strings.Join(users, ", "))
		}
	}
	for _, t := range targets {
		t.s.EnumTypes = slices.DeleteFunc(t.s.EnumTypes, func(e *metadatapb.EnumTypeMetadata) bool { return e.Name == t.name })
		t.s.CompositeTypes = slices.DeleteFunc(t.s.CompositeTypes, func(x *metadatapb.CompositeTypeMetadata) bool { return x.Name == t.name })
	}
	return nil
}

// eachTable visits every table of the catalog.
func (c *Catalog) eachTable(visit func(*metadatapb.SchemaMetadata, *metadatapb.TableMetadata)) {
	for _, s := range c.l.meta.Schemas {
		for _, table := range s.Tables {
			visit(s, table)
		}
	}
}

// A dependent is a view or materialized view that reads a relation, or a
// column of one.
type dependent struct {
	schema *metadatapb.SchemaMetadata
	name   string
	kind   relKind
}

// dependents returns the views and materialized views that read a relation,
// or, when column is not empty, that column of it.
func (c *Catalog) dependents(schema, name, column string) []dependent {
	var out []dependent
	reads := func(deps []*metadatapb.DependencyColumn, definition string) bool {
		for _, dep := range deps {
			if dep.Schema == schema && dep.Table == name {
				return column == "" || readsColumn(definition, column)
			}
		}
		return false
	}
	for _, s := range c.l.meta.Schemas {
		for _, v := range s.Views {
			if !(s.Name == schema && v.Name == name) && reads(v.DependencyColumns, v.Definition) {
				out = append(out, dependent{s, v.Name, relView})
			}
		}
		for _, v := range s.MaterializedViews {
			if !(s.Name == schema && v.Name == name) && reads(v.DependencyColumns, v.Definition) {
				out = append(out, dependent{s, v.Name, relMaterializedView})
			}
		}
	}
	return out
}

// readsColumn reports whether a view's query names a column, or holds a star
// that reads every column.
func readsColumn(definition, column string) bool {
	stmt, err := parser.ParseOne(definition)
	if err != nil {
		return true
	}
	found := false
	tree.WalkStmt(&nameVisitor{visit: func(name *tree.UnresolvedName) {
		if name.Star || name.Parts[0] == column {
			found = true
		}
	}}, stmt.AST)
	return found
}

type nameVisitor struct{ visit func(*tree.UnresolvedName) }

func (v *nameVisitor) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	switch e := expr.(type) {
	case *tree.UnresolvedName:
		v.visit(e)
	case tree.UnqualifiedStar:
		v.visit(&tree.UnresolvedName{Star: true, NumParts: 1})
	}
	return true, expr
}

func (*nameVisitor) VisitPost(expr tree.Expr) tree.Expr { return expr }

// referencing returns the foreign keys of other tables that reference a
// table, or, when column is not empty, that column of it.
func (c *Catalog) referencing(schema, name, column string) []*metadatapb.ForeignKeyMetadata {
	var out []*metadatapb.ForeignKeyMetadata
	c.eachTable(func(s *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, fk := range table.ForeignKeys {
			if fk.ReferencedSchema == schema && fk.ReferencedTable == name && (column == "" || slices.Contains(fk.ReferencedColumns, column)) {
				if column == "" && s.Name == schema && table.Name == name {
					continue
				}
				out = append(out, fk)
			}
		}
	})
	return out
}

// dropForeignKeys removes foreign keys from the tables that hold them.
func (c *Catalog) dropForeignKeys(fks []*metadatapb.ForeignKeyMetadata) {
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		table.ForeignKeys = slices.DeleteFunc(table.ForeignKeys, func(fk *metadatapb.ForeignKeyMetadata) bool { return slices.Contains(fks, fk) })
	})
}

// sequenceUsers returns the columns whose default reads a sequence.
func (c *Catalog) sequenceUsers(schema, name string) []*metadatapb.ColumnMetadata {
	var out []*metadatapb.ColumnMetadata
	c.eachTable(func(s *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, column := range table.Columns {
			for _, seq := range usedSequences(column.Default) {
				if seq == schema+"."+name {
					out = append(out, column)
				}
			}
		}
	})
	return out
}

// dropRelations drops tables, views, materialized views or sequences, with
// what CASCADE takes along. inSchema drops them for DROP SCHEMA, where
// everything of the schema goes.
func (c *Catalog) dropRelations(names tree.TableNames, ifExists bool, behavior tree.DropBehavior, kind relKind, inSchema bool) error {
	type target struct {
		s    *metadatapb.SchemaMetadata
		name string
		kind relKind
	}
	var targets []target
	for i := range names {
		tn := &names[i]
		s := c.l.schemas[schemaOf(tn)]
		got := relation(s, string(tn.ObjectName))
		if got == relNone {
			if ifExists {
				continue
			}
			return refuse(codeUndefinedTable, "relation %q does not exist", tn.String())
		}
		if got != kind {
			return refuse(codeWrongObjectType, "%q is not a %s", tn.ObjectName, kind)
		}
		targets = append(targets, target{s, string(tn.ObjectName), got})
	}
	dropping := func(schema, name string) bool {
		return slices.ContainsFunc(targets, func(t target) bool { return t.s.Name == schema && t.name == name })
	}
	cascade := behavior == tree.DropCascade
	for i := 0; i < len(targets); i++ {
		t := targets[i]
		for _, d := range c.dependents(t.s.Name, t.name, "") {
			if dropping(d.schema.Name, d.name) {
				continue
			}
			if !cascade {
				return refuse(codeDependentObjects, "cannot drop relation %q because view %q depends on it", t.name, d.name)
			}
			targets = append(targets, target{d.schema, d.name, d.kind})
		}
		switch t.kind {
		case relTable:
			for _, fk := range c.referencing(t.s.Name, t.name, "") {
				if owner := c.foreignKeyTable(fk); owner != nil && dropping(owner.schema, owner.table) {
					continue
				}
				if !cascade && !inSchema {
					return refuse(codeDependentObjects, "%q is referenced by foreign key %q", t.name, fk.Name)
				}
			}
		case relSequence:
			if users := c.sequenceUsers(t.s.Name, t.name); len(users) > 0 {
				if !cascade {
					return refuse(codeDependentObjects, "cannot drop sequence %s because other objects depend on it", t.name)
				}
				for _, column := range users {
					column.Default = ""
				}
			}
		}
	}
	for _, t := range targets {
		switch t.kind {
		case relTable:
			c.dropForeignKeys(c.referencing(t.s.Name, t.name, ""))
			t.s.Tables = slices.DeleteFunc(t.s.Tables, func(x *metadatapb.TableMetadata) bool { return x.Name == t.name })
			// A sequence a column of the table owns goes with it.
			t.s.Sequences = slices.DeleteFunc(t.s.Sequences, func(x *metadatapb.SequenceMetadata) bool { return x.OwnerTable == t.name })
		case relView:
			t.s.Views = slices.DeleteFunc(t.s.Views, func(x *metadatapb.ViewMetadata) bool { return x.Name == t.name })
		case relMaterializedView:
			t.s.MaterializedViews = slices.DeleteFunc(t.s.MaterializedViews, func(x *metadatapb.MaterializedViewMetadata) bool { return x.Name == t.name })
		case relSequence:
			t.s.Sequences = slices.DeleteFunc(t.s.Sequences, func(x *metadatapb.SequenceMetadata) bool { return x.Name == t.name })
		}
	}
	return nil
}

type tableRef struct{ schema, table string }

// foreignKeyTable returns the table that holds a foreign key.
func (c *Catalog) foreignKeyTable(fk *metadatapb.ForeignKeyMetadata) *tableRef {
	var out *tableRef
	c.eachTable(func(s *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		if slices.Contains(table.ForeignKeys, fk) {
			out = &tableRef{s.Name, table.Name}
		}
	})
	return out
}

func (c *Catalog) renameRelation(n *tree.RenameTable) error {
	kind := relTable
	switch {
	case n.IsMaterialized:
		kind = relMaterializedView
	case n.IsView:
		kind = relView
	case n.IsSequence:
		kind = relSequence
	}
	s := c.l.schemas[objectSchema(n.Name)]
	name := n.Name.Object()
	got := relation(s, name)
	if got == relNone {
		if n.IfExists {
			return nil
		}
		return refuse(codeUndefinedTable, "relation %q does not exist", n.Name.String())
	}
	if got != kind && !(kind == relTable && got != relNone) {
		return refuse(codeWrongObjectType, "%q is not a %s", name, kind)
	}
	if deps := c.dependents(s.Name, name, ""); len(deps) > 0 {
		return refuse(codeDependentObjects, "cannot rename relation %q because view %q depends on it", name, deps[0].name)
	}
	// A new name without a schema keeps the relation in its own.
	targetSchema := n.NewName.Schema()
	if targetSchema == "" {
		targetSchema = s.Name
	}
	target, err := c.schemaFor(targetSchema)
	if err != nil {
		return err
	}
	to := n.NewName.Object()
	if relation(target, to) != relNone || typeTaken(target, to) {
		return refuse(codeDuplicateRelation, "relation %q already exists", to)
	}
	c.moveRelation(s, name, got, target, to)
	return nil
}

func (c *Catalog) setRelationSchema(n *tree.AlterTableSetSchema) error {
	s := c.l.schemas[objectSchema(n.Name)]
	name := n.Name.Object()
	got := relation(s, name)
	if got == relNone {
		if n.IfExists {
			return nil
		}
		return refuse(codeUndefinedTable, "relation %q does not exist", n.Name.String())
	}
	target, err := c.schemaFor(string(n.Schema))
	if err != nil {
		return err
	}
	if relation(target, name) != relNone {
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	if deps := c.dependents(s.Name, name, ""); len(deps) > 0 {
		return refuse(codeDependentObjects, "cannot set schema on relation %q because view %q depends on it", name, deps[0].name)
	}
	c.moveRelation(s, name, got, target, name)
	return nil
}

// moveRelation renames a relation, or moves it to another schema, with the
// foreign keys, indexes and sequence owners that name it.
func (c *Catalog) moveRelation(s *metadatapb.SchemaMetadata, name string, kind relKind, target *metadatapb.SchemaMetadata, to string) {
	switch kind {
	case relTable:
		table := schemaTable(s, name)
		s.Tables = slices.DeleteFunc(s.Tables, func(x *metadatapb.TableMetadata) bool { return x == table })
		table.Name = to
		target.Tables = append(target.Tables, table)
		reindex(table.Indexes, target.Name, to, "", "")
		c.eachTable(func(_ *metadatapb.SchemaMetadata, other *metadatapb.TableMetadata) {
			for _, fk := range other.ForeignKeys {
				if fk.ReferencedSchema == s.Name && fk.ReferencedTable == name {
					fk.ReferencedSchema, fk.ReferencedTable = target.Name, to
				}
			}
		})
		for _, sequence := range s.Sequences {
			if sequence.OwnerTable == name {
				sequence.OwnerTable = to
			}
		}
	case relView:
		view := findView(s, name)
		s.Views = slices.DeleteFunc(s.Views, func(x *metadatapb.ViewMetadata) bool { return x == view })
		view.Name = to
		target.Views = append(target.Views, view)
	case relMaterializedView:
		view := findMaterializedView(s, name)
		s.MaterializedViews = slices.DeleteFunc(s.MaterializedViews, func(x *metadatapb.MaterializedViewMetadata) bool { return x == view })
		view.Name = to
		target.MaterializedViews = append(target.MaterializedViews, view)
		reindex(view.Indexes, target.Name, to, "", "")
	case relSequence:
		sequence := findSequence(s, name)
		s.Sequences = slices.DeleteFunc(s.Sequences, func(x *metadatapb.SequenceMetadata) bool { return x == sequence })
		sequence.Name = to
		target.Sequences = append(target.Sequences, sequence)
		// CockroachDB holds a default's sequence by ID, so the default
		// follows it.
		for _, column := range c.sequenceUsers(s.Name, name) {
			column.Default = renameSequenceIn(column.Default, s.Name+"."+name, target.Name+"."+to)
		}
	}
}

// reindex writes indexes' definitions again for their table's name, and for
// a renamed column when from is not empty.
func reindex(indexes []*metadatapb.IndexMetadata, schema, table, from, to string) {
	for _, index := range indexes {
		stmt, err := parser.ParseOne(index.Definition)
		if err != nil {
			continue
		}
		n, ok := stmt.AST.(*tree.CreateIndex)
		if !ok {
			continue
		}
		n.Table = tree.MakeTableNameWithSchema("", tree.Name(schema), tree.Name(table))
		n.Table.ExplicitCatalog = false
		if from != "" {
			for i := range n.Columns {
				if string(n.Columns[i].Column) == from {
					n.Columns[i].Column = tree.Name(to)
				}
				if n.Columns[i].Expr != nil {
					n.Columns[i].Expr = renameColumnIn(n.Columns[i].Expr, from, to)
				}
			}
			for i := range n.Storing {
				if string(n.Storing[i]) == from {
					n.Storing[i] = tree.Name(to)
				}
			}
			if n.Predicate != nil {
				n.Predicate = renameColumnIn(n.Predicate, from, to)
			}
			for i, expr := range index.Expressions {
				if expr == from {
					index.Expressions[i] = to
				}
			}
		}
		index.Name = string(n.Name)
		index.Definition = tree.AsString(n)
	}
}

// renameSequenceIn returns a default with the sequence its nextval, currval
// or setval names renamed.
func renameSequenceIn(def, from, to string) string {
	expr, err := parser.ParseExpr(def)
	if err != nil {
		return def
	}
	out, err := tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		call, ok := e.(*tree.FuncExpr)
		if !ok || len(call.Exprs) == 0 {
			return true, e, nil
		}
		name, ok := call.Func.FunctionReference.(*tree.UnresolvedName)
		if !ok || !slices.Contains([]string{"nextval", "currval", "setval"}, strings.ToLower(name.Parts[0])) {
			return true, e, nil
		}
		arg := unparen(call.Exprs[0])
		if cast, ok := arg.(*tree.CastExpr); ok {
			arg = unparen(cast.Expr)
		}
		str, ok := arg.(*tree.StrVal)
		if !ok {
			return true, e, nil
		}
		if sequence, err := parser.ParseQualifiedTableName(str.RawString()); err == nil && schemaOf(sequence)+"."+string(sequence.ObjectName) == from {
			renamed := *call
			renamed.Exprs = append(tree.Exprs{tree.NewStrVal(to)}, call.Exprs[1:]...)
			return false, &renamed, nil
		}
		return true, e, nil
	})
	if err != nil {
		return def
	}
	return tree.AsString(out)
}

// renameColumnIn returns an expression with a column renamed.
func renameColumnIn(expr tree.Expr, from, to string) tree.Expr {
	out, err := tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		if name, ok := e.(*tree.UnresolvedName); ok && name.NumParts == 1 && !name.Star && name.Parts[0] == from {
			return false, &tree.UnresolvedName{NumParts: 1, Parts: tree.NameParts{to}}, nil
		}
		return true, e, nil
	})
	if err != nil {
		return expr
	}
	return out
}

func (c *Catalog) alterTable(n *tree.AlterTable) error {
	tn := n.Table.ToTableName()
	table, err := c.table(&tn)
	if err != nil {
		if n.IfExists && err.(*Error).Code == codeUndefinedTable {
			return nil
		}
		return err
	}
	schema := schemaOf(&tn)
	c.pending = table
	defer func() { c.pending = nil }()
	for _, cmd := range n.Cmds {
		if err := c.alterTableCmd(schema, table, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (c *Catalog) alterTableCmd(schema string, table *metadatapb.TableMetadata, cmd tree.AlterTableCmd) error {
	columns := map[string]bool{}
	for _, column := range table.Columns {
		columns[column.Name] = true
	}
	column := func(name tree.Name) (*metadatapb.ColumnMetadata, error) {
		if col := findColumn(table, string(name)); col != nil {
			return col, nil
		}
		return nil, refuse(codeUndefinedColumn, "column %q does not exist", name)
	}
	switch cmd := cmd.(type) {
	case *tree.AlterTableAddColumn:
		def := cmd.ColumnDef
		if columns[string(def.Name)] {
			if cmd.IfNotExists {
				return nil
			}
			return refuse(codeDuplicateColumn, "column %q of relation %q already exists", def.Name, table.Name)
		}
		if err := c.checkType(def.Type); err != nil {
			return err
		}
		if err := c.checkSequences(def.DefaultExpr.Expr); err != nil {
			return err
		}
		columns[string(def.Name)] = true
		if err := c.checkTableDef(table.Name, columns, def); err != nil {
			return err
		}
		return c.l.addColumn(schema, table, def)
	case *tree.AlterTableAddConstraint:
		name := constraintDefName(cmd.ConstraintDef)
		if _, unique := cmd.ConstraintDef.(*tree.UniqueConstraintTableDef); unique && name != "" && slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool { return i.Name == name }) {
			return refuse(codeDuplicateRelation, "index with name %q already exists", name)
		}
		if name != "" && constraintNamed(table, name) {
			return refuse(codeDuplicateObject, "duplicate constraint name: %q", name)
		}
		if err := c.checkTableDef(table.Name, columns, cmd.ConstraintDef); err != nil {
			return err
		}
		if u, ok := cmd.ConstraintDef.(*tree.UniqueConstraintTableDef); ok && u.PrimaryKey {
			if _, primary, _ := tableShape(table); primary != nil {
				return refuse("42611", "multiple primary keys for table %q are not allowed", table.Name)
			}
			if err := nullableKey(table, u.Columns); err != nil {
				return err
			}
			return c.alterPrimaryKey(schema, table, string(u.Name), indexOf(&u.IndexTableDef))
		}
		// The catalog holds no validation state: NOT VALID is one CockroachDB
		// takes.
		cmd.ValidationBehavior = tree.ValidationDefault
		return c.l.addTableDef(schema, table, cmd.ConstraintDef)
	case *tree.AlterTableAlterPrimaryKey:
		if err := checkIndexColumns(columns, cmd.Columns); err != nil {
			return err
		}
		if err := nullableKey(table, cmd.Columns); err != nil {
			return err
		}
		return c.alterPrimaryKey(schema, table, string(cmd.Name), &tree.CreateIndex{Columns: cmd.Columns, Sharded: cmd.Sharded, StorageParams: cmd.StorageParams})
	case *tree.AlterTableValidateConstraint:
		if !constraintNamed(table, string(cmd.Constraint)) {
			return refuse(codeUndefinedObject, "constraint %q of relation %q does not exist", cmd.Constraint, table.Name)
		}
	case *tree.AlterTableSetDefault:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		if err := c.checkSequences(cmd.Default); err != nil {
			return err
		}
		col.Default = ""
		if cmd.Default != nil {
			col.Default = sdlExpr(cmd.Default)
		}
	case *tree.AlterTableSetNotNull:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		col.Nullable = false
	case *tree.AlterTableDropNotNull:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		if primaryKeyHolds(table, col.Name) {
			return refuse(codeInvalidTableDef, "column %q is in a primary index", col.Name)
		}
		col.Nullable = true
	case *tree.AlterTableSetVisible:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		col.IsInvisible = !cmd.Visible
	case *tree.AlterTableAlterColumnType:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		if err := c.checkType(cmd.ToType); err != nil {
			return err
		}
		if deps := c.dependents(schema, table.Name, col.Name); len(deps) > 0 {
			return refuse(codeDependentObjects, "cannot alter type of column %q because view %q depends on it", col.Name, deps[0].name)
		}
		typ, collation := sdlType(cmd.ToType)
		col.Type, col.Collation = typ, collation
		if cmd.Collation != "" {
			col.Collation = cmd.Collation
		}
	case *tree.AlterTableRenameColumn:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		if deps := c.dependents(schema, table.Name, col.Name); len(deps) > 0 {
			return refuse(codeDependentObjects, "cannot rename column %q because view %q depends on it", col.Name, deps[0].name)
		}
		if columns[string(cmd.NewName)] {
			return refuse(codeDuplicateColumn, "column %q of relation %q already exists", cmd.NewName, table.Name)
		}
		c.renameColumn(schema, table, col.Name, string(cmd.NewName))
	case *tree.AlterTableDropColumn:
		col := findColumn(table, string(cmd.Column))
		if col == nil {
			if cmd.IfExists {
				return nil
			}
			return refuse(codeUndefinedColumn, "column %q does not exist", cmd.Column)
		}
		return c.dropColumn(schema, table, col.Name, cmd.DropBehavior)
	case *tree.AlterTableRenameConstraint:
		from, to := string(cmd.Constraint), string(cmd.NewName)
		if !constraintNamed(table, from) {
			return refuse(codeUndefinedObject, "constraint %q of relation %q does not exist", from, table.Name)
		}
		if constraintNamed(table, to) {
			return refuse(codeDuplicateObject, "duplicate constraint name: %q", to)
		}
		for _, check := range table.CheckConstraints {
			if check.Name == from {
				check.Name = to
			}
		}
		for _, fk := range table.ForeignKeys {
			if fk.Name == from {
				fk.Name = to
			}
		}
		for _, index := range table.Indexes {
			if index.Name == from {
				index.Name = to
				reindexNamed(index, to)
			}
		}
	case *tree.AlterTableDropConstraint:
		name := string(cmd.Constraint)
		if !constraintNamed(table, name) {
			if cmd.IfExists {
				return nil
			}
			return refuse(codeUndefinedObject, "constraint %q of relation %q does not exist", name, table.Name)
		}
		for _, index := range table.Indexes {
			if index.Name == name && index.Primary {
				return refuse(codeFeatureNotSupport, "cannot drop primary key constraint %q", name)
			}
			if index.Name == name {
				return refuse(codeFeatureNotSupport, "cannot drop UNIQUE constraint %q using ALTER TABLE DROP CONSTRAINT, use DROP INDEX CASCADE instead", name)
			}
		}
		table.CheckConstraints = slices.DeleteFunc(table.CheckConstraints, func(x *metadatapb.CheckConstraintMetadata) bool { return x.Name == name })
		table.ForeignKeys = slices.DeleteFunc(table.ForeignKeys, func(x *metadatapb.ForeignKeyMetadata) bool { return x.Name == name })
		table.Indexes = slices.DeleteFunc(table.Indexes, func(x *metadatapb.IndexMetadata) bool { return x.Name == name })
	default:
		return &Error{Message: fmt.Sprintf("ALTER TABLE %s is not supported", tree.AsString(cmd))}
	}
	return nil
}

// nullableKey refuses a primary key on a nullable column.
func nullableKey(table *metadatapb.TableMetadata, elems tree.IndexElemList) error {
	for _, elem := range elems {
		if column := findColumn(table, string(elem.Column)); column != nil && column.Nullable {
			return refuse("42P15", "cannot use nullable column %q in primary key", column.Name)
		}
	}
	return nil
}

func constraintDefName(def tree.ConstraintTableDef) string {
	switch d := def.(type) {
	case *tree.CheckConstraintTableDef:
		return string(d.Name)
	case *tree.UniqueConstraintTableDef:
		return string(d.Name)
	case *tree.ForeignKeyConstraintTableDef:
		return string(d.Name)
	}
	return ""
}

// constraintNamed reports whether a table has a constraint, or an index that
// is one, by name.
func constraintNamed(table *metadatapb.TableMetadata, name string) bool {
	return constraintTaken(table, name) || slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool { return i.Name == name && (i.Unique || i.Primary) })
}

func primaryKeyHolds(table *metadatapb.TableMetadata, column string) bool {
	for _, index := range table.Indexes {
		if index.Primary && slices.Contains(index.Expressions, column) {
			return true
		}
	}
	return false
}

// reindexNamed writes an index's definition again under its new name.
func reindexNamed(index *metadatapb.IndexMetadata, name string) {
	stmt, err := parser.ParseOne(index.Definition)
	if err != nil {
		return
	}
	if n, ok := stmt.AST.(*tree.CreateIndex); ok {
		n.Name = tree.Name(name)
		index.Definition = tree.AsString(n)
	}
}

// alterPrimaryKey replaces a table's primary key. CockroachDB keeps the old
// key as a unique index, unless it was the hidden rowid, which it drops with
// its column.
func (c *Catalog) alterPrimaryKey(schema string, table *metadatapb.TableMetadata, name string, key *tree.CreateIndex) error {
	var old *metadatapb.IndexMetadata
	for _, index := range table.Indexes {
		if index.Primary {
			old = index
		}
	}
	table.Indexes = slices.DeleteFunc(table.Indexes, func(x *metadatapb.IndexMetadata) bool { return x == old })
	if old != nil {
		if len(old.Expressions) == 1 && isImplicitRowIDName(table, old.Expressions[0]) {
			table.Columns = slices.DeleteFunc(table.Columns, func(x *metadatapb.ColumnMetadata) bool { return x.Name == old.Expressions[0] })
		} else {
			stmt, err := parser.ParseOne(old.Definition)
			if err == nil {
				if n, ok := stmt.AST.(*tree.CreateIndex); ok {
					n.Name = tree.Name(table.Name + "_" + strings.Join(old.Expressions, "_") + "_key")
					table.Indexes = append(table.Indexes, sdlIndex(schema, table.Name, table.Indexes, n))
				}
			}
		}
	}
	if name == "" && old != nil {
		name = old.Name
	}
	addPrimaryKey(schema, table, name, key)
	return nil
}

// renameColumn renames a column, in the table's indexes, checks and foreign
// keys and in those of other tables that reference it.
func (c *Catalog) renameColumn(schema string, table *metadatapb.TableMetadata, from, to string) {
	findColumn(table, from).Name = to
	reindex(table.Indexes, schema, table.Name, from, to)
	for _, column := range table.Columns {
		if column.Generation == nil {
			continue
		}
		if expr, err := parser.ParseExpr(column.Generation.Expression); err == nil {
			column.Generation.Expression = tree.AsString(renameColumnIn(expr, from, to))
		}
	}
	for _, check := range table.CheckConstraints {
		expr, err := parser.ParseExpr(check.Expression)
		if err == nil {
			check.Expression = "(" + tree.AsString(renameColumnIn(expr, from, to)) + ")"
		}
	}
	for _, fk := range table.ForeignKeys {
		for i, col := range fk.Columns {
			if col == from {
				fk.Columns[i] = to
			}
		}
	}
	for _, fk := range c.referencing(schema, table.Name, "") {
		for i, col := range fk.ReferencedColumns {
			if col == from {
				fk.ReferencedColumns[i] = to
			}
		}
	}
	for _, fk := range table.ForeignKeys {
		if fk.ReferencedSchema == schema && fk.ReferencedTable == table.Name {
			for i, col := range fk.ReferencedColumns {
				if col == from {
					fk.ReferencedColumns[i] = to
				}
			}
		}
	}
	for _, s := range c.l.meta.Schemas {
		if s.Name != schema {
			continue
		}
		for _, sequence := range s.Sequences {
			if sequence.OwnerTable == table.Name && sequence.OwnerColumn == from {
				sequence.OwnerColumn = to
			}
		}
	}
}

// dropColumn drops a column, with the indexes and checks that read it, and,
// with CASCADE, the views that do.
func (c *Catalog) dropColumn(schema string, table *metadatapb.TableMetadata, name string, behavior tree.DropBehavior) error {
	if primaryKeyHolds(table, name) {
		return refuse("42P10", "column %q is referenced by the primary key", name)
	}
	cascade := behavior == tree.DropCascade
	for _, column := range table.Columns {
		if column.Generation == nil || column.Name == name {
			continue
		}
		if expr, err := parser.ParseExpr(column.Generation.Expression); err == nil && exprReads(expr, name) {
			if cascade {
				// v25.2.3's declarative schema changer fails this.
				return refuse("XX000", "internal error: building declarative schema changer plan for ALTER TABLE: unable to make progress")
			}
			return refuse("42P10", "column %q is referenced by computed column %q", name, column.Name)
		}
	}
	deps := c.dependents(schema, table.Name, name)
	if len(deps) > 0 && !cascade {
		return refuse(codeDependentObjects, "cannot drop column %q because view %q depends on it", name, deps[0].name)
	}
	if fks := c.referencing(schema, table.Name, name); len(fks) > 0 {
		if !cascade {
			return refuse(codeDependentObjects, "column %q is referenced by foreign key %q", name, fks[0].Name)
		}
		c.dropForeignKeys(fks)
	}
	for _, d := range deps {
		names := tree.TableNames{tree.MakeTableNameWithSchema("", tree.Name(d.schema.Name), tree.Name(d.name))}
		if err := c.dropRelations(names, true, tree.DropCascade, d.kind, false); err != nil {
			return err
		}
	}
	table.Columns = slices.DeleteFunc(table.Columns, func(x *metadatapb.ColumnMetadata) bool { return x.Name == name })
	for i, column := range table.Columns {
		column.Position = int32(i + 1)
	}
	table.Indexes = slices.DeleteFunc(table.Indexes, func(index *metadatapb.IndexMetadata) bool {
		return !index.Primary && indexReadsColumn(index, name)
	})
	table.CheckConstraints = slices.DeleteFunc(table.CheckConstraints, func(check *metadatapb.CheckConstraintMetadata) bool {
		expr, err := parser.ParseExpr(check.Expression)
		return err == nil && checkColumns(map[string]bool{}, expr) != nil && exprReads(expr, name)
	})
	table.ForeignKeys = slices.DeleteFunc(table.ForeignKeys, func(fk *metadatapb.ForeignKeyMetadata) bool { return slices.Contains(fk.Columns, name) })
	return nil
}

// indexReads reports whether an index keys, stores or filters on a column.
func indexReadsColumn(index *metadatapb.IndexMetadata, column string) bool {
	stmt, err := parser.ParseOne(index.Definition)
	if err != nil {
		return slices.Contains(index.Expressions, column)
	}
	n, ok := stmt.AST.(*tree.CreateIndex)
	if !ok {
		return false
	}
	for _, elem := range n.Columns {
		if string(elem.Column) == column || (elem.Expr != nil && exprReads(elem.Expr, column)) {
			return true
		}
	}
	if slices.Contains(n.Storing, tree.Name(column)) {
		return true
	}
	return n.Predicate != nil && exprReads(n.Predicate, column)
}

func exprReads(expr tree.Expr, column string) bool {
	found := false
	tree.WalkExprConst(&columnCollector{visit: func(name string) {
		if name == column {
			found = true
		}
	}}, expr)
	return found
}

// index finds an index by table@name, or by name in its schema.
func (c *Catalog) index(n *tree.TableIndexName, ifExists bool) (*metadatapb.IndexMetadata, *[]*metadatapb.IndexMetadata, error) {
	schema := schemaOf(&n.Table)
	s := c.l.schemas[schema]
	name := string(n.Index)
	if s != nil {
		var lists []*[]*metadatapb.IndexMetadata
		if n.Table.ObjectName != "" {
			if l := c.l.relationIndexes(schema, string(n.Table.ObjectName)); l != nil {
				lists = append(lists, l)
			} else if !ifExists {
				return nil, nil, refuse(codeUndefinedTable, "relation %q does not exist", n.Table.String())
			}
		} else {
			for _, table := range s.Tables {
				lists = append(lists, &table.Indexes)
			}
			for _, view := range s.MaterializedViews {
				lists = append(lists, &view.Indexes)
			}
		}
		for _, list := range lists {
			for _, index := range *list {
				if index.Name == name {
					return index, list, nil
				}
			}
		}
	}
	return nil, nil, refuse(codeUndefinedObject, "index %q does not exist", name)
}

func (c *Catalog) renameIndex(n *tree.RenameIndex) error {
	index, list, err := c.index(n.Index, n.IfExists)
	if err != nil {
		if n.IfExists {
			return nil
		}
		return err
	}
	to := string(n.NewName)
	s := c.l.schemas[schemaOf(&n.Index.Table)]
	if slices.ContainsFunc(*list, func(i *metadatapb.IndexMetadata) bool { return i.Name == to }) || relation(s, to) != relNone {
		return refuse(codeDuplicateRelation, "relation %q already exists", to)
	}
	index.Name = to
	reindexNamed(index, to)
	return nil
}

func (c *Catalog) alterIndexVisible(n *tree.AlterIndexVisible) error {
	index, _, err := c.index(&n.Index, n.IfExists)
	if err != nil {
		if n.IfExists {
			return nil
		}
		return err
	}
	stmt, perr := parser.ParseOne(index.Definition)
	if perr != nil {
		return nil
	}
	if ci, ok := stmt.AST.(*tree.CreateIndex); ok {
		ci.Invisibility = n.Invisibility
		index.Definition = tree.AsString(ci)
	}
	return nil
}

func (c *Catalog) dropIndex(n *tree.DropIndex) error {
	for _, name := range n.IndexList {
		index, list, err := c.index(name, n.IfExists)
		if err != nil {
			if n.IfExists {
				continue
			}
			return err
		}
		if index.Primary {
			return refuse(codeFeatureNotSupport, "index %q is in use as a primary key constraint", index.Name)
		}
		if index.IsConstraint && n.DropBehavior != tree.DropCascade {
			return refuse(codeDependentObjects, "index %q is in use as unique constraint", index.Name)
		}
		*list = slices.DeleteFunc(*list, func(x *metadatapb.IndexMetadata) bool { return x == index })
		// The foreign keys that relied on the index's key need CASCADE,
		// which takes them.
		var orphans []*metadatapb.ForeignKeyMetadata
		c.eachTable(func(s *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
			if &table.Indexes != list {
				return
			}
			for _, fk := range c.referencing(s.Name, table.Name, "") {
				if !slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool {
					return i.Unique && sameSet(i.Expressions, fk.ReferencedColumns)
				}) {
					orphans = append(orphans, fk)
				}
			}
		})
		if len(orphans) > 0 && n.DropBehavior != tree.DropCascade {
			return refuse(codeDependentObjects, "%q is referenced by foreign key %q", index.Name, orphans[0].Name)
		}
		c.dropForeignKeys(orphans)
	}
	return nil
}
