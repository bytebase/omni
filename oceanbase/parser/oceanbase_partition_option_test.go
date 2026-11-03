package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// partitionOptionCase is a statement, whether OceanBase 4.3.5 parses it (false means
// the engine answers 1064), and for a refused one the text the error is at: its last
// occurrence in the statement.
type partitionOptionCase struct {
	sql    string
	parses bool
	at     string
}

// partitionOptionCases: a partition takes `[ID n] [ENGINE = InnoDB]` and a subpartition
// `[ENGINE = InnoDB]`, nothing else; ADD PARTITION and REORGANIZE PARTITION ... INTO take
// only RANGE and LIST partitions (parsePartitionDef). Read against
// 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02, the RANGE and LIST subpartitions on
// 2026-11-03. obpart_r is partitioned by RANGE, obpart_l
// by LIST, obpart_h by HASH and obpart_k by KEY.
var partitionOptionCases = []partitionOptionCase{
	// COMMENT, on a partition of every kind and on a subpartition.
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) COMMENT 'x')", false, "COMMENT"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) COMMENT = 'x')", false, "COMMENT"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) COMMENT 'x' ENGINE = InnoDB)", false, "COMMENT"},
	{"CREATE TABLE t (a INT) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1) COMMENT 'x')", false, "COMMENT"},
	{"CREATE TABLE t (a INT) PARTITION BY HASH(a) (PARTITION p0 COMMENT 'x', PARTITION p1)", false, "COMMENT"},
	{"CREATE TABLE t (a INT) PARTITION BY KEY(a) (PARTITION p0 COMMENT 'x')", false, "COMMENT"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 COMMENT 'x', SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))", false, "COMMENT"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 COMMENT 'x', SUBPARTITION s1))", false, "COMMENT"},
	{"ALTER TABLE obpart_r ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) COMMENT 'x')", false, "COMMENT"},
	{"ALTER TABLE obpart_r REORGANIZE PARTITION p0 INTO (PARTITION p0a VALUES LESS THAN (5) COMMENT 'x', PARTITION p0b VALUES LESS THAN (10))", false, "COMMENT"},
	{"ALTER TABLE obpart_r PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) COMMENT 'x')", false, "COMMENT"},
	// The other table options MySQL takes on a partition.
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) MAX_ROWS = 10)", false, "MAX_ROWS"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) MIN_ROWS = 10)", false, "MIN_ROWS"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) DATA DIRECTORY = '/tmp')", false, "DATA"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) INDEX DIRECTORY = '/tmp')", false, "INDEX"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = ts1)", false, "TABLESPACE"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE ts1)", false, "TABLESPACE"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) PCTFREE = 5)", false, "PCTFREE"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 MAX_ROWS = 5))", false, "MAX_ROWS"},
	// ENGINE takes '=' and InnoDB, once.
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE InnoDB)", false, "InnoDB"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = MyISAM)", false, "MyISAM"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = 'InnoDB')", false, "'InnoDB'"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = `InnoDB`)", false, "`InnoDB`"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB ENGINE = InnoDB)", false, "ENGINE"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 ENGINE InnoDB))", false, "InnoDB"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB)", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = innodb)", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = INNODB)", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE=InnoDB)", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 ENGINE = InnoDB, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB (SUBPARTITION s0 ENGINE = InnoDB, SUBPARTITION s1))", true, ""},
	// ID n, before ENGINE, on a partition only.
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ID 5)", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ID 5 ENGINE = InnoDB)", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY HASH(a) (PARTITION p0 ID 5, PARTITION p1 ENGINE = InnoDB)", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1) ID 5 ENGINE = InnoDB)", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) ID 4 ENGINE = InnoDB (SUBPARTITION s0 ENGINE = InnoDB))", true, ""},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ENGINE = InnoDB ID 5)", false, "ID 5"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ID = 5)", false, "= 5"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ID 5 ID 6)", false, "ID 6"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10) ID -5)", false, "-5"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 ID 3))", false, "ID 3"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 ID 3) (PARTITION p0 VALUES LESS THAN (10))", false, "ID 3"},
	{"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN MAXVALUE)", true, ""},
	// A RANGE or LIST subpartition carries VALUES, then ENGINE = InnoDB alone.
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5) ENGINE = InnoDB, SUBPARTITION s1 VALUES LESS THAN (9)))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5) COMMENT 'x'))", false, "COMMENT"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1) ENGINE = InnoDB, SUBPARTITION s1 VALUES IN (2)) (PARTITION p0 VALUES LESS THAN (10))", true, ""},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1) MAX_ROWS = 5) (PARTITION p0 VALUES LESS THAN (10))", false, "MAX_ROWS"},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (1) ID 3))", false, "ID 3"},

	// ADD PARTITION takes a list of RANGE or LIST partitions, with VALUES each.
	{"ALTER TABLE obpart_r ADD PARTITION PARTITIONS 2", false, "PARTITIONS"},
	{"ALTER TABLE obpart_h ADD PARTITION PARTITIONS 2", false, "PARTITIONS"},
	{"ALTER TABLE obpart_k ADD PARTITION PARTITIONS 1", false, "PARTITIONS"},
	{"ALTER TABLE obpart_h ADD PARTITION (PARTITION p2)", false, ")"},
	{"ALTER TABLE obpart_k ADD PARTITION (PARTITION p2)", false, ")"},
	{"ALTER TABLE obpart_r ADD PARTITION (PARTITION p1)", false, ")"},
	{"ALTER TABLE obpart_r ADD PARTITION (PARTITION p1 VALUES LESS THAN (20), PARTITION p2)", false, ")"},
	{"ALTER TABLE obpart_r ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) ENGINE = InnoDB)", true, ""},
	{"ALTER TABLE obpart_r ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) ID 7)", true, ""},
	// These two parse and are then refused, with 1149 and 1480.
	{"ALTER TABLE obpart_r ADD PARTITION (PARTITION p1 VALUES IN (20))", true, ""},
	{"ALTER TABLE obpart_l ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))", true, ""},
	// REORGANIZE PARTITION ... INTO reads the same list; this one parses and is then
	// refused with 1235.
	{"ALTER TABLE obpart_k REORGANIZE PARTITION p0 INTO (PARTITION p0a, PARTITION p0b)", false, ", PARTITION p0b"},
	{"ALTER TABLE obpart_r REORGANIZE PARTITION p0 INTO (PARTITION p0a VALUES LESS THAN (5), PARTITION p0b VALUES LESS THAN (10))", true, ""},
}

// TestOceanBasePartitionOptionGrammar holds the parser to partitionOptionCases. A
// refused statement fails where OceanBase's "near" says.
// TestOceanBasePartitionOptionGrammarOracle asks the engine the same statements.
func TestOceanBasePartitionOptionGrammar(t *testing.T) {
	for _, c := range partitionOptionCases {
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

// TestOceanBasePartitionOptionGrammarOracle asks OceanBase about partitionOptionCases
// through PREPARE, which parses without changing anything: 1064 is a syntax error, any
// other answer means the statement parsed.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBasePartitionOptionGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obpart",
		"USE obpart",
		"DROP TABLE IF EXISTS obpart_r, obpart_l, obpart_h, obpart_k",
		"CREATE TABLE obpart_r (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE obpart_l (a INT) PARTITION BY LIST(a) (PARTITION p0 VALUES IN (1))",
		"CREATE TABLE obpart_h (a INT) PARTITION BY HASH(a) PARTITIONS 2",
		"CREATE TABLE obpart_k (a INT) PARTITION BY KEY(a) PARTITIONS 2",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obpart")
	})
	for _, c := range partitionOptionCases {
		_, err := db.ExecContext(ctx, "PREPARE obpart_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
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
