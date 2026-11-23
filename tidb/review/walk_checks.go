package review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// A table's CHECK constraints are one key of its shape, read only when
// Present, its kind each constraint's name and the columns it reads
// (ConstraintInfo.ConstraintCols). The sync gives it; any ALTER TABLE of
// the table leaves it Unknown, and a cleared table path Absent, which says
// nothing.
//
// TiDB v8.5.5 refuses dropping a column a CHECK reading other columns too
// reads, and renaming a column any CHECK reads (3959,
// IsColumnDroppableWithCheckConstraint and
// IsColumnRenameableWithCheckConstraint in pkg/ddl/constraint.go), whether
// tidb_enable_check_constraint is on or off once the CHECK exists
// (verified). A CHECK reading only the column goes with it. DROP CHECK and
// DROP CONSTRAINT in a statement of several subcommands fail too, but with
// 8200 only while the switch is on, and are a no-op while it is off, so the
// walk, which does not know the switch, claims nothing once one of them
// comes first.
const spaceCheckSet = "check set"

func checkSetKey(db, table string) walk.Key {
	return walk.Key{Space: spaceCheckSet, Parent: tablePath(db, table), Name: shapeName}
}

// tableCheck is a CHECK constraint as the walk keeps it: its name and the
// keys of the columns it reads.
type tableCheck struct {
	name string
	cols []string
}

// syncedChecks reads a synced table's CHECK constraints, false when the
// walk cannot read one of them.
func syncedChecks(t *metadata.TableMetadata) (string, bool) {
	var checks []tableCheck
	for _, c := range t.GetCheckConstraints() {
		name, ok := plainName(spaceCheckSet, c.GetName())
		if !ok {
			return "", false
		}
		cols, ok := checkColumns(c.GetExpression())
		if !ok {
			return "", false
		}
		checks = append(checks, tableCheck{name: name, cols: cols})
	}
	return encodeChecks(checks), true
}

// checkColumns returns the keys of the distinct columns a CHECK
// expression reads, false when the walk cannot read it.
func checkColumns(expr string) ([]string, bool) {
	const prefix = "SELECT ("
	list, err := parser.Parse(prefix + expr + ")")
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil, false
	}
	var cols []string
	ok := true
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ColumnRef:
			k, readable := plainName(spaceColumn, n.Column)
			if n.Star || !readable {
				ok = false
				return false
			}
			if !slices.Contains(cols, k) {
				cols = append(cols, k)
			}
		case *ast.SubqueryExpr:
			ok = false
			return false
		}
		return ok
	})
	return cols, ok
}

func encodeChecks(checks []tableCheck) string {
	parts := make([]string, 0, len(checks))
	for _, c := range checks {
		parts = append(parts, c.name+"\x1e"+strings.Join(c.cols, "\x1f"))
	}
	return strings.Join(parts, "\x1d")
}

func decodeChecks(kind string) []tableCheck {
	if kind == "" {
		return nil
	}
	var out []tableCheck
	for _, part := range strings.Split(kind, "\x1d") {
		name, cols, _ := strings.Cut(part, "\x1e")
		c := tableCheck{name: name}
		if cols != "" {
			c.cols = strings.Split(cols, "\x1f")
		}
		out = append(out, c)
	}
	return out
}

// checkUsesFails fails an ALTER TABLE that drops a column a CHECK reads
// with other columns, or renames a column a CHECK reads (3959). TiDB judges
// every subcommand in order against the table as it was, so the first one
// refused is the one reported; the walk stops at a DROP CHECK or DROP
// CONSTRAINT, whose answer turns on the switch.
func (b *builder) checkUsesFails(v *ast.AlterTableStmt, r tableRef) bool {
	f, kind := b.l.Lookup(checkSetKey(r.db, r.name))
	if f != walk.Present {
		return false
	}
	checks := decodeChecks(kind)
	if len(checks) == 0 {
		return false
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		var col string
		renames := false
		switch c.Type {
		case ast.ATDropConstraint:
			switch dropKindOf(b.w.sql, c) {
			case dropCheck, dropConstraint:
				return false
			}
			continue
		case ast.ATDropColumn:
			col = c.Name
		case ast.ATRenameColumn:
			col, renames = c.Name, true
		default:
			continue
		}
		k, ok := b.key(spaceColumn, col)
		if !ok || b.lookupFact(colKey(r.db, r.name, k)) != walk.Present {
			continue
		}
		for _, chk := range checks {
			if !slices.Contains(chk.cols, k) || (!renames && len(chk.cols) < 2) {
				continue
			}
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm,
				"Check constraint "+review.Code(ident(chk.name))+" uses column "+review.Code(ident(col))+", which can't be dropped or renamed.")
			return true
		}
	}
	return false
}
