package driver

import (
	"testing"

	"github.com/bytebase/omni/driver"
)

func TestAffectedRows(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		result    *driver.Result
		want      int64
	}{
		{name: "nil", statement: "INSERT INTO t VALUES (1)", result: nil},
		{
			name:      "an insert",
			statement: "INSERT INTO t VALUES (1)",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int64, Int64: 2}}},
			},
			want: 2,
		},
		{
			name:      "a count beside another column",
			statement: "INSERT INTO t VALUES (1)",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_inserted_rows", TypeName: "LONG"}, {Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int64, Int64: 2}, {Kind: driver.Int64, Int64: 3}}},
			},
			want: 3,
		},
		{
			name:      "more than one row",
			statement: "INSERT INTO t VALUES (1)",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int64, Int64: 2}}, {{Kind: driver.Int64, Int64: 3}}},
			},
			want: 5,
		},
		{
			name:      "an int32 count",
			statement: "UPDATE t SET a = 1",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "INT"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int32, Int64: 4}}},
			},
			want: 4,
		},
		{
			name:      "a null count",
			statement: "DELETE FROM t",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "LONG"}},
				Rows:    [][]driver.Value{{{Kind: driver.Null}}},
			},
		},
		{
			name:      "no count column",
			statement: "CREATE TABLE t (a INT)",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "id", TypeName: "INT"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int32, Int64: 1}}},
			},
		},
		{
			name:      "a read is not a write",
			statement: "SELECT 100 AS num_affected_rows",
			result: &driver.Result{
				Columns: []driver.Column{{Name: "num_affected_rows", TypeName: "INT"}},
				Rows:    [][]driver.Value{{{Kind: driver.Int32, Int64: 100}}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := affectedRows(test.statement, test.result); got != test.want {
				t.Errorf("affectedRows = %d, want %d", got, test.want)
			}
		})
	}
}
