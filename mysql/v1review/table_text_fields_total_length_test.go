package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestTableTextFieldsTotalLength(t *testing.T) {
	columns := func(types ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, typ := range types {
			out = append(out, &metadata.ColumnMetadata{Type: typ})
		}
		return out
	}
	final := &v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{Name: "t", Columns: columns("TEXT", "varchar(20)", "int(11)", "decimal(10,2)")},
			{Name: "parted", Columns: columns("varchar(300)"), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}}},
		},
	}}}}
	type got struct{ Name, Value, Text string }
	for _, tt := range []struct {
		name    string
		sql     string
		maximum int
		want    []got
	}{
		{name: "create judged by the final table", sql: "CREATE TABLE t (a INT)", maximum: 65565, want: []got{{"t", "65566", "CREATE TABLE t (a INT)"}}},
		{name: "alter judged by the final table", sql: "ALTER TABLE t ADD COLUMN x INT", maximum: 100, want: []got{{"t", "65566", "ALTER TABLE t ADD COLUMN x INT"}}},
		{name: "within the limit", sql: "ALTER TABLE t ADD COLUMN x INT", maximum: 65566},
		{name: "table name compared case-insensitively", sql: "ALTER TABLE T ADD COLUMN x INT", maximum: 0, want: []got{{"T", "65566", "ALTER TABLE T ADD COLUMN x INT"}}},
		{name: "partition finds its table", sql: "ALTER TABLE p0 ADD COLUMN x INT", maximum: 0, want: []got{{"p0", "300", "ALTER TABLE p0 ADD COLUMN x INT"}}},
		{name: "table missing from final", sql: "CREATE TABLE gone (a TEXT)", maximum: 0},
		{name: "database not part of the name", sql: "ALTER TABLE db.t ADD COLUMN x INT", maximum: 0, want: []got{{"t", "65566", "ALTER TABLE db.t ADD COLUMN x INT"}}},
		{name: "drop not judged", sql: "DROP TABLE t", maximum: 0},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE t ADD COLUMN x INT; END", maximum: 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatal(err)
			}
			var findings []got
			for _, f := range TableTextFieldsTotalLength(list.Items, final, tt.maximum) {
				if f.Code != v1review.TableTextTooLong {
					t.Fatalf("code %v", f.Code)
				}
				findings = append(findings, got{f.Name, f.Value, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, findings); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTextLength(t *testing.T) {
	for typ, want := range map[string]int64{
		"CHAR": 1, "binary": 1, "tinytext": 255, "TINYBLOB": 255, "text": 65_535, "blob": 65_535,
		"mediumtext": 16_777_215, "mediumblob": 16_777_215, "longtext": 4_294_967_295, "longblob": 4_294_967_295,
		"int": 0, "char(123)": 123, "VARCHAR(123)": 123, "int(11)": 11, "decimal(10,2)": 0,
		"varchar(99999999999999999999)": 0, "enum('a')": 0, "": 0,
	} {
		if got := textLength(typ); got != want {
			t.Errorf("textLength(%q) = %d, want %d", typ, got, want)
		}
	}
}

func TestTableTextFieldsTotalLengthNilFinalPanicsOnlyWhenJudging(t *testing.T) {
	list, err := parser.Parse("DROP TABLE t")
	if err != nil {
		t.Fatal(err)
	}
	if got := TableTextFieldsTotalLength(list.Items, nil, 0); got != nil {
		t.Errorf("got %v", got)
	}
	list, err = parser.Parse("CREATE TABLE t (a INT)")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recover() == nil {
			t.Error("want a panic")
		}
	}()
	TableTextFieldsTotalLength(list.Items, nil, 0)
}
