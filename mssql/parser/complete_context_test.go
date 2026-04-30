package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestCollectCompletionUpdateTargetScopeUsesQualifiedDMLTarget(t *testing.T) {
	sql := "UPDATE School.dbo.Student SET "
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Candidates == nil || ctx.Scope == nil {
		t.Fatal("expected completion context")
	}
	if !ctx.Candidates.HasRule("columnref") {
		t.Fatal("expected columnref candidate")
	}
	if ctx.Scope.DMLTarget == nil {
		t.Fatal("expected DML target")
	}
	if got, want := *ctx.Scope.DMLTarget, (RangeReference{Kind: RangeReferenceDMLTarget, Database: "School", Schema: "dbo", Object: "Student"}); !sameReferenceName(got, want) || got.Kind != want.Kind {
		t.Fatalf("DML target = %+v, want %+v", got, want)
	}
	if got, want := refObjects(ctx.Scope.LocalReferences), []string{"Student"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
}

func TestCollectCompletionDeleteTargetScopeUsesQualifiedDMLTarget(t *testing.T) {
	sql := "DELETE FROM School.dbo.Student WHERE "
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil || ctx.Scope.DMLTarget == nil {
		t.Fatal("expected DML target")
	}
	if got, want := *ctx.Scope.DMLTarget, (RangeReference{Kind: RangeReferenceDMLTarget, Database: "School", Schema: "dbo", Object: "Student"}); !sameReferenceName(got, want) || got.Kind != want.Kind {
		t.Fatalf("DML target = %+v, want %+v", got, want)
	}
}

func TestCollectCompletionSetOperationUsesCursorArmScope(t *testing.T) {
	sql := "SELECT Id FROM Employees UNION SELECT EmployeeId FROM Address UNION SELECT  FROM MySchema.SalaryLevel"
	cursor := strings.Index(sql, " FROM MySchema.SalaryLevel")
	if cursor < 0 {
		t.Fatal("test SQL missing cursor marker")
	}
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := refObjects(ctx.Scope.LocalReferences), []string{"SalaryLevel"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if len(ctx.Scope.OuterReferences) != 0 {
		t.Fatalf("outer references = %v, want none", ctx.Scope.OuterReferences)
	}
}

func TestCollectCompletionValuesDerivedTableAliasColumns(t *testing.T) {
	sql := "SELECT v. FROM (VALUES (1, 'a')) AS v(Id, Label)"
	cursor := strings.Index(sql, " FROM")
	if cursor < 0 {
		t.Fatal("test SQL missing cursor marker")
	}
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := len(ctx.Scope.LocalReferences), 1; got != want {
		t.Fatalf("local reference count = %d, want %d", got, want)
	}
	ref := ctx.Scope.LocalReferences[0]
	if ref.Kind != RangeReferenceValues {
		t.Fatalf("reference kind = %v, want values", ref.Kind)
	}
	if ref.Alias != "v" {
		t.Fatalf("alias = %q, want v", ref.Alias)
	}
	if got, want := ref.Columns, []string{"Id", "Label"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("columns = %v, want %v", got, want)
	}
}

func TestCollectCompletionCTEDefinitionsAndReferenceScope(t *testing.T) {
	sql := "WITH cte(Id, Label) AS (SELECT Id, Name FROM dbo.Employees) SELECT  FROM cte"
	cursor := strings.Index(sql, " FROM cte")
	if cursor < 0 {
		t.Fatal("test SQL missing cursor marker")
	}
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := len(ctx.CTEs), 1; got != want {
		t.Fatalf("CTE count = %d, want %d", got, want)
	}
	if ctx.CTEs[0].Kind != RangeReferenceCTE || ctx.CTEs[0].Object != "cte" {
		t.Fatalf("CTE = %+v, want cte reference", ctx.CTEs[0])
	}
	if got, want := ctx.CTEs[0].Columns, []string{"Id", "Label"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("CTE columns = %v, want %v", got, want)
	}
	if got, want := len(ctx.Scope.LocalReferences), 1; got != want {
		t.Fatalf("local reference count = %d, want %d", got, want)
	}
	if ref := ctx.Scope.LocalReferences[0]; ref.Kind != RangeReferenceCTE || ref.Object != "cte" {
		t.Fatalf("local reference = %+v, want CTE cte", ref)
	}
}

func TestCollectCompletionNestedSelectOuterReferences(t *testing.T) {
	sql := "SELECT * FROM Employees e WHERE EXISTS (SELECT 1 FROM Address a WHERE a.EmployeeId = e.)"
	cursor := strings.Index(sql, "e.)") + len("e.")
	ctx := CollectCompletion(sql, cursor)
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if got, want := refObjects(ctx.Scope.LocalReferences), []string{"Address"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("local references = %v, want %v", got, want)
	}
	if got, want := len(ctx.Scope.OuterReferences), 1; got != want {
		t.Fatalf("outer reference levels = %d, want %d", got, want)
	}
	if got, want := refObjects(ctx.Scope.OuterReferences[0]), []string{"Employees"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("outer references = %v, want %v", got, want)
	}
}

func TestCollectCompletionMergeScope(t *testing.T) {
	sql := "MERGE INTO dbo.Target AS t USING School.dbo.Source AS s ON s.Id = t.Id WHEN MATCHED THEN UPDATE SET Name = "
	ctx := CollectCompletion(sql, len(sql))
	if ctx == nil || ctx.Scope == nil {
		t.Fatal("expected completion scope")
	}
	if ctx.Scope.MergeTarget == nil {
		t.Fatal("expected merge target")
	}
	if got, want := *ctx.Scope.MergeTarget, (RangeReference{Kind: RangeReferenceMergeTarget, Schema: "dbo", Object: "Target", Alias: "t"}); !sameReferenceName(got, want) || got.Kind != want.Kind {
		t.Fatalf("merge target = %+v, want %+v", got, want)
	}
	if ctx.Scope.MergeSource == nil {
		t.Fatal("expected merge source")
	}
	if got, want := *ctx.Scope.MergeSource, (RangeReference{Kind: RangeReferenceMergeSource, Database: "School", Schema: "dbo", Object: "Source", Alias: "s"}); !sameReferenceName(got, want) || got.Kind != want.Kind {
		t.Fatalf("merge source = %+v, want %+v", got, want)
	}
}

func TestCollectCompletionPrefixIntentAndQualifier(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		prefix    string
		kind      ObjectKind
		qualifier MultipartName
	}{
		{
			name:   "sequence prefix",
			input:  "SELECT NEXT VALUE FOR Employee|",
			prefix: "Employee",
			kind:   ObjectKindSequence,
		},
		{
			name:   "table schema qualifier",
			input:  "SELECT * FROM dbo.Us|",
			prefix: "Us",
			kind:   ObjectKindTable,
			qualifier: MultipartName{
				Schema: "dbo",
			},
		},
		{
			name:   "view schema qualifier",
			input:  "DROP VIEW MySchema.|",
			prefix: "",
			kind:   ObjectKindView,
			qualifier: MultipartName{
				Schema: "MySchema",
			},
		},
		{
			name:   "procedure parameter type",
			input:  "CREATE PROCEDURE p @Name | AS SELECT 1",
			prefix: "",
			kind:   ObjectKindType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := strings.Replace(tt.input, "|", "", 1)
			cursor := strings.Index(tt.input, "|")
			ctx := CollectCompletion(sql, cursor)
			if ctx == nil {
				t.Fatal("expected completion context")
			}
			if ctx.Prefix != tt.prefix {
				t.Fatalf("prefix = %q, want %q", ctx.Prefix, tt.prefix)
			}
			if ctx.Intent == nil {
				t.Fatal("expected completion intent")
			}
			if !hasObjectKind(ctx.Intent.ObjectKinds, tt.kind) {
				t.Fatalf("object kinds = %v, want %v", ctx.Intent.ObjectKinds, tt.kind)
			}
			if ctx.Intent.Qualifier != tt.qualifier {
				t.Fatalf("qualifier = %+v, want %+v", ctx.Intent.Qualifier, tt.qualifier)
			}
		})
	}
}

func refObjects(refs []RangeReference) []string {
	var out []string
	for _, ref := range refs {
		out = append(out, ref.Object)
	}
	return out
}

func hasObjectKind(kinds []ObjectKind, want ObjectKind) bool {
	for _, got := range kinds {
		if got == want {
			return true
		}
	}
	return false
}

func sameReferenceName(a, b RangeReference) bool {
	return a.Server == b.Server && a.Database == b.Database && a.Schema == b.Schema && a.Object == b.Object && a.Alias == b.Alias
}
