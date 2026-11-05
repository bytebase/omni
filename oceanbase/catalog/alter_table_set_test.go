package catalog

import (
	"context"
	"fmt"
	"testing"
)

// alterTableSetScenarios each change a table's options through OceanBase's ALTER TABLE
// ... SET, or through options written one after another without commas, on
// alterTableSetBase.
var alterTableSetScenarios = []string{
	"ALTER TABLE t SET COMMENT 'x'",
	"ALTER TABLE t SET COMMENT = 'y'",
	"ALTER TABLE t SET PCTFREE = 6 COMMENT 'z'",
	"ALTER TABLE t SET PCTFREE = 7, COMMENT 'w'",
	"ALTER TABLE t SET PCTFREE = 8, SET COMMENT 'v'",
	"ALTER TABLE t PCTFREE = 9 COMMENT 'u'",
	"ALTER TABLE t SET DEFAULT CHARACTER SET = gbk COLLATE gbk_bin",
	"ALTER TABLE t SET ROW_FORMAT = COMPRESSED",
	"ALTER TABLE t SET ENGINE = InnoDB ROW_FORMAT = COMPACT",
	"ALTER TABLE t SET COMMENT 'one', SET COMMENT 'two'",
}

const alterTableSetBase = "CREATE TABLE t (a INT PRIMARY KEY)"

// TestAlterTableSetAppliesOptions pins that an option set through SET, or after
// another without a comma, applies as it does on its own.
func TestAlterTableSetAppliesOptions(t *testing.T) {
	for _, tc := range []struct {
		alter                               string
		comment, charset, collation, rowFmt string
	}{
		{alter: "ALTER TABLE t SET COMMENT 'x'", comment: "x"},
		{alter: "ALTER TABLE t SET PCTFREE = 6 COMMENT 'z'", comment: "z"},
		{alter: "ALTER TABLE t PCTFREE = 9 COMMENT 'u'", comment: "u"},
		{alter: "ALTER TABLE t SET COMMENT 'one', SET COMMENT 'two'", comment: "two"},
		{alter: "ALTER TABLE t SET DEFAULT CHARACTER SET = gbk COLLATE gbk_bin", charset: "gbk", collation: "gbk_bin"},
		{alter: "ALTER TABLE t SET ROW_FORMAT = COMPRESSED COMMENT 'c'", comment: "c", rowFmt: "COMPRESSED"},
	} {
		t.Run(tc.alter, func(t *testing.T) {
			c := New()
			results, err := c.Exec("CREATE DATABASE d; USE d; "+alterTableSetBase+"; "+tc.alter+";", nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("Exec: %v", r.Error)
				}
			}
			tbl := c.GetDatabase("d").GetTable("t")
			if tc.comment != "" && tbl.Comment != tc.comment {
				t.Errorf("comment = %q, want %q", tbl.Comment, tc.comment)
			}
			if tc.charset != "" && tbl.Charset != tc.charset {
				t.Errorf("charset = %q, want %q", tbl.Charset, tc.charset)
			}
			if tc.collation != "" && tbl.Collation != tc.collation {
				t.Errorf("collation = %q, want %q", tbl.Collation, tc.collation)
			}
			if tc.rowFmt != "" && tbl.RowFormat != tc.rowFmt {
				t.Errorf("row format = %q, want %q", tbl.RowFormat, tc.rowFmt)
			}
		})
	}
}

// TestOracle_OceanBaseAlterTableSet applies each of alterTableSetScenarios on
// OceanBase, reads the table back, and requires the catalog's result and the readback
// to canonicalize equal, but for the obEchoedRowFormat gap. A SET of
// DYNAMIC_PARTITION_POLICY, which the catalog does not model, must still execute in the
// catalog and leave a readback it can load.
//
//	go test ./oceanbase/catalog/ -run OceanBase -count=1
func TestOracle_OceanBaseAlterTableSet(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for i, alter := range alterTableSetScenarios {
		t.Run(alter, func(t *testing.T) {
			dbName := fmt.Sprintf("obset_%d", i)
			if _, err := storeAndRead(ctx, o, dbName, alterTableSetBase, "t"); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+dbName) })
			if _, err := o.db.ExecContext(ctx, alter); err != nil {
				t.Fatalf("OceanBase refused %q: %v", alter, err)
			}
			var readback string
			if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE "+dbName+".t").Scan(new(string), &readback); err != nil {
				t.Fatal(err)
			}
			assertDiffEmptyOrRowFormatEcho(t, o, readback, alterTableSetBase+";\n"+alter+";", "t")
		})
	}
	t.Run("dynamic partition policy", func(t *testing.T) {
		const create = "CREATE TABLE t (a DATETIME NOT NULL, PRIMARY KEY (a)) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY', PRECREATE_TIME = '3day') PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN ('2026-01-01'))"
		const alter = "ALTER TABLE t SET DYNAMIC_PARTITION_POLICY (ENABLE = FALSE)"
		if _, err := storeAndRead(ctx, o, "obset_dp", create, "t"); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _, _ = o.db.ExecContext(ctx, "DROP DATABASE IF EXISTS obset_dp") })
		if _, err := o.db.ExecContext(ctx, alter); err != nil {
			t.Fatalf("OceanBase refused %q: %v", alter, err)
		}
		var readback string
		if err := o.db.QueryRowContext(ctx, "SHOW CREATE TABLE obset_dp.t").Scan(new(string), &readback); err != nil {
			t.Fatal(err)
		}
		sc := serverCharsetFor(o.version)
		loadOneTable(t, sc, create+";\n"+alter+";", "t")
		loadOneTable(t, sc, readback, "t")
	})
}
