package catalog

type ConstraintType int

const (
	ConPrimaryKey ConstraintType = iota
	ConUniqueKey
	ConForeignKey
	ConCheck
)

type Constraint struct {
	Name        string
	Type        ConstraintType
	Table       *Table
	Columns     []string
	IndexName   string
	RefDatabase string
	RefTable    string
	RefColumns  []string
	OnDelete    string
	OnUpdate    string
	MatchType   string
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
}
