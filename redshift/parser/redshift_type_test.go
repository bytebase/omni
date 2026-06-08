package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftTypeNamesParse(t *testing.T) {
	tree, err := Parse(`
CREATE TABLE redshift_types (
    c_varchar_max VARCHAR(MAX),
    c_character_varying CHARACTER VARYING(100),
    c_nvarchar NVARCHAR(200),
    c_bpchar BPCHAR(25),
    c_varbyte VARBYTE(1024),
    c_varbinary VARBINARY(512),
    c_binary_varying BINARY VARYING(256),
    c_super SUPER,
    c_geometry GEOMETRY,
    c_geography GEOGRAPHY,
    c_hllsketch HLLSKETCH
);`)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	stmt := firstCreateStmt(t, tree)
	wantTypes := map[string]string{
		"c_varchar_max":       "pg_catalog.varchar",
		"c_character_varying": "pg_catalog.varchar",
		"c_nvarchar":          "pg_catalog.varchar",
		"c_bpchar":            "pg_catalog.bpchar",
		"c_varbyte":           "pg_catalog.varbyte",
		"c_varbinary":         "pg_catalog.varbyte",
		"c_binary_varying":    "pg_catalog.varbyte",
		"c_super":             "pg_catalog.super",
		"c_geometry":          "pg_catalog.geometry",
		"c_geography":         "pg_catalog.geography",
		"c_hllsketch":         "pg_catalog.hllsketch",
	}
	for _, item := range stmt.TableElts.Items {
		col, ok := item.(*nodes.ColumnDef)
		if !ok {
			t.Fatalf("expected ColumnDef, got %T", item)
		}
		want, ok := wantTypes[col.Colname]
		if !ok {
			t.Fatalf("unexpected column %q", col.Colname)
		}
		if got := typeNameString(col.TypeName); got != want {
			t.Fatalf("column %q type = %q, want %q", col.Colname, got, want)
		}
	}
}

func TestRedshiftVarcharMaxTypmod(t *testing.T) {
	tree, err := Parse("CREATE TABLE events (event_data VARCHAR(MAX));")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	stmt := firstCreateStmt(t, tree)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)
	if col.TypeName.Typmods == nil || len(col.TypeName.Typmods.Items) != 1 {
		t.Fatalf("expected one typmod, got %#v", col.TypeName.Typmods)
	}
	value, ok := col.TypeName.Typmods.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected MAX typmod to be String, got %T", col.TypeName.Typmods.Items[0])
	}
	if value.Str != "max" {
		t.Fatalf("expected MAX typmod %q, got %q", "max", value.Str)
	}
}

func TestRedshiftVarcharMaxInFunctionSignatureParse(t *testing.T) {
	sql := "ALTER FUNCTION transform_json(json_data varchar(max), options varchar) OWNER TO data_team;"
	if _, err := Parse(sql); err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
}

func typeNameString(typeName *nodes.TypeName) string {
	if typeName == nil || typeName.Names == nil {
		return ""
	}
	var parts []string
	for _, item := range typeName.Names.Items {
		str, ok := item.(*nodes.String)
		if !ok {
			continue
		}
		parts = append(parts, str.Str)
	}
	return strings.Join(parts, ".")
}
