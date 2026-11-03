package deparse

import "testing"

// TestNameText pins how an item's source text becomes its name: TiDB drops
// every "*/" and every "/*!" opener with its version, wherever they stand.
// Each want is the column name TiDB v8.5.5 gives a view over the item (see
// also the oracle cases in catalog's view_names_test.go).
func TestNameText(t *testing.T) {
	for _, tc := range []struct{ text, want string }{
		{"a + 1", "a + 1"},
		{"a /* c */ + 1", "a /* c  + 1"},
		{"b /*+ hint */ + 4", "b /*+ hint  + 4"},
		{"a -- note\n + 1", "a -- note\n + 1"},
		{"a/*!80000 +2", "a +2"},
		{"a /*!+3*/ + 1", "a +3 + 1"},
		{"a /**/ + 1", "a /* + 1"},
		// "*/*" is no exception: its "*/" goes too.
		{"a */*c*/3", "a *c3"},
		{"a*/*c*/b", "a*cb"},
		{"a */**/3", "a *3"},
		{"a /*! + 1 */ */*c*/3", "a  + 1  *c3"},
		{"a /*!50100 */*/*c*/3", "a  *c3"},
		// Nor is a string, a quoted identifier or a line comment.
		{"a + '*/'", "a + ''"},
		{"a + 'x*/y'", "a + 'xy'"},
		{"a + \"*/\"", "a + \"\""},
		{"a + 'it''s */'", "a + 'it''s '"},
		{"concat('*/', b)", "concat('', b)"},
		{"`x*/y` + 1", "`xy` + 1"},
		{"`p/*!q` + 1", "`pq` + 1"},
		{"a -- x */\n + 1", "a -- x \n + 1"},
		// The version is an optional M and five or six digits; anything else
		// after "/*!" stays.
		{"a + '/*!50100 x'", "a + ' x'"},
		{"a + '/*!1234 x'", "a + '1234 x'"},
		{"a + '/*!123456 x'", "a + ' x'"},
		{"a + '/*!1234567 x'", "a + '7 x'"},
		{"a + '/*!M12345 x'", "a + ' x'"},
		{"a + '/*!12345M x'", "a + 'M x'"},
		// One pass, left to right: what a drop leaves is not read again.
		{"a + '/**/!x'", "a + '/*!x'"},
		{"a + 'x*/*!50100y'", "a + 'x*!50100y'"},
		{"a + '**//'", "a + '*/'"},
		{"a + '/*!/**/x'", "a + '/*x'"},
		// Other openers are ordinary text.
		{"a /*T![clustered_index] + 1 */ + 2", "a /*T![clustered_index] + 1  + 2"},
		{"a /*M! + 1 */ + 2", "a /*M! + 1  + 2"},
		{"a + '/**!x'", "a + '/**!x'"},
	} {
		if got := nameText(tc.text); got != tc.want {
			t.Errorf("nameText(%q) = %q, want %q", tc.text, got, tc.want)
		}
	}
}

// TestItemText: TiDB names an item after its text through the comments that
// follow it, up to the next token.
func TestItemText(t *testing.T) {
	const sql = "SELECT /* a */ a+1, /* b */ b+1 /* c */ ,c+1 FROM t"
	stmt := parseSelectOrFatal(t, sql, sql)
	var got []string
	for _, target := range stmt.TargetList {
		text, ok := itemText(sql, target)
		if !ok {
			t.Fatalf("no text for %T", target)
		}
		got = append(got, text)
	}
	want := []string{"a+1", "b+1 /* c */", "c+1"}
	for i := range want {
		if i >= len(got) || got[i] != want[i] {
			t.Fatalf("item texts %q, want %q", got, want)
		}
	}
}

// TestAutoAliasWithoutSource: with no source text, DeparseSelect rebuilds an
// unaliased subquery's name from the AST. It used to keep only the select
// list, FROM and WHERE of a plain query, so a WITH, a set operation or a
// trailing clause vanished from the name; and writing a set operation cleared
// its WITH from the AST, which dropped it from the definition too.
func TestAutoAliasWithoutSource(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{
			"SELECT (WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)",
			"select (with `c` as (select 1 AS `x`) select `x` from `c` union select 2) AS `(WITH c AS (SELECT 1 AS x) SELECT x FROM c UNION SELECT 2)`",
		},
		{
			"SELECT (SELECT a FROM t GROUP BY a HAVING a > 0 LIMIT 1), (SELECT a FROM t ORDER BY a DESC LIMIT 2, 1)",
			"select (select `a` from `t` group by `a` having (`a` > 0) limit 1) AS `(SELECT a FROM t GROUP BY a HAVING a > 0 LIMIT 1)`,(select `a` from `t` order by `a` desc limit 2,1) AS `(SELECT a FROM t ORDER BY a DESC LIMIT 2, 1)`",
		},
		{
			"SELECT EXISTS (SELECT DISTINCT a FROM t UNION ALL SELECT 1)",
			"select exists(select distinct `a` from `t` union all select 1) AS `EXISTS(SELECT DISTINCT a FROM t UNION ALL SELECT 1)`",
		},
	} {
		stmt := parseSelectOrFatal(t, tc.sql, tc.sql)
		if got := DeparseSelect(stmt); got != tc.want {
			t.Errorf("DeparseSelect(%s)\n got: %s\nwant: %s", tc.sql, got, tc.want)
		}
		if got := DeparseSelect(stmt); got != tc.want {
			t.Errorf("DeparseSelect(%s) again\n got: %s\nwant: %s", tc.sql, got, tc.want)
		}
	}
}
