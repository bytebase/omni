package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// subpartitionValueLayout lists a table's partitions with their subpartitions and each
// subpartition's bound: "p0:s0<5>,s1<MAXVALUE>".
func subpartitionValueLayout(tbl *Table) []string {
	if tbl.Partitioning == nil {
		return nil
	}
	var out []string
	for _, pd := range tbl.Partitioning.Partitions {
		subs := make([]string, 0, len(pd.SubPartitions))
		for _, sp := range pd.SubPartitions {
			subs = append(subs, sp.Name+"<"+sp.ValueExpr+">")
		}
		out = append(out, pd.Name+":"+strings.Join(subs, ","))
	}
	return out
}

// rangeListSubpartitionCases are RANGE and LIST subpartitioned tables, each with the
// layout OceanBase gives it and the partition clause 4.3.5.6 prints for it in SHOW
// CREATE TABLE (read on 2026-11-02). The engine names a template's subpartitions after
// the partition and the entry, as it does for HASH.
var rangeListSubpartitionCases = []struct {
	name, ddl, readback string
	layout              []string
}{
	{"range template",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range(`b`) subpartition template (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE))\n(partition `p0` values less than (10),\npartition `p1` values less than (20))",
		[]string{"p0:p0ss0<5>,p0ss1<MAXVALUE>", "p1:p1ss0<5>,p1ss1<MAXVALUE>"}},
	{"list template",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (DEFAULT)) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by list(`b`) subpartition template (\nsubpartition `s0` values in (1,2),\nsubpartition `s1` values in (DEFAULT))\n(partition `p0` values less than (10))",
		[]string{"p0:p0ss0<1,2>,p0ss1<DEFAULT>"}},
	{"range per partition",
		"CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES IN (1) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (10)), PARTITION p1 VALUES IN (2) (SUBPARTITION s2 VALUES LESS THAN (100)))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by list(`a`) subpartition by range(`b`)\n(partition `p0` values in (1) (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (10)),\npartition `p1` values in (2) (\nsubpartition `s2` values less than (100)))",
		[]string{"p0:s0<5>,s1<10>", "p1:s2<100>"}},
	{"range columns template",
		"CREATE TABLE t (a INT, b VARCHAR(10), c DATE) PARTITION BY RANGE COLUMNS(a) SUBPARTITION BY RANGE COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN ('m', '2020-01-01'), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE, MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` varchar(10) DEFAULT NULL,\n  `c` date DEFAULT NULL\n)\n partition by range columns(`a`) subpartition by range columns(`b`, `c`) subpartition template (\nsubpartition `s0` values less than ('m','2020-01-01'),\nsubpartition `s1` values less than (MAXVALUE,MAXVALUE))\n(partition `p0` values less than (10))",
		[]string{"p0:p0ss0<'m','2020-01-01'>,p0ss1<MAXVALUE>"}},
	{"list columns template over hash",
		"CREATE TABLE t (a INT, b VARCHAR(10), c INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY LIST COLUMNS(b, c) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (('x', 1), ('y', 2)), SUBPARTITION s1 VALUES IN (DEFAULT))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` varchar(10) DEFAULT NULL,\n  `c` int(11) DEFAULT NULL\n)\n partition by hash(`a`) subpartition by list columns(`b`, `c`) subpartition template (\nsubpartition `s0` values in (('x',1),('y',2)),\nsubpartition `s1` values in (DEFAULT))\n(partition `p0`,\npartition `p1`)",
		[]string{"p0:p0ss0<('x',1),('y',2)>,p0ss1<DEFAULT>", "p1:p1ss0<('x',1),('y',2)>,p1ss1<DEFAULT>"}},
	{"list per partition over key",
		"CREATE TABLE t (a INT, b INT) PARTITION BY KEY(a) SUBPARTITION BY LIST(b) (PARTITION p0 (SUBPARTITION s0 VALUES IN (1), SUBPARTITION s1 VALUES IN (2)), PARTITION p1 (SUBPARTITION s2 VALUES IN (3)))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by key(`a`) subpartition by list(`b`)\n(partition `p0` (\nsubpartition `s0` values in (1),\nsubpartition `s1` values in (2)),\npartition `p1` (\nsubpartition `s2` values in (3)))",
		[]string{"p0:s0<1>,s1<2>", "p1:s2<3>"}},
	{"range expression template",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b + 1) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5)) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range((`b` + 1)) subpartition template (\nsubpartition `s0` values less than (5))\n(partition `p0` values less than (10))",
		[]string{"p0:p0ss0<5>"}},
	{"range columns template over list columns",
		"CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS(a) SUBPARTITION BY RANGE COLUMNS(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES IN (1, 2), PARTITION p1 VALUES IN (3))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by list columns(`a`) subpartition by range columns(`b`) subpartition template (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE))\n(partition `p0` values in (1,2),\npartition `p1` values in (3))",
		[]string{"p0:p0ss0<5>,p0ss1<MAXVALUE>", "p1:p1ss0<5>,p1ss1<MAXVALUE>"}},
	// A partition that lists no subpartitions gets one, to MAXVALUE or DEFAULT, which
	// OceanBase names P8192, P8193, ... across the table rather than <partition>sp0.
	{"range, a partition without",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5)), PARTITION p1 VALUES LESS THAN (20))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range(`b`)\n(partition `p0` values less than (10) (\nsubpartition `s0` values less than (5)),\npartition `p1` values less than (20) (\nsubpartition `P8192` values less than (MAXVALUE)))",
		[]string{"p0:s0<5>", "p1:P8192<MAXVALUE>"}},
	{"range, none listed",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range(`b`)\n(partition `p0` values less than (10) (\nsubpartition `P8192` values less than (MAXVALUE)))",
		[]string{"p0:P8192<MAXVALUE>"}},
	{"list, a partition without",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (1)), PARTITION p1 VALUES LESS THAN (20))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by list(`b`)\n(partition `p0` values less than (10) (\nsubpartition `s0` values in (1)),\npartition `p1` values less than (20) (\nsubpartition `P8192` values in (DEFAULT)))",
		[]string{"p0:s0<1>", "p1:P8192<DEFAULT>"}},
	{"list columns, none listed",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST COLUMNS(b) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by list columns(`b`)\n(partition `p0` values less than (10) (\nsubpartition `P8192` values in (DEFAULT)))",
		[]string{"p0:P8192<DEFAULT>"}},
	{"range over hash, a partition without",
		"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY RANGE(b) (PARTITION p0 (SUBPARTITION s0 VALUES LESS THAN (5)), PARTITION p1)",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by hash(`a`) subpartition by range(`b`)\n(partition `p0` (\nsubpartition `s0` values less than (5)),\npartition `p1` (\nsubpartition `P8192` values less than (MAXVALUE)))",
		[]string{"p0:s0<5>", "p1:P8192<MAXVALUE>"}},
	{"list over hash partitions, none listed",
		"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY LIST(b)",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by hash(`a`) subpartition by list(`b`)\n(partition `p0` (\nsubpartition `P8192` values in (DEFAULT)),\npartition `p1` (\nsubpartition `P8193` values in (DEFAULT)))",
		[]string{"p0:P8192<DEFAULT>", "p1:P8193<DEFAULT>"}},
	{"range template to MAXVALUE without parentheses",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN MAXVALUE) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range(`b`) subpartition template (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE))\n(partition `p0` values less than (10))",
		[]string{"p0:p0ss0<5>,p0ss1<MAXVALUE>"}},
	{"range columns template over a date",
		"CREATE TABLE t (a INT, b DATE) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN ('2020-01-01'), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` date DEFAULT NULL\n)\n partition by range(`a`) subpartition by range columns(`b`) subpartition template (\nsubpartition `s0` values less than ('2020-01-01'),\nsubpartition `s1` values less than (MAXVALUE))\n(partition `p0` values less than (10))",
		[]string{"p0:p0ss0<'2020-01-01'>,p0ss1<MAXVALUE>"}},
	{"range bound folded",
		"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b + 1) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5 + 1)))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range((`b` + 1))\n(partition `p0` values less than (10) (\nsubpartition `s0` values less than (6)))",
		[]string{"p0:s0<6>"}},
	{"list over list",
		"CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES IN (1, 2) (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (3)))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by list(`a`) subpartition by list(`b`)\n(partition `p0` values in (1,2) (\nsubpartition `s0` values in (1,2),\nsubpartition `s1` values in (3)))",
		[]string{"p0:s0<1,2>,s1<3>"}},
	{"list columns per partition with DEFAULT",
		"CREATE TABLE t (a INT, b VARCHAR(10), c INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST COLUMNS(b, c) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (('x', 1), ('y', 2)), SUBPARTITION s1 VALUES IN (DEFAULT)))",
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` varchar(10) DEFAULT NULL,\n  `c` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by list columns(`b`, `c`)\n(partition `p0` values less than (10) (\nsubpartition `s0` values in (('x',1),('y',2)),\nsubpartition `s1` values in (DEFAULT)))",
		[]string{"p0:s0<('x',1),('y',2)>,s1<DEFAULT>"}},
}

// rangeListAddPartitionScenarios add partitions to RANGE and LIST subpartitioned tables.
// Each records the layout OceanBase gives the table and the readback it prints after.
var rangeListAddPartitionScenarios = []struct {
	name     string
	ddl      []string
	readback string
	layout   []string
}{
	{"range template takes the bounds",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range(`b`) subpartition template (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE))\n(partition `p0` values less than (10),\npartition `p2` values less than (30))",
		[]string{"p0:p0ss0<5>,p0ss1<MAXVALUE>", "p2:p2ss0<5>,p2ss1<MAXVALUE>"}},
	{"list template after a partition's own",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES IN (1, 2), SUBPARTITION s1 VALUES IN (DEFAULT)) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION z0 VALUES IN (5), SUBPARTITION z1 VALUES IN (DEFAULT)))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by list(`b`)\n(partition `p0` values less than (10) (\nsubpartition `p0ss0` values in (1,2),\nsubpartition `p0ss1` values in (DEFAULT)),\npartition `p1` values less than (20) (\nsubpartition `z0` values in (5),\nsubpartition `z1` values in (DEFAULT)),\npartition `p2` values less than (30) (\nsubpartition `p2ss0` values in (1,2),\nsubpartition `p2ss1` values in (DEFAULT)))",
		[]string{"p0:p0ss0<1,2>,p0ss1<DEFAULT>", "p1:z0<5>,z1<DEFAULT>", "p2:p2ss0<1,2>,p2ss1<DEFAULT>"}},
	{"range per partition adds one to MAXVALUE",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
		},
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range(`b`)\n(partition `p0` values less than (10) (\nsubpartition `s0` values less than (5),\nsubpartition `s1` values less than (MAXVALUE)),\npartition `p1` values less than (20) (\nsubpartition `p1sp0` values less than (MAXVALUE)))",
		[]string{"p0:s0<5>,s1<MAXVALUE>", "p1:p1sp0<MAXVALUE>"}},
	{"range per partition adds its own",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY RANGE(b) (PARTITION p0 VALUES IN (1) (SUBPARTITION s0 VALUES LESS THAN (5)))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES IN (3) (SUBPARTITION s3 VALUES LESS THAN (7)))",
		},
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by list(`a`) subpartition by range(`b`)\n(partition `p0` values in (1) (\nsubpartition `s0` values less than (5)),\npartition `p2` values in (3) (\nsubpartition `s3` values less than (7)))",
		[]string{"p0:s0<5>", "p2:s3<7>"}},
	{"list per partition adds one to DEFAULT after a made-up one",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY LIST(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES IN (1)), PARTITION p1 VALUES LESS THAN (20))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by list(`b`)\n(partition `p0` values less than (10) (\nsubpartition `s0` values in (1)),\npartition `p1` values less than (20) (\nsubpartition `P8192` values in (DEFAULT)),\npartition `p2` values less than (30) (\nsubpartition `p2sp0` values in (DEFAULT)))",
		[]string{"p0:s0<1>", "p1:P8192<DEFAULT>", "p2:p2sp0<DEFAULT>"}},
	{"range columns none listed adds one to MAXVALUE",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE COLUMNS(b) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
		},
		"CREATE TABLE `t` (\n  `a` int(11) DEFAULT NULL,\n  `b` int(11) DEFAULT NULL\n)\n partition by range(`a`) subpartition by range columns(`b`)\n(partition `p0` values less than (10) (\nsubpartition `P8192` values less than (MAXVALUE)),\npartition `p1` values less than (20) (\nsubpartition `p1sp0` values less than (MAXVALUE)))",
		[]string{"p0:P8192<MAXVALUE>", "p1:p1sp0<MAXVALUE>"}},
}

// requirePartitionsEqual loads two forms of a table and requires the same layout with
// bounds and no partition diff between them, in either direction.
func requirePartitionsEqual(t *testing.T, a, b string, layout []string) {
	t.Helper()
	sc := serverCharsetFor(MySQL80)
	n := NormalizerFor(MySQL80)
	ca, ta := loadOneTable(t, sc, a, "t")
	cb, tb := loadOneTable(t, sc, b, "t")
	for _, side := range []struct {
		what string
		tbl  *Table
	}{{"first", ta}, {"second", tb}} {
		if got := subpartitionValueLayout(side.tbl); strings.Join(got, " ") != strings.Join(layout, " ") {
			t.Errorf("the %s form reads as %v, want %v\n%s", side.what, got, layout, map[string]string{"first": a, "second": b}[side.what])
		}
	}
	for _, d := range []*SchemaDiff{DiffWithNormalizer(ca, cb, n), DiffWithNormalizer(cb, ca, n), DiffWithNormalizer(cb, cb, n)} {
		for _, te := range d.Tables {
			if te.PartitionChanged {
				t.Errorf("the partitioning differs:\n  %s\n  %s", a, b)
			}
		}
	}
}

// TestRangeListSubpartitionReadback pins that each RANGE or LIST subpartitioned table
// and the form OceanBase prints for it read as the same layout, bounds included, and
// compare equal. TestOracle_OceanBaseRangeListSubpartitions asks the engine.
func TestRangeListSubpartitionReadback(t *testing.T) {
	for _, tc := range rangeListSubpartitionCases {
		t.Run(tc.name, func(t *testing.T) {
			requirePartitionsEqual(t, tc.ddl, tc.readback, tc.layout)
		})
	}
}

// TestRangeListAddPartition pins the subpartitions a partition added to a RANGE or LIST
// subpartitioned table gets: its own, the template's with their bounds, or one to
// MAXVALUE.
func TestRangeListAddPartition(t *testing.T) {
	for _, sc := range rangeListAddPartitionScenarios {
		t.Run(sc.name, func(t *testing.T) {
			requirePartitionsEqual(t, strings.Join(sc.ddl, ";\n")+";", sc.readback, sc.layout)
		})
	}
}

// TestRangeListSubpartitionRendering pins that the catalog's own rendering of a RANGE or
// LIST subpartitioned table reads back to the same table.
func TestRangeListSubpartitionRendering(t *testing.T) {
	for _, tc := range rangeListSubpartitionCases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+tc.ddl+";", nil); err != nil {
				t.Fatal(err)
			}
			requirePartitionsEqual(t, tc.ddl, c.ShowCreateTable("d", "t"), tc.layout)
		})
	}
}

// TestOracle_OceanBaseRangeListSubpartitions creates each of rangeListSubpartitionCases
// and runs each of rangeListAddPartitionScenarios on OceanBase, and requires the engine
// to print the recorded partition clause and give the recorded layout, and the catalog to
// read both the DDL and the readback as that layout.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/catalog/ -run OceanBase -count=1
func TestOracle_OceanBaseRangeListSubpartitions(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	check := func(t *testing.T, dbName string, ddl []string, recorded string, layout []string) {
		t.Helper()
		if _, err := storeAndRead(ctx, o, dbName, ddl[0], "t"); err != nil {
			t.Fatalf("OceanBase would not take the table: %v", err)
		}
		t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName) })
		for _, stmt := range ddl[1:] {
			if _, err := o.db.ExecContext(ctx, stmt); err != nil {
				t.Fatalf("OceanBase refused %q: %v", stmt, err)
			}
		}
		var readback string
		if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE "+dbName+".t").Scan(new(string), &readback); err != nil {
			t.Fatal(err)
		}
		clause := recorded[strings.Index(recorded, "\n partition by"):]
		if !strings.Contains(readback, clause) {
			t.Errorf("OceanBase prints\n%s\nthe case records%s", readback, clause)
		}
		rows, err := o.db.QueryContext(ctx, "SELECT PARTITION_NAME, SUBPARTITION_NAME FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION", dbName)
		if err != nil {
			t.Fatal(err)
		}
		var engine []string
		for rows.Next() {
			var part, sub string
			if err := rows.Scan(&part, &sub); err != nil {
				t.Fatal(err)
			}
			if n := len(engine); n > 0 && strings.HasPrefix(engine[n-1], part+":") {
				engine[n-1] += "," + sub
			} else {
				engine = append(engine, part+":"+sub)
			}
		}
		_ = rows.Close()
		var names []string
		for _, l := range layout {
			part, subs, _ := strings.Cut(l, ":")
			names = append(names, part+":"+stripBounds(subs))
		}
		if strings.Join(engine, " ") != strings.Join(names, " ") {
			t.Errorf("OceanBase layout = %v, the case records %v", engine, names)
		}
		requirePartitionsEqual(t, strings.Join(ddl, ";\n")+";", readback, layout)
	}
	for i, tc := range rangeListSubpartitionCases {
		t.Run(tc.name, func(t *testing.T) {
			check(t, fmt.Sprintf("obrlsub_%d", i), []string{tc.ddl}, tc.readback, tc.layout)
		})
	}
	for i, sc := range rangeListAddPartitionScenarios {
		t.Run(sc.name, func(t *testing.T) {
			check(t, fmt.Sprintf("obrladd_%d", i), sc.ddl, sc.readback, sc.layout)
		})
	}
}

// stripBounds drops the <bound> after each subpartition name of a layout entry's list.
func stripBounds(subs string) string {
	var b strings.Builder
	depth := 0
	for _, r := range subs {
		switch {
		case r == '<':
			depth++
		case r == '>':
			depth--
		case depth == 0:
			b.WriteRune(r)
		}
	}
	return b.String()
}
