package driver

import (
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
