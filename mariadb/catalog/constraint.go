package catalog

type ConstraintType int

const (
	ConPrimaryKey ConstraintType = iota
	ConUniqueKey
	ConForeignKey
	ConCheck
)

type Constraint struct {
	Name          string
	Type          ConstraintType
	Table         *Table
	Columns       []string
	IndexName     string
	RefDatabase   string
	RefTable      string
	RefColumns    []string
	OnDelete      string
	OnUpdate      string
	MatchType     string
	CheckExpr     string
	NotEnforced   bool
	CheckAnalyzed AnalyzedExpr // Phase 3: analyzed CHECK expression body
	// Implicit marks a constraint the catalog made rather than one the
	// statement declared: the json_valid check a JSON column carries. It goes
	// with its column instead of standing on its own, and no statement can
	// name it. The origin is recorded rather than guessed from the SQL,
	// because a statement may declare a check with the very same name and
	// expression and that one is ordinary: measured on mariadb:11.8.8,
	// `CREATE TABLE t (j LONGTEXT, CONSTRAINT j CHECK (json_valid(j)))` then
	// `ALTER TABLE t DROP CONSTRAINT j` is accepted, while the same DROP
	// against a JSON column answers 1091.
	Implicit bool
	// Column names the column a CHECK belongs to, where it belongs to one: a
	// check declared inside a column definition, and the json_valid check a
	// JSON column carries. MariaDB writes such a check on its column and a
	// check declared at table level at table level, whatever either is named,
	// so the owner is recorded rather than read off the name. Measured on
	// mariadb:11.8.8: `CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a >
	// 0))` writes the check at table level, and after `CREATE TABLE t (j
	// JSON, a INT, CONSTRAINT k CHECK (a > 0)); ALTER TABLE t CHANGE j k
	// JSON` the table carries two checks named k, the column's and the
	// table's, and writes both.
	//
	// It is empty for a check the statement declared at table level, which
	// belongs to no column.
	Column string
}

// isColumnCheck reports whether the check belongs to a column of its table
// rather than to the table. Such a check is written on its column, goes with
// it when the column is dropped or its definition replaced, and is not a name
// a statement can use: measured on mariadb:11.8.8, `ALTER TABLE t DROP
// CONSTRAINT a` against `CREATE TABLE t (a INT CHECK (a > 0))` answers 1091,
// and a statement may declare a table-level check of the same name.
func isColumnCheck(con *Constraint) bool {
	return con != nil && con.Type == ConCheck && con.Column != ""
}

// columnOwnsCheck reports whether a check of the table already belongs to the
// named column. A column carries at most one: measured on mariadb:11.8.8, a
// JSON column that declares its own check carries that one instead of
// json_valid, `CREATE TABLE t (j JSON CHECK (j IS NOT NULL))` reporting
// CHECK (`j` is not null) and nothing else.
func columnOwnsCheck(tbl *Table, name string) bool {
	key := toLower(name)
	for _, con := range tbl.Constraints {
		if isColumnCheck(con) && toLower(con.Column) == key {
			return true
		}
	}
	return false
}
