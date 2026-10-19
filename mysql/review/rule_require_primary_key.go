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
// rule returns each top-level one for checkColumnDrops to judge against
// every target, unless the statement also drops or adds a primary key or
// renames the table, or an earlier statement made, dropped, or renamed a
// table of its name; with it go what the change does to columns, what it
// drops by name, and the databases it drops.
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
		// drop is an ALTER TABLE whose DROP COLUMN the targets judge.
		drop *columnDrop
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
	out := columnDrops{touched: newSeen(), removals: make(map[[2]string][]removal)}
	// gone holds the tables an earlier statement made, dropped, or
	// renamed, under each name, and goneDatabases the databases one
	// dropped: a table of such a name may not be the synced one. anyColumn
	// holds the tables a statement touched any column, or the primary key,
	// of.
	gone, anyColumn := newSeen(), newSeen()
	goneDatabases := make(map[string]bool)
	isGone := func(ref *ast.TableRef) bool {
		_, made := gone.first(ref, "")
		// A table named without a database, the targets tell.
		return made || ref.Schema != "" && goneDatabases[nameKey(ref.Schema)]
	}
	// A statement's touches come at twice its position, before its own
	// column drops, at the next.
	touch := func(position int, touches ...columnTouch) {
		for _, t := range touches {
			if t.column == "" {
				anyColumn.add(t.table, "", position)
			} else {
				out.touched.add(t.table, t.column, position)
			}
		}
	}
	// remove records what a statement drops by name, of a table the change
	// has not made, dropped, or renamed before.
	remove := func(position int, ref *ast.TableRef, kind dropKind, name string) {
		if ref == nil || isGone(ref) {
			return
		}
		k := [2]string{nameKey(ref.Name), nameKey(name)}
		out.removals[k] = append(out.removals[k], removal{ref, kind, position})
	}
	for i := range stmts {
		s := &stmts[i]
		switch v := s.node.(type) {
		case *ast.CreateTableStmt:
			gone.add(v.Table, "", i)
			touch(2*i, createTouches(v)...)
			// IF NOT EXISTS of a table the change made does nothing, and a
			// temporary table leaves the table of its name as it is.
			if !v.IfNotExists && !v.Temporary {
				settle(v.Table)
			}
			if v.Table == nil || v.Temporary || v.IfNotExists || v.Like != nil || createsPrimaryKey(v) {
				continue
			}
			add(&pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "creates table " + table(v.Table) + " without a primary key"})
		case *ast.AlterTableStmt:
			if v.Table == nil {
				continue
			}
			adds, renames, addsColumn := false, false, false
			var dropped *ast.AlterTableCmd
			var drops []*ast.AlterTableCmd
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				kind := dropKindOf(sql, cmd)
				switch kind {
				case dropColumn:
					drops = append(drops, cmd)
				case dropIndex, dropForeignKey, dropCheck, dropConstraint:
					remove(i, v.Table, kind, cmd.Name)
				}
				if kind != dropColumn {
					touch(2*i, alterTouches(sql, v.Table, cmd)...)
				}
				switch {
				case addsPrimaryKey(cmd):
					adds = true
				case cmd.Type == ast.ATAddColumn:
					addsColumn = true
				case cmd.Type == ast.ATRenameTable:
					renames = true
				// DROP PRIMARY KEY, or DROP INDEX `PRIMARY`; a check or foreign
				// key quoted as PRIMARY is not the primary key.
				case dropped == nil && kind == dropPrimaryKey:
					dropped = cmd
				}
			}
			if adds || renames {
				settle(v.Table)
			}
			if dropped != nil && !adds && !renames {
				add(&pending{ref: v.Table, index: s.index, rng: rangeOf(dropped.Loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
			}
			if _, unsure := anyColumn.first(v.Table, ""); len(drops) > 0 && dropped == nil && !adds && !renames && !unsure && !isGone(v.Table) {
				add(&pending{ref: v.Table, drop: &columnDrop{statement: s.index, position: i, table: v.Table, drops: drops, addsColumn: addsColumn}})
			}
			for _, cmd := range drops {
				remove(i, v.Table, dropColumn, cmd.Name)
				touch(2*i+1, columnTouch{v.Table, cmd.Name})
			}
			for _, cmd := range v.Commands {
				if cmd != nil && cmd.Type == ast.ATRenameTable {
					gone.add(v.Table, "", i)
					gone.add(renamedTo(v.Table, cmd), "", i)
				}
			}
		case *ast.CreateIndexStmt:
			if functionalParts(v.Columns) {
				touch(2*i, columnTouch{table: v.Table})
			}
		case *ast.DropIndexStmt:
			if !sameName(v.Name, "PRIMARY") {
				remove(i, v.Table, dropIndex, v.Name)
				break
			}
			touch(2*i, columnTouch{table: v.Table})
			add(&pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
		case *ast.DropTableStmt:
			// DROP TEMPORARY TABLE drops no permanent table, and leaves the
			// one of its name as it was.
			if v.Temporary {
				break
			}
			for _, ref := range v.Tables {
				remove(i, ref, dropNone, "")
				gone.add(ref, "", i)
				settle(ref)
			}
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair != nil {
					gone.add(pair.Old, "", i)
					gone.add(pair.New, "", i)
					settle(pair.Old)
					settle(pair.New)
				}
			}
		case *ast.DropDatabaseStmt:
			goneDatabases[nameKey(v.Name)] = true
			out.databases = append(out.databases, v.Name)
			// A table without a database may be in the one dropped; for a
			// column drop, the targets tell.
			for _, p := range list {
				if p.ref != nil && (p.ref.Schema == "" && p.drop == nil || p.ref.Schema != "" && sameName(p.ref.Schema, v.Name)) {
					p.settled = true
				}
			}
		}
	}
	for _, p := range list {
		switch {
		case p.settled:
		case p.drop != nil:
			out.list = append(out.list, *p.drop)
		default:
			r.report(review.RequirePrimaryKey, p.index, p.rng, p.message)
		}
	}
	return out
}

// columnDrops is what the targets judge of a change's column drops: the
// ALTER TABLE statements, what the change did to columns, what it drops by
// name, which a drop no longer depends on, and the databases it drops,
// which may hold a table named without one.
type columnDrops struct {
	list []columnDrop
	// touched is the columns a statement named, or made something depend
	// on, at twice the statement's position, its own column drops at the
	// next.
	touched   seen
	removals  map[[2]string][]removal
	databases []string
}

// columnDrop is a top-level ALTER TABLE that drops columns of a table the
// change did not make, drop, or rename before it, nor touch any column or
// the primary key of. Whether it drops the table's primary key depends on
// the key's columns, which the target's synced schema holds.
type columnDrop struct {
	// statement is the index in Result.Statements, and position the
	// statement's place in the parsed list.
	statement int
	position  int
	table     *ast.TableRef
	// drops are the statement's DROP COLUMN subcommands, in order, and
	// addsColumn whether it adds a column too.
	drops      []*ast.AlterTableCmd
	addsColumn bool
}

// columnTouch is a column of a table a statement names, or may make
// something depend on: a key, a generated column, a check, or a foreign
// key. An empty column is any column of the table, or its primary key.
type columnTouch struct {
	table  *ast.TableRef
	column string
}

// removal is a column, index, or constraint of a table a top-level
// statement drops by name, of the kind dropKindOf tells; dropNone is the
// table itself, with all it holds.
type removal struct {
	table    *ast.TableRef
	kind     dropKind
	position int
}

// seen records the first position at which a statement did something to
// a table of a name, and to a column of it or, with "", to the table.
// first answers for a reference the earliest of them that may have meant
// its table: one named without a database may mean a table of any.
type seen struct {
	byName   map[[2]string]int
	bySchema map[[3]string]int
}

func newSeen() seen {
	return seen{byName: make(map[[2]string]int), bySchema: make(map[[3]string]int)}
}

func (s seen) add(ref *ast.TableRef, column string, position int) {
	if ref == nil {
		return
	}
	k := [2]string{nameKey(ref.Name), nameKey(column)}
	if _, ok := s.byName[k]; !ok {
		s.byName[k] = position
	}
	ks := [3]string{k[0], k[1], nameKey(ref.Schema)}
	if _, ok := s.bySchema[ks]; !ok {
		s.bySchema[ks] = position
	}
}

func (s seen) first(ref *ast.TableRef, column string) (int, bool) {
	k := [2]string{nameKey(ref.Name), nameKey(column)}
	if ref.Schema == "" {
		p, ok := s.byName[k]
		return p, ok
	}
	p, ok := s.bySchema[[3]string{k[0], k[1], ""}]
	if q, found := s.bySchema[[3]string{k[0], k[1], nameKey(ref.Schema)}]; found && (!ok || q < p) {
		p, ok = q, true
	}
	return p, ok
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
// A change that drops the target's database leaves no table named without
// a database to judge: before a drop, the drop fails, and after, the
// table is gone.
func checkColumnDrops(stmts []statement, c columnDrops, x schemaIndex, r *reporter) {
	firstUse := beforeUse(stmts)
	dropsTarget := slices.ContainsFunc(c.databases, func(db string) bool { return x.same(db, x.schema.GetName()) })
	for _, d := range c.list {
		if dropsTarget && d.table.Schema == "" {
			continue
		}
		// removed reports whether the change dropped, by the statement or
		// before it, the object of a table of the target, of one of the
		// kinds, or the table.
		removed := func(tableName, name string, kinds ...dropKind) bool {
			for _, key := range [][2]string{{nameKey(tableName), nameKey(name)}, {nameKey(tableName), ""}} {
				for _, rm := range c.removals[key] {
					inTarget := rm.table.Schema == "" && rm.position <= firstUse || rm.table.Schema != "" && x.same(rm.table.Schema, x.schema.GetName())
					if rm.position <= d.position && inTarget && x.same(rm.table.Name, tableName) && (rm.kind == dropNone || key[1] != "" && slices.Contains(kinds, rm.kind)) {
						return true
					}
				}
			}
			return false
		}
		touched := func(column string) bool {
			p, ok := c.touched.first(d.table, column)
			return ok && p <= 2*d.position
		}
		checkColumnDrop(d, firstUse < d.position, x, touched, removed, r)
	}
}

// checkColumnDrop reports an ALTER TABLE that drops every column of the
// primary key of the target's synced table: MySQL drops an index with the
// last of its columns, and a drop of only some leaves the key on the rest.
// The table is found as MySQL resolves it, and not at all when a USE
// before the statement may have moved the current database. The drop is
// not judged when an earlier statement, or another subcommand of the
// statement, touched a column of the key or one the statement drops, or
// when the server refuses it: it drops a column the table does not have
// without IF EXISTS, or one twice, or every column of the table and adds
// none, or something the synced schema holds, and the change has not
// dropped, depends on a column it drops. A foreign key of another
// database, which the synced schema does not hold, refuses it too.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html#alter-table-add-drop-column
func checkColumnDrop(d columnDrop, moved bool, x schemaIndex, touched func(string) bool, removed func(string, string, ...dropKind) bool, r *reporter) {
	t := x.table(d.table, moved)
	key := primaryKeyColumns(t)
	if key == nil {
		return
	}
	drops := func(column string) bool {
		return slices.ContainsFunc(d.drops, func(cmd *ast.AlterTableCmd) bool { return sameName(cmd.Name, column) })
	}
	for _, column := range key {
		if !drops(column) || touched(column) {
			return
		}
	}
	for i, cmd := range d.drops {
		synced := slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return sameName(c.GetName(), cmd.Name) })
		twice := slices.ContainsFunc(d.drops[:i], func(prev *ast.AlterTableCmd) bool { return sameName(prev.Name, cmd.Name) })
		if !synced && !cmd.IfExists || twice || touched(cmd.Name) {
			return
		}
	}
	if !d.addsColumn && !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return !drops(c.GetName()) }) {
		return
	}
	for _, cmd := range d.drops {
		if x.keepsColumn(t, cmd.Name, drops, removed) {
			return
		}
	}
	at := d.drops[slices.IndexFunc(d.drops, func(cmd *ast.AlterTableCmd) bool {
		return slices.ContainsFunc(key, func(c string) bool { return sameName(c, cmd.Name) })
	})]
	r.report(review.RequirePrimaryKey, d.statement, rangeOf(at.Loc), "removes the primary key of "+table(d.table)+", and the change adds none back")
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

// keepsColumn reports whether something the synced schema holds, and the
// change has not dropped, may keep the server from dropping a column of
// the table, given the columns the statement drops: a foreign key of the
// table that holds it, or one of the database that references it; a
// generated column the statement keeps that may read it; a check that
// reads it and another column; a functional index part that may read it;
// or the table's partitioning or subpartitioning, by an expression that
// may read it or by KEY (), which reads the primary key. A check that
// reads the column alone goes with it. A foreign key names the table it
// references without its database, so one that may reference a table of
// another database counts.
func (x schemaIndex) keepsColumn(t *metadata.TableMetadata, column string, drops func(string) bool, removed func(string, string, ...dropKind) bool) bool {
	has := func(list []string) bool {
		return slices.ContainsFunc(list, func(c string) bool { return sameName(c, column) })
	}
	for _, fk := range t.GetForeignKeys() {
		if has(fk.GetColumns()) && !removed(t.GetName(), fk.GetName(), dropForeignKey, dropConstraint) {
			return true
		}
	}
	for _, s := range x.schema.GetSchemas() {
		for _, other := range s.GetTables() {
			for _, fk := range other.GetForeignKeys() {
				db := fk.GetReferencedSchema()
				if (db == "" || x.same(db, x.schema.GetName())) && x.same(fk.GetReferencedTable(), t.GetName()) && has(fk.GetReferencedColumns()) &&
					!removed(other.GetName(), fk.GetName(), dropForeignKey, dropConstraint) {
					return true
				}
			}
		}
	}
	for _, c := range t.GetColumns() {
		if g := c.GetGeneration(); g != nil && !drops(c.GetName()) && !removed(t.GetName(), c.GetName(), dropColumn) && mentions(g.GetExpression(), column) {
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
