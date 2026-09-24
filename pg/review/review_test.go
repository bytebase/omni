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
			name:    "unterminated string",
			sql:     "'unterminated",
			targets: 1,
			ranges:  []review.Range{{Start: 0, End: 13}},
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.Syntax, 0, review.Range{}, `unterminated quoted string at or near "'unterminated"`}}
			},
		},
		{
			name:    "unterminated comment after a statement, the message on one line",
			sql:     "DELETE FROM t;\n  /* unterminated\ncomment",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.Syntax, 1, review.Range{Start: 17, End: 17}, `unterminated /* comment at or near "/* unterminated\ncomment"`}}
			},
		},
		{
			name:    "unterminated string past comments and metacommands",
			sql:     "DELETE FROM t;\n-- header\n\\restrict abc\n/* block */ 'unterminated",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				at := strings.Index(sql, "'unterminated")
				return []finding{{review.Syntax, 1, review.Range{Start: at, End: at}, `unterminated quoted string at or near "'unterminated"`}}
			},
		},
		{
			name:    "unterminated identifier",
			sql:     "DELETE FROM t; \"abc",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.Syntax, 1, review.Range{Start: 15, End: 15}, `unterminated quoted identifier at or near ""abc"`}}
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
			name:    "require where skips a plain EXPLAIN but not EXPLAIN ANALYZE",
			sql:     "EXPLAIN UPDATE t SET a = 1; EXPLAIN (ANALYZE false) DELETE FROM t; EXPLAIN (ANALYZE off) DELETE FROM t; EXPLAIN (ANALYZE 0) DELETE FROM t;\nEXPLAIN ANALYZE DELETE FROM u; EXPLAIN (ANALYZE, BUFFERS) DELETE FROM v; EXPLAIN (ANALYZE true) DELETE FROM w; EXPLAIN (ANALYZE on) DELETE FROM x; EXPLAIN (ANALYZE 1) DELETE FROM y; EXPLAIN (ANALYZE t) DELETE FROM z; EXPLAIN (ANALYZE ye) DELETE FROM zz; EXPLAIN (ANALYZE o) DELETE FROM t; EXPLAIN (ANALYZE fal) DELETE FROM t; EXPLAIN (ANALYZE n) DELETE FROM t; EXPLAIN (ANALYZE maybe) DELETE FROM t; EXPLAIN (ANALYZE false, ANALYZE true) DELETE FROM last; EXPLAIN (ANALYZE, ANALYZE off) DELETE FROM t; EXPLAIN (ANALYZE maybe, ANALYZE true) DELETE FROM t; EXPLAIN (ANALYZE 2) DELETE FROM t; EXPLAIN (ANALYZE -1) DELETE FROM t; EXPLAIN (ANALYZE true, BUFFERS maybe) DELETE FROM t; EXPLAIN (ANALYZE, BOGUS) DELETE FROM t; EXPLAIN (ANALYZE, FORMAT foo) DELETE FROM t; EXPLAIN (ANALYZE, GENERIC_PLAN) DELETE FROM t; EXPLAIN (ANALYZE, FORMAT json, SERIALIZE binary, BUFFERS off, WAL 1, TIMING 'off') DELETE FROM opts;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.RequireWhere, 4, span(t, sql, "DELETE FROM u"), "DELETE FROM u has no WHERE clause"},
					{review.RequireWhere, 5, span(t, sql, "DELETE FROM v"), "DELETE FROM v has no WHERE clause"},
					{review.RequireWhere, 6, span(t, sql, "DELETE FROM w"), "DELETE FROM w has no WHERE clause"},
					{review.RequireWhere, 7, span(t, sql, "DELETE FROM x"), "DELETE FROM x has no WHERE clause"},
					{review.RequireWhere, 8, span(t, sql, "DELETE FROM y"), "DELETE FROM y has no WHERE clause"},
					{review.RequireWhere, 9, span(t, sql, "DELETE FROM z"), "DELETE FROM z has no WHERE clause"},
					{review.RequireWhere, 10, span(t, sql, "DELETE FROM zz"), "DELETE FROM zz has no WHERE clause"},
					{review.RequireWhere, 15, span(t, sql, "DELETE FROM last"), "DELETE FROM last has no WHERE clause"},
					{review.RequireWhere, 24, span(t, sql, "DELETE FROM opts"), "DELETE FROM opts has no WHERE clause"},
				}
			},
		},
		{
			name:    "require is null",
			sql:     "SELECT 1 FROM t WHERE a = NULL OR NULL != b OR c IS NULL OR d = NULL::int OR (e + 1) <> NULL OR f IS DISTINCT FROM NULL OR g OPERATOR(custom.=) NULL OR h OPERATOR(pg_catalog.<>) NULL OR i = (NULL::text COLLATE \"C\") OR j = NULL COLLATE \"C\" OR k = (NULL::text)::custom OR l = (NULL COLLATE \"C\")::text OR m = NULL COLLATE \"C\" COLLATE \"POSIX\";\nUPDATE \"T\" SET x = NULL WHERE y\n  =\n NULL;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.RequireIsNull, 0, span(t, sql, "a = NULL"), `"a = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, "NULL != b"), `"NULL != b" is never true; use IS NOT NULL`},
					{review.RequireIsNull, 0, span(t, sql, "d = NULL::int"), `"d = NULL::int" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, "(e + 1) <> NULL"), `"(e + 1) <> NULL" is never true; use IS NOT NULL`},
					{review.RequireIsNull, 0, span(t, sql, "h OPERATOR(pg_catalog.<>) NULL"), `"h OPERATOR(pg_catalog.<>) NULL" is never true; use IS NOT NULL`},
					{review.RequireIsNull, 0, span(t, sql, `i = (NULL::text COLLATE "C")`), `"i = (NULL::text COLLATE "C")" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, `j = NULL COLLATE "C"`), `"j = NULL COLLATE "C"" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, `l = (NULL COLLATE "C")::text`), `"l = (NULL COLLATE "C")::text" is never true; use IS NULL`},
					{review.RequireIsNull, 0, span(t, sql, `m = NULL COLLATE "C" COLLATE "POSIX"`), `"m = NULL COLLATE "C" COLLATE "POSIX"" is never true; use IS NULL`},
					{review.RequireIsNull, 1, span(t, sql, "y\n  =\n NULL"), `"y = NULL" is never true; use IS NULL`},
				}
			},
		},
		{
			name:    "require is null looks at predicates only",
			sql:     "SELECT a = NULL AS c, CASE WHEN b = NULL THEN 1 END, CASE d WHEN NULL THEN 1 END, count(*) FILTER (WHERE e = NULL), JSON_ARRAYAGG(x) FILTER (WHERE ee = NULL), (SELECT f = NULL FROM u WHERE g = NULL) FROM t JOIN u ON h = NULL WHERE i IN (SELECT j = NULL FROM v WHERE k = NULL) GROUP BY 1 HAVING l = NULL;\nCREATE TABLE n (x int CHECK (x <> NULL)); CREATE INDEX i ON t (a) WHERE m = NULL; INSERT INTO t VALUES (1) ON CONFLICT (a) WHERE o = NULL DO UPDATE SET a = 1 WHERE p = NULL; UPDATE t SET q = NULL WHERE r = NULL;\nCREATE TABLE d (a boolean DEFAULT (NULL = NULL), b boolean GENERATED ALWAYS AS (NULL = NULL) STORED); CREATE TRIGGER tg BEFORE UPDATE ON t FOR EACH ROW WHEN (NEW.s = NULL) EXECUTE FUNCTION f(); CREATE PUBLICATION pub FOR TABLE t WHERE (u = NULL);\nSELECT 1 WHERE (SELECT v = NULL FROM t LIMIT 1) AND NOT (SELECT w = NULL FROM t) OR (SELECT x = NULL FROM t) IS TRUE OR y = (SELECT z = NULL FROM t) OR EXISTS (SELECT aa = NULL FROM t) OR (SELECT bb = NULL FROM t WHERE cc = NULL) OR COALESCE((SELECT dd = NULL FROM t), false) OR (SELECT ee = NULL FROM t)::boolean;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				var out []finding
				for _, expr := range []string{"b = NULL", "e = NULL", "ee = NULL", "g = NULL", "h = NULL", "k = NULL", "l = NULL"} {
					out = append(out, finding{review.RequireIsNull, 0, span(t, sql, expr), `"` + expr + `" is never true; use IS NULL`})
				}
				out = append(out, finding{review.RequireIsNull, 1, span(t, sql, "x <> NULL"), `"x <> NULL" is never true; use IS NOT NULL`})
				// A comparison in a subquery's select list is a value even
				// when the subquery is the predicate; only cc, in the
				// subquery's own WHERE, is reported from that statement.
				for i, expr := range []string{"m = NULL", "o = NULL", "p = NULL", "r = NULL", "NEW.s = NULL", "u = NULL", "cc = NULL"} {
					stmt := []int{2, 3, 3, 4, 6, 7, 8}[i]
					out = append(out, finding{review.RequireIsNull, stmt, span(t, sql, expr), `"` + expr + `" is never true; use IS NULL`})
				}
				return out
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
			sql:     "DROP CAST (int AS text); DROP TRANSFORM FOR int LANGUAGE plpgsql; DROP OPERATOR CLASS s.c USING btree; DROP OPERATOR FAMILY f USING btree; DROP OPERATOR s.+ (int, int), - (NONE, int); DROP AGGREGATE agg(int, text), s.agg2(*); DROP ROUTINE r; DROP RULE r ON s.t; DROP POLICY p ON t; DROP TYPE s.ty, ty2; DROP FUNCTION f(int[], text[][]); DROP FUNCTION g(t.c%TYPE);",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 0, span(t, sql, "DROP CAST (int AS text)"), "drops cast from int4 to text"},
					{review.DisallowDropObject, 1, span(t, sql, "DROP TRANSFORM FOR int LANGUAGE plpgsql"), "drops transform for int4 language plpgsql"},
					{review.DisallowDropObject, 2, span(t, sql, "DROP OPERATOR CLASS s.c USING btree"), "drops operator class s.c using btree"},
					{review.DisallowDropObject, 3, span(t, sql, "DROP OPERATOR FAMILY f USING btree"), "drops operator family f using btree"},
					{review.DisallowDropObject, 4, span(t, sql, "DROP OPERATOR s.+ (int, int), - (NONE, int)"), `drops operator s.+(int4, int4), -(NONE, int4)`},
					{review.DisallowDropObject, 5, span(t, sql, "DROP AGGREGATE agg(int, text), s.agg2(*)"), "drops aggregate agg(int4, text), s.agg2(*)"},
					{review.DisallowDropObject, 6, span(t, sql, "DROP ROUTINE r"), "drops routine r"},
					{review.DisallowDropObject, 7, span(t, sql, "DROP RULE r ON s.t"), "drops rule r on s.t"},
					{review.DisallowDropObject, 8, span(t, sql, "DROP POLICY p ON t"), "drops policy p on t"},
					{review.DisallowDropObject, 9, span(t, sql, "DROP TYPE s.ty, ty2"), "drops type s.ty, ty2"},
					{review.DisallowDropObject, 10, span(t, sql, "DROP FUNCTION f(int[], text[][])"), "drops function f(int4[], text[][])"},
					{review.DisallowDropObject, 11, span(t, sql, "DROP FUNCTION g(t.c%TYPE)"), "drops function g(t.c%TYPE)"},
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
			name:    "identifiers that need quoting keep their quotes and their spaces",
			sql:     "DROP TABLE \"select\", \"user\", \"Order\", \"with  two\", \"order\", \"off\", \"text\"; TRUNCATE \"select\".\"user\"; DROP FUNCTION f(), s.p(); DROP PROCEDURE p(); DROP ROUTINE r(); DROP AGGREGATE a(*);\nSELECT 1 FROM t WHERE \"a  b\"\n =\n  NULL;\nSELECT 1 FROM t WHERE \"c\nd\" = NULL; DROP TABLE \"e\r\nf\";\nSELECT 1 FROM t WHERE $$a  b$$ = NULL OR $tag$c\n d$tag$   =  NULL OR 'e  f' = NULL OR $1 = NULL OR E'g\\'  h' = NULL OR 'i''  j' = NULL OR k -- note\n /* block\n comment */ = NULL OR l -- note\r = NULL;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 0, span(t, sql, `DROP TABLE "select", "user", "Order", "with  two", "order", "off", "text"`), `drops table "select", "user", "Order", "with  two", "order", off, text`},
					{review.DisallowTruncate, 1, span(t, sql, `TRUNCATE "select"."user"`), `truncates "select"."user"`},
					{review.DisallowDropObject, 2, span(t, sql, "DROP FUNCTION f(), s.p()"), "drops function f(), s.p()"},
					{review.DisallowDropObject, 3, span(t, sql, "DROP PROCEDURE p()"), "drops procedure p()"},
					{review.DisallowDropObject, 4, span(t, sql, "DROP ROUTINE r()"), "drops routine r()"},
					{review.DisallowDropObject, 5, span(t, sql, "DROP AGGREGATE a(*)"), "drops aggregate a(*)"},
					{review.RequireIsNull, 6, span(t, sql, "\"a  b\"\n =\n  NULL"), `""a  b" = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 7, span(t, sql, "\"c\nd\" = NULL"), `""c\nd" = NULL" is never true; use IS NULL`},
					{review.DisallowDropObject, 8, span(t, sql, "DROP TABLE \"e\r\nf\""), `drops table "e\r\nf"`},
					{review.RequireIsNull, 9, span(t, sql, "$$a  b$$ = NULL"), `"$$a  b$$ = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "$tag$c\n d$tag$   =  NULL"), `"$tag$c\n d$tag$ = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "'e  f' = NULL"), `"'e  f' = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "$1 = NULL"), `"$1 = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "E'g\\'  h' = NULL"), `"E'g\'  h' = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "'i''  j' = NULL"), `"'i''  j' = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "k -- note\n /* block\n comment */ = NULL"), `"k = NULL" is never true; use IS NULL`},
					{review.RequireIsNull, 9, span(t, sql, "l -- note\r = NULL"), `"l = NULL" is never true; use IS NULL`},
				}
			},
		},
		{
			name:    "ALTER TYPE DROP ATTRIBUTE anchors the subcommand",
			sql:     "SELECT 1; ALTER TYPE t DROP ATTRIBUTE a;\n  ALTER TYPE s.t ADD ATTRIBUTE b int, DROP ATTRIBUTE c CASCADE ;",
			targets: 1,
			want: func(t *testing.T, sql string) []finding {
				return []finding{
					{review.DisallowDropObject, 1, span(t, sql, "DROP ATTRIBUTE a"), "drops attribute a of type t"},
					{review.DisallowDropObject, 2, span(t, sql, "DROP ATTRIBUTE c CASCADE"), "drops attribute c of type s.t"},
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
			name:    "a line comment ends at a bare carriage return, as in the lexer",
			sql:     "SELECT 1;--c\rDELETE FROM t;",
			targets: 1,
			ranges:  []review.Range{{Start: 0, End: 9}, {Start: 9, End: 27}},
			want: func(t *testing.T, sql string) []finding {
				return []finding{{review.RequireWhere, 1, span(t, sql, "DELETE FROM t"), "DELETE FROM t has no WHERE clause"}}
			},
		},
		{
			name:    "form feed and vertical tab are whitespace",
			sql:     "DELETE FROM t WHERE a = 1;\f\v\n\f",
			targets: 1,
			ranges:  []review.Range{{Start: 0, End: 26}},
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

func TestUnparsedFinding(t *testing.T) {
	// The parser reports an unterminated construct itself now; the range
	// check behind this finding stays as a guard against text the parser
	// might drop, so its anchor and message are tested directly.
	sql := "SELECT 1;\n-- header\n\\restrict abc\n/* block\n */ 'unterminated\n/* open"
	ranges := statementRanges(sql)
	if len(ranges) != 2 {
		t.Fatalf("ranges = %v, want 2", ranges)
	}
	at := strings.Index(sql, "'unterminated")
	got := unparsedFinding(sql, 1, ranges[1])
	want := &review.Finding{Rule: review.Syntax, Statement: 1, Range: review.Range{Start: at, End: at}, Message: `syntax error at or near "'unterminated"`}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unparsedFinding mismatch (-want +got):\n%s", diff)
	}
	// An unterminated block comment is the text at fault, not trivia.
	open := strings.Index(sql, "/* open")
	if got := skipTrivia(sql, open, len(sql)); got != open {
		t.Errorf("skipTrivia past an open comment = %d, want %d", got, open)
	}
}

func TestSyntaxMessageIsOneLine(t *testing.T) {
	result, err := Review(context.Background(), "SELECT 1 AS x \"bad\nname\"", review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 || result.Findings[0].Rule != review.Syntax {
		t.Fatalf("Findings = %+v, want one syntax finding", result.Findings)
	}
	if msg := result.Findings[0].Message; strings.ContainsAny(msg, "\n\r") || !strings.Contains(msg, `\n`) {
		t.Fatalf("Message = %q, want the line break written as \\n", msg)
	}
}

func TestReviewRecoversFromPanics(t *testing.T) {
	// The parser accepts DROP AGGREGATE b() with a typed-nil argument that
	// makes ast.Walk panic; once it rejects the form like PostgreSQL, this
	// is a syntax finding. Either way Review returns, without panicking.
	result, err := Review(context.Background(), "DROP AGGREGATE b();", review.Options{}, nil)
	if err == nil && (len(result.Findings) != 1 || result.Findings[0].Rule != review.Syntax) {
		t.Fatalf("Review = %+v, %v; want an error or one syntax finding", result, err)
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
