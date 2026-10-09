package review

import (
	"slices"
	"strings"

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
// whole, so what the drops find is reported only when every subcommand is
// known to run. A statement the server may refuse ends the scan: what
// follows it is not known.
//
// pg: src/backend/commands/tablecmds.c — ATController (AT_PASS_DROP first)
func (s *scan) alterTable(st *statement, v *ast.AlterTableStmt) {
	if v.Cmds == nil || v.Relation == nil {
		return
	}
	// ALTER of a relation the target certainly lacks is refused, and
	// under IF EXISTS does nothing; ALTER TABLE of a relation that is no
	// table is refused.
	if s.index != nil && s.missing(v.Relation) {
		if !v.Missing_ok {
			s.stop()
		}
		return
	}
	if ast.ObjectType(v.ObjType) == ast.OBJECT_TABLE {
		if _, kind, ok := s.lookup(v.Relation); ok && (kind == kindView || kind == kindMatView || kind == kindSequence || kind == kindIndex || kind == kindCompositeType) {
			s.stop()
			return
		}
	}
	// A table the change created is no longer what its CREATE said; a
	// relation it created as another kind takes no ALTER TABLE.
	if c, ok := s.madeAt(v.Relation); ok {
		if ast.ObjectType(v.ObjType) == ast.OBJECT_TABLE && c.kind != kindTable && c.kind != kindForeignTable || c.def != nil && c.def.refusesAlter(v) {
			s.stop()
			return
		}
		c.def = nil
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
	// a statement the server refuses drops nothing.
	var out droppedKeys
	withhold := false
	// droppedColumns names the dropped columns as the statement does, and
	// droppedOrigins as the synced table does, "" for one the scan cannot
	// follow there.
	droppedColumns := make(map[string]bool)
	droppedOrigins := make(map[string]bool)
	droppedNames := make(map[string]bool)
	// Every drop of the statement sees the table as it was before the
	// statement; a column drop leaves it unsettled only after them all.
	droppedAColumn := false
	// reopened are the pending tables the statement drops the added key of.
	var reopened []*pendingTable
	for _, cmd := range drops {
		cascade := cmd.Behavior == int(ast.DROP_CASCADE)
		t, known := before, knownBefore && isTable
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropConstraint:
			// A second drop of the same constraint fails, and so may a drop
			// of one an earlier column drop of the statement took with it.
			if droppedNames[cmd.Name] && !cmd.Missing_ok {
				s.stop()
				return
			}
			if len(droppedOrigins) > 0 && !cmd.Missing_ok && knownBefore && s.mayDependOn(before, cmd.Name, droppedOrigins) {
				withhold = true
			}
			droppedNames[cmd.Name] = true
			s.retireReferences(v.Relation, cmd.Name)
			// A constraint the change certainly dropped, or renamed, is not
			// there to drop.
			if known && !cmd.Missing_ok && (s.dropped[[3]string{t.schema, t.table, cmd.Name}] || s.goneConstraints[[3]string{t.schema, t.table, cmd.Name}]) {
				s.stop()
				return
			}
			if synced, ok := s.syncedConstraint(t, cmd.Name); known && ok && s.dropConstraint(st, v, cmd, synced, t, cascade, &out) {
				s.stop()
				return
			}
			s.constraintChanged(schema, name, cmd.Name)
			if cascade {
				s.dropKeyReferences(v.Relation, t, known, cmd.Name)
			}
			for _, p := range s.pending {
				// A key a new foreign key references is not dropped without
				// CASCADE.
				if p.settled && p.key == cmd.Name && s.namesPending(p, v.Relation) &&
					(cascade || !s.newlyReferenced(tableRef{p.schema, v.Relation.Relname}, nil, true)) {
					reopened = append(reopened, p)
				}
			}
		case ast.AT_DropColumn:
			if droppedColumns[cmd.Name] && !cmd.Missing_ok {
				s.stop()
				return
			}
			// DROP COLUMN IF EXISTS of a column the table lacks does nothing.
			if cmd.Missing_ok && known && !s.hasColumn(t, cmd.Name) && !s.columnRenamed(t, cmd.Name) {
				continue
			}
			droppedColumns[cmd.Name] = true
			// Something the change made may depend on the column, which
			// refuses the statement without CASCADE.
			if !cascade && s.columnsUsed[name] {
				s.stop()
				return
			}
			// Dropping a column of the key a table the change made declared
			// drops the key, unless something the change made depends on
			// the column, which refuses the statement without CASCADE.
			for _, p := range s.pending {
				if !p.settled || p.key == "" || !slices.Contains(p.keyColumns, cmd.Name) || slices.Contains(reopened, p) || !s.namesPending(p, v.Relation) {
					continue
				}
				own := tableRef{p.schema, v.Relation.Relname}
				if !cascade && (s.newlyReferencedColumn(own, cmd.Name, p.keyColumns, true) || s.dependsOnNew(s.newReads, own, nil)) {
					s.stop()
					return
				}
				if dropsKeyColumns(v, p) {
					reopened = append(reopened, p)
				}
			}
			// What the column is in the synced table.
			origin, followed := cmd.Name, known
			if known {
				origin, followed = s.originOf(t, cmd.Name)
				if !followed {
					droppedOrigins[""] = true
				} else if origin != "" {
					droppedOrigins[origin] = true
				}
			}
			if known && s.dropColumn(cmd, t, cascade, &out) {
				s.stop()
				return
			}
			droppedAColumn = true
			if known && followed && origin != "" {
				// The table's own foreign keys on the column go with it.
				for _, fk := range s.index.foreignKeys {
					if fk.owner == t && slices.ContainsFunc(fk.local, func(c string) bool { return columnName(c) == origin }) {
						s.constraintChanged(t.schema, t.table, fk.name)
						s.dropped[[3]string{t.schema, t.table, fk.name}] = true
					}
				}
			}
			if cascade {
				sure := known && followed
				s.dropReaders(t, sure, v.Relation, origin)
				s.dropReferences(v.Relation, t, known, func(fk foreignKeyRef) (bool, bool) {
					if !sure {
						return true, false
					}
					certainly := origin != "" && slices.ContainsFunc(fk.columns, func(c string) bool { return columnName(c) == origin })
					return certainly || slices.ContainsFunc(fk.columns, func(c string) bool { return columnName(c) == "" }), certainly
				})
			}
		}
	}
	withhold = withhold || out.uncertain
	if knownBefore && isTable {
		refused, uncertain := s.refuses(before, droppedColumns, droppedOrigins, droppedNames, others)
		// A subcommand certain to fail fails the statement, whatever its
		// drops do.
		if refused {
			s.stop()
			return
		}
		withhold = withhold || uncertain
	}
	// A statement the server may refuse leaves what follows unknown.
	if withhold {
		s.stop()
		return
	}
	if droppedAColumn {
		s.unsettled[[2]string{schema, name}] = true
		for column := range droppedColumns {
			s.setColumn(schema, name, column, false, false)
		}
	}
	for _, p := range reopened {
		p.settled, p.key = false, ""
	}
	for _, name := range out.freed {
		s.free(schema, name)
	}
	for _, f := range out.findings {
		s.r.add(f)
	}
	if out.lost != nil {
		s.lostKey(st, v, out.lost, out.table)
	}
	addedHere := make(map[string]bool)
	for _, cmd := range others {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddConstraint:
			if c, ok := cmd.Def.(*ast.Constraint); ok {
				if c.Conname != "" {
					s.constraintChanged(schema, name, c.Conname)
				}
				switch {
				case !makesIndex(c):
				case c.Conname == "":
					s.generate(v.Relation)
				default:
					// A key's index takes the constraint's name, and a
					// partition's a generated one.
					s.touchName(v.Relation.Schemaname, c.Conname)
					if s.mayHavePartitions(v.Relation) {
						s.generate(nil)
					}
				}
				if c.Indexname != "" {
					// USING INDEX renames the index to the constraint.
					s.touchName(v.Relation.Schemaname, c.Indexname)
					s.touchName(v.Relation.Schemaname, c.Conname)
					s.renamedKeys[c.Indexname] = true
				}
				if c.Contype == ast.CONSTR_PRIMARY {
					// A key on an index the statement names has the index's
					// columns, which the scan does not read.
					var columns []string
					if c.Indexname == "" {
						columns = nameParts(c.Keys)
					}
					s.keyed(v.Relation, s.keyName(v.Relation, c), columns)
				}
				if c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE {
					s.newKey(v.Relation)
				}
				s.references(v.Relation, c)
			}
		case ast.AT_AddColumn:
			if cd, ok := cmd.Def.(*ast.ColumnDef); ok {
				// ADD COLUMN IF NOT EXISTS of a column the table has, or an
				// earlier subcommand added, does nothing, its constraints
				// included.
				if cmd.Missing_ok && (addedHere[cd.Colname] || knownBefore && !droppedColumns[cd.Colname] && s.hasColumn(before, cd.Colname)) {
					continue
				}
				addedHere[cd.Colname] = true
				if columnGeneratesNames(cd) {
					s.generate(v.Relation)
				}
				if s.mayHavePartitions(v.Relation) {
					s.generate(nil)
				}
				// A generated column depends on the columns it reads.
				if slices.ContainsFunc(constraintsOf(cd.Constraints), func(c *ast.Constraint) bool { return c.Contype == ast.CONSTR_GENERATED }) || cd.Generated != 0 {
					s.unsettled[[2]string{schema, name}] = true
				}
				s.setColumn(schema, name, cd.Colname, true, false)
				for _, c := range constraintsOf(cd.Constraints) {
					if c.Conname != "" {
						s.constraintChanged(schema, name, c.Conname)
					}
					if makesIndex(c) && c.Conname != "" {
						s.touchName(v.Relation.Schemaname, c.Conname)
					}
					if c.Contype == ast.CONSTR_PRIMARY {
						s.keyed(v.Relation, s.keyName(v.Relation, c), []string{cd.Colname})
					}
					if c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE {
						s.newKey(v.Relation)
					}
				}
				s.references(v.Relation, constraintsOf(cd.Constraints)...)
			}
		case ast.AT_AddIdentity:
			s.generate(v.Relation)
		case ast.AT_AddIndexConstraint, ast.AT_AddIndex:
			if idx, ok := cmd.Def.(*ast.IndexStmt); ok {
				s.renamedKeys[idx.Idxname] = true
				if idx.Primary {
					s.keyed(v.Relation, "", nil)
				}
				if idx.Primary || idx.Unique {
					s.newKey(v.Relation)
				}
			}
		case ast.AT_AttachPartition:
			if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
				// The partition gets the parent's indexes under generated
				// names.
				s.generate(pc.Name)
				s.unsettle(pc.Name)
				for _, p := range s.matching(pc.Name) {
					p.settle()
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
	// uncertain marks a drop the server may refuse.
	uncertain bool
	// freed lists the relations the drops took with them: a key's index,
	// and a dropped column's indexes and owned sequences.
	freed []string
}

// dropConstraint checks a DROP CONSTRAINT against the synced table, which
// lists the constraint as synced. It returns true when the server refuses
// the statement: the constraint does not exist and IF EXISTS is absent,
// or it is a key a foreign key references and CASCADE is absent.
func (s *scan) dropConstraint(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, synced string, t tableRef, cascade bool, out *droppedKeys) (refused bool) {
	table := s.index.schemas[t.schema].tables[t.table]
	con, ok := constraint(table, synced)
	if !ok {
		// The snapshot leaves out NOT NULL constraints, so a constraint it
		// does not list may still exist, or may not, and the statement may
		// fail: its findings are withheld.
		if !cmd.Missing_ok {
			out.uncertain = true
		}
		return false
	}
	// A check constraint a table that may be its parent has too, under
	// the name and with the expression, may be inherited, which the server
	// does not drop from the child.
	if con.kind == "check constraint" && s.index.mayInheritCheck(t, synced) {
		out.uncertain = true
		return false
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
	// A key's index has the constraint's name.
	if s.index.constraintIndexes[tableRef{t.schema, synced}] {
		out.freed = append(out.freed, cmd.Name)
	}
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
	// column is what the dropped column is in the synced table, through
	// any rename.
	column, ok := s.originOf(t, cmd.Name)
	switch {
	case !ok:
		return false
	case column == "":
		// A column the change added is not the key's; one it dropped is
		// not there.
		exists, _ := s.columnNow(t, cmd.Name)
		return !exists && !cmd.Missing_ok
	}
	if !slices.ContainsFunc(table.GetColumns(), func(c *metadata.ColumnMetadata) bool { return c.GetName() == column }) {
		return !cmd.Missing_ok
	}
	defer func() {
		if refused {
			return
		}
		// The column's indexes and owned sequences go with it.
		for _, i := range table.GetIndexes() {
			if slices.Contains(keyColumns(i), column) && !s.isTouched(&ast.RangeVar{Schemaname: t.schema, Relname: i.GetName()}) {
				out.freed = append(out.freed, i.GetName())
			}
		}
		for _, q := range s.index.ownedSequences[columnRef{t.schema, t.table, column}] {
			if !s.reowned[[2]string{t.schema, q}] && !s.reowned[[2]string{"", q}] && !s.isTouched(&ast.RangeVar{Schemaname: t.schema, Relname: q}) {
				out.freed = append(out.freed, q)
			}
		}
	}()
	pk := primaryKey(table)
	var pkColumns []string
	if pk != nil {
		pkColumns = keyColumns(pk)
	}
	if !cascade && (s.index.referencesColumn(t, column, s.dropped) || generatedFrom(table, column) || slices.ContainsFunc(s.index.readsColumn[columnRef{t.schema, t.table, column}], s.viewHolds) ||
		s.newlyReferencedColumn(t, column, pkColumns, pk != nil) || s.newlyReferencedColumn(t, cmd.Name, pkColumns, pk != nil) || s.dependsOnNew(s.newReads, t, nil)) {
		return true
	}
	if pk == nil {
		return false
	}
	if !s.holdsConstraint(t, pk.GetName()) || !slices.Contains(keyColumns(pk), column) {
		return false
	}
	if out.lost == nil {
		out.lost, out.table = cmd, t
	}
	return false
}

// generatedFrom reports whether another column of the synced table may be
// generated from the column: its generation expression refers to it, or,
// when the scan cannot parse the expression, names it.
func generatedFrom(t *metadata.TableMetadata, column string) bool {
	for _, c := range t.GetColumns() {
		if c.GetName() == column || c.GetGeneration() == nil {
			continue
		}
		expr := c.GetGeneration().GetExpression()
		if refs, ok := expressionColumns(expr); ok {
			if slices.Contains(refs, column) {
				return true
			}
		} else if mentions(expr, column) {
			return true
		}
	}
	return false
}

// expressionColumns parses an expression and returns the columns it
// refers to without a qualifier. ok is false when it does not parse.
func expressionColumns(expr string) ([]string, bool) {
	sql := "SELECT " + expr
	stmts, failed := parse(sql, statementRanges(sql))
	if failed != nil || len(stmts) != 1 {
		return nil, false
	}
	sel, ok := stmts[0].node.(*ast.SelectStmt)
	if !ok || sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return nil, false
	}
	return columnsIn(sel.TargetList.Items[0]), true
}

// mentions reports whether an expression may name a column: the name as
// a whole word, or quoted.
func mentions(expr, column string) bool {
	if strings.Contains(expr, `"`+strings.ReplaceAll(column, `"`, `""`)+`"`) {
		return true
	}
	word := func(b byte) bool {
		return b == '_' || b == '$' || b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= 0x80
	}
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

// dropReaders records the views a cascading DROP COLUMN takes, with the
// views reading them: for a table the scan resolved (known), the views
// reading the column are certainly dropped and their names freed; the
// other views of the table, or every view reading a table of that name
// otherwise, may be, and their names are touched.
func (s *scan) dropReaders(t tableRef, known bool, rv *ast.RangeVar, column string) {
	if s.index == nil {
		return
	}
	var certain, maybe []tableRef
	if known {
		certain = slices.Clone(s.index.readsColumn[columnRef{t.schema, t.table, column}])
		for _, view := range append(slices.Clone(s.index.readers[t]), s.index.mayRead[t]...) {
			if !slices.Contains(certain, view) {
				maybe = append(maybe, view)
			}
		}
	} else {
		for schema := range s.index.schemas {
			if rv.Schemaname == "" || rv.Schemaname == schema {
				t := tableRef{schema, rv.Relname}
				maybe = append(append(maybe, s.index.readers[t]...), s.index.mayRead[t]...)
			}
		}
	}
	s.cascadeViews(certain, maybe)
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
		s.dropReferences(rv, t, true, func(fk foreignKeyRef) (bool, bool) {
			same := sameColumns(fk.columns, con.columns)
			return same, same
		})
	default:
		s.dropReferences(rv, t, true, nil)
	}
}

// mayDependOn reports whether a constraint of the synced table may involve
// one of the columns, so that dropping the column took it: a check or an
// exclusion constraint, whose columns the snapshot does not list, or a
// key or foreign key on one of them.
func (s *scan) mayDependOn(t tableRef, name string, columns map[string]bool) bool {
	// A dropped column the scan cannot follow may be any.
	if columns[""] {
		return true
	}
	if synced, ok := s.syncedConstraint(t, name); ok {
		name = synced
	}
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
// altering a column it lacks, adding a constraint under a name it uses, a
// second primary key, or a key or foreign key on a column it lacks, or a
// foreign key the referenced table has no key for. uncertain reports a
// foreign key whose referenced key the scan cannot check. Other failures,
// such as a type error or rows a new constraint rejects, are not known
// before the statement runs.
//
// The server adds every column of the statement before any constraint,
// and every key before any foreign key, whatever their order.
//
// pg: src/backend/commands/tablecmds.c — ATController (AT_PASS_ADD_COL,
// AT_PASS_ADD_INDEX, AT_PASS_ADD_OTHERCONSTR)
func (s *scan) refuses(t tableRef, droppedColumns, droppedOrigins, droppedConstraints map[string]bool, others []*ast.AlterTableCmd) (refused, uncertain bool) {
	table := s.index.schemas[t.schema].tables[t.table]
	// has reports whether the table has the column once the drops ran.
	has := func(column string) bool {
		return !droppedColumns[column] && s.hasColumn(t, column)
	}
	// uses reports whether a constraint name certainly names a constraint
	// of the table.
	uses := func(name string) bool {
		synced, ok := s.syncedConstraint(t, name)
		if droppedConstraints[name] || !ok || s.dropped[[3]string{t.schema, t.table, name}] {
			return false
		}
		// A column drop of the statement may have taken the constraint.
		if len(droppedOrigins) > 0 && s.mayDependOn(t, name, droppedOrigins) {
			return false
		}
		_, ok = constraint(table, synced)
		return ok
	}
	// holds reports whether the table certainly keeps a constraint of the
	// synced table, under whatever name it has now.
	holds := func(synced string) bool {
		name, ok := s.constraintNow(t, synced)
		return ok && uses(name)
	}
	// The columns and keys the statement adds, and the constraints it adds
	// with them.
	addedColumns := make(map[string]bool)
	var addedKeys [][]string
	addsPrimary, addsIndex := false, false
	var constraints []*ast.Constraint
	for _, cmd := range others {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddColumn:
			cd, ok := cmd.Def.(*ast.ColumnDef)
			if !ok {
				continue
			}
			if has(cd.Colname) || addedColumns[cd.Colname] {
				if cmd.Missing_ok {
					// The column exists: the subcommand does nothing.
					continue
				}
				return true, false
			}
			addedColumns[cd.Colname] = true
			for _, c := range constraintsOf(cd.Constraints) {
				if c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE {
					addedKeys = append(addedKeys, []string{cd.Colname})
				}
				addsPrimary = addsPrimary || c.Contype == ast.CONSTR_PRIMARY
				constraints = append(constraints, c)
			}
		case ast.AT_AddConstraint:
			c, ok := cmd.Def.(*ast.Constraint)
			if !ok {
				continue
			}
			if c.Indexname != "" {
				addsIndex = true
			} else if c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE {
				addedKeys = append(addedKeys, nameParts(c.Keys))
			}
			addsPrimary = addsPrimary || c.Contype == ast.CONSTR_PRIMARY
			constraints = append(constraints, c)
		case ast.AT_AddIndexConstraint, ast.AT_AddIndex:
			addsIndex = true
		}
	}
	// ownKey reports whether the table, once the statement ran, has its
	// primary key (no columns) or a unique key on exactly the columns, for
	// a foreign key referencing its own table. known is false when an
	// index the statement or an earlier one adds may be the key.
	ownKey := func(columns []string) (found, known bool) {
		if addsIndex || droppedOrigins[""] || s.newKeys[[2]string{t.schema, t.table}] || s.newKeys[[2]string{"", t.table}] {
			return false, false
		}
		if len(columns) == 0 {
			pk := primaryKey(table)
			return addsPrimary || pk != nil && holds(pk.GetName()), true
		}
		for _, key := range addedKeys {
			if sameKey(key, columns) {
				return true, true
			}
		}
		for _, i := range table.GetIndexes() {
			if !s.isKey(t, i) || droppedConstraints[i.GetName()] {
				continue
			}
			if key := keyColumns(i); !slices.ContainsFunc(key, func(c string) bool { return droppedOrigins[c] }) && sameKey(key, columns) {
				return true, true
			}
		}
		return false, true
	}
	addedNames := make(map[string]bool)
	addedKey := false
	keyed := func() bool {
		pk := primaryKey(table)
		return addedKey || pk != nil && holds(pk.GetName())
	}
	for _, cmd := range others {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_ColumnDefault, ast.AT_DropNotNull, ast.AT_SetNotNull, ast.AT_AlterColumnType,
			ast.AT_SetStatistics, ast.AT_SetStorage, ast.AT_SetCompression, ast.AT_SetOptions, ast.AT_ResetOptions,
			ast.AT_AddIdentity, ast.AT_SetIdentity, ast.AT_DropIdentity, ast.AT_SetExpression, ast.AT_DropExpression:
			// A column the statement adds may or may not exist when this
			// subcommand's pass runs.
			if cmd.Name == "" || addedColumns[cmd.Name] {
				continue
			}
			if !has(cmd.Name) {
				return true, false
			}
			// A column of the primary key the table keeps stays NOT NULL.
			if ast.AlterTableType(cmd.Subtype) == ast.AT_DropNotNull {
				if pk := primaryKey(table); pk != nil && holds(pk.GetName()) {
					if origin, ok := s.originOf(t, cmd.Name); ok && origin != "" && slices.Contains(keyColumns(pk), origin) {
						return true, false
					}
				}
			}
		}
	}
	for _, c := range constraints {
		// A key or foreign key names columns the table must have, and a
		// check names nothing but its columns.
		for _, list := range []*ast.List{c.Keys, c.FkAttrs} {
			for _, column := range nameParts(list) {
				if !has(column) && !addedColumns[column] {
					return true, false
				}
			}
		}
		if c.Contype == ast.CONSTR_CHECK && slices.ContainsFunc(columnsIn(c.RawExpr), func(column string) bool { return !has(column) && !addedColumns[column] }) {
			return true, false
		}
		// A key names a column once.
		if (c.Contype == ast.CONSTR_PRIMARY || c.Contype == ast.CONSTR_UNIQUE) && hasDuplicate(nameParts(c.Keys)) {
			return true, false
		}
		if c.Conname != "" {
			if uses(c.Conname) || addedNames[c.Conname] {
				return true, false
			}
			addedNames[c.Conname] = true
		}
		// A key's index is the one USING INDEX names, or one under the
		// constraint's name in the table's schema.
		switch {
		case c.Indexname != "":
			r, u := s.refusesIndex(t, c.Indexname, droppedOrigins, droppedConstraints)
			if r {
				return true, false
			}
			uncertain = uncertain || u
			// The index takes the constraint's name.
			if c.Conname != "" && c.Conname != c.Indexname {
				taken, known := s.indexNameTaken(t, c.Conname, droppedOrigins, droppedConstraints)
				if taken {
					return true, false
				}
				uncertain = uncertain || !known
			}
		case makesIndex(c) && c.Conname != "":
			taken, known := s.indexNameTaken(t, c.Conname, droppedOrigins, droppedConstraints)
			if taken {
				return true, false
			}
			uncertain = uncertain || !known
		}
		if c.Contype == ast.CONSTR_PRIMARY {
			if keyed() {
				return true, false
			}
			addedKey = true
		}
		if c.Contype == ast.CONSTR_FOREIGN {
			r, u := s.refusesReference(t, c, func(column string) bool { return has(column) || addedColumns[column] }, ownKey)
			if r {
				return true, false
			}
			uncertain = uncertain || u
		}
	}
	return false, uncertain
}

// refusesReference reports whether the server certainly refuses a foreign
// key the statement adds to t for what it references: a relation the
// target lacks or that is not a table, a referenced column the table
// lacks, a count of referenced columns unlike the referencing ones, or no
// primary key, or no unique key on exactly the columns named. uncertain
// reports a referenced table whose keys the scan cannot tell: one the
// change made or altered, or gave a key or a unique index. own and ownKey
// say whether t has a column and a key once the statement ran, for a
// foreign key referencing its own table.
//
// pg: src/backend/commands/tablecmds.c — ATAddForeignKeyConstraint,
// transformFkeyGetPrimaryKey, transformFkeyCheckAttrs
func (s *scan) refusesReference(t tableRef, c *ast.Constraint, own func(string) bool, ownKey func([]string) (bool, bool)) (refused, uncertain bool) {
	if c.Pktable == nil {
		return false, true
	}
	columns := nameParts(c.PkAttrs)
	// A foreign key of a column names no referencing columns: it is the one.
	local := max(len(nameParts(c.FkAttrs)), 1)
	if len(columns) > 0 && len(columns) != local {
		return true, false
	}
	// No key has a column twice, and no foreign key names one twice.
	if hasDuplicate(columns) || hasDuplicate(nameParts(c.FkAttrs)) {
		return true, false
	}
	if s.missing(c.Pktable) {
		return true, false
	}
	schema, kind, ok := s.lookup(c.Pktable)
	if !ok {
		// A table the change created and said all of is checked against
		// what its CREATE gave it.
		if def, made := s.madeDef(c.Pktable); made {
			return def.refusedBy(columns, local), false
		}
		return false, true
	}
	ref := tableRef{schema, c.Pktable.Relname}
	if ref == t {
		if slices.ContainsFunc(columns, func(column string) bool { return !own(column) }) {
			return true, false
		}
		found, known := ownKey(columns)
		return known && !found, !known
	}
	switch kind {
	case kindTable:
	case kindView, kindMatView, kindForeignTable, kindSequence, kindIndex, kindCompositeType:
		return true, false
	default:
		return false, true
	}
	if s.isUnsettled(ref) {
		return false, true
	}
	if slices.ContainsFunc(columns, func(column string) bool { return !s.hasColumn(ref, column) }) {
		return true, false
	}
	if s.newKeys[[2]string{ref.schema, ref.table}] || s.newKeys[[2]string{"", ref.table}] {
		return false, true
	}
	table := s.index.schemas[ref.schema].tables[ref.table]
	if len(columns) == 0 {
		pk := primaryKey(table)
		if pk == nil || !s.indexHolds(ref, pk) {
			return true, false
		}
		key := keyColumns(pk)
		return key != nil && len(key) != local, false
	}
	for _, i := range table.GetIndexes() {
		if s.isKey(ref, i) && sameKey(keyColumns(i), columns) {
			return false, false
		}
	}
	return true, false
}

// refusesIndex reports whether the server certainly refuses ADD CONSTRAINT
// ... USING INDEX of a synced table's index: no index of the table has
// the name, or the index is a constraint's already, or not unique, or
// partial, or on an expression, or sorted other than ascending, or the
// statement drops it with a column. uncertain reports a name an earlier
// statement touched, or that the snapshot lists twice.
//
// pg: src/backend/parser/parse_utilcmd.c — transformIndexConstraint
func (s *scan) refusesIndex(t tableRef, name string, droppedColumns, droppedConstraints map[string]bool) (refused, uncertain bool) {
	if droppedColumns[""] || s.isTouched(&ast.RangeVar{Schemaname: t.schema, Relname: name}) {
		return false, true
	}
	switch kind, ok := s.index.schemas[t.schema].relations[name]; {
	case ok && kind == kindAmbiguous:
		return false, true
	case !ok || kind != kindIndex:
		return true, false
	}
	table := s.index.schemas[t.schema].tables[t.table]
	for _, i := range table.GetIndexes() {
		if i.GetName() != name {
			continue
		}
		key := keyColumns(i)
		return droppedConstraints[name] || s.index.constraintIndexes[tableRef{t.schema, name}] || !i.GetUnique() || key == nil ||
			strings.Contains(strings.ToUpper(i.GetDefinition()), " WHERE ") || slices.Contains(i.GetDescending(), true) ||
			slices.ContainsFunc(key, func(c string) bool { return droppedColumns[c] }), false
	}
	return true, false
}

// indexNameTaken reports whether the index of a key the statement adds
// under a name cannot take it: a relation of the table's schema has it,
// other than an index the statement drops with its constraint or its
// column. known is false when an earlier statement touched the name, or
// the statement may drop the index of that name.
func (s *scan) indexNameTaken(t tableRef, name string, droppedColumns, droppedConstraints map[string]bool) (taken, known bool) {
	rv := &ast.RangeVar{Schemaname: t.schema, Relname: name}
	switch {
	case s.freed[[2]string{t.schema, name}]:
		return false, true
	case s.isTouched(rv) || droppedColumns[""]:
		return false, false
	}
	if _, ok := s.index.schemas[t.schema].relations[name]; !ok || droppedConstraints[name] {
		return false, true
	}
	for _, i := range s.index.schemas[t.schema].tables[t.table].GetIndexes() {
		if i.GetName() != name || len(droppedColumns) == 0 {
			continue
		}
		key := keyColumns(i)
		if key == nil {
			return false, false
		}
		if slices.ContainsFunc(key, func(c string) bool { return droppedColumns[c] }) {
			return false, true
		}
	}
	return true, true
}

// columnsIn returns the column names an expression refers to without a
// qualifier.
func columnsIn(expr ast.Node) []string {
	var out []string
	ast.Inspect(expr, func(n ast.Node) bool {
		if ref, ok := n.(*ast.ColumnRef); ok && ref.Fields != nil && len(ref.Fields.Items) == 1 {
			if name, ok := ref.Fields.Items[0].(*ast.String); ok {
				out = append(out, name.Str)
			}
		}
		return true
	})
	return out
}

// hasDuplicate reports whether a list names something twice.
func hasDuplicate(names []string) bool {
	for i, name := range names {
		if slices.Contains(names[:i], name) {
			return true
		}
	}
	return false
}

// isKey reports whether a synced table's index is a key a foreign key can
// reference, still there: a primary key or a unique index that is not
// partial.
func (s *scan) isKey(t tableRef, i *metadata.IndexMetadata) bool {
	return (i.GetUnique() || i.GetPrimary()) && s.indexHolds(t, i) && !strings.Contains(strings.ToUpper(i.GetDefinition()), " WHERE ")
}

// sameKey reports whether a key's columns, nil when unknown, are exactly
// the columns, in any order.
func sameKey(key, columns []string) bool {
	return key != nil && len(key) == len(columns) && !slices.ContainsFunc(columns, func(column string) bool { return !slices.Contains(key, column) })
}

// indexHolds reports whether a synced table's index is still there as
// synced: no statement dropped its constraint, or changed or reused its
// name.
func (s *scan) indexHolds(t tableRef, i *metadata.IndexMetadata) bool {
	return !s.dropped[[3]string{t.schema, t.table, i.GetName()}] && s.constraintKnown(t, i.GetName()) &&
		!s.isTouched(&ast.RangeVar{Schemaname: t.schema, Relname: i.GetName()})
}

// hasColumn reports whether the synced table has the column after the
// statements so far: as the last statement that added, dropped, or renamed
// it left it, or as synced.
func (s *scan) hasColumn(t tableRef, column string) bool {
	if exists, changed := s.columnNow(t, column); changed {
		return exists
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
