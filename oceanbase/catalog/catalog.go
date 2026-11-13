package catalog

import "fmt"

// Catalog is the in-memory MySQL catalog.
type Catalog struct {
	databases        map[string]*Database // lowered name -> Database
	currentDB        string
	defaultCharset   string
	defaultCollation string
	charsetClient    string
	collationConn    string
	foreignKeyChecks bool // SET foreign_key_checks (default true)
	session          SessionState
	madeUpNames      int // names made for constraints written without one (madeUpConstraintName)
}

// madeUpConstraintName names a foreign key or a CHECK written without a name. OceanBase
// names it after its table and the time in microseconds, child_OBFK_1794529609863994 or
// t_OBCHECK_1794529978914345, so no generated name collides with a written one and none
// follows the table through a rename. The catalog has no clock: a count stands in for the
// time, in the time's 16 digits, and no two are alike in one catalog. The diff pairs such
// a constraint with the engine's by definition, not by this name (forgetMadeUpRenames).
// TestScenario_C1 and TestScenario_PS hold the shape to the engine.
//
// taken reports a name a constraint already has. The count passes over those, so a schema
// loaded with the names an earlier catalog made, as its SHOW CREATE TABLE writes them,
// still takes an unnamed constraint.
func (c *Catalog) madeUpConstraintName(table, infix string, taken func(name string) bool) string {
	for {
		c.madeUpNames++
		if name := fmt.Sprintf("%s%s%016d", table, infix, c.madeUpNames); !taken(name) {
			return name
		}
	}
}

// The infix of the name OceanBase makes up for a foreign key and for a CHECK.
const (
	madeUpForeignKeyInfix = "_OBFK_"
	madeUpCheckInfix      = "_OBCHECK_"
)

type SessionState struct {
	ForeignKeyChecks             bool
	CharsetClient                string
	CollationConnection          string
	ExplicitDefaultsForTimestamp bool
	SQLMode                      string
	// Version is the MySQL major version whose stored form the diff/generate
	// canonicalizer reproduces. MySQL80 is the only one.
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
		Version:                      MySQL80,
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

// SetVersion fixes the MySQL major version whose stored form the diff/generate
// canonicalizer reproduces for this catalog. MySQL80 is the only one.
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
