package review

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/review"
)

// span is the byte range of the first occurrence of text in sql.
func span(t *testing.T, sql, text string) review.Range {
	t.Helper()
	i := strings.Index(sql, text)
	if i < 0 {
		t.Fatalf("%q not in %q", text, sql)
	}
	return review.Range{Start: i, End: i + len(text)}
}

// nthSpan is the byte range of the n-th (0-based) occurrence of text.
func nthSpan(t *testing.T, sql, text string, n int) review.Range {
	t.Helper()
	off := 0
	for ; n > 0; n-- {
		i := strings.Index(sql[off:], text)
		if i < 0 {
			t.Fatalf("occurrence %d of %q not in %q", n, text, sql)
		}
		off += i + len(text)
	}
	r := span(t, sql[off:], text)
	return review.Range{Start: r.Start + off, End: r.End + off}
}

func TestReview(t *testing.T) {
	type finding struct {
		rule      review.Rule
		statement int
		rng       review.Range
		message   string
	}
	tests := []struct {
		name    string
		sql     string
		opts    review.Options
		targets int
		want    func(t *testing.T, sql string) []finding
		ranges  []review.Range
	}{
		{
			name:    "syntax error is the only finding",
			sql:     "DELETE FROM t; SELEC 2; DELETE FROM u;",
			targets: 2,
			ranges:  []review.Range{{Start: 0, End: 14}, {Start: 14, End: 23}, {Start: 23, End: 38}},
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.Syntax, 1, review.Range{Start: 15, End: 15}, `syntax error at or near "SELEC"`}}
			},
		},
		{
			name:    "syntax error with the rule switched off is still reported",
			sql:     "SELEC 1",
			opts:    review.Options{Rules: []review.Rule{review.RequireWhere}},
			targets: 1,
			ranges:  []review.Range{{Start: 0, End: 7}},
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.Syntax, 0, review.Range{}, `syntax error at or near "SELEC"`}}
			},
		},
		{
			name:    "require where",
			sql:     "UPDATE t SET a = 1; DELETE FROM s.t USING u; DELETE FROM t WHERE a = 1;\nWITH d AS (DELETE FROM t RETURNING *) SELECT 1; MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN DELETE;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.RequireWhere, 0, span(t, sql, "UPDATE t SET a = 1"), "UPDATE t has no WHERE clause"},
					{review.RequireWhere, 1, span(t, sql, "DELETE FROM s.t USING u"), "DELETE FROM s.t has no WHERE clause"},
					{review.RequireWhere, 3, span(t, sql, "DELETE FROM t RETURNING *"), "DELETE FROM t has no WHERE clause"},
				}
			},
		},
		{
			name:    "require is null",
			sql:     "SELECT 1 FROM t WHERE a = NULL OR NULL != b OR c IS NULL OR d = NULL::int OR (e + 1) <> NULL OR f IS DISTINCT FROM NULL OR g OPERATOR(custom.=) NULL OR h OPERATOR(pg_catalog.<>) NULL;\nUPDATE \"T\" SET x = NULL WHERE y\n  =\n NULL;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.RequireIsNull, 0, span(t, sql, "a = NULL"), `"a = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, "NULL != b"), `"NULL != b" is never true; use IS NOT NULL`},
					{review.RequireIsNull, 0, span(t, sql, "d = NULL::int"), `"d = NULL::int" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, "(e + 1) <> NULL"), `"(e + 1) <> NULL" is never true; use IS NOT NULL`},
					{review.RequireIsNull, 0, span(t, sql, "h OPERATOR(pg_catalog.<>) NULL"), `"h OPERATOR(pg_catalog.<>) NULL" is never true; use IS NOT NULL`},
					{review.RequireIsNull, 1, span(t, sql, "y\n  =\n NULL"), `"y = NULL" is never true; use IS NULL`},
				}
			},
		},
		{
			name:    "disallow drop object",
			sql:     `DROP TABLE IF EXISTS a.t1, "T2" CASCADE; DROP VIEW v; DROP TRIGGER trg ON s.t; DROP FUNCTION f(int, text); DROP DATABASE d; DROP ROLE r; DROP SCHEMA s; ALTER TABLE t DROP COLUMN c, DROP CONSTRAINT k, DROP c2; DROP INDEX CONCURRENTLY i; DROP OWNED BY r;`,
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 0, span(t, sql, `DROP TABLE IF EXISTS a.t1, "T2" CASCADE`), `drops table a.t1, "T2"`},
					{review.DisallowDropObject, 1, span(t, sql, "DROP VIEW v"), "drops view v"},
					{review.DisallowDropObject, 2, span(t, sql, "DROP TRIGGER trg ON s.t"), "drops trigger trg on s.t"},
					{review.DisallowDropObject, 3, span(t, sql, "DROP FUNCTION f(int, text)"), "drops function f(int4, text)"},
					{review.DisallowDropObject, 4, span(t, sql, "DROP DATABASE d"), "drops database d"},
					{review.DisallowDropObject, 5, span(t, sql, "DROP ROLE r"), "drops role r"},
					{review.DisallowDropObject, 6, span(t, sql, "DROP SCHEMA s"), "drops schema s"},
					{review.DisallowDropObject, 7, span(t, sql, "DROP COLUMN c"), "drops column c of t"},
					{review.DisallowDropObject, 7, span(t, sql, "DROP c2"), "drops column c2 of t"},
					{review.DisallowDropObject, 8, span(t, sql, "DROP INDEX CONCURRENTLY i"), "drops index i"},
					{review.DisallowDropObject, 9, span(t, sql, "DROP OWNED BY r"), "drops objects owned by r"},
				}
			},
		},
		{
			name:    "disallow drop object names every object shape",
			sql:     "DROP CAST (int AS text); DROP TRANSFORM FOR int LANGUAGE plpgsql; DROP OPERATOR CLASS s.c USING btree; DROP OPERATOR FAMILY f USING btree; DROP OPERATOR s.+ (int, int), - (NONE, int); DROP AGGREGATE agg(int, text), s.agg2(*); DROP ROUTINE r; DROP RULE r ON s.t; DROP POLICY p ON t; DROP TYPE s.ty, ty2; DROP FUNCTION f(int[], text[][]);",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 0, span(t, sql, "DROP CAST (int AS text)"), "drops cast from int4 to text"},
					{review.DisallowDropObject, 1, span(t, sql, "DROP TRANSFORM FOR int LANGUAGE plpgsql"), "drops transform for int4 language plpgsql"},
					{review.DisallowDropObject, 2, span(t, sql, "DROP OPERATOR CLASS s.c USING btree"), "drops operator class s.c using btree"},
					{review.DisallowDropObject, 3, span(t, sql, "DROP OPERATOR FAMILY f USING btree"), "drops operator family f using btree"},
					{review.DisallowDropObject, 4, span(t, sql, "DROP OPERATOR s.+ (int, int), - (NONE, int)"), `drops operator s."+"(int4, int4), "-"(NONE, int4)`},
					{review.DisallowDropObject, 5, span(t, sql, "DROP AGGREGATE agg(int, text), s.agg2(*)"), "drops aggregate agg(int4, text), s.agg2(*)"},
					{review.DisallowDropObject, 6, span(t, sql, "DROP ROUTINE r"), "drops routine r"},
					{review.DisallowDropObject, 7, span(t, sql, "DROP RULE r ON s.t"), "drops rule r on s.t"},
					{review.DisallowDropObject, 8, span(t, sql, "DROP POLICY p ON t"), "drops policy p on t"},
					{review.DisallowDropObject, 9, span(t, sql, "DROP TYPE s.ty, ty2"), "drops type s.ty, ty2"},
					{review.DisallowDropObject, 10, span(t, sql, "DROP FUNCTION f(int[], text[][])"), "drops function f(int4[], text[][])"},
				}
			},
		},
		{
			name:    "statements inside a SQL-standard routine body are reported where they are",
			sql:     "CREATE PROCEDURE p() LANGUAGE SQL BEGIN ATOMIC DROP TABLE t; TRUNCATE u; ALTER TABLE v RENAME TO w; DELETE FROM x; END;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 0, span(t, sql, "DROP TABLE t"), "drops table t"},
					{review.DisallowTruncate, 0, span(t, sql, "TRUNCATE u"), "truncates u"},
					{review.DisallowRename, 0, span(t, sql, "ALTER TABLE v RENAME TO w"), "renames table v to w"},
					{review.RequireWhere, 0, span(t, sql, "DELETE FROM x"), "DELETE FROM x has no WHERE clause"},
				}
			},
		},
		{
			name:    "disallow truncate",
			sql:     "TRUNCATE a.t1, \"T2\" RESTART IDENTITY; TRUNCATE TABLE ONLY t CASCADE;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowTruncate, 0, span(t, sql, `TRUNCATE a.t1, "T2" RESTART IDENTITY`), `truncates a.t1, "T2"`},
					{review.DisallowTruncate, 1, span(t, sql, "TRUNCATE TABLE ONLY t CASCADE"), "truncates t"},
				}
			},
		},
		{
			name:    "disallow rename",
			sql:     "ALTER TABLE t RENAME TO u; ALTER TABLE IF EXISTS s.t RENAME COLUMN a TO b; ALTER VIEW v RENAME COLUMN c TO d; ALTER INDEX i RENAME TO j; ALTER TABLE t RENAME CONSTRAINT k TO l; ALTER TABLE t RENAME a TO b;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowRename, 0, span(t, sql, "ALTER TABLE t RENAME TO u"), "renames table t to u"},
					{review.DisallowRename, 1, span(t, sql, "ALTER TABLE IF EXISTS s.t RENAME COLUMN a TO b"), "renames column a of s.t to b"},
					{review.DisallowRename, 2, span(t, sql, "ALTER VIEW v RENAME COLUMN c TO d"), "renames column c of v to d"},
					{review.DisallowRename, 5, span(t, sql, "ALTER TABLE t RENAME a TO b"), "renames column a of t to b"},
				}
			},
		},
		{
			name:    "online migration requested",
			sql:     "ALTER TABLE t ADD COLUMN c int;",
			opts:    review.Options{Change: review.Change{OnlineMigration: true}},
			targets: 3,
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.OnlineMigration, -1, review.Range{}, "the change requests online migration, which PostgreSQL does not support"}}
			},
		},
		{
			name:    "online migration not requested",
			sql:     "ALTER TABLE t ADD COLUMN c int;",
			opts:    review.Options{Change: review.Change{PriorBackup: true}},
			targets: 1,
			want:    func(t *testing.T, sql string) []finding { return nil },
		},
		{
			name:    "rules filter: only the listed rules run, unknown rules are ignored",
			sql:     "DELETE FROM t; TRUNCATE t; DROP TABLE t;",
			opts:    review.Options{Rules: []review.Rule{review.DisallowTruncate, review.Rule("NOT_A_RULE"), review.WalkThrough}},
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.DisallowTruncate, 1, span(t, sql, "TRUNCATE t"), "truncates t"}}
			},
		},
		{
			name:    "empty rules list runs nothing but syntax",
			sql:     "DELETE FROM t;",
			opts:    review.Options{Rules: []review.Rule{}},
			targets: 1,
			want:    func(t *testing.T, sql string) []finding { return nil },
		},
		{
			name:    "findings are sorted by statement then offset",
			sql:     "DROP TABLE t;\nDELETE FROM u;\nSELECT 1 WHERE b = NULL AND a = NULL;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 0, span(t, sql, "DROP TABLE t"), "drops table t"},
					{review.RequireWhere, 1, span(t, sql, "DELETE FROM u"), "DELETE FROM u has no WHERE clause"},
					{review.RequireIsNull, 2, span(t, sql, "b = NULL"), `"b = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 2, span(t, sql, "a = NULL"), `"a = NULL" is never true; use IS NULL`},
				}
			},
		},
		{
			name:    "comments and metacommands are not statements",
			sql:     "\\restrict abc\n-- a comment\n/* block */ ;\nDELETE FROM t;\n\\unrestrict abc\n",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.RequireWhere, 0, span(t, sql, "DELETE FROM t"), "DELETE FROM t has no WHERE clause"}}
			},
		},
		{
			name:    "no statements",
			sql:     "  -- nothing\n",
			targets: 1,
			want:    func(t *testing.T, sql string) []finding { return nil },
		},
		{
			name:    "no targets",
			sql:     "DELETE FROM t",
			targets: 0,
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.RequireWhere, 0, span(t, sql, "DELETE FROM t"), "DELETE FROM t has no WHERE clause"}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targets := make([]review.Target, tt.targets)
			result, err := Review(context.Background(), tt.sql, tt.opts, targets)
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			if tt.ranges != nil {
				if diff := cmp.Diff(tt.ranges, result.Statements); diff != "" {
					t.Errorf("Statements mismatch (-want +got):\n%s", diff)
				}
			}
			if len(result.Failures) != 0 {
				t.Errorf("Failures = %v, want none", result.Failures)
			}
			var want []review.Finding
			for _, f := range tt.want(t, tt.sql) {
				want = append(want, review.Finding{Rule: f.rule, Statement: f.statement, Range: f.rng, Message: f.message, Targets: allTargets(tt.targets)})
			}
			if diff := cmp.Diff(want, result.Findings, cmp.Transformer("nilSlice", func(f []review.Finding) []review.Finding {
				if f == nil {
					return []review.Finding{}
				}
				return f
			})); diff != "" {
				t.Errorf("Findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReviewStatementsListEveryStatement(t *testing.T) {
	sql := "SELECT 1;\n\nSELECT 2\n;SELECT 3"
	result, err := Review(context.Background(), sql, review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []review.Range{{Start: 0, End: 9}, {Start: 9, End: 21}, {Start: 21, End: 29}}
	if diff := cmp.Diff(want, result.Statements); diff != "" {
		t.Errorf("Statements mismatch (-want +got):\n%s", diff)
	}
}

func TestReviewCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Review(ctx, "SELECT 1", review.Options{}, nil); err == nil {
		t.Fatal("Review with a canceled context returned no error")
	}
}

func TestRulesAreDistinct(t *testing.T) {
	seen := make(map[review.Rule]bool)
	for _, rule := range Rules {
		if seen[rule] {
			t.Errorf("rule %s listed twice", rule)
		}
		seen[rule] = true
	}
}
