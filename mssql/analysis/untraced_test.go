package analysis

import (
	"strings"
	"testing"
)

// TestResolveQuerySpanLinkedServer pins the answer for a statement that reads
// another server: ResolveQuerySpan returns a *LinkedServerError and no span,
// because the connected instance's metadata describes no relation the server
// resolves (TestOracleLinkedServerRead).
func TestResolveQuerySpanLinkedServer(t *testing.T) {
	tests := []struct {
		statement string
		want      *LinkedServerError
	}{
		{"SELECT * FROM srv.db.dbo.t", &LinkedServerError{Server: "srv", Table: "db.dbo.t"}},
		{"SELECT a FROM srv.db.dbo.t AS x", &LinkedServerError{Server: "srv", Table: "db.dbo.t"}},
		{"SELECT a FROM [my srv].db.dbo.t", &LinkedServerError{Server: "my srv", Table: "db.dbo.t"}},
		// A read inside a subquery or a join is the same read.
		{"SELECT a FROM t WHERE b IN (SELECT x FROM srv.db.dbo.t)", &LinkedServerError{Server: "srv", Table: "db.dbo.t"}},
		{"SELECT a FROM t JOIN srv.db.dbo.t AS s ON s.a = t.a", &LinkedServerError{Server: "srv", Table: "db.dbo.t"}},
		// A temp table filled from the remote table copies its rows.
		{"SELECT a INTO #copied FROM srv.db.dbo.t", &LinkedServerError{Server: "srv", Table: "db.dbo.t"}},
		// Pass-through queries spell no column of the remote result.
		{"SELECT * FROM OPENQUERY(srv, 'SELECT a FROM t')", &LinkedServerError{Server: "srv", PassThrough: true}},
		{"SELECT o.a FROM OPENQUERY([my srv], 'SELECT a FROM t') AS o", &LinkedServerError{Server: "my srv", PassThrough: true}},
		// An ad-hoc provider has no linked server to name, and a pass-through
		// call the parser reads no server from is still a pass-through call.
		{"SELECT * FROM OPENROWSET('SQLNCLI', 'x', 'SELECT 1')", &LinkedServerError{}},
		{"SELECT * FROM OPENQUERY('srv', 'SELECT 1')", &LinkedServerError{PassThrough: true}},
		{"SELECT * FROM OPENDATASOURCE('SQLNCLI', 'x').db.dbo.t", &LinkedServerError{}},
	}
	for _, tc := range tests {
		_, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog()})
		linked, ok := err.(*LinkedServerError)
		if !ok {
			t.Errorf("%q: got %v, want *LinkedServerError", tc.statement, err)
			continue
		}
		if linked.Server != tc.want.Server || linked.Table != tc.want.Table || linked.PassThrough != tc.want.PassThrough {
			t.Errorf("%q: got %+v, want %+v", tc.statement, *linked, *tc.want)
		}
	}
}

// TestLinkedServerErrorText pins the sentence a caller shows for each shape.
func TestLinkedServerErrorText(t *testing.T) {
	for _, tc := range []struct {
		err  *LinkedServerError
		want string
	}{
		{
			err:  &LinkedServerError{Server: "srv", Table: "db.dbo.t"},
			want: `linked server "srv": the connected instance's metadata does not describe "db.dbo.t", so the columns it returns cannot be traced`,
		},
		{
			err:  &LinkedServerError{Server: "srv"},
			want: `linked server "srv": a pass-through query's columns are the remote server's, so they cannot be traced`,
		},
		{
			err:  &LinkedServerError{PassThrough: true},
			want: "a pass-through query's columns are the remote server's, so they cannot be traced",
		},
		{
			err:  &LinkedServerError{},
			want: "an ad-hoc provider read's columns are the provider's, not the connected instance's, so they cannot be traced",
		},
	} {
		if got := tc.err.Error(); got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

// TestResolveQuerySpanLinkedServerWrite pins where the refusal stops: a write
// through a linked server returns no rows, so its span names the table it
// writes, with Server set, and needs no lineage to mask.
func TestResolveQuerySpanLinkedServerWrite(t *testing.T) {
	server := ColumnResource{Server: "srv", Database: "db", Schema: "dbo", Table: "t"}
	for statement, want := range map[string]ColumnSet{
		"UPDATE srv.db.dbo.t SET a = 1":             set(server),
		"INSERT srv.db.dbo.t (a) VALUES (1)":        set(server),
		"DELETE srv.db.dbo.t WHERE a = 1":           set(server),
		"INSERT t SELECT a FROM srv.db.dbo.t":       set(server, column("t", "")),
		"UPDATE t SET a = b FROM srv.db.dbo.t AS b": set(server, column("t", "")),
	} {
		span, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog()})
		if err != nil {
			t.Errorf("%q: got %v, want a span", statement, err)
			continue
		}
		if len(span.Results) != 0 {
			t.Errorf("%q: got results %+v, want none", statement, span.Results)
		}
		for resource := range want {
			if !span.SourceColumns[resource] {
				t.Errorf("%q: got sources %v, want %s", statement, span.SourceColumns, resource.String())
			}
		}
	}
}

// TestResolveQuerySpanTempTableUntraced walks a batch the way a masking caller
// does, one statement per ResolveQuerySpan call with the temp tables carried
// between them, and pins the answer for each read: the relation is named and
// the rows are untraced, whatever wrote the rows.
func TestResolveQuerySpanTempTableUntraced(t *testing.T) {
	tempTables := map[string]*TempTable{}
	opts := ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog(), TempTables: tempTables}
	steps := []struct {
		statement string
		relations []string
	}{
		// The writes return no rows, so there is nothing to withhold. Each
		// shape SQL Server accepts is here (TestOracleTempTableRows).
		{"CREATE TABLE #t (a INT, secret NVARCHAR(20))", nil},
		{"SELECT a, b INTO #into FROM t", nil},
		{"INSERT #t VALUES (1, 'x')", nil},
		{"INSERT #t (a) SELECT b FROM t1", nil},
		{"UPDATE #t SET b = t.b FROM t WHERE t.a = #t.a", nil},
		{"MERGE #into AS tgt USING (SELECT a FROM t) AS src ON tgt.a = src.a WHEN NOT MATCHED THEN INSERT (a) VALUES (src.a);", nil},
		// A read of what those statements wrote.
		{"SELECT b FROM #into", []string{"#into"}},
		{"SELECT * FROM #into", []string{"#into"}},
		{"SELECT x.a FROM #into AS x", []string{"#into"}},
		{"SELECT a FROM t WHERE a IN (SELECT a FROM #into)", []string{"#into"}},
		{"WITH c AS (SELECT a FROM #t) SELECT a FROM c", []string{"#t"}},
		// The relation is named as the batch declared it, whatever case the
		// statement wrote.
		{"SELECT a FROM #INTO", []string{"#into"}},
		// A query of only the server's metadata returns before the FROM clause
		// is resolved, and still reads the temp table it joins.
		{"SELECT o.name FROM #into AS x CROSS JOIN sys.objects AS o", []string{"#into"}},
		{"SELECT o.name FROM @tv AS v CROSS JOIN sys.objects AS o", []string{"@tv"}},
		// The same table through tempdb, and a name the batch did not declare:
		// its columns are unknown too, so SELECT * names none and the rows are
		// untraced.
		{"SELECT a FROM tempdb..#t", []string{"#t"}},
		{"SELECT * FROM #undeclared", []string{"#undeclared"}},
		// A table variable is the same shape.
		{"DECLARE @tv TABLE (id INT)", nil},
		{"SELECT id FROM @tv", []string{"@tv"}},
		// A statement that reads no temp table is traced as before.
		{"SELECT b FROM t", nil},
	}
	for _, step := range steps {
		span, err := ResolveQuerySpan(step.statement, opts)
		if err != nil {
			t.Fatalf("%q: %v", step.statement, err)
		}
		var got []string
		if span.Untraced != nil {
			for _, relation := range span.Untraced.Relations {
				if relation.Column != "" {
					t.Errorf("%q: relation %+v names a column, want a whole table", step.statement, relation)
				}
				got = append(got, relation.String())
			}
		}
		if strings.Join(got, ",") != strings.Join(step.relations, ",") {
			t.Errorf("%q: untraced relations %v, want %v", step.statement, got, step.relations)
		}
	}
	if got := tempTables["#into"]; got == nil || strings.Join(got.Columns, ",") != "a,b" {
		t.Errorf("#into: got %+v, want columns a,b", got)
	}
	// A column of a temp table the caller did not carry is an error, as a
	// column of a table the catalog does not have is.
	if _, err := ResolveQuerySpan("SELECT a FROM #undeclared", opts); err == nil {
		t.Error("SELECT a FROM #undeclared: got a span, want an error")
	}
}

// TestUntracedTempTableKeepsColumns pins what a span still gives a caller for a
// temp-table read: the result columns, each pointing at the temp table, with
// Untraced saying that chain stops there.
func TestUntracedTempTableKeepsColumns(t *testing.T) {
	tempTables := map[string]*TempTable{"#t": {Name: "#t", Columns: []string{"a", "secret"}}}
	span, err := ResolveQuerySpan("SELECT secret FROM #t", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog(), TempTables: tempTables})
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(span.Results) != 1 || span.Results[0].Name != "secret" {
		t.Fatalf("got results %+v", span.Results)
	}
	if want := set(ColumnResource{Table: "#t", Column: "secret"}); !equalSet(span.Results[0].SourceColumns, want) {
		t.Errorf("got sources %v, want %v", span.Results[0].SourceColumns, want)
	}
	want := "the statement's rows cannot be traced to columns: it reads #t (whose rows the batch's own statements wrote)"
	if span.Untraced == nil || span.Untraced.Error() != want {
		t.Errorf("got %v, want %q", span.Untraced, want)
	}
}

// TestResolveQuerySpanVariableUntraced pins the answer for a statement that
// reads a local variable: the rows are untraced, because a statement of the
// batch gave the variable its value from a column this call did not see.
func TestResolveQuerySpanVariableUntraced(t *testing.T) {
	tests := []struct {
		statement string
		variables []string
	}{
		{"SELECT @v AS v", []string{"@v"}},
		{"SELECT @v + a AS x FROM t", []string{"@v"}},
		{"SELECT a FROM t WHERE a = @v", []string{"@v"}},
		{"SELECT a FROM t WHERE b BETWEEN @lo AND @hi", []string{"@hi", "@lo"}},
		{"SELECT @v, @w AS w", []string{"@v", "@w"}},
		// A rowset function's argument is read like any other value.
		{"SELECT * FROM OPENJSON(@json)", []string{"@json"}},
		{"SELECT * FROM STRING_SPLIT(@list, ',')", []string{"@list"}},
		// A variable's XML method names the variable in the node, not in a
		// child, and its node columns come from that document.
		{"SELECT c FROM @x.nodes('/r') AS n(c)", []string{"@x"}},
		{"SELECT COUNT(*) FROM @x.nodes('/r') AS n(c)", []string{"@x"}},
		// A system function's value is the server's state, not a column.
		{"SELECT @@VERSION AS v", nil},
		{"SELECT a FROM t WHERE a = @@SPID", nil},
		// A statement that returns no rows has nothing to withhold.
		{"DECLARE @v INT = (SELECT a FROM t)", nil},
		{"SET @v = (SELECT a FROM t)", nil},
		{"SELECT @v = a FROM t", nil},
		// A statement that writes a table has no result set to withhold.
		{"SELECT a INTO #t FROM t", nil},
		// A literal is not a variable.
		{"SELECT a FROM t", nil},
	}
	for _, tc := range tests {
		span, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog()})
		if err != nil {
			t.Errorf("%q: %v", tc.statement, err)
			continue
		}
		var got []string
		if span.Untraced != nil {
			got = span.Untraced.Variables
		}
		if strings.Join(got, ",") != strings.Join(tc.variables, ",") {
			t.Errorf("%q: untraced variables %v, want %v", tc.statement, got, tc.variables)
		}
	}
}

// TestResolveQuerySpanSystemMetadata pins the answer for a query of the
// server's own metadata, which classifies SelectInfoSchema and so returns
// before the result columns are extracted: what it reads from a variable still
// has no lineage.
func TestResolveQuerySpanSystemMetadata(t *testing.T) {
	for _, tc := range []struct {
		statement string
		variables []string
	}{
		{"SELECT @v FROM sys.objects", []string{"@v"}},
		{"SELECT * FROM sys.objects WHERE name = @v", []string{"@v"}},
		{"SELECT @v FROM master.dbo.spt_values", []string{"@v"}},
		{"SELECT * FROM sys.objects", nil},
		{"SELECT name FROM master.dbo.spt_values", nil},
	} {
		span, err := ResolveQuerySpan(tc.statement, ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog()})
		if err != nil {
			t.Errorf("%q: %v", tc.statement, err)
			continue
		}
		if span.Type != SelectInfoSchema {
			t.Errorf("%q: type %v, want SelectInfoSchema", tc.statement, span.Type)
		}
		var got []string
		if span.Untraced != nil {
			got = span.Untraced.Variables
		}
		if strings.Join(got, ",") != strings.Join(tc.variables, ",") {
			t.Errorf("%q: untraced variables %v, want %v", tc.statement, got, tc.variables)
		}
	}
}

// TestResolveQuerySpanVariableAndTempTable pins a statement that reads both,
// and the sentence the caller shows.
func TestResolveQuerySpanVariableAndTempTable(t *testing.T) {
	tempTables := map[string]*TempTable{"#t": {Name: "#t", Columns: []string{"a"}}}
	span, err := ResolveQuerySpan("SELECT a FROM #t WHERE a = @v", ResolveOptions{DefaultDatabase: "db", Catalog: testCatalog(), TempTables: tempTables})
	if err != nil {
		t.Fatalf("%v", err)
	}
	want := "the statement's rows cannot be traced to columns: it reads #t (whose rows the batch's own statements wrote) and @v (whose value a statement of the batch assigned)"
	if span.Untraced == nil || span.Untraced.Error() != want {
		t.Errorf("got %v, want %q", span.Untraced, want)
	}
}

// TestClassifySelectAssign pins the classification a statement that assigns a
// variable gets: SQL Server returns no result set for it (TestOracleVariableRows),
// so the span carries none, and the value's tables stay in the access set.
func TestClassifySelectAssign(t *testing.T) {
	for _, statement := range []string{
		"SELECT @v = a FROM t",
		"SELECT @v = (SELECT a FROM t1)",
		"SELECT @v = a, @w = b FROM t",
	} {
		span, err := ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "db", IgnoreCaseSensitive: true, Catalog: testCatalog()})
		if err != nil {
			t.Fatalf("%q: %v", statement, err)
		}
		if len(span.Results) != 0 {
			t.Errorf("%q: got results %+v, want none", statement, span.Results)
		}
		if span.Untraced != nil {
			t.Errorf("%q: got untraced %v, want none: the statement returns no rows", statement, span.Untraced)
		}
	}
	span := resolve(t, "SELECT @v = a FROM t")
	if !equalSet(span.SourceColumns, set(ColumnResource{Database: "db", Schema: "dbo", Table: "t"})) {
		t.Errorf("got sources %v, want the table the value reads", span.SourceColumns)
	}
}
