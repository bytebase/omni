package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// defineSRS records a spatial reference system CREATE SPATIAL REFERENCE
// SYSTEM added, which a column's SRID may then name. A system the server
// ships stays as it is, and IF NOT EXISTS leaves one that exists alone.
func (c *Catalog) defineSRS(s *nodes.CreateSpatialRefSysStmt) {
	if builtIn, _ := lookupSRS(int(s.SRID)); builtIn {
		return
	}
	if _, defined := c.userSRS[int(s.SRID)]; defined && s.IfNotExists {
		return
	}
	if c.userSRS == nil {
		c.userSRS = make(map[int]bool)
	}
	c.userSRS[int(s.SRID)] = strings.HasPrefix(strings.ToUpper(strings.TrimSpace(s.Definition)), "GEOGCS")
}

// dropSRS forgets a spatial reference system CREATE SPATIAL REFERENCE SYSTEM
// added, unless a column's SRID names it (3716, mysql:8.0.46).
func (c *Catalog) dropSRS(s *nodes.DropSpatialRefSysStmt) error {
	srid := int(s.SRID)
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			for _, col := range tbl.Columns {
				if col.HasSRID && col.SRID == srid {
					return &Error{Code: 3716, SQLState: "SR005", Message: fmt.Sprintf("Can't modify SRID %d. There is at least one column depending on it.", srid)}
				}
			}
		}
	}
	delete(c.userSRS, srid)
	return nil
}

// lookupSRS reports whether srid names a spatial reference system the
// server has, built in or added, and whether it is geographic.
func (c *Catalog) lookupSRS(srid int) (exists, geographic bool) {
	if geographic, ok := c.userSRS[srid]; ok {
		return true, geographic
	}
	return lookupSRS(srid)
}

// checkColumnSRIDs refuses a column whose SRID names no spatial reference
// system the server has (3548, mysql:8.0.46: GEOMETRY SRID 9999).
func (c *Catalog) checkColumnSRIDs(defs []*nodes.ColumnDef) error {
	for _, def := range defs {
		if def == nil || def.TypeName == nil || !def.TypeName.HasSRID {
			continue
		}
		if exists, _ := c.lookupSRS(def.TypeName.SRID); !exists {
			return &Error{Code: 3548, SQLState: "SR001", Message: fmt.Sprintf("There's no spatial reference system with SRID %d.", def.TypeName.SRID)}
		}
	}
	return nil
}

// checkKeysOverRedefinedColumns holds the keys a table has to the columns an
// ALTER TABLE defined anew (before holds the columns it had): a SPATIAL key's
// column stays NOT NULL (1252), a FULLTEXT key's a character type (1283), and
// no engine but InnoDB takes a SPATIAL key on a geographic SRID (1178)
// (mysql:8.0.46).
func (c *Catalog) checkKeysOverRedefinedColumns(tbl *Table, before map[*Column]bool) error {
	for _, idx := range tbl.Indexes {
		if !idx.Spatial && !idx.Fulltext {
			continue
		}
		for _, ic := range idx.Columns {
			col := tbl.GetColumn(ic.Name)
			if col == nil || before[col] {
				continue
			}
			switch {
			case idx.Spatial && col.Nullable:
				return &Error{Code: 1252, SQLState: "42000", Message: "All parts of a SPATIAL index must be NOT NULL"}
			case idx.Fulltext && !isFulltextColumnType(col.DataType):
				return errCannotPartOfFulltext(col.Name)
			}
		}
	}
	return c.checkGeographicSpatialKeys(tbl)
}

// checkGeographicSpatialKeys refuses a SPATIAL key on a column of a
// geographic SRID in a table whose engine is not InnoDB (1178): MyISAM takes
// SRID 0 and 3857 and refuses 4326 (mysql:8.0.46).
func (c *Catalog) checkGeographicSpatialKeys(tbl *Table) error {
	if isInnoDBEngine(tbl.Engine) {
		return nil
	}
	for _, idx := range tbl.Indexes {
		if !idx.Spatial {
			continue
		}
		for _, ic := range idx.Columns {
			if col := tbl.GetColumn(ic.Name); col != nil && col.HasSRID {
				if _, geographic := c.lookupSRS(col.SRID); geographic {
					return &Error{Code: 1178, SQLState: "42000", Message: "The storage engine for the table doesn't support geographic spatial reference systems"}
				}
			}
		}
	}
	return nil
}

// spatialKeysForGeometry makes SPATIAL each plain key of an InnoDB table on
// one NOT NULL geometry column alone, as InnoDB stores it: ADD INDEX idx (loc)
// on a POINT NOT NULL column reads back as SPATIAL KEY `idx` (`loc`)
// (mysql:8.0.46).
func spatialKeysForGeometry(tbl *Table) {
	if !isInnoDBEngine(tbl.Engine) {
		return
	}
	for _, idx := range tbl.Indexes {
		if idx.Primary || idx.Unique || idx.Fulltext || idx.Spatial || len(idx.Columns) != 1 {
			continue
		}
		ic := idx.Columns[0]
		if ic.Expr != "" || ic.Length > 0 {
			continue
		}
		if col := tbl.GetColumn(ic.Name); col != nil && !col.Nullable && isGeometryColumnType(col.DataType) {
			idx.Spatial = true
			idx.IndexType = "SPATIAL"
		}
	}
}
