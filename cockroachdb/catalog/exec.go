package catalog

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/cockroachdb/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

// Catalog is one database's schema, changed by DDL as CockroachDB changes it.
// It starts with the schema public, as a new database does.
//
// Exec accepts by default. It refuses a statement with CockroachDB's SQLSTATE
// where TestSimulation's rows, or a check of the pinned server, show
// CockroachDB refusing it: an object missing, taken or still needed by
// another. A statement whose effect on other objects the metadata cannot
// hold is refused with an Error that carries no code. It does not check what
// a statement's expressions compute.
type Catalog struct {
	l *sdlLoader
	// pending is the table ALTER TABLE is changing, whose columns a foreign
	// key it adds names.
	pending *metadatapb.TableMetadata
	// creating is the schema of the table CREATE TABLE is declaring, whose
	// keys pending holds; empty otherwise.
	creating string
	// computed holds the computed columns of the table CREATE TABLE or
	// ALTER TABLE is declaring.
	computed map[string]bool
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
		line += strings.Count(seg.Text[:leadLength(seg.Text)], "\n")
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

// leadLength is the length of the whitespace and comments a statement's
// text starts with.
func leadLength(text string) int {
	i := 0
	for i < len(text) {
		switch {
		case strings.ContainsRune(" \t\r\n", rune(text[i])):
			i++
		case strings.HasPrefix(text[i:], "--"):
			end := strings.IndexByte(text[i:], '\n')
			if end < 0 {
				return len(text)
			}
			i += end + 1
		case strings.HasPrefix(text[i:], "/*"):
			// Block comments nest.
			depth := 0
			for i < len(text) {
				if strings.HasPrefix(text[i:], "/*") {
					depth, i = depth+1, i+2
				} else if strings.HasPrefix(text[i:], "*/") {
					depth, i = depth-1, i+2
					if depth == 0 {
						break
					}
				} else {
					i++
				}
			}
		default:
			return i
		}
	}
	return i
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
	case *tree.Select:
		if strings.Contains(strings.ToLower(tree.AsString(n)), "search_path") {
			return &Error{Message: "setting search_path is not supported"}
		}
	case *tree.AlterTableOwner:
		// Owners are not held; the relation must exist.
		if relation(c.l.schemas[objectSchema(n.Name)], n.Name.Object()) == relNone && !n.IfExists {
			return refuse(codeUndefinedTable, "relation %q does not exist", n.Name.String())
		}
	case *tree.SetVar:
		if strings.EqualFold(n.Name, "search_path") || strings.EqualFold(n.Name, "schema") {
			// The catalog resolves every name in public.
			return &Error{Message: "SET search_path is not supported"}
		}
	case *tree.CreateExtension:
		if slices.ContainsFunc(c.l.meta.Extensions, func(e *metadatapb.ExtensionMetadata) bool { return e.Name == n.Name.String() }) {
			// CockroachDB takes it again, with IF NOT EXISTS or not.
			return nil
		}
	case *tree.CreateRoutine:
		return c.createRoutine(n)
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
		truncated := map[tableRef]bool{}
		for i := range n.Tables {
			table, err := c.table(&n.Tables[i])
			if err != nil {
				return err
			}
			truncated[tableRef{schemaOf(&n.Tables[i]), table.Name}] = true
		}
		for i := range n.Tables {
			for _, fk := range c.referencing(schemaOf(&n.Tables[i]), string(n.Tables[i].ObjectName), "") {
				if owner := c.foreignKeyTable(fk); owner != nil && !truncated[*owner] && n.DropBehavior != tree.DropCascade {
					return refuse("XXUUU", "%q is referenced by foreign key %q", n.Tables[i].ObjectName, fk.Name)
				}
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

// uniqueConstraintNames refuses a CREATE TABLE that names two constraints
// alike: 42P07 when both are unique indexes, 42710 otherwise.
func uniqueConstraintNames(table string, defs tree.TableDefs) error {
	indexes := map[string]bool{}
	add := func(name tree.Name, index bool) error {
		if name == "" {
			return nil
		}
		if wasIndex, ok := indexes[string(name)]; ok {
			if wasIndex && index {
				return refuse(codeDuplicateRelation, "duplicate index name: %q", name)
			}
			return refuse(codeDuplicateObject, "relation %q: duplicate constraint name: %q", table, name)
		}
		indexes[string(name)] = index
		return nil
	}
	for _, def := range defs {
		var err error
		switch d := def.(type) {
		case *tree.ColumnTableDef:
			for _, check := range d.CheckExprs {
				if err = add(check.ConstraintName, false); err != nil {
					return err
				}
			}
			if d.Unique.IsUnique || d.PrimaryKey.IsPrimaryKey {
				err = add(d.Unique.ConstraintName, true)
			}
			if err == nil && d.References.Table != nil {
				err = add(d.References.ConstraintName, false)
			}
		case *tree.CheckConstraintTableDef:
			err = add(d.Name, false)
		case *tree.ForeignKeyConstraintTableDef:
			err = add(d.Name, false)
		case *tree.UniqueConstraintTableDef:
			err = add(d.Name, true)
		case *tree.IndexTableDef:
			err = add(d.Name, true)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// createRoutine checks a CREATE FUNCTION or PROCEDURE as CockroachDB does
// before the loader stores it.
func (c *Catalog) createRoutine(n *tree.CreateRoutine) error {
	schema := "public"
	if n.Name.ExplicitSchema {
		schema = string(n.Name.SchemaName)
	}
	s := c.l.schemas[schema]
	if s == nil {
		return refuse(codeUndefinedSchema, "cannot create %q because the target database or schema does not exist", n.Name.String())
	}
	var args []string
	for _, param := range n.Params {
		if err := c.checkType(param.Type); err != nil {
			return err
		}
		if err := c.checkSequences(param.DefaultVal); err != nil {
			return err
		}
		if param.Class != tree.RoutineParamOut {
			args = append(args, signatureType(param.Type))
		}
	}
	if n.ReturnType != nil {
		if err := c.checkType(n.ReturnType.Type); err != nil {
			return err
		}
	}
	name := string(n.Name.ObjectName)
	signature := name + "(" + strings.Join(args, ", ") + ")"
	for _, function := range s.Functions {
		if function.Signature != signature {
			continue
		}
		if !n.Replace {
			return refuse("42723", "function %q already exists with same argument types", name)
		}
		stmt, err := parser.ParseOne(function.Definition)
		if err != nil {
			break
		}
		old, ok := stmt.AST.(*tree.CreateRoutine)
		if !ok {
			break
		}
		if old.IsProcedure != n.IsProcedure {
			return refuse(codeWrongObjectType, "cannot change routine kind")
		}
		if old.ReturnType != nil && n.ReturnType != nil && signatureType(old.ReturnType.Type) != signatureType(n.ReturnType.Type) {
			return refuse("42P13", "cannot change return type of existing function")
		}
	}
	return c.l.load(n)
}

// signatureType writes a type as a routine's signature holds it, a
// user-defined one under its schema.
func signatureType(ref tree.ResolvableTypeReference) string {
	if _, known := tree.GetStaticallyKnownType(ref); known {
		return ref.SQLString()
	}
	return schemaQualified(ref).SQLString()
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
	if s.Name == "public" {
		// Every unqualified name the catalog resolves is in public.
		return &Error{Message: "renaming schema public is not supported"}
	}
	for _, typ := range append(enumNames(s), compositeNames(s)...) {
		if err := c.routineTypeUser(s.Name, typ); err != nil {
			return err
		}
		if err := c.retypeBlocked(s.Name, typ); err != nil {
			return err
		}
	}
	if len(s.Functions) > 0 {
		// A routine's stored definition names its schema, which the
		// catalog does not rewrite.
		return &Error{Message: fmt.Sprintf("renaming schema %q, which holds routines, is not supported", s.Name)}
	}
	for _, sequence := range sequenceNames(s) {
		if err := c.viewSequenceUser(s.Name, sequence, relSequence); err != nil {
			return err
		}
	}
	for _, relation := range append(append(tableNames(s), viewNames(s)...), sequenceNames(s)...) {
		if readers := c.routineReaders(s.Name, relation); len(readers) > 0 {
			return refuse(codeDependentObjects, "cannot rename schema because function %q depends on relation %q", readers[0], relation)
		}
	}
	for _, other := range c.l.meta.Schemas {
		for _, deps := range viewDependencies(other) {
			if slices.ContainsFunc(deps, func(d *metadatapb.DependencyColumn) bool { return d.Schema == s.Name }) {
				// The view's query names the schema, which the catalog
				// does not rewrite.
				return &Error{Message: fmt.Sprintf("renaming schema %q, which a view reads, is not supported", s.Name)}
			}
		}
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
	for _, table := range s.Tables {
		reindex(table.Indexes, to, table.Name, "", "")
	}
	for _, view := range s.MaterializedViews {
		reindex(view.Indexes, to, view.Name, "", "")
	}
	// CockroachDB holds a default's sequence and a column's type by ID, so
	// they follow the schema.
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, column := range table.Columns {
			for _, sequence := range s.Sequences {
				column.Default = renameSequenceIn(column.Default, [2]string{from, sequence.Name}, [2]string{to, sequence.Name})
			}
		}
	})
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
		for _, typ := range append(slices.Clone(enumNames(s)), compositeNames(s)...) {
			for _, user := range c.typeUsers(name, typ) {
				if !strings.HasPrefix(user, tree.NameString(name)+".") {
					return refuse(codeFeatureNotSupport, "unimplemented: cannot drop type %q because other objects (%s) still depend on it", typ, user)
				}
			}
		}
		// Each relation goes with what CASCADE takes along, in other
		// schemas too.
		in := func(names ...string) tree.TableNames {
			var out tree.TableNames
			for _, n := range names {
				out = append(out, tree.MakeTableNameWithSchema("", tree.Name(name), tree.Name(n)))
			}
			return out
		}
		var tables, views, matviews, sequences []string
		for _, x := range s.Tables {
			tables = append(tables, x.Name)
		}
		for _, x := range s.Views {
			views = append(views, x.Name)
		}
		for _, x := range s.MaterializedViews {
			matviews = append(matviews, x.Name)
		}
		for _, x := range s.Sequences {
			sequences = append(sequences, x.Name)
		}
		for _, drop := range []struct {
			names []string
			kind  relKind
		}{{views, relView}, {matviews, relMaterializedView}, {tables, relTable}, {sequences, relSequence}} {
			if err := c.dropRelations(in(drop.names...), true, tree.DropCascade, drop.kind, true); err != nil {
				return err
			}
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
			if err := c.checkDefault(d.DefaultExpr.Expr, "CREATE TABLE"); err != nil {
				return err
			}
		}
	}
	primaryKeys := 0
	for _, def := range n.Defs {
		switch d := def.(type) {
		case *tree.ColumnTableDef:
			if d.PrimaryKey.IsPrimaryKey {
				primaryKeys++
			}
		case *tree.UniqueConstraintTableDef:
			if d.PrimaryKey {
				primaryKeys++
			}
		}
	}
	if primaryKeys > 1 {
		return refuse("XXUUU", "multiple primary keys for table %q are not allowed", name)
	}
	if err := uniqueConstraintNames(name, n.Defs); err != nil {
		return err
	}
	// The new table's columns and keys, for the foreign keys it declares,
	// its own among them.
	pending := &metadatapb.TableMetadata{Name: name}
	for _, def := range n.Defs {
		switch d := def.(type) {
		case *tree.ColumnTableDef:
			typ, _ := sdlType(d.Type)
			pending.Columns = append(pending.Columns, &metadatapb.ColumnMetadata{Name: string(d.Name), Type: typ})
			if d.PrimaryKey.IsPrimaryKey || d.Unique.IsUnique {
				pending.Indexes = append(pending.Indexes, &metadatapb.IndexMetadata{Unique: true, Primary: d.PrimaryKey.IsPrimaryKey, Expressions: []string{string(d.Name)}})
			}
		case *tree.UniqueConstraintTableDef:
			if d.Predicate == nil {
				key := &metadatapb.IndexMetadata{Unique: true, Primary: d.PrimaryKey}
				for _, elem := range d.Columns {
					key.Expressions = append(key.Expressions, string(elem.Column))
				}
				pending.Indexes = append(pending.Indexes, key)
			}
		}
	}
	c.pending, c.creating = pending, s.Name
	c.computed = map[string]bool{}
	for _, def := range n.Defs {
		if d, ok := def.(*tree.ColumnTableDef); ok && d.Computed.Computed {
			c.computed[string(d.Name)] = true
		}
	}
	defer func() { c.pending, c.creating, c.computed = nil, "", nil }()
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
		if d.Computed.Computed {
			if err := checkColumns(columns, d.Computed.Expr); err != nil {
				return err
			}
			if d.DefaultExpr.Expr != nil {
				return refuse(codeInvalidTableDef, "STORED COMPUTED COLUMN cannot have default values")
			}
			if !d.Computed.Virtual {
				for column := range c.computed {
					if exprReads(d.Computed.Expr, column) {
						return refuse(codeInvalidTableDef, "STORED COMPUTED COLUMN expression cannot reference computed columns")
					}
				}
			}
		}
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
		return checkIndexDef(columns, &d.IndexTableDef)
	case *tree.IndexTableDef:
		return checkIndexDef(columns, d)
	case *tree.ForeignKeyConstraintTableDef:
		return c.checkForeignKey(table, columns, d)
	}
	return nil
}

// checkDefault refuses a default that reads a column or calls nextval on a
// sequence the catalog does not hold.
func (c *Catalog) checkDefault(expr tree.Expr, in string) error {
	if expr == nil {
		return nil
	}
	reads := false
	tree.WalkExprConst(&columnCollector{visit: func(string) { reads = true }}, expr)
	if reads {
		return refuse("42601", "variable sub-expressions are not allowed in DEFAULT (in %s)", in)
	}
	return c.checkSequences(expr)
}

// checkSequences refuses a default that calls nextval on a sequence the
// catalog does not hold.
func (c *Catalog) checkSequences(expr tree.Expr) error {
	if expr == nil {
		return nil
	}
	for _, tn := range calledSequences(expr) {
		schema, sequence := schemaOf(tn), string(tn.ObjectName)
		switch relation(c.l.schemas[schema], sequence) {
		case relSequence:
		case relNone:
			return refuse(codeUndefinedTable, "relation %q does not exist", sequence)
		default:
			return refuse(codeWrongObjectType, "%q is not a sequence", sequence)
		}
	}
	return nil
}

// callsSequence reports whether a default calls nextval, currval or setval
// on a sequence.
func callsSequence(def, schema, name string) bool {
	if def == "" {
		return false
	}
	expr, err := parser.ParseExpr(def)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(calledSequences(expr), func(tn *tree.TableName) bool {
		return schemaOf(tn) == schema && string(tn.ObjectName) == name
	})
}

// calledSequences returns the sequences an expression's nextval, currval and
// setval name, each with its schema and name apart.
func calledSequences(expr tree.Expr) []*tree.TableName {
	var out []*tree.TableName
	_, _ = tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
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
		if annotation, ok := arg.(*tree.AnnotateTypeExpr); ok {
			arg = unparen(annotation.Expr)
		}
		if str, ok := arg.(*tree.StrVal); ok {
			if tn, err := parser.ParseQualifiedTableName(str.RawString()); err == nil {
				out = append(out, tn)
			}
		}
		return true, e, nil
	})
	return out
}

// checkIndexDef refuses an index a table definition declares that keys,
// stores or filters on a column the table does not have.
func checkIndexDef(columns map[string]bool, d *tree.IndexTableDef) error {
	if err := checkIndexColumns(columns, d.Columns); err != nil {
		return err
	}
	for _, col := range d.Storing {
		if !columns[string(col)] {
			return refuse(codeUndefinedColumn, "column %q does not exist", col)
		}
	}
	if d.Predicate != nil {
		return checkColumns(columns, d.Predicate)
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
	_, _ = tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		// A column named alone or under its relation (t.a).
		if name, ok := e.(*tree.UnresolvedName); ok && !name.Star && name.NumParts <= 2 && missing == "" && !columns[name.Parts[0]] {
			missing = name.Parts[0]
		}
		return true, e, nil
	})
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
	var referenced *metadatapb.TableMetadata
	if c.creating != "" && schemaOf(&d.Table) == c.creating && string(d.Table.ObjectName) == table {
		// A table's foreign key to itself, in its CREATE TABLE.
		referenced = c.pending
	} else {
		var err error
		if referenced, err = c.table(&d.Table); err != nil {
			return err
		}
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
	if len(to) != len(d.FromCols) {
		return refuse("42601", "%d columns must reference exactly %d columns in referenced table (found %d)", len(d.FromCols), len(d.FromCols), len(to))
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
	return refFamily(ref)
}

// refFamily is the family of a type reference: a user-defined type is a
// family of its own, and an array's family is its element's.
func refFamily(ref tree.ResolvableTypeReference) string {
	if array, ok := ref.(*tree.ArrayTypeReference); ok {
		return "array:" + refFamily(array.ElementType)
	}
	t, ok := tree.GetStaticallyKnownType(ref)
	if !ok {
		return "udt:" + schemaQualified(ref).SQLString()
	}
	if t.Family() == types.ArrayFamily {
		return "array:" + t.ArrayContents().Family().String()
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
	for _, x := range b {
		if !slices.Contains(a, x) {
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
	if kind == relNone && typeTaken(s, name) {
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	if err := c.checkViewQuery(n.AsSource); err != nil {
		return err
	}
	if kind != relNone {
		var cyclic bool
		storedRelations(n.AsSource, func(tn *tree.TableName) {
			cyclic = cyclic || c.reaches(schemaOf(tn), string(tn.ObjectName), s.Name, name)
		})
		if cyclic && n.Replace {
			return refuse("42P17", "cyclic view dependency for relation %s", name)
		}
	}
	if output := outputNames(n.AsSource); output != nil {
		names := output
		if len(n.ColumnNames) > 0 {
			if len(n.ColumnNames) != len(output) {
				return refuse("42601", "CREATE VIEW specifies %d column names, but data source has %d column", len(n.ColumnNames), len(output))
			}
			names = nil
			for _, col := range n.ColumnNames {
				names = append(names, string(col))
			}
		}
		seen := map[string]bool{}
		for _, col := range names {
			if seen[col] {
				return refuse(codeDuplicateColumn, "duplicate column name: %q", col)
			}
			seen[col] = true
		}
	}
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
	s := c.l.schemas[schemaOf(tn)]
	var relationColumns []*metadatapb.ColumnMetadata
	switch name := string(tn.ObjectName); relation(s, name) {
	case relTable:
		relationColumns = schemaTable(s, name).Columns
	case relView:
		relationColumns = findView(s, name).Columns
	case relMaterializedView:
		relationColumns = findMaterializedView(s, name).Columns
	}
	if len(relationColumns) == 0 {
		return nil
	}
	columns := map[string]bool{}
	for _, column := range relationColumns {
		columns[column.Name] = true
	}
	exprs := []tree.Expr{}
	for _, expr := range sel.Exprs {
		exprs = append(exprs, expr.Expr)
	}
	for _, window := range sel.Window {
		exprs = append(exprs, window.Partitions...)
		for _, order := range window.OrderBy {
			exprs = append(exprs, order.Expr)
		}
	}
	if query.Limit != nil {
		exprs = append(exprs, query.Limit.Count, query.Limit.Offset)
	}
	// ORDER BY may name an output column too.
	ordered := maps.Clone(columns)
	for _, expr := range sel.Exprs {
		if expr.As != "" {
			ordered[string(expr.As)] = true
		}
	}
	for _, order := range query.OrderBy {
		if order.Expr != nil {
			if err := checkColumns(ordered, order.Expr); err != nil {
				return err
			}
		}
	}
	// GROUP BY and DISTINCT ON may name an output column too.
	for _, expr := range append(slices.Clone(sel.GroupBy), sel.DistinctOn...) {
		if err := checkColumns(ordered, expr); err != nil {
			return err
		}
	}
	if sel.Where != nil {
		exprs = append(exprs, sel.Where.Expr)
	}
	if sel.Having != nil {
		exprs = append(exprs, sel.Having.Expr)
	}
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		if err := checkColumns(columns, expr); err != nil {
			return err
		}
	}
	return nil
}

// checkViewQuery refuses a view query that reads a relation the catalog does
// not hold, or, on one table, a column the table does not have.
func (c *Catalog) checkViewQuery(query *tree.Select) error {
	var missing error
	storedRelations(query, func(tn *tree.TableName) {
		// A cast's type comes through as a name too.
		if s := c.l.schemas[schemaOf(tn)]; missing == nil && relation(s, string(tn.ObjectName)) == relNone && (s == nil || !typeTaken(s, string(tn.ObjectName))) {
			missing = refuse(codeUndefinedTable, "relation %q does not exist", tn.String())
		}
	})
	if missing != nil {
		return missing
	}
	return c.checkViewColumns(query)
}

// reaches reports whether a relation is the target, or a view that reads it,
// directly or through other views.
func (c *Catalog) reaches(schema, name, targetSchema, target string) bool {
	seen := map[string]bool{}
	var walk func(schema, name string) bool
	walk = func(schema, name string) bool {
		if schema == targetSchema && name == target {
			return true
		}
		id := getObjectID(schema, name)
		if seen[id] {
			return false
		}
		seen[id] = true
		s := c.l.schemas[schema]
		var deps []*metadatapb.DependencyColumn
		if v := findView(s, name); s != nil && v != nil {
			deps = v.DependencyColumns
		} else if v := findMaterializedView(s, name); s != nil && v != nil {
			deps = v.DependencyColumns
		}
		for _, d := range deps {
			if walk(d.Schema, d.Table) {
				return true
			}
		}
		return false
	}
	return walk(schema, name)
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
	if relation(s, name) != relNone || typeTaken(s, name) {
		if n.IfNotExists {
			return nil
		}
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	if err := c.checkSequenceOptions(s.Name, n.Options); err != nil {
		return err
	}
	if err := c.l.load(n); err != nil {
		return err
	}
	return checkBounds(findSequence(s, name), true)
}

// checkBounds refuses a sequence whose start is outside its bounds.
// typed checks the bounds against the sequence's type too, as CREATE and an
// ALTER that states AS do; another ALTER widens a bound at the type's limit
// to INT8's (applySequenceOptions), which CockroachDB takes.
func checkBounds(sequence *metadatapb.SequenceMetadata, typed bool) error {
	minimum, err1 := strconv.ParseInt(sequence.MinValue, 10, 64)
	maximum, err2 := strconv.ParseInt(sequence.MaxValue, 10, 64)
	start, err3 := strconv.ParseInt(sequence.Start, 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil
	}
	if cache, err := strconv.ParseInt(sequence.CacheSize, 10, 64); err == nil && cache <= 0 {
		return refuse(codeInvalidParameter, "CACHE (%d) must be greater than zero", cache)
	}
	limit := sequenceLimit(sequence.DataType)
	switch {
	case sequence.Increment == "0":
		return refuse(codeInvalidParameter, "INCREMENT must not be zero")
	case typed && maximum > limit:
		return refuse(codeInvalidParameter, "MAXVALUE (%d) must be less than (%d) for type %s", maximum, limit, sequence.DataType)
	case typed && minimum < -limit-1:
		return refuse(codeInvalidParameter, "MINVALUE (%d) must be greater than (%d) for type %s", minimum, -limit-1, sequence.DataType)
	case minimum > maximum:
		return refuse(codeInvalidParameter, "MINVALUE (%d) must be less than MAXVALUE (%d)", minimum, maximum)
	case start > maximum:
		return refuse(codeInvalidParameter, "START value (%d) cannot be greater than MAXVALUE (%d)", start, maximum)
	case start < minimum:
		return refuse(codeInvalidParameter, "START value (%d) cannot be less than MINVALUE (%d)", start, minimum)
	}
	return nil
}

func (c *Catalog) checkSequenceOptions(schema string, options tree.SequenceOptions) error {
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
		if owner := ptr(option.ColumnItemVal.TableName.ToTableName()); schemaOf(owner) != schema {
			// The metadata names an owner by its table alone, in the
			// sequence's schema.
			return &Error{Message: fmt.Sprintf("a sequence owned by a table in another schema (%s) is not supported", owner)}
		}
	}
	return nil
}

func (c *Catalog) alterSequence(n *tree.AlterSequence) error {
	s := c.l.schemas[objectSchema(n.Name)]
	switch relation(s, n.Name.Object()) {
	case relSequence:
	case relNone:
		if n.IfExists {
			return nil
		}
		return refuse(codeUndefinedTable, "relation %q does not exist", n.Name.String())
	default:
		return refuse(codeWrongObjectType, "%q is not a sequence", n.Name.Object())
	}
	if err := c.checkSequenceOptions(s.Name, n.Options); err != nil {
		return err
	}
	// RESTART moves the sequence's value, which the catalog does not hold,
	// within its bounds.
	var restart *int64
	for _, o := range n.Options {
		if o.Name == tree.SeqOptRestart && o.IntVal != nil {
			restart = o.IntVal
		}
	}
	n.Options = slices.DeleteFunc(n.Options, func(o tree.SequenceOption) bool { return o.Name == tree.SeqOptRestart })
	if err := c.l.load(n); err != nil {
		return err
	}
	sequence := findSequence(s, n.Name.Object())
	typed := slices.ContainsFunc(n.Options, func(o tree.SequenceOption) bool { return o.Name == tree.SeqOptAs })
	if err := checkBounds(sequence, typed); err != nil {
		return err
	}
	if restart != nil {
		if maximum, err := strconv.ParseInt(sequence.MaxValue, 10, 64); err == nil && *restart > maximum {
			return refuse(codeInvalidParameter, "RESTART value (%d) cannot be greater than MAXVALUE (%d)", *restart, maximum)
		}
		if minimum, err := strconv.ParseInt(sequence.MinValue, 10, 64); err == nil && *restart < minimum {
			return refuse(codeInvalidParameter, "RESTART value (%d) cannot be less than MINVALUE (%d)", *restart, minimum)
		}
	}
	return nil
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
	labels := map[tree.Name]bool{}
	for _, elem := range n.CompositeTypeList {
		if labels[elem.Label] {
			return refuse("42P17", "composite type definition contains duplicate label %q", elem.Label)
		}
		labels[elem.Label] = true
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
		// An expression holds the label by its value, which follows the
		// rename; one that holds it where the catalog cannot tell it is the
		// enum's is refused, as is a view's query, which the catalog does
		// not rewrite.
		typ := typeName(s.Name, name)
		if err := c.viewTypeUser(s.Name, name, "renaming a value of"); err != nil {
			return err
		}

		from, to := string(cmd.OldVal), string(cmd.NewVal)
		var stuck error
		c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
			var enumColumns []string
			for _, column := range table.Columns {
				if column.Type == typ || column.Type == typ+"[]" {
					enumColumns = append(enumColumns, column.Name)
				}
			}
			relabel := func(text string) string {
				if text == "" {
					return text
				}
				return relabelIn(text, typ, from, to, enumColumns)
			}
			for _, column := range table.Columns {
				if column.Type == typ || column.Type == typ+"[]" {
					column.Default = relabel(column.Default)
					if expr, err := parser.ParseExpr(column.Default); err == nil && column.Default != "" && uncastLabel(expr, from) && stuck == nil {
						stuck = &Error{Message: fmt.Sprintf("renaming enum value %q, which the default of column %q holds uncast, is not supported", from, column.Name)}
					}
				} else if namesType(column.Default, typ) {
					// A cast to the enum in another column's default
					// (('sad'::mood)::STRING).
					column.Default = relabel(column.Default)
				}
			}
			reads := func(text string) bool {
				expr, err := parser.ParseExpr(text)
				return err == nil && (namesType(text, typ) || slices.ContainsFunc(enumColumns, func(col string) bool { return exprReads(expr, col) }))
			}
			// What still holds the label bare, beside an enum column, is
			// a form relabelIn does not know: refused.
			relabelRead := func(text string) string {
				out := relabel(text)
				if expr, err := parser.ParseExpr(out); err == nil && uncastLabel(expr, from) && stuck == nil {
					stuck = &Error{Message: fmt.Sprintf("renaming enum value %q, which an expression of table %q holds uncast, is not supported", from, table.Name)}
				}
				return out
			}
			for _, column := range table.Columns {
				if column.Generation != nil && reads(column.Generation.Expression) {
					column.Generation.Expression = relabelRead(column.Generation.Expression)
				}
			}
			for _, check := range table.CheckConstraints {
				if reads(check.Expression) {
					check.Expression = relabelRead(check.Expression)
				}
			}
		})
		if stuck != nil {
			return stuck
		}
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
		if err := c.viewTypeUser(s.Name, name, "removing a value of"); err != nil {
			return err
		}
		enum.Values = slices.Delete(enum.Values, at, at+1)
	case *tree.AlterTypeRename:
		if err := c.routineTypeUser(s.Name, name); err != nil {
			return err
		}
		if err := c.retypeBlocked(s.Name, name); err != nil {
			return err
		}
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
		if err := c.routineTypeUser(s.Name, name); err != nil {
			return err
		}
		if err := c.retypeBlocked(s.Name, name); err != nil {
			return err
		}
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
			column.Default = retypeIn(column.Default, old, now)
			if column.Generation != nil {
				column.Generation.Expression = retypeIn(column.Generation.Expression, old, now)
			}
		}
		for _, check := range table.CheckConstraints {
			check.Expression = retypeIn(check.Expression, old, now)
		}
	})
}

// retypeBlocked refuses renaming or moving a type that a view's query or an
// index's definition names, which the catalog does not rewrite.
func (c *Catalog) retypeBlocked(schema, name string) error {
	typ := typeName(schema, name)
	for _, s := range c.l.meta.Schemas {
		var definitions []string
		for _, v := range s.Views {
			definitions = append(definitions, v.Definition)
		}
		for _, v := range s.MaterializedViews {
			definitions = append(definitions, v.Definition)
			for _, index := range v.Indexes {
				definitions = append(definitions, index.Definition)
			}
		}
		for _, table := range s.Tables {
			for _, index := range table.Indexes {
				definitions = append(definitions, index.Definition)
			}
		}
		for _, definition := range definitions {
			if strings.Contains(definition, "::") && (strings.Contains(definition, typ) || strings.Contains(definition, tree.NameString(name))) {
				return &Error{Message: fmt.Sprintf("renaming or moving type %q, which a view or an index names, is not supported", typ)}
			}
		}
	}
	return nil
}

// retypeIn returns an expression with its casts and annotations to a
// user-defined type renamed, as CockroachDB, which holds the type by ID,
// writes them afterwards.
func retypeIn(text, old, now string) string {
	if text == "" {
		return text
	}
	expr, err := parser.ParseExpr(text)
	if err != nil {
		return text
	}
	ref, err := parser.GetTypeFromValidSQLSyntax(now)
	if err != nil {
		return text
	}
	// renamed returns the type a cast or annotation names, renamed, or nil
	// when it names another: the type itself or an array of it.
	var renamed func(t tree.ResolvableTypeReference) tree.ResolvableTypeReference
	renamed = func(t tree.ResolvableTypeReference) tree.ResolvableTypeReference {
		if array, ok := t.(*tree.ArrayTypeReference); ok {
			if element := renamed(array.ElementType); element != nil {
				return &tree.ArrayTypeReference{ElementType: element}
			}
			return nil
		}
		if _, known := tree.GetStaticallyKnownType(t); !known && schemaQualified(t).SQLString() == old {
			return ref
		}
		return nil
	}
	changed := false
	out, err := tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		switch x := e.(type) {
		case *tree.CastExpr:
			if t := renamed(x.Type); t != nil {
				cast := *x
				cast.Type, changed = t, true
				return true, &cast, nil
			}
		case *tree.AnnotateTypeExpr:
			if t := renamed(x.Type); t != nil {
				annotation := *x
				annotation.Type, changed = t, true
				return true, &annotation, nil
			}
		}
		return true, e, nil
	})
	if err != nil || !changed {
		return text
	}
	return tree.AsString(out)
}

// relabelIn returns a default of enum type typ with a label renamed: the
// default itself when it is the label, and the label cast or annotated to
// the type anywhere in it. Other strings are left as written.
func relabelIn(text, typ, from, to string, enumColumns []string) string {
	expr, err := parser.ParseExpr(text)
	if err != nil {
		return text
	}
	if str, ok := unparen(expr).(*tree.StrVal); ok && str.RawString() == from {
		return tree.AsString(tree.NewStrVal(to))
	}
	isEnum := func(t tree.ResolvableTypeReference) bool {
		_, known := tree.GetStaticallyKnownType(t)
		return !known && schemaQualified(t).SQLString() == typ
	}
	label := func(inner tree.Expr, t tree.ResolvableTypeReference) bool {
		str, ok := unparen(inner).(*tree.StrVal)
		return ok && isEnum(t) && str.RawString() == from
	}
	// relabelArray renames the label among an array's elements cast to an
	// array of the enum.
	relabelArray := func(inner tree.Expr, t tree.ResolvableTypeReference) (tree.Expr, bool) {
		elem, ok := t.(*tree.ArrayTypeReference)
		array, isArray := unparen(inner).(*tree.Array)
		if !ok || !isArray || !isEnum(elem.ElementType) {
			return nil, false
		}
		out := *array
		out.Exprs = make(tree.Exprs, len(array.Exprs))
		for i, e := range array.Exprs {
			out.Exprs[i] = e
			if str, ok := unparen(e).(*tree.StrVal); ok && str.RawString() == from {
				out.Exprs[i] = tree.NewStrVal(to)
			}
		}
		return &out, true
	}
	isEnumColumn := func(e tree.Expr) bool {
		name, ok := unparen(e).(*tree.UnresolvedName)
		return ok && !name.Star && slices.Contains(enumColumns, name.Parts[0])
	}
	bare := func(e tree.Expr) bool {
		str, ok := unparen(e).(*tree.StrVal)
		return ok && str.RawString() == from
	}
	out, err := tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		switch x := e.(type) {
		case *tree.RangeCond:
			// m BETWEEN 'sad' AND 'happy': the bounds are the enum's.
			if isEnumColumn(x.Left) {
				out := *x
				if bare(x.From) {
					out.From = tree.NewStrVal(to)
				}
				if bare(x.To) {
					out.To = tree.NewStrVal(to)
				}
				return true, &out, nil
			}
		case *tree.CaseExpr:
			// CASE m WHEN 'sad' ...: the WHEN values are the enum's.
			if x.Expr != nil && isEnumColumn(x.Expr) {
				out := *x
				out.Whens = make([]*tree.When, len(x.Whens))
				for i, w := range x.Whens {
					when := *w
					if bare(w.Cond) {
						when.Cond = tree.NewStrVal(to)
					}
					out.Whens[i] = &when
				}
				return true, &out, nil
			}
		case *tree.ComparisonExpr:
			// A label compared with an enum column is the enum's (m = 'sad').
			cmp := *x
			switch {
			case isEnumColumn(x.Left) && bare(x.Right):
				cmp.Right = tree.NewStrVal(to)
				return true, &cmp, nil
			case isEnumColumn(x.Right) && bare(x.Left):
				cmp.Left = tree.NewStrVal(to)
				return true, &cmp, nil
			case isEnumColumn(x.Left):
				// m IN ('sad', 'happy'): the tuple's labels are the enum's.
				if tuple, ok := unparen(x.Right).(*tree.Tuple); ok {
					out := *tuple
					out.Exprs = make(tree.Exprs, len(tuple.Exprs))
					for i, e := range tuple.Exprs {
						out.Exprs[i] = e
						if bare(e) {
							out.Exprs[i] = tree.NewStrVal(to)
						}
					}
					cmp.Right = &out
					return true, &cmp, nil
				}
			}
		case *tree.CastExpr:
			if inner, ok := relabelArray(x.Expr, x.Type); ok {
				cast := *x
				cast.Expr = inner
				return false, &cast, nil
			}
			if label(x.Expr, x.Type) {
				cast := *x
				cast.Expr = tree.NewStrVal(to)
				return false, &cast, nil
			}
		case *tree.AnnotateTypeExpr:
			if inner, ok := relabelArray(x.Expr, x.Type); ok {
				annotation := *x
				annotation.Expr = inner
				return false, &annotation, nil
			}
			if label(x.Expr, x.Type) {
				annotation := *x
				annotation.Expr = tree.NewStrVal(to)
				return false, &annotation, nil
			}
		}
		return true, e, nil
	})
	if err != nil {
		return text
	}
	return tree.AsString(out)
}

// viewDependencies returns the relations each view and materialized view of
// a schema reads.
func viewDependencies(s *metadatapb.SchemaMetadata) [][]*metadatapb.DependencyColumn {
	var out [][]*metadatapb.DependencyColumn
	for _, v := range s.Views {
		out = append(out, v.DependencyColumns)
	}
	for _, v := range s.MaterializedViews {
		out = append(out, v.DependencyColumns)
	}
	return out
}

func tableNames(s *metadatapb.SchemaMetadata) []string {
	var out []string
	for _, x := range s.Tables {
		out = append(out, x.Name)
	}
	return out
}

func viewNames(s *metadatapb.SchemaMetadata) []string {
	var out []string
	for _, x := range s.Views {
		out = append(out, x.Name)
	}
	for _, x := range s.MaterializedViews {
		out = append(out, x.Name)
	}
	return out
}

func sequenceNames(s *metadatapb.SchemaMetadata) []string {
	var out []string
	for _, x := range s.Sequences {
		out = append(out, x.Name)
	}
	return out
}

// routineReaders returns the routines, schema-qualified, whose bodies read
// a relation.
func (c *Catalog) routineReaders(schema, name string) []string {
	id := getObjectID(schema, name)
	var out []string
	for _, s := range c.l.meta.Schemas {
		for _, function := range s.Functions {
			if slices.Contains(routineRelations(function.Definition), id) {
				out = append(out, typeName(s.Name, function.Name))
			}
		}
	}
	return out
}

// routineColumn refuses renaming or dropping a column a routine that reads
// its table may name, whose stored body the catalog does not rewrite.
func (c *Catalog) routineColumn(schema, table, column string) error {
	for _, reader := range c.routineReaders(schema, table) {
		for _, s := range c.l.meta.Schemas {
			for _, function := range s.Functions {
				if typeName(s.Name, function.Name) == reader && identifierIn(routineBody(function.Definition), column) {
					return &Error{Message: fmt.Sprintf("changing column %q, which routine %s may read, is not supported", column, reader)}
				}
			}
		}
	}
	return nil
}

// identifierIn reports whether a text holds a name as a whole word.
func identifierIn(text, name string) bool {
	for i := strings.Index(text, name); i >= 0; {
		before := i == 0 || !isIdentChar(text[i-1])
		after := i+len(name) >= len(text) || !isIdentChar(text[i+len(name)])
		if before && after {
			return true
		}
		next := strings.Index(text[i+1:], name)
		if next < 0 {
			break
		}
		i += next + 1
	}
	return false
}

func isIdentChar(b byte) bool {
	return b == '_' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

// routineDefaultsName reports whether a routine's parameter defaults hold a
// name.
func routineDefaultsName(definition, name string) bool {
	stmt, err := parser.ParseOne(definition)
	if err != nil {
		return false
	}
	routine, ok := stmt.AST.(*tree.CreateRoutine)
	if !ok {
		return false
	}
	for _, param := range routine.Params {
		if param.DefaultVal != nil && identifierIn(tree.AsString(param.DefaultVal), name) {
			return true
		}
	}
	return false
}

// routineNamesType reports whether a routine's parameters or result have a
// type, as the catalog names a user-defined type.
func routineNamesType(definition, typ, name string) bool {
	// The body and the parameters' defaults are stored as written; a type
	// they name is not rewritten.
	if identifierIn(routineBody(definition), name) || routineDefaultsName(definition, name) {
		return true
	}
	stmt, err := parser.ParseOne(definition)
	if err != nil {
		return false
	}
	routine, ok := stmt.AST.(*tree.CreateRoutine)
	if !ok {
		return false
	}
	is := func(ref tree.ResolvableTypeReference) bool {
		if _, known := tree.GetStaticallyKnownType(ref); known {
			return false
		}
		name := schemaQualified(ref).SQLString()
		return name == typ || name == typ+"[]"
	}
	for _, param := range routine.Params {
		if is(param.Type) {
			return true
		}
	}
	return routine.ReturnType != nil && is(routine.ReturnType.Type)
}

func enumNames(s *metadatapb.SchemaMetadata) []string {
	var out []string
	for _, e := range s.EnumTypes {
		out = append(out, e.Name)
	}
	return out
}

func compositeNames(s *metadatapb.SchemaMetadata) []string {
	var out []string
	for _, t := range s.CompositeTypes {
		out = append(out, t.Name)
	}
	return out
}

// labelUser returns a table a column default of which holds an enum label.
func (c *Catalog) labelUser(typ, label string) string {
	user := ""
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		var texts []string
		for _, column := range table.Columns {
			if column.Type == typ || column.Type == typ+"[]" || namesType(column.Default, typ) {
				texts = append(texts, column.Default)
			}
			if column.Generation != nil {
				texts = append(texts, column.Generation.Expression)
			}
		}
		for _, check := range table.CheckConstraints {
			texts = append(texts, check.Expression)
		}
		for _, text := range texts {
			if user == "" && text != "" {
				expr, err := parser.ParseExpr(text)
				if err == nil && readsLabel(expr, label) {
					user = table.Name
				}
			}
		}
	})
	return user
}

// uncastLabel reports whether an expression holds a label as a bare string
// other than as a comparison's operand, where its type is not the catalog's
// to tell.
func uncastLabel(expr tree.Expr, label string) bool {
	operands := map[tree.Expr]bool{}
	found := false
	_, _ = tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		switch x := e.(type) {
		case *tree.ComparisonExpr:
			operands[unparen(x.Left)], operands[unparen(x.Right)] = true, true
		case *tree.CastExpr, *tree.AnnotateTypeExpr:
			return false, e, nil
		case *tree.StrVal:
			if x.RawString() == label && !operands[e] {
				found = true
			}
		}
		return true, e, nil
	})
	return found
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
			uses := column.Type == typ || column.Type == typ+"[]" || namesType(column.Default, typ)
			if column.Generation != nil {
				uses = uses || namesType(column.Generation.Expression, typ)
			}
			if uses {
				users = append(users, typeName(s.Name, table.Name))
				return
			}
		}
		for _, check := range table.CheckConstraints {
			if namesType(check.Expression, typ) {
				users = append(users, typeName(s.Name, table.Name))
				return
			}
		}
		for _, index := range table.Indexes {
			if indexNamesType(index.Definition, typ) {
				users = append(users, typeName(s.Name, table.Name))
				return
			}
		}
	})
	for _, s := range c.l.meta.Schemas {
		for _, composite := range s.CompositeTypes {
			for _, attr := range composite.Attributes {
				if attr.Type == typ || attr.Type == typ+"[]" {
					users = append(users, typeName(s.Name, composite.Name))
				}
			}
		}
		for _, v := range s.Views {
			if viewNamesType(v.Definition, typ) {
				users = append(users, typeName(s.Name, v.Name))
			}
		}
		for _, v := range s.MaterializedViews {
			if viewNamesType(v.Definition, typ) {
				users = append(users, typeName(s.Name, v.Name))
			}
		}
		// A routine names its parameter and result types by their schema.
		for _, function := range s.Functions {
			if routineNamesType(function.Definition, typ, name) {
				users = append(users, typeName(s.Name, function.Name))
			}
		}
	}
	return users
}

// viewTypeUser refuses changing a type a view's query, an index's
// expression or a routine names, which the catalog does not rewrite.
func (c *Catalog) viewTypeUser(schema, name, what string) error {
	typ := typeName(schema, name)
	for _, s := range c.l.meta.Schemas {
		for _, function := range s.Functions {
			if routineNamesType(function.Definition, typ, name) {
				return &Error{Message: fmt.Sprintf("%s type %q, which routine %q names, is not supported", what, typ, function.Name)}
			}
		}
		for _, table := range s.Tables {
			for _, index := range table.Indexes {
				if indexNamesType(index.Definition, typ) {
					return &Error{Message: fmt.Sprintf("%s type %q, which index %q names, is not supported", what, typ, index.Name)}
				}
			}
		}
		for _, v := range s.Views {
			if viewNamesType(v.Definition, typ) {
				return &Error{Message: fmt.Sprintf("%s type %q, which view %q names, is not supported", what, typ, v.Name)}
			}
		}
		for _, v := range s.MaterializedViews {
			if viewNamesType(v.Definition, typ) {
				return &Error{Message: fmt.Sprintf("%s type %q, which materialized view %q names, is not supported", what, typ, v.Name)}
			}
		}
	}
	return nil
}

// viewSequenceUser refuses renaming, moving or dropping a sequence a view's
// query calls, which the catalog does not track or rewrite.
func (c *Catalog) viewSequenceUser(schema, name string, kind relKind) error {
	if kind != relSequence {
		return nil
	}
	for _, s := range c.l.meta.Schemas {
		for _, function := range s.Functions {
			// Its body or a parameter's default.
			if identifierIn(function.Definition, name) {
				return &Error{Message: fmt.Sprintf("changing sequence %q, which routine %q may call, is not supported", name, function.Name)}
			}
		}
		var definitions []string
		for _, v := range s.Views {
			definitions = append(definitions, v.Definition)
		}
		for _, v := range s.MaterializedViews {
			definitions = append(definitions, v.Definition)
		}
		for _, definition := range definitions {
			stmt, err := parser.ParseOne(definition)
			if err != nil {
				continue
			}
			calls := false
			tree.WalkStmt(&exprVisitor{visit: func(e tree.Expr) {
				for _, tn := range calledSequences(e) {
					calls = calls || schemaOf(tn) == schema && string(tn.ObjectName) == name
				}
			}}, stmt.AST)
			if calls {
				return &Error{Message: fmt.Sprintf("changing sequence %q, which a view calls, is not supported", name)}
			}
		}
	}
	return nil
}

// indexNamesType reports whether an index's key expressions or predicate
// cast or annotate to a type.
func indexNamesType(definition, typ string) bool {
	stmt, err := parser.ParseOne(definition)
	if err != nil {
		return false
	}
	n, ok := stmt.AST.(*tree.CreateIndex)
	if !ok {
		return false
	}
	for _, elem := range n.Columns {
		if elem.Expr != nil && namesType(tree.AsString(elem.Expr), typ) {
			return true
		}
	}
	return n.Predicate != nil && namesType(tree.AsString(n.Predicate), typ)
}

// viewNamesType reports whether a view's query casts or annotates to a type.
func viewNamesType(definition, typ string) bool {
	stmt, err := parser.ParseOne(definition)
	if err != nil {
		return false
	}
	found := false
	tree.WalkStmt(&exprVisitor{visit: func(e tree.Expr) {
		switch x := e.(type) {
		case *tree.CastExpr, *tree.AnnotateTypeExpr:
			found = found || namesType(tree.AsString(x), typ)
		}
	}}, stmt.AST)
	return found
}

type exprVisitor struct{ visit func(tree.Expr) }

func (v *exprVisitor) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	v.visit(expr)
	return true, expr
}

func (*exprVisitor) VisitPost(expr tree.Expr) tree.Expr { return expr }

// namesType reports whether an expression casts or annotates to a type.
func namesType(text, typ string) bool {
	return text != "" && retypeIn(text, typ, typ+"_") != text
}

// routineTypeUser refuses renaming or moving a type a routine names, whose
// stored definition the catalog does not rewrite.
func (c *Catalog) routineTypeUser(schema, name string) error {
	typ := typeName(schema, name)
	for _, s := range c.l.meta.Schemas {
		for _, function := range s.Functions {
			if routineNamesType(function.Definition, typ, name) {
				return &Error{Message: fmt.Sprintf("renaming or moving type %q, which routine %q names, is not supported", typ, function.Name)}
			}
		}
	}
	return nil
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
	// A column alias list or a NATURAL join reads columns by position or
	// by name without naming them.
	if hasColumnAliases(stmt.AST) {
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

// hasColumnAliases reports whether a query renames a relation's columns
// (FROM t AS u(x, y)) or joins NATURALly, reading columns it does not name.
func hasColumnAliases(stmt tree.Statement) bool {
	var selects func(tree.SelectStatement) bool
	var tables func(tree.TableExpr) bool
	tables = func(te tree.TableExpr) bool {
		switch t := te.(type) {
		case *tree.AliasedTableExpr:
			if len(t.As.Cols) > 0 {
				return true
			}
			if sub, ok := t.Expr.(*tree.Subquery); ok {
				if sel, ok := sub.Select.(*tree.ParenSelect); ok {
					return selects(sel.Select.Select)
				}
			}
		case *tree.JoinTableExpr:
			switch t.Cond.(type) {
			case tree.NaturalJoinCond, *tree.UsingJoinCond:
				// A join on columns it names as a list, or on all.
				return true
			}
			return tables(t.Left) || tables(t.Right)
		case *tree.ParenTableExpr:
			return tables(t.Expr)
		}
		return false
	}
	selects = func(sel tree.SelectStatement) bool {
		switch s := sel.(type) {
		case *tree.SelectClause:
			return slices.ContainsFunc(s.From.Tables, tables)
		case *tree.UnionClause:
			return selects(s.Left.Select) || selects(s.Right.Select)
		case *tree.ParenSelect:
			return selects(s.Select.Select)
		}
		return false
	}
	if n, ok := stmt.(*tree.Select); ok {
		return selects(n.Select)
	}
	return false
}

type nameVisitor struct{ visit func(*tree.UnresolvedName) }

func (v *nameVisitor) VisitPre(expr tree.Expr) (bool, tree.Expr) {
	switch e := expr.(type) {
	case *tree.FuncExpr:
		// The walker does not descend into OVER (...).
		if w := e.WindowDef; w != nil {
			for _, p := range w.Partitions {
				tree.WalkExprConst(v, p)
			}
			for _, o := range w.OrderBy {
				if o.Expr != nil {
					tree.WalkExprConst(v, o.Expr)
				}
			}
		}
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

// supportsKey reports whether a unique index, not partial, keys exactly a
// foreign key's referenced columns.
func supportsKey(index *metadatapb.IndexMetadata, fk *metadatapb.ForeignKeyMetadata) bool {
	return index.Unique && !strings.Contains(index.Definition, " WHERE ") && sameSet(index.Expressions, fk.ReferencedColumns)
}

// keysOn returns every foreign key that references a table, its own
// included.
func (c *Catalog) keysOn(schema, name string) []*metadatapb.ForeignKeyMetadata {
	var out []*metadatapb.ForeignKeyMetadata
	c.eachTable(func(_ *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
		for _, fk := range table.ForeignKeys {
			if fk.ReferencedSchema == schema && fk.ReferencedTable == name {
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
		for _, reader := range c.routineReaders(t.s.Name, t.name) {
			if inSchema && strings.HasPrefix(reader, tree.NameString(t.s.Name)+".") {
				// It goes with its schema.
				continue
			}
			if !cascade {
				return refuse(codeDependentObjects, "cannot drop relation %q because function %q depends on it", t.name, reader)
			}
			return &Error{Message: fmt.Sprintf("dropping relation %q with the routine %q that reads it is not supported", t.name, reader)}
		}
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
			// A sequence a column of the table owns goes with it.
			for _, sequence := range t.s.Sequences {
				if sequence.OwnerTable == t.name && !dropping(t.s.Name, sequence.Name) {
					targets = append(targets, target{t.s, sequence.Name, relSequence})
				}
			}
			for _, fk := range c.referencing(t.s.Name, t.name, "") {
				if owner := c.foreignKeyTable(fk); owner != nil && dropping(owner.schema, owner.table) {
					continue
				}
				if !cascade && !inSchema {
					return refuse(codeDependentObjects, "%q is referenced by foreign key %q", t.name, fk.Name)
				}
			}
		case relSequence:
			if err := c.viewSequenceUser(t.s.Name, t.name, relSequence); err != nil {
				return err
			}
			var users []*metadatapb.ColumnMetadata
			c.eachTable(func(s *metadatapb.SchemaMetadata, table *metadatapb.TableMetadata) {
				if dropping(s.Name, table.Name) {
					return
				}
				for _, column := range table.Columns {
					if callsSequence(column.Default, t.s.Name, t.name) {
						users = append(users, column)
					}
				}
			})
			if len(users) > 0 {
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
	if readers := c.routineReaders(s.Name, name); len(readers) > 0 {
		return refuse(codeDependentObjects, "cannot rename relation %q because function %q depends on it", name, readers[0])
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
	if err := ownedAcross(s, name, got, target); err != nil {
		return err
	}
	if err := c.viewSequenceUser(s.Name, name, got); err != nil {
		return err
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
	want := relTable
	switch {
	case n.IsMaterialized:
		want = relMaterializedView
	case n.IsView:
		want = relView
	case n.IsSequence:
		want = relSequence
	}
	if want != relTable && got != want {
		return refuse(codeWrongObjectType, "%q is not a %s", name, want)
	}
	target, err := c.schemaFor(string(n.Schema))
	if err != nil {
		return err
	}
	if relation(target, name) != relNone || typeTaken(target, name) {
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
	}
	if deps := c.dependents(s.Name, name, ""); len(deps) > 0 {
		return refuse(codeDependentObjects, "cannot set schema on relation %q because view %q depends on it", name, deps[0].name)
	}
	if readers := c.routineReaders(s.Name, name); len(readers) > 0 {
		return refuse(codeDependentObjects, "cannot set schema on relation %q because function %q depends on it", name, readers[0])
	}
	if err := ownedAcross(s, name, got, target); err != nil {
		return err
	}
	if err := c.viewSequenceUser(s.Name, name, got); err != nil {
		return err
	}
	c.moveRelation(s, name, got, target, name)
	return nil
}

// ownedAcross refuses moving a table that owns a sequence, or a sequence a
// table owns, to another schema: CockroachDB leaves the other where it is,
// owned across schemas, which the metadata cannot state.
func ownedAcross(s *metadatapb.SchemaMetadata, name string, kind relKind, target *metadatapb.SchemaMetadata) error {
	if target == s {
		return nil
	}
	if kind == relTable && slices.ContainsFunc(s.Sequences, func(x *metadatapb.SequenceMetadata) bool { return x.OwnerTable == name }) {
		return &Error{Message: fmt.Sprintf("moving table %q, which owns a sequence, to another schema is not supported", name)}
	}
	if sequence := findSequence(s, name); kind == relSequence && sequence.OwnerTable != "" {
		return &Error{Message: fmt.Sprintf("moving sequence %q, which a table owns, to another schema is not supported", name)}
	}
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
		// A stored expression names its table (t.a) by the new name.
		for _, check := range table.CheckConstraints {
			check.Expression = requalifyIn(check.Expression, name, to)
		}
		for _, column := range table.Columns {
			if column.Generation != nil {
				column.Generation.Expression = requalifyIn(column.Generation.Expression, name, to)
			}
		}
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
			column.Default = renameSequenceIn(column.Default, [2]string{s.Name, name}, [2]string{target.Name, to})
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
					if i < len(index.Expressions) {
						index.Expressions[i] = sdlExpr(n.Columns[i].Expr)
					}
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
func renameSequenceIn(def string, from, to [2]string) string {
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
		if annotation, ok := arg.(*tree.AnnotateTypeExpr); ok {
			arg = unparen(annotation.Expr)
		}
		str, ok := arg.(*tree.StrVal)
		if !ok {
			return true, e, nil
		}
		if sequence, err := parser.ParseQualifiedTableName(str.RawString()); err == nil && schemaOf(sequence) == from[0] && string(sequence.ObjectName) == from[1] {
			renamed := *call
			renamed.Exprs = append(tree.Exprs{tree.NewStrVal(tree.NameString(to[0]) + "." + tree.NameString(to[1]))}, call.Exprs[1:]...)
			return false, &renamed, nil
		}
		return true, e, nil
	})
	if err != nil {
		return def
	}
	return tree.AsString(out)
}

// requalifyIn returns an expression with the relation its qualified column
// references name renamed.
func requalifyIn(text, from, to string) string {
	expr, err := parser.ParseExpr(text)
	if err != nil {
		return text
	}
	changed := false
	out, err := tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		if name, ok := e.(*tree.UnresolvedName); ok && name.NumParts == 2 && !name.Star && name.Parts[1] == from {
			renamed := *name
			renamed.Parts[1], changed = to, true
			return false, &renamed, nil
		}
		return true, e, nil
	})
	if err != nil || !changed {
		return text
	}
	return tree.AsString(out)
}

// renameColumnIn returns an expression with a column renamed.
func renameColumnIn(expr tree.Expr, from, to string) tree.Expr {
	out, err := tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		// A column named alone or under its relation (t.a); the qualifier
		// stays.
		if name, ok := e.(*tree.UnresolvedName); ok && name.NumParts <= 2 && !name.Star && name.Parts[0] == from {
			renamed := *name
			renamed.Parts[0] = to
			return false, &renamed, nil
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
	c.computed = map[string]bool{}
	for _, column := range table.Columns {
		if column.Generation != nil {
			c.computed[column.Name] = true
		}
	}
	defer func() { c.pending, c.computed = nil, nil }()
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
		if def.PrimaryKey.IsPrimaryKey {
			return refuse("42611", "multiple primary keys for table %q are not allowed", table.Name)
		}
		if def.Family.Create {
			return &Error{Message: fmt.Sprintf("column %s: CREATE FAMILY is not supported", def.Name)}
		}
		var names []tree.Name
		for _, check := range def.CheckExprs {
			names = append(names, check.ConstraintName)
		}
		if def.Unique.IsUnique {
			names = append(names, def.Unique.ConstraintName)
		}
		if def.References.Table != nil {
			names = append(names, def.References.ConstraintName)
		}
		seen := map[tree.Name]bool{}
		for _, name := range names {
			if name != "" && (seen[name] || constraintNamed(table, string(name)) || slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool { return i.Name == string(name) })) {
				return refuse(codeDuplicateObject, "duplicate constraint name: %q", name)
			}
			seen[name] = true
		}
		if err := c.checkType(def.Type); err != nil {
			return err
		}
		if err := c.checkDefault(def.DefaultExpr.Expr, "ADD COLUMN"); err != nil {
			return err
		}
		columns[string(def.Name)] = true
		if def.Computed.Computed {
			// Its expression may not read itself.
			c.computed[string(def.Name)] = true
		}
		// The foreign key the column declares checks the column's type.
		typ, _ := sdlType(def.Type)
		c.pending = &metadatapb.TableMetadata{Name: table.Name, Columns: append(slices.Clone(table.Columns), &metadatapb.ColumnMetadata{Name: string(def.Name), Type: typ})}
		err := c.checkTableDef(table.Name, columns, def)
		c.pending = table
		if err != nil {
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
		if !constraintTaken(table, string(cmd.Constraint)) {
			return refuse(codeWrongObjectType, "constraint %q of relation %q is not a foreign key, check, or unique without index constraint", cmd.Constraint, table.Name)
		}
	case *tree.AlterTableSetDefault:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		if col.Generation != nil {
			if cmd.Default == nil {
				return refuse("42601", "column %q of relation %q is a computed column", col.Name, table.Name)
			}
			return refuse("42601", "computed column %q cannot also have a DEFAULT expression", col.Name)
		}
		if err := c.checkDefault(cmd.Default, "SET DEFAULT"); err != nil {
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
		if cmd.Collation != "" {
			if family(typ) != types.StringFamily.String() && family(typ) != types.CollatedStringFamily.String() {
				return refuse("42601", "COLLATE can only be used with string types")
			}
			collation = cmd.Collation
		}
		if cmd.Using != nil {
			if err := checkColumns(columns, cmd.Using); err != nil {
				return err
			}
		}
		// A change of family or of collation rewrites the stored data.
		rewrites := family(typ) != family(col.Type) || collation != col.Collation
		if col.Generation != nil && family(typ) != family(col.Type) {
			return refuse("XXUUU", "expected STORED COMPUTED COLUMN expression to have type %s", typ)
		}
		inForeignKey := len(c.referencing(schema, table.Name, col.Name)) > 0 ||
			slices.ContainsFunc(table.ForeignKeys, func(fk *metadatapb.ForeignKeyMetadata) bool { return slices.Contains(fk.Columns, col.Name) })
		for _, other := range table.Columns {
			if other.Generation == nil || other == col {
				continue
			}
			if expr, err := parser.ParseExpr(other.Generation.Expression); err == nil && exprReads(expr, col.Name) {
				return &Error{Message: fmt.Sprintf("changing the type of column %q, which computed column %q reads, is not supported", col.Name, other.Name)}
			}
		}
		inIndex := slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool { return indexReadsColumn(i, col.Name) })
		inCheck := slices.ContainsFunc(table.CheckConstraints, func(k *metadatapb.CheckConstraintMetadata) bool {
			expr, err := parser.ParseExpr(k.Expression)
			return err == nil && exprReads(expr, col.Name)
		})
		if inIndex && rewrites {
			return refuse(codeFeatureNotSupport, "unimplemented: ALTER COLUMN TYPE requiring rewrite of on-disk data is currently not supported for columns that are part of an index")
		}
		if (inForeignKey || inCheck) && rewrites {
			return refuse(codeFeatureNotSupport, "unimplemented: ALTER COLUMN TYPE requiring a rewrite of on-disk data is not supported for columns that have constraints")
		}
		col.Type, col.Collation = typ, collation
	case *tree.AlterTableRenameColumn:
		col, err := column(cmd.Column)
		if err != nil {
			return err
		}
		if deps := c.dependents(schema, table.Name, col.Name); len(deps) > 0 {
			return refuse(codeDependentObjects, "cannot rename column %q because view %q depends on it", col.Name, deps[0].name)
		}
		if err := c.routineColumn(schema, table.Name, col.Name); err != nil {
			return err
		}
		for _, index := range table.Indexes {
			if strings.Contains(index.Definition, " PARTITION BY ") && identifierIn(index.Definition[strings.Index(index.Definition, " PARTITION BY "):], col.Name) {
				return &Error{Message: fmt.Sprintf("renaming column %q, which index %q partitions by, is not supported", col.Name, index.Name)}
			}
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
		if slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool { return i.Name == to }) {
			return refuse(codeDuplicateRelation, "relation %s already exists", to)
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
					// CockroachDB keeps it as t_<columns>_key, numbered past
					// a name another index holds (t_a_key1).
					n.Name = ""
					table.Indexes = append(table.Indexes, sdlIndex(schema, table.Name, table.Indexes, n))
				}
			}
		}
	}
	if name == "" && old != nil {
		name = old.Name
	} else if slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool { return i.Name == name }) {
		return refuse(codeDuplicateRelation, "relation %q already exists", name)
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
	if err := c.routineColumn(schema, table.Name, name); err != nil {
		return err
	}
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
	table.ForeignKeys = slices.DeleteFunc(table.ForeignKeys, func(fk *metadatapb.ForeignKeyMetadata) bool { return slices.Contains(fk.Columns, name) })
	// A foreign key that relied on an index the column took along needs
	// CASCADE, which takes it.
	var orphans []*metadatapb.ForeignKeyMetadata
	for _, fk := range c.keysOn(schema, table.Name) {
		if !slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool {
			return supportsKey(i, fk)
		}) {
			orphans = append(orphans, fk)
		}
	}
	if len(orphans) > 0 && !cascade {
		return &Error{Message: fmt.Sprintf("dropping column %q drops an index foreign key %q relies on", name, orphans[0].Name)}
	}
	c.dropForeignKeys(orphans)
	table.CheckConstraints = slices.DeleteFunc(table.CheckConstraints, func(check *metadatapb.CheckConstraintMetadata) bool {
		expr, err := parser.ParseExpr(check.Expression)
		return err == nil && checkColumns(map[string]bool{}, expr) != nil && exprReads(expr, name)
	})
	table.ForeignKeys = slices.DeleteFunc(table.ForeignKeys, func(fk *metadatapb.ForeignKeyMetadata) bool { return slices.Contains(fk.Columns, name) })
	// A sequence the column owns goes with it.
	var owned tree.TableNames
	for _, sequence := range c.l.schemas[schema].GetSequences() {
		if sequence.OwnerTable == table.Name && sequence.OwnerColumn == name {
			owned = append(owned, tree.MakeTableNameWithSchema("", tree.Name(schema), tree.Name(sequence.Name)))
		}
	}
	return c.dropRelations(owned, true, behavior, relSequence, false)
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
	_, _ = tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		// A column named alone or under its relation (t.a).
		if name, ok := e.(*tree.UnresolvedName); ok && !name.Star && name.NumParts <= 2 && name.Parts[0] == column {
			found = true
		}
		return true, e, nil
	})
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
		var found *metadatapb.IndexMetadata
		var in *[]*metadatapb.IndexMetadata
		for _, list := range lists {
			for _, index := range *list {
				if index.Name != name {
					continue
				}
				if found != nil {
					return nil, nil, refuse("42P08", "index name %q is ambiguous", name)
				}
				found, in = index, list
			}
		}
		if found != nil {
			return found, in, nil
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
			for _, fk := range c.keysOn(s.Name, table.Name) {
				if !slices.ContainsFunc(table.Indexes, func(i *metadatapb.IndexMetadata) bool {
					return supportsKey(i, fk)
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
