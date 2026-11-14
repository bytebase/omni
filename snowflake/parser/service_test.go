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
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|FROM SPECIFICATION_TEMPLATE|USING|" {
		t.Fatalf("clauses = %q", got)
	}
	if pool := s.Clauses[0].Names; len(pool) != 1 || pool[0].String() != "pool1" {
		t.Errorf("pool = %+v", pool)
	}
	if spec := s.Clauses[1].Options; len(spec) != 1 || spec[0].Name != "" || spec[0].Lit.Value != "spec: {{ x }}" {
		t.Errorf("specification = %+v", spec)
	}
	using := s.Clauses[2].Options
	if len(using) != 3 || using[0].Name != "x" || using[0].Lit.Value != "v" || using[1].Name != `"Mixed"` ||
		using[1].Lit.Value != "2" || using[2].Words != "TRUE" {
		t.Errorf("USING = %+v", using)
	}
	if props := s.Clauses[3]; len(props.Options) != 3 || len(props.Tags) != 1 {
		t.Errorf("properties = %+v", props)
	}

	// The guide's examples write properties before the specification.
	s = mustObjectStmt(t, "CREATE SERVICE s IN COMPUTE POOL p MIN_INSTANCES=1 MAX_INSTANCES=1 FROM SPECIFICATION_TEMPLATE $$x$$ USING (image_url=>'i', SERVER_PORT=>8000) COMMENT = 'c'")
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL||FROM SPECIFICATION_TEMPLATE|USING|" {
		t.Errorf("clauses = %q", got)
	}

	for in, want := range map[string]string{
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'spec'":                          "FROM SPECIFICATION/",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_FILE = 'a.yaml'":                 "FROM/SPECIFICATION_FILE",
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION_TEMPLATE_FILE = 'a' USING(a=>1)": "FROM/SPECIFICATION_TEMPLATE_FILE",
	} {
		s := mustObjectStmt(t, in)
		if c := s.Clauses[1]; len(c.Options) != 1 || c.Keywords+"/"+c.Options[0].Name != want {
			t.Errorf("%q: FROM = %+v", in, c)
		}
	}
}

func TestParseAlterService(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER SERVICE IF EXISTS s SUSPEND":                                 "SUSPEND",
		"ALTER SERVICE s RESUME":                                            "RESUME",
		"ALTER SERVICE s FROM SPECIFICATION 'x'":                            "FROM SPECIFICATION",
		"ALTER SERVICE s FROM SPECIFICATION_TEMPLATE 'x' USING (a => 1)":    "FROM SPECIFICATION_TEMPLATE|USING",
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
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|FROM SPECIFICATION|" {
		t.Fatalf("clauses = %q", got)
	}
	if w, ok := optWords(s.Clauses[2].Options, "NAME"); !ok || w != "DB.S.JOB" {
		t.Errorf("NAME = %q", w)
	}
	// The usage notes and examples write the properties before the specification.
	s = mustObjectStmt(t, "EXECUTE JOB SERVICE IN COMPUTE POOL p NAME = j ASYNC = TRUE FROM SPECIFICATION 'spec' REPLICAS = 2")
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL||FROM SPECIFICATION|" {
		t.Errorf("clauses = %q", got)
	}
	s = mustObjectStmt(t, "execute job service in compute pool p from specification_template 'x' using (a => 'b')")
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|FROM SPECIFICATION_TEMPLATE|USING" {
		t.Errorf("clauses = %q", got)
	}

	in = "EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 'spec' FROM (SELECT a FROM t) " +
		"MODEL = db.s.m VERSION = v1 FUNCTION = 'predict'"
	s = mustObjectStmt(t, in)
	if s.Kind != "INFERENCE JOB SERVICE" {
		t.Fatalf("Kind = %q", s.Kind)
	}
	if got := strings.Join(clauseKeywords(s), "|"); got != "IN COMPUTE POOL|WITH SPECIFICATION|FROM|" {
		t.Fatalf("clauses = %q", got)
	}
	if s.Clauses[2].Query == nil {
		t.Error("FROM has no query")
	}
	if o := s.Clauses[1].Options; len(o) != 1 || o[0].Lit == nil || o[0].Lit.Value != "spec" {
		t.Errorf("WITH SPECIFICATION = %+v", o)
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
		"CREATE SERVICE s IN COMPUTE POOL db.p FROM SPECIFICATION 'x'",                                  // the pool is a single name
		"EXECUTE JOB SERVICE IN COMPUTE POOL db.p FROM SPECIFICATION 'x'",                               // the pool is a single name
		"CREATE SERVICE s IN COMPUTE POOL p FROM SPECIFICATION 'x' MIN_INSTANCES = 1 MIN_INSTANCES = 2", // set twice
		"ALTER SERVICE s SET COMMENT = 'c' WITH TAG (t = 'v')",
		"ALTER SERVICE s SET QUERY_WAREHOUSE = db.wh",                                    // a warehouse is a single name
		"EXECUTE JOB SERVICE IN COMPUTE POOL p NAME = j FROM SPECIFICATION 's' NAME = k", // set twice
		"EXECUTE JOB SERVICE IN COMPUTE POOL p NAME = j",
		"CREATE SERVICE s IN COMPUTE POOL p MIN_INSTANCES = 1 FROM SPECIFICATION 'x' MIN_INSTANCES = 2", // set twice
		"CREATE SERVICE s IN COMPUTE POOL p TAG (t = 'v') FROM SPECIFICATION 'x' TAG (u = 'w')",
		"CREATE SERVICE s IN COMPUTE POOL p NOT_A_PROPERTY = 1 FROM SPECIFICATION 'x'",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM (SELECT 1) MODEL = m VERSION = a.b",
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

func TestParseServiceStageAndRestore(t *testing.T) {
	for in, want := range map[string]string{
		"CREATE SERVICE s IN COMPUTE POOL p FROM @tutorial_stage SPECIFICATION_FILE = 'echo_spec.yaml' MIN_INSTANCES = 2": "@tutorial_stage/SPECIFICATION_FILE",
		"CREATE SERVICE s IN COMPUTE POOL p FROM @db.s.st/dir SPECIFICATION_TEMPLATE_FILE = 't.yaml' USING (a => 1)":      "@db.s.st/dir/SPECIFICATION_TEMPLATE_FILE",
		"ALTER SERVICE s FROM @st SPECIFICATION_FILE = 'a.yaml'":                                                          "@st/SPECIFICATION_FILE",
		"EXECUTE JOB SERVICE IN COMPUTE POOL p FROM @st SPECIFICATION_FILE = 'a.yaml' NAME = j":                           "@st/SPECIFICATION_FILE",
	} {
		s := mustObjectStmt(t, in)
		var from *ast.ObjectClause
		for _, c := range s.Clauses {
			if c.Keywords == "FROM" {
				from = c
			}
		}
		if from == nil || len(from.Options) != 2 || from.Options[0].Words+"/"+from.Options[1].Name != want {
			t.Errorf("%q: FROM = %+v", in, from)
		}
	}

	in := "EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION $$x$$ FROM @db.s.input_stage/images/ MODEL = db.s.m ASYNC = FALSE"
	s := mustObjectStmt(t, in)
	if from := s.Clauses[2]; from.Keywords != "FROM" || from.Query != nil || len(from.Options) != 1 ||
		from.Options[0].Words != "@db.s.input_stage/images/" {
		t.Errorf("inference FROM = %+v", from)
	}

	s = mustObjectStmt(t, "ALTER SERVICE IF EXISTS s RESTORE VOLUME v INSTANCES 0, 2 FROM SNAPSHOT snap")
	if got := strings.Join(clauseKeywords(s), "|"); got != "RESTORE VOLUME|INSTANCES|FROM SNAPSHOT" {
		t.Fatalf("clauses = %q", got)
	}
	if ids := s.Clauses[1].Options; len(ids) != 2 || ids[0].Lit.Ival != 0 || ids[1].Lit.Ival != 2 {
		t.Errorf("INSTANCES = %+v", ids)
	}
	if s.Clauses[0].Names[0].String() != "v" || s.Clauses[2].Names[0].String() != "snap" {
		t.Errorf("names = %+v / %+v", s.Clauses[0].Names, s.Clauses[2].Names)
	}

	for _, in := range []string{
		"CREATE SERVICE s IN COMPUTE POOL p FROM @st",                                   // the file is required
		"CREATE SERVICE s IN COMPUTE POOL p FROM @ SPECIFICATION_FILE = 'a'",            // no stage
		"CREATE SERVICE s IN COMPUTE POOL p FROM @st SPECIFICATION 'x'",                 // inline text takes no stage
		"CREATE SERVICE s IN COMPUTE POOL p FROM @st SPECIFICATION_TEMPLATE_FILE = 'a'", // USING is required
		"CREATE SERVICE s IN COMPUTE POOL p FROM @st \"SPECIFICATION_FILE\" = 'a'",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM @ MODEL = m",
		"EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL p WITH SPECIFICATION 's' FROM @st",
		"ALTER SERVICE s RESTORE VOLUME v INSTANCES FROM SNAPSHOT snap",
		"ALTER SERVICE s RESTORE VOLUME v INSTANCES a FROM SNAPSHOT snap",
		"ALTER SERVICE s RESTORE VOLUME v INSTANCES 1, FROM SNAPSHOT snap",
		"ALTER SERVICE s RESTORE VOLUME v INSTANCES 1 FROM snap",
		"ALTER SERVICE s RESTORE VOLUME db.v INSTANCES 1 FROM SNAPSHOT snap",
		"ALTER SERVICE s RESTORE VOLUME v INSTANCES 1",
	} {
		mustRejectObject(t, in)
	}
}
