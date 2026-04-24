package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

func (c *Catalog) createDatabase(stmt *nodes.CreateDatabaseStmt) error {
	name := stmt.Name
	key := toLower(name)
	if c.databases[key] != nil {
		if stmt.IfNotExists {
			return nil
		}
		return errDupDatabase(name)
	}
	charset := c.defaultCharset
	collation := c.defaultCollation
	charsetExplicit := false
	collationExplicit := false
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "character set", "charset":
			charset = normalizeCharsetName(opt.Value)
			charsetExplicit = true
		case "collate":
			collation = opt.Value
			collationExplicit = true
		}
	}
	// When charset is specified without explicit collation, derive the default collation.
	if charsetExplicit && !collationExplicit {
		if dc, ok := defaultCollationForCharset[toLower(charset)]; ok {
			collation = dc
		}
	}
	c.databases[key] = newDatabase(name, charset, collation)
	return nil
}

func (c *Catalog) dropDatabase(stmt *nodes.DropDatabaseStmt) error {
	name := stmt.Name
	key := toLower(name)
	if c.databases[key] == nil {
		if stmt.IfExists {
			return nil
		}
		return errUnknownDatabase(name)
	}
	delete(c.databases, key)
	if toLower(c.currentDB) == key {
		c.currentDB = ""
	}
	return nil
}

func (c *Catalog) useDatabase(stmt *nodes.UseStmt) error {
	name := stmt.Database
	key := toLower(name)
	if c.databases[key] == nil {
		return errUnknownDatabase(name)
	}
	c.currentDB = name
	return nil
}

func (c *Catalog) alterDatabase(stmt *nodes.AlterDatabaseStmt) error {
	name := stmt.Name
	if name == "" {
		name = c.currentDB
	}
	db, err := c.resolveDatabase(name)
	if err != nil {
		return err
	}
	charsetExplicit := false
	collationExplicit := false
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "character set", "charset":
			db.Charset = normalizeCharsetName(opt.Value)
			charsetExplicit = true
		case "collate":
			db.Collation = opt.Value
			collationExplicit = true
		}
	}
	// When charset is changed without explicit collation, derive the default collation.
	if charsetExplicit && !collationExplicit {
		if dc, ok := defaultCollationForCharset[toLower(db.Charset)]; ok {
			db.Collation = dc
		}
	}
	return nil
}
