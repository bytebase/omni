package catalog

import "testing"

// TestCommentOnTypeName holds COMMENT ON TYPE and DOMAIN, whose object the
// parser gives as a TypeName, to what PostgreSQL 17.12 answers.
func TestCommentOnTypeName(t *testing.T) {
	setup := "CREATE TYPE mood AS ENUM ('x'); CREATE DOMAIN d AS int; "
	for _, tc := range []struct {
		sql  string
		code string // "" for success
	}{
		{"COMMENT ON TYPE mood IS 'a'", ""},
		{"COMMENT ON TYPE int IS 'b'", ""},
		{"COMMENT ON TYPE mood[] IS 'c'", ""},
		{"COMMENT ON TYPE mood(1) IS 'd'", CodeSyntaxError},
		{"COMMENT ON DOMAIN mood IS 'e'", CodeWrongObjectType},
		{"COMMENT ON DOMAIN d IS 'f'", ""},
		{"COMMENT ON TYPE d IS 'g'", ""},
		{"COMMENT ON TYPE nope IS 'h'", CodeUndefinedObject},
		{"COMMENT ON TYPE IS IS 'i'", CodeUndefinedObject},
		{"COMMENT ON DOMAIN public.d IS NULL", ""},
		{"CREATE TYPE \"Mood\" AS ENUM ('x'); COMMENT ON TYPE \"Mood\" IS 'j'", ""},
		{"COMMENT ON TYPE \"Mood\" IS 'k'", CodeUndefinedObject},
		{"CREATE TYPE public.int4 AS ENUM ('x'); COMMENT ON TYPE pg_catalog.int4 IS 'l'", ""},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			results, err := New().Exec(setup+tc.sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			var got error
			for _, r := range results {
				if r.Error != nil {
					got = r.Error
					break
				}
			}
			if tc.code == "" {
				if got != nil {
					t.Fatalf("got %v, want success", got)
				}
				return
			}
			assertCode(t, got, tc.code)
		})
	}
}

// TestCommentOnTypeNameTarget: a type comment lands on the type named as
// written: a quoted name keeps its case, and pg_catalog's type is not the
// search path's one that shadows it.
func TestCommentOnTypeNameTarget(t *testing.T) {
	c := New()
	if _, err := c.Exec(`CREATE TYPE "Mood" AS ENUM ('x'); CREATE TYPE mood AS ENUM ('y'); CREATE TYPE public.int4 AS ENUM ('z');
COMMENT ON TYPE "Mood" IS 'upper'; COMMENT ON TYPE mood IS 'lower'; COMMENT ON TYPE pg_catalog.int4 IS 'builtin'`, nil); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ schema, name, want string }{
		{"public", "Mood", "upper"}, {"public", "mood", "lower"}, {"pg_catalog", "int4", "builtin"}, {"public", "int4", ""},
	} {
		s := c.schemaByName[tc.schema]
		bt := c.typeByName[typeKey{ns: s.OID, name: tc.name}]
		if bt == nil {
			t.Fatalf("no type %s.%s", tc.schema, tc.name)
		}
		if got := c.comments[commentKey{ObjType: 't', ObjOID: bt.OID}]; got != tc.want {
			t.Errorf("comment on %s.%s = %q, want %q", tc.schema, tc.name, got, tc.want)
		}
	}
}
