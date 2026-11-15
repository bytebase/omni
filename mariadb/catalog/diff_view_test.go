package catalog

import (
	"strings"
	"testing"
)

// Hermetic unit tests for the view differ. They assert the diff actions and what
// the key folds: own-database qualifiers, ALGORITHM, SQL SECURITY, CHECK OPTION,
// the column list, and a definer that is not compared.

// loadViewCatalog loads SDL describing one or more views (and their tables) at a given version
// and returns the catalog. The schema is wrapped in a database so view identity is (db, name).
func loadViewCatalog(t *testing.T, version Version, sdl string) *Catalog {
	t.Helper()
	wrapped := "CREATE DATABASE vt DEFAULT CHARSET=utf8mb4;\nUSE vt;\n" + sdl
	cat, err := LoadSDLWithVersion(wrapped, version)
	if err != nil {
		t.Fatalf("LoadSDLWithVersion failed for %q: %v", sdl, err)
	}
	return cat
}

func getView(t *testing.T, c *Catalog, name string) *View {
	t.Helper()
	db := c.GetDatabase("vt")
	if db == nil {
		t.Fatalf("database vt not found")
	}
	v := db.Views[toLower(name)]
	if v == nil {
		t.Fatalf("view %q not found", name)
	}
	return v
}

// findViewDiff returns the diff entry for a named view, or nil.
func findViewDiff(d *SchemaDiff, name string) *ViewDiffEntry {
	for i := range d.Views {
		if strings.EqualFold(d.Views[i].Name, name) {
			return &d.Views[i]
		}
	}
	return nil
}

func TestDiffViews_SelfDiffEmpty(t *testing.T) {
	for _, version := range []Version{MariaDB118} {
		sdl := `CREATE TABLE t (a INT, b VARCHAR(20), c INT);
CREATE VIEW v1 AS SELECT a, b FROM t;
CREATE VIEW v2 (p, q) AS SELECT a, b FROM t;
CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v3 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;`
		cat := loadViewCatalog(t, version, sdl)
		n := NormalizerFor(version)
		if d := DiffWithNormalizer(cat, cat, n); !d.IsEmpty() {
			t.Errorf("[v%d] self-diff not empty: views=%d", version, len(d.Views))
		}
	}
}

func TestDiffViews_Add(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT, b INT);`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	e := findViewDiff(d, "v1")
	if e == nil || e.Action != DiffAdd {
		t.Fatalf("want v1 ADD, got %+v", d.Views)
	}
	if e.To == nil || e.From != nil {
		t.Errorf("ADD entry should carry To and no From")
	}
}

func TestDiffViews_Drop(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	e := findViewDiff(d, "v1")
	if e == nil || e.Action != DiffDrop {
		t.Fatalf("want v1 DROP, got %+v", d.Views)
	}
	if e.From == nil || e.To != nil {
		t.Errorf("DROP entry should carry From and no To")
	}
}

func TestDiffViews_ModifyBody(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT b FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	e := findViewDiff(d, "v1")
	if e == nil || e.Action != DiffModify {
		t.Fatalf("want v1 MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_ModifyAlgorithm(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE ALGORITHM=UNDEFINED VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE ALGORITHM=MERGE VIEW v1 AS SELECT a FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("algorithm change should be a MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_ModifySQLSecurity(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE SQL SECURITY DEFINER VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE SQL SECURITY INVOKER VIEW v1 AS SELECT a FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("sql security change should be a MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_ModifyCheckOption(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("check-option change should be a MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_AddRemoveCheckOption(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("adding a check option should be a MODIFY, got %+v", d.Views)
	}
}

// DEFINER differences must not produce a diff.
func TestDiffViews_DefinerIgnored(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE DEFINER=`+"`root`@`%`"+` VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT);
CREATE DEFINER=`+"`admin`@`localhost`"+` VIEW v1 AS SELECT a FROM t;`)
	// Sanity: the definers really differ on the loaded views.
	if getView(t, from, "v1").Definer == getView(t, to, "v1").Definer {
		t.Fatalf("test setup: definers should differ")
	}
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	if !d.IsEmpty() {
		t.Errorf("DEFINER-only change must NOT diff (ignore-in-diff), got %d view changes", len(d.Views))
	}
}

// An explicit column-list change is a modify.
func TestDiffViews_ExplicitColumns80(t *testing.T) {
	from := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 (p, q) AS SELECT a, b FROM t;`)
	to := loadViewCatalog(t, MariaDB118, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 (x, y) AS SELECT a, b FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MariaDB118))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("explicit column-list change should be a MODIFY on 8.0, got %+v", d.Views)
	}
}

func TestCanonicalView_UserVsEngineBodyConverge(t *testing.T) {
	// MariaDB stores a view's references into its own database without the
	// database, so a body that writes them and one that does not are one view.
	qualified := loadCat(t, "CREATE DATABASE vt; USE vt; CREATE TABLE t (a INT, b INT); CREATE VIEW v1 AS SELECT vt.t.a, vt.t.b FROM vt.t;")
	plain := loadCat(t, "CREATE DATABASE vt; USE vt; CREATE TABLE t (a INT, b INT); CREATE VIEW v1 AS SELECT t.a, t.b FROM t;")
	if d := Diff(qualified, plain); !d.IsEmpty() {
		t.Errorf("own-database qualifiers made a diff:\n%s", describeDiff(d))
	}
}
