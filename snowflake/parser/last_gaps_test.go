package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseObjectBodies(t *testing.T) {
	s := mustObjectStmt(t, "CREATE STORAGE LIFECYCLE POLICY IF NOT EXISTS p AS (event_ts TIMESTAMP, id NUMBER) "+
		"RETURNS BOOLEAN -> event_ts < CURRENT_TIMESTAMP() AND id > 1 ARCHIVE_TIER = COOL ARCHIVE_FOR_DAYS = 180 WITH TAG (t = 'v')")
	if got := strings.Join(clauseKeywords(s), "|"); got != "AS|" {
		t.Fatalf("clauses = %q", got)
	}
	if c := s.Clauses[0]; len(c.Params) != 2 || c.Type == nil || c.Body == nil {
		t.Errorf("AS = %+v", c)
	}
	s = mustObjectStmt(t, "ALTER STORAGE LIFECYCLE POLICY p SET BODY -> a > 1")
	if c := s.Clauses[0]; c.Keywords != "SET BODY" || c.Body == nil {
		t.Errorf("SET BODY = %+v", c)
	}

	s = mustObjectStmt(t, "CREATE OR REPLACE DATA MOVEMENT RULE r TYPE = 'COPY_INTO_EXTERNAL_STAGE' MAX_ROWS AS () RETURNS INTEGER -> (1000) COMMENT = 'c'")
	if got := strings.Join(clauseKeywords(s), "|"); got != "|MAX_ROWS AS|" {
		t.Fatalf("clauses = %q", got)
	}
	if c := s.Clauses[1]; c.Params == nil || len(c.Params) != 0 || c.Type == nil {
		t.Errorf("MAX_ROWS AS = %+v", c)
	} else if _, ok := c.Body.(*ast.ParenExpr); !ok {
		t.Errorf("body = %T", c.Body)
	}
	s = mustObjectStmt(t, "ALTER DATA MOVEMENT RULE IF EXISTS r SET MAX_ROWS AS () RETURNS INTEGER -> (500)")
	if c := s.Clauses[0]; c.Keywords != "SET MAX_ROWS AS" || c.Loc.Start != len("ALTER DATA MOVEMENT RULE IF EXISTS r ") {
		t.Errorf("SET MAX_ROWS AS = %+v", c)
	}

	s = mustObjectStmt(t, "CREATE OR REPLACE PROVISIONED THROUGHPUT pt CLOUD_PROVIDER = 'aws' MODEL = 'm' PTUS = 64 TERM_START = 'a' TERM_END = 'b'")
	if s.Kind != "PROVISIONED THROUGHPUT" || len(s.Clauses) != 1 || len(s.Clauses[0].Options) != 5 {
		t.Errorf("PROVISIONED THROUGHPUT = %+v", s)
	}

	for _, in := range []string{
		"CREATE STORAGE LIFECYCLE POLICY p AS () RETURNS BOOLEAN -> 1",
		"CREATE STORAGE LIFECYCLE POLICY p AS (a INT) RETURNS INT -> 1",
		"CREATE STORAGE LIFECYCLE POLICY p AS (a INT) RETURNS BOOLEAN",
		"CREATE STORAGE LIFECYCLE POLICY p AS (a INT) RETURNS BOOLEAN -> 1 ARCHIVE_TIER = WARM",
		"ALTER STORAGE LIFECYCLE POLICY p SET BODY 1",
		"CREATE DATA MOVEMENT RULE r TYPE = 'x' MAX_ROWS AS () RETURNS INTEGER -> 1", // the body is parenthesized
		"CREATE DATA MOVEMENT RULE r TYPE = 'x' MAX_ROWS AS () RETURNS INTEGER -> (1) + 1",
		"CREATE DATA MOVEMENT RULE r TYPE = 'x' MAX_ROWS AS (a INT) RETURNS INTEGER -> (1)",
		"CREATE DATA MOVEMENT RULE r TYPE = 'x' MAX_ROWS AS () RETURNS NUMBER -> (1)",
		"CREATE DATA MOVEMENT RULE r MAX_ROWS AS () RETURNS INTEGER -> (1)", // TYPE is required
		"CREATE DATA MOVEMENT RULE r TYPE = x MAX_ROWS AS () RETURNS INTEGER -> (1)",
		"CREATE DATA MOVEMENT RULE r TYPE = 'x' MAX_ROWS AS () RETURNS INTEGER -> (1) TYPE = 'y'",
		"ALTER DATA MOVEMENT RULE r SET TYPE = 'x'",
		"ALTER DATA MOVEMENT RULE r UNSET TYPE",
		"CREATE PROVISIONED THROUGHPUT pt CLOUD_PROVIDER = 'aws' MODEL = 'm' PTUS = 64 TERM_START = 'a'",
		"CREATE PROVISIONED THROUGHPUT IF NOT EXISTS pt CLOUD_PROVIDER = 'aws' MODEL = 'm' PTUS = 64 TERM_START = 'a' TERM_END = 'b'",
		"CREATE PROVISIONED THROUGHPUT pt CLOUD_PROVIDER = 'aws' MODEL = 'm' PTUS = 'x' TERM_START = 'a' TERM_END = 'b'",
		"CREATE PROVISIONED THROUGHPUT pt CLOUD_PROVIDER = 'aws' MODEL = 'm' PTUS = 1 TERM_START = 'a' TERM_END = 'b' TAG (t = 'v')",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseStreamPipeSessionCallNetworkGaps(t *testing.T) {
	for in, want := range map[string]ast.StreamSourceKind{
		"CREATE STREAM s ON EVENT TABLE t COMMENT = 'c'":         ast.StreamOnEventTable,
		"CREATE STREAM s ON DYNAMIC TABLE db.s.t":                ast.StreamOnDynamicTable,
		"CREATE STREAM s ON EXTERNAL TABLE t INSERT_ONLY = TRUE": ast.StreamOnExternalTable,
	} {
		st, ok := mustParseOne(t, in).(*ast.CreateStreamStmt)
		if !ok || st.SourceKind != want {
			t.Errorf("%q: got %+v", in, st)
		}
	}
	st := mustParseOne(t, "CREATE STREAM s CLONE t COPY GRANTS COPY = (GRANTS)").(*ast.CreateStreamStmt)
	if st.Clone == nil || !st.CopyGrants || len(st.Options) != 1 || st.Options[0].Name != "COPY" {
		t.Errorf("CLONE = %+v", st)
	}

	pipe := mustParseOne(t, "ALTER PIPE p SET ERROR_INTEGRATION = i, PIPE_EXECUTION_PAUSED = TRUE, COMMENT = 'c'").(*ast.AlterPipeStmt)
	if len(pipe.Options) != 3 {
		t.Errorf("ALTER PIPE = %+v", pipe)
	}
	sess := mustParseOne(t, "ALTER SESSION SET QUERY_TAG = 'x', TIMEZONE = 'UTC' LOCK_TIMEOUT = 1").(*ast.AlterSessionStmt)
	if len(sess.Options) != 3 {
		t.Errorf("ALTER SESSION = %+v", sess)
	}
	call := mustParseOne(t, "CALL p(1, x => 2) INTO :v").(*ast.CallStmt)
	if call.Into.Name != "v" || len(call.Args) != 2 || call.Loc.End != len("CALL p(1, x => 2) INTO :v") {
		t.Errorf("CALL INTO = %+v", call)
	}
	for in, want := range map[string]ast.AlterPolicyAction{
		"ALTER NETWORK POLICY p ADD ALLOWED_NETWORK_RULE_LIST = ('a', 'b')": ast.AlterPolicyAdd,
		"ALTER NETWORK POLICY p REMOVE BLOCKED_NETWORK_RULE_LIST = ('a')":   ast.AlterPolicyRemove,
	} {
		ap, ok := mustParseOne(t, in).(*ast.AlterPolicyStmt)
		if !ok || ap.Action != want || len(ap.Options) != 1 {
			t.Errorf("%q: got %+v", in, ap)
		}
	}

	for _, in := range []string{
		"CREATE STREAM s ON EVENT TABLE t AT (OFFSET => -1)",
		"CREATE STREAM s ON DYNAMIC TABLE t BEFORE (STATEMENT => 'x')",
		"CREATE STREAM s ON EVENT t",
		"CREATE STREAM s CLONE t COPY = (TAGS)",
		"CREATE STREAM s CLONE t COPY = ()",
		"ALTER PIPE p SET COMMENT = 'c',",
		"ALTER PIPE p SET , COMMENT = 'c'",
		"ALTER SESSION SET QUERY_TAG = 'x',",
		"CALL p(1,)",
		"CALL p() INTO v",
		"CALL p() INTO :v.w",
		"CALL p() INTO :",
		"ALTER NETWORK POLICY IF EXISTS p ADD ALLOWED_NETWORK_RULE_LIST = ('a')",
		"ALTER NETWORK POLICY p ADD ALLOWED_NETWORK_RULE_LIST = ()",
		"ALTER NETWORK POLICY p ADD ALLOWED_IP_LIST = ('1.1.1.1')",
		"ALTER NETWORK POLICY p REMOVE ALLOWED_NETWORK_RULE_LIST = (a)",
		"ALTER SESSION POLICY p ADD ALLOWED_NETWORK_RULE_LIST = ('a')",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseLastDocForms(t *testing.T) {
	// 1. ALTER CORTEX SEARCH SERVICE ... SET AUTO_SUSPEND = NULL
	cs := mustObjectStmt(t, "ALTER CORTEX SEARCH SERVICE s SET AUTO_SUSPEND = NULL COMMENT = 'c'")
	if len(cs.Clauses) != 1 || cs.Clauses[0].Options[0].Words != "NULL" {
		t.Errorf("CORTEX SEARCH AUTO_SUSPEND = NULL: %+v", cs)
	}

	// 2. CREATE [TEMP|TEMPORARY] [SECURE] AGENT
	ag := mustObjectStmt(t, "CREATE OR REPLACE TEMP SECURE AGENT db.s.a COPY GRANTS COMMENT = 'c' FROM SPECIFICATION $$models: {}$$")
	if !ag.OrReplace || !ag.Temporary || !ag.Secure || ag.Kind != "AGENT" {
		t.Errorf("TEMP SECURE AGENT = %+v", ag)
	}
	ag = mustObjectStmt(t, "CREATE TEMPORARY AGENT IF NOT EXISTS a FROM SPECIFICATION 's'")
	if !ag.Temporary || ag.Secure || !ag.IfNotExists {
		t.Errorf("TEMPORARY AGENT = %+v", ag)
	}
	ag = mustObjectStmt(t, "CREATE SECURE AGENT a COMMENT = 'c' FROM @db.s.st/dir")
	if ag.Temporary || !ag.Secure {
		t.Errorf("SECURE AGENT stage = %+v", ag)
	}

	// 3. CREATE [TEMP|TEMPORARY] [SECURE] AGGREGATE FUNCTION ... LANGUAGE PYTHON
	fn := mustParseOne(t, "CREATE OR REPLACE TEMP SECURE AGGREGATE FUNCTION db.s.f(x INT DEFAULT 1) COPY GRANTS "+
		"RETURNS INT NOT NULL LANGUAGE PYTHON STRICT IMMUTABLE RUNTIME_VERSION = '3.11' "+
		"PACKAGES = ('numpy') IMPORTS = ('@s/f.py') ARTIFACT_REPOSITORY = 'snowflake.snowpark.pypi_shared_repository' "+
		"EXTERNAL_ACCESS_INTEGRATIONS = (e1) SECRETS = ('k' = db.s.sec) COPY = (GRANTS) COMMENT = 'c' HANDLER = 'Agg' AS $$class Agg: pass$$").(*ast.CreateRoutineStmt)
	if !fn.OrReplace || !fn.Temporary || !fn.Secure || !fn.Aggregate || !fn.CopyGrants || !fn.ReturnNotNull {
		t.Errorf("AGGREGATE FUNCTION = %+v", fn)
	}
	fn = mustParseOne(t, "CREATE AGGREGATE FUNCTION IF NOT EXISTS f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' IMPORTS = ('@s/f.py') HANDLER = 'Agg'").(*ast.CreateRoutineStmt)
	if !fn.Aggregate || !fn.IfNotExists || fn.Body != "" {
		t.Errorf("AGGREGATE FUNCTION staged = %+v", fn)
	}

	// 4. ALTER {FAILOVER|REPLICATION} GROUP external volume actions
	fg := mustAlterGroup(t, "ALTER FAILOVER GROUP IF EXISTS fg ADD v1, v2 TO ALLOWED_EXTERNAL_VOLUMES")
	if !fg.Failover || fg.Action != ast.AlterGroupAdd || fg.ListTarget != "ALLOWED_EXTERNAL_VOLUMES" || len(fg.Names) != 2 {
		t.Errorf("FAILOVER ADD ALLOWED_EXTERNAL_VOLUMES = %+v", fg)
	}
	fg = mustAlterGroup(t, "ALTER FAILOVER GROUP fg MOVE EXTERNAL VOLUMES v1, v2 TO FAILOVER GROUP fg2")
	if fg.Action != ast.AlterGroupMove || fg.MoveKind != "EXTERNAL VOLUMES" || fg.MoveTo.String() != "fg2" {
		t.Errorf("FAILOVER MOVE EXTERNAL VOLUMES = %+v", fg)
	}
	rg := mustAlterGroup(t, "ALTER REPLICATION GROUP rg REMOVE v1 FROM ALLOWED_EXTERNAL_VOLUMES")
	if rg.Failover || rg.Action != ast.AlterGroupRemove || rg.ListTarget != "ALLOWED_EXTERNAL_VOLUMES" {
		t.Errorf("REPLICATION REMOVE ALLOWED_EXTERNAL_VOLUMES = %+v", rg)
	}
	rg = mustAlterGroup(t, "ALTER REPLICATION GROUP rg MOVE EXTERNAL VOLUMES v1 TO REPLICATION GROUP rg2")
	if rg.Failover || rg.MoveKind != "EXTERNAL VOLUMES" || rg.MoveTo.String() != "rg2" {
		t.Errorf("REPLICATION MOVE EXTERNAL VOLUMES = %+v", rg)
	}

	// 5. CREATE DATABASE ROLE ... CLONE
	dr := mustCreateRole(t, "CREATE OR REPLACE DATABASE ROLE db.r1 CLONE db.r0")
	if !dr.Database || !dr.OrReplace || dr.Clone == nil || dr.Clone.String() != "db.r0" {
		t.Errorf("DATABASE ROLE CLONE = %+v", dr)
	}
	dr = mustCreateRole(t, "CREATE DATABASE ROLE IF NOT EXISTS r1 CLONE r0")
	if !dr.Database || !dr.IfNotExists || dr.Clone == nil || dr.Clone.String() != "r0" {
		t.Errorf("DATABASE ROLE IF NOT EXISTS CLONE = %+v", dr)
	}

	// 6 & 7. CREATE ICEBERG TABLE ... FROM BACKUP SET and Parquet virtual columns
	it := mustParseOne(t, "CREATE ICEBERG TABLE db.s.t FROM BACKUP SET db.s.bs IDENTIFIER 'b1'").(*ast.CreateTableStmt)
	if !it.Iceberg || it.FromBackup == nil || it.FromBackup.ID != "b1" {
		t.Errorf("ICEBERG FROM BACKUP SET = %+v", it)
	}
	it = mustParseOne(t, "CREATE ICEBERG TABLE t (a INT, b INT AS a + 1) BASE_LOCATION = 'p/'").(*ast.CreateTableStmt)
	if !it.Iceberg || len(it.Columns) != 2 || it.Columns[1].VirtualExpr == nil {
		t.Errorf("ICEBERG virtual col = %+v", it)
	}

	for _, in := range []string{
		// CORTEX SEARCH SERVICE AUTO_SUSPEND
		"CREATE CORTEX SEARCH SERVICE s ON a ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = '1 hour' AUTO_SUSPEND = NULL AS SELECT 1",
		"ALTER CORTEX SEARCH SERVICE s SET AUTO_SUSPEND = 'NULL'",
		`ALTER CORTEX SEARCH SERVICE s SET AUTO_SUSPEND = "NULL"`,
		"ALTER CORTEX SEARCH SERVICE s SET AUTO_SUSPEND = NONE",
		// AGENT modifiers
		"CREATE SECURE TEMP AGENT a FROM SPECIFICATION 's'",
		"CREATE TEMP AGENT a FROM @st/dir",
		"CREATE TRANSIENT AGENT a FROM SPECIFICATION 's'",
		"CREATE LOCAL TEMP AGENT a FROM SPECIFICATION 's'",
		"CREATE OR ALTER TEMP AGENT a FROM SPECIFICATION 's'",
		"CREATE TEMP EXTERNAL AGENT e",
		"CREATE SECURE EXTERNAL AGENT e",
		// AGGREGATE FUNCTION
		"CREATE AGGREGATE FUNCTION f() RETURNS INT AS '1'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE JAVA RUNTIME_VERSION = '11' HANDLER = 'h' AS 's'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON HANDLER = 'h' AS 's'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' AS 's'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' HANDLER = 'h'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON VOLATILE IMMUTABLE RUNTIME_VERSION = '3.11' HANDLER = 'h' AS 's'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON STRICT CALLED ON NULL INPUT RUNTIME_VERSION = '3.11' HANDLER = 'h' AS 's'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' HANDLER = 'h' COPY = (TAGS) AS 's'",
		"CREATE AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' HANDLER = 'h' UNKNOWN_PROP = 's' AS 's'",
		"CREATE OR ALTER AGGREGATE FUNCTION f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' HANDLER = 'h' AS 's'",
		"CREATE OR REPLACE AGGREGATE FUNCTION IF NOT EXISTS f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' HANDLER = 'h' AS 's'",
		`CREATE "AGGREGATE" FUNCTION f() RETURNS INT LANGUAGE PYTHON RUNTIME_VERSION = '3.11' HANDLER = 'h' AS 's'`,
		// FAILOVER / REPLICATION GROUP external volumes
		"ALTER FAILOVER GROUP fg ADD v1 TO ALLOWED_EXTERNAL_VOLUMES IGNORE EDITION CHECK",
		"ALTER FAILOVER GROUP fg ADD db.v1 TO ALLOWED_EXTERNAL_VOLUMES",
		`ALTER FAILOVER GROUP fg ADD v1 TO "ALLOWED_EXTERNAL_VOLUMES"`,
		"ALTER FAILOVER GROUP fg MOVE EXTERNAL VOLUMES v1 TO REPLICATION GROUP rg",
		"ALTER REPLICATION GROUP rg MOVE EXTERNAL VOLUMES v1 TO FAILOVER GROUP fg",
		`ALTER FAILOVER GROUP fg MOVE EXTERNAL "VOLUMES" v1 TO FAILOVER GROUP fg2`,
		"ALTER FAILOVER GROUP fg MOVE EXTERNAL VOLUMES db.v1 TO FAILOVER GROUP fg2",
		// DATABASE ROLE CLONE
		"CREATE ROLE r1 CLONE r0",
		"CREATE OR ALTER DATABASE ROLE r1 CLONE r0",
		"CREATE OR REPLACE DATABASE ROLE IF NOT EXISTS r1 CLONE r0",
		"CREATE DATABASE ROLE a.b.c CLONE r0",
		"CREATE DATABASE ROLE r1 CLONE a.b.c",
		"CREATE DATABASE ROLE r1 CLONE r0 COMMENT = 'c'",
		// ICEBERG TABLE FROM BACKUP SET and virtual columns
		"CREATE OR REPLACE ICEBERG TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'",
		"CREATE TRANSIENT ICEBERG TABLE t FROM BACKUP SET bs IDENTIFIER 'b1'",
		"CREATE ICEBERG TABLE IF NOT EXISTS t FROM BACKUP SET bs IDENTIFIER 'b1'",
		`CREATE ICEBERG TABLE t FROM "BACKUP" SET bs IDENTIFIER 'b1'`,
		"CREATE ICEBERG TABLE t (a INT, b AS 1) BASE_LOCATION = 'p/'",
		"CREATE TRANSIENT ICEBERG TABLE t (a INT, b INT AS 1) BASE_LOCATION = 'p/'",
		"CREATE ICEBERG TABLE t (a INT, b INT AS 1) AS SELECT 1, 2",
		"CREATE ICEBERG TABLE t (a INT, b INT AS 1 DEFAULT 2) BASE_LOCATION = 'p/'",
		"CREATE TABLE t (a INT, b INT AS 1)",
	} {
		mustRejectObject(t, in)
	}
}
