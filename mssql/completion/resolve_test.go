package completion

import (
	"context"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/mssql/parser"
)

func TestResolveTokensToKeywords(t *testing.T) {
	// Build a CandidateSet with some token candidates.
	cs := &parser.CandidateSet{
		Tokens: []int{},
	}
	// Use Collect on "SELECT " to get real token candidates.
	cs = parser.Collect("SELECT ", 7)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}

	result := resolve(cs, nil, "SELECT ", 7)
	if len(result) == 0 {
		t.Fatal("expected non-empty result from token resolution")
	}

	// All results with nil catalog should be keywords (or types from type_name rule).
	for _, c := range result {
		if c.Type != CandidateKeyword && c.Type != CandidateType_ {
			t.Errorf("expected keyword or type candidate with nil catalog, got type %d for %q", c.Type, c.Text)
		}
	}
}

func TestResolveNilCandidateSet(t *testing.T) {
	result := resolve(nil, nil, "", 0)
	if result != nil {
		t.Errorf("expected nil for nil CandidateSet, got %v", result)
	}
}

func TestResolveTypeNameRule(t *testing.T) {
	// Directly test resolveRule for type_name.
	result := resolveRule("type_name", nil, "", 0)
	if len(result) == 0 {
		t.Fatal("expected type candidates for type_name rule")
	}

	// All should be CandidateType_.
	for _, c := range result {
		if c.Type != CandidateType_ {
			t.Errorf("expected CandidateType_ for type_name rule, got %d for %q", c.Type, c.Text)
		}
	}

	// Check some expected types are present.
	typeSet := make(map[string]bool)
	for _, c := range result {
		typeSet[c.Text] = true
	}

	expected := []string{
		"INT", "BIGINT", "SMALLINT", "TINYINT",
		"VARCHAR", "NVARCHAR", "CHAR", "NCHAR",
		"TEXT", "NTEXT",
		"DATETIME", "DATETIME2", "DATE", "TIME",
		"DECIMAL", "NUMERIC", "FLOAT", "REAL",
		"BIT",
		"MONEY", "SMALLMONEY",
		"UNIQUEIDENTIFIER",
		"XML",
		"VARBINARY", "IMAGE",
		"SQL_VARIANT",
	}
	for _, e := range expected {
		if !typeSet[e] {
			t.Errorf("expected type %q in type_name candidates", e)
		}
	}

	if len(result) != len(expected) {
		t.Errorf("expected %d type candidates, got %d", len(expected), len(result))
	}
}

func TestResolveRuleNilCatalog(t *testing.T) {
	// Rules that require catalog should return nil when catalog is nil.
	rules := []string{
		"table_ref", "columnref", "schema_ref", "func_name",
		"proc_ref", "index_ref", "trigger_ref", "database_ref",
		"sequence_ref", "login_ref", "user_ref", "role_ref",
	}
	for _, rule := range rules {
		result := resolveRule(rule, nil, "", 0)
		if result != nil {
			t.Errorf("expected nil for rule %q with nil catalog, got %v", rule, result)
		}
	}
}

func TestResolveUnknownRule(t *testing.T) {
	result := resolveRule("unknown_rule", nil, "", 0)
	if result != nil {
		t.Errorf("expected nil for unknown rule, got %v", result)
	}
}

func TestResolveDedupKeywords(t *testing.T) {
	// Create a CandidateSet with duplicate tokens.
	cs := &parser.CandidateSet{
		Tokens: []int{},
	}
	// Collect from empty input — should get top-level keywords, all deduped.
	cs = parser.Collect("", 0)
	result := resolve(cs, nil, "", 0)

	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]bool)
	for _, c := range result {
		k := key{c.Text, c.Type}
		if seen[k] {
			t.Errorf("duplicate candidate after resolve: %q (type %d)", c.Text, c.Type)
		}
		seen[k] = true
	}
}

func TestResolveTypeNameWithCandidateSet(t *testing.T) {
	// Simulate a CandidateSet that includes a type_name rule.
	cs := &parser.CandidateSet{
		Tokens: []int{},
		Rules:  []parser.RuleCandidate{{Rule: "type_name"}},
	}
	result := resolve(cs, nil, "", 0)

	// Should contain type candidates.
	hasType := false
	for _, c := range result {
		if c.Type == CandidateType_ {
			hasType = true
			break
		}
	}
	if !hasType {
		t.Error("expected type candidates when type_name rule is in CandidateSet")
	}
}

func TestCompleteWithCatalog(t *testing.T) {
	cat := catalog.New()
	err := cat.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: "Company",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{
						Name: "Employees",
						Columns: []*metadata.ColumnMetadata{
							{Name: "Id", Type: "int"},
							{Name: "Name", Type: "varchar(100)"},
						},
					},
				},
				Views: []*metadata.ViewMetadata{
					{Name: "ActiveEmployees"},
				},
				Sequences: []*metadata.SequenceMetadata{
					{Name: "EmpSeq"},
				},
				Procedures: []*metadata.ProcedureMetadata{
					{Name: "SyncEmployees"},
				},
			},
			{
				Name: "hr",
				Tables: []*metadata.TableMetadata{
					{
						Name: "Salaries",
						Columns: []*metadata.ColumnMetadata{
							{Name: "EmpId", Type: "int"},
							{Name: "Amount", Type: "money"},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = cat.LoadMetadataInto(context.Background(), &metadata.DatabaseSchemaMetadata{
		Name: "School",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "dbo",
				Tables: []*metadata.TableMetadata{
					{
						Name: "Student",
						Columns: []*metadata.ColumnMetadata{
							{Name: "Id", Type: "int"},
							{Name: "ParentName", Type: "varchar(100)"},
						},
					},
				},
			},
		},
	}, "School")
	if err != nil {
		t.Fatal(err)
	}

	hasCandidate := func(cs []Candidate, text string, typ CandidateType) bool {
		for _, c := range cs {
			if c.Text == text && c.Type == typ {
				return true
			}
		}
		return false
	}

	// Unqualified FROM
	fromCands := Complete("SELECT * FROM ", len("SELECT * FROM "), cat)
	if !hasCandidate(fromCands, "Employees", CandidateTable) {
		t.Fatalf("missing Employees table in %#v", fromCands)
	}
	if !hasCandidate(fromCands, "ActiveEmployees", CandidateView) {
		t.Fatalf("missing ActiveEmployees view in %#v", fromCands)
	}
	if !hasCandidate(fromCands, "hr", CandidateSchema) {
		t.Fatalf("missing hr schema in %#v", fromCands)
	}

	// Schema-qualified FROM
	hrCands := Complete("SELECT * FROM hr.", len("SELECT * FROM hr."), cat)
	if !hasCandidate(hrCands, "Salaries", CandidateTable) {
		t.Fatalf("missing hr.Salaries in %#v", hrCands)
	}
	if hasCandidate(hrCands, "Employees", CandidateTable) {
		t.Fatalf("unexpected dbo.Employees in hr. completion: %#v", hrCands)
	}

	// Cross-database table completion
	schoolCands := Complete("SELECT * FROM School.dbo.", len("SELECT * FROM School.dbo."), cat)
	if !hasCandidate(schoolCands, "Student", CandidateTable) {
		t.Fatalf("missing School.dbo.Student in %#v", schoolCands)
	}

	// Cross-database column completion
	schoolCols := Complete("SELECT School.dbo.Student.", len("SELECT School.dbo.Student."), cat)
	if !hasCandidate(schoolCols, "ParentName", CandidateColumn) {
		t.Fatalf("missing ParentName in %#v", schoolCols)
	}

	// Alias-qualified column completion
	aliasCols := Complete("SELECT e. FROM Employees AS e", len("SELECT e."), cat)
	if !hasCandidate(aliasCols, "Id", CandidateColumn) || !hasCandidate(aliasCols, "Name", CandidateColumn) {
		t.Fatalf("missing alias columns in %#v", aliasCols)
	}

	// CTE column completion
	cteCols := Complete("WITH cte AS (SELECT Id, Name FROM Employees) SELECT cte. FROM cte", len("WITH cte AS (SELECT Id, Name FROM Employees) SELECT cte."), cat)
	if !hasCandidate(cteCols, "Id", CandidateColumn) || !hasCandidate(cteCols, "Name", CandidateColumn) {
		t.Fatalf("missing CTE columns in %#v", cteCols)
	}

	// Procedure completion
	procCands := Complete("EXEC ", len("EXEC "), cat)
	if !hasCandidate(procCands, "SyncEmployees", CandidateProcedure) {
		t.Fatalf("missing SyncEmployees procedure in %#v", procCands)
	}

	// Built-in function completion
	selectCands := Complete("SELECT  FROM Employees", len("SELECT "), cat)
	if !hasCandidate(selectCands, "GETDATE", CandidateFunction) {
		t.Fatalf("missing GETDATE built-in function in %#v", selectCands)
	}

	// Pseudocolumn $action in MERGE OUTPUT
	mergeOut := Complete("MERGE INTO Employees AS d USING Employees AS s ON d.Id = s.Id WHEN MATCHED THEN UPDATE SET d.Name = s.Name OUTPUT ", len("MERGE INTO Employees AS d USING Employees AS s ON d.Id = s.Id WHEN MATCHED THEN UPDATE SET d.Name = s.Name OUTPUT "), cat)
	if !hasCandidate(mergeOut, "$action", CandidateColumn) {
		t.Fatalf("missing $action pseudo-column in %#v", mergeOut)
	}

	// Canonical built-in functions (including aggregates & STRING_AGG)
	for _, prefix := range []string{"COU", "STRING_"} {
		sql := "SELECT " + prefix
		cands := Complete(sql, len(sql), cat)
		want := "COUNT"
		if prefix == "STRING_" {
			want = "STRING_AGG"
		}
		if !hasCandidate(cands, want, CandidateFunction) {
			t.Fatalf("missing %s function for prefix %q in %#v", want, prefix, cands)
		}
	}
	for _, invalid := range []string{"MATH_SIGN", "STRINGAGG"} {
		sql := "SELECT " + invalid
		if hasCandidate(Complete(sql, len(sql), cat), invalid, CandidateFunction) {
			t.Fatalf("unexpected non-standard function %q", invalid)
		}
	}

	// Contextual keyword GO as identifier must not trigger batch isolation
	goIdent := Complete("SELECT GO +  FROM Employees", len("SELECT GO + "), cat)
	if !hasCandidate(goIdent, "Employees", CandidateTable) && !hasCandidate(goIdent, "Id", CandidateColumn) {
		t.Fatalf("expected Employees columns in expression after GO identifier, got %#v", goIdent)
	}

	// Following statements after cursor are excluded from statement scope
	isolated := Complete("SELECT e.; SELECT Id FROM Secret AS e", len("SELECT e."), cat)
	if hasCandidate(isolated, "Id", CandidateColumn) {
		t.Fatalf("following statement leaked columns across semicolon: %#v", isolated)
	}

	// ALTER TABLE column completion scopes only to the altered table
	alterCols := Complete("ALTER TABLE Employees ADD CONSTRAINT pk PRIMARY KEY (", len("ALTER TABLE Employees ADD CONSTRAINT pk PRIMARY KEY ("), cat)
	if !hasCandidate(alterCols, "Id", CandidateColumn) || !hasCandidate(alterCols, "Name", CandidateColumn) {
		t.Fatalf("missing Employees columns in ALTER TABLE PK: %#v", alterCols)
	}
	if hasCandidate(alterCols, "DepartmentId", CandidateColumn) || hasCandidate(alterCols, "EmployeeId", CandidateColumn) {
		t.Fatalf("unexpected other-table columns in ALTER TABLE PK: %#v", alterCols)
	}
}
