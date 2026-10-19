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
// rule returns each top-level one for checkColumnDrop to judge against
// every target, unless the statement also drops or adds a primary key or
// renames the table, or an earlier statement made, dropped, or renamed a
// table of its name.
//
// Without simulating the change, the rule follows the table by name. A
// later statement that may give it a key, drop it or its database, rename
// it, or create a table of its name settles it, and it is not reported.
func checkRequirePrimaryKey(sql string, stmts []statement, r *reporter) []columnDrop {
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
	settle := func(ref *ast.TableRef) {
		for _, p := range list {
			if !p.settled && sameTable(p.ref, ref) {
				p.settled = true
			}
		}
	}
	// gone are the tables an earlier statement made, dropped, or renamed,
	// under each name, and goneDatabases the databases one dropped: a
	// table of such a name may not be the synced one. touched is what
	// earlier statements did to columns that a DROP COLUMN depends on.
	var gone []*ast.TableRef
	var goneDatabases []string
	var touched []columnTouch
	isGone := func(ref *ast.TableRef) bool {
		return slices.ContainsFunc(gone, func(t *ast.TableRef) bool { return sameTable(t, ref) }) ||
			slices.ContainsFunc(goneDatabases, func(db string) bool { return ref.Schema == "" || sameName(ref.Schema, db) })
	}
	for i := range stmts {
		s := &stmts[i]
		switch v := s.node.(type) {
		case *ast.CreateTableStmt:
			gone = append(gone, v.Table)
			touched = append(touched, createTouches(v)...)
			// IF NOT EXISTS of a table the change made does nothing, and a
			// temporary table leaves the table of its name as it is.
			if !v.IfNotExists && !v.Temporary {
				settle(v.Table)
			}
			if v.Table == nil || v.Temporary || v.IfNotExists || v.Like != nil || createsPrimaryKey(v) {
				continue
			}
			list = append(list, &pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "creates table " + table(v.Table) + " without a primary key"})
		case *ast.AlterTableStmt:
			if v.Table == nil {
				continue
			}
			adds, renames := false, false
			var dropped *ast.AlterTableCmd
			var drops []*ast.AlterTableCmd
			// own is what the statement's other subcommands do to columns.
			var own []columnTouch
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				if cmd.Type == ast.ATDropColumn {
					drops = append(drops, cmd)
				} else {
					own = append(own, alterTouches(sql, v.Table, cmd)...)
				}
				switch {
				case addsPrimaryKey(cmd):
					adds = true
				case cmd.Type == ast.ATRenameTable:
					renames = true
					gone = append(gone, v.Table, renamedTo(v.Table, cmd))
				// DROP PRIMARY KEY, or DROP INDEX `PRIMARY`; a check or foreign
				// key quoted as PRIMARY is not the primary key.
				case dropped == nil && dropKindOf(sql, cmd) == dropPrimaryKey:
					dropped = cmd
				}
			}
			if adds || renames {
				settle(v.Table)
			}
			if dropped != nil && !adds && !renames {
				list = append(list, &pending{ref: v.Table, index: s.index, rng: rangeOf(dropped.Loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
			}
			if len(drops) > 0 && dropped == nil && !adds && !renames && !isGone(v.Table) {
				if d, ok := newColumnDrop(s.index, i, v.Table, drops, append(slices.Clone(touched), own...)); ok {
					list = append(list, &pending{ref: v.Table, drop: d})
				}
			}
			touched = append(touched, own...)
			for _, cmd := range drops {
				touched = append(touched, columnTouch{v.Table, cmd.Name})
			}
		case *ast.CreateIndexStmt:
			if functionalParts(v.Columns) {
				touched = append(touched, columnTouch{table: v.Table})
			}
		case *ast.DropIndexStmt:
			if sameName(v.Name, "PRIMARY") {
				touched = append(touched, columnTouch{table: v.Table})
				list = append(list, &pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
			}
		case *ast.DropTableStmt:
			gone = append(gone, v.Tables...)
			// DROP TEMPORARY TABLE drops no permanent table.
			if v.Temporary {
				break
			}
			for _, ref := range v.Tables {
				settle(ref)
			}
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair != nil {
					gone = append(gone, pair.Old, pair.New)
					settle(pair.Old)
					settle(pair.New)
				}
			}
		case *ast.DropDatabaseStmt:
			goneDatabases = append(goneDatabases, v.Name)
			// A table without a database may be in the one dropped.
			for _, p := range list {
				if p.ref != nil && (p.ref.Schema == "" || sameName(p.ref.Schema, v.Name)) {
					p.settled = true
				}
			}
		}
	}
	var drops []columnDrop
	for _, p := range list {
		switch {
		case p.settled:
		case p.drop != nil:
			drops = append(drops, *p.drop)
		default:
			r.report(review.RequirePrimaryKey, p.index, p.rng, p.message)
		}
	}
	return drops
}

// columnDrop is a top-level ALTER TABLE that drops columns of a table the
// change did not make, drop, or rename before it. Whether it drops the
// table's primary key depends on the key's columns, which the target's
// synced schema holds.
type columnDrop struct {
	// statement is the index in Result.Statements, and position the
	// statement's place in the parsed list.
	statement int
	position  int
	table     *ast.TableRef
	// drops are the statement's DROP COLUMN subcommands, in order, and
	// touched the columns of the table an earlier statement, or another
	// subcommand of this one, named or made something depend on.
	drops   []*ast.AlterTableCmd
	touched []string
}

// columnTouch is a column of a table a statement names, or may make
// something depend on: a key, a generated column, a check, or a foreign
// key. An empty column is any column of the table, or its primary key.
type columnTouch struct {
	table  *ast.TableRef
	column string
}

// newColumnDrop returns the column drop of an ALTER TABLE, with the
// columns of its table the touches name. ok is false when one may have
// touched any column, or the key, of the table.
func newColumnDrop(statement, position int, ref *ast.TableRef, drops []*ast.AlterTableCmd, touches []columnTouch) (*columnDrop, bool) {
	d := &columnDrop{statement: statement, position: position, table: ref, drops: drops}
	for _, t := range touches {
		if !sameTable(t.table, ref) {
			continue
		}
		if t.column == "" {
			return nil, false
		}
		d.touched = append(d.touched, t.column)
	}
	return d, true
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

// checkColumnDrop reports an ALTER TABLE that drops every column of the
// primary key of the target's synced table: MySQL drops an index with the
// last of its columns, and a drop of only some leaves the key on the rest.
// The table is found as MySQL resolves it, and not at all when a USE
// before the statement may have moved the current database. The drop is
// not judged when an earlier statement, or another subcommand of the
// statement, touched a column of the key or one the statement drops, or
// when the server refuses it: it drops every column of the table, or
// something the synced schema holds depends on a column it drops. A
// foreign key of another database, which the synced schema does not hold,
// refuses it too.
//
// mysql: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html#alter-table-add-drop-column
func checkColumnDrop(stmts []statement, d columnDrop, x schemaIndex, r *reporter) {
	t := x.table(d.table, usesDatabase(stmts, d.position))
	key := primaryKeyColumns(t)
	if key == nil {
		return
	}
	drops := func(column string) bool {
		return slices.ContainsFunc(d.drops, func(cmd *ast.AlterTableCmd) bool { return sameName(cmd.Name, column) })
	}
	touched := func(column string) bool {
		return slices.ContainsFunc(d.touched, func(c string) bool { return sameName(c, column) })
	}
	for _, column := range key {
		if !drops(column) || touched(column) {
			return
		}
	}
	if !slices.ContainsFunc(t.GetColumns(), func(c *metadata.ColumnMetadata) bool { return !drops(c.GetName()) }) {
		return
	}
	for _, cmd := range d.drops {
		if touched(cmd.Name) || x.keepsColumn(t, cmd.Name, drops) {
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

// keepsColumn reports whether something the synced schema holds may keep
// the server from dropping a column of the table, given the columns the
// statement drops: a foreign key of the table that holds it, or one of
// the database that references it; a generated column the statement
// keeps that may read it; a check that reads it and another column; a
// functional index part that may read it; or the table's partitioning or
// subpartitioning, by an expression that may read it or by KEY (), which
// reads the primary key. A check that reads the column alone goes with
// it. A foreign key names the table it references without its database,
// so one that may reference a table of another database counts.
func (x schemaIndex) keepsColumn(t *metadata.TableMetadata, column string, drops func(string) bool) bool {
	has := func(list []string) bool {
		return slices.ContainsFunc(list, func(c string) bool { return sameName(c, column) })
	}
	if slices.ContainsFunc(t.GetForeignKeys(), func(fk *metadata.ForeignKeyMetadata) bool { return has(fk.GetColumns()) }) {
		return true
	}
	for _, s := range x.schema.GetSchemas() {
		for _, other := range s.GetTables() {
			for _, fk := range other.GetForeignKeys() {
				db := fk.GetReferencedSchema()
				if (db == "" || x.same(db, x.schema.GetName())) && x.same(fk.GetReferencedTable(), t.GetName()) && has(fk.GetReferencedColumns()) {
					return true
				}
			}
		}
	}
	for _, c := range t.GetColumns() {
		if g := c.GetGeneration(); g != nil && !drops(c.GetName()) && mentions(g.GetExpression(), column) {
			return true
		}
	}
	for _, c := range t.GetCheckConstraints() {
		reads, ok := expressionColumns(c.GetExpression())
		if !ok && mentions(c.GetExpression(), column) || has(reads) && slices.ContainsFunc(reads, func(r string) bool { return !sameName(r, column) }) {
			return true
		}
	}
	for _, idx := range t.GetIndexes() {
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
