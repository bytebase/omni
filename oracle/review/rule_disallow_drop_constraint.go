package review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// constraintKinds names the constraint kinds DisallowDropConstraint
// covers.
var constraintKinds = map[ast.ConstraintType]string{
	ast.CONSTRAINT_PRIMARY: "primary key",
	ast.CONSTRAINT_UNIQUE:  "unique constraint",
	ast.CONSTRAINT_FOREIGN: "foreign key",
	ast.CONSTRAINT_CHECK:   "check constraint",
}

// scanConstraints follows the change's constraints, statement by
// statement, for DisallowDropConstraint and RequirePrimaryKey. It does not
// simulate the change: it reads what each statement certainly does to a
// table's constraints, and where a statement may have done more, it stops
// judging what it no longer knows.
//
// DisallowDropConstraint reports ALTER TABLE ... DROP PRIMARY KEY and DROP
// UNIQUE (columns), from the SQL alone, and DROP CONSTRAINT of a name
// whose kind is known when the statement runs: a primary key, unique,
// foreign key, or check constraint the change itself defined on the table
// earlier and has not dropped since, or else the synced table's
// constraint of the name. Oracle names constraints per schema, so a
// constraint the change defines under a name already taken makes its
// statement fail (ORA-02264): from such a collision on, or one the scan
// cannot rule out, no DROP CONSTRAINT is judged. Neither is one of a name
// the change cannot place, nor one on a table whose constraints the change
// may have changed in a way not followed: a table it created with IF NOT
// EXISTS, renamed, restored with FLASHBACK, or dropped a column of, which
// takes the column's constraints along. A foreign key is not judged after
// a drop that CASCADE carries to other tables' foreign keys. DROP
// CONSTRAINT inside PL/SQL runs only as the text of an EXECUTE IMMEDIATE,
// which is not read.
//
// oracle: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-TABLE.html (drop_constraint_clause)
func scanConstraints(sql string, stmts []statement, x schemaIndex, on map[review.Rule]bool, r *reporter) {
	s := &constraintScan{
		sql:           sql,
		x:             x,
		on:            on,
		r:             r,
		tables:        make(map[tableKey]*tableState),
		freed:         make(map[string]bool),
		maybe:         make(map[string]bool),
		droppedSynced: make(map[tableKey]bool),
	}
	if x.schema != nil {
		s.synced = x.constraintNames()
	}
	schemas := schemasAt(stmts, x.db())
	for pos := range stmts {
		st := &stmts[pos]
		cur := schemas[pos]
		switch v := st.node.(type) {
		case *ast.CreateTableStmt:
			s.createTable(st, v, rangeOf(st.loc), cur)
		case *ast.CreateSchemaStmt:
			if v.Stmts == nil {
				continue
			}
			for _, item := range v.Stmts.Items {
				if ct, ok := item.(*ast.CreateTableStmt); ok {
					s.createTable(st, ct, rangeOf(ct.Loc), v.SchemaName)
				}
			}
		case *ast.AlterTableStmt:
			s.alterTable(st, v, cur)
		case *ast.DropStmt:
			s.drop(v, cur)
		case *ast.AdminDDLStmt:
			// DROP USER ... CASCADE drops the schema's tables, and the
			// foreign keys other schemas' tables have on them.
			if strings.EqualFold(v.Action, "DROP") && v.ObjectType == ast.OBJECT_USER && v.Name != nil {
				s.settleWhere(func(k tableKey) bool { return k.schema == v.Name.Name || k.schema == "" || k.schema == unknownSchema })
				s.cascaded = true
			}
		case *ast.RenameStmt:
			if v.OldName != nil && v.NewName != nil {
				s.renamed(keyIn(v.OldName, cur), keyIn(v.NewName, cur))
			}
		case *ast.FlashbackTableStmt:
			// FLASHBACK TABLE ... TO BEFORE DROP restores a table, its
			// constraints renamed; the other forms change its rows.
			for _, name := range v.Tables {
				keys := []tableKey{keyIn(name, cur)}
				if v.Rename != "" {
					keys = append(keys, tableKey{schema: keys[0].schema, name: v.Rename})
				}
				for _, k := range keys {
					s.settle(k)
					s.touch(k, func(t *tableState) { t.replaced, t.unsure = true, true })
				}
			}
		}
	}
	s.reportPending()
}

// constraintScan is the state scanConstraints carries from statement to
// statement.
type constraintScan struct {
	sql string
	x   schemaIndex
	on  map[review.Rule]bool
	r   *reporter

	// synced is the target database's constraints by name, nil without a
	// synced schema.
	synced map[string]syncedConstraint
	// defined lists the named constraints the change defined, in order.
	defined []*definedConstraint
	// tables is what the scan knows of each table the change touched.
	tables map[tableKey]*tableState
	// freed holds the names the change dropped a constraint of; maybe the
	// names a constraint may have from a statement the scan cannot follow.
	freed map[string]bool
	maybe map[string]bool
	// droppedSynced holds the synced tables the change dropped, which
	// took their constraints along.
	droppedSynced map[tableKey]bool
	// lost is set from a statement on which what the change's
	// constraints are is not known: a constraint-name collision, which
	// fails its statement, or a DROP CLUSTER, which may drop tables it
	// does not name.
	lost bool
	// cascaded is set after a drop CASCADE carries to the foreign keys of
	// tables it does not name.
	cascaded bool
	// pending lists the tables RequirePrimaryKey reports at the end of
	// the change unless a later statement settles them.
	pending []*pendingTable
}

// definedConstraint is a named constraint the change defined.
type definedConstraint struct {
	table tableKey
	name  string
	kind  ast.ConstraintType
	alive bool
}

// tableState is what the scan knows of a table of a name.
type tableState struct {
	// replaced: the table of the name is not the synced table of the
	// name, whose constraints are no longer its.
	replaced bool
	// unsure: the table's constraints may have changed in a way the scan
	// does not follow.
	unsure bool
	// created: the change made the table of this name.
	created bool
	// pkDropped and uniqueDropped: the change dropped the table's primary
	// key, or a unique constraint it does not name.
	pkDropped     bool
	uniqueDropped bool
	// pkFirst is the first thing the change did to the primary key of a
	// table it did not make: "add" or "drop".
	pkFirst string
}

// state returns the state of exactly the key, made on first use.
func (s *constraintScan) state(k tableKey) *tableState {
	t, ok := s.tables[k]
	if !ok {
		t = &tableState{}
		s.tables[k] = t
	}
	return t
}

// touch applies a change to every table the key may name.
func (s *constraintScan) touch(k tableKey, f func(*tableState)) {
	f(s.state(k))
	for o, t := range s.tables {
		if o != k && o.mayBe(k) {
			f(t)
		}
	}
}

// anyState reports whether a table the key may name is in the state.
func (s *constraintScan) anyState(k tableKey, f func(*tableState) bool) bool {
	for o, t := range s.tables {
		if o.mayBe(k) && f(t) {
			return true
		}
	}
	return false
}

// createTable follows a CREATE TABLE: a new table with the constraints it
// declares, or, with IF NOT EXISTS, maybe the table that was there.
func (s *constraintScan) createTable(st *statement, v *ast.CreateTableStmt, rng review.Range, cur string) {
	if v.Name == nil || v.Name.Name == "" || v.Name.DBLink != "" {
		return
	}
	k := keyIn(v.Name, cur)
	s.settle(k)
	declared := tableConstraints(v)
	if v.IfNotExists {
		s.touch(k, func(t *tableState) { t.replaced, t.unsure = true, true })
		for _, c := range declared {
			if c.name != "" {
				s.maybe[c.name] = true
			}
		}
		return
	}
	// The name was free when the table was made: what the change defined
	// on a table of the name before is gone with it.
	for _, d := range s.defined {
		if d.alive && d.table.mayBe(k) {
			d.alive = false
		}
	}
	s.touch(k, func(t *tableState) { t.replaced, t.unsure = true, true })
	s.tables[k] = &tableState{replaced: true, created: true}
	for _, c := range declared {
		s.define(k, c.name, c.kind)
	}
	if s.on[review.RequirePrimaryKey] && judgedTable(v) && !createsPrimaryKey(v) {
		s.pending = append(s.pending, &pendingTable{key: k, statement: st.index, rng: rng, message: "creates table " + objectName(v.Name) + " without a primary key"})
	}
}

// alterTable follows the subcommands of an ALTER TABLE, in order.
func (s *constraintScan) alterTable(st *statement, v *ast.AlterTableStmt, cur string) {
	if v.Name == nil || v.Name.DBLink != "" || v.Actions == nil {
		return
	}
	k := keyIn(v.Name, cur)
	for _, item := range v.Actions.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		switch cmd.Action {
		case ast.AT_ADD_CONSTRAINT:
			if c := cmd.Constraint; c != nil {
				s.define(k, c.Name, c.Type)
				if c.Type == ast.CONSTRAINT_PRIMARY {
					s.gainedPrimaryKey(k)
				}
			}
		case ast.AT_ADD_COLUMN, ast.AT_MODIFY_COLUMN:
			for _, col := range commandColumns(cmd) {
				for _, c := range columnConstraints(col) {
					s.define(k, c.Name, c.Type)
					if c.Type == ast.CONSTRAINT_PRIMARY {
						s.gainedPrimaryKey(k)
					}
				}
				// MODIFY ... NULL drops the column's NOT NULL constraint,
				// which Oracle keeps as a check constraint.
				if cmd.Action == ast.AT_MODIFY_COLUMN && col.Null {
					s.touch(k, func(t *tableState) { t.unsure = true })
				}
			}
		case ast.AT_DROP_CONSTRAINT:
			switch {
			case strings.EqualFold(cmd.Subtype, "PRIMARY KEY"):
				s.dropPrimaryKey(st, v, cmd, k)
			case strings.EqualFold(cmd.Subtype, "UNIQUE"):
				s.dropUnique(st, v, cmd, k)
			case cmd.Constraint != nil && cmd.Constraint.Name != "":
				s.dropNamed(st, v, cmd, k, cmd.Constraint.Name)
			}
			if cascades(s.sql, cmd.Loc) {
				s.cascaded = true
			}
		case ast.AT_DROP_COLUMN, ast.AT_SET_UNUSED, ast.AT_DROP_UNUSED_COLUMNS:
			s.touch(k, func(t *tableState) { t.unsure = true })
			if cascades(s.sql, cmd.Loc) {
				s.cascaded = true
			}
		case ast.AT_RENAME_CONSTRAINT:
			for _, d := range s.defined {
				if d.alive && d.name == cmd.ColumnName && d.table.mayBe(k) {
					d.alive = false
				}
			}
			s.freed[cmd.ColumnName] = true
			s.maybe[cmd.NewName] = true
			s.touch(k, func(t *tableState) { t.unsure = true })
		case ast.AT_RENAME:
			if cmd.NewName != "" {
				s.renamed(k, tableKey{schema: k.schema, name: cmd.NewName})
			}
		}
	}
}

// drop follows a DROP TABLE, which takes the table's constraints along,
// and a DROP CLUSTER, which with INCLUDING TABLES drops tables it does not
// name.
func (s *constraintScan) drop(v *ast.DropStmt, cur string) {
	switch v.ObjectType {
	case ast.OBJECT_TABLE:
		if v.Names == nil {
			return
		}
		for _, item := range v.Names.Items {
			name, ok := item.(*ast.ObjectName)
			if !ok || name.DBLink != "" {
				continue
			}
			k := keyIn(name, cur)
			s.settle(k)
			if !s.anyState(k, func(t *tableState) bool { return t.replaced }) {
				s.droppedSynced[k] = true
			}
			for _, d := range s.defined {
				if d.alive && d.table.mayBe(k) {
					d.alive = false
				}
			}
			s.touch(k, func(t *tableState) { *t = tableState{replaced: true} })
		}
		if v.Cascade {
			s.cascaded = true
		}
	case ast.OBJECT_CLUSTER:
		s.settleWhere(func(tableKey) bool { return true })
		s.lost = true
		s.cascaded = true
	}
}

// renamed follows a table renamed from one name to another: the table
// keeps its constraints under the new name, which the scan no longer
// follows, and what the change did to its primary key.
func (s *constraintScan) renamed(from, to tableKey) {
	s.settle(from)
	s.settle(to)
	for _, d := range s.defined {
		if d.alive && d.table.is(from) {
			d.table = to
		}
	}
	old := s.state(from)
	created, pkFirst := old.created, old.pkFirst
	s.touch(from, func(t *tableState) { *t = tableState{replaced: true, unsure: true} })
	s.touch(to, func(t *tableState) { t.replaced, t.unsure = true, true })
	moved := s.state(to)
	moved.created, moved.pkFirst = created, pkFirst
}

// define follows a named constraint the change defines. A name already
// taken in the schema, or one the scan cannot rule out, fails the
// statement, and the scan loses track.
func (s *constraintScan) define(k tableKey, name string, kind ast.ConstraintType) {
	if name == "" {
		return
	}
	if s.inUse(name) {
		s.lost = true
	}
	s.defined = append(s.defined, &definedConstraint{table: k, name: name, kind: kind, alive: true})
}

// inUse reports whether a constraint of the name may exist: one the change
// defined and has not dropped, one the scan cannot place, or a synced one
// the change has not dropped, by name, with its table, or as the primary
// key it dropped.
func (s *constraintScan) inUse(name string) bool {
	if s.maybe[name] {
		return true
	}
	for _, d := range s.defined {
		if d.alive && d.name == name {
			return true
		}
	}
	c, ok := s.synced[name]
	if !ok || s.freed[name] || s.droppedSynced[c.table] {
		return false
	}
	if t, ok := s.tables[c.table]; ok && c.kind == ast.CONSTRAINT_PRIMARY && t.pkDropped && !t.replaced {
		return false
	}
	return true
}

// findDefined returns the live constraint the change defined on exactly
// the table under the name, nil for none.
func (s *constraintScan) findDefined(k tableKey, name string) *definedConstraint {
	for i := len(s.defined) - 1; i >= 0; i-- {
		d := s.defined[i]
		if d.alive && d.name == name && d.table.is(k) {
			return d
		}
	}
	return nil
}

// dropPrimaryKey follows ALTER TABLE ... DROP PRIMARY KEY, which drops
// the table's primary key, whatever its name.
func (s *constraintScan) dropPrimaryKey(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, k tableKey) {
	if s.on[review.DisallowDropConstraint] {
		s.r.report(review.DisallowDropConstraint, st.index, rangeOf(cmd.Loc), "drops the primary key of "+objectName(v.Name))
	}
	for _, d := range s.defined {
		if d.alive && d.kind == ast.CONSTRAINT_PRIMARY && d.table.mayBe(k) {
			d.alive = false
		}
	}
	s.lostPrimaryKey(st, v, cmd, k)
}

// dropUnique follows ALTER TABLE ... DROP UNIQUE (columns), which drops
// the unique constraint on the columns, whatever its name.
func (s *constraintScan) dropUnique(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, k tableKey) {
	if s.on[review.DisallowDropConstraint] {
		s.r.report(review.DisallowDropConstraint, st.index, rangeOf(cmd.Loc), "drops the unique constraint on "+uniqueColumns(s.sql, cmd)+" of "+objectName(v.Name))
	}
	s.touch(k, func(t *tableState) { t.uniqueDropped = true })
}

// dropNamed follows ALTER TABLE ... DROP CONSTRAINT name.
func (s *constraintScan) dropNamed(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, k tableKey, name string) {
	defer func() { s.freed[name] = true }()
	// judged reports whether a drop of a constraint of the kind is judged:
	// the scan knows the table's constraints, and nothing dropped one of
	// the kind out of its sight.
	judged := func(kind ast.ConstraintType) bool {
		switch {
		case s.lost, s.anyState(k, func(t *tableState) bool { return t.unsure }):
			return false
		case kind == ast.CONSTRAINT_PRIMARY:
			return !s.anyState(k, func(t *tableState) bool { return t.pkDropped })
		case kind == ast.CONSTRAINT_UNIQUE:
			return !s.anyState(k, func(t *tableState) bool { return t.uniqueDropped })
		case kind == ast.CONSTRAINT_FOREIGN:
			return !s.cascaded
		}
		return true
	}
	if d := s.findDefined(k, name); d != nil {
		d.alive = false
		if judged(d.kind) {
			s.dropped(st, v, cmd, k, name, d.kind)
		}
		return
	}
	if s.freed[name] || s.maybe[name] || !s.x.own(k) || s.anyState(k, func(t *tableState) bool { return t.replaced }) {
		return
	}
	t := s.x.table(k)
	if t == nil {
		return
	}
	if kind, ok := constraintKind(t, name); ok && judged(kind) {
		s.dropped(st, v, cmd, k, name, kind)
	}
}

// dropped reports the drop of a constraint of a known kind: under
// DisallowDropConstraint when the kind is one it covers, and to
// RequirePrimaryKey when it is the primary key.
func (s *constraintScan) dropped(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, k tableKey, name string, kind ast.ConstraintType) {
	label, ok := constraintKinds[kind]
	if !ok {
		return
	}
	if s.on[review.DisallowDropConstraint] {
		s.r.report(review.DisallowDropConstraint, st.index, rangeOf(cmd.Loc), "drops "+label+" "+ident(name)+" of "+objectName(v.Name))
	}
	if kind == ast.CONSTRAINT_PRIMARY {
		s.lostPrimaryKey(st, v, cmd, k)
	}
}

// uniqueColumns writes the column list of DROP UNIQUE (columns), which the
// parser does not keep.
func uniqueColumns(sql string, cmd *ast.AlterTableCmd) string {
	var names []string
	open := false
	for _, tok := range tokens(sql, cmd.Loc) {
		switch {
		case tok.Type == '(':
			open = true
		case tok.Type == ')':
			if open {
				return "(" + strings.Join(names, ", ") + ")"
			}
		case open && tok.Type != ',':
			names = append(names, ident(tok.Str))
		}
	}
	return "its columns"
}

// tableConstraint is a named or unnamed constraint a CREATE TABLE
// declares.
type tableConstraint struct {
	name string
	kind ast.ConstraintType
}

// tableConstraints lists a CREATE TABLE's constraints, of its columns and
// of the table. The parser keeps no name for a named NOT NULL constraint.
func tableConstraints(v *ast.CreateTableStmt) []tableConstraint {
	var out []tableConstraint
	if v.Columns != nil {
		for _, item := range v.Columns.Items {
			if col, ok := item.(*ast.ColumnDef); ok {
				for _, c := range columnConstraints(col) {
					out = append(out, tableConstraint{c.Name, c.Type})
				}
			}
		}
	}
	if v.Constraints != nil {
		for _, item := range v.Constraints.Items {
			if c, ok := item.(*ast.TableConstraint); ok {
				out = append(out, tableConstraint{c.Name, c.Type})
			}
		}
	}
	return out
}

// columnConstraints lists a column definition's constraints.
func columnConstraints(col *ast.ColumnDef) []*ast.ColumnConstraint {
	if col == nil || col.Constraints == nil {
		return nil
	}
	var out []*ast.ColumnConstraint
	for _, item := range col.Constraints.Items {
		if c, ok := item.(*ast.ColumnConstraint); ok {
			out = append(out, c)
		}
	}
	return out
}

// commandColumns lists the column definitions of an ADD or MODIFY.
func commandColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	var out []*ast.ColumnDef
	if cmd.ColumnDef != nil {
		out = append(out, cmd.ColumnDef)
	}
	if cmd.ColumnDefs != nil {
		for _, item := range cmd.ColumnDefs.Items {
			if col, ok := item.(*ast.ColumnDef); ok {
				out = append(out, col)
			}
		}
	}
	return out
}
