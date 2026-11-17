package driver

import (
	"strings"
	"testing"
)

func TestGetQualifiedTblName(t *testing.T) {
	for _, tc := range []struct {
		name    string
		catalog string
		schema  string
		table   string
		want    string
		wantErr bool
	}{
		{name: "all three parts", catalog: "cat", schema: "sch", table: "tbl", want: "`cat`.`sch`.`tbl`"},
		{name: "no catalog", schema: "sch", table: "tbl", want: "`sch`.`tbl`"},
		{name: "a backtick in a part is doubled", catalog: "cat", schema: "sch`ema", table: "t`bl", want: "`cat`.`sch``ema`.`t``bl`"},
		{name: "no table", catalog: "cat", schema: "sch", wantErr: true},
		{name: "no schema", catalog: "cat", table: "tbl", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getQualifiedTblName(tc.catalog, tc.schema, tc.table)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("getQualifiedTblName accepted %q/%q", tc.schema, tc.table)
				}
				return
			}
			if err != nil {
				t.Fatalf("getQualifiedTblName: %v", err)
			}
			if got != tc.want {
				t.Errorf("getQualifiedTblName = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestQualifyViewDDL(t *testing.T) {
	for _, tc := range []struct {
		name string
		ddl  string
		want string
	}{
		{
			name: "a schema-qualified view gains its catalog",
			ddl:  "CREATE VIEW omni_record.recent (\n  id INT)\nAS SELECT 1",
			want: "CREATE VIEW `workspace`.`omni_record`.`recent` (\n  id INT)\nAS SELECT 1",
		},
		{
			name: "an already qualified view is replaced whole",
			ddl:  "CREATE VIEW `workspace`.`omni_record`.`recent` AS SELECT 1",
			want: "CREATE VIEW `workspace`.`omni_record`.`recent` AS SELECT 1",
		},
		{
			name: "a quoted identifier may hold the separator spellings",
			ddl:  "CREATE VIEW sch.`sales AS archived (old)` (\n  id INT)\nAS SELECT 1",
			want: "CREATE VIEW `workspace`.`omni_record`.`recent` (\n  id INT)\nAS SELECT 1",
		},
		{
			name: "another statement is left alone",
			ddl:  "CREATE MATERIALIZED VIEW `workspace`.`omni_record`.`totals` AS SELECT 1",
			want: "CREATE MATERIALIZED VIEW `workspace`.`omni_record`.`totals` AS SELECT 1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := qualifyViewDDL(tc.ddl, "`workspace`.`omni_record`.`recent`")
			if got != tc.want {
				t.Errorf("qualifyViewDDL = %q, want %q", got, tc.want)
			}
			if !strings.HasPrefix(got, "CREATE ") {
				t.Errorf("qualifyViewDDL lost the statement: %q", got)
			}
		})
	}
}
