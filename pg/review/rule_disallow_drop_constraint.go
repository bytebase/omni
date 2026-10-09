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
// others, so the scan does too.
//
// pg: src/backend/commands/tablecmds.c — ATController (AT_PASS_DROP first)
func (s *scan) alterTable(st *statement, v *ast.AlterTableStmt) {
	if v.Cmds == nil || v.Relation == nil {
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
	// What the drops find is kept until every drop is known to run: a
	// statement the server refuses drops nothing.
	var out droppedKeys
	for _, cmd := range drops {
		cascade := cmd.Behavior == int(ast.DROP_CASCADE)
		t, known := s.table(v.Relation)
		known = known && isTable
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropConstraint:
			if known && s.constraintKnown(name, cmd.Name) && s.dropConstraint(st, v, cmd, t, cascade, &out) {
				s.stop()
				return
			}
			s.constraints[[2]string{name, cmd.Name}] = true
		case ast.AT_DropColumn:
			if known && s.dropColumn(cmd, t, cascade, &out) {
				s.stop()
				return
			}
			s.unsettled[name] = true
			if cascade {
				s.dropReaders(t, v.Relation)
			}
		}
		if cascade {
			s.dropReferences(name)
		}
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
					s.constraints[[2]string{name, c.Conname}] = true
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
				s.references(c)
			}
		case ast.AT_AddColumn:
			if cd, ok := cmd.Def.(*ast.ColumnDef); ok {
				s.columns[[2]string{name, cd.Colname}] = true
				for _, c := range constraintsOf(cd.Constraints) {
					if c.Conname != "" {
						s.constraints[[2]string{name, c.Conname}] = true
					}
					if c.Contype == ast.CONSTR_PRIMARY {
						s.keyed(v.Relation)
					}
				}
				s.references(constraintsOf(cd.Constraints)...)
			}
		case ast.AT_AddIndexConstraint, ast.AT_AddIndex:
			if idx, ok := cmd.Def.(*ast.IndexStmt); ok {
				s.renamedKeys[idx.Idxname] = true
				if idx.Primary {
					s.keyed(v.Relation)
				}
			}
		case ast.AT_AttachPartition:
			if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
				s.unsettled[pc.Name.Relname] = true
				for _, p := range s.matching(pc.Name) {
					p.settled = true
				}
			}
		case ast.AT_DetachPartition:
			if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
				s.unsettled[pc.Name.Relname] = true
			}
		case ast.AT_AddInherit, ast.AT_DropInherit, ast.AT_AddOf, ast.AT_DropOf:
			s.unsettled[name] = true
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
			s.referencedByChange[t.table] {
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
	column := [2]string{t.table, cmd.Name}
	switch {
	case s.renamedTo[column]:
		// A column renamed to this name: what it was is not followed.
		return false
	case s.renamedFrom[column] && !s.columns[column]:
		return !cmd.Missing_ok
	}
	if !slices.ContainsFunc(table.GetColumns(), func(c *metadata.ColumnMetadata) bool { return c.GetName() == cmd.Name }) {
		// A column the change added is not the key's.
		return !cmd.Missing_ok && !s.columns[[2]string{t.table, cmd.Name}]
	}
	if !cascade && (s.index.referencesColumn(t, cmd.Name, s.dropped) || s.index.readsColumn[columnRef{t.schema, t.table, cmd.Name}] ||
		s.referencedByChange[t.table] || s.readByChange[t.table]) {
		return true
	}
	pk := primaryKey(table)
	if pk == nil || !s.constraintKnown(t.table, pk.GetName()) {
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
			s.touchName(schema, view)
		}
	}
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
