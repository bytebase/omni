package driver

import (
	"testing"

	"github.com/bytebase/omni/driver"
)

func TestAffectedRows(t *testing.T) {
	tests := []struct {
		name   string
		result *driver.Result
		want   int64
	}{
		{name: "nil", result: nil},
		{
			name: "an insert",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int64, Int64: 2}}},
			},
			want: 2,
		},
		{
			name: "a count beside another column",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_inserted_rows", TypeName: "LONG"}, {Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int64, Int64: 2}, {Kind: driver.Int64, Int64: 3}}},
			},
			want: 3,
		},
		{
			name: "more than one row",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int64, Int64: 2}}, {{Kind: driver.Int64, Int64: 3}}},
			},
			want: 5,
		},
		{
			name: "an int32 count",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "INT"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int32, Int64: 4}}},
			},
			want: 4,
		},
		{
			name: "a null count",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Null}}},
			},
		},
		{
			name: "no count column",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "id", TypeName: "INT"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int32, Int64: 1}}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := affectedRows(test.result); got != test.want {
				t.Errorf("affectedRows = %d, want %d", got, test.want)
			}
		})
	}
}
