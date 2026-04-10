package ast

import "testing"

func TestTypeKind_String(t *testing.T) {
	cases := []struct {
		kind TypeKind
		want string
	}{
		{TypeUnknown, "Unknown"},
		{TypeInt, "Int"},
		{TypeNumber, "Number"},
		{TypeFloat, "Float"},
		{TypeBoolean, "Boolean"},
		{TypeDate, "Date"},
		{TypeDateTime, "DateTime"},
		{TypeTime, "Time"},
		{TypeTimestamp, "Timestamp"},
		{TypeTimestampLTZ, "TimestampLTZ"},
		{TypeTimestampNTZ, "TimestampNTZ"},
		{TypeTimestampTZ, "TimestampTZ"},
		{TypeChar, "Char"},
		{TypeVarchar, "Varchar"},
		{TypeBinary, "Binary"},
		{TypeVarbinary, "Varbinary"},
		{TypeVariant, "Variant"},
		{TypeObject, "Object"},
		{TypeArray, "Array"},
		{TypeGeography, "Geography"},
		{TypeGeometry, "Geometry"},
		{TypeVector, "Vector"},
		{TypeKind(999), "Unknown"},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			if got := c.kind.String(); got != c.want {
				t.Errorf("TypeKind(%d).String() = %q, want %q", c.kind, got, c.want)
			}
		})
	}
}

func TestTypeName_Tag(t *testing.T) {
	var n TypeName
	if (&n).Tag() != T_TypeName {
		t.Errorf("Tag() = %v, want T_TypeName", (&n).Tag())
	}
}

func TestTypeName_WalkerVisitsElementType(t *testing.T) {
	// Verify the walker descends into ElementType.
	inner := &TypeName{Kind: TypeVarchar, Name: "VARCHAR"}
	outer := &TypeName{Kind: TypeArray, Name: "ARRAY", ElementType: inner, VectorDim: -1}

	var visited []NodeTag
	Inspect(outer, func(n Node) bool {
		if n != nil {
			visited = append(visited, n.Tag())
		}
		return true
	})

	// Should visit outer (TypeName) then inner (TypeName).
	if len(visited) != 2 {
		t.Fatalf("visited %d nodes, want 2: %+v", len(visited), visited)
	}
	if visited[0] != T_TypeName || visited[1] != T_TypeName {
		t.Errorf("visited = %+v, want [T_TypeName, T_TypeName]", visited)
	}
}
