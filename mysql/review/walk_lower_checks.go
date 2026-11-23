package review

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// checkUsesFails fails an ALTER TABLE that renames a column a CHECK
// constraint of the table uses, or drops one that a CHECK constraint uses
// with another column (3959, verified on mysql:8.0.46; a constraint on the
// dropped column alone goes with it). The walk knows the columns of the
// constraints the sync lists, while the table and the constraint are
// there; a statement that drops the constraint too is left alone.
func (b *builder) checkUsesFails(r tableRef, v *ast.AlterTableStmt) bool {
	checks := b.w.checkCols[relKey(r.db, r.name)]
	if len(checks) == 0 {
		return false
	}
	droppedChecks := make(map[string]bool)
	for _, c := range v.Commands {
		if c != nil && c.Type == ast.ATDropConstraint {
			if k, ok := b.key(spaceCheckConstraint, c.Name); ok {
				droppedChecks[k] = true
			}
		}
	}
	names := make([]string, 0, len(checks))
	for name := range checks {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		col, renames := "", false
		switch c.Type {
		case ast.ATDropColumn:
			col = c.Name
		case ast.ATRenameColumn:
			col, renames = c.Name, true
			if k1, ok1 := b.key(spaceColumn, c.Name); ok1 {
				if k2, ok2 := b.key(spaceColumn, c.NewName); ok2 && k1 == k2 {
					continue
				}
			}
		case ast.ATChangeColumn:
			if c.Column == nil {
				continue
			}
			col, renames = c.Name, true
			if k1, ok1 := b.key(spaceColumn, c.Name); ok1 {
				if k2, ok2 := b.key(spaceColumn, c.Column.Name); ok2 && k1 == k2 {
					continue
				}
			}
		default:
			continue
		}
		k, ok := b.key(spaceColumn, col)
		if !ok {
			continue
		}
		for _, name := range names {
			cols := checks[name]
			if droppedChecks[name] || !containsString(cols, k) || (!renames && len(cols) < 2) {
				continue
			}
			ck := chkKey(r.db, r.name, name)
			if f, _ := b.l.Lookup(ck); f != walk.Present {
				continue
			}
			if f, _ := b.l.Lookup(colKey(r.db, r.name, k)); f != walk.Present {
				continue
			}
			b.failReferenced(rangeOf(c.Loc), []walk.Key{ck}, "Check constraint "+review.Code(name)+" uses column "+review.Code(ident(col))+", which can't be dropped or renamed.")
			return true
		}
	}
	return false
}

func containsString(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// checkColumns returns the keys of the distinct columns a synced CHECK
// expression reads, false when the walk can't read it.
func checkColumns(expr string, foldTables bool) ([]string, bool) {
	cols := extractPartitionExprCols(expr)
	if len(cols) == 0 {
		return nil, false
	}
	var out []string
	for _, c := range cols {
		k, ok := nameKey(spaceColumn, strings.Trim(c, "`"), foldTables)
		if !ok {
			return nil, false
		}
		if !containsString(out, k) {
			out = append(out, k)
		}
	}
	return out, true
}
