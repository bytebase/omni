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
