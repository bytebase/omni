package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// addPartitionScenario is a table and the ALTERs that add partitions to it, with the
// layout OceanBase leaves: each partition with its subpartitions, as
// information_schema.PARTITIONS lists them.
type addPartitionScenario struct {
	name   string
	ddl    []string
	layout []string
}

// addPartitionScenarios are read against 5.7.25-OceanBase_CE-v4.3.5.6 on 2026-11-02, by
// TestOracle_OceanBaseAddPartitionSubpartitions. A table whose CREATE lists no
// subpartitions per partition has a template: SUBPARTITION TEMPLATE names each
// subpartition after its partition and the entry (p0ss0), SUBPARTITIONS N, or no count
// at all, is a template of p0 ... p<N-1> (p0sp0), and every partition added later takes
// it, even after one is added with subpartitions of its own. A table whose CREATE lists
// subpartitions for any partition has none: a partition that lists none at CREATE gets N,
// or one, and a partition added later gets the single subpartition <partition>sp0.
var addPartitionScenarios = []addPartitionScenario{
	{
		"template, then a partition without subpartitions",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
		},
		[]string{"p0:p0ss0,p0ss1", "p1:p1ss0,p1ss1"},
	},
	{
		"template, then two partitions at once",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN (30))",
		},
		[]string{"p0:p0ss0,p0ss1", "p1:p1ss0,p1ss1", "p2:p2ss0,p2ss1"},
	},
	{
		"template, then a partition with its own subpartitions",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION x0, SUBPARTITION x1, SUBPARTITION x2))",
		},
		[]string{"p0:p0ss0,p0ss1", "p1:x0,x1,x2"},
	},
	{
		"template, then a partition with its own, then one without",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION x0))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		[]string{"p0:p0ss0,p0ss1", "p1:x0", "p2:p2ss0,p2ss1"},
	},
	{
		"template on LIST",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY LIST(a) SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES IN (1, 2))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES IN (3))",
		},
		[]string{"p0:p0ss0,p0ss1", "p1:p1ss0,p1ss1"},
	},
	{
		"template on generated HASH partitions",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY HASH(a) PARTITIONS 2 SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1)",
		},
		[]string{"p0:p0ss0,p0ss1", "p1:p1ss0,p1ss1"},
	},
	{
		"SUBPARTITIONS N",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
		},
		[]string{"p0:p0sp0,p0sp1", "p1:p1sp0,p1sp1"},
	},
	{
		"SUBPARTITIONS N, then a partition with its own, then one without",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION y0, SUBPARTITION y1, SUBPARTITION y2))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		[]string{"p0:p0sp0,p0sp1", "p1:y0,y1,y2", "p2:p2sp0,p2sp1"},
	},
	{
		"subpartitions listed per partition",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30) (SUBPARTITION s2, SUBPARTITION s3))",
		},
		[]string{"p0:s0,s1", "p1:p1sp0", "p2:s2,s3"},
	},
	{
		"partitions that list different numbers",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION s2))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		[]string{"p0:s0,s1", "p1:s2", "p2:p2sp0"},
	},
	{
		"a partition that lists none beside one that does",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0), PARTITION p1 VALUES LESS THAN (20))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		[]string{"p0:s0", "p1:p1sp0", "p2:p2sp0"},
	},
	{
		"SUBPARTITIONS N beside a partition that lists its own",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1, SUBPARTITION s2), PARTITION p1 VALUES LESS THAN (20))",
			"ALTER TABLE t ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))",
		},
		[]string{"p0:s0,s1,s2", "p1:p1sp0,p1sp1", "p2:p2sp0"},
	},
	{
		"no count and no lists",
		[]string{
			"CREATE TABLE t (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
		},
		[]string{"p0:p0sp0", "p1:p1sp0"},
	},
	{
		"no subpartitioning",
		[]string{
			"CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10))",
			"ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20))",
		},
		[]string{"p0:", "p1:"},
	},
}

// partitionLayout lists each partition of tbl with its subpartitions.
func partitionLayout(tbl *Table) []string {
	if tbl == nil || tbl.Partitioning == nil {
		return nil
	}
	var out []string
	for _, pd := range tbl.Partitioning.Partitions {
		var subs []string
		for _, sp := range pd.SubPartitions {
			subs = append(subs, sp.Name)
		}
		out = append(out, pd.Name+":"+strings.Join(subs, ","))
	}
	return out
}

// TestAddPartitionSubpartitions holds the catalog to the layouts OceanBase leaves after
// ALTER TABLE ... ADD PARTITION. The template used to be expanded at CREATE and then
// forgotten, so an added partition had no subpartitions at all, and the ones a partition
// listed itself were dropped. It also requires SHOW CREATE TABLE to show every
// subpartition once a partition lists its own, as OceanBase's does, rather than a count.
func TestAddPartitionSubpartitions(t *testing.T) {
	for _, sc := range addPartitionScenarios {
		t.Run(sc.name, func(t *testing.T) {
			c := setupWithDB(t)
			for _, ddl := range sc.ddl {
				mustExec(t, c, ddl)
			}
			got := partitionLayout(c.GetDatabase("testdb").GetTable("t"))
			if strings.Join(got, " ") != strings.Join(sc.layout, " ") {
				t.Errorf("layout = %v, want %v", got, sc.layout)
			}
			// The rendered table loads back to the same layout, so a diff against
			// OceanBase's readback compares what the engine holds.
			show := c.ShowCreateTable("testdb", "t")
			_, back := loadOneTable(t, "utf8mb4", show, "t")
			if got := partitionLayout(back); strings.Join(got, " ") != strings.Join(sc.layout, " ") {
				t.Errorf("SHOW CREATE TABLE loads back as %v, want %v\n%s", got, sc.layout, show)
			}
		})
	}
}

// TestAddPartitionKeepsTemplate pins that the template outlives a partition added with
// subpartitions of its own: OceanBase then prints every partition's subpartitions in SHOW
// CREATE TABLE, but a partition added after that without any still takes the template's.
// It also pins that a cloned table keeps a template of its own.
func TestAddPartitionKeepsTemplate(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, addPartitionScenarios[0].ddl[0])
	tbl := c.GetDatabase("testdb").GetTable("t")
	if n := len(tbl.Partitioning.SubPartTemplate); n != 2 {
		t.Fatalf("template = %d entries after CREATE, want 2", n)
	}
	clone := cloneTable(tbl)
	clone.Partitioning.SubPartTemplate[0].Name = "changed"
	if got := tbl.Partitioning.SubPartTemplate[0].Name; got != "s0" {
		t.Errorf("changing the clone's template renamed the original's entry to %q", got)
	}
	mustExec(t, c, "ALTER TABLE t ADD PARTITION (PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION x0))")
	if n := len(c.GetDatabase("testdb").GetTable("t").Partitioning.SubPartTemplate); n != 2 {
		t.Errorf("template = %d entries after a partition with its own subpartitions, want 2", n)
	}
}

// TestReorganizeSubpartitionedTable pins that REORGANIZE PARTITION on a subpartitioned
// table is refused, as 4.3.5.6 refuses it with 1235 "split subpartition is not
// supported", splitting or merging, with a template or without. It used to succeed and
// drop the subpartitions of every partition it made.
func TestReorganizeSubpartitionedTable(t *testing.T) {
	for _, ddl := range []string{
		"CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))",
		"CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0, SUBPARTITION s1), PARTITION p1 VALUES LESS THAN (20) (SUBPARTITION s2, SUBPARTITION s3))",
		"CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITIONS 2 (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))",
	} {
		for _, alter := range reorganizeSubpartitionedAlters {
			c := setupWithDB(t)
			mustExec(t, c, ddl)
			results, err := c.Exec(alter, nil)
			for _, r := range results {
				if r.Error != nil {
					err = r.Error
				}
			}
			var ce *Error
			if !errors.As(err, &ce) || ce.Code != 1235 {
				t.Errorf("%s\n%s: err = %v, want 1235", ddl, alter, err)
			}
		}
	}
	// A table with no subpartitions still reorganizes.
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a)) PARTITION BY RANGE(a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))")
	mustExec(t, c, reorganizeSubpartitionedAlters[0])
	if got := partitionLayout(c.GetDatabase("testdb").GetTable("t")); strings.Join(got, " ") != "p0: p1a: p1b:" {
		t.Errorf("layout = %v, want [p0: p1a: p1b:]", got)
	}
}

var reorganizeSubpartitionedAlters = []string{
	"ALTER TABLE t REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20))",
	"ALTER TABLE t REORGANIZE PARTITION p0, p1 INTO (PARTITION p01 VALUES LESS THAN (20))",
}

// TestOracle_OceanBaseAddPartitionSubpartitions runs each scenario on OceanBase and
// requires the layout it reports, and the layout the catalog builds from its SHOW CREATE
// TABLE, to be the one the scenario records.
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/catalog/ -run OceanBase -count=1
func TestOracle_OceanBaseAddPartitionSubpartitions(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	t.Run("reorganize", func(t *testing.T) {
		if _, err := storeAndRead(ctx, o, "obreorg", "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b)) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0, SUBPARTITION s1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))", "t"); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS obreorg") })
		for _, alter := range reorganizeSubpartitionedAlters {
			_, err := o.db.ExecContext(ctx, alter)
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) || myErr.Number != 1235 {
				t.Errorf("OceanBase answered %v for %q, want 1235", err, alter)
			}
		}
	})
	for i, sc := range addPartitionScenarios {
		t.Run(sc.name, func(t *testing.T) {
			dbName := fmt.Sprintf("obaddpart_%d", i)
			readback, err := storeAndRead(ctx, o, dbName, sc.ddl[0], "t")
			if err != nil {
				t.Fatalf("OceanBase would not take the table: %v", err)
			}
			for _, ddl := range sc.ddl[1:] {
				if _, err := o.db.ExecContext(ctx, ddl); err != nil {
					t.Fatalf("OceanBase would not take %q: %v", ddl, err)
				}
			}
			t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName) })

			rows, err := o.db.QueryContext(ctx, "SELECT PARTITION_NAME, COALESCE(SUBPARTITION_NAME, '') FROM information_schema.PARTITIONS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' ORDER BY PARTITION_ORDINAL_POSITION, SUBPARTITION_ORDINAL_POSITION", dbName)
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
			if strings.Join(engine, " ") != strings.Join(sc.layout, " ") {
				t.Errorf("OceanBase layout = %v, the scenario records %v", engine, sc.layout)
			}

			if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE "+dbName+".t").Scan(new(string), &readback); err != nil {
				t.Fatal(err)
			}
			_, stored := loadOneTable(t, serverCharsetFor(o.version), readback, "t")
			if got := partitionLayout(stored); strings.Join(got, " ") != strings.Join(sc.layout, " ") {
				t.Errorf("the catalog reads OceanBase's SHOW CREATE TABLE as %v, want %v\n%s", got, sc.layout, readback)
			}
		})
	}
}
