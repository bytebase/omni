package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseCreateService(t *testing.T) {
	in := "CREATE SERVICE IF NOT EXISTS db.s.svc IN COMPUTE POOL pool1 FROM SPECIFICATION_TEMPLATE $$spec: {{ x }}$$ " +
		"USING (x => 'v', \"Mixed\" => 2, y => true) MIN_INSTANCES = 1 EXTERNAL_ACCESS_INTEGRATIONS = (eai1, eai2) " +
		"WITH TAG (t = 'v') COMMENT = 'c'"
	s := mustObjectStmt(t, in)
	if s.Verb != "CREATE" || s.Kind != "SERVICE" || !s.IfNotExists || s.Name.String() != "db.s.svc" {
		t.Fatalf("got %+v", s)
	}
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|FROM|USING|" {
		t.Fatalf("clauses = %q", got)
	}
	if pool := s.Clauses[0].Names; len(pool) != 1 || pool[0].String() != "pool1" {
		t.Errorf("pool = %+v", pool)
	}
	if spec, ok := optLitString(s.Clauses[1].Options, "SPECIFICATION_TEMPLATE"); !ok || spec != "spec: {{ x }}" {
		t.Errorf("specification = %q", spec)
	}
	using := s.Clauses[2].Options
	if len(using) != 3 || using[0].Name != "X" || using[0].Lit.Value != "v" || using[1].Name != "Mixed" ||
		using[1].Lit.Value != "2" || using[2].Words != "TRUE" {
		t.Errorf("USING = %+v", using)
	}
	if props := s.Clauses[3]; len(props.Options) != 3 || len(props.Tags) != 1 {
		t.Errorf("properties = %+v", props)
	}

	for in, want := range map[string]string{
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'spec'":                          "SPECIFICATION",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_FILE = 'a.yaml'":                 "SPECIFICATION_FILE",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE_FILE = 'a' USING(a=>1)": "SPECIFICATION_TEMPLATE_FILE",
	} {
		s := mustObjectStmt(t, in)
		if o := s.Clauses[1].Options; len(o) != 1 || o[0].Name != want {
			t.Errorf("%q: FROM = %+v", in, o)
		}
	}
}

func TestParseAlterService(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER SERVICE IF EXISTS s SUSPEND":                                 "SUSPEND",
		"ALTER SERVICE s RESUME":                                            "RESUME",
		"ALTER SERVICE s FROM SPECIFICATION 'x'":                            "FROM",
		"ALTER SERVICE s FROM SPECIFICATION_TEMPLATE 'x' USING (a => 1)":    "FROM|USING",
		"ALTER SERVICE s SET MIN_INSTANCES = 1 QUERY_WAREHOUSE = wh":        "SET",
		"ALTER SERVICE s SET TAG t = 'v', u = 'w'":                          "SET TAG",
		"ALTER SERVICE s UNSET MIN_INSTANCES, EXTERNAL_ACCESS_INTEGRATIONS": "UNSET",
	} {
		s := mustObjectStmt(t, in)
		if got := strings.Join(clauseKeywords(s), "|"); s.Verb != "ALTER" || got != want {
			t.Errorf("%q: clauses = %q, want %q", in, got, want)
		}
	}
}

func TestParseExecuteJobService(t *testing.T) {
	in := "EXECUTE JOB SERVICE IN COMPUTE POOL p FROM SPECIFICATION 'spec' NAME = db.s.job ASYNC = TRUE REPLICAS = 2"
	s := mustObjectStmt(t, in)
	if s.Verb != "EXECUTE" || s.Kind != "JOB SERVICE" || s.Name != nil {
		t.Fatalf("got %+v", s)
	}
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|FROM|" {
		t.Fatalf("clauses = %q", got)
	}
	if w, ok := optWords(s.Clauses[2].Options, "NAME"); !ok || w != "DB.S.JOB" {
		t.Errorf("NAME = %q", w)
	}
	s = mustObjectStmt(t, "execute job service in compute pool p from specification_template 'x' using (a => 'b')")
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|FROM|USING" {
		t.Errorf("clauses = %q", got)
	}

	in = "EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 'spec' FROM (SELECT a FROM t) " +
		"MODEL = db.s.m VERSION = v1 FUNCTION = 'predict'"
	s = mustObjectStmt(t, in)
	if s.Kind != "INFERENCE JOB SERVICE" {
		t.Fatalf("Kind = %q", s.Kind)
	}
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|WITH|FROM|" {
		t.Fatalf("clauses = %q", got)
	}
	if s.Clauses[2].Query == nil {
		t.Error("FROM has no query")
	}
	if _, ok := optLitString(s.Clauses[1].Options, "SPECIFICATION"); !ok {
		t.Error("WITH has no SPECIFICATION")
	}
	mustObjectStmt(t, "EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM (WITH c AS (SELECT 1) SELECT * FROM c) MODEL = m")
}

func TestStartsExecuteJobService(t *testing.T) {
	for in, want := range map[string]bool{
		"JOB SERVICE":           true,
		"INFERENCE JOB SERVICE": true,
		"JOB":                   false,
		"INFERENCE SERVICE":     false,
		"TASK t":                false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsExecuteJobService(); got != want {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}

func TestParseImageArtifactRepository(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OR REPLACE IMAGE REPOSITORY r ENCRYPTION = (TYPE = 'SNOWFLAKE_SSE') COMMENT = 'c' WITH TAG (t = 'v')")
	if s.Kind != "IMAGE REPOSITORY" || !s.OrReplace || len(s.Clauses) != 1 || len(s.Clauses[0].Options) != 2 {
		t.Fatalf("got %+v", s)
	}
	s = mustObjectStmt(t, "CREATE IMAGE REPOSITORY IF NOT EXISTS r")
	if !s.IfNotExists || len(s.Clauses) != 0 {
		t.Fatalf("got %+v", s)
	}
	s = mustObjectStmt(t, "CREATE ARTIFACT REPOSITORY r TYPE = PYPI API_INTEGRATION = 'i' PACKAGE_RETENTION = ON_OBJECT_CREATION")
	if w, ok := optWords(s.Clauses[0].Options, "TYPE"); s.Kind != "ARTIFACT REPOSITORY" || !ok || w != "PYPI" {
		t.Fatalf("got %+v", s)
	}
	for in, want := range map[string]string{
		"ALTER IMAGE REPOSITORY IF EXISTS r SET COMMENT = 'c'":    "SET",
		"ALTER IMAGE REPOSITORY r UNSET TAG t, u":                 "UNSET TAG",
		"ALTER ARTIFACT REPOSITORY r SET INDEX_URL = 'u'":         "SET",
		"ALTER ARTIFACT REPOSITORY r SET TAG t = 'v'":             "SET TAG",
		"ALTER ARTIFACT REPOSITORY r UNSET AUTHENTICATION_SECRET": "UNSET",
	} {
		s := mustObjectStmt(t, in)
		if got := strings.Join(clauseKeywords(s), "|"); got != want {
			t.Errorf("%q: clauses = %q, want %q", in, got, want)
		}
	}
	for in, want := range map[string]string{
		"DROP SERVICE IF EXISTS s":        "SERVICE",
		"DROP IMAGE REPOSITORY r":         "IMAGE REPOSITORY",
		"DROP ARTIFACT REPOSITORY db.s.r": "ARTIFACT REPOSITORY",
	} {
		d, ok := mustParseOne(t, in).(*ast.DropStmt)
		if !ok || d.Kind != ast.DropObject || d.Object != want {
			t.Errorf("%q: got %+v", in, d)
		}
	}
}

func TestParseServiceRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE SERVICE s FROM SPECIFICATION 'x'",                                     // no compute pool
		"CREATE SERVICE s IN COMPUTE POOL p",                                          // no specification
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION x",                     // text is a string
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_FILE = x",              // path is a string
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE 'x'",          // USING is required
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x' USING (a => 1)",    // no USING without a template
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE 'x' USING ()", // an argument is required
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE 'x' USING (a = 1)",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE 'x' USING (a => (1))",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE 'x' USING (a => 1 b => 2)",
		"CREATE OR REPLACE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x'", // not documented
		"ALTER SERVICE s",
		"ALTER SERVICE s RENAME TO t", // not documented
		"ALTER SERVICE s SUSPEND RESUME",
		"ALTER SERVICE s FROM SPECIFICATION_TEMPLATE 'x'",
		"EXECUTE JOB SERVICE IN COMPUTE POOL p",
		"EXECUTE JOB SERVICE IN COMPUTE POOL p FROM SPECIFICATION_FILE = 'x'", // needs a stage
		"EXECUTE JOB SERVICE p FROM SPECIFICATION 'x'",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM (SELECT 1)", // MODEL is required
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p FROM (SELECT 1) MODEL = m",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM t MODEL = m",
		"CREATE ARTIFACT REPOSITORY r",               // TYPE is required
		"CREATE ARTIFACT REPOSITORY r COMMENT = 'c'", // TYPE is required
		"ALTER IMAGE REPOSITORY r",
		"ALTER IMAGE REPOSITORY r RENAME TO s",
		"DROP IMAGE REPOSITORY r CASCADE",

		// Properties the pages do not document, or values of another shape.
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x' NOT_A_PROPERTY = 1",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x' MIN_INSTANCES = 'a'",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x' AUTO_RESUME = 'x'",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x' EXTERNAL_ACCESS_INTEGRATIONS = ('a')",
		"ALTER SERVICE s SET LOG_LEVEL = info",
		"ALTER SERVICE s SET MIN_INSTANCES = 1 TAG t = 'v'", // SET TAG is a form of its own
		"ALTER SERVICE s UNSET NOT_A_PROPERTY",
		"EXECUTE JOB SERVICE IN COMPUTE POOL p FROM SPECIFICATION 's' MIN_INSTANCES = 1",
		"EXECUTE JOB SERVICE IN COMPUTE POOL p FROM SPECIFICATION 's' REPLICAS = x",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM (SELECT 1) MODEL = m QUERY_WAREHOUSE = w",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM (SELECT 1) MODEL = 'm'",
		"CREATE IMAGE REPOSITORY r ENCRYPTION = (TYPE = 'OTHER')",
		"CREATE IMAGE REPOSITORY r ENCRYPTION = (TYPE = SNOWFLAKE_SSE)",
		"CREATE IMAGE REPOSITORY r ENCRYPTION = (TYPE = 'SNOWFLAKE_SSE', KMS_KEY_ID = 'k')",
		"CREATE IMAGE REPOSITORY r TYPE = PYPI",
		"ALTER IMAGE REPOSITORY r SET ENCRYPTION = (TYPE = 'SNOWFLAKE_SSE')",
		"ALTER IMAGE REPOSITORY r UNSET ENCRYPTION",
		"CREATE ARTIFACT REPOSITORY r TYPE = NPM",
		"CREATE ARTIFACT REPOSITORY r TYPE = PYPI PACKAGE_RETENTION = NEVER",
		"CREATE ARTIFACT REPOSITORY r TYPE = PYPI AUTHENTICATION_SECRET = 's'",
		"ALTER ARTIFACT REPOSITORY r SET TYPE = PYPI",
		"ALTER ARTIFACT REPOSITORY r SET API_INTEGRATION = 'i'",
		"ALTER ARTIFACT REPOSITORY r UNSET COMMENT, PACKAGE_RETENTION",
		"ALTER ARTIFACT REPOSITORY r UNSET INDEX_URL",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseShowServiceScopes(t *testing.T) {
	for in, want := range map[string]string{
		"SHOW NODES IN COMPUTE POOL p LIKE 'n%' LIMIT 1":  "COMPUTE POOL",
		"SHOW SERVICES IN COMPUTE POOL p":                 "COMPUTE POOL",
		"SHOW ENDPOINTS IN SERVICE db.s.svc":              "SERVICE",
		"SHOW SERVICE CONTAINERS IN SERVICE svc":          "SERVICE",
		"SHOW IMAGES IN IMAGE REPOSITORY r":               "IMAGE REPOSITORY",
		"SHOW PACKAGES LIKE 'x' IN ARTIFACT REPOSITORY r": "ARTIFACT REPOSITORY",
	} {
		s, ok := mustParseOne(t, in).(*ast.ShowStmt)
		if !ok || s.Scope != ast.ShowScopeOther || s.ScopeText != want || s.ScopeName == nil {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	s := mustParseOne(t, "SHOW NODES IN COMPUTE POOL p LIKE 'n%'").(*ast.ShowStmt)
	if !s.HasLike || s.Like != "n%" {
		t.Errorf("LIKE = %q", s.Like)
	}
	// A schema named SERVICE is still a bare schema scope.
	for _, in := range []string{"SHOW TABLES IN service", "SHOW TABLES IN service LIMIT 1"} {
		s := mustParseOne(t, in).(*ast.ShowStmt)
		if s.Scope != ast.ShowScopeSchema || s.ScopeName.String() != "service" {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	for _, in := range []string{
		"SHOW NODES IN COMPUTE POOL",
		"SHOW TABLES IN COMPUTE POOL p LIKE 'x'", // only SHOW NODES writes LIKE after its scope
		"SHOW NODES LIKE 'a' IN COMPUTE POOL p LIKE 'b'",
	} {
		mustRejectObject(t, in)
	}
}
