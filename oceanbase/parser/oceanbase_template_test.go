package parser

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// templateCase is a statement and whether OceanBase 4.3.5 parses it: false means the
// engine answers 1064.
type templateCase struct {
	sql    string
	parses bool
}

// templateCases: with a SUBPARTITION TEMPLATE no partition may list subpartitions of
// its own, whether or not they match the template. Without a template, partitions may
// list different numbers of subpartitions. Read against 5.7.25-OceanBase_CE-v4.3.5.6
// on 2026-11-02.
var templateCases = []templateCase{
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))", true},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)", true},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION z0, SUBPARTITION z1, SUBPARTITION z2), PARTITION p1 VALUES LESS THAN (20))", false},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION z0, SUBPARTITION z1))", false},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION z0, SUBPARTITION z1), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z2, SUBPARTITION z3))", false},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0))", false},
	{"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION s2))", true},
}

// TestOceanBaseTemplateGrammar holds the parser to templateCases. A partition that
// lists subpartitions under a template fails at its first SUBPARTITION.
func TestOceanBaseTemplateGrammar(t *testing.T) {
	for _, c := range templateCases {
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
		if want := strings.Index(c.sql, "(SUBPARTITION z") + 1; pe.Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%v)", c.sql, pe.Position, want, err)
		}
	}
}

// TestOceanBaseTemplateGrammarOracle asks OceanBase about templateCases through
// PREPARE, which parses without creating anything.
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseTemplateGrammarOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{"CREATE DATABASE IF NOT EXISTS obtemplate", "USE obtemplate"} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obtemplate")
	})
	for _, c := range templateCases {
		_, err := db.ExecContext(ctx, "PREPARE obtemplate_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
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
