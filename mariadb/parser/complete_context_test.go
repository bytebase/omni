package parser

import (
	"strings"
	"testing"
	"time"
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
	case ObjectKindConstraint:
		return "constraint"
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

// TestCollectCompletionSignals holds MariaDB's signal set to a corpus of
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
			name: "later cte is not in scope inside an earlier body",
			sql:  "WITH c AS (SELECT * FROM real_table WHERE |), real_table(y) AS (SELECT 1) SELECT * FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "real_table"}}, ctes: map[string][]string{"c": nil, "real_table": {"y"}}},
		},
		{
			name: "recursive cte names itself",
			sql:  "WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n FROM r WHERE |) SELECT * FROM r",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceCTE, table: "r"}}, ctes: map[string][]string{"r": {"n"}}},
		},
		{
			name: "explain update",
			sql:  "EXPLAIN UPDATE t SET x = |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "explain options before delete",
			sql:  "EXPLAIN ANALYZE FORMAT = JSON DELETE FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "join condition sees only its operands",
			sql:  "SELECT * FROM t1 JOIN t2 ON | JOIN t3 ON t2.a = t3.a",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "t1"}, {table: "t2"}}},
		},
		{
			name: "select list sees every joined source",
			sql:  "SELECT | FROM t1 JOIN t2 ON t1.a = t2.a JOIN t3 ON t2.a = t3.a",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1"}, {table: "t2"}, {table: "t3"}}},
		},
		{
			name: "custom delimiter separates statements",
			sql:  "DELIMITER $$\nSELECT * FROM old_table $$\nUPDATE current_table SET x = | $$",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "current_table"}}, dml: "->current_table"},
		},
		{
			name: "second foreign key reference",
			sql:  "ALTER TABLE child ADD FOREIGN KEY (a) REFERENCES p(a), ADD FOREIGN KEY (b) REFERENCES q(|)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "q", alias: "q"}}},
		},
		{
			name: "index hint for group by",
			sql:  "SELECT b.| FROM t1 a USE INDEX FOR GROUP BY (i) JOIN t2 b ON a.id = b.id",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"b"}, local: []ref{{table: "t1", alias: "a"}, {table: "t2", alias: "b"}}},
		},
		{
			name: "odbc outer join escape",
			sql:  "SELECT a.| FROM { OJ t1 a LEFT JOIN t2 b ON a.id = b.id }",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"a"}, local: []ref{{table: "t1", alias: "a"}, {table: "t2", alias: "b"}}},
		},
		{
			name: "cursor at a custom terminator",
			sql:  "DELIMITER $$\nSELECT 1$$|",
			want: wantSignals{},
		},
		// MariaDB rejects the INSERT row alias ("not supported in MariaDB"),
		// so these cases are MySQL's and OceanBase's, not this corpus's.
		{
			name: "update inside a routine body",
			sql:  "CREATE PROCEDURE p() BEGIN UPDATE t SET x = |; END",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "delete after a statement in a routine body",
			sql:  "CREATE PROCEDURE p() BEGIN SELECT 1; DELETE FROM d WHERE x = |; END",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "d"}}, dml: "->d"},
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
			name: "second statement keeps its own parentheses",
			sql:  "SELECT 1; SELECT * FROM (SELECT a FROM t) d WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}},
		},
		{
			name: "delete target list with using",
			sql:  "DELETE FROM t1, t2 USING t1 JOIN t2 ON t1.x = t2.x WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t1"}, {kind: RangeReferenceDMLTarget, table: "t2"}, {table: "t1"}, {table: "t2"}}, dml: "->t1"},
		},
		{
			name: "set operation order by names the first arm",
			sql:  "SELECT 1 AS first_name UNION SELECT 2 AS second_name ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, aliases: []string{"first_name"}},
		},
		// MariaDB rejects `TABLE t`: MySQL 8.0.19's TABLE query primary has no
		// MariaDB counterpart, so the reader's TABLE frame has nothing to read.
		{
			name: "delete target with star",
			sql:  "DELETE FROM t1.* USING t1 JOIN t2 ON t1.x = t2.x WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t1"}, {table: "t1"}, {table: "t2"}}, dml: "->t1"},
		},
		{
			name: "delete modifier before from",
			sql:  "DELETE LOW_PRIORITY FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "single table delete alias",
			sql:  "DELETE FROM db.t AS a WHERE a.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"a"}, local: []ref{{kind: RangeReferenceDMLTarget, database: "db", table: "t", alias: "a"}}, dml: "->db.t a"},
		},
		{
			name: "extra parentheses around a derived table",
			sql:  "SELECT * FROM ((SELECT a AS x FROM t)) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}},
		},
		{
			name: "foreign key references its parent",
			sql:  "ALTER TABLE child ADD FOREIGN KEY (parent_id) REFERENCES db.parent (|)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{database: "db", table: "parent", alias: "parent"}}},
		},
		{
			name: "foreign key names its own columns from the altered table",
			sql:  "ALTER TABLE child ADD FOREIGN KEY (|) REFERENCES db.parent (id)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "child", alias: "child"}}},
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
			name: "constraint name position",
			sql:  "ALTER TABLE t DROP FOREIGN KEY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindConstraint}, local: []ref{{table: "t", alias: "t"}}},
		},
		{
			name: "dropped constraint name position",
			sql:  "ALTER TABLE t DROP CONSTRAINT |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindConstraint}, local: []ref{{table: "t", alias: "t"}}},
		},
		{
			name: "index name position",
			sql:  "DROP INDEX | ON t",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindIndex}, local: []ref{{table: "t", alias: "t"}}},
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
			name: "expression operand is not an implicit alias",
			sql:  "SELECT * FROM (SELECT a + b FROM t) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}, cols: map[string][]string{"d": {}}},
		},
		{
			name: "expression operand is not a select alias",
			sql:  "SELECT a + b FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}},
		},
		{
			name: "order by after a from-less select list",
			sql:  "SELECT 1 AS total ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, aliases: []string{"total"}},
		},
		{
			name: "update subquery sees the update target",
			sql:  "UPDATE t SET x = (SELECT a FROM u WHERE |)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}, outer: [][]ref{{{table: "t"}}}, dml: "->t"},
		},
		{
			name: "delete subquery sees the delete target",
			sql:  "DELETE FROM t WHERE id IN (SELECT a FROM u WHERE |)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}, outer: [][]ref{{{table: "t"}}}, dml: "->t"},
		},
		{
			name: "insert select subquery does not see the write target",
			sql:  "INSERT INTO t1 SELECT (SELECT a FROM u WHERE |) FROM t2",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "u"}}, outer: [][]ref{{{table: "t2"}}}, dml: "->t1"},
		},
		{
			name: "scoped index hint keeps the join after it",
			sql:  "SELECT b.| FROM t1 a USE INDEX FOR JOIN (i) JOIN t2 b ON a.id = b.id",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"b"}, local: []ref{{table: "t1", alias: "a"}, {table: "t2", alias: "b"}}},
		},
		{
			name: "partition selection before the alias",
			sql:  "SELECT a.| FROM t PARTITION (p0) AS a",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"a"}, local: []ref{{table: "t", alias: "a"}}},
		},
		{
			name: "update modifiers before the target",
			sql:  "UPDATE LOW_PRIORITY IGNORE t SET x = |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "function arguments keep the order by clause",
			sql:  "SELECT a AS total FROM t ORDER BY COALESCE(a, |)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"total"}},
		},
		{
			name: "case expression then is not a statement boundary",
			sql:  "UPDATE t SET x = CASE WHEN y = 1 THEN | ELSE 0 END",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "insert select duplicate key update sees the destination",
			sql:  "INSERT INTO t (a) SELECT a FROM s ON DUPLICATE KEY UPDATE a = |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "s"}}, outer: [][]ref{{{table: "t"}}}, dml: "->t"},
		},
		{
			name: "distinctrow is a select option",
			sql:  "SELECT * FROM (SELECT DISTINCTROW a FROM t) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}, cols: map[string][]string{"d": {"a"}}},
		},
		{
			name: "interval unit is not an implicit alias",
			sql:  "SELECT * FROM (SELECT NOW() + INTERVAL 1 DAY FROM t) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}, cols: map[string][]string{"d": {}}},
		},
		{
			name: "interval expression keeps its explicit alias",
			sql:  "SELECT NOW() + INTERVAL 1 DAY AS x FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"x"}},
		},
		{
			name: "delete star target before from",
			sql:  "DELETE a.* FROM db.t a WHERE a.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"a"}, local: []ref{{kind: RangeReferenceDMLTarget, table: "a"}, {database: "db", table: "t", alias: "a"}}, dml: "->a"},
		},
		{
			name: "adjacent text literals concatenate",
			sql:  "SELECT * FROM (SELECT 'a' 'b') d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}, cols: map[string][]string{"d": {}}},
		},
		{
			name: "explicit string alias",
			sql:  "SELECT 'a' AS 'b' FROM t ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "t"}}, aliases: []string{"b"}},
		},
		{
			name: "alias after an interval expression",
			sql:  "SELECT * FROM (SELECT NOW() + INTERVAL 1 DAY next_day FROM t) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}, cols: map[string][]string{"d": {"next_day"}}},
		},
		{
			name: "case expression inside a routine body",
			sql:  "CREATE PROCEDURE p() BEGIN UPDATE t SET x = CASE WHEN y = 1 THEN | ELSE 0 END; END",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "parenthesized set operation names the first arm",
			sql:  "(SELECT a AS x FROM t) UNION (SELECT b AS y FROM u) ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, aliases: []string{"x"}},
		},
		{
			name: "comma item is not part of a join condition",
			sql:  "SELECT * FROM t1, t2 JOIN t3 ON |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "t2"}, {table: "t3"}}},
		},
		{
			name: "nested join condition sees its operands",
			sql:  "SELECT * FROM t1 JOIN (t2 JOIN t3 ON |) ON t1.x = 2",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "t2"}, {table: "t3"}}},
		},
		{
			name: "wrapped with body names its output after the clause",
			sql:  "WITH c AS ((WITH d AS (SELECT 1 AS wrong) SELECT 2 AS actual)) SELECT c.| FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"actual"}}, cols: map[string][]string{"c": {"actual"}}},
		},
		{
			name: "unterminated derived table hides its siblings",
			sql:  "SELECT * FROM t, (SELECT |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}},
		},
		// MariaDB rejects the INSERT row alias ("not supported in MariaDB"),
		// so these cases are MySQL's and OceanBase's, not this corpus's.
		{
			name: "chained join keeps the whole left operand",
			sql:  "SELECT * FROM t1 JOIN t2 ON t1.a = t2.a JOIN t3 ON t1.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"t1"}, local: []ref{{table: "t1"}, {table: "t2"}, {table: "t3"}}},
		},
		{
			name: "recursive cte with an unfinished body",
			sql:  "WITH RECURSIVE c(x) AS (SELECT 1 UNION ALL SELECT c.| FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"x"}}, cols: map[string][]string{"c": {"x"}}},
		},
		{
			name: "unfinished references column list names the parent",
			sql:  "ALTER TABLE child ADD FOREIGN KEY (a) REFERENCES parent (|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "parent", alias: "parent"}}},
		},
		{
			name: "trigger row aliases name the trigger table",
			sql:  "CREATE TRIGGER tr BEFORE INSERT ON db.t FOR EACH ROW SET NEW.a = NEW.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"NEW"}, local: []ref{{database: "db", table: "t", alias: "NEW"}}},
		},
		{
			name: "correlated on condition keeps the outer frame",
			sql:  "SELECT * FROM outer_table o WHERE EXISTS (SELECT * FROM t1 a JOIN t2 b ON a.x = o.|)",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"o"}, local: []ref{{table: "t1", alias: "a"}, {table: "t2", alias: "b"}}, outer: [][]ref{{{table: "outer_table", alias: "o"}}}},
		},
		{
			name: "procedural case branch starts its own statement",
			sql:  "CREATE PROCEDURE p() BEGIN CASE WHEN 1 THEN UPDATE t SET x = |; END CASE; END",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "values body names generated columns",
			sql:  "WITH c AS (VALUES ROW((SELECT 1 AS wrong), 2)) SELECT c.| FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"column_0", "column_1"}}, cols: map[string][]string{"c": {"column_0", "column_1"}}},
		},
		{
			name: "definer clause before the trigger object",
			sql:  "CREATE DEFINER = 'root'@'localhost' TRIGGER tr BEFORE INSERT ON db.t FOR EACH ROW SET NEW.a = NEW.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"NEW"}, local: []ref{{database: "db", table: "t", alias: "NEW"}}},
		},
		{
			name: "set operation after a cte names its arm",
			sql:  "WITH z AS (SELECT 1 AS wrong) (SELECT a AS actual FROM t) UNION (SELECT b FROM u) ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, ctes: map[string][]string{"z": {"wrong"}}, aliases: []string{"actual"}},
		},
		{
			name: "on condition through expression parentheses",
			sql:  "SELECT * FROM t1 a JOIN t2 b ON (c.|) JOIN t3 c ON b.x = c.x",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{table: "t1", alias: "a"}, {table: "t2", alias: "b"}}},
		},
		{
			name: "cursor before an unfinished references list",
			sql:  "ALTER TABLE child ADD FOREIGN KEY (|) REFERENCES parent (a",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn}, local: []ref{{table: "child", alias: "child"}}},
		},
		{
			name: "case terminator is not an implicit alias",
			sql:  "SELECT * FROM (SELECT CASE WHEN a = 1 THEN b ELSE c END FROM t) d WHERE d.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"d"}, local: []ref{{kind: RangeReferenceSubquery, table: "d", alias: "d"}}, cols: map[string][]string{"d": {}}},
		},
		{
			name: "intersect arms keep their own scope",
			sql:  "SELECT a FROM t1 INTERSECT SELECT b FROM t2 WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t2"}}},
		},
		{
			name: "except arms keep their own scope",
			sql:  "SELECT a FROM t1 EXCEPT SELECT b FROM t2 WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t2"}}},
		},
		{
			name: "minus is a function, not a set operator",
			sql:  "SELECT MINUS(a) FROM t WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t"}}},
		},
		{
			name: "delete returning",
			sql:  "DELETE FROM t RETURNING |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "insert returning",
			sql:  "INSERT INTO t (a) VALUES (1) RETURNING |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "replace returning",
			sql:  "REPLACE INTO t (a) VALUES (1) RETURNING |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{kind: RangeReferenceDMLTarget, table: "t"}}, dml: "->t"},
		},
		{
			name: "inside a text literal",
			sql:  "SELECT '|' FROM t",
			want: wantSignals{local: []ref{{table: "t"}}},
		},
		{
			name: "join condition ends at the next clause",
			sql:  "SELECT * FROM t1, t2 JOIN t3 ON t2.a = t3.a WHERE |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction}, local: []ref{{table: "t1"}, {table: "t2"}, {table: "t3"}}},
		},
		{
			name: "trigger body statement after the first names the row",
			sql:  "CREATE TRIGGER tr BEFORE INSERT ON db.t FOR EACH ROW BEGIN SET NEW.a = 1; SET NEW.b = NEW.|",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"NEW"}, local: []ref{{database: "db", table: "t", alias: "NEW"}}},
		},
		{
			name: "parenthesized first arm names the set operation result",
			sql:  "(SELECT a AS first_name FROM t) UNION SELECT b AS second_name FROM u ORDER BY |",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindDatabase, ObjectKindTable, ObjectKindView, ObjectKindColumn, ObjectKindFunction, ObjectKindSelectAlias}, local: []ref{{table: "u"}}, aliases: []string{"first_name"}},
		},
		{
			name: "nested cte body names its columns",
			sql:  "WITH c AS (WITH d AS (SELECT 1 AS x) SELECT x AS y FROM d) SELECT c.| FROM c",
			want: wantSignals{prefix: "", kinds: []ObjectKind{ObjectKindTable, ObjectKindView, ObjectKindColumn}, qualifier: []string{"c"}, local: []ref{{kind: RangeReferenceCTE, table: "c"}}, ctes: map[string][]string{"c": {"y"}}, cols: map[string][]string{"c": {"y"}}},
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

// TestCollectCompletionNestedWithClauses pins that a statement nested in WITH
// bodies is read once per clause, not once per enclosing clause: a body's output
// columns are its own query's, so stepping over a clause while reading another
// must not pay for the bodies inside. The bound is generous — the repeated read
// this guards was quadratic and took seconds here — because a completion request
// runs on every keystroke.
func TestCollectCompletionNestedWithClauses(t *testing.T) {
	const depth = 4800
	sql := "SELECT * FROM (" + strings.Repeat("WITH c AS (", depth) + "SELECT 1" +
		strings.Repeat(") SELECT * FROM c", depth) + ") d WHERE "
	start := time.Now()
	ctx := CollectCompletion(sql, len(sql))
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("read %d nested WITH clauses in %s, want each clause read once", depth, elapsed)
	}
	if ctx == nil || ctx.Candidates == nil {
		t.Fatal("expected a context")
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

// TestCollectCompletionReferenceColumns pins the columns the reader reports for
// a reference, which is what a caller cannot compute from metadata when the SQL
// itself names them.
func TestCollectCompletionReferenceColumns(t *testing.T) {
	cases := []struct {
		name  string
		sql   string
		ctes  map[string][]string
		local map[string][]string
		outer map[string][]string
	}{
		{
			name:  "nested with body names its output after the clause",
			sql:   "WITH c AS (WITH inner_cte(x) AS (SELECT 1) SELECT x AS y FROM inner_cte) SELECT c.| FROM c",
			ctes:  map[string][]string{"c": {"y"}, "inner_cte": {"x"}},
			local: map[string][]string{"c": {"y"}},
		},
		{
			name:  "a shadowing with clause does not rewrite the outer reference",
			sql:   "WITH c(x) AS (SELECT 1) SELECT * FROM c t WHERE EXISTS (WITH c(y) AS (SELECT 2) SELECT * FROM c WHERE |)",
			ctes:  map[string][]string{"c": {"y"}},
			local: map[string][]string{"c": {"y"}},
			outer: map[string][]string{"t": {"x"}},
		},
		{
			name:  "expression items name no column",
			sql:   "SELECT * FROM (SELECT a + b, a + t.b, a * b, count(*), t.a AS z FROM t) d WHERE d.|",
			local: map[string][]string{"d": {"z"}},
		},
		{
			name:  "qualified names name their last part",
			sql:   "SELECT * FROM (SELECT a, t.b, db.t.c FROM t) d WHERE d.|",
			local: map[string][]string{"d": {"a", "b", "c"}},
		},
		{
			name:  "a table function's arguments see the sources before it",
			sql:   "SELECT * FROM t, JSON_TABLE((SELECT t.`x`), '$' COLUMNS (x INT PATH '$')) jt WHERE EXISTS (SELECT 1 FROM JSON_TABLE((SELECT t.|), '$' COLUMNS (y INT PATH '$')) j) AND 1",
			outer: map[string][]string{"t": {}},
		},
		{
			name:  "json table declares its columns",
			sql:   "SELECT jt.| FROM JSON_TABLE('[1]', '$[*]' COLUMNS (x INT PATH '$', NESTED PATH '$.a' COLUMNS (y INT PATH '$'))) AS jt",
			local: map[string][]string{"jt": {"x", "y"}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cursor := strings.Index(tc.sql, "|")
			if cursor < 0 {
				t.Fatalf("case has no cursor marker")
			}
			ctx := CollectCompletion(tc.sql[:cursor]+tc.sql[cursor+1:], cursor)
			for _, cte := range ctx.CTEs {
				want, ok := tc.ctes[cte.Table]
				if !ok {
					t.Errorf("unexpected cte %q", cte.Table)
					continue
				}
				if got := strings.Join(cte.Columns, ","); got != strings.Join(want, ",") {
					t.Errorf("cte %q columns = %q, want %q", cte.Table, got, strings.Join(want, ","))
				}
			}
			check := func(where string, refs []RangeReference, want map[string][]string) {
				for _, ref := range refs {
					if _, ok := want[ref.Table]; !ok {
						continue
					}
					if got, expected := strings.Join(ref.Columns, ","), strings.Join(want[ref.Table], ","); got != expected {
						t.Errorf("%s reference %q columns = %q, want %q", where, ref.Table, got, expected)
					}
				}
			}
			check("local", ctx.Scope.LocalReferences, tc.local)
			for _, level := range ctx.Scope.OuterReferences {
				check("outer", level, tc.outer)
			}
		})
	}
}

// TestCollectCompletionNestingLimits holds the reader to the parser's nesting
// budget: editor input is not guaranteed to parse, and a reader that keeps
// descending past the limit kills the process rather than returning a scope.
func TestCollectCompletionNestingLimits(t *testing.T) {
	parens := strings.Repeat("(", maxNestingDepth+2) + "t" + strings.Repeat(")", maxNestingDepth+2)
	jsonColumns := "SELECT * FROM JSON_TABLE('[1]', '$[*]' COLUMNS (" +
		strings.Repeat("NESTED PATH '$.a' COLUMNS (", maxNestingDepth+2) +
		"x INT PATH '$'" +
		strings.Repeat("))", maxNestingDepth+2) +
		") jt WHERE "
	// Each level of a WITH body's output columns re-enters the reader that
	// read the level around it.
	withBodies := "SELECT * FROM (" +
		strings.Repeat("WITH c AS (", maxNestingDepth+2) + "SELECT 1" +
		strings.Repeat(") SELECT * FROM c", maxNestingDepth+2) + ") d WHERE "
	braces := strings.Repeat("{ OJ ", maxNestingDepth+2) + "t" + strings.Repeat(" }", maxNestingDepth+2)
	for name, sql := range map[string]string{
		"parenthesized table sources": "SELECT * FROM " + parens + " WHERE ",
		"odbc outer join escapes":     "SELECT * FROM " + braces + " WHERE ",
		"json_table nested columns":   jsonColumns,
		"nested with bodies":          withBodies,
	} {
		t.Run(name, func(t *testing.T) {
			ctx := CollectCompletion(sql, len(sql))
			if ctx == nil || ctx.Candidates == nil {
				t.Fatal("expected a context")
			}
		})
	}
}
