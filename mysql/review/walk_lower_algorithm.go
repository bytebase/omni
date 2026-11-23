package review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// alterClauses reads the ALGORITHM and LOCK clauses of an ALTER TABLE,
// upper case, "" when not written; the last one written counts.
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

// lockAlgorithmFault returns the fault of an ALGORITHM and LOCK pair MySQL
// refuses whatever the statement does (verified on mysql:8.0.46): INSTANT
// with a LOCK other than DEFAULT (1221), and COPY with LOCK=NONE (1846).
func lockAlgorithmFault(algorithm, lock string) string {
	switch {
	case algorithm == "INSTANT" && (lock == "NONE" || lock == "SHARED" || lock == "EXCLUSIVE"):
		return review.Code("ALGORITHM=INSTANT") + " can't take " + review.Code("LOCK="+lock) + "."
	case algorithm == "COPY" && lock == "NONE":
		return review.Code("ALGORITHM=COPY") + " needs a lock, and " + review.Code("LOCK=NONE") + " allows none."
	}
	return ""
}

// unsupported is the text for an operation InnoDB can't run with the
// ALGORITHM or LOCK the statement asks for (1845, 1846).
func unsupported(clause, operation string) string {
	return "InnoDB can't " + operation + " with " + review.Code(clause) + "."
}

// neverInstant names an ALTER TABLE subcommand InnoDB never runs with
// ALGORITHM=INSTANT (1845, each verified on mysql:8.0.46), "" for one it
// may: adding or dropping an index or a key, renaming one or changing its
// visibility, adding a CHECK constraint or a foreign key, the table options
// that rebuild the table, CONVERT TO, and FORCE.
func neverInstant(c *ast.AlterTableCmd) string {
	switch c.Type {
	case ast.ATAddIndex:
		return "add an index"
	case ast.ATAddConstraint:
		if c.Constraint != nil && c.Constraint.Type == ast.ConstrCheck {
			// A CHECK that isn't enforced judges no row, and runs
			// INSTANT (mysql:8.0.46).
			if c.Constraint.NotEnforced {
				return ""
			}
			return "add a " + review.Code("CHECK") + " constraint"
		}
		return "add a key"
	case ast.ATDropIndex:
		return "drop an index"
	case ast.ATDropConstraint:
		if strings.EqualFold(c.Name, "PRIMARY") {
			return "drop the primary key"
		}
	case ast.ATRenameIndex:
		return "rename an index"
	case ast.ATAlterIndexVisible, ast.ATAlterIndexInvisible:
		return "change an index's visibility"
	case ast.ATConvertCharset:
		return "convert the table's character set"
	case ast.ATForce:
		return "rebuild the table"
	case ast.ATTableOption:
		if c.Option == nil {
			return ""
		}
		switch n := optionName(c.Option); n {
		case "ENGINE", "COMMENT", "AUTO_INCREMENT", "ROW_FORMAT":
			return "change " + review.Code(n)
		}
	}
	return ""
}

// alterAlgorithmFails fails an ALTER TABLE of an InnoDB table whose
// ALGORITHM or LOCK clause can't run what it does (verified on
// mysql:8.0.46):
//
//   - ALGORITHM=INSTANT with an operation InnoDB never runs instantly
//     (neverInstant, 1845);
//   - a column type change, which rebuilds the table, with INSTANT, INPLACE
//     or LOCK=NONE (1846);
//   - dropping the primary key without adding one, with INSTANT, INPLACE or
//     LOCK=NONE (1846);
//   - adding a CHECK constraint, which copies the table, with INPLACE (1845)
//     or LOCK=NONE (1846), and a FULLTEXT index with LOCK=NONE (1846).
func (b *builder) alterAlgorithmFails(v *ast.AlterTableStmt, r tableRef) bool {
	algorithm, lock := alterClauses(v)
	if algorithm == "" && lock == "" {
		return false
	}
	if msg := lockAlgorithmFault(algorithm, lock); msg != "" {
		b.fail(b.anchor(), nil, review.UnsupportedForm, msg)
		return true
	}
	eng, known := b.w.engines[relKey(r.db, r.name)]
	if !known || !isInnoDBEngine(eng) {
		return false
	}
	if _, sets := alterEngine(v); sets {
		if e, _ := alterEngine(v); !isInnoDBEngine(e) {
			return false
		}
	}
	rebuilds := func(op string) bool {
		for _, clause := range []string{"ALGORITHM=INSTANT", "ALGORITHM=INPLACE", "LOCK=NONE"} {
			if clause == "ALGORITHM="+algorithm || clause == "LOCK="+lock {
				b.fail(b.anchor(), nil, review.UnsupportedForm, unsupported(clause, op))
				return true
			}
		}
		return false
	}
	addsPrimary, dropsPrimary := false, false
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch {
		case (c.Type == ast.ATDropIndex || c.Type == ast.ATDropConstraint) && strings.EqualFold(c.Name, "PRIMARY"):
			dropsPrimary = true
		case (c.Type == ast.ATAddConstraint || c.Type == ast.ATAddIndex) && c.Constraint != nil && c.Constraint.Type == ast.ConstrPrimaryKey:
			addsPrimary = true
		case c.Type == ast.ATAddColumn || c.Type == ast.ATModifyColumn || c.Type == ast.ATChangeColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if colHasPrimaryKey(col) {
					addsPrimary = true
				}
			}
		}
	}
	if dropsPrimary && !addsPrimary {
		if f, _ := b.l.Lookup(idxKey(r.db, r.name, primaryIndex)); f == walk.Present && rebuilds("drop the primary key without adding one") {
			return true
		}
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		if op := neverInstant(c); op != "" && algorithm == "INSTANT" {
			b.fail(b.anchor(), nil, review.UnsupportedForm, unsupported("ALGORITHM=INSTANT", op))
			return true
		}
		switch c.Type {
		case ast.ATModifyColumn, ast.ATChangeColumn:
			if c.Column == nil {
				continue
			}
			from := c.Column.Name
			if c.Type == ast.ATChangeColumn {
				from = c.Name
			}
			k, ok := b.key(spaceColumn, from)
			if !ok {
				continue
			}
			old, now := b.w.colTypes[colKey(r.db, r.name, k)], astBaseType(c.Column.TypeName)
			// FLOAT(p) is DOUBLE for p over 24, which the base type
			// doesn't tell, so the two are left alone.
			floating := (old == "float" || old == "double") && (now == "float" || now == "double")
			if old != "" && now != "" && old != now && !floating && rebuilds("change a column's type") {
				return true
			}
		case ast.ATAddConstraint, ast.ATAddIndex:
			if c.Constraint == nil {
				continue
			}
			if c.Constraint.Type == ast.ConstrFulltextIndex && lock == "NONE" {
				b.fail(b.anchor(), nil, review.UnsupportedForm, unsupported("LOCK=NONE", "add a "+review.Code("FULLTEXT")+" index"))
				return true
			}
			if c.Constraint.Type != ast.ConstrCheck || c.Constraint.NotEnforced {
				continue
			}
			switch {
			case algorithm == "INPLACE":
				b.fail(b.anchor(), nil, review.UnsupportedForm, unsupported("ALGORITHM=INPLACE", "add a "+review.Code("CHECK")+" constraint"))
				return true
			case lock == "NONE":
				b.fail(b.anchor(), nil, review.UnsupportedForm, unsupported("LOCK=NONE", "add a "+review.Code("CHECK")+" constraint"))
				return true
			}
		}
	}
	return false
}

// indexAlgorithmFails fails a CREATE INDEX or DROP INDEX on an InnoDB table
// with ALGORITHM=INSTANT (1845) or an ALGORITHM and LOCK pair MySQL refuses
// (1221, 1846).
func (b *builder) indexAlgorithmFails(r tableRef, algorithm, lock, op string) bool {
	algorithm, lock = strings.ToUpper(algorithm), strings.ToUpper(lock)
	if msg := lockAlgorithmFault(algorithm, lock); msg != "" {
		b.fail(b.anchor(), nil, review.UnsupportedForm, msg)
		return true
	}
	if algorithm != "INSTANT" {
		return false
	}
	if eng, known := b.w.engines[relKey(r.db, r.name)]; !known || !isInnoDBEngine(eng) {
		return false
	}
	b.fail(b.anchor(), nil, review.UnsupportedForm, unsupported("ALGORITHM=INSTANT", op))
	return true
}

// colHasPrimaryKey reports a column definition declaring PRIMARY KEY.
func colHasPrimaryKey(col *ast.ColumnDef) bool {
	for _, cc := range col.Constraints {
		if cc != nil && cc.Type == ast.ColConstrPrimaryKey {
			return true
		}
	}
	return false
}
