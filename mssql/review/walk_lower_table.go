package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// tableNeeds adds the check that a reference names a relation of one of
// kinds: a table, or for some statements a view too.
func (b *builder) tableNeeds(r objRef, kinds []string, absent string, present func(string) string) {
	if r.where != placed {
		return
	}
	space, parent := r.relSpace()
	if r.temp {
		kinds = nil
	}
	b.needs(space, parent, r.n, kinds, note(r.rng, absent, present))
}

// ofKind adds the check that a name, when it exists, is of one of kinds:
// a wrong kind fails even under IF EXISTS (3705, verified).
func (b *builder) ofKind(space, parent string, n name, kinds []string, nt *walkNote) {
	if s := b.see(space, parent, n); s.fact == walk.Present {
		b.needs(space, parent, n, kinds, nt)
	}
}

// createdTable records that the statement created the table a reference
// names. When the table is the only relation of its class, its children
// are new: they are cleared, fill records them, and createdTable reports
// true.
func (b *builder) createdTable(r objRef, fill func(path string)) bool {
	if r.where != placed {
		b.arriveObject(r, kindTable)
		if r.where == unplaced {
			b.markAllTablesDirty()
		}
		return false
	}
	space, parent := r.relSpace()
	fresh, _ := b.classFact(r.n.classKey(space, parent))
	b.arriveObject(r, kindTable)
	path := r.tablePath()
	b.markTableDirty(path)
	if fresh != walk.Absent {
		b.forgetChildren(path)
		return false
	}
	b.clearChildren(path)
	fill(path)
	return true
}

// builtinSystemTypes are SQL Server's built-in system data type names
// (including ANSI synonyms). User-defined alias types (CREATE TYPE ... FROM)
// may only derive from built-in system types, and column-type checks only
// fire for recognized built-in types so alias/UDT columns are never misjudged.
var builtinSystemTypes = map[string]bool{
	"bigint": true, "int": true, "integer": true, "smallint": true, "tinyint": true, "bit": true,
	"decimal": true, "dec": true, "numeric": true, "money": true, "smallmoney": true,
	"float": true, "real": true, "double precision": true,
	"date": true, "time": true, "datetime": true, "datetime2": true, "datetimeoffset": true, "smalldatetime": true,
	"char": true, "character": true, "varchar": true, "char varying": true, "character varying": true,
	"nchar": true, "nchar varying": true, "national char": true, "national character": true,
	"nvarchar": true, "national char varying": true, "national character varying": true,
	"text": true, "ntext": true, "national text": true,
	"binary": true, "varbinary": true, "binary varying": true, "image": true,
	"uniqueidentifier": true, "xml": true, "sql_variant": true, "hierarchyid": true,
	"geometry": true, "geography": true, "timestamp": true, "rowversion": true,
	"sysname": true, "json": true, "vector": true,
}

func isBuiltinSystemType(name string) bool {
	return builtinSystemTypes[strings.ToLower(strings.TrimSpace(name))]
}

// isSystemBaseType reports whether dt names a built-in SQL Server system type
// (unqualified or in schema sys).
func isSystemBaseType(dt *ast.DataType) bool {
	if dt == nil {
		return true
	}
	if dt.Schema != "" && !strings.EqualFold(dt.Schema, "sys") {
		return false
	}
	return isBuiltinSystemType(dt.Name)
}

func typeDisplay(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	if dt.Schema != "" {
		return ident(dt.Schema) + "." + ident(dt.Name)
	}
	return ident(dt.Name)
}

// nonIndexableKeyType reports whether a synced base type cannot be used as a
// key column in a B-tree/columnstore index, statistics, PK, or UNIQUE
// constraint (1977, 1978, 1919, 35343).
func nonIndexableKeyType(typ string) bool {
	return typ == "xml" || typ == "geometry" || typ == "geography"
}

// dupColumnInLists reports the first exact column name that appears twice
// across the given column lists (*ast.IndexColumn or *ast.String items).
func dupColumnInLists(lists ...*ast.List) (string, review.Range, bool) {
	seen := make(map[string]bool)
	for _, list := range lists {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			var name string
			var loc ast.Loc
			switch v := item.(type) {
			case *ast.IndexColumn:
				name, loc = v.Name, v.Loc
			case *ast.String:
				name = v.Str
			default:
				continue
			}
			if name == "" {
				continue
			}
			if seen[name] {
				return name, rangeOf(loc), true
			}
			seen[name] = true
		}
	}
	return "", review.Range{}, false
}

// checkIndexColList checks that every non-pseudo column in list exists on a
// known table (1911/1769) and, when checkKeyType is true, is not of a
// non-indexable system type such as xml, geometry, or geography.
func (b *builder) checkIndexColList(r objRef, path string, list *ast.List, checkKeyType bool) {
	if list == nil {
		return
	}
	for _, item := range list.Items {
		var name string
		var loc ast.Loc
		switch v := item.(type) {
		case *ast.IndexColumn:
			name, loc = v.Name, v.Loc
		case *ast.String:
			name = v.Str
		default:
			continue
		}
		if name == "" || isPseudoCol(name) {
			continue
		}
		rng := b.anchor()
		if loc.Start >= 0 && loc.End > loc.Start {
			rng = rangeOf(loc)
		}
		if cn, ok := keyed(name); ok {
			b.needs(spaceColumn, path, cn, nil, note(rng, msgMissingIn("column", ident(name), r.written), nil))
		}
		if checkKeyType {
			if typ := b.knownColType(r, name); nonIndexableKeyType(typ) {
				b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(name))+" of type "+review.Code(typ)+" cannot be used as an index or statistics key column.")
			}
		}
	}
}

// columnsPresent reports whether every column an index's lists name is
// certainly on the table.
func (b *builder) columnsPresent(path string, lists ...*ast.List) bool {
	for _, list := range lists {
		if list == nil {
			continue
		}
		for _, item := range list.Items {
			var col string
			switch v := item.(type) {
			case *ast.IndexColumn:
				col = v.Name
			case *ast.String:
				col = v.Str
			default:
				return false
			}
			cn, ok := keyed(col)
			if !ok || b.see(spaceColumn, path, cn).fact != walk.Present {
				return false
			}
		}
	}
	return true
}

// constraintPrefixes are the generated names of unnamed constraints,
// PK__<table>__<hex> and the like (verified), as class prefixes.
var constraintPrefixes = map[ast.ConstraintType]string{
	ast.ConstraintPrimaryKey: classPrefix("PK__"),
	ast.ConstraintUnique:     classPrefix("UQ__"),
	ast.ConstraintForeignKey: classPrefix("FK__"),
	ast.ConstraintCheck:      classPrefix("CK__"),
	ast.ConstraintDefault:    classPrefix("DF__"),
}

var walkConstraintKinds = map[ast.ConstraintType]string{
	ast.ConstraintPrimaryKey: kindPK,
	ast.ConstraintUnique:     kindUQ,
	ast.ConstraintForeignKey: kindFK,
	ast.ConstraintCheck:      kindCK,
	ast.ConstraintDefault:    kindDF,
}

// tableDef is what a CREATE TABLE or an ALTER TABLE ... ADD defines.
type tableDef struct {
	columns     []*ast.ColumnDef
	constraints []*ast.ConstraintDef
	// defaults holds the columns with an unnamed DEFAULT.
	defaults int
	indexes  []*ast.InlineIndexDef
}

func (d *tableDef) addColumn(c *ast.ColumnDef) {
	d.columns = append(d.columns, c)
	if c.DefaultExpr != nil {
		d.defaults++
	}
	if c.Constraints != nil {
		for _, item := range c.Constraints.Items {
			if cd, ok := item.(*ast.ConstraintDef); ok {
				d.constraints = append(d.constraints, cd)
			}
		}
	}
}

// marked reports whether a name a generated constraint name may hold is
// marked: the table's, a column's the definition adds, or one its
// constraints name.
func (d *tableDef) marked(r objRef) bool {
	if !r.n.plain {
		return true
	}
	isMarked := func(spelling string) bool {
		n, ok := keyed(spelling)
		return !ok || !n.plain
	}
	for _, c := range d.columns {
		if isMarked(c.Name) {
			return true
		}
	}
	for _, c := range d.constraints {
		if c.Columns == nil {
			continue
		}
		for _, item := range c.Columns.Items {
			switch v := item.(type) {
			case *ast.String:
				if isMarked(v.Str) {
					return true
				}
			case *ast.IndexColumn:
				if isMarked(v.Name) {
					return true
				}
			}
		}
	}
	return false
}

// keyedConstraint is a constraint the walk keeps: one of the five kinds.
func keyedConstraint(c *ast.ConstraintDef) bool {
	_, ok := walkConstraintKinds[c.Type]
	return ok
}

// defTextFails checks what a definition's text alone settles: no column,
// constraint, or index name twice in the statement (2705, 8168), at most
// one primary key (8110), no constraint named like the table (2714), and
// no identifier too long (103). verb is what the statement does with the
// definition, declares or adds. It reports whether the statement
// certainly fails.
func (b *builder) defTextFails(r objRef, d *tableDef, verb string) bool {
	seen := make(map[string]bool)
	for _, c := range d.columns {
		if b.tooLong(c.Name, false, rangeOf(c.Loc)) {
			return true
		}
		if seen[c.Name] {
			b.fail(rangeOf(c.Loc), review.DuplicateName, msgTwice(verb, "column", ident(c.Name)))
			return true
		}
		seen[c.Name] = true
	}
	// One statement may not name two constraints or indexes alike, of
	// any kinds: a CHECK, DEFAULT, or FOREIGN KEY and an inline INDEX of
	// one name fail too (8168, verified). names holds the kind word of
	// each name's first use, which the message keeps when the second use
	// is of the same kind.
	names := make(map[string]string)
	twice := func(first, kind string) string {
		if first != kind {
			return ""
		}
		return kind
	}
	pks := 0
	for _, c := range d.constraints {
		if !keyedConstraint(c) {
			continue
		}
		if c.Type == ast.ConstraintPrimaryKey {
			pks++
		}
		if c.Name == "" {
			continue
		}
		if b.tooLong(c.Name, false, rangeOf(c.Loc)) {
			return true
		}
		if first, ok := names[c.Name]; ok {
			b.fail(rangeOf(c.Loc), review.DuplicateName, msgTwice(verb, twice(first, "constraint"), ident(c.Name)))
			return true
		}
		names[c.Name] = "constraint"
		if !r.temp && c.Name == r.n.spelling {
			b.fail(rangeOf(c.Loc), review.DuplicateName, msgTwice("names", "", ident(c.Name)))
			return true
		}
	}
	for _, i := range d.indexes {
		if i.Name == "" {
			continue
		}
		if b.tooLong(i.Name, false, rangeOf(i.Loc)) {
			return true
		}
		if first, ok := names[i.Name]; ok {
			b.fail(rangeOf(i.Loc), review.DuplicateName, msgTwice(verb, twice(first, "index"), ident(i.Name)))
			return true
		}
		names[i.Name] = "index"
	}
	if pks > 1 {
		b.fail(b.anchor(), review.NameTaken, "This statement gives "+review.Code(r.written)+" two primary keys.")
		return true
	}
	for _, i := range d.indexes {
		if !i.Columnstore && (i.Columns == nil || len(i.Columns.Items) == 0) {
			b.fail(rangeOf(i.Loc), review.UnsupportedForm, "Index "+review.Code(ident(i.Name))+" must specify at least one column.")
			return true
		}
		if dup, dupRng, hasDup := dupColumnInLists(i.Columns, i.IncludeCols); hasDup {
			if dupRng.Start == 0 && dupRng.End == 0 {
				dupRng = rangeOf(i.Loc)
			}
			b.fail(dupRng, review.DuplicateName, msgTwice("lists", "column", ident(dup)))
			return true
		}
	}
	return false
}

// defNameChecks adds what a definition needs from the ledger: each
// constraint name free in the schema (2714), and each foreign key's
// target a table (1767, 1768). A temporary table's constraint names are
// tempdb's, and its foreign keys are skipped.
func (b *builder) defNameChecks(r objRef, d *tableDef) {
	if r.temp {
		return
	}
	for _, c := range d.constraints {
		if !keyedConstraint(c) {
			continue
		}
		if c.Type == ast.ConstraintForeignKey && c.RefTable != nil {
			b.foreignTarget(r, c.RefTable)
		}
		cn, ok := keyed(c.Name)
		if c.Name == "" || !ok || r.where != placed {
			continue
		}
		crng := rangeOf(c.Loc)
		present := taken(walkConstraintKinds[c.Type], ident(c.Name))
		b.lacks(spaceObject, r.parent, cn, note(crng, "", present))
		b.lacks(spaceRelation, r.parent, cn, note(crng, "", present))
	}
}

// foreignTarget adds the check that a foreign key's target is a table
// (1767 missing, 1768 another kind). A table's foreign key to itself is
// checked by the statement that creates it.
func (b *builder) foreignTarget(r objRef, tr *ast.TableRef) {
	t, ok := b.ref(tr, false)
	if !ok || t.where != placed || t.temp {
		return
	}
	if !differ(t.schema.spelling, r.schema.spelling) && !differ(t.n.spelling, r.n.spelling) {
		// It may name the table itself, as dbo.X does dbo.x under a
		// case-insensitive collation (verified).
		return
	}
	b.tableNeeds(t, []string{kindTable}, msgMissing("table", t.written), notKind(t.written, "table", kindTable))
}

// defEffects records what a definition adds to a table whose children
// the walk knows, and its constraints' names in the schema.
func (b *builder) defEffects(r objRef, path string, known bool, d *tableDef) {
	if known {
		for _, c := range d.columns {
			if n, ok := keyed(c.Name); ok {
				b.arrive(spaceColumn, path, n, "")
			}
		}
		for _, i := range d.indexes {
			if n, ok := keyed(i.Name); ok {
				b.arrive(spaceIndex, path, n, kindIndex)
			}
		}
	}
	marked := d.marked(r)
	if d.defaults > 0 {
		b.generated(r, path, known, ast.ConstraintDefault, marked)
	}
	for _, c := range d.constraints {
		if !keyedConstraint(c) || r.temp && c.Type == ast.ConstraintForeignKey {
			// A temporary table skips its foreign keys (verified).
			continue
		}
		kind := walkConstraintKinds[c.Type]
		if c.Type == ast.ConstraintPrimaryKey && known {
			b.Set(primaryKey(path), walk.Present, "")
		}
		cn, ok := keyed(c.Name)
		if c.Name == "" || !ok {
			b.generated(r, path, known, c.Type, marked)
			continue
		}
		b.constraintArrives(r, path, known, cn, kind)
	}
}

func primaryKey(path string) walk.Key {
	return walk.Key{Space: spacePrimary, Parent: path, Name: primaryName}
}

// constraintArrives records a named constraint: in the schema's object
// space, the table's constraint space, and for a PK or UQ its index
// space. A temporary table's constraint names are tempdb's.
func (b *builder) constraintArrives(r objRef, path string, known bool, cn name, kind string) {
	if !r.temp {
		b.arriveObject(objRef{n: cn, schema: r.schema, parent: r.parent, where: r.where}, kind)
	}
	if !known {
		return
	}
	b.arrive(spaceConstraint, path, cn, kind)
	if kind == kindPK || kind == kindUQ {
		b.arrive(spaceIndex, path, cn, kind)
	}
}

// generated records an unnamed constraint, whose generated name the walk
// does not derive.
//
// The name holds the table's name, and a column's for some kinds. When
// one of those is marked, so may the name be, and then it may equal a
// name of any class: every name of the spaces may have been added.
func (b *builder) generated(r objRef, path string, known bool, t ast.ConstraintType, marked bool) {
	prefix := constraintPrefixes[t]
	if marked {
		prefix = ""
	}
	if !r.temp {
		switch r.where {
		case placed:
			b.MayAdd(walk.Scope{Space: spaceObject, Parent: r.parent, Prefix: prefix})
		case unplaced:
			b.MayAdd(walk.Scope{Space: spaceObject, Deep: true, Prefix: prefix})
		}
	}
	if !known {
		return
	}
	b.MayAdd(walk.Scope{Space: spaceConstraint, Parent: path, Prefix: prefix})
	if t == ast.ConstraintPrimaryKey || t == ast.ConstraintUnique {
		b.MayAdd(walk.Scope{Space: spaceIndex, Parent: path, Prefix: prefix})
	}
}

// unsupportedTable reports whether a CREATE TABLE makes more than the
// walk follows: a graph or file table, or a system-versioned or ledger
// table, which add history tables and hidden objects.
func unsupportedTable(v *ast.CreateTableStmt) bool {
	if v.IsNode || v.IsEdge || v.IsFileTable || v.PeriodStartCol != "" {
		return true
	}
	if v.TableOptions != nil {
		for _, item := range v.TableOptions.Items {
			if o, ok := item.(*ast.TableOption); ok {
				switch strings.ToUpper(o.Name) {
				case "SYSTEM_VERSIONING", "LEDGER", "MEMORY_OPTIMIZED":
					return true
				}
			}
		}
	}
	return false
}

// createTable follows CREATE TABLE. One object: every check holds against
// the state before the statement. An unsupported table kind (graph, file,
// system-versioned, ledger, or memory-optimized) still checks schema
// existence and object-name availability before ending the walk.
func (b *builder) createTable(v *ast.CreateTableStmt) bool {
	r, ok := b.ref(v.Name, true)
	if !ok {
		return false
	}
	if r.where == elsewhere {
		return true
	}
	b.Order = walk.Lookups
	if b.refTooLong(v.Name, r) {
		return true
	}
	if unsupportedTable(v) {
		b.schemaNeeds(r)
		b.lacksObject(r, kindTable)
		b.Last = true
		return true
	}
	d := &tableDef{}
	if v.Columns != nil {
		for _, item := range v.Columns.Items {
			if c, ok := item.(*ast.ColumnDef); ok {
				d.addColumn(c)
			}
		}
	}
	if v.Constraints != nil {
		for _, item := range v.Constraints.Items {
			if c, ok := item.(*ast.ConstraintDef); ok {
				d.constraints = append(d.constraints, c)
			}
		}
	}
	if v.Indexes != nil {
		for _, item := range v.Indexes.Items {
			if i, ok := item.(*ast.InlineIndexDef); ok {
				d.indexes = append(d.indexes, i)
			}
		}
	}
	if b.defTextFails(r, d, "declares") {
		return true
	}
	// The table's expressions are bound before its name is resolved, once
	// its column types are found: CREATE TABLE nope.t (c NVARCHAR(10)
	// DEFAULT (N'x' COLLATE a)) is 448 (measured, mssql/catalog).
	if systemTypes(d) && b.tableCollateFails(d, nil) {
		return true
	}
	b.schemaNeeds(r)
	b.lacksObject(r, kindTable)
	b.defNameChecks(r, d)
	fresh := b.createdTable(r, func(path string) {
		b.Set(primaryKey(path), walk.Absent, "")
		b.defEffects(r, path, true, d)
	})
	if !fresh {
		b.defEffects(r, r.tablePath(), false, d)
	}
	return true
}

// alterTable follows ALTER TABLE, one action or a list of one kind of
// action. Each item is a part of its own, judged against the table as it
// was (walk.Passes): a check on a name another item changes is left out.
func (b *builder) alterTable(v *ast.AlterTableStmt) bool {
	r, ok := b.ref(v.Name, false)
	if !ok {
		return false
	}
	if r.where == elsewhere {
		return true
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	if v.Actions == nil {
		return false
	}
	var actions []*ast.AlterTableAction
	for _, item := range v.Actions.Items {
		a, ok := item.(*ast.AlterTableAction)
		if !ok {
			return false
		}
		actions = append(actions, a)
	}
	b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), notKind(r.written, "table", kindTable))
	known := b.knownTable(r)
	path := r.tablePath()
	switch {
	case r.where == unplaced:
		// The table may be any of its name; its children are not known
		// while its class is forgotten.
		b.forgetTable(r)
		b.markAllTablesDirty()
	case !known:
		b.forgetMeant(r, childSpaces...)
		b.markAllTablesDirty()
	}
	if b.alterTextFails(r, actions) {
		return true
	}
	for _, a := range actions {
		b.Next()
		if !b.alterAction(r, path, known, a) {
			b.Last = true
			return true
		}
	}
	return true
}

// alterTextFails checks what an ALTER TABLE's text alone settles: a name
// given twice, two primary keys, and names too long.
func (b *builder) alterTextFails(r objRef, actions []*ast.AlterTableAction) bool {
	d := &tableDef{}
	// A column and a constraint are one namespace here (8168).
	dropped := map[ast.AlterTableActionType]map[string]string{}
	for _, a := range actions {
		switch a.Type {
		case ast.ATAddColumn:
			if a.Column != nil {
				d.addColumn(a.Column)
			}
		case ast.ATAddIndex:
			if a.Index != nil {
				d.indexes = append(d.indexes, a.Index)
			}
		case ast.ATAddConstraint:
			if a.Constraint != nil {
				d.constraints = append(d.constraints, a.Constraint)
			}
		case ast.ATDropColumn, ast.ATDropConstraint, ast.ATDropIndex:
			nm, kind := a.ColName, "column"
			if a.Type == ast.ATDropConstraint {
				kind = "constraint"
				if a.Constraint != nil {
					nm = a.Constraint.Name
				}
			}
			if a.Type == ast.ATDropIndex {
				kind = "index"
				nm = a.IndexName
			}
			space := a.Type
			if space == ast.ATDropConstraint {
				space = ast.ATDropColumn
			}
			if dropped[space] == nil {
				dropped[space] = map[string]string{}
			}
			if prev, twice := dropped[space][nm]; twice {
				if prev != kind {
					kind = ""
				}
				b.fail(rangeOf(a.Loc), review.DuplicateName, msgTwice("drops", kind, ident(nm)))
				return true
			}
			dropped[space][nm] = kind
		}
	}
	return b.defTextFails(r, d, "adds")
}

// alterAction lowers one ALTER TABLE action, false for one the walk does
// not follow.
func (b *builder) alterAction(r objRef, path string, known bool, a *ast.AlterTableAction) bool {
	rng := rangeOf(a.Loc)
	switch a.Type {
	case ast.ATAddIndex, ast.ATDropIndex, ast.ATAddConstraint, ast.ATDropConstraint:
		// The table's indexes may change.
		b.markTableDirty(path)
	}
	switch a.Type {
	case ast.ATAddColumn:
		if a.Column == nil {
			return false
		}
		cn, ok := keyed(a.Column.Name)
		if !ok {
			return false
		}
		if known {
			b.lacks(spaceColumn, path, cn, note(rng, "", text(msgExistsIn("column", ident(cn.spelling), r.written))))
		}
		d := &tableDef{}
		d.addColumn(a.Column)
		b.defNameChecks(r, d)
		b.addConstraintChecks(r, path, known, d)
		// The column's expressions are bound once the table and the free
		// name are found (4902, 2705 first; measured, mssql/catalog).
		if known && b.see(spaceColumn, path, cn).fact == walk.Absent && systemTypes(d) && b.tableCollateFails(d, b.columnType(r)) {
			return true
		}
		b.defEffects(r, path, known, d)
		b.markTableDirty(path)
	case ast.ATAddIndex:
		// ADD INDEX name ...: the index must not already name one of the
		// table's indexes (1913), and its name then arrives in that space.
		if a.Index == nil {
			return false
		}
		in, ok := keyed(a.Index.Name)
		if !ok {
			return false
		}
		if known {
			b.lacks(spaceIndex, path, in, note(rng, "", text(msgExistsIn("index", ident(a.Index.Name), r.written))))
			b.checkIndexColList(r, path, a.Index.Columns, true)
			b.checkIndexColList(r, path, a.Index.IncludeCols, false)
		}
		b.defEffects(r, path, known, &tableDef{indexes: []*ast.InlineIndexDef{a.Index}})
	case ast.ATAddConstraint:
		if a.Constraint == nil {
			return false
		}
		c := a.Constraint
		if c.Type == ast.ConstraintEdge {
			if b.knownNonEdgeTable(r) {
				b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" is not a graph edge table.")
			}
			b.Last = true
			return true
		}
		if known {
			switch c.Type {
			case ast.ConstraintPrimaryKey, ast.ConstraintUnique:
				if dup, dupRng, hasDup := dupColumnInLists(c.Columns); hasDup {
					if dupRng.Start == 0 && dupRng.End == 0 {
						dupRng = rng
					}
					b.fail(dupRng, review.DuplicateName, msgTwice("lists", "column", ident(dup)))
					return true
				}
				b.checkIndexColList(r, path, c.Columns, true)
			case ast.ConstraintForeignKey:
				if c.Columns == nil || len(c.Columns.Items) == 0 {
					b.fail(rng, review.UnsupportedForm, "Table-level foreign key constraint must specify referencing columns.")
					return true
				}
				b.checkIndexColList(r, path, c.Columns, true)
			}
		}
		d := &tableDef{constraints: []*ast.ConstraintDef{c}}
		b.defNameChecks(r, d)
		b.addConstraintChecks(r, path, known, d)
		if b.constraintBinds(r, path, known, c) && b.tableCollateFails(d, b.columnType(r)) {
			return true
		}
		b.defEffects(r, path, known, d)
	case ast.ATDropIndex:
		// DROP INDEX name as a table element: the index must be one of the
		// table's, and it is gone afterwards. Only the named index leaves: the
		// other indexes of the table stay known. A key's index carries the
		// constraint's name and cannot be dropped on its own, as DROP INDEX
		// reports for the same name (IndexBacksConstraint).
		in, ok := keyed(a.IndexName)
		if !ok {
			return false
		}
		if !known {
			return true
		}
		b.needs(spaceIndex, path, in, nil, note(rng, msgMissingIn("index", ident(a.IndexName), r.written), nil))
		ixName := review.Code(ident(a.IndexName))
		b.ofKind(spaceIndex, path, in, []string{kindIndex, kindXMLIndex}, note(rng, "", text("Index "+ixName+" enforces constraint "+ixName+".")).as(review.IndexBacksConstraint))
		b.leave(spaceIndex, path, in)
	case ast.ATDropColumn:
		cn, ok := keyed(a.ColName)
		if !ok {
			return false
		}
		if !known {
			return true
		}
		s := b.see(spaceColumn, path, cn)
		switch {
		case !a.IfExists:
			b.needs(spaceColumn, path, cn, nil, note(rng, msgMissingIn("column", ident(cn.spelling), r.written), nil))
		case s.fact == walk.Unknown:
			b.maybe()
		}
		// A foreign key's column can't be dropped (5074), unless the
		// statement drops the key too, which the walk does not judge.
		if st := b.syncedTable(r); st != nil && st.fkColumns[a.ColName] != "" && b.singleAction() {
			b.fail(rng, review.UnsupportedForm, "Foreign key "+review.Code(ident(st.fkColumns[a.ColName]))+" depends on column "+review.Code(ident(a.ColName))+".")
			return true
		}
		b.leave(spaceColumn, path, cn)
		b.markTableDirty(path)
	case ast.ATDropConstraint:
		if a.Constraint == nil {
			return false
		}
		cn, ok := keyed(a.Constraint.Name)
		if !ok {
			return false
		}
		return b.dropConstraint(r, path, known, cn, a.IfExists, rng)
	case ast.ATAlterColumn, ast.ATAlterColumnAddDrop:
		cn, ok := keyed(a.ColName)
		if !ok {
			return false
		}
		if known {
			b.needs(spaceColumn, path, cn, nil, note(rng, msgMissingIn("column", ident(cn.spelling), r.written), nil))
		}
		if a.Type == ast.ATAlterColumnAddDrop {
			b.columnOptions(r, a, rng)
		}
		b.markTableDirty(path)
	case ast.ATSwitchPartition:
		if a.TargetName != nil {
			if tr, ok := b.ref(a.TargetName, false); ok && tr.where == placed {
				b.tableNeeds(tr, []string{kindTable}, msgMissing("table", tr.written), notKind(tr.written, "table", kindTable))
				if known && r.temp == tr.temp && r.schema.spelling == tr.schema.spelling && r.n.spelling == tr.n.spelling {
					b.fail(tr.rng, review.UnsupportedForm, "Cannot switch table "+review.Code(r.written)+" to itself.")
				}
			}
		}
		b.Last = true
	case ast.ATAddPeriod:
		if a.Names != nil && len(a.Names.Items) == 2 {
			s0, ok0 := a.Names.Items[0].(*ast.String)
			s1, ok1 := a.Names.Items[1].(*ast.String)
			if ok0 && ok1 && known {
				cn0, okCN0 := keyed(s0.Str)
				if okCN0 {
					b.needs(spaceColumn, path, cn0, nil, note(rng, msgMissingIn("column", ident(s0.Str), r.written), nil))
				}
				if cn1, ok := keyed(s1.Str); ok {
					b.needs(spaceColumn, path, cn1, nil, note(rng, msgMissingIn("column", ident(s1.Str), r.written), nil))
				}
				if okCN0 && s0.Str == s1.Str && b.see(spaceColumn, path, cn0).fact != walk.Absent {
					b.fail(rng, review.DuplicateName, msgTwice("lists", "column", ident(s0.Str)))
				}
			}
		}
		b.Last = true
	case ast.ATCheckConstraint, ast.ATNocheckConstraint:
		if a.Names == nil {
			return false
		}
		for _, item := range a.Names.Items {
			s, ok := item.(*ast.String)
			if !ok {
				return false
			}
			if strings.EqualFold(s.Str, "ALL") {
				continue
			}
			cn, ok := keyed(s.Str)
			if !ok {
				return false
			}
			if known {
				b.needs(spaceConstraint, path, cn, nil, note(rng, msgMissingIn("constraint", ident(cn.spelling), r.written), nil))
			}
		}
	case ast.ATEnableTrigger:
		b.Session(func() { b.w.triggers = true })
	case ast.ATDisableTrigger, ast.ATRebuild, ast.ATEnableChangeTracking, ast.ATDisableChangeTracking:
	case ast.ATSet:
		if a.Options != nil {
			for _, item := range a.Options.Items {
				s, ok := item.(*ast.String)
				if !ok {
					continue
				}
				opt := strings.ToUpper(s.Str)
				if opt == "FILETABLE_DIRECTORY" {
					if st := b.syncedTable(r); st != nil && !st.fileTable {
						b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" isn't a FileTable.")
						return true
					}
				}
				if strings.Contains(opt, "SYSTEM_VERSIONING") {
					if opt == "SYSTEM_VERSIONING" && b.periodless(r, rng) {
						return true
					}
					return false
				}
			}
		}
	case ast.ATDropPeriod:
		return b.periodless(r, rng)
	default:
		return false
	}
	return true
}

// periodless fails a statement that needs the table's SYSTEM_TIME period,
// or system versioning, which needs the period, on a synced table no
// column of which can be a period column (13510, 13591, 13593).
func (b *builder) periodless(r objRef, rng review.Range) bool {
	st := b.syncedTable(r)
	if st == nil || !st.noPeriod {
		return false
	}
	b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" has no `SYSTEM_TIME` period.")
	return true
}

// addConstraintChecks adds what an added column's or constraint's
// constraints need from the table: no second primary key (1779), and a
// PK's or UQ's name free among its indexes (1913).
func (b *builder) addConstraintChecks(r objRef, path string, known bool, d *tableDef) {
	if !known {
		return
	}
	for _, c := range d.constraints {
		switch c.Type {
		case ast.ConstraintDefault:
			// DEFAULT ... FOR a column the table has (1752) and that has
			// no default yet (1781, verified).
			if c.Columns != nil && len(c.Columns.Items) == 1 {
				if s, ok := c.Columns.Items[0].(*ast.String); ok {
					if cn, ok := keyed(s.Str); ok && !b.addsColumn(s.Str) {
						b.needs(spaceColumn, path, cn, nil, note(rangeOf(c.Loc), msgMissingIn("column", ident(s.Str), r.written), nil))
					}
					// The synced default holds while its constraint does.
					if st := b.syncedTable(r); st != nil && st.defaults[s.Str] != "" {
						if dn, ok := keyed(st.defaults[s.Str]); ok {
							if got := b.see(spaceConstraint, path, dn); got.fact == walk.Present {
								msg := "Column " + review.Code(ident(s.Str)) + " of " + review.Code(r.written) + " already has a default."
								b.Fail([]walk.Key{got.key}, walk.Present, &walkNote{rng: rangeOf(c.Loc), absent: msg, present: text(msg), fixed: review.UnsupportedForm})
							}
						}
					}
				}
			}
		case ast.ConstraintPrimaryKey:
			if f, _ := b.l.Lookup(primaryKey(path)); f == walk.Present {
				b.Need(primaryKey(path), walk.Absent, nil, note(rangeOf(c.Loc), "", text(review.Code(r.written)+" already has a primary key.")))
			}
			fallthrough
		case ast.ConstraintUnique:
			if cn, ok := keyed(c.Name); ok && c.Name != "" {
				b.lacks(spaceIndex, path, cn, note(rangeOf(c.Loc), "", text(msgExistsIn("index", ident(c.Name), r.written))))
			}
		}
	}
}

// dropConstraint follows DROP CONSTRAINT [IF EXISTS] n: n must be a
// constraint of the table (3728, 3733). The constraint's name leaves the
// schema, and a PK or UQ takes its index along; a constraint of unknown
// kind may have been either.
func (b *builder) dropConstraint(r objRef, path string, known bool, cn name, ifExists bool, rng review.Range) bool {
	kind := ""
	if known {
		s := b.see(spaceConstraint, path, cn)
		kind = s.kind
		switch {
		case !ifExists:
			b.needs(spaceConstraint, path, cn, nil, note(rng, msgMissingIn("constraint", ident(cn.spelling), r.written), nil))
		case s.fact == walk.Absent:
			// IF EXISTS skips a name that is no constraint of the table,
			// which another object of the schema may hold.
			return true
		case s.fact == walk.Unknown:
			b.maybe()
		}
	} else if ifExists {
		b.maybe()
	}
	if !r.temp {
		b.leaveObject(objRef{n: cn, schema: r.schema, parent: r.parent, where: r.where}, kind)
	}
	if !known {
		return true
	}
	b.leave(spaceConstraint, path, cn)
	switch kind {
	case kindPK:
		b.leave(spaceIndex, path, cn)
		b.Set(primaryKey(path), walk.Absent, "")
	case kindUQ:
		b.leave(spaceIndex, path, cn)
	case kindFK, kindCK, kindDF:
	default:
		b.MayRemove(walk.Scope{Space: spaceIndex, Parent: path, Prefix: cn.class})
		b.across(spaceIndex, path, cn, false)
		b.MayRemove(walk.Scope{Space: spacePrimary, Parent: path})
	}
	return true
}

// dropTables follows DROP TABLE [IF EXISTS] a, .... The items run in
// order and a missing one does not keep the others from being dropped
// (verified): each is a part judged against the state before the
// statement (walk.Passes), and a name written twice fails (3701). A table
// takes its children along, and its constraints' and triggers' names
// leave the schema, which the walk does not list: every object of the
// schema may be gone.
func (b *builder) dropTables(v *ast.DropStmt) bool {
	return b.dropEach(v, spaceObject, "table", func(r objRef) bool {
		present := notKind(r.written, "table", kindTable)
		if !v.IfExists {
			b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), present)
		} else if r.where == placed {
			space, parent := r.relSpace()
			b.ofKind(space, parent, r.n, []string{kindTable}, note(r.rng, "", present))
			if !r.temp {
				b.ofKind(spaceObject, r.parent, r.n, []string{kindTable}, note(r.rng, "", present))
			}
			if b.see(space, parent, r.n).fact == walk.Unknown {
				b.maybe()
			}
		}
		if len(v.Names.Items) == 1 {
			if from := b.referencer(r); from != "" {
				b.fail(r.rng, review.UnsupportedForm, msgReferenced(r.written, from))
				return true
			}
		}
		known := b.knownTable(r)
		path := r.tablePath()
		b.leaveObject(r, kindTable)
		if known {
			b.markTableDirty(path)
			b.clearChildren(path)
		} else if r.where != elsewhere {
			b.markAllTablesDirty()
		}
		if r.where == placed && !r.temp {
			b.MayRemove(walk.Scope{Space: spaceObject, Parent: r.parent})
		}
		if !r.temp && b.mayDropLedger(r, known) {
			// A ledger table's history table and ledger view, in any
			// schema, are renamed with it (verified).
			for _, space := range []string{spaceObject, spaceRelation} {
				b.MayRemove(walk.Scope{Space: space, Deep: true})
				b.MayAdd(walk.Scope{Space: space, Deep: true, Prefix: droppedLedger})
			}
		}
		if !r.temp && r.where != elsewhere && b.w.snap.cdc {
			// Change data capture's capture job drops a dropped table's
			// change table, functions, and procedures in schema cdc once
			// it reads the drop, before the next batch in auto-commit
			// mode (verified).
			for _, space := range []string{spaceObject, spaceRelation} {
				b.MayRemove(walk.Scope{Space: space, Parent: cdcSchema})
			}
		}
		return true
	})
}

// cdcSchema is the path of schema cdc, where change data capture keeps
// its objects.
var cdcSchema = func() string {
	n, _ := keyed("cdc")
	return walk.Path(n.class)
}()

// droppedLedger begins the class of the names a dropped ledger table and
// its history table and ledger view take: MSSQL_DroppedLedgerTable_...
var droppedLedger = classPrefix("MSSQL_DroppedLedger")

// mayDropLedger reports whether the table a reference means may be a
// synced ledger table: the one the walk knows it means, or any when the
// walk cannot tell.
func (b *builder) mayDropLedger(r objRef, known bool) bool {
	s := b.w.snap
	if s == nil || !s.ledger {
		return false
	}
	path := r.tablePath()
	if !known || s.ambiguous[path] {
		return true
	}
	t := s.tables[path]
	return t != nil && t.ledger
}

// dropEach lowers a DROP of a list of objects of a kind, one part per
// item, with the item's checks and effects from lower, which returns false
// for an item the walk does not follow. An unqualified name is looked up
// in the lookup space. A name written twice fails, unless IF EXISTS skips
// it.
func (b *builder) dropEach(v *ast.DropStmt, lookup, kind string, lower func(r objRef) bool) bool {
	if v.Names == nil {
		return false
	}
	seen := make(map[string]bool)
	for _, item := range v.Names.Items {
		tr, ok := item.(*ast.TableRef)
		if !ok {
			return false
		}
		b.Next()
		r, ok := b.refIn(tr, false, lookup)
		if !ok {
			return false
		}
		if b.refTooLong(tr, r) {
			return true
		}
		// The parts joined by a byte no name the walk keys holds, as
		// dbo.[x.y] and [dbo.x].y are two objects.
		written := strings.Join([]string{tr.Server, tr.Database, tr.Schema, tr.Object}, "\x00")
		if seen[written] && !v.IfExists {
			b.fail(r.rng, review.DuplicateName, msgTwice("drops", kind, r.written))
			return true
		}
		seen[written] = true
		if r.where == elsewhere {
			continue
		}
		if !lower(r) {
			return false
		}
	}
	return true
}

// truncate follows TRUNCATE TABLE, which needs a table (4701, 4708) and
// fires no trigger.
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	r, ok := b.ref(v.Table, false)
	if !ok {
		return false
	}
	if b.refTooLong(v.Table, r) {
		return true
	}
	b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), notKind(r.written, "table", kindTable))
	if from := b.referencer(r); from != "" {
		b.fail(b.anchor(), review.UnsupportedForm, msgReferenced(r.written, from))
	}
	return true
}

// msgReferenced says that another table's foreign key references a table.
func msgReferenced(table, from string) string {
	return "Table " + review.Code(table) + " is referenced by a foreign key of " + review.Code(from) + "."
}

// indexOn resolves the table of an index statement and adds the check
// that it is a table or view (1088 missing, 1914 another kind).
func (b *builder) indexOn(tr *ast.TableRef) (objRef, bool) {
	r, ok := b.ref(tr, false)
	if !ok {
		return r, false
	}
	if r.where == unplaced {
		b.forgetTable(r)
	}
	b.tableNeeds(r, []string{kindTable, kindView}, msgMissing("table", r.written), notKind(r.written, "table or view", kindTable, kindView))
	return r, true
}

// createIndex follows CREATE INDEX: the name must be free on the table
// (1913, which statistics share), or with DROP_EXISTING = ON exist
// (7999). A view's indexes are not synced, so only a table's are judged.
func (b *builder) createIndex(v *ast.CreateIndexStmt) bool {
	if dup, dupRng, hasDup := dupColumnInLists(v.Columns, v.IncludeCols, v.OrderCols); hasDup {
		if dupRng.Start == 0 && dupRng.End == 0 {
			dupRng = b.anchor()
		}
		b.fail(dupRng, review.DuplicateName, msgTwice("lists", "column", ident(dup)))
		return true
	}
	in, ok := keyed(v.Name)
	if !ok {
		return false
	}
	rebuild, ok := dropExisting(v.Options)
	if !ok {
		return false
	}
	r, ok := b.indexOn(v.Table)
	if !ok {
		return false
	}
	if b.refTooLong(v.Table, r) || b.tooLong(v.Name, false, b.anchor()) {
		return true
	}
	known := b.knownTable(r)
	if !known {
		if !rebuild && b.unboundView(r) {
			b.fail(b.anchor(), review.UnsupportedForm, "SQL Server doesn't accept an index on view "+review.Code(r.written)+", which isn't schema-bound.")
			return true
		}
		if r.where == placed {
			b.forgetMeant(r, spaceIndex)
		}
		b.markAllTablesDirty()
		return true
	}
	path := r.tablePath()
	rng := b.anchor()
	b.checkIndexColList(r, path, v.Columns, true)
	b.checkIndexColList(r, path, v.IncludeCols, false)
	b.checkIndexColList(r, path, v.OrderCols, true)
	b.markTableDirty(path)
	if rebuild {
		// The index is rebuilt under its name, and a PK's or UQ's keeps
		// its constraint.
		b.needs(spaceIndex, path, in, nil, note(rng, msgMissingIn("index", ident(v.Name), r.written), nil))
		return true
	}
	b.lacks(spaceIndex, path, in, note(rng, "", text(msgExistsIn("index", ident(v.Name), r.written))))
	// A table has one clustered index, rowstore or columnstore (1902,
	// 35372), which SQL Server checks once the columns are found (1911).
	if st := b.syncedTable(r); st != nil && st.clustered != "" && v.Clustered != nil && *v.Clustered &&
		b.see(spaceIndex, path, in).fact == walk.Absent && b.columnsPresent(path, v.Columns, v.IncludeCols, v.OrderCols) {
		b.fail(rng, review.UnsupportedForm, "Table "+review.Code(r.written)+" already has clustered index "+review.Code(ident(st.clustered))+".")
		return true
	}
	b.arrive(spaceIndex, path, in, kindIndex)
	return true
}

// createSpecializedIndex follows CREATE [PRIMARY|SELECTIVE] XML INDEX and
// CREATE SPATIAL INDEX: unlike ordinary indexes, these can only be created
// on a user table (1088 missing, 6334 non-table), their target column must
// exist (1911) and have a compatible type (6330 XML, 12002 spatial), a
// secondary XML index's USING XML INDEX parent must exist as an XML index
// (6333), and their index name must be free on the table (1913) or exist
// when DROP_EXISTING = ON (7999).
func (b *builder) createSpecializedIndex(name string, table *ast.TableRef, colName, usingIndex string, opts *ast.List, kind string, isXML, isSpatial bool) bool {
	in, ok := keyed(name)
	if !ok {
		return false
	}
	rebuild, ok := dropExisting(opts)
	if !ok {
		return false
	}
	r, ok := b.ref(table, false)
	if !ok {
		return false
	}
	if r.where == unplaced {
		b.forgetTable(r)
	}
	b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), notKind(r.written, "table", kindTable))
	if b.refTooLong(table, r) || b.tooLong(name, false, b.anchor()) {
		return true
	}
	known := b.knownTable(r)
	if !known {
		if r.where == placed {
			b.forgetMeant(r, spaceIndex)
		}
		b.markAllTablesDirty()
		return true
	}
	path := r.tablePath()
	rng := b.anchor()
	b.markTableDirty(path)
	if colName != "" && !isPseudoCol(colName) {
		if cn, ok := keyed(colName); ok {
			b.needs(spaceColumn, path, cn, nil, note(rng, msgMissingIn("column", ident(colName), r.written), nil))
		}
	}
	usingFails := false
	if usingIndex != "" {
		if pn, ok := keyed(usingIndex); ok {
			b.needs(spaceIndex, path, pn, []string{kindXMLIndex}, note(rng, msgMissingIn("index", ident(usingIndex), r.written), notKind(ident(usingIndex), "XML index", kindXMLIndex)))
			ps := b.see(spaceIndex, path, pn)
			if ps.fact == walk.Absent || (ps.fact == walk.Present && ps.kind != "" && ps.kind != kindXMLIndex) {
				usingFails = true
			}
		}
	}
	ixSight := b.see(spaceIndex, path, in)
	ixFails := (!rebuild && ixSight.fact == walk.Present) || (rebuild && ixSight.fact == walk.Absent)
	if colName != "" && !usingFails {
		if typ := b.knownColType(r, colName); typ != "" && isBuiltinSystemType(typ) {
			if isSpatial && typ != "geometry" && typ != "geography" {
				b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(colName))+" of type "+review.Code(typ)+" cannot be used in a spatial index.")
			} else if isXML && !ixFails && typ != "xml" {
				b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(colName))+" of type "+review.Code(typ)+" cannot be used in an XML index.")
			}
		}
	}
	if rebuild {
		b.needs(spaceIndex, path, in, nil, note(rng, msgMissingIn("index", ident(name), r.written), nil))
		return true
	}
	b.lacks(spaceIndex, path, in, note(rng, "", text(msgExistsIn("index", ident(name), r.written))))
	b.arrive(spaceIndex, path, in, kind)
	return true
}

// createStatistics follows CREATE STATISTICS: checks duplicate column names
// (1909), target table or view existence (1088, 1929), column existence and
// indexability (1911, 1977, 1978), and collision with an existing index
// name (1927), then marks the statistics name as May-added in spaceIndex.
func (b *builder) createStatistics(v *ast.CreateStatisticsStmt) bool {
	if dup, dupRng, hasDup := dupColumnInLists(v.Columns); hasDup {
		if dupRng.Start == 0 && dupRng.End == 0 {
			dupRng = b.anchor()
		}
		b.fail(dupRng, review.DuplicateName, msgTwice("lists", "column", ident(dup)))
		return true
	}
	sn, ok := keyed(v.Name)
	if !ok {
		return false
	}
	r, ok := b.ref(v.Table, false)
	if !ok {
		return false
	}
	if r.where == elsewhere {
		return true
	}
	if b.refTooLong(v.Table, r) || b.tooLong(v.Name, false, b.anchor()) {
		return true
	}
	b.tableNeeds(r, []string{kindTable, kindView}, msgMissing("table", r.written), notKind(r.written, "table or view", kindTable, kindView))
	if !b.knownTable(r) {
		if r.where == placed {
			b.forgetMeant(r, spaceIndex)
		}
		return true
	}
	path := r.tablePath()
	b.checkIndexColList(r, path, v.Columns, true)
	b.lacks(spaceIndex, path, sn, note(b.anchor(), "", text(msgExistsIn("index", ident(v.Name), r.written))))
	// Statistics share the index namespace (1913) but are not indexes
	// themselves, so record a MayAdd rather than arrive.
	b.MayAdd(walk.Scope{Space: spaceIndex, Parent: path, Prefix: sn.class})
	b.across(spaceIndex, path, sn, true)
	return true
}

// updateStatistics follows UPDATE STATISTICS: its target must exist as a
// table or indexed view (2706).
func (b *builder) updateStatistics(v *ast.UpdateStatisticsStmt) bool {
	r, ok := b.ref(v.Table, false)
	if !ok {
		return false
	}
	if r.where == elsewhere {
		return true
	}
	if b.refTooLong(v.Table, r) {
		return true
	}
	b.tableNeeds(r, []string{kindTable, kindView}, msgMissing("table", r.written), notKind(r.written, "table or view", kindTable, kindView))
	return true
}

// dropExisting reads the DROP_EXISTING option, which the parser gives as
// NAME=VALUE: the last one written counts (verified). ok is false for a
// value other than ON or OFF, which the server rejects (102, 153).
func dropExisting(opts *ast.List) (on, ok bool) {
	if opts == nil {
		return false, true
	}
	ok = true
	for _, item := range opts.Items {
		s, isString := item.(*ast.String)
		if !isString {
			continue
		}
		name, value, _ := strings.Cut(s.Str, "=")
		if !strings.EqualFold(name, "DROP_EXISTING") {
			continue
		}
		on, ok = strings.EqualFold(value, "ON"), strings.EqualFold(value, "ON") || strings.EqualFold(value, "OFF")
	}
	return on, ok
}

// alterIndex follows ALTER INDEX, which renames nothing: the table must
// exist (1088), a named index on it (2727), and FOR (ADD/REMOVE ...) path
// alterations require a selective XML index (6333).
func (b *builder) alterIndex(v *ast.AlterIndexStmt) bool {
	r, ok := b.indexOn(v.Table)
	if !ok {
		return false
	}
	if strings.EqualFold(v.IndexName, "ALL") || !b.knownTable(r) && !b.unboundView(r) {
		return true
	}
	in, ok := keyed(v.IndexName)
	if !ok {
		return false
	}
	var kinds []string
	var present func(string) string
	if strings.EqualFold(v.Action, "FOR") {
		kinds = []string{kindXMLIndex}
		present = notKind(ident(v.IndexName), "XML index", kindXMLIndex)
	}
	b.needs(spaceIndex, r.tablePath(), in, kinds, note(b.anchor(), msgMissingIn("index", ident(v.IndexName), r.written), present))
	return true
}

// dropIndexes follows DROP INDEX [IF EXISTS] i ON t, ..., and the old
// DROP INDEX t.i: the table must exist and hold the index (3701), which
// must not enforce a constraint, under IF EXISTS too (3723).
func (b *builder) dropIndexes(v *ast.DropStmt) bool {
	if v.Names == nil {
		return false
	}
	seen := make(map[string]bool)
	for i, item := range v.Names.Items {
		ix, ok := item.(*ast.TableRef)
		if !ok {
			return false
		}
		var on *ast.TableRef
		if v.OnTables != nil && i < len(v.OnTables.Items) {
			on, _ = v.OnTables.Items[i].(*ast.TableRef)
		}
		if on == nil {
			// DROP INDEX [schema.]table.index.
			if ix.Server != "" || ix.Schema == "" {
				return false
			}
			on = &ast.TableRef{Schema: ix.Database, Object: ix.Schema, Loc: ix.Loc}
		} else if ix.Schema != "" || ix.Database != "" || ix.Server != "" {
			return false
		}
		in, ok := keyed(ix.Object)
		if !ok {
			return false
		}
		b.Next()
		r, ok := b.ref(on, false)
		if !ok {
			return false
		}
		key := strings.Join([]string{on.Server, on.Database, on.Schema, on.Object, ix.Object}, "\x00")
		if seen[key] && !v.IfExists {
			b.fail(rangeOf(ix.Loc), review.DuplicateName, msgTwice("drops", "index", ident(ix.Object)))
			return true
		}
		seen[key] = true
		if r.where == elsewhere {
			continue
		}
		if r.where == unplaced {
			b.forgetTable(r)
			b.markAllTablesDirty()
			if v.IfExists {
				b.maybe()
			}
			continue
		}
		rng := rangeOf(ix.Loc)
		if !v.IfExists {
			b.tableNeeds(r, []string{kindTable, kindView}, msgMissingIn("index", ident(ix.Object), r.written), notKind(r.written, "table or view", kindTable, kindView))
		}
		known := b.knownTable(r) || b.unboundView(r)
		if !known {
			if v.IfExists {
				b.maybe()
			}
			b.forgetMeant(r, spaceIndex)
			b.markAllTablesDirty()
			continue
		}
		path := r.tablePath()
		b.markTableDirty(path)
		s := b.see(spaceIndex, path, in)
		switch {
		case !v.IfExists:
			b.needs(spaceIndex, path, in, nil, note(rng, msgMissingIn("index", ident(ix.Object), r.written), nil))
		case s.fact == walk.Unknown:
			b.maybe()
		}
		// Only a clustered index is dropped ONLINE (3745) or with MAXDOP or
		// MOVE TO (3748).
		if st := b.syncedTable(r); st != nil && len(v.Names.Items) == 1 && s.fact == walk.Present && st.nonclustered[in.class] {
			if opt := clusteredDropOption(v.Options); opt != "" {
				b.fail(rng, review.UnsupportedForm, "Index "+review.Code(ident(ix.Object))+" isn't clustered, so it can't be dropped with "+review.Code(opt)+".")
				return true
			}
		}
		// A key's index has the name of the constraint it enforces.
		ixName := review.Code(ident(ix.Object))
		b.ofKind(spaceIndex, path, in, []string{kindIndex, kindXMLIndex}, note(rng, "", text("Index "+ixName+" enforces constraint "+ixName+".")).as(review.IndexBacksConstraint))
		b.leave(spaceIndex, path, in)
		if s.fact != walk.Present || s.kind != kindIndex {
			// A primary XML index takes its secondary ones along.
			b.MayRemove(walk.Scope{Space: spaceIndex, Parent: path})
		}
	}
	return true
}

// clusteredDropOption returns the first DROP INDEX option only a clustered
// index takes: ONLINE = ON, MAXDOP, or MOVE TO.
func clusteredDropOption(opts *ast.List) string {
	if opts == nil {
		return ""
	}
	for _, item := range opts.Items {
		s, ok := item.(*ast.String)
		if !ok {
			continue
		}
		k, val, _ := strings.Cut(s.Str, "=")
		switch k = strings.ToUpper(strings.TrimSpace(k)); {
		case k == "ONLINE" && strings.EqualFold(strings.TrimSpace(val), "ON"):
			return "ONLINE = ON"
		case k == "MAXDOP", k == "MOVE TO":
			return k
		}
	}
	return ""
}

// maskingFunctions are the functions dynamic data masking takes.
var maskingFunctions = []string{"default(", "email(", "random(", "partial(", "datetime("}

// columnOptions checks ALTER COLUMN ... ADD or DROP of a column property:
// a masking function SQL Server has (16006), and, on a synced column of a
// system type, ROWGUIDCOL only on uniqueidentifier (2761, 4926), HIDDEN
// only on a datetime2 period column (13735), and PERSISTED only on a
// computed column (4919, verified).
func (b *builder) columnOptions(r objRef, a *ast.AlterTableAction, rng review.Range) {
	if a.Options == nil {
		return
	}
	typ := b.knownColType(r, a.ColName)
	system := typ != "" && isBuiltinSystemType(typ)
	for _, item := range a.Options.Items {
		o, ok := item.(*ast.AlterColumnOption)
		if !ok {
			continue
		}
		add := strings.EqualFold(o.Action, "ADD")
		switch opt := strings.ToUpper(o.Option); {
		case add && opt == "MASKED":
			fn := strings.ToLower(strings.TrimSpace(o.MaskFunction))
			known := false
			for _, m := range maskingFunctions {
				known = known || strings.HasPrefix(fn, m)
			}
			if !known {
				b.fail(rng, review.UnsupportedForm, "SQL Server doesn't accept the masking function "+review.Code(o.MaskFunction)+".")
				return
			}
		case opt == "NOT FOR REPLICATION" && typ != "" && !b.syncedTable(r).identity[a.ColName]:
			b.fail(rng, review.UnsupportedForm, "Column "+review.Code(ident(a.ColName))+" isn't an identity column.")
			return
		case opt == "ROWGUIDCOL" && system && typ != "uniqueidentifier",
			add && opt == "HIDDEN" && system && typ != "datetime2",
			add && opt == "PERSISTED" && typ != "":
			b.fail(rng, review.UnsupportedForm, "SQL Server doesn't accept "+review.Code(o.Action+" "+o.Option)+" on column "+review.Code(ident(a.ColName))+" of type "+review.Code(typ)+".")
			return
		}
	}
}

// singleAction reports whether the statement is an ALTER TABLE of one
// action.
func (b *builder) singleAction() bool {
	at, ok := b.s.node.(*ast.AlterTableStmt)
	return ok && at.Actions != nil && len(at.Actions.Items) == 1
}

// addsColumn reports whether the statement, an ALTER TABLE, may add a
// column of the name, which its other items see.
func (b *builder) addsColumn(col string) bool {
	at, ok := b.s.node.(*ast.AlterTableStmt)
	if !ok || at.Actions == nil {
		return true
	}
	for _, item := range at.Actions.Items {
		a, ok := item.(*ast.AlterTableAction)
		if !ok {
			return true
		}
		if a.Type == ast.ATAddColumn && (a.Column == nil || !differ(a.Column.Name, col)) {
			return true
		}
	}
	return false
}
