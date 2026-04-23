package catalog

// Catalog is the in-memory TiDB catalog.
type Catalog struct {
	databases         map[string]*Database         // lowered name -> Database
	placementPolicies map[string]*PlacementPolicy  // lowered name -> PlacementPolicy (TiDB)
	currentDB         string
	defaultCharset    string
	defaultCollation  string
	foreignKeyChecks  bool // SET foreign_key_checks (default true)
}

func New() *Catalog {
	return &Catalog{
		databases:         make(map[string]*Database),
		placementPolicies: make(map[string]*PlacementPolicy),
		defaultCharset:    "utf8mb4",
		defaultCollation:  "utf8mb4_bin", // TiDB v8.5 default (MySQL 8.0 uses utf8mb4_0900_ai_ci)
		foreignKeyChecks:  true,
	}
}

// ForeignKeyChecks returns whether FK validation is enabled.
func (c *Catalog) ForeignKeyChecks() bool { return c.foreignKeyChecks }

// SetForeignKeyChecks enables or disables FK validation.
func (c *Catalog) SetForeignKeyChecks(v bool) { c.foreignKeyChecks = v }

func (c *Catalog) SetCurrentDatabase(name string) { c.currentDB = name }
func (c *Catalog) CurrentDatabase() string         { return c.currentDB }

func (c *Catalog) GetDatabase(name string) *Database {
	return c.databases[toLower(name)]
}

func (c *Catalog) Databases() []*Database {
	result := make([]*Database, 0, len(c.databases))
	for _, db := range c.databases {
		result = append(result, db)
	}
	return result
}

func (c *Catalog) resolveDatabase(name string) (*Database, error) {
	if name == "" {
		name = c.currentDB
	}
	if name == "" {
		return nil, errNoDatabaseSelected()
	}
	db := c.GetDatabase(name)
	if db == nil {
		return nil, errUnknownDatabase(name)
	}
	return db, nil
}
