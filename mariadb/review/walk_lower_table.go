package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// alterTable follows ALTER TABLE. MariaDB judges DROP, RENAME, CHANGE, and
// MODIFY against the table as it was and ADD after them, and resolves
// index columns against the final ones (verified): the drops and renames
// are lowered first, then the rest, then a RENAME TO of the table itself,
// each subcommand a part of its own (walk.Passes).
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
		// A temporary table's columns and keys are not followed.
		return false
	case f == walk.Unknown:
		for _, c := range v.Commands {
			if !alterNoFacts(c) {
				return false
			}
		}
		return true
	case f == walk.Present && (kind == relSequence || b.w.versioned[relKey(r.db, r.name)]):
		// Some ALTER TABLEs of a sequence succeed and others fail (4086),
		// and a system-versioned table refuses most (4119) and hides
		// columns: the walk follows neither (MariaDB spec §2.3; see
		// walker.versioned).
		return false
	}
	b.tableNeeds(r)
	if f == walk.Absent || kind == relView {
		return true
	}
	// The table keeps foreign keys as its engine does; a statement that
	// changes the engine claims nothing of them.
	rel := relKey(r.db, r.name)
	engine, sets := alterEngine(v)
	b.fks = !sets && !b.w.ignoresFKs[rel]
	if sets {
		ignores := !keepsFKs(engine)
		b.Session(func() {
			if ignores {
				b.w.ignoresFKs[rel] = true
			} else {
				delete(b.w.ignoresFKs, rel)
			}
		})
	}
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

// alterEngine returns the ENGINE an ALTER TABLE sets, if it sets one.
func alterEngine(v *ast.AlterTableStmt) (string, bool) {
	for _, c := range v.Commands {
		if c.Type == ast.ATTableOption && c.Option != nil && strings.EqualFold(c.Option.Name, "ENGINE") {
			return engineOf([]*ast.TableOption{c.Option}), true
		}
	}
	return "", false
}

// alterEarly reports whether a subcommand is judged against the table as
// it was: a drop or a rename.
func alterEarly(t ast.AlterTableCmdType) bool {
	switch t {
	case ast.ATDropColumn, ast.ATDropIndex, ast.ATDropConstraint, ast.ATChangeColumn, ast.ATRenameColumn, ast.ATRenameIndex:
		return true
	}
	return false
}

// alterNoFacts reports whether a subcommand changes no name the walk
// keeps: table options, partitions, tablespaces, and CHECK constraints,
// whose names the walk does not follow. MySQL's WITH VALIDATION, ALTER
// CHECK, and SECONDARY_LOAD are not here: omni's parser takes them, and
// MariaDB rejects them (1064, verified), so they end the walk.
func alterNoFacts(c *ast.AlterTableCmd) bool {
	switch c.Type {
	case ast.ATTableOption, ast.ATConvertCharset, ast.ATAlgorithm, ast.ATLock, ast.ATForce, ast.ATOrderBy,
		ast.ATEnableKeys, ast.ATDisableKeys, ast.ATDiscardTablespace, ast.ATImportTablespace,
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
		// IF EXISTS makes a missing column a note (verified).
		if !c.IfExists {
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
	case ast.ATAlterColumnDefault:
		// MySQL's SET [IN]VISIBLE and SET/DROP NOT NULL, which omni's
		// parser takes, MariaDB rejects (1064, verified): they end the
		// walk.
		return b.columnNeeds(r, c.Name, rng)
	case ast.ATAddConstraint, ast.ATAddIndex:
		if c.Constraint == nil {
			return false
		}
		b.addConstraint(r, c.Constraint)
	case ast.ATDropIndex:
		return b.dropIndexName(r, c.Name, rng, c.IfExists)
	case ast.ATDropConstraint:
		switch dropKindOf(b.w.sql, c) {
		case dropPrimaryKey:
			return b.dropIndexName(r, "PRIMARY", rng, c.IfExists)
		case dropForeignKey:
			k, ok := b.key(spaceForeignKey, c.Name)
			if !ok {
				return false
			}
			fk, own := fkKeys(r, k)
			if !b.fks {
				b.Set(fk, walk.Unknown, "")
				b.Set(own, walk.Unknown, "")
				break
			}
			if c.IfExists {
				// A note when the table has no such foreign key, which
				// leaves another table's of the name (verified).
				switch f, _ := b.l.Lookup(own); f {
				case walk.Present:
					b.Set(fk, walk.Absent, "")
				case walk.Unknown:
					b.Set(fk, walk.Unknown, "")
				}
				b.Set(own, walk.Absent, "")
				break
			}
			// The name must be the table's: another table's foreign key of
			// the name fails as a missing one does (1091, verified).
			msg := "Foreign key " + ident(c.Name) + " does not exist in table " + r.written
			b.need(fk, walk.Present, nil, rng, msg, "")  // 1091
			b.need(own, walk.Present, nil, rng, msg, "") // 1091
			b.Set(fk, walk.Absent, "")
			b.Set(own, walk.Absent, "")
		case dropConstraint:
			// The foreign key, unique key, or check of the name.
			if k, ok := b.key(spaceForeignKey, c.Name); ok {
				fk, own := fkKeys(r, k)
				b.Set(fk, walk.Unknown, "")
				b.Set(own, walk.Unknown, "")
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
	default:
		// MySQL's ALTER INDEX VISIBLE and INVISIBLE, which omni's parser
		// takes, MariaDB rejects (1064, verified): they end the walk.
		return alterNoFacts(c)
	}
	return true
}

// alterTextFails fails an ALTER TABLE its text alone refutes: a column
// added or dropped twice (1060, 1091), an index added or dropped twice
// (1061, 1091), two primary keys (1068), an index named PRIMARY (1280), or
// a name longer than 64 characters (1059). A drop with IF EXISTS is not
// counted: a second one is a note (verified).
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
			if !c.IfExists && twice("drop column", spaceColumn, c.Name, rng, "Column "+ident(c.Name)+" is dropped twice from table "+r.written) {
				return true
			}
		case ast.ATDropIndex:
			if !c.IfExists && twice("drop index", spaceIndex, c.Name, rng, "Index "+ident(c.Name)+" is dropped twice from table "+r.written) {
				return true
			}
		case ast.ATDropConstraint:
			if c.IfExists {
				continue
			}
			switch dropKindOf(b.w.sql, c) {
			case dropPrimaryKey:
				if twice("drop index", spaceIndex, "PRIMARY", rng, "Index PRIMARY is dropped twice from table "+r.written) {
					return true
				}
			case dropForeignKey:
				if b.fks && twice("drop foreign key", spaceForeignKey, c.Name, rng, "Foreign key "+ident(c.Name)+" is dropped twice from table "+r.written) {
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
			if b.tooLong(con.Name, rng) {
				return true
			}
			switch con.Type {
			case ast.ConstrPrimaryKey:
				primaries++
			case ast.ConstrForeignKey:
				// 1061, for the index each would add (verified).
				if b.fks && twice("add foreign key", spaceForeignKey, con.Name, rng, "Foreign key "+ident(con.Name)+" is added twice to table "+r.written) {
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

// renameColumn follows CHANGE a b and RENAME COLUMN a TO b.
func (b *builder) renameColumn(r tableRef, from, to string, rng review.Range) bool {
	fk, ok1 := b.key(spaceColumn, from)
	tk, ok2 := b.key(spaceColumn, to)
	if !ok1 || !ok2 {
		return false
	}
	if b.tooLong(to, rng) {
		return true
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
		if k, ok := b.key(spaceForeignKey, c.Name); ok {
			b.need(objKey(spaceForeignKey, r.db, k), walk.Absent, nil, rng, "", "Foreign key "+ident(c.Name)+" already exists in the database") // 1826; 1022 on 5.7
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
// PRIMARY. IF EXISTS makes a missing index a note (verified).
func (b *builder) dropIndexName(r tableRef, name string, rng review.Range, ifExists bool) bool {
	k, ok := b.key(spaceIndex, name)
	if !ok {
		return false
	}
	ik := idxKey(r.db, r.name, k)
	switch {
	case ifExists:
	case k == primaryIndex:
		b.need(ik, walk.Present, nil, rng, "Table "+r.written+" has no primary key", "") // 1091
	default:
		b.need(ik, walk.Present, nil, rng, "Index "+ident(name)+" does not exist in table "+r.written, "") // 1091
	}
	b.leave(ik)
	return true
}

// alterRename follows ALTER TABLE ... RENAME [TO] t2.
func (b *builder) alterRename(r tableRef, c *ast.AlterTableCmd) bool {
	dest := c.NewTable
	if dest == nil {
		dest = &ast.TableRef{Name: c.NewName, Loc: c.Loc}
	}
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

// moveTable records a relation renamed: what it holds goes with it, and
// so, to another database, do its foreign keys and triggers, under names
// the walk does not list.
func (b *builder) moveTable(from, to tableRef, kind string) {
	toPath := tablePath(to.db, to.name)
	b.Move(tablePath(from.db, from.name), toPath)
	b.Set(relKey(from.db, from.name), walk.Absent, "")
	b.Set(relKey(to.db, to.name), walk.Present, kind)
	if from.db != to.db {
		for _, space := range []string{spaceForeignKey, spaceTrigger} {
			b.MayRemove(walk.Scope{Space: space, Parent: walk.Path(from.db)})
			b.MayAdd(walk.Scope{Space: space, Parent: walk.Path(to.db)})
		}
	}
	// The foreign keys named <from>_ibfk_... become <to>_ibfk_...
	// (verified).
	old, renamed := generatedFKs(from.name), generatedFKs(to.name)
	b.MayRemove(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(from.db), Prefix: old})
	b.MayAdd(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(to.db), Prefix: renamed})
	b.MayRemove(walk.Scope{Space: spaceTableForeignKey, Parent: toPath, Prefix: old})
	b.MayAdd(walk.Scope{Space: spaceTableForeignKey, Parent: toPath, Prefix: renamed})
	fk, tk := relKey(from.db, from.name), relKey(to.db, to.name)
	b.Session(func() {
		if b.w.triggered[fk] {
			b.w.triggered[tk] = true
		}
		ignores := b.w.ignoresFKs[fk]
		delete(b.w.ignoresFKs, fk)
		if ignores {
			b.w.ignoresFKs[tk] = true
		} else {
			delete(b.w.ignoresFKs, tk)
		}
		versioned := b.w.versioned[fk]
		delete(b.w.versioned, fk)
		if versioned {
			b.w.versioned[tk] = true
		} else {
			delete(b.w.versioned, tk)
		}
	})
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

// dropTable follows DROP [TEMPORARY] TABLE [IF EXISTS], which looks every
// table up before it drops any (walk.Lookups; MariaDB drops the tables
// that exist and fails on the first that does not, so the statement fails
// either way, verified). Without TEMPORARY a temporary table or sequence
// of the name is dropped instead of the base relation. DROP TABLE drops a
// sequence too (verified).
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
		if b.tooLong(t.Name, r.rng) {
			return true
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
			b.droppedOwned(r)
			continue
		}
		f, kind := b.l.Lookup(k)
		if !v.IfExists {
			// A view is no table to DROP TABLE (1965, verified).
			b.need(k, walk.Present, []string{relTable, relSequence}, r.rng, "Table "+r.written+" does not exist", r.written+" is a view")
		}
		switch {
		case f == walk.Present && kind == relView:
		case f == walk.Present && kind == "", f == walk.Unknown && v.IfExists:
			// It may be a view, which DROP TABLE IF EXISTS leaves with a
			// note (verified).
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
			b.droppedOwned(r)
		default:
			b.Set(k, walk.Absent, "")
			b.Clear(path)
			b.droppedOwned(r)
		}
	}
	return true
}

// droppedOwned records that a dropped table's foreign keys and triggers
// went with it, under names the walk does not list.
func (b *builder) droppedOwned(r tableRef) {
	b.MayRemove(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(r.db)})
	b.MayRemove(walk.Scope{Space: spaceTrigger, Parent: walk.Path(r.db)})
}

// truncate follows TRUNCATE: each table must exist and be a table (1146,
// a view too). TRUNCATE of a sequence fails (1031), which the walk does
// not claim: it ends there (MariaDB spec §2.3).
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	for _, t := range v.Tables {
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			continue
		}
		if f, temp, kind := b.resolve(r); !temp && f != walk.Unknown {
			if kind == relSequence {
				return false
			}
			b.need(relKey(r.db, r.name), walk.Present, []string{relTable}, r.rng, "Table "+r.written+" does not exist", "Table "+r.written+" does not exist")
		}
	}
	return true
}

// createSequence follows CREATE [OR REPLACE] [TEMPORARY] SEQUENCE [IF NOT
// EXISTS] (MariaDB spec §2.3). A sequence is a relation: a table, view, or
// sequence of its name fails it (1050), but for OR REPLACE, which replaces
// a table or a sequence and fails on a view (1965), and IF NOT EXISTS,
// which makes it a note (verified). A temporary sequence shares the
// temporary tables' names (1050) and needs no database, as a temporary
// table does (verified). The walk does not follow a sequence's columns.
func (b *builder) createSequence(v *ast.CreateSequenceStmt) bool {
	r, st := b.ref(v.Name)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if b.tooLong(v.Name.Name, r.rng) {
		return true
	}
	b.Order = walk.Lookups
	if v.Temporary {
		tk := tempKey(r.db, r.name)
		prior, _ := b.l.Lookup(tk)
		switch {
		case v.OrReplace:
		case !v.IfNotExists:
			b.need(tk, walk.Absent, nil, r.rng, "", "Table "+r.written+" already exists") // 1050
		case prior != walk.Absent:
			if prior == walk.Unknown {
				b.Set(tk, walk.Present, "")
			}
			return true
		}
		b.Set(tk, walk.Present, relSequence)
		return true
	}
	if v.Name.Schema != "" {
		b.databaseNeeds(r.db, v.Name.Schema, r.rng) // 1049
	}
	k := relKey(r.db, r.name)
	prior, _ := b.l.Lookup(k)
	switch {
	case v.OrReplace:
		if prior == walk.Present {
			b.need(k, walk.Present, []string{relTable, relSequence}, r.rng, "", r.written+" is a view")
		}
	case !v.IfNotExists:
		b.need(k, walk.Absent, nil, r.rng, "", "Table "+r.written+" already exists") // 1050
		prior = walk.Absent
	case prior == walk.Present:
		return true // a note
	}
	b.Next()
	old := prior == walk.Unknown && !v.OrReplace
	if old {
		b.Optional()
	}
	if v.OrReplace && prior != walk.Absent {
		b.droppedOwned(r)
	}
	b.Set(k, walk.Present, relSequence)
	b.Forget(walk.Scope{Parent: tablePath(r.db, r.name), Deep: true})
	if old {
		b.Next()
		b.Set(k, walk.Present, "")
	}
	return true
}

// dropSequence follows DROP [TEMPORARY] SEQUENCE [IF EXISTS], which, as
// DROP TABLE does, looks every name up first and drops those that exist
// before it fails on one that does not (verified). Each must be a
// sequence: a missing one fails with 4091, a table with 4090, and a view
// with 1965. Without TEMPORARY, a temporary sequence of the name is
// dropped instead of the base one, and a temporary table of the name is
// passed over; with it, only a temporary sequence is found (4091, a
// temporary table too). All verified.
func (b *builder) dropSequence(v *ast.DropSequenceStmt) bool {
	b.Order = walk.Lookups
	for _, t := range v.Sequences {
		b.Next()
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			return false
		}
		if b.tooLong(t.Name, r.rng) {
			return true
		}
		tk := tempKey(r.db, r.name)
		tf, tkind := b.l.Lookup(tk)
		if v.Temporary {
			if !v.IfExists {
				b.need(tk, walk.Present, []string{relSequence}, r.rng, "Temporary sequence "+r.written+" does not exist", "Temporary table "+r.written+" is not a sequence") // 4091
			}
			switch {
			case tf == walk.Present && tkind == relSequence:
				b.Set(tk, walk.Absent, "")
			case tf == walk.Present && tkind == "", tf == walk.Unknown:
				b.Set(tk, walk.Unknown, "")
			}
			continue
		}
		k := relKey(r.db, r.name)
		path := walk.Scope{Parent: tablePath(r.db, r.name), Deep: true}
		switch {
		case tf == walk.Present && tkind == relSequence:
			b.Set(tk, walk.Absent, "")
			continue
		case tf == walk.Present && tkind == "", tf == walk.Unknown:
			b.Set(tk, walk.Unknown, "")
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
			continue
		}
		f, kind := b.l.Lookup(k)
		if !v.IfExists {
			present := r.written + " is not a sequence" // 4090
			if kind == relView {
				present = r.written + " is a view" // 1965
			}
			b.need(k, walk.Present, []string{relSequence}, r.rng, "Sequence "+r.written+" does not exist", present) // 4091
		}
		switch {
		case f == walk.Present && kind != relSequence && kind != "":
			// A note with IF EXISTS, which drops nothing.
		case f == walk.Present && kind == "", f == walk.Unknown && v.IfExists:
			b.Set(k, walk.Unknown, "")
			b.Forget(path)
		default:
			b.Set(k, walk.Absent, "")
			b.Clear(path)
		}
	}
	return true
}
