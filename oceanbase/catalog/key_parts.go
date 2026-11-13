package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// keyKind is what a key's columns must be: anything, text (FULLTEXT), a geometry
// (SPATIAL) or a vector (VECTOR).
type keyKind int

const (
	keyPlain keyKind = iota
	keyFulltext
	keySpatial
	keyVector
)

// checkKeyParts refuses a key whose parts or block size OceanBase refuses, before
// the key's name is looked up (1061), in CREATE TABLE, ALTER TABLE and CREATE INDEX
// alike (TestOracle_OceanBaseKeyParts):
//
//   - BLOCK_SIZE outside 1024 to 1048576 (5258), before the parts;
//   - a column named twice (1060), whether or not the table has it;
//   - part by part: a functional part that is a bare column (3762), a prefix on a
//     column that is not a string (1089), and a column of the wrong type for a
//     FULLTEXT (1283), SPATIAL (1687) or VECTOR (7601) key.
//
// A column the table lacks ends the parts read here; validateIndexColumns answers it.
func checkKeyParts(tbl *Table, parts []*IndexColumn, kind keyKind, opts []*nodes.IndexOption) error {
	for _, opt := range opts {
		if strings.EqualFold(opt.Name, "BLOCK_SIZE") {
			if lit, ok := opt.Value.(*nodes.IntLit); ok && (lit.Value < 1024 || lit.Value > 1048576) {
				return &Error{Code: 5258, SQLState: "HY000", Message: "Invalid block size, block size should between 1024 and 1048576"}
			}
		}
	}
	seen := map[string]bool{}
	for _, p := range parts {
		if p.Name == "" {
			continue
		}
		if seen[toLower(p.Name)] {
			return errDupColumn(p.Name)
		}
		seen[toLower(p.Name)] = true
	}
	for _, p := range parts {
		if p.Name == "" {
			if name, ok := bareColumnPart(p.ExprNode); ok {
				if tbl.GetColumn(name) == nil {
					return nil // ((missing)) is 1072, validateIndexColumns's
				}
				return &Error{Code: 3762, SQLState: "HY000", Message: "Functional index on a column is not supported. Consider using a regular index instead."}
			}
			continue
		}
		col := tbl.GetColumn(p.Name)
		if col == nil {
			// A missing column is 1072 before any later part is read: FULLTEXT (zz, c)
			// is 1072 with c an INT. validateIndexColumns answers it.
			return nil
		}
		typ := toLower(col.DataType)
		if p.Length > 0 && !isStringKeyType(typ) {
			return &Error{Code: 1089, SQLState: "HY000", Message: "Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys"}
		}
		switch kind {
		case keyFulltext:
			if err := fulltextColumnError(col); err != nil {
				return err
			}
		case keySpatial:
			if !isSpatialKeyType(typ) {
				return &Error{Code: 1687, SQLState: "42000", Message: "A SPATIAL index may only contain a geometrical type column."}
			}
		case keyVector:
			if typ != "vector" {
				return &Error{Code: 7601, SQLState: "HY000", Message: fmt.Sprintf("Column '%s' cannot be part of VECTOR index", col.Name)}
			}
		}
	}
	return nil
}

// bareColumnPart returns the column of a functional key part that is only a column,
// as ((c)).
func bareColumnPart(expr nodes.ExprNode) (string, bool) {
	for {
		switch e := expr.(type) {
		case *nodes.ParenExpr:
			expr = e.Expr
		case *nodes.ColumnRef:
			return e.Column, true
		default:
			return "", false
		}
	}
}

func isStringKeyType(typ string) bool {
	switch typ {
	case "char", "varchar", "binary", "varbinary", "tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

func isSpatialKeyType(typ string) bool {
	switch typ {
	case "geometry", "point", "linestring", "polygon", "multipoint", "multilinestring",
		"multipolygon", "geometrycollection", "geomcollection":
		return true
	}
	return false
}

// constraintKeyKind is the kind of key a table constraint defines, and whether it
// defines one.
func constraintKeyKind(con *nodes.Constraint) (keyKind, bool) {
	switch con.Type {
	case nodes.ConstrPrimaryKey, nodes.ConstrUnique, nodes.ConstrIndex:
		return keyPlain, true
	case nodes.ConstrFulltextIndex:
		return keyFulltext, true
	case nodes.ConstrSpatialIndex:
		return keySpatial, true
	case nodes.ConstrVector:
		return keyVector, true
	}
	return keyPlain, false
}
