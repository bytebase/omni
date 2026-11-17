package parser

import (
	"strings"
	"testing"
)

// ref is the part of a RangeReference a test states: the kind and the names the
// parser reads off the text. Columns are checked separately because only a body
// the parser can read statically has them.
type ref struct {
	kind     RangeReferenceKind
	database string
	table    string
	alias    string
}

type wantSignals struct {
	prefix    string
	kinds     []ObjectKind
	qualifier []string
	local     []ref
	outer     [][]ref
	ctes      map[string][]string
	aliases   []string
	dml       string
	// cols are the statically readable columns of a local reference, keyed
	// by the name it is referred to by.
	cols map[string][]string
}

func kindNames(kinds []ObjectKind) string {
	names := make([]string, 0, len(kinds))
	for _, k := range kinds {
		names = append(names, objectKindName(k))
	}
	return strings.Join(names, ",")
}

func objectKindName(k ObjectKind) string {
	switch k {
	case ObjectKindDatabase:
		return "database"
	case ObjectKindTable:
		return "table"
	case ObjectKindView:
		return "view"
	case ObjectKindColumn:
		return "column"
	case ObjectKindIndex:
		return "index"
	case ObjectKindFunction:
		return "function"
	case ObjectKindProcedure:
		return "procedure"
	case ObjectKindTrigger:
		return "trigger"
	case ObjectKindEvent:
		return "event"
	case ObjectKindType:
		return "type"
	case ObjectKindCharset:
		return "charset"
	case ObjectKindEngine:
		return "engine"
	case ObjectKindVariable:
		return "variable"
	case ObjectKindSelectAlias:
		return "selectalias"
	default:
		return "unknown"
	}
}

func refNames(refs []RangeReference) string {
	names := make([]string, 0, len(refs))
	for _, r := range refs {
		name := r.Table
		if r.Database != "" {
			name = r.Database + "." + name
		}
		if r.Alias != "" {
			name += " " + r.Alias
		}
		if r.Kind == RangeReferenceDMLTarget {
			name = "->" + name
		}
		names = append(names, name)
	}
	return strings.Join(names, ",")
}

// TestCollectCompletionSignals holds the MySQL signal set to a corpus of
// scripts x cursor positions: the prefix, the object kinds and qualifier the
// position implies, the relation scope with its aliases and derived scopes, the
// CTE definitions with their columns, the select aliases, and the DML target.
func TestCollectCompletionSignals(t *testing.T) {
	// The cursor is the "|" in sql; it is removed before the call so a case
	// reads as the editor shows it.
	cases := []struct {
		name string
		sql  string
		want wantSignals
	}{
		{
			name: "select list",
			sql:  "SELECT | FROM t",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "after from",
			sql:  "SELECT a FROM |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView}},
		},
		{
			name: "where column",
			sql:  "SELECT * FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "table qualified column",
			sql:  "SELECT a.| FROM t a",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"a"}, local: []ref{{table: "t", alias: "a"}}},
		},
		{
			name: "database qualified table",
			sql:  "SELECT * FROM db.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView}, qualifier: []string{"db"}},
		},
		{
			name: "database table column",
			sql:  "SELECT db.t.| FROM db.t",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindColumn}, qualifier: []string{"db", "t"}, local: []ref{{database: "db", table: "t"}}},
		},
		{
			name: "joined aliases",
			sql:  "SELECT * FROM t1 a JOIN t2 AS b ON a.id = b.id WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1", alias: "a"}, {table: "t2", alias: "b"}}},
		},
		{
			name: "cross join without on",
			sql:  "SELECT * FROM t1, t2 WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1"}, {table: "t2"}}},
		},
		{
			name: "straight join",
			sql:  "SELECT * FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1"}, {table: "t2"}}},
		},
		{
			name: "derived table columns from select list",
			sql:  "SELECT * FROM (SELECT a, b AS c FROM t) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}},
		},
		{
			name: "derived table alias columns",
			sql:  "SELECT * FROM (SELECT a, b FROM t) AS d (x, y) WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}},
		},
		{
			name: "unterminated derived table body",
			sql:  "SELECT * FROM (SELECT a FROM |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView}},
		},
		{
			name: "cte from select list",
			sql:  "WITH c AS (SELECT id, name FROM t) SELECT * FROM c WHERE c.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"id", "name"}}},
		},
		{
			name: "cte with declared columns",
			sql:  "WITH c (x, y) AS (SELECT id, name FROM t) SELECT | FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"x", "y"}}},
		},
		{
			name: "recursive cte",
			sql:  "WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT * FROM r WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceCTE, table: "r"}}, ctes: map[string][]string{"r": {"n"}}},
		},
		{
			name: "two ctes second sees first",
			sql:  "WITH a AS (SELECT x FROM t), b AS (SELECT x FROM a) SELECT * FROM b WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceCTE, table: "b"}}, ctes: map[string][]string{"a": {"x"}, "b": {"x"}}},
		},
		{
			name: "cte body cursor",
			sql:  "WITH c AS (SELECT id FROM t WHERE |) SELECT * FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}, ctes: map[string][]string{"c": {"id"}}},
		},
		{
			name: "correlated subquery sees outer block",
			sql:  "SELECT * FROM t WHERE a IN (SELECT b FROM u WHERE |)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}, outer: [][]ref{{{table: "t"}}}},
		},
		{
			name: "derived table body does not see its siblings",
			sql:  "SELECT * FROM (SELECT a FROM t, u WHERE |) d",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}, {table: "u"}}},
		},
		{
			name: "lateral derived table sees earlier reference",
			sql:  "SELECT * FROM t, LATERAL (SELECT a FROM u WHERE u.x = t.x AND |) d",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}, outer: [][]ref{{{table: "t"}}}},
		},
		{
			name: "update target",
			sql:  "UPDATE t SET | WHERE id = 1",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "update join target",
			sql:  "UPDATE t1 a JOIN t2 b ON a.id = b.id SET a.|",
			want: wantSignals{prefix: "", kinds: nil, qualifier: []string{"a"}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t1", alias: "a"}, {table: "t2", alias: "b"}}, dml: "->t1 a"},
		},
		{
			name: "delete target",
			sql:  "DELETE FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "insert target",
			sql:  "INSERT INTO db.t (a) VALUES (|)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, database: "db", table: "t"}}, dml: "->db.t"},
		},
		{
			name: "insert select keeps the query scope",
			sql:  "INSERT INTO t SELECT * FROM u WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}, dml: "->t"},
		},
		{
			name: "order by offers select aliases",
			sql:  "SELECT a AS x, b + 1 AS y FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"x", "y"}},
		},
		{
			name: "group by offers select aliases",
			sql:  "SELECT a AS x, COUNT(*) AS n FROM t GROUP BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"x", "n"}},
		},
		{
			name: "where does not offer select aliases",
			sql:  "SELECT a AS x FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}, aliases: []string{"x"}},
		},
		{
			name: "select alias without as",
			sql:  "SELECT a x, b FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"x"}},
		},
		{
			name: "prefix",
			sql:  "SELECT * FROM us|",
			want: wantSignals{prefix: "us", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView}, local: []ref{{table: "us"}}},
		},
		{
			name: "quoted prefix",
			sql:  "SELECT * FROM `us|",
			want: wantSignals{prefix: "us", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView}, local: []ref{{table: "us"}}},
		},
		{
			name: "quoted reference",
			sql:  "SELECT * FROM `my table` WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "my table"}}},
		},
		{
			name: "index hint is not an alias",
			sql:  "SELECT * FROM t USE INDEX (i) WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "partition clause is not an alias",
			sql:  "SELECT * FROM t PARTITION (p0) WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "no leak between statements",
			sql:  "SELECT * FROM earlier; SELECT * FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "second statement after cursor of first",
			sql:  "SELECT | FROM t; DELETE FROM u",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "union arms keep their own scope",
			sql:  "SELECT a FROM t UNION SELECT b FROM u WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}},
		},
		{
			name: "set-operation arms do not see each other",
			sql:  "(SELECT c1 FROM t1) UNION (SELECT | FROM t2)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t2"}}},
		},
		{
			name: "aliases come from the cursor block only",
			sql:  "SELECT 1 AS outer_alias FROM (SELECT a AS inner_alias FROM t ORDER BY |) d",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"inner_alias"}},
		},
		{
			name: "parenthesized join contributes its members",
			sql:  "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1"}, {table: "t2"}}},
		},
		{
			name: "parenthesized join alias hides its members",
			// The grammar does not accept an alias after a parenthesized
			// join yet — MySQL does — so no rule reaches the caret. The
			// token reader still names the reference.
			sql:  "SELECT * FROM (t1 JOIN t2 ON t1.a = t2.a) j WHERE |",
			want: wantSignals{prefix: "", local: []ref{{kind: RangeReferenceJoinAlias, table: "j", alias: "j"}}},
		},
		{
			name: "qualified select item alias",
			sql:  "SELECT a.b c FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"c"}},
		},
		{
			name: "select modifier is not an alias",
			sql:  "SELECT DISTINCT a FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}},
		},
		{
			name: "derived table columns from qualified items",
			sql:  "SELECT * FROM (SELECT a.b AS c, d FROM t) x WHERE x.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"x"}, local: []ref{{kind: RangeReferenceSubquery, table: "x", alias: "x"}}},
		},
		{
			name: "ddl type name",
			sql:  "CREATE TABLE t (a |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindType}},
		},
		{
			name: "index name position",
			sql:  "DROP INDEX | ON t",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindIndex}},
		},
		{
			name: "engine name position",
			sql:  "CREATE TABLE t (a INT) ENGINE = |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindEngine}},
		},
		{
			name: "charset position",
			// The table options accept a charset, but Collect offers no rule
			// here yet: the position keeps its qualifier and no kind.
			sql:  "CREATE TABLE t (a INT) CHARACTER SET |",
			want: wantSignals{prefix: "", kinds: nil},
		},
		{
			name: "view position",
			sql:  "DROP VIEW |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindView}},
		},
		{
			name: "database position",
			sql:  "DROP DATABASE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase}},
		},
		{
			name: "unaliased nested groups keep their members",
			sql:  "SELECT * FROM ((t1 JOIN t2 ON t1.a = t2.a)) WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1"}, {table: "t2"}}},
		},
		{
			name: "inner subquery does not see a sibling's FROM",
			sql:  "SELECT (SELECT |), (SELECT a FROM u) FROM t",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, outer: [][]ref{{{table: "t"}}}},
		},
		{
			name: "correlated references survive an empty local scope",
			sql:  "SELECT * FROM t WHERE EXISTS (SELECT |)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, outer: [][]ref{{{table: "t"}}}},
		},
		{
			name: "insert select has no local source",
			sql:  "INSERT INTO dst (x) SELECT |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, dml: "->dst"},
		},
		{
			name: "json table columns",
			sql:  "SELECT jt.| FROM JSON_TABLE('[1,2]', '$[*]' COLUMNS (x INT PATH '$')) AS jt",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"jt"}, local: []ref{{kind: RangeReferenceFunction, table: "jt", alias: "jt"}}, cols: map[string][]string{"jt": {"x"}}},
		},
		{
			name: "json table nested columns",
			sql:  "SELECT jt.| FROM JSON_TABLE('[1]', '$' COLUMNS (a INT PATH '$.a', NESTED PATH '$.n[*]' COLUMNS (b INT PATH '$'))) AS jt",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"jt"}, local: []ref{{kind: RangeReferenceFunction, table: "jt", alias: "jt"}}, cols: map[string][]string{"jt": {"a", "b"}}},
		},
		{
			name: "delete target list keeps the join scope",
			sql:  "DELETE a FROM t1 a JOIN t2 b USING (id) WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "a"}, {table: "t1", alias: "a"}, {table: "t2", alias: "b"}}, dml: "->a"},
		},
		{
			name: "delete using form",
			sql:  "DELETE FROM t1 USING t1 JOIN t2 ON t1.a = t2.a WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t1"}, {table: "t1"}, {table: "t2"}}, dml: "->t1"},
		},
		{
			name: "closed subquery does not make an outer where an alias position",
			sql:  "SELECT x AS total FROM t WHERE id IN (SELECT id FROM u ORDER BY id) AND |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}, aliases: []string{"total"}},
		},
		{
			name: "arithmetic item keeps its alias",
			sql:  "WITH c AS (SELECT price * qty AS total FROM items) SELECT * FROM c WHERE c.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"total"}}, cols: map[string][]string{"c": {"total"}}},
		},
		{
			name: "function item keeps its alias",
			sql:  "SELECT SUM(a) total FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"total"}},
		},
		{
			name: "nested with shadows the statement's cte",
			sql:  "WITH c(x) AS (SELECT a FROM t) SELECT * FROM (WITH c(y) AS (SELECT b FROM u) SELECT * FROM c WHERE c.|) d",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"y"}}},
		},
		{
			name: "quoted prefix with a space",
			sql:  "SELECT * FROM `my ta|",
			want: wantSignals{prefix: "my ta", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView}, local: []ref{{table: "my ta"}}},
		},
		{
			name: "unicode prefix",
			sql:  "SELECT 名|",
			want: wantSignals{prefix: "名", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}},
		},
		{
			name: "inside a text literal",
			sql:  "SELECT '|' FROM t",
			want: wantSignals{local: []ref{{table: "t"}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cursor := strings.Index(tc.sql, "|")
			if cursor < 0 {
				t.Fatalf("case has no cursor marker")
			}
			sql := tc.sql[:cursor] + tc.sql[cursor+1:]
			ctx := CollectCompletion(sql, cursor)
			if ctx.Prefix != tc.want.prefix {
				t.Errorf("prefix = %q, want %q", ctx.Prefix, tc.want.prefix)
			}
			if got := kindNames(intentKinds(ctx)); got != kindNames(tc.want.kinds) {
				t.Errorf("object kinds = %s, want %s", got, kindNames(tc.want.kinds))
			}
			if got, want := strings.Join(intentQualifier(ctx), "."), strings.Join(tc.want.qualifier, "."); got != want {
				t.Errorf("qualifier = %q, want %q", got, want)
			}
			if got := refNames(ctx.Scope.LocalReferences); got != refNames(wantRefs(tc.want.local)) {
				t.Errorf("local references = %q, want %q", got, refNames(wantRefs(tc.want.local)))
			}
			for _, ref := range ctx.Scope.LocalReferences {
				name := ref.Table
				if name == "" {
					name = ref.Alias
				}
				want, checked := tc.want.cols[name]
				if !checked {
					continue
				}
				if got := strings.Join(ref.Columns, ","); got != strings.Join(want, ",") {
					t.Errorf("local reference %q columns = %q, want %q", name, got, want)
				}
			}
			if len(ctx.Scope.OuterReferences) != len(tc.want.outer) {
				t.Errorf("outer levels = %d, want %d", len(ctx.Scope.OuterReferences), len(tc.want.outer))
			} else {
				for i := range tc.want.outer {
					if got, want := refNames(ctx.Scope.OuterReferences[i]), refNames(wantRefs(tc.want.outer[i])); got != want {
						t.Errorf("outer[%d] = %q, want %q", i, got, want)
					}
				}
			}
			if len(tc.want.ctes) != len(ctx.CTEs) {
				t.Errorf("ctes = %d, want %d", len(ctx.CTEs), len(tc.want.ctes))
			}
			for _, cte := range ctx.CTEs {
				want, ok := tc.want.ctes[cte.Table]
				if !ok {
					t.Errorf("unexpected cte %q", cte.Table)
					continue
				}
				if got := strings.Join(cte.Columns, ","); got != strings.Join(want, ",") {
					t.Errorf("cte %q columns = %q, want %q", cte.Table, got, strings.Join(want, ","))
				}
			}
			var aliases []string
			for _, a := range ctx.SelectAliases {
				aliases = append(aliases, a.Name)
			}
			if got := strings.Join(aliases, ","); got != strings.Join(tc.want.aliases, ",") {
				t.Errorf("select aliases = %q, want %q", got, strings.Join(tc.want.aliases, ","))
			}
			if tc.want.dml == "" {
				if ctx.Scope.DMLTarget != nil {
					t.Errorf("dml target = %q, want none", refNames([]RangeReference{*ctx.Scope.DMLTarget}))
				}
			} else if ctx.Scope.DMLTarget == nil {
				t.Errorf("dml target = none, want %q", tc.want.dml)
			} else if got := refNames([]RangeReference{*ctx.Scope.DMLTarget}); got != tc.want.dml {
				t.Errorf("dml target = %q, want %q", got, tc.want.dml)
			}
			if ctx.Candidates == nil {
				t.Error("candidates = nil, want the grammar's own candidate set")
			}
			if ctx.StatementLoc.Start < 0 || ctx.StatementLoc.End > len(sql) || ctx.StatementLoc.Start > cursor {
				t.Errorf("statement loc = %v, want a range holding the cursor in %d bytes", ctx.StatementLoc, len(sql))
			}
		})
	}
}

func intentKinds(ctx *CompletionContext) []ObjectKind {
	if ctx.Intent == nil {
		return nil
	}
	return ctx.Intent.ObjectKinds
}

func intentQualifier(ctx *CompletionContext) []string {
	if ctx.Intent == nil {
		return nil
	}
	return ctx.Intent.Qualifier.Parts
}

func wantRefs(in []ref) []RangeReference {
	refs := make([]RangeReference, 0, len(in))
	for _, r := range in {
		refs = append(refs, RangeReference{Kind: r.kind, Database: r.database, Table: r.table, Alias: r.alias})
	}
	return refs
}

// TestCollectCompletionCandidates pins that the grammar's own candidate set
// travels with the signals, since a caller still needs the keywords and the
// rule names it reads.
func TestCollectCompletionCandidates(t *testing.T) {
	sql := "SELECT | FROM t"
	ctx := CollectCompletion(sql, strings.Index(sql, "|"))
	if ctx.Candidates == nil || !ctx.Candidates.HasRule("columnref") || !ctx.Candidates.HasRule("func_name") {
		t.Fatalf("rules = %v, want columnref and func_name", ctx.Candidates.Rules)
	}
	if len(ctx.Candidates.Tokens) == 0 {
		t.Fatal("tokens = none, want the keywords the position allows")
	}
	for _, tok := range ctx.Candidates.Tokens {
		if TokenName(tok) == "" {
			t.Fatalf("token %d has no keyword name", tok)
		}
	}
}

// TestCollectCompletionNestingLimit pins that a group nested past the
// parser's own limit is skipped instead of descended into: completion runs on
// editor input, so it must not walk a group the grammar already rejects.
func TestCollectCompletionNestingLimit(t *testing.T) {
	within := "SELECT * FROM " + strings.Repeat("(", 8) + "t" + strings.Repeat(")", 8) + " WHERE |"
	if ctx := CollectCompletion(strings.Replace(within, "|", "", 1), strings.Index(within, "|")); len(ctx.Scope.LocalReferences) != 1 || ctx.Scope.LocalReferences[0].Table != "t" {
		t.Fatalf("local references = %q, want t", refNames(ctx.Scope.LocalReferences))
	}
	beyond := "SELECT * FROM " + strings.Repeat("(", maxNestingDepth+8) + "t" + strings.Repeat(")", maxNestingDepth+8) + " WHERE |"
	ctx := CollectCompletion(strings.Replace(beyond, "|", "", 1), strings.Index(beyond, "|"))
	if len(ctx.Scope.LocalReferences) != 0 {
		t.Fatalf("local references = %q, want none past the nesting limit", refNames(ctx.Scope.LocalReferences))
	}
}

// TestCollectCompletionEveryOffset checks that the reader survives every cursor
// position of an incomplete script: completion runs on text that is still being
// typed, so a panic here kills the editor's request.
func TestCollectCompletionEveryOffset(t *testing.T) {
	scripts := []string{
		"WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM r) SELECT * FROM r",
		"SELECT * FROM (SELECT a FROM t1, t2 LEFT JOIN t3 ON t2.a = t3.a) d WHERE d.a IN (SELECT b FROM u)",
		"UPDATE t1 a JOIN t2 b ON a.id = b.id SET a.x = 1 WHERE a.y = 2",
		"DELETE t1, t2 FROM t1 JOIN t2 ON t1.a = t2.a WHERE t1.b = 1",
		"INSERT INTO t (a, b) VALUES (1, 2) ON DUPLICATE KEY UPDATE a = 3",
		"CREATE TABLE t (a INT, b VARCHAR(10), PRIMARY KEY (a), KEY i (b)) ENGINE = InnoDB",
		"SELECT 'unterminated FROM t WHERE `x",
		strings.Repeat("(", 200) + "SELECT 1",
	}
	for _, script := range scripts {
		for offset := 0; offset <= len(script); offset++ {
			CollectCompletion(script, offset)
		}
	}
}
