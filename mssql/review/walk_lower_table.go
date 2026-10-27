package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// tableNeeds adds the check that a reference names a relation of one of
// kinds: a table, or for some statements a view too.
func (b *builder) tableNeeds(r objRef, kinds []string, absent, present string) {
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
		return false
	}
	space, parent := r.relSpace()
	fresh, _ := b.classFact(r.n.classKey(space, parent))
	b.arriveObject(r, kindTable)
	path := r.tablePath()
	if fresh != walk.Absent {
		b.forgetChildren(path)
		return false
	}
	b.clearChildren(path)
	fill(path)
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

// keyedConstraint is a constraint the walk keeps: one of the five kinds.
func keyedConstraint(c *ast.ConstraintDef) bool {
	_, ok := walkConstraintKinds[c.Type]
	return ok
}

// defTextFails checks what a definition's text alone settles: no column,
// constraint, or index name twice in the statement (2705, 8168), at most
// one primary key (8110), no constraint named like the table (2714), and
// no identifier too long (103). It reports whether the statement
// certainly fails.
func (b *builder) defTextFails(r objRef, d *tableDef) bool {
	seen := make(map[string]bool)
	for _, c := range d.columns {
		if b.tooLong(c.Name, false, rangeOf(c.Loc)) {
			return true
		}
		if seen[c.Name] {
			b.fail(rangeOf(c.Loc), "column "+ident(c.Name)+" is named twice in the statement")
			return true
		}
		seen[c.Name] = true
	}
	// One statement may not name two constraints or indexes alike, of
	// any kinds: a CHECK, DEFAULT, or FOREIGN KEY and an inline INDEX of
	// one name fail too (8168, verified).
	names := make(map[string]bool)
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
		if names[c.Name] {
			b.fail(rangeOf(c.Loc), msgNameTwice(c.Name))
			return true
		}
		names[c.Name] = true
		if !r.temp && c.Name == r.n.spelling {
			b.fail(rangeOf(c.Loc), msgExists(ident(c.Name)))
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
		if names[i.Name] {
			b.fail(rangeOf(i.Loc), msgNameTwice(i.Name))
			return true
		}
		names[i.Name] = true
	}
	if pks > 1 {
		b.fail(b.anchor(), "the statement adds two primary keys to "+r.written)
		return true
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
		b.lacks(spaceObject, r.parent, cn, note(crng, "", msgExists(ident(c.Name))))
		b.lacks(spaceRelation, r.parent, cn, note(crng, "", msgExists(ident(c.Name))))
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
	b.tableNeeds(t, []string{kindTable}, "foreign key target "+t.written+" does not exist", "foreign key target "+t.written+" is not a table")
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
	if d.defaults > 0 {
		b.generated(r, path, known, ast.ConstraintDefault)
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
			b.generated(r, path, known, c.Type)
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
func (b *builder) generated(r objRef, path string, known bool, t ast.ConstraintType) {
	prefix := constraintPrefixes[t]
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
// the state before the statement.
func (b *builder) createTable(v *ast.CreateTableStmt) bool {
	if unsupportedTable(v) {
		return false
	}
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
	if b.defTextFails(r, d) {
		return true
	}
	b.schemaNeeds(r)
	b.lacksObject(r)
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
	b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), msgNotKind(r.written, "table"))
	known := b.knownTable(r)
	path := r.tablePath()
	if r.where == unplaced && !r.temp {
		// The table may be any of its name; its children are not known
		// while its class is forgotten.
		b.forgetName(spaceRelation, r.n)
	}
	if b.alterTextFails(r, actions) {
		return true
	}
	for _, a := range actions {
		b.Next()
		if !b.alterAction(r, path, known, a) {
			return false
		}
	}
	return true
}

// alterTextFails checks what an ALTER TABLE's text alone settles: a name
// given twice, two primary keys, and names too long.
func (b *builder) alterTextFails(r objRef, actions []*ast.AlterTableAction) bool {
	d := &tableDef{}
	dropped := map[ast.AlterTableActionType]map[string]bool{}
	for _, a := range actions {
		switch a.Type {
		case ast.ATAddColumn:
			if a.Column != nil {
				d.addColumn(a.Column)
			}
		case ast.ATAddConstraint:
			if a.Constraint != nil {
				d.constraints = append(d.constraints, a.Constraint)
			}
		case ast.ATDropColumn, ast.ATDropConstraint:
			nm := a.ColName
			if a.Type == ast.ATDropConstraint && a.Constraint != nil {
				nm = a.Constraint.Name
			}
			if dropped[a.Type] == nil {
				dropped[a.Type] = map[string]bool{}
			}
			if dropped[a.Type][nm] {
				b.fail(rangeOf(a.Loc), msgNameTwice(nm))
				return true
			}
			dropped[a.Type][nm] = true
		}
	}
	return b.defTextFails(r, d)
}

// alterAction lowers one ALTER TABLE action, false for one the walk does
// not follow.
func (b *builder) alterAction(r objRef, path string, known bool, a *ast.AlterTableAction) bool {
	rng := rangeOf(a.Loc)
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
			b.lacks(spaceColumn, path, cn, note(rng, "", "column "+ident(cn.spelling)+" already exists in "+r.written))
		}
		d := &tableDef{}
		d.addColumn(a.Column)
		b.defNameChecks(r, d)
		b.addConstraintChecks(r, path, known, d)
		b.defEffects(r, path, known, d)
	case ast.ATAddConstraint:
		if a.Constraint == nil {
			return false
		}
		d := &tableDef{constraints: []*ast.ConstraintDef{a.Constraint}}
		b.defNameChecks(r, d)
		b.addConstraintChecks(r, path, known, d)
		b.defEffects(r, path, known, d)
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
			b.needs(spaceColumn, path, cn, nil, note(rng, "column "+ident(cn.spelling)+" does not exist in "+r.written, ""))
		case s.fact == walk.Unknown:
			b.maybe()
		}
		b.leave(spaceColumn, path, cn)
	case ast.ATDropConstraint:
		if a.Constraint == nil {
			return false
		}
		cn, ok := keyed(a.Constraint.Name)
		if !ok {
			return false
		}
		return b.dropConstraint(r, path, known, cn, a.IfExists, rng)
	case ast.ATAlterColumn:
		cn, ok := keyed(a.ColName)
		if !ok {
			return false
		}
		if known {
			b.needs(spaceColumn, path, cn, nil, note(rng, "column "+ident(cn.spelling)+" does not exist in "+r.written, ""))
		}
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
				b.needs(spaceConstraint, path, cn, nil, note(rng, ident(cn.spelling)+" is not a constraint of "+r.written, ""))
			}
		}
	case ast.ATEnableTrigger:
		b.Session(func() { b.w.triggers = true })
	case ast.ATDisableTrigger, ast.ATRebuild, ast.ATEnableChangeTracking, ast.ATDisableChangeTracking, ast.ATAlterColumnAddDrop:
	case ast.ATSet:
		if a.Options != nil {
			for _, item := range a.Options.Items {
				if s, ok := item.(*ast.String); ok && strings.Contains(strings.ToUpper(s.Str), "SYSTEM_VERSIONING") {
					return false
				}
			}
		}
	default:
		return false
	}
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
		case ast.ConstraintPrimaryKey:
			if f, _ := b.l.Lookup(primaryKey(path)); f == walk.Present {
				b.Need(primaryKey(path), walk.Absent, nil, note(rangeOf(c.Loc), "", "table "+r.written+" already has a primary key"))
			}
			fallthrough
		case ast.ConstraintUnique:
			if cn, ok := keyed(c.Name); ok && c.Name != "" {
				b.lacks(spaceIndex, path, cn, note(rangeOf(c.Loc), "", "index "+ident(c.Name)+" already exists on "+r.written))
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
			b.needs(spaceConstraint, path, cn, nil, note(rng, ident(cn.spelling)+" is not a constraint of "+r.written, ""))
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
	return b.dropEach(v, func(r objRef) bool {
		if !v.IfExists {
			b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), msgNotKind(r.written, "table"))
		} else if r.where == placed {
			space, parent := r.relSpace()
			b.ofKind(space, parent, r.n, []string{kindTable}, note(r.rng, "", msgNotKind(r.written, "table")))
			if !r.temp {
				b.ofKind(spaceObject, r.parent, r.n, []string{kindTable}, note(r.rng, "", msgNotKind(r.written, "table")))
			}
			if b.see(space, parent, r.n).fact == walk.Unknown {
				b.maybe()
			}
		}
		known := b.knownTable(r)
		path := r.tablePath()
		b.leaveObject(r, kindTable)
		if known {
			b.clearChildren(path)
		}
		if r.where == placed && !r.temp {
			b.MayRemove(walk.Scope{Space: spaceObject, Parent: r.parent})
		}
		return true
	})
}

// dropEach lowers a DROP of a list, one part per item, with the item's
// checks and effects from lower, which returns false for an item the walk
// does not follow. A name written twice fails, unless IF EXISTS skips it.
func (b *builder) dropEach(v *ast.DropStmt, lower func(r objRef) bool) bool {
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
		r, ok := b.ref(tr, false)
		if !ok {
			return false
		}
		if b.refTooLong(tr, r) {
			return true
		}
		written := strings.Join([]string{tr.Server, tr.Database, tr.Schema, tr.Object}, ".")
		if seen[written] && !v.IfExists {
			b.fail(r.rng, r.written+" is dropped twice")
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
	b.tableNeeds(r, []string{kindTable}, msgMissing("table", r.written), msgNotKind(r.written, "table"))
	return true
}

// indexOn resolves the table of an index statement and adds the check
// that it is a table or view (1088 missing, 1914 another kind).
func (b *builder) indexOn(tr *ast.TableRef) (objRef, bool) {
	r, ok := b.ref(tr, false)
	if !ok {
		return r, false
	}
	if r.where == unplaced && !r.temp {
		b.forgetName(spaceRelation, r.n)
	}
	b.tableNeeds(r, []string{kindTable, kindView}, msgMissing("table", r.written), msgNotKind(r.written, "table or view"))
	return r, true
}

// createIndex follows CREATE INDEX: the name must be free on the table
// (1913, which statistics share), or with DROP_EXISTING = ON exist
// (7999). A view's indexes are not synced, so only a table's are judged.
func (b *builder) createIndex(v *ast.CreateIndexStmt) bool {
	in, ok := keyed(v.Name)
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
		return true
	}
	path := r.tablePath()
	rng := b.anchor()
	if dropExisting(v.Options) {
		// The index is rebuilt under its name, and a PK's or UQ's keeps
		// its constraint.
		b.needs(spaceIndex, path, in, nil, note(rng, "index "+ident(v.Name)+" does not exist on "+r.written, ""))
		return true
	}
	b.lacks(spaceIndex, path, in, note(rng, "", "index "+ident(v.Name)+" already exists on "+r.written))
	b.arrive(spaceIndex, path, in, kindIndex)
	return true
}

func dropExisting(opts *ast.List) bool {
	if opts == nil {
		return false
	}
	for _, item := range opts.Items {
		if s, ok := item.(*ast.String); ok && strings.EqualFold(strings.ReplaceAll(s.Str, " ", ""), "DROP_EXISTING=ON") {
			return true
		}
	}
	return false
}

// alterIndex follows ALTER INDEX, which renames nothing: the table must
// exist (1088), and a named index on it (2727).
func (b *builder) alterIndex(v *ast.AlterIndexStmt) bool {
	r, ok := b.indexOn(v.Table)
	if !ok {
		return false
	}
	if strings.EqualFold(v.IndexName, "ALL") || !b.knownTable(r) {
		return true
	}
	in, ok := keyed(v.IndexName)
	if !ok {
		return false
	}
	b.needs(spaceIndex, r.tablePath(), in, nil, note(b.anchor(), "index "+ident(v.IndexName)+" does not exist on "+r.written, ""))
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
		key := strings.Join([]string{on.Database, on.Schema, on.Object, ix.Object}, ".")
		if seen[key] && !v.IfExists {
			b.fail(rangeOf(ix.Loc), "index "+ident(ix.Object)+" is dropped twice")
			return true
		}
		seen[key] = true
		if r.where == elsewhere {
			continue
		}
		if r.where == unplaced && !r.temp {
			b.forgetName(spaceRelation, r.n)
			if v.IfExists {
				b.maybe()
			}
			continue
		}
		rng := rangeOf(ix.Loc)
		if !v.IfExists {
			b.tableNeeds(r, []string{kindTable, kindView}, "index "+ident(ix.Object)+" does not exist on "+r.written, msgNotKind(r.written, "table or view"))
		}
		known := b.knownTable(r)
		if !known {
			if v.IfExists {
				b.maybe()
			}
			continue
		}
		path := r.tablePath()
		s := b.see(spaceIndex, path, in)
		switch {
		case !v.IfExists:
			b.needs(spaceIndex, path, in, nil, note(rng, "index "+ident(ix.Object)+" does not exist on "+r.written, ""))
		case s.fact == walk.Unknown:
			b.maybe()
		}
		b.ofKind(spaceIndex, path, in, []string{kindIndex, kindXMLIndex}, note(rng, "", "index "+ident(ix.Object)+" of "+r.written+" enforces a constraint"))
		b.leave(spaceIndex, path, in)
		if s.fact != walk.Present || s.kind != kindIndex {
			// A primary XML index takes its secondary ones along.
			b.MayRemove(walk.Scope{Space: spaceIndex, Parent: path})
		}
	}
	return true
}
