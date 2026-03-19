package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

func (c *Catalog) createView(stmt *nodes.CreateViewStmt) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	if _, exists := db.Views[key]; exists {
		if !stmt.OrReplace {
			return errDupTable(stmt.Name.Name)
		}
	}
	db.Views[key] = &View{
		Name:        stmt.Name.Name,
		Database:    db,
		Algorithm:   stmt.Algorithm,
		Definer:     stmt.Definer,
		SqlSecurity: stmt.SqlSecurity,
		CheckOption: stmt.CheckOption,
		Columns:     stmt.Columns,
	}
	return nil
}

func (c *Catalog) dropView(stmt *nodes.DropViewStmt) error {
	for _, ref := range stmt.Views {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		if _, exists := db.Views[key]; !exists {
			if stmt.IfExists {
				continue
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		delete(db.Views, key)
	}
	return nil
}
