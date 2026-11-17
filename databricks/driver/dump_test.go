package driver

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
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

func TestParseTableProperties(t *testing.T) {
	got := parseTableProperties("Owner=omni, Created-by=delta, delta.minReaderVersion=1, path=s3://bucket")
	want := map[string]string{
		"Owner":                  "omni",
		"Created-by":             "delta",
		"delta.minReaderVersion": "1",
		"path":                   "s3",
	}
	// The pattern is Bytebase's: a value ends at the first character it does
	// not cover. "s3://bucket" therefore reads as "s3", which is what a
	// recording must confirm.
	for key, value := range want {
		if got[key] != value {
			t.Errorf("properties[%q] = %q, want %q", key, got[key], value)
		}
	}
}

func TestFormatTblProperties(t *testing.T) {
	got := formatTblProperties(&tableAdditionalInfo{tblProps: map[string]string{
		"z": "last",
		"a": "first",
		"m": "middle",
	}})
	want := "a = first,\nm = middle,\nz = last\n"
	if got != want {
		t.Errorf("formatTblProperties = %q, want %q", got, want)
	}
	if got := formatTblProperties(nil); got != "" {
		t.Errorf("formatTblProperties(nil) = %q, want empty", got)
	}
}

func TestGenMaterializedViewDDL(t *testing.T) {
	got := genMaterializedViewDDL("`cat`.`sch`.`mv`", &metadatapb.MaterializedViewMetadata{
		Comment: "a materialized view",
	}, "a = first,\nz = last\n")
	if !strings.HasPrefix(got, "CREATE MATERIALIZED VIEW `cat`.`sch`.`mv`\n") {
		t.Errorf("DDL = %q, want CREATE MATERIALIZED VIEW", got)
	}
	if !strings.Contains(got, "COMMENT a materialized view\n") {
		t.Errorf("DDL = %q, want Bytebase's unquoted comment", got)
	}
	if !strings.Contains(got, "TBLPROPERTIES a = first,\nz = last\n\n") {
		t.Errorf("DDL = %q, want the properties in key order", got)
	}
}
