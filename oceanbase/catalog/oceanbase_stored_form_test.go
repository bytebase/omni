package catalog

import (
	"context"
	"strings"
	"testing"
)

// What OceanBase stores where it parts from MySQL, held to the engine. Each test here is
// the evidence for one rule of the normalizer: the engine's answers are on its tape.

// rowFormats is what OceanBase stores for each ROW_FORMAT a table is written with.
// IgnoreRowFormat rests on the first three: DYNAMIC is what a table that declares
// nothing stores, so a readback's DYNAMIC is not a declaration.
var rowFormats = []struct{ written, stored string }{
	{"", "DYNAMIC"},
	{"DEFAULT", "DYNAMIC"},
	{"DYNAMIC", "DYNAMIC"},
	{"COMPACT", "COMPACT"},
	{"REDUNDANT", "REDUNDANT"},
	{"COMPRESSED", "COMPRESSED"},
	{"CONDENSED", "CONDENSED"},
}

// TestOracle_OceanBaseRowFormat holds rowFormats to the engine, and requires each
// written form to canonicalize equal to what the engine stores for it. Before
// IgnoreRowFormat knew DYNAMIC, every table of every schema differed from its own
// readback by one MODIFY.
//
// The second half is the way back: a table that drops its ROW_FORMAT is altered with
// ROW_FORMAT=DEFAULT, and that has to store what a table declaring none stores.
func TestOracle_OceanBaseRowFormat(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, rf := range rowFormats {
		name, create := "unstated", "CREATE TABLE t (id INT PRIMARY KEY)"
		if rf.written != "" {
			name, create = rf.written, create+" ROW_FORMAT="+rf.written
		}
		t.Run(name, gapped(func(t testing.TB) {
			readback, err := storeAndRead(ctx, o, "obrf_"+strings.ToLower(name), create, "t")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(readback, " ROW_FORMAT = "+rf.stored+" ") {
				t.Errorf("recorded ROW_FORMAT = %s, but the engine stored\n%s", rf.stored, readback)
			}
			assertDiffEmptyAgainstStored(t, o, readback, create, "t")
		}))
	}
	n := NormalizerFor(o.version)
	for _, from := range []string{"COMPACT", "COMPRESSED"} {
		t.Run("from "+from+" back to unstated", gapped(func(t testing.TB) {
			assertApplyCorrect(t, o, n, migrationProbe{
				id:      "rowformat-reset-" + strings.ToLower(from),
				table:   "t",
				fromDDL: "CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=" + from,
				toDDL:   "CREATE TABLE t (id INT PRIMARY KEY)",
			})
		}))
	}
}

// TestOracle_OceanBaseEngine holds CanonicalEngine to the engine: OceanBase takes any
// ENGINE a table is written with, a name no engine has included, prints none in SHOW
// CREATE TABLE and reports InnoDB for every table. So the written ENGINE is not schema.
func TestOracle_OceanBaseEngine(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, engine := range []string{"", "InnoDB", "MyISAM", "MEMORY", "nosuch"} {
		name, create := "unstated", "CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))"
		if engine != "" {
			name, create = engine, create+" ENGINE="+engine
		}
		t.Run(name, func(t *testing.T) {
			db := "obeng_" + strings.ToLower(name)
			readback, err := storeAndRead(ctx, o, db, create, "t")
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(strings.ToUpper(readback), "ENGINE") {
				t.Errorf("the readback names an engine:\n%s", readback)
			}
			var reported string
			if err := o.db.QueryRowContext(ctx, "SELECT ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't'", db).Scan(&reported); err != nil {
				t.Fatal(err)
			}
			if reported != "InnoDB" {
				t.Errorf("information_schema reports ENGINE %q, recorded InnoDB", reported)
			}
			assertDiffEmptyAgainstStored(t, o, readback, create, "t")
		})
	}
}

// TestOracle_OceanBaseIndexOrderAndType holds CanonicalIndexColumn and canonicalIndexType
// to the engine. OceanBase takes DESC on an index column and stores the column ascending,
// and SHOW CREATE TABLE prints no USING, whichever type the index was written with. The
// dictionary does keep HASH, where no readback of the table shows it.
func TestOracle_OceanBaseIndexOrderAndType(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	const create = "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY ab (a, b DESC), KEY ha (a) USING HASH, KEY bb (b) USING BTREE, UNIQUE KEY ub (b DESC))"
	readback, err := storeAndRead(ctx, o, "obidx_order_type", create, "t")
	if err != nil {
		t.Fatal(err)
	}
	for _, word := range []string{"DESC", "USING", "HASH", "BTREE"} {
		if strings.Contains(strings.ToUpper(readback), word) {
			t.Errorf("the readback says %s:\n%s", word, readback)
		}
	}
	rows, err := o.db.QueryContext(ctx, "SELECT INDEX_NAME, COLUMN_NAME, COLLATION, INDEX_TYPE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'obidx_order_type' AND TABLE_NAME = 't' ORDER BY BINARY INDEX_NAME, SEQ_IN_INDEX")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var index, column, collation, kind string
		if err := rows.Scan(&index, &column, &collation, &kind); err != nil {
			t.Fatal(err)
		}
		got = append(got, index+" "+column+" "+collation+" "+kind)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	want := []string{"PRIMARY id A BTREE", "ab a A BTREE", "ab b A BTREE", "bb b A BTREE", "ha a A HASH", "ub b A BTREE"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("information_schema.STATISTICS:\n%s\nrecorded:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	assertDiffEmptyAgainstStored(t, o, readback, create, "t")
}

// TestOracle_OceanBaseViewOptions holds canonicalViewAlgorithm and
// canonicalViewSQLSecurity to the engine: OceanBase takes ALGORITHM and SQL SECURITY on
// a view and keeps neither. SHOW CREATE VIEW prints none of them, and
// information_schema.VIEWS reports SECURITY_TYPE NONE for every view. The check option
// is kept, and stays part of a view's key.
func TestOracle_OceanBaseViewOptions(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS obview_options", "CREATE DATABASE obview_options", "USE obview_options",
		"CREATE TABLE t (a INT)",
		"CREATE VIEW v_plain AS SELECT a FROM t",
		"CREATE ALGORITHM=MERGE VIEW v_merge AS SELECT a FROM t",
		"CREATE ALGORITHM=TEMPTABLE VIEW v_temp AS SELECT a FROM t",
		"CREATE SQL SECURITY INVOKER VIEW v_inv AS SELECT a FROM t",
		"CREATE SQL SECURITY DEFINER VIEW v_def AS SELECT a FROM t",
		"CREATE VIEW v_check AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION",
	} {
		if _, err := o.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	for _, tc := range []struct{ view, checkOption string }{
		{"v_plain", "NONE"}, {"v_merge", "NONE"}, {"v_temp", "NONE"}, {"v_inv", "NONE"}, {"v_def", "NONE"}, {"v_check", "CASCADED"},
	} {
		t.Run(tc.view, func(t *testing.T) {
			var name, ddl, charset, collation string
			if err := o.db.QueryRowContext(ctx, "SHOW CREATE VIEW "+tc.view).Scan(&name, &ddl, &charset, &collation); err != nil {
				t.Fatal(err)
			}
			for _, word := range []string{"ALGORITHM", "SQL SECURITY", "DEFINER", "INVOKER"} {
				if strings.Contains(strings.ToUpper(ddl), word) {
					t.Errorf("the readback says %s: %s", word, ddl)
				}
			}
			if kept := strings.Contains(strings.ToUpper(ddl), "WITH CHECK OPTION"); kept != (tc.checkOption != "NONE") {
				t.Errorf("the readback keeps the check option: %v, recorded %s: %s", kept, tc.checkOption, ddl)
			}
			var security, check string
			if err := o.db.QueryRowContext(ctx, "SELECT SECURITY_TYPE, CHECK_OPTION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'obview_options' AND TABLE_NAME = ?", tc.view).Scan(&security, &check); err != nil {
				t.Fatal(err)
			}
			if security != "NONE" || check != tc.checkOption {
				t.Errorf("information_schema.VIEWS says SECURITY_TYPE %s, CHECK_OPTION %s; recorded NONE, %s", security, check, tc.checkOption)
			}
		})
	}
}
