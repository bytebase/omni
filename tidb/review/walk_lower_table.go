package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// alterTable follows ALTER TABLE. TiDB judges every subcommand against
// the table as it was, and runs none first (spec §3.3, verified: ADD a,
// DROP a fails with 1060 where a exists, a swap of two columns with 1060,
// and AFTER a column the statement adds with 1054): every check holds
// against the state before the statement (walk.Lookups). Each subcommand
// is a part of its own. The effects apply drops and renames first, then
// the rest, as on MySQL: two subcommands on one column or index fail
// (8200), so the order matters only to the scopes a dropped column's
// indexes leave.
func (b *builder) alterTable(v *ast.AlterTableStmt) bool {
	r, st := b.ref(v.Table)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	f, temp, kind := b.resolve(r)
	switch {
	case temp:
		// TiDB refuses ALTER TABLE of a local temporary table (8200), which
		// the walk does not claim.
		return false
	case f == walk.Unknown:
		for _, c := range v.Commands {
			if !alterNoFacts(c) {
				return false
			}
		}
		return true
	}
	// A view or a sequence fails with 1347, whatever the subcommands
	// (verified).
	b.tableNeeds(r)
	if f == walk.Absent || kind == relView || kind == relSequence {
		return true
	}
	b.Order = walk.Lookups
	b.fks = true
	if b.alterTextFails(v, r) {
		return true
	}
	var rename *ast.AlterTableCmd
	for _, early := range []bool{true, false} {
		for _, c := range v.Commands {
			if c.Type == ast.ATRenameTable {
				rename = c
				continue
			}
			if alterEarly(c.Type) != early {
				continue
			}
			b.Next()
			b.early = early
			if !b.alterCmd(r, c) {
				return false
			}
		}
	}
	if rename != nil {
		b.Next()
		return b.alterRename(r, rename)
	}
	return true
}

// alterEarly reports whether a subcommand's effects apply first: a drop or
// a rename.
func alterEarly(t ast.AlterTableCmdType) bool {
	switch t {
	case ast.ATDropColumn, ast.ATDropIndex, ast.ATDropConstraint, ast.ATChangeColumn, ast.ATRenameColumn, ast.ATRenameIndex:
		return true
	}
	return false
}

// alterNoFacts reports whether a subcommand changes no name the walk
// keeps: table options, partitions, tablespaces, and CHECK constraints,
// whose names the walk does not follow.
func alterNoFacts(c *ast.AlterTableCmd) bool {
	switch c.Type {
	case ast.ATTableOption, ast.ATConvertCharset, ast.ATAlgorithm, ast.ATLock, ast.ATForce, ast.ATOrderBy,
		ast.ATEnableKeys, ast.ATDisableKeys, ast.ATWithValidation, ast.ATWithoutValidation, ast.ATAlterCheckEnforced,
		ast.ATSecondaryLoad, ast.ATSecondaryUnload, ast.ATDiscardTablespace, ast.ATImportTablespace,
		ast.ATAddPartition, ast.ATDropPartition, ast.ATCoalescePartition, ast.ATReorganizePartition,
		ast.ATExchangePartition, ast.ATTruncatePartition, ast.ATAnalyzePartition, ast.ATCheckPartition,
		ast.ATOptimizePartition, ast.ATRebuildPartition, ast.ATRepairPartition, ast.ATDiscardPartitionTablespace,
		ast.ATImportPartitionTablespace, ast.ATRemovePartitioning, ast.ATPartitionBy:
		return true
	case ast.ATAddConstraint:
		return c.Constraint != nil && c.Constraint.Type == ast.ConstrCheck
	}
	return false
}

func (b *builder) alterCmd(r tableRef, c *ast.AlterTableCmd) bool {
	rng := rangeOf(c.Loc)
	switch c.Type {
	case ast.ATAddColumn:
		cols := c.Columns
		if c.Column != nil {
			cols = append([]*ast.ColumnDef{c.Column}, cols...)
		}
		for _, col := range cols {
			if !b.addColumn(r, col, c.After, rangeOf(col.Loc)) {
				return false
			}
		}
	case ast.ATDropColumn:
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			return false
		}
		ck := colKey(r.db, r.name, k)
		if c.IfExists {
			// DROP COLUMN IF EXISTS of a missing column is a note, and the
			// statement goes on without it: ADD COLUMN zz, DROP COLUMN IF
			// EXISTS zz adds zz (verified).
			switch f, _ := b.l.Lookup(ck); f {
			case walk.Absent:
				return true
			case walk.Unknown:
				b.Set(ck, walk.Unknown, "")
				b.MayRemove(walk.Scope{Space: spaceIndex, Parent: tablePath(r.db, r.name)})
				return true
			}
		} else {
			b.need(ck, walk.Present, nil, rng, "Column "+ident(c.Name)+" does not exist in table "+r.written, "") // 1091
		}
		b.leave(ck)
		// An index loses the column, and goes with its last one.
		b.MayRemove(walk.Scope{Space: spaceIndex, Parent: tablePath(r.db, r.name)})
	case ast.ATModifyColumn:
		if c.Column == nil || !b.columnNeeds(r, c.Column.Name, rng) {
			return false
		}
		b.after(r, c.After, rng)
		b.alterColumnKeys(r, c.Column, rng)
	case ast.ATChangeColumn:
		if c.Column == nil || !b.renameColumn(r, c.Name, c.Column.Name, rng) {
			return false
		}
		b.after(r, c.After, rng)
		b.alterColumnKeys(r, c.Column, rng)
	case ast.ATRenameColumn:
		return b.renameColumn(r, c.Name, c.NewName, rng)
	case ast.ATAlterColumnDefault, ast.ATAlterColumnSetNotNull, ast.ATAlterColumnDropNotNull,
		ast.ATAlterColumnVisible, ast.ATAlterColumnInvisible:
		return b.columnNeeds(r, c.Name, rng)
	case ast.ATAddConstraint, ast.ATAddIndex:
		if c.Constraint == nil {
			return false
		}
		b.addConstraint(r, c.Constraint)
	case ast.ATDropIndex:
		return b.dropIndexName(r, c.Name, rng)
	case ast.ATDropConstraint:
		switch dropKindOf(b.w.sql, c) {
		case dropPrimaryKey:
			return b.dropIndexName(r, "PRIMARY", rng)
		case dropForeignKey:
			k, ok := b.key(spaceTableForeignKey, c.Name)
			if !ok {
				return false
			}
			own := fkKey(r, k)
			// The name must be the table's (1091, verified).
			b.need(own, walk.Present, nil, rng, "Foreign key "+ident(c.Name)+" does not exist in table "+r.written, "") // 1091
			b.Set(own, walk.Absent, "")
		case dropCheck:
		case dropConstraint:
			// TiDB drops only a check of the name, a note while checks are
			// off (verified); the walk claims nothing of it, and leaves the
			// foreign key and index of the name Unknown, as on MySQL.
			if k, ok := b.key(spaceTableForeignKey, c.Name); ok {
				b.Set(fkKey(r, k), walk.Unknown, "")
				b.Set(idxKey(r.db, r.name, k), walk.Unknown, "")
			}
		default:
			return false
		}
	case ast.ATRenameIndex:
		from, ok1 := b.key(spaceIndex, c.Name)
		to, ok2 := b.key(spaceIndex, c.NewName)
		if !ok1 || !ok2 || from == primaryIndex || to == primaryIndex {
			return false
		}
		fk, tk := idxKey(r.db, r.name, from), idxKey(r.db, r.name, to)
		b.need(fk, walk.Present, nil, rng, "Index "+ident(c.Name)+" does not exist in table "+r.written, "") // 1091
		if from != to {
			b.need(tk, walk.Absent, nil, rng, "", "Index "+ident(c.NewName)+" already exists in table "+r.written) // 1061
			b.leave(fk)
			b.arrive(tk)
		}
	case ast.ATAlterIndexVisible, ast.ATAlterIndexInvisible:
		k, ok := b.key(spaceIndex, c.Name)
		if !ok {
			return false
		}
		b.need(idxKey(r.db, r.name, k), walk.Present, nil, rng, "Index "+ident(c.Name)+" does not exist in table "+r.written, "")
	default:
		return alterNoFacts(c)
	}
	return true
}

// alterTextFails fails an ALTER TABLE its text alone refutes: a column
// added or dropped twice, an index added or dropped twice, two primary
// keys, all with 8200 on TiDB ("Unsupported operate same column", verified),
// an index named PRIMARY (1280), a name longer than 64 characters (1059),
// or RENAME TO beside another subcommand (8200, verified).
func (b *builder) alterTextFails(v *ast.AlterTableStmt, r tableRef) bool {
	seen := make(map[[2]string]bool)
	twice := func(kind, space, name string, rng review.Range, message string) bool {
		k, ok := b.key(space, name)
		if !ok {
			return false
		}
		if seen[[2]string{kind, k}] {
			b.fail(rng, nil, message)
			return true
		}
		seen[[2]string{kind, k}] = true
		return false
	}
	// RENAME TO takes no other subcommand but ALGORITHM and LOCK, which
	// TiDB sets aside first (verified: COMMENT, ENGINE, FORCE, ADD COLUMN,
	// RENAME INDEX, and a second RENAME TO fail with 8200). A RENAME TO of
	// the table's own name, in any case, changes nothing and takes any
	// (verified).
	renames, others := false, 0
	for _, c := range v.Commands {
		switch c.Type {
		case ast.ATRenameTable:
			if k, ok := b.key(spaceRelation, c.NewName); !ok || k == r.name {
				continue
			}
			if renames {
				others++
			}
			renames = true
		case ast.ATAlgorithm, ast.ATLock:
		default:
			others++
		}
	}
	if renames && others > 0 {
		b.fail(b.anchor(), nil, "RENAME TO cannot be combined with other changes to table "+r.written)
		return true
	}
	primaries := 0
	for _, c := range v.Commands {
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAddColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if b.tooLong(col.Name, rangeOf(col.Loc)) ||
					twice("add column", spaceColumn, col.Name, rangeOf(col.Loc), "Column "+ident(col.Name)+" is added twice to table "+r.written) {
					return true
				}
				for _, cc := range col.Constraints {
					if cc.Type == ast.ColConstrPrimaryKey {
						primaries++
					}
				}
			}
		case ast.ATDropColumn:
			// DROP COLUMN IF EXISTS twice of a missing column is two notes
			// (verified).
			if !c.IfExists && twice("drop column", spaceColumn, c.Name, rng, "Column "+ident(c.Name)+" is dropped twice from table "+r.written) {
				return true
			}
		case ast.ATDropIndex:
			if twice("drop index", spaceIndex, c.Name, rng, "Index "+ident(c.Name)+" is dropped twice from table "+r.written) {
				return true
			}
		case ast.ATDropConstraint:
			switch dropKindOf(b.w.sql, c) {
			case dropPrimaryKey:
				if twice("drop index", spaceIndex, "PRIMARY", rng, "Index PRIMARY is dropped twice from table "+r.written) {
					return true
				}
			case dropForeignKey:
				if b.fks && twice("drop foreign key", spaceTableForeignKey, c.Name, rng, "Foreign key "+ident(c.Name)+" is dropped twice from table "+r.written) {
					return true
				}
			}
		case ast.ATChangeColumn, ast.ATRenameColumn, ast.ATModifyColumn:
			// A column a CHANGE or RENAME names is one more the statement
			// adds (1060, verified).
			to := c.NewName
			if c.Column != nil {
				to = c.Column.Name
				for _, cc := range c.Column.Constraints {
					if cc.Type == ast.ColConstrPrimaryKey {
						primaries++
					}
				}
			}
			from, ok1 := b.key(spaceColumn, c.Name)
			dest, ok2 := b.key(spaceColumn, to)
			if c.Type != ast.ATModifyColumn && ok1 && ok2 && from != dest &&
				twice("add column", spaceColumn, to, rng, "Column "+ident(to)+" is added twice to table "+r.written) {
				return true
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			con := c.Constraint
			if con == nil {
				continue
			}
			// A CHECK constraint's name may be long (namesTooLong).
			if con.Type != ast.ConstrCheck && b.tooLong(keyName(con), rng) {
				return true
			}
			switch con.Type {
			case ast.ConstrPrimaryKey:
				primaries++
			case ast.ConstrForeignKey:
				// 1826 in CREATE TABLE (verified).
				if b.fks && twice("add foreign key", spaceTableForeignKey, keyName(con), rng, "Foreign key "+ident(keyName(con))+" is added twice to table "+r.written) {
					return true
				}
			case ast.ConstrCheck:
			default:
				if k, ok := b.key(spaceIndex, con.Name); ok && k == primaryIndex {
					b.fail(rng, nil, "Index name PRIMARY is reserved for the primary key")
					return true
				}
				if twice("add index", spaceIndex, con.Name, rng, "Index "+ident(con.Name)+" is added twice to table "+r.written) {
					return true
				}
			}
		}
	}
	if primaries > 1 {
		b.fail(b.anchor(), nil, "Multiple primary keys for table "+r.written)
		return true
	}
	return false
}

// addColumn follows ADD COLUMN: the column must be new (1060), and a
// column AFTER names must exist (1054).
func (b *builder) addColumn(r tableRef, col *ast.ColumnDef, after string, rng review.Range) bool {
	k, ok := b.key(spaceColumn, col.Name)
	if !ok {
		return false
	}
	ck := colKey(r.db, r.name, k)
	b.need(ck, walk.Absent, nil, rng, "", "Column "+ident(col.Name)+" already exists in table "+r.written) // 1060
	b.after(r, after, rng)
	b.Set(ck, walk.Present, "")
	b.alterColumnKeys(r, col, rng)
	return true
}

// after checks the column a FIRST/AFTER clause names (1054).
func (b *builder) after(r tableRef, after string, rng review.Range) {
	if k, ok := b.key(spaceColumn, after); ok {
		b.need(colKey(r.db, r.name, k), walk.Present, nil, rng, "Column "+ident(after)+" does not exist in table "+r.written, "") // 1054
	}
}

// columnNeeds checks that a column the subcommand names exists (1054).
func (b *builder) columnNeeds(r tableRef, name string, rng review.Range) bool {
	k, ok := b.key(spaceColumn, name)
	if !ok {
		return false
	}
	b.need(colKey(r.db, r.name, k), walk.Present, nil, rng, "Column "+ident(name)+" does not exist in table "+r.written, "") // 1054
	return true
}

// renameColumn follows CHANGE a b and RENAME COLUMN a TO b. TiDB checks
// no length here: CHANGE and RENAME COLUMN to a name of 65 characters
// succeed (verified).
func (b *builder) renameColumn(r tableRef, from, to string, rng review.Range) bool {
	fk, ok1 := b.key(spaceColumn, from)
	tk, ok2 := b.key(spaceColumn, to)
	if !ok1 || !ok2 {
		return false
	}
	fc, tc := colKey(r.db, r.name, fk), colKey(r.db, r.name, tk)
	b.need(fc, walk.Present, nil, rng, "Column "+ident(from)+" does not exist in table "+r.written, "") // 1054
	if fk != tk {
		b.need(tc, walk.Absent, nil, rng, "", "Column "+ident(to)+" already exists in table "+r.written) // 1060
		b.leave(fc)
		b.arrive(tc)
	}
	return true
}

// alterColumnKeys records the keys a column definition of ALTER TABLE
// declares inline: a primary key the table must not have yet (1068).
func (b *builder) alterColumnKeys(r tableRef, col *ast.ColumnDef, rng review.Range) {
	for _, cc := range col.Constraints {
		if cc.Type == ast.ColConstrPrimaryKey {
			b.need(idxKey(r.db, r.name, primaryIndex), walk.Absent, nil, rng, "", "Table "+r.written+" already has a primary key") // 1068
		}
	}
	b.columnKeys(r, col)
}

// addConstraint follows ADD of a key, index, or foreign key.
func (b *builder) addConstraint(r tableRef, c *ast.Constraint) {
	rng := rangeOf(c.Loc)
	switch c.Type {
	case ast.ConstrPrimaryKey:
		b.need(idxKey(r.db, r.name, primaryIndex), walk.Absent, nil, rng, "", "Table "+r.written+" already has a primary key") // 1068
	case ast.ConstrCheck:
		return
	case ast.ConstrForeignKey:
		if !b.fks {
			break
		}
		// The name must be new on the table; another table's is taken
		// (verified).
		if k, ok := b.key(spaceTableForeignKey, keyName(c)); ok {
			b.need(fkKey(r, k), walk.Absent, nil, rng, "", "Foreign key "+ident(keyName(c))+" already exists in table "+r.written) // 1826
		}
		b.foreignKeyTarget(r, c)
	default:
		if k, ok := b.key(spaceIndex, c.Name); ok {
			b.need(idxKey(r.db, r.name, k), walk.Absent, nil, rng, "", "Index "+ident(c.Name)+" already exists in table "+r.written) // 1061
		}
	}
	b.constraint(r, c)
}

// dropIndexName follows DROP INDEX, of the primary key when the name is
// PRIMARY.
func (b *builder) dropIndexName(r tableRef, name string, rng review.Range) bool {
	k, ok := b.key(spaceIndex, name)
	if !ok {
		return false
	}
	ik := idxKey(r.db, r.name, k)
	if k == primaryIndex {
		b.need(ik, walk.Present, nil, rng, "Table "+r.written+" has no primary key", "") // 1091
	} else {
		b.need(ik, walk.Present, nil, rng, "Index "+ident(name)+" does not exist in table "+r.written, "") // 1091
	}
	b.leave(ik)
	return true
}

// alterRename follows ALTER TABLE ... RENAME [TO] t2.
func (b *builder) alterRename(r tableRef, c *ast.AlterTableCmd) bool {
	// TiDB's AST gives RENAME TO only a name: the qualified form is a
	// syntax error in its parser, and a quoted `db.t` is one name.
	dest := &ast.TableRef{Name: c.NewName, Loc: c.Loc}
	to, st := b.ref(dest)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if b.tooLong(dest.Name, to.rng) {
		return true
	}
	if to.db == r.db && to.name == r.name {
		return true
	}
	if dest.Schema != "" {
		b.databaseNeeds(to.db, dest.Schema, to.rng) // 1049; 1025 on 5.7
	}
	b.need(relKey(to.db, to.name), walk.Absent, nil, to.rng, "", "Table "+to.written+" already exists") // 1050
	b.moveTable(r, to, relTable)
	return true
}

// moveTable records a relation renamed: what it holds goes with it, its
// foreign keys under their names (verified: RENAME TABLE child TO kid
// keeps fk_1).
func (b *builder) moveTable(from, to tableRef, kind string) {
	b.Move(tablePath(from.db, from.name), tablePath(to.db, to.name))
	b.Set(relKey(from.db, from.name), walk.Absent, "")
	b.Set(relKey(to.db, to.name), walk.Present, kind)
}

// renameTable follows RENAME TABLE, whose pairs run in order: a swap
// through a third name succeeds (verified). Each pair is a part, and a
// pair's check on a name another pair changes is left out.
func (b *builder) renameTable(v *ast.RenameTableStmt) bool {
	moved := make(map[walk.Key]string)
	for _, p := range v.Pairs {
		b.Next()
		from, st1 := b.ref(p.Old)
		to, st2 := b.ref(p.New)
		switch {
		case st1 == refFailed || st2 == refFailed:
			return true
		case st1 == refUnknown || st2 == refUnknown:
			return false
		}
		if b.tooLong(p.New.Name, to.rng) {
			return true
		}
		fk := relKey(from.db, from.name)
		f, temp, kind := b.resolve(from)
		if k, ok := moved[fk]; ok {
			f, temp, kind = walk.Present, false, k
		}
		if _, tt, _ := b.resolve(to); temp || tt || f == walk.Unknown {
			return false
		}
		b.need(fk, walk.Present, nil, from.rng, "Table "+from.written+" does not exist", "") // 1146; 1017 on 5.7
		if p.New.Schema != "" {
			b.databaseNeeds(to.db, p.New.Schema, to.rng) // 1049; 1025 on 5.7
		}
		if kind == relView && from.db != to.db {
			b.fail(to.rng, []walk.Key{fk}, "View "+from.written+" cannot move to another database") // 1450 (verified)
			return true
		}
		b.need(relKey(to.db, to.name), walk.Absent, nil, to.rng, "", "Table "+to.written+" already exists") // 1050
		b.moveTable(from, to, kind)
		delete(moved, fk)
		moved[relKey(to.db, to.name)] = kind
	}
	return true
}

// dropTable follows DROP [TEMPORARY] TABLE [IF EXISTS], judged against
// the state before it (walk.Lookups). TiDB drops every base table of the
// list it can before it fails on a missing name, a view, or a sequence
// (verified); the statement fails either way. Without
// TEMPORARY a temporary table of the name is dropped instead of the base
// table. TiDB checks no length here: DROP TABLE IF EXISTS of a name of 65
// characters is a note (verified).
func (b *builder) dropTable(v *ast.DropTableStmt) bool {
	b.Order = walk.Lookups
	for _, t := range v.Tables {
		b.Next()
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			return false
		}
		tk := tempKey(r.db, r.name)
		tf, _ := b.l.Lookup(tk)
		if v.Temporary {
			if !v.IfExists {
				b.need(tk, walk.Present, nil, r.rng, "Temporary table "+r.written+" does not exist", "") // 1051
			}
			b.Set(tk, walk.Absent, "")
			continue
		}
		path := walk.Scope{Parent: tablePath(r.db, r.name), Deep: true}
		k := relKey(r.db, r.name)
		switch tf {
		case walk.Present:
			b.Set(tk, walk.Absent, "")
			continue
		case walk.Unknown:
			b.Set(tk, walk.Unknown, "")
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
			continue
		}
		f, kind := b.l.Lookup(k)
		if !v.IfExists {
			// A view is no table to DROP TABLE (1051).
			b.need(k, walk.Present, []string{relTable}, r.rng, "Table "+r.written+" does not exist", "Table "+r.written+" does not exist")
		}
		switch {
		case f == walk.Present && (kind == relView || kind == relSequence):
			// DROP TABLE IF EXISTS leaves a view or a sequence with a note
			// (verified).
		case f == walk.Present && kind == "", f == walk.Unknown && v.IfExists:
			// It may be a view or a sequence.
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
		default:
			b.Set(k, walk.Absent, "")
			b.Clear(path)
		}
	}
	return true
}

// truncate follows TRUNCATE: each table must exist and be a table (1146,
// a view or a sequence too, verified).
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	for _, t := range v.Tables {
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			continue
		}
		if f, temp, _ := b.resolve(r); !temp && f != walk.Unknown {
			b.need(relKey(r.db, r.name), walk.Present, []string{relTable}, r.rng, "Table "+r.written+" does not exist", "Table "+r.written+" does not exist")
		}
	}
	return true
}
