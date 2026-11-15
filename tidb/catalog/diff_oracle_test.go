package catalog

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// diffOraclePairs are pairs of schema scripts. TiDB decides whether a pair
// is one schema by what it stores (storedSchema), and Diff must agree both
// ways, on catalogs built by Exec and on catalogs loaded from the snapshots
// Bytebase's sync takes (syncTiDBSnapshot).
//
// syncDrops names what the snapshot leaves out, for a pair whose difference
// only that shows: the snapshot path must then find no difference, so the
// entry fails once the sync records it.
var diffOraclePairs = []struct {
	name, a, b string
	syncDrops  string
}{
	// One schema, spelled two ways.
	{name: "display width", a: "CREATE TABLE t (a int(11), b tinyint(4))", b: "CREATE TABLE t (a int, b tinyint)"},
	{name: "boolean", a: "CREATE TABLE t (a boolean)", b: "CREATE TABLE t (a tinyint(1))"},
	{name: "charset echo", a: "CREATE TABLE t (a varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)", b: "CREATE TABLE t (a varchar(10)) DEFAULT CHARSET=utf8mb4"},
	{name: "numeric default quoting", a: "CREATE TABLE t (a int DEFAULT 5, b decimal(5,2) DEFAULT 1.5)", b: "CREATE TABLE t (a int DEFAULT '5', b decimal(5,2) DEFAULT '1.5')"},
	{name: "default null", a: "CREATE TABLE t (a int DEFAULT NULL, b varchar(3) NULL)", b: "CREATE TABLE t (a int, b varchar(3))"},
	{name: "key name", a: "CREATE TABLE t (a int, KEY (a))", b: "CREATE TABLE t (a int, KEY a (a))"},
	{name: "key order", a: "CREATE TABLE t (a int, b int, KEY ka (a), KEY kb (b))", b: "CREATE TABLE t (a int, b int, KEY kb (b), KEY ka (a))"},
	{name: "inline primary key", a: "CREATE TABLE t (id int PRIMARY KEY)", b: "CREATE TABLE t (id int NOT NULL, PRIMARY KEY (id))"},
	{name: "engine and row format", a: "CREATE TABLE t (a int) ENGINE=MyISAM ROW_FORMAT=COMPACT", b: "CREATE TABLE t (a int)"},
	{name: "key type", a: "CREATE TABLE t (a int, KEY k (a) USING HASH)", b: "CREATE TABLE t (a int, KEY k (a))"},
	{name: "fulltext kept as nothing", a: "CREATE TABLE t (a text, FULLTEXT KEY ft (a))", b: "CREATE TABLE t (a text)"},
	{name: "linear hash", a: "CREATE TABLE t (a int) PARTITION BY LINEAR HASH (a) PARTITIONS 4", b: "CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 4"},
	{name: "auto increment counter", a: "CREATE TABLE t (id int AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT=100", b: "CREATE TABLE t (id int AUTO_INCREMENT PRIMARY KEY)"},
	{name: "check spelling", a: "CREATE TABLE t (a int, CONSTRAINT c CHECK (a>0))", b: "CREATE TABLE t (a int, CONSTRAINT c CHECK (`a` > 0))"},
	{name: "view qualifier", a: "CREATE TABLE t (a int); CREATE VIEW v AS SELECT a FROM t", b: "CREATE TABLE t (a int); CREATE VIEW v AS SELECT t.a FROM omni_diff_b.t"},
	{name: "sequence defaults", a: "CREATE SEQUENCE s", b: "CREATE SEQUENCE s START 1 MINVALUE 1 INCREMENT 1 CACHE 1000 NOCYCLE"},
	{name: "foreign key backing index", a: "CREATE TABLE p (id int PRIMARY KEY); CREATE TABLE c (id int PRIMARY KEY, pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))", b: "CREATE TABLE p (id int PRIMARY KEY); CREATE TABLE c (id int PRIMARY KEY, pid int, KEY fk (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))"},

	// Two schemas.
	{name: "decimal default scale", a: "CREATE TABLE t (a decimal(5,2) DEFAULT 1.5)", b: "CREATE TABLE t (a decimal(5,2) DEFAULT 1.50)"},
	{name: "column type", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a bigint)"},
	{name: "nullability", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int NOT NULL)"},
	{name: "column default", a: "CREATE TABLE t (a int DEFAULT 1)", b: "CREATE TABLE t (a int DEFAULT 2)"},
	{name: "column comment", a: "CREATE TABLE t (a int COMMENT 'x')", b: "CREATE TABLE t (a int COMMENT 'it''s')"},
	{name: "column collation", a: "CREATE TABLE t (a varchar(5))", b: "CREATE TABLE t (a varchar(5) COLLATE utf8mb4_general_ci)"},
	{name: "on update", a: "CREATE TABLE t (a timestamp DEFAULT CURRENT_TIMESTAMP)", b: "CREATE TABLE t (a timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)"},
	{name: "generated expression", a: "CREATE TABLE t (a int, b int AS (a + 1))", b: "CREATE TABLE t (a int, b int AS (a + 2))"},
	{name: "added column", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int, b int)"},
	{name: "key columns", a: "CREATE TABLE t (a int, b int, KEY k (a))", b: "CREATE TABLE t (a int, b int, KEY k (a, b))"},
	{name: "key uniqueness", a: "CREATE TABLE t (a int, KEY k (a))", b: "CREATE TABLE t (a int, UNIQUE KEY k (a))"},
	{name: "key visibility", a: "CREATE TABLE t (a int, KEY k (a))", b: "CREATE TABLE t (a int, KEY k (a) INVISIBLE)", syncDrops: "index visibility"},
	{name: "key comment", a: "CREATE TABLE t (a int, KEY k (a))", b: "CREATE TABLE t (a int, KEY k (a) COMMENT 'c')"},
	{name: "key prefix", a: "CREATE TABLE t (a varchar(20), KEY k (a))", b: "CREATE TABLE t (a varchar(20), KEY k (a(5)))"},
	{name: "expression key", a: "CREATE TABLE t (a int, KEY k ((a + 1)))", b: "CREATE TABLE t (a int, KEY k ((a + 2)))"},
	{name: "clustered primary key", a: "CREATE TABLE t (id int, PRIMARY KEY (id) CLUSTERED)", b: "CREATE TABLE t (id int, PRIMARY KEY (id) NONCLUSTERED)"},
	{name: "auto random bits", a: "CREATE TABLE t (id bigint AUTO_RANDOM(5) PRIMARY KEY)", b: "CREATE TABLE t (id bigint AUTO_RANDOM(6) PRIMARY KEY)"},
	{name: "shard row id bits", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int) SHARD_ROW_ID_BITS=4"},
	{name: "auto id cache", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int) AUTO_ID_CACHE=100", syncDrops: "AUTO_ID_CACHE"},
	{name: "ttl", a: "CREATE TABLE t (a datetime) TTL = a + INTERVAL 1 DAY", b: "CREATE TABLE t (a datetime) TTL = a + INTERVAL 2 DAY", syncDrops: "TTL"},
	{name: "table comment", a: "CREATE TABLE t (a int) COMMENT 'x'", b: "CREATE TABLE t (a int) COMMENT 'y'"},
	{name: "table collation", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int) COLLATE utf8mb4_general_ci"},
	{name: "partition count", a: "CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 4", b: "CREATE TABLE t (a int) PARTITION BY HASH (a) PARTITIONS 8"},
	{name: "range bound", a: "CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))", b: "CREATE TABLE t (a int) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (20))"},
	{name: "foreign key action", a: "CREATE TABLE p (id int PRIMARY KEY); CREATE TABLE c (pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))", b: "CREATE TABLE p (id int PRIMARY KEY); CREATE TABLE c (pid int, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id) ON DELETE CASCADE)"},
	{name: "check expression", a: "CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 0))", b: "CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 1))"},
	{name: "check enforcement", a: "CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 0))", b: "CREATE TABLE t (a int, CONSTRAINT c CHECK (a > 0) NOT ENFORCED)", syncDrops: "NOT ENFORCED"},
	{name: "view body", a: "CREATE TABLE t (a int, b int); CREATE VIEW v AS SELECT a FROM t", b: "CREATE TABLE t (a int, b int); CREATE VIEW v AS SELECT b AS a FROM t"},
	{name: "view check option", a: "CREATE TABLE t (a int); CREATE VIEW v AS SELECT a FROM t WITH CASCADED CHECK OPTION", b: "CREATE TABLE t (a int); CREATE VIEW v AS SELECT a FROM t WITH LOCAL CHECK OPTION", syncDrops: "view CHECK OPTION"},
	{name: "added view", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int); CREATE VIEW v AS SELECT a FROM t"},
	{name: "sequence increment", a: "CREATE SEQUENCE s", b: "CREATE SEQUENCE s INCREMENT 2"},
	{name: "sequence cache", a: "CREATE SEQUENCE s", b: "CREATE SEQUENCE s NOCACHE"},
	{name: "sequence comment", a: "CREATE SEQUENCE s", b: "CREATE SEQUENCE s COMMENT 'c'"},
	{name: "added sequence", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int); CREATE SEQUENCE s"},
	{name: "added table", a: "CREATE TABLE t (a int)", b: "CREATE TABLE t (a int); CREATE TABLE u (a int)"},
}

const (
	diffOracleDBA = "omni_diff_a"
	diffOracleDBB = "omni_diff_b"
	// diffOracleDB is the one database both sides load into in the catalog.
	diffOracleDB = "omni_diff"
)

// TestOracle_DiffMatchesServer checks diffOraclePairs against TiDB.
func TestOracle_DiffMatchesServer(t *testing.T) {
	ctr, done := startContainer(t)
	defer done()
	exec := func(t *testing.T, stmt string) {
		t.Helper()
		if _, err := ctr.db.ExecContext(ctr.ctx, stmt); err != nil {
			t.Fatalf("TiDB: %q: %v", stmt, err)
		}
	}
	exec(t, "SET GLOBAL tidb_enable_check_constraint = ON")
	defer exec(t, "SET GLOBAL tidb_enable_check_constraint = OFF")
	defer exec(t, "USE test")
	for _, p := range diffOraclePairs {
		t.Run(p.name, func(t *testing.T) {
			stored := map[string]string{}
			for _, side := range [][2]string{{diffOracleDBA, p.a}, {diffOracleDBB, p.b}} {
				db, script := side[0], side[1]
				exec(t, "DROP DATABASE IF EXISTS "+db)
				exec(t, "CREATE DATABASE "+db)
				exec(t, "USE "+db)
				for _, stmt := range splitStatements(script) {
					if stmt = strings.TrimSpace(stmt); stmt != "" {
						exec(t, stmt)
					}
				}
				stored[db] = storedSchema(t, ctr, db)
			}
			same := stored[diffOracleDBA] == stored[diffOracleDBB]

			a, b := execCatalog(t, p.a), execCatalog(t, p.b)
			checkDiffVerdict(t, "Exec", a, b, same, stored)

			snapA := loadSnapshot(t, syncTiDBSnapshotWithSequences(t, ctr, diffOracleDBA))
			snapB := loadSnapshot(t, syncTiDBSnapshotWithSequences(t, ctr, diffOracleDBB))
			if p.syncDrops != "" {
				if d := Diff(snapA, snapB); !d.IsEmpty() {
					t.Errorf("snapshot: the sync was taken to drop %s, but the snapshots differ: %+v", p.syncDrops, d)
				}
			} else {
				checkDiffVerdict(t, "snapshot", snapA, snapB, same, stored)
			}
			// A schema and the snapshot of it are one schema.
			for _, side := range []struct {
				name     string
				exec     *Catalog
				snapshot *Catalog
			}{{"a", a, snapA}, {"b", b, snapB}} {
				if p.syncDrops != "" {
					continue
				}
				if d := Diff(side.exec, side.snapshot); !d.IsEmpty() {
					t.Errorf("%s: schema text and its snapshot differ: %s", side.name, describeDiff(d))
				}
				if d := Diff(side.snapshot, side.exec); !d.IsEmpty() {
					t.Errorf("%s: snapshot and its schema text differ: %s", side.name, describeDiff(d))
				}
			}
		})
	}
}

func checkDiffVerdict(t *testing.T, path string, a, b *Catalog, same bool, stored map[string]string) {
	t.Helper()
	for _, d := range []*SchemaDiff{Diff(a, b), Diff(b, a)} {
		if d.IsEmpty() != same {
			t.Errorf("%s: Diff empty = %v, TiDB stores the same = %v: %s\nA:\n%s\nB:\n%s",
				path, d.IsEmpty(), same, describeDiff(d), stored[diffOracleDBA], stored[diffOracleDBB])
		}
	}
}

// execCatalog builds a catalog from a script, in diffOracleDB, with the
// session the oracle runs TiDB in.
func execCatalog(t *testing.T, script string) *Catalog {
	t.Helper()
	c := New()
	c.SetCheckConstraints(true)
	results, err := c.Exec("CREATE DATABASE "+diffOracleDB+"; USE "+diffOracleDB+";\n"+
		strings.ReplaceAll(script, diffOracleDBB+".", diffOracleDB+"."), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("catalog: line %d: %v", r.Line, r.Error)
		}
	}
	return c
}

func loadSnapshot(t *testing.T, meta *metadata.DatabaseSchemaMetadata) *Catalog {
	t.Helper()
	return loadSnapshotIn(t, meta, diffOracleDB)
}

// syncTiDBSnapshotWithSequences is syncTiDBSnapshot with the database's
// sequences, which Bytebase's sync also reads.
func syncTiDBSnapshotWithSequences(t *testing.T, ctr *mysqlContainer, db string) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	meta := syncTiDBSnapshot(t, ctr, db)
	meta.Schemas[0].Sequences = syncSequences(t, ctr, db)
	return meta
}

func describeDiff(d *SchemaDiff) string {
	var parts []string
	for _, e := range d.Tables {
		s := e.Action.String() + " table " + e.Name
		for _, c := range e.Columns {
			s += "; " + c.Action.String() + " column " + c.Name
		}
		for _, i := range e.Indexes {
			s += "; " + i.Action.String() + " key " + i.Name
		}
		for _, f := range e.ForeignKeys {
			s += "; " + f.Action.String() + " foreign key " + f.Name
		}
		for _, c := range e.Checks {
			s += "; " + c.Action.String() + " check " + c.Name
		}
		if e.PartitionChanged {
			s += "; partitioning"
		}
		if e.OptionsChanged {
			s += "; options"
		}
		parts = append(parts, s)
	}
	for _, e := range d.Views {
		parts = append(parts, e.Action.String()+" view "+e.Name)
	}
	for _, e := range d.Sequences {
		parts = append(parts, e.Action.String()+" sequence "+e.Name)
	}
	return strings.Join(parts, "\n")
}

// storedSchema is what TiDB stores for a database: SHOW CREATE of every
// table, view and sequence by name, and a view's CHECK_OPTION, which SHOW
// CREATE VIEW leaves out. The database's own name, the AUTO_INCREMENT and
// AUTO_RANDOM_BASE counters, and the order of keys are left out.
func storedSchema(t *testing.T, ctr *mysqlContainer, db string) string {
	t.Helper()
	rows, err := ctr.db.QueryContext(ctr.ctx, "SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME", db)
	if err != nil {
		t.Fatal(err)
	}
	type object struct{ name, kind string }
	var objects []object
	for rows.Next() {
		var o object
		if err := rows.Scan(&o.name, &o.kind); err != nil {
			t.Fatal(err)
		}
		objects = append(objects, o)
	}
	rows.Close()
	var parts []string
	for _, o := range objects {
		var text string
		switch o.kind {
		case "VIEW":
			var checkOption string
			if err := ctr.db.QueryRowContext(ctr.ctx, "SELECT CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", db, o.name).Scan(&checkOption); err != nil {
				t.Fatal(err)
			}
			text = showCreateOnTiDB(t, ctr, "view:"+o.name) + " CHECK_OPTION=" + checkOption
		case "SEQUENCE":
			text = showCreateSequenceOnTiDB(t, ctr, o.name)
		default:
			text = sortKeyLines(storedCounters.ReplaceAllString(showCreateOnTiDB(t, ctr, "table:"+o.name), ""))
		}
		parts = append(parts, strings.ReplaceAll(text, quoteIdent(db)+".", ""))
	}
	slices.Sort(parts)
	return strings.Join(parts, "\n")
}

var storedCounters = regexp.MustCompile(` AUTO_INCREMENT=\d+| /\*T!\[auto_rand_base\] AUTO_RANDOM_BASE=\d+ \*/`)
