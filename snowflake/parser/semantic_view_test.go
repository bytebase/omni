package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func mustCreateTag(t *testing.T, input string) *ast.CreateTagStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateTagStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateTagStmt", input, node)
	}
	return stmt
}

func mustAlterTag(t *testing.T, input string) *ast.AlterTagStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterTagStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterTagStmt", input, node)
	}
	return stmt
}

func mustCreateSemanticView(t *testing.T, input string) *ast.CreateSemanticViewStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateSemanticViewStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateSemanticViewStmt", input, node)
	}
	return stmt
}

func mustAlterSemanticView(t *testing.T, input string) *ast.AlterSemanticViewStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.AlterSemanticViewStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterSemanticViewStmt", input, node)
	}
	return stmt
}

func mustCreateDataset(t *testing.T, input string) *ast.CreateDatasetStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CreateDatasetStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateDatasetStmt", input, node)
	}
	return stmt
}

// section returns the named section of a semantic view, or nil.
func section(stmt *ast.CreateSemanticViewStmt, keyword string) *ast.SemanticViewSection {
	for _, s := range stmt.Sections {
		if s.Keyword == keyword {
			return s
		}
	}
	return nil
}

// assertLoc asserts that a node's Loc spans the whole statement (Start at 0,
// End at len(input)) for a single-statement input with no trailing semicolon.
func assertFullLoc(t *testing.T, loc ast.Loc, input string) {
	t.Helper()
	if loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", loc.Start)
	}
	if loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d (len input)", loc.End, len(input))
	}
}

// ---------------------------------------------------------------------------
// CREATE TAG — docs (truth1):
//
//	CREATE [OR REPLACE] TAG [IF NOT EXISTS] <name>
//	  [ALLOWED_VALUES '<v>' [, ...]]
//	  [PROPAGATE = {...} [ON_CONFLICT = {...}]]
//	  [COMMENT = '<string>']
// ---------------------------------------------------------------------------

func TestParseCreateTag_Modifiers(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG cost_center")
		if stmt.Name.String() != "cost_center" {
			t.Errorf("Name = %q, want cost_center", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.IfNotExists {
			t.Errorf("unexpected modifier: %+v", stmt)
		}
		if stmt.AllowedValues != nil || len(stmt.Options) != 0 {
			t.Errorf("expected no allowed values/options, got %+v / %+v", stmt.AllowedValues, stmt.Options)
		}
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE OR REPLACE TAG t")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG IF NOT EXISTS t")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("or replace if not exists is permissive", func(t *testing.T) {
		// Docs say OR REPLACE and IF NOT EXISTS are mutually exclusive, but the
		// engine parses permissively and defers the check to the semantic layer.
		stmt := mustCreateTag(t, "CREATE OR REPLACE TAG IF NOT EXISTS t")
		if !stmt.OrReplace || !stmt.IfNotExists {
			t.Errorf("OrReplace=%v IfNotExists=%v, want both", stmt.OrReplace, stmt.IfNotExists)
		}
	})

	t.Run("qualified name", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG db.sch.t")
		if stmt.Name.String() != "db.sch.t" {
			t.Errorf("Name = %q, want db.sch.t", stmt.Name.String())
		}
	})

	t.Run("quoted name", func(t *testing.T) {
		stmt := mustCreateTag(t, `CREATE TAG "My Tag"`)
		if stmt.Name.String() != `"My Tag"` {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})

	t.Run("loc spans statement", func(t *testing.T) {
		input := "CREATE TAG cost_center"
		stmt := mustCreateTag(t, input)
		assertFullLoc(t, stmt.Loc, input)
	})
}

func TestParseCreateTag_AllowedValues(t *testing.T) {
	t.Run("single value", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG t ALLOWED_VALUES 'red'")
		if len(stmt.AllowedValues) != 1 || stmt.AllowedValues[0] != "red" {
			t.Errorf("AllowedValues = %v, want [red]", stmt.AllowedValues)
		}
	})

	t.Run("multiple values", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG t ALLOWED_VALUES 'finance', 'engineering', 'sales'")
		want := []string{"finance", "engineering", "sales"}
		if len(stmt.AllowedValues) != len(want) {
			t.Fatalf("AllowedValues = %v, want %v", stmt.AllowedValues, want)
		}
		for i, v := range want {
			if stmt.AllowedValues[i] != v {
				t.Errorf("AllowedValues[%d] = %q, want %q", i, stmt.AllowedValues[i], v)
			}
		}
	})

	t.Run("allowed values then comment", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG t ALLOWED_VALUES 'a', 'b' COMMENT = 'note'")
		if len(stmt.AllowedValues) != 2 {
			t.Fatalf("AllowedValues = %v", stmt.AllowedValues)
		}
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "note" {
			t.Errorf("COMMENT option wrong: %+v", stmt.Options)
		}
	})
}

func TestParseCreateTag_Options(t *testing.T) {
	t.Run("comment only", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG cost_center COMMENT = 'cost_center tag'")
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit.Value != "cost_center tag" {
			t.Errorf("COMMENT option wrong: %+v", stmt.Options)
		}
	})

	t.Run("propagate and on_conflict (docs/corpus, beyond legacy g4)", func(t *testing.T) {
		// example_02 from the create-tag corpus. PROPAGATE / ON_CONFLICT are not in
		// the legacy g4 create_tag rule (flagged divergence: docs win).
		stmt := mustCreateTag(t, "CREATE TAG my_tag ALLOWED_VALUES 'blue', 'red' PROPAGATE = ON_DEPENDENCY ON_CONFLICT = ALLOWED_VALUES_SEQUENCE")
		if len(stmt.AllowedValues) != 2 {
			t.Fatalf("AllowedValues = %v", stmt.AllowedValues)
		}
		prop := findOption(stmt.Options, "PROPAGATE")
		if prop == nil || prop.Words != "ON_DEPENDENCY" {
			t.Errorf("PROPAGATE option wrong: %+v", stmt.Options)
		}
		oc := findOption(stmt.Options, "ON_CONFLICT")
		if oc == nil || oc.Words != "ALLOWED_VALUES_SEQUENCE" {
			t.Errorf("ON_CONFLICT option wrong: %+v", stmt.Options)
		}
	})

	t.Run("propagate string on_conflict", func(t *testing.T) {
		stmt := mustCreateTag(t, "CREATE TAG t PROPAGATE = ON_DATA_MOVEMENT ON_CONFLICT = 'custom'")
		oc := findOption(stmt.Options, "ON_CONFLICT")
		if oc == nil || oc.Lit == nil || oc.Lit.Value != "custom" {
			t.Errorf("ON_CONFLICT string wrong: %+v", stmt.Options)
		}
	})
}

func TestParseCreateTag_Negatives(t *testing.T) {
	// Missing name.
	assertParseError(t, "CREATE TAG")
	// ALLOWED_VALUES with no value.
	assertParseError(t, "CREATE TAG t ALLOWED_VALUES")
	// ALLOWED_VALUES trailing comma with no value.
	assertParseError(t, "CREATE TAG t ALLOWED_VALUES 'a',")
}

// ---------------------------------------------------------------------------
// ALTER TAG — docs (truth1):
//
//	RENAME TO <new_name>
//	{ADD|DROP} ALLOWED_VALUES '<v>' [, ...]
//	SET [ALLOWED_VALUES ...] [PROPAGATE ...] [COMMENT ...]
//	UNSET {ALLOWED_VALUES | PROPAGATE | ON_CONFLICT | COMMENT | DCM PROJECT}
//	SET MASKING POLICY <p> [, MASKING POLICY <p2> ...] [FORCE]
//	UNSET MASKING POLICY <p> [, MASKING POLICY <p2> ...]
// ---------------------------------------------------------------------------

func TestParseAlterTag_Rename(t *testing.T) {
	stmt := mustAlterTag(t, "ALTER TAG t RENAME TO t2")
	if stmt.Action != ast.AlterTagRename {
		t.Errorf("Action = %v, want Rename", stmt.Action)
	}
	if stmt.NewName.String() != "t2" {
		t.Errorf("NewName = %q, want t2", stmt.NewName.String())
	}
}

func TestParseAlterTag_IfExists(t *testing.T) {
	stmt := mustAlterTag(t, "ALTER TAG IF EXISTS t RENAME TO t2")
	if !stmt.IfExists {
		t.Error("IfExists not set")
	}
}

func TestParseAlterTag_AllowedValues(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t ADD ALLOWED_VALUES 'x', 'y'")
		if stmt.Action != ast.AlterTagAddAllowedValues {
			t.Errorf("Action = %v, want AddAllowedValues", stmt.Action)
		}
		if len(stmt.AllowedValues) != 2 {
			t.Errorf("AllowedValues = %v", stmt.AllowedValues)
		}
	})

	t.Run("drop", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t DROP ALLOWED_VALUES 'x'")
		if stmt.Action != ast.AlterTagDropAllowedValues {
			t.Errorf("Action = %v, want DropAllowedValues", stmt.Action)
		}
	})
}

func TestParseAlterTag_Set(t *testing.T) {
	t.Run("set allowed values", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t SET ALLOWED_VALUES 'a', 'b'")
		if stmt.Action != ast.AlterTagSet {
			t.Errorf("Action = %v, want Set", stmt.Action)
		}
		if len(stmt.AllowedValues) != 2 {
			t.Errorf("AllowedValues = %v", stmt.AllowedValues)
		}
	})

	t.Run("set comment", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t SET COMMENT = 'note'")
		if stmt.Action != ast.AlterTagSet {
			t.Errorf("Action = %v, want Set", stmt.Action)
		}
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit.Value != "note" {
			t.Errorf("COMMENT wrong: %+v", stmt.Options)
		}
	})

	t.Run("set propagate and comment", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t SET PROPAGATE = ON_DEPENDENCY COMMENT = 'n'")
		if findOption(stmt.Options, "PROPAGATE") == nil {
			t.Errorf("PROPAGATE missing: %+v", stmt.Options)
		}
		if findOption(stmt.Options, "COMMENT") == nil {
			t.Errorf("COMMENT missing: %+v", stmt.Options)
		}
	})
}

func TestParseAlterTag_Unset(t *testing.T) {
	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t UNSET COMMENT")
		if stmt.Action != ast.AlterTagUnset {
			t.Errorf("Action = %v, want Unset", stmt.Action)
		}
		if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "COMMENT" {
			t.Errorf("UnsetProps = %v, want [COMMENT]", stmt.UnsetProps)
		}
	})

	t.Run("unset allowed values", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t UNSET ALLOWED_VALUES")
		if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "ALLOWED_VALUES" {
			t.Errorf("UnsetProps = %v, want [ALLOWED_VALUES]", stmt.UnsetProps)
		}
	})

	t.Run("unset dcm project (two words)", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t UNSET DCM PROJECT")
		if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "DCM PROJECT" {
			t.Errorf("UnsetProps = %v, want [DCM PROJECT]", stmt.UnsetProps)
		}
	})
}

func TestParseAlterTag_MaskingPolicy(t *testing.T) {
	t.Run("set single", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t SET MASKING POLICY p")
		if stmt.Action != ast.AlterTagSetMaskingPolicy {
			t.Errorf("Action = %v, want SetMaskingPolicy", stmt.Action)
		}
		if len(stmt.MaskingPolicies) != 1 || stmt.MaskingPolicies[0].String() != "p" {
			t.Errorf("MaskingPolicies = %v", stmt.MaskingPolicies)
		}
		if stmt.Force {
			t.Error("Force should be false")
		}
	})

	t.Run("set multiple", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t SET MASKING POLICY p1, MASKING POLICY p2")
		if len(stmt.MaskingPolicies) != 2 {
			t.Errorf("MaskingPolicies = %v, want 2", stmt.MaskingPolicies)
		}
	})

	t.Run("set with force", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t SET MASKING POLICY p FORCE")
		if !stmt.Force {
			t.Error("Force not set")
		}
	})

	t.Run("unset single", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t UNSET MASKING POLICY p")
		if stmt.Action != ast.AlterTagUnsetMaskingPolicy {
			t.Errorf("Action = %v, want UnsetMaskingPolicy", stmt.Action)
		}
		if len(stmt.MaskingPolicies) != 1 {
			t.Errorf("MaskingPolicies = %v", stmt.MaskingPolicies)
		}
	})

	t.Run("unset multiple", func(t *testing.T) {
		stmt := mustAlterTag(t, "ALTER TAG t UNSET MASKING POLICY p1, MASKING POLICY p2")
		if len(stmt.MaskingPolicies) != 2 {
			t.Errorf("MaskingPolicies = %v, want 2", stmt.MaskingPolicies)
		}
	})
}

func TestParseAlterTag_Negatives(t *testing.T) {
	// Bare ALTER TAG with no name.
	assertParseError(t, "ALTER TAG")
	// No action.
	assertParseError(t, "ALTER TAG t")
	// RENAME without TO.
	assertParseError(t, "ALTER TAG t RENAME t2")
	// ADD without ALLOWED_VALUES.
	assertParseError(t, "ALTER TAG t ADD 'x'")
	// SET with nothing settable.
	assertParseError(t, "ALTER TAG t SET")
	// SET MASKING without POLICY.
	assertParseError(t, "ALTER TAG t SET MASKING p")
}

// ---------------------------------------------------------------------------
// CREATE SEMANTIC VIEW — docs (truth1):
//
//	CREATE [OR REPLACE] SEMANTIC VIEW [IF NOT EXISTS] <name>
//	  TABLES (...) [RELATIONSHIPS (...)] [FACTS (...)] [DIMENSIONS (...)]
//	  [METRICS (...)] [COMMENT = '...'] [AI_SQL_GENERATION '...']
//	  [AI_QUESTION_CATEGORIZATION '...'] [AI_VERIFIED_QUERIES (...)]
//	  [[WITH] TAG (...)] [COPY GRANTS]
// ---------------------------------------------------------------------------

func TestParseCreateSemanticView_Minimal(t *testing.T) {
	t.Run("tables and metrics", func(t *testing.T) {
		input := "CREATE SEMANTIC VIEW sv TABLES (orders) METRICS (orders.total AS SUM(amount))"
		stmt := mustCreateSemanticView(t, input)
		if stmt.Name.String() != "sv" {
			t.Errorf("Name = %q, want sv", stmt.Name.String())
		}
		tbl := section(stmt, "TABLES")
		if tbl == nil || tbl.Body != "orders" {
			t.Errorf("TABLES body = %q, want orders", bodyOf(tbl))
		}
		m := section(stmt, "METRICS")
		if m == nil || m.Body != "orders.total AS SUM(amount)" {
			t.Errorf("METRICS body = %q", bodyOf(m))
		}
		assertFullLoc(t, stmt.Loc, input)
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE OR REPLACE SEMANTIC VIEW sv TABLES (t) DIMENSIONS (t.d AS d)")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW IF NOT EXISTS sv TABLES (t) DIMENSIONS (t.d AS d)")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("qualified name", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW db.sch.sv TABLES (t) METRICS (t.m AS COUNT(*))")
		if stmt.Name.String() != "db.sch.sv" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})
}

func TestParseCreateSemanticView_AllSections(t *testing.T) {
	// Every documented section, with nested parens inside section bodies to
	// exercise the balanced-group reader.
	input := `CREATE SEMANTIC VIEW sv
		TABLES (
			o AS orders PRIMARY KEY (id),
			c AS customers PRIMARY KEY (id)
		)
		RELATIONSHIPS (
			orders_to_customers AS o (cust_id) REFERENCES c (id)
		)
		FACTS (
			o.amount AS amount
		)
		DIMENSIONS (
			c.name AS name,
			o.date AS order_date
		)
		METRICS (
			o.total AS SUM(o.amount)
		)`
	stmt := mustCreateSemanticView(t, input)
	for _, kw := range []string{"TABLES", "RELATIONSHIPS", "FACTS", "DIMENSIONS", "METRICS"} {
		if section(stmt, kw) == nil {
			t.Errorf("missing section %s", kw)
		}
	}
	// Balanced-group reader must capture nested parens whole.
	tbl := section(stmt, "TABLES")
	if !strings.Contains(tbl.Body, "PRIMARY KEY (id)") {
		t.Errorf("TABLES body did not capture nested parens: %q", tbl.Body)
	}
	rel := section(stmt, "RELATIONSHIPS")
	if !strings.Contains(rel.Body, "REFERENCES c (id)") {
		t.Errorf("RELATIONSHIPS body wrong: %q", rel.Body)
	}
}

func TestParseCreateSemanticView_TrailingClauses(t *testing.T) {
	t.Run("comment", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) COMMENT = 'a view'")
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit == nil || c.Lit.Value != "a view" {
			t.Errorf("COMMENT wrong: %+v", stmt.Options)
		}
	})

	t.Run("ai_sql_generation bare string", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) AI_SQL_GENERATION 'use the metrics'")
		ai := findOption(stmt.Options, "AI_SQL_GENERATION")
		if ai == nil || ai.Lit == nil || ai.Lit.Value != "use the metrics" {
			t.Errorf("AI_SQL_GENERATION wrong: %+v", stmt.Options)
		}
	})

	t.Run("ai_question_categorization bare string", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) AI_QUESTION_CATEGORIZATION 'categorize'")
		ai := findOption(stmt.Options, "AI_QUESTION_CATEGORIZATION")
		if ai == nil || ai.Lit == nil || ai.Lit.Value != "categorize" {
			t.Errorf("AI_QUESTION_CATEGORIZATION wrong: %+v", stmt.Options)
		}
	})

	t.Run("ai_verified_queries section", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) AI_VERIFIED_QUERIES (q AS 'SELECT 1')")
		if section(stmt, "AI_VERIFIED_QUERIES") == nil {
			t.Errorf("missing AI_VERIFIED_QUERIES section: %+v", stmt.Sections)
		}
	})

	t.Run("ai_verified_queries after scalar ai options (doc order)", func(t *testing.T) {
		// Per the docs AI_VERIFIED_QUERIES follows the scalar AI_SQL_GENERATION /
		// AI_QUESTION_CATEGORIZATION options; the order-tolerant body loop must
		// still capture it as a section.
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) COMMENT = 'c' AI_SQL_GENERATION 'g' AI_VERIFIED_QUERIES (q AS 'SELECT 1')")
		if section(stmt, "AI_VERIFIED_QUERIES") == nil {
			t.Errorf("missing AI_VERIFIED_QUERIES section: %+v", stmt.Sections)
		}
		if findOption(stmt.Options, "AI_SQL_GENERATION") == nil {
			t.Errorf("missing AI_SQL_GENERATION option: %+v", stmt.Options)
		}
		if findOption(stmt.Options, "COMMENT") == nil {
			t.Errorf("missing COMMENT option: %+v", stmt.Options)
		}
	})

	t.Run("with tag", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) WITH TAG (cost = 'high')")
		if len(stmt.Tags) != 1 || stmt.Tags[0].Name.String() != "cost" || stmt.Tags[0].Value != "high" {
			t.Errorf("Tags wrong: %+v", stmt.Tags)
		}
	})

	t.Run("bare tag", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) TAG (cost = 'high')")
		if len(stmt.Tags) != 1 {
			t.Errorf("Tags wrong: %+v", stmt.Tags)
		}
	})

	t.Run("copy grants", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) COPY GRANTS")
		if !stmt.CopyGrants {
			t.Error("CopyGrants not set")
		}
	})

	t.Run("comment then tag then copy grants", func(t *testing.T) {
		stmt := mustCreateSemanticView(t, "CREATE SEMANTIC VIEW sv TABLES (t) METRICS (t.m AS COUNT(*)) COMMENT = 'c' WITH TAG (k = 'v') COPY GRANTS")
		if findOption(stmt.Options, "COMMENT") == nil {
			t.Error("COMMENT missing")
		}
		if len(stmt.Tags) != 1 {
			t.Error("Tags missing")
		}
		if !stmt.CopyGrants {
			t.Error("CopyGrants missing")
		}
	})
}

func TestParseCreateSemanticView_Negatives(t *testing.T) {
	// Missing VIEW keyword.
	assertParseError(t, "CREATE SEMANTIC sv TABLES (t)")
	// Missing name.
	assertParseError(t, "CREATE SEMANTIC VIEW TABLES (t)")
	// Unbalanced section parens.
	assertParseError(t, "CREATE SEMANTIC VIEW sv TABLES (t")
}

// ---------------------------------------------------------------------------
// ALTER SEMANTIC VIEW — docs (truth1):
//
//	RENAME TO <new_name>
//	SET COMMENT = '<string>'
//	UNSET COMMENT
//	SET TAG <tag> = '<value>' [, ...]
//	UNSET TAG <tag> [, ...]
// ---------------------------------------------------------------------------

func TestParseAlterSemanticView_Rename(t *testing.T) {
	stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv RENAME TO sv2")
	if stmt.Action != ast.AlterSemanticViewRename {
		t.Errorf("Action = %v, want Rename", stmt.Action)
	}
	if stmt.NewName.String() != "sv2" {
		t.Errorf("NewName = %q", stmt.NewName.String())
	}
}

func TestParseAlterSemanticView_IfExists(t *testing.T) {
	stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW IF EXISTS sv RENAME TO sv2")
	if !stmt.IfExists {
		t.Error("IfExists not set")
	}
}

func TestParseAlterSemanticView_SetUnsetComment(t *testing.T) {
	t.Run("set comment", func(t *testing.T) {
		stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv SET COMMENT = 'note'")
		if stmt.Action != ast.AlterSemanticViewSet {
			t.Errorf("Action = %v, want Set", stmt.Action)
		}
		c := findOption(stmt.Options, "COMMENT")
		if c == nil || c.Lit.Value != "note" {
			t.Errorf("COMMENT wrong: %+v", stmt.Options)
		}
	})

	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv UNSET COMMENT")
		if stmt.Action != ast.AlterSemanticViewUnset {
			t.Errorf("Action = %v, want Unset", stmt.Action)
		}
		if len(stmt.UnsetProps) != 1 || stmt.UnsetProps[0] != "COMMENT" {
			t.Errorf("UnsetProps = %v", stmt.UnsetProps)
		}
	})
}

func TestParseAlterSemanticView_Tags(t *testing.T) {
	t.Run("set tag", func(t *testing.T) {
		stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv SET TAG cost = 'high'")
		if stmt.Action != ast.AlterSemanticViewSetTag {
			t.Errorf("Action = %v, want SetTag", stmt.Action)
		}
		if len(stmt.Tags) != 1 || stmt.Tags[0].Value != "high" {
			t.Errorf("Tags = %+v", stmt.Tags)
		}
	})

	t.Run("set multiple tags", func(t *testing.T) {
		stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv SET TAG a = '1', b = '2'")
		if len(stmt.Tags) != 2 {
			t.Errorf("Tags = %+v, want 2", stmt.Tags)
		}
	})

	t.Run("unset tag", func(t *testing.T) {
		stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv UNSET TAG cost")
		if stmt.Action != ast.AlterSemanticViewUnsetTag {
			t.Errorf("Action = %v, want UnsetTag", stmt.Action)
		}
		if len(stmt.UnsetTags) != 1 || stmt.UnsetTags[0].String() != "cost" {
			t.Errorf("UnsetTags = %v", stmt.UnsetTags)
		}
	})

	t.Run("unset multiple tags", func(t *testing.T) {
		stmt := mustAlterSemanticView(t, "ALTER SEMANTIC VIEW sv UNSET TAG a, b")
		if len(stmt.UnsetTags) != 2 {
			t.Errorf("UnsetTags = %v, want 2", stmt.UnsetTags)
		}
	})
}

func TestParseAlterSemanticView_Negatives(t *testing.T) {
	// Missing VIEW keyword.
	assertParseError(t, "ALTER SEMANTIC sv RENAME TO sv2")
	// No action.
	assertParseError(t, "ALTER SEMANTIC VIEW sv")
	// SET with nothing settable.
	assertParseError(t, "ALTER SEMANTIC VIEW sv SET")
	// RENAME without TO.
	assertParseError(t, "ALTER SEMANTIC VIEW sv RENAME sv2")
}

// ---------------------------------------------------------------------------
// CREATE DATASET — docs (truth1) + legacy g4 (truth2):
//
//	CREATE [OR REPLACE] DATASET [IF NOT EXISTS] <name>
//
// (Docs render IF NOT EXISTS before DATASET; the post-keyword spelling is
// accepted here, matching the legacy grammar and the rest of the engine — a
// flagged divergence.)
// ---------------------------------------------------------------------------

func TestParseCreateDataset(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		input := "CREATE DATASET my_dataset"
		stmt := mustCreateDataset(t, input)
		if stmt.Name.String() != "my_dataset" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.IfNotExists {
			t.Errorf("unexpected modifier: %+v", stmt)
		}
		assertFullLoc(t, stmt.Loc, input)
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateDataset(t, "CREATE OR REPLACE DATASET my_dataset")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})

	t.Run("if not exists (post-keyword, legacy spelling)", func(t *testing.T) {
		stmt := mustCreateDataset(t, "CREATE DATASET IF NOT EXISTS my_dataset")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("qualified name", func(t *testing.T) {
		stmt := mustCreateDataset(t, "CREATE DATASET db.sch.ds")
		if stmt.Name.String() != "db.sch.ds" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})
}

func TestParseCreateDataset_Negatives(t *testing.T) {
	// Missing name.
	assertParseError(t, "CREATE DATASET")
}

// ---------------------------------------------------------------------------
// Official docs corpus — every CREATE TAG statement in the create-tag corpus
// must parse with zero errors. The official docs are the authoritative oracle
// (truth1). There is no semantic-view / dataset corpus; those forms are covered
// from the docs above.
//
// example_03.sql uses `CREATE OR ALTER TAG`, a preview feature whose OR ALTER
// prefix the shared parseCreateStmt parser does not yet recognize (owned by
// parser.go / parseCreateStmt — out of this node's writes-scope). Such
// statements are skipped here and tracked as a flagged divergence (see
// orAlterLimited).
// ---------------------------------------------------------------------------

func TestCreateTag_OfficialCorpus(t *testing.T) {
	const dir = "testdata/official/create-tag"
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
			assertTagStatementsParse(t, string(data))
		})
	}
}

// assertTagStatementsParse parses sql and asserts every CREATE / ALTER TAG
// statement parses with no errors and to the expected AST type. OR ALTER preview
// statements are skipped (see orAlterLimited).
func assertTagStatementsParse(t *testing.T, sql string) {
	t.Helper()
	for _, seg := range Split(sql) {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)
		var want func(ast.Node) bool
		switch {
		case strings.HasPrefix(upper, "CREATE"):
			want = func(n ast.Node) bool { _, ok := n.(*ast.CreateTagStmt); return ok }
		case strings.HasPrefix(upper, "ALTER"):
			want = func(n ast.Node) bool { _, ok := n.(*ast.AlterTagStmt); return ok }
		default:
			continue // context statement owned by another DAG node
		}
		if orAlterLimited(upper) {
			// Known dependency limitation: must currently fail to parse. If it
			// starts parsing, the OR ALTER gap was closed — surface that so the
			// filter can be removed.
			if _, errs := parseSingle(seg.Text, seg.ByteStart); len(errs) == 0 {
				t.Logf("note: OR ALTER statement now parses, drop it from orAlterLimited: %q", text)
			}
			continue
		}
		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if len(errs) > 0 {
			t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
			continue
		}
		if !want(node) {
			t.Errorf("statement %q parsed to unexpected type %T", text, node)
		}
	}
}

// bodyOf returns the section body or a sentinel for nil (test diagnostics only).
func bodyOf(s *ast.SemanticViewSection) string {
	if s == nil {
		return "<nil>"
	}
	return s.Body
}
