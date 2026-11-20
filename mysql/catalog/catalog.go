package catalog

import (
	"path"
	"strings"
	"time"
)

type userVarKind int

const (
	userVarNull userVarKind = iota
	userVarString
	userVarNumeric
)

type userVarValue struct {
	kind userVarKind
	str  string
}

// Catalog is the in-memory MySQL catalog.
type Catalog struct {
	databases                 map[string]*Database      // lowered name -> Database
	tablespaces               map[string]bool           // exact general tablespace names
	tablespaceBlockSizes      map[string]uint64         // general tablespace name -> resolved FILE_BLOCK_SIZE bytes
	tablespaceEncrypted       map[string]bool           // general tablespace name -> ENCRYPTION='Y'
	undoTablespaces           map[string]bool           // exact undo tablespace names -> active
	tablespaceFiles           map[string]string         // tablespace name -> normalized datafile path
	tablespaceDirs            map[string]bool           // known InnoDB datafile directories
	dataDir                   string                    // InnoDB/server datadir (default /var/lib/mysql)
	userVars                  map[string]userVarValue   // lowered user var name (without @) -> value
	systemTimeZone            string                    // @@system_time_zone (default UTC)
	systemTimeZoneLoc         *time.Location            // optional host Location when @@system_time_zone is an abbreviation
	namedTimeZones            map[string]*time.Location // lowercased mysql.time_zone_name entries -> Location
	globalTimeZone            string
	globalSQLMode             string
	globalCharsetClient       string
	globalCollationConn       string
	globalForeignKeyChecks    bool
	globalGenerateGIPK        bool
	globalShowGIPK            bool
	globalExplicitDefaultsTS  bool
	globalDefaultTableEncrypt bool
	currentDB                 string
	defaultCharset            string
	defaultCollation          string
	charsetClient             string
	collationConn             string
	innodbPageSize            uint64 // configured innodb_page_size (0 => 16384)
	foreignKeyChecks          bool   // SET foreign_key_checks (default true)
	loadingSDL                bool   // true while LoadSDLWithVersion executes statements
	generateGIPK              bool   // SET sql_generate_invisible_primary_key
	showGIPK                  bool   // SET show_gipk_in_create_table_and_information_schema
	session                   SessionState
}

type SessionState struct {
	ForeignKeyChecks             bool
	GenerateInvisiblePrimaryKey  bool
	ShowGIPK                     bool
	CharsetClient                string
	CollationConnection          string
	ExplicitDefaultsForTimestamp bool
	DefaultTableEncryption       bool
	SQLMode                      string
	TimeZone                     string
	SystemTimeZone               string
	// Version is the MySQL major version whose stored form the diff/generate
	// canonicalizer reproduces. MySQL80 is the only one.
	Version Version
}

func New() *Catalog {
	session := defaultSessionState()
	return &Catalog{
		databases:                make(map[string]*Database),
		tablespaces:              make(map[string]bool),
		tablespaceBlockSizes:     make(map[string]uint64),
		tablespaceEncrypted:      make(map[string]bool),
		undoTablespaces:          map[string]bool{"innodb_undo_001": true, "innodb_undo_002": true},
		tablespaceFiles:          map[string]string{"innodb_undo_001": "/var/lib/mysql/undo_001", "innodb_undo_002": "/var/lib/mysql/undo_002"},
		tablespaceDirs:           make(map[string]bool),
		dataDir:                  "/var/lib/mysql",
		userVars:                 make(map[string]userVarValue),
		systemTimeZone:           session.SystemTimeZone,
		globalTimeZone:           "SYSTEM",
		globalSQLMode:            "DEFAULT",
		globalCharsetClient:      session.CharsetClient,
		globalCollationConn:      session.CollationConnection,
		globalForeignKeyChecks:   session.ForeignKeyChecks,
		globalGenerateGIPK:       session.GenerateInvisiblePrimaryKey,
		globalShowGIPK:           session.ShowGIPK,
		globalExplicitDefaultsTS: session.ExplicitDefaultsForTimestamp,
		defaultCharset:           "utf8mb4",
		defaultCollation:         "utf8mb4_0900_ai_ci",
		charsetClient:            session.CharsetClient,
		collationConn:            session.CollationConnection,
		foreignKeyChecks:         session.ForeignKeyChecks,
		generateGIPK:             session.GenerateInvisiblePrimaryKey,
		showGIPK:                 session.ShowGIPK,
		session:                  session,
	}
}

func defaultSessionState() SessionState {
	return SessionState{
		ForeignKeyChecks:             true,
		CharsetClient:                "utf8mb4",
		CollationConnection:          "utf8mb4_0900_ai_ci",
		ExplicitDefaultsForTimestamp: true,
		SQLMode:                      "DEFAULT",
		TimeZone:                     "SYSTEM",
		SystemTimeZone:               "UTC",
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

// SystemTimeZone returns the server's configured system_time_zone (default "UTC")
// used whenever the session or global time_zone is "SYSTEM".
func (c *Catalog) SystemTimeZone() string {
	if tz := strings.TrimSpace(c.session.SystemTimeZone); tz != "" {
		return tz
	}
	if tz := strings.TrimSpace(c.systemTimeZone); tz != "" {
		return tz
	}
	return "UTC"
}

// SetSystemTimeZone configures the server's system_time_zone resolved when
// session time_zone is "SYSTEM".
func (c *Catalog) SetSystemTimeZone(tz string) {
	tz = strings.TrimSpace(tz)
	if tz == "" {
		tz = "UTC"
	}
	c.systemTimeZone = tz
	c.session.SystemTimeZone = tz
	c.systemTimeZoneLoc = nil
}

// SystemTimeZoneLocation returns the host *time.Location associated with
// @@system_time_zone when explicitly configured, or nil.
func (c *Catalog) SystemTimeZoneLocation() *time.Location {
	return c.systemTimeZoneLoc
}

// SetSystemTimeZoneLocation configures the host *time.Location used whenever
// session time_zone is "SYSTEM", preserving local/DST semantics even when
// @@system_time_zone reports an OS abbreviation such as "EST" or "EDT".
func (c *Catalog) SetSystemTimeZoneLocation(loc *time.Location) {
	c.systemTimeZoneLoc = loc
}

// SetNamedTimeZones replaces the catalog's loaded named time zone table
// (modeling mysql.time_zone_name) by loading each IANA zone name via
// time.LoadLocation. Unloadable or empty zone names are ignored.
func (c *Catalog) SetNamedTimeZones(zones ...string) {
	m := make(map[string]*time.Location, len(zones))
	for _, z := range zones {
		name := strings.TrimSpace(z)
		if name == "" {
			continue
		}
		if loc, err := time.LoadLocation(name); err == nil && loc != nil {
			m[toLower(name)] = loc
		}
	}
	c.namedTimeZones = m
}

// SetNamedTimeZone registers or removes a single named time zone entry
// (modeling a row in mysql.time_zone_name). Passing loc == nil deletes it.
func (c *Catalog) SetNamedTimeZone(name string, loc *time.Location) {
	key := toLower(strings.TrimSpace(name))
	if key == "" {
		return
	}
	if loc == nil {
		delete(c.namedTimeZones, key)
		return
	}
	if c.namedTimeZones == nil {
		c.namedTimeZones = make(map[string]*time.Location)
	}
	c.namedTimeZones[key] = loc
}

func (c *Catalog) lookupNamedTimeZone(name string) (*time.Location, bool) {
	if c == nil || len(c.namedTimeZones) == 0 {
		return nil, false
	}
	loc, ok := c.namedTimeZones[toLower(strings.TrimSpace(name))]
	return loc, ok && loc != nil
}

func (c *Catalog) effectiveSessionTimeZone() string {
	tz := strings.TrimSpace(c.session.TimeZone)
	if tz == "" || strings.EqualFold(tz, "SYSTEM") {
		return c.SystemTimeZone()
	}
	return tz
}

func (c *Catalog) resolveEffectiveTimeZone() resolvedTimeZone {
	tz := strings.TrimSpace(c.session.TimeZone)
	if tz == "" || strings.EqualFold(tz, "SYSTEM") {
		if c.systemTimeZoneLoc != nil {
			return resolvedTimeZone{
				raw:   c.SystemTimeZone(),
				valid: true,
				loc:   c.systemTimeZoneLoc,
			}
		}
		return c.resolveTimeZoneString(c.SystemTimeZone())
	}
	return c.resolveTimeZoneString(tz)
}

// DataDir returns the configured MySQL/InnoDB data directory (default "/var/lib/mysql").
func (c *Catalog) DataDir() string {
	if d := strings.TrimSpace(c.dataDir); d != "" {
		return d
	}
	return "/var/lib/mysql"
}

// SetDataDir sets the server's data directory. Replacing the active datadir
// automatically replaces the implicit datadir allowance without removing
// directories explicitly configured via AddTablespaceDir or SetTablespaceDirs.
func (c *Catalog) SetDataDir(dir string) {
	norm := normalizeTablespaceConfigDir(dir)
	if norm == "" {
		return
	}
	c.dataDir = norm
}

// AddTablespaceDir registers one or more known InnoDB tablespace directories
// (accepting individual paths or ';'-separated innodb_directories lists).
func (c *Catalog) AddTablespaceDir(dirs ...string) {
	if c.tablespaceDirs == nil {
		c.tablespaceDirs = make(map[string]bool)
	}
	for _, raw := range dirs {
		for _, part := range strings.Split(raw, ";") {
			if norm := normalizeTablespaceConfigDir(part); norm != "" {
				c.tablespaceDirs[norm] = true
			}
		}
	}
}

// SetTablespaceDirs replaces the set of explicitly configured InnoDB
// tablespace directories (DataDir remains allowed implicitly).
func (c *Catalog) SetTablespaceDirs(dirs ...string) {
	c.tablespaceDirs = make(map[string]bool)
	c.AddTablespaceDir(dirs...)
}

func normalizeTablespaceConfigDir(dir string) string {
	s := strings.TrimSpace(strings.ReplaceAll(dir, `\`, "/"))
	if s == "" {
		return ""
	}
	return path.Clean(s)
}

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
