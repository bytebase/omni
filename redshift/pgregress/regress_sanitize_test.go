package pgregress

import "testing"

func TestSQLForParserPsqlVariableContexts(t *testing.T) {
	tests := []struct {
		name string
		stmt ExtractedStmt
		want string
	}{
		{
			name: "standalone psql variable macro",
			stmt: ExtractedStmt{SQL: ":init_range_parted", HasPsqlVar: true},
			want: "SELECT 1",
		},
		{
			name: "commented standalone psql variable macro",
			stmt: ExtractedStmt{SQL: "-- reset test fixture\n:init_range_parted", HasPsqlVar: true},
			want: "SELECT 1",
		},
		{
			name: "leading psql variable macro before SQL",
			stmt: ExtractedStmt{SQL: ":init_range_parted;\nUPDATE range_parted SET c = 95 RETURNING *", HasPsqlVar: true},
			want: "UPDATE range_parted SET c = 95 RETURNING *",
		},
		{
			name: "large object oid variable",
			stmt: ExtractedStmt{SQL: "COMMENT ON LARGE OBJECT :newloid IS 'copied object'", HasPsqlVar: true},
			want: "COMMENT ON LARGE OBJECT 1 IS 'copied object'",
		},
		{
			name: "explain variable query",
			stmt: ExtractedStmt{SQL: "EXPLAIN :query", HasPsqlVar: true},
			want: "EXPLAIN SELECT 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sqlForParser(tt.stmt)
			if got != tt.want {
				t.Fatalf("sqlForParser() = %q, want %q", got, tt.want)
			}
		})
	}
}
