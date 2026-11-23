package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// alterClauses reads the ALGORITHM and LOCK clauses of an ALTER TABLE,
// upper case, "" when not written; the last one written counts (verified
// on mariadb:11.8.8).
func alterClauses(v *ast.AlterTableStmt) (algorithm, lock string) {
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATAlgorithm:
			algorithm = strings.ToUpper(c.Name)
		case ast.ATLock:
			lock = strings.ToUpper(c.Name)
		}
	}
	return algorithm, lock
}

// unsupported is the text for an operation InnoDB can't run with the
// ALGORITHM or LOCK the statement asks for (1845, 1846).
func unsupported(clause, operation string) string {
	return "InnoDB can't " + operation + " with " + review.Code(clause) + "."
}

// columnsOf returns the columns an ADD, MODIFY, or CHANGE defines.
func columnsOf(c *ast.AlterTableCmd) []*ast.ColumnDef {
	cols := c.Columns
	if c.Column != nil {
		cols = append([]*ast.ColumnDef{c.Column}, cols...)
	}
	return cols
}

// colDeclares reports whether a column definition declares a constraint
// of a type.
func colDeclares(col *ast.ColumnDef, t ast.ColumnConstraintType) bool {
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == t {
			return true
		}
	}
	return false
}

// fixedTypes are the base types whose change between two of them the
// walk judges: a change from one to another of them is a change of type
// to MariaDB. FLOAT and DOUBLE are left out, since FLOAT(p) is DOUBLE
// for p over 24, and so is every type the walk may not name as MariaDB
// lists it.
var fixedTypes = map[string]bool{
	"tinyint": true, "smallint": true, "mediumint": true, "int": true, "bigint": true, "decimal": true,
	"bit": true, "char": true, "varchar": true, "binary": true, "varbinary": true,
	"tinytext": true, "text": true, "mediumtext": true, "longtext": true,
	"tinyblob": true, "blob": true, "mediumblob": true, "longblob": true,
	"date": true, "time": true, "datetime": true, "timestamp": true, "year": true,
	"enum": true, "set": true, "json": true,
}

// typeChanges reports whether a column of base type from changes to one
// of base type to, as far as the walk tells.
func typeChanges(from, to string) bool {
	return fixedTypes[from] && fixedTypes[to] && from != to
}

// builtinEngines are the storage engines mariadb:11.8.8 has built in,
// by every name it takes for them.
var builtinEngines = map[string]bool{
	"innodb": true, "innobase": true, "myisam": true, "aria": true, "maria": true, "memory": true,
	"heap": true, "csv": true, "mrg_myisam": true, "merge": true, "sequence": true, "performance_schema": true,
}

// namedOnce reports whether only one subcommand of an ALTER TABLE names a
// column, as the one it adds, drops, changes, or renames.
func (b *builder) namedOnce(v *ast.AlterTableStmt, k string) bool {
	return b.namedCount(v, k) == 1
}

// namedCount returns how many subcommands of an ALTER TABLE name a column,
// as the one it adds, drops, changes, or renames; a name the walk can't
// key counts for every column.
func (b *builder) namedCount(v *ast.AlterTableStmt, k string) int {
	n := 0
	is := func(name string) bool {
		nk, ok := b.key(spaceColumn, name)
		return !ok || nk == k
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		hit := false
		switch c.Type {
		case ast.ATAddColumn, ast.ATModifyColumn:
			for _, col := range columnsOf(c) {
				hit = hit || col != nil && is(col.Name)
			}
		case ast.ATChangeColumn:
			hit = is(c.Name) || c.Column != nil && is(c.Column.Name)
		case ast.ATRenameColumn:
			hit = is(c.Name) || is(c.NewName)
		case ast.ATDropColumn, ast.ATAlterColumnDefault:
			hit = is(c.Name)
		}
		if hit {
			n++
		}
	}
	return n
}

// alterAlgorithmFails fails an ALTER TABLE of a synced InnoDB table the
// change has not touched whose ALGORITHM or LOCK clause can't run what it
// does (each verified on mariadb:11.8.8):
//
//   - ALGORITHM=INSTANT adding or dropping an index (1846), a FULLTEXT one
//     (1845), or changing ROW_FORMAT (1846);
//   - ALGORITHM=INSTANT or NOCOPY rebuilding the table: FORCE, or ENGINE of
//     its own engine (1845);
//   - a change of ENGINE with INSTANT, NOCOPY, or INPLACE (1845);
//   - with INPLACE, NOCOPY, or INSTANT: a column type change (1846),
//     dropping the primary key without adding one (1846), adding a CHECK
//     constraint (1845), adding a STORED generated column (1845), and adding
//     a foreign key while foreign_key_checks is on (1846);
//   - LOCK=NONE adding an AUTO_INCREMENT column (1846).
//
// MariaDB runs ALGORITHM=COPY with LOCK=NONE, and every other LOCK the
// walk sees, online.
func (b *builder) alterAlgorithmFails(v *ast.AlterTableStmt, r tableRef) bool {
	algorithm, lock := alterClauses(v)
	if algorithm == "" && lock == "" {
		return false
	}
	st := b.syncedOf(r)
	if st == nil || st.engine != "innodb" {
		return false
	}
	keys := []walk.Key{relKey(r.db, r.name)}
	fail := func(clause, op string) bool {
		b.fail(b.anchor(), keys, review.UnsupportedForm, unsupported(clause, op))
		return true
	}
	instant := algorithm == "INSTANT"
	noCopy := algorithm == "NOCOPY" || instant
	inPlace := algorithm == "INPLACE" || noCopy
	addsPrimary, dropsPrimary, dropsIndex, addsIndex := false, false, false, false
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATDropIndex, ast.ATDropConstraint:
			switch dropKindOf(b.w.sql, c) {
			case dropPrimaryKey:
				dropsPrimary = true
			case dropIndex:
				dropsIndex = true
			case dropConstraint:
				// It may drop a unique key.
				dropsIndex = true
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			if c.Constraint != nil && c.Constraint.Type == ast.ConstrPrimaryKey {
				addsPrimary = true
			}
			if c.Constraint != nil && c.Constraint.Type != ast.ConstrCheck && c.Constraint.Type != ast.ConstrForeignKey {
				addsIndex = true
			}
		case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
			for _, col := range columnsOf(c) {
				if col != nil && colDeclares(col, ast.ColConstrPrimaryKey) {
					addsPrimary = true
				}
			}
		}
	}
	if engine, sets := alterEngine(v); sets {
		e := strings.ToLower(engine)
		if e == "innobase" {
			e = "innodb"
		}
		if e != "innodb" {
			// Under a sql_mode the session set, MariaDB may put its default
			// engine in place of one it lacks.
			if inPlace && (builtinEngines[e] || !b.w.sqlModeSet && mariadbEngines[e]) {
				return fail("ALGORITHM="+algorithm, "change the table's engine")
			}
			return false
		}
		if noCopy {
			return fail("ALGORITHM="+algorithm, "rebuild the table")
		}
	}
	if dropsPrimary && !addsPrimary && inPlace {
		if f, _ := b.l.Lookup(idxKey(r.db, r.name, primaryIndex)); f == walk.Present {
			return fail("ALGORITHM="+algorithm, "drop the primary key without adding one")
		}
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATForce:
			if noCopy {
				return fail("ALGORITHM="+algorithm, "rebuild the table")
			}
		case ast.ATTableOption:
			if instant && c.Option != nil && optionName(c.Option) == "ROW_FORMAT" {
				return fail("ALGORITHM=INSTANT", "change "+review.Code("ROW_FORMAT"))
			}
		case ast.ATDropIndex:
			// An index dropped and added back as it was is no change.
			if instant && dropKindOf(b.w.sql, c) == dropIndex && !c.IfExists && !addsIndex {
				if k, ok := b.key(spaceIndex, c.Name); ok {
					if f, _ := b.l.Lookup(idxKey(r.db, r.name, k)); f == walk.Present {
						return fail("ALGORITHM=INSTANT", "drop an index")
					}
				}
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			con := c.Constraint
			if con == nil {
				continue
			}
			switch con.Type {
			case ast.ConstrCheck:
				if inPlace && !con.NotEnforced && !con.IfNotExists {
					return fail("ALGORITHM="+algorithm, "add a "+review.Code("CHECK")+" constraint")
				}
			case ast.ConstrForeignKey:
				if inPlace && !con.IfNotExists && b.w.fkChecks == review.SwitchOn {
					return fail("ALGORITHM="+algorithm, "add a foreign key while "+review.Code("foreign_key_checks")+" is on")
				}
			case ast.ConstrFulltextIndex:
				if instant && !con.IfNotExists {
					return fail("ALGORITHM=INSTANT", "add a "+review.Code("FULLTEXT")+" index")
				}
			case ast.ConstrIndex, ast.ConstrUnique:
				// Dropping an index and adding it back as it was is no
				// change at all (verified), so a statement that drops one
				// is left alone.
				if instant && !con.IfNotExists && !dropsIndex {
					return fail("ALGORITHM=INSTANT", "add an index")
				}
			}
		case ast.ATAddColumn:
			if c.IfNotExists {
				continue
			}
			for _, col := range columnsOf(c) {
				if col == nil {
					continue
				}
				if inPlace && col.Generated != nil && col.Generated.Stored && col.Generated.Expr != nil {
					return fail("ALGORITHM="+algorithm, "add a "+review.Code("STORED")+" generated column")
				}
				if lock == "NONE" && (col.AutoIncrement || colDeclares(col, ast.ColConstrAutoIncrement)) {
					return fail("LOCK=NONE", "add an "+review.Code("AUTO_INCREMENT")+" column")
				}
			}
		case ast.ATModifyColumn, ast.ATChangeColumn:
			if !inPlace || c.Column == nil || c.IfExists {
				continue
			}
			from := c.Column.Name
			if c.Type == ast.ATChangeColumn {
				from = c.Name
			}
			k, ok := b.key(spaceColumn, from)
			if !ok || !b.namedOnce(v, k) {
				continue
			}
			if c.Type == ast.ATChangeColumn {
				if to, ok := b.key(spaceColumn, c.Column.Name); !ok || to != k && !b.namedOnce(v, to) {
					continue
				}
			}
			if _, gen := st.generated[k]; gen || c.Column.Generated != nil {
				continue
			}
			if typeChanges(b.ledgerType(r, k), astBaseType(c.Column.TypeName)) {
				return fail("ALGORITHM="+algorithm, "change a column's type")
			}
		}
	}
	return false
}

// indexAlgorithmFails fails a CREATE INDEX on a synced InnoDB table the
// change has not touched with ALGORITHM=INSTANT (1846, or 1845 for a
// FULLTEXT index; verified on mariadb:11.8.8).
func (b *builder) indexAlgorithmFails(v *ast.CreateIndexStmt, r tableRef) bool {
	if !strings.EqualFold(v.Algorithm, "INSTANT") || v.IfNotExists || v.OrReplace {
		return false
	}
	st := b.syncedOf(r)
	if st == nil || st.engine != "innodb" {
		return false
	}
	op := "add an index"
	if v.Fulltext {
		op = "add a " + review.Code("FULLTEXT") + " index"
	}
	b.fail(b.anchor(), []walk.Key{relKey(r.db, r.name)}, review.UnsupportedForm, unsupported("ALGORITHM=INSTANT", op))
	return true
}
