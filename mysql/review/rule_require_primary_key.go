package review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/review"
)

// checkRequirePrimaryKey reports two things, each only when the change
// ends that way: a table a CREATE TABLE made without a primary key, and a
// table whose primary key an ALTER TABLE drops without adding one back. A
// DROP PRIMARY KEY only runs on a table that has one, so neither needs the
// synced schema, and a table that had none before the change is never
// reported. A table made by CREATE TABLE ... LIKE, whose keys come from
// the table it copies, a CREATE TABLE IF NOT EXISTS, which may make
// nothing, and a temporary table are not judged, and neither is a
// statement inside a routine body, which runs only when called. CREATE
// TABLE ... AS SELECT copies no key from its query, so it is judged by
// what it declares.
//
// An ALTER TABLE that drops columns takes the primary key with it when it
// drops every column of the key, which only the synced schema tells. The
// rule returns, for checkColumnDrops to judge against every target, each
// top-level statement that drops something by name, and what the change
// does to tables and columns by name: which table a name means, and
// whether a statement runs, depend on the target.
//
// Without simulating the change, the rule follows the table by name. A
// later statement that may give it a key, drop it or its database, rename
// it, or create a table of its name settles it, and it is not reported.
func checkRequirePrimaryKey(sql string, stmts []statement, r *reporter) columnDrops {
	type pending struct {
		ref     *ast.TableRef
		index   int
		rng     review.Range
		message string
		settled bool
	}
	var list []*pending
	// byName holds the pending tables by name, which settle looks up.
	byName := make(map[string][]*pending)
	add := func(p *pending) {
		list = append(list, p)
		if p.ref != nil {
			byName[nameKey(p.ref.Name)] = append(byName[nameKey(p.ref.Name)], p)
		}
	}
	settle := func(ref *ast.TableRef) {
		if ref == nil {
			return
		}
		for _, p := range byName[nameKey(ref.Name)] {
			if !p.settled && sameTable(p.ref, ref) {
				p.settled = true
			}
		}
	}
	out := columnDrops{
		removals: make(map[[2]string][]removal), touched: make(seen), anyColumn: make(seen), gone: make(seen), settles: make(seen),
		indexes: make(seen), referenced: make(seen), foreignKeys: make(map[string]int),
	}
	// A statement's touches come at twice its position, before its own
	// column drops, at the next. A touch of a table other than the one
	// the statement makes or alters is one a foreign key references.
	touch := func(position int, own *ast.TableRef, touches ...columnTouch) {
		for _, t := range touches {
			if t.column == "" {
				out.anyColumn.add(t.table, "", position)
			} else {
				out.touched.add(t.table, t.column, position)
			}
			if t.table != own {
				out.referenced.add(t.table, "", position)
			}
		}
	}
	// named records the names of the indexes and foreign keys a statement
	// makes or renames.
	named := func(position int, ref *ast.TableRef, c *ast.Constraint) {
		if c == nil || c.Name == "" {
			return
		}
		if c.Type == ast.ConstrForeignKey {
			if _, ok := out.foreignKeys[nameKey(c.Name)]; !ok {
				out.foreignKeys[nameKey(c.Name)] = position
			}
			return
		}
		out.indexes.add(ref, c.Name, position)
	}
	// record adds a statement that drops by name, with what it drops.
	record := func(d dropStatement) {
		k := len(out.statements)
		out.statements = append(out.statements, d)
		removed := func(ref *ast.TableRef, kind dropKind, name string) {
			key := [2]string{nameKey(ref.Name), nameKey(name)}
			out.removals[key] = append(out.removals[key], removal{ref, kind, k})
		}
		for _, ref := range d.tables {
			if d.dropsTable {
				removed(ref, dropNone, "")
			}
			for _, cmd := range d.columns {
				removed(ref, dropColumn, cmd.Name)
			}
			for _, n := range d.named {
				removed(ref, n.kind, n.name)
			}
			for _, column := range d.plain {
				removed(ref, dropGeneration, column)
			}
		}
	}
	for i := range stmts {
		s := &stmts[i]
		switch v := s.node.(type) {
		case *ast.CreateTableStmt:
			out.gone.add(v.Table, "", i)
			touch(2*i, v.Table, createTouches(v)...)
			for _, c := range v.Constraints {
				named(2*i, v.Table, c)
			}
			// IF NOT EXISTS of a table the change made does nothing, and a
			// temporary table leaves the table of its name as it is.
			if !v.IfNotExists && !v.Temporary {
				settle(v.Table)
				out.settles.add(v.Table, "", i)
			}
			if v.Table == nil || v.Temporary || v.IfNotExists || v.Like != nil || createsPrimaryKey(v) {
				continue
			}
			add(&pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "creates table " + table(v.Table) + " without a primary key"})
		case *ast.AlterTableStmt:
			if v.Table == nil {
				continue
			}
			adds, renames := false, false
			var dropped *ast.AlterTableCmd
			d := dropStatement{statement: s.index, position: i, tables: []*ast.TableRef{v.Table}}
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				kind := dropKindOf(sql, cmd)
				switch kind {
				case dropColumn:
					d.columns = append(d.columns, cmd)
				case dropIndex, dropForeignKey, dropCheck, dropConstraint:
					d.named = append(d.named, namedDrop{kind, cmd.Name, cmd.IfExists})
				}
				if kind != dropColumn {
					touch(2*i, v.Table, alterTouches(sql, v.Table, cmd)...)
					d.note(cmd)
					if cmd.Type == ast.ATAddConstraint || cmd.Type == ast.ATAddIndex {
						named(2*i, v.Table, cmd.Constraint)
					}
					if cmd.Type == ast.ATRenameIndex {
						out.indexes.add(v.Table, cmd.Name, 2*i)
						out.indexes.add(v.Table, cmd.NewName, 2*i)
					}
				}
				switch {
				case addsPrimaryKey(cmd):
					adds = true
				case cmd.Type == ast.ATRenameTable:
					renames = true
					out.gone.add(renamedTo(v.Table, cmd), "", i)
					out.settles.add(renamedTo(v.Table, cmd), "", i)
				// DROP PRIMARY KEY, or DROP INDEX `PRIMARY`; a check or foreign
				// key quoted as PRIMARY is not the primary key.
				case dropped == nil && kind == dropPrimaryKey:
					dropped = cmd
				}
			}
			if adds || renames {
				settle(v.Table)
				out.settles.add(v.Table, "", i)
			}
			if renames {
				out.gone.add(v.Table, "", i)
			}
			if dropped != nil && !adds && !renames {
				add(&pending{ref: v.Table, index: s.index, rng: rangeOf(dropped.Loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
			}
			if len(d.columns) > 0 || len(d.named) > 0 || len(d.plain) > 0 {
				d.judged = len(d.columns) > 0 && dropped == nil && !adds && !renames
				record(d)
			}
			for _, cmd := range d.columns {
				touch(2*i+1, v.Table, columnTouch{v.Table, cmd.Name})
			}
		case *ast.CreateIndexStmt:
			if functionalParts(v.Columns) {
				touch(2*i, v.Table, columnTouch{table: v.Table})
			}
			out.indexes.add(v.Table, v.IndexName, 2*i)
		case *ast.DropIndexStmt:
			if !sameName(v.Name, "PRIMARY") {
				record(dropStatement{statement: s.index, position: i, tables: []*ast.TableRef{v.Table}, named: []namedDrop{{kind: dropIndex, name: v.Name}}})
				break
			}
			touch(2*i, v.Table, columnTouch{table: v.Table})
			add(&pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
		case *ast.DropTableStmt:
			// DROP TEMPORARY TABLE drops no permanent table, and leaves the
			// one of its name as it was.
			if v.Temporary {
				break
			}
			record(dropStatement{statement: s.index, position: i, tables: v.Tables, dropsTable: true, ifExists: v.IfExists})
			for _, ref := range v.Tables {
				out.gone.add(ref, "", i)
				out.settles.add(ref, "", i)
				settle(ref)
			}
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair != nil {
					for _, ref := range []*ast.TableRef{pair.Old, pair.New} {
						out.gone.add(ref, "", i)
						out.settles.add(ref, "", i)
						settle(ref)
					}
				}
			}
		case *ast.DropDatabaseStmt:
			out.databases = append(out.databases, v.Name)
			// A table without a database may be in the one dropped.
			for _, p := range list {
				if p.ref != nil && (p.ref.Schema == "" || sameName(p.ref.Schema, v.Name)) {
					p.settled = true
				}
			}
		}
	}
	for _, p := range list {
		if !p.settled {
			r.report(review.RequirePrimaryKey, p.index, p.rng, p.message)
		}
	}
	return out
}

// columnDrops is what the targets judge of a change's column drops: the
// top-level statements that drop something by name, in order, what they
// drop, what the change did to tables and columns, and the databases it
// drops.
type columnDrops struct {
	statements []dropStatement
	// removals holds what the statements drop, by table and object name;
	// a table a DROP TABLE drops is under the name "".
	removals map[[2]string][]removal
	// touched is the columns a statement named, or made something depend
	// on, and anyColumn the tables a statement touched any column, or the
	// primary key, of, each at twice the statement's position, its own
	// column drops at the next. gone is the tables a statement made,
	// dropped, or renamed, under each name, and settles those it may have
	// given a key, dropped, renamed, or made, at the statement's position.
	touched, anyColumn, gone, settles seen
	// indexes is the indexes a statement made or renamed, under each name,
	// by table, at twice the statement's position; foreignKeys the first
	// such position of each foreign key name, which a database holds once;
	// and referenced the tables a foreign key a statement made references.
	indexes, referenced seen
	foreignKeys         map[string]int
	databases           []string
}

// dropGeneration is a MODIFY or CHANGE COLUMN that gives a column a
// definition without a generation expression, which then reads nothing.
const dropGeneration dropKind = -1

// any reports whether the change has a column drop to judge.
func (c columnDrops) any() bool {
	return slices.ContainsFunc(c.statements, func(d dropStatement) bool { return d.judged })
}

// dropStatement is a top-level statement that drops something by name:
// an ALTER TABLE with a DROP subcommand of a column, an index, or a
// constraint, a DROP INDEX, or a DROP TABLE.
type dropStatement struct {
	// statement is the index in Result.Statements, and position the
	// statement's place in the parsed list.
	statement int
	position  int
	// tables are the tables it alters, or, with dropsTable, drops; with
	// ifExists, a table the database lacks is passed over.
	tables     []*ast.TableRef
	dropsTable bool
	ifExists   bool
	// columns are the DROP COLUMN subcommands, in order, and named the
	// drops of an index or a constraint.
	columns []*ast.AlterTableCmd
	named   []namedDrop
	// The other subcommands: adds are the columns they add, renames the
	// columns they give a new name, needs the columns they change, and
	// plain those they give a definition without a generation expression;
	// keys are the indexes and foreign keys they add, and indexRenames the
	// indexes they rename.
	adds, needs, plain    []string
	renames, indexRenames [][2]string
	keys                  []*ast.Constraint
	// judged is an ALTER TABLE whose column drops may take the primary
	// key: it drops columns, and no key, adds no key, and renames no
	// table.
	judged bool
}

// namedDrop is a DROP subcommand of an index or a constraint, of the kind
// dropKindOf tells.
type namedDrop struct {
	kind     dropKind
	name     string
	ifExists bool
}

// note records what an ALTER TABLE subcommand other than a drop adds,
// renames, or changes, which the table must have, or must not.
func (d *dropStatement) note(cmd *ast.AlterTableCmd) {
	switch cmd.Type {
	case ast.ATAddColumn:
		for _, col := range append([]*ast.ColumnDef{cmd.Column}, cmd.Columns...) {
			if col != nil {
				d.adds = append(d.adds, col.Name)
			}
		}
	case ast.ATAlterColumnDefault, ast.ATAlterColumnSetNotNull, ast.ATAlterColumnDropNotNull, ast.ATAlterColumnVisible, ast.ATAlterColumnInvisible:
		d.needs = append(d.needs, cmd.Name)
	case ast.ATModifyColumn, ast.ATChangeColumn, ast.ATRenameColumn:
		d.needs = append(d.needs, cmd.Name)
		if cmd.Type != ast.ATModifyColumn && !sameName(cmd.NewName, cmd.Name) {
			d.renames = append(d.renames, [2]string{cmd.Name, cmd.NewName})
		}
		if cmd.Column != nil && cmd.Column.Generated == nil {
			d.plain = append(d.plain, cmd.Name)
		}
	case ast.ATAddConstraint, ast.ATAddIndex:
		if c := cmd.Constraint; c != nil && c.Type != ast.ConstrPrimaryKey && c.Type != ast.ConstrCheck {
			d.keys = append(d.keys, c)
		}
	case ast.ATRenameIndex:
		d.indexRenames = append(d.indexRenames, [2]string{cmd.Name, cmd.NewName})
	}
}

// columnTouch is a column of a table a statement names, or may make
// something depend on: a key, a generated column, a check, or a foreign
// key. An empty column is any column of the table, or its primary key.
type columnTouch struct {
	table  *ast.TableRef
	column string
}

// removal is a column, index, or constraint of a table a statement of
// columnDrops.statements drops by name, of the kind dropKindOf tells;
// dropNone is the table itself, with all it holds.
type removal struct {
	table     *ast.TableRef
	kind      dropKind
	statement int
}

// seen records, by table name and a column of the table or "" for the
// table, the first and last positions at which a statement did something
// to a table of each spelling of the name, in the database it named or
// none.
type seen map[[2]string]map[tableName]span

type tableName struct{ schema, name string }

type span struct{ first, last int }

func (s seen) add(ref *ast.TableRef, column string, position int) {
	if ref == nil {
		return
	}
	k := [2]string{nameKey(ref.Name), nameKey(column)}
	if s[k] == nil {
		s[k] = make(map[tableName]span)
	}
	n := tableName{ref.Schema, ref.Name}
	sp, ok := s[k][n]
	if !ok {
		sp.first = position
	}
	sp.last = position
	s[k][n] = sp
}

// within returns the first and last positions at which a statement did
// something to a table that may be the one of the name, as match tells.
func (s seen) within(name, column string, match func(tableName) bool) (first, last int, ok bool) {
	for n, sp := range s[[2]string{nameKey(name), nameKey(column)}] {
		if !match(n) {
			continue
		}
		if !ok || sp.first < first {
			first = sp.first
		}
		if !ok || sp.last > last {
			last = sp.last
		}
		ok = true
	}
	return first, last, ok
}

// nameKey is the key two names sameName matches share.
func nameKey(name string) string {
	return strings.ToLower(strings.ToUpper(name))
}

// renamedTo is the table an ALTER TABLE ... RENAME subcommand renames to.
// A new name without a database stays in the table's.
func renamedTo(from *ast.TableRef, cmd *ast.AlterTableCmd) *ast.TableRef {
	to := &ast.TableRef{Name: cmd.NewName}
	if cmd.NewTable != nil {
		copied := *cmd.NewTable
		to = &copied
	}
	if to.Schema == "" {
		to.Schema = from.Schema
	}
	return to
}

// alterTouches returns what an ALTER TABLE subcommand does that a later
// DROP COLUMN of the table depends on: the columns it adds, drops,
// changes, or renames; the columns a foreign key it adds holds, and those
// it references; and any column, when it adds or drops a primary key,
// adds a generated column, a check, or a functional index, or partitions
// the table, any of which may read a column the server then refuses to
// drop.
func alterTouches(sql string, ref *ast.TableRef, cmd *ast.AlterTableCmd) []columnTouch {
	anyColumn := []columnTouch{{table: ref}}
	if addsPrimaryKey(cmd) || dropKindOf(sql, cmd) == dropPrimaryKey {
		return anyColumn
	}
	var out []columnTouch
	add := func(names ...string) {
		for _, name := range names {
			if name != "" {
				out = append(out, columnTouch{ref, name})
			}
		}
	}
	switch cmd.Type {
	case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
		add(cmd.Name, cmd.NewName)
		for _, col := range append([]*ast.ColumnDef{cmd.Column}, cmd.Columns...) {
			if col == nil {
				continue
			}
			if col.Generated != nil {
				return anyColumn
			}
			add(col.Name)
			for _, c := range col.Constraints {
				switch {
				case c == nil:
				case c.Type == ast.ColConstrCheck:
					return anyColumn
				case c.Type == ast.ColConstrReferences:
					out = append(out, referenceTouches(c.RefTable, c.RefColumns)...)
				}
			}
		}
	case ast.ATDropColumn, ast.ATRenameColumn, ast.ATAlterColumnDefault, ast.ATAlterColumnSetNotNull, ast.ATAlterColumnDropNotNull,
		ast.ATAlterColumnVisible, ast.ATAlterColumnInvisible:
		add(cmd.Name, cmd.NewName)
	case ast.ATAddConstraint, ast.ATAddIndex:
		c := cmd.Constraint
		switch {
		case c == nil:
		case c.Type == ast.ConstrCheck:
			return anyColumn
		case c.Type == ast.ConstrForeignKey:
			add(c.Columns...)
			out = append(out, referenceTouches(c.RefTable, c.RefColumns)...)
		case functionalParts(c.IndexColumns):
			return anyColumn
		}
	case ast.ATPartitionBy, ast.ATRemovePartitioning:
		return anyColumn
	}
	return out
}

// createTouches returns the columns of other tables the foreign keys of
// a CREATE TABLE reference.
func createTouches(v *ast.CreateTableStmt) []columnTouch {
	var out []columnTouch
	for _, c := range v.Constraints {
		if c != nil && c.Type == ast.ConstrForeignKey {
			out = append(out, referenceTouches(c.RefTable, c.RefColumns)...)
		}
	}
	for _, col := range v.Columns {
		if col == nil {
			continue
		}
		for _, c := range col.Constraints {
			if c != nil && c.Type == ast.ColConstrReferences {
				out = append(out, referenceTouches(c.RefTable, c.RefColumns)...)
			}
		}
	}
	return out
}

// referenceTouches returns the columns a foreign key references, any
// column of the table when it lists none.
func referenceTouches(ref *ast.TableRef, columns []string) []columnTouch {
	if ref == nil {
		return nil
	}
	if len(columns) == 0 {
		return []columnTouch{{table: ref}}
	}
	out := make([]columnTouch, 0, len(columns))
	for _, c := range columns {
		out = append(out, columnTouch{ref, c})
	}
	return out
}

// functionalParts reports whether a key part of an index is an expression
// rather than a column.
func functionalParts(parts []*ast.IndexColumn) bool {
	return slices.ContainsFunc(parts, func(p *ast.IndexColumn) bool {
		if p == nil {
			return false
		}
		_, column := p.Expr.(*ast.ColumnRef)
		return p.Functional || !column
	})
}

// checkColumnDrops judges the change's column drops against one target.
// A change that drops the target's database leaves none of its tables to
// judge: a drop before it fails, and the table it drops is gone after.
func checkColumnDrops(stmts []statement, c columnDrops, x schemaIndex, r *reporter) {
	if !c.any() || slices.ContainsFunc(c.databases, func(db string) bool { return x.same(db, x.schema.GetName()) }) {
		return
	}
	j := &judge{c: c, x: x, firstUse: beforeUse(stmts), runs: make([]int8, len(c.statements)), referencing: make(map[string][]referencingKey)}
	for _, s := range x.schema.GetSchemas() {
		for _, t := range s.GetTables() {
			for _, fk := range t.GetForeignKeys() {
				if db := fk.GetReferencedSchema(); db == "" || x.same(db, x.schema.GetName()) {
					k := nameKey(fk.GetReferencedTable())
					j.referencing[k] = append(j.referencing[k], referencingKey{t, fk})
				}
			}
		}
	}
	for i, d := range c.statements {
		if d.judged && j.statementRuns(i) {
			j.checkColumnDrop(d, r)
		}
	}
}

// judge answers for one target what the change's statements that drop by
// name do: runs memoizes, by statement, whether it is known to run, 1, or
// not, -1.
type judge struct {
	c           columnDrops
	x           schemaIndex
	firstUse    int
	runs        []int8
	referencing map[string][]referencingKey
}

// referencingKey is a foreign key of a table of the synced schema, by the
// table it references.
type referencingKey struct {
	table *metadata.TableMetadata
	fk    *metadata.ForeignKeyMetadata
}

// mayBe reports whether a table a statement named may be the target's
// table of the name: named alike, as the target compares names, in the
// target's database or in none, which a USE may have moved.
func (j *judge) mayBe(name string) func(tableName) bool {
	return func(n tableName) bool {
		return j.x.same(n.name, name) && (n.schema == "" || j.x.same(n.schema, j.x.schema.GetName()))
	}
}

// synced returns the synced table a reference of a statement certainly
// means, nil when a USE before it may have moved the current database, or
// an earlier statement made, dropped, or renamed a table it may mean.
func (j *judge) synced(ref *ast.TableRef, position int) *metadata.TableMetadata {
	t := j.x.table(ref, j.firstUse < position)
	if t == nil {
		return nil
	}
	if first, _, ok := j.c.gone.within(ref.Name, "", j.mayBe(ref.Name)); ok && first < position {
		return nil
	}
	return t
}

// touched reports whether a statement at a position up to the given one,
// doubled, touched a column of the target's table of the name.
func (j *judge) touched(name, column string, upTo int) bool {
	first, _, ok := j.c.touched.within(name, column, j.mayBe(name))
	return ok && first <= upTo
}

// statementRuns reports whether the statement is known to run on the
// target: each table it names is the synced one, and nothing the synced
// schema, or an earlier statement, tells keeps a subcommand from running.
// The server refuses an ALTER TABLE that drops a column the table does not
// have without IF EXISTS, or one twice, or every column of the table and
// adds none, or a column something depends on; that changes a column the
// table does not have, or one the statement drops, or makes a virtual
// generated column an ordinary one; that adds a column the table has and
// the statement keeps, or one twice; that adds an index under a name the
// table uses, or one twice, or a key on a column the table lacks once the
// statement ran, or a foreign key under a name the database uses, or one
// whose referenced columns lead no index; that renames an index the table
// does not have; or that drops an index or a constraint the table does
// not have without IF EXISTS. It refuses a DROP TABLE of a table the
// database lacks without IF EXISTS, or one a foreign key of another table
// references. Other failures, such as a type error or rows a key rejects,
// are not known before the statement runs.
func (j *judge) statementRuns(i int) bool {
	if j.runs[i] != 0 {
		return j.runs[i] > 0
	}
	j.runs[i] = -1
	if j.knownToRun(i) {
		j.runs[i] = 1
	}
	return j.runs[i] > 0
}

func (j *judge) knownToRun(i int) bool {
	d := j.c.statements[i]
	if d.dropsTable {
		return j.dropsTables(d, i)
	}
	ref := d.tables[0]
	t := j.synced(ref, d.position)
	if t == nil {
		return false
	}
	if first, _, ok := j.c.anyColumn.within(ref.Name, "", j.mayBe(ref.Name)); ok && first <= 2*d.position {
		return false
	}
	is := func(name string) func(string) bool { return func(n string) bool { return sameName(n, name) } }
	hasColumn := func(column string) bool {
		return slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return sameName(c.GetName(), column) })
	}
	drops := func(column string) bool {
		return slices.ContainsFunc(d.columns, func(cmd *ast.AlterTableCmd) bool { return sameName(cmd.Name, column) })
	}
	for k, cmd := range d.columns {
		twice := slices.ContainsFunc(d.columns[:k], func(prev *ast.AlterTableCmd) bool { return sameName(prev.Name, cmd.Name) })
		if !hasColumn(cmd.Name) && !cmd.IfExists || twice || j.touched(ref.Name, cmd.Name, 2*d.position) {
			return false
		}
	}
	// What an earlier statement touched, the synced schema no longer tells.
	for _, column := range d.needs {
		if !hasColumn(column) || drops(column) || j.touched(ref.Name, column, 2*d.position-1) {
			return false
		}
	}
	for _, column := range d.plain {
		for _, c := range t.GetColumns() {
			if g := c.GetGeneration(); sameName(c.GetName(), column) && g != nil && g.GetType() != metadata.GenerationMetadata_TYPE_STORED {
				return false
			}
		}
	}
	made := slices.Clone(d.adds)
	for _, r := range d.renames {
		made = append(made, r[1])
	}
	for k, column := range made {
		if hasColumn(column) && !drops(column) || slices.ContainsFunc(made[:k], is(column)) || j.touched(ref.Name, column, 2*d.position-1) {
			return false
		}
	}
	if len(d.columns) > 0 && len(d.adds) == 0 && !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return !drops(c.GetName()) }) {
		return false
	}
	for _, n := range d.named {
		if !n.ifExists && !j.holds(t, n, i) {
			return false
		}
	}
	// The columns and indexes the table has once the statement ran.
	renamedAway := func(column string) bool {
		return slices.ContainsFunc(d.renames, func(r [2]string) bool { return sameName(r[0], column) })
	}
	hasAfter := func(column string) bool {
		return hasColumn(column) && !drops(column) && !renamedAway(column) || slices.ContainsFunc(made, is(column))
	}
	droppedIndex := func(name string) bool {
		return slices.ContainsFunc(d.named, func(n namedDrop) bool {
			return (n.kind == dropIndex || n.kind == dropConstraint) && sameName(n.name, name)
		}) ||
			slices.ContainsFunc(d.indexRenames, func(r [2]string) bool { return sameName(r[0], name) })
	}
	syncedIndex := func(name string) bool {
		return slices.ContainsFunc(t.GetIndexes(), func(idx *metadata.IndexMetadata) bool { return sameName(idx.GetName(), name) })
	}
	indexTouched := func(name string) bool {
		first, _, ok := j.c.indexes.within(ref.Name, name, j.mayBe(ref.Name))
		return ok && first < 2*d.position
	}
	var indexes []string
	for _, r := range d.indexRenames {
		if sameName(r[0], "PRIMARY") || !syncedIndex(r[0]) || indexTouched(r[0]) || j.removed(t.GetName(), r[0], i, false, dropIndex, dropConstraint) {
			return false
		}
		indexes = append(indexes, r[1])
	}
	for _, c := range d.keys {
		for _, column := range c.Columns {
			if !hasAfter(column) || j.touched(ref.Name, column, 2*d.position-1) {
				return false
			}
		}
		if c.Type == ast.ConstrForeignKey {
			if !j.referencesIndex(c, t, d.position, i) || c.Name != "" && j.foreignKeyTaken(c.Name, t, d, i) {
				return false
			}
			continue
		}
		if c.Name != "" {
			indexes = append(indexes, c.Name)
		}
	}
	for k, name := range indexes {
		if syncedIndex(name) && !droppedIndex(name) || indexTouched(name) || slices.ContainsFunc(indexes[:k], is(name)) {
			return false
		}
	}
	for _, cmd := range d.columns {
		if j.keepsColumn(t, cmd.Name, drops, i) {
			return false
		}
	}
	return true
}

// dropsTables reports whether a DROP TABLE is known to run: each table it
// names is a synced one, or, under IF EXISTS, one the database certainly
// lacks; no foreign key a statement before made may reference one; and no
// foreign key of a table the statement keeps references one. MySQL drops
// the tables of the statement together.
func (j *judge) dropsTables(d dropStatement, i int) bool {
	var tables []*metadata.TableMetadata
	for _, ref := range d.tables {
		if ref == nil {
			return false
		}
		switch t := j.synced(ref, d.position); {
		case t != nil:
			tables = append(tables, t)
		case !d.ifExists || !j.lacks(ref, d.position):
			return false
		}
		if first, _, ok := j.c.referenced.within(ref.Name, "", j.mayBe(ref.Name)); ok && first < 2*d.position {
			return false
		}
	}
	for _, t := range tables {
		for _, r := range j.referencing[nameKey(t.GetName())] {
			if r.table == t || slices.Contains(tables, r.table) || !j.x.same(r.fk.GetReferencedTable(), t.GetName()) ||
				j.removed(r.table.GetName(), r.fk.GetName(), i, false, dropForeignKey, dropConstraint) {
				continue
			}
			return false
		}
	}
	return true
}

// lacks reports whether a reference of a statement certainly names no
// table or view of the target: it names the target's database, which no
// USE before it moved, the database holds no relation of the name, and no
// statement before it made, dropped, or renamed a table it may mean.
func (j *judge) lacks(ref *ast.TableRef, position int) bool {
	if ref.Schema == "" && j.firstUse < position || ref.Schema != "" && !j.x.same(ref.Schema, j.x.schema.GetName()) || j.x.hasRelation(ref.Name) {
		return false
	}
	first, _, ok := j.c.gone.within(ref.Name, "", j.mayBe(ref.Name))
	return !ok || first >= position
}

// referencesIndex reports whether the table a foreign key of a statement
// references is a synced one other than the table it alters, with an
// index that the referenced columns lead and no statement before dropped,
// and no statement before touched those columns or any of the table.
func (j *judge) referencesIndex(c *ast.Constraint, own *metadata.TableMetadata, position, i int) bool {
	if c.RefTable == nil || len(c.RefColumns) == 0 {
		return false
	}
	rt := j.synced(c.RefTable, position)
	if rt == nil || rt == own {
		return false
	}
	if first, _, ok := j.c.anyColumn.within(c.RefTable.Name, "", j.mayBe(c.RefTable.Name)); ok && first < 2*position {
		return false
	}
	for _, column := range c.RefColumns {
		if j.touched(c.RefTable.Name, column, 2*position-1) {
			return false
		}
	}
	return slices.ContainsFunc(rt.GetIndexes(), func(idx *metadata.IndexMetadata) bool {
		parts := idx.GetExpressions()
		if len(parts) < len(c.RefColumns) {
			return false
		}
		for k, column := range c.RefColumns {
			if !sameName(parts[k], column) {
				return false
			}
		}
		return !j.removed(rt.GetName(), idx.GetName(), i, false, dropIndex, dropConstraint)
	})
}

// foreignKeyTaken reports whether a foreign key name may name one the
// database holds when the statement runs: one of a synced table the
// statement does not drop, and no statement before dropped, or one a
// statement before made.
func (j *judge) foreignKeyTaken(name string, own *metadata.TableMetadata, d dropStatement, i int) bool {
	if first, ok := j.c.foreignKeys[nameKey(name)]; ok && first < 2*d.position {
		return true
	}
	for _, s := range j.x.schema.GetSchemas() {
		for _, t := range s.GetTables() {
			for _, fk := range t.GetForeignKeys() {
				if !sameName(fk.GetName(), name) || j.removed(t.GetName(), name, i, false, dropForeignKey, dropConstraint) {
					continue
				}
				if t == own && slices.ContainsFunc(d.named, func(n namedDrop) bool {
					return (n.kind == dropForeignKey || n.kind == dropConstraint) && sameName(n.name, name)
				}) {
					continue
				}
				return true
			}
		}
	}
	return false
}

// holds reports whether the synced table has an index or a constraint of
// the name a DROP subcommand drops, of a kind it drops, which no
// statement before the given one is known to have dropped.
func (j *judge) holds(t *metadata.TableMetadata, n namedDrop, i int) bool {
	before := func(kinds ...dropKind) bool { return j.removed(t.GetName(), n.name, i, false, kinds...) }
	index := slices.ContainsFunc(t.GetIndexes(), func(idx *metadata.IndexMetadata) bool { return !idx.GetPrimary() && sameName(idx.GetName(), n.name) })
	fk := slices.ContainsFunc(t.GetForeignKeys(), func(fk *metadata.ForeignKeyMetadata) bool { return sameName(fk.GetName(), n.name) })
	check := slices.ContainsFunc(t.GetCheckConstraints(), func(c *metadata.CheckConstraintMetadata) bool { return sameName(c.GetName(), n.name) })
	switch n.kind {
	case dropIndex:
		return index && !before(dropIndex, dropConstraint)
	case dropForeignKey:
		return fk && !before(dropForeignKey, dropConstraint)
	case dropCheck:
		return check && !before(dropCheck, dropConstraint)
	case dropConstraint:
		return (index || fk || check) && !before(dropIndex, dropForeignKey, dropCheck, dropConstraint)
	}
	return false
}

// removed reports whether a statement known to run, up to the given one
// or, with !upTo, before it, dropped the object of the name of the
// target's table, of one of the kinds, or the table. The given statement
// itself counts as running: what it drops is judged with it.
func (j *judge) removed(tableName, name string, i int, upTo bool, kinds ...dropKind) bool {
	for _, key := range [][2]string{{nameKey(tableName), nameKey(name)}, {nameKey(tableName), ""}} {
		for _, rm := range j.c.removals[key] {
			if rm.statement > i || rm.statement == i && !upTo || !j.x.same(rm.table.Name, tableName) {
				continue
			}
			if rm.kind != dropNone && (key[1] == "" || !slices.Contains(kinds, rm.kind)) {
				continue
			}
			if rm.statement == i || j.statementRuns(rm.statement) {
				return true
			}
		}
	}
	return false
}

// checkColumnDrop reports an ALTER TABLE, known to run, that drops every
// column of the primary key of the target's synced table: MySQL drops an
// index with the last of its columns, and a drop of only some leaves the
// key on the rest. A later statement that may give the table a key, drop
// it, rename it, or make a table of its name settles it. A foreign key of
// another database, which the synced schema does not hold, refuses the
// drop too.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html#alter-table-add-drop-column
func (j *judge) checkColumnDrop(d dropStatement, r *reporter) {
	ref := d.tables[0]
	key := primaryKeyColumns(j.synced(ref, d.position))
	if key == nil {
		return
	}
	dropsKey := func(cmd *ast.AlterTableCmd) bool {
		return slices.ContainsFunc(key, func(c string) bool { return sameName(c, cmd.Name) })
	}
	for _, column := range key {
		if !slices.ContainsFunc(d.columns, func(cmd *ast.AlterTableCmd) bool { return sameName(cmd.Name, column) }) {
			return
		}
	}
	if _, last, ok := j.c.settles.within(ref.Name, "", j.mayBe(ref.Name)); ok && last > d.position {
		return
	}
	at := d.columns[slices.IndexFunc(d.columns, dropsKey)]
	r.report(review.RequirePrimaryKey, d.statement, rangeOf(at.Loc), "removes the primary key of "+table(ref)+", and the change adds none back")
}

// primaryKeyColumns returns the columns of the synced table's primary
// key, nil when it has none, or a part of it names no column of the table.
func primaryKeyColumns(t *metadata.TableMetadata) []string {
	i := slices.IndexFunc(t.GetIndexes(), (*metadata.IndexMetadata).GetPrimary)
	if i < 0 {
		return nil
	}
	key := t.GetIndexes()[i].GetExpressions()
	for _, part := range key {
		if !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return sameName(c.GetName(), part) }) {
			return nil
		}
	}
	if len(key) == 0 {
		return nil
	}
	return key
}

// keepsColumn reports whether something the synced schema holds, and no
// statement up to the given one is known to have dropped, may keep the
// server from dropping a column of the table, given the columns the
// statement drops: a foreign key of the table that holds it, or one of the
// database that references it, or the table holding one; a generated
// column the statement keeps that may read it; a check that reads it and
// another column; a functional index part that may read it; or the table's
// partitioning or subpartitioning, by an expression that may read it or by
// KEY (), which reads the primary key. A check that reads the column alone
// goes with it. A foreign key names the table it references without its
// database, so one that may reference a table of another database counts.
func (j *judge) keepsColumn(t *metadata.TableMetadata, column string, drops func(string) bool, i int) bool {
	has := func(list []string) bool {
		return slices.ContainsFunc(list, func(c string) bool { return sameName(c, column) })
	}
	removed := func(table, name string, kinds ...dropKind) bool { return j.removed(table, name, i, true, kinds...) }
	for _, fk := range t.GetForeignKeys() {
		if has(fk.GetColumns()) && !removed(t.GetName(), fk.GetName(), dropForeignKey, dropConstraint) {
			return true
		}
	}
	for _, ref := range j.referencing[nameKey(t.GetName())] {
		if j.x.same(ref.fk.GetReferencedTable(), t.GetName()) && has(ref.fk.GetReferencedColumns()) && !removed(ref.table.GetName(), ref.fk.GetName(), dropForeignKey, dropConstraint) {
			return true
		}
	}
	for _, c := range t.GetColumns() {
		if g := c.GetGeneration(); g != nil && !drops(c.GetName()) && !removed(t.GetName(), c.GetName(), dropColumn, dropGeneration) && mentions(g.GetExpression(), column) {
			return true
		}
	}
	for _, c := range t.GetCheckConstraints() {
		if removed(t.GetName(), c.GetName(), dropCheck, dropConstraint) {
			continue
		}
		reads, ok := expressionColumns(c.GetExpression())
		if !ok && mentions(c.GetExpression(), column) || has(reads) && slices.ContainsFunc(reads, func(r string) bool { return !sameName(r, column) }) {
			return true
		}
	}
	for _, idx := range t.GetIndexes() {
		if removed(t.GetName(), idx.GetName(), dropIndex, dropConstraint) {
			continue
		}
		for _, part := range idx.GetExpressions() {
			if !sameName(part, column) && mentions(part, column) {
				return true
			}
		}
	}
	for _, p := range t.GetPartitions() {
		if p.GetExpression() == "" || mentions(p.GetExpression(), column) {
			return true
		}
		for _, sub := range p.GetSubpartitions() {
			if sub.GetExpression() == "" || mentions(sub.GetExpression(), column) {
				return true
			}
		}
	}
	return false
}

// expressionColumns parses an expression of the synced schema and returns
// the columns it reads. ok is false when it does not parse.
func expressionColumns(expr string) (columns []string, ok bool) {
	list, err := parser.Parse("SELECT " + expr)
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil, false
	}
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		if c, ok := n.(*ast.ColumnRef); ok {
			columns = append(columns, c.Column)
		}
		return true
	})
	return columns, true
}

// mentions reports whether an expression may name a column: the name,
// compared without case, as a whole word or quoted in backticks.
func mentions(expr, column string) bool {
	expr, column = strings.ToLower(expr), strings.ToLower(column)
	if column == "" {
		return false
	}
	if strings.Contains(expr, "`"+strings.ReplaceAll(column, "`", "``")+"`") {
		return true
	}
	word := func(b byte) bool { return isWordByte(b) || b >= 0x80 }
	for i := 0; ; {
		j := strings.Index(expr[i:], column)
		if j < 0 {
			return false
		}
		j += i
		end := j + len(column)
		if (j == 0 || !word(expr[j-1])) && (end == len(expr) || !word(expr[end])) {
			return true
		}
		i = j + 1
	}
}

// createsPrimaryKey reports whether a CREATE TABLE declares a primary key,
// as a table constraint or on a column.
func createsPrimaryKey(v *ast.CreateTableStmt) bool {
	for _, c := range v.Constraints {
		if c != nil && c.Type == ast.ConstrPrimaryKey {
			return true
		}
	}
	for _, col := range v.Columns {
		if columnIsPrimaryKey(col) {
			return true
		}
	}
	return false
}

// addsPrimaryKey reports whether an ALTER TABLE subcommand may give the
// table a primary key: ADD PRIMARY KEY, or a column added, modified, or
// changed with PRIMARY KEY.
func addsPrimaryKey(cmd *ast.AlterTableCmd) bool {
	switch cmd.Type {
	case ast.ATAddConstraint, ast.ATAddIndex:
		return cmd.Constraint != nil && cmd.Constraint.Type == ast.ConstrPrimaryKey
	case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
		if columnIsPrimaryKey(cmd.Column) {
			return true
		}
		for _, col := range cmd.Columns {
			if columnIsPrimaryKey(col) {
				return true
			}
		}
	}
	return false
}

func columnIsPrimaryKey(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	for _, c := range col.Constraints {
		if c != nil && c.Type == ast.ColConstrPrimaryKey {
			return true
		}
	}
	return false
}
