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

func TestOrderViewObjects(t *testing.T) {
	objects := []dumpObject{
		{catalog: "cat", schema: "a", name: "v", kind: "VIEW"},
		{catalog: "cat", schema: "z", name: "t", kind: "VIEW"},
		{catalog: "cat", schema: "b", name: "mid", kind: "VIEW"},
		{catalog: "cat", schema: "c", name: "base", kind: "VIEW"},
	}
	got := orderViewObjects(objects, map[string][]string{
		"cat.a.v":   {"cat.z.t", "cat.other.table"},
		"cat.z.t":   {"cat.b.mid"},
		"cat.b.mid": {"cat.c.base"},
	})
	want := []string{"cat.c.base", "cat.b.mid", "cat.z.t", "cat.a.v"}
	for i, object := range got {
		if object.key() != want[i] {
			t.Errorf("object %d = %q, want %q", i, object.key(), want[i])
		}
	}

	cycle := []dumpObject{
		{catalog: "cat", schema: "a", name: "one", kind: "VIEW"},
		{catalog: "cat", schema: "a", name: "two", kind: "VIEW"},
	}
	got = orderViewObjects(cycle, map[string][]string{
		"cat.a.one": {"cat.a.two"},
		"cat.a.two": {"cat.a.one"},
	})
	if len(got) != 2 || got[0].name != "one" || got[1].name != "two" {
		t.Errorf("cycle order = %v, want both views in name order", got)
	}
}
