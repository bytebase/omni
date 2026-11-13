package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableNoDuplicateIndex judges the indexes the top-level statements create,
// in order, as one list across the whole script; statements nested in a
// routine body are not judged, and nothing is ever removed from the list
// (DROP INDEX, DROP TABLE and renames are not read). A CREATE TABLE adds
// its table-level PRIMARY KEY, UNIQUE, INDEX, FULLTEXT and SPATIAL
// constraints, not column-level keys; an ALTER TABLE adds the same kinds
// from its ADD actions; a CREATE INDEX adds its index. An index is keyed by
// its kind (unique, fulltext or spatial, else none, so a primary key and a
// plain index on the same columns are duplicates), its index type as
// written before the key parts (BTREE when none is; a USING after them is
// not read), its table name without the database, and its plain column
// names (an expression key part is left out), joined with "-". The first
// index whose key an earlier index already had is the duplicate.
//
// Every CREATE TABLE, ALTER TABLE and CREATE INDEX that names its table,
// from the one that adds the duplicate on, is one TableDuplicateIndex,
// even when it adds nothing, as V1 did. Name is the duplicate's index name
// as written ("" when unnamed); Table is the table the reporting statement
// names, which can differ from the duplicate's; Start and End anchor the
// duplicate's constraint in a CREATE TABLE, and otherwise the whole ALTER
// TABLE or CREATE INDEX statement that added it; Statement is the index of
// the statement that added the duplicate. Ported from Bytebase's V1 rule
// TABLE_NO_DUPLICATE_INDEX.
func TableNoDuplicateIndex(stmts []ast.Node) []v1review.Finding {
	d := duplicateIndexes{seen: make(map[string]bool)}
	var findings []v1review.Finding
	for i, stmt := range stmts {
		var table string
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			table = n.Table.Name
			for _, c := range n.Constraints {
				if c != nil {
					d.addConstraint(table, c, c.Loc, i)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			table = n.Table.Name
			for _, cmd := range n.Commands {
				if cmd != nil && (cmd.Type == ast.ATAddConstraint || cmd.Type == ast.ATAddIndex) && cmd.Constraint != nil {
					d.addConstraint(table, cmd.Constraint, n.Loc, i)
				}
			}
		case *ast.CreateIndexStmt:
			if n.Table == nil {
				continue
			}
			table = n.Table.Name
			index := duplicateIndex{name: n.IndexName, indexType: n.IndexType, table: table, columns: duplicateIndexColumns(n.Columns), loc: n.Loc, stmt: i}
			switch {
			case n.Unique:
				index.kind = "unique"
			case n.Fulltext:
				index.kind = "fulltext"
			case n.Spatial:
				index.kind = "spatial"
			}
			d.add(index)
		default:
			continue
		}
		if dup := d.duplicate; dup != nil {
			start, end := span(dup.loc)
			findings = append(findings, v1review.Finding{Code: v1review.TableDuplicateIndex, Name: dup.name, Table: table, Start: start, End: end, Statement: dup.stmt})
		}
	}
	return findings
}

type duplicateIndex struct {
	name      string
	kind      string // unique, fulltext, spatial, or "" for any other
	indexType string
	table     string
	columns   []string
	loc       ast.Loc // the node the finding anchors
	stmt      int
}

// duplicateIndexes is the list of indexes seen so far, kept as the keys
// seen and the first index whose key was already seen.
type duplicateIndexes struct {
	seen      map[string]bool
	duplicate *duplicateIndex
}

func (d *duplicateIndexes) addConstraint(table string, c *ast.Constraint, loc ast.Loc, stmt int) {
	index := duplicateIndex{name: c.Name, indexType: c.IndexType, table: table, columns: c.Columns, loc: loc, stmt: stmt}
	switch c.Type {
	case ast.ConstrPrimaryKey, ast.ConstrIndex:
	case ast.ConstrUnique:
		index.kind = "unique"
	case ast.ConstrFulltextIndex:
		index.kind = "fulltext"
	case ast.ConstrSpatialIndex:
		index.kind = "spatial"
	default:
		return
	}
	d.add(index)
}

func (d *duplicateIndexes) add(index duplicateIndex) {
	var parts []string
	if index.kind != "" {
		parts = append(parts, index.kind)
	}
	indexType := index.indexType
	if indexType == "" {
		indexType = "BTREE"
	}
	parts = append(parts, indexType, index.table)
	parts = append(parts, index.columns...)
	key := strings.Join(parts, "-")
	if d.seen[key] && d.duplicate == nil {
		d.duplicate = &index
	}
	d.seen[key] = true
}

// duplicateIndexColumns returns the names of the key parts that are plain
// column references, nil when there are none.
func duplicateIndexColumns(cols []*ast.IndexColumn) []string {
	var names []string
	for _, col := range cols {
		if ref, ok := col.Expr.(*ast.ColumnRef); ok {
			names = append(names, ref.Column)
		}
	}
	return names
}
