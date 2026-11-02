package catalog

// Catalog is the in-memory MySQL catalog.
type Catalog struct {
	databases        map[string]*Database // lowered name -> Database
	currentDB        string
	defaultCharset   string
	defaultCollation string
	charsetClient    string
	collationConn    string
	foreignKeyChecks bool // SET foreign_key_checks (default true)
	generateGIPK     bool // SET sql_generate_invisible_primary_key
	showGIPK         bool // SET show_gipk_in_create_table_and_information_schema
	session          SessionState
}

type SessionState struct {
	ForeignKeyChecks             bool
	GenerateInvisiblePrimaryKey  bool
	ShowGIPK                     bool
	CharsetClient                string
	CollationConnection          string
	ExplicitDefaultsForTimestamp bool
	SQLMode                      string
	// Version is the MySQL major version (MySQL57/MySQL80) whose stored form the
	// diff/generate canonicalizer must reproduce. It defaults to MySQL80 so an
	// unset catalog keeps the historical 8.0 behavior; bytebase sets it from the
	// synced database's server version before diffing a 5.7 schema so a bare
	// CHARSET=utf8mb4 does not phantom-diff against the 8.0 default collation, and
	// so generated DDL never emits a collation (utf8mb4_0900_ai_ci) that 5.7 lacks.
	Version Version
}

func New() *Catalog {
	session := defaultSessionState()
	return &Catalog{
		databases:        make(map[string]*Database),
		defaultCharset:   "utf8mb4",
		defaultCollation: "utf8mb4_0900_ai_ci",
		charsetClient:    session.CharsetClient,
		collationConn:    session.CollationConnection,
		foreignKeyChecks: session.ForeignKeyChecks,
		generateGIPK:     session.GenerateInvisiblePrimaryKey,
		showGIPK:         session.ShowGIPK,
		session:          session,
	}
}

func defaultSessionState() SessionState {
	return SessionState{
		ForeignKeyChecks:             true,
		CharsetClient:                "utf8mb4",
		CollationConnection:          "utf8mb4_0900_ai_ci",
		ExplicitDefaultsForTimestamp: true,
		SQLMode:                      "DEFAULT",
		// Default to the modern 8.0 stored form (Version's zero value is MySQL57, so
		// this MUST be set explicitly to preserve the historical default behavior).
		Version: MySQL80,
	}
}

// ForeignKeyChecks returns whether FK validation is enabled.
func (c *Catalog) ForeignKeyChecks() bool { return c.session.ForeignKeyChecks }

// SetForeignKeyChecks enables or disables FK validation.
func (c *Catalog) SetForeignKeyChecks(v bool) {
	c.foreignKeyChecks = v
	c.session.ForeignKeyChecks = v
}

// Version returns the MySQL major version whose stored form this catalog's diff and
// generate paths canonicalize against (default MySQL80).
func (c *Catalog) Version() Version { return c.session.Version }

// SetVersion fixes the MySQL major version (MySQL57/MySQL80) whose stored form the
// diff/generate canonicalizer reproduces for this catalog. bytebase's
// mysqlDiffSDLMigration sets it from the synced database's server version so a
// 5.7-synced schema canonicalizes (and generates DDL) as 5.7 rather than 8.0.
func (c *Catalog) SetVersion(v Version) { c.session.Version = v }

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
