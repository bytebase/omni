package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestCollectionTypes pins the tree of OceanBase's collection types and holds a
// MAP to the nesting limit. Whether each form parses is held to OceanBase by the
// data_type rows of testdata/compat.tsv.
func TestCollectionTypes(t *testing.T) {
	column := func(sql string) *nodes.DataType {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0].(*nodes.CreateTableStmt).Columns[0].TypeName
	}
	// ARRAY(ARRAY(INT)[]): the suffix applies to the nested ARRAY.
	if dt := column("CREATE TABLE t (a ARRAY(ARRAY(INT)[]))"); dt.Name != "ARRAY" || dt.ArrayElem.Name != "ARRAY" ||
		dt.ArrayElem.ArrayElem.Name != "ARRAY" || dt.ArrayElem.ArrayElem.ArrayElem.Name != "INT" {
		t.Errorf("ARRAY(ARRAY(INT)[]) = %+v", dt)
	}
	if dt := column("CREATE TABLE t (a MAP(INT, VARCHAR(3)))"); dt.Name != "MAP" || dt.MapKey.Name != "INT" || dt.ArrayElem.Name != "VARCHAR" {
		t.Errorf("MAP(INT, VARCHAR(3)) = %+v", dt)
	}
	deep := "CREATE TABLE t (a " + strings.Repeat("MAP(INT, ", maxNestingDepth+1) + "INT" + strings.Repeat(")", maxNestingDepth+1) + ")"
	if _, err := Parse(deep); !isNestingError(err) {
		t.Errorf("a MAP nested past the limit = %v, want the nesting error", err)
	}
}
