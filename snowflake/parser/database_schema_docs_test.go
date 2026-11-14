package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Documented DATABASE / SCHEMA forms (compat.tsv docs rows):
// https://docs.snowflake.com/en/sql-reference/sql/create-database
// https://docs.snowflake.com/en/sql-reference/sql/create-database-catalog-linked
// https://docs.snowflake.com/en/sql-reference/sql/alter-database
// https://docs.snowflake.com/en/sql-reference/sql/alter-database-catalog-linked
// https://docs.snowflake.com/en/sql-reference/sql/create-schema
// https://docs.snowflake.com/en/sql-reference/sql/alter-schema
// https://docs.snowflake.com/en/sql-reference/sql/create-versioned-schema
// ---------------------------------------------------------------------------

// dbParseOK parses sql strictly and returns its single statement.
func dbParseOK(t *testing.T, sql string) ast.Node {
	t.Helper()
	f, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(f.Stmts) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(f.Stmts))
	}
	return f.Stmts[0]
}

// optionNames returns the names of opts.
func optionNames(opts []*ast.CopyOption) []string {
	var names []string
	for _, o := range opts {
		names = append(names, o.Name)
	}
	return names
}

func TestCreateDatabase_OpenEndedProps(t *testing.T) {
	sql := "CREATE OR ALTER DATABASE d DATA_RETENTION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = v CATALOG = c " +
		"LOG_LEVEL = 'INFO' OBJECT_VISIBILITY = PRIVILEGED ICEBERG_VERSION_DEFAULT = 2 COMMENT = 'x'"
	stmt := dbParseOK(t, sql).(*ast.CreateDatabaseStmt)
	if stmt.Props.DataRetention == nil || *stmt.Props.DataRetention != 1 {
		t.Errorf("DataRetention = %v, want 1", stmt.Props.DataRetention)
	}
	if stmt.Props.Comment == nil || *stmt.Props.Comment != "x" {
		t.Errorf("Comment = %v, want x", stmt.Props.Comment)
	}
	got := strings.Join(optionNames(stmt.Props.Options), ",")
	if want := "EXTERNAL_VOLUME,CATALOG,LOG_LEVEL,OBJECT_VISIBILITY,ICEBERG_VERSION_DEFAULT"; got != want {
		t.Fatalf("options = %s, want %s", got, want)
	}
	if v := stmt.Props.Options[3].Words; v != "PRIVILEGED" {
		t.Errorf("OBJECT_VISIBILITY = %q, want PRIVILEGED", v)
	}
	if o := stmt.Props.Options[1]; sql[o.Loc.Start:o.Loc.End] != "CATALOG = c" {
		t.Errorf("CATALOG option Loc covers %q", sql[o.Loc.Start:o.Loc.End])
	}
}

func TestCreateDatabase_CatalogLinked(t *testing.T) {
	sql := "CREATE DATABASE d LINKED_CATALOG = (CATALOG = 'cat', ALLOWED_NAMESPACES = ('a', 'b'), " +
		"ALLOWED_WRITE_OPERATIONS = NONE, SYNC_INTERVAL_SECONDS = 60) EXTERNAL_VOLUME = 'v' " +
		"CATALOG_CASE_SENSITIVITY = CASE_INSENSITIVE WITH TAG (t = 'x') WITH CONTACT (SUPPORT = c)"
	stmt := dbParseOK(t, sql).(*ast.CreateDatabaseStmt)
	if len(stmt.Props.Options) != 3 || stmt.Props.Options[0].Name != "LINKED_CATALOG" {
		t.Fatalf("options = %v", optionNames(stmt.Props.Options))
	}
	group := stmt.Props.Options[0].Group
	if got := strings.Join(optionNames(group), ","); got != "CATALOG,ALLOWED_NAMESPACES,ALLOWED_WRITE_OPERATIONS,SYNC_INTERVAL_SECONDS" {
		t.Fatalf("LINKED_CATALOG group = %s", got)
	}
	if len(group[1].List) != 2 || group[1].List[1].Value != "b" {
		t.Errorf("ALLOWED_NAMESPACES list = %v", group[1].List)
	}
	if len(stmt.Tags) != 1 || len(stmt.Props.Contacts) != 1 || stmt.Props.Contacts[0].Words != "C" {
		t.Errorf("tags = %d, contacts = %v", len(stmt.Tags), stmt.Props.Contacts)
	}
}

func TestCreateDatabase_Sources(t *testing.T) {
	stmt := dbParseOK(t, "CREATE DATABASE d FROM LISTING 'GZ1Z2Z3'").(*ast.CreateDatabaseStmt)
	if stmt.FromListing == nil || *stmt.FromListing != "GZ1Z2Z3" {
		t.Errorf("FromListing = %v", stmt.FromListing)
	}
	stmt = dbParseOK(t, "CREATE DATABASE d FROM BACKUP SET bs IDENTIFIER 'id1'").(*ast.CreateDatabaseStmt)
	if stmt.FromBackupSet == nil || stmt.FromBackupSet.Normalize() != "BS" || stmt.BackupID != "id1" {
		t.Errorf("FromBackupSet = %v, BackupID = %q", stmt.FromBackupSet, stmt.BackupID)
	}
	stmt = dbParseOK(t, "CREATE DATABASE d AS REPLICA OF acct.primary_db DATA_RETENTION_TIME_IN_DAYS = 1").(*ast.CreateDatabaseStmt)
	if stmt.ReplicaOf == nil || stmt.ReplicaOf.Normalize() != "ACCT.PRIMARY_DB" || stmt.Props.DataRetention == nil {
		t.Errorf("ReplicaOf = %v", stmt.ReplicaOf)
	}
	stmt = dbParseOK(t, "CREATE DATABASE d CLONE s AT (OFFSET => -60) IGNORE TABLES WITH INSUFFICIENT DATA RETENTION IGNORE HYBRID TABLES").(*ast.CreateDatabaseStmt)
	if stmt.Clone == nil || !stmt.IgnoreInsufficientRetention || !stmt.IgnoreHybridTables {
		t.Errorf("clone = %v, ignore = %v/%v", stmt.Clone, stmt.IgnoreInsufficientRetention, stmt.IgnoreHybridTables)
	}
}

func TestCreateSchema_DocumentedForms(t *testing.T) {
	stmt := dbParseOK(t, `CREATE SCHEMA s CLONE t IGNORE HYBRID TABLES WITH MANAGED ACCESS CATALOG = c WITH CONTACT (SUPPORT = c1, ACCESS_APPROVAL = db.sc."Mixed")`).(*ast.CreateSchemaStmt)
	if stmt.Clone == nil || stmt.IgnoreInsufficientRetention || !stmt.IgnoreHybridTables || !stmt.ManagedAccess {
		t.Errorf("clone = %v, ignore hybrid = %v, managed = %v", stmt.Clone, stmt.IgnoreHybridTables, stmt.ManagedAccess)
	}
	if len(stmt.Props.Contacts) != 2 || stmt.Props.Contacts[1].Words != `DB.SC."Mixed"` {
		t.Errorf("contacts = %v", stmt.Props.Contacts)
	}
	stmt = dbParseOK(t, "CREATE SCHEMA s FROM BACKUP SET bs IDENTIFIER 'id1'").(*ast.CreateSchemaStmt)
	if stmt.FromBackupSet == nil || stmt.BackupID != "id1" {
		t.Errorf("FromBackupSet = %v, BackupID = %q", stmt.FromBackupSet, stmt.BackupID)
	}
}

func TestCreateVersionedSchema(t *testing.T) {
	stmt := dbParseOK(t, "CREATE OR ALTER VERSIONED SCHEMA s WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 COMMENT = 'x'").(*ast.CreateSchemaStmt)
	if !stmt.Versioned || !stmt.OrAlter || !stmt.ManagedAccess || stmt.Props.Comment == nil {
		t.Errorf("versioned = %v, orAlter = %v, managed = %v", stmt.Versioned, stmt.OrAlter, stmt.ManagedAccess)
	}
	for _, sql := range []string{
		"CREATE VERSIONED SCHEMA s",
		"CREATE OR REPLACE VERSIONED SCHEMA s",
		"CREATE OR ALTER TRANSIENT VERSIONED SCHEMA s",
		"CREATE OR ALTER VERSIONED SCHEMA s CLONE t",
		"CREATE OR ALTER VERSIONED SCHEMA s WITH TAG (t = 'x')",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

func TestAlterDatabase_DocumentedActions(t *testing.T) {
	stmt := dbParseOK(t, "ALTER DATABASE d SET FEATURE POLICY p FORCE").(*ast.AlterDatabaseStmt)
	if stmt.Action != ast.AlterDBSetFeaturePolicy || stmt.FeaturePolicy.Normalize() != "P" || !stmt.Force {
		t.Errorf("SET FEATURE POLICY: action = %v, policy = %v, force = %v", stmt.Action, stmt.FeaturePolicy, stmt.Force)
	}
	if stmt := dbParseOK(t, "ALTER DATABASE d UNSET FEATURE POLICY").(*ast.AlterDatabaseStmt); stmt.Action != ast.AlterDBUnsetFeaturePolicy {
		t.Errorf("UNSET FEATURE POLICY: action = %v", stmt.Action)
	}
	if stmt := dbParseOK(t, "ALTER DATABASE IF EXISTS d SUSPEND DISCOVERY").(*ast.AlterDatabaseStmt); stmt.Action != ast.AlterDBSuspendDiscovery {
		t.Errorf("SUSPEND DISCOVERY: action = %v", stmt.Action)
	}
	if stmt := dbParseOK(t, "ALTER DATABASE d RESUME DISCOVERY").(*ast.AlterDatabaseStmt); stmt.Action != ast.AlterDBResumeDiscovery {
		t.Errorf("RESUME DISCOVERY: action = %v", stmt.Action)
	}

	stmt = dbParseOK(t, "ALTER DATABASE d SET DATA_RETENTION_TIME_IN_DAYS = 1 CONTACT SUPPORT = c1, ACCESS_APPROVAL = c2 ENABLE_DATA_COMPACTION = TRUE").(*ast.AlterDatabaseStmt)
	if stmt.Action != ast.AlterDBSet || len(stmt.SetProps.Contacts) != 2 || len(stmt.SetProps.Options) != 1 {
		t.Errorf("SET CONTACT: props = %+v", stmt.SetProps)
	}
	stmt = dbParseOK(t, "ALTER DATABASE d UNSET CONTACT support, COMMENT, OBJECT_VISIBILITY").(*ast.AlterDatabaseStmt)
	if got := strings.Join(stmt.UnsetProps, ","); got != "CONTACT SUPPORT,COMMENT,OBJECT_VISIBILITY" {
		t.Errorf("UNSET props = %s", got)
	}
	stmt = dbParseOK(t, "ALTER DATABASE d SET TAG t1 = 'a', t2 = 'b'").(*ast.AlterDatabaseStmt)
	if stmt.Action != ast.AlterDBSetTag || len(stmt.Tags) != 2 {
		t.Errorf("SET TAG: action = %v, tags = %d", stmt.Action, len(stmt.Tags))
	}
	stmt = dbParseOK(t, "ALTER DATABASE d UNSET TAG t1, t2").(*ast.AlterDatabaseStmt)
	if stmt.Action != ast.AlterDBUnsetTag || len(stmt.UnsetTags) != 2 {
		t.Errorf("UNSET TAG: action = %v, tags = %d", stmt.Action, len(stmt.UnsetTags))
	}
}

func TestAlterDatabase_UpdateLinkedCatalog(t *testing.T) {
	tests := []struct {
		sql, op, target string
		namespaces      int
		options         int
	}{
		{"ALTER DATABASE d UPDATE LINKED_CATALOG ADD ('a', 'b') TO ALLOWED_NAMESPACES", "ADD", "ALLOWED_NAMESPACES", 2, 0},
		{"ALTER DATABASE d UPDATE LINKED_CATALOG REMOVE ('a') FROM BLOCKED_NAMESPACES", "REMOVE", "BLOCKED_NAMESPACES", 1, 0},
		{"ALTER DATABASE d UPDATE LINKED_CATALOG UNSET ALLOWED_NAMESPACES", "UNSET", "ALLOWED_NAMESPACES", 0, 0},
		{"ALTER DATABASE d UPDATE LINKED_CATALOG SET SYNC_INTERVAL_SECONDS = 60", "SET", "", 0, 1},
	}
	for _, tt := range tests {
		stmt := dbParseOK(t, tt.sql).(*ast.AlterDatabaseStmt)
		u := stmt.LinkedCatalog
		if stmt.Action != ast.AlterDBUpdateLinkedCatalog || u == nil {
			t.Fatalf("%s: action = %v", tt.sql, stmt.Action)
		}
		if u.Op != tt.op || u.Target != tt.target || len(u.Namespaces) != tt.namespaces || len(u.Options) != tt.options {
			t.Errorf("%s: got %+v", tt.sql, u)
		}
	}
}

func TestAlterSchema_DocumentedActions(t *testing.T) {
	stmt := dbParseOK(t, "ALTER SCHEMA s SET LOG_LEVEL = 'INFO' CONTACT SUPPORT = c STORAGE_SERIALIZATION_POLICY = OPTIMIZED").(*ast.AlterSchemaStmt)
	if stmt.Action != ast.AlterSchemaSet || len(stmt.SetProps.Options) != 2 || len(stmt.SetProps.Contacts) != 1 {
		t.Errorf("SET: props = %+v", stmt.SetProps)
	}
	stmt = dbParseOK(t, "ALTER SCHEMA s UNSET CONTACT SUPPORT").(*ast.AlterSchemaStmt)
	if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "CONTACT SUPPORT" {
		t.Errorf("UNSET CONTACT: props = %v", stmt.UnsetProps)
	}
	stmt = dbParseOK(t, "ALTER SCHEMA s UNSET TAG t").(*ast.AlterSchemaStmt)
	if stmt.Action != ast.AlterSchemaUnsetTag || len(stmt.UnsetTags) != 1 {
		t.Errorf("UNSET TAG: action = %v", stmt.Action)
	}
}

func TestDatabaseSchema_Rejects(t *testing.T) {
	for _, sql := range []string{
		"CREATE DATABASE a zzz",
		"CREATE DATABASE a CONTACT (a = a)",
		"CREATE DATABASE a WITH CONTACT (a a)",
		"CREATE DATABASE a LINKED_CATALOG = (CATALOG 's')",
		"CREATE DATABASE a LINKED_CATALOG = (CATALOG = 's' zzz)",
		"CREATE DATABASE a LINKED_CATALOG = (CATALOG = 's', ALLOWED_NAMESPACES = ())",
		"CREATE DATABASE a FROM LISTING x",
		"CREATE DATABASE a FROM BACKUP SET b",
		"CREATE DATABASE a CLONE b IGNORE TABLES",
		"ALTER DATABASE a SET",
		"ALTER DATABASE a SET CONTACT a = a, = a",
		"ALTER DATABASE a SET 'FEATURE' POLICY p",
		"ALTER DATABASE a SUSPEND",
		"ALTER DATABASE a UPDATE LINKED_CATALOG SET",
		"ALTER DATABASE a UPDATE LINKED_CATALOG ADD () TO ALLOWED_NAMESPACES",
		"ALTER DATABASE a UNSET DATA_RETENTION_TIME_IN_DAYS zzz",
		"ALTER SCHEMA a SET",
		"ALTER SCHEMA a UNSET CONTACT a zzz",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

func TestDatabaseSchema_Lookahead(t *testing.T) {
	newP := func(sql string) *Parser {
		p := &Parser{lexer: NewLexer(sql), input: sql}
		p.advance()
		return p
	}
	if !newP("CATALOG = c").startsDBSchemaOption() || newP("CATALOG c").startsDBSchemaOption() ||
		newP("WITH TAG (a = 'b')").startsDBSchemaOption() || newP("'k' = v").startsDBSchemaOption() {
		t.Error("startsDBSchemaOption")
	}
	if !newP("CONTACT support = c").startsInlineContact() || newP("CONTACT = c").startsInlineContact() {
		t.Error("startsInlineContact")
	}
	if !newP("FEATURE POLICY p").startsFeaturePolicy() || newP("FEATURE = 1").startsFeaturePolicy() {
		t.Error("startsFeaturePolicy")
	}
	if !tokIsWord(newP("discovery").cur, "DISCOVERY") || tokIsWord(newP("'DISCOVERY'").cur, "DISCOVERY") ||
		tokIsWord(newP(`"DISCOVERY"`).cur, "DISCOVERY") {
		t.Error("tokIsWord")
	}
}
