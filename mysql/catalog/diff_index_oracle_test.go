package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the index differ (correctness-protocol.md gates 1 & 2), against the LIVE
// MySQL engines (5.7 :13307, 8.0 :13306). For every index variant the node covers, two
// properties are proven mechanically:
//
//  1. Idempotence (the spine): a schema in its real stored form (the engine's SHOW CREATE
//     readback, reloaded into a catalog) self-diffs empty, and the user form vs the engine's
//     stored form diffs empty — e.g. an unnamed KEY (a) (which MySQL auto-names `a`) collapses
//     onto the readback's named form. A non-empty diff is a normalization bug.
//  2. The FK-implicit cross-reference: a foreign key with no explicit backing index produces NO
//     index diff (the auto-created KEY is excluded on both sides), proven via a multi-table
//     round-trip (single-table readback cannot resolve the FK's referenced table).
//
// The harness reuses assertDiffEmptyAgainstReadback / connectOracle / serverCharsetFor / both /
// only from the diff + normalize oracle tests; it skips cleanly when the engines are unreachable.

// indexDiffProbes enumerates the single-table index FORMS that must round-trip empty: a user-form
// DDL whose engine-stored form differs (auto-name, prefix echo, version-flagged DESC/KBS, ...) but
// must canonicalize equal. Cross-table FK forms are proven separately (they need both tables).
func indexDiffProbes() []diffProbe {
	return []diffProbe{
		// PRIMARY KEY — single and composite; columns are forced NOT NULL by the PK.
		{"pk-simple", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)", "t", both()},
		{"pk-composite", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))", "t", both()},
		// UNIQUE — MySQL's unique constraint is a unique index; single + composite.
		{"unique-single", "CREATE TABLE t (id INT PRIMARY KEY, a INT, UNIQUE KEY ua (a))", "t", both()},
		{"unique-composite", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, UNIQUE KEY uab (a,b))", "t", both()},
		// Secondary (non-unique) — single, composite, multiple.
		{"secondary", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ka (a), KEY kab (a,b))", "t", both()},
		// Unnamed indexes → MySQL auto-name (a, a_2): the user form vs stored named form must
		// collapse. The headline auto-name canonicalization case.
		{"unnamed-autoname", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY (a), KEY (a,b))", "t", both()},
		// USING BTREE echo (InnoDB coerces HASH→BTREE; BTREE is echoed verbatim).
		{"using-btree", "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ua (a) USING BTREE) ENGINE=InnoDB", "t", both()},
		// Prefix length col(N).
		{"prefix", "CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(100), c CHAR(40), KEY sp (s(10)), KEY cp (c(5)))", "t", both()},
		// Index COMMENT (escaped).
		{"comment", "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) COMMENT 'idx ''c'' here')", "t", both()},
		// Index KEY_BLOCK_SIZE — echoed on 5.7, absorbed on 8.0; must round-trip on BOTH.
		{"key-block-size", "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a) KEY_BLOCK_SIZE=4) ENGINE=InnoDB", "t", both()},
		// Descending — stored on 8.0, silently ascending on 5.7; round-trips on both.
		{"descending", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC))", "t", both()},
		// Composite mixed ASC/DESC + prefix together.
		{"composite-mixed", "CREATE TABLE t (id INT PRIMARY KEY, a INT, s VARCHAR(50), KEY mix (a DESC, s(8)))", "t", both()},
		// Invisible — 8.0 only (/*!80000 INVISIBLE */).
		{"invisible", "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY va (a) /*!80000 INVISIBLE */)", "t", only(MySQL80)},
		// Functional / expression index — 8.0 only.
		{"functional", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY fa ((a + 1)), KEY fab ((a + b)))", "t", only(MySQL80)},
		// FULLTEXT.
		{"fulltext", "CREATE TABLE t (id INT PRIMARY KEY, body TEXT, ftcol VARCHAR(200), FULLTEXT KEY ftb (body), FULLTEXT KEY ftc (ftcol))", "t", both()},
		// SPATIAL — needs a NOT NULL geometry column. 8.0 supports SPATIAL on InnoDB; 5.7 needs
		// MyISAM. Prove on 8.0 (InnoDB default) to keep the form simple.
		{"spatial", "CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL, SPATIAL KEY sg (g))", "t", only(MySQL80)},
		// Mixed bag: PK + unique + secondary + invisible on one table (8.0).
		{"mixed-all-80", "CREATE TABLE t (id INT NOT NULL, a INT, b INT, c VARCHAR(40), PRIMARY KEY (id), UNIQUE KEY ua (a), KEY kb (b DESC), KEY kc (c(5)) /*!80000 INVISIBLE */)", "t", only(MySQL80)},
		// Mixed bag for 5.7 (no DESC/invisible/functional).
		{"mixed-all-57", "CREATE TABLE t (id INT NOT NULL, a INT, b INT, c VARCHAR(40), PRIMARY KEY (id), UNIQUE KEY ua (a), KEY kb (b), KEY kc (c(5)))", "t", both()},
	}
}

// TestOracle_IndexDiffIdempotence proves gates 1 & 2 for every single-table index form: user DDL
// vs its engine readback diffs EMPTY, and the stored form self-diffs empty, on every supported
// version.
func TestOracle_IndexDiffIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, probe := range indexDiffProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := "idx_" + strings.ReplaceAll(probe.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, probe.create, probe.table)
			})
		}
	}
}

// fkIndexCase is a parent+child schema with a foreign key, used to prove the FK-implicit
// cross-reference round-trips empty.
type fkIndexCase struct {
	id       string
	childDDL string // table `c`, references table `p`
	versions []Version
}

// fkIndexCases enumerate the FK + index interplay forms. Each builds parent `p` + child `c`,
// reads BOTH back, reloads as one schema, and the diff must be empty — the FK backing index
// (auto-created by MySQL, synthesized by the loader) must not phantom-diff.
func fkIndexCases() []fkIndexCase {
	return []fkIndexCase{
		// Named FK, no explicit backing index → MySQL auto-creates KEY `fk`.
		{"fk-named-implicit", "CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Unnamed FK → constraint `c_ibfk_1`, backing index named after the column `pid`.
		{"fk-unnamed-implicit", "CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Explicit covering index present (distinct name) → MySQL reuses it, no separate backing
		// index; the explicit index is user-managed and must round-trip.
		{"fk-explicit-covering", "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Composite FK with its own auto-created composite backing index.
		{"fk-composite", "CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))", both()},
	}
}

// TestOracle_IndexDiffFKImplicit proves the FK-implicit cross-reference: a schema with a foreign
// key and its (auto-created) backing index self-diffs empty when loaded from the real engine
// readbacks of BOTH tables. This is the multi-table analog of the idempotence probe — the
// single-table harness cannot resolve the FK's referenced table.
func TestOracle_IndexDiffFKImplicit(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)
		ctx := context.Background()

		// parent `p` has both a single-col PK (for fk-named/unnamed/explicit) and a composite PK
		// (for fk-composite); pick the parent matching the child.
		parentFor := func(id string) string {
			if id == "fk-composite" {
				return "CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))"
			}
			return "CREATE TABLE p (id INT NOT NULL PRIMARY KEY, name VARCHAR(50))"
		}

		for _, fc := range fkIndexCases() {
			if !containsVersion(fc.versions, version) {
				continue
			}
			t.Run(o.name+"/"+fc.id, func(t *testing.T) {
				dbName := "idxfk_" + strings.ReplaceAll(fc.id, "-", "_")
				stmts := []string{
					"DROP DATABASE IF EXISTS " + dbName,
					fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, sc),
					"USE " + dbName,
					parentFor(fc.id),
					fc.childDDL,
				}
				for _, s := range stmts {
					if _, err := o.db.ExecContext(ctx, s); err != nil {
						t.Skipf("[%s] setup failed (may be expected): %q: %v", o.name, s, err)
					}
				}
				// Read both tables back and reload as a single schema.
				var b strings.Builder
				fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, sc, dbName)
				for _, tbl := range []string{"p", "c"} {
					var nm, ddl string
					row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, tbl))
					if err := row.Scan(&nm, &ddl); err != nil {
						t.Fatalf("[%s] SHOW CREATE %s: %v", o.name, tbl, err)
					}
					b.WriteString(ddl + ";\n")
				}
				cat, err := LoadSQL(b.String())
				if err != nil {
					t.Fatalf("[%s] reload of FK readback failed: %v\n%s", o.name, err, b.String())
				}
				if d := DiffWithNormalizer(cat, cat, n); !d.IsEmpty() {
					t.Errorf("[%s] FK schema self-diff not empty (phantom on backing index): %s", o.name, describeDiff(d))
				}
			})
		}
	}
}
