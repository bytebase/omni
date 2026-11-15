package catalog

import (
	"testing"
)

func TestDiffReportsEntries(t *testing.T) {
	from := execCatalog(t, `
CREATE TABLE p (id int PRIMARY KEY);
CREATE TABLE t (id int PRIMARY KEY, a int, b int, pid int, KEY ka (a), KEY kb (b),
  CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id), CONSTRAINT ck CHECK (a > 0));
CREATE TABLE gone (a int);
CREATE VIEW v AS SELECT a FROM t;
CREATE SEQUENCE s;`)
	to := execCatalog(t, `
CREATE TABLE p (id int PRIMARY KEY);
CREATE TABLE t (id int PRIMARY KEY, a bigint, c int, pid int, KEY ka (a, c),
  CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id) ON DELETE CASCADE, CONSTRAINT ck2 CHECK (c > 0)) COMMENT 'x';
CREATE TABLE added (a int);
CREATE VIEW v AS SELECT c AS a FROM t;
CREATE SEQUENCE s INCREMENT 2;`)
	got := describeDiff(Diff(from, to))
	want := `ADD table added
DROP table gone
MODIFY table t; MODIFY column a; DROP column b; ADD column c; MODIFY key ka; DROP key kb; MODIFY foreign key fk; DROP check ck; ADD check ck2; options
MODIFY view v
MODIFY sequence s`
	if got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
	if d := Diff(to, to); !d.IsEmpty() {
		t.Errorf("a catalog differs from itself: %s", describeDiff(d))
	}
}

// A table's default charset is the default of columns added later; the
// columns it has keep theirs.
func TestDiffTableCharsetLeavesColumnsAsTheyAre(t *testing.T) {
	from := execCatalog(t, "CREATE TABLE t (a varchar(5)) CHARSET latin1")
	to := execCatalog(t, "CREATE TABLE t (a varchar(5) CHARSET latin1) CHARSET utf8mb4")
	if got, want := describeDiff(Diff(from, to)), "MODIFY table t; options"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	to = execCatalog(t, "CREATE TABLE t (a varchar(5)) CHARSET utf8mb4")
	if got, want := describeDiff(Diff(from, to)), "MODIFY table t; MODIFY column a; options"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDiffLeavesOutTemporaryTables(t *testing.T) {
	from := execCatalog(t, "CREATE TABLE t (a int)")
	to := execCatalog(t, "CREATE TABLE t (a int); CREATE TEMPORARY TABLE tmp (a int)")
	if d := Diff(from, to); !d.IsEmpty() {
		t.Errorf("temporary table diffed: %s", describeDiff(d))
	}
}
