package catalog

import (
	"strings"
	"testing"
)

// Hermetic unit tests for the view differ + generator (no live engine). They assert the
// structural diff actions, the canonical-key folding (body db-qualifier stripping, ALGORITHM /
// SQL SECURITY / CHECK OPTION / explicit-column comparison, DEFINER-ignored), and the generated
// DDL shape + view-on-view ordering. The oracle gates (idempotence + apply-correctness against
// the real 5.7 / 8.0 engines) live in diff_view_oracle_test.go and migration_view_oracle_test.go.

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
	for _, version := range []Version{MySQL57, MySQL80} {
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
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT, b INT);`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	e := findViewDiff(d, "v1")
	if e == nil || e.Action != DiffAdd {
		t.Fatalf("want v1 ADD, got %+v", d.Views)
	}
	if e.To == nil || e.From != nil {
		t.Errorf("ADD entry should carry To and no From")
	}
}

func TestDiffViews_Drop(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	e := findViewDiff(d, "v1")
	if e == nil || e.Action != DiffDrop {
		t.Fatalf("want v1 DROP, got %+v", d.Views)
	}
	if e.From == nil || e.To != nil {
		t.Errorf("DROP entry should carry From and no To")
	}
}

func TestDiffViews_ModifyBody(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 AS SELECT b FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	e := findViewDiff(d, "v1")
	if e == nil || e.Action != DiffModify {
		t.Fatalf("want v1 MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_ModifyAlgorithm(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE ALGORITHM=UNDEFINED VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE ALGORITHM=MERGE VIEW v1 AS SELECT a FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("algorithm change should be a MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_ModifySQLSecurity(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE SQL SECURITY DEFINER VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE SQL SECURITY INVOKER VIEW v1 AS SELECT a FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("sql security change should be a MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_ModifyCheckOption(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH LOCAL CHECK OPTION;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("check-option change should be a MODIFY, got %+v", d.Views)
	}
}

func TestDiffViews_AddRemoveCheckOption(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("adding a check option should be a MODIFY, got %+v", d.Views)
	}
}

// DEFINER differences must NOT produce a diff (ignore-in-diff).
func TestDiffViews_DefinerIgnored(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE DEFINER=`+"`root`@`%`"+` VIEW v1 AS SELECT a FROM t;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT);
CREATE DEFINER=`+"`admin`@`localhost`"+` VIEW v1 AS SELECT a FROM t;`)
	// Sanity: the definers really differ on the loaded views.
	if getView(t, from, "v1").Definer == getView(t, to, "v1").Definer {
		t.Fatalf("test setup: definers should differ")
	}
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	if !d.IsEmpty() {
		t.Errorf("DEFINER-only change must NOT diff (ignore-in-diff), got %d view changes", len(d.Views))
	}
}

// On 8.0, an explicit column-list change is a modify; the list is a stored attribute.
func TestDiffViews_ExplicitColumns80(t *testing.T) {
	from := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 (p, q) AS SELECT a, b FROM t;`)
	to := loadViewCatalog(t, MySQL80, `CREATE TABLE t (a INT, b INT);
CREATE VIEW v1 (x, y) AS SELECT a, b FROM t;`)
	d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	if e := findViewDiff(d, "v1"); e == nil || e.Action != DiffModify {
		t.Fatalf("explicit column-list change should be a MODIFY on 8.0, got %+v", d.Views)
	}
}

// canonicalViewBody strips the view's OWN database qualifier but keeps a cross-database one.
func TestCanonicalViewBody_StripsOwnDatabase(t *testing.T) {
	db := &Database{Name: "vt"}
	v := &View{
		Name:       "v1",
		Database:   db,
		Definition: "select `vt`.`t`.`a` AS `a` from `vt`.`t`",
	}
	got := canonicalViewBody(v)
	want := "select `t`.`a` AS `a` from `t`"
	if got != want {
		t.Errorf("own-db strip:\n  got:  %q\n  want: %q", got, want)
	}

	// A cross-database reference keeps its qualifier.
	v2 := &View{
		Name:       "v2",
		Database:   db,
		Definition: "select `other`.`t`.`a` AS `a` from `other`.`t`",
	}
	if got := canonicalViewBody(v2); got != v2.Definition {
		t.Errorf("cross-db reference must be preserved:\n  got:  %q\n  want: %q", got, v2.Definition)
	}
}

// Regression (review blocker): the strip must be position-aware, never a blind substring replace.
func TestCanonicalViewBody_PositionAware(t *testing.T) {
	db := &Database{Name: "vt"}

	// A relation literally named the same as the database (`vt`.`vt`): only the LEADING database
	// slot is stripped — the table qualifier `vt` must survive.
	v := &View{Name: "w", Database: db, Definition: "select `vt`.`vt`.`a` AS `a` from `vt`.`vt`"}
	got := canonicalViewBody(v)
	want := "select `vt`.`a` AS `a` from `vt`"
	if got != want {
		t.Errorf("db-named relation: own-db slot only must be stripped:\n  got:  %q\n  want: %q", got, want)
	}

	// A string literal that contains the token `vt`. must be left untouched (different literal
	// values must NOT collapse to the same canonical body).
	v1 := &View{Name: "s", Database: db, Definition: "select '`vt`.x' AS `s` from `vt`.`t`"}
	v2 := &View{Name: "s", Database: db, Definition: "select '`vt`.y' AS `s` from `vt`.`t`"}
	b1, b2 := canonicalViewBody(v1), canonicalViewBody(v2)
	if b1 == b2 {
		t.Errorf("string-literal content must not be stripped/collapsed: both = %q", b1)
	}
	if !strings.Contains(b1, "'`vt`.x'") {
		t.Errorf("string literal mangled: %q", b1)
	}
}

// User-form body (db-unqualified) and engine-form body (db-qualified) collapse to the same
// canonical key — the core idempotence property, here exercised without a live engine.
func TestCanonicalView_UserVsEngineBodyConverge(t *testing.T) {
	db := &Database{Name: "vt"}
	user := &View{Name: "v1", Database: db, Algorithm: "UNDEFINED", SqlSecurity: "DEFINER",
		Definition: "select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`"}
	engine := &View{Name: "v1", Database: db, Algorithm: "UNDEFINED", SqlSecurity: "DEFINER",
		Definition: "select `vt`.`t`.`a` AS `a`,`vt`.`t`.`b` AS `b` from `vt`.`t`"}
	n := NormalizerFor(MySQL80)
	if canonicalView(user, n) != canonicalView(engine, n) {
		t.Errorf("user vs engine body did not converge:\n  user:   %q\n  engine: %q",
			canonicalView(user, n), canonicalView(engine, n))
	}
}
