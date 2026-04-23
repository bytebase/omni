package catalog

import nodes "github.com/bytebase/omni/tidb/ast"

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
	placementPolicy := ""
	tiFlashReplica := 0
	var tiFlashLocationLabels []string
	charsetExplicit := false
	collationExplicit := false
	for _, opt := range stmt.Options {
		switch toLower(opt.Name) {
		case "character set", "charset":
			charset = opt.Value
			charsetExplicit = true
		case "collate":
			collation = opt.Value
			collationExplicit = true
		case "placement policy":
			placementPolicy = opt.Value
		case "tiflash replica":
			tiFlashReplica = opt.TiFlashReplica
			tiFlashLocationLabels = append([]string(nil), opt.TiFlashLocationLabels...)
		}
	}
	if err := c.validatePolicyRef(placementPolicy); err != nil {
		return err
	}
	// When charset is specified without explicit collation, derive the default collation.
	if charsetExplicit && !collationExplicit {
		if dc, ok := defaultCollationForCharset[toLower(charset)]; ok {
			collation = dc
		}
	}
	db := newDatabase(name, charset, collation)
	db.PlacementPolicy = resolvePolicyRef(placementPolicy)
	db.TiFlashReplica = tiFlashReplica
	db.TiFlashLocationLabels = tiFlashLocationLabels
	c.databases[key] = db
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
			db.Charset = opt.Value
			charsetExplicit = true
		case "collate":
			db.Collation = opt.Value
			collationExplicit = true
		case "placement policy":
			if err := c.validatePolicyRef(opt.Value); err != nil {
				return err
			}
			db.PlacementPolicy = resolvePolicyRef(opt.Value)
		case "tiflash replica":
			db.TiFlashReplica = opt.TiFlashReplica
			db.TiFlashLocationLabels = append([]string(nil), opt.TiFlashLocationLabels...)
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
