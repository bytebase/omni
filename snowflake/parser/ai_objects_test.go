package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// wantClauses checks an ObjectStmt's clause keywords.
func wantClauses(t *testing.T, in string, s *ast.ObjectStmt, want ...string) {
	t.Helper()
	if got := clauseKeywords(s); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("%q: clauses = %q, want %q", in, got, want)
	}
}

func TestParseCortexSearch(t *testing.T) {
	in := "CREATE OR REPLACE CORTEX SEARCH SERVICE db.s.svc ON body PRIMARY KEY (id, k) ATTRIBUTES region, kind " +
		"WAREHOUSE = wh TARGET_LAG = '1 hour' EMBEDDING_MODEL = 'm' REFRESH_MODE = INCREMENTAL AUTO_SUSPEND = 60 " +
		"AS SELECT body, region, kind FROM t"
	s := mustObjectStmt(t, in)
	if s.Kind != "CORTEX SEARCH SERVICE" || !s.OrReplace || s.Name.String() != "db.s.svc" {
		t.Fatalf("got %+v", s)
	}
	wantClauses(t, in, s, "ON", "PRIMARY KEY", "ATTRIBUTES", "", "AS")
	if len(s.Clauses[1].Options[0].Group) != 2 || len(s.Clauses[2].Names) != 2 || s.Clauses[4].Query == nil {
		t.Errorf("clauses = %+v", s.Clauses)
	}

	in = "CREATE CORTEX SEARCH SERVICE svc TEXT INDEXES a, b VECTOR INDEXES v1 (model = 'e5'), v2, v3 (query_model = 'q') " +
		"ATTRIBUTES a WAREHOUSE = wh TARGET_LAG = '1 day' AS (SELECT a, b FROM t)"
	s = mustObjectStmt(t, in)
	wantClauses(t, in, s, "TEXT INDEXES", "VECTOR INDEXES", "ATTRIBUTES", "", "AS")
	if v := s.Clauses[1]; len(v.Options) != 3 || v.Options[0].Name != "V1" || v.Options[0].Group[0].Name != "MODEL" ||
		!v.Options[1].Bare || v.Options[1].Group != nil || v.Options[2].Group[0].Name != "QUERY_MODEL" {
		t.Errorf("VECTOR INDEXES = %+v", v)
	}

	for in, want := range map[string]string{
		"ALTER CORTEX SEARCH SERVICE IF EXISTS s SUSPEND":                                                     "SUSPEND",
		"ALTER CORTEX SEARCH SERVICE s RESUME SERVING":                                                        "RESUME SERVING",
		"ALTER CORTEX SEARCH SERVICE s REFRESH":                                                               "REFRESH",
		"ALTER CORTEX SEARCH SERVICE s SET TARGET_LAG = '1 hour' PRIMARY KEY = (a, b) REQUEST_LOGGING = TRUE": "SET",
		"ALTER CORTEX SEARCH SERVICE s UNSET PRIMARY KEY":                                                     "UNSET PRIMARY KEY",
		"ALTER CORTEX SEARCH SERVICE s SET ATTRIBUTES (a, b)":                                                 "SET ATTRIBUTES",
		"ALTER CORTEX SEARCH SERVICE s UNSET ATTRIBUTES":                                                      "UNSET ATTRIBUTES",
		"ALTER CORTEX SEARCH SERVICE s ADD SCORING PROFILE IF NOT EXISTS p '{}'":                              "ADD SCORING PROFILE IF NOT EXISTS",
		"ALTER CORTEX SEARCH SERVICE s DROP SCORING PROFILE IF EXISTS p":                                      "DROP SCORING PROFILE IF EXISTS",
		"ALTER CORTEX SEARCH SERVICE s SET TAG t = 'v'":                                                       "SET TAG",
		"ALTER CORTEX SEARCH SERVICE s UNSET TAG t, u":                                                        "UNSET TAG",
	} {
		wantClauses(t, in, mustObjectStmt(t, in), want)
	}

	for _, in := range []string{
		"CREATE CORTEX SEARCH SERVICE s ON a ATTRIBUTES b WAREHOUSE = wh AS SELECT 1",                // TARGET_LAG is required
		"CREATE CORTEX SEARCH SERVICE s ON a WAREHOUSE = wh TARGET_LAG = '1 hour' AS SELECT 1",       // ATTRIBUTES is required
		"CREATE CORTEX SEARCH SERVICE s ON a ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = '1 hour'",      // AS is required
		"CREATE CORTEX SEARCH SERVICE s ON a ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = 1 AS SELECT 1", // a string
		"CREATE CORTEX SEARCH SERVICE s ON a ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = 'x' X = 1 AS SELECT 1",
		"CREATE CORTEX SEARCH SERVICE IF NOT EXISTS s TEXT INDEXES a VECTOR INDEXES v ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = 'x' AS SELECT 1",
		"CREATE CORTEX SEARCH SERVICE s TEXT INDEXES a VECTOR INDEXES v (dim = 1) ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = 'x' AS SELECT 1",
		"CREATE CORTEX SEARCH SERVICE s TEXT INDEXES a ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = 'x' AS SELECT 1",
		"ALTER CORTEX SEARCH SERVICE s SUSPEND BOTH",
		"ALTER CORTEX SEARCH SERVICE s SET PRIMARY KEY = ()",
		"ALTER CORTEX SEARCH SERVICE s SET REFRESH_MODE = FULL",   // CREATE's alone
		"ALTER CORTEX SEARCH SERVICE s SET EMBEDDING_MODEL = 'm'", // CREATE's alone
		"ALTER CORTEX SEARCH SERVICE s UNSET COMMENT",
		"ALTER CORTEX SEARCH SERVICE IF EXISTS s ADD SCORING PROFILE p '{}'",
		"ALTER CORTEX SEARCH SERVICE s ADD SCORING PROFILE p",
		"DROP CORTEX SEARCH SERVICE IF EXISTS s", // the page writes no IF EXISTS
		"ALTER CORTEX SEARCH SERVICE s SET COMMENT = 'a' COMMENT = 'b'",
		"CREATE CORTEX SEARCH SERVICE s ON t.a ATTRIBUTES b WAREHOUSE = wh TARGET_LAG = 'x' AS SELECT 1",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAgent(t *testing.T) {
	for in, want := range map[string][]string{
		"CREATE OR REPLACE AGENT a COPY GRANTS SECURE = TRUE COMMENT = 'c' FROM SPECIFICATION $$models: {}$$": {"COPY GRANTS", "", "FROM SPECIFICATION"},
		"CREATE AGENT IF NOT EXISTS a FROM SPECIFICATION 's'":                                                 {"FROM SPECIFICATION"},
		"ALTER AGENT a SET COMMENT = 'c' PROFILE = '{}'":                                                      {"SET"},
		"ALTER AGENT a MODIFY LIVE VERSION SET SPECIFICATION = 's'":                                           {"MODIFY LIVE VERSION SET"},
		"ALTER AGENT a SET TAG t = 'v'":                                                                       {"SET TAG"},
		"CREATE OR REPLACE EXTERNAL AGENT e WITH VERSION v1 COMMENT = 'c'":                                    {"WITH VERSION", ""},
		"ALTER EXTERNAL AGENT IF EXISTS e ADD VERSION v2":                                                     {"ADD VERSION"},
		"ALTER EXTERNAL AGENT e SET COMMENT = 'c'":                                                            {"SET"},
	} {
		wantClauses(t, in, mustObjectStmt(t, in), want...)
	}
	s := mustObjectStmt(t, "CREATE AGENT a FROM SPECIFICATION 's'")
	if spec := s.Clauses[0].Options[0]; spec.Name != "" || spec.Lit.Value != "s" {
		t.Errorf("FROM = %+v", spec)
	}
	for _, in := range []string{"DROP AGENT IF EXISTS a", "DROP EXTERNAL AGENT e"} {
		if d, ok := mustParseOne(t, in).(*ast.DropStmt); !ok || d.Kind != ast.DropObject {
			t.Errorf("%q: got %T", in, d)
		}
	}
	for _, in := range []string{
		"CREATE AGENT a",
		"CREATE AGENT a PROFILE = '{}' FROM @s/agents/a", // the stage form takes COMMENT alone
		"CREATE AGENT a FROM s/agents/a",
		"CREATE AGENT a FROM SPECIFICATION",
		"CREATE AGENT a FROM SPECIFICATION x",
		"CREATE AGENT a SECURE = 'yes' FROM SPECIFICATION 's'",
		"CREATE AGENT a MODEL = x FROM SPECIFICATION 's'",
		"ALTER AGENT IF EXISTS a SET COMMENT = 'c'",
		"ALTER AGENT a SET SPECIFICATION = 's'",
		"ALTER AGENT a MODIFY LIVE VERSION SET COMMENT = 's'",
		"ALTER AGENT a UNSET COMMENT",
		"ALTER AGENT a RENAME TO b",
		"CREATE EXTERNAL AGENT e WITH VERSION",
		"ALTER EXTERNAL AGENT e SET VERSION = v",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseModel(t *testing.T) {
	for in, want := range map[string][]string{
		"CREATE OR REPLACE MODEL m WITH VERSION v1 FROM MODEL db.s.src VERSION v0": {"WITH VERSION", "FROM MODEL", "VERSION"},
		"CREATE MODEL IF NOT EXISTS m FROM MODEL src":                              {"FROM MODEL"},
		"ALTER MODEL IF EXISTS m SET COMMENT = 'c' DEFAULT_VERSION = 'v1'":         {"SET"},
		"ALTER MODEL m SET TAG t = 'v'":                                            {"SET TAG"},
		"ALTER MODEL m UNSET TAG t":                                                {"UNSET TAG"},
		"ALTER MODEL m VERSION v1 SET ALIAS = best":                                {"VERSION", "SET"},
		"ALTER MODEL m VERSION v1 UNSET ALIAS":                                     {"VERSION", "UNSET"},
		"ALTER MODEL m RENAME TO n":                                                {"RENAME TO"},
		"ALTER MODEL m ADD VERSION v2 FROM MODEL src VERSION v0":                   {"ADD VERSION", "FROM MODEL", "VERSION"},
		"ALTER MODEL IF EXISTS m DROP VERSION v1":                                  {"DROP VERSION"},
		"ALTER MODEL m MODIFY VERSION v1 SET COMMENT = 'c' METADATA = '{}'":        {"MODIFY VERSION", "SET"},
		"ALTER MODEL MONITOR IF EXISTS mm SUSPEND":                                 {"SUSPEND"},
		"ALTER MODEL MONITOR mm SET REFRESH_INTERVAL = '1 day' WAREHOUSE = wh":     {"SET"},
		"ALTER MODEL MONITOR mm SET BASELINE = 'b'":                                {"SET"},
		"ALTER MODEL MONITOR mm UNSET BASELINE, COMMENT":                           {"UNSET"},
		"ALTER MODEL MONITOR mm ADD SEGMENT_COLUMN = 'c'":                          {"ADD"},
		"ALTER MODEL MONITOR mm DROP SEGMENT_COLUMN = 'c'":                         {"DROP"},
		"CREATE MODEL MONITOR mm WITH MODEL = db.s.m GATEWAY = g FUNCTION = 'predict' " +
			"WAREHOUSE = wh REFRESH_INTERVAL = '1 day' AGGREGATION_WINDOW = '1 day' ID_COLUMNS = ('id')": {"WITH"},
		"CREATE OR REPLACE MODEL MONITOR mm WITH MODEL = m VERSION = 'v1' FUNCTION = 'predict' SOURCE = src " +
			"WAREHOUSE = wh REFRESH_INTERVAL = '1 day' AGGREGATION_WINDOW = '1 day' TIMESTAMP_COLUMN = ts BASELINE = b": {"WITH"},
	} {
		s := mustObjectStmt(t, in)
		wantClauses(t, in, s, want...)
	}
	s := mustObjectStmt(t, "ALTER MODEL m VERSION v1 SET ALIAS = best")
	if v, c := s.Clauses[0], s.Clauses[1]; v.Names[0].String() != "v1" || c.Options[0].Words != "BEST" {
		t.Errorf("VERSION ... SET = %+v, %+v", v, c)
	}
	for in, kind := range map[string]string{
		"DROP MODEL MONITOR IF EXISTS mm": "MODEL MONITOR",
		"DROP MODEL db.s.m":               "MODEL",
		"DROP EXPERIMENT e":               "EXPERIMENT",
		"DROP CORTEX SEARCH SERVICE s":    "CORTEX SEARCH SERVICE",
	} {
		if d, ok := mustParseOne(t, in).(*ast.DropStmt); !ok || d.Object != kind {
			t.Errorf("%q: got %+v", in, d)
		}
	}
	for _, in := range []string{
		"CREATE MODEL m",
		"CREATE MODEL m FROM MODEL",
		"CREATE MODEL m FROM src",
		"CREATE MODEL m FROM @",
		"ALTER MODEL m SET ALIAS = a",
		"ALTER MODEL m VERSION v1 SET COMMENT = 'c'",
		"ALTER MODEL m VERSION v1 UNSET COMMENT",
		"ALTER MODEL m MODIFY VERSION v1 SET ALIAS = a",
		"ALTER MODEL m MODIFY VERSION v1 SET",
		"ALTER MODEL IF EXISTS m RENAME TO n",
		"ALTER MODEL m UNSET COMMENT",
		"ALTER MODEL m DROP VERSION",
		"DROP MODEL IF EXISTS m", // the page writes no IF EXISTS
		"ALTER MODEL m UNSET ALIAS",
		"ALTER MODEL m MODIFY VERSION v1 SET COMMENT = 'a' COMMENT = 'b'",
		"ALTER MODEL MONITOR db.s.mm SUSPEND", // a single <monitor_name>
		"CREATE MODEL MONITOR mm WITH MODEL = m VERSION = 'v' FUNCTION = 'f' SOURCE = db.s.t WAREHOUSE = wh " +
			"REFRESH_INTERVAL = 'x' AGGREGATION_WINDOW = 'y' TIMESTAMP_COLUMN = ts", // a single <source_name>
		"CREATE MODEL MONITOR mm WITH MODEL = m FUNCTION = 'f' WAREHOUSE = wh REFRESH_INTERVAL = 'x' AGGREGATION_WINDOW = 'y'", // SOURCE ... or GATEWAY
		"CREATE MODEL MONITOR mm WITH MODEL = m GATEWAY = g FUNCTION = f WAREHOUSE = wh REFRESH_INTERVAL = 'x' AGGREGATION_WINDOW = 'y'",
		"CREATE MODEL MONITOR mm MODEL = m GATEWAY = g FUNCTION = 'f' WAREHOUSE = wh REFRESH_INTERVAL = 'x' AGGREGATION_WINDOW = 'y'",
		"ALTER MODEL MONITOR mm SET BASELINE = b",
		"ALTER MODEL MONITOR mm UNSET WAREHOUSE",
		"ALTER MODEL MONITOR mm ADD SEGMENT_COLUMN = c",
		"ALTER MODEL MONITOR mm ADD COLUMN = 'c'",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseExperiment(t *testing.T) {
	for in, want := range map[string][]string{
		"CREATE OR REPLACE EXPERIMENT e":                         nil,
		"CREATE EXPERIMENT IF NOT EXISTS db.s.e":                 nil,
		"ALTER EXPERIMENT e ADD RUN r":                           {"ADD RUN"},
		"ALTER EXPERIMENT e COMMIT RUN r":                        {"COMMIT RUN"},
		"ALTER EXPERIMENT e COMMIT RUN r WITH STATUS = 'FAILED'": {"COMMIT RUN", "WITH"},
		"ALTER EXPERIMENT e MODIFY RUN r SET STATUS = 'PAUSED'":  {"MODIFY RUN", "SET"},
		"ALTER EXPERIMENT e DROP RUN r":                          {"DROP RUN"},
	} {
		wantClauses(t, in, mustObjectStmt(t, in), want...)
	}
	for _, in := range []string{
		"CREATE EXPERIMENT e COMMENT = 'c'",
		"ALTER EXPERIMENT IF EXISTS e ADD RUN r",
		"ALTER EXPERIMENT e ADD RUN",
		"ALTER EXPERIMENT e MODIFY RUN r",
		"ALTER EXPERIMENT e MODIFY RUN r SET STATUS = PAUSED",
		"ALTER EXPERIMENT e COMMIT RUN r WITH STATUS = 'BOGUS'",
		"ALTER EXPERIMENT e COMMIT RUN r WITH STATUS = 'RUNNING'", // MODIFY's
		"ALTER EXPERIMENT e MODIFY RUN r SET STATUS = 'FINISHED'", // COMMIT's
		"ALTER EXPERIMENT e COMMIT RUN r WITH COMMENT = 'c'",
		"ALTER EXPERIMENT e SET STATUS = 'x'",
		"DROP EXPERIMENT IF EXISTS e", // the page writes no IF EXISTS
	} {
		mustRejectObject(t, in)
	}
}

func TestParseOnlineFeatureTable(t *testing.T) {
	in := "CREATE OR REPLACE ONLINE FEATURE TABLE f PRIMARY KEY (id, k) TARGET_LAG = '1 minute' WAREHOUSE = wh " +
		"REFRESH_MODE = AUTO WITH COMMENT = 'c' WITH TAG (t = 'v') FROM src"
	s := mustObjectStmt(t, in)
	wantClauses(t, in, s, "PRIMARY KEY", "", "FROM")
	if c := s.Clauses[1]; len(c.Options) != 4 || len(c.Tags) != 1 {
		t.Errorf("properties = %+v", c)
	}
	for in, want := range map[string]string{
		"ALTER ONLINE FEATURE TABLE IF EXISTS f SUSPEND":        "SUSPEND",
		"ALTER ONLINE FEATURE TABLE f REFRESH":                  "REFRESH",
		"ALTER ONLINE FEATURE TABLE f RENAME TO g":              "RENAME TO",
		"ALTER ONLINE FEATURE TABLE f SET TARGET_LAG = '1 day'": "SET",
		"ALTER ONLINE FEATURE TABLE f UNSET TAG t":              "UNSET TAG",
	} {
		wantClauses(t, in, mustObjectStmt(t, in), want)
	}
	for _, in := range []string{
		"CREATE ONLINE FEATURE TABLE f PRIMARY KEY (id) WAREHOUSE = wh FROM src", // TARGET_LAG is required
		"CREATE ONLINE FEATURE TABLE f PRIMARY KEY (id) TARGET_LAG = '1 day' WAREHOUSE = wh",
		"CREATE ONLINE FEATURE TABLE f TARGET_LAG = '1 day' WAREHOUSE = wh FROM src",
		"CREATE ONLINE FEATURE TABLE IF NOT EXISTS f PRIMARY KEY (id) TARGET_LAG = '1 day' WAREHOUSE = wh FROM src",
		"CREATE ONLINE FEATURE TABLE f PRIMARY KEY (id) TARGET_LAG = '1 day' WAREHOUSE = wh REFRESH_MODE = LAZY FROM src",
		"CREATE ONLINE FEATURE TABLE f PRIMARY KEY (id) TARGET_LAG = '1 day' WAREHOUSE = wh WITH FROM src",
		"ALTER ONLINE FEATURE TABLE f UNSET COMMENT",
		"ALTER ONLINE FEATURE TABLE f SET LAG = '1 day'",
		"CREATE ONLINE FEATURE TABLE f PRIMARY KEY (id) TARGET_LAG = '1 day' WAREHOUSE = wh FROM db.s.src", // a single <source>
		"CREATE ONLINE FEATURE TABLE f PRIMARY KEY (id) TARGET_LAG = '1 day' WAREHOUSE = wh COMMENT = 'a' WITH COMMENT = 'b' FROM src",
		"ALTER ONLINE FEATURE TABLE f SET REFRESH_MODE = AUTO", // CREATE's alone
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAiFunction(t *testing.T) {
	in := "CREATE OR REPLACE AI FUNCTION db.s.f(body VARCHAR, n NUMBER(10, 2)) RETURNS VARCHAR COMMENT = 'c' AS $$ AI_COMPLETE('m', body) $$"
	node := mustParseOne(t, in)
	f, ok := node.(*ast.CreateAiFunctionStmt)
	if !ok || !f.OrReplace || f.Name.String() != "db.s.f" || len(f.Args) != 2 || f.ReturnType == nil ||
		len(f.Options) != 1 || !f.BodyDollar || f.Loc.End != len(in) {
		t.Fatalf("got %#v", node)
	}
	in = "CREATE AI FUNCTION IF NOT EXISTS f() RETURNS VARCHAR FROM EXPERIMENT db.s.e RUN iter_7"
	f, ok = mustParseOne(t, in).(*ast.CreateAiFunctionStmt)
	if !ok || !f.IfNotExists || f.Args != nil || f.Experiment.String() != "db.s.e" || f.Run.Name != "iter_7" || f.Loc.End != len(in) {
		t.Fatalf("got %#v", f)
	}
	for _, in := range []string{
		"CREATE AI FUNCTION f() RETURNS VARCHAR",
		"CREATE AI FUNCTION f RETURNS VARCHAR AS 's'",
		"CREATE AI FUNCTION f() AS 's'",
		"CREATE AI FUNCTION f(a) RETURNS VARCHAR AS 's'",
		"CREATE AI FUNCTION f(a INT DEFAULT 1) RETURNS VARCHAR AS 's'",
		"CREATE AI FUNCTION f() RETURNS VARCHAR LANGUAGE = SQL AS 's'",
		"CREATE AI FUNCTION f() RETURNS VARCHAR COMMENT = 'c' FROM EXPERIMENT e RUN r",
		"CREATE AI FUNCTION f() RETURNS VARCHAR FROM EXPERIMENT e",
		"CREATE OR REPLACE AI FUNCTION IF NOT EXISTS f() RETURNS VARCHAR AS 's'",
		"CREATE OR ALTER AI FUNCTION f() RETURNS VARCHAR AS 's'",
		"CREATE SECURE AI FUNCTION f() RETURNS VARCHAR AS 's'",
		"CREATE TEMP AI FUNCTION f() RETURNS VARCHAR AS 's'",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAlterDataset(t *testing.T) {
	in := "ALTER DATASET db.s.d ADD VERSION 'v1' FROM (SELECT seq4() AS id, uniform(1, 10, random()) AS part FROM TABLE(GENERATOR(ROWCOUNT => 10))) " +
		"PARTITION BY part COMMENT = 'c' METADATA = '{}'"
	d, ok := mustParseOne(t, in).(*ast.AlterDatasetStmt)
	if !ok || d.Drop || d.Name.String() != "db.s.d" || d.Version.Name != "v1" || !d.Version.Quoted ||
		d.Query == nil || d.PartitionBy == nil || len(d.Options) != 2 || d.Loc.End != len(in) {
		t.Fatalf("got %#v", d)
	}
	d, ok = mustParseOne(t, "ALTER DATASET IF EXISTS d DROP VERSION v1").(*ast.AlterDatasetStmt)
	if !ok || !d.Drop || !d.IfExists || d.Version.Name != "v1" || d.Version.Quoted {
		t.Fatalf("got %#v", d)
	}
	mustParseOne(t, "ALTER DATASET d ADD VERSION v1 FROM SELECT 1")
	// An unparenthesized query stops before the version's clauses.
	for in, n := range map[string]int{
		"ALTER DATASET d ADD VERSION v1 FROM SELECT 1 COMMENT = 'c'":                            1,
		"ALTER DATASET d ADD VERSION v1 FROM SELECT * FROM t PARTITION BY part METADATA = '{}'": 1,
		"ALTER DATASET d ADD VERSION v1 FROM SELECT a comment FROM t metadata = 'm'":            1,
		"ALTER DATASET d ADD VERSION v1 FROM SELECT 1 AS comment":                               0,
	} {
		d, ok := mustParseOne(t, in).(*ast.AlterDatasetStmt)
		if !ok || len(d.Options) != n || d.Loc.End != len(in) {
			t.Errorf("%q: got %#v", in, d)
		}
	}
	for _, in := range []string{
		"ALTER DATASET IF EXISTS d ADD VERSION v1 FROM SELECT 1",
		"ALTER DATASET d ADD VERSION v1",
		"ALTER DATASET d ADD VERSION v1 FROM (SELECT 1) METADATA = 1",
		"ALTER DATASET d ADD VERSION v1 FROM (SELECT 1) TAG = 'x'",
		"ALTER DATASET d DROP VERSION",
		"ALTER DATASET d RENAME TO e",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseCreateDatasetIfNotExistsFirst(t *testing.T) {
	for in, ifNotExists := range map[string]bool{
		"CREATE IF NOT EXISTS DATASET d":   true,
		"CREATE DATASET IF NOT EXISTS d":   true,
		"CREATE OR REPLACE DATASET db.s.d": false,
	} {
		d, ok := mustParseOne(t, in).(*ast.CreateDatasetStmt)
		if !ok || d.IfNotExists != ifNotExists || d.Loc.End != len(in) {
			t.Errorf("%q: got %#v", in, d)
		}
	}
	for _, in := range []string{
		"CREATE OR REPLACE IF NOT EXISTS DATASET d",
		"CREATE IF NOT EXISTS DATASET IF NOT EXISTS d",
		"CREATE IF NOT EXISTS TABLE t (a INT)",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseShowInModelScopes(t *testing.T) {
	for in, scope := range map[string]string{
		"SHOW VERSIONS LIKE 'v%' IN MODEL db.s.m":           "MODEL",
		"SHOW FUNCTIONS IN MODEL m":                         "MODEL",
		"SHOW VERSIONS IN DATASET d LIMIT 10":               "DATASET",
		"SHOW RUNS IN EXPERIMENT e":                         "EXPERIMENT",
		"SHOW RUN METRICS IN EXPERIMENT e LIMIT 5 FROM 'x'": "EXPERIMENT",
		"SHOW MODEL MONITORS IN GATEWAY g":                  "GATEWAY",
	} {
		s, ok := mustParseOne(t, in).(*ast.ShowStmt)
		if !ok || s.Scope != ast.ShowScopeOther || s.ScopeText != scope || s.ScopeName == nil {
			t.Errorf("%q: got %#v", in, s)
		}
	}
	// A scope word with no name after it is a schema's name.
	s, ok := mustParseOne(t, "SHOW TABLES IN model").(*ast.ShowStmt)
	if !ok || s.Scope != ast.ShowScopeSchema || s.ScopeName.String() != "model" {
		t.Errorf("SHOW TABLES IN model: got %#v", s)
	}
}
