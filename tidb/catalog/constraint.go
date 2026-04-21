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

	// TiDB-specific constraint metadata.
	// Clustered is nil when unset, &true for CLUSTERED PK, &false for NONCLUSTERED.
	Clustered *bool
}
