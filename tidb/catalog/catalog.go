package catalog

// Catalog is the in-memory TiDB catalog.
type Catalog struct {
	databases         map[string]*Database        // lowered name -> Database
	placementPolicies map[string]*PlacementPolicy // lowered name -> PlacementPolicy (TiDB)
	rangePolicies     map[string]string           // lowered range name -> policy ALTER RANGE set (TiDB)
	currentDB         string
	defaultCharset    string
	defaultCollation  string
	foreignKeyChecks  bool // SET foreign_key_checks (default true)
	// checkConstraints is tidb_enable_check_constraint, OFF by default. While it
	// is off, v8.5.5 parses CHECK and keeps nothing: a CHECK clause in CREATE or
	// ALTER TABLE adds no constraint, and DROP CHECK, DROP CONSTRAINT and ALTER
	// CHECK succeed whatever they name (read on 2026-11-06).
	checkConstraints bool
	// nonStrict records a session sql_mode with neither STRICT_TRANS_TABLES nor
	// STRICT_ALL_TABLES nor TRADITIONAL. The server's default is strict, which
	// the zero value stands for.
	nonStrict bool
	// droppingTables holds, while DROP TABLE runs, the "db.table" keys of the
	// tables it drops, whose foreign keys hold no parent back.
	droppingTables map[string]bool
	// pendingRetargets holds the foreign key retargeting the renames of the
	// running ALTER TABLE queued, applied once the statement holds.
	pendingRetargets []func()
	// extendedStats is tidb_enable_extended_stats, OFF by default in TiDB v8.5.5.
	extendedStats bool
	// tiflashServerCount models the cluster's available TiFlash store count for
	// SET TIFLASH REPLICA validation. -1 (the default) leaves replica counts
	// unconstrained for offline schema modeling; >= 0 enforces TiDB's check that
	// replica count <= server count (1105).
	tiflashServerCount int
}

func New() *Catalog {
	return &Catalog{
		databases:          make(map[string]*Database),
		placementPolicies:  make(map[string]*PlacementPolicy),
		rangePolicies:      make(map[string]string),
		defaultCharset:     "utf8mb4",
		defaultCollation:   "utf8mb4_bin", // TiDB v8.5 default (MySQL 8.0 uses utf8mb4_0900_ai_ci)
		foreignKeyChecks:   true,
		tiflashServerCount: -1,
	}
}

// CheckConstraints reports whether tidb_enable_check_constraint is on.
func (c *Catalog) CheckConstraints() bool { return c.checkConstraints }

// SetCheckConstraints sets tidb_enable_check_constraint.
func (c *Catalog) SetCheckConstraints(v bool) { c.checkConstraints = v }

// ExtendedStats reports whether tidb_enable_extended_stats is on.
func (c *Catalog) ExtendedStats() bool { return c.extendedStats }

// SetExtendedStats sets tidb_enable_extended_stats.
func (c *Catalog) SetExtendedStats(v bool) { c.extendedStats = v }

// TiFlashServerCount returns the configured TiFlash server count (-1 means unconstrained).
func (c *Catalog) TiFlashServerCount() int { return c.tiflashServerCount }

// SetTiFlashServerCount sets the TiFlash server count used to validate SET TIFLASH REPLICA.
func (c *Catalog) SetTiFlashServerCount(n int) { c.tiflashServerCount = n }

// ForeignKeyChecks returns whether FK validation is enabled.
func (c *Catalog) ForeignKeyChecks() bool { return c.foreignKeyChecks }

// SetForeignKeyChecks enables or disables FK validation.
func (c *Catalog) SetForeignKeyChecks(v bool) { c.foreignKeyChecks = v }

func (c *Catalog) SetCurrentDatabase(name string) { c.currentDB = name }
func (c *Catalog) CurrentDatabase() string        { return c.currentDB }

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
