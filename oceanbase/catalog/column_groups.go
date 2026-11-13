package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// ColumnGroups are the column groups OceanBase keeps for a table beside its rows:
// ALL COLUMNS, the whole row in columnar form, and EACH COLUMN, a group per column.
// A table written with neither has its default group alone, which SHOW CREATE TABLE
// does not print. OceanBase 4.3.5.6 refuses a group that names its columns,
// g(a, b), with 1235 wherever one is written, so these two are all a table holds
// (TestOracle_OceanBaseColumnGroups).
type ColumnGroups struct {
	AllColumns bool
	EachColumn bool
}

// IsEmpty reports whether the table has its default group alone.
func (g ColumnGroups) IsEmpty() bool { return !g.AllColumns && !g.EachColumn }

// String is the list as SHOW CREATE TABLE prints it, all columns first whatever
// order the groups were written or added in, or "" for none.
func (g ColumnGroups) String() string {
	var parts []string
	if g.AllColumns {
		parts = append(parts, "all columns")
	}
	if g.EachColumn {
		parts = append(parts, "each column")
	}
	return strings.Join(parts, ", ")
}

// readColumnGroups reads a column group list the parser kept as written, the
// text inside WITH COLUMN GROUP ( ... ) or {ADD | DROP} COLUMN GROUP ( ... ),
// with the parser's own lexer. The parser has taken the list, so an element is
// ALL COLUMNS, EACH COLUMN, or a name and its columns in parentheses.
//
// OceanBase reads the elements in written order: a named group is 1235 and a
// group written twice is 1060, whichever comes first. A second EACH COLUMN is
// named 'single column group', and one already on the table 'each column
// group'.
func readColumnGroups(text string) (ColumnGroups, error) {
	var groups ColumnGroups
	var element []obparser.Token
	depth := 0
	read := func() error {
		named := false
		for _, tok := range element {
			named = named || tok.Type == '('
		}
		switch {
		case len(element) == 0: // nothing the parser takes
		case named:
			return errColumnGroupNotSupported()
		case strings.EqualFold(element[0].Str, "ALL"):
			if groups.AllColumns {
				return errDupColumnGroup("all column group")
			}
			groups.AllColumns = true
		default: // EACH COLUMN
			if groups.EachColumn {
				return errDupColumnGroup("single column group")
			}
			groups.EachColumn = true
		}
		element = element[:0]
		return nil
	}
	for _, tok := range obparser.Tokenize(text) {
		switch {
		case tok.Type == '(':
			depth++
		case tok.Type == ')':
			depth--
		case tok.Type == ',' && depth == 0:
			if err := read(); err != nil {
				return ColumnGroups{}, err
			}
			continue
		}
		element = append(element, tok)
	}
	if err := read(); err != nil {
		return ColumnGroups{}, err
	}
	return groups, nil
}

// columnGroupDelayed reports whether an {ADD | DROP} COLUMN GROUP command ends
// with DELAYED, which the parser keeps as the command's option.
func columnGroupDelayed(cmd *nodes.AlterTableCmd) bool {
	return cmd.Option != nil && strings.EqualFold(cmd.Option.Name, "DELAYED")
}

// alterAddColumnGroup is ALTER TABLE ... ADD COLUMN GROUP (...) [DELAYED]. After
// the list, OceanBase refuses DELAYED with 1235 unless the list holds EACH COLUMN,
// and then a group the table has with 1060, ALL COLUMNS before EACH COLUMN in
// either order (TestOracle_OceanBaseColumnGroups).
func alterAddColumnGroup(tbl *Table, cmd *nodes.AlterTableCmd) error {
	groups, err := readColumnGroups(cmd.Name)
	if err != nil {
		return err
	}
	if columnGroupDelayed(cmd) && !groups.EachColumn {
		return errColumnGroupNotSupported()
	}
	if groups.AllColumns && tbl.ColumnGroups.AllColumns {
		return errDupColumnGroup("all column group")
	}
	if groups.EachColumn && tbl.ColumnGroups.EachColumn {
		return errDupColumnGroup("each column group")
	}
	tbl.ColumnGroups.AllColumns = tbl.ColumnGroups.AllColumns || groups.AllColumns
	tbl.ColumnGroups.EachColumn = tbl.ColumnGroups.EachColumn || groups.EachColumn
	return nil
}

// alterDropColumnGroup is ALTER TABLE ... DROP COLUMN GROUP (...) [DELAYED].
// OceanBase refuses DELAYED with 1235 before it reads the list, and then a group
// the table lacks with 4185, ALL COLUMNS before EACH COLUMN in either order
// (TestOracle_OceanBaseColumnGroups).
func alterDropColumnGroup(tbl *Table, cmd *nodes.AlterTableCmd) error {
	if columnGroupDelayed(cmd) {
		return errColumnGroupNotSupported()
	}
	groups, err := readColumnGroups(cmd.Name)
	if err != nil {
		return err
	}
	if groups.AllColumns && !tbl.ColumnGroups.AllColumns {
		return errNoSuchColumnGroup("all column group")
	}
	if groups.EachColumn && !tbl.ColumnGroups.EachColumn {
		return errNoSuchColumnGroup("each column group")
	}
	tbl.ColumnGroups.AllColumns = tbl.ColumnGroups.AllColumns && !groups.AllColumns
	tbl.ColumnGroups.EachColumn = tbl.ColumnGroups.EachColumn && !groups.EachColumn
	return nil
}

// errColumnGroupNotSupported is 1235, OceanBase's answer to a named column group
// and to a DELAYED it does not take.
func errColumnGroupNotSupported() error {
	return &Error{Code: 1235, SQLState: "0A000", Message: "Not supported feature or function"}
}

// errDupColumnGroup is 1060 for a column group written twice or already on the
// table. OceanBase names the group by its kind, as 'all column group'.
func errDupColumnGroup(name string) error {
	return &Error{Code: ErrDupColumn, SQLState: "42S21", Message: "Duplicate column group name '" + name + "'"}
}

// errNoSuchColumnGroup is 4185: DROP COLUMN GROUP names a group the table lacks.
func errNoSuchColumnGroup(name string) error {
	return &Error{Code: 4185, SQLState: "HY000", Message: "Column group '" + name + "' not found"}
}
