package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestFigureColName_XmlExpr(t *testing.T) {
	tests := []struct {
		name string
		op   nodes.XmlExprOp
		want string
	}{
		{"xmlconcat", nodes.IS_XMLCONCAT, "xmlconcat"},
		{"xmlelement", nodes.IS_XMLELEMENT, "xmlelement"},
		{"xmlforest", nodes.IS_XMLFOREST, "xmlforest"},
		{"xmlparse", nodes.IS_XMLPARSE, "xmlparse"},
		{"xmlpi", nodes.IS_XMLPI, "xmlpi"},
		{"xmlroot", nodes.IS_XMLROOT, "xmlroot"},
		{"xmlserialize", nodes.IS_XMLSERIALIZE, "xmlserialize"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &nodes.XmlExpr{Op: tt.op}
			got := figureColName(n)
			if got != tt.want {
				t.Errorf("figureColName(XmlExpr{Op: %v}) = %q, want %q", tt.op, got, tt.want)
			}
		})
	}
}

func TestFigureColName_XmlSerialize(t *testing.T) {
	n := &nodes.XmlSerialize{}
	got := figureColName(n)
	if got != "xmlserialize" {
		t.Errorf("figureColName(XmlSerialize) = %q, want %q", got, "xmlserialize")
	}
}

func TestFigureColName_GroupingFunc(t *testing.T) {
	n := &nodes.GroupingFunc{}
	got := figureColName(n)
	if got != "grouping" {
		t.Errorf("figureColName(GroupingFunc) = %q, want %q", got, "grouping")
	}
}

func TestFigureColName_AIndirection(t *testing.T) {
	// A_Indirection with a field name as last element.
	n := &nodes.A_Indirection{
		Arg: &nodes.ColumnRef{},
		Indirection: &nodes.List{
			Items: []nodes.Node{
				&nodes.String{Str: "field_name"},
			},
		},
	}
	got := figureColName(n)
	if got != "field_name" {
		t.Errorf("figureColName(A_Indirection) = %q, want %q", got, "field_name")
	}

	// A_Indirection with non-string last element should return "?column?".
	n2 := &nodes.A_Indirection{
		Arg:         &nodes.ColumnRef{},
		Indirection: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
	}
	got2 := figureColName(n2)
	if got2 != "?column?" {
		t.Errorf("figureColName(A_Indirection with star) = %q, want %q", got2, "?column?")
	}
}
