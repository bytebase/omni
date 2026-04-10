package ast

import "testing"

// ---------------------------------------------------------------------------
// Ident tests
// ---------------------------------------------------------------------------

func TestIdent_Normalize(t *testing.T) {
	cases := []struct {
		name  string
		ident Ident
		want  string
	}{
		{"unquoted lowercase", Ident{Name: "foo"}, "FOO"},
		{"unquoted uppercase", Ident{Name: "FOO"}, "FOO"},
		{"unquoted mixed", Ident{Name: "FoO"}, "FOO"},
		{"quoted lowercase", Ident{Name: "foo", Quoted: true}, "foo"},
		{"quoted uppercase", Ident{Name: "FOO", Quoted: true}, "FOO"},
		{"quoted mixed", Ident{Name: "FoO", Quoted: true}, "FoO"},
		{"empty unquoted", Ident{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ident.Normalize(); got != c.want {
				t.Errorf("Normalize() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestIdent_String(t *testing.T) {
	cases := []struct {
		name  string
		ident Ident
		want  string
	}{
		{"unquoted", Ident{Name: "foo"}, "foo"},
		{"quoted simple", Ident{Name: "foo", Quoted: true}, `"foo"`},
		{"quoted with space", Ident{Name: "my table", Quoted: true}, `"my table"`},
		{"quoted with inner quote", Ident{Name: `a"b`, Quoted: true}, `"a""b"`},
		{"empty unquoted", Ident{}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ident.String(); got != c.want {
				t.Errorf("String() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestIdent_IsEmpty(t *testing.T) {
	cases := []struct {
		name  string
		ident Ident
		want  bool
	}{
		{"zero value", Ident{}, true},
		{"has name", Ident{Name: "x"}, false},
		{"quoted empty name", Ident{Quoted: true}, false},
		{"has loc only", Ident{Loc: Loc{Start: 0, End: 1}}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ident.IsEmpty(); got != c.want {
				t.Errorf("IsEmpty() = %v, want %v", got, c.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ObjectName tests
// ---------------------------------------------------------------------------

func TestObjectName_Normalize(t *testing.T) {
	cases := []struct {
		name string
		obj  ObjectName
		want string
	}{
		{"1-part unquoted", ObjectName{Name: Ident{Name: "table"}}, "TABLE"},
		{"1-part quoted", ObjectName{Name: Ident{Name: "Table", Quoted: true}}, "Table"},
		{"2-part", ObjectName{
			Schema: Ident{Name: "schema"},
			Name:   Ident{Name: "table"},
		}, "SCHEMA.TABLE"},
		{"3-part", ObjectName{
			Database: Ident{Name: "db"},
			Schema:   Ident{Name: "schema"},
			Name:     Ident{Name: "table"},
		}, "DB.SCHEMA.TABLE"},
		{"3-part mixed quoting", ObjectName{
			Database: Ident{Name: "My DB", Quoted: true},
			Schema:   Ident{Name: "schema"},
			Name:     Ident{Name: "Quoted Table", Quoted: true},
		}, "My DB.SCHEMA.Quoted Table"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.obj.Normalize(); got != c.want {
				t.Errorf("Normalize() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestObjectName_String(t *testing.T) {
	cases := []struct {
		name string
		obj  ObjectName
		want string
	}{
		{"1-part", ObjectName{Name: Ident{Name: "table"}}, "table"},
		{"2-part", ObjectName{
			Schema: Ident{Name: "schema"},
			Name:   Ident{Name: "table"},
		}, "schema.table"},
		{"3-part quoted", ObjectName{
			Database: Ident{Name: "My DB", Quoted: true},
			Schema:   Ident{Name: "schema"},
			Name:     Ident{Name: "Quoted Table", Quoted: true},
		}, `"My DB".schema."Quoted Table"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.obj.String(); got != c.want {
				t.Errorf("String() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestObjectName_Parts(t *testing.T) {
	cases := []struct {
		name    string
		obj     ObjectName
		wantLen int
	}{
		{"1-part", ObjectName{Name: Ident{Name: "t"}}, 1},
		{"2-part", ObjectName{Schema: Ident{Name: "s"}, Name: Ident{Name: "t"}}, 2},
		{"3-part", ObjectName{Database: Ident{Name: "d"}, Schema: Ident{Name: "s"}, Name: Ident{Name: "t"}}, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			parts := c.obj.Parts()
			if len(parts) != c.wantLen {
				t.Errorf("Parts() len = %d, want %d", len(parts), c.wantLen)
			}
		})
	}
}

func TestObjectName_Matches(t *testing.T) {
	mkObj := func(db, schema, name string, dbQ, schemaQ, nameQ bool) ObjectName {
		o := ObjectName{Name: Ident{Name: name, Quoted: nameQ}}
		if schema != "" {
			o.Schema = Ident{Name: schema, Quoted: schemaQ}
		}
		if db != "" {
			o.Database = Ident{Name: db, Quoted: dbQ}
		}
		return o
	}

	cases := []struct {
		name string
		a, b ObjectName
		want bool
	}{
		{"1-part matches same name (case-folded)",
			mkObj("", "", "foo", false, false, false),
			mkObj("db", "schema", "FOO", false, false, false),
			true},
		{"1-part does NOT match different name",
			mkObj("", "", "foo", false, false, false),
			mkObj("", "", "bar", false, false, false),
			false},
		{"2-part matches same schema.name",
			mkObj("", "schema", "table", false, false, false),
			mkObj("db", "SCHEMA", "TABLE", false, false, false),
			true},
		{"2-part does NOT match different schema",
			mkObj("", "schema1", "table", false, false, false),
			mkObj("", "schema2", "table", false, false, false),
			false},
		{"3-part exact match",
			mkObj("db", "schema", "table", false, false, false),
			mkObj("DB", "SCHEMA", "TABLE", false, false, false),
			true},
		{"3-part does NOT match different db",
			mkObj("db1", "schema", "table", false, false, false),
			mkObj("db2", "schema", "table", false, false, false),
			false},
		{"quoted vs unquoted are DIFFERENT",
			mkObj("", "", "foo", false, false, true),
			mkObj("", "", "foo", false, false, false),
			false},
		{"quoted vs unquoted same normalized value",
			mkObj("", "", "FOO", false, false, true),
			mkObj("", "", "foo", false, false, false),
			true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.a.Matches(c.b); got != c.want {
				t.Errorf("Matches() = %v, want %v\n  a=%+v\n  b=%+v", got, c.want, c.a, c.b)
			}
		})
	}
}

func TestObjectName_Tag(t *testing.T) {
	var n ObjectName
	if (&n).Tag() != T_ObjectName {
		t.Errorf("Tag() = %v, want T_ObjectName", (&n).Tag())
	}
}
