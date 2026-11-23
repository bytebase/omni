package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// A column's key carries its base data type as its kind: the type the
// sync lists (ColumnMetadata.Type, "int(11) unsigned" read as int), or
// the one the change declares. "" is a type the walk does not know. The
// walk judges a key part by its column's type only where MariaDB refuses
// every type of a class, and only on a table it knows to be InnoDB:
// sql_table.cc checks the type before the engine, but an engine may
// declare FULLTEXT or prefix keys of its own.

// baseType returns the lower-case base type of a type as the sync lists
// it or the parser names it.
func baseType(raw string) string {
	s := strings.TrimSpace(raw)
	if i := strings.IndexAny(s, " ("); i >= 0 {
		s = s[:i]
	}
	s = strings.ToLower(s)
	switch s {
	case "integer", "int4":
		return "int"
	case "int1", "bool", "boolean":
		return "tinyint"
	case "int2":
		return "smallint"
	case "int3", "middleint":
		return "mediumint"
	case "int8":
		return "bigint"
	case "dec", "numeric", "fixed":
		return "decimal"
	}
	return s
}

func astBaseType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	return baseType(dt.Name)
}

// numericOrTemporal reports whether a base type is a number, a bit
// string, or a date or time: none is a string, a geometry, or a vector.
func numericOrTemporal(t string) bool {
	switch t {
	case "tinyint", "smallint", "mediumint", "int", "bigint",
		"decimal", "float", "double", "bit",
		"date", "time", "datetime", "timestamp", "year":
		return true
	}
	return false
}

// characterType reports whether a base type is a character string.
func characterType(t string) bool {
	switch t {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		return true
	}
	return false
}

// keyKind is what a key's type makes MariaDB demand of its parts.
type keyKind uint8

const (
	keyPlain keyKind = iota
	keyFulltext
	keySpatial
	keyVector
)

func constraintKeyKind(c *ast.Constraint) keyKind {
	switch {
	case c.Type == ast.ConstrFulltextIndex:
		return keyFulltext
	case c.Type == ast.ConstrSpatialIndex, strings.EqualFold(c.IndexType, "RTREE"):
		return keySpatial
	case c.Type == ast.ConstrVector:
		return keyVector
	}
	return keyPlain
}

func createIndexKeyKind(v *ast.CreateIndexStmt) keyKind {
	switch {
	case v.Fulltext:
		return keyFulltext
	case v.Spatial, strings.EqualFold(v.IndexType, "RTREE"):
		return keySpatial
	case v.Vector:
		return keyVector
	}
	return keyPlain
}

// keyPartTypesFail fails a key whose plain column part has a type the
// key refuses (verified on mariadb:11.8.8): a FULLTEXT key a column that
// is no character string (1283), a SPATIAL or RTREE key one that is no
// geometry and a VECTOR key one that is no vector (1210, "Incorrect
// arguments to SPATIAL INDEX"), and any other key a prefix of a column
// that is no string (1089). typeOf returns a column's base type, ""
// when the walk does not know it.
func (b *builder) keyPartTypesFail(kind keyKind, cols []*ast.IndexColumn, typeOf func(k string) string) bool {
	for _, ic := range cols {
		if ic == nil || ic.Functional || ic.WithoutOverlaps {
			continue
		}
		ref, ok := ic.Expr.(*ast.ColumnRef)
		if !ok || ref.Star || ref.Table != "" || ref.Schema != "" {
			continue
		}
		k, ok := b.key(spaceColumn, ref.Column)
		if !ok {
			continue
		}
		t := typeOf(k)
		if t == "" {
			continue
		}
		rng := rangeOf(ic.Loc)
		switch kind {
		case keyFulltext:
			if numericOrTemporal(t) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(ref.Column))+" cannot be part of a "+review.Code("FULLTEXT")+" index.") // 1283
				return true
			}
		case keySpatial:
			if numericOrTemporal(t) || characterType(t) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(ref.Column))+" cannot be part of a "+review.Code("SPATIAL")+" index.") // 1210
				return true
			}
		case keyVector:
			if numericOrTemporal(t) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(ref.Column))+" cannot be part of a "+review.Code("VECTOR")+" index.") // 1210
				return true
			}
		default:
			if ic.HasPrefix && ic.Length > 0 && numericOrTemporal(t) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(ref.Column))+" cannot use an index prefix length.") // 1089
				return true
			}
		}
	}
	return false
}

// ledgerType returns the type the ledger gives a column of a table.
func (b *builder) ledgerType(r tableRef, k string) string {
	if f, kind := b.l.Lookup(colKey(r.db, r.name, k)); f == walk.Present {
		return kind
	}
	return ""
}

// alterTypes returns the column types an ALTER TABLE sets for its own
// later subcommands, which the ledger, as it was before the statement,
// does not show: a column it adds, modifies, changes, renames, or drops
// has the type it declares, or none the walk knows.
func (b *builder) alterTypes(cmds []*ast.AlterTableCmd) map[string]string {
	types := make(map[string]string)
	set := func(name, t string) {
		if k, ok := b.key(spaceColumn, name); ok {
			types[k] = t
		}
	}
	for _, c := range cmds {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATAddColumn:
			cols := c.Columns
			if c.Column != nil {
				cols = append([]*ast.ColumnDef{c.Column}, cols...)
			}
			for _, col := range cols {
				if col == nil {
					continue
				}
				// IF NOT EXISTS may leave the column it finds as it was.
				if c.IfNotExists {
					set(col.Name, "")
				} else {
					set(col.Name, astBaseType(col.TypeName))
				}
			}
		case ast.ATModifyColumn:
			if c.Column != nil {
				set(c.Column.Name, astBaseType(c.Column.TypeName))
			}
		case ast.ATChangeColumn:
			set(c.Name, "")
			if c.Column != nil {
				set(c.Column.Name, astBaseType(c.Column.TypeName))
			}
		case ast.ATRenameColumn:
			set(c.Name, "")
			set(c.NewName, "")
		case ast.ATDropColumn:
			set(c.Name, "")
		}
	}
	return types
}
