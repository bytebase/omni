package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	gomysql "github.com/go-sql-driver/mysql"
)

// subpartitionCase is a statement, whether OceanBase 4.3.5 parses it (false means the
// engine answers 1064), and for a refused one the text the error is at: its last
// occurrence in the statement.
type subpartitionCase struct {
	sql    string
	parses bool
	at     string
}

// subpartitionCases: OceanBase subpartitions by RANGE and LIST as well as HASH and KEY,
// with subpartitions of the subpartitioning's kind in a template and of one kind in a
// partition's list (parseSubpartitionClause). Read against 5.7.25-OceanBase_CE-v4.3.5.6
// on 2026-11-02.
var subpartitionCases = []subpartitionCase{
	// What SHOW CREATE TABLE prints for RANGE and LIST subpartitioning.
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` int(11) DEFAULT NULL)\n partition by range(`a`) subpartition by range(`b`) subpartition template (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE))\n(partition `p0` values less than (10),\npartition `p1` values less than (20))", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` int(11) DEFAULT NULL)\n partition by range(`a`) subpartition by list(`b`) subpartition template (\nsubpartition `s0` values in (1,2),\nsubpartition `s1` values in (DEFAULT))\n(partition `p0` values less than (10))", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` int(11) DEFAULT NULL)\n partition by list(`a`) subpartition by range(`b`)\n(partition `p0` values in (1) (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (10)),\npartition `p1` values in (2) (\nsubpartition `s2` values less than (100)))", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` varchar(10) DEFAULT NULL, `c` date DEFAULT NULL)\n partition by range columns(`a`) subpartition by range columns(`b`, `c`) subpartition template (\nsubpartition `s0` values less than ('m','2020-01-01'),\nsubpartition `s1` values less than (MAXVALUE,MAXVALUE))\n(partition `p0` values less than (10))", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` varchar(10) DEFAULT NULL, `c` int(11) DEFAULT NULL)\n partition by hash(`a`) subpartition by list columns(`b`, `c`) subpartition template (\nsubpartition `s0` values in (('x',1),('y',2)),\nsubpartition `s1` values in (DEFAULT))\n(partition `p0`,\npartition `p1`)", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` int(11) DEFAULT NULL)\n partition by key(`a`) subpartition by list(`b`)\n(partition `p0` (\nsubpartition `s0` values in (1),\nsubpartition `s1` values in (2)),\npartition `p1` (\nsubpartition `s2` values in (3)))", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` int(11) DEFAULT NULL)\n partition by range(`a`) subpartition by range((`b` + 1)) subpartition template (\nsubpartition `s0` values less than (5))\n(partition `p0` values less than (10))", true, ""},
	{"CREATE TABLE `t` (`a` int(11) DEFAULT NULL, `b` int(11) DEFAULT NULL)\n partition by range(`a`) subpartition by range(`b`)\n(partition `p0` values less than (10) (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE)),\npartition `p1` values less than (20) (\nsubpartition `p1sp0` values less than (MAXVALUE)))", true, ""},

	// The same as users write it.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN MAXVALUE) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS(a) SUBPARTITION BY RANGE COLUMNS(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5) ENGINE = InnoDB) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) PARTITIONS 2 SUBPARTITION BY LIST(b) (PARTITION p0 (SUBPARTITION s0 VALUES IN (1)), PARTITION p1 (SUBPARTITION s1 VALUES IN (2)))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5)), PARTITION p1 VALUES LESS THAN (20))", true, ""},
	// These parse and are then refused with 1210: a subpartition of the wrong kind for
	// SUBPARTITION BY, and a RANGE table with no partitions listed.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (5)))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (3)))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION s1 VALUES LESS THAN (5)))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE))", true, ""},
	// PARTITIONS n may follow the subpartition clause, once.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 PARTITIONS 3", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0) PARTITIONS 2", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) PARTITIONS 1 SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) PARTITIONS 1 (PARTITION p0 VALUES LESS THAN (10))", true, ""},

	// A template entry is of the subpartitioning's kind.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0) (PARTITION p0 VALUES LESS THAN (10))", false, ") (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))", false, "VALUES LESS THAN (5)"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))", false, "LESS THAN (5)"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (5)) (PARTITION p0 VALUES LESS THAN (10))", false, "IN (5)"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN 1) (PARTITION p0 VALUES LESS THAN (10))", false, "1) (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN 5) (PARTITION p0 VALUES LESS THAN (10))", false, "5) (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN ()) (PARTITION p0 VALUES LESS THAN (10))", false, ")) (PARTITION"},
	// A partition's subpartitions are of one kind, set by the first.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1))", false, "))"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES IN (6)))", false, "IN (6)"},
	// SUBPARTITION BY takes no LINEAR, no KEY ALGORITHM, no empty column list, and no
	// SUBPARTITIONS n for RANGE or LIST or beside a template.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LINEAR HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", false, "LINEAR"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY KEY ALGORITHM = 2 (b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", false, "ALGORITHM"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY KEY() SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", false, ") SUBPARTITIONS"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS() SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))", false, ") SUBPARTITION TEMPLATE"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (3)))", false, "SUBPARTITIONS"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITIONS 2 SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))", false, "SUBPARTITIONS"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 SUBPARTITION TEMPLATE (SUBPARTITION s0) (PARTITION p0 VALUES LESS THAN (10))", false, "SUBPARTITION TEMPLATE"},

	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (1), SUBPARTITION s1 VALUES LESS THAN (5)))", false, "LESS THAN (5)))"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1 VALUES IN (5)))", false, "VALUES IN (5)"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1 VALUES IN (5)))", false, "VALUES IN (5)"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1), SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))", false, ") (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (DEFAULT, 1)) (PARTITION p0 VALUES LESS THAN (10))", false, ", 1))"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE () (PARTITION p0 VALUES LESS THAN (10))", false, ") (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (PARTITION p0 VALUES LESS THAN (10))", false, "PARTITION p0"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LINEAR KEY(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", false, "LINEAR"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITIONS 2 SUBPARTITION TEMPLATE (SUBPARTITION s0) (PARTITION p0 VALUES LESS THAN (10))", false, "SUBPARTITION TEMPLATE"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))", false, "SUBPARTITIONS"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY HASH(b) SUBPARTITIONS 2 PARTITIONS 2", false, "PARTITIONS 2"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 PARTITIONS 1 (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 PARTITIONS 2", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY LIST(b)", true, ""},
	// These parse and are then refused, with 1653 and 1504.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (MAXVALUE, 1)) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 0 (PARTITION p0 VALUES LESS THAN (10))", true, ""},

	// VALUES needs its bound.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN) (PARTITION p0 VALUES LESS THAN (10))", false, ") (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN) (PARTITION p0 VALUES LESS THAN (10))", false, ") (PARTITION"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN))", false, "))"},
	// MAXVALUE beside other bounds.
	{"CREATE TABLE t (a INT, b INT, c INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5, MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))", true, ""},

	// LESS takes THAN.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS (5)) (PARTITION p0 VALUES LESS THAN (10))", false, "(5)"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS (5)))", false, "(5)"},

	// The trailing PARTITIONS takes its count.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 PARTITIONS", false, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 PARTITIONS (PARTITION p0 VALUES LESS THAN (10))", false, "(PARTITION"},

	// ADD PARTITION takes a partition's own RANGE or LIST subpartitions.
	{"ALTER TABLE obsub_rl ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0 VALUES IN (5), SUBPARTITION z1 VALUES IN (DEFAULT)))", true, ""},
	{"ALTER TABLE obsub_rr ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0 VALUES LESS THAN (7)))", true, ""},
}

// TestOceanBaseSubpartitionGrammar holds the parser to subpartitionCases. A refused
// statement fails where OceanBase's "near" says. TestOceanBaseSubpartitionGrammarOracle
// asks the engine the same statements.
func TestOceanBaseSubpartitionGrammar(t *testing.T) {
	for _, c := range subpartitionCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.parses {
			t.Errorf("Parse(%q) error = %v, want parses %v", c.sql, err, c.parses)
			continue
		}
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) error %v is not a *ParseError", c.sql, err)
			continue
		}
		if want := strings.LastIndex(c.sql, c.at); pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d at %q (%v)", c.sql, pe.Position, want, c.at, err)
		}
	}
}

// TestOceanBaseSubpartitionAST pins the clause RANGE and LIST subpartitioning give: the
// type with its expression or columns, and each subpartition's bound.
func TestOceanBaseSubpartitionAST(t *testing.T) {
	list, err := Parse("CREATE TABLE t (a INT, b INT, c INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN ((1, 2)), SUBPARTITION s1 VALUES IN (DEFAULT)) (PARTITION p0 VALUES LESS THAN (10))")
	if err != nil {
		t.Fatal(err)
	}
	pc := list.Items[0].(*nodes.CreateTableStmt).Partitions
	if pc.SubPartType != nodes.PartitionList || strings.Join(pc.SubPartColumns, ",") != "b,c" || pc.SubPartExpr != nil {
		t.Errorf("subpartitioning = %v %v %v, want LIST COLUMNS (b, c)", pc.SubPartType, pc.SubPartColumns, pc.SubPartExpr)
	}
	if len(pc.SubPartTemplate) != 2 || pc.SubPartTemplate[0].Values == nil || pc.SubPartTemplate[1].Values == nil {
		t.Fatalf("template = %+v, want two entries with values", pc.SubPartTemplate)
	}

	list, err = Parse("CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES IN (1) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN MAXVALUE))")
	if err != nil {
		t.Fatal(err)
	}
	pc = list.Items[0].(*nodes.CreateTableStmt).Partitions
	if pc.SubPartType != nodes.PartitionRange || pc.SubPartExpr == nil {
		t.Errorf("subpartitioning = %v %v, want RANGE (b)", pc.SubPartType, pc.SubPartExpr)
	}
	subs := pc.Partitions[0].SubPartitions
	if len(subs) != 2 || subs[0].Values == nil {
		t.Fatalf("subpartitions = %+v, want two with values", subs)
	}
	if s, ok := subs[1].Values.(*nodes.String); !ok || s.Str != "MAXVALUE" {
		t.Errorf("s1 values = %#v, want MAXVALUE", subs[1].Values)
	}
}

// TestOceanBaseSubpartitionGrammarOracle asks OceanBase about subpartitionCases through
// PREPARE, which parses without creating anything: 1064 is a syntax error, any other
// answer means the statement parsed.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseSubpartitionGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obsub",
		"USE obsub",
		"DROP TABLE IF EXISTS obsub_rl, obsub_rr",
		"CREATE TABLE obsub_rl (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (DEFAULT)) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE obsub_rr (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5)))",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obsub")
	})
	for _, c := range subpartitionCases {
		_, err := db.ExecContext(ctx, "PREPARE obsub_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
		engineParses := err == nil
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
			}
			engineParses = myErr.Number != 1064
		}
		if engineParses != c.parses {
			t.Errorf("OceanBase parses=%v, want %v for %q: %v", engineParses, c.parses, c.sql, err)
		}
		if _, perr := Parse(c.sql); (perr == nil) != engineParses {
			t.Errorf("DIVERGENCE omni parses=%v, OceanBase parses=%v for %q: %v", perr == nil, engineParses, c.sql, perr)
		}
	}
}
