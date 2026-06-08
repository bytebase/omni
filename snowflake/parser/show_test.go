package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// SHOW — generic object classes
// ---------------------------------------------------------------------------

func TestParseShow_Basic(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantClass string
		wantTerse bool
	}{
		{"databases", "SHOW DATABASES", "DATABASES", false},
		{"terse databases", "SHOW TERSE DATABASES", "DATABASES", true},
		{"tables", "SHOW TABLES", "TABLES", false},
		{"views", "SHOW VIEWS", "VIEWS", false},
		{"materialized views", "SHOW MATERIALIZED VIEWS", "MATERIALIZED VIEWS", false},
		{"masking policies", "SHOW MASKING POLICIES", "MASKING POLICIES", false},
		{"row access policies", "SHOW ROW ACCESS POLICIES", "ROW ACCESS POLICIES", false},
		{"external tables", "SHOW EXTERNAL TABLES", "EXTERNAL TABLES", false},
		{"user functions", "SHOW USER FUNCTIONS", "USER FUNCTIONS", false},
		{"warehouses", "SHOW WAREHOUSES", "WAREHOUSES", false},
		{"streams", "SHOW STREAMS", "STREAMS", false},
		{"tasks", "SHOW TASKS", "TASKS", false},
		// Open-ended class: a class the lexer emits as a bare identifier (not a
		// keyword) must still parse — proving no fixed class vocabulary.
		{"unknown future class", "SHOW WIDGETS", "WIDGETS", false},
		{"unknown two-word class", "SHOW CORTEX SERVICES", "CORTEX SERVICES", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := mustParseOne(t, tt.input)
			stmt, ok := node.(*ast.ShowStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.ShowStmt", node)
			}
			if stmt.ObjectClass != tt.wantClass {
				t.Errorf("ObjectClass = %q, want %q", stmt.ObjectClass, tt.wantClass)
			}
			if stmt.Terse != tt.wantTerse {
				t.Errorf("Terse = %v, want %v", stmt.Terse, tt.wantTerse)
			}
		})
	}
}

func TestParseShow_Options(t *testing.T) {
	t.Run("history", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW DATABASES HISTORY")
		if !stmt.History {
			t.Errorf("History = false, want true")
		}
	})

	t.Run("like", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TABLES LIKE '%PART%'")
		if !stmt.HasLike || stringContent(stmt.Like) != "%PART%" {
			t.Errorf("Like = %q (has=%v), want '%%PART%%'", stmt.Like, stmt.HasLike)
		}
	})

	t.Run("in database scope", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TABLES IN DATABASE mydb")
		if stmt.Scope != ast.ShowScopeDatabase {
			t.Errorf("Scope = %v, want Database", stmt.Scope)
		}
		if stmt.ScopeName == nil || stmt.ScopeName.String() != "mydb" {
			t.Errorf("ScopeName = %v, want mydb", stmt.ScopeName)
		}
	})

	t.Run("in account scope", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TABLES IN ACCOUNT")
		if stmt.Scope != ast.ShowScopeAccount {
			t.Errorf("Scope = %v, want Account", stmt.Scope)
		}
		if stmt.ScopeName != nil {
			t.Errorf("ScopeName = %v, want nil", stmt.ScopeName)
		}
	})

	t.Run("in bare schema scope", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TABLES IN tpch_sf1")
		if stmt.Scope != ast.ShowScopeSchema {
			t.Errorf("Scope = %v, want Schema (bare name)", stmt.Scope)
		}
		if stmt.ScopeName == nil || stmt.ScopeName.String() != "tpch_sf1" {
			t.Errorf("ScopeName = %v, want tpch_sf1", stmt.ScopeName)
		}
	})

	t.Run("like then in", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TERSE TABLES LIKE '%PART%' IN tpch_sf1")
		if !stmt.Terse {
			t.Errorf("Terse = false, want true")
		}
		if !stmt.HasLike {
			t.Errorf("HasLike = false, want true")
		}
		if stmt.Scope != ast.ShowScopeSchema {
			t.Errorf("Scope = %v, want Schema", stmt.Scope)
		}
	})

	t.Run("starts with", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TERSE TABLES IN tpch_sf1 STARTS WITH 'LINE'")
		if !stmt.HasStarts || stringContent(stmt.StartsWith) != "LINE" {
			t.Errorf("StartsWith = %q (has=%v), want 'LINE'", stmt.StartsWith, stmt.HasStarts)
		}
	})

	t.Run("limit from", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TERSE TABLES IN tpch_sf1 LIMIT 3 FROM 'J'")
		if !stmt.HasLimit || stmt.Limit != "3" {
			t.Errorf("Limit = %q (has=%v), want 3", stmt.Limit, stmt.HasLimit)
		}
		if stringContent(stmt.LimitFrom) != "J" {
			t.Errorf("LimitFrom = %q, want 'J'", stmt.LimitFrom)
		}
	})

	t.Run("with privileges", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW DATABASES WITH PRIVILEGES USAGE, MODIFY")
		if len(stmt.Privileges) != 2 {
			t.Fatalf("Privileges len = %d, want 2", len(stmt.Privileges))
		}
		if stmt.Privileges[0].Name != "USAGE" || stmt.Privileges[1].Name != "MODIFY" {
			t.Errorf("Privileges = %v, want [USAGE MODIFY]", privNames(stmt.Privileges))
		}
	})

	t.Run("history then like", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW TABLES HISTORY LIKE 'test_show_tables_history'")
		if !stmt.History {
			t.Errorf("History = false, want true")
		}
		if !stmt.HasLike {
			t.Errorf("HasLike = false, want true")
		}
	})
}

// ---------------------------------------------------------------------------
// SHOW ... ->> <query>  (result pipe / FLOW)
// ---------------------------------------------------------------------------

func TestParseShow_Pipe(t *testing.T) {
	t.Run("show in account pipe select", func(t *testing.T) {
		// The official corpus pipes `... ->> SELECT ... FROM $1 ...`. The $1
		// result-set reference is not yet parseable by the expression/table-ref
		// parser (see TestParseShow_PipeDollarLimitation); here we verify the pipe
		// machinery itself with a table reference the parser supports.
		stmt := mustParseShow(t, `SHOW TABLES IN ACCOUNT ->> SELECT name FROM result_scan ORDER BY 1`)
		if stmt.Scope != ast.ShowScopeAccount {
			t.Errorf("Scope = %v, want Account", stmt.Scope)
		}
		if stmt.Pipe == nil {
			t.Fatalf("Pipe is nil, want a piped SELECT statement")
		}
		if _, ok := stmt.Pipe.(*ast.SelectStmt); !ok {
			t.Errorf("Pipe = %T, want *ast.SelectStmt", stmt.Pipe)
		}
	})
}

// TestParseShow_PipeDollarLimitation documents that a SHOW ... ->> SELECT ...
// FROM $1 result pipe cannot yet parse, because the table-reference / expression
// parser (T3/T5) does not implement the $N result-set reference. The SHOW pipe
// machinery is correct — it delegates the piped query to parseStmt — so when the
// table-ref parser gains $N support these official examples
// (show-tables/example_07..09) will parse with no change to show.go. Flagged
// divergence (shared root cause with the SET $var limitation).
func TestParseShow_PipeDollarLimitation(t *testing.T) {
	result := ParseBestEffort(`SHOW TABLES ->> SELECT * FROM $1 ORDER BY 1`)
	if len(result.Errors) == 0 {
		t.Skip("table-ref parser now supports $N — wire show-tables/example_07..09 into the corpus filter")
	}
}

// ---------------------------------------------------------------------------
// SHOW GRANTS
// ---------------------------------------------------------------------------

func TestParseShowGrants(t *testing.T) {
	t.Run("on database", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS ON DATABASE sales")
		if !stmt.IsGrants {
			t.Errorf("IsGrants = false, want true")
		}
		if stmt.GrantsOn == nil {
			t.Fatalf("GrantsOn is nil")
		}
		if stmt.GrantsOn.ObjectType != "DATABASE" || stmt.GrantsOn.Name.String() != "sales" {
			t.Errorf("GrantsOn = %+v, want DATABASE sales", stmt.GrantsOn)
		}
	})

	t.Run("on table", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS ON TABLE my_tbl")
		if stmt.GrantsOn == nil || stmt.GrantsOn.ObjectType != "TABLE" {
			t.Errorf("GrantsOn = %+v, want TABLE my_tbl", stmt.GrantsOn)
		}
	})

	t.Run("on account", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS ON ACCOUNT")
		if stmt.GrantsOn == nil || stmt.GrantsOn.Kind != ast.GrantTargetAccount {
			t.Errorf("GrantsOn = %+v, want ACCOUNT", stmt.GrantsOn)
		}
	})

	t.Run("to role", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS TO ROLE analyst")
		if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeRole {
			t.Fatalf("GrantsTo = %+v, want ROLE analyst", stmt.GrantsTo)
		}
		if stmt.GrantsTo.Name.String() != "analyst" {
			t.Errorf("GrantsTo.Name = %q, want analyst", stmt.GrantsTo.Name.String())
		}
	})

	t.Run("to user", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS TO USER user1")
		if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeUser {
			t.Errorf("GrantsTo = %+v, want USER user1", stmt.GrantsTo)
		}
	})

	t.Run("of role", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS OF ROLE analyst")
		if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeRole {
			t.Errorf("GrantsTo = %+v, want OF ROLE analyst", stmt.GrantsTo)
		}
	})

	t.Run("bare", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS")
		if !stmt.IsGrants {
			t.Errorf("IsGrants = false, want true")
		}
		if stmt.GrantsOn != nil || stmt.GrantsTo != nil {
			t.Errorf("bare SHOW GRANTS should have no On/To, got On=%v To=%v", stmt.GrantsOn, stmt.GrantsTo)
		}
	})

	t.Run("future in schema", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW FUTURE GRANTS IN SCHEMA sales.public")
		if !stmt.IsGrants || !stmt.Future {
			t.Errorf("IsGrants/Future = %v/%v, want true/true", stmt.IsGrants, stmt.Future)
		}
		if stmt.GrantsOn == nil || stmt.GrantsOn.Kind != ast.GrantTargetAllIn {
			t.Fatalf("GrantsOn = %+v, want container in schema", stmt.GrantsOn)
		}
		if stmt.GrantsOn.Container != ast.GrantContainerSchema {
			t.Errorf("Container = %v, want Schema", stmt.GrantsOn.Container)
		}
	})

	t.Run("future in database", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW FUTURE GRANTS IN DATABASE mydb")
		if stmt.GrantsOn == nil || stmt.GrantsOn.Container != ast.GrantContainerDatabase {
			t.Errorf("GrantsOn = %+v, want IN DATABASE mydb", stmt.GrantsOn)
		}
	})

	// Documented forms the legacy ANTLR grammar omits (docs = truth1):
	// SHOW GRANTS ... [LIMIT <rows>], SHOW FUTURE GRANTS TO ROLE / TO DATABASE ROLE.
	t.Run("grants to role limit", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS TO ROLE analyst LIMIT 5")
		if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeRole {
			t.Fatalf("GrantsTo = %+v, want ROLE analyst", stmt.GrantsTo)
		}
		if !stmt.HasLimit || stmt.Limit != "5" {
			t.Errorf("Limit = %q (has=%v), want 5", stmt.Limit, stmt.HasLimit)
		}
	})

	t.Run("grants on account limit", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW GRANTS ON ACCOUNT LIMIT 10")
		if !stmt.HasLimit || stmt.Limit != "10" {
			t.Errorf("Limit = %q (has=%v), want 10", stmt.Limit, stmt.HasLimit)
		}
	})

	t.Run("future grants to role", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW FUTURE GRANTS TO ROLE r1")
		if !stmt.IsGrants || !stmt.Future {
			t.Errorf("IsGrants/Future = %v/%v, want true/true", stmt.IsGrants, stmt.Future)
		}
		if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeRole {
			t.Fatalf("GrantsTo = %+v, want TO ROLE r1", stmt.GrantsTo)
		}
		if stmt.GrantsTo.Name.String() != "r1" {
			t.Errorf("GrantsTo.Name = %q, want r1", stmt.GrantsTo.Name.String())
		}
	})

	t.Run("future grants to database role", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW FUTURE GRANTS TO DATABASE ROLE mydb.dr")
		if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeDatabaseRole {
			t.Fatalf("GrantsTo = %+v, want TO DATABASE ROLE", stmt.GrantsTo)
		}
	})

	t.Run("future grants in schema limit", func(t *testing.T) {
		stmt := mustParseShow(t, "SHOW FUTURE GRANTS IN SCHEMA s LIMIT 3")
		if stmt.GrantsOn == nil || stmt.GrantsOn.Container != ast.GrantContainerSchema {
			t.Errorf("GrantsOn = %+v, want IN SCHEMA", stmt.GrantsOn)
		}
		if !stmt.HasLimit || stmt.Limit != "3" {
			t.Errorf("Limit = %q (has=%v), want 3", stmt.Limit, stmt.HasLimit)
		}
	})
}

// TestParseShow_OtherScope covers the SHOW PARAMETERS-style IN/FOR scopes that
// are neither ACCOUNT/DATABASE/SCHEMA/TABLE/VIEW: SESSION, USER <name>,
// WAREHOUSE <name>, TASK <name>, APPLICATION <name>. These are captured as
// ShowScopeOther with the scope keyword(s) in ScopeText, instead of being
// mis-modeled as a bare schema name. (Legacy show_parameters enumerates these.)
func TestParseShow_OtherScope(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantText string
		wantName string // ScopeName.String(); "" if none
	}{
		{"for session", "SHOW PARAMETERS IN SESSION", "SESSION", ""},
		{"for user named", "SHOW PARAMETERS FOR USER u1", "USER", "u1"},
		{"in warehouse named", "SHOW PARAMETERS IN WAREHOUSE wh", "WAREHOUSE", "wh"},
		{"in task named", "SHOW PARAMETERS IN TASK t1", "TASK", "t1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := mustParseShow(t, tt.input)
			if stmt.Scope != ast.ShowScopeOther {
				t.Fatalf("Scope = %v, want Other", stmt.Scope)
			}
			if stmt.ScopeText != tt.wantText {
				t.Errorf("ScopeText = %q, want %q", stmt.ScopeText, tt.wantText)
			}
			if tt.wantName == "" {
				if stmt.ScopeName != nil {
					t.Errorf("ScopeName = %v, want nil", stmt.ScopeName)
				}
			} else if stmt.ScopeName == nil || stmt.ScopeName.String() != tt.wantName {
				t.Errorf("ScopeName = %v, want %q", stmt.ScopeName, tt.wantName)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

func TestParseDescribe(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantType  string
		wantName  string // ObjectName.String(); "" when literal
		wantShort bool
		wantTypeO string // TYPE = <kw>
		wantSig   []string
	}{
		{"table", "DESCRIBE TABLE obj", "TABLE", "obj", false, "", nil},
		{"desc short", "DESC TABLE obj", "TABLE", "obj", true, "", nil},
		{"table qualified", "DESCRIBE TABLE db.sch.obj", "TABLE", "db.sch.obj", false, "", nil},
		{"table type columns", "DESCRIBE TABLE obj TYPE = COLUMNS", "TABLE", "obj", false, "COLUMNS", nil},
		{"view", "DESCRIBE VIEW obj", "VIEW", "obj", false, "", nil},
		{"materialized view", "DESCRIBE MATERIALIZED VIEW obj", "MATERIALIZED VIEW", "obj", false, "", nil},
		{"external table", "DESCRIBE EXTERNAL TABLE obj", "EXTERNAL TABLE", "obj", false, "", nil},
		{"file format", "DESCRIBE FILE FORMAT obj", "FILE FORMAT", "obj", false, "", nil},
		{"masking policy", "DESCRIBE MASKING POLICY obj", "MASKING POLICY", "obj", false, "", nil},
		{"row access policy", "DESCRIBE ROW ACCESS POLICY obj", "ROW ACCESS POLICY", "obj", false, "", nil},
		{"function no args", "DESCRIBE FUNCTION obj()", "FUNCTION", "obj", false, "", []string{}},
		{"function one arg", "DESCRIBE FUNCTION obj(INT)", "FUNCTION", "obj", false, "", []string{"INT"}},
		{"function two args", "DESCRIBE FUNCTION obj(INT, STRING)", "FUNCTION", "obj", false, "", []string{"INT", "STRING"}},
		{"procedure", "DESCRIBE PROCEDURE obj(STRING, INT)", "PROCEDURE", "obj", false, "", []string{"STRING", "INT"}},
		{"schema", "DESCRIBE SCHEMA obj", "SCHEMA", "obj", false, "", nil},
		{"sequence", "DESCRIBE SEQUENCE obj", "SEQUENCE", "obj", false, "", nil},
		{"stream", "DESCRIBE STREAM obj", "STREAM", "obj", false, "", nil},
		{"stage", "DESCRIBE STAGE obj", "STAGE", "obj", false, "", nil},
		{"task", "DESCRIBE TASK obj", "TASK", "obj", false, "", nil},
		{"user", "DESCRIBE USER obj", "USER", "obj", false, "", nil},
		{"warehouse", "DESCRIBE WAREHOUSE obj", "WAREHOUSE", "obj", false, "", nil},
		{"database", "DESCRIBE DATABASE obj", "DATABASE", "obj", false, "", nil},
		// Open-ended type: an unknown future object type must still parse.
		{"unknown type", "DESCRIBE WIDGET obj", "WIDGET", "obj", false, "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := mustParseOne(t, tt.input)
			stmt, ok := node.(*ast.DescribeStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.DescribeStmt", node)
			}
			if stmt.ObjectType != tt.wantType {
				t.Errorf("ObjectType = %q, want %q", stmt.ObjectType, tt.wantType)
			}
			if stmt.Short != tt.wantShort {
				t.Errorf("Short = %v, want %v", stmt.Short, tt.wantShort)
			}
			if tt.wantName != "" {
				if stmt.Name == nil || stmt.Name.String() != tt.wantName {
					t.Errorf("Name = %v, want %q", stmt.Name, tt.wantName)
				}
			}
			if stmt.TypeOption != tt.wantTypeO {
				t.Errorf("TypeOption = %q, want %q", stmt.TypeOption, tt.wantTypeO)
			}
			if !eqStrings(sigNames(stmt.Signature), tt.wantSig) {
				t.Errorf("Signature = %v, want %v", sigNames(stmt.Signature), tt.wantSig)
			}
		})
	}
}

func TestParseDescribe_Special(t *testing.T) {
	t.Run("search optimization on", func(t *testing.T) {
		node := mustParseOne(t, "DESCRIBE SEARCH OPTIMIZATION ON obj")
		stmt := node.(*ast.DescribeStmt)
		if stmt.ObjectType != "SEARCH OPTIMIZATION" {
			t.Errorf("ObjectType = %q, want SEARCH OPTIMIZATION", stmt.ObjectType)
		}
		if stmt.Name == nil || stmt.Name.String() != "obj" {
			t.Errorf("Name = %v, want obj", stmt.Name)
		}
	})

	t.Run("result string literal", func(t *testing.T) {
		node := mustParseOne(t, "DESCRIBE RESULT '01a2b3c4'")
		stmt := node.(*ast.DescribeStmt)
		if stmt.ObjectType != "RESULT" {
			t.Errorf("ObjectType = %q, want RESULT", stmt.ObjectType)
		}
		if stmt.NameLiteral == nil || stmt.NameLiteral.Kind != ast.LitString {
			t.Errorf("NameLiteral = %+v, want a string literal", stmt.NameLiteral)
		}
	})

	t.Run("transaction number literal", func(t *testing.T) {
		node := mustParseOne(t, "DESCRIBE TRANSACTION 1")
		stmt := node.(*ast.DescribeStmt)
		if stmt.ObjectType != "TRANSACTION" {
			t.Errorf("ObjectType = %q, want TRANSACTION", stmt.ObjectType)
		}
		if stmt.NameLiteral == nil || stmt.NameLiteral.Kind != ast.LitInt {
			t.Errorf("NameLiteral = %+v, want an int literal", stmt.NameLiteral)
		}
	})
}

func TestParseShowDescribe_Errors(t *testing.T) {
	cases := []string{
		"SHOW",                            // no object class
		"SHOW TERSE",                      // TERSE then nothing
		"SHOW GRANTS ON",                  // ON with no target
		"SHOW GRANTS TO",                  // TO with no grantee
		"SHOW GRANTS TO ROLE",             // ROLE with no name
		"SHOW FUTURE GRANTS",              // FUTURE GRANTS without IN
		"SHOW FUTURE GRANTS IN",           // IN with no container
		"SHOW TABLES LIKE",                // LIKE with no pattern
		"SHOW TABLES STARTS WITH",         // STARTS WITH with no string
		"SHOW TABLES LIMIT",               // LIMIT with no number
		"DESCRIBE",                        // nothing
		"DESC",                            // nothing
		"DESCRIBE TABLE",                  // type with no name
		"DESCRIBE SEARCH OPTIMIZATION ON", // ON with no name
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			result := ParseBestEffort(c)
			if len(result.Errors) == 0 {
				t.Errorf("expected parse error for %q, got none (stmts=%d)", c, len(result.File.Stmts))
			}
		})
	}
}

// mustParseShow parses input and asserts it yields a single *ast.ShowStmt.
func mustParseShow(t *testing.T, input string) *ast.ShowStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.ShowStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.ShowStmt", node)
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Official + legacy documentation corpus — every SHOW / DESCRIBE / DESC / USE /
// SET / UNSET / COMMENT / TRUNCATE statement must parse with zero errors. The
// official docs corpus is the authoritative oracle (truth1); the legacy corpus
// is the regression baseline (truth2). Statements owned by other DAG nodes
// (CREATE / DROP / INSERT / SELECT context lines) are skipped.
// ---------------------------------------------------------------------------

// utilityCorpusDirs are official-docs corpora whose utility statements are all
// owned by this node. (Other statements in these files — e.g. the CREATE/DROP/
// INSERT/SELECT setup lines in truncate-table — are context and are skipped.)
var utilityCorpusDirs = []string{
	"testdata/official/show-databases",
	"testdata/official/show-tables",
	"testdata/official/show-grants",
	"testdata/official/set",
	"testdata/official/unset",
	"testdata/official/truncate-table",
}

// utilityLegacyFiles are legacy corpus files of utility statements owned by this
// node.
var utilityLegacyFiles = []string{
	"testdata/legacy/show.sql",
	"testdata/legacy/use.sql",
	"testdata/legacy/comment.sql",
	"testdata/legacy/describe.sql",
}

func TestUtility_OfficialCorpus(t *testing.T) {
	for _, dir := range utilityCorpusDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertUtilityStatementsParse(t, string(data))
			})
		}
	}
}

func TestUtility_LegacyCorpus(t *testing.T) {
	for _, path := range utilityLegacyFiles {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			assertUtilityStatementsParse(t, string(data))
		})
	}
}

// utilityDollarLimited matches the corpus statements that are owned by this node
// but cannot yet parse because the shared table-reference parser (T5) does not
// implement the $N result-set reference (FROM $1). They are excluded from the
// zero-error corpus assertion and tracked as a flagged divergence; their root
// cause is the dependency, not this node's grammar. When the table-ref parser
// gains $-support, these will parse with no change to this node.
//
// NOTE: the expression-position $var case (SET (min, max) = (50, 2 * $min),
// set/example_06) used to live here but now parses since the expression parser
// gained $-support — it has been removed from this filter.
func utilityDollarLimited(upper string) bool {
	// SHOW ... ->> SELECT ... FROM $1 ...   (show-tables/example_07..09)
	if strings.Contains(upper, "->>") && strings.Contains(upper, "$1") {
		return true
	}
	return false
}

// assertUtilityStatementsParse parses sql and asserts that every utility
// statement in it (SHOW / DESCRIBE / DESC / USE / SET / UNSET / COMMENT /
// TRUNCATE) parses with no errors and to the expected AST type. Statements
// owned by other DAG nodes are skipped, as are the $-limited statements above.
func assertUtilityStatementsParse(t *testing.T, sql string) {
	t.Helper()
	for _, seg := range Split(sql) {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)
		kind, want := utilityStmtKind(upper)
		if kind == "" {
			continue // context statement owned by another DAG node
		}
		if utilityDollarLimited(upper) {
			// Known dependency limitation; must currently fail to parse. If it
			// starts parsing, the limitation was lifted — surface that so the
			// filter can be removed.
			if _, errs := parseSingle(seg.Text, seg.ByteStart); len(errs) == 0 {
				t.Logf("note: $-limited statement now parses, drop it from utilityDollarLimited: %q", text)
			}
			continue
		}
		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if len(errs) > 0 {
			t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
			continue
		}
		if !want(node) {
			t.Errorf("statement %q (%s) parsed to unexpected type %T", text, kind, node)
		}
	}
}

// utilityStmtKind classifies an uppercased statement text by its leading
// keyword, returning the kind label and a predicate that checks the parsed node
// is the matching AST type. Returns ("", nil) for non-utility statements.
func utilityStmtKind(upper string) (string, func(ast.Node) bool) {
	switch {
	case strings.HasPrefix(upper, "SHOW"):
		return "SHOW", func(n ast.Node) bool { _, ok := n.(*ast.ShowStmt); return ok }
	case strings.HasPrefix(upper, "DESCRIBE"), hasWordPrefix(upper, "DESC"):
		return "DESCRIBE", func(n ast.Node) bool { _, ok := n.(*ast.DescribeStmt); return ok }
	case hasWordPrefix(upper, "USE"):
		return "USE", func(n ast.Node) bool { _, ok := n.(*ast.UseStmt); return ok }
	case strings.HasPrefix(upper, "UNSET"):
		return "UNSET", func(n ast.Node) bool { _, ok := n.(*ast.UnsetStmt); return ok }
	case hasWordPrefix(upper, "SET"):
		return "SET", func(n ast.Node) bool { _, ok := n.(*ast.SetStmt); return ok }
	case strings.HasPrefix(upper, "COMMENT"):
		return "COMMENT", func(n ast.Node) bool { _, ok := n.(*ast.CommentStmt); return ok }
	case strings.HasPrefix(upper, "TRUNCATE"):
		return "TRUNCATE", func(n ast.Node) bool { _, ok := n.(*ast.TruncateStmt); return ok }
	}
	return "", nil
}

// hasWordPrefix reports whether upper begins with the given keyword as a whole
// word (followed by a space, '(' , or end-of-string). This distinguishes "SET
// V1 = 1" / "SET (a) = (1)" from a longer token, and "DESC TABLE t" from a
// column named DESC... (not a concern here but keeps the match precise).
func hasWordPrefix(upper, kw string) bool {
	if !strings.HasPrefix(upper, kw) {
		return false
	}
	if len(upper) == len(kw) {
		return true
	}
	next := upper[len(kw)]
	return next == ' ' || next == '\t' || next == '\n' || next == '('
}
