package completion

import (
	"strings"
	"testing"
)

func TestExtractTableRefs(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []TableRef
	}{
		// --- Simple SELECT ---
		{
			name: "simple SELECT FROM",
			sql:  "SELECT * FROM t",
			want: []TableRef{{Table: "t"}},
		},
		// --- Alias ---
		{
			name: "SELECT FROM with AS alias",
			sql:  "SELECT * FROM t AS x",
			want: []TableRef{{Table: "t", Alias: "x"}},
		},
		{
			name: "SELECT FROM with bare alias",
			sql:  "SELECT * FROM t x",
			want: []TableRef{{Table: "t", Alias: "x"}},
		},
		// --- JOIN ---
		{
			name: "FROM with JOIN",
			sql:  "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id",
			want: []TableRef{{Table: "t1"}, {Table: "t2"}},
		},
		{
			name: "FROM with LEFT JOIN",
			sql:  "SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id",
			want: []TableRef{{Table: "t1"}, {Table: "t2"}},
		},
		// --- Schema-qualified ---
		{
			name: "schema-qualified table",
			sql:  "SELECT * FROM dbo.t",
			want: []TableRef{{Schema: "dbo", Table: "t"}},
		},
		// --- UPDATE ---
		{
			name: "UPDATE table",
			sql:  "UPDATE t SET x = 1",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "UPDATE schema-qualified",
			sql:  "UPDATE dbo.t SET x = 1",
			want: []TableRef{{Schema: "dbo", Table: "t"}},
		},
		// --- INSERT ---
		{
			name: "INSERT INTO table",
			sql:  "INSERT INTO t (a) VALUES (1)",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "INSERT INTO schema-qualified",
			sql:  "INSERT INTO dbo.t (a) VALUES (1)",
			want: []TableRef{{Schema: "dbo", Table: "t"}},
		},
		// --- DELETE ---
		{
			name: "DELETE FROM table",
			sql:  "DELETE FROM t WHERE id = 1",
			want: []TableRef{{Table: "t"}},
		},
		// --- MERGE ---
		{
			name: "MERGE INTO table",
			sql:  "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.x = s.x;",
			want: []TableRef{{Table: "t"}, {Table: "s"}},
		},
		// --- Multiple tables ---
		{
			name: "multiple JOINs",
			sql:  "SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id LEFT JOIN t3 ON t2.id = t3.id",
			want: []TableRef{{Table: "t1"}, {Table: "t2"}, {Table: "t3"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTableRefs(tt.sql, len(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("extractTableRefs(%q): got %d refs, want %d\n  got:  %+v\n  want: %+v",
					tt.sql, len(got), len(tt.want), got, tt.want)
			}
			for i := range tt.want {
				if got[i].Schema != tt.want[i].Schema ||
					got[i].Table != tt.want[i].Table ||
					got[i].Alias != tt.want[i].Alias {
					t.Errorf("extractTableRefs(%q)[%d] = %+v, want %+v",
						tt.sql, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExtractTableRefsUsesCompletionScopeCurrentSetOperationArm(t *testing.T) {
	sql := "SELECT Id FROM Employees UNION SELECT EmployeeId FROM Address UNION SELECT  FROM MySchema.SalaryLevel"
	cursor := strings.Index(sql, " FROM MySchema.SalaryLevel")
	got := extractTableRefs(sql, cursor)
	want := []TableRef{{Schema: "MySchema", Table: "SalaryLevel"}}
	if len(got) != len(want) {
		t.Fatalf("extractTableRefs(%q): got %d refs, want %d\n  got:  %+v\n  want: %+v", sql, len(got), len(want), got, want)
	}
	for i := range want {
		if got[i].Schema != want[i].Schema || got[i].Table != want[i].Table || got[i].Alias != want[i].Alias {
			t.Fatalf("extractTableRefs(%q)[%d] = %+v, want %+v", sql, i, got[i], want[i])
		}
	}
}

func TestExtractTableRefsLexerFallback(t *testing.T) {
	// Incomplete SQL that fails AST parsing should fall back to lexer extraction.
	tests := []struct {
		name string
		sql  string
		want []TableRef
	}{
		{
			name: "incomplete SELECT FROM",
			sql:  "SELECT * FROM t WHERE ",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "incomplete JOIN",
			sql:  "SELECT * FROM t1 JOIN t2 ON ",
			want: []TableRef{{Table: "t1"}, {Table: "t2"}},
		},
		{
			name: "incomplete UPDATE",
			sql:  "UPDATE t SET ",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "incomplete INSERT",
			sql:  "INSERT INTO t ",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "incomplete DELETE",
			sql:  "DELETE FROM t WHERE ",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "incomplete MERGE",
			sql:  "MERGE INTO t ",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "schema-qualified in fallback",
			sql:  "SELECT * FROM dbo.t WHERE ",
			want: []TableRef{{Schema: "dbo", Table: "t"}},
		},
		{
			name: "alias in fallback",
			sql:  "SELECT * FROM t AS x WHERE ",
			want: []TableRef{{Table: "t", Alias: "x"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTableRefs(tt.sql, len(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("extractTableRefs(%q): got %d refs, want %d\n  got:  %+v\n  want: %+v",
					tt.sql, len(got), len(tt.want), got, tt.want)
			}
			for i := range tt.want {
				if got[i].Schema != tt.want[i].Schema ||
					got[i].Table != tt.want[i].Table ||
					got[i].Alias != tt.want[i].Alias {
					t.Errorf("extractTableRefs(%q)[%d] = %+v, want %+v",
						tt.sql, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExtractTableRefsLexerDirect(t *testing.T) {
	// Test the lexer-based extraction directly.
	tests := []struct {
		name string
		sql  string
		want []TableRef
	}{
		{
			name: "simple FROM",
			sql:  "SELECT * FROM t",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "schema.table",
			sql:  "SELECT * FROM dbo.t",
			want: []TableRef{{Schema: "dbo", Table: "t"}},
		},
		{
			name: "JOIN",
			sql:  "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id",
			want: []TableRef{{Table: "t1"}, {Table: "t2"}},
		},
		{
			name: "INTO",
			sql:  "INSERT INTO t (a) VALUES (1)",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "UPDATE",
			sql:  "UPDATE t SET x = 1",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "MERGE INTO",
			sql:  "MERGE INTO t USING s ON t.id = s.id",
			want: []TableRef{{Table: "t"}},
		},
		{
			name: "MERGE without INTO",
			sql:  "MERGE t USING s ON t.id = s.id",
			want: []TableRef{{Table: "t"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTableRefsLexer(tt.sql, len(tt.sql))
			if len(got) != len(tt.want) {
				t.Fatalf("extractTableRefsLexer(%q): got %d refs, want %d\n  got:  %+v\n  want: %+v",
					tt.sql, len(got), len(tt.want), got, tt.want)
			}
			for i := range tt.want {
				if got[i].Schema != tt.want[i].Schema ||
					got[i].Table != tt.want[i].Table ||
					got[i].Alias != tt.want[i].Alias {
					t.Errorf("extractTableRefsLexer(%q)[%d] = %+v, want %+v",
						tt.sql, i, got[i], tt.want[i])
				}
			}
		})
	}
}
