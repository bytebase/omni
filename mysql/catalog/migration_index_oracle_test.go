package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle apply-correctness proof for the index generator (correctness-protocol.md gate 2),
// against the LIVE MySQL engines. For representative (from, to) pairs covering every index
// variant, the generated ALTER TABLE ... ADD/DROP index DDL applied to a real `from` database
// must yield a table whose canonical form equals `to`. Idempotence (gate 1) for index forms is
// covered by extending the migration idempotence corpus is not needed here because the index
// diff probes already prove the no-op plan is empty via TestOracle_IndexDiffIdempotence; this
// file proves the EMITTED DDL is correct.
//
// The harness reuses assertApplyCorrect / catalogFromReadback / loadOneTable / connectOracle /
// both / only from migration_oracle_test.go; it skips cleanly when the engines are unreachable.

// indexMigrationProbes enumerates the index ADD/DROP/MODIFY FORMS the generator covers. Each is a
// single-table (from, to) pair (FK-backing indexes are proven separately — they need a parent and
// are owned by the FK node). An empty fromDDL means a pure CREATE, which exercises the
// "new table → ADD INDEX for every non-PK index" path (formatCreateTable inlines only the PK).
func indexMigrationProbes() []migrationProbe {
	return []migrationProbe{
		// ---- CREATE TABLE carrying indexes (new table; non-PK indexes added via ALTER) ----
		{"create-with-secondary", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kab (a,b))", both()},
		{"create-with-unique", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, UNIQUE KEY ua (a), UNIQUE KEY uab (a,b))", both()},
		{"create-with-prefix", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY sp (s(10)))", both()},
		{"create-with-comment", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'hello')", both()},
		{"create-with-using-btree", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ua (a) USING BTREE) ENGINE=InnoDB", both()},
		{"create-with-kbs", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) KEY_BLOCK_SIZE=4) ENGINE=InnoDB", both()},
		{"create-with-composite-pk", "t", "",
			"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b), KEY kc (c))", both()},
		{"create-with-fulltext", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY ftb (body))", both()},
		// 8.0-only forms.
		{"create-with-descending", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))", only(MySQL80)},
		{"create-with-invisible", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)", only(MySQL80)},
		{"create-with-functional", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY fa ((a + 1)))", only(MySQL80)},
		{"create-with-spatial", "t", "",
			"CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL, SPATIAL KEY sg (g))", only(MySQL80)},

		// ---- ADD INDEX (modified table) ----
		{"add-secondary", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))", both()},
		{"add-unique", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))", both()},
		{"add-composite", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))", both()},
		{"add-prefix", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100))",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY sp (s(10)))", both()},
		{"add-comment", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'c')", both()},
		{"add-fulltext", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, body TEXT)",
			"CREATE TABLE t (id INT PRIMARY KEY, body TEXT, FULLTEXT KEY ftb (body))", both()},
		{"add-primary-key", "t",
			"CREATE TABLE t (id INT NOT NULL, a INT)",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)", both()},
		// 8.0-only.
		{"add-descending", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))", only(MySQL80)},
		{"add-invisible", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)", only(MySQL80)},
		{"add-functional", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY fa ((a + 1)))", only(MySQL80)},

		// ---- DROP INDEX (modified table) ----
		{"drop-secondary", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)", both()},
		{"drop-unique", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)", both()},
		{"drop-primary-key", "t",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT NOT NULL, a INT)", both()},
		{"drop-one-of-many", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kb (b))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kb (b))", both()},

		// ---- MODIFY index (DROP+ADD): same name, changed shape ----
		{"modify-columns", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))", both()},
		{"modify-kind-to-unique", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY k (a))", both()},
		{"modify-prefix-length", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(10)))",
			"CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), KEY k (s(20)))", both()},
		{"modify-comment", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'old')",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) COMMENT 'new')", both()},
		{"modify-pk-columns", "t",
			"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a))",
			"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))", both()},
		// 8.0-only: visibility flip and direction flip are MODIFYs.
		{"modify-visibility", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) /*!80000 INVISIBLE */)", only(MySQL80)},
		{"modify-direction", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a DESC))", only(MySQL80)},

		// ---- index DROP + column DROP interplay (ordering: DROP INDEX before DROP COLUMN) ----
		// MySQL auto-drops a single-column index when its column is dropped; an explicit DROP INDEX
		// after the column is gone fails errno 1091. The index drop must precede the column drop.
		{"drop-column-and-its-single-col-index", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))",
			"CREATE TABLE t (id INT PRIMARY KEY)", both()},
		// Composite index on (a,b); dropping BOTH columns auto-drops the index → the explicit DROP
		// INDEX must still come first.
		{"drop-columns-and-their-composite-index", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY kab (a,b))",
			"CREATE TABLE t (id INT PRIMARY KEY)", both()},
		// Drop one column of a composite index but keep the index name on the remaining column:
		// k (a,b) → k (b) while dropping a. This is a MODIFY (DROP+ADD) racing a column DROP; the
		// index DROP (PhasePre) precedes the column DROP, and the re-ADD (PhaseMain) lands after.
		{"modify-index-shrink-with-column-drop", "t",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY k (a,b))",
			"CREATE TABLE t (id INT PRIMARY KEY, b INT, KEY k (b))", both()},
		// Add a column AND an index on it in the same plan (ADD COLUMN then ADD INDEX).
		{"add-column-and-index-on-it", "t",
			"CREATE TABLE t (id INT PRIMARY KEY)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))", both()},
	}
}

// TestOracle_IndexMigrationApplyCorrectness proves gate 2 for every index probe: the generated
// DDL transforms a real `from` database into a `to`-equal one (compared via canonical readback).
func TestOracle_IndexMigrationApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range indexMigrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// TestOracle_IndexGenerateFKOnlyChangeEmitsNoIndexOps proves the generate-side of the FK-implicit
// contract: when the only difference between `from` and `to` is a foreign key (its backing index
// rides along), the index generator emits NO ADD/DROP index ops — the FK node owns that index.
// This guards against the index node fighting the FK node (duplicate-key-name on ADD, errno 1553
// on DROP). Proven via a multi-table schema loaded from real engine readbacks.
func TestOracle_IndexGenerateFKOnlyChangeEmitsNoIndexOps(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)
		ctx := context.Background()

		// Load a 2-table schema (parent + child) from real readbacks for a given child DDL.
		loadSchema := func(dbName, childDDL string) *Catalog {
			stmts := []string{
				"DROP DATABASE IF EXISTS " + dbName,
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, sc),
				"USE " + dbName,
				"CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))",
				childDDL,
			}
			for _, s := range stmts {
				if _, err := o.db.ExecContext(ctx, s); err != nil {
					t.Skipf("[%s] setup: %q: %v", o.name, s, err)
				}
			}
			var b strings.Builder
			fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, sc, dbName)
			for _, tbl := range []string{"p", "c"} {
				var nm, ddl string
				row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, tbl))
				if err := row.Scan(&nm, &ddl); err != nil {
					t.Fatalf("[%s] show create %s: %v", o.name, tbl, err)
				}
				b.WriteString(ddl + ";\n")
			}
			cat, err := LoadSQL(b.String())
			if err != nil {
				t.Fatalf("[%s] reload: %v\n%s", o.name, err, b.String())
			}
			return cat
		}

		assertNoIndexOps := func(t *testing.T, from, to *Catalog) {
			t.Helper()
			diff := DiffWithNormalizer(from, to, n)
			plan := GenerateMigrationWithNormalizer(from, to, diff, n)
			for _, op := range plan.Ops {
				if op.Type == OpAddIndex || op.Type == OpDropIndex {
					t.Errorf("[%s] unexpected index op on FK-only change: %s", o.name, op.SQL)
				}
			}
		}

		// from: child with FK (backing index auto-created); to: same child WITHOUT the FK (and
		// hence — in the user's target — without the backing index). The index plan must be empty:
		// the FK node drops the FK; the leftover backing index is the FK node's concern.
		t.Run(o.name+"/drop-fk-no-index-op", func(t *testing.T) {
			from := loadSchema("idxgfk1", "CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))")
			// `to`: child with neither FK nor index (a plain table).
			toWrapped := fmt.Sprintf("CREATE DATABASE idxgfk1 DEFAULT CHARSET=%s;\nUSE idxgfk1;\nCREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50));\nCREATE TABLE c (id INT PRIMARY KEY, pid INT);", sc)
			to, err := LoadSQL(toWrapped)
			if err != nil {
				t.Fatalf("to load: %v", err)
			}
			assertNoIndexOps(t, from, to)
		})

		// from: child with neither; to: child with FK (backing index rides with it). Adding the FK
		// is the FK node's job; the index plan must be empty.
		t.Run(o.name+"/add-fk-no-index-op", func(t *testing.T) {
			from := loadSchema("idxgfk2", "CREATE TABLE c (id INT PRIMARY KEY, pid INT)")
			toCat := loadSchema("idxgfk2b", "CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))")
			assertNoIndexOps(t, from, toCat)
		})
	}
}
