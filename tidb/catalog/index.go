package catalog

type Index struct {
	Name     string
	Table    *Table
	Columns  []*IndexColumn
	Unique   bool
	Primary  bool
	Fulltext bool
	Spatial  bool
	// Vector marks TiDB's VECTOR INDEX, which SHOW CREATE TABLE writes with
	// that word in place of KEY (ConstraintVector in v8.5.5
	// pkg/parser/ast/ddl.go restores `VECTOR INDEX`).
	Vector       bool
	IndexType    string // BTREE, HASH, FULLTEXT, SPATIAL
	Comment      string
	Visible      bool
	KeyBlockSize int
	// addedByAlter marks a primary key ALTER TABLE added, which TiDB builds
	// NONCLUSTERED unless told CLUSTERED.
	addedByAlter bool
}

type IndexColumn struct {
	Name       string
	Expr       string
	Length     int
	Descending bool
}
