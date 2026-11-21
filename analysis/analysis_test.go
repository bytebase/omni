package analysis_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestQueryTypeValuesAndString(t *testing.T) {
	// Numeric values are pinned to Bytebase's base.QueryType and span corpus
	// fixtures.
	cases := []struct {
		qt   analysis.QueryType
		num  int
		name string
	}{
		{analysis.QueryTypeUnknown, 0, "UNKNOWN"},
		{analysis.Select, 1, "SELECT"},
		{analysis.Explain, 2, "EXPLAIN"},
		{analysis.SelectInfoSchema, 3, "SELECT_INFO_SCHEMA"},
		{analysis.DDL, 4, "DDL"},
		{analysis.DML, 5, "DML"},
		{analysis.QueryType(999), 999, "UNKNOWN"},
	}
	for _, tc := range cases {
		if int(tc.qt) != tc.num {
			t.Errorf("QueryType %s = %d, want %d", tc.name, int(tc.qt), tc.num)
		}
		if got := tc.qt.String(); got != tc.name {
			t.Errorf("QueryType(%d).String() = %q, want %q", tc.num, got, tc.name)
		}
	}
}

func TestColumnResourceStringAndSet(t *testing.T) {
	cases := []struct {
		res  analysis.ColumnResource
		want string
	}{
		{
			res:  analysis.ColumnResource{},
			want: "",
		},
		{
			res:  analysis.ColumnResource{Database: "db", Table: "t", Column: "c"},
			want: "db.t.c",
		},
		{
			res:  analysis.ColumnResource{Database: "db", Schema: "public", Table: "t", Column: "c"},
			want: "db.public.t.c",
		},
		{
			res:  analysis.ColumnResource{Server: "srv", Database: "db", Schema: "dbo", Table: "t", Column: "c"},
			want: "srv.db.dbo.t.c",
		},
		{
			res:  analysis.ColumnResource{Instance: "inst1", Server: "link1", Database: "hr", Table: "emp", Column: "sal"},
			want: "link1.hr.emp.sal",
		},
	}
	for _, tc := range cases {
		if got := tc.res.String(); got != tc.want {
			t.Errorf("ColumnResource(%+v).String() = %q, want %q", tc.res, got, tc.want)
		}
	}

	set := analysis.ColumnSet{
		{Database: "db", Schema: "s", Table: "t", Column: "a"}: true,
	}
	if !set[analysis.ColumnResource{Database: "db", Schema: "s", Table: "t", Column: "a"}] {
		t.Fatal("expected ColumnResource key match in ColumnSet")
	}
	if set[analysis.ColumnResource{Instance: "other", Database: "db", Schema: "s", Table: "t", Column: "a"}] {
		t.Fatal("expected distinct Instance to be a distinct ColumnSet key")
	}
}

func TestResourceNotFoundError(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	cause := errors.New("boom")

	e1 := &analysis.ResourceNotFoundError{
		Database: strPtr("db1"),
		Schema:   strPtr("s1"),
		Table:    strPtr("t1"),
		Column:   strPtr("c1"),
		Err:      cause,
	}
	want1 := "resource not found: database: db1, schema: s1, table: t1, column: c1, err: boom"
	if got := e1.Error(); got != want1 {
		t.Errorf("Error() = %q, want %q", got, want1)
	}
	if !errors.Is(e1, cause) {
		t.Errorf("errors.Is(e1, cause) = false, want true")
	}

	e2 := &analysis.ResourceNotFoundError{
		Server:       strPtr("srv"),
		DatabaseLink: strPtr("lnk"),
		Database:     strPtr("db"),
		Function:     strPtr("fn"),
	}
	want2 := "resource not found: server: srv, database link: lnk, database: db, function: fn"
	if got := e2.Error(); got != want2 {
		t.Errorf("Error() = %q, want %q", got, want2)
	}
	if e2.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", e2.Unwrap())
	}
}

func TestTypeNotSupportedError(t *testing.T) {
	cause := errors.New("pivot")
	e := &analysis.TypeNotSupportedError{
		Err:   cause,
		Type:  "PIVOT",
		Name:  "sales",
		Extra: "dynamic",
	}
	want := "type not supported: PIVOT, err: pivot, name: sales, extra: dynamic"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if !errors.Is(e, cause) {
		t.Errorf("errors.Is(e, cause) = false, want true")
	}
}

func TestFunctionNotSupportedError(t *testing.T) {
	cause := errors.New("bad body")
	e := &analysis.FunctionNotSupportedError{
		Err:        cause,
		Function:   "my_fn",
		Definition: "CREATE FUNCTION my_fn()",
	}
	want := "function is not supported with data masking: my_fn, err: bad body, definition: CREATE FUNCTION my_fn()"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if !errors.Is(e, cause) {
		t.Errorf("errors.Is(e, cause) = false, want true")
	}
}

func TestUnresolvedColumnsError(t *testing.T) {
	e := &analysis.UnresolvedColumnsError{
		Relations: []analysis.ColumnResource{
			{Database: "db2", Schema: "public", Table: "z"},
			{Database: "db1", Schema: "public", Table: "a"},
			{Database: "db2", Schema: "public", Table: "m"},
			{Schema: "public", Table: "local"},
		},
	}
	wantMsg := "the synced schema describes no columns for db1.public.a, db2.public.m, db2.public.z, public.local"
	if got := e.Error(); got != wantMsg {
		t.Errorf("Error() = %q, want %q", got, wantMsg)
	}
	wantDBs := []string{"db1", "db2"}
	if got := e.Databases(); !reflect.DeepEqual(got, wantDBs) {
		t.Errorf("Databases() = %v, want %v", got, wantDBs)
	}
}

func TestSentinelErrors(t *testing.T) {
	if got := analysis.ErrMixUserSystemTables.Error(); got != "cannot access user and system tables at the same time" {
		t.Errorf("ErrMixUserSystemTables = %q", got)
	}
	if got := analysis.ErrExplainAnalyzeExecutes.Error(); got != "EXPLAIN ANALYZE executes the statement it plans, which an explain request never does: run it as a query, or drop ANALYZE to plan the statement" {
		t.Errorf("ErrExplainAnalyzeExecutes = %q", got)
	}
	if got := analysis.ErrNotExplainable.Error(); got != "EXPLAIN plans only a query or a data change" {
		t.Errorf("ErrNotExplainable = %q", got)
	}
}

func TestAliasAssignabilityContract(t *testing.T) {
	type (
		queryType                 = analysis.QueryType
		columnResource            = analysis.ColumnResource
		columnSet                 = analysis.ColumnSet
		querySpanResult           = analysis.QuerySpanResult
		resolvedSpan              = analysis.ResolvedSpan
		resourceNotFoundError     = analysis.ResourceNotFoundError
		typeNotSupportedError     = analysis.TypeNotSupportedError
		functionNotSupportedError = analysis.FunctionNotSupportedError
		unresolvedColumnsError    = analysis.UnresolvedColumnsError
		changedResource           = analysis.ChangedResource
		statementChanges          = analysis.StatementChanges
		byteRange                 = analysis.ByteRange
		statement                 = analysis.Statement
	)

	var qt analysis.QueryType = queryType(analysis.Select)
	if qt != analysis.Select {
		t.Fatalf("aliased QueryType = %v, want %v", qt, analysis.Select)
	}

	col := columnResource{Database: "db", Table: "t", Column: "c"}
	var cs analysis.ColumnSet = columnSet{col: true}
	var res analysis.QuerySpanResult = querySpanResult{
		Name:          "c",
		SourceColumns: cs,
		IsPlainField:  true,
	}
	table := "missing"
	var nf error = &resourceNotFoundError{Table: &table}
	var targetNF *analysis.ResourceNotFoundError
	if !errors.As(nf, &targetNF) || targetNF.Table == nil || *targetNF.Table != "missing" {
		t.Fatalf("errors.As failed for aliased ResourceNotFoundError: %v", nf)
	}

	var tns error = &typeNotSupportedError{Type: "PIVOT"}
	var targetTNS *analysis.TypeNotSupportedError
	if !errors.As(tns, &targetTNS) || targetTNS.Type != "PIVOT" {
		t.Fatalf("errors.As failed for aliased TypeNotSupportedError: %v", tns)
	}

	var fns error = &functionNotSupportedError{Function: "fn"}
	var targetFNS *analysis.FunctionNotSupportedError
	if !errors.As(fns, &targetFNS) || targetFNS.Function != "fn" {
		t.Fatalf("errors.As failed for aliased FunctionNotSupportedError: %v", fns)
	}

	var uc error = &unresolvedColumnsError{Relations: []columnResource{{Database: "db", Table: "t"}}}
	var targetUC *analysis.UnresolvedColumnsError
	if !errors.As(uc, &targetUC) || len(targetUC.Relations) != 1 {
		t.Fatalf("errors.As failed for aliased UnresolvedColumnsError: %v", uc)
	}

	var span *analysis.ResolvedSpan = &resolvedSpan{
		Type:                   qt,
		SourceColumns:          cs,
		Results:                []querySpanResult{res},
		NotFoundError:          targetNF,
		UnresolvedColumnsError: targetUC,
	}
	if span.Type != analysis.Select || len(span.Results) != 1 {
		t.Fatalf("unexpected aliased ResolvedSpan: %+v", span)
	}

	var changes analysis.StatementChanges = statementChanges{
		Resources: []changedResource{{Database: "db", Table: "t", AffectsData: true}},
		DMLCount:  1,
	}
	if len(changes.Resources) != 1 || !changes.Resources[0].AffectsData {
		t.Fatalf("unexpected aliased StatementChanges: %+v", changes)
	}

	var br analysis.ByteRange = byteRange{Start: 0, End: 8}
	var stmt analysis.Statement = statement{Text: "SELECT 1", Start: br.Start, End: br.End}
	if stmt.Text != "SELECT 1" || stmt.End != 8 {
		t.Fatalf("unexpected aliased Statement/ByteRange: %+v %+v", stmt, br)
	}

	aliasedMix := analysis.ErrMixUserSystemTables
	if !errors.Is(aliasedMix, analysis.ErrMixUserSystemTables) || aliasedMix != analysis.ErrMixUserSystemTables {
		t.Fatal("aliased ErrMixUserSystemTables identity mismatch")
	}
}
