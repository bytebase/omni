package catalog

type Index struct {
	Name         string
	Table        *Table
	Columns      []*IndexColumn
	Unique       bool
	Primary      bool
	Fulltext     bool
	Spatial      bool
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
