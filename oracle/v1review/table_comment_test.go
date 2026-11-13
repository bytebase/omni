package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/v1review"
)

// tableFinding is a finding of a rule judging statements together as the
// tests compare it: the code, the qualified name, the statement's index,
// and the text its range anchors.
type tableFinding struct {
	Code      v1review.Code
	Name      string
	Statement int
	Text      string
}

func runTableRule(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []tableFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []tableFinding
	for _, f := range rule(list.Items) {
		got = append(got, tableFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestTableComment(t *testing.T) {
	missing := func(name string, statement int, text string) tableFinding {
		return tableFinding{v1review.TableCommentMissing, name, statement, text}
	}
	tooLong := func(name string, statement int, text string) tableFinding {
		return tableFinding{v1review.TableCommentTooLong, name, statement, text}
	}
	required := &v1review.TableCommentConfig{Required: true, MaxLength: 5}
	for _, tt := range []struct {
		name    string
		sql     string
		current string
		config  *v1review.TableCommentConfig
		want    []tableFinding
	}{
		{name: "commented", sql: "CREATE TABLE t (a INT);\nCOMMENT ON TABLE t IS 'ok';", current: "S", config: required},
		{
			name:    "no comment, qualified by the current schema",
			sql:     "CREATE TABLE t (a INT);",
			current: "S",
			config:  required,
			want:    []tableFinding{missing("S.T", 0, "CREATE TABLE t (a INT)")},
		},
		{
			name:   "no current schema leaves the name bare",
			sql:    "CREATE TABLE t (a INT);",
			config: required,
			want:   []tableFinding{missing("T", 0, "CREATE TABLE t (a INT)")},
		},
		{name: "comment before the create counts", sql: "COMMENT ON TABLE t IS 'ok';\nCREATE TABLE t (a INT);", current: "S", config: required},
		{name: "schema written or current is the same table", sql: "CREATE TABLE s.t (a INT);\nCOMMENT ON TABLE t IS 'ok';", current: "S", config: required},
		{
			name:    "last comment wins, empty is missing",
			sql:     "CREATE TABLE t (a INT);\nCOMMENT ON TABLE t IS 'ok';\nCOMMENT ON TABLE t IS '';",
			current: "S",
			config:  required,
			want:    []tableFinding{missing("S.T", 0, "CREATE TABLE t (a INT)")},
		},
		{
			name:    "too long in bytes anchors the create",
			sql:     "CREATE TABLE t (a INT);\nCOMMENT ON TABLE t IS 'ééé';",
			current: "S",
			config:  required,
			want:    []tableFinding{tooLong("S.T", 0, "CREATE TABLE t (a INT)")},
		},
		{name: "zero length is no limit", sql: "CREATE TABLE t (a INT);\nCOMMENT ON TABLE t IS 'a long comment';", config: &v1review.TableCommentConfig{Required: true}},
		{name: "not required", sql: "CREATE TABLE t (a INT);", config: &v1review.TableCommentConfig{MaxLength: 5}},
		{
			name:    "a table created twice is judged twice on its last create",
			sql:     "CREATE TABLE b (a INT);\nCREATE TABLE a (a INT);\nCREATE TABLE b (x INT);",
			current: "S",
			config:  required,
			want: []tableFinding{
				missing("S.B", 2, "CREATE TABLE b (x INT)"),
				missing("S.A", 1, "CREATE TABLE a (a INT)"),
				missing("S.B", 2, "CREATE TABLE b (x INT)"),
			},
		},
		{
			name:    "quoted names keep their case",
			sql:     "CREATE TABLE \"t\" (a INT);\nCOMMENT ON TABLE t IS 'ok';",
			current: "S",
			config:  required,
			want:    []tableFinding{missing("S.t", 0, "CREATE TABLE \"t\" (a INT)")},
		},
		{name: "a column comment satisfies the table", sql: "CREATE TABLE t (a INT);\nCOMMENT ON COLUMN t.a IS 'ok';", current: "S", config: required},
		{
			name:   "a column comment is judged as the table's",
			sql:    "CREATE TABLE t (a INT);\nCOMMENT ON COLUMN t.a IS 'a long comment';",
			config: &v1review.TableCommentConfig{MaxLength: 1},
			want:   []tableFinding{tooLong("T", 0, "CREATE TABLE t (a INT)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runTableRule(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return TableComment(stmts, tt.current, tt.config) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTableCommentNilConfigPanicsOnlyWhenJudging(t *testing.T) {
	list, err := parser.Parse("COMMENT ON TABLE t IS 'x';\nALTER TABLE t ADD b INT;")
	if err != nil {
		t.Fatal(err)
	}
	if got := TableComment(list.Items, "S", nil); got != nil {
		t.Errorf("got %v, want none", got)
	}
	list, err = parser.Parse("CREATE TABLE t (a INT);")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recover() == nil {
			t.Error("want a panic")
		}
	}()
	TableComment(list.Items, "S", nil)
}
