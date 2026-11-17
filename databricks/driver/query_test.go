package driver

import (
	"math"
	"testing"

	dbsql "github.com/databricks/databricks-sdk-go/service/sql"

	"github.com/bytebase/omni/driver"
)

func TestToValue(t *testing.T) {
	for _, tc := range []struct {
		name     string
		typeName string
		value    string
		want     driver.Value
		wantErr  bool
	}{
		{name: "bool true", typeName: "BOOLEAN", value: "true", want: driver.Value{Kind: driver.Bool, Bool: true}},
		{name: "bool null", typeName: "BOOLEAN", value: "", want: driver.Value{Kind: driver.Null}},
		{name: "binary", typeName: "BINARY", value: "abc", want: driver.Value{Kind: driver.Bytes, Bytes: []byte("abc")}},
		{name: "short", typeName: "SHORT", value: "12", want: driver.Value{Kind: driver.Int32, Int64: 12}},
		{name: "int", typeName: "INT", value: "-3", want: driver.Value{Kind: driver.Int32, Int64: -3}},
		{name: "long", typeName: "LONG", value: "9223372036854775807", want: driver.Value{Kind: driver.Int64, Int64: math.MaxInt64}},
		{name: "float", typeName: "FLOAT", value: "1.5", want: driver.Value{Kind: driver.Float32, Float64: float64(float32(1.5))}},
		{name: "double", typeName: "DOUBLE", value: "1.25", want: driver.Value{Kind: driver.Float64, Float64: 1.25}},
		{name: "an empty string is the empty string", typeName: "STRING", value: "", want: driver.Value{Kind: driver.String, String: ""}},
		{name: "an empty non-string is null", typeName: "DATE", value: "", want: driver.Value{Kind: driver.Null}},
		{name: "a decimal stays text", typeName: "DECIMAL", value: "1.20", want: driver.Value{Kind: driver.String, String: "1.20"}},
		{name: "a date stays text", typeName: "DATE", value: "2026-11-17", want: driver.Value{Kind: driver.String, String: "2026-11-17"}},
		{name: "a timestamp stays text", typeName: "TIMESTAMP", value: "2026-11-17 10:00:00", want: driver.Value{Kind: driver.String, String: "2026-11-17 10:00:00"}},
		{name: "an array stays text", typeName: "ARRAY", value: "[1,2]", want: driver.Value{Kind: driver.String, String: "[1,2]"}},
		{name: "an int that is not a number", typeName: "INT", value: "x", wantErr: true},
		{name: "a bool that is not a bool", typeName: "BOOLEAN", value: "maybe", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toValue(tc.typeName, tc.value)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("toValue(%q, %q) accepted the value", tc.typeName, tc.value)
				}
				return
			}
			if err != nil {
				t.Fatalf("toValue: %v", err)
			}
			if got.Kind != tc.want.Kind {
				t.Fatalf("toValue(%q, %q).Kind = %s, want %s", tc.typeName, tc.value, got.Kind, tc.want.Kind)
			}
			switch got.Kind {
			case driver.Bool:
				if got.Bool != tc.want.Bool {
					t.Errorf("Bool = %v, want %v", got.Bool, tc.want.Bool)
				}
			case driver.Bytes:
				if string(got.Bytes) != string(tc.want.Bytes) {
					t.Errorf("Bytes = %q, want %q", got.Bytes, tc.want.Bytes)
				}
			case driver.Int32, driver.Int64:
				if got.Int64 != tc.want.Int64 {
					t.Errorf("Int64 = %d, want %d", got.Int64, tc.want.Int64)
				}
			case driver.Float32, driver.Float64:
				if got.Float64 != tc.want.Float64 {
					t.Errorf("Float64 = %v, want %v", got.Float64, tc.want.Float64)
				}
			case driver.String:
				if got.String != tc.want.String {
					t.Errorf("String = %q, want %q", got.String, tc.want.String)
				}
			}
		})
	}
}

func queryResponse(columns []dbsql.ColumnInfo, rows [][]string, truncated bool) *dbsql.StatementResponse {
	return &dbsql.StatementResponse{
		Manifest: &dbsql.ResultManifest{
			Schema:    &dbsql.ResultSchema{Columns: columns},
			Truncated: truncated,
		},
		Result: &dbsql.ResultData{DataArray: rows},
	}
}

func TestResultFromResponse(t *testing.T) {
	t.Run("a query's columns and rows", func(t *testing.T) {
		result, err := resultFromResponse(queryResponse(
			[]dbsql.ColumnInfo{
				{Name: "id", TypeName: dbsql.ColumnInfoTypeNameInt},
				{Name: "name", TypeName: dbsql.ColumnInfoTypeNameString},
			},
			[][]string{{"1", "alpha"}, {"", ""}},
			false,
		), "SELECT id, name FROM t", true, 0)
		if err != nil {
			t.Fatalf("resultFromResponse: %v", err)
		}
		if len(result.Columns) != 2 || result.Columns[0].Name != "id" || result.Columns[0].TypeName != "INT" {
			t.Errorf("Columns = %+v", result.Columns)
		}
		if len(result.Rows) != 2 {
			t.Fatalf("Rows = %d, want 2", len(result.Rows))
		}
		if result.Rows[0][0].Kind != driver.Int32 || result.Rows[0][0].Int64 != 1 {
			t.Errorf("row 0 cell 0 = %+v, want int32 1", result.Rows[0][0])
		}
		if result.Rows[0][1].Kind != driver.String || result.Rows[0][1].String != "alpha" {
			t.Errorf("row 0 cell 1 = %+v, want alpha", result.Rows[0][1])
		}
		if result.Rows[1][0].Kind != driver.Null {
			t.Errorf("row 1 cell 0 = %+v, want NULL", result.Rows[1][0])
		}
		if result.Rows[1][1].Kind != driver.String || result.Rows[1][1].String != "" {
			t.Errorf("row 1 cell 1 = %+v, want the empty string", result.Rows[1][1])
		}
	})

	t.Run("a query that returned no rows still states its columns", func(t *testing.T) {
		result, err := resultFromResponse(queryResponse(
			[]dbsql.ColumnInfo{{Name: "id", TypeName: dbsql.ColumnInfoTypeNameLong}},
			nil,
			false,
		), "SELECT id FROM t", true, 0)
		if err != nil {
			t.Fatalf("resultFromResponse: %v", err)
		}
		if len(result.Columns) != 1 || len(result.Rows) != 0 {
			t.Errorf("result = %+v, want one column and no rows", result)
		}
	})

	t.Run("a query without columns is an error", func(t *testing.T) {
		if _, err := resultFromResponse(&dbsql.StatementResponse{}, "SELECT 1", true, 0); err == nil {
			t.Fatal("resultFromResponse accepted a query with no columns")
		}
	})

	t.Run("a command states no columns, rows or count", func(t *testing.T) {
		result, err := resultFromResponse(&dbsql.StatementResponse{}, "CREATE TABLE t (a int)", false, 0)
		if err != nil {
			t.Fatalf("resultFromResponse: %v", err)
		}
		if len(result.Columns) != 0 || len(result.Rows) != 0 || result.AffectedRows != nil {
			t.Errorf("result = %+v, want no rows and no count", result)
		}
	})

	t.Run("a row that does not match the columns is an error", func(t *testing.T) {
		_, err := resultFromResponse(queryResponse(
			[]dbsql.ColumnInfo{{Name: "id", TypeName: dbsql.ColumnInfoTypeNameInt}},
			[][]string{{"1", "extra"}},
			false,
		), "SELECT id FROM t", true, 0)
		if err == nil {
			t.Fatal("resultFromResponse accepted a row with a different width")
		}
	})

	t.Run("MaxResultBytes stops after the row that crosses it", func(t *testing.T) {
		result, err := resultFromResponse(queryResponse(
			[]dbsql.ColumnInfo{{Name: "id", TypeName: dbsql.ColumnInfoTypeNameInt}},
			[][]string{{"1"}, {"2"}},
			false,
		), "SELECT id FROM t", true, 1)
		if err != nil {
			t.Fatalf("resultFromResponse: %v", err)
		}
		if !result.Truncated {
			t.Error("Truncated = false, want true")
		}
		if len(result.Rows) != 1 {
			t.Errorf("Rows = %d, want the first row only", len(result.Rows))
		}
	})

	t.Run("the server's truncation is kept", func(t *testing.T) {
		result, err := resultFromResponse(queryResponse(
			[]dbsql.ColumnInfo{{Name: "id", TypeName: dbsql.ColumnInfoTypeNameInt}},
			[][]string{{"1"}},
			true,
		), "SELECT id FROM t", true, 0)
		if err != nil {
			t.Fatalf("resultFromResponse: %v", err)
		}
		if !result.Truncated {
			t.Error("Truncated = false, want the server's truncation")
		}
	})
}
