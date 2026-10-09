package review

import (
	"slices"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// alterTable follows ALTER TABLE. DisallowDropConstraint reports each DROP
// CONSTRAINT that names a primary key, foreign key, unique, or check
// constraint the table has in the synced schema, when nothing earlier in
// the change may have changed which table the name means or what the
// constraint is. The type comes from the schema, since the statement does
// not say it. A constraint the change itself added is not the schema's, so
// dropping it is not reported. RequirePrimaryKey learns here that a table
// lost its primary key, or gained one.
//
// The server runs every DROP subcommand of the statement before the
// others, so the scan does too, and the statement succeeds or fails as a
// whole, so what the drops find is reported only when no subcommand is
// known to fail.
//
// pg: src/backend/commands/tablecmds.c — ATController (AT_PASS_DROP first)
func (s *scan) alterTable(st *statement, v *ast.AlterTableStmt) {
	if v.Cmds == nil || v.Relation == nil {
		return
	}
	// ALTER of a relation the target certainly lacks is refused.
	if !v.Missing_ok && s.index != nil && s.missing(v.Relation) {
		s.stop()
		return
	}
	isTable := ast.ObjectType(v.ObjType) == ast.OBJECT_TABLE
	var drops, others []*ast.AlterTableCmd
	for _, item := range v.Cmds.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			continue
		}
		// ALTER TYPE ... CASCADE alters the tables of the type.
		if ast.ObjectType(v.ObjType) == ast.OBJECT_TYPE && cmd.Behavior == int(ast.DROP_CASCADE) {
			s.stop()
			return
		}
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropConstraint, ast.AT_DropColumn:
			drops = append(drops, cmd)
		default:
			others = append(others, cmd)
		}
	}
	name := v.Relation.Relname
	schema := s.schemaOf(v.Relation)
	before, knownBefore := s.table(v.Relation)
	// What the drops find is kept until every subcommand is known to run:
	// a statement the server refuses drops nothing. A DROP the server may
	// refuse withholds the statement's findings without ending the scan.
	var out droppedKeys
	withhold := false
	droppedColumns := make(map[string]bool)
	droppedNames := make(map[string]bool)
	for _, cmd := range drops {
		cascade := cmd.Behavior == int(ast.DROP_CASCADE)
		t, known := s.table(v.Relation)
		known = known && isTable
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropConstraint:
			// A second drop of the same constraint fails, and so may a drop
			// of one an earlier column drop of the statement took with it.
			if droppedNames[cmd.Name] && !cmd.Missing_ok {
				s.stop()
				return
			}
			if len(droppedColumns) > 0 && !cmd.Missing_ok && knownBefore && s.mayDependOn(before, cmd.Name, droppedColumns) {
				withhold = true
			}
			droppedNames[cmd.Name] = true
			s.retireReferences(v.Relation, cmd.Name)
			if known && s.constraintKnown(t, cmd.Name) && s.dropConstraint(st, v, cmd, t, cascade, &out) {
				s.stop()
				return
			}
			s.constraints[[3]string{schema, name, cmd.Name}] = true
			if cascade {
				s.dropKeyReferences(v.Relation, t, known, cmd.Name)
			}
		case ast.AT_DropColumn:
			if droppedColumns[cmd.Name] && !cmd.Missing_ok {
				s.stop()
				return
			}
			droppedColumns[cmd.Name] = true
			if known && s.dropColumn(cmd, t, cascade, &out) {
				s.stop()
				return
			}
			s.unsettled[[2]string{schema, name}] = true
			if cascade {
				s.dropReaders(t, v.Relation)
				s.dropReferences(v.Relation, t, known, func(fk foreignKeyRef) bool {
					return !known || slices.ContainsFunc(fk.columns, func(c string) bool { n := columnName(c); return n == cmd.Name || n == "" })
				})
			}
		}
	}
	if withhold || knownBefore && isTable && s.refuses(before, drops, others) {
		if !withhold {
			s.stop()
			return
		}
		out = droppedKeys{}
	}
	for _, f := range out.findings {
		s.r.add(f)
	}
	if out.lost != nil {
		s.lostKey(st, v, out.lost, out.table)
	}
	for _, cmd := range others {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddConstraint:
			if c, ok := cmd.Def.(*ast.Constraint); ok {
				if c.Conname != "" {
					s.constraints[[3]string{schema, name, c.Conname}] = true
				}
				switch {
				case !makesIndex(c):
				case c.Conname == "":
					s.generated = true
				default:
					// A key's index takes the constraint's name.
					s.touchName(v.Relation.Schemaname, c.Conname)
				}
				if c.Indexname != "" {
					// USING INDEX renames the index to the constraint.
					s.touchName(v.Relation.Schemaname, c.Indexname)
					s.touchName(v.Relation.Schemaname, c.Conname)
					s.renamedKeys[c.Indexname] = true
				}
				if c.Contype == ast.CONSTR_PRIMARY {
					s.keyed(v.Relation)
				}
				s.references(v.Relation, c)
			}
		case ast.AT_AddColumn:
			if cd, ok := cmd.Def.(*ast.ColumnDef); ok {
				// ADD COLUMN IF NOT EXISTS of a column the table has does
				// nothing, its constraints included.
				if cmd.Missing_ok && knownBefore && !droppedColumns[cd.Colname] && s.hasColumn(before, cd.Colname) {
					continue
				}
				s.generated = true
				s.columns[[3]string{schema, name, cd.Colname}] = true
				for _, c := range constraintsOf(cd.Constraints) {
					if c.Conname != "" {
						s.constraints[[3]string{schema, name, c.Conname}] = true
					}
					if c.Contype == ast.CONSTR_PRIMARY {
						s.keyed(v.Relation)
					}
				}
				s.references(v.Relation, constraintsOf(cd.Constraints)...)
			}
		case ast.AT_AddIdentity:
			s.generated = true
		case ast.AT_AddIndexConstraint, ast.AT_AddIndex:
			if idx, ok := cmd.Def.(*ast.IndexStmt); ok {
				s.renamedKeys[idx.Idxname] = true
				if idx.Primary {
					s.keyed(v.Relation)
				}
			}
		case ast.AT_AttachPartition:
			if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
				s.unsettle(pc.Name)
				for _, p := range s.matching(pc.Name) {
					p.settled = true
				}
			}
		case ast.AT_DetachPartition:
			if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
				s.unsettle(pc.Name)
			}
		case ast.AT_AddInherit, ast.AT_DropInherit, ast.AT_AddOf, ast.AT_DropOf:
			s.unsettled[[2]string{schema, name}] = true
		}
	}
}

// droppedKeys is what the drops of one statement found: the findings of
// DisallowDropConstraint, and the subcommand that removed the table's
// primary key, if one did.
type droppedKeys struct {
	findings []review.Finding
	lost     *ast.AlterTableCmd
	table    tableRef
}

// dropConstraint checks a DROP CONSTRAINT against the synced table. It
// returns true when the server refuses the statement: the constraint does
// not exist and IF EXISTS is absent, or it is a key a foreign key
// references and CASCADE is absent.
func (s *scan) dropConstraint(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, t tableRef, cascade bool, out *droppedKeys) (refused bool) {
	table := s.index.schemas[t.schema].tables[t.table]
	con, ok := constraint(table, cmd.Name)
	if !ok {
		// The snapshot leaves out NOT NULL constraints, so a constraint it
		// does not list may still exist; either way nothing is reported.
		return !cmd.Missing_ok
	}
	if !cascade && (con.kind == "primary key" || con.kind == "unique constraint") {
		// When the key's columns cannot be read, any foreign key on the
		// table may reference it, and so may one the change added.
		if con.columns != nil && s.index.referenced(t, con.columns, s.dropped) || con.columns == nil && s.index.referencedAny(t, s.dropped) ||
			s.newlyReferenced(t, con.columns, con.kind == "primary key") {
			return true
		}
	}
	s.dropped[[3]string{t.schema, t.table, cmd.Name}] = true
	if s.on[review.DisallowDropConstraint] && con.kind != "" {
		out.findings = append(out.findings, review.Finding{
			Rule:      review.DisallowDropConstraint,
			Statement: st.index,
			Range:     rangeOf(cmd.Loc),
			Message:   "drops " + con.kind + " " + ident(cmd.Name) + " of " + relation(v.Relation),
		})
	}
	if con.kind == "primary key" && out.lost == nil {
		out.lost, out.table = cmd, t
	}
	return false
}

// dropColumn checks a DROP COLUMN against the synced table and notes a
// column of its primary key, which goes with the column. It returns true
// when the server refuses the statement: the column does not exist and IF
// EXISTS is absent, or, without CASCADE, something outside the table
// depends on it: a foreign key referencing a key with the column, or a
// view reading it. A foreign key or a view the change created may depend
// on any column of the tables it names.
//
// pg: src/backend/commands/tablecmds.c — ATExecDropColumn
func (s *scan) dropColumn(cmd *ast.AlterTableCmd, t tableRef, cascade bool, out *droppedKeys) (refused bool) {
	table := s.index.schemas[t.schema].tables[t.table]
	added := hasColumnIn(s.columns, t, cmd.Name)
	switch {
	case hasColumnIn(s.renamedTo, t, cmd.Name):
		// A column renamed to this name: what it was is not followed.
		return false
	case hasColumnIn(s.renamedFrom, t, cmd.Name) && !added:
		return !cmd.Missing_ok
	}
	if !slices.ContainsFunc(table.GetColumns(), func(c *metadata.ColumnMetadata) bool { return c.GetName() == cmd.Name }) {
		// A column the change added is not the key's.
		return !cmd.Missing_ok && !added
	}
	pk := primaryKey(table)
	var pkColumns []string
	if pk != nil {
		pkColumns = keyColumns(pk)
	}
	if !cascade && (s.index.referencesColumn(t, cmd.Name, s.dropped) || slices.ContainsFunc(s.index.readsColumn[columnRef{t.schema, t.table, cmd.Name}], s.viewHolds) ||
		s.newlyReferencedColumn(t, cmd.Name, pkColumns, pk != nil) || dependsOnNew(s.newReads, t)) {
		return true
	}
	if pk == nil || !s.constraintKnown(t, pk.GetName()) {
		return false
	}
	if !slices.Contains(keyColumns(pk), cmd.Name) {
		return false
	}
	if out.lost == nil {
		out.lost, out.table = cmd, t
	}
	return false
}

// dropReaders touches the views that read the table, which a cascading
// DROP COLUMN may drop: in the table's schema when the scan resolved it,
// and in every schema the name may be in otherwise.
func (s *scan) dropReaders(t tableRef, rv *ast.RangeVar) {
	if s.index == nil {
		return
	}
	for schema := range s.index.schemas {
		if t.schema != "" && t.schema != schema || t.schema == "" && rv.Schemaname != "" && rv.Schemaname != schema {
			continue
		}
		for _, view := range s.index.readers[tableRef{schema, rv.Relname}] {
			s.touchName(view.schema, view.table)
		}
	}
}

// dropKeyReferences records the foreign keys a cascading DROP CONSTRAINT
// takes: those on the dropped key when the scan knows the table and the
// key's columns, none when the constraint is not a key, and every one
// referencing the table otherwise.
func (s *scan) dropKeyReferences(rv *ast.RangeVar, t tableRef, known bool, name string) {
	if !known {
		s.dropReferences(rv, t, false, nil)
		return
	}
	con, ok := constraint(s.index.schemas[t.schema].tables[t.table], name)
	switch {
	case ok && con.kind != "primary key" && con.kind != "unique constraint" && con.kind != "":
		// A foreign key or check constraint has no foreign keys on it.
	case ok && con.columns != nil:
		s.dropReferences(rv, t, true, func(fk foreignKeyRef) bool { return sameColumns(fk.columns, con.columns) })
	default:
		s.dropReferences(rv, t, true, nil)
	}
}

// mayDependOn reports whether a constraint of the synced table may involve
// one of the columns, so that dropping the column took it: a check or an
// exclusion constraint, whose columns the snapshot does not list, or a
// key or foreign key on one of them.
func (s *scan) mayDependOn(t tableRef, name string, columns map[string]bool) bool {
	table := s.index.schemas[t.schema].tables[t.table]
	for _, fk := range table.GetForeignKeys() {
		if fk.GetName() == name {
			return slices.ContainsFunc(fk.GetColumns(), func(c string) bool { return columns[columnName(c)] || columnName(c) == "" })
		}
	}
	con, ok := constraint(table, name)
	if !ok || con.columns == nil {
		return true
	}
	return slices.ContainsFunc(con.columns, func(c string) bool { return columns[c] })
}

// refuses reports whether a subcommand of the statement is certain to
// fail against the synced table, after its drops: adding a column it has,
// altering a column it lacks, or adding a constraint under a name it
// uses, or a second primary key. Other failures, such as a type error,
// are not known before the statement runs.
func (s *scan) refuses(t tableRef, drops, others []*ast.AlterTableCmd) bool {
	table := s.index.schemas[t.schema].tables[t.table]
	droppedColumns := make(map[string]bool)
	droppedConstraints := make(map[string]bool)
	for _, cmd := range drops {
		if ast.AlterTableType(cmd.Subtype) == ast.AT_DropColumn {
			droppedColumns[cmd.Name] = true
		} else {
			droppedConstraints[cmd.Name] = true
		}
	}
	// has reports whether the table has the column once the drops ran.
	has := func(column string) bool {
		return !droppedColumns[column] && s.hasColumn(t, column)
	}
	// uses reports whether a constraint name certainly names a constraint
	// of the table.
	uses := func(name string) bool {
		if droppedConstraints[name] || !s.constraintKnown(t, name) || s.dropped[[3]string{t.schema, t.table, name}] {
			return false
		}
		_, ok := constraint(table, name)
		return ok
	}
	// What the statement's own additions add counts too.
	addedColumns := make(map[string]bool)
	addedNames := make(map[string]bool)
	addedKey := false
	keyed := func() bool {
		pk := primaryKey(table)
		return addedKey || pk != nil && uses(pk.GetName())
	}
	addKey := func(c *ast.Constraint) bool {
		if c.Conname != "" {
			if uses(c.Conname) || addedNames[c.Conname] {
				return true
			}
			addedNames[c.Conname] = true
		}
		if c.Contype == ast.CONSTR_PRIMARY {
			if keyed() {
				return true
			}
			addedKey = true
		}
		return false
	}
	for _, cmd := range others {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddColumn:
			cd, ok := cmd.Def.(*ast.ColumnDef)
			if !ok {
				continue
			}
			if (has(cd.Colname) || addedColumns[cd.Colname]) && !cmd.Missing_ok {
				return true
			}
			addedColumns[cd.Colname] = true
			for _, c := range constraintsOf(cd.Constraints) {
				if addKey(c) {
					return true
				}
			}
		case ast.AT_ColumnDefault, ast.AT_DropNotNull, ast.AT_SetNotNull, ast.AT_AlterColumnType,
			ast.AT_SetStatistics, ast.AT_SetStorage, ast.AT_SetCompression, ast.AT_SetOptions, ast.AT_ResetOptions,
			ast.AT_AddIdentity, ast.AT_SetIdentity, ast.AT_DropIdentity, ast.AT_SetExpression, ast.AT_DropExpression:
			// A column the statement adds may or may not exist when this
			// subcommand's pass runs.
			if cmd.Name == "" || addedColumns[cmd.Name] {
				continue
			}
			if !has(cmd.Name) {
				return true
			}
		case ast.AT_AddConstraint:
			c, ok := cmd.Def.(*ast.Constraint)
			if !ok {
				continue
			}
			if addKey(c) {
				return true
			}
		}
	}
	return false
}

// hasColumn reports whether the synced table has the column after the
// statements so far: added, renamed to it, or synced and not renamed away.
func (s *scan) hasColumn(t tableRef, column string) bool {
	switch {
	case hasColumnIn(s.columns, t, column) || hasColumnIn(s.renamedTo, t, column):
		return true
	case hasColumnIn(s.renamedFrom, t, column):
		return false
	}
	table := s.index.schemas[t.schema].tables[t.table]
	return slices.ContainsFunc(table.GetColumns(), func(c *metadata.ColumnMetadata) bool { return c.GetName() == column })
}

// makesIndex reports whether a constraint is backed by an index.
func makesIndex(c *ast.Constraint) bool {
	return c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE || c.Contype == ast.CONSTR_EXCLUSION
}

func constraintsOf(list *ast.List) []*ast.Constraint {
	if list == nil {
		return nil
	}
	var out []*ast.Constraint
	for _, item := range list.Items {
		if c, ok := item.(*ast.Constraint); ok {
			out = append(out, c)
		}
	}
	return out
}
